package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", health)

	http.HandleFunc("/accounts", getAccountsHandler)
	http.HandleFunc("/balance", getBalanceHandler)
	http.HandleFunc("/blocknumber", getBlockNumberHandler)
	http.HandleFunc("/block", getBlockHandler)
	http.HandleFunc("/transactiontrace", getTransactionTraceHandler)

	http.HandleFunc("/smartcontracts", getSmartContractsHandler)
	http.HandleFunc("/richestusers", getRichestUsersHandler)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
