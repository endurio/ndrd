// Copyright (c) 2018-2019 The Endurio developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package types

type Token uint8

const (
	Token0 = Token(0x0)
	Token1 = Token(0x1)

	TokenInvalid = Token(0xFF)
)

func (t Token) String() string {
	switch t {
	case Token0:
		return "NDR"
	case Token1:
		return "STB"
	default:
		return "UnknownToken"
	}
}
