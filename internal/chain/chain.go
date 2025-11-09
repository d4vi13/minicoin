package chain

import (
	"crypto/sha512"
	"container/list"
	"time"
)

var chain List 

type ChainNode struct {
	clientId uint
	transactionTime time.Time
	transactionValue int64
	nodeHash [sha512.Size256]byte
}

func AddChainNode() error {

}

func isNodeHashCorrupted() bool {

}

func IsChainTainted() bool {

}

func GetClientBalance(clientId uint) int64 {

}
