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
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	pb "google.golang.org/grpc/binarylog/grpc_binarylog_v1"
)

var (
	defaultSink Sink = &noopSink***REMOVED******REMOVED*** // TODO(blog): change this default (file in /tmp).
)

// SetDefaultSink sets the sink where binary logs will be written to.
//
// Not thread safe. Only set during initialization.
func SetDefaultSink(s Sink) ***REMOVED***
	if defaultSink != nil ***REMOVED***
		defaultSink.Close()
	***REMOVED***
	defaultSink = s
***REMOVED***

// Sink writes log entry into the binary log sink.
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
func newWriterSink(w io.Writer) *writerSink ***REMOVED***
	return &writerSink***REMOVED***out: w***REMOVED***
***REMOVED***

type writerSink struct ***REMOVED***
	out io.Writer
***REMOVED***

func (ws *writerSink) Write(e *pb.GrpcLogEntry) error ***REMOVED***
	b, err := proto.Marshal(e)
	if err != nil ***REMOVED***
		grpclogLogger.Infof("binary logging: failed to marshal proto message: %v", err)
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

type bufWriteCloserSink struct ***REMOVED***
	mu     sync.Mutex
	closer io.Closer
	out    *writerSink   // out is built on buf.
	buf    *bufio.Writer // buf is kept for flush.

	writeStartOnce sync.Once
	writeTicker    *time.Ticker
***REMOVED***

func (fs *bufWriteCloserSink) Write(e *pb.GrpcLogEntry) error ***REMOVED***
	// Start the write loop when Write is called.
	fs.writeStartOnce.Do(fs.startFlushGoroutine)
	fs.mu.Lock()
	if err := fs.out.Write(e); err != nil ***REMOVED***
		fs.mu.Unlock()
		return err
	***REMOVED***
	fs.mu.Unlock()
	return nil
***REMOVED***

const (
	bufFlushDuration = 60 * time.Second
)

func (fs *bufWriteCloserSink) startFlushGoroutine() ***REMOVED***
	fs.writeTicker = time.NewTicker(bufFlushDuration)
	go func() ***REMOVED***
		for range fs.writeTicker.C ***REMOVED***
			fs.mu.Lock()
			fs.buf.Flush()
			fs.mu.Unlock()
		***REMOVED***
	***REMOVED***()
***REMOVED***

func (fs *bufWriteCloserSink) Close() error ***REMOVED***
	if fs.writeTicker != nil ***REMOVED***
		fs.writeTicker.Stop()
	***REMOVED***
	fs.mu.Lock()
	fs.buf.Flush()
	fs.closer.Close()
	fs.out.Close()
	fs.mu.Unlock()
	return nil
***REMOVED***

func newBufWriteCloserSink(o io.WriteCloser) Sink ***REMOVED***
	bufW := bufio.NewWriter(o)
	return &bufWriteCloserSink***REMOVED***
		closer: o,
		out:    newWriterSink(bufW),
		buf:    bufW,
	***REMOVED***
***REMOVED***

// NewTempFileSink creates a temp file and returns a Sink that writes to this
// file.
func NewTempFileSink() (Sink, error) ***REMOVED***
	tempFile, err := ioutil.TempFile("/tmp", "grpcgo_binarylog_*.txt")
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to create temp file: %v", err)
	***REMOVED***
	return newBufWriteCloserSink(tempFile), nil
***REMOVED***
