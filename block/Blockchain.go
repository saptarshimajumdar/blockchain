package block

import (
	"blockchain/utils"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
	MINING_TIMER_SEC  = 10
)

type Block struct {
	Nonce        int            `json:"nonce"`
	PrevHash     [32]byte       `json:"prev_hash"`
	Timestamp    int64          `json:"timestamp"`
	Transactions []*Transaction `json:"transactions"`
}

type Blockchain struct {
	TransactionPool   []*Transaction `json:"TransactionPool"`
	Chain             []*Block       `json:"Chain"`
	BlockchainAddress string         `json:"BlockchainAddress"`
	Port              uint16
	mux               sync.Mutex
}

type Transaction struct {
	SenderBlockchainAddress    string  `json:"SenderBlockchainAddress"`
	RecipientBlockchainAddress string  `json:"RecipientBlockchainAddress"`
	Value                      float32 `json:"Value"`
}

type TransactionRequest struct {
	SenderBlockchainAddress    *string  `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string  `json:"recipient_blockchain_address"`
	SenderPublicKey            *string  `json:"sender_public_key"`
	Value                      *float32 `json:"value"`
	Signature                  *string  `json:"signature"`
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

func (t *Transaction) PrintTransaction() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_blockchain_address      %s\n", t.SenderBlockchainAddress)
	fmt.Printf(" recipient_blockchain_address   %s\n", t.RecipientBlockchainAddress)
	fmt.Printf(" value                          %.1f\n", t.Value)
}

func (bc *Blockchain) GetTransactionPool() []*Transaction {
	return bc.TransactionPool
}

func (b *Block) PrintBlock() {
	fmt.Printf("Nonce: %d\n", b.Nonce)
	fmt.Printf("PrevHash: %x\n", b.PrevHash)
	fmt.Printf("Timestamp: %d\n", b.Timestamp)
	fmt.Println("Transactions:")
	for _, tx := range b.Transactions {
		tx.PrintTransaction()
	}
	fmt.Println("__________________________________________________________________")
}

func (bc *Blockchain) PrintBlockchain() {
	fmt.Println("Blockchain:")
	for _, block := range bc.Chain {
		block.PrintBlock()
	}
}
func (bc *Blockchain) CreateBlock(nonce int, prevHash [32]byte) *Block {
	b := NewBlock(nonce, prevHash, bc.TransactionPool)
	bc.Chain = append(bc.Chain, b)
	bc.TransactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) CreateTransaction(sender string, recipient string, value float32,
	senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	isTransacted := bc.AddTransaction(sender, recipient, value, senderPublicKey, s)

	return isTransacted
}

// sender = sender address etc.
func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32,
	senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := Transaction{sender, recipient, value}

	if sender == MINING_SENDER {
		bc.TransactionPool = append(bc.TransactionPool, &t)
		return true
	}

	if bc.VerifyTransactionSignature(senderPublicKey, s, &t) {
		// if bc.CalculateTotalAmount(sender) < value {
		// 	log.Println("ERROR: Not enough Balance in wallet")
		// 	return false
		// }
		bc.TransactionPool = append(bc.TransactionPool, &t)
		return true
	} else {
		log.Println("ERROR: Verification of Transaction Failed")
	}
	return false
}

func (bc *Blockchain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey,
	s *utils.Signature, t *Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)

}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.TransactionPool {
		transactions = append(transactions,
			NewTransaction(t.SenderBlockchainAddress,
				t.RecipientBlockchainAddress,
				t.Value))
	}
	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, prevHash [32]byte, txns []*Transaction, difficulty int) bool {
	zeroes := strings.Repeat("0", difficulty)
	guessBlock := Block{nonce, prevHash, time.Now().UnixMilli(), txns}
	guessHash := fmt.Sprintf("%x", guessBlock.Hash())
	// fmt.Println(guessHash)
	return guessHash[:difficulty] == zeroes
}

// func to get the nonce value by trial and error
func (bc *Blockchain) ProofOfWork() int {
	txns := bc.CopyTransactionPool()
	prevHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, prevHash, txns, MINING_DIFFICULTY) {
		nonce++
	}
	return nonce
}

func (bc *Blockchain) Mining() bool {

	bc.mux.Lock()
	defer bc.mux.Unlock()

	if len(bc.TransactionPool) == 0 {
		return false
	}

	//while rewarding the miner there is no transaction
	bc.AddTransaction(MINING_SENDER, bc.BlockchainAddress, MINING_REWARD, nil, nil)
	nonce := bc.ProofOfWork()
	prevHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, prevHash)
	log.Println("action=Mining, status=success")
	return true
}

func (bc *Blockchain) StartMining() {
	bc.Mining()
	_ = time.AfterFunc(time.Second*MINING_TIMER_SEC, bc.StartMining)
}

// total transactions for the bcAdress node
func (bc *Blockchain) CalculateTotalAmount(bcAddress string) float32 {
	var amt float32 = 0.0
	for _, b := range bc.Chain {
		for _, t := range b.Transactions {
			v := t.Value
			if bcAddress == t.RecipientBlockchainAddress {
				amt += v
			} else if bcAddress == t.SenderBlockchainAddress {
				amt -= v
			}
		}
	}
	return amt
}

func NewBlock(nonce int, prevHash [32]byte, txns []*Transaction) *Block {
	return &Block{
		Timestamp:    time.Now().UnixMilli(),
		Nonce:        nonce,
		PrevHash:     prevHash,
		Transactions: txns,
	}
}

func NewBlockChain(BlockchainAddress string, port uint16) *Blockchain {
	b := new(Block)
	bc := new(Blockchain)
	bc.BlockchainAddress = BlockchainAddress
	bc.CreateBlock(0, b.Hash())
	bc.Port = port
	return bc
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (b *Block) Hash() [32]byte {
	m, err := json.Marshal(b)
	// fmt.Println(string(m))
	if err != nil {
		return [32]byte{}
	}
	return sha256.Sum256([]byte(m))
}

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderBlockchainAddress == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.SenderPublicKey == nil ||
		tr.Value == nil ||
		tr.Signature == nil {
		return false
	}
	return true
}
