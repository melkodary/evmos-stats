// Package main is used to start the server
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"onchain-stats/client"
	"onchain-stats/service"
	"time"
)

const BaseURL = "http://localhost:8545"

func GetSmartContractsHandler(w http.ResponseWriter, r *http.Request) {
	contractInteractions, err := service.GetSmartContracts(100, 200)

	if err != nil {
		http.Error(w, "Error fetching smart contracts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(contractInteractions); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

func GetRichestUsersHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: get blocks range from query params
	richestUsers, err := service.CalculateRichestUsers(200)

	if err != nil {
		http.Error(w, "Error fetching richest users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(richestUsers); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

func GetAccountsHandler(w http.ResponseWriter, r *http.Request) {
	accounts, err := service.GetAccounts()
	if err != nil {
		http.Error(w, "Error fetching accounts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(accounts); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

func GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	block := r.URL.Query().Get("block")
	if block == "" {
		block = "latest"
	}

	if address == "" {
		http.Error(w, "Missing address", http.StatusBadRequest)
		return
	}

	balance, err := service.GetBalance(address, block)
	if err != nil {
		http.Error(w, "Error fetching balance: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(balance); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

func GetBlockHandler(w http.ResponseWriter, r *http.Request) {
	blockNumber := r.URL.Query().Get("blockNumber")
	if blockNumber == "" {
		http.Error(w, "Missing blockNumber", http.StatusBadRequest)
		return
	}

	block, err := service.GetBlock(blockNumber)
	if err != nil {
		http.Error(w, "Error fetching block: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(block); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

func GetBlockNumberHandler(w http.ResponseWriter, r *http.Request) {
	blockNumber, err := service.GetLatestBlock()
	if err != nil {
		http.Error(w, "Error fetching block number: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(blockNumber); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

func GetTransactionTraceHandler(w http.ResponseWriter, r *http.Request) {
	txHash := r.URL.Query().Get("txHash")
	if txHash == "" {
		http.Error(w, "Missing txHash", http.StatusBadRequest)
		return
	}

	trace, err := service.GetTransactionTrace(txHash)
	if err != nil {
		http.Error(w, "Error fetching transaction trace: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(trace); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

func Health(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprintf(w, "Hello, World!"); err != nil {
		http.Error(w, "Error writing response: "+err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	service.SetClient(&client.EvmosClient{BaseURL: BaseURL})

	http.HandleFunc("/", Health)

	http.HandleFunc("/accounts", GetAccountsHandler)
	http.HandleFunc("/balance", GetBalanceHandler)
	http.HandleFunc("/blocknumber", GetBlockNumberHandler)
	http.HandleFunc("/block", GetBlockHandler)
	http.HandleFunc("/transactiontrace", GetTransactionTraceHandler)

	http.HandleFunc("/smartcontracts", GetSmartContractsHandler)
	http.HandleFunc("/richestusers", GetRichestUsersHandler)

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
		Handler:      nil, // Use the default http.DefaultServeMux
	}

	fmt.Println("Server is running on port 8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
