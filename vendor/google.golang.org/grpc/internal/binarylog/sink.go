/*
 *
 * Copyright 2018 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package binarylog

import (
	"bufio"
	"encoding/binary"
	"io"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	pb "google.golang.org/grpc/binarylog/grpc_binarylog_v1"
)

var (
	// DefaultSink is the sink where the logs will be written to. It's exported
	// for the binarylog package to update.
	DefaultSink Sink = &noopSink***REMOVED******REMOVED*** // TODO(blog): change this default (file in /tmp).
)

// Sink writes log entry into the binary log sink.
//
// sink is a copy of the exported binarylog.Sink, to avoid circular dependency.
type Sink interface ***REMOVED***
	// Write will be called to write the log entry into the sink.
	//
	// It should be thread-safe so it can be called in parallel.
	Write(*pb.GrpcLogEntry) error
	// Close will be called when the Sink is replaced by a new Sink.
	Close() error
***REMOVED***

type noopSink struct***REMOVED******REMOVED***

func (ns *noopSink) Write(*pb.GrpcLogEntry) error ***REMOVED*** return nil ***REMOVED***
func (ns *noopSink) Close() error                 ***REMOVED*** return nil ***REMOVED***

// newWriterSink creates a binary log sink with the given writer.
//
// Write() marshals the proto message and writes it to the given writer. Each
// message is prefixed with a 4 byte big endian unsigned integer as the length.
//
// No buffer is done, Close() doesn't try to close the writer.
func newWriterSink(w io.Writer) Sink ***REMOVED***
	return &writerSink***REMOVED***out: w***REMOVED***
***REMOVED***

type writerSink struct ***REMOVED***
	out io.Writer
***REMOVED***

func (ws *writerSink) Write(e *pb.GrpcLogEntry) error ***REMOVED***
	b, err := proto.Marshal(e)
	if err != nil ***REMOVED***
		grpclogLogger.Errorf("binary logging: failed to marshal proto message: %v", err)
		return err
	***REMOVED***
	hdr := make([]byte, 4)
	binary.BigEndian.PutUint32(hdr, uint32(len(b)))
	if _, err := ws.out.Write(hdr); err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err := ws.out.Write(b); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func (ws *writerSink) Close() error ***REMOVED*** return nil ***REMOVED***

type bufferedSink struct ***REMOVED***
	mu             sync.Mutex
	closer         io.Closer
	out            Sink          // out is built on buf.
	buf            *bufio.Writer // buf is kept for flush.
	flusherStarted bool

	writeTicker *time.Ticker
	done        chan struct***REMOVED******REMOVED***
***REMOVED***

func (fs *bufferedSink) Write(e *pb.GrpcLogEntry) error ***REMOVED***
	fs.mu.Lock()
	defer fs.mu.Unlock()
	if !fs.flusherStarted ***REMOVED***
		// Start the write loop when Write is called.
		fs.startFlushGoroutine()
		fs.flusherStarted = true
	***REMOVED***
	if err := fs.out.Write(e); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

const (
	bufFlushDuration = 60 * time.Second
)

func (fs *bufferedSink) startFlushGoroutine() ***REMOVED***
	fs.writeTicker = time.NewTicker(bufFlushDuration)
	go func() ***REMOVED***
		for ***REMOVED***
			select ***REMOVED***
			case <-fs.done:
				return
			case <-fs.writeTicker.C:
			***REMOVED***
			fs.mu.Lock()
			if err := fs.buf.Flush(); err != nil ***REMOVED***
				grpclogLogger.Warningf("failed to flush to Sink: %v", err)
			***REMOVED***
			fs.mu.Unlock()
		***REMOVED***
	***REMOVED***()
***REMOVED***

func (fs *bufferedSink) Close() error ***REMOVED***
	fs.mu.Lock()
	defer fs.mu.Unlock()
	if fs.writeTicker != nil ***REMOVED***
		fs.writeTicker.Stop()
	***REMOVED***
	close(fs.done)
	if err := fs.buf.Flush(); err != nil ***REMOVED***
		grpclogLogger.Warningf("failed to flush to Sink: %v", err)
	***REMOVED***
	if err := fs.closer.Close(); err != nil ***REMOVED***
		grpclogLogger.Warningf("failed to close the underlying WriterCloser: %v", err)
	***REMOVED***
	if err := fs.out.Close(); err != nil ***REMOVED***
		grpclogLogger.Warningf("failed to close the Sink: %v", err)
	***REMOVED***
	return nil
***REMOVED***

// NewBufferedSink creates a binary log sink with the given WriteCloser.
//
// Write() marshals the proto message and writes it to the given writer. Each
// message is prefixed with a 4 byte big endian unsigned integer as the length.
//
// Content is kept in a buffer, and is flushed every 60 seconds.
//
// Close closes the WriteCloser.
func NewBufferedSink(o io.WriteCloser) Sink ***REMOVED***
	bufW := bufio.NewWriter(o)
	return &bufferedSink***REMOVED***
		closer: o,
		out:    newWriterSink(bufW),
		buf:    bufW,
		done:   make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
***REMOVED***
