// Package service provides the business logic for the API endpoints.
package service

import (
	"fmt"
	"math/big"
	"sort"
	"sync"
)

type EvmosClientInterface interface {
	GetBlockNumber() (string, error)
	GetTransactionTrace(txHash string) (map[string]interface{}, error)
	GetCode(address, blockNumber string) (string, error)
	GetBlocksInRange(start, end int) ([]map[string]interface{}, error)
	GetBalance(address, block string) (string, error)
	GetAccounts() ([]string, error)
	GetBlock(blockNumber string) (map[string]interface{}, error)
}

type kv struct {
	Key   string
	Value *big.Int
}

var evmosClient EvmosClientInterface

// SetClient Utilized for testing purposes, but can be used to set a custom client
func SetClient(client EvmosClientInterface) {
	evmosClient = client
}

func GetLatestBlock() (string, error) {
	return evmosClient.GetBlockNumber()
}

func GetTransactionTrace(txHash string) (map[string]interface{}, error) {
	return evmosClient.GetTransactionTrace(txHash)
}

// IsContractAddress checks if the given address is a contract address or an EOA.
func IsContractAddress(address string) (bool, error) {
	code, err := evmosClient.GetCode(address, "latest")
	if err != nil {
		return false, err
	}

	return code != "0x", nil
}

// ExtractSmartContracts processes a list of blocks to identify and count interactions with smart contracts.
// It iterates through each block's transactions, checking if the transaction is a contract creation or an interaction with an existing contract.
// It also traces internal contract calls within each transaction.
func ExtractSmartContracts(blocks []map[string]interface{}) (map[string]int, error) {
	contractInteractions := make(map[string]int)
	for _, block := range blocks {
		transactions := block["transactions"].([]interface{})
		for _, tx := range transactions {
			txMap := tx.(map[string]interface{})
			txHash := txMap["hash"].(string)
			to := txMap["to"]

			// it's a contract creation
			if to == nil {
				contractAddress := txMap["contractAddress"]
				if contractAddress != nil && contractAddress != "" {
					contractAddrStr := contractAddress.(string)
					contractInteractions[contractAddrStr]++
				}
			} else {
				toAddress := to.(string)
				isContract, err := IsContractAddress(toAddress)
				if err != nil {
					return nil, err
				}
				if isContract {
					contractInteractions[toAddress]++
				}
			}

			// Add internal contract interactions via transaction trace
			trace, err := GetTransactionTrace(txHash)
			if err != nil {
				return nil, err
			}

			// TODO: assuming trace["calls"] is always present.
			// Not sure since every block has empty transaction list and as such couldn't test this.
			internalCalls := trace["calls"].([]interface{})
			for _, call := range internalCalls {
				callMap := call.(map[string]interface{})
				contractAddress := callMap["to"].(string)
				contractInteractions[contractAddress]++
			}
		}
	}

	return contractInteractions, nil
}

// ExtractWallets processes a list of blocks to identify unique wallets that have interacted with the blockchain.
// It iterates through each block's transactions, checking the sender and receiver of each transaction.
func ExtractWallets(blocks []map[string]interface{}) []string {
	wallets := make(map[string]struct{})
	for _, block := range blocks {
		transactions := block["transactions"].([]interface{})
		for _, tx := range transactions {
			txMap := tx.(map[string]interface{})

			from := txMap["from"].(string)
			if from != "" {
				wallets[from] = struct{}{}
			}

			to := txMap["to"]
			if to != nil && to.(string) != "" {
				toAddress := to.(string)

				isContract, err := IsContractAddress(toAddress)
				if err != nil {
					continue
				}

				// it's an EOA (not a contract)
				if !isContract {
					wallets[toAddress] = struct{}{}
				}
			}
		}
	}

	walletList := make([]string, 0, len(wallets))
	for wallet := range wallets {
		walletList = append(walletList, wallet)
	}
	return walletList
}

func GetSmartContracts(startBlock, endBlock int) ([]kv, error) {
	blocks, err := evmosClient.GetBlocksInRange(startBlock, endBlock)
	if err != nil {
		return nil, err
	}

	contractInteractions, err := ExtractSmartContracts(blocks)
	if err != nil {
		return nil, err
	}

	// Sort contracts by number of interactions
	var sortedContracts []kv
	for k, v := range contractInteractions {
		sortedContracts = append(sortedContracts, kv{k, big.NewInt(int64(v))})
	}

	sort.Slice(sortedContracts, func(i, j int) bool {
		return sortedContracts[i].Value.Cmp(sortedContracts[j].Value) > 0
	})

	return sortedContracts, nil
}

func GetWalletBalances(wallets []string, blockNumber string) (map[string]*big.Int, error) {
	balances := make(map[string]*big.Int)

	var wg sync.WaitGroup
	balanceChannel := make(chan kv, len(wallets))
	workerPool := make(chan struct{}, 8) // Limit to 8 concurrent goroutines

	// Parallelize balance fetching with worker pool
	for _, wallet := range wallets {
		wg.Add(1)
		workerPool <- struct{}{}
		go func(wallet string) {
			defer wg.Done()
			defer func() { <-workerPool }()

			balance, err := evmosClient.GetBalance(wallet, blockNumber)
			if err == nil {
				balanceInt := new(big.Int)
				balanceInt.SetString(balance[2:], 16) // Convert hex string to big.Int
				balanceChannel <- kv{wallet, balanceInt}
			}
		}(wallet)
	}

	// Close the channel once all balances are fetched
	go func() {
		wg.Wait()
		close(balanceChannel)
	}()

	// Collect results from channel
	for walletBalance := range balanceChannel {
		balances[walletBalance.Key] = walletBalance.Value
	}

	return balances, nil
}

// CalculateRichestUsers calculates the richest users based on their wallet balances at the end block.
// It only needs the last block, since the last block contains the most up-to-date balances of all wallets.
func CalculateRichestUsers(block int) ([]kv, error) {
	blocks, err := evmosClient.GetBlocksInRange(block, block)
	if err != nil {
		return nil, err
	}

	wallets := ExtractWallets(blocks)
	balances, err := GetWalletBalances(wallets, fmt.Sprintf("0x%x", block))

	if err != nil {
		return nil, err
	}

	// Sort wallets by balance
	var sortedWallets []kv
	for k, v := range balances {
		sortedWallets = append(sortedWallets, kv{k, v})
	}

	sort.Slice(sortedWallets, func(i, j int) bool {
		return sortedWallets[i].Value.Cmp(sortedWallets[j].Value) > 0
	})

	return sortedWallets, nil
}

func GetAccounts() ([]string, error) {
	return evmosClient.GetAccounts()
}

func GetBalance(address, block string) (string, error) {
	balance, err := evmosClient.GetBalance(address, block)
	if err != nil {
		return "", err
	}

	return balance, nil
}

func GetBlock(blockNumber string) (map[string]interface{}, error) {
	block, err := evmosClient.GetBlock(blockNumber)
	if err != nil {
		return nil, err
	}

	return block, nil
}
