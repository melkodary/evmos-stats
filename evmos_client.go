package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type EmvosClient struct {
	BaseURL string
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

func (c *EmvosClient) GetBalance(address string, blockNumber string) (string, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"method":  "eth_getBalance",
		"params":  []interface{}{address, blockNumber},
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

func (c *EmvosClient) GetBlocksInRange(start, end int) ([]map[string]interface{}, error) {
	var blocks []map[string]interface{}
	for i := start; i <= end; i++ {
		blockNumber := fmt.Sprintf("0x%x", i)
		block, err := c.GetBlock(blockNumber)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}
