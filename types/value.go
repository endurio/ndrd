// Copyright (c) 2018-2019 The Endurio developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package types

type Value struct {
	Amount
	Token
}

var (
	ValueInvalid = Value{-1, TokenInvalid}
)
