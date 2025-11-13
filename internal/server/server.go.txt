package server

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/d4vi13/minicoin/internal/api"
	"github.com/d4vi13/minicoin/internal/chain"
)

// Serve starts the blockchain server on specified port
func Serve(port int, failIn int) {
	address := fmt.Sprintf(":%d", port)

	// Create TCP listener on all interfaces
	listener, err := net.Listen("tcp4", address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// Initialize blockchain with failure threshold
	chain.Init(failIn)

	// Main server loop - accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		// Handle each client in separate goroutine for concurrency
		go handleClient(conn)
	}
}

// handleClient processes individual client connections
func handleClient(conn net.Conn) {
	var req api.ClientRequest
	var res api.ServerResponse

	log.Printf("Serving %s\n", conn.RemoteAddr().String())

	// Ensure connection is closed when function exits
	defer conn.Close()
	
	// Receive client request from connection
	err := api.RecvPackage(&req, conn)
	if err != nil {
		log.Println("Failed to get client request: ", err)
		return
	}

	// Check blockchain integrity before processing requests
	isBlockchainTainted := chain.IsChainTainted()
	
	// Only process non-integrity requests if blockchain is not tainted
	if !isBlockchainTainted || req.Type == api.ClientCheckBlockchainIntegrity {
		// Route request to appropriate handler based on type
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
		// Blockchain is tainted - reject all non-integrity check requests
		res.Type = api.ServerFailedResponse
		res.IsBlockchainCorrupted = true
		log.Println("Blockchain corrupted! Operation denied!")
	}

	os.Stderr.WriteString("\n")

	// Send response back to client
	err = api.SendPackage(api.ServerResponsePkg, res, conn)
	if err != nil {
		log.Printf("Failed to send response %v", err)
	}
}

// handleTransaction processes client transaction requests
func handleTransaction(clientId uint32, value int64) api.ServerResponse {
	var res api.ServerResponse

	log.Printf("Got transaction request:")
	log.Printf("\t Requesting Client: [%d]", clientId)
	log.Printf("\t Requested Value: [%d]", value)
	
	// Attempt to add transaction to blockchain
	chainErr := chain.AddTransaction(clientId, value)
	
	// Convert chain error to API response error
	translateChainError(&res, chainErr)

	switch res.FailType {
	case api.ServerClientOverdraw:
		log.Printf("\t Client [%d] does not have enough balance!", clientId)
	case api.ServerClientUnknown:
		log.Printf("\t Client [%d] not found!", clientId)
	case api.ServerNoFail:
		// Transaction successful - get updated balance
		balance, _ := chain.GetClientBalance(clientId)
		res.ClientBalance = balance
	}

	return res
}

// handleCheckBlockchain processes blockchain integrity check requests
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

// handleCheckBalance processes client balance inquiry requests
func handleCheckBalance(clientId uint32) api.ServerResponse {
	var res api.ServerResponse

	log.Printf("Got check balance request:")
	log.Printf("\t Requesting Client: [%d]", clientId)
	
	// Retrieve client balance from blockchain
	balance, chainErr := chain.GetClientBalance(clientId)
	
	// Convert chain error to API response error
	translateChainError(&res, chainErr)
	res.ClientBalance = balance
	
	if res.Type == api.ServerFailedResponse {
		log.Printf("\t Client [%d] not found!", clientId)
	} else {
		log.Printf("\t Client [%d] with balance [%d]", clientId, balance)
	}

	return res
}

// translateChainError converts chain errors to API response errors
func translateChainError(res *api.ServerResponse, chainErr chain.ChainError) {
	if chainErr == chain.SUCCESS {
		res.Type = api.ServerSuccessResponse
		res.FailType = api.ServerNoFail
		return
	}

	// Mark response as failed for all error cases
	res.Type = api.ServerFailedResponse
	
	// Map specific chain errors to API error types
	if chainErr == chain.CLIENT_OVERDRAW {
		res.FailType = api.ServerClientOverdraw
		return
	}
	if chainErr == chain.CLIENT_NOT_FOUND {
		res.FailType = api.ServerClientUnknown  
		return
	}
	if chainErr == chain.BLOCKCHAIN_TAINTED {
		res.IsBlockchainCorrupted = true
		return
	}
}
