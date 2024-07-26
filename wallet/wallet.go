package wallet

import (
	"blockchain/utils"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	PrivateKey        *ecdsa.PrivateKey `json:"private_key"`
	PublicKey         *ecdsa.PublicKey  `json:"public_key"`
	BlockchainAddress string            `json:"blockchain_address"`
}

type Transaction struct {
	senderPrivateKey           *ecdsa.PrivateKey
	senderPublicKey            *ecdsa.PublicKey
	SenderBlockchainAddress    string  `json:"SenderBlockchainAddress"`
	RecipientBlockchainAddress string  `json:"RecipientBlockchainAddress"`
	Value                      float32 `json:"Value"`
}

type TransactionRequest struct {
	SenderPrivateKey           *string `json:"sender_private_key"`
	SenderBlockchainAddress    *string `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string `json:"recipient_blockchain_address"`
	SenderPublicKey            *string `json:"sender_public_key"`
	Value                      *string `json:"value"`
}

func NewTransaction(privKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey,
	senderBlockchainAddress string, recipientBlockchainAddress string, value float32) *Transaction {
	return &Transaction{privKey, publicKey, senderBlockchainAddress, recipientBlockchainAddress, value}
}

func (t *Transaction) GenerateSignature() *utils.Signature {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	return &utils.Signature{R: r, S: s}
}

func NewWallet() *Wallet {
	w := new(Wallet)
	privKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.PrivateKey = privKey
	w.PublicKey = &privKey.PublicKey

	//address calculation
	h2 := sha256.New()
	h2.Write(w.PublicKey.X.Bytes())
	h2.Write(w.PublicKey.Y.Bytes())
	digest2 := h2.Sum(nil)
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)
	chsum := digest6[:4]
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], chsum[:])
	address := base58.Encode(dc8)

	w.BlockchainAddress = address
	return w
}

func (w *Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PrivateKey        string `json:"private_key"`
		PublicKey         string `json:"public_key"`
		BlockchainAddress string `json:"blockchain_address"`
	}{
		PrivateKey:        w.PrivateKeyStr(),
		PublicKey:         w.PublicKeyStr(),
		BlockchainAddress: w.GetBlockchainAddress(),
	})
}

func (w *Wallet) GetBlockchainAddress() string {
	return w.BlockchainAddress
}

func (w *Wallet) GetPrivateKey() *ecdsa.PrivateKey {
	return w.PrivateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.PrivateKey.D.Bytes())
}

func (w *Wallet) GetPublicKey() *ecdsa.PublicKey {
	return w.PublicKey
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x %x", w.PublicKey.X.Bytes(), w.PublicKey.Y.Bytes())
}

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderPrivateKey == nil ||
		tr.SenderBlockchainAddress == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.SenderPublicKey == nil ||
		tr.Value == nil {
		return false
	}
	return true
}
