package main

import (
	"blockchain/utils"
	"blockchain/wallet"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
)

const tempDir = "./templates"

type WalletServer struct {
	Port    uint16 `json:"Port"`
	Gateway string `json:"Gateway"`
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port, gateway}
}

func (ws *WalletServer) GetPort() uint16 {
	return ws.Port
}

func (ws *WalletServer) GetGateway() string {
	return ws.Gateway
}

func (ws *WalletServer) Index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		tmplFile := path.Join(tempDir, "index.html")
		t, err := template.ParseFiles(tmplFile)
		if err != nil {
			log.Printf("ERROR: Failed to parse template file %s: %v", tmplFile, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = t.Execute(w, nil)
		if err != nil {
			log.Printf("ERROR: Failed to execute template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		log.Printf("ERROR: Invalid HTTP method")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (ws *WalletServer) Wallet(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		myWallet := wallet.NewWallet()
		m, _ := json.Marshal(myWallet)
		io.WriteString(w, string(m[:]))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid http method")
	}
}
func (ws *WalletServer) CreateTransaction(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var t wallet.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		if !t.Validate() {
			log.Println("ERROR: missing field(s)")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		// publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		// privateKey := utils.PrivateKeyFromString(*t.SenderPrivateKey, publicKey)
		// value, err := strconv.ParseFloat(*t.Value, 32)
		// if err != nil {
		// 	log.Println("ERROR: parse error")
		// 	io.WriteString(w, string(utils.JsonStatus("fail")))
		// 	return
		// }
		// value32 := float32(value)

		// w.Header().Add("Content-Type", "application/json")
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid http method")
	}
}

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/transaction", ws.CreateTransaction)
	addr := "0.0.0.0:" + strconv.Itoa(int(ws.GetPort()))
	log.Printf("Wallet server running on http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
