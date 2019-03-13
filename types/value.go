// Copyright (c) 2018-2019 The Endurio developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package types

type Value struct {
	Amount
	Token
}

var (
	ValueEmpty = Value{}
	ValueDummy = Value{-1, TokenInvalid}
)

// RangeCheck checks whether the value is in it's valid range.
// Returns 0 for a valid range, negative for lower than minimum,
// and positive value for higher than maximum.
func (v *Value) RangeCheck() int {
	if v.Amount < 0 {
		return -1
	}
	if v.Amount > MaxAtom {
		return 1
	}
	return 0
}
