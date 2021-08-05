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
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"sort"
	"sync"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
)

// GRPCServer is the interface provided by a gRPC server. It is implemented by
// *grpc.Server, but could also be implemented by other concrete types. It acts
// as a registry, for accumulating the services exposed by the server.
type GRPCServer interface ***REMOVED***
	grpc.ServiceRegistrar
	GetServiceInfo() map[string]grpc.ServiceInfo
***REMOVED***

var _ GRPCServer = (*grpc.Server)(nil)

type serverReflectionServer struct ***REMOVED***
	rpb.UnimplementedServerReflectionServer
	s GRPCServer

	initSymbols  sync.Once
	serviceNames []string
	symbols      map[string]*dpb.FileDescriptorProto // map of fully-qualified names to files
***REMOVED***

// Register registers the server reflection service on the given gRPC server.
func Register(s GRPCServer) ***REMOVED***
	rpb.RegisterServerReflectionServer(s, &serverReflectionServer***REMOVED***
		s: s,
	***REMOVED***)
***REMOVED***

// protoMessage is used for type assertion on proto messages.
// Generated proto message implements function Descriptor(), but Descriptor()
// is not part of interface proto.Message. This interface is needed to
// call Descriptor().
type protoMessage interface ***REMOVED***
	Descriptor() ([]byte, []int)
***REMOVED***

func (s *serverReflectionServer) getSymbols() (svcNames []string, symbolIndex map[string]*dpb.FileDescriptorProto) ***REMOVED***
	s.initSymbols.Do(func() ***REMOVED***
		serviceInfo := s.s.GetServiceInfo()

		s.symbols = map[string]*dpb.FileDescriptorProto***REMOVED******REMOVED***
		s.serviceNames = make([]string, 0, len(serviceInfo))
		processed := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
		for svc, info := range serviceInfo ***REMOVED***
			s.serviceNames = append(s.serviceNames, svc)
			fdenc, ok := parseMetadata(info.Metadata)
			if !ok ***REMOVED***
				continue
			***REMOVED***
			fd, err := decodeFileDesc(fdenc)
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			s.processFile(fd, processed)
		***REMOVED***
		sort.Strings(s.serviceNames)
	***REMOVED***)

	return s.serviceNames, s.symbols
***REMOVED***

func (s *serverReflectionServer) processFile(fd *dpb.FileDescriptorProto, processed map[string]struct***REMOVED******REMOVED***) ***REMOVED***
	filename := fd.GetName()
	if _, ok := processed[filename]; ok ***REMOVED***
		return
	***REMOVED***
	processed[filename] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

	prefix := fd.GetPackage()

	for _, msg := range fd.MessageType ***REMOVED***
		s.processMessage(fd, prefix, msg)
	***REMOVED***
	for _, en := range fd.EnumType ***REMOVED***
		s.processEnum(fd, prefix, en)
	***REMOVED***
	for _, ext := range fd.Extension ***REMOVED***
		s.processField(fd, prefix, ext)
	***REMOVED***
	for _, svc := range fd.Service ***REMOVED***
		svcName := fqn(prefix, svc.GetName())
		s.symbols[svcName] = fd
		for _, meth := range svc.Method ***REMOVED***
			name := fqn(svcName, meth.GetName())
			s.symbols[name] = fd
		***REMOVED***
	***REMOVED***

	for _, dep := range fd.Dependency ***REMOVED***
		fdenc := proto.FileDescriptor(dep)
		fdDep, err := decodeFileDesc(fdenc)
		if err != nil ***REMOVED***
			continue
		***REMOVED***
		s.processFile(fdDep, processed)
	***REMOVED***
***REMOVED***

func (s *serverReflectionServer) processMessage(fd *dpb.FileDescriptorProto, prefix string, msg *dpb.DescriptorProto) ***REMOVED***
	msgName := fqn(prefix, msg.GetName())
	s.symbols[msgName] = fd

	for _, nested := range msg.NestedType ***REMOVED***
		s.processMessage(fd, msgName, nested)
	***REMOVED***
	for _, en := range msg.EnumType ***REMOVED***
		s.processEnum(fd, msgName, en)
	***REMOVED***
	for _, ext := range msg.Extension ***REMOVED***
		s.processField(fd, msgName, ext)
	***REMOVED***
	for _, fld := range msg.Field ***REMOVED***
		s.processField(fd, msgName, fld)
	***REMOVED***
	for _, oneof := range msg.OneofDecl ***REMOVED***
		oneofName := fqn(msgName, oneof.GetName())
		s.symbols[oneofName] = fd
	***REMOVED***
***REMOVED***

func (s *serverReflectionServer) processEnum(fd *dpb.FileDescriptorProto, prefix string, en *dpb.EnumDescriptorProto) ***REMOVED***
	enName := fqn(prefix, en.GetName())
	s.symbols[enName] = fd

	for _, val := range en.Value ***REMOVED***
		valName := fqn(enName, val.GetName())
		s.symbols[valName] = fd
	***REMOVED***
***REMOVED***

func (s *serverReflectionServer) processField(fd *dpb.FileDescriptorProto, prefix string, fld *dpb.FieldDescriptorProto) ***REMOVED***
	fldName := fqn(prefix, fld.GetName())
	s.symbols[fldName] = fd
***REMOVED***

func fqn(prefix, name string) string ***REMOVED***
	if prefix == "" ***REMOVED***
		return name
	***REMOVED***
	return prefix + "." + name
***REMOVED***

// fileDescForType gets the file descriptor for the given type.
// The given type should be a proto message.
func (s *serverReflectionServer) fileDescForType(st reflect.Type) (*dpb.FileDescriptorProto, error) ***REMOVED***
	m, ok := reflect.Zero(reflect.PtrTo(st)).Interface().(protoMessage)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("failed to create message from type: %v", st)
	***REMOVED***
	enc, _ := m.Descriptor()

	return decodeFileDesc(enc)
***REMOVED***

// decodeFileDesc does decompression and unmarshalling on the given
// file descriptor byte slice.
func decodeFileDesc(enc []byte) (*dpb.FileDescriptorProto, error) ***REMOVED***
	raw, err := decompress(enc)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to decompress enc: %v", err)
	***REMOVED***

	fd := new(dpb.FileDescriptorProto)
	if err := proto.Unmarshal(raw, fd); err != nil ***REMOVED***
		return nil, fmt.Errorf("bad descriptor: %v", err)
	***REMOVED***
	return fd, nil
***REMOVED***

// decompress does gzip decompression.
func decompress(b []byte) ([]byte, error) ***REMOVED***
	r, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("bad gzipped descriptor: %v", err)
	***REMOVED***
	out, err := ioutil.ReadAll(r)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("bad gzipped descriptor: %v", err)
	***REMOVED***
	return out, nil
***REMOVED***

func typeForName(name string) (reflect.Type, error) ***REMOVED***
	pt := proto.MessageType(name)
	if pt == nil ***REMOVED***
		return nil, fmt.Errorf("unknown type: %q", name)
	***REMOVED***
	st := pt.Elem()

	return st, nil
***REMOVED***

func fileDescContainingExtension(st reflect.Type, ext int32) (*dpb.FileDescriptorProto, error) ***REMOVED***
	m, ok := reflect.Zero(reflect.PtrTo(st)).Interface().(proto.Message)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("failed to create message from type: %v", st)
	***REMOVED***

	var extDesc *proto.ExtensionDesc
	for id, desc := range proto.RegisteredExtensions(m) ***REMOVED***
		if id == ext ***REMOVED***
			extDesc = desc
			break
		***REMOVED***
	***REMOVED***

	if extDesc == nil ***REMOVED***
		return nil, fmt.Errorf("failed to find registered extension for extension number %v", ext)
	***REMOVED***

	return decodeFileDesc(proto.FileDescriptor(extDesc.Filename))
***REMOVED***

func (s *serverReflectionServer) allExtensionNumbersForType(st reflect.Type) ([]int32, error) ***REMOVED***
	m, ok := reflect.Zero(reflect.PtrTo(st)).Interface().(proto.Message)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("failed to create message from type: %v", st)
	***REMOVED***

	exts := proto.RegisteredExtensions(m)
	out := make([]int32, 0, len(exts))
	for id := range exts ***REMOVED***
		out = append(out, id)
	***REMOVED***
	return out, nil
***REMOVED***

// fileDescWithDependencies returns a slice of serialized fileDescriptors in
// wire format ([]byte). The fileDescriptors will include fd and all the
// transitive dependencies of fd with names not in sentFileDescriptors.
func fileDescWithDependencies(fd *dpb.FileDescriptorProto, sentFileDescriptors map[string]bool) ([][]byte, error) ***REMOVED***
	r := [][]byte***REMOVED******REMOVED***
	queue := []*dpb.FileDescriptorProto***REMOVED***fd***REMOVED***
	for len(queue) > 0 ***REMOVED***
		currentfd := queue[0]
		queue = queue[1:]
		if sent := sentFileDescriptors[currentfd.GetName()]; len(r) == 0 || !sent ***REMOVED***
			sentFileDescriptors[currentfd.GetName()] = true
			currentfdEncoded, err := proto.Marshal(currentfd)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			r = append(r, currentfdEncoded)
		***REMOVED***
		for _, dep := range currentfd.Dependency ***REMOVED***
			fdenc := proto.FileDescriptor(dep)
			fdDep, err := decodeFileDesc(fdenc)
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			queue = append(queue, fdDep)
		***REMOVED***
	***REMOVED***
	return r, nil
***REMOVED***

// fileDescEncodingByFilename finds the file descriptor for given filename,
// finds all of its previously unsent transitive dependencies, does marshalling
// on them, and returns the marshalled result.
func (s *serverReflectionServer) fileDescEncodingByFilename(name string, sentFileDescriptors map[string]bool) ([][]byte, error) ***REMOVED***
	enc := proto.FileDescriptor(name)
	if enc == nil ***REMOVED***
		return nil, fmt.Errorf("unknown file: %v", name)
	***REMOVED***
	fd, err := decodeFileDesc(enc)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return fileDescWithDependencies(fd, sentFileDescriptors)
***REMOVED***

// parseMetadata finds the file descriptor bytes specified meta.
// For SupportPackageIsVersion4, m is the name of the proto file, we
// call proto.FileDescriptor to get the byte slice.
// For SupportPackageIsVersion3, m is a byte slice itself.
func parseMetadata(meta interface***REMOVED******REMOVED***) ([]byte, bool) ***REMOVED***
	// Check if meta is the file name.
	if fileNameForMeta, ok := meta.(string); ok ***REMOVED***
		return proto.FileDescriptor(fileNameForMeta), true
	***REMOVED***

	// Check if meta is the byte slice.
	if enc, ok := meta.([]byte); ok ***REMOVED***
		return enc, true
	***REMOVED***

	return nil, false
***REMOVED***

// fileDescEncodingContainingSymbol finds the file descriptor containing the
// given symbol, finds all of its previously unsent transitive dependencies,
// does marshalling on them, and returns the marshalled result. The given symbol
// can be a type, a service or a method.
func (s *serverReflectionServer) fileDescEncodingContainingSymbol(name string, sentFileDescriptors map[string]bool) ([][]byte, error) ***REMOVED***
	_, symbols := s.getSymbols()
	fd := symbols[name]
	if fd == nil ***REMOVED***
		// Check if it's a type name that was not present in the
		// transitive dependencies of the registered services.
		if st, err := typeForName(name); err == nil ***REMOVED***
			fd, err = s.fileDescForType(st)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if fd == nil ***REMOVED***
		return nil, fmt.Errorf("unknown symbol: %v", name)
	***REMOVED***

	return fileDescWithDependencies(fd, sentFileDescriptors)
***REMOVED***

// fileDescEncodingContainingExtension finds the file descriptor containing
// given extension, finds all of its previously unsent transitive dependencies,
// does marshalling on them, and returns the marshalled result.
func (s *serverReflectionServer) fileDescEncodingContainingExtension(typeName string, extNum int32, sentFileDescriptors map[string]bool) ([][]byte, error) ***REMOVED***
	st, err := typeForName(typeName)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	fd, err := fileDescContainingExtension(st, extNum)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return fileDescWithDependencies(fd, sentFileDescriptors)
***REMOVED***

// allExtensionNumbersForTypeName returns all extension numbers for the given type.
func (s *serverReflectionServer) allExtensionNumbersForTypeName(name string) ([]int32, error) ***REMOVED***
	st, err := typeForName(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	extNums, err := s.allExtensionNumbersForType(st)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return extNums, nil
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
			b, err := s.fileDescEncodingByFilename(req.FileByFilename, sentFileDescriptors)
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
			svcNames, _ := s.getSymbols()
			serviceResponses := make([]*rpb.ServiceResponse, len(svcNames))
			for i, n := range svcNames ***REMOVED***
				serviceResponses[i] = &rpb.ServiceResponse***REMOVED***
					Name: n,
				***REMOVED***
			***REMOVED***
			out.MessageResponse = &rpb.ServerReflectionResponse_ListServicesResponse***REMOVED***
				ListServicesResponse: &rpb.ListServiceResponse***REMOVED***
					Service: serviceResponses,
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
