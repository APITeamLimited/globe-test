// Copyright 2017 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"bytes"
	"net"
	"sync"
	"time"
)

// PreparedMessage caches on the wire representations of a message payload.
// Use PreparedMessage to efficiently send a message payload to multiple
// connections. PreparedMessage is especially useful when compression is used
// because the CPU and memory expensive compression operation can be executed
// once for a given set of compression options.
type PreparedMessage struct ***REMOVED***
	messageType int
	data        []byte
	mu          sync.Mutex
	frames      map[prepareKey]*preparedFrame
***REMOVED***

// prepareKey defines a unique set of options to cache prepared frames in PreparedMessage.
type prepareKey struct ***REMOVED***
	isServer         bool
	compress         bool
	compressionLevel int
***REMOVED***

// preparedFrame contains data in wire representation.
type preparedFrame struct ***REMOVED***
	once sync.Once
	data []byte
***REMOVED***

// NewPreparedMessage returns an initialized PreparedMessage. You can then send
// it to connection using WritePreparedMessage method. Valid wire
// representation will be calculated lazily only once for a set of current
// connection options.
func NewPreparedMessage(messageType int, data []byte) (*PreparedMessage, error) ***REMOVED***
	pm := &PreparedMessage***REMOVED***
		messageType: messageType,
		frames:      make(map[prepareKey]*preparedFrame),
		data:        data,
	***REMOVED***

	// Prepare a plain server frame.
	_, frameData, err := pm.frame(prepareKey***REMOVED***isServer: true, compress: false***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// To protect against caller modifying the data argument, remember the data
	// copied to the plain server frame.
	pm.data = frameData[len(frameData)-len(data):]
	return pm, nil
***REMOVED***

func (pm *PreparedMessage) frame(key prepareKey) (int, []byte, error) ***REMOVED***
	pm.mu.Lock()
	frame, ok := pm.frames[key]
	if !ok ***REMOVED***
		frame = &preparedFrame***REMOVED******REMOVED***
		pm.frames[key] = frame
	***REMOVED***
	pm.mu.Unlock()

	var err error
	frame.once.Do(func() ***REMOVED***
		// Prepare a frame using a 'fake' connection.
		// TODO: Refactor code in conn.go to allow more direct construction of
		// the frame.
		mu := make(chan struct***REMOVED******REMOVED***, 1)
		mu <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
		var nc prepareConn
		c := &Conn***REMOVED***
			conn:                   &nc,
			mu:                     mu,
			isServer:               key.isServer,
			compressionLevel:       key.compressionLevel,
			enableWriteCompression: true,
			writeBuf:               make([]byte, defaultWriteBufferSize+maxFrameHeaderSize),
		***REMOVED***
		if key.compress ***REMOVED***
			c.newCompressionWriter = compressNoContextTakeover
		***REMOVED***
		err = c.WriteMessage(pm.messageType, pm.data)
		frame.data = nc.buf.Bytes()
	***REMOVED***)
	return pm.messageType, frame.data, err
***REMOVED***

type prepareConn struct ***REMOVED***
	buf bytes.Buffer
	net.Conn
***REMOVED***

func (pc *prepareConn) Write(p []byte) (int, error)        ***REMOVED*** return pc.buf.Write(p) ***REMOVED***
func (pc *prepareConn) SetWriteDeadline(t time.Time) error ***REMOVED*** return nil ***REMOVED***
