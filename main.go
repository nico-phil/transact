package main

import (
	"github.com/nico-phil/transact/block"
)


func main(){
	blockChain := block.NewBlockchain()


	blockChain.CreateTransaction("recipeint address", "senderaddress", 1)
	blockChain.CreateTransaction("recipeint address", "senderaddress", 2)
	
	blockChain.CreateBlock(1, blockChain.LastBlock().Hash(), blockChain.TransactionPool)


	blockChain.Print()

	
	
}