package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/nico-phil/transact/utils"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey *ecdsa.PublicKey
	BlockchainAddress string
}

func NewWallet() (*Wallet, error){
	w := new(Wallet)
	//1
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return &Wallet{}, err
	}

	w.PrivateKey = privateKey
	w.PublicKey = &privateKey.PublicKey

	//2. perform SHA-256 hashing on the public key(32 bytes)
	h2 := sha256.New()
	h2.Write(w.PublicKey.X.Bytes())
	h2.Write(w.PublicKey.Y.Bytes())
	digest2 := h2.Sum(nil)


	//3 
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)


	// 4
	v4 := make([]byte, 21)
	v4[0] = 0x00
	copy(v4[1:], digest3)

	// 5
	h5 := sha256.New()
	h5.Write(v4)
	digest5 := h3.Sum(nil)

	// 6
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)
	
	// 7
	chekSum := digest6[:4]


	// 8
	digest8 := make([]byte, 25)
	copy(digest8[:21], v4)
	copy(digest8[21:], chekSum)

	w.BlockchainAddress =  base58.Encode(digest8)
	
	return w, nil
}

func(w *Wallet) PulicKeyStr() string{
	return fmt.Sprintf("%x%x", w.PublicKey.X.Bytes(), w.PublicKey.Y.Bytes())
}

func(w *Wallet) PrivateKeyStr() string{
	return fmt.Sprintf("%x", w.PrivateKey.D.Bytes())
}

func(w *Wallet) Print(){
	fmt.Printf("public key: %x%x\n", w.PublicKey.X.Bytes(), w.PublicKey.Y.Bytes())
	fmt.Printf("private key: %x\n", w.PrivateKey.D.Bytes())
	fmt.Printf("blockchain address: %s\n", w.BlockchainAddress)
	
}

type Transaction struct {
	SenderPublicKey *ecdsa.PublicKey 
	SenderPrivateKey *ecdsa.PrivateKey 
	SenderBlockchainAddress string 
	RecipientBlockchainAddress string 
	Value float64 
}

func NewTransaction(
	sendPublicKey *ecdsa.PublicKey, 
	senderPrivateKey *ecdsa.PrivateKey, 
	senderBlockchainAddress string, 
	recipientblockchainAddress string, 
	value float64) *Transaction{


	return &Transaction{
		SenderPublicKey: sendPublicKey,
		SenderPrivateKey: senderPrivateKey,
		SenderBlockchainAddress: senderBlockchainAddress,
		RecipientBlockchainAddress: recipientblockchainAddress,
		Value: value,
	}

}

func(t *Transaction) MarshalJSON() ([]byte, error){
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

func(t *Transaction) GenerateSignature() (*utils.Signature, error){
	m, err := json.Marshal(t)
	if err != nil {
		return &utils.Signature{}, err
	}
	hash := sha256.Sum256(m)
	r, s, err := ecdsa.Sign(rand.Reader, t.SenderPrivateKey, hash[:])
	if err != nil {
		return &utils.Signature{}, err
	}

	signature := utils.Signature{R: r, S:s}
	return &signature, nil
}


func(t *Transaction) Print(){
	fmt.Printf("sender_public_key: %s", t.SenderPublicKey)
	fmt.Printf("sender_private_key: %s", t.SenderPrivateKey)
	fmt.Printf("sender_blockchain_address: %s", t.SenderBlockchainAddress)
	fmt.Printf("recipent_blockchain_address: %s", t.RecipientBlockchainAddress)
	fmt.Printf("value: %d", t.Value)
}

