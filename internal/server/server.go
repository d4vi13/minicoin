package server

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/d4vi13/minicoin/internal/api"
	"github.com/d4vi13/minicoin/internal/chain"
)

func Serve(port int, failIn int) {
	address := fmt.Sprintf(":%d", port)

	listener, err := net.Listen("tcp4", address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	chain.Init(failIn)

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
		log.Println("Failed to get client request: ", err)
		return
	}

	isBlockchainTainted := chain.IsChainTainted()
	if !isBlockchainTainted || req.Type == api.ClientCheckBlockchainIntegrity {
		switch req.Type {
		case api.ClientTransaction:
			res = handleTransaction(req.Identifier, req.TransactionValue)
		case api.ClientCheckBalance:
			res = handleCheckBalance(req.Identifier)
		case api.ClientCheckBlockchainIntegrity:
			res = handleCheckBlockchain(req.Identifier, isBlockchainTainted)
		default:
			log.Println("Request is not client transaction")
		}
	} else {
		res.Type = api.ServerFailedResponse
		res.IsBlockchainCorrupted = true
		log.Println("Blockchain corrupted! Operation denied!")
	}

	os.Stderr.WriteString("\n")

	err = api.SendPackage(api.ServerResponsePkg, res, conn)
	if err != nil {
		log.Printf("Failed to send response %v", err)
	}
}

func handleTransaction(clientId uint32, value int64) api.ServerResponse {
	var res api.ServerResponse

	log.Printf("Got transaction request:")
	log.Printf("\t Requesting Client: [%d]", clientId)
	log.Printf("\t Requested Value: [%d]", value)
	chainErr := chain.AddTransaction(clientId, value)
	translateChainError(&res, chainErr)

	switch res.FailType {
	case api.ServerClientOverdraw:
		log.Printf("\t Client [%d] does not have enough balance!", clientId)
	case api.ServerClientUnkown:
		log.Printf("\t Client [%d] not found!", clientId)
	case api.ServerNoFail:
		balance, _ := chain.GetClientBalance(clientId)
		res.ClientBalance = balance
	}

	return res
}

func handleCheckBlockchain(clientId uint32, blockchainTainted bool) api.ServerResponse {
	var res api.ServerResponse

	log.Printf("Got check blockchain request:")
	log.Printf("\t Requesting Client: [%d]", clientId)
	res.Type = api.ServerSuccessResponse
	res.IsBlockchainCorrupted = blockchainTainted
	if res.IsBlockchainCorrupted {
		log.Println("\t Blockchain is corrupted!")
	} else {
		log.Println("\t Blockchain is not corrupted!")
	}

	return res
}

func handleCheckBalance(clientId uint32) api.ServerResponse {
	var res api.ServerResponse

	log.Printf("Got check balance request:")
	log.Printf("\t Requesting Client: [%d]", clientId)
	balance, chainErr := chain.GetClientBalance(clientId)
	translateChainError(&res, chainErr)
	res.ClientBalance = balance
	if res.Type == api.ServerFailedResponse {
		log.Printf("\t Client [%d] not found!", clientId)
	} else {
		log.Printf("\t Client [%d] with balance [%d]", clientId, balance)
	}

	return res
}

func translateChainError(res *api.ServerResponse, chainErr chain.ChainError) {
	if chainErr == chain.SUCCESS {
		res.Type = api.ServerSuccessResponse
		res.FailType = api.ServerNoFail
		return
	}

	res.Type = api.ServerFailedResponse
	if chainErr == chain.CLIENT_OVERDRAW {
		res.FailType = api.ServerClientOverdraw
		return
	}
	if chainErr == chain.CLIENT_NOT_FOUND {
		res.FailType = api.ServerClientUnkown
		return
	}
	if chainErr == chain.BLOCKCHAIN_TAINTED {
		res.IsBlockchainCorrupted = true
		return
	}
}
