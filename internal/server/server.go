package server

import (
	"fmt"
	"log"
	"net"

	"github.com/d4vi13/minicoin/internal/api"
	"github.com/d4vi13/minicoin/internal/chain"
)

func Serve(port int) {
	address := fmt.Sprintf(":%d", port)

	listener, err := net.Listen("tcp4", address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	chain.Init()

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
		res = handleTransaction(req.Identifier, req.TransactionValue)
	case api.ClientCheckBalance:
		log.Printf("Got check balance request, client [%d]", req.Identifier)
		res = handleCheckBalance(req.Identifier)
	case api.ClientCheckBlockchainIntegrity:
		log.Printf("Got check block chain request, client [%d]", req.Identifier)
		res = handleCheckBlockchain()
	default:
		log.Println("Request is not client transaction")
	}

	err = api.SendPackage(api.ServerResponsePkg, res, conn)
	if err != nil {
		log.Printf("Failed to send response %v", err)
	}
}

func handleTransaction(clientId uint, value int64) (api.ServerResponse) {
	var res api.ServerResponse

	chainErr := chain.AddTransaction(clientId, value)
	translateChainError(&res, chainErr)

	return res
}

func handleCheckBlockchain() api.ServerResponse {
	var res api.ServerResponse

	res.Type = api.ServerSuccessResponse
	res.IsBlockchainCorrupted = chain.IsChainTainted()
	if res.IsBlockchainCorrupted {
		log.Println("Blockchain is corrupted!")
	} else {
		log.Println("Blockchain is not corrupted!")
	}

	return res
}

func handleCheckBalance(clientId uint) api.ServerResponse {
	var res api.ServerResponse

	balance, chainErr := chain.GetClientBalance(clientId)
	translateChainError(&res, chainErr)
	res.ClientBalance = balance

	return res
}

func translateChainError(res *api.ServerResponse, chainErr chain.ChainError) {
	if chainErr == chain.SUCCESS {
		res.Type = api.ServerSuccessResponse
		res.FailType = api.ServerNoFail
		return
	}

	res.Type = api.ServerFailedResponse
	if (chainErr == chain.CLIENT_OVERDRAW) {
		res.FailType = api.ServerClientOverdraw
		return
	}
	if (chainErr == chain.CLIENT_NOT_FOUND) {
		res.FailType = api.ServerClientUnkown
		return
	}
	if (chainErr == chain.BLOCKCHAIN_TAINTED) {
		res.IsBlockchainCorrupted = true
		return
	}
}
