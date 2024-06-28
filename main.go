package main

import (
	"blockchain/block"
	"blockchain/wallet" // Relative import path to the wallet package
	"fmt"
)

func main() {
	walletM := wallet.NewWallet()
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()

	t := wallet.NewTransaction(walletA.GetPrivateKey(), walletA.GetPublicKey(), walletA.GetBlockchainAddress(), walletB.GetBlockchainAddress(), 3.0)

	blockchain := block.NewBlockChain(walletM.GetBlockchainAddress())
	isAdded := blockchain.AddTransaction(walletA.GetBlockchainAddress(), walletB.GetBlockchainAddress(), 3.0, walletA.GetPublicKey(), t.GenerateSignature())
	fmt.Println("added?", isAdded)
	blockchain.Mining()
	blockchain.PrintBlockchain()
}
