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
	"net"
	"strings"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	pb "google.golang.org/grpc/binarylog/grpc_binarylog_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type callIDGenerator struct ***REMOVED***
	id uint64
***REMOVED***

func (g *callIDGenerator) next() uint64 ***REMOVED***
	id := atomic.AddUint64(&g.id, 1)
	return id
***REMOVED***

// reset is for testing only, and doesn't need to be thread safe.
func (g *callIDGenerator) reset() ***REMOVED***
	g.id = 0
***REMOVED***

var idGen callIDGenerator

// MethodLogger is the sub-logger for each method.
type MethodLogger struct ***REMOVED***
	headerMaxLen, messageMaxLen uint64

	callID          uint64
	idWithinCallGen *callIDGenerator

	sink Sink // TODO(blog): make this plugable.
***REMOVED***

func newMethodLogger(h, m uint64) *MethodLogger ***REMOVED***
	return &MethodLogger***REMOVED***
		headerMaxLen:  h,
		messageMaxLen: m,

		callID:          idGen.next(),
		idWithinCallGen: &callIDGenerator***REMOVED******REMOVED***,

		sink: defaultSink, // TODO(blog): make it plugable.
	***REMOVED***
***REMOVED***

// Log creates a proto binary log entry, and logs it to the sink.
func (ml *MethodLogger) Log(c LogEntryConfig) ***REMOVED***
	m := c.toProto()
	timestamp, _ := ptypes.TimestampProto(time.Now())
	m.Timestamp = timestamp
	m.CallId = ml.callID
	m.SequenceIdWithinCall = ml.idWithinCallGen.next()

	switch pay := m.Payload.(type) ***REMOVED***
	case *pb.GrpcLogEntry_ClientHeader:
		m.PayloadTruncated = ml.truncateMetadata(pay.ClientHeader.GetMetadata())
	case *pb.GrpcLogEntry_ServerHeader:
		m.PayloadTruncated = ml.truncateMetadata(pay.ServerHeader.GetMetadata())
	case *pb.GrpcLogEntry_Message:
		m.PayloadTruncated = ml.truncateMessage(pay.Message)
	***REMOVED***

	ml.sink.Write(m)
***REMOVED***

func (ml *MethodLogger) truncateMetadata(mdPb *pb.Metadata) (truncated bool) ***REMOVED***
	if ml.headerMaxLen == maxUInt ***REMOVED***
		return false
	***REMOVED***
	var (
		bytesLimit = ml.headerMaxLen
		index      int
	)
	// At the end of the loop, index will be the first entry where the total
	// size is greater than the limit:
	//
	// len(entry[:index]) <= ml.hdr && len(entry[:index+1]) > ml.hdr.
	for ; index < len(mdPb.Entry); index++ ***REMOVED***
		entry := mdPb.Entry[index]
		if entry.Key == "grpc-trace-bin" ***REMOVED***
			// "grpc-trace-bin" is a special key. It's kept in the log entry,
			// but not counted towards the size limit.
			continue
		***REMOVED***
		currentEntryLen := uint64(len(entry.Value))
		if currentEntryLen > bytesLimit ***REMOVED***
			break
		***REMOVED***
		bytesLimit -= currentEntryLen
	***REMOVED***
	truncated = index < len(mdPb.Entry)
	mdPb.Entry = mdPb.Entry[:index]
	return truncated
***REMOVED***

func (ml *MethodLogger) truncateMessage(msgPb *pb.Message) (truncated bool) ***REMOVED***
	if ml.messageMaxLen == maxUInt ***REMOVED***
		return false
	***REMOVED***
	if ml.messageMaxLen >= uint64(len(msgPb.Data)) ***REMOVED***
		return false
	***REMOVED***
	msgPb.Data = msgPb.Data[:ml.messageMaxLen]
	return true
***REMOVED***

// LogEntryConfig represents the configuration for binary log entry.
type LogEntryConfig interface ***REMOVED***
	toProto() *pb.GrpcLogEntry
***REMOVED***

// ClientHeader configs the binary log entry to be a ClientHeader entry.
type ClientHeader struct ***REMOVED***
	OnClientSide bool
	Header       metadata.MD
	MethodName   string
	Authority    string
	Timeout      time.Duration
	// PeerAddr is required only when it's on server side.
	PeerAddr net.Addr
***REMOVED***

func (c *ClientHeader) toProto() *pb.GrpcLogEntry ***REMOVED***
	// This function doesn't need to set all the fields (e.g. seq ID). The Log
	// function will set the fields when necessary.
	clientHeader := &pb.ClientHeader***REMOVED***
		Metadata:   mdToMetadataProto(c.Header),
		MethodName: c.MethodName,
		Authority:  c.Authority,
	***REMOVED***
	if c.Timeout > 0 ***REMOVED***
		clientHeader.Timeout = ptypes.DurationProto(c.Timeout)
	***REMOVED***
	ret := &pb.GrpcLogEntry***REMOVED***
		Type: pb.GrpcLogEntry_EVENT_TYPE_CLIENT_HEADER,
		Payload: &pb.GrpcLogEntry_ClientHeader***REMOVED***
			ClientHeader: clientHeader,
		***REMOVED***,
	***REMOVED***
	if c.OnClientSide ***REMOVED***
		ret.Logger = pb.GrpcLogEntry_LOGGER_CLIENT
	***REMOVED*** else ***REMOVED***
		ret.Logger = pb.GrpcLogEntry_LOGGER_SERVER
	***REMOVED***
	if c.PeerAddr != nil ***REMOVED***
		ret.Peer = addrToProto(c.PeerAddr)
	***REMOVED***
	return ret
***REMOVED***

// ServerHeader configs the binary log entry to be a ServerHeader entry.
type ServerHeader struct ***REMOVED***
	OnClientSide bool
	Header       metadata.MD
	// PeerAddr is required only when it's on client side.
	PeerAddr net.Addr
***REMOVED***

func (c *ServerHeader) toProto() *pb.GrpcLogEntry ***REMOVED***
	ret := &pb.GrpcLogEntry***REMOVED***
		Type: pb.GrpcLogEntry_EVENT_TYPE_SERVER_HEADER,
		Payload: &pb.GrpcLogEntry_ServerHeader***REMOVED***
			ServerHeader: &pb.ServerHeader***REMOVED***
				Metadata: mdToMetadataProto(c.Header),
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	if c.OnClientSide ***REMOVED***
		ret.Logger = pb.GrpcLogEntry_LOGGER_CLIENT
	***REMOVED*** else ***REMOVED***
		ret.Logger = pb.GrpcLogEntry_LOGGER_SERVER
	***REMOVED***
	if c.PeerAddr != nil ***REMOVED***
		ret.Peer = addrToProto(c.PeerAddr)
	***REMOVED***
	return ret
***REMOVED***

// ClientMessage configs the binary log entry to be a ClientMessage entry.
type ClientMessage struct ***REMOVED***
	OnClientSide bool
	// Message can be a proto.Message or []byte. Other messages formats are not
	// supported.
	Message interface***REMOVED******REMOVED***
***REMOVED***

func (c *ClientMessage) toProto() *pb.GrpcLogEntry ***REMOVED***
	var (
		data []byte
		err  error
	)
	if m, ok := c.Message.(proto.Message); ok ***REMOVED***
		data, err = proto.Marshal(m)
		if err != nil ***REMOVED***
			grpclogLogger.Infof("binarylogging: failed to marshal proto message: %v", err)
		***REMOVED***
	***REMOVED*** else if b, ok := c.Message.([]byte); ok ***REMOVED***
		data = b
	***REMOVED*** else ***REMOVED***
		grpclogLogger.Infof("binarylogging: message to log is neither proto.message nor []byte")
	***REMOVED***
	ret := &pb.GrpcLogEntry***REMOVED***
		Type: pb.GrpcLogEntry_EVENT_TYPE_CLIENT_MESSAGE,
		Payload: &pb.GrpcLogEntry_Message***REMOVED***
			Message: &pb.Message***REMOVED***
				Length: uint32(len(data)),
				Data:   data,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	if c.OnClientSide ***REMOVED***
		ret.Logger = pb.GrpcLogEntry_LOGGER_CLIENT
	***REMOVED*** else ***REMOVED***
		ret.Logger = pb.GrpcLogEntry_LOGGER_SERVER
	***REMOVED***
	return ret
***REMOVED***

// ServerMessage configs the binary log entry to be a ServerMessage entry.
type ServerMessage struct ***REMOVED***
	OnClientSide bool
	// Message can be a proto.Message or []byte. Other messages formats are not
	// supported.
	Message interface***REMOVED******REMOVED***
***REMOVED***

func (c *ServerMessage) toProto() *pb.GrpcLogEntry ***REMOVED***
	var (
		data []byte
		err  error
	)
	if m, ok := c.Message.(proto.Message); ok ***REMOVED***
		data, err = proto.Marshal(m)
		if err != nil ***REMOVED***
			grpclogLogger.Infof("binarylogging: failed to marshal proto message: %v", err)
		***REMOVED***
	***REMOVED*** else if b, ok := c.Message.([]byte); ok ***REMOVED***
		data = b
	***REMOVED*** else ***REMOVED***
		grpclogLogger.Infof("binarylogging: message to log is neither proto.message nor []byte")
	***REMOVED***
	ret := &pb.GrpcLogEntry***REMOVED***
		Type: pb.GrpcLogEntry_EVENT_TYPE_SERVER_MESSAGE,
		Payload: &pb.GrpcLogEntry_Message***REMOVED***
			Message: &pb.Message***REMOVED***
				Length: uint32(len(data)),
				Data:   data,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	if c.OnClientSide ***REMOVED***
		ret.Logger = pb.GrpcLogEntry_LOGGER_CLIENT
	***REMOVED*** else ***REMOVED***
		ret.Logger = pb.GrpcLogEntry_LOGGER_SERVER
	***REMOVED***
	return ret
***REMOVED***

// ClientHalfClose configs the binary log entry to be a ClientHalfClose entry.
type ClientHalfClose struct ***REMOVED***
	OnClientSide bool
***REMOVED***

func (c *ClientHalfClose) toProto() *pb.GrpcLogEntry ***REMOVED***
	ret := &pb.GrpcLogEntry***REMOVED***
		Type:    pb.GrpcLogEntry_EVENT_TYPE_CLIENT_HALF_CLOSE,
		Payload: nil, // No payload here.
	***REMOVED***
	if c.OnClientSide ***REMOVED***
		ret.Logger = pb.GrpcLogEntry_LOGGER_CLIENT
	***REMOVED*** else ***REMOVED***
		ret.Logger = pb.GrpcLogEntry_LOGGER_SERVER
	***REMOVED***
	return ret
***REMOVED***

// ServerTrailer configs the binary log entry to be a ServerTrailer entry.
type ServerTrailer struct ***REMOVED***
	OnClientSide bool
	Trailer      metadata.MD
	// Err is the status error.
	Err error
	// PeerAddr is required only when it's on client side and the RPC is trailer
	// only.
	PeerAddr net.Addr
***REMOVED***

func (c *ServerTrailer) toProto() *pb.GrpcLogEntry ***REMOVED***
	st, ok := status.FromError(c.Err)
	if !ok ***REMOVED***
		grpclogLogger.Info("binarylogging: error in trailer is not a status error")
	***REMOVED***
	var (
		detailsBytes []byte
		err          error
	)
	stProto := st.Proto()
	if stProto != nil && len(stProto.Details) != 0 ***REMOVED***
		detailsBytes, err = proto.Marshal(stProto)
		if err != nil ***REMOVED***
			grpclogLogger.Infof("binarylogging: failed to marshal status proto: %v", err)
		***REMOVED***
	***REMOVED***
	ret := &pb.GrpcLogEntry***REMOVED***
		Type: pb.GrpcLogEntry_EVENT_TYPE_SERVER_TRAILER,
		Payload: &pb.GrpcLogEntry_Trailer***REMOVED***
			Trailer: &pb.Trailer***REMOVED***
				Metadata:      mdToMetadataProto(c.Trailer),
				StatusCode:    uint32(st.Code()),
				StatusMessage: st.Message(),
				StatusDetails: detailsBytes,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	if c.OnClientSide ***REMOVED***
		ret.Logger = pb.GrpcLogEntry_LOGGER_CLIENT
	***REMOVED*** else ***REMOVED***
		ret.Logger = pb.GrpcLogEntry_LOGGER_SERVER
	***REMOVED***
	if c.PeerAddr != nil ***REMOVED***
		ret.Peer = addrToProto(c.PeerAddr)
	***REMOVED***
	return ret
***REMOVED***

// Cancel configs the binary log entry to be a Cancel entry.
type Cancel struct ***REMOVED***
	OnClientSide bool
***REMOVED***

func (c *Cancel) toProto() *pb.GrpcLogEntry ***REMOVED***
	ret := &pb.GrpcLogEntry***REMOVED***
		Type:    pb.GrpcLogEntry_EVENT_TYPE_CANCEL,
		Payload: nil,
	***REMOVED***
	if c.OnClientSide ***REMOVED***
		ret.Logger = pb.GrpcLogEntry_LOGGER_CLIENT
	***REMOVED*** else ***REMOVED***
		ret.Logger = pb.GrpcLogEntry_LOGGER_SERVER
	***REMOVED***
	return ret
***REMOVED***

// metadataKeyOmit returns whether the metadata entry with this key should be
// omitted.
func metadataKeyOmit(key string) bool ***REMOVED***
	switch key ***REMOVED***
	case "lb-token", ":path", ":authority", "content-encoding", "content-type", "user-agent", "te":
		return true
	case "grpc-trace-bin": // grpc-trace-bin is special because it's visiable to users.
		return false
	***REMOVED***
	return strings.HasPrefix(key, "grpc-")
***REMOVED***

func mdToMetadataProto(md metadata.MD) *pb.Metadata ***REMOVED***
	ret := &pb.Metadata***REMOVED******REMOVED***
	for k, vv := range md ***REMOVED***
		if metadataKeyOmit(k) ***REMOVED***
			continue
		***REMOVED***
		for _, v := range vv ***REMOVED***
			ret.Entry = append(ret.Entry,
				&pb.MetadataEntry***REMOVED***
					Key:   k,
					Value: []byte(v),
				***REMOVED***,
			)
		***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

func addrToProto(addr net.Addr) *pb.Address ***REMOVED***
	ret := &pb.Address***REMOVED******REMOVED***
	switch a := addr.(type) ***REMOVED***
	case *net.TCPAddr:
		if a.IP.To4() != nil ***REMOVED***
			ret.Type = pb.Address_TYPE_IPV4
		***REMOVED*** else if a.IP.To16() != nil ***REMOVED***
			ret.Type = pb.Address_TYPE_IPV6
		***REMOVED*** else ***REMOVED***
			ret.Type = pb.Address_TYPE_UNKNOWN
			// Do not set address and port fields.
			break
		***REMOVED***
		ret.Address = a.IP.String()
		ret.IpPort = uint32(a.Port)
	case *net.UnixAddr:
		ret.Type = pb.Address_TYPE_UNIX
		ret.Address = a.String()
	default:
		ret.Type = pb.Address_TYPE_UNKNOWN
	***REMOVED***
	return ret
***REMOVED***
