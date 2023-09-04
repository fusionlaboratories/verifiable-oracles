package balance_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/qredo/verifiable-oracles/pkg/balance"
)

func TestZeroValueWei(t *testing.T) {
	b := &balance.Balance{}

	assert.Equal(t, big.NewInt(0), b.Wei())
}

func TestZeroValueEth(t *testing.T) {
	b := &balance.Balance{}

	assert.Equal(t, big.NewFloat(0), b.Eth())
}

func TestBalanceFromWeiMakesCopy(t *testing.T) {
	a := big.NewInt(1)
	b := balance.BalanceFromWei(a)

	assert.Equal(t, a, b.Wei())
	a.Add(a, a)
	assert.NotEqual(t, a, b.Wei())
}

func TestWeiReturnsFreshCopy(t *testing.T) {
	b := balance.BalanceFromWei(big.NewInt(1))
	c := b.Wei()

	assert.Equal(t, b.Wei(), c)
	c.Add(c, c)
	assert.NotEqual(t, b.Wei(), c)
}

func TestEthReturnsFreshCopy(t *testing.T) {
	b := balance.BalanceFromWei(big.NewInt(1))
	c := b.Eth()

	assert.Equal(t, b.Eth(), c)
	c.Add(c, c)
	assert.NotEqual(t, b.Eth(), c)
}

func TestBalance(t *testing.T) {
	tests := map[string]struct {
		wei  *big.Int
		want string
	}{
		"0_WEI":             {wei: big.NewInt(0), want: "0"},
		"1_WEI":             {wei: big.NewInt(1), want: "1e-18"},
		"1_000_WEI":         {wei: big.NewInt(1_000), want: "1e-15"},
		"1_000_000_WEI":     {wei: big.NewInt(1_000_000), want: "1e-12"},
		"1_000_000_000_WEI": {wei: big.NewInt(1_000_000_000), want: "1e-09"},
		"1_ETH":             {wei: big.NewInt(1_000_000_000_000_000_000), want: "1"},
		"5e-06_ETH":         {wei: big.NewInt(5_000_000_000_000), want: "5e-06"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			b := balance.BalanceFromWei(tc.wei)

			assert.Equal(t, tc.wei, b.Wei())
			assert.Equal(t, tc.want, b.Eth().String())
		})
	}
}
