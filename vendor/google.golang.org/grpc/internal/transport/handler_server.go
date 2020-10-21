/*
 *
 * Copyright 2016 gRPC authors.
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

// This file is the implementation of a gRPC server using HTTP/2 which
// uses the standard Go http2 Server implementation (via the
// http.Handler interface), rather than speaking low-level HTTP/2
// frames itself. It is the implementation of *grpc.Server.ServeHTTP.

package transport

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/http2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/internal/grpcutil"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

// NewServerHandlerTransport returns a ServerTransport handling gRPC
// from inside an http.Handler. It requires that the http Server
// supports HTTP/2.
func NewServerHandlerTransport(w http.ResponseWriter, r *http.Request, stats stats.Handler) (ServerTransport, error) ***REMOVED***
	if r.ProtoMajor != 2 ***REMOVED***
		return nil, errors.New("gRPC requires HTTP/2")
	***REMOVED***
	if r.Method != "POST" ***REMOVED***
		return nil, errors.New("invalid gRPC request method")
	***REMOVED***
	contentType := r.Header.Get("Content-Type")
	// TODO: do we assume contentType is lowercase? we did before
	contentSubtype, validContentType := grpcutil.ContentSubtype(contentType)
	if !validContentType ***REMOVED***
		return nil, errors.New("invalid gRPC request content-type")
	***REMOVED***
	if _, ok := w.(http.Flusher); !ok ***REMOVED***
		return nil, errors.New("gRPC requires a ResponseWriter supporting http.Flusher")
	***REMOVED***

	st := &serverHandlerTransport***REMOVED***
		rw:             w,
		req:            r,
		closedCh:       make(chan struct***REMOVED******REMOVED***),
		writes:         make(chan func()),
		contentType:    contentType,
		contentSubtype: contentSubtype,
		stats:          stats,
	***REMOVED***

	if v := r.Header.Get("grpc-timeout"); v != "" ***REMOVED***
		to, err := decodeTimeout(v)
		if err != nil ***REMOVED***
			return nil, status.Errorf(codes.Internal, "malformed time-out: %v", err)
		***REMOVED***
		st.timeoutSet = true
		st.timeout = to
	***REMOVED***

	metakv := []string***REMOVED***"content-type", contentType***REMOVED***
	if r.Host != "" ***REMOVED***
		metakv = append(metakv, ":authority", r.Host)
	***REMOVED***
	for k, vv := range r.Header ***REMOVED***
		k = strings.ToLower(k)
		if isReservedHeader(k) && !isWhitelistedHeader(k) ***REMOVED***
			continue
		***REMOVED***
		for _, v := range vv ***REMOVED***
			v, err := decodeMetadataHeader(k, v)
			if err != nil ***REMOVED***
				return nil, status.Errorf(codes.Internal, "malformed binary metadata: %v", err)
			***REMOVED***
			metakv = append(metakv, k, v)
		***REMOVED***
	***REMOVED***
	st.headerMD = metadata.Pairs(metakv...)

	return st, nil
***REMOVED***

// serverHandlerTransport is an implementation of ServerTransport
// which replies to exactly one gRPC request (exactly one HTTP request),
// using the net/http.Handler interface. This http.Handler is guaranteed
// at this point to be speaking over HTTP/2, so it's able to speak valid
// gRPC.
type serverHandlerTransport struct ***REMOVED***
	rw         http.ResponseWriter
	req        *http.Request
	timeoutSet bool
	timeout    time.Duration

	headerMD metadata.MD

	closeOnce sync.Once
	closedCh  chan struct***REMOVED******REMOVED*** // closed on Close

	// writes is a channel of code to run serialized in the
	// ServeHTTP (HandleStreams) goroutine. The channel is closed
	// when WriteStatus is called.
	writes chan func()

	// block concurrent WriteStatus calls
	// e.g. grpc/(*serverStream).SendMsg/RecvMsg
	writeStatusMu sync.Mutex

	// we just mirror the request content-type
	contentType string
	// we store both contentType and contentSubtype so we don't keep recreating them
	// TODO make sure this is consistent across handler_server and http2_server
	contentSubtype string

	stats stats.Handler
***REMOVED***

func (ht *serverHandlerTransport) Close() error ***REMOVED***
	ht.closeOnce.Do(ht.closeCloseChanOnce)
	return nil
***REMOVED***

func (ht *serverHandlerTransport) closeCloseChanOnce() ***REMOVED*** close(ht.closedCh) ***REMOVED***

func (ht *serverHandlerTransport) RemoteAddr() net.Addr ***REMOVED*** return strAddr(ht.req.RemoteAddr) ***REMOVED***

// strAddr is a net.Addr backed by either a TCP "ip:port" string, or
// the empty string if unknown.
type strAddr string

func (a strAddr) Network() string ***REMOVED***
	if a != "" ***REMOVED***
		// Per the documentation on net/http.Request.RemoteAddr, if this is
		// set, it's set to the IP:port of the peer (hence, TCP):
		// https://golang.org/pkg/net/http/#Request
		//
		// If we want to support Unix sockets later, we can
		// add our own grpc-specific convention within the
		// grpc codebase to set RemoteAddr to a different
		// format, or probably better: we can attach it to the
		// context and use that from serverHandlerTransport.RemoteAddr.
		return "tcp"
	***REMOVED***
	return ""
***REMOVED***

func (a strAddr) String() string ***REMOVED*** return string(a) ***REMOVED***

// do runs fn in the ServeHTTP goroutine.
func (ht *serverHandlerTransport) do(fn func()) error ***REMOVED***
	select ***REMOVED***
	case <-ht.closedCh:
		return ErrConnClosing
	case ht.writes <- fn:
		return nil
	***REMOVED***
***REMOVED***

func (ht *serverHandlerTransport) WriteStatus(s *Stream, st *status.Status) error ***REMOVED***
	ht.writeStatusMu.Lock()
	defer ht.writeStatusMu.Unlock()

	headersWritten := s.updateHeaderSent()
	err := ht.do(func() ***REMOVED***
		if !headersWritten ***REMOVED***
			ht.writePendingHeaders(s)
		***REMOVED***

		// And flush, in case no header or body has been sent yet.
		// This forces a separation of headers and trailers if this is the
		// first call (for example, in end2end tests's TestNoService).
		ht.rw.(http.Flusher).Flush()

		h := ht.rw.Header()
		h.Set("Grpc-Status", fmt.Sprintf("%d", st.Code()))
		if m := st.Message(); m != "" ***REMOVED***
			h.Set("Grpc-Message", encodeGrpcMessage(m))
		***REMOVED***

		if p := st.Proto(); p != nil && len(p.Details) > 0 ***REMOVED***
			stBytes, err := proto.Marshal(p)
			if err != nil ***REMOVED***
				// TODO: return error instead, when callers are able to handle it.
				panic(err)
			***REMOVED***

			h.Set("Grpc-Status-Details-Bin", encodeBinHeader(stBytes))
		***REMOVED***

		if md := s.Trailer(); len(md) > 0 ***REMOVED***
			for k, vv := range md ***REMOVED***
				// Clients don't tolerate reading restricted headers after some non restricted ones were sent.
				if isReservedHeader(k) ***REMOVED***
					continue
				***REMOVED***
				for _, v := range vv ***REMOVED***
					// http2 ResponseWriter mechanism to send undeclared Trailers after
					// the headers have possibly been written.
					h.Add(http2.TrailerPrefix+k, encodeMetadataHeader(k, v))
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***)

	if err == nil ***REMOVED*** // transport has not been closed
		if ht.stats != nil ***REMOVED***
			// Note: The trailer fields are compressed with hpack after this call returns.
			// No WireLength field is set here.
			ht.stats.HandleRPC(s.Context(), &stats.OutTrailer***REMOVED***
				Trailer: s.trailer.Copy(),
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	ht.Close()
	return err
***REMOVED***

// writePendingHeaders sets common and custom headers on the first
// write call (Write, WriteHeader, or WriteStatus)
func (ht *serverHandlerTransport) writePendingHeaders(s *Stream) ***REMOVED***
	ht.writeCommonHeaders(s)
	ht.writeCustomHeaders(s)
***REMOVED***

// writeCommonHeaders sets common headers on the first write
// call (Write, WriteHeader, or WriteStatus).
func (ht *serverHandlerTransport) writeCommonHeaders(s *Stream) ***REMOVED***
	h := ht.rw.Header()
	h["Date"] = nil // suppress Date to make tests happy; TODO: restore
	h.Set("Content-Type", ht.contentType)

	// Predeclare trailers we'll set later in WriteStatus (after the body).
	// This is a SHOULD in the HTTP RFC, and the way you add (known)
	// Trailers per the net/http.ResponseWriter contract.
	// See https://golang.org/pkg/net/http/#ResponseWriter
	// and https://golang.org/pkg/net/http/#example_ResponseWriter_trailers
	h.Add("Trailer", "Grpc-Status")
	h.Add("Trailer", "Grpc-Message")
	h.Add("Trailer", "Grpc-Status-Details-Bin")

	if s.sendCompress != "" ***REMOVED***
		h.Set("Grpc-Encoding", s.sendCompress)
	***REMOVED***
***REMOVED***

// writeCustomHeaders sets custom headers set on the stream via SetHeader
// on the first write call (Write, WriteHeader, or WriteStatus).
func (ht *serverHandlerTransport) writeCustomHeaders(s *Stream) ***REMOVED***
	h := ht.rw.Header()

	s.hdrMu.Lock()
	for k, vv := range s.header ***REMOVED***
		if isReservedHeader(k) ***REMOVED***
			continue
		***REMOVED***
		for _, v := range vv ***REMOVED***
			h.Add(k, encodeMetadataHeader(k, v))
		***REMOVED***
	***REMOVED***

	s.hdrMu.Unlock()
***REMOVED***

func (ht *serverHandlerTransport) Write(s *Stream, hdr []byte, data []byte, opts *Options) error ***REMOVED***
	headersWritten := s.updateHeaderSent()
	return ht.do(func() ***REMOVED***
		if !headersWritten ***REMOVED***
			ht.writePendingHeaders(s)
		***REMOVED***
		ht.rw.Write(hdr)
		ht.rw.Write(data)
		ht.rw.(http.Flusher).Flush()
	***REMOVED***)
***REMOVED***

func (ht *serverHandlerTransport) WriteHeader(s *Stream, md metadata.MD) error ***REMOVED***
	if err := s.SetHeader(md); err != nil ***REMOVED***
		return err
	***REMOVED***

	headersWritten := s.updateHeaderSent()
	err := ht.do(func() ***REMOVED***
		if !headersWritten ***REMOVED***
			ht.writePendingHeaders(s)
		***REMOVED***

		ht.rw.WriteHeader(200)
		ht.rw.(http.Flusher).Flush()
	***REMOVED***)

	if err == nil ***REMOVED***
		if ht.stats != nil ***REMOVED***
			// Note: The header fields are compressed with hpack after this call returns.
			// No WireLength field is set here.
			ht.stats.HandleRPC(s.Context(), &stats.OutHeader***REMOVED***
				Header:      md.Copy(),
				Compression: s.sendCompress,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***

func (ht *serverHandlerTransport) HandleStreams(startStream func(*Stream), traceCtx func(context.Context, string) context.Context) ***REMOVED***
	// With this transport type there will be exactly 1 stream: this HTTP request.

	ctx := ht.req.Context()
	var cancel context.CancelFunc
	if ht.timeoutSet ***REMOVED***
		ctx, cancel = context.WithTimeout(ctx, ht.timeout)
	***REMOVED*** else ***REMOVED***
		ctx, cancel = context.WithCancel(ctx)
	***REMOVED***

	// requestOver is closed when the status has been written via WriteStatus.
	requestOver := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		select ***REMOVED***
		case <-requestOver:
		case <-ht.closedCh:
		case <-ht.req.Context().Done():
		***REMOVED***
		cancel()
		ht.Close()
	***REMOVED***()

	req := ht.req

	s := &Stream***REMOVED***
		id:             0, // irrelevant
		requestRead:    func(int) ***REMOVED******REMOVED***,
		cancel:         cancel,
		buf:            newRecvBuffer(),
		st:             ht,
		method:         req.URL.Path,
		recvCompress:   req.Header.Get("grpc-encoding"),
		contentSubtype: ht.contentSubtype,
	***REMOVED***
	pr := &peer.Peer***REMOVED***
		Addr: ht.RemoteAddr(),
	***REMOVED***
	if req.TLS != nil ***REMOVED***
		pr.AuthInfo = credentials.TLSInfo***REMOVED***State: *req.TLS, CommonAuthInfo: credentials.CommonAuthInfo***REMOVED***SecurityLevel: credentials.PrivacyAndIntegrity***REMOVED******REMOVED***
	***REMOVED***
	ctx = metadata.NewIncomingContext(ctx, ht.headerMD)
	s.ctx = peer.NewContext(ctx, pr)
	if ht.stats != nil ***REMOVED***
		s.ctx = ht.stats.TagRPC(s.ctx, &stats.RPCTagInfo***REMOVED***FullMethodName: s.method***REMOVED***)
		inHeader := &stats.InHeader***REMOVED***
			FullMethod:  s.method,
			RemoteAddr:  ht.RemoteAddr(),
			Compression: s.recvCompress,
		***REMOVED***
		ht.stats.HandleRPC(s.ctx, inHeader)
	***REMOVED***
	s.trReader = &transportReader***REMOVED***
		reader:        &recvBufferReader***REMOVED***ctx: s.ctx, ctxDone: s.ctx.Done(), recv: s.buf, freeBuffer: func(*bytes.Buffer) ***REMOVED******REMOVED******REMOVED***,
		windowHandler: func(int) ***REMOVED******REMOVED***,
	***REMOVED***

	// readerDone is closed when the Body.Read-ing goroutine exits.
	readerDone := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		defer close(readerDone)

		// TODO: minimize garbage, optimize recvBuffer code/ownership
		const readSize = 8196
		for buf := make([]byte, readSize); ; ***REMOVED***
			n, err := req.Body.Read(buf)
			if n > 0 ***REMOVED***
				s.buf.put(recvMsg***REMOVED***buffer: bytes.NewBuffer(buf[:n:n])***REMOVED***)
				buf = buf[n:]
			***REMOVED***
			if err != nil ***REMOVED***
				s.buf.put(recvMsg***REMOVED***err: mapRecvMsgError(err)***REMOVED***)
				return
			***REMOVED***
			if len(buf) == 0 ***REMOVED***
				buf = make([]byte, readSize)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	// startStream is provided by the *grpc.Server's serveStreams.
	// It starts a goroutine serving s and exits immediately.
	// The goroutine that is started is the one that then calls
	// into ht, calling WriteHeader, Write, WriteStatus, Close, etc.
	startStream(s)

	ht.runStream()
	close(requestOver)

	// Wait for reading goroutine to finish.
	req.Body.Close()
	<-readerDone
***REMOVED***

func (ht *serverHandlerTransport) runStream() ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case fn := <-ht.writes:
			fn()
		case <-ht.closedCh:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (ht *serverHandlerTransport) IncrMsgSent() ***REMOVED******REMOVED***

func (ht *serverHandlerTransport) IncrMsgRecv() ***REMOVED******REMOVED***

func (ht *serverHandlerTransport) Drain() ***REMOVED***
	panic("Drain() is not implemented")
***REMOVED***

// mapRecvMsgError returns the non-nil err into the appropriate
// error value as expected by callers of *grpc.parser.recvMsg.
// In particular, in can only be:
//   * io.EOF
//   * io.ErrUnexpectedEOF
//   * of type transport.ConnectionError
//   * an error from the status package
func mapRecvMsgError(err error) error ***REMOVED***
	if err == io.EOF || err == io.ErrUnexpectedEOF ***REMOVED***
		return err
	***REMOVED***
	if se, ok := err.(http2.StreamError); ok ***REMOVED***
		if code, ok := http2ErrConvTab[se.Code]; ok ***REMOVED***
			return status.Error(code, se.Error())
		***REMOVED***
	***REMOVED***
	if strings.Contains(err.Error(), "body closed by handler") ***REMOVED***
		return status.Error(codes.Canceled, err.Error())
	***REMOVED***
	return connectionErrorf(true, err, err.Error())
***REMOVED***
