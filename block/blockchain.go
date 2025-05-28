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
	TransactionPool []*Transaction
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

func(bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte, transaction []*Transaction ){
	block :=  NewBlock(nonce, previousHash, transaction)
	bc.Chain = append(bc.Chain, block)
	
	// update trnsactionPool
}

func(bc *Blockchain) CreateTransaction(recipientAddress, senderAddress string, value float64){
	transaction := NewTransaction(recipientAddress, senderAddress, value)
	bc.TransactionPool = append(bc.TransactionPool, transaction)
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
	Transactions []*Transaction
	Timestamp int64
}

func NewBlock(nonce int, prevHash [32]byte, transacrions []*Transaction) *Block {
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
	fmt.Printf("Timestamp : %d \n", b.Timestamp)
	for _, t := range b.Transactions {
		t.Print()
	}
}


func(b *Block) Hash() [32]byte{
	m, _ := json.Marshal(b)
	hash := sha256.Sum256(m)
	return hash
}

type Transaction struct {
	RecipientBlockchainAddress string
	SenderBlockchainAddress string
	Value float64
}

func NewTransaction(recipentAddress, senderAddress string, value float64) *Transaction {
	return &Transaction{
		RecipientBlockchainAddress: recipentAddress,
		SenderBlockchainAddress: senderAddress,
		Value: value,
	}
}

func( t *Transaction) Print(){
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf("recipient_address: %s \n", t.RecipientBlockchainAddress)
	fmt.Printf("sender_address: %s \n", t.SenderBlockchainAddress)
	fmt.Printf("value : %f \n", t.Value)
}





