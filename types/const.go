// Copyright (c) 2013-2014 The btcsuite developers
// Copyright (c) 2018-2019 The Endurio developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package types

const (
	// AtomPerCent is the number of satoshi in one bitcoin cent.
	AtomPerCent = 1e6

	// AtomPerCoin is the number of satoshi in one bitcoin (1 BTC).
	AtomPerCoin = 1e8

	// MaxAtom is the maximum transaction amount allowed in satoshi.
	MaxAtom = 21e6 * AtomPerCoin
)