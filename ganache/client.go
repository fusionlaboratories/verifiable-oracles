package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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

func (a Account) Account() common.Address {
	return a.account
}

func (a Account) WeiAt(ctx context.Context, blockNumber *big.Int) (*big.Int, error) {
	if blockNumber == nil {
		return nil, errors.New("blockNumber is nil")
	}

	return a.client.BalanceAt(ctx, a.account, blockNumber)
}

func (a Account) EthAt(ctx context.Context, blockNumber *big.Int) (*big.Float, error) {
	wei, err := a.WeiAt(ctx, blockNumber)

	if err != nil {
		return nil, err
	}

	return EthFromWei(wei), nil
}

func (a Account) PublicKey() *ecdsa.PublicKey {
	return a.publicKey
}

func EthFromWei(wei *big.Int) *big.Float {
	weiFloat := new(big.Float)
	weiFloat.SetString(wei.String())
	eth := new(big.Float).Quo(weiFloat, big.NewFloat(math.Pow10(18)))
	return eth
}

func BlockNumberFromUint64(blockNumber uint64) *big.Int {
	return new(big.Int).SetUint64(blockNumber)
}

func (a Account) SendEthTo(ctx context.Context, toAddress common.Address, amount *big.Int, gasLimit uint64) (*types.Transaction, error) {
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

	// Sign transaction
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
	fmt.Printf("Successfully connected to ganache on %s\n", url)

	// What is the current block
	blockNumberUint64, err := client.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	blockNumber := BlockNumberFromUint64(blockNumberUint64)
	fmt.Printf("Current block number is %d\n", blockNumber)

	alice, err := NewAccount(client, "0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1", "0x4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d")
	if err != nil {
		panic(err)
	}

	aliceEth, err := alice.EthAt(context.Background(), blockNumber)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Account %s has %s ETH\n", alice.Account(), aliceEth.String())

	// getting the 2nd account
	bob, err := NewAccount(client, "0xFFcf8FDEE72ac11b5c542428B35EEF5769C409f0", "0x6cbed15c793ce57650b9877cf6fa156fbef513c4e6134f022a85b1ffdd59b2a1")
	if err != nil {
		panic(err)
	}

	bobEth, err := bob.EthAt(context.Background(), blockNumber)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Account %s has %s ETH\n", bob.Account(), bobEth.String())

	// Try to send ETH from alice to bob
	signedTx, err := alice.SendEthTo(context.Background(), bob.Account(), big.NewInt(1000000000000000000), 21000)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Signed Tx: %s\n", signedTx.Hash())

	nextBlock := big.NewInt(0).Add(blockNumber, big.NewInt(1))
	fmt.Printf("Next block number is %d\n", nextBlock)

	aliceEth_, err := alice.EthAt(context.Background(), nextBlock)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Account %s has %s ETH\n", alice.Account(), aliceEth_.String())

	bobEth_, err := bob.EthAt(context.Background(), nextBlock)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Account %s has %s ETH\n", bob.Account(), bobEth_.String())
}
