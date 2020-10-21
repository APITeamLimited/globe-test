// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protodesc

import (
	"google.golang.org/protobuf/internal/encoding/defval"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/filedesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	"google.golang.org/protobuf/types/descriptorpb"
)

// resolver is a wrapper around a local registry of declarations within the file
// and the remote resolver. The remote resolver is restricted to only return
// descriptors that have been imported.
type resolver struct ***REMOVED***
	local   descsByName
	remote  Resolver
	imports importSet

	allowUnresolvable bool
***REMOVED***

func (r *resolver) resolveMessageDependencies(ms []filedesc.Message, mds []*descriptorpb.DescriptorProto) (err error) ***REMOVED***
	for i, md := range mds ***REMOVED***
		m := &ms[i]
		for j, fd := range md.GetField() ***REMOVED***
			f := &m.L2.Fields.List[j]
			if f.L1.Cardinality == protoreflect.Required ***REMOVED***
				m.L2.RequiredNumbers.List = append(m.L2.RequiredNumbers.List, f.L1.Number)
			***REMOVED***
			if fd.OneofIndex != nil ***REMOVED***
				k := int(fd.GetOneofIndex())
				if !(0 <= k && k < len(md.GetOneofDecl())) ***REMOVED***
					return errors.New("message field %q has an invalid oneof index: %d", f.FullName(), k)
				***REMOVED***
				o := &m.L2.Oneofs.List[k]
				f.L1.ContainingOneof = o
				o.L1.Fields.List = append(o.L1.Fields.List, f)
			***REMOVED***

			if f.L1.Kind, f.L1.Enum, f.L1.Message, err = r.findTarget(f.Kind(), f.Parent().FullName(), partialName(fd.GetTypeName()), f.IsWeak()); err != nil ***REMOVED***
				return errors.New("message field %q cannot resolve type: %v", f.FullName(), err)
			***REMOVED***
			if fd.DefaultValue != nil ***REMOVED***
				v, ev, err := unmarshalDefault(fd.GetDefaultValue(), f, r.allowUnresolvable)
				if err != nil ***REMOVED***
					return errors.New("message field %q has invalid default: %v", f.FullName(), err)
				***REMOVED***
				f.L1.Default = filedesc.DefaultValue(v, ev)
			***REMOVED***
		***REMOVED***

		if err := r.resolveMessageDependencies(m.L1.Messages.List, md.GetNestedType()); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := r.resolveExtensionDependencies(m.L1.Extensions.List, md.GetExtension()); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (r *resolver) resolveExtensionDependencies(xs []filedesc.Extension, xds []*descriptorpb.FieldDescriptorProto) (err error) ***REMOVED***
	for i, xd := range xds ***REMOVED***
		x := &xs[i]
		if x.L1.Extendee, err = r.findMessageDescriptor(x.Parent().FullName(), partialName(xd.GetExtendee()), false); err != nil ***REMOVED***
			return errors.New("extension field %q cannot resolve extendee: %v", x.FullName(), err)
		***REMOVED***
		if x.L1.Kind, x.L2.Enum, x.L2.Message, err = r.findTarget(x.Kind(), x.Parent().FullName(), partialName(xd.GetTypeName()), false); err != nil ***REMOVED***
			return errors.New("extension field %q cannot resolve type: %v", x.FullName(), err)
		***REMOVED***
		if xd.DefaultValue != nil ***REMOVED***
			v, ev, err := unmarshalDefault(xd.GetDefaultValue(), x, r.allowUnresolvable)
			if err != nil ***REMOVED***
				return errors.New("extension field %q has invalid default: %v", x.FullName(), err)
			***REMOVED***
			x.L2.Default = filedesc.DefaultValue(v, ev)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (r *resolver) resolveServiceDependencies(ss []filedesc.Service, sds []*descriptorpb.ServiceDescriptorProto) (err error) ***REMOVED***
	for i, sd := range sds ***REMOVED***
		s := &ss[i]
		for j, md := range sd.GetMethod() ***REMOVED***
			m := &s.L2.Methods.List[j]
			m.L1.Input, err = r.findMessageDescriptor(m.Parent().FullName(), partialName(md.GetInputType()), false)
			if err != nil ***REMOVED***
				return errors.New("service method %q cannot resolve input: %v", m.FullName(), err)
			***REMOVED***
			m.L1.Output, err = r.findMessageDescriptor(s.FullName(), partialName(md.GetOutputType()), false)
			if err != nil ***REMOVED***
				return errors.New("service method %q cannot resolve output: %v", m.FullName(), err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// findTarget finds an enum or message descriptor if k is an enum, message,
// group, or unknown. If unknown, and the name could be resolved, the kind
// returned kind is set based on the type of the resolved descriptor.
func (r *resolver) findTarget(k protoreflect.Kind, scope protoreflect.FullName, ref partialName, isWeak bool) (protoreflect.Kind, protoreflect.EnumDescriptor, protoreflect.MessageDescriptor, error) ***REMOVED***
	switch k ***REMOVED***
	case protoreflect.EnumKind:
		ed, err := r.findEnumDescriptor(scope, ref, isWeak)
		if err != nil ***REMOVED***
			return 0, nil, nil, err
		***REMOVED***
		return k, ed, nil, nil
	case protoreflect.MessageKind, protoreflect.GroupKind:
		md, err := r.findMessageDescriptor(scope, ref, isWeak)
		if err != nil ***REMOVED***
			return 0, nil, nil, err
		***REMOVED***
		return k, nil, md, nil
	case 0:
		// Handle unspecified kinds (possible with parsers that operate
		// on a per-file basis without knowledge of dependencies).
		d, err := r.findDescriptor(scope, ref)
		if err == protoregistry.NotFound && (r.allowUnresolvable || isWeak) ***REMOVED***
			return k, filedesc.PlaceholderEnum(ref.FullName()), filedesc.PlaceholderMessage(ref.FullName()), nil
		***REMOVED*** else if err == protoregistry.NotFound ***REMOVED***
			return 0, nil, nil, errors.New("%q not found", ref.FullName())
		***REMOVED*** else if err != nil ***REMOVED***
			return 0, nil, nil, err
		***REMOVED***
		switch d := d.(type) ***REMOVED***
		case protoreflect.EnumDescriptor:
			return protoreflect.EnumKind, d, nil, nil
		case protoreflect.MessageDescriptor:
			return protoreflect.MessageKind, nil, d, nil
		default:
			return 0, nil, nil, errors.New("unknown kind")
		***REMOVED***
	default:
		if ref != "" ***REMOVED***
			return 0, nil, nil, errors.New("target name cannot be specified for %v", k)
		***REMOVED***
		if !k.IsValid() ***REMOVED***
			return 0, nil, nil, errors.New("invalid kind: %d", k)
		***REMOVED***
		return k, nil, nil, nil
	***REMOVED***
***REMOVED***

// findDescriptor finds the descriptor by name,
// which may be a relative name within some scope.
//
// Suppose the scope was "fizz.buzz" and the reference was "Foo.Bar",
// then the following full names are searched:
//	* fizz.buzz.Foo.Bar
//	* fizz.Foo.Bar
//	* Foo.Bar
func (r *resolver) findDescriptor(scope protoreflect.FullName, ref partialName) (protoreflect.Descriptor, error) ***REMOVED***
	if !ref.IsValid() ***REMOVED***
		return nil, errors.New("invalid name reference: %q", ref)
	***REMOVED***
	if ref.IsFull() ***REMOVED***
		scope, ref = "", ref[1:]
	***REMOVED***
	var foundButNotImported protoreflect.Descriptor
	for ***REMOVED***
		// Derive the full name to search.
		s := protoreflect.FullName(ref)
		if scope != "" ***REMOVED***
			s = scope + "." + s
		***REMOVED***

		// Check the current file for the descriptor.
		if d, ok := r.local[s]; ok ***REMOVED***
			return d, nil
		***REMOVED***

		// Check the remote registry for the descriptor.
		d, err := r.remote.FindDescriptorByName(s)
		if err == nil ***REMOVED***
			// Only allow descriptors covered by one of the imports.
			if r.imports[d.ParentFile().Path()] ***REMOVED***
				return d, nil
			***REMOVED***
			foundButNotImported = d
		***REMOVED*** else if err != protoregistry.NotFound ***REMOVED***
			return nil, errors.Wrap(err, "%q", s)
		***REMOVED***

		// Continue on at a higher level of scoping.
		if scope == "" ***REMOVED***
			if d := foundButNotImported; d != nil ***REMOVED***
				return nil, errors.New("resolved %q, but %q is not imported", d.FullName(), d.ParentFile().Path())
			***REMOVED***
			return nil, protoregistry.NotFound
		***REMOVED***
		scope = scope.Parent()
	***REMOVED***
***REMOVED***

func (r *resolver) findEnumDescriptor(scope protoreflect.FullName, ref partialName, isWeak bool) (protoreflect.EnumDescriptor, error) ***REMOVED***
	d, err := r.findDescriptor(scope, ref)
	if err == protoregistry.NotFound && (r.allowUnresolvable || isWeak) ***REMOVED***
		return filedesc.PlaceholderEnum(ref.FullName()), nil
	***REMOVED*** else if err == protoregistry.NotFound ***REMOVED***
		return nil, errors.New("%q not found", ref.FullName())
	***REMOVED*** else if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ed, ok := d.(protoreflect.EnumDescriptor)
	if !ok ***REMOVED***
		return nil, errors.New("resolved %q, but it is not an enum", d.FullName())
	***REMOVED***
	return ed, nil
***REMOVED***

func (r *resolver) findMessageDescriptor(scope protoreflect.FullName, ref partialName, isWeak bool) (protoreflect.MessageDescriptor, error) ***REMOVED***
	d, err := r.findDescriptor(scope, ref)
	if err == protoregistry.NotFound && (r.allowUnresolvable || isWeak) ***REMOVED***
		return filedesc.PlaceholderMessage(ref.FullName()), nil
	***REMOVED*** else if err == protoregistry.NotFound ***REMOVED***
		return nil, errors.New("%q not found", ref.FullName())
	***REMOVED*** else if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	md, ok := d.(protoreflect.MessageDescriptor)
	if !ok ***REMOVED***
		return nil, errors.New("resolved %q, but it is not an message", d.FullName())
	***REMOVED***
	return md, nil
***REMOVED***

// partialName is the partial name. A leading dot means that the name is full,
// otherwise the name is relative to some current scope.
// See google.protobuf.FieldDescriptorProto.type_name.
type partialName string

func (s partialName) IsFull() bool ***REMOVED***
	return len(s) > 0 && s[0] == '.'
***REMOVED***

func (s partialName) IsValid() bool ***REMOVED***
	if s.IsFull() ***REMOVED***
		return protoreflect.FullName(s[1:]).IsValid()
	***REMOVED***
	return protoreflect.FullName(s).IsValid()
***REMOVED***

const unknownPrefix = "*."

// FullName converts the partial name to a full name on a best-effort basis.
// If relative, it creates an invalid full name, using a "*." prefix
// to indicate that the start of the full name is unknown.
func (s partialName) FullName() protoreflect.FullName ***REMOVED***
	if s.IsFull() ***REMOVED***
		return protoreflect.FullName(s[1:])
	***REMOVED***
	return protoreflect.FullName(unknownPrefix + s)
***REMOVED***

func unmarshalDefault(s string, fd protoreflect.FieldDescriptor, allowUnresolvable bool) (protoreflect.Value, protoreflect.EnumValueDescriptor, error) ***REMOVED***
	var evs protoreflect.EnumValueDescriptors
	if fd.Enum() != nil ***REMOVED***
		evs = fd.Enum().Values()
	***REMOVED***
	v, ev, err := defval.Unmarshal(s, fd.Kind(), evs, defval.Descriptor)
	if err != nil && allowUnresolvable && evs != nil && protoreflect.Name(s).IsValid() ***REMOVED***
		v = protoreflect.ValueOfEnum(0)
		if evs.Len() > 0 ***REMOVED***
			v = protoreflect.ValueOfEnum(evs.Get(0).Number())
		***REMOVED***
		ev = filedesc.PlaceholderEnumValue(fd.Enum().FullName().Parent().Append(protoreflect.Name(s)))
	***REMOVED*** else if err != nil ***REMOVED***
		return v, ev, err
	***REMOVED***
	if fd.Syntax() == protoreflect.Proto3 ***REMOVED***
		return v, ev, errors.New("cannot be specified under proto3 semantics")
	***REMOVED***
	if fd.Kind() == protoreflect.MessageKind || fd.Kind() == protoreflect.GroupKind || fd.Cardinality() == protoreflect.Repeated ***REMOVED***
		return v, ev, errors.New("cannot be specified on composite types")
	***REMOVED***
	return v, ev, nil
***REMOVED***
