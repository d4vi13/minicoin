package chain

import (
	"container/list"
	"crypto/sha512"
	"encoding/binary"
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

type NodeContent struct {
	clientId         uint
	transactionTime  time.Time
	transactionValue int64
}

type ChainNode struct {
	nodeContent NodeContent
	nodeHash    [sha512.Size256]byte
}

func Init() error {
	chain.Init()

	landmarkNode := make(ChainNode)
	chain.PushBack(landmarkNode)

	byteBuffer := make([]byte, unsafe.Sizeof(landmarkNode.nodeContent))
	_, err := binary.Encode(byteBuffer, binary.LittleEndian, landmarkNode.nodeContent)
	if err != nil {
		return err
	}

}

func (node *ChainNode) Hash() {
	copy(node.nodeHash[:])
}

func addChainNode(clientId uint, transactionValue int64) ChainError {
	newNode := make(ChainNode)

	if transactionValue < 0 {
		balance, err := GetClientBalance(clientId)
		if err != SUCCESS {
			return err
		}
		if (balance + transactionValue) < 0 {
			return CLIENT_OVERDRAW
		}
	}

	(*newNode).clientId = clientId
	(*newNode).transactionTime = time.Now()
	(*newNode).transactionValue = transactionValue

	chain.PushBack()
}

func isNodeHashCorrupted() bool {

}

func IsChainTainted() bool {

}

func GetClientBalance(clientId uint) (int64, ChainError) {

}
