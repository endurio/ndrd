// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2018-2019 The Endurio developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

const (
	// OdrVersion is the current latest supported order version.
	OdrVersion = 1
)

// MsgOdr implements the Message interface and represents an order message.
// It is used to deliver transaction information in response to a getdata
// message (MsgGetData) for a given order.
type MsgOdr struct {
	*MsgTx
}

// Copy creates a deep copy of a transaction so that the original does not get
// modified when the copy is manipulated.
func (msg *MsgOdr) Copy() *MsgOdr {
	return &MsgOdr{
		MsgTx: msg.MsgTx.Copy(),
	}
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgOdr) Command() string {
	return CmdOdr
}

// NewMsgOdr returns a new bitcoin tx message that conforms to the Message
// interface.  The return instance has a default version of TxVersion and there
// are no transaction inputs or outputs.  Also, the lock time is set to zero
// to indicate the transaction is valid immediately as opposed to some time in
// future.
func NewMsgOdr(version int32) *MsgOdr {
	return &MsgOdr{
		MsgTx: &MsgTx{
			Version: version,
			TxIn:    make([]*TxIn, 0, defaultTxInOutAlloc),
			TxOut:   make([]*TxOut, 0, defaultTxInOutAlloc),
		},
	}
}
