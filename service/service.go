package service

import (
	"fmt"
	"math/big"
	"onchain-stats/client"
	"sort"
	"sync"
)

type kv struct {
	Key   string
	Value *big.Int
}

const BASE_URL = "http://localhost:8545"

var evmos_client *client.EvmosClient = &client.EvmosClient{BaseURL: BASE_URL}

func GetLatestBlock() (string, error) {
	return evmos_client.GetBlockNumber()
}

func GetTransactionTrace(txHash string) (map[string]interface{}, error) {
	return evmos_client.GetTransactionTrace(txHash)
}

func IsContractAddress(address string) (bool, error) {
	code, err := evmos_client.GetCode(address, "latest")
	if err != nil {
		return false, err
	}

	return code != "0x", nil
}

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
	blocks, err := evmos_client.GetBlocksInRange(startBlock, endBlock)
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

			balance, err := evmos_client.GetBalance(wallet, blockNumber)
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

func CalculateRichestUsers(startBlock, endBlock int) ([]kv, error) {
	blocks, err := evmos_client.GetBlocksInRange(startBlock, endBlock)
	if err != nil {
		return nil, err
	}

	wallets := ExtractWallets(blocks)
	balances, err := GetWalletBalances(wallets, fmt.Sprintf("0x%x", endBlock))

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
	return evmos_client.GetAccounts()
}

func GetBalance(address, block string) (string, error) {
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
