// Copyright (c) 2013-2014 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain

import (
	"math"
	"sync"
	"time"
)

var (
	// maxPriceEntries is the maximum number of entries allowed in the
	// median time data.  This is a variable as opposed to a constant so the
	// test code can modify it.
	maxPriceEntries = 1024
)

// Price ...
type Price float64

// PriceDesc ...
type PriceDesc struct {
	Price
	Timestamp time.Time
}

// FeedPriceSource provides a mechanism to add several time samples which are
// used to determine a median time which is then used as an offset to the local
// clock.
type FeedPriceSource interface {
	// LastPrice ...
	LastPrice() *PriceDesc
	PriceToMine() Price
	FeedPrice(price Price)
}

// feedPrice ...
type feedPrice struct {
	mtx           sync.Mutex
	blockTime     time.Duration
	lastPriceDesc *PriceDesc
}

// Ensure the feedPrice type implements the FeedPriceSource interface.
var _ FeedPriceSource = (*feedPrice)(nil)

// LastPrice ...
func (fp *feedPrice) LastPrice() *PriceDesc {
	return fp.lastPriceDesc
}

func (fp *feedPrice) PriceToMine() Price {
	if fp.lastPriceDesc == nil {
		return Price(math.NaN())
	}
	duration := time.Since(fp.lastPriceDesc.Timestamp)
	if duration < 0 || fp.blockTime < duration {
		return Price(math.NaN())
	}
	return fp.lastPriceDesc.Price
}

func (fp *feedPrice) FeedPrice(price Price) {
	fp.lastPriceDesc = &PriceDesc{
		Price:     price - 1.0,
		Timestamp: time.Now(),
	}
}

// NewFeedPrice ...
func NewFeedPrice(blockTime time.Duration) FeedPriceSource {
	return &feedPrice{
		blockTime: blockTime,
	}
}
