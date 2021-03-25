// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"bufio"
	"errors"
	"io"
	"net/rpc"
)

var errRpcJsonNeedsTermWhitespace = errors.New("rpc requires JsonHandle with TermWhitespace=true")

// Rpc provides a rpc Server or Client Codec for rpc communication.
type Rpc interface ***REMOVED***
	ServerCodec(conn io.ReadWriteCloser, h Handle) rpc.ServerCodec
	ClientCodec(conn io.ReadWriteCloser, h Handle) rpc.ClientCodec
***REMOVED***

// RPCOptions holds options specific to rpc functionality
type RPCOptions struct ***REMOVED***
	// RPCNoBuffer configures whether we attempt to buffer reads and writes during RPC calls.
	//
	// Set RPCNoBuffer=true to turn buffering off.
	// Buffering can still be done if buffered connections are passed in, or
	// buffering is configured on the handle.
	RPCNoBuffer bool
***REMOVED***

// rpcCodec defines the struct members and common methods.
type rpcCodec struct ***REMOVED***
	c io.Closer
	r io.Reader
	w io.Writer
	f ioFlusher

	dec *Decoder
	enc *Encoder
	// bw  *bufio.Writer
	// br  *bufio.Reader
	h Handle

	cls atomicClsErr
***REMOVED***

func newRPCCodec(conn io.ReadWriteCloser, h Handle) rpcCodec ***REMOVED***
	// return newRPCCodec2(bufio.NewReader(conn), bufio.NewWriter(conn), conn, h)
	return newRPCCodec2(conn, conn, conn, h)
***REMOVED***

func newRPCCodec2(r io.Reader, w io.Writer, c io.Closer, h Handle) rpcCodec ***REMOVED***
	// defensive: ensure that jsonH has TermWhitespace turned on.
	if jsonH, ok := h.(*JsonHandle); ok && !jsonH.TermWhitespace ***REMOVED***
		panic(errRpcJsonNeedsTermWhitespace)
	***REMOVED***
	// always ensure that we use a flusher, and always flush what was written to the connection.
	// we lose nothing by using a buffered writer internally.
	f, ok := w.(ioFlusher)
	bh := basicHandle(h)
	if !bh.RPCNoBuffer ***REMOVED***
		if bh.WriterBufferSize <= 0 ***REMOVED***
			if !ok ***REMOVED***
				bw := bufio.NewWriter(w)
				f, w = bw, bw
			***REMOVED***
		***REMOVED***
		if bh.ReaderBufferSize <= 0 ***REMOVED***
			if _, ok = w.(ioPeeker); !ok ***REMOVED***
				if _, ok = w.(ioBuffered); !ok ***REMOVED***
					br := bufio.NewReader(r)
					r = br
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return rpcCodec***REMOVED***
		c:   c,
		w:   w,
		r:   r,
		f:   f,
		h:   h,
		enc: NewEncoder(w, h),
		dec: NewDecoder(r, h),
	***REMOVED***
***REMOVED***

func (c *rpcCodec) write(obj1, obj2 interface***REMOVED******REMOVED***, writeObj2 bool) (err error) ***REMOVED***
	if c.c != nil ***REMOVED***
		cls := c.cls.load()
		if cls.closed ***REMOVED***
			return cls.errClosed
		***REMOVED***
	***REMOVED***
	err = c.enc.Encode(obj1)
	if err == nil ***REMOVED***
		if writeObj2 ***REMOVED***
			err = c.enc.Encode(obj2)
		***REMOVED***
	***REMOVED***
	if c.f != nil ***REMOVED***
		if err == nil ***REMOVED***
			err = c.f.Flush()
		***REMOVED*** else ***REMOVED***
			_ = c.f.Flush() // swallow flush error, so we maintain prior error on write
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (c *rpcCodec) swallow(err *error) ***REMOVED***
	defer panicToErr(c.dec, err)
	c.dec.swallow()
***REMOVED***

func (c *rpcCodec) read(obj interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	if c.c != nil ***REMOVED***
		cls := c.cls.load()
		if cls.closed ***REMOVED***
			return cls.errClosed
		***REMOVED***
	***REMOVED***
	//If nil is passed in, we should read and discard
	if obj == nil ***REMOVED***
		// var obj2 interface***REMOVED******REMOVED***
		// return c.dec.Decode(&obj2)
		c.swallow(&err)
		return
	***REMOVED***
	return c.dec.Decode(obj)
***REMOVED***

func (c *rpcCodec) Close() error ***REMOVED***
	if c.c == nil ***REMOVED***
		return nil
	***REMOVED***
	cls := c.cls.load()
	if cls.closed ***REMOVED***
		return cls.errClosed
	***REMOVED***
	cls.errClosed = c.c.Close()
	cls.closed = true
	c.cls.store(cls)
	return cls.errClosed
***REMOVED***

func (c *rpcCodec) ReadResponseBody(body interface***REMOVED******REMOVED***) error ***REMOVED***
	return c.read(body)
***REMOVED***

// -------------------------------------

type goRpcCodec struct ***REMOVED***
	rpcCodec
***REMOVED***

func (c *goRpcCodec) WriteRequest(r *rpc.Request, body interface***REMOVED******REMOVED***) error ***REMOVED***
	return c.write(r, body, true)
***REMOVED***

func (c *goRpcCodec) WriteResponse(r *rpc.Response, body interface***REMOVED******REMOVED***) error ***REMOVED***
	return c.write(r, body, true)
***REMOVED***

func (c *goRpcCodec) ReadResponseHeader(r *rpc.Response) error ***REMOVED***
	return c.read(r)
***REMOVED***

func (c *goRpcCodec) ReadRequestHeader(r *rpc.Request) error ***REMOVED***
	return c.read(r)
***REMOVED***

func (c *goRpcCodec) ReadRequestBody(body interface***REMOVED******REMOVED***) error ***REMOVED***
	return c.read(body)
***REMOVED***

// -------------------------------------

// goRpc is the implementation of Rpc that uses the communication protocol
// as defined in net/rpc package.
type goRpc struct***REMOVED******REMOVED***

// GoRpc implements Rpc using the communication protocol defined in net/rpc package.
//
// Note: network connection (from net.Dial, of type io.ReadWriteCloser) is not buffered.
//
// For performance, you should configure WriterBufferSize and ReaderBufferSize on the handle.
// This ensures we use an adequate buffer during reading and writing.
// If not configured, we will internally initialize and use a buffer during reads and writes.
// This can be turned off via the RPCNoBuffer option on the Handle.
//   var handle codec.JsonHandle
//   handle.RPCNoBuffer = true // turns off attempt by rpc module to initialize a buffer
//
// Example 1: one way of configuring buffering explicitly:
//   var handle codec.JsonHandle // codec handle
//   handle.ReaderBufferSize = 1024
//   handle.WriterBufferSize = 1024
//   var conn io.ReadWriteCloser // connection got from a socket
//   var serverCodec = GoRpc.ServerCodec(conn, handle)
//   var clientCodec = GoRpc.ClientCodec(conn, handle)
//
// Example 2: you can also explicitly create a buffered connection yourself,
// and not worry about configuring the buffer sizes in the Handle.
//   var handle codec.Handle     // codec handle
//   var conn io.ReadWriteCloser // connection got from a socket
//   var bufconn = struct ***REMOVED***      // bufconn here is a buffered io.ReadWriteCloser
//       io.Closer
//       *bufio.Reader
//       *bufio.Writer
//   ***REMOVED******REMOVED***conn, bufio.NewReader(conn), bufio.NewWriter(conn)***REMOVED***
//   var serverCodec = GoRpc.ServerCodec(bufconn, handle)
//   var clientCodec = GoRpc.ClientCodec(bufconn, handle)
//
var GoRpc goRpc

func (x goRpc) ServerCodec(conn io.ReadWriteCloser, h Handle) rpc.ServerCodec ***REMOVED***
	return &goRpcCodec***REMOVED***newRPCCodec(conn, h)***REMOVED***
***REMOVED***

func (x goRpc) ClientCodec(conn io.ReadWriteCloser, h Handle) rpc.ClientCodec ***REMOVED***
	return &goRpcCodec***REMOVED***newRPCCodec(conn, h)***REMOVED***
***REMOVED***
