package main

import (
	"github.com/Shaheen-AlQaraghuli/wallet-go/cmd/server"
)

// @title Wallet Service API.
// @version 1.0
// @description Wallet Service Documentation.
// @contact.name Wallet Service Owners
// @BasePath    /api/
func main(){
	server.StartServer()
}