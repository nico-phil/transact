package main

import (
	"fmt"
	"net/http"
	"text/template"

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

func(ws *WalletServer) Run() error{
	router := http.NewServeMux()

	router.HandleFunc("/", ws.Index)
	router.HandleFunc("/wallet", ws.CreateWalletHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d",ws.Port), router)
}