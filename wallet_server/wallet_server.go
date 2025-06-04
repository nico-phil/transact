package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/nico-phil/transact/block"
	"github.com/nico-phil/transact/utils"
	"github.com/nico-phil/transact/wallet"
)

type wrapper map[string]any

type WalletServer struct {
	Port   int
	Client *http.Client
}

func NewWalletServer(port int) *WalletServer {
	return &WalletServer{Port: port, Client: &http.Client{}}
}

func (ws *WalletServer) Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		fmt.Println("error", err)
	}
	t.Execute(w, "")
}

func (ws *WalletServer) CreateWalletHandler(w http.ResponseWriter, r *http.Request) {
	newWallet, err := wallet.NewWallet()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, wrapper{"error": "failed to create wallet"})
	}
	data := wrapper{
		"private_key":        newWallet.PrivateKeyStr(),
		"public_key":         newWallet.PulicKeyStr(),
		"blockchain_address": newWallet.BlockchainAddress,
	}

	utils.WriteJSON(w, http.StatusCreated, data)
}

func (ws *WalletServer) CreateTransactions(w http.ResponseWriter, r *http.Request) {
	var tr wallet.TransactionRequest
	err := utils.ReadJSON(r, &tr)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, wrapper{"message": "fail to read json"})
		return
	}

	senderPublicKey := utils.PublicKeyFromString(tr.SenderPublicKey)
	senderPrivateKey := utils.PrivateKeyFromString(tr.SenderPrivateKey, *senderPublicKey)

	newTransaction := wallet.NewTransaction(senderPublicKey, senderPrivateKey,
		tr.SenderBlockchainAddress, tr.RecipientblockchainAddress, tr.Value)

	signature, err := newTransaction.GenerateSignature()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, wrapper{"error": "failed to generate signature"})
		return
	}

	blockchainTransaction := block.TransactionRequest{
		SenderBlockchainAddress:    tr.SenderBlockchainAddress,
		RecipientblockchainAddress: tr.RecipientblockchainAddress,
		Value:                      tr.Value,
		SenderPublicKey:            tr.SenderPublicKey,
		Signature:                  signature.String(),
	}

	url := "http://localhost:3000/blockchain/transactions"
	m, _ := json.Marshal(blockchainTransaction)
	buf := bytes.NewBuffer(m)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, buf)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, wrapper{"error": "failed to generate req"})
		return
	}

	response, err := ws.Client.Do(req)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, wrapper{"error": "request failed"})
		return
	}

	if response.StatusCode != http.StatusCreated {
		utils.WriteJSON(w, http.StatusInternalServerError, wrapper{"error": "request failed", " status_code": response.StatusCode})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, wrapper{"message": "sucess"})
}

func(ws *WalletServer) GetBalance(w http.ResponseWriter, r *http.Request){
	url := ""
	req, err:= http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, wrapper{"error": "cannot create new req"})
		return
	}
	
	amountResponse, err := ws.Client.Do(req)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, wrapper{"error": "reqeut failed"})
		return
	}
	
	if amountResponse.StatusCode != http.StatusOK {
		utils.WriteJSON(w, http.StatusNotFound, wrapper{"error": "request failed with bad code status"})
		return
	}

	defer amountResponse.Body.Close()
	
	type AmountResponse struct {
		Amount float64 `json:"amount"`
	}

	var aR AmountResponse
	json.NewDecoder(amountResponse.Body).Decode(aR)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, wrapper{"error": "fail to read body"})
	}
	
	utils.WriteJSON(w, http.StatusCreated, wrapper{"amount": aR.Amount})
}

func (ws *WalletServer) Run() error {
	router := http.NewServeMux()

	router.HandleFunc("/", ws.Index)
	router.HandleFunc("/wallet", ws.CreateWalletHandler)
	router.HandleFunc("POST /wallet/transactions", ws.CreateTransactions)
	router.HandleFunc("POST /wallet/balance", ws.GetBalance)
	return http.ListenAndServe(fmt.Sprintf(":%d", ws.Port), router)
}
