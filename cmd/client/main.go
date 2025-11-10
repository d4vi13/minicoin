package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/d4vi13/minicoin/internal/client"
)

func Parse(addr *string, port *int, clientId *uint, action *int, value *string) error {
	const (
		ADDR_DEFAULT   = "localhost"
		PORT_DEFAULT   = 8080
		ID_DEFAULT     = 0
		ACTION_DEFAULT = -1
		VALUE_DEFAULT  = ""
	)

	actionDesc := fmt.Sprintf("define action: Transaction [%d] Check Blockchain"+
		"[%d] Get Balance [%d]", client.TRANSACTION, client.CHECK_BLOCKCHAIN,
		client.GET_BALANCE)

	flag.StringVar(addr, "addr", ADDR_DEFAULT, "define server address")
	flag.IntVar(port, "port", PORT_DEFAULT, "define server port")
	flag.UintVar(clientId, "id", ID_DEFAULT, "define client")
	flag.IntVar(action, "action", ACTION_DEFAULT, actionDesc)
	flag.StringVar(value, "value", VALUE_DEFAULT, "define transaction value")

	flag.Parse()

	if *clientId == ID_DEFAULT {
		return fmt.Errorf("Invalid client id!")
	}
	if *action == -1 {
		return fmt.Errorf("No action given!")
	}
	if *action >= client.MAX_ACTION {
		return fmt.Errorf("Invalid action given!")
	}
	if *action == client.TRANSACTION && *value == VALUE_DEFAULT {
		return fmt.Errorf("Transaction set and no value given!")
	}

	return nil
}

func main() {
	var addr string
	var action int
	var id uint
	var value int64
	var port int

	var valueString string
	err := Parse(&addr, &port, &id, &action, &valueString)
	if err != nil {
		flag.PrintDefaults()
		fmt.Println()
		log.Fatal(err)
	}

	value, err = strconv.ParseInt(valueString, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	var minicoinClient client.Client
	minicoinClient.Init(id, addr, port)
	minicoinClient.HandleAction(action, value)
}
