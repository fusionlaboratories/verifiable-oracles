package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/qredo/verifiable-oracles/pkg/balance"
)

type Account struct {
	client     *ethclient.Client
	account    common.Address
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

// Create a new Account (or Wallet)
func NewAccount(client *ethclient.Client, address string, hexkey string) (*Account, error) {
	// TODO: This check does not make sense, as the invariant can still be
	// violated _if_ the structure gets initialized with _new_ keyword...
	// TODO: Refactor the code, to reduce the amount of error checks
	// - https://go.dev/blog/errors-are-values
	// TODO: Consider using errors.Join
	// - https://pkg.go.dev/errors#Join
	if client == nil {
		return nil, errors.New("client is nil")
	}

	// Trim "0x" prefix, because crypto.HexToECDSA doesn't seem to handle it
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(hexkey, "0x"))
	if err != nil {
		return nil, err
	}

	// generating public key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("cannot assert type: publicKey is not of *ecdsa.PublicKey")
	}

	acc := &Account{
		client:     client,
		account:    common.HexToAddress(address),
		privateKey: privateKey,
		publicKey:  publicKeyECDSA,
	}

	return acc, nil
}

func (a *Account) Account() common.Address {
	return a.account
}

func (a *Account) BalanceAt(ctx context.Context, blockNumber *big.Int) (*balance.Balance, error) {
	if blockNumber == nil {
		return nil, fmt.Errorf("blockNumber is nil")
	}

	wei, err := a.client.BalanceAt(ctx, a.account, blockNumber)
	if err != nil {
		return nil, err
	}

	return balance.BalanceFromWei(wei), nil
}

func (a *Account) PublicKey() *ecdsa.PublicKey {
	return a.publicKey
}

func BlockNumberFromUint64(blockNumber uint64) *big.Int {
	return new(big.Int).SetUint64(blockNumber)
}

func (a *Account) SendEthTo(ctx context.Context, toAddress common.Address, amount *big.Int, gasLimit uint64) (*types.Transaction, error) {
	// Create transaction
	fromAddress := a.account
	nonce, err := a.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, err
	}

	gasPrice, err := a.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	var data []byte
	tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, data)

	chainID, err := a.client.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), a.privateKey)
	if err != nil {
		return nil, err
	}

	return signedTx, a.client.SendTransaction(ctx, signedTx)
}

func main() {
	host := "127.0.0.1"
	port := "8545"
	url := fmt.Sprintf("http://%s:%s", host, port)

	client, err := ethclient.Dial(url)
	if err != nil {
		panic(err)
	}
	defer client.Close() // Close the client when done

	fmt.Printf("Successfully connected to ganache on %s\n", url)

	ctx := context.Background()

	// Get current block number
	blockNumberUint64, err := client.BlockNumber(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current block number is %d\n", blockNumberUint64)

	// Create accounts
	accounts := []struct {
		address string
		key     string
	}{
		{"0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1", "0x4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d"},
		{"0xFFcf8FDEE72ac11b5c542428B35EEF5769C409f0", "0x6cbed15c793ce57650b9877cf6fa156fbef513c4e6134f022a85b1ffdd59b2a1"},
	}

	blockNumber := BlockNumberFromUint64(blockNumberUint64)

	for _, accInfo := range accounts {
		acc, err := NewAccount(client, accInfo.address, accInfo.key)
		if err != nil {
			panic(err)
		}

		balance, err := acc.BalanceAt(ctx, blockNumber)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Account %s has %s ETH (Wei: %s)\n", acc.Account().Hex(), balance.Eth().String(), balance.Wei().String())

		bobAddress := common.HexToAddress("0xFFcf8FDEE72ac11b5c542428B35EEF5769C409f0")
		amount := big.NewInt(1000000000000000000)
		gasLimit := uint64(21000)

		signedTx, err := acc.SendEthTo(ctx, bobAddress, amount, gasLimit)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Signed Tx: %s\n", signedTx.Hash())
	}

	nextBlock := new(big.Int).Add(blockNumber, big.NewInt(1))
	fmt.Printf("Next block number is %s\n", nextBlock)

	for _, accInfo := range accounts {
		acc, err := NewAccount(client, accInfo.address, accInfo.key)
		if err != nil {
			panic(err)
		}

		balance, err := acc.BalanceAt(ctx, nextBlock)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Account %s has %s ETH (Wei: %s)\n", acc.Account().Hex(), balance.Eth().String(), balance.Wei().String())
	}
}
