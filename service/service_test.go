package service

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockEvmosClient struct {
	accounts         []string
	block            map[string]interface{}
	blockNumber      string
	transactionTrace map[string]interface{}
	code             map[string]string
	blocksInRange    []map[string]interface{}
	balances         map[string]string
}

func (m *MockEvmosClient) GetAccounts() ([]string, error) {
	return m.accounts, nil
}

func (m *MockEvmosClient) GetBlock(blockNumber string) (map[string]interface{}, error) {
	return m.block, nil
}

func (m *MockEvmosClient) GetBlockNumber() (string, error) {
	return m.blockNumber, nil
}

func (m *MockEvmosClient) GetTransactionTrace(txHash string) (map[string]interface{}, error) {
	return m.transactionTrace, nil
}

func (m *MockEvmosClient) GetCode(address, blockNumber string) (string, error) {
	if code, exists := m.code[address]; exists {
		return code, nil
	}
	return "0x", nil
}

func (m *MockEvmosClient) GetBlocksInRange(startBlock, endBlock int) ([]map[string]interface{}, error) {
	return m.blocksInRange, nil
}

func (m *MockEvmosClient) GetBalance(address, block string) (string, error) {
	if balance, exists := m.balances[address]; exists {
		return balance, nil
	}
	return "0x0", nil
}

func TestGetLatestBlock(t *testing.T) {
	client := &MockEvmosClient{
		blockNumber: "0x1",
	}
	SetClient(client)
	block, err := GetLatestBlock()
	assert.NoError(t, err)
	assert.Equal(t, "0x1", block)
}

func TestGetSmartContracts(t *testing.T) {
	client := &MockEvmosClient{
		blocksInRange: []map[string]interface{}{
			{
				"transactions": []interface{}{
					map[string]interface{}{
						"hash":            "0xTxHash1",
						"to":              "0xContractAddress1",
						"contractAddress": nil,
					},
					map[string]interface{}{
						"hash":            "0xTxHash2",
						"to":              nil,
						"contractAddress": "0xContractAddress2",
					},
					map[string]interface{}{
						"hash":            "0xTxHash3",
						"to":              "0xContractAddress3",
						"contractAddress": nil,
					},
					map[string]interface{}{
						"hash":            "0xTxHash4",
						"to":              "0xContractAddress4",
						"contractAddress": nil,
					},
					map[string]interface{}{
						"hash":            "0xTxHash5",
						"to":              "0xContractAddress1",
						"contractAddress": nil,
					},
					map[string]interface{}{
						"hash":            "0xTxHash6",
						"to":              "0xContractAddress3",
						"contractAddress": nil,
					},
					map[string]interface{}{
						"hash":            "0xTxHash7",
						"to":              "0xContractAddress4",
						"contractAddress": nil,
					},
					map[string]interface{}{
						"hash":            "0xTxHash8",
						"to":              "0xContractAddress4",
						"contractAddress": nil,
					},
				},
			},
		},
		transactionTrace: map[string]interface{}{
			"calls": []interface{}{
				map[string]interface{}{
					"to": "0xContractAddress2",
				},
			},
		},
		code: map[string]string{
			"0xContractAddress1": "0x6001600101",
			"0xContractAddress2": "0x6001600102",
			"0xContractAddress3": "0x6001600103",
			"0xContractAddress4": "0x6001600104",
		},
	}
	SetClient(client)

	expectedContracts := []kv{
		{"0xContractAddress2", big.NewInt(9)},
		{"0xContractAddress4", big.NewInt(3)},
		{"0xContractAddress3", big.NewInt(2)},
		{"0xContractAddress1", big.NewInt(2)},
	}

	contracts, err := GetSmartContracts(100, 200)

	assert.NoError(t, err)
	assert.Equal(t, len(expectedContracts), len(contracts))

	for i, contract := range contracts {
		assert.Equal(t, expectedContracts[i].Key, contract.Key)
		if expectedContracts[i].Value.Cmp(contract.Value) != 0 {
			t.Errorf("Value mismatch for contract %s: expected %s, got %s", contract.Key, expectedContracts[i].Value.String(), contract.Value.String())
		}
	}

	for i, expectedContract := range expectedContracts {
		assert.Equal(t, expectedContract.Key, contracts[i].Key)
	}
}

func TestCalculateRichestUsers(t *testing.T) {
	client := &MockEvmosClient{
		blocksInRange: []map[string]interface{}{
			{
				"transactions": []interface{}{
					map[string]interface{}{
						"hash": "0xTxHash1",
						"from": "0xWallet1",
						"to":   "0xWallet2",
					},
					map[string]interface{}{
						"hash": "0xTxHash2",
						"from": "0xWallet3",
						"to":   "0xWallet4",
					},
				},
			},
		},
		code: map[string]string{
			"0xWallet1": "0x",
			"0xWallet2": "0x",
			"0xWallet3": "0x",
			"0xWallet4": "0x",
		},
		balances: map[string]string{
			"0xWallet1": "0x5",
			"0xWallet2": "0x3",
			"0xWallet3": "0x8",
			"0xWallet4": "0x1",
		},
	}

	SetClient(client)

	expectedWallets := []kv{
		{"0xWallet3", big.NewInt(8)},
		{"0xWallet1", big.NewInt(5)},
		{"0xWallet2", big.NewInt(3)},
		{"0xWallet4", big.NewInt(1)},
	}

	wallets, err := CalculateRichestUsers(200)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedWallets), len(wallets))

	expectedMap := make(map[string]*big.Int)
	for _, wallet := range expectedWallets {
		expectedMap[wallet.Key] = wallet.Value
	}

	for _, wallet := range wallets {
		expectedValue, exists := expectedMap[wallet.Key]
		assert.True(t, exists, "Unexpected wallet: %s", wallet.Key)
		assert.Equal(t, 0, expectedValue.Cmp(wallet.Value), "Value mismatch for wallet %s: expected %s, got %s", wallet.Key, expectedValue.String(), wallet.Value.String())
	}
}
