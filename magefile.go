// +build mage

package main

import (
  "github.com/magefile/mage/sh"
  "fmt"
)

func Client() error {
  fmt.Println("Building client...")
  return sh.Run("go", "build", "-o", "client", "cmd/client/main.go")
}

func Server() error {
  fmt.Println("Building server...")
  return sh.Run("go", "build", "-o", "server", "cmd/server/main.go")
}

func Clean() error {
  fmt.Println("Cleaning...")
  return sh.Run("rm", "-rf", "client", "server")
}
