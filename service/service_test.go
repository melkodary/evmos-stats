package service

import "testing"

type MockEvmosClient struct{}

func (*MockEvmosClient) GetAccounts() ([]string, error) {
	return []string{}, nil
}

// TODO - create a random block
func (*MockEvmosClient) GetBlock(blockNumber string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (m *MockEvmosClient) GetBlockNumber() (string, error) {
	return "0x1", nil
}

// TODO: create a random transaction
func (m *MockEvmosClient) GetTransactionTrace(txHash string) (map[string]interface{}, error) {
	return map[string]interface{}{"calls": []interface{}{}}, nil
}

func (m *MockEvmosClient) GetCode(address, blockNumber string) (string, error) {
	return "0x", nil
}

func (m *MockEvmosClient) GetBlocksInRange(start, end int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

func (m *MockEvmosClient) GetBalance(address, block string) (string, error) {
	return "0x1", nil
}

func TestGetLatestBlock(t *testing.T) {
	SetClient(&MockEvmosClient{})
	block, err := GetLatestBlock()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if block != "0x1" {
		t.Errorf("Expected block number 0x1, got %v", block)
	}
}
