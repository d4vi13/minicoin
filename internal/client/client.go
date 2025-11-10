package client

import (
	"fmt"
	"log"
	"net"

	"github.com/d4vi13/minicoin/internal/api"
)

const (
	TRANSACTION = iota
	CHECK_BLOCKCHAIN
	GET_BALANCE
	MAX_ACTION
)

const (
	TYPE = "tcp4"
)

type Client struct {
	identifier uint
	connection *net.TCPConn
}

func (client *Client) Init(clientId uint, name string, port int) {
	client.identifier = clientId

	serverAddr := fmt.Sprintf("%s:%d", name, port)
	tcpServer, err := net.ResolveTCPAddr(TYPE, serverAddr)
	if err != nil {
		log.Fatal("ResolveTCPAddr failed:", err.Error())
	}

	client.connection, err = net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		log.Fatal("Dial failed:", err.Error())
	}
}

func (client *Client) HandleAction(action int, value int64) {
	defer client.connection.Close()

	switch action {
	case TRANSACTION:
		var req api.ClientRequest

		req.Type = api.ClientTransaction
		req.Identifier = client.identifier
		req.TransactionValue = value
		api.SendPackage(api.ClientRequestPkg, req, client.connection)
	default:
		log.Println("Action is not transaction")
	}

}

func (client *Client) transaction(value int64) {

}

func (client *Client) isBlockchainCorrupted() bool {
	return false
}

func (client *Client) getMyBalance() int64 {
	return 0
}
