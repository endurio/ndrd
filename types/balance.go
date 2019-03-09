// Copyright (c) 2018-2019 The Endurio developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package types

import "math/big"

// Balance contains all Balances carried by an tx output.
// Acts like a map of token => amount, but practically almost always a single to 2 tokens array.
type Balance struct {
	a0, a1 Amount
}

var (
	BalanceEmpty = Balance{}
)

func NewBalance(a0, a1 Amount) *Balance {
	return &Balance{a0, a1}
}

func (b Balance) Map() []Amount {
	return []Amount{b.a0, b.a1}
}

func (b Balance) Amount(token Token) Amount {
	if token == Token0 {
		return b.a0
	}
	return b.a1
}

func (b *Balance) SetAmount(token Token, amount Amount) {
	if token == Token0 {
		b.a0 = amount
	} else {
		b.a1 = amount
	}
}

func (b *Balance) Clone() *Balance {
	return &Balance{b.a0, b.a1}
}

func (b Balance) Big() *BigBalance {
	var bv BigBalance
	bv.a0.SetUint64(uint64(b.a0))
	bv.a1.SetUint64(uint64(b.a1))
	return &bv
}

func (b *Balance) Add(balance *Balance) *Balance {
	b.a0 += balance.a0
	b.a1 += balance.a1
	return b
}

func (b *Balance) Sub(balance *Balance) *Balance {
	b.a0 -= balance.a0
	b.a1 -= balance.a1
	return b
}

type BigBalance struct {
	a0, a1 big.Int
}

var (
	BigBalanceEmtpy = BigBalance{}
)

func (b *BigBalance) Clone() *BigBalance {
	var bv BigBalance
	bv.a0.Set(&b.a0)
	bv.a1.Set(&b.a1)
	return &bv
}

func (b *BigBalance) Add(Balance *BigBalance) *BigBalance {
	b.a0.Add(&b.a0, &Balance.a0)
	b.a1.Add(&b.a1, &Balance.a1)
	return b
}

func (b *BigBalance) Sub(Balance *BigBalance) *BigBalance {
	b.a0.Sub(&b.a0, &Balance.a0)
	b.a1.Sub(&b.a1, &Balance.a1)
	return b
}
