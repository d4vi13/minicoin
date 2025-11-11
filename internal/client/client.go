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
	var res api.ServerResponse
	var err error

	defer client.connection.Close()

	res, err = client.request(value)
	if err != nil {
		log.Printf("Transaction failed %v", err)
		return
	}
	err = api.CheckServerResponse(res)
	if err != nil {
		log.Println(err)
		return
	}

	switch action {
	case TRANSACTION:
		log.Printf("Transaction successful and balance is [%v]", res.ClientBalance)
	case CHECK_BLOCKCHAIN:
		if res.IsBlockchainCorrupted == true {
			log.Println("BLOCKCHAIN CORRUPTED!")
		} else {
			log.Println("Blockchain is fine!")
		}
	case GET_BALANCE:
		log.Printf("Client [%d] with balance [%v]", client.identifier, res.ClientBalance)
	default:
		log.Println("Action is invalid")
	}
}

func (client *Client) request(value int64) (api.ServerResponse, error) {
	var req api.ClientRequest
	var res api.ServerResponse

	req.Type = api.ClientTransaction
	req.Identifier = client.identifier
	req.TransactionValue = value

	err := api.SendPackage(api.ClientRequestPkg, req, client.connection)
	if err != nil {
		return api.ServerResponse{}, fmt.Errorf("Failed to send transaction request %v", err)
	}

	err = api.RecvPackage(&res, client.connection)
	if err != nil {
		return api.ServerResponse{}, fmt.Errorf("Failed to send transaction request %v", err)
	}

	return res, nil
}
