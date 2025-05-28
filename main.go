package main

import (
	"github.com/nico-phil/transact/block"
)


func main(){
	blockChain := block.NewBlockchain()


	blockChain.CreateBlock(0, blockChain.LastBlock().Hash(), []string{"tr1", "tr2", "tr3"})

	blockChain.Print()

	
	
}