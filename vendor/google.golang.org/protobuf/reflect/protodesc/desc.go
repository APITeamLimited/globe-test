// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package protodesc provides functionality for converting
// FileDescriptorProto messages to/from protoreflect.FileDescriptor values.
//
// The google.protobuf.FileDescriptorProto is a protobuf message that describes
// the type information for a .proto file in a form that is easily serializable.
// The protoreflect.FileDescriptor is a more structured representation of
// the FileDescriptorProto message where references and remote dependencies
// can be directly followed.
package protodesc

import (
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/filedesc"
	"google.golang.org/protobuf/internal/pragma"
	"google.golang.org/protobuf/internal/strs"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	"google.golang.org/protobuf/types/descriptorpb"
)

// Resolver is the resolver used by NewFile to resolve dependencies.
// The enums and messages provided must belong to some parent file,
// which is also registered.
//
// It is implemented by protoregistry.Files.
type Resolver interface ***REMOVED***
	FindFileByPath(string) (protoreflect.FileDescriptor, error)
	FindDescriptorByName(protoreflect.FullName) (protoreflect.Descriptor, error)
***REMOVED***

// FileOptions configures the construction of file descriptors.
type FileOptions struct ***REMOVED***
	pragma.NoUnkeyedLiterals

	// AllowUnresolvable configures New to permissively allow unresolvable
	// file, enum, or message dependencies. Unresolved dependencies are replaced
	// by placeholder equivalents.
	//
	// The following dependencies may be left unresolved:
	//	• Resolving an imported file.
	//	• Resolving the type for a message field or extension field.
	//	If the kind of the field is unknown, then a placeholder is used for both
	//	the Enum and Message accessors on the protoreflect.FieldDescriptor.
	//	• Resolving an enum value set as the default for an optional enum field.
	//	If unresolvable, the protoreflect.FieldDescriptor.Default is set to the
	//	first value in the associated enum (or zero if the also enum dependency
	//	is also unresolvable). The protoreflect.FieldDescriptor.DefaultEnumValue
	//	is populated with a placeholder.
	//	• Resolving the extended message type for an extension field.
	//	• Resolving the input or output message type for a service method.
	//
	// If the unresolved dependency uses a relative name,
	// then the placeholder will contain an invalid FullName with a "*." prefix,
	// indicating that the starting prefix of the full name is unknown.
	AllowUnresolvable bool
***REMOVED***

// NewFile creates a new protoreflect.FileDescriptor from the provided
// file descriptor message. See FileOptions.New for more information.
func NewFile(fd *descriptorpb.FileDescriptorProto, r Resolver) (protoreflect.FileDescriptor, error) ***REMOVED***
	return FileOptions***REMOVED******REMOVED***.New(fd, r)
***REMOVED***

// NewFiles creates a new protoregistry.Files from the provided
// FileDescriptorSet message. See FileOptions.NewFiles for more information.
func NewFiles(fd *descriptorpb.FileDescriptorSet) (*protoregistry.Files, error) ***REMOVED***
	return FileOptions***REMOVED******REMOVED***.NewFiles(fd)
***REMOVED***

// New creates a new protoreflect.FileDescriptor from the provided
// file descriptor message. The file must represent a valid proto file according
// to protobuf semantics. The returned descriptor is a deep copy of the input.
//
// Any imported files, enum types, or message types referenced in the file are
// resolved using the provided registry. When looking up an import file path,
// the path must be unique. The newly created file descriptor is not registered
// back into the provided file registry.
func (o FileOptions) New(fd *descriptorpb.FileDescriptorProto, r Resolver) (protoreflect.FileDescriptor, error) ***REMOVED***
	if r == nil ***REMOVED***
		r = (*protoregistry.Files)(nil) // empty resolver
	***REMOVED***

	// Handle the file descriptor content.
	f := &filedesc.File***REMOVED***L2: &filedesc.FileL2***REMOVED******REMOVED******REMOVED***
	switch fd.GetSyntax() ***REMOVED***
	case "proto2", "":
		f.L1.Syntax = protoreflect.Proto2
	case "proto3":
		f.L1.Syntax = protoreflect.Proto3
	default:
		return nil, errors.New("invalid syntax: %q", fd.GetSyntax())
	***REMOVED***
	f.L1.Path = fd.GetName()
	if f.L1.Path == "" ***REMOVED***
		return nil, errors.New("file path must be populated")
	***REMOVED***
	f.L1.Package = protoreflect.FullName(fd.GetPackage())
	if !f.L1.Package.IsValid() && f.L1.Package != "" ***REMOVED***
		return nil, errors.New("invalid package: %q", f.L1.Package)
	***REMOVED***
	if opts := fd.GetOptions(); opts != nil ***REMOVED***
		opts = proto.Clone(opts).(*descriptorpb.FileOptions)
		f.L2.Options = func() protoreflect.ProtoMessage ***REMOVED*** return opts ***REMOVED***
	***REMOVED***

	f.L2.Imports = make(filedesc.FileImports, len(fd.GetDependency()))
	for _, i := range fd.GetPublicDependency() ***REMOVED***
		if !(0 <= i && int(i) < len(f.L2.Imports)) || f.L2.Imports[i].IsPublic ***REMOVED***
			return nil, errors.New("invalid or duplicate public import index: %d", i)
		***REMOVED***
		f.L2.Imports[i].IsPublic = true
	***REMOVED***
	for _, i := range fd.GetWeakDependency() ***REMOVED***
		if !(0 <= i && int(i) < len(f.L2.Imports)) || f.L2.Imports[i].IsWeak ***REMOVED***
			return nil, errors.New("invalid or duplicate weak import index: %d", i)
		***REMOVED***
		f.L2.Imports[i].IsWeak = true
	***REMOVED***
	imps := importSet***REMOVED***f.Path(): true***REMOVED***
	for i, path := range fd.GetDependency() ***REMOVED***
		imp := &f.L2.Imports[i]
		f, err := r.FindFileByPath(path)
		if err == protoregistry.NotFound && (o.AllowUnresolvable || imp.IsWeak) ***REMOVED***
			f = filedesc.PlaceholderFile(path)
		***REMOVED*** else if err != nil ***REMOVED***
			return nil, errors.New("could not resolve import %q: %v", path, err)
		***REMOVED***
		imp.FileDescriptor = f

		if imps[imp.Path()] ***REMOVED***
			return nil, errors.New("already imported %q", path)
		***REMOVED***
		imps[imp.Path()] = true
	***REMOVED***
	for i := range fd.GetDependency() ***REMOVED***
		imp := &f.L2.Imports[i]
		imps.importPublic(imp.Imports())
	***REMOVED***

	// Handle source locations.
	for _, loc := range fd.GetSourceCodeInfo().GetLocation() ***REMOVED***
		var l protoreflect.SourceLocation
		// TODO: Validate that the path points to an actual declaration?
		l.Path = protoreflect.SourcePath(loc.GetPath())
		s := loc.GetSpan()
		switch len(s) ***REMOVED***
		case 3:
			l.StartLine, l.StartColumn, l.EndLine, l.EndColumn = int(s[0]), int(s[1]), int(s[0]), int(s[2])
		case 4:
			l.StartLine, l.StartColumn, l.EndLine, l.EndColumn = int(s[0]), int(s[1]), int(s[2]), int(s[3])
		default:
			return nil, errors.New("invalid span: %v", s)
		***REMOVED***
		// TODO: Validate that the span information is sensible?
		// See https://github.com/protocolbuffers/protobuf/issues/6378.
		if false && (l.EndLine < l.StartLine || l.StartLine < 0 || l.StartColumn < 0 || l.EndColumn < 0 ||
			(l.StartLine == l.EndLine && l.EndColumn <= l.StartColumn)) ***REMOVED***
			return nil, errors.New("invalid span: %v", s)
		***REMOVED***
		l.LeadingDetachedComments = loc.GetLeadingDetachedComments()
		l.LeadingComments = loc.GetLeadingComments()
		l.TrailingComments = loc.GetTrailingComments()
		f.L2.Locations.List = append(f.L2.Locations.List, l)
	***REMOVED***

	// Step 1: Allocate and derive the names for all declarations.
	// This copies all fields from the descriptor proto except:
	//	google.protobuf.FieldDescriptorProto.type_name
	//	google.protobuf.FieldDescriptorProto.default_value
	//	google.protobuf.FieldDescriptorProto.oneof_index
	//	google.protobuf.FieldDescriptorProto.extendee
	//	google.protobuf.MethodDescriptorProto.input
	//	google.protobuf.MethodDescriptorProto.output
	var err error
	sb := new(strs.Builder)
	r1 := make(descsByName)
	if f.L1.Enums.List, err = r1.initEnumDeclarations(fd.GetEnumType(), f, sb); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if f.L1.Messages.List, err = r1.initMessagesDeclarations(fd.GetMessageType(), f, sb); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if f.L1.Extensions.List, err = r1.initExtensionDeclarations(fd.GetExtension(), f, sb); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if f.L1.Services.List, err = r1.initServiceDeclarations(fd.GetService(), f, sb); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Step 2: Resolve every dependency reference not handled by step 1.
	r2 := &resolver***REMOVED***local: r1, remote: r, imports: imps, allowUnresolvable: o.AllowUnresolvable***REMOVED***
	if err := r2.resolveMessageDependencies(f.L1.Messages.List, fd.GetMessageType()); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := r2.resolveExtensionDependencies(f.L1.Extensions.List, fd.GetExtension()); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := r2.resolveServiceDependencies(f.L1.Services.List, fd.GetService()); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Step 3: Validate every enum, message, and extension declaration.
	if err := validateEnumDeclarations(f.L1.Enums.List, fd.GetEnumType()); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := validateMessageDeclarations(f.L1.Messages.List, fd.GetMessageType()); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := validateExtensionDeclarations(f.L1.Extensions.List, fd.GetExtension()); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return f, nil
***REMOVED***

type importSet map[string]bool

func (is importSet) importPublic(imps protoreflect.FileImports) ***REMOVED***
	for i := 0; i < imps.Len(); i++ ***REMOVED***
		if imp := imps.Get(i); imp.IsPublic ***REMOVED***
			is[imp.Path()] = true
			is.importPublic(imp.Imports())
		***REMOVED***
	***REMOVED***
***REMOVED***

// NewFiles creates a new protoregistry.Files from the provided
// FileDescriptorSet message. The descriptor set must include only
// valid files according to protobuf semantics. The returned descriptors
// are a deep copy of the input.
func (o FileOptions) NewFiles(fds *descriptorpb.FileDescriptorSet) (*protoregistry.Files, error) ***REMOVED***
	files := make(map[string]*descriptorpb.FileDescriptorProto)
	for _, fd := range fds.File ***REMOVED***
		if _, ok := files[fd.GetName()]; ok ***REMOVED***
			return nil, errors.New("file appears multiple times: %q", fd.GetName())
		***REMOVED***
		files[fd.GetName()] = fd
	***REMOVED***
	r := &protoregistry.Files***REMOVED******REMOVED***
	for _, fd := range files ***REMOVED***
		if err := o.addFileDeps(r, fd, files); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return r, nil
***REMOVED***
func (o FileOptions) addFileDeps(r *protoregistry.Files, fd *descriptorpb.FileDescriptorProto, files map[string]*descriptorpb.FileDescriptorProto) error ***REMOVED***
	// Set the entry to nil while descending into a file's dependencies to detect cycles.
	files[fd.GetName()] = nil
	for _, dep := range fd.Dependency ***REMOVED***
		depfd, ok := files[dep]
		if depfd == nil ***REMOVED***
			if ok ***REMOVED***
				return errors.New("import cycle in file: %q", dep)
			***REMOVED***
			continue
		***REMOVED***
		if err := o.addFileDeps(r, depfd, files); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	// Delete the entry once dependencies are processed.
	delete(files, fd.GetName())
	f, err := o.New(fd, r)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return r.RegisterFile(f)
***REMOVED***
