package block

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Blockchain struct {
	Chain []*Block
}

func NewBlockchain() *Blockchain{
	block := &Block{}
	hash := block.Hash()
	block.PrevHash = hash
	return &Blockchain{
		Chain: []*Block{block},
	}
}

func( bc *Blockchain) LastBlock() *Block {
	return bc.Chain[len(bc.Chain) - 1]
}

func(bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte, transaction []string ){
	block :=  NewBlock(nonce, previousHash, transaction)
	bc.Chain = append(bc.Chain, block)
}

func(bc *Blockchain) Print(){
	for i, block := range bc.Chain {
		fmt.Printf("%s chain %d %s \n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
}


type Block struct {
	Nonce int
	PrevHash [32]byte
	Transactions []string
	Timestamp int64
}

func NewBlock(nonce int, prevHash [32]byte, transacrions []string) *Block {
	return &Block{
		Nonce: nonce,
		PrevHash: prevHash,
		Transactions: transacrions,
		Timestamp: time.Now().Unix(),
	}
}

func(b *Block) Print(){
	fmt.Printf("Nonce: %d \n", b.Nonce)
	fmt.Printf("Previous Hash: %x \n", b.PrevHash)
	fmt.Printf("transactions : %v \n", b.Transactions)
	fmt.Printf("%s  \n", strings.Repeat("=", 20))
}


func(b *Block) Hash() [32]byte{
	m, _ := json.Marshal(b)
	hash := sha256.Sum256(m)
	return hash

}





