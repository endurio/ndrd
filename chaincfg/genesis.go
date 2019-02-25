// Copyright (c) 2014-2016 The endurio developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

import (
	"time"

	"github.com/endurio/ndrd/chaincfg/chainhash"
	"github.com/endurio/ndrd/wire"
)

// genesisCoinbaseTx is the coinbase transaction for the genesis blocks for
// the main network, regression test network, and test network (version 3).
var genesisCoinbaseTx = wire.MsgTx{
	Version: 1,
	TxIn: []*wire.TxIn{
		{
			PreviousOutPoint: wire.OutPoint{
				Hash:  chainhash.Hash{},
				Index: 0xffffffff,
			},
			SignatureScript: []byte{
				0x56, 0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6c, 0x69,
				0x74, 0x79, 0x20, 0x63, 0x61, 0x6e, 0x20, 0x6e,
				0x65, 0x69, 0x74, 0x68, 0x65, 0x72, 0x20, 0x62,
				0x65, 0x20, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
				0x64, 0x20, 0x6e, 0x6f, 0x72, 0x20, 0x64, 0x65,
				0x73, 0x74, 0x72, 0x6f, 0x79, 0x65, 0x64, 0x3b,
				0x20, 0x76, 0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6c,
				0x69, 0x74, 0x79, 0x20, 0x63, 0x61, 0x6e, 0x20,
				0x6f, 0x6e, 0x6c, 0x79, 0x20, 0x62, 0x65, 0x20,
				0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72,
				0x72, 0x65, 0x64, 0x20, 0x66, 0x72, 0x6f, 0x6d,
				0x20, 0x6f, 0x6e, 0x65, 0x20, 0x74, 0x6f, 0x6b,
				0x65, 0x6e, 0x20, 0x74, 0x6f, 0x20, 0x61, 0x6e,
				0x6f, 0x74, 0x68, 0x65, 0x72, 0x2e,
			},
			Sequence: 0xffffffff,
		},
	},
	TxOut: []*wire.TxOut{
		{
			Value: PreminedSTB,
			PkScript: []byte{
				0x76,                                           // OP_DUP
				0xA9,                                           // OP_HASH160,
				0x14,                                           // bytes to push
				0x4E, 0xDA, 0xBF, 0xD1, 0x85, 0x94, 0xEC, 0xCA, // PKH: ANxpQRnzNGhWEL3qvfn8siWsi1ai846fWr
				0x70, 0xC7, 0x1C, 0x31, 0xF1, 0xF5, 0x69, 0xA6,
				0x9A, 0x3C, 0xFA, 0x17,
				0x88, // OP_EQUALVERIFY,
				0xAC, // OP_CHECKSIG,
			},
		},
		{
			Value: PreminedNDR,
			PkScript: []byte{
				wire.OP_NDR,
				0x76,                                           // OP_DUP
				0xA9,                                           // OP_HASH160,
				0x14,                                           // bytes to push
				0x4E, 0xDA, 0xBF, 0xD1, 0x85, 0x94, 0xEC, 0xCA, // PKH: ANxpQRnzNGhWEL3qvfn8siWsi1ai846fWr
				0x70, 0xC7, 0x1C, 0x31, 0xF1, 0xF5, 0x69, 0xA6,
				0x9A, 0x3C, 0xFA, 0x17,
				0x88, // OP_EQUALVERIFY,
				0xAC, // OP_CHECKSIG,
			},
		},
	},
	LockTime: 0,
}

// genesisHash is the hash of the first block in the block chain for the main
// network (genesis block).
var genesisHash = genesisBlock.BlockHash()

// genesisMerkleRoot is the hash of the first transaction in the genesis block
// for the main network.
var genesisMerkleRoot = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x3b, 0xa3, 0xed, 0xfd, 0x7a, 0x7b, 0x12, 0xb2,
	0x7a, 0xc7, 0x2c, 0x3e, 0x67, 0x76, 0x8f, 0x61,
	0x7f, 0xc8, 0x1b, 0xc3, 0x88, 0x8a, 0x51, 0x32,
	0x3a, 0x9f, 0xb8, 0xaa, 0x4b, 0x1e, 0x5e, 0x4a,
})

// genesisBlock defines the genesis block of the block chain which serves as the
// public transaction ledger for the main network.
var genesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},         // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: genesisMerkleRoot,        // 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
		Timestamp:  time.Unix(0x495fab29, 0), // 2009-01-03 18:15:05 +0000 UTC
		Bits:       0x1d00ffff,               // 486604799 [00000000ffff0000000000000000000000000000000000000000000000000000]
		Nonce:      0x7c2bac1d,               // 2083236893
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// regTestGenesisHash is the hash of the first block in the block chain for the
// regression test network (genesis block).
var regTestGenesisHash = regTestGenesisBlock.BlockHash()

// regTestGenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the regression test network.  It is the same as the merkle root for
// the main network.
var regTestGenesisMerkleRoot = genesisMerkleRoot

// regTestGenesisBlock defines the genesis block of the block chain which serves
// as the public transaction ledger for the regression test network.
var regTestGenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},         // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: regTestGenesisMerkleRoot, // 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
		Timestamp:  time.Unix(1296688602, 0), // 2011-02-02 23:16:42 +0000 UTC
		Bits:       0x207fffff,               // 545259519 [7fffff0000000000000000000000000000000000000000000000000000000000]
		Nonce:      2,
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// testNet3GenesisHash is the hash of the first block in the block chain for the
// test network (version 3).
var testNet3GenesisHash = testNet3GenesisBlock.BlockHash()

// testNet3GenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the test network (version 3).  It is the same as the merkle root
// for the main network.
var testNet3GenesisMerkleRoot = genesisMerkleRoot

// testNet3GenesisBlock defines the genesis block of the block chain which
// serves as the public transaction ledger for the test network (version 3).
var testNet3GenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},          // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: testNet3GenesisMerkleRoot, // 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
		Timestamp:  time.Unix(1296688602, 0),  // 2011-02-02 23:16:42 +0000 UTC
		Bits:       0x1d00ffff,                // 486604799 [00000000ffff0000000000000000000000000000000000000000000000000000]
		Nonce:      0x18aea41a,                // 414098458
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// simNetGenesisHash is the hash of the first block in the block chain for the
// simulation test network.
var simNetGenesisHash = simNetGenesisBlock.BlockHash()

// simNetGenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the simulation test network.  It is the same as the merkle root for
// the main network.
var simNetGenesisMerkleRoot = genesisMerkleRoot

// simNetGenesisCoinbaseTx is the coinbase transaction for the genesis blocks for
// the sim net.
var simNetGenesisCoinbaseTx = wire.MsgTx{
	Version: 1,
	TxIn: []*wire.TxIn{
		{
			PreviousOutPoint: wire.OutPoint{
				Hash:  chainhash.Hash{},
				Index: 0xffffffff,
			},
			SignatureScript: []byte{
				0x59, 0x6f, 0x75, 0x20, 0x6f, 0x66, 0x74, 0x65,
				0x6e, 0x20, 0x66, 0x65, 0x65, 0x6c, 0x20, 0x74,
				0x69, 0x72, 0x65, 0x64, 0x2c, 0x20, 0x6e, 0x6f,
				0x74, 0x20, 0x62, 0x65, 0x63, 0x61, 0x75, 0x73,
				0x65, 0x20, 0x79, 0x6f, 0x75, 0x27, 0x76, 0x65,
				0x20, 0x64, 0x6f, 0x6e, 0x65, 0x20, 0x74, 0x6f,
				0x6f, 0x20, 0x6d, 0x75, 0x63, 0x68, 0x2c, 0x20,
				0x62, 0x75, 0x74, 0x20, 0x62, 0x65, 0x63, 0x61,
				0x75, 0x73, 0x65, 0x20, 0x79, 0x6f, 0x75, 0x27,
				0x76, 0x65, 0x20, 0x64, 0x6f, 0x6e, 0x65, 0x20,
				0x74, 0x6f, 0x6f, 0x20, 0x6c, 0x69, 0x74, 0x74,
				0x6c, 0x65, 0x20, 0x6f, 0x66, 0x20, 0x77, 0x68,
				0x61, 0x74, 0x20, 0x73, 0x70, 0x61, 0x72, 0x6b,
				0x73, 0x20, 0x61, 0x20, 0x6c, 0x69, 0x67, 0x68,
				0x74, 0x20, 0x69, 0x6e, 0x20, 0x79, 0x6f, 0x75,
				0x2e,
			},
			Sequence: 0xffffffff,
		},
	},
	TxOut: []*wire.TxOut{
		{
			Value: PreminedSTB,
			PkScript: []byte{
				0x76,                                           // OP_DUP
				0xA9,                                           // OP_HASH160,
				0x14,                                           // bytes to push
				0xEC, 0xCA, 0xD9, 0xB4, 0x1F, 0x2B, 0xC2, 0x40, // PKH
				0x70, 0xF9, 0xDE, 0xC1, 0x7F, 0xD9, 0xAC, 0x0B,
				0x0D, 0x1D, 0xBC, 0xE9,
				0x88, // OP_EQUALVERIFY,
				0xAC, // OP_CHECKSIG,
			},
		},
		{
			Value: PreminedNDR,
			PkScript: []byte{
				wire.OP_NDR,
				0x76,                                           // OP_DUP
				0xA9,                                           // OP_HASH160,
				0x14,                                           // bytes to push
				0xEC, 0xCA, 0xD9, 0xB4, 0x1F, 0x2B, 0xC2, 0x40, // PKH
				0x70, 0xF9, 0xDE, 0xC1, 0x7F, 0xD9, 0xAC, 0x0B,
				0x0D, 0x1D, 0xBC, 0xE9,
				0x88, // OP_EQUALVERIFY,
				0xAC, // OP_CHECKSIG,
			},
		},
	},
	LockTime: 0,
}

// simNetGenesisBlock defines the genesis block of the block chain which serves
// as the public transaction ledger for the simulation test network.
var simNetGenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},         // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: simNetGenesisMerkleRoot,  // 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
		Timestamp:  time.Unix(1401292357, 0), // 2014-05-28 15:52:37 +0000 UTC
		Bits:       0x207fffff,               // 545259519 [7fffff0000000000000000000000000000000000000000000000000000000000]
		Nonce:      2,
	},
	Transactions: []*wire.MsgTx{&simNetGenesisCoinbaseTx},
}
