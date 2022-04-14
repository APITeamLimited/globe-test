package grpcext

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync/atomic"
	"testing"

	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

func TestInvoke(t *testing.T) ***REMOVED***
	t.Parallel()

	helloReply := func(in, out *dynamicpb.Message, _ ...grpc.CallOption) error ***REMOVED***
		err := protojson.Unmarshal([]byte(`***REMOVED***"reply":"text reply"***REMOVED***`), out)
		require.NoError(t, err)

		return nil
	***REMOVED***

	c := Conn***REMOVED***raw: invokemock(helloReply)***REMOVED***
	r := Request***REMOVED***
		MethodDescriptor: methodFromProto("SayHello"),
		Message:          []byte(`***REMOVED***"greeting":"text request"***REMOVED***`),
	***REMOVED***
	res, err := c.Invoke(context.Background(), "/hello.HelloService/SayHello", metadata.New(nil), r)
	require.NoError(t, err)

	assert.Equal(t, codes.OK, res.Status)
	assert.Equal(t, map[string]interface***REMOVED******REMOVED******REMOVED***"reply": "text reply"***REMOVED***, res.Message)
	assert.Empty(t, res.Error)
***REMOVED***

func TestInvokeWithCallOptions(t *testing.T) ***REMOVED***
	t.Parallel()

	reply := func(in, out *dynamicpb.Message, opts ...grpc.CallOption) error ***REMOVED***
		assert.Len(t, opts, 3) // two by default plus one injected
		return nil
	***REMOVED***

	c := Conn***REMOVED***raw: invokemock(reply)***REMOVED***
	r := Request***REMOVED***
		MethodDescriptor: methodFromProto("NoOp"),
		Message:          []byte(`***REMOVED******REMOVED***`),
	***REMOVED***
	res, err := c.Invoke(context.Background(), "/hello.HelloService/NoOp", metadata.New(nil), r, grpc.UseCompressor("fakeone"))
	require.NoError(t, err)
	assert.NotNil(t, res)
***REMOVED***

func TestInvokeReturnError(t *testing.T) ***REMOVED***
	t.Parallel()

	helloReply := func(in, out *dynamicpb.Message, _ ...grpc.CallOption) error ***REMOVED***
		return fmt.Errorf("test error")
	***REMOVED***

	c := Conn***REMOVED***raw: invokemock(helloReply)***REMOVED***
	r := Request***REMOVED***
		MethodDescriptor: methodFromProto("SayHello"),
		Message:          []byte(`***REMOVED***"greeting":"text request"***REMOVED***`),
	***REMOVED***
	res, err := c.Invoke(context.Background(), "/hello.HelloService/SayHello", metadata.New(nil), r)
	require.NoError(t, err)

	assert.Equal(t, codes.Unknown, res.Status)
	assert.NotEmpty(t, res.Error)
	assert.Equal(t, map[string]interface***REMOVED******REMOVED******REMOVED***"reply": ""***REMOVED***, res.Message)
***REMOVED***

func TestConnInvokeInvalid(t *testing.T) ***REMOVED***
	t.Parallel()

	var (
		// valid arguments
		ctx        = context.Background()
		url        = "not-empty-url-for-method"
		md         = metadata.New(nil)
		methodDesc = methodFromProto("SayHello")
		payload    = []byte(`***REMOVED***"greeting":"test"***REMOVED***`)
	)

	req := Request***REMOVED***
		MethodDescriptor: methodDesc,
		Message:          payload,
	***REMOVED***

	tests := []struct ***REMOVED***
		name   string
		ctx    context.Context
		md     metadata.MD
		url    string
		req    Request
		experr string
	***REMOVED******REMOVED***
		***REMOVED***
			name:   "EmptyMethod",
			ctx:    ctx,
			url:    "",
			md:     md,
			req:    req,
			experr: "url is required",
		***REMOVED***,
		***REMOVED***
			name:   "NullMethodDescriptor",
			ctx:    ctx,
			url:    url,
			md:     nil,
			req:    Request***REMOVED***Message: payload***REMOVED***,
			experr: "method descriptor is required",
		***REMOVED***,
		***REMOVED***
			name:   "NullMessage",
			ctx:    ctx,
			url:    url,
			md:     nil,
			req:    Request***REMOVED***MethodDescriptor: methodDesc***REMOVED***,
			experr: "message is required",
		***REMOVED***,
		***REMOVED***
			name:   "EmptyMessage",
			ctx:    ctx,
			url:    url,
			md:     nil,
			req:    Request***REMOVED***MethodDescriptor: methodDesc, Message: []byte***REMOVED******REMOVED******REMOVED***,
			experr: "message is required",
		***REMOVED***,
	***REMOVED***

	for _, tt := range tests ***REMOVED***
		tt := tt
		t.Run(tt.name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			c := Conn***REMOVED******REMOVED***
			res, err := c.Invoke(tt.ctx, tt.url, tt.md, tt.req)
			require.Error(t, err)
			require.Nil(t, res)
			assert.Contains(t, err.Error(), tt.experr)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestResolveFileDescriptors(t *testing.T) ***REMOVED***
	t.Parallel()

	tests := []struct ***REMOVED***
		name                string
		pkgs                []string
		services            []string
		expectedDescriptors int
	***REMOVED******REMOVED***
		***REMOVED***
			name:                "SuccessSamePackage",
			pkgs:                []string***REMOVED***"mypkg"***REMOVED***,
			services:            []string***REMOVED***"Service1", "Service2", "Service3"***REMOVED***,
			expectedDescriptors: 3,
		***REMOVED***,
		***REMOVED***
			name:                "SuccessMultiPackages",
			pkgs:                []string***REMOVED***"mypkg1", "mypkg2", "mypkg3"***REMOVED***,
			services:            []string***REMOVED***"Service", "Service", "Service"***REMOVED***,
			expectedDescriptors: 3,
		***REMOVED***,
		***REMOVED***
			name:                "DeduplicateServices",
			pkgs:                []string***REMOVED***"mypkg1"***REMOVED***,
			services:            []string***REMOVED***"Service1", "Service2", "Service1"***REMOVED***,
			expectedDescriptors: 2,
		***REMOVED***,
		***REMOVED***
			name:                "NoServices",
			services:            []string***REMOVED******REMOVED***,
			expectedDescriptors: 0,
		***REMOVED***,
	***REMOVED***

	for _, tt := range tests ***REMOVED***
		tt := tt
		t.Run(tt.name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			var (
				lsr  = &reflectpb.ListServiceResponse***REMOVED******REMOVED***
				mock = &getServiceFileDescriptorMock***REMOVED******REMOVED***
			)
			for i, service := range tt.services ***REMOVED***
				// if only one package is defined then
				// the package is the same for every service
				pkg := tt.pkgs[0]
				if len(tt.pkgs) > 1 ***REMOVED***
					pkg = tt.pkgs[i]
				***REMOVED***

				lsr.Service = append(lsr.Service, &reflectpb.ServiceResponse***REMOVED***
					Name: fmt.Sprintf("%s.%s", pkg, service),
				***REMOVED***)
				mock.pkgs = append(mock.pkgs, pkg)
				mock.names = append(mock.names, service)
			***REMOVED***

			rc := ReflectionClient***REMOVED******REMOVED***
			fdset, err := rc.resolveServiceFileDescriptors(mock, lsr)
			require.NoError(t, err)
			assert.Len(t, fdset.File, tt.expectedDescriptors)
		***REMOVED***)
	***REMOVED***
***REMOVED***

type getServiceFileDescriptorMock struct ***REMOVED***
	pkgs  []string
	names []string
	nreqs int64
***REMOVED***

func (m *getServiceFileDescriptorMock) Send(req *reflectpb.ServerReflectionRequest) error ***REMOVED***
	// TODO: check that the sent message is expected,
	// otherwise return an error
	return nil
***REMOVED***

func (m *getServiceFileDescriptorMock) Recv() (*reflectpb.ServerReflectionResponse, error) ***REMOVED***
	n := atomic.AddInt64(&m.nreqs, 1)
	ptr := func(s string) (sptr *string) ***REMOVED***
		return &s
	***REMOVED***
	index := n - 1
	fdp := &descriptorpb.FileDescriptorProto***REMOVED***
		Package: ptr(m.pkgs[index]),
		Name:    ptr(m.names[index]),
	***REMOVED***
	b, err := proto.Marshal(fdp)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	srr := &reflectpb.ServerReflectionResponse***REMOVED***
		MessageResponse: &reflectpb.ServerReflectionResponse_FileDescriptorResponse***REMOVED***
			FileDescriptorResponse: &reflectpb.FileDescriptorResponse***REMOVED***
				FileDescriptorProto: [][]byte***REMOVED***b***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	return srr, nil
***REMOVED***

func methodFromProto(method string) protoreflect.MethodDescriptor ***REMOVED***
	parser := protoparse.Parser***REMOVED***
		InferImportPaths: false,
		Accessor: protoparse.FileAccessor(func(filename string) (io.ReadCloser, error) ***REMOVED***
			b := `
syntax = "proto3";

package hello;

service HelloService ***REMOVED***
  rpc SayHello(HelloRequest) returns (HelloResponse);
  rpc NoOp(Empty) returns (Empty);
  rpc LotsOfReplies(HelloRequest) returns (stream HelloResponse);
  rpc LotsOfGreetings(stream HelloRequest) returns (HelloResponse);
  rpc BidiHello(stream HelloRequest) returns (stream HelloResponse);
***REMOVED***

message HelloRequest ***REMOVED***
  string greeting = 1;
***REMOVED***

message HelloResponse ***REMOVED***
  string reply = 1;
***REMOVED***

message Empty ***REMOVED***

***REMOVED***`
			return io.NopCloser(bytes.NewBufferString(b)), nil
		***REMOVED***),
	***REMOVED***

	fds, err := parser.ParseFiles("any-path")
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	fd, err := protodesc.NewFile(fds[0].AsFileDescriptorProto(), nil)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	services := fd.Services()
	if services.Len() == 0 ***REMOVED***
		panic("no available services")
	***REMOVED***
	return services.Get(0).Methods().ByName(protoreflect.Name(method))
***REMOVED***

// invokemock is a mock for the grpc connection supporting only unary requests.
type invokemock func(in, out *dynamicpb.Message, opts ...grpc.CallOption) error

func (im invokemock) Invoke(ctx context.Context, url string, payload interface***REMOVED******REMOVED***, reply interface***REMOVED******REMOVED***, opts ...grpc.CallOption) error ***REMOVED***
	in, ok := payload.(*dynamicpb.Message)
	if !ok ***REMOVED***
		return fmt.Errorf("unexpected type for payload")
	***REMOVED***
	out, ok := reply.(*dynamicpb.Message)
	if !ok ***REMOVED***
		return fmt.Errorf("unexpected type for reply")
	***REMOVED***
	return im(in, out, opts...)
***REMOVED***

func (invokemock) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) ***REMOVED***
	panic("not implemented")
***REMOVED***

func (invokemock) Close() error ***REMOVED***
	return nil
***REMOVED***
