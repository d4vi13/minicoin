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

type ChainError int

const (
	CLIENT_NOT_FOUND = iota
	CLIENT_OVERDRAW
	BLOCKCHAIN_TAINTED
	SUCCESS
)

var chain list.List
var failIn int

type ChainNode struct {
	clientId         uint32
	transactionTime  int64
	transactionValue int64
	nodeHash         [sha512.Size256]byte
}

func Init(nodeAdditionsUntilFailure int) {
	chain.Init()
	landmarkNode := new(ChainNode)
	chain.PushBack(landmarkNode)
	failIn = nodeAdditionsUntilFailure
}

func (node *ChainNode) Hash(prevHash []byte) error {
	copy(node.nodeHash[:], prevHash)

	byteBuffer := make([]byte, unsafe.Sizeof(*node))
	_, err := binary.Encode(byteBuffer, binary.LittleEndian, node)
	if err != nil {
		return err
	}

	newHash := sha512.Sum512_256(byteBuffer)
	copy(node.nodeHash[:], newHash[:])

	return nil
}

func addChainNode(clientId uint32, transactionValue int64) {
	newNode := new(ChainNode)

	if failIn >= 0 {
		failIn -= 1
	}

	newNode.clientId = uint32(clientId)
	newNode.transactionTime = time.Now().UnixNano()
	newNode.transactionValue = transactionValue

	prevHash := chain.Back().Value.(*ChainNode).nodeHash
	err := newNode.Hash(prevHash[:])
	if err != nil {
		log.Println("deu ruim: ", err)
	}

	log.Printf("Adding node: ")
	log.Printf("\t Client id: [%v]", (*newNode).clientId)
	time := time.Unix(0, (*newNode).transactionTime).Format(time.UnixDate)
	log.Printf("\t Transaction time: [%v]", time)
	log.Printf("\t Node hash: [%x]", (*newNode).nodeHash)

	if failIn == 0 {
		newNode.nodeHash[0] += 1
		log.Println("Altering hash...")
	}

	chain.PushBack(newNode)
}

func isNodeHashIntegral(elem *list.Element) bool {
	prev := elem.Prev()
	if prev == nil {
		return true
	}

	nodePrev := prev.Value.(*ChainNode)
	nodeElem := elem.Value.(*ChainNode)

	nodeCopy := *nodeElem
	err := nodeCopy.Hash(nodePrev.nodeHash[:])
	if err != nil {
		log.Println("Falha na verificacao do hash: ", err)
	}

	return bytes.Equal(nodeCopy.nodeHash[:], nodeElem.nodeHash[:])
}

func AddTransaction(clientId uint32, transactionValue int64) ChainError {
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

func IsChainTainted() bool {
	cont := 1

	for elem := chain.Front(); elem != nil; elem = elem.Next() {
		if isNodeHashIntegral(elem) == false {
			log.Printf("\t Node [%v] hash tainted", cont)
			return true
		}
		cont += 1
	}

	return false
}

func GetClientBalance(clientId uint32) (int64, ChainError) {
	balance := int64(0)
	clientFound := false

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
