package api

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/d4vi13/minicoin/internal/api"
)

const (
	DELIM = 0
)

type PackageType int

const (
	ClientRequestPkg PackageType = iota
	ServerResponsePkg
)

type ServerResponseType int

const (
	ServerSuccessResponse ServerResponseType = iota
	ServerFailedResponse
)

type ServerFailType int

const (
	ServerNoFail ServerFailType = iota
	ServerClientUnkown
	ServerClientOverdraw
)

type ClientRequestType int

const (
	ClientCheckBalance ClientRequestType = iota
	ClientTransaction
	ClientCheckBlockchainIntegrity
)

// Defines interface for communication
type PackageHeader struct {
	PkgType PackageType `json:"pkgType"`
}

// Defines interface for client request
type ClientRequest struct {
	Type             ClientRequestType `json:"type"`
	Identifier       uint              `json:"identifier"`
	TransactionValue int64             `json:"transactionValue"`
}

// Defines interface for server response
type ServerResponse struct {
	Type                  ServerResponseType `json:"type"`
	FailType              ServerFailType     `json:"failType"`
	ClientBalance         int64              `json:"clientBalance"`
	IsBlockchainCorrupted bool               `json:"isBlockchainCorrupted"`
}

func CheckServerResponse(serverResp ServerResponse) error {
	if serverResp.Type == ServerSuccessResponse {
		return nil
	}

	switch serverResp.FailType {
	case api.ServerNoFail:
		return fmt.Errorf("Server failed but no fail type was specified")
	case api.ServerClientUnkown:
		return fmt.Errorf("Server did not recognize client")
	case api.ServerClientOverdraw:
		return fmt.Errorf("Server returned not enough balance")
	default:
		return fmt.Errorf("Unknown error in CheckServerResponse")
	}
}

func send(data any, conn net.Conn) error {
	tmp, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Write(tmp)
	if err != nil {
		return err
	}

	byteBuffer := make([]byte, 1)
	byteBuffer[0] = DELIM
	_, err = conn.Write(byteBuffer)
	if err != nil {
		return err
	}

	return nil
}

func SendPackage(pkgType PackageType, payload any, conn net.Conn) error {
	var pkg PackageHeader

	pkg.PkgType = pkgType
	err := send(pkg, conn)
	if err != nil {
		return fmt.Errorf("Failed to send package header: %v", err)
	}

	err = send(payload, conn)
	if err != nil {
		return fmt.Errorf("Failed to send payload: %v", err)
	}

	return nil
}

func recv(conn net.Conn) ([]byte, error) {
	var buffer []byte
	byteBuffer := make([]byte, 1)

	for {
		n, err := conn.Read(byteBuffer)
		if err != nil {
			return nil, err
		}

		if n > 0 {
			if byteBuffer[0] == DELIM {
				return buffer, nil
			}

			buffer = append(buffer, byteBuffer[0])
		}
	}
}

func RecvPackage(payload any, conn net.Conn) error {
	var pkg PackageHeader

	tmp, err := recv(conn)
	if err != nil {
		return fmt.Errorf("Failed to read package header")
	}

	err = json.Unmarshal(tmp, &pkg)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal header")
	}

	tmp, err = recv(conn)
	if err != nil {
		return fmt.Errorf("Failed to read payload")
	}

	err = json.Unmarshal(tmp, payload)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal payload")
	}

	return nil
}
