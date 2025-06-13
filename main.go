package main

import (
	"fmt"
	"net"
	"os"

	"github.com/nico-phil/transact/block"
	"github.com/nico-phil/transact/wallet"
)

func main() {
	blockChain := block.NewBlockchain("")

	walletA, _ := wallet.NewWallet()
	walletB, _ := wallet.NewWallet()

	t := wallet.NewTransaction(walletA.PublicKey, walletA.PrivateKey, walletA.BlockchainAddress, walletB.BlockchainAddress, 1.0)

	sig, err := t.GenerateSignature()
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("signature: %s", sig)

	isAdded := blockChain.AddTransaction(t.RecipientBlockchainAddress, t.SenderBlockchainAddress, t.Value, t.SenderPublicKey, sig)

	if isAdded {
		blockChain.Mining()
	}

	blockChain.Print()

}

func GetHost() string {
	// defaultHost := "127.0.0.1"
	hostname, err := os.Hostname()
	if err != nil {
		return "127.0.0.1"
	}

	fmt.Println("hostname", hostname)
	address, err := net.LookupHost(hostname)
	if err != nil {
		return "127.0.0.1"
	}

	fmt.Println("address", address)
	return address[2]
}