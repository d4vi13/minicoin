package main

import (
	"flag"

	"github.com/d4vi13/minicoin/internal/server"
)

func Parse(port *int, failIn *int) {
	const (
		PORT_DEFAULT    = 8080
		FAIL_IN_DEFAULT = -1
	)

	flag.IntVar(port, "port", PORT_DEFAULT, "Set server port")
	flag.IntVar(failIn, "fail-in", FAIL_IN_DEFAULT, "Set node insertions until failure")
	flag.Parse()
}

func main() {
	var port int
	var failIn int

	Parse(&port, &failIn)
	server.Serve(port, failIn)
}
