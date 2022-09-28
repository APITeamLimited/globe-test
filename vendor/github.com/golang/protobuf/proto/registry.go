// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"sync"

	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/runtime/protoimpl"
)

// filePath is the path to the proto source file.
type filePath = string // e.g., "google/protobuf/descriptor.proto"

// fileDescGZIP is the compressed contents of the encoded FileDescriptorProto.
type fileDescGZIP = []byte

var fileCache sync.Map // map[filePath]fileDescGZIP

// RegisterFile is called from generated code to register the compressed
// FileDescriptorProto with the file path for a proto source file.
//
// Deprecated: Use protoregistry.GlobalFiles.RegisterFile instead.
func RegisterFile(s filePath, d fileDescGZIP) ***REMOVED***
	// Decompress the descriptor.
	zr, err := gzip.NewReader(bytes.NewReader(d))
	if err != nil ***REMOVED***
		panic(fmt.Sprintf("proto: invalid compressed file descriptor: %v", err))
	***REMOVED***
	b, err := ioutil.ReadAll(zr)
	if err != nil ***REMOVED***
		panic(fmt.Sprintf("proto: invalid compressed file descriptor: %v", err))
	***REMOVED***

	// Construct a protoreflect.FileDescriptor from the raw descriptor.
	// Note that DescBuilder.Build automatically registers the constructed
	// file descriptor with the v2 registry.
	protoimpl.DescBuilder***REMOVED***RawDescriptor: b***REMOVED***.Build()

	// Locally cache the raw descriptor form for the file.
	fileCache.Store(s, d)
***REMOVED***

// FileDescriptor returns the compressed FileDescriptorProto given the file path
// for a proto source file. It returns nil if not found.
//
// Deprecated: Use protoregistry.GlobalFiles.FindFileByPath instead.
func FileDescriptor(s filePath) fileDescGZIP ***REMOVED***
	if v, ok := fileCache.Load(s); ok ***REMOVED***
		return v.(fileDescGZIP)
	***REMOVED***

	// Find the descriptor in the v2 registry.
	var b []byte
	if fd, _ := protoregistry.GlobalFiles.FindFileByPath(s); fd != nil ***REMOVED***
		b, _ = Marshal(protodesc.ToFileDescriptorProto(fd))
	***REMOVED***

	// Locally cache the raw descriptor form for the file.
	if len(b) > 0 ***REMOVED***
		v, _ := fileCache.LoadOrStore(s, protoimpl.X.CompressGZIP(b))
		return v.(fileDescGZIP)
	***REMOVED***
	return nil
***REMOVED***

// enumName is the name of an enum. For historical reasons, the enum name is
// neither the full Go name nor the full protobuf name of the enum.
// The name is the dot-separated combination of just the proto package that the
// enum is declared within followed by the Go type name of the generated enum.
type enumName = string // e.g., "my.proto.package.GoMessage_GoEnum"

// enumsByName maps enum values by name to their numeric counterpart.
type enumsByName = map[string]int32

// enumsByNumber maps enum values by number to their name counterpart.
type enumsByNumber = map[int32]string

var enumCache sync.Map     // map[enumName]enumsByName
var numFilesCache sync.Map // map[protoreflect.FullName]int

// RegisterEnum is called from the generated code to register the mapping of
// enum value names to enum numbers for the enum identified by s.
//
// Deprecated: Use protoregistry.GlobalTypes.RegisterEnum instead.
func RegisterEnum(s enumName, _ enumsByNumber, m enumsByName) ***REMOVED***
	if _, ok := enumCache.Load(s); ok ***REMOVED***
		panic("proto: duplicate enum registered: " + s)
	***REMOVED***
	enumCache.Store(s, m)

	// This does not forward registration to the v2 registry since this API
	// lacks sufficient information to construct a complete v2 enum descriptor.
***REMOVED***

// EnumValueMap returns the mapping from enum value names to enum numbers for
// the enum of the given name. It returns nil if not found.
//
// Deprecated: Use protoregistry.GlobalTypes.FindEnumByName instead.
func EnumValueMap(s enumName) enumsByName ***REMOVED***
	if v, ok := enumCache.Load(s); ok ***REMOVED***
		return v.(enumsByName)
	***REMOVED***

	// Check whether the cache is stale. If the number of files in the current
	// package differs, then it means that some enums may have been recently
	// registered upstream that we do not know about.
	var protoPkg protoreflect.FullName
	if i := strings.LastIndexByte(s, '.'); i >= 0 ***REMOVED***
		protoPkg = protoreflect.FullName(s[:i])
	***REMOVED***
	v, _ := numFilesCache.Load(protoPkg)
	numFiles, _ := v.(int)
	if protoregistry.GlobalFiles.NumFilesByPackage(protoPkg) == numFiles ***REMOVED***
		return nil // cache is up-to-date; was not found earlier
	***REMOVED***

	// Update the enum cache for all enums declared in the given proto package.
	numFiles = 0
	protoregistry.GlobalFiles.RangeFilesByPackage(protoPkg, func(fd protoreflect.FileDescriptor) bool ***REMOVED***
		walkEnums(fd, func(ed protoreflect.EnumDescriptor) ***REMOVED***
			name := protoimpl.X.LegacyEnumName(ed)
			if _, ok := enumCache.Load(name); !ok ***REMOVED***
				m := make(enumsByName)
				evs := ed.Values()
				for i := evs.Len() - 1; i >= 0; i-- ***REMOVED***
					ev := evs.Get(i)
					m[string(ev.Name())] = int32(ev.Number())
				***REMOVED***
				enumCache.LoadOrStore(name, m)
			***REMOVED***
		***REMOVED***)
		numFiles++
		return true
	***REMOVED***)
	numFilesCache.Store(protoPkg, numFiles)

	// Check cache again for enum map.
	if v, ok := enumCache.Load(s); ok ***REMOVED***
		return v.(enumsByName)
	***REMOVED***
	return nil
***REMOVED***

// walkEnums recursively walks all enums declared in d.
func walkEnums(d interface ***REMOVED***
	Enums() protoreflect.EnumDescriptors
	Messages() protoreflect.MessageDescriptors
***REMOVED***, f func(protoreflect.EnumDescriptor)) ***REMOVED***
	eds := d.Enums()
	for i := eds.Len() - 1; i >= 0; i-- ***REMOVED***
		f(eds.Get(i))
	***REMOVED***
	mds := d.Messages()
	for i := mds.Len() - 1; i >= 0; i-- ***REMOVED***
		walkEnums(mds.Get(i), f)
	***REMOVED***
***REMOVED***

// messageName is the full name of protobuf message.
type messageName = string

var messageTypeCache sync.Map // map[messageName]reflect.Type

// RegisterType is called from generated code to register the message Go type
// for a message of the given name.
//
// Deprecated: Use protoregistry.GlobalTypes.RegisterMessage instead.
func RegisterType(m Message, s messageName) ***REMOVED***
	mt := protoimpl.X.LegacyMessageTypeOf(m, protoreflect.FullName(s))
	if err := protoregistry.GlobalTypes.RegisterMessage(mt); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	messageTypeCache.Store(s, reflect.TypeOf(m))
***REMOVED***

// RegisterMapType is called from generated code to register the Go map type
// for a protobuf message representing a map entry.
//
// Deprecated: Do not use.
func RegisterMapType(m interface***REMOVED******REMOVED***, s messageName) ***REMOVED***
	t := reflect.TypeOf(m)
	if t.Kind() != reflect.Map ***REMOVED***
		panic(fmt.Sprintf("invalid map kind: %v", t))
	***REMOVED***
	if _, ok := messageTypeCache.Load(s); ok ***REMOVED***
		panic(fmt.Errorf("proto: duplicate proto message registered: %s", s))
	***REMOVED***
	messageTypeCache.Store(s, t)
***REMOVED***

// MessageType returns the message type for a named message.
// It returns nil if not found.
//
// Deprecated: Use protoregistry.GlobalTypes.FindMessageByName instead.
func MessageType(s messageName) reflect.Type ***REMOVED***
	if v, ok := messageTypeCache.Load(s); ok ***REMOVED***
		return v.(reflect.Type)
	***REMOVED***

	// Derive the message type from the v2 registry.
	var t reflect.Type
	if mt, _ := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(s)); mt != nil ***REMOVED***
		t = messageGoType(mt)
	***REMOVED***

	// If we could not get a concrete type, it is possible that it is a
	// pseudo-message for a map entry.
	if t == nil ***REMOVED***
		d, _ := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(s))
		if md, _ := d.(protoreflect.MessageDescriptor); md != nil && md.IsMapEntry() ***REMOVED***
			kt := goTypeForField(md.Fields().ByNumber(1))
			vt := goTypeForField(md.Fields().ByNumber(2))
			t = reflect.MapOf(kt, vt)
		***REMOVED***
	***REMOVED***

	// Locally cache the message type for the given name.
	if t != nil ***REMOVED***
		v, _ := messageTypeCache.LoadOrStore(s, t)
		return v.(reflect.Type)
	***REMOVED***
	return nil
***REMOVED***

func goTypeForField(fd protoreflect.FieldDescriptor) reflect.Type ***REMOVED***
	switch k := fd.Kind(); k ***REMOVED***
	case protoreflect.EnumKind:
		if et, _ := protoregistry.GlobalTypes.FindEnumByName(fd.Enum().FullName()); et != nil ***REMOVED***
			return enumGoType(et)
		***REMOVED***
		return reflect.TypeOf(protoreflect.EnumNumber(0))
	case protoreflect.MessageKind, protoreflect.GroupKind:
		if mt, _ := protoregistry.GlobalTypes.FindMessageByName(fd.Message().FullName()); mt != nil ***REMOVED***
			return messageGoType(mt)
		***REMOVED***
		return reflect.TypeOf((*protoreflect.Message)(nil)).Elem()
	default:
		return reflect.TypeOf(fd.Default().Interface())
	***REMOVED***
***REMOVED***

func enumGoType(et protoreflect.EnumType) reflect.Type ***REMOVED***
	return reflect.TypeOf(et.New(0))
***REMOVED***

func messageGoType(mt protoreflect.MessageType) reflect.Type ***REMOVED***
	return reflect.TypeOf(MessageV1(mt.Zero().Interface()))
***REMOVED***

// MessageName returns the full protobuf name for the given message type.
//
// Deprecated: Use protoreflect.MessageDescriptor.FullName instead.
func MessageName(m Message) messageName ***REMOVED***
	if m == nil ***REMOVED***
		return ""
	***REMOVED***
	if m, ok := m.(interface***REMOVED*** XXX_MessageName() messageName ***REMOVED***); ok ***REMOVED***
		return m.XXX_MessageName()
	***REMOVED***
	return messageName(protoimpl.X.MessageDescriptorOf(m).FullName())
***REMOVED***

// RegisterExtension is called from the generated code to register
// the extension descriptor.
//
// Deprecated: Use protoregistry.GlobalTypes.RegisterExtension instead.
func RegisterExtension(d *ExtensionDesc) ***REMOVED***
	if err := protoregistry.GlobalTypes.RegisterExtension(d); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

type extensionsByNumber = map[int32]*ExtensionDesc

var extensionCache sync.Map // map[messageName]extensionsByNumber

// RegisteredExtensions returns a map of the registered extensions for the
// provided protobuf message, indexed by the extension field number.
//
// Deprecated: Use protoregistry.GlobalTypes.RangeExtensionsByMessage instead.
func RegisteredExtensions(m Message) extensionsByNumber ***REMOVED***
	// Check whether the cache is stale. If the number of extensions for
	// the given message differs, then it means that some extensions were
	// recently registered upstream that we do not know about.
	s := MessageName(m)
	v, _ := extensionCache.Load(s)
	xs, _ := v.(extensionsByNumber)
	if protoregistry.GlobalTypes.NumExtensionsByMessage(protoreflect.FullName(s)) == len(xs) ***REMOVED***
		return xs // cache is up-to-date
	***REMOVED***

	// Cache is stale, re-compute the extensions map.
	xs = make(extensionsByNumber)
	protoregistry.GlobalTypes.RangeExtensionsByMessage(protoreflect.FullName(s), func(xt protoreflect.ExtensionType) bool ***REMOVED***
		if xd, ok := xt.(*ExtensionDesc); ok ***REMOVED***
			xs[int32(xt.TypeDescriptor().Number())] = xd
		***REMOVED*** else ***REMOVED***
			// TODO: This implies that the protoreflect.ExtensionType is a
			// custom type not generated by protoc-gen-go. We could try and
			// convert the type to an ExtensionDesc.
		***REMOVED***
		return true
	***REMOVED***)
	extensionCache.Store(s, xs)
	return xs
***REMOVED***
