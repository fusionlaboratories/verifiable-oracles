//go:build integration

// TODO: Fix the workspace, so this file is properly included in the workspace
// TODO: Reorganise the source code
package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
)

func TestGanacheLive(t *testing.T) {
	host := "127.0.0.1"
	port := "8545"
	url := fmt.Sprintf("http://%s:%s", host, port)

	t.Logf("Trying to connect to %s", url)
	client, err := ethclient.Dial(url)
	if err != nil {
		t.Error(err)
	}

	// What is the current block
	blockNumberUint64, err := client.BlockNumber(context.Background())
	if err != nil {
		t.Error(err)
	}
	blockNumber := BlockNumberFromUint64(blockNumberUint64)
	t.Logf("Current block number is %d\n", blockNumber)
}
