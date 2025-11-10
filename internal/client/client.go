package client

import (
	"net"
)

const (
	TRANSACTION = iota
	CHECK_BLOCKCHAIN
	GET_BALANCE
	MAX_ACTION
)

type Client struct {
	Identifier uint
	Connection *net.TCPConn
}

func (client *Client) Init(clientId int, addr string, port int) {

}

func (client *Client) connect(serverAddress string) {

}

func (client *Client) HandleAction(action int, value int) {

}

func (client *Client) transaction(value int64) {

}

func (client *Client) isBlockchainCorrupted() bool {
	return false
}

func (client *Client) getMyBalance() int64 {
	return 0
}
