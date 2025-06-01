package main

import (
	"fmt"
	"net/http"
)

type wrapper map[string]any

type BlockchainServer struct {
	Port int
}

func NewBlockchainServer(port int) *BlockchainServer {
	return &BlockchainServer{Port:port }
}

func(bs *BlockchainServer) GetBlockchainHandler(w http.ResponseWriter, r *http.Request){

}

func(bs *BlockchainServer) Run() error{
	router := http.NewServeMux()
	
	router.HandleFunc("/", bs.GetBlockchainHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", bs.Port), router)
}

