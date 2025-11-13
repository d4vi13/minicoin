//go:build mage
// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
)

func All() error {
	err := Client()
	if err != nil {
		return err
	}
	return Server()
}

func Client() error {
	fmt.Println("Building client...")
	return sh.Run("go", "build", "-o", "client", "cmd/client/main.go")
}

func Server() error {
	fmt.Println("Building server...")
	return sh.Run("go", "build", "-o", "minicoin-server", "cmd/server/main.go")
}

func Clean() error {
	fmt.Println("Cleaning...")
	return sh.Run("rm", "-rf", "client", "minicoin-server")
}
