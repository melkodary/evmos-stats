package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const BASE_URL = "http://localhost:8545"

type EmvosClient struct {
	BaseURL string
}

func NewEmvosClient(baseURL string) *EmvosClient {
	return &EmvosClient{BaseURL: baseURL}
}

func (c *EmvosClient) GetBlockNumber() (string, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
		"id":      1,
		"jsonrpc": "2.0",
	})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(c.BaseURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Result string `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Result, nil
}

func (c *EmvosClient) GetBlock(blockNumber string) (map[string]interface{}, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{blockNumber, true},
		"id":      1,
		"jsonrpc": "2.0",
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(c.BaseURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Result map[string]interface{} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	blockDetails := map[string]interface{}{
		"hash":             result.Result["hash"],
		"transactions":     result.Result["transactions"],
		"gasUsed":          result.Result["gasUsed"],
		"gasLimit":         result.Result["gasLimit"],
		"timestamp":        result.Result["timestamp"],
		"totalDifficulty":  result.Result["totalDifficulty"],
		"transactionsRoot": result.Result["transactionsRoot"],
	}

	return blockDetails, nil
}

func (c *EmvosClient) GetBalance(address string, blockParameter string) (string, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"method":  "eth_getBalance",
		"params":  []interface{}{address, blockParameter},
		"id":      1,
		"jsonrpc": "2.0",
	})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(c.BaseURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Result string `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Result, nil
}

func (c *EmvosClient) GetAccounts() ([]string, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"method":  "eth_accounts",
		"params":  []interface{}{},
		"id":      1,
		"jsonrpc": "2.0",
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(c.BaseURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Result []string `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Result, nil
}

func (c *EmvosClient) GetTransactionTrace(txHash string) (map[string]interface{}, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"method":  "debug_traceTransaction",
		"params":  []interface{}{txHash, map[string]string{"tracer": "callTracer"}},
		"id":      1,
		"jsonrpc": "2.0",
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(c.BaseURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Result map[string]interface{} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Result, nil
}

func (c *EmvosClient) GetTransactionReceipt(txHash string) (map[string]interface{}, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"method":  "eth_getTransactionReceipt",
		"params":  []interface{}{txHash},
		"id":      1,
		"jsonrpc": "2.0",
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(c.BaseURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Result map[string]interface{} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Result, nil
}

func main() {
	client := NewEmvosClient(BASE_URL)

	http.HandleFunc("/blocknumber", func(w http.ResponseWriter, r *http.Request) {
		blockNumber, err := client.GetBlockNumber()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Latest Block Number: %s", blockNumber)
	})

	http.HandleFunc("/block", func(w http.ResponseWriter, r *http.Request) {
		blockNumber := r.URL.Query().Get("blockNumber")
		if blockNumber == "" {
			http.Error(w, "blockNumber query parameter is required", http.StatusBadRequest)
			return
		}

		blockDetails, err := client.GetBlock(blockNumber)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(blockDetails)
	})

	http.HandleFunc("/balance", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "address query parameter is required", http.StatusBadRequest)
			return
		}

		blockParameter := r.URL.Query().Get("blockParameter")
		if blockParameter == "" {
			blockParameter = "latest"
		}

		balance, err := client.GetBalance(address, blockParameter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Balance: %s", balance)
	})

	http.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		accounts, err := client.GetAccounts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(accounts)
	})

	http.HandleFunc("/transactiontrace", func(w http.ResponseWriter, r *http.Request) {
		txHash := r.URL.Query().Get("txHash")
		if txHash == "" {
			http.Error(w, "txHash query parameter is required", http.StatusBadRequest)
			return
		}

		trace, err := client.GetTransactionTrace(txHash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(trace)
	})

	http.HandleFunc("/transactionreceipt", func(w http.ResponseWriter, r *http.Request) {
		txHash := r.URL.Query().Get("txHash")
		if txHash == "" {
			http.Error(w, "txHash query parameter is required", http.StatusBadRequest)
			return
		}

		receipt, err := client.GetTransactionReceipt(txHash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(receipt)
	})

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
