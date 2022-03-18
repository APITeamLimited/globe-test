// Package sourceinfo provides the ability to register and query source code info
// for file descriptors that are compiled into the binary. This data is registered
// by code generated from the protoc-gen-gosrcinfo plugin.
//
// The standard descriptors bundled into the compiled binary are stripped of source
// code info, to reduce binary size and reduce runtime memory footprint. However,
// the source code info can be very handy and worth the size cost when used with
// gRPC services and the server reflection service. Without source code info, the
// descriptors that a client downloads from the reflection service have no comments.
// But the presence of comments, and the ability to show them to humans, can greatly
// improve the utility of user agents that use the reflection service.
//
// When the protoc-gen-gosrcinfo plugin is used, the desc.Load* methods, which load
// descriptors for compiled-in elements, will automatically include source code
// info, using the data registered with this package.
//
// In order to make the reflection service use this functionality, you will need to
// be using v1.45 or higher of the Go runtime for gRPC (google.golang.org/grpc). The
// following snippet demonstrates how to do this in your server. Do this instead of
// using the reflection.Register function:
//
//    refSvr := reflection.NewServer(reflection.ServerOptions***REMOVED***
//        Services:           grpcServer,
//        DescriptorResolver: sourceinfo.GlobalFiles,
//        ExtensionResolver:  sourceinfo.GlobalFiles,
//    ***REMOVED***)
//    grpc_reflection_v1alpha.RegisterServerReflectionServer(grpcServer, refSvr)
//
package sourceinfo

import (
	"fmt"
	"sync"

	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

var (
	// GlobalFiles is a registry of descriptors that include source code info, if the
	// file they belong to were processed with protoc-gen-gosrcinfo.
	//
	// If is mean to serve as a drop-in alternative to protoregistry.GlobalFiles that
	// can include source code info in the returned descriptors.
	GlobalFiles Resolver = registry***REMOVED******REMOVED***

	mu               sync.RWMutex
	sourceInfoByFile = map[string]*descriptorpb.SourceCodeInfo***REMOVED******REMOVED***
	fileDescriptors  = map[protoreflect.FileDescriptor]protoreflect.FileDescriptor***REMOVED******REMOVED***
)

type Resolver interface ***REMOVED***
	protodesc.Resolver
	protoregistry.ExtensionTypeResolver
	RangeExtensionsByMessage(message protoreflect.FullName, f func(protoreflect.ExtensionType) bool)
***REMOVED***

// RegisterSourceInfo registers the given source code info for the file descriptor
// with the given path/name.
//
// This is automatically used from generated code if using the protoc-gen-gosrcinfo
// plugin.
func RegisterSourceInfo(file string, srcInfo *descriptorpb.SourceCodeInfo) ***REMOVED***
	mu.Lock()
	defer mu.Unlock()
	sourceInfoByFile[file] = srcInfo
***REMOVED***

// SourceInfoForFile queries for any registered source code info for the file
// descriptor with the given path/name. It returns nil if no source code info
// was registered.
func SourceInfoForFile(file string) *descriptorpb.SourceCodeInfo ***REMOVED***
	mu.RLock()
	defer mu.RUnlock()
	return sourceInfoByFile[file]
***REMOVED***

func getFile(fd protoreflect.FileDescriptor) protoreflect.FileDescriptor ***REMOVED***
	if fd == nil ***REMOVED***
		return nil
	***REMOVED***

	mu.RLock()
	result := fileDescriptors[fd]
	mu.RUnlock()

	if result != nil ***REMOVED***
		return result
	***REMOVED***

	mu.Lock()
	defer mu.Unlock()
	// double-check, in case it was added to map while upgrading lock
	result = fileDescriptors[fd]
	if result != nil ***REMOVED***
		return result
	***REMOVED***

	srcInfo := sourceInfoByFile[fd.Path()]
	if len(srcInfo.GetLocation()) > 0 ***REMOVED***
		result = &fileDescriptor***REMOVED***
			FileDescriptor: fd,
			locs: &sourceLocations***REMOVED***
				orig: srcInfo.Location,
			***REMOVED***,
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// nothing to do; don't bother wrapping
		result = fd
	***REMOVED***
	fileDescriptors[fd] = result
	return result
***REMOVED***

type registry struct***REMOVED******REMOVED***

var _ protodesc.Resolver = &registry***REMOVED******REMOVED***

func (r registry) FindFileByPath(path string) (protoreflect.FileDescriptor, error) ***REMOVED***
	fd, err := protoregistry.GlobalFiles.FindFileByPath(path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return getFile(fd), nil
***REMOVED***

func (r registry) FindDescriptorByName(name protoreflect.FullName) (protoreflect.Descriptor, error) ***REMOVED***
	d, err := protoregistry.GlobalFiles.FindDescriptorByName(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	switch d := d.(type) ***REMOVED***
	case protoreflect.FileDescriptor:
		return getFile(d), nil
	case protoreflect.MessageDescriptor:
		return messageDescriptor***REMOVED***d***REMOVED***, nil
	case protoreflect.ExtensionTypeDescriptor:
		return extensionDescriptor***REMOVED***d***REMOVED***, nil
	case protoreflect.FieldDescriptor:
		return fieldDescriptor***REMOVED***d***REMOVED***, nil
	case protoreflect.OneofDescriptor:
		return oneOfDescriptor***REMOVED***d***REMOVED***, nil
	case protoreflect.EnumDescriptor:
		return enumDescriptor***REMOVED***d***REMOVED***, nil
	case protoreflect.EnumValueDescriptor:
		return enumValueDescriptor***REMOVED***d***REMOVED***, nil
	case protoreflect.ServiceDescriptor:
		return serviceDescriptor***REMOVED***d***REMOVED***, nil
	case protoreflect.MethodDescriptor:
		return methodDescriptor***REMOVED***d***REMOVED***, nil
	default:
		return nil, fmt.Errorf("unrecognized descriptor type: %T", d)
	***REMOVED***
***REMOVED***

func (r registry) FindExtensionByName(field protoreflect.FullName) (protoreflect.ExtensionType, error) ***REMOVED***
	xt, err := protoregistry.GlobalTypes.FindExtensionByName(field)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return extensionType***REMOVED***xt***REMOVED***, nil
***REMOVED***

func (r registry) FindExtensionByNumber(message protoreflect.FullName, field protoreflect.FieldNumber) (protoreflect.ExtensionType, error) ***REMOVED***
	xt, err := protoregistry.GlobalTypes.FindExtensionByNumber(message, field)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return extensionType***REMOVED***xt***REMOVED***, nil
***REMOVED***

func (r registry) RangeExtensionsByMessage(message protoreflect.FullName, fn func(protoreflect.ExtensionType) bool) ***REMOVED***
	protoregistry.GlobalTypes.RangeExtensionsByMessage(message, func(xt protoreflect.ExtensionType) bool ***REMOVED***
		return fn(extensionType***REMOVED***xt***REMOVED***)
	***REMOVED***)
***REMOVED***
