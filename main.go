package main

import (
	"github.com/nico-phil/transact/block"
)

func main() {
	blockChain := block.NewBlockchain()

	blockChain.CreateTransaction("recipeint address", "senderaddress", 1)


	// blockChain.CreateBlock(1, blockChain.LastBlock().Hash(), blockChain.TransactionPool)

	_ = blockChain.ProofOfWork()

	blockChain.Print()

}
