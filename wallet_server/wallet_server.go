package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/nico-phil/transact/block"
	"github.com/nico-phil/transact/utils"
	"github.com/nico-phil/transact/wallet"
)

type wrapper map[string]any

type WalletServer struct {
	Port int
}

func NewWalletServer(port int) *WalletServer {
	return &WalletServer{Port:port }
}

func(ws *WalletServer) Index(w http.ResponseWriter, r *http.Request){
	t , err := template.ParseFiles("template/index.html", )
	if err != nil {
		fmt.Println("error", err)
	}
	t.Execute(w, "")
}

func(ws *WalletServer) CreateWalletHandler(w http.ResponseWriter, r *http.Request){
	newWallet, err := wallet.NewWallet()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, wrapper{"error": "failed to create wallet"})
	}
	data := wrapper{
		"private_key":  newWallet.PrivateKeyStr(),
		"public_key": newWallet.PulicKeyStr(),
		"blockchain_address": newWallet.BlockchainAddress,
	}

	utils.WriteJSON(w, http.StatusCreated, data)
}

func(ws *WalletServer) CreateTransactions(w http.ResponseWriter,  r *http.Request){
	var tr wallet.TransactionRequest
	err := utils.ReadJSON(r, &tr)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, wrapper{"message":"fail to read json"})
		return
	}

	senderPublicKey := utils.PublicKeyFromString(tr.SenderPublicKey)
	senderPrivateKey := utils.PrivateKeyFromString(tr.SenderPrivateKey, *senderPublicKey )
	
	
	newTransaction := wallet.NewTransaction(senderPublicKey, senderPrivateKey, 
		tr.SenderBlockchainAddress, tr.RecipientblockchainAddress, tr.Value)

	signature, err := newTransaction.GenerateSignature()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, wrapper{"error": "failed to generate signature"})
		return
	}

	blockchainTransaction :=  block.TransactionRequest {
		SenderBlockchainAddress: tr.SenderBlockchainAddress,
		RecipientblockchainAddress: tr.RecipientblockchainAddress,
		Value: tr.Value,
		SenderPublicKey: tr.SenderPublicKey,
		Signature: signature.String(),
	}


	utils.WriteJSON(w, http.StatusCreated, wrapper{"transaction": blockchainTransaction})
}

func(ws *WalletServer) Run() error{
	router := http.NewServeMux()

	router.HandleFunc("/", ws.Index)
	router.HandleFunc("/wallet", ws.CreateWalletHandler)
	router.HandleFunc("POST /transactions", ws.CreateTransactions)
	return http.ListenAndServe(fmt.Sprintf(":%d",ws.Port), router)
}