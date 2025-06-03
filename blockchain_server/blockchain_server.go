package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nico-phil/transact/block"
	"github.com/nico-phil/transact/utils"
	"github.com/nico-phil/transact/wallet"
)

var cache = map[string]*block.Blockchain{}

type wrapper map[string]any

type BlockchainServer struct {
	Port int
}

func NewBlockchainServer(port int) *BlockchainServer {
	return &BlockchainServer{Port: port}
}

func (bs *BlockchainServer) GetBlockchain() *block.Blockchain {
	blockchain, ok := cache["blockchain"]
	if !ok {
		w, _ := wallet.NewWallet()
		blockchain = block.NewBlockchain(w.BlockchainAddress)
		log.Printf("miner_wallet_private_key %v", w.PrivateKeyStr())
		log.Printf("miner_wallet_public_key %v", w.PulicKeyStr())
		log.Printf("miner_blockchain_address %v", w.BlockchainAddress)
	}

	return blockchain
}

func (bs *BlockchainServer) GetchainsHandler(w http.ResponseWriter, r *http.Request) {
	bc := bs.GetBlockchain()

	data := wrapper{
		"chain": bc.Chain,
	}
	utils.WriteJSON(w, http.StatusOK, data)

}

func (bs *BlockchainServer) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	blockchain := bs.GetBlockchain()
	var tr block.TransactionRequest
	err := utils.ReadJSON(r, &tr)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, wrapper{"error": "cannot read json"})
		return
	}

	publicKey := utils.PublicKeyFromString(tr.SenderPublicKey)
	signature := utils.SignatureFromString(tr.Signature)
	
	isAdded := blockchain.AddTransaction(tr.RecipientblockchainAddress, tr.SenderBlockchainAddress, tr.Value, publicKey, signature)
	if(!isAdded) {
		utils.WriteJSON(w, http.StatusInternalServerError, wrapper{"error": "cannot add transaction to transaction poll"})
		return
	}

	fmt.Println("in blockchain:", isAdded)
	fmt.Println("transactionPool:", blockchain.TransactionPool)

	
	utils.WriteJSON(w, http.StatusCreated, wrapper{"tractions": tr})

}

func(bs *BlockchainServer) GetTransactions(w http.ResponseWriter, r *http.Request){
	blockchain := bs.GetBlockchain()
	transactions := blockchain.TransactionPool
	if len(transactions) == 0 {
		transactions =[]*block.Transaction{}
	}

	utils.WriteJSON(w, http.StatusOK, wrapper{"transactions": transactions})
}

func (bs *BlockchainServer) Run() error {
	router := http.NewServeMux()

	router.HandleFunc("/chains", bs.GetchainsHandler)
	router.HandleFunc("POST /blockchain/transactions", bs.CreateTransaction)
	router.HandleFunc("GET /blockchain/transactions", bs.GetTransactions)
	return http.ListenAndServe(fmt.Sprintf(":%d", bs.Port), router)
}
