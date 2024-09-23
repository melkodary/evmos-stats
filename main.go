package main

import (
	"fmt"
	"net/http"
	"onchain-stats/api"
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

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
