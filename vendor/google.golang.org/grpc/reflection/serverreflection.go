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

/*
Package reflection implements server reflection service.

The service implemented is defined in:
https://github.com/grpc/grpc/blob/master/src/proto/grpc/reflection/v1alpha/reflection.proto.

To register server reflection on a gRPC server:
	import "google.golang.org/grpc/reflection"

	s := grpc.NewServer()
	pb.RegisterYourOwnServer(s, &server***REMOVED******REMOVED***)

	// Register reflection service on gRPC server.
	reflection.Register(s)

	s.Serve(lis)

*/
package reflection // import "google.golang.org/grpc/reflection"

import (
	"io"
	"sort"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// GRPCServer is the interface provided by a gRPC server. It is implemented by
// *grpc.Server, but could also be implemented by other concrete types. It acts
// as a registry, for accumulating the services exposed by the server.
type GRPCServer interface ***REMOVED***
	grpc.ServiceRegistrar
	ServiceInfoProvider
***REMOVED***

var _ GRPCServer = (*grpc.Server)(nil)

// Register registers the server reflection service on the given gRPC server.
func Register(s GRPCServer) ***REMOVED***
	svr := NewServer(ServerOptions***REMOVED***Services: s***REMOVED***)
	rpb.RegisterServerReflectionServer(s, svr)
***REMOVED***

// ServiceInfoProvider is an interface used to retrieve metadata about the
// services to expose.
//
// The reflection service is only interested in the service names, but the
// signature is this way so that *grpc.Server implements it. So it is okay
// for a custom implementation to return zero values for the
// grpc.ServiceInfo values in the map.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type ServiceInfoProvider interface ***REMOVED***
	GetServiceInfo() map[string]grpc.ServiceInfo
***REMOVED***

// ExtensionResolver is the interface used to query details about extensions.
// This interface is satisfied by protoregistry.GlobalTypes.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type ExtensionResolver interface ***REMOVED***
	protoregistry.ExtensionTypeResolver
	RangeExtensionsByMessage(message protoreflect.FullName, f func(protoreflect.ExtensionType) bool)
***REMOVED***

// ServerOptions represents the options used to construct a reflection server.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type ServerOptions struct ***REMOVED***
	// The source of advertised RPC services. If not specified, the reflection
	// server will report an empty list when asked to list services.
	//
	// This value will typically be a *grpc.Server. But the set of advertised
	// services can be customized by wrapping a *grpc.Server or using an
	// alternate implementation that returns a custom set of service names.
	Services ServiceInfoProvider
	// Optional resolver used to load descriptors. If not specified,
	// protoregistry.GlobalFiles will be used.
	DescriptorResolver protodesc.Resolver
	// Optional resolver used to query for known extensions. If not specified,
	// protoregistry.GlobalTypes will be used.
	ExtensionResolver ExtensionResolver
***REMOVED***

// NewServer returns a reflection server implementation using the given options.
// This can be used to customize behavior of the reflection service. Most usages
// should prefer to use Register instead.
//
// Experimental
//
// Notice: This function is EXPERIMENTAL and may be changed or removed in a
// later release.
func NewServer(opts ServerOptions) rpb.ServerReflectionServer ***REMOVED***
	if opts.DescriptorResolver == nil ***REMOVED***
		opts.DescriptorResolver = protoregistry.GlobalFiles
	***REMOVED***
	if opts.ExtensionResolver == nil ***REMOVED***
		opts.ExtensionResolver = protoregistry.GlobalTypes
	***REMOVED***
	return &serverReflectionServer***REMOVED***
		s:            opts.Services,
		descResolver: opts.DescriptorResolver,
		extResolver:  opts.ExtensionResolver,
	***REMOVED***
***REMOVED***

type serverReflectionServer struct ***REMOVED***
	rpb.UnimplementedServerReflectionServer
	s            ServiceInfoProvider
	descResolver protodesc.Resolver
	extResolver  ExtensionResolver
***REMOVED***

// fileDescWithDependencies returns a slice of serialized fileDescriptors in
// wire format ([]byte). The fileDescriptors will include fd and all the
// transitive dependencies of fd with names not in sentFileDescriptors.
func (s *serverReflectionServer) fileDescWithDependencies(fd protoreflect.FileDescriptor, sentFileDescriptors map[string]bool) ([][]byte, error) ***REMOVED***
	var r [][]byte
	queue := []protoreflect.FileDescriptor***REMOVED***fd***REMOVED***
	for len(queue) > 0 ***REMOVED***
		currentfd := queue[0]
		queue = queue[1:]
		if sent := sentFileDescriptors[currentfd.Path()]; len(r) == 0 || !sent ***REMOVED***
			sentFileDescriptors[currentfd.Path()] = true
			fdProto := protodesc.ToFileDescriptorProto(currentfd)
			currentfdEncoded, err := proto.Marshal(fdProto)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			r = append(r, currentfdEncoded)
		***REMOVED***
		for i := 0; i < currentfd.Imports().Len(); i++ ***REMOVED***
			queue = append(queue, currentfd.Imports().Get(i))
		***REMOVED***
	***REMOVED***
	return r, nil
***REMOVED***

// fileDescEncodingContainingSymbol finds the file descriptor containing the
// given symbol, finds all of its previously unsent transitive dependencies,
// does marshalling on them, and returns the marshalled result. The given symbol
// can be a type, a service or a method.
func (s *serverReflectionServer) fileDescEncodingContainingSymbol(name string, sentFileDescriptors map[string]bool) ([][]byte, error) ***REMOVED***
	d, err := s.descResolver.FindDescriptorByName(protoreflect.FullName(name))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return s.fileDescWithDependencies(d.ParentFile(), sentFileDescriptors)
***REMOVED***

// fileDescEncodingContainingExtension finds the file descriptor containing
// given extension, finds all of its previously unsent transitive dependencies,
// does marshalling on them, and returns the marshalled result.
func (s *serverReflectionServer) fileDescEncodingContainingExtension(typeName string, extNum int32, sentFileDescriptors map[string]bool) ([][]byte, error) ***REMOVED***
	xt, err := s.extResolver.FindExtensionByNumber(protoreflect.FullName(typeName), protoreflect.FieldNumber(extNum))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return s.fileDescWithDependencies(xt.TypeDescriptor().ParentFile(), sentFileDescriptors)
***REMOVED***

// allExtensionNumbersForTypeName returns all extension numbers for the given type.
func (s *serverReflectionServer) allExtensionNumbersForTypeName(name string) ([]int32, error) ***REMOVED***
	var numbers []int32
	s.extResolver.RangeExtensionsByMessage(protoreflect.FullName(name), func(xt protoreflect.ExtensionType) bool ***REMOVED***
		numbers = append(numbers, int32(xt.TypeDescriptor().Number()))
		return true
	***REMOVED***)
	sort.Slice(numbers, func(i, j int) bool ***REMOVED***
		return numbers[i] < numbers[j]
	***REMOVED***)
	if len(numbers) == 0 ***REMOVED***
		// maybe return an error if given type name is not known
		if _, err := s.descResolver.FindDescriptorByName(protoreflect.FullName(name)); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return numbers, nil
***REMOVED***

// listServices returns the names of services this server exposes.
func (s *serverReflectionServer) listServices() []*rpb.ServiceResponse ***REMOVED***
	serviceInfo := s.s.GetServiceInfo()
	resp := make([]*rpb.ServiceResponse, 0, len(serviceInfo))
	for svc := range serviceInfo ***REMOVED***
		resp = append(resp, &rpb.ServiceResponse***REMOVED***Name: svc***REMOVED***)
	***REMOVED***
	sort.Slice(resp, func(i, j int) bool ***REMOVED***
		return resp[i].Name < resp[j].Name
	***REMOVED***)
	return resp
***REMOVED***

// ServerReflectionInfo is the reflection service handler.
func (s *serverReflectionServer) ServerReflectionInfo(stream rpb.ServerReflection_ServerReflectionInfoServer) error ***REMOVED***
	sentFileDescriptors := make(map[string]bool)
	for ***REMOVED***
		in, err := stream.Recv()
		if err == io.EOF ***REMOVED***
			return nil
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		out := &rpb.ServerReflectionResponse***REMOVED***
			ValidHost:       in.Host,
			OriginalRequest: in,
		***REMOVED***
		switch req := in.MessageRequest.(type) ***REMOVED***
		case *rpb.ServerReflectionRequest_FileByFilename:
			var b [][]byte
			fd, err := s.descResolver.FindFileByPath(req.FileByFilename)
			if err == nil ***REMOVED***
				b, err = s.fileDescWithDependencies(fd, sentFileDescriptors)
			***REMOVED***
			if err != nil ***REMOVED***
				out.MessageResponse = &rpb.ServerReflectionResponse_ErrorResponse***REMOVED***
					ErrorResponse: &rpb.ErrorResponse***REMOVED***
						ErrorCode:    int32(codes.NotFound),
						ErrorMessage: err.Error(),
					***REMOVED***,
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				out.MessageResponse = &rpb.ServerReflectionResponse_FileDescriptorResponse***REMOVED***
					FileDescriptorResponse: &rpb.FileDescriptorResponse***REMOVED***FileDescriptorProto: b***REMOVED***,
				***REMOVED***
			***REMOVED***
		case *rpb.ServerReflectionRequest_FileContainingSymbol:
			b, err := s.fileDescEncodingContainingSymbol(req.FileContainingSymbol, sentFileDescriptors)
			if err != nil ***REMOVED***
				out.MessageResponse = &rpb.ServerReflectionResponse_ErrorResponse***REMOVED***
					ErrorResponse: &rpb.ErrorResponse***REMOVED***
						ErrorCode:    int32(codes.NotFound),
						ErrorMessage: err.Error(),
					***REMOVED***,
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				out.MessageResponse = &rpb.ServerReflectionResponse_FileDescriptorResponse***REMOVED***
					FileDescriptorResponse: &rpb.FileDescriptorResponse***REMOVED***FileDescriptorProto: b***REMOVED***,
				***REMOVED***
			***REMOVED***
		case *rpb.ServerReflectionRequest_FileContainingExtension:
			typeName := req.FileContainingExtension.ContainingType
			extNum := req.FileContainingExtension.ExtensionNumber
			b, err := s.fileDescEncodingContainingExtension(typeName, extNum, sentFileDescriptors)
			if err != nil ***REMOVED***
				out.MessageResponse = &rpb.ServerReflectionResponse_ErrorResponse***REMOVED***
					ErrorResponse: &rpb.ErrorResponse***REMOVED***
						ErrorCode:    int32(codes.NotFound),
						ErrorMessage: err.Error(),
					***REMOVED***,
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				out.MessageResponse = &rpb.ServerReflectionResponse_FileDescriptorResponse***REMOVED***
					FileDescriptorResponse: &rpb.FileDescriptorResponse***REMOVED***FileDescriptorProto: b***REMOVED***,
				***REMOVED***
			***REMOVED***
		case *rpb.ServerReflectionRequest_AllExtensionNumbersOfType:
			extNums, err := s.allExtensionNumbersForTypeName(req.AllExtensionNumbersOfType)
			if err != nil ***REMOVED***
				out.MessageResponse = &rpb.ServerReflectionResponse_ErrorResponse***REMOVED***
					ErrorResponse: &rpb.ErrorResponse***REMOVED***
						ErrorCode:    int32(codes.NotFound),
						ErrorMessage: err.Error(),
					***REMOVED***,
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				out.MessageResponse = &rpb.ServerReflectionResponse_AllExtensionNumbersResponse***REMOVED***
					AllExtensionNumbersResponse: &rpb.ExtensionNumberResponse***REMOVED***
						BaseTypeName:    req.AllExtensionNumbersOfType,
						ExtensionNumber: extNums,
					***REMOVED***,
				***REMOVED***
			***REMOVED***
		case *rpb.ServerReflectionRequest_ListServices:
			out.MessageResponse = &rpb.ServerReflectionResponse_ListServicesResponse***REMOVED***
				ListServicesResponse: &rpb.ListServiceResponse***REMOVED***
					Service: s.listServices(),
				***REMOVED***,
			***REMOVED***
		default:
			return status.Errorf(codes.InvalidArgument, "invalid MessageRequest: %v", in.MessageRequest)
		***REMOVED***

		if err := stream.Send(out); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
***REMOVED***
