package client

import (
	"net"
)

type Client struct {
	Identifier uint
	Connection *TCPConn
}

func Connect(serverAddress string) *TCPConn {

}

func connectByIp(serverIp string) *TCPConn {

}

func connectByName(hostName string) *TCPConn {

}

func (client *Client) Transaction(value int64) {

} 

func (client *Client) CheckBlockchainIntegrity() bool {

} 

func (client *Client) GetMyBalance() int64 {

}