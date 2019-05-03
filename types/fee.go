// Copyright (c) 2013, 2014 The btcsuite developers
// Copyright (c) 2018-2019 The Endurio developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package types

// Fee represents the fee of an transaction.
type Fee Balance

var (
	FeeEmpty = Fee{}
	FeeDummy = Fee{-1, -1}
)

func (f *Fee) Balance() *Balance {
	return (*Balance)(f)
}

func (f *Fee) Add(g *Fee) *Fee {
	return f.Balance().Add(g.Balance()).Fee()
}

// Price represents the fee over KB of transaction size.
// Price.a(i) = Tx.Fee.a(i) / Tx.SizeInKB
type Price Balance

// MinTokenPrice is the miner-configured price rate accepted for each tokens.
// MinTokenPrice.a(i) == 0 when miner does not accept fee in token(i)
// MinTokenPrice.a(i) > 0 then the tx is accepted when it pays at least a(i) fee of token(i)
// Multiple prices can be accepted, so the tx can be paid by multiple tokens.
// For each pair of accepted tokens, a(i)/a(j) should be equal market_price(j)/market_price(i).
type MinTokenPrice Price

// Rate calculates the price rate paid for a tx.
// Tx with higher rate will have higher priority.
// A rate >= 1.0 is sufficient for an tx to be accepted.
func (p Price) Rate(r MinTokenPrice) (rate float64) {
	if r.a0 > 0 {
		rate += float64(p.a0) / float64(r.a0)
	}
	if r.a1 > 0 {
		rate += float64(p.a1) / float64(r.a1)
	}
	return rate
}
