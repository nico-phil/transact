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
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_DIFFICULTY = 3
	MINER_REWARDS     = 1.0
)

type Blockchain struct {
	Chain             []*Block
	TransactionPool   []*Transaction
	BlockchainAddress string
}

func NewBlockchain(blockchainAddress string) *Blockchain {
	block := &Block{}
	hash := block.Hash()
	block.PrevHash = hash
	return &Blockchain{
		Chain:             []*Block{block},
		BlockchainAddress: blockchainAddress,
	}
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte, transaction []*Transaction) {
	block := NewBlock(nonce, previousHash, transaction)
	bc.Chain = append(bc.Chain, block)

	// update trnsactionPool
	bc.TransactionPool = []*Transaction{}
}

func (bc *Blockchain) AddTransaction(recipientAddress, senderAddress string, value float64, senderPublickey *ecdsa.PublicKey, signature *utils.Signature) bool {
	transaction := NewTransaction(recipientAddress, senderAddress, value)

	if senderAddress == MINING_SENDER {
		bc.TransactionPool = append(bc.TransactionPool, transaction)
		return true
	}

	isVerify := bc.VerifyTransactionSignature(senderPublickey, signature, transaction)
	if !isVerify {
		return false
	}

	bc.TransactionPool = append(bc.TransactionPool, transaction)
	return true
}

func (bc *Blockchain) ProofOfWork() int {
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

	return nonce

}

func (bc *Blockchain) Mining() bool {
	if len(bc.TransactionPool) == 0 {
		return false
	}

	isAdded := bc.AddTransaction(bc.BlockchainAddress, MINING_SENDER, MINER_REWARDS, nil, nil)
	if !isAdded {
		return false
	}
	
	nonce := bc.ProofOfWork()
	if nonce == 0 {
		return false
	}

	bc.CreateBlock(nonce, bc.LastBlock().PrevHash, bc.TransactionPool)
	return true
}

func(bc *Blockchain) StartMining(){
	bc.Mining()
	time.AfterFunc(time.Second * 30, bc.StartMining)
}

func (bc *Blockchain) VerifyTransactionSignature(senderPublickey *ecdsa.PublicKey, signature *utils.Signature, transaction *Transaction) bool {
	m, _ := json.Marshal(transaction)
	h := sha256.Sum256(m)
	isVerify := ecdsa.Verify(senderPublickey, h[:], signature.R, signature.S)
	return isVerify
}

func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float64 {
	var value float64 = 0
	for _, c := range bc.Chain {
		for _, t := range c.Transactions {
			if t.SenderBlockchainAddress == blockchainAddress {
				value = value - t.Value
			}

			if t.RecipientBlockchainAddress == blockchainAddress {
				value = value + t.Value
			}
		}
	}

	return value
}

func (bc *Blockchain) Print() {
	for i, block := range bc.Chain {
		fmt.Printf("%s chain %d %s \n", strings.Repeat("=", 40), i, strings.Repeat("=", 40))
		block.Print()
	}
}

type Block struct {
	Nonce        int
	PrevHash     [32]byte
	Transactions []*Transaction
	Timestamp    int64
}

func NewBlock(nonce int, prevHash [32]byte, transacrions []*Transaction) *Block {
	return &Block{
		Nonce:        nonce,
		PrevHash:     prevHash,
		Transactions: transacrions,
		Timestamp:    time.Now().Unix(),
	}
}

func (b *Block) MarshalJSON() ([]byte, error) {
	var bl = struct {
		Nonce        int            `json:"nonce"`
		PrevHash     string         `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
		Timestamp    int64          `json:"timestamp"`
	}{
		Nonce:        b.Nonce,
		PrevHash:     fmt.Sprintf("%x", b.PrevHash),
		Transactions: b.Transactions,
		Timestamp:    b.Timestamp,
	}

	return json.Marshal(bl)
}

func (b *Block) Print() {
	fmt.Printf("Nonce: %d \n", b.Nonce)
	fmt.Printf("Previous Hash: %x \n", b.PrevHash)
	fmt.Printf("Timestamp : %d \n", b.Timestamp)
	for _, t := range b.Transactions {
		t.Print()
	}
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	hash := sha256.Sum256(m)
	return hash
}

type Transaction struct {
	RecipientBlockchainAddress string
	SenderBlockchainAddress    string
	Value                      float64
}

func NewTransaction(recipentAddress, senderAddress string, value float64) *Transaction {
	return &Transaction{
		RecipientBlockchainAddress: recipentAddress,
		SenderBlockchainAddress:    senderAddress,
		Value:                      value,
	}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	var tr = struct {
		SenderBlockchainAddress    string  `json:"sender_blockchain_address"`
		RecipientBlockchainAddress string  `json:"recipient_blockchain_address"`
		Value                      float64 `json:"value"`
	}{
		SenderBlockchainAddress:    t.SenderBlockchainAddress,
		RecipientBlockchainAddress: t.RecipientBlockchainAddress,
		Value:                      t.Value,
	}

	m, err := json.Marshal(tr)
	return m, err
}

func (t *Transaction) Print() {
	fmt.Printf("     %s transactions %s \n", strings.Repeat("-", 20), strings.Repeat("-", 20))
	fmt.Printf("recipient_address: %s \n", t.RecipientBlockchainAddress)
	fmt.Printf("sender_address: %s \n", t.SenderBlockchainAddress)
	fmt.Printf("value : %f \n", t.Value)
}

type TransactionRequest struct {
	SenderBlockchainAddress    string  `json:"sender_blockchain_address"`
	RecipientblockchainAddress string  `json:"recipient_blockchain_address"`
	Value                      float64 `json:"value"`
	SenderPublicKey            string  `json:"sender_public_key"`
	Signature                  string  `json:"signature"`
}
