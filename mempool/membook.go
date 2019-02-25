// Copyright (c) 2013-2016 The endurio developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mempool

import (
	"container/list"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/endurio/ndrd/blockchain"
	"github.com/endurio/ndrd/btcjson"
	"github.com/endurio/ndrd/chaincfg"
	"github.com/endurio/ndrd/chaincfg/chainhash"
	"github.com/endurio/ndrd/mining"
	"github.com/endurio/ndrd/txscript"
	"github.com/endurio/ndrd/util"
	"github.com/endurio/ndrd/wire"
)

// OdrDesc is a descriptor containing an order in the mempool along with
// additional metadata.
type OdrDesc struct {
	mining.OdrDesc
}

// Price returns the order price.
func (oD *OdrDesc) Price() float64 {
	return float64(oD.Payout) / float64(oD.Amount)
}

// OrderBookResult returns OrderBookResult object for the order
func (oD *OdrDesc) OrderBookResult() *btcjson.GetOrderBookResult {
	return &btcjson.GetOrderBookResult{
		Bid:    oD.Bid,
		Price:  oD.Price(),
		Amount: oD.Amount.ToBTC(),
	}
}

// OdrBook ...
type OdrBook struct {
	// The following variables must only be used atomically.
	lastUpdated int64 // last time pool was updated

	mtx sync.RWMutex
	cfg Config

	bids      *list.List
	asks      *list.List
	book      map[chainhash.Hash]*list.Element
	outpoints map[wire.OutPoint]*list.Element
}

// Ensure the OdrBook type implements the mining.OdrSource interface.
var _ mining.OdrSource = (*OdrBook)(nil)

// isOrderInBook returns whether or not the passed order already
// exists in the main book.
//
// This function MUST be called with the meobook lock held (for reads).
func (ob *OdrBook) isOrderInBook(hash *chainhash.Hash) bool {
	if _, exists := ob.book[*hash]; exists {
		return true
	}

	return false
}

// IsOrderInBook returns whether or not the passed order already
// exists in the main book.
//
// This function is safe for concurrent access.
func (ob *OdrBook) IsOrderInBook(hash *chainhash.Hash) bool {
	// Protect concurrent access.
	ob.mtx.RLock()
	inBook := ob.isOrderInBook(hash)
	ob.mtx.RUnlock()

	return inBook
}

// haveOrder returns whether or not the passed order already exists
// in the main book or in the orphan pool.
//
// This function MUST be called with the mempool lock held (for reads).
func (ob *OdrBook) haveOrder(hash *chainhash.Hash) bool {
	return ob.isOrderInBook(hash)
}

// HaveOrder returns whether or not the passed order already exists
// in the main book or in the orphan pool.
//
// This function is safe for concurrent access.
func (ob *OdrBook) HaveOrder(hash *chainhash.Hash) bool {
	// Protect concurrent access.
	ob.mtx.RLock()
	haveTx := ob.haveOrder(hash)
	ob.mtx.RUnlock()

	return haveTx
}

// removeOrder is the internal function which implements the public
// RemoveOrder.  See the comment for RemoveOrder for more details.
//
// This function MUST be called with the mempool lock held (for writes).
func (ob *OdrBook) removeOrder(order *util.Odr) {
	txHash := order.Hash()

	// Remove the order if needed.
	if element, exists := ob.book[*txHash]; exists {
		// Mark the referenced outpoints as unspent by the pool.
		for _, txIn := range order.TxIn {
			delete(ob.outpoints, txIn.PreviousOutPoint)
		}
		delete(ob.book, *txHash)
		if element.Value.(*OdrDesc).Bid {
			ob.bids.Remove(element)
		} else {
			ob.asks.Remove(element)
		}

		atomic.StoreInt64(&ob.lastUpdated, time.Now().Unix())
	}
}

// RemoveOrder removes the passed order from the mempool. When the
// removeRedeemers flag is set, any orders that redeem outputs from the
// removed order will also be removed recursively from the mempool, as
// they would otherwise become orphans.
//
// This function is safe for concurrent access.
func (ob *OdrBook) RemoveOrder(order *util.Odr) {
	// Protect concurrent access.
	ob.mtx.Lock()
	ob.removeOrder(order)
	ob.mtx.Unlock()
}

// RemoveDoubleSpends removes all orders which spend outputs spent by the
// passed order from the memory pool.  Removing those orders then
// leads to removing all orders which rely on them, recursively.  This is
// necessary when a block is connected to the main chain because the block may
// contain orders which were previously unknown to the memory pool.
//
// This function is safe for concurrent access.
func (ob *OdrBook) RemoveDoubleSpends(tx *util.Tx) {
	// Protect concurrent access.
	ob.mtx.Lock()
	for _, txIn := range tx.MsgTx().TxIn {
		if elementRedeemer, ok := ob.outpoints[txIn.PreviousOutPoint]; ok {
			odrDesc := elementRedeemer.Value.(*OdrDesc)
			if !odrDesc.Hash().IsEqual(tx.Hash()) {
				ob.removeOrder(odrDesc.Odr)
			}
		}
	}
	ob.mtx.Unlock()
}

func insertOrder(orders *list.List, orderDesc *OdrDesc) *list.Element {
	price := orderDesc.Price()
	// ordered insert
	e := orders.Front()
	for ; e != nil; e = e.Next() {
		o := e.Value.(*OdrDesc)

		if orderDesc.Bid {
			// biding orders sorted from highest bidder down
			if price > o.Price() {
				break
			}
		} else {
			// asking orders sorted from lowser asker up
			if price < o.Price() {
				break
			}
		}
	}
	if e == nil {
		return orders.PushBack(orderDesc)
	}
	return orders.InsertBefore(orderDesc, e)
}

// addOrder adds the passed odr to the memory pool.  It should
// not be called directly as it doesn't perform any validation.  This is a
// helper for maybeAcceptOrder.
//
// This function MUST be called with the mempool lock held (for writes).
func (ob *OdrBook) addOrder(odr *util.Odr, stb, ndr int64, height int32) *OdrDesc {
	odrDesc := &OdrDesc{
		OdrDesc: mining.OdrDesc{
			Odr:    odr,
			Added:  time.Now(),
			Height: height,
			Bid:    ndr > 0,
			Amount: util.Amount(abs(ndr)),
			Payout: util.Amount(abs(stb)),
		},
	}

	var element *list.Element
	if odrDesc.Bid {
		element = insertOrder(ob.bids, odrDesc)
	} else {
		element = insertOrder(ob.asks, odrDesc)
	}

	ob.book[*odrDesc.Hash()] = element
	for _, txIn := range odrDesc.TxIn {
		ob.outpoints[txIn.PreviousOutPoint] = element
	}

	atomic.StoreInt64(&ob.lastUpdated, time.Now().Unix())
	return odrDesc
}

// checkBookDoubleSpend checks whether or not the passed order is
// attempting to spend coins already spent by other orders in the pool.
// Note it does not check for double spends against orders already in the
// main chain.
//
// This function MUST be called with the mempool lock held (for reads).
func (ob *OdrBook) checkBookDoubleSpend(order *util.Odr) error {
	for _, txIn := range order.TxIn {
		if element, exists := ob.outpoints[txIn.PreviousOutPoint]; exists {
			str := fmt.Sprintf("output %v already spent by order %v in the memory pool",
				txIn.PreviousOutPoint, element.Value.(*OdrDesc).Hash())
			return txRuleError(wire.RejectDuplicate, str)
		}
	}

	return nil
}

// CheckSpend checks whether the passed outpoint is already spent by a
// order in the mempool. If that's the case the spending order will
// be returned, if not nil will be returned.
func (ob *OdrBook) CheckSpend(op wire.OutPoint) *util.Tx {
	ob.mtx.RLock()
	txR := ob.outpoints[op]
	ob.mtx.RUnlock()

	return txR.Value.(*OdrDesc).Tx
}

// fetchInputUtxos loads utxo details about the input transactions referenced by
// the passed transaction.  First, it loads the details form the viewpoint of
// the main chain, then it adjusts them based upon the contents of the
// transaction pool.
//
// This function MUST be called with the mempool lock held (for reads).
func (ob *OdrBook) fetchInputUtxos(order *util.Odr) (*blockchain.UtxoViewpoint, error) {
	utxoView, err := ob.cfg.FetchUtxoView(order.Tx)
	if err != nil {
		return nil, err
	}

	return utxoView, nil
}

// FetchOrder returns the requested order from the order pool.
// This only fetches from the main order pool and does not include
// orphans.
//
// This function is safe for concurrent access.
func (ob *OdrBook) FetchOrder(txHash *chainhash.Hash) (*util.Odr, error) {
	// Protect concurrent access.
	ob.mtx.RLock()
	element, exists := ob.book[*txHash]
	ob.mtx.RUnlock()

	if exists {
		return element.Value.(*OdrDesc).Odr, nil
	}

	return nil, fmt.Errorf("order is not in the pool")
}

// maybeAcceptOrder is the internal function which implements the public
// MaybeAcceptOrder.  See the comment for MaybeAcceptOrder for
// more details.
//
// This function MUST be called with the mempool lock held (for writes).
func (ob *OdrBook) maybeAcceptOrder(order *util.Odr) (*OdrDesc, error) {
	txHash := order.Hash()

	// If a transaction has iwtness data, and segwit isn't active yet, If
	// segwit isn't active yet, then we won't accept it into the mempool as
	// it can't be mined yet.
	if order.HasWitness() {
		segwitActive, err := ob.cfg.IsDeploymentActive(chaincfg.DeploymentSegwit)
		if err != nil {
			return nil, err
		}

		if !segwitActive {
			str := fmt.Sprintf("order %v has witness data, "+
				"but segwit isn't active yet", txHash)
			return nil, txRuleError(wire.RejectNonstandard, str)
		}
	}

	// Don't accept the transaction if it already exists in the pool.  This
	// applies to orphan transactions as well when the reject duplicate
	// orphans flag is set.  This check is intended to be a quick check to
	// weed out duplicates.
	if ob.isOrderInBook(txHash) {
		str := fmt.Sprintf("already have order %v", txHash)
		return nil, txRuleError(wire.RejectDuplicate, str)
	}

	// Perform preliminary sanity checks on the transaction.  This makes
	// use of blockchain which contains the invariant rules for what
	// transactions are allowed into blocks.
	err := blockchain.CheckTransactionSanity(order.Tx)
	if err != nil {
		if cerr, ok := err.(blockchain.RuleError); ok {
			return nil, chainRuleError(cerr)
		}
		return nil, err
	}

	// An order must not be a coinbase transaction.
	if blockchain.IsCoinBase(order.Tx) {
		str := fmt.Sprintf("order %v is a coinbase", txHash)
		return nil, txRuleError(wire.RejectInvalid, str)
	}

	// Get the current height of the main chain.
	bestHeight := ob.cfg.BestHeight()
	nextBlockHeight := bestHeight + 1

	medianTimePast := ob.cfg.MedianTimePast()

	// Don't allow non-standard transactions if the network parameters
	// forbid their acceptance.
	if !ob.cfg.Policy.AcceptNonStd {
		err = checkTransactionStandard(order.Tx, nextBlockHeight,
			medianTimePast, ob.cfg.Policy.MinRelayTxFee,
			ob.cfg.Policy.MaxTxVersion)
		if err != nil {
			// Attempt to extract a reject code from the error so
			// it can be retained.  When not possible, fall back to
			// a non standard error.
			rejectCode, found := extractRejectCode(err)
			if !found {
				rejectCode = wire.RejectNonstandard
			}
			str := fmt.Sprintf("order %v is not standard: %v",
				txHash, err)
			return nil, txRuleError(rejectCode, str)
		}
	}
	// The transaction may not use any of the same outputs as other
	// transactions already in the pool as that would ultimately result in a
	// double spend.  This check is intended to be quick and therefore only
	// detects double spends within the transaction pool itself.  The
	// transaction could still be double spending coins from the main chain
	// at this point.  There is a more in-depth check that happens later
	// after fetching the referenced transaction inputs from the main chain
	// which examines the actual spend data and prevents double spends.
	err = ob.checkBookDoubleSpend(order)
	if err != nil {
		return nil, err
	}

	utxoView, err := ob.fetchInputUtxos(order)
	if err != nil {
		if cerr, ok := err.(blockchain.RuleError); ok {
			return nil, chainRuleError(cerr)
		}
		return nil, err
	}

	// Don't allow the transaction into the mempool unless its sequence
	// lock is active, meaning that it'll be allowed into the next block
	// with respect to its defined relative lock times.
	sequenceLock, err := ob.cfg.CalcSequenceLock(order.Tx, utxoView)
	if err != nil {
		if cerr, ok := err.(blockchain.RuleError); ok {
			return nil, chainRuleError(cerr)
		}
		return nil, err
	}
	if !blockchain.SequenceLockActive(sequenceLock, nextBlockHeight,
		medianTimePast) {
		return nil, txRuleError(wire.RejectNonstandard,
			"order's sequence locks on inputs not met")
	}

	// Perform several checks on the transaction inputs using the invariant
	// rules in blockchain for what transactions are allowed into blocks.
	// Also returns the fees associated with the transaction which will be
	// used later.
	balances, err := blockchain.CheckTransactionInputs(order.Tx, nextBlockHeight,
		utxoView, ob.cfg.ChainParams)
	if err != nil {
		if cerr, ok := err.(blockchain.RuleError); ok {
			return nil, chainRuleError(cerr)
		}
		return nil, err
	}

	if balances[wire.STB] <= 0 && balances[wire.NDR] <= 0 {
		str := fmt.Sprintf("Not an order: %v", txHash)
		return nil, blockchain.RuleError{
			ErrorCode:   blockchain.ErrNotAnOrder,
			Description: str,
		}
	}

	// Don't allow transactions with non-standard inputs if the network
	// parameters forbid their acceptance.
	if !ob.cfg.Policy.AcceptNonStd {
		err := checkInputsStandard(order.Tx, utxoView)
		if err != nil {
			// Attempt to extract a reject code from the error so
			// it can be retained.  When not possible, fall back to
			// a non standard error.
			rejectCode, found := extractRejectCode(err)
			if !found {
				rejectCode = wire.RejectNonstandard
			}
			str := fmt.Sprintf("order %v has a non-standard "+
				"input: %v", txHash, err)
			return nil, txRuleError(rejectCode, str)
		}
	}

	// NOTE: if you modify this code to accept non-standard transactions,
	// you should add code here to check that the transaction does a
	// reasonable nuober of ECDSA signature verifications.

	// Don't allow transactions with an excessive nuober of signature
	// operations which would result in making it impossible to mine.  Since
	// the coinbase address itself can contain signature operations, the
	// maximum allowed signature operations per transaction is less than
	// the maximum allowed signature operations per block.
	// TODO(roasbeef): last bool should be conditional on segwit activation
	sigOpCost, err := blockchain.GetSigOpCost(order.Tx, false, utxoView, true, true)
	if err != nil {
		if cerr, ok := err.(blockchain.RuleError); ok {
			return nil, chainRuleError(cerr)
		}
		return nil, err
	}
	if sigOpCost > ob.cfg.Policy.MaxSigOpCostPerTx {
		str := fmt.Sprintf("order %v sigop cost is too high: %d > %d",
			txHash, sigOpCost, ob.cfg.Policy.MaxSigOpCostPerTx)
		return nil, txRuleError(wire.RejectNonstandard, str)
	}

	// Verify crypto signatures for each input and reject the transaction if
	// any don't verify.
	err = blockchain.ValidateTransactionScripts(order.Tx, utxoView,
		txscript.StandardVerifyFlags, ob.cfg.SigCache,
		ob.cfg.HashCache)
	if err != nil {
		if cerr, ok := err.(blockchain.RuleError); ok {
			return nil, chainRuleError(cerr)
		}
		return nil, err
	}

	// Add to transaction pool.
	oD := ob.addOrder(order, balances[wire.STB], balances[wire.NDR], bestHeight)

	log.Debugf("Accepted order %v (book size: %v)", txHash, len(ob.book))

	return oD, nil
}

// MaybeAcceptOrder is the main workhorse for handling insertion of new
// free-standing orders into a order book.  It includes functionality
// such as rejecting duplicate orders, ensuring orders follow all
// rules, detecting orphan orders, and insertion into the order book.
//
// This function is safe for concurrent access.
func (ob *OdrBook) MaybeAcceptOrder(order *util.Odr) (*OdrDesc, error) {
	// Protect concurrent access.
	ob.mtx.Lock()
	oD, err := ob.maybeAcceptOrder(order)
	ob.mtx.Unlock()

	return oD, err
}

// ProcessOrder is the main workhorse for handling insertion of new
// free-standing orders into the memory pool.  It includes functionality
// such as rejecting duplicate orders, ensuring transactions follow all
// rules, orphan order handling, and insertion into the memory pool.
//
// It returns a slice of orders added to the mempool.  When the
// error is nil, the list will include the passed order itself along
// with any additional orphan orders that were added as a result of
// the passed one being accepted.
//
// This function is safe for concurrent access.
func (ob *OdrBook) ProcessOrder(order *util.Odr) (*OdrDesc, error) {
	log.Tracef("Processing order %v", order.Hash())

	// Protect concurrent access.
	ob.mtx.Lock()
	defer ob.mtx.Unlock()

	// Potentially accept the transaction to the memory pool.
	oD, err := ob.maybeAcceptOrder(order)
	if err != nil {
		return nil, err
	}

	return oD, nil
}

// Count returns the nuober of orders in the main book.  It does not
// include the orphan pool.
//
// This function is safe for concurrent access.
func (ob *OdrBook) Count() int {
	ob.mtx.RLock()
	count := len(ob.book)
	ob.mtx.RUnlock()

	return count
}

// TxHashes returns a slice of hashes for all of the orders in the memory
// pool.
//
// This function is safe for concurrent access.
func (ob *OdrBook) TxHashes() []*chainhash.Hash {
	ob.mtx.RLock()
	hashes := make([]*chainhash.Hash, len(ob.book))
	i := 0
	for hash := range ob.book {
		hashCopy := hash
		hashes[i] = &hashCopy
		i++
	}
	ob.mtx.RUnlock()

	return hashes
}

// OdrDescs returns a slice of descriptors for all the orders in the book.
// The descriptors are to be treated as read only.
//
// This function is safe for concurrent access.
func (ob *OdrBook) OdrDescs() []*OdrDesc {
	ob.mtx.RLock()
	descs := make([]*OdrDesc, len(ob.book))
	i := 0
	for _, element := range ob.book {
		descs[i] = element.Value.(*OdrDesc)
		i++
	}
	ob.mtx.RUnlock()

	return descs
}

// MiningDescs returns a slice of mining descriptors for the orders
// in the book, to spend as much as possible the amount of STB.
//
// This is part of the mining.OdrSource interface implementation and is safe for
// concurrent access as required by the interface contract.
func (ob *OdrBook) MiningDescs(payout *big.Int) []*mining.OdrDesc {
	ob.mtx.RLock()
	defer ob.mtx.RUnlock()

	orders := ob.asks
	if payout.Sign() < 0 {
		orders = ob.bids
		payout = new(big.Int).Abs(payout)
	}

	descs, err := getOrdersForPayout(orders, payout)
	if err != nil {
		log.Errorf("Unable to get orders for mining: %v", err)
		return nil
	}
	result := make([]*mining.OdrDesc, len(descs))
	i := 0
	for _, desc := range descs {
		result[i] = &desc.OdrDesc
		i++
	}
	return result
}

// Get orders to spend as much as possible of STB amount.
// not thread safe
func getOrdersForPayout(orders *list.List, payout *big.Int) ([]*OdrDesc, error) {
	result := make([]*OdrDesc, 0, orders.Len())
	remain := new(big.Int).Set(payout)

	for e := orders.Front(); e != nil; e = e.Next() {
		odrDesc := e.Value.(*OdrDesc)
		remain.Sub(remain, big.NewInt(abs(int64(odrDesc.Payout))))
		if remain.Sign() < 0 {
			break
		}
		result = append(result, odrDesc)
	}

	return result, nil
}

// RawMembookVerbose returns all of the entries in the mempool as a fully
// populated btcjson result.
//
// This function is safe for concurrent access.
func (ob *OdrBook) RawMembookVerbose() map[string]*btcjson.GetRawMembookVerboseResult {
	ob.mtx.RLock()
	defer ob.mtx.RUnlock()

	result := make(map[string]*btcjson.GetRawMembookVerboseResult,
		len(ob.book))

	for _, element := range ob.book {
		// Calculate the current priority based on the inputs to
		// the transaction.  Use zero if one or more of the
		// input transactions can't be found for some reason.
		desc := element.Value.(*OdrDesc)
		odr := desc.Odr

		mpd := &btcjson.GetRawMembookVerboseResult{
			Size:    int32(odr.SerializeSize()),
			Vsize:   int32(GetTxVirtualSize(odr.Tx)),
			Depends: make([]string, 0),
		}
		for _, txIn := range odr.TxIn {
			hash := &txIn.PreviousOutPoint.Hash
			if ob.haveOrder(hash) {
				mpd.Depends = append(mpd.Depends,
					hash.String())
			}
		}

		result[odr.Hash().String()] = mpd
	}

	return result
}

// Get orders to cover as much as possible the market depth of NDR.
// not thread safe
func getOrdersForDepth(orders *list.List, depth float64) []*OdrDesc {
	result := make([]*OdrDesc, 0, orders.Len())

	var total float64

	for e := orders.Front(); e != nil; e = e.Next() {
		odrDesc := e.Value.(*OdrDesc)
		result = append(result, odrDesc)

		// limit the list by market depth param
		if depth > 0 {
			total += odrDesc.Amount.ToBTC()
			if total >= depth {
				return result
			}
		}
	}

	return result
}

// OrderBook returns all of the entries in the order book as a fully
// populated btcjson result.
//
// This function is safe for concurrent access.
func (ob *OdrBook) OrderBook(depth float64) ([]*btcjson.GetOrderBookResult, error) {
	ob.mtx.RLock()
	defer ob.mtx.RUnlock()

	asks := getOrdersForDepth(ob.asks, depth)
	bids := getOrdersForDepth(ob.bids, depth)

	result := make([]*btcjson.GetOrderBookResult, len(asks)+len(bids))

	var idx int
	// asks list is reverted
	for i := len(asks) - 1; i >= 0; i-- {
		odrDesc := asks[i]
		result[idx] = odrDesc.OrderBookResult()
		idx++
	}
	for _, odrDesc := range bids {
		result[idx] = odrDesc.OrderBookResult()
		idx++
	}

	return result, nil
}

// LastUpdated returns the last time a order was added to or removed from
// the main book.  It does not include the orphan pool.
//
// This function is safe for concurrent access.
func (ob *OdrBook) LastUpdated() time.Time {
	return time.Unix(atomic.LoadInt64(&ob.lastUpdated), 0)
}

// NewMemBook returns a new order book for validating and storing standalone
// orders until they are matched and mined into a block.
func NewMemBook(cfg *Config) *OdrBook {
	return &OdrBook{
		cfg:       *cfg,
		book:      make(map[chainhash.Hash]*list.Element),
		bids:      list.New(),
		asks:      list.New(),
		outpoints: make(map[wire.OutPoint]*list.Element),
	}
}

func abs(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}
