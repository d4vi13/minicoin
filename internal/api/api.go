package api

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"reflect"
	"unsafe"
)

const (
	DELIM = 0
)

type PackageType int32
type ServerResponseType int32
type ServerFailType int32
type ClientRequestType int32

const (
	ClientRequestPkg PackageType = iota
	ServerResponsePkg
)

const (
	ServerSuccessResponse ServerResponseType = iota
	ServerFailedResponse
)

const (
	ServerNoFail ServerFailType = iota
	ServerClientUnknown
	ServerClientOverdraw
	BlockchainTainted
)

const (
	ClientCheckBalance ClientRequestType = iota
	ClientTransaction
	ClientCheckBlockchainIntegrity
)

// Defines interface for communication
type PackageHeader struct {
	PkgType PackageType
}

// Defines interface for client request
type ClientRequest struct {
	Type             ClientRequestType
	Identifier       uint32
	TransactionValue int64
}

// Defines interface for server response
type ServerResponse struct {
	Type                  ServerResponseType
	FailType              ServerFailType
	ClientBalance         int64
	IsBlockchainCorrupted bool
}

func send(data any, conn net.Conn) error {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.BigEndian, data)
	if err != nil {
		return err
	}

	bytes := buf.Bytes()
	_, err = conn.Write(bytes)
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

func recv(conn net.Conn, numB uint32) ([]byte, error) {
	byteBuffer := make([]byte, numB)

	var buffer []byte
	for {
		_, err := conn.Read(byteBuffer)
		if err != nil {
			return nil, err
		}
		buffer = append(buffer, byteBuffer...)

		if uint32(len(buffer)) == numB {
			return buffer, nil
		}
	}
}

func RecvPackage(payload any, conn net.Conn) error {
	var pkg PackageHeader

	pkgSize := uint32(unsafe.Sizeof(pkg))
	tmp, err := recv(conn, pkgSize)
	if err != nil {
		return fmt.Errorf("Failed to read package header")
	}
	reader := bytes.NewReader(tmp)
	err = binary.Read(reader, binary.BigEndian, &pkg)
	if err != nil {
		return fmt.Errorf("Failed to write to pkg var")
	}

	reflection := reflect.ValueOf(payload)
	reflection = reflection.Elem()
	payloadSize := uint32(reflection.Type().Size())
	tmp, err = recv(conn, payloadSize)
	if err != nil {
		return fmt.Errorf("Failed to read payload")
	}

	reader = bytes.NewReader(tmp)
	err = binary.Read(reader, binary.BigEndian, payload)
	if err != nil {
		return fmt.Errorf("Failed to write to payload var: %v", err)
	}

	return nil
}
