package block

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/nico-phil/transact/utils"
)

const (
	MINING_SENDER = "THE BLOCKCHAIN"
	MINING_DIFFICULTY = 3
	MINER_REWARDS = 1.0

)

type Blockchain struct {
	Chain []*Block
	TransactionPool []*Transaction
	BlockchainAddress string
	
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
	bc.TransactionPool = []*Transaction{}
}

func(bc *Blockchain) AddTransaction(recipientAddress, senderAddress string, value float64, senderPublickey *ecdsa.PublicKey, signature *utils.Signature) bool{
	transaction := NewTransaction(recipientAddress, senderAddress, value)

	isVerify :=  bc.VerifyTransactionSignature(senderPublickey, signature, transaction )
	if !isVerify {
		return false
	}

	bc.TransactionPool = append(bc.TransactionPool, transaction)
	return true
}

func(bc *Blockchain) ProofOfWork() int{
	nonce := 0
	zeros := strings.Repeat("0", 3)
	guessBlock := Block{Nonce: nonce, PrevHash: bc.LastBlock().PrevHash, 
		Transactions: bc.TransactionPool, Timestamp: time.Now().Unix()}
	m, _ := json.Marshal(guessBlock)

	guessBlockHash := sha256.Sum256(m)
	guessBlockStr := fmt.Sprintf("%x", guessBlockHash)
	for guessBlockStr[:3] != zeros {
		nonce++
		guessBlock.Nonce = nonce
		m, _ = json.Marshal(guessBlock)
		guessBlockHash = sha256.Sum256(m)
		guessBlockStr = fmt.Sprintf("%x", guessBlockHash)

	}

	bc.CreateBlock(nonce, guessBlock.PrevHash, guessBlock.Transactions)

	return nonce

}

func(bc *Blockchain) Mining() bool{
	if len(bc.TransactionPool) ==0 {
		return false
	}

	// bc.AddTransaction("THE MINER BLOCKCHAIN ADDRESS", "THE BLOCKCHAIN", MINER_REWARDS, )
	_ = bc.ProofOfWork()


	return true
	
	//send reward to miner
	//remove money from user A
	// send money to user B
}

func(bc *Blockchain) VerifyTransactionSignature(senderPublickey *ecdsa.PublicKey, signature *utils.Signature, transaction *Transaction) bool{
	m, _ := json.Marshal(transaction)
	h := sha256.Sum256(m)
	isVerify := ecdsa.Verify(senderPublickey, h[:], signature.R, signature.S)
	fmt.Println("transaction verification:", isVerify)
	return isVerify
}

func(bc *Blockchain) Print(){
	for i, block := range bc.Chain {
		fmt.Printf("%s chain %d %s \n", strings.Repeat("=", 40), i, strings.Repeat("=", 40))
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

func(t *Transaction) MarshalJSON()([]byte, error){
	var tr = struct{
		SenderBlockchainAddress string `json:"sender_blockchain_address"`
		RecipientBlockchainAddress string `json:"recipient_blockchain_address"`
		Value float64 `json:"value"`
	}{
		SenderBlockchainAddress: t.SenderBlockchainAddress,
		RecipientBlockchainAddress: t.RecipientBlockchainAddress,
		Value: t.Value, 
	}
	
	m, err := json.Marshal(tr)
	return m, err
}


func( t *Transaction) Print(){
	fmt.Printf("     %s transactions %s \n", strings.Repeat("-", 20), strings.Repeat("-", 20))
	fmt.Printf("recipient_address: %s \n", t.RecipientBlockchainAddress)
	fmt.Printf("sender_address: %s \n", t.SenderBlockchainAddress)
	fmt.Printf("value : %f \n", t.Value)
}





