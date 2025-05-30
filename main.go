package main

import (
	"github.com/nico-phil/transact/wallet"
)

func main() {
	// blockChain := block.NewBlockchain()

	walletA, _ := wallet.NewWallet()
	walletB, _ := wallet.NewWallet()

	t := wallet.NewTransaction(walletA.PublicKey, walletA.PrivateKey, walletA.BlockchainAddress,walletB.BlockchainAddress, 1.0 )

	t.GenerateSignature()

	// blockChain.CreateTransaction("recipeint address", "senderaddress", 1)
	// blockChain.Mining()

	// blockChain.Print()
	
	// fmt.Println("pool", len(blockChain.TransactionPool))

	

}
