package main

import (
	"flag"
	"fmt"
	"github.com/d4vi13/minicoin/internal/client"
	"log"
	"strconv"
)

func Parse(addr *string, port *int, clientId *int, action *int, value *string) error {
	const (
		ADDR_DEFAULT   = "localhost"
		PORT_DEFAULT   = 8080
		ID_DEFAULT     = -1
		ACTION_DEFAULT = -1
		VALUE_DEFAULT  = ""
	)

	actionDesc := fmt.Sprintf("define action: Transaction [%d] Check Blockchain"+
		"[%d] Get Balance [%d]", client.TRANSACTION, client.CHECK_BLOCKCHAIN,
		client.GET_BALANCE)

	flag.StringVar(addr, "addr", ADDR_DEFAULT, "define server address")
	flag.IntVar(port, "port", PORT_DEFAULT, "define server port")
	flag.IntVar(clientId, "id", ID_DEFAULT, "define client")
	flag.IntVar(action, "action", ACTION_DEFAULT, actionDesc)
	flag.StringVar(value, "value", VALUE_DEFAULT, "define transaction value")

	flag.Parse()

	if *clientId == ID_DEFAULT {
		return fmt.Errorf("No id given!")
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
	var id int
	var value int
	var port int

	var valueString string
	err := Parse(&addr, &port, &id, &action, &valueString)
	if err != nil {
		flag.PrintDefaults()
    fmt.Println()
		log.Fatal(err)
	}

	value, err = strconv.Atoi(valueString)
	if err != nil {
		log.Fatal(err)
	}

	var minicoinClient client.Client
	minicoinClient.Init(id, addr, port)
	minicoinClient.HandleAction(action, value)
}
