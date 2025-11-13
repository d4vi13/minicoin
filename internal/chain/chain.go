// Package chain implemts blockchain logic
package chain

import (
	"bytes"
	"container/list"
	"crypto/sha512"
	"encoding/binary"
	"log"
	"time"
	"unsafe"
)

// ChainError represents error codes for blockchain operations
type ChainError int

// Blockchain error constants
const (
	CLIENT_NOT_FOUND = iota
	CLIENT_OVERDRAW
	BLOCKCHAIN_TAINTED
	SUCCESS
)

// Global blockchain and failure counter
var chain list.List
var failIn int

// ChainNode represents a single block in the blockchain
type ChainNode struct {
	clientId         uint32    // Client identifier
	transactionTime  int64     // Transaction timestamp
	transactionValue int64     // Transaction amount
	nodeHash         [sha512.Size256]byte  // Block hash
}

// Init initializes the blockchain with a genesis block
func Init(nodeAdditionsUntilFailure int) {
	chain.Init()
	landmarkNode := new(ChainNode)
	chain.PushBack(landmarkNode)
	failIn = nodeAdditionsUntilFailure
}

// Hash calculates the block hash using previous block's hash
func (node *ChainNode) Hash(prevHash []byte) error {
	copy(node.nodeHash[:], prevHash)

	// Serialize node data for hashing
	byteBuffer := make([]byte, unsafe.Sizeof(*node))
	_, err := binary.Encode(byteBuffer, binary.LittleEndian, node)
	if err != nil {
		return err
	}

	// Generate SHA512/256 hash
	newHash := sha512.Sum512_256(byteBuffer)
	copy(node.nodeHash[:], newHash[:])

	return nil
}

// addChainNode creates and adds a new block to the chain
func addChainNode(clientId uint32, transactionValue int64) {
	newNode := new(ChainNode)

	// Countdown to simulated failure
	if failIn >= 0 {
		failIn -= 1
	}

	newNode.clientId = uint32(clientId)
	newNode.transactionTime = time.Now().UnixNano()
	newNode.transactionValue = transactionValue

	// Calculate hash using previous block's hash
	prevHash := chain.Back().Value.(*ChainNode).nodeHash
	err := newNode.Hash(prevHash[:])
	if err != nil {
		log.Println("deu ruim: ", err)
	}

	// Log block details
	log.Printf("Adding node: ")
	log.Printf("\t Client id: [%v]", (*newNode).clientId)
	time := time.Unix(0, (*newNode).transactionTime).Format(time.UnixDate)
	log.Printf("\t Transaction time: [%v]", time)
	log.Printf("\t Node hash: [%x]", (*newNode).nodeHash)

	// Simulate hash corruption for testing
	if failIn == 0 {
		newNode.nodeHash[0] += 1
		log.Println("Altering hash...")
	}

	chain.PushBack(newNode)
}

// isNodeHashIntegral verifies block integrity by recalculating hash
func isNodeHashIntegral(elem *list.Element) bool {
	prev := elem.Prev()
	if prev == nil {
		return true
	}

	nodePrev := prev.Value.(*ChainNode)
	nodeElem := elem.Value.(*ChainNode)

	// Recalculate hash for verification
	nodeCopy := *nodeElem
	err := nodeCopy.Hash(nodePrev.nodeHash[:])
	if err != nil {
		log.Println("Falha na verificacao do hash: ", err)
	}

	return bytes.Equal(nodeCopy.nodeHash[:], nodeElem.nodeHash[:])
}

// AddTransaction adds a new transaction to the blockchain
func AddTransaction(clientId uint32, transactionValue int64) ChainError {
	// Check for overdraft on withdrawals
	if transactionValue < 0 {
		balance, err := GetClientBalance(clientId)
		if err != SUCCESS {
			return err
		}
		if (balance + transactionValue) < 0 {
			return CLIENT_OVERDRAW
		}
	}

	addChainNode(clientId, transactionValue)

	return SUCCESS
}

// IsChainTainted verifies integrity of entire blockchain
func IsChainTainted() bool {
	cont := 1

	// Check each block's hash integrity
	for elem := chain.Front(); elem != nil; elem = elem.Next() {
		if isNodeHashIntegral(elem) == false {
			log.Printf("\t Node [%v] hash tainted", cont)
			return true
		}
		cont += 1
	}

	return false
}

// GetClientBalance calculates total balance for a client
func GetClientBalance(clientId uint32) (int64, ChainError) {
	balance := int64(0)
	clientFound := false

	// Sum all transactions for the client
	for elem := chain.Front(); elem != nil; elem = elem.Next() {
		node := elem.Value.(*ChainNode)
		if node.clientId == uint32(clientId) {
			balance += node.transactionValue
			clientFound = true
		}
	}

	if clientFound == false {
		return balance, CLIENT_NOT_FOUND
	}

	return balance, SUCCESS
}
