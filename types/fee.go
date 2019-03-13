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
type Price Balance

func (p Price) Rate(r PriceRate) (rate float64) {
	if r.a0 > 0 {
		rate += float64(p.a0) / float64(r.a0)
	}
	if r.a1 > 0 {
		rate += float64(p.a1) / float64(r.a1)
	}
	return rate
}

type PriceRate Price
