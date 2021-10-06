/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2020 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	grpcstats "google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"

	//nolint: staticcheck
	protoV1 "github.com/golang/protobuf/proto"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/stats"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

//nolint: lll
var (
	errInvokeRPCInInitContext = common.NewInitContextError("invoking RPC methods in the init context is not supported")
	errConnectInInitContext   = common.NewInitContextError("connecting to a gRPC server in the init context is not supported")
)

// Client represents a gRPC client that can be used to make RPC requests
type Client struct ***REMOVED***
	mds         map[string]protoreflect.MethodDescriptor
	conn        *grpc.ClientConn
	reflectOnce sync.Once
***REMOVED***

// XClient represents the Client constructor (e.g. `new grpc.Client()`) and
// creates a new gPRC client object that can load protobuf definitions, connect
// to servers and invoke RPC methods.
func (*GRPC) XClient(ctxPtr *context.Context) interface***REMOVED******REMOVED*** ***REMOVED***
	rt := common.GetRuntime(*ctxPtr)

	return common.Bind(rt, &Client***REMOVED******REMOVED***, ctxPtr)
***REMOVED***

// MethodInfo holds information on any parsed method descriptors that can be used by the goja VM
type MethodInfo struct ***REMOVED***
	grpc.MethodInfo `json:"-" js:"-"`
	Package         string
	Service         string
	FullMethod      string
***REMOVED***

// Response is a gRPC response that can be used by the goja VM
type Response struct ***REMOVED***
	Status   codes.Code
	Message  interface***REMOVED******REMOVED***
	Headers  map[string][]string
	Trailers map[string][]string
	Error    interface***REMOVED******REMOVED***
***REMOVED***

func walkFileDescriptors(seen map[string]struct***REMOVED******REMOVED***, fd *desc.FileDescriptor) []*descriptorpb.FileDescriptorProto ***REMOVED***
	fds := []*descriptorpb.FileDescriptorProto***REMOVED******REMOVED***

	if _, ok := seen[fd.GetName()]; ok ***REMOVED***
		return fds
	***REMOVED***
	seen[fd.GetName()] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	fds = append(fds, fd.AsFileDescriptorProto())

	for _, dep := range fd.GetDependencies() ***REMOVED***
		deps := walkFileDescriptors(seen, dep)
		fds = append(fds, deps...)
	***REMOVED***

	return fds
***REMOVED***

// Load will parse the given proto files and make the file descriptors available to request.
func (c *Client) Load(ctxPtr *context.Context, importPaths []string, filenames ...string) ([]MethodInfo, error) ***REMOVED***
	if lib.GetState(*ctxPtr) != nil ***REMOVED***
		return nil, errors.New("load must be called in the init context")
	***REMOVED***

	initEnv := common.GetInitEnv(*ctxPtr)
	if initEnv == nil ***REMOVED***
		return nil, errors.New("missing init environment")
	***REMOVED***

	// If no import paths are specified, use the current working directory
	if len(importPaths) == 0 ***REMOVED***
		importPaths = append(importPaths, initEnv.CWD.Path)
	***REMOVED***

	parser := protoparse.Parser***REMOVED***
		ImportPaths:      importPaths,
		InferImportPaths: false,
		Accessor: protoparse.FileAccessor(func(filename string) (io.ReadCloser, error) ***REMOVED***
			absFilePath := initEnv.GetAbsFilePath(filename)
			return initEnv.FileSystems["file"].Open(absFilePath)
		***REMOVED***),
	***REMOVED***

	fds, err := parser.ParseFiles(filenames...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	fdset := &descriptorpb.FileDescriptorSet***REMOVED******REMOVED***

	seen := make(map[string]struct***REMOVED******REMOVED***)
	for _, fd := range fds ***REMOVED***
		fdset.File = append(fdset.File, walkFileDescriptors(seen, fd)...)
	***REMOVED***
	return c.convertToMethodInfo(fdset)
***REMOVED***

func (c *Client) convertToMethodInfo(fdset *descriptorpb.FileDescriptorSet) ([]MethodInfo, error) ***REMOVED***
	files, err := protodesc.NewFiles(fdset)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var rtn []MethodInfo
	if c.mds == nil ***REMOVED***
		// This allows us to call load() multiple times, without overwriting the
		// previously loaded definitions.
		c.mds = make(map[string]protoreflect.MethodDescriptor)
	***REMOVED***
	appendMethodInfo := func(
		fd protoreflect.FileDescriptor,
		sd protoreflect.ServiceDescriptor,
		md protoreflect.MethodDescriptor,
	) ***REMOVED***
		name := fmt.Sprintf("/%s/%s", sd.FullName(), md.Name())
		c.mds[name] = md
		rtn = append(rtn, MethodInfo***REMOVED***
			MethodInfo: grpc.MethodInfo***REMOVED***
				Name:           string(md.Name()),
				IsClientStream: md.IsStreamingClient(),
				IsServerStream: md.IsStreamingServer(),
			***REMOVED***,
			Package:    string(fd.Package()),
			Service:    string(sd.Name()),
			FullMethod: name,
		***REMOVED***)
	***REMOVED***
	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool ***REMOVED***
		sds := fd.Services()
		for i := 0; i < sds.Len(); i++ ***REMOVED***
			sd := sds.Get(i)
			mds := sd.Methods()
			for j := 0; j < mds.Len(); j++ ***REMOVED***
				md := mds.Get(j)
				appendMethodInfo(fd, sd, md)
			***REMOVED***
		***REMOVED***
		return true
	***REMOVED***)
	return rtn, nil
***REMOVED***

type transportCreds struct ***REMOVED***
	credentials.TransportCredentials
	errc chan<- error
***REMOVED***

func (t transportCreds) ClientHandshake(ctx context.Context,
	addr string, in net.Conn) (net.Conn, credentials.AuthInfo, error) ***REMOVED***
	out, auth, err := t.TransportCredentials.ClientHandshake(ctx, addr, in)
	if err != nil ***REMOVED***
		t.errc <- err
	***REMOVED***

	return out, auth, err
***REMOVED***

// Connect is a block dial to the gRPC server at the given address (host:port)
// nolint: funlen
func (c *Client) Connect(ctxPtr *context.Context, addr string, params map[string]interface***REMOVED******REMOVED***) (bool, error) ***REMOVED***
	state := lib.GetState(*ctxPtr)
	if state == nil ***REMOVED***
		return false, errConnectInInitContext
	***REMOVED***
	isPlaintext, reflect, timeout := false, false, 60*time.Second

	for k, v := range params ***REMOVED***
		switch k ***REMOVED***
		case "plaintext":
			isPlaintext, _ = v.(bool)
		case "timeout":
			var err error
			timeout, err = types.GetDurationValue(v)
			if err != nil ***REMOVED***
				return false, fmt.Errorf("invalid timeout value: %w", err)
			***REMOVED***
		case "reflect":
			var ok bool
			reflect, ok = v.(bool)
			if !ok ***REMOVED***
				return false, fmt.Errorf("invalid reflect value: '%#v', it needs to be boolean", v)
			***REMOVED***

		default:
			return false, fmt.Errorf("unknown connect param: %q", k)
		***REMOVED***
	***REMOVED***

	// (rogchap) Even with FailOnNonTempDialError, if there is a TLS error this will timeout
	// rather than report the error, so we can't rely on WithBlock. By running in a goroutine
	// we can then wait on the error channel instead, which could happen before the Dial
	// returns. We only need to close the channel to un-block in a non-error scenario;
	// otherwise it can be GCd without closing as we return on an error on the channel.
	errc := make(chan error, 1)
	go func() ***REMOVED***
		opts := []grpc.DialOption***REMOVED***
			grpc.WithBlock(),
			grpc.FailOnNonTempDialError(true),
			grpc.WithStatsHandler(c),
		***REMOVED***

		if ua := state.Options.UserAgent; ua.Valid ***REMOVED***
			opts = append(opts, grpc.WithUserAgent(ua.ValueOrZero()))
		***REMOVED***

		if !isPlaintext ***REMOVED***
			tlsCfg := state.TLSConfig.Clone()
			tlsCfg.NextProtos = []string***REMOVED***"h2"***REMOVED***

			// TODO(rogchap): Would be good to add support for custom RootCAs (self signed)

			// (rogchap) we create a wrapper for transport credentials so that we can report
			// on any TLS errors.
			creds := transportCreds***REMOVED***
				credentials.NewTLS(tlsCfg),
				errc,
			***REMOVED***
			opts = append(opts, grpc.WithTransportCredentials(creds))
		***REMOVED***

		if isPlaintext ***REMOVED***
			opts = append(opts, grpc.WithInsecure())
		***REMOVED***

		dialer := func(ctx context.Context, addr string) (net.Conn, error) ***REMOVED***
			return state.Dialer.DialContext(ctx, "tcp", addr)
		***REMOVED***
		opts = append(opts, grpc.WithContextDialer(dialer))

		ctx, cancel := context.WithTimeout(*ctxPtr, timeout)
		defer cancel()

		var err error
		c.conn, err = grpc.DialContext(ctx, addr, opts...)
		if err != nil ***REMOVED***
			errc <- err

			return
		***REMOVED***
		close(errc)
	***REMOVED***()
	if err := <-errc; err != nil ***REMOVED***
		return false, err
	***REMOVED***
	if reflect ***REMOVED***
		var err error
		c.reflectOnce.Do(func() ***REMOVED***
			err = c.reflect(ctxPtr)
		***REMOVED***)
		if err != nil ***REMOVED***
			return false, fmt.Errorf("reflect: %w", err)
		***REMOVED***
	***REMOVED***
	return true, nil
***REMOVED***

// reflect will use the grpc reflection api to make the file descriptors available to request.
// It is called in the connect function the first time the Client.Connect function is called.
func (c *Client) reflect(ctxPtr *context.Context) error ***REMOVED***
	client := reflectpb.NewServerReflectionClient(c.conn)
	methodClient, err := client.ServerReflectionInfo(*ctxPtr)
	if err != nil ***REMOVED***
		return fmt.Errorf("can't get server info: %w", err)
	***REMOVED***
	req := &reflectpb.ServerReflectionRequest***REMOVED***
		MessageRequest: &reflectpb.ServerReflectionRequest_ListServices***REMOVED******REMOVED***,
	***REMOVED***
	if err = methodClient.Send(req); err != nil ***REMOVED***
		return fmt.Errorf("can't send list services request: %w", err)
	***REMOVED***
	resp, err := methodClient.Recv()
	if err != nil ***REMOVED***
		return fmt.Errorf("please enable server reflection on your grpc server: %w", err)
	***REMOVED***
	listResp := resp.GetListServicesResponse()
	if listResp == nil ***REMOVED***
		return fmt.Errorf("can't list services, nil response")
	***REMOVED***
	fdset := &descriptorpb.FileDescriptorSet***REMOVED******REMOVED***
	for _, service := range listResp.GetService() ***REMOVED***
		req = &reflectpb.ServerReflectionRequest***REMOVED***
			MessageRequest: &reflectpb.ServerReflectionRequest_FileContainingSymbol***REMOVED***
				FileContainingSymbol: service.GetName(),
			***REMOVED***,
		***REMOVED***
		if err = methodClient.Send(req); err != nil ***REMOVED***
			return fmt.Errorf("can't send file descriptors request on service %q: %w", service, err)
		***REMOVED***
		resp, err = methodClient.Recv()
		if err != nil ***REMOVED***
			return fmt.Errorf("can't receive file descriptors request on service %q: %w", service, err)
		***REMOVED***
		fdResp := resp.GetFileDescriptorResponse()
		for _, raw := range fdResp.GetFileDescriptorProto() ***REMOVED***
			var fdp descriptorpb.FileDescriptorProto
			if err = proto.Unmarshal(raw, &fdp); err != nil ***REMOVED***
				return fmt.Errorf("can't unmarshal proto on service %q: %w", service, err)
			***REMOVED***
			fdset.File = append(fdset.File, &fdp)
		***REMOVED***
	***REMOVED***
	_, err = c.convertToMethodInfo(fdset)
	if err != nil ***REMOVED***
		err = fmt.Errorf("can't convert method info: %w", err)
	***REMOVED***
	return err
***REMOVED***

// Invoke creates and calls a unary RPC by fully qualified method name
//nolint: funlen,gocognit,gocyclo
func (c *Client) Invoke(ctxPtr *context.Context,
	method string, req goja.Value, params map[string]interface***REMOVED******REMOVED***) (*Response, error) ***REMOVED***
	ctx := *ctxPtr
	rt := common.GetRuntime(ctx)
	state := lib.GetState(ctx)
	if state == nil ***REMOVED***
		return nil, errInvokeRPCInInitContext
	***REMOVED***
	if c.conn == nil ***REMOVED***
		return nil, errors.New("no gRPC connection, you must call connect first")
	***REMOVED***
	if method == "" ***REMOVED***
		return nil, errors.New("method to invoke cannot be empty")
	***REMOVED***
	if method[0] != '/' ***REMOVED***
		method = "/" + method
	***REMOVED***
	md := c.mds[method]
	if md == nil ***REMOVED***
		return nil, fmt.Errorf("method %q not found in file descriptors", method)
	***REMOVED***

	tags := state.CloneTags()
	timeout := 60 * time.Second

	ctx = metadata.NewOutgoingContext(ctx, metadata.New(nil))
	for k, v := range params ***REMOVED***
		switch k ***REMOVED***
		case "headers":
			rawHeaders, ok := v.(map[string]interface***REMOVED******REMOVED***)
			if !ok ***REMOVED***
				return nil, errors.New("headers must be an object with key-value pairs")
			***REMOVED***
			for hk, kv := range rawHeaders ***REMOVED***
				// TODO(rogchap): Should we manage a string slice?
				strVal, ok := kv.(string)
				if !ok ***REMOVED***
					return nil, fmt.Errorf("header %q value must be a string", hk)
				***REMOVED***
				ctx = metadata.AppendToOutgoingContext(ctx, hk, strVal)
			***REMOVED***
		case "tags":
			rawTags, ok := v.(map[string]interface***REMOVED******REMOVED***)
			if !ok ***REMOVED***
				return nil, errors.New("tags must be an object with key-value pairs")
			***REMOVED***
			for tk, tv := range rawTags ***REMOVED***
				strVal, ok := tv.(string)
				if !ok ***REMOVED***
					return nil, fmt.Errorf("tag %q value must be a string", tk)
				***REMOVED***
				tags[tk] = strVal
			***REMOVED***
		case "timeout":
			var err error
			timeout, err = types.GetDurationValue(v)
			if err != nil ***REMOVED***
				return nil, fmt.Errorf("invalid timeout value: %w", err)
			***REMOVED***
		default:
			return nil, fmt.Errorf("unknown param: %q", k)
		***REMOVED***
	***REMOVED***
	if state.Options.SystemTags.Has(stats.TagURL) ***REMOVED***
		tags["url"] = fmt.Sprintf("%s%s", c.conn.Target(), method)
	***REMOVED***
	parts := strings.Split(method[1:], "/")
	if state.Options.SystemTags.Has(stats.TagService) ***REMOVED***
		tags["service"] = parts[0]
	***REMOVED***
	if state.Options.SystemTags.Has(stats.TagMethod) ***REMOVED***
		tags["method"] = parts[1]
	***REMOVED***

	// Only set the name system tag if the user didn't explicitly set it beforehand
	if _, ok := tags["name"]; !ok && state.Options.SystemTags.Has(stats.TagName) ***REMOVED***
		tags["name"] = method
	***REMOVED***

	ctx = withTags(ctx, tags)

	reqdm := dynamicpb.NewMessage(md.Input())
	***REMOVED***
		b, err := req.ToObject(rt).MarshalJSON()
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("unable to serialise request object: %v", err)
		***REMOVED***
		if err := protojson.Unmarshal(b, reqdm); err != nil ***REMOVED***
			return nil, fmt.Errorf("unable to serialise request object to protocol buffer: %v", err)
		***REMOVED***
	***REMOVED***

	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resp := dynamicpb.NewMessage(md.Output())
	header, trailer := metadata.New(nil), metadata.New(nil)
	err := c.conn.Invoke(reqCtx, method, reqdm, resp, grpc.Header(&header), grpc.Trailer(&trailer))

	var response Response
	response.Headers = header
	response.Trailers = trailer

	marshaler := protojson.MarshalOptions***REMOVED***EmitUnpopulated: true***REMOVED***

	if err != nil ***REMOVED***
		sterr := status.Convert(err)
		response.Status = sterr.Code()

		// (rogchap) when you access a JSON property in goja, you are actually accessing the underling
		// Go type (struct, map, slice etc); because these are dynamic messages the Unmarshaled JSON does
		// not map back to a "real" field or value (as a normal Go type would). If we don't marshal and then
		// unmarshal back to a map, you will get "undefined" when accessing JSON properties, even when
		// JSON.Stringify() shows the object to be correctly present.

		raw, _ := marshaler.Marshal(sterr.Proto())
		errMsg := make(map[string]interface***REMOVED******REMOVED***)
		_ = json.Unmarshal(raw, &errMsg)
		response.Error = errMsg
	***REMOVED***

	if resp != nil ***REMOVED***
		// (rogchap) there is a lot of marshaling/unmarshaling here, but if we just pass the dynamic message
		// the default Marshaller would be used, which would strip any zero/default values from the JSON.
		// eg. given this message:
		// message Point ***REMOVED***
		//    double x = 1;
		// 	  double y = 2;
		// 	  double z = 3;
		// ***REMOVED***
		// and a value like this:
		// msg := Point***REMOVED***X: 6, Y: 4, Z: 0***REMOVED***
		// would result in JSON output:
		// ***REMOVED***"x":6,"y":4***REMOVED***
		// rather than the desired:
		// ***REMOVED***"x":6,"y":4,"z":0***REMOVED***
		raw, _ := marshaler.Marshal(resp)
		msg := make(map[string]interface***REMOVED******REMOVED***)
		_ = json.Unmarshal(raw, &msg)
		response.Message = msg
	***REMOVED***
	return &response, nil
***REMOVED***

// Close will close the client gRPC connection
func (c *Client) Close() error ***REMOVED***
	if c == nil || c.conn == nil ***REMOVED***
		return nil
	***REMOVED***
	err := c.conn.Close()
	c.conn = nil

	return err
***REMOVED***

// TagConn implements the stats.Handler interface
func (*Client) TagConn(ctx context.Context, _ *grpcstats.ConnTagInfo) context.Context ***REMOVED***
	// noop
	return ctx
***REMOVED***

// HandleConn implements the stats.Handler interface
func (*Client) HandleConn(context.Context, grpcstats.ConnStats) ***REMOVED***
	// noop
***REMOVED***

// TagRPC implements the stats.Handler interface
func (*Client) TagRPC(ctx context.Context, _ *grpcstats.RPCTagInfo) context.Context ***REMOVED***
	// noop
	return ctx
***REMOVED***

// HandleRPC implements the stats.Handler interface
func (c *Client) HandleRPC(ctx context.Context, stat grpcstats.RPCStats) ***REMOVED***
	state := lib.GetState(ctx)
	tags := getTags(ctx)
	switch s := stat.(type) ***REMOVED***
	case *grpcstats.OutHeader:
		if state.Options.SystemTags.Has(stats.TagIP) && s.RemoteAddr != nil ***REMOVED***
			if ip, _, err := net.SplitHostPort(s.RemoteAddr.String()); err == nil ***REMOVED***
				tags["ip"] = ip
			***REMOVED***
		***REMOVED***
	case *grpcstats.End:
		if state.Options.SystemTags.Has(stats.TagStatus) ***REMOVED***
			tags["status"] = strconv.Itoa(int(status.Code(s.Error)))
		***REMOVED***

		mTags := map[string]string(tags)
		sampleTags := stats.IntoSampleTags(&mTags)
		stats.PushIfNotDone(ctx, state.Samples, stats.ConnectedSamples***REMOVED***
			Samples: []stats.Sample***REMOVED***
				***REMOVED***
					Metric: state.BuiltinMetrics.GRPCReqDuration,
					Tags:   sampleTags,
					Value:  stats.D(s.EndTime.Sub(s.BeginTime)),
					Time:   s.EndTime,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***)
	***REMOVED***

	// (rogchap) Re-using --http-debug flag as gRPC is technically still HTTP
	if state.Options.HTTPDebug.String != "" ***REMOVED***
		logger := state.Logger.WithField("source", "http-debug")
		httpDebugOption := state.Options.HTTPDebug.String
		debugStat(stat, logger, httpDebugOption)
	***REMOVED***
***REMOVED***

func debugStat(stat grpcstats.RPCStats, logger logrus.FieldLogger, httpDebugOption string) ***REMOVED***
	switch s := stat.(type) ***REMOVED***
	case *grpcstats.OutHeader:
		logger.Infof("Out Header:\nFull Method: %s\nRemote Address: %s\n%s\n",
			s.FullMethod, s.RemoteAddr, formatMetadata(s.Header))
	case *grpcstats.OutTrailer:
		if len(s.Trailer) > 0 ***REMOVED***
			logger.Infof("Out Trailer:\n%s\n", formatMetadata(s.Trailer))
		***REMOVED***
	case *grpcstats.OutPayload:
		if httpDebugOption == "full" ***REMOVED***
			logger.Infof("Out Payload:\nWire Length: %d\nSent Time: %s\n%s\n\n",
				s.WireLength, s.SentTime, formatPayload(s.Payload))
		***REMOVED***
	case *grpcstats.InHeader:
		if len(s.Header) > 0 ***REMOVED***
			logger.Infof("In Header:\nWire Length: %d\n%s\n", s.WireLength, formatMetadata(s.Header))
		***REMOVED***
	case *grpcstats.InTrailer:
		if len(s.Trailer) > 0 ***REMOVED***
			logger.Infof("In Trailer:\nWire Length: %d\n%s\n", s.WireLength, formatMetadata(s.Trailer))
		***REMOVED***
	case *grpcstats.InPayload:
		if httpDebugOption == "full" ***REMOVED***
			logger.Infof("In Payload:\nWire Length: %d\nReceived Time: %s\n%s\n\n",
				s.WireLength, s.RecvTime, formatPayload(s.Payload))
		***REMOVED***
	***REMOVED***
***REMOVED***

func formatMetadata(md metadata.MD) string ***REMOVED***
	var sb strings.Builder
	for k, v := range md ***REMOVED***
		sb.WriteString(k)
		sb.WriteString(": ")
		sb.WriteString(strings.Join(v, ", "))
		sb.WriteRune('\n')
	***REMOVED***

	return sb.String()
***REMOVED***

func formatPayload(payload interface***REMOVED******REMOVED***) string ***REMOVED***
	msg, ok := payload.(proto.Message)
	if !ok ***REMOVED***
		// check to see if we are dealing with a APIv1 message
		msgV1, ok := payload.(protoV1.Message)
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		msg = protoV1.MessageV2(msgV1)
	***REMOVED***

	marshaler := prototext.MarshalOptions***REMOVED***
		Multiline: true,
		Indent:    "  ",
	***REMOVED***
	b, err := marshaler.Marshal(msg)
	if err != nil ***REMOVED***
		return ""
	***REMOVED***
	return string(b)
***REMOVED***
