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

func (v Balance) Map() []Amount {
	return []Amount{v.a0, v.a1}
}

func (v Balance) Amount(token Token) Amount {
	if token == Token0 {
		return v.a0
	}
	return v.a1
}

func (v *Balance) SetAmount(token Token, amount Amount) {
	if token == Token0 {
		v.a0 = amount
	} else {
		v.a1 = amount
	}
}

func (v *Balance) Clone() *Balance {
	return &Balance{v.a0, v.a1}
}

func (v Balance) Big() *BigBalance {
	var bv BigBalance
	bv.a0.SetUint64(uint64(v.a0))
	bv.a1.SetUint64(uint64(v.a1))
	return &bv
}

func (v *Balance) Add(Balance *Balance) *Balance {
	v.a0 += Balance.a0
	v.a1 += Balance.a1
	return v
}

func (v *Balance) Sub(Balance *Balance) *Balance {
	v.a0 -= Balance.a0
	v.a1 -= Balance.a1
	return v
}

type BigBalance struct {
	a0, a1 big.Int
}

var (
	BigBalanceEmtpy = BigBalance{}
)

func (v *BigBalance) Clone() *BigBalance {
	var bv BigBalance
	bv.a0.Set(&v.a0)
	bv.a1.Set(&v.a1)
	return &bv
}

func (v *BigBalance) Add(Balance *BigBalance) *BigBalance {
	v.a0.Add(&v.a0, &Balance.a0)
	v.a1.Add(&v.a1, &Balance.a1)
	return v
}

func (v *BigBalance) Sub(Balance *BigBalance) *BigBalance {
	v.a0.Sub(&v.a0, &Balance.a0)
	v.a1.Sub(&v.a1, &Balance.a1)
	return v
}
