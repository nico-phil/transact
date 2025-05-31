package main

import (
	"fmt"

	"github.com/nico-phil/transact/block"
	"github.com/nico-phil/transact/wallet"
)

func main() {
	blockChain := block.NewBlockchain()

	walletA, _ := wallet.NewWallet()
	walletB, _ := wallet.NewWallet()

	t := wallet.NewTransaction(walletA.PublicKey, walletA.PrivateKey, walletA.BlockchainAddress,walletB.BlockchainAddress, 1.0 )

	sig, _ := t.GenerateSignature()

	blockchainTransaction := block.Transaction{
		SenderBlockchainAddress: t.SenderBlockchainAddress,
		RecipientBlockchainAddress: t.RecipientBlockchainAddress,
		Value: t.Value,
	}
	
	isVerify := blockChain.VerifyTransactionSignature(t.SenderPublicKey, sig, &blockchainTransaction )

	// blockChain.CreateTransaction("recipeint address", "senderaddress", 1)
	// blockChain.Mining()

	// blockChain.Print()
	
	fmt.Println("isVerify", isVerify)

	

}
