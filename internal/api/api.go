package api

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"reflect"
	"unsafe"
)

// Package type identifiers
type PackageType int32
type ServerResponseType int32
type ServerFailType int32
type ClientRequestType int32

// Package type constants
const (
	ClientRequestPkg  PackageType = iota // Client-to-server request package
	ServerResponsePkg                    // Server-to-client response package
)

// Server response type constants
const (
	ServerSuccessResponse ServerResponseType = iota // Request completed successfully
	ServerFailedResponse                            // Request failed
)

// Server failure type constants
const (
	ServerNoFail         ServerFailType = iota // No failure occurred
	ServerClientUnknown                        // Client identifier not recognized
	ServerClientOverdraw                       // Client has insufficient balance
	BlockchainTainted                          // Blockchain integrity compromised
)

// Client request type constants
const (
	ClientCheckBalance             ClientRequestType = iota // Check client balance
	ClientTransaction                                       // Perform transaction
	ClientCheckBlockchainIntegrity                          // Verify blockchain integrity
)

// PackageHeader defines the header for all network packages
type PackageHeader struct {
	PkgType PackageType // Identifies the type of package being sent
}

// ClientRequest defines the structure for client requests to server
type ClientRequest struct {
	Type             ClientRequestType // Type of request being made
	Identifier       uint32            // Client identifier
	TransactionValue int64             // Transaction amount (if applicable)
}

// ServerResponse defines the structure for server responses to clients
type ServerResponse struct {
	Type                  ServerResponseType // Success or failure response
	FailType              ServerFailType     // Specific failure reason (if failed)
	ClientBalance         int64              // Current client balance
	IsBlockchainCorrupted bool               // Blockchain integrity status
}

// send serializes and sends any data structure over the network connection
func send(data any, conn net.Conn) error {
	buf := new(bytes.Buffer)

	// Serialize data to binary format using big-endian byte order
	err := binary.Write(buf, binary.BigEndian, data)
	if err != nil {
		return err
	}

	// Convert buffer to byte slice and send over network
	bytes := buf.Bytes()
	_, err = conn.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

// SendPackage sends a complete package with header and payload
func SendPackage(pkgType PackageType, payload any, conn net.Conn) error {
	var pkg PackageHeader

	// Set package type and send header
	pkg.PkgType = pkgType
	err := send(pkg, conn)
	if err != nil {
		return fmt.Errorf("Failed to send package header: %v", err)
	}

	// Send the actual payload data
	err = send(payload, conn)
	if err != nil {
		return fmt.Errorf("Failed to send payload: %v", err)
	}

	return nil
}

// recv reads exactly numB bytes from the network connection
func recv(conn net.Conn, numB uint32) ([]byte, error) {
	// Create buffer for reading
	byteBuffer := make([]byte, numB)

	var buffer []byte
	// Continue reading until we have exactly numB bytes
	for {
		_, err := conn.Read(byteBuffer)
		if err != nil {
			return nil, err
		}
		buffer = append(buffer, byteBuffer...)

		// Return when we have collected the exact number of bytes needed
		if uint32(len(buffer)) == numB {
			return buffer, nil
		}
	}
}

// RecvPackage receives a complete package and deserializes the payload
func RecvPackage(payload any, conn net.Conn) error {
	var pkg PackageHeader

	// Calculate package header size using unsafe (risky approach)
	pkgSize := uint32(unsafe.Sizeof(pkg))

	// Receive package header
	tmp, err := recv(conn, pkgSize)
	if err != nil {
		return fmt.Errorf("Failed to read package header")
	}

	// Deserialize package header from received bytes
	reader := bytes.NewReader(tmp)
	err = binary.Read(reader, binary.BigEndian, &pkg)
	if err != nil {
		return fmt.Errorf("Failed to write to pkg var")
	}

	// Use reflection to determine payload size
	reflection := reflect.ValueOf(payload)
	reflection = reflection.Elem()
	payloadSize := uint32(reflection.Type().Size())

	// Receive payload data
	tmp, err = recv(conn, payloadSize)
	if err != nil {
		return fmt.Errorf("Failed to read payload")
	}

	// Deserialize payload into the provided variable
	reader = bytes.NewReader(tmp)
	err = binary.Read(reader, binary.BigEndian, payload)
	if err != nil {
		return fmt.Errorf("Failed to write to payload var: %v", err)
	}

	return nil
}
