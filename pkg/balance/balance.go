package balance

import "math/big"

type Balance struct {
	// NOTE: This potentially can be replaces by a value field
	wei *big.Int
	eth *big.Float
}

func BalanceFromWei(wei *big.Int) *Balance {
	fresh := &big.Int{}
	fresh.Set(wei)

	return &Balance{wei: fresh}
}

// Get wei balance and initalize it if needed
func (b *Balance) getWei() *big.Int {
	if b.wei == nil {
		b.wei = big.NewInt(0)
	}

	return b.wei
}

// Get eth balance and initialize it if needed
func (b *Balance) getEth() *big.Float {
	if b.eth == nil {
		b.eth = ethFromWei(b.getWei())
	}

	return b.eth
}

// Return a copy of WEI balance
func (b *Balance) Wei() *big.Int {
	fresh := &big.Int{}
	fresh.Set(b.getWei())
	return fresh
}

// Return a copy of ETH balance
func (b *Balance) Eth() *big.Float {
	fresh := big.NewFloat(0)
	fresh.Set(b.getEth())
	return fresh
}

// Convert wei to eth
func ethFromWei(wei *big.Int) *big.Float {
	result := big.NewFloat(0)
	result.SetInt(wei)
	result.Quo(result, big.NewFloat(1e18))
	return result
}
