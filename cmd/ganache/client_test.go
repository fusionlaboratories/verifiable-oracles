package main

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestNewAccount(t *testing.T) {
	host := "127.0.0.1"
	port := "8545"
	url := fmt.Sprintf("http://%s:%s", host, port)

	client, err := ethclient.Dial(url)
	if err != nil {
		t.Fatalf("Failed to connect to ganache on %s: %v", url, err)
	}
	defer client.Close()

	tests := []struct {
		name      string
		client    *ethclient.Client
		address   string
		hexkey    string
		expectNil bool // Whether the fields in the expected result should be nil
	}{
		{
			name:      "Valid Account",
			client:    client,
			address:   "0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1",
			hexkey:    "0x4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d",
			expectNil: false,
		},
		{
			name:      "Invalid Account",
			client:    client,
			address:   "invalid address",
			hexkey:    "invalid hexkey",
			expectNil: true,
		},
		{
			name:      "Invalid Client",
			client:    nil,
			address:   "0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1",
			hexkey:    "0x4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d",
			expectNil: true,
		},
		// Add more test cases here
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, _ := NewAccount(test.client, test.address, test.hexkey)

			if test.expectNil {
				if result != nil {
					t.Errorf("Expected nil, but got %+v", result)
				}
			} else {
				if result == nil {
					t.Error("Expected non-nil result, but got nil")
				} else {
					if result.client == nil || result.account == (common.Address{}) || result.privateKey == nil || result.publicKey == nil {
						t.Errorf("Fields in result are nil: %+v", result)
					}
				}
			}
		})
	}
}
