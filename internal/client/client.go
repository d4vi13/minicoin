// Package client implements client actions logic
package client

import (
	"fmt"
	"log"
	"net"

	"github.com/d4vi13/minicoin/internal/api"
)

// Client request action types
const (
	TRANSACTION = iota
	CHECK_BLOCKCHAIN
	GET_BALANCE
	MAX_ACTION
)

// Network connection type
const (
	TYPE = "tcp4"
)

// Client represents a blockchain client connection
type Client struct {
	identifier uint32
	connection *net.TCPConn
}

// Init initializes client connection to server
func (client *Client) Init(clientId uint32, name string, port int) {
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

// HandleAction processes different client actions
func (client *Client) HandleAction(action int, value int64) {
	var res api.ServerResponse
	var err error

	defer client.connection.Close()

	switch action {
	case TRANSACTION:
		res, err = client.request(api.ClientTransaction, value)
		if err != nil {
			log.Printf("Transaction failed %v", err)
			return
		}
		err = CheckServerResponse(res)
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("Transaction successful and balance is [%v]", res.ClientBalance)
	case CHECK_BLOCKCHAIN:
		res, err = client.request(api.ClientCheckBlockchainIntegrity, value)
		if err != nil {
			log.Printf("Transaction failed %v", err)
			return
		}
		if res.IsBlockchainCorrupted == true {
			log.Println("Blockchain corrupted!")
		} else {
			log.Println("Blockchain is fine!")
		}
	case GET_BALANCE:
		res, err = client.request(api.ClientCheckBalance, value)
		if err != nil {
			log.Printf("Transaction failed %v", err)
			return
		}
		err = CheckServerResponse(res)
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("Client [%d] with balance [%v]", client.identifier, res.ClientBalance)
	default:
		log.Println("Action is invalid")
	}
}

// request sends request to server and returns response
func (client *Client) request(typeReq api.ClientRequestType, value int64) (api.ServerResponse, error) {
	var req api.ClientRequest
	var res api.ServerResponse

	req.Type = typeReq
	req.Identifier = client.identifier
	req.TransactionValue = value

	err := api.SendPackage(api.ClientRequestPkg, req, client.connection)
	if err != nil {
		return api.ServerResponse{}, fmt.Errorf("Failed to send request: %v", err)
	}

	err = api.RecvPackage(&res, client.connection)
	if err != nil {
		return api.ServerResponse{}, fmt.Errorf("Failed to recv response: %v", err)
	}

	return res, nil
}

// CheckServerResponse validates server response and returns appropriate error
func CheckServerResponse(serverResp api.ServerResponse) error {
	if serverResp.Type == api.ServerSuccessResponse {
		return nil
	}

	if serverResp.IsBlockchainCorrupted {
		return fmt.Errorf("Blockchain corrupted!")
	}

	switch serverResp.FailType {
	case api.ServerNoFail:
		return fmt.Errorf("Server failed but no fail type was specified!")
	case api.ServerClientUnkown:
		return fmt.Errorf("Server did not recognize client!")
	case api.ServerClientOverdraw:
		return fmt.Errorf("Server returned not enough balance!")
	default:
		return fmt.Errorf("Unknown error in CheckServerResponse!")
	}
}
