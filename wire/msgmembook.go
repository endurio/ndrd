// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2018-2019 The Endurio developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"fmt"
	"io"
)

// MsgMemBook implements the Message interface and represents a bitcoin membook
// message.  It is used to request a list of orders still in the active
// order book of a relay.
//
// This message has no payload and was not added until protocol versions
// starting with BIP0035Version.
type MsgMemBook struct{}

// BtcDecode decodes r using the bitcoin protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgMemBook) BtcDecode(r io.Reader, pver uint32, enc MessageEncoding) error {
	if pver < BIP0035Version {
		str := fmt.Sprintf("membook message invalid for protocol "+
			"version %d", pver)
		return messageError("MsgMemBook.BtcDecode", str)
	}

	return nil
}

// BtcEncode encodes the receiver to w using the bitcoin protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgMemBook) BtcEncode(w io.Writer, pver uint32, enc MessageEncoding) error {
	if pver < BIP0035Version {
		str := fmt.Sprintf("membook message invalid for protocol "+
			"version %d", pver)
		return messageError("MsgMemBook.BtcEncode", str)
	}

	return nil
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgMemBook) Command() string {
	return CmdMemBook
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgMemBook) MaxPayloadLength(pver uint32) uint32 {
	return 0
}

// NewMsgMemBook returns a new bitcoin pong message that conforms to the Message
// interface.  See MsgPong for details.
func NewMsgMemBook() *MsgMemBook {
	return &MsgMemBook{}
}
