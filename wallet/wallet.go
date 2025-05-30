package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"github.com/btcsuite/btcutil/base58"
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

