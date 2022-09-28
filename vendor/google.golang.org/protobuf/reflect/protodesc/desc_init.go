// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protodesc

import (
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/filedesc"
	"google.golang.org/protobuf/internal/strs"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"google.golang.org/protobuf/types/descriptorpb"
)

type descsByName map[protoreflect.FullName]protoreflect.Descriptor

func (r descsByName) initEnumDeclarations(eds []*descriptorpb.EnumDescriptorProto, parent protoreflect.Descriptor, sb *strs.Builder) (es []filedesc.Enum, err error) ***REMOVED***
	es = make([]filedesc.Enum, len(eds)) // allocate up-front to ensure stable pointers
	for i, ed := range eds ***REMOVED***
		e := &es[i]
		e.L2 = new(filedesc.EnumL2)
		if e.L0, err = r.makeBase(e, parent, ed.GetName(), i, sb); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if opts := ed.GetOptions(); opts != nil ***REMOVED***
			opts = proto.Clone(opts).(*descriptorpb.EnumOptions)
			e.L2.Options = func() protoreflect.ProtoMessage ***REMOVED*** return opts ***REMOVED***
		***REMOVED***
		for _, s := range ed.GetReservedName() ***REMOVED***
			e.L2.ReservedNames.List = append(e.L2.ReservedNames.List, protoreflect.Name(s))
		***REMOVED***
		for _, rr := range ed.GetReservedRange() ***REMOVED***
			e.L2.ReservedRanges.List = append(e.L2.ReservedRanges.List, [2]protoreflect.EnumNumber***REMOVED***
				protoreflect.EnumNumber(rr.GetStart()),
				protoreflect.EnumNumber(rr.GetEnd()),
			***REMOVED***)
		***REMOVED***
		if e.L2.Values.List, err = r.initEnumValuesFromDescriptorProto(ed.GetValue(), e, sb); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return es, nil
***REMOVED***

func (r descsByName) initEnumValuesFromDescriptorProto(vds []*descriptorpb.EnumValueDescriptorProto, parent protoreflect.Descriptor, sb *strs.Builder) (vs []filedesc.EnumValue, err error) ***REMOVED***
	vs = make([]filedesc.EnumValue, len(vds)) // allocate up-front to ensure stable pointers
	for i, vd := range vds ***REMOVED***
		v := &vs[i]
		if v.L0, err = r.makeBase(v, parent, vd.GetName(), i, sb); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if opts := vd.GetOptions(); opts != nil ***REMOVED***
			opts = proto.Clone(opts).(*descriptorpb.EnumValueOptions)
			v.L1.Options = func() protoreflect.ProtoMessage ***REMOVED*** return opts ***REMOVED***
		***REMOVED***
		v.L1.Number = protoreflect.EnumNumber(vd.GetNumber())
	***REMOVED***
	return vs, nil
***REMOVED***

func (r descsByName) initMessagesDeclarations(mds []*descriptorpb.DescriptorProto, parent protoreflect.Descriptor, sb *strs.Builder) (ms []filedesc.Message, err error) ***REMOVED***
	ms = make([]filedesc.Message, len(mds)) // allocate up-front to ensure stable pointers
	for i, md := range mds ***REMOVED***
		m := &ms[i]
		m.L2 = new(filedesc.MessageL2)
		if m.L0, err = r.makeBase(m, parent, md.GetName(), i, sb); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if opts := md.GetOptions(); opts != nil ***REMOVED***
			opts = proto.Clone(opts).(*descriptorpb.MessageOptions)
			m.L2.Options = func() protoreflect.ProtoMessage ***REMOVED*** return opts ***REMOVED***
			m.L1.IsMapEntry = opts.GetMapEntry()
			m.L1.IsMessageSet = opts.GetMessageSetWireFormat()
		***REMOVED***
		for _, s := range md.GetReservedName() ***REMOVED***
			m.L2.ReservedNames.List = append(m.L2.ReservedNames.List, protoreflect.Name(s))
		***REMOVED***
		for _, rr := range md.GetReservedRange() ***REMOVED***
			m.L2.ReservedRanges.List = append(m.L2.ReservedRanges.List, [2]protoreflect.FieldNumber***REMOVED***
				protoreflect.FieldNumber(rr.GetStart()),
				protoreflect.FieldNumber(rr.GetEnd()),
			***REMOVED***)
		***REMOVED***
		for _, xr := range md.GetExtensionRange() ***REMOVED***
			m.L2.ExtensionRanges.List = append(m.L2.ExtensionRanges.List, [2]protoreflect.FieldNumber***REMOVED***
				protoreflect.FieldNumber(xr.GetStart()),
				protoreflect.FieldNumber(xr.GetEnd()),
			***REMOVED***)
			var optsFunc func() protoreflect.ProtoMessage
			if opts := xr.GetOptions(); opts != nil ***REMOVED***
				opts = proto.Clone(opts).(*descriptorpb.ExtensionRangeOptions)
				optsFunc = func() protoreflect.ProtoMessage ***REMOVED*** return opts ***REMOVED***
			***REMOVED***
			m.L2.ExtensionRangeOptions = append(m.L2.ExtensionRangeOptions, optsFunc)
		***REMOVED***
		if m.L2.Fields.List, err = r.initFieldsFromDescriptorProto(md.GetField(), m, sb); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if m.L2.Oneofs.List, err = r.initOneofsFromDescriptorProto(md.GetOneofDecl(), m, sb); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if m.L1.Enums.List, err = r.initEnumDeclarations(md.GetEnumType(), m, sb); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if m.L1.Messages.List, err = r.initMessagesDeclarations(md.GetNestedType(), m, sb); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if m.L1.Extensions.List, err = r.initExtensionDeclarations(md.GetExtension(), m, sb); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return ms, nil
***REMOVED***

func (r descsByName) initFieldsFromDescriptorProto(fds []*descriptorpb.FieldDescriptorProto, parent protoreflect.Descriptor, sb *strs.Builder) (fs []filedesc.Field, err error) ***REMOVED***
	fs = make([]filedesc.Field, len(fds)) // allocate up-front to ensure stable pointers
	for i, fd := range fds ***REMOVED***
		f := &fs[i]
		if f.L0, err = r.makeBase(f, parent, fd.GetName(), i, sb); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		f.L1.IsProto3Optional = fd.GetProto3Optional()
		if opts := fd.GetOptions(); opts != nil ***REMOVED***
			opts = proto.Clone(opts).(*descriptorpb.FieldOptions)
			f.L1.Options = func() protoreflect.ProtoMessage ***REMOVED*** return opts ***REMOVED***
			f.L1.IsWeak = opts.GetWeak()
			f.L1.HasPacked = opts.Packed != nil
			f.L1.IsPacked = opts.GetPacked()
		***REMOVED***
		f.L1.Number = protoreflect.FieldNumber(fd.GetNumber())
		f.L1.Cardinality = protoreflect.Cardinality(fd.GetLabel())
		if fd.Type != nil ***REMOVED***
			f.L1.Kind = protoreflect.Kind(fd.GetType())
		***REMOVED***
		if fd.JsonName != nil ***REMOVED***
			f.L1.StringName.InitJSON(fd.GetJsonName())
		***REMOVED***
	***REMOVED***
	return fs, nil
***REMOVED***

func (r descsByName) initOneofsFromDescriptorProto(ods []*descriptorpb.OneofDescriptorProto, parent protoreflect.Descriptor, sb *strs.Builder) (os []filedesc.Oneof, err error) ***REMOVED***
	os = make([]filedesc.Oneof, len(ods)) // allocate up-front to ensure stable pointers
	for i, od := range ods ***REMOVED***
		o := &os[i]
		if o.L0, err = r.makeBase(o, parent, od.GetName(), i, sb); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if opts := od.GetOptions(); opts != nil ***REMOVED***
			opts = proto.Clone(opts).(*descriptorpb.OneofOptions)
			o.L1.Options = func() protoreflect.ProtoMessage ***REMOVED*** return opts ***REMOVED***
		***REMOVED***
	***REMOVED***
	return os, nil
***REMOVED***

func (r descsByName) initExtensionDeclarations(xds []*descriptorpb.FieldDescriptorProto, parent protoreflect.Descriptor, sb *strs.Builder) (xs []filedesc.Extension, err error) ***REMOVED***
	xs = make([]filedesc.Extension, len(xds)) // allocate up-front to ensure stable pointers
	for i, xd := range xds ***REMOVED***
		x := &xs[i]
		x.L2 = new(filedesc.ExtensionL2)
		if x.L0, err = r.makeBase(x, parent, xd.GetName(), i, sb); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if opts := xd.GetOptions(); opts != nil ***REMOVED***
			opts = proto.Clone(opts).(*descriptorpb.FieldOptions)
			x.L2.Options = func() protoreflect.ProtoMessage ***REMOVED*** return opts ***REMOVED***
			x.L2.IsPacked = opts.GetPacked()
		***REMOVED***
		x.L1.Number = protoreflect.FieldNumber(xd.GetNumber())
		x.L1.Cardinality = protoreflect.Cardinality(xd.GetLabel())
		if xd.Type != nil ***REMOVED***
			x.L1.Kind = protoreflect.Kind(xd.GetType())
		***REMOVED***
		if xd.JsonName != nil ***REMOVED***
			x.L2.StringName.InitJSON(xd.GetJsonName())
		***REMOVED***
	***REMOVED***
	return xs, nil
***REMOVED***

func (r descsByName) initServiceDeclarations(sds []*descriptorpb.ServiceDescriptorProto, parent protoreflect.Descriptor, sb *strs.Builder) (ss []filedesc.Service, err error) ***REMOVED***
	ss = make([]filedesc.Service, len(sds)) // allocate up-front to ensure stable pointers
	for i, sd := range sds ***REMOVED***
		s := &ss[i]
		s.L2 = new(filedesc.ServiceL2)
		if s.L0, err = r.makeBase(s, parent, sd.GetName(), i, sb); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if opts := sd.GetOptions(); opts != nil ***REMOVED***
			opts = proto.Clone(opts).(*descriptorpb.ServiceOptions)
			s.L2.Options = func() protoreflect.ProtoMessage ***REMOVED*** return opts ***REMOVED***
		***REMOVED***
		if s.L2.Methods.List, err = r.initMethodsFromDescriptorProto(sd.GetMethod(), s, sb); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return ss, nil
***REMOVED***

func (r descsByName) initMethodsFromDescriptorProto(mds []*descriptorpb.MethodDescriptorProto, parent protoreflect.Descriptor, sb *strs.Builder) (ms []filedesc.Method, err error) ***REMOVED***
	ms = make([]filedesc.Method, len(mds)) // allocate up-front to ensure stable pointers
	for i, md := range mds ***REMOVED***
		m := &ms[i]
		if m.L0, err = r.makeBase(m, parent, md.GetName(), i, sb); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if opts := md.GetOptions(); opts != nil ***REMOVED***
			opts = proto.Clone(opts).(*descriptorpb.MethodOptions)
			m.L1.Options = func() protoreflect.ProtoMessage ***REMOVED*** return opts ***REMOVED***
		***REMOVED***
		m.L1.IsStreamingClient = md.GetClientStreaming()
		m.L1.IsStreamingServer = md.GetServerStreaming()
	***REMOVED***
	return ms, nil
***REMOVED***

func (r descsByName) makeBase(child, parent protoreflect.Descriptor, name string, idx int, sb *strs.Builder) (filedesc.BaseL0, error) ***REMOVED***
	if !protoreflect.Name(name).IsValid() ***REMOVED***
		return filedesc.BaseL0***REMOVED******REMOVED***, errors.New("descriptor %q has an invalid nested name: %q", parent.FullName(), name)
	***REMOVED***

	// Derive the full name of the child.
	// Note that enum values are a sibling to the enum parent in the namespace.
	var fullName protoreflect.FullName
	if _, ok := parent.(protoreflect.EnumDescriptor); ok ***REMOVED***
		fullName = sb.AppendFullName(parent.FullName().Parent(), protoreflect.Name(name))
	***REMOVED*** else ***REMOVED***
		fullName = sb.AppendFullName(parent.FullName(), protoreflect.Name(name))
	***REMOVED***
	if _, ok := r[fullName]; ok ***REMOVED***
		return filedesc.BaseL0***REMOVED******REMOVED***, errors.New("descriptor %q already declared", fullName)
	***REMOVED***
	r[fullName] = child

	// TODO: Verify that the full name does not already exist in the resolver?
	// This is not as critical since most usages of NewFile will register
	// the created file back into the registry, which will perform this check.

	return filedesc.BaseL0***REMOVED***
		FullName:   fullName,
		ParentFile: parent.ParentFile().(*filedesc.File),
		Parent:     parent,
		Index:      idx,
	***REMOVED***, nil
***REMOVED***
