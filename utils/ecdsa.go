package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func (s *Signature) String() string {
	return fmt.Sprintf("%x%x", s.R, s.S)
}

func PublicKeyFromString(s string) *ecdsa.PublicKey {
	bx, _ := hex.DecodeString(s[:64])
	by, _ := hex.DecodeString(s[64:])

	var bix big.Int
	var biy big.Int

	bix.SetBytes(bx)
	biy.SetBytes(by)

	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: &bix, Y: &biy}
}

func PrivateKeyFromString(s string, publicKey ecdsa.PublicKey) *ecdsa.PrivateKey {
	b, _ := hex.DecodeString(s[:])

	var d big.Int
	d.SetBytes(b)
	return &ecdsa.PrivateKey{PublicKey: publicKey, D: &d}
}
