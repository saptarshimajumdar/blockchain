package main

import (
	"blockchain/block"
	"blockchain/wallet"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

var cache map[string]*block.Blockchain = make(map[string]*block.Blockchain)

type BlockchainServer struct {
	port uint16 `json : "port"`
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port}
}

func (bcs *BlockchainServer) GetPort() uint16 {
	return bcs.port
}

func (bcs *BlockchainServer) GetBlockchain() *block.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		minersWallet := wallet.NewWallet()
		bc = block.NewBlockChain(minersWallet.GetBlockchainAddress(), bcs.GetPort())
		cache["blockchain"] = bc
		log.Printf("private_key: %v", minersWallet.PrivateKeyStr())
		log.Printf("public_key: %v", minersWallet.PublicKeyStr())
		log.Printf("blockchain_address: %v", minersWallet.GetBlockchainAddress())
	}
	return bc
}

func (bcs *BlockchainServer) GetChain(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockchain()
		m, _ := json.Marshal(bc)
		io.WriteString(w, string(m[:]))
	default:
		log.Printf("ERROR : Invalid http method")
	}
}

func (bcs *BlockchainServer) Run() {
	http.HandleFunc("/", bcs.GetChain)
	addr := "0.0.0.0:" + strconv.Itoa(int(bcs.GetPort()))
	log.Printf("Server is running on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
