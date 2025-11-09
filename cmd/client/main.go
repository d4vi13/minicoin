package main

import (
  "github.com/d4vi13/minicoin/internal/client"
  "github.com/d4vi13/minicoin/internal/api"
)

func main() {
  client.Start()

  var pkgTest api.Package
  pkgTest.PkgType = api.ClientRequestPkg
}
