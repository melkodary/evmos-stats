// Package main is used to start the server
package main

import (
	"fmt"
	"net/http"
	"onchain-stats/api"
	"time"
)

func main() {
	http.HandleFunc("/", api.Health)

	http.HandleFunc("/accounts", api.GetAccountsHandler)
	http.HandleFunc("/balance", api.GetBalanceHandler)
	http.HandleFunc("/blocknumber", api.GetBlockNumberHandler)
	http.HandleFunc("/block", api.GetBlockHandler)
	http.HandleFunc("/transactiontrace", api.GetTransactionTraceHandler)

	http.HandleFunc("/smartcontracts", api.GetSmartContractsHandler)
	http.HandleFunc("/richestusers", api.GetRichestUsersHandler)

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
