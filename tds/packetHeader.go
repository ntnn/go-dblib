// SPDX-FileCopyrightText: 2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// Default size of a packet-header.
const (
	PacketHeaderSize = 8
)

//go:generate stringer -type=PacketHeaderType

// PacketHeaderType is the type for bitmask values of a packet-header
// type.
type PacketHeaderType uint8

// Types of packet-header.
const (
	TDS_BUF_LANG PacketHeaderType = iota + 1
	TDS_BUF_LOGIN
	TDS_BUF_RPC
	TDS_BUF_RESPONSE
	TDS_BUF_UNFMT
	TDS_BUF_ATTN
	TDS_BUF_BULK
	TDS_BUF_SETUP
	TDS_BUF_CLOSE
	TDS_BUF_ERROR
	TDS_BUF_PROTACK
	TDS_BUF_ECHO
	TDS_BUF_LOGOUT
	TDS_BUF_ENDPARAM
	TDS_BUF_NORMAL
	TDS_BUF_URGENT
	TDS_BUF_MIGRATE
	TDS_BUF_HELLO
	TDS_BUF_CMDSEQ_NORMAL
	TDS_BUF_CMDSEQ_LOGIN
	TDS_BUF_CMDSEQ_LIVENESS
	TDS_BUF_CMDSEQ_RESERVED1
	TDS_BUF_CMDSEQ_RESERVED2
)

//go:generate stringer -type=PacketHeaderStatus

// PacketHeaderStatus is the type for bitmask values of a packet-header
// status.
type PacketHeaderStatus uint8

// States of packet-header.
const (
	// Last buffer in a request or response
	TDS_BUFSTAT_EOM PacketHeaderStatus = 0x1
	// Acknowledgment of last receiver attention
	TDS_BUFSTAT_ATTNACK PacketHeaderStatus = 0x2
	// Attention request
	TDS_BUFSTAT_ATTN PacketHeaderStatus = 0x4
	// Event notification
	TDS_BUFSTAT_EVENT PacketHeaderStatus = 0x8
	// Buffer is encrypted
	TDS_BUFSTAT_SEAL PacketHeaderStatus = 0x10
	// Buffer is encrypted (SQL Anywhere CMDSQ protocol)
	TDS_BUFSTAT_ENCRYPT PacketHeaderStatus = 0x20
	// Buffer is encrypted with symmetric key for on demand command
	// encryption
	TDS_BUFSTAT_SYMENCRYPT PacketHeaderStatus = 0x40
)

//PacketHeader contains information about the packet-header.
type PacketHeader struct {
	// Message type, e.g. for login or language command
	MsgType PacketHeaderType
	// Status, e.g. encrypted or EOM
	Status PacketHeaderStatus
	// Length of package in bytes
	Length uint16
	// Channel the packet belongs to when multiplexing
	Channel uint16
	// PacketNr for ordering when multiplexing
	PacketNr uint8
	// Allowed window size before ACK is received
	Window uint8
}

// NewPacketHeader creates a packet-header struct.
func NewPacketHeader(packetSize int) PacketHeader {
	return PacketHeader{
		Length: uint16(packetSize),
	}
}

func (header PacketHeader) String() string {
	return fmt.Sprintf(
		"MsgType: %s, Status: %s, Length: %d, Channel: %d, PacketNr: %d, Window: %d",
		header.MsgType, header.Status, header.Length, header.Channel, header.PacketNr, header.Window,
	)
}

// WriteTo writes a byte slice and returns the amount of bytes written.
func (header PacketHeader) WriteTo(w io.Writer) (int64, error) {
	bs := make([]byte, PacketHeaderSize)
	n, err := header.Read(bs)
	if err != nil || n != PacketHeaderSize {
		return n, fmt.Errorf("failed to write header information to byte slice: %w", err)
	}

	m, err := w.Write(bs)
	return int64(m), err
}

// Read satisfies the io.Reader interface.
func (header PacketHeader) Read(bs []byte) (int64, error) {
	if len(bs) != PacketHeaderSize {
		return 0, fmt.Errorf("target buffer has unexpected length, expected %d bytes, buffer length is %d",
			PacketHeaderSize, len(bs))
	}

	bs[0] = byte(header.MsgType)
	bs[1] = byte(header.Status)
	// header is always big endian
	binary.BigEndian.PutUint16(bs[2:4], header.Length)
	binary.BigEndian.PutUint16(bs[4:6], header.Channel)
	bs[6] = byte(header.PacketNr)
	bs[7] = byte(header.Window)
	return PacketHeaderSize, nil
}

// ReadFrom reads a byte slice, writes it to the header and returns the
// amount of written bytes.
func (header *PacketHeader) ReadFrom(r io.Reader) (int64, error) {
	bs := make([]byte, PacketHeaderSize)
	n, err := r.Read(bs)
	if err != nil || n != PacketHeaderSize {
		if n == 0 && errors.Is(err, io.EOF) {
			return 0, ErrEOFAfterZeroRead
		}
		return int64(n), fmt.Errorf("read %d of %d expected bytes from reader: %w", n, PacketHeaderSize, err)
	}

	m, err := header.Write(bs)
	return int64(m), err
}

// Write satisfies the io.Writer interface.
func (header *PacketHeader) Write(bs []byte) (int64, error) {
	if len(bs) != PacketHeaderSize {
		return 0, fmt.Errorf("passed buffer has unexpected length, expected 8 bytes, buffer length is %d", len(bs))
	}

	header.MsgType = PacketHeaderType(bs[0])
	header.Status = PacketHeaderStatus(bs[1])
	uvarint := binary.BigEndian.Uint16(bs[2:4])
	header.Length = uint16(uvarint)
	uvarint = binary.BigEndian.Uint16(bs[4:6])
	header.Channel = uint16(uvarint)
	header.PacketNr = uint8(bs[6])
	header.Window = uint8(bs[7])

	return PacketHeaderSize, nil
}
