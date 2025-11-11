package server

import (
	"fmt"
	"log"
	"net"

	"github.com/d4vi13/minicoin/internal/api"
)

func Serve(port int) {
	address := fmt.Sprintf(":%d", port)

	listener, err := net.Listen("tcp4", address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	var req api.ClientRequest
	var res api.ServerResponse

	log.Printf("Serving %s\n", conn.RemoteAddr().String())

	defer conn.Close()
	err := api.RecvPackage(&req, conn)
	if err != nil {
		log.Println("Failed to get client request!")
		return
	}

	switch req.Type {
	case api.ClientTransaction:
		log.Printf("Got transaction, value is [%v], client [%d]", req.TransactionValue, req.Identifier)
		res, err = handleTransaction(req.Identifier, req.TransactionValue)
		if err != nil {
			log.Printf("Failed to handle client transaction %v", err)
		}
		// case api.ClientCheckBalance:
		//
		// case api.ClientCheckBlockchainIntegrity:
	default:
		log.Println("Request is not client transaction")
	}

	err = api.SendPackage(api.ServerResponsePkg, res, conn)
	if err != nil {
		log.Printf("Failed to send response %v", err)
	}
}

func handleTransaction(clientId uint, value int64) (api.ServerResponse, error) {

	var res api.ServerResponse

	res.Type = api.ServerSuccessResponse
	res.FailType = api.ServerNoFail
	res.ClientBalance = 999

	return res, nil
}

func handleCheckBlockchain() (api.ServerResponse, error) {
	return api.ServerResponse{}, nil
}

func handleCheckBalance() (api.ServerResponse, error) {
	return api.ServerResponse{}, nil
}
