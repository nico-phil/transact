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

	sig, err:= t.GenerateSignature()
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("signature: %s", sig)

	blockchainTransaction := block.Transaction{
		SenderBlockchainAddress: t.SenderBlockchainAddress,
		RecipientBlockchainAddress: t.RecipientBlockchainAddress,
		Value: t.Value,
	}
	
	isVerify := blockChain.VerifyTransactionSignature(t.SenderPublicKey, sig, &blockchainTransaction )
	
	fmt.Println("isVerify", isVerify)

	

}
