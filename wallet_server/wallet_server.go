package main

import (
	"fmt"
	"net/http"
)

type WalletServer struct {
	Port int
}

func NewWalletServer(port int) *WalletServer {
	return &WalletServer{Port:port }
}

func(ws *WalletServer) HelloHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("hello from wallet server")
}

func(ws *WalletServer) Run() error{
	router := http.NewServeMux()

	router.HandleFunc("/", ws.HelloHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d",ws.Port), router)
}