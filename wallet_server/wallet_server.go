package main

import (
	"fmt"
	"net/http"
	"text/template"
)

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

func(ws *WalletServer) Run() error{
	router := http.NewServeMux()

	router.HandleFunc("/", ws.Index)
	return http.ListenAndServe(fmt.Sprintf(":%d",ws.Port), router)
}