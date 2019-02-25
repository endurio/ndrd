// Copyright (c) 2017 The endurio developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package netsync

import (
	"github.com/endurio/ndrd/blockchain"
	"github.com/endurio/ndrd/chaincfg"
	"github.com/endurio/ndrd/chaincfg/chainhash"
	"github.com/endurio/ndrd/mempool"
	"github.com/endurio/ndrd/peer"
	"github.com/endurio/ndrd/chainutil"
	"github.com/endurio/ndrd/wire"
)

// PeerNotifier exposes methods to notify peers of status changes to
// transactions, blocks, etc. Currently server (in the main package) implements
// this interface.
type PeerNotifier interface {
	AnnounceNewTransactions(newTxs []*mempool.TxDesc)

	AnnounceNewOrders(newOrders []*mempool.OdrDesc)

	UpdatePeerHeights(latestBlkHash *chainhash.Hash, latestHeight int32, updateSource *peer.Peer)

	RelayInventory(invVect *wire.InvVect, data interface{})

	TransactionConfirmed(tx *chainutil.Tx)

	OrderFilled(odr *chainutil.Odr)
}

// Config is a configuration struct used to initialize a new SyncManager.
type Config struct {
	PeerNotifier PeerNotifier
	Chain        *blockchain.BlockChain
	TxMemPool    *mempool.TxPool
	OdrMemBook   *mempool.OdrBook
	ChainParams  *chaincfg.Params

	DisableCheckpoints bool
	MaxPeers           int

	FeeEstimator *mempool.FeeEstimator
}
