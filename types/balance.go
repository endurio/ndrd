// Copyright (c) 2018-2019 The Endurio developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package types

import (
	"fmt"
	"math/big"
)

// Balance contains all token amounts carried by an tx output.
// Designed as map of Token => Amount, but currently be a struct of 2 fields for optimization.
type Balance struct {
	a0, a1 Amount
}

var (
	BalanceEmpty = Balance{}
	BalanceDummy = Balance{-1, -1}
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

func (b *Balance) Value(token Token) Value {
	if token == Token0 {
		return Value{b.a0, Token0}
	}
	return Value{b.a1, Token1}
}

func (b *Balance) Values() []Value {
	values := make([]Value, 2)
	if b.a0 != 0 {
		values = append(values, Value{b.a0, Token0})
	}
	if b.a1 != 0 {
		values = append(values, Value{b.a1, Token1})
	}
	return values
}

func (b *Balance) Fee() *Fee {
	return (*Fee)(b)
}

func (b *Balance) String() string {
	return fmt.Sprintf("(%v,%v)", b.a0, b.a1)
}

func (b *Balance) Clone() *Balance {
	return &Balance{b.a0, b.a1}
}

func (b *Balance) AddValue(v Value) *Balance {
	switch v.Token {
	case Token0:
		b.a0 += v.Amount
	case Token1:
		b.a1 += v.Amount
	default:
		panic(fmt.Sprintf("Unknown Token: %v", v.Token))
	}
	return b
}

func (b *Balance) SubValue(v Value) *Balance {
	switch v.Token {
	case Token0:
		b.a0 -= v.Amount
	case Token1:
		b.a1 -= v.Amount
	default:
		panic(fmt.Sprintf("Unknown Token: %v", v.Token))
	}
	return b
}

func (b *Balance) Add(c *Balance) *Balance {
	b.a0 += c.a0
	b.a1 += c.a1
	return b
}

func (b *Balance) Sub(c *Balance) *Balance {
	b.a0 -= c.a0
	b.a1 -= c.a1
	return b
}

func (b *Balance) Neg() *Balance {
	b.a0 = -b.a0
	b.a1 = -b.a1
	return b
}

// RangeCheck checks whether the value is in it's valid range.
// Returns 0 for a valid range, negative for lower than minimum,
// and positive value for higher than maximum.
func (b *Balance) RangeCheck() int {
	var check int
	if b.a0 < 0 {
		check = -1
	} else if b.a0 > MaxAtom {
		check = 1
	}
	if b.a1 < 0 {
		check += -2
	} else if b.a1 > MaxAtom {
		check += 2
	}
	return check
}

// SafeAdd perform Add with overflows check.
func (b *Balance) SafeAdd(c *Balance) error {
	if c.a0 != 0 {
		result := b.a0 + c.a0
		if (c.a0 > 0 && result < b.a0) ||
			(c.a0 < 0 && result > b.a0) {
			return fmt.Errorf("balances addition overflows		token: %v, a: %v, b: %v",
				Token0, b.a0, c.a0)
		}
	}
	if c.a1 != 0 {
		result := b.a1 + c.a1
		if (c.a1 > 0 && result < b.a1) ||
			(c.a1 < 0 && result > b.a1) {
			return fmt.Errorf("balances addition overflows		token: %v, a: %v, b: %v",
				Token1, b.a1, c.a1)
		}
	}
	b.Add(c)
	return nil
}

// Cover returns c >= b
func (b *Balance) Cover(c *Balance) bool {
	return b.a0 >= c.a0 && b.a1 >= c.a1
}

func (b Balance) Big() *BigBalance {
	var bb BigBalance
	bb.a0.SetUint64(uint64(b.a0))
	bb.a1.SetUint64(uint64(b.a1))
	return &bb
}

type BigBalance struct {
	a0, a1 big.Int
}

var (
	BigBalanceEmtpy = BigBalance{}
)

func (b *BigBalance) Clone() *BigBalance {
	var bb BigBalance
	bb.a0.Set(&b.a0)
	bb.a1.Set(&b.a1)
	return &bb
}

func (b *BigBalance) Add(bb *BigBalance) *BigBalance {
	b.a0.Add(&b.a0, &bb.a0)
	b.a1.Add(&b.a1, &bb.a1)
	return b
}

func (b *BigBalance) Sub(bb *BigBalance) *BigBalance {
	b.a0.Sub(&b.a0, &bb.a0)
	b.a1.Sub(&b.a1, &bb.a1)
	return b
}
