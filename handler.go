package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func getSmartContractsHandler(w http.ResponseWriter, r *http.Request) {
	contractInteractions, err := GetSmartContracts(100, 200)

	if err != nil {
		http.Error(w, "Error fetching smart contracts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contractInteractions)
}

func getRichestUsersHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: get blocks range from query params
	sortedWallets, err := CalculateRichestUsers(100, 200)

	if err != nil {
		http.Error(w, "Error fetching richest users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sortedWallets)
}

func getAccountsHandler(w http.ResponseWriter, r *http.Request) {
	accounts, err := GetAccounts()
	if err != nil {
		http.Error(w, "Error fetching accounts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

func getBalanceHandler(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	block := r.URL.Query().Get("block")
	if block == "" {
		block = "latest"
	}

	if address == "" {
		http.Error(w, "Missing address", http.StatusBadRequest)
		return
	}

	balance, err := GetBalance(address, block)
	if err != nil {
		http.Error(w, "Error fetching balance: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
}

func getBlockHandler(w http.ResponseWriter, r *http.Request) {
	blockNumber := r.URL.Query().Get("blockNumber")
	if blockNumber == "" {
		http.Error(w, "Missing blockNumber", http.StatusBadRequest)
		return
	}

	block, err := GetBlock(blockNumber)
	if err != nil {
		http.Error(w, "Error fetching block: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(block)
}

func getBlockNumberHandler(w http.ResponseWriter, r *http.Request) {
	blockNumber, err := GetLatestBlock()
	if err != nil {
		http.Error(w, "Error fetching block number: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blockNumber)
}

func getTransactionTraceHandler(w http.ResponseWriter, r *http.Request) {
	txHash := r.URL.Query().Get("txHash")
	if txHash == "" {
		http.Error(w, "Missing txHash", http.StatusBadRequest)
		return
	}

	trace, err := GetTransactionTrace(txHash)
	if err != nil {
		http.Error(w, "Error fetching transaction trace: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trace)
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
