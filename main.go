package main

import (
	"fmt"

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
