// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain

import (
	"math"
	"math/big"
)

// calculate the median price for the nextPrice chained after the block node.
func calcMedianPrice(node *blockNode, epoch int32) float64 {
	if node == nil || node.height+1 < epoch {
		// 1 epoch must pass before the first absorption
		return math.NaN()
	}

	var count int
	var sum float64
	for n := node; n != nil && node.height-n.height < epoch; n = n.parent {
		priceDrv := float64(n.priceDerivation)
		if math.IsNaN(priceDrv) {
			continue
		}
		count++
		sum += priceDrv
	}

	if count*3 < int(epoch)*2 {
		// less than super majority
		return math.NaN()
	}

	// TODO: remove bad price here, e.g. take only 2/3 of valid prices

	return sum / float64(count)
}

// CalcNextAbsorptionRate calculates the absorption for the block
// after the end of the current best chain based on the price history.
//
// This function is safe for concurrent access.
func (b *BlockChain) CalcNextAbsorptionRate() Price {
	b.chainLock.RLock()
	b.stateLock.RLock()
	defer b.stateLock.RUnlock()
	defer b.chainLock.RUnlock()
	return Price(b.calcNextAbsorptionRate())
}

// calcNextAbsorptionRate calculates the absorption for the block
// after the end of the current best chain based on the price history.
func (b *BlockChain) calcNextAbsorptionRate() float64 {
	tip := b.bestChain.Tip()
	// Genesis block.
	if tip == nil {
		return math.NaN()
	}

	epoch := b.chainParams.BlockPerTimespan

	medianPrice := calcMedianPrice(tip, epoch)
	if math.IsNaN(medianPrice) {
		return math.NaN()
	}

	lastAbsnHeight := b.stateSnapshot.LastAbsnHeight
	if tip.height-lastAbsnHeight >= epoch {
		// passive condition: 1 epoch without any active absorption
		// or absorption never occurs, wait for the first epoch pass
		return medianPrice
	}

	// check for active condition
	lastAbsnNode := b.bestChain.NodeByHeight(lastAbsnHeight)
	lastAbsnMedianPrice := calcMedianPrice(lastAbsnNode, epoch)
	if math.IsNaN(lastAbsnMedianPrice) {
		return medianPrice
	}

	priceRate := medianPrice / lastAbsnMedianPrice
	if priceRate >= 2 || priceRate <= -0.5 {
		// active absorption
		return medianPrice
	}

	return math.NaN()
}

// CalcNextAbsorption calculates the next absorption amount of STB
// after the current best chain tip.
func (b *BlockChain) CalcNextAbsorption() *big.Int {
	b.chainLock.RLock()
	b.stateLock.RLock()
	defer b.stateLock.RUnlock()
	defer b.chainLock.RUnlock()
	return b.calcNextAbsorption()
}

func (b *BlockChain) calcNextAbsorption() *big.Int {
	tip := b.bestChain.Tip()
	// Genesis block.
	if tip == nil {
		return nil
	}

	epoch := b.chainParams.BlockPerTimespan
	lastAbsnHeight := b.stateSnapshot.LastAbsnHeight

	remainBlockToAbsorb := epoch - (tip.height - lastAbsnHeight)
	if remainBlockToAbsorb <= 0 {
		// absorption only occurs for 1 week
		return nil
	}

	lastAbsnSupply := b.stateSnapshot.LastAbsnSupply
	alreadyAbsorbed := new(big.Int).Set(&b.stateSnapshot.TotalSupply)
	alreadyAbsorbed.Sub(alreadyAbsorbed, &lastAbsnSupply)

	lastAbsnNode := b.bestChain.NodeByHeight(lastAbsnHeight)
	lastAbsnRate := calcMedianPrice(lastAbsnNode, epoch)

	lastAbsnFloat := new(big.Float).SetInt(&lastAbsnSupply)
	lastAbsnFloat.Mul(lastAbsnFloat, big.NewFloat(lastAbsnRate))
	lastAbsn, _ := lastAbsnFloat.Int(nil)

	remainAbsn := new(big.Int).Sub(lastAbsn, alreadyAbsorbed)
	remainAbsn.Div(remainAbsn, big.NewInt(int64(remainBlockToAbsorb)))

	return remainAbsn
}
