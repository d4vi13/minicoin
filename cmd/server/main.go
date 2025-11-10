package main

import (
	"flag"

	"github.com/d4vi13/minicoin/internal/server"
)

func Parse(port *int) {
	const (
		PORT_DEFAULT = 8080
	)

	flag.IntVar(port, "port", PORT_DEFAULT, "Set server port")
}

func main() {
	var port int

	Parse(&port)
	server.Serve(port)
}
