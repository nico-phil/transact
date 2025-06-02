package main

import (
	"log"

	"github.com/nico-phil/transact/config"
)

func main() {
	blockchainServer := NewBlockchainServer(config.GetServerPort())
	err := blockchainServer.Run()
	if err != nil {
		log.Fatal("cannot start blockchain server")
	}
}
