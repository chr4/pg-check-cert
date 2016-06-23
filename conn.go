/*
 * Copyright (c) 2011-2013, 'pq' Contributors Portions Copyright (C) 2011 Blake
 * Mizerany
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"bufio"
	"net"
	"time"
)

type parameterStatus struct {
	// server version in the same format as server_version_num, or 0 if
	// unavailable
	serverVersion int

	// the current location based on the TimeZone value of the session, if
	// available
	currentLocation *time.Location
}

type transactionStatus byte

const (
	txnStatusIdle                transactionStatus = 'I'
	txnStatusIdleInTransaction   transactionStatus = 'T'
	txnStatusInFailedTransaction transactionStatus = 'E'
)

type conn struct {
	c         net.Conn
	buf       *bufio.Reader
	namei     int
	scratch   [512]byte
	txnStatus transactionStatus

	parameterStatus parameterStatus

	saveMessageType   byte
	saveMessageBuffer []byte

	// If true, this connection is bad and all public-facing functions should
	// return ErrBadConn.
	bad bool

	// If set, this connection should never use the binary format when
	// receiving query results from prepared statements.  Only provided for
	// debugging.
	disablePreparedBinaryResult bool

	// Whether to always send []byte parameters over as binary.  Enables single
	// round-trip mode for non-prepared Query calls.
	binaryParameters bool
}

func (c *conn) writeBuf(b byte) *writeBuf {
	c.scratch[0] = b
	return &writeBuf{
		buf: c.scratch[:5],
		pos: 1,
	}
}

func (cn *conn) sendStartupPacket(m *writeBuf) {
	// sanity check
	if m.buf[0] != 0 {
		panic("oops")
	}

	_, err := cn.c.Write((m.wrap())[1:])
	if err != nil {
		panic(err)
	}
}
