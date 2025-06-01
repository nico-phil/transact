package main

import (
	"log"

	"github.com/nico-phil/transact/config"
)


func main(){
	WalletServer := NewWalletServer(config.GetServerPort())
	err := WalletServer.Run()
	if err != nil {
		log.Fatal("cannot start wallet server")
	}
}