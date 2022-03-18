// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package protoregistry provides data structures to register and lookup
// protobuf descriptor types.
//
// The Files registry contains file descriptors and provides the ability
// to iterate over the files or lookup a specific descriptor within the files.
// Files only contains protobuf descriptors and has no understanding of Go
// type information that may be associated with each descriptor.
//
// The Types registry contains descriptor types for which there is a known
// Go type associated with that descriptor. It provides the ability to iterate
// over the registered types or lookup a type by name.
package protoregistry

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"google.golang.org/protobuf/internal/encoding/messageset"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/flags"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// conflictPolicy configures the policy for handling registration conflicts.
//
// It can be over-written at compile time with a linker-initialized variable:
//	go build -ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn"
//
// It can be over-written at program execution with an environment variable:
//	GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn ./main
//
// Neither of the above are covered by the compatibility promise and
// may be removed in a future release of this module.
var conflictPolicy = "panic" // "panic" | "warn" | "ignore"

// ignoreConflict reports whether to ignore a registration conflict
// given the descriptor being registered and the error.
// It is a variable so that the behavior is easily overridden in another file.
var ignoreConflict = func(d protoreflect.Descriptor, err error) bool ***REMOVED***
	const env = "GOLANG_PROTOBUF_REGISTRATION_CONFLICT"
	const faq = "https://developers.google.com/protocol-buffers/docs/reference/go/faq#namespace-conflict"
	policy := conflictPolicy
	if v := os.Getenv(env); v != "" ***REMOVED***
		policy = v
	***REMOVED***
	switch policy ***REMOVED***
	case "panic":
		panic(fmt.Sprintf("%v\nSee %v\n", err, faq))
	case "warn":
		fmt.Fprintf(os.Stderr, "WARNING: %v\nSee %v\n\n", err, faq)
		return true
	case "ignore":
		return true
	default:
		panic("invalid " + env + " value: " + os.Getenv(env))
	***REMOVED***
***REMOVED***

var globalMutex sync.RWMutex

// GlobalFiles is a global registry of file descriptors.
var GlobalFiles *Files = new(Files)

// GlobalTypes is the registry used by default for type lookups
// unless a local registry is provided by the user.
var GlobalTypes *Types = new(Types)

// NotFound is a sentinel error value to indicate that the type was not found.
//
// Since registry lookup can happen in the critical performance path, resolvers
// must return this exact error value, not an error wrapping it.
var NotFound = errors.New("not found")

// Files is a registry for looking up or iterating over files and the
// descriptors contained within them.
// The Find and Range methods are safe for concurrent use.
type Files struct ***REMOVED***
	// The map of descsByName contains:
	//	EnumDescriptor
	//	EnumValueDescriptor
	//	MessageDescriptor
	//	ExtensionDescriptor
	//	ServiceDescriptor
	//	*packageDescriptor
	//
	// Note that files are stored as a slice, since a package may contain
	// multiple files. Only top-level declarations are registered.
	// Note that enum values are in the top-level since that are in the same
	// scope as the parent enum.
	descsByName map[protoreflect.FullName]interface***REMOVED******REMOVED***
	filesByPath map[string]protoreflect.FileDescriptor
***REMOVED***

type packageDescriptor struct ***REMOVED***
	files []protoreflect.FileDescriptor
***REMOVED***

// RegisterFile registers the provided file descriptor.
//
// If any descriptor within the file conflicts with the descriptor of any
// previously registered file (e.g., two enums with the same full name),
// then the file is not registered and an error is returned.
//
// It is permitted for multiple files to have the same file path.
func (r *Files) RegisterFile(file protoreflect.FileDescriptor) error ***REMOVED***
	if r == GlobalFiles ***REMOVED***
		globalMutex.Lock()
		defer globalMutex.Unlock()
	***REMOVED***
	if r.descsByName == nil ***REMOVED***
		r.descsByName = map[protoreflect.FullName]interface***REMOVED******REMOVED******REMOVED***
			"": &packageDescriptor***REMOVED******REMOVED***,
		***REMOVED***
		r.filesByPath = make(map[string]protoreflect.FileDescriptor)
	***REMOVED***
	path := file.Path()
	if prev := r.filesByPath[path]; prev != nil ***REMOVED***
		r.checkGenProtoConflict(path)
		err := errors.New("file %q is already registered", file.Path())
		err = amendErrorWithCaller(err, prev, file)
		if r == GlobalFiles && ignoreConflict(file, err) ***REMOVED***
			err = nil
		***REMOVED***
		return err
	***REMOVED***

	for name := file.Package(); name != ""; name = name.Parent() ***REMOVED***
		switch prev := r.descsByName[name]; prev.(type) ***REMOVED***
		case nil, *packageDescriptor:
		default:
			err := errors.New("file %q has a package name conflict over %v", file.Path(), name)
			err = amendErrorWithCaller(err, prev, file)
			if r == GlobalFiles && ignoreConflict(file, err) ***REMOVED***
				err = nil
			***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	var err error
	var hasConflict bool
	rangeTopLevelDescriptors(file, func(d protoreflect.Descriptor) ***REMOVED***
		if prev := r.descsByName[d.FullName()]; prev != nil ***REMOVED***
			hasConflict = true
			err = errors.New("file %q has a name conflict over %v", file.Path(), d.FullName())
			err = amendErrorWithCaller(err, prev, file)
			if r == GlobalFiles && ignoreConflict(d, err) ***REMOVED***
				err = nil
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	if hasConflict ***REMOVED***
		return err
	***REMOVED***

	for name := file.Package(); name != ""; name = name.Parent() ***REMOVED***
		if r.descsByName[name] == nil ***REMOVED***
			r.descsByName[name] = &packageDescriptor***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
	p := r.descsByName[file.Package()].(*packageDescriptor)
	p.files = append(p.files, file)
	rangeTopLevelDescriptors(file, func(d protoreflect.Descriptor) ***REMOVED***
		r.descsByName[d.FullName()] = d
	***REMOVED***)
	r.filesByPath[path] = file
	return nil
***REMOVED***

// Several well-known types were hosted in the google.golang.org/genproto module
// but were later moved to this module. To avoid a weak dependency on the
// genproto module (and its relatively large set of transitive dependencies),
// we rely on a registration conflict to determine whether the genproto version
// is too old (i.e., does not contain aliases to the new type declarations).
func (r *Files) checkGenProtoConflict(path string) ***REMOVED***
	if r != GlobalFiles ***REMOVED***
		return
	***REMOVED***
	var prevPath string
	const prevModule = "google.golang.org/genproto"
	const prevVersion = "cb27e3aa (May 26th, 2020)"
	switch path ***REMOVED***
	case "google/protobuf/field_mask.proto":
		prevPath = prevModule + "/protobuf/field_mask"
	case "google/protobuf/api.proto":
		prevPath = prevModule + "/protobuf/api"
	case "google/protobuf/type.proto":
		prevPath = prevModule + "/protobuf/ptype"
	case "google/protobuf/source_context.proto":
		prevPath = prevModule + "/protobuf/source_context"
	default:
		return
	***REMOVED***
	pkgName := strings.TrimSuffix(strings.TrimPrefix(path, "google/protobuf/"), ".proto")
	pkgName = strings.Replace(pkgName, "_", "", -1) + "pb" // e.g., "field_mask" => "fieldmaskpb"
	currPath := "google.golang.org/protobuf/types/known/" + pkgName
	panic(fmt.Sprintf(""+
		"duplicate registration of %q\n"+
		"\n"+
		"The generated definition for this file has moved:\n"+
		"\tfrom: %q\n"+
		"\tto:   %q\n"+
		"A dependency on the %q module must\n"+
		"be at version %v or higher.\n"+
		"\n"+
		"Upgrade the dependency by running:\n"+
		"\tgo get -u %v\n",
		path, prevPath, currPath, prevModule, prevVersion, prevPath))
***REMOVED***

// FindDescriptorByName looks up a descriptor by the full name.
//
// This returns (nil, NotFound) if not found.
func (r *Files) FindDescriptorByName(name protoreflect.FullName) (protoreflect.Descriptor, error) ***REMOVED***
	if r == nil ***REMOVED***
		return nil, NotFound
	***REMOVED***
	if r == GlobalFiles ***REMOVED***
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	***REMOVED***
	prefix := name
	suffix := nameSuffix("")
	for prefix != "" ***REMOVED***
		if d, ok := r.descsByName[prefix]; ok ***REMOVED***
			switch d := d.(type) ***REMOVED***
			case protoreflect.EnumDescriptor:
				if d.FullName() == name ***REMOVED***
					return d, nil
				***REMOVED***
			case protoreflect.EnumValueDescriptor:
				if d.FullName() == name ***REMOVED***
					return d, nil
				***REMOVED***
			case protoreflect.MessageDescriptor:
				if d.FullName() == name ***REMOVED***
					return d, nil
				***REMOVED***
				if d := findDescriptorInMessage(d, suffix); d != nil && d.FullName() == name ***REMOVED***
					return d, nil
				***REMOVED***
			case protoreflect.ExtensionDescriptor:
				if d.FullName() == name ***REMOVED***
					return d, nil
				***REMOVED***
			case protoreflect.ServiceDescriptor:
				if d.FullName() == name ***REMOVED***
					return d, nil
				***REMOVED***
				if d := d.Methods().ByName(suffix.Pop()); d != nil && d.FullName() == name ***REMOVED***
					return d, nil
				***REMOVED***
			***REMOVED***
			return nil, NotFound
		***REMOVED***
		prefix = prefix.Parent()
		suffix = nameSuffix(name[len(prefix)+len("."):])
	***REMOVED***
	return nil, NotFound
***REMOVED***

func findDescriptorInMessage(md protoreflect.MessageDescriptor, suffix nameSuffix) protoreflect.Descriptor ***REMOVED***
	name := suffix.Pop()
	if suffix == "" ***REMOVED***
		if ed := md.Enums().ByName(name); ed != nil ***REMOVED***
			return ed
		***REMOVED***
		for i := md.Enums().Len() - 1; i >= 0; i-- ***REMOVED***
			if vd := md.Enums().Get(i).Values().ByName(name); vd != nil ***REMOVED***
				return vd
			***REMOVED***
		***REMOVED***
		if xd := md.Extensions().ByName(name); xd != nil ***REMOVED***
			return xd
		***REMOVED***
		if fd := md.Fields().ByName(name); fd != nil ***REMOVED***
			return fd
		***REMOVED***
		if od := md.Oneofs().ByName(name); od != nil ***REMOVED***
			return od
		***REMOVED***
	***REMOVED***
	if md := md.Messages().ByName(name); md != nil ***REMOVED***
		if suffix == "" ***REMOVED***
			return md
		***REMOVED***
		return findDescriptorInMessage(md, suffix)
	***REMOVED***
	return nil
***REMOVED***

type nameSuffix string

func (s *nameSuffix) Pop() (name protoreflect.Name) ***REMOVED***
	if i := strings.IndexByte(string(*s), '.'); i >= 0 ***REMOVED***
		name, *s = protoreflect.Name((*s)[:i]), (*s)[i+1:]
	***REMOVED*** else ***REMOVED***
		name, *s = protoreflect.Name((*s)), ""
	***REMOVED***
	return name
***REMOVED***

// FindFileByPath looks up a file by the path.
//
// This returns (nil, NotFound) if not found.
func (r *Files) FindFileByPath(path string) (protoreflect.FileDescriptor, error) ***REMOVED***
	if r == nil ***REMOVED***
		return nil, NotFound
	***REMOVED***
	if r == GlobalFiles ***REMOVED***
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	***REMOVED***
	if fd, ok := r.filesByPath[path]; ok ***REMOVED***
		return fd, nil
	***REMOVED***
	return nil, NotFound
***REMOVED***

// NumFiles reports the number of registered files.
func (r *Files) NumFiles() int ***REMOVED***
	if r == nil ***REMOVED***
		return 0
	***REMOVED***
	if r == GlobalFiles ***REMOVED***
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	***REMOVED***
	return len(r.filesByPath)
***REMOVED***

// RangeFiles iterates over all registered files while f returns true.
// The iteration order is undefined.
func (r *Files) RangeFiles(f func(protoreflect.FileDescriptor) bool) ***REMOVED***
	if r == nil ***REMOVED***
		return
	***REMOVED***
	if r == GlobalFiles ***REMOVED***
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	***REMOVED***
	for _, file := range r.filesByPath ***REMOVED***
		if !f(file) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// NumFilesByPackage reports the number of registered files in a proto package.
func (r *Files) NumFilesByPackage(name protoreflect.FullName) int ***REMOVED***
	if r == nil ***REMOVED***
		return 0
	***REMOVED***
	if r == GlobalFiles ***REMOVED***
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	***REMOVED***
	p, ok := r.descsByName[name].(*packageDescriptor)
	if !ok ***REMOVED***
		return 0
	***REMOVED***
	return len(p.files)
***REMOVED***

// RangeFilesByPackage iterates over all registered files in a given proto package
// while f returns true. The iteration order is undefined.
func (r *Files) RangeFilesByPackage(name protoreflect.FullName, f func(protoreflect.FileDescriptor) bool) ***REMOVED***
	if r == nil ***REMOVED***
		return
	***REMOVED***
	if r == GlobalFiles ***REMOVED***
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	***REMOVED***
	p, ok := r.descsByName[name].(*packageDescriptor)
	if !ok ***REMOVED***
		return
	***REMOVED***
	for _, file := range p.files ***REMOVED***
		if !f(file) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// rangeTopLevelDescriptors iterates over all top-level descriptors in a file
// which will be directly entered into the registry.
func rangeTopLevelDescriptors(fd protoreflect.FileDescriptor, f func(protoreflect.Descriptor)) ***REMOVED***
	eds := fd.Enums()
	for i := eds.Len() - 1; i >= 0; i-- ***REMOVED***
		f(eds.Get(i))
		vds := eds.Get(i).Values()
		for i := vds.Len() - 1; i >= 0; i-- ***REMOVED***
			f(vds.Get(i))
		***REMOVED***
	***REMOVED***
	mds := fd.Messages()
	for i := mds.Len() - 1; i >= 0; i-- ***REMOVED***
		f(mds.Get(i))
	***REMOVED***
	xds := fd.Extensions()
	for i := xds.Len() - 1; i >= 0; i-- ***REMOVED***
		f(xds.Get(i))
	***REMOVED***
	sds := fd.Services()
	for i := sds.Len() - 1; i >= 0; i-- ***REMOVED***
		f(sds.Get(i))
	***REMOVED***
***REMOVED***

// MessageTypeResolver is an interface for looking up messages.
//
// A compliant implementation must deterministically return the same type
// if no error is encountered.
//
// The Types type implements this interface.
type MessageTypeResolver interface ***REMOVED***
	// FindMessageByName looks up a message by its full name.
	// E.g., "google.protobuf.Any"
	//
	// This return (nil, NotFound) if not found.
	FindMessageByName(message protoreflect.FullName) (protoreflect.MessageType, error)

	// FindMessageByURL looks up a message by a URL identifier.
	// See documentation on google.protobuf.Any.type_url for the URL format.
	//
	// This returns (nil, NotFound) if not found.
	FindMessageByURL(url string) (protoreflect.MessageType, error)
***REMOVED***

// ExtensionTypeResolver is an interface for looking up extensions.
//
// A compliant implementation must deterministically return the same type
// if no error is encountered.
//
// The Types type implements this interface.
type ExtensionTypeResolver interface ***REMOVED***
	// FindExtensionByName looks up a extension field by the field's full name.
	// Note that this is the full name of the field as determined by
	// where the extension is declared and is unrelated to the full name of the
	// message being extended.
	//
	// This returns (nil, NotFound) if not found.
	FindExtensionByName(field protoreflect.FullName) (protoreflect.ExtensionType, error)

	// FindExtensionByNumber looks up a extension field by the field number
	// within some parent message, identified by full name.
	//
	// This returns (nil, NotFound) if not found.
	FindExtensionByNumber(message protoreflect.FullName, field protoreflect.FieldNumber) (protoreflect.ExtensionType, error)
***REMOVED***

var (
	_ MessageTypeResolver   = (*Types)(nil)
	_ ExtensionTypeResolver = (*Types)(nil)
)

// Types is a registry for looking up or iterating over descriptor types.
// The Find and Range methods are safe for concurrent use.
type Types struct ***REMOVED***
	typesByName         typesByName
	extensionsByMessage extensionsByMessage

	numEnums      int
	numMessages   int
	numExtensions int
***REMOVED***

type (
	typesByName         map[protoreflect.FullName]interface***REMOVED******REMOVED***
	extensionsByMessage map[protoreflect.FullName]extensionsByNumber
	extensionsByNumber  map[protoreflect.FieldNumber]protoreflect.ExtensionType
)

// RegisterMessage registers the provided message type.
//
// If a naming conflict occurs, the type is not registered and an error is returned.
func (r *Types) RegisterMessage(mt protoreflect.MessageType) error ***REMOVED***
	// Under rare circumstances getting the descriptor might recursively
	// examine the registry, so fetch it before locking.
	md := mt.Descriptor()

	if r == GlobalTypes ***REMOVED***
		globalMutex.Lock()
		defer globalMutex.Unlock()
	***REMOVED***

	if err := r.register("message", md, mt); err != nil ***REMOVED***
		return err
	***REMOVED***
	r.numMessages++
	return nil
***REMOVED***

// RegisterEnum registers the provided enum type.
//
// If a naming conflict occurs, the type is not registered and an error is returned.
func (r *Types) RegisterEnum(et protoreflect.EnumType) error ***REMOVED***
	// Under rare circumstances getting the descriptor might recursively
	// examine the registry, so fetch it before locking.
	ed := et.Descriptor()

	if r == GlobalTypes ***REMOVED***
		globalMutex.Lock()
		defer globalMutex.Unlock()
	***REMOVED***

	if err := r.register("enum", ed, et); err != nil ***REMOVED***
		return err
	***REMOVED***
	r.numEnums++
	return nil
***REMOVED***

// RegisterExtension registers the provided extension type.
//
// If a naming conflict occurs, the type is not registered and an error is returned.
func (r *Types) RegisterExtension(xt protoreflect.ExtensionType) error ***REMOVED***
	// Under rare circumstances getting the descriptor might recursively
	// examine the registry, so fetch it before locking.
	//
	// A known case where this can happen: Fetching the TypeDescriptor for a
	// legacy ExtensionDesc can consult the global registry.
	xd := xt.TypeDescriptor()

	if r == GlobalTypes ***REMOVED***
		globalMutex.Lock()
		defer globalMutex.Unlock()
	***REMOVED***

	field := xd.Number()
	message := xd.ContainingMessage().FullName()
	if prev := r.extensionsByMessage[message][field]; prev != nil ***REMOVED***
		err := errors.New("extension number %d is already registered on message %v", field, message)
		err = amendErrorWithCaller(err, prev, xt)
		if !(r == GlobalTypes && ignoreConflict(xd, err)) ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if err := r.register("extension", xd, xt); err != nil ***REMOVED***
		return err
	***REMOVED***
	if r.extensionsByMessage == nil ***REMOVED***
		r.extensionsByMessage = make(extensionsByMessage)
	***REMOVED***
	if r.extensionsByMessage[message] == nil ***REMOVED***
		r.extensionsByMessage[message] = make(extensionsByNumber)
	***REMOVED***
	r.extensionsByMessage[message][field] = xt
	r.numExtensions++
	return nil
***REMOVED***

func (r *Types) register(kind string, desc protoreflect.Descriptor, typ interface***REMOVED******REMOVED***) error ***REMOVED***
	name := desc.FullName()
	prev := r.typesByName[name]
	if prev != nil ***REMOVED***
		err := errors.New("%v %v is already registered", kind, name)
		err = amendErrorWithCaller(err, prev, typ)
		if !(r == GlobalTypes && ignoreConflict(desc, err)) ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if r.typesByName == nil ***REMOVED***
		r.typesByName = make(typesByName)
	***REMOVED***
	r.typesByName[name] = typ
	return nil
***REMOVED***

// FindEnumByName looks up an enum by its full name.
// E.g., "google.protobuf.Field.Kind".
//
// This returns (nil, NotFound) if not found.
func (r *Types) FindEnumByName(enum protoreflect.FullName) (protoreflect.EnumType, error) ***REMOVED***
	if r == nil ***REMOVED***
		return nil, NotFound
	***REMOVED***
	if r == GlobalTypes ***REMOVED***
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	***REMOVED***
	if v := r.typesByName[enum]; v != nil ***REMOVED***
		if et, _ := v.(protoreflect.EnumType); et != nil ***REMOVED***
			return et, nil
		***REMOVED***
		return nil, errors.New("found wrong type: got %v, want enum", typeName(v))
	***REMOVED***
	return nil, NotFound
***REMOVED***

// FindMessageByName looks up a message by its full name,
// e.g. "google.protobuf.Any".
//
// This returns (nil, NotFound) if not found.
func (r *Types) FindMessageByName(message protoreflect.FullName) (protoreflect.MessageType, error) ***REMOVED***
	if r == nil ***REMOVED***
		return nil, NotFound
	***REMOVED***
	if r == GlobalTypes ***REMOVED***
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	***REMOVED***
	if v := r.typesByName[message]; v != nil ***REMOVED***
		if mt, _ := v.(protoreflect.MessageType); mt != nil ***REMOVED***
			return mt, nil
		***REMOVED***
		return nil, errors.New("found wrong type: got %v, want message", typeName(v))
	***REMOVED***
	return nil, NotFound
***REMOVED***

// FindMessageByURL looks up a message by a URL identifier.
// See documentation on google.protobuf.Any.type_url for the URL format.
//
// This returns (nil, NotFound) if not found.
func (r *Types) FindMessageByURL(url string) (protoreflect.MessageType, error) ***REMOVED***
	// This function is similar to FindMessageByName but
	// truncates anything before and including '/' in the URL.
	if r == nil ***REMOVED***
		return nil, NotFound
	***REMOVED***
	if r == GlobalTypes ***REMOVED***
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	***REMOVED***
	message := protoreflect.FullName(url)
	if i := strings.LastIndexByte(url, '/'); i >= 0 ***REMOVED***
		message = message[i+len("/"):]
	***REMOVED***

	if v := r.typesByName[message]; v != nil ***REMOVED***
		if mt, _ := v.(protoreflect.MessageType); mt != nil ***REMOVED***
			return mt, nil
		***REMOVED***
		return nil, errors.New("found wrong type: got %v, want message", typeName(v))
	***REMOVED***
	return nil, NotFound
***REMOVED***

// FindExtensionByName looks up a extension field by the field's full name.
// Note that this is the full name of the field as determined by
// where the extension is declared and is unrelated to the full name of the
// message being extended.
//
// This returns (nil, NotFound) if not found.
func (r *Types) FindExtensionByName(field protoreflect.FullName) (protoreflect.ExtensionType, error) ***REMOVED***
	if r == nil ***REMOVED***
		return nil, NotFound
	***REMOVED***
	if r == GlobalTypes ***REMOVED***
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	***REMOVED***
	if v := r.typesByName[field]; v != nil ***REMOVED***
		if xt, _ := v.(protoreflect.ExtensionType); xt != nil ***REMOVED***
			return xt, nil
		***REMOVED***

		// MessageSet extensions are special in that the name of the extension
		// is the name of the message type used to extend the MessageSet.
		// This naming scheme is used by text and JSON serialization.
		//
		// This feature is protected by the ProtoLegacy flag since MessageSets
		// are a proto1 feature that is long deprecated.
		if flags.ProtoLegacy ***REMOVED***
			if _, ok := v.(protoreflect.MessageType); ok ***REMOVED***
				field := field.Append(messageset.ExtensionName)
				if v := r.typesByName[field]; v != nil ***REMOVED***
					if xt, _ := v.(protoreflect.ExtensionType); xt != nil ***REMOVED***
						if messageset.IsMessageSetExtension(xt.TypeDescriptor()) ***REMOVED***
							return xt, nil
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***

		return nil, errors.New("found wrong type: got %v, want extension", typeName(v))
	***REMOVED***
	return nil, NotFound
***REMOVED***

// FindExtensionByNumber looks up a extension field by the field number
// within some parent message, identified by full name.
//
// This returns (nil, NotFound) if not found.
func (r *Types) FindExtensionByNumber(message protoreflect.FullName, field protoreflect.FieldNumber) (protoreflect.ExtensionType, error) ***REMOVED***
	if r == nil ***REMOVED***
		return nil, NotFound
	***REMOVED***
	if r == GlobalTypes ***REMOVED***
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	***REMOVED***
	if xt, ok := r.extensionsByMessage[message][field]; ok ***REMOVED***
		return xt, nil
	***REMOVED***
	return nil, NotFound
***REMOVED***

// NumEnums reports the number of registered enums.
func (r *Types) NumEnums() int ***REMOVED***
	if r == nil ***REMOVED***
		return 0
	***REMOVED***
	if r == GlobalTypes ***REMOVED***
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	***REMOVED***
	return r.numEnums
***REMOVED***

// RangeEnums iterates over all registered enums while f returns true.
// Iteration order is undefined.
func (r *Types) RangeEnums(f func(protoreflect.EnumType) bool) ***REMOVED***
	if r == nil ***REMOVED***
		return
	***REMOVED***
	if r == GlobalTypes ***REMOVED***
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	***REMOVED***
	for _, typ := range r.typesByName ***REMOVED***
		if et, ok := typ.(protoreflect.EnumType); ok ***REMOVED***
			if !f(et) ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// NumMessages reports the number of registered messages.
func (r *Types) NumMessages() int ***REMOVED***
	if r == nil ***REMOVED***
		return 0
	***REMOVED***
	if r == GlobalTypes ***REMOVED***
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	***REMOVED***
	return r.numMessages
***REMOVED***

// RangeMessages iterates over all registered messages while f returns true.
// Iteration order is undefined.
func (r *Types) RangeMessages(f func(protoreflect.MessageType) bool) ***REMOVED***
	if r == nil ***REMOVED***
		return
	***REMOVED***
	if r == GlobalTypes ***REMOVED***
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	***REMOVED***
	for _, typ := range r.typesByName ***REMOVED***
		if mt, ok := typ.(protoreflect.MessageType); ok ***REMOVED***
			if !f(mt) ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// NumExtensions reports the number of registered extensions.
func (r *Types) NumExtensions() int ***REMOVED***
	if r == nil ***REMOVED***
		return 0
	***REMOVED***
	if r == GlobalTypes ***REMOVED***
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	***REMOVED***
	return r.numExtensions
***REMOVED***

// RangeExtensions iterates over all registered extensions while f returns true.
// Iteration order is undefined.
func (r *Types) RangeExtensions(f func(protoreflect.ExtensionType) bool) ***REMOVED***
	if r == nil ***REMOVED***
		return
	***REMOVED***
	if r == GlobalTypes ***REMOVED***
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	***REMOVED***
	for _, typ := range r.typesByName ***REMOVED***
		if xt, ok := typ.(protoreflect.ExtensionType); ok ***REMOVED***
			if !f(xt) ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// NumExtensionsByMessage reports the number of registered extensions for
// a given message type.
func (r *Types) NumExtensionsByMessage(message protoreflect.FullName) int ***REMOVED***
	if r == nil ***REMOVED***
		return 0
	***REMOVED***
	if r == GlobalTypes ***REMOVED***
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	***REMOVED***
	return len(r.extensionsByMessage[message])
***REMOVED***

// RangeExtensionsByMessage iterates over all registered extensions filtered
// by a given message type while f returns true. Iteration order is undefined.
func (r *Types) RangeExtensionsByMessage(message protoreflect.FullName, f func(protoreflect.ExtensionType) bool) ***REMOVED***
	if r == nil ***REMOVED***
		return
	***REMOVED***
	if r == GlobalTypes ***REMOVED***
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	***REMOVED***
	for _, xt := range r.extensionsByMessage[message] ***REMOVED***
		if !f(xt) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func typeName(t interface***REMOVED******REMOVED***) string ***REMOVED***
	switch t.(type) ***REMOVED***
	case protoreflect.EnumType:
		return "enum"
	case protoreflect.MessageType:
		return "message"
	case protoreflect.ExtensionType:
		return "extension"
	default:
		return fmt.Sprintf("%T", t)
	***REMOVED***
***REMOVED***

func amendErrorWithCaller(err error, prev, curr interface***REMOVED******REMOVED***) error ***REMOVED***
	prevPkg := goPackage(prev)
	currPkg := goPackage(curr)
	if prevPkg == "" || currPkg == "" || prevPkg == currPkg ***REMOVED***
		return err
	***REMOVED***
	return errors.New("%s\n\tpreviously from: %q\n\tcurrently from:  %q", err, prevPkg, currPkg)
***REMOVED***

func goPackage(v interface***REMOVED******REMOVED***) string ***REMOVED***
	switch d := v.(type) ***REMOVED***
	case protoreflect.EnumType:
		v = d.Descriptor()
	case protoreflect.MessageType:
		v = d.Descriptor()
	case protoreflect.ExtensionType:
		v = d.TypeDescriptor()
	***REMOVED***
	if d, ok := v.(protoreflect.Descriptor); ok ***REMOVED***
		v = d.ParentFile()
	***REMOVED***
	if d, ok := v.(interface***REMOVED*** GoPackagePath() string ***REMOVED***); ok ***REMOVED***
		return d.GoPackagePath()
	***REMOVED***
	return ""
***REMOVED***
