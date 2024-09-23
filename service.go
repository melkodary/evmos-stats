package main

import (
	"fmt"
	"math/big"
	"sort"
)

const BASE_URL = "http://localhost:8545"

var evmos_client *EmvosClient = &EmvosClient{BaseURL: BASE_URL}

func GetLatestBlock() (string, error) {
	return evmos_client.GetBlockNumber()
}

func GetTransactionTrace(txHash string) (map[string]interface{}, error) {
	return evmos_client.GetTransactionTrace(txHash)
}

func ExtractSmartContracts(blocks []map[string]interface{}) map[string]int {
	contractInteractions := make(map[string]int)
	for _, block := range blocks {
		transactions := block["transactions"].([]interface{})
		for _, tx := range transactions {
			txMap := tx.(map[string]interface{})
			to := txMap["to"].(string)
			if to != "" {
				contractInteractions[to]++
			}
		}
	}
	return contractInteractions
}

func ExtractWallets(blocks []map[string]interface{}) []string {
	wallets := make(map[string]struct{})
	for _, block := range blocks {
		transactions := block["transactions"].([]interface{})
		for _, tx := range transactions {
			txMap := tx.(map[string]interface{})
			to := txMap["to"].(string)
			if to != "" {
				wallets[to] = struct{}{}
			}
		}
	}

	walletList := make([]string, 0, len(wallets))
	for wallet := range wallets {
		walletList = append(walletList, wallet)
	}
	return walletList
}

func GetSmartContracts(startBlock, endBlock int) (map[string]int, error) {
	blocks, err := evmos_client.GetBlocksInRange(startBlock, endBlock)
	if err != nil {
		return nil, err
	}

	contractInteractions := ExtractSmartContracts(blocks)

	return contractInteractions, nil
}

func GetWalletBalances(wallets []string, blockNumber string) (map[string]*big.Int, error) {
	balances := make(map[string]*big.Int)
	for _, wallet := range wallets {
		balance, err := evmos_client.GetBalance(wallet, blockNumber)
		if err != nil {
			return nil, err
		}
		balanceInt := new(big.Int)
		balanceInt.SetString(balance[2:], 16) // Convert hex string to big.Int
		balances[wallet] = balanceInt
	}
	return balances, nil
}

func CalculateRichestUsers(startBlock, endBlock int) (map[string]*big.Int, error) {
	blocks, err := evmos_client.GetBlocksInRange(startBlock, endBlock)
	if err != nil {
		return nil, err
	}

	wallets := ExtractWallets(blocks)
	balances, err := GetWalletBalances(wallets, fmt.Sprintf("0x%x", endBlock))

	if err != nil {
		return nil, err
	}

	type kv struct {
		Key   string
		Value *big.Int
	}

	var sortedWallets []kv
	for k, v := range balances {
		sortedWallets = append(sortedWallets, kv{k, v})
	}

	sort.Slice(sortedWallets, func(i, j int) bool {
		return sortedWallets[i].Value.Cmp(sortedWallets[j].Value) > 0
	})

	return balances, nil
}

func GetAccounts() ([]string, error) {
	return evmos_client.GetAccounts()
}

func GetBalance(address, block string) (string, error) {
	if block == "" {
		block = "latest"
	}

	balance, err := evmos_client.GetBalance(address, block)
	if err != nil {
		return "", err
	}

	return balance, nil
}

func GetBlock(blockNumber string) (map[string]interface{}, error) {
	block, err := evmos_client.GetBlock(blockNumber)
	if err != nil {
		return nil, err
	}

	return block, nil
}
