package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/APITeamLimited/globe-test/worker/js/common"
	"github.com/APITeamLimited/globe-test/worker/js/modules"
	"github.com/APITeamLimited/globe-test/worker/libWorker/netext/grpcext"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"

	"github.com/dop251/goja"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

// Client represents a gRPC client that can be used to make RPC requests
type Client struct ***REMOVED***
	mds  map[string]protoreflect.MethodDescriptor
	conn *grpcext.Conn
	vu   modules.VU
	addr string
***REMOVED***

// Load will parse the given proto files and make the file descriptors available to request.
func (c *Client) Load(importPaths []string, filenames ...string) ([]MethodInfo, error) ***REMOVED***
	if c.vu.State() != nil ***REMOVED***
		return nil, errors.New("load must be called in the init context")
	***REMOVED***

	initEnv := c.vu.InitEnv()
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

// Connect is a block dial to the gRPC server at the given address (host:port)
func (c *Client) Connect(addr string, params map[string]interface***REMOVED******REMOVED***) (bool, error) ***REMOVED***
	state := c.vu.State()
	if state == nil ***REMOVED***
		return false, common.NewInitContextError("connecting to a gRPC server in the init context is not supported")
	***REMOVED***

	p, err := c.parseConnectParams(params)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	opts := grpcext.DefaultOptions(c.vu)

	var tcred credentials.TransportCredentials
	if !p.IsPlaintext ***REMOVED***
		tlsCfg := state.TLSConfig.Clone()
		tlsCfg.NextProtos = []string***REMOVED***"h2"***REMOVED***

		// TODO(rogchap): Would be good to add support for custom RootCAs (self signed)
		tcred = credentials.NewTLS(tlsCfg)
	***REMOVED*** else ***REMOVED***
		tcred = insecure.NewCredentials()
	***REMOVED***
	opts = append(opts, grpc.WithTransportCredentials(tcred))

	if ua := state.Options.UserAgent; ua.Valid ***REMOVED***
		opts = append(opts, grpc.WithUserAgent(ua.ValueOrZero()))
	***REMOVED***

	ctx, cancel := context.WithTimeout(c.vu.Context(), p.Timeout)
	defer cancel()

	c.addr = addr
	c.conn, err = grpcext.Dial(ctx, addr, opts...)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if !p.UseReflectionProtocol ***REMOVED***
		return true, nil
	***REMOVED***
	fdset, err := c.conn.Reflect(ctx)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	_, err = c.convertToMethodInfo(fdset)
	if err != nil ***REMOVED***
		return false, fmt.Errorf("can't convert method info: %w", err)
	***REMOVED***

	return true, err
***REMOVED***

// Invoke creates and calls a unary RPC by fully qualified method name
func (c *Client) Invoke(
	method string,
	req goja.Value,
	params map[string]interface***REMOVED******REMOVED***,
) (*grpcext.Response, error) ***REMOVED***
	state := c.vu.State()
	if state == nil ***REMOVED***
		return nil, common.NewInitContextError("invoking RPC methods in the init context is not supported")
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
	methodDesc := c.mds[method]
	if methodDesc == nil ***REMOVED***
		return nil, fmt.Errorf("method %q not found in file descriptors", method)
	***REMOVED***

	p, err := c.parseParams(params)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	b, err := req.ToObject(c.vu.Runtime()).MarshalJSON()
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("unable to serialise request object: %w", err)
	***REMOVED***

	md := metadata.New(nil)
	for param, strval := range p.Metadata ***REMOVED***
		md.Append(param, strval)
	***REMOVED***

	ctx, cancel := context.WithTimeout(c.vu.Context(), p.Timeout)
	defer cancel()

	tags := state.CloneTags()
	for k, v := range p.Tags ***REMOVED***
		tags[k] = v
	***REMOVED***

	if state.Options.SystemTags.Has(workerMetrics.TagURL) ***REMOVED***
		tags["url"] = fmt.Sprintf("%s%s", c.addr, method)
	***REMOVED***
	parts := strings.Split(method[1:], "/")
	if state.Options.SystemTags.Has(workerMetrics.TagService) ***REMOVED***
		tags["service"] = parts[0]
	***REMOVED***
	if state.Options.SystemTags.Has(workerMetrics.TagMethod) ***REMOVED***
		tags["method"] = parts[1]
	***REMOVED***

	// Only set the name system tag if the user didn't explicitly set it beforehand
	if _, ok := tags["name"]; !ok && state.Options.SystemTags.Has(workerMetrics.TagName) ***REMOVED***
		tags["name"] = method
	***REMOVED***

	reqmsg := grpcext.Request***REMOVED***
		MethodDescriptor: methodDesc,
		Message:          b,
		Tags:             tags,
	***REMOVED***

	return c.conn.Invoke(ctx, method, md, reqmsg)
***REMOVED***

// Close will close the client gRPC connection
func (c *Client) Close() error ***REMOVED***
	if c.conn == nil ***REMOVED***
		return nil
	***REMOVED***
	err := c.conn.Close()
	c.conn = nil

	return err
***REMOVED***

// MethodInfo holds information on any parsed method descriptors that can be used by the goja VM
type MethodInfo struct ***REMOVED***
	Package         string
	Service         string
	FullMethod      string
	grpc.MethodInfo `json:"-" js:"-"`
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
		messages := fd.Messages()
		for i := 0; i < messages.Len(); i++ ***REMOVED***
			message := messages.Get(i)
			_, errFind := protoregistry.GlobalTypes.FindMessageByName(message.FullName())
			if errors.Is(errFind, protoregistry.NotFound) ***REMOVED***
				err = protoregistry.GlobalTypes.RegisterMessage(dynamicpb.NewMessageType(message))
				if err != nil ***REMOVED***
					return false
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return true
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return rtn, nil
***REMOVED***

type params struct ***REMOVED***
	Metadata map[string]string
	Tags     map[string]string
	Timeout  time.Duration
***REMOVED***

func (c *Client) parseParams(raw map[string]interface***REMOVED******REMOVED***) (params, error) ***REMOVED***
	p := params***REMOVED***
		Timeout: 1 * time.Minute,
	***REMOVED***
	for k, v := range raw ***REMOVED***
		switch k ***REMOVED***
		case "headers":
			c.vu.State().Logger.Warn("The headers property is deprecated, replace it with the metadata property, please.")
			fallthrough
		case "metadata":
			p.Metadata = make(map[string]string)

			rawHeaders, ok := v.(map[string]interface***REMOVED******REMOVED***)
			if !ok ***REMOVED***
				return p, errors.New("metadata must be an object with key-value pairs")
			***REMOVED***
			for hk, kv := range rawHeaders ***REMOVED***
				// TODO(rogchap): Should we manage a string slice?
				strval, ok := kv.(string)
				if !ok ***REMOVED***
					return p, fmt.Errorf("metadata %q value must be a string", hk)
				***REMOVED***
				p.Metadata[hk] = strval
			***REMOVED***
		case "tags":
			p.Tags = make(map[string]string)

			rawTags, ok := v.(map[string]interface***REMOVED******REMOVED***)
			if !ok ***REMOVED***
				return p, errors.New("tags must be an object with key-value pairs")
			***REMOVED***
			for tk, tv := range rawTags ***REMOVED***
				strVal, ok := tv.(string)
				if !ok ***REMOVED***
					return p, fmt.Errorf("tag %q value must be a string", tk)
				***REMOVED***
				p.Tags[tk] = strVal
			***REMOVED***
		case "timeout":
			var err error
			p.Timeout, err = types.GetDurationValue(v)
			if err != nil ***REMOVED***
				return p, fmt.Errorf("invalid timeout value: %w", err)
			***REMOVED***
		default:
			return p, fmt.Errorf("unknown param: %q", k)
		***REMOVED***
	***REMOVED***
	return p, nil
***REMOVED***

type connectParams struct ***REMOVED***
	IsPlaintext           bool
	UseReflectionProtocol bool
	Timeout               time.Duration
***REMOVED***

func (c *Client) parseConnectParams(raw map[string]interface***REMOVED******REMOVED***) (connectParams, error) ***REMOVED***
	params := connectParams***REMOVED***
		IsPlaintext:           false,
		UseReflectionProtocol: false,
		Timeout:               time.Minute,
	***REMOVED***
	for k, v := range raw ***REMOVED***
		switch k ***REMOVED***
		case "plaintext":
			var ok bool
			params.IsPlaintext, ok = v.(bool)
			if !ok ***REMOVED***
				return params, fmt.Errorf("invalid plaintext value: '%#v', it needs to be boolean", v)
			***REMOVED***
		case "timeout":
			var err error
			params.Timeout, err = types.GetDurationValue(v)
			if err != nil ***REMOVED***
				return params, fmt.Errorf("invalid timeout value: %w", err)
			***REMOVED***
		case "reflect":
			var ok bool
			params.UseReflectionProtocol, ok = v.(bool)
			if !ok ***REMOVED***
				return params, fmt.Errorf("invalid reflect value: '%#v', it needs to be boolean", v)
			***REMOVED***

		default:
			return params, fmt.Errorf("unknown connect param: %q", k)
		***REMOVED***
	***REMOVED***
	return params, nil
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
