// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protodesc

import (
	"strings"
	"unicode"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/filedesc"
	"google.golang.org/protobuf/internal/flags"
	"google.golang.org/protobuf/internal/strs"
	"google.golang.org/protobuf/reflect/protoreflect"

	"google.golang.org/protobuf/types/descriptorpb"
)

func validateEnumDeclarations(es []filedesc.Enum, eds []*descriptorpb.EnumDescriptorProto) error ***REMOVED***
	for i, ed := range eds ***REMOVED***
		e := &es[i]
		if err := e.L2.ReservedNames.CheckValid(); err != nil ***REMOVED***
			return errors.New("enum %q reserved names has %v", e.FullName(), err)
		***REMOVED***
		if err := e.L2.ReservedRanges.CheckValid(); err != nil ***REMOVED***
			return errors.New("enum %q reserved ranges has %v", e.FullName(), err)
		***REMOVED***
		if len(ed.GetValue()) == 0 ***REMOVED***
			return errors.New("enum %q must contain at least one value declaration", e.FullName())
		***REMOVED***
		allowAlias := ed.GetOptions().GetAllowAlias()
		foundAlias := false
		for i := 0; i < e.Values().Len(); i++ ***REMOVED***
			v1 := e.Values().Get(i)
			if v2 := e.Values().ByNumber(v1.Number()); v1 != v2 ***REMOVED***
				foundAlias = true
				if !allowAlias ***REMOVED***
					return errors.New("enum %q has conflicting non-aliased values on number %d: %q with %q", e.FullName(), v1.Number(), v1.Name(), v2.Name())
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if allowAlias && !foundAlias ***REMOVED***
			return errors.New("enum %q allows aliases, but none were found", e.FullName())
		***REMOVED***
		if e.Syntax() == protoreflect.Proto3 ***REMOVED***
			if v := e.Values().Get(0); v.Number() != 0 ***REMOVED***
				return errors.New("enum %q using proto3 semantics must have zero number for the first value", v.FullName())
			***REMOVED***
			// Verify that value names in proto3 do not conflict if the
			// case-insensitive prefix is removed.
			// See protoc v3.8.0: src/google/protobuf/descriptor.cc:4991-5055
			names := map[string]protoreflect.EnumValueDescriptor***REMOVED******REMOVED***
			prefix := strings.Replace(strings.ToLower(string(e.Name())), "_", "", -1)
			for i := 0; i < e.Values().Len(); i++ ***REMOVED***
				v1 := e.Values().Get(i)
				s := strs.EnumValueName(strs.TrimEnumPrefix(string(v1.Name()), prefix))
				if v2, ok := names[s]; ok && v1.Number() != v2.Number() ***REMOVED***
					return errors.New("enum %q using proto3 semantics has conflict: %q with %q", e.FullName(), v1.Name(), v2.Name())
				***REMOVED***
				names[s] = v1
			***REMOVED***
		***REMOVED***

		for j, vd := range ed.GetValue() ***REMOVED***
			v := &e.L2.Values.List[j]
			if vd.Number == nil ***REMOVED***
				return errors.New("enum value %q must have a specified number", v.FullName())
			***REMOVED***
			if e.L2.ReservedNames.Has(v.Name()) ***REMOVED***
				return errors.New("enum value %q must not use reserved name", v.FullName())
			***REMOVED***
			if e.L2.ReservedRanges.Has(v.Number()) ***REMOVED***
				return errors.New("enum value %q must not use reserved number %d", v.FullName(), v.Number())
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func validateMessageDeclarations(ms []filedesc.Message, mds []*descriptorpb.DescriptorProto) error ***REMOVED***
	for i, md := range mds ***REMOVED***
		m := &ms[i]

		// Handle the message descriptor itself.
		isMessageSet := md.GetOptions().GetMessageSetWireFormat()
		if err := m.L2.ReservedNames.CheckValid(); err != nil ***REMOVED***
			return errors.New("message %q reserved names has %v", m.FullName(), err)
		***REMOVED***
		if err := m.L2.ReservedRanges.CheckValid(isMessageSet); err != nil ***REMOVED***
			return errors.New("message %q reserved ranges has %v", m.FullName(), err)
		***REMOVED***
		if err := m.L2.ExtensionRanges.CheckValid(isMessageSet); err != nil ***REMOVED***
			return errors.New("message %q extension ranges has %v", m.FullName(), err)
		***REMOVED***
		if err := (*filedesc.FieldRanges).CheckOverlap(&m.L2.ReservedRanges, &m.L2.ExtensionRanges); err != nil ***REMOVED***
			return errors.New("message %q reserved and extension ranges has %v", m.FullName(), err)
		***REMOVED***
		for i := 0; i < m.Fields().Len(); i++ ***REMOVED***
			f1 := m.Fields().Get(i)
			if f2 := m.Fields().ByNumber(f1.Number()); f1 != f2 ***REMOVED***
				return errors.New("message %q has conflicting fields: %q with %q", m.FullName(), f1.Name(), f2.Name())
			***REMOVED***
		***REMOVED***
		if isMessageSet && !flags.ProtoLegacy ***REMOVED***
			return errors.New("message %q is a MessageSet, which is a legacy proto1 feature that is no longer supported", m.FullName())
		***REMOVED***
		if isMessageSet && (m.Syntax() != protoreflect.Proto2 || m.Fields().Len() > 0 || m.ExtensionRanges().Len() == 0) ***REMOVED***
			return errors.New("message %q is an invalid proto1 MessageSet", m.FullName())
		***REMOVED***
		if m.Syntax() == protoreflect.Proto3 ***REMOVED***
			if m.ExtensionRanges().Len() > 0 ***REMOVED***
				return errors.New("message %q using proto3 semantics cannot have extension ranges", m.FullName())
			***REMOVED***
			// Verify that field names in proto3 do not conflict if lowercased
			// with all underscores removed.
			// See protoc v3.8.0: src/google/protobuf/descriptor.cc:5830-5847
			names := map[string]protoreflect.FieldDescriptor***REMOVED******REMOVED***
			for i := 0; i < m.Fields().Len(); i++ ***REMOVED***
				f1 := m.Fields().Get(i)
				s := strings.Replace(strings.ToLower(string(f1.Name())), "_", "", -1)
				if f2, ok := names[s]; ok ***REMOVED***
					return errors.New("message %q using proto3 semantics has conflict: %q with %q", m.FullName(), f1.Name(), f2.Name())
				***REMOVED***
				names[s] = f1
			***REMOVED***
		***REMOVED***

		for j, fd := range md.GetField() ***REMOVED***
			f := &m.L2.Fields.List[j]
			if m.L2.ReservedNames.Has(f.Name()) ***REMOVED***
				return errors.New("message field %q must not use reserved name", f.FullName())
			***REMOVED***
			if !f.Number().IsValid() ***REMOVED***
				return errors.New("message field %q has an invalid number: %d", f.FullName(), f.Number())
			***REMOVED***
			if !f.Cardinality().IsValid() ***REMOVED***
				return errors.New("message field %q has an invalid cardinality: %d", f.FullName(), f.Cardinality())
			***REMOVED***
			if m.L2.ReservedRanges.Has(f.Number()) ***REMOVED***
				return errors.New("message field %q must not use reserved number %d", f.FullName(), f.Number())
			***REMOVED***
			if m.L2.ExtensionRanges.Has(f.Number()) ***REMOVED***
				return errors.New("message field %q with number %d in extension range", f.FullName(), f.Number())
			***REMOVED***
			if fd.Extendee != nil ***REMOVED***
				return errors.New("message field %q may not have extendee: %q", f.FullName(), fd.GetExtendee())
			***REMOVED***
			if f.L1.IsProto3Optional ***REMOVED***
				if f.Syntax() != protoreflect.Proto3 ***REMOVED***
					return errors.New("message field %q under proto3 optional semantics must be specified in the proto3 syntax", f.FullName())
				***REMOVED***
				if f.Cardinality() != protoreflect.Optional ***REMOVED***
					return errors.New("message field %q under proto3 optional semantics must have optional cardinality", f.FullName())
				***REMOVED***
				if f.ContainingOneof() != nil && f.ContainingOneof().Fields().Len() != 1 ***REMOVED***
					return errors.New("message field %q under proto3 optional semantics must be within a single element oneof", f.FullName())
				***REMOVED***
			***REMOVED***
			if f.IsWeak() && !flags.ProtoLegacy ***REMOVED***
				return errors.New("message field %q is a weak field, which is a legacy proto1 feature that is no longer supported", f.FullName())
			***REMOVED***
			if f.IsWeak() && (f.Syntax() != protoreflect.Proto2 || !isOptionalMessage(f) || f.ContainingOneof() != nil) ***REMOVED***
				return errors.New("message field %q may only be weak for an optional message", f.FullName())
			***REMOVED***
			if f.IsPacked() && !isPackable(f) ***REMOVED***
				return errors.New("message field %q is not packable", f.FullName())
			***REMOVED***
			if err := checkValidGroup(f); err != nil ***REMOVED***
				return errors.New("message field %q is an invalid group: %v", f.FullName(), err)
			***REMOVED***
			if err := checkValidMap(f); err != nil ***REMOVED***
				return errors.New("message field %q is an invalid map: %v", f.FullName(), err)
			***REMOVED***
			if f.Syntax() == protoreflect.Proto3 ***REMOVED***
				if f.Cardinality() == protoreflect.Required ***REMOVED***
					return errors.New("message field %q using proto3 semantics cannot be required", f.FullName())
				***REMOVED***
				if f.Enum() != nil && !f.Enum().IsPlaceholder() && f.Enum().Syntax() != protoreflect.Proto3 ***REMOVED***
					return errors.New("message field %q using proto3 semantics may only depend on a proto3 enum", f.FullName())
				***REMOVED***
			***REMOVED***
		***REMOVED***
		seenSynthetic := false // synthetic oneofs for proto3 optional must come after real oneofs
		for j := range md.GetOneofDecl() ***REMOVED***
			o := &m.L2.Oneofs.List[j]
			if o.Fields().Len() == 0 ***REMOVED***
				return errors.New("message oneof %q must contain at least one field declaration", o.FullName())
			***REMOVED***
			if n := o.Fields().Len(); n-1 != (o.Fields().Get(n-1).Index() - o.Fields().Get(0).Index()) ***REMOVED***
				return errors.New("message oneof %q must have consecutively declared fields", o.FullName())
			***REMOVED***

			if o.IsSynthetic() ***REMOVED***
				seenSynthetic = true
				continue
			***REMOVED***
			if !o.IsSynthetic() && seenSynthetic ***REMOVED***
				return errors.New("message oneof %q must be declared before synthetic oneofs", o.FullName())
			***REMOVED***

			for i := 0; i < o.Fields().Len(); i++ ***REMOVED***
				f := o.Fields().Get(i)
				if f.Cardinality() != protoreflect.Optional ***REMOVED***
					return errors.New("message field %q belongs in a oneof and must be optional", f.FullName())
				***REMOVED***
				if f.IsWeak() ***REMOVED***
					return errors.New("message field %q belongs in a oneof and must not be a weak reference", f.FullName())
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if err := validateEnumDeclarations(m.L1.Enums.List, md.GetEnumType()); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := validateMessageDeclarations(m.L1.Messages.List, md.GetNestedType()); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := validateExtensionDeclarations(m.L1.Extensions.List, md.GetExtension()); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func validateExtensionDeclarations(xs []filedesc.Extension, xds []*descriptorpb.FieldDescriptorProto) error ***REMOVED***
	for i, xd := range xds ***REMOVED***
		x := &xs[i]
		// NOTE: Avoid using the IsValid method since extensions to MessageSet
		// may have a field number higher than normal. This check only verifies
		// that the number is not negative or reserved. We check again later
		// if we know that the extendee is definitely not a MessageSet.
		if n := x.Number(); n < 0 || (protowire.FirstReservedNumber <= n && n <= protowire.LastReservedNumber) ***REMOVED***
			return errors.New("extension field %q has an invalid number: %d", x.FullName(), x.Number())
		***REMOVED***
		if !x.Cardinality().IsValid() || x.Cardinality() == protoreflect.Required ***REMOVED***
			return errors.New("extension field %q has an invalid cardinality: %d", x.FullName(), x.Cardinality())
		***REMOVED***
		if xd.JsonName != nil ***REMOVED***
			if xd.GetJsonName() != strs.JSONCamelCase(string(x.Name())) ***REMOVED***
				return errors.New("extension field %q may not have an explicitly set JSON name: %q", x.FullName(), xd.GetJsonName())
			***REMOVED***
		***REMOVED***
		if xd.OneofIndex != nil ***REMOVED***
			return errors.New("extension field %q may not be part of a oneof", x.FullName())
		***REMOVED***
		if md := x.ContainingMessage(); !md.IsPlaceholder() ***REMOVED***
			if !md.ExtensionRanges().Has(x.Number()) ***REMOVED***
				return errors.New("extension field %q extends %q with non-extension field number: %d", x.FullName(), md.FullName(), x.Number())
			***REMOVED***
			isMessageSet := md.Options().(*descriptorpb.MessageOptions).GetMessageSetWireFormat()
			if isMessageSet && !isOptionalMessage(x) ***REMOVED***
				return errors.New("extension field %q extends MessageSet and must be an optional message", x.FullName())
			***REMOVED***
			if !isMessageSet && !x.Number().IsValid() ***REMOVED***
				return errors.New("extension field %q has an invalid number: %d", x.FullName(), x.Number())
			***REMOVED***
		***REMOVED***
		if xd.GetOptions().GetWeak() ***REMOVED***
			return errors.New("extension field %q cannot be a weak reference", x.FullName())
		***REMOVED***
		if x.IsPacked() && !isPackable(x) ***REMOVED***
			return errors.New("extension field %q is not packable", x.FullName())
		***REMOVED***
		if err := checkValidGroup(x); err != nil ***REMOVED***
			return errors.New("extension field %q is an invalid group: %v", x.FullName(), err)
		***REMOVED***
		if md := x.Message(); md != nil && md.IsMapEntry() ***REMOVED***
			return errors.New("extension field %q cannot be a map entry", x.FullName())
		***REMOVED***
		if x.Syntax() == protoreflect.Proto3 ***REMOVED***
			switch x.ContainingMessage().FullName() ***REMOVED***
			case (*descriptorpb.FileOptions)(nil).ProtoReflect().Descriptor().FullName():
			case (*descriptorpb.EnumOptions)(nil).ProtoReflect().Descriptor().FullName():
			case (*descriptorpb.EnumValueOptions)(nil).ProtoReflect().Descriptor().FullName():
			case (*descriptorpb.MessageOptions)(nil).ProtoReflect().Descriptor().FullName():
			case (*descriptorpb.FieldOptions)(nil).ProtoReflect().Descriptor().FullName():
			case (*descriptorpb.OneofOptions)(nil).ProtoReflect().Descriptor().FullName():
			case (*descriptorpb.ExtensionRangeOptions)(nil).ProtoReflect().Descriptor().FullName():
			case (*descriptorpb.ServiceOptions)(nil).ProtoReflect().Descriptor().FullName():
			case (*descriptorpb.MethodOptions)(nil).ProtoReflect().Descriptor().FullName():
			default:
				return errors.New("extension field %q cannot be declared in proto3 unless extended descriptor options", x.FullName())
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// isOptionalMessage reports whether this is an optional message.
// If the kind is unknown, it is assumed to be a message.
func isOptionalMessage(fd protoreflect.FieldDescriptor) bool ***REMOVED***
	return (fd.Kind() == 0 || fd.Kind() == protoreflect.MessageKind) && fd.Cardinality() == protoreflect.Optional
***REMOVED***

// isPackable checks whether the pack option can be specified.
func isPackable(fd protoreflect.FieldDescriptor) bool ***REMOVED***
	switch fd.Kind() ***REMOVED***
	case protoreflect.StringKind, protoreflect.BytesKind, protoreflect.MessageKind, protoreflect.GroupKind:
		return false
	***REMOVED***
	return fd.IsList()
***REMOVED***

// checkValidGroup reports whether fd is a valid group according to the same
// rules that protoc imposes.
func checkValidGroup(fd protoreflect.FieldDescriptor) error ***REMOVED***
	md := fd.Message()
	switch ***REMOVED***
	case fd.Kind() != protoreflect.GroupKind:
		return nil
	case fd.Syntax() != protoreflect.Proto2:
		return errors.New("invalid under proto2 semantics")
	case md == nil || md.IsPlaceholder():
		return errors.New("message must be resolvable")
	case fd.FullName().Parent() != md.FullName().Parent():
		return errors.New("message and field must be declared in the same scope")
	case !unicode.IsUpper(rune(md.Name()[0])):
		return errors.New("message name must start with an uppercase")
	case fd.Name() != protoreflect.Name(strings.ToLower(string(md.Name()))):
		return errors.New("field name must be lowercased form of the message name")
	***REMOVED***
	return nil
***REMOVED***

// checkValidMap checks whether the field is a valid map according to the same
// rules that protoc imposes.
// See protoc v3.8.0: src/google/protobuf/descriptor.cc:6045-6115
func checkValidMap(fd protoreflect.FieldDescriptor) error ***REMOVED***
	md := fd.Message()
	switch ***REMOVED***
	case md == nil || !md.IsMapEntry():
		return nil
	case fd.FullName().Parent() != md.FullName().Parent():
		return errors.New("message and field must be declared in the same scope")
	case md.Name() != protoreflect.Name(strs.MapEntryName(string(fd.Name()))):
		return errors.New("incorrect implicit map entry name")
	case fd.Cardinality() != protoreflect.Repeated:
		return errors.New("field must be repeated")
	case md.Fields().Len() != 2:
		return errors.New("message must have exactly two fields")
	case md.ExtensionRanges().Len() > 0:
		return errors.New("message must not have any extension ranges")
	case md.Enums().Len()+md.Messages().Len()+md.Extensions().Len() > 0:
		return errors.New("message must not have any nested declarations")
	***REMOVED***
	kf := md.Fields().Get(0)
	vf := md.Fields().Get(1)
	switch ***REMOVED***
	case kf.Name() != "key" || kf.Number() != 1 || kf.Cardinality() != protoreflect.Optional || kf.ContainingOneof() != nil || kf.HasDefault():
		return errors.New("invalid key field")
	case vf.Name() != "value" || vf.Number() != 2 || vf.Cardinality() != protoreflect.Optional || vf.ContainingOneof() != nil || vf.HasDefault():
		return errors.New("invalid value field")
	***REMOVED***
	switch kf.Kind() ***REMOVED***
	case protoreflect.BoolKind: // bool
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind: // int32
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind: // int64
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind: // uint32
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind: // uint64
	case protoreflect.StringKind: // string
	default:
		return errors.New("invalid key kind: %v", kf.Kind())
	***REMOVED***
	if e := vf.Enum(); e != nil && e.Values().Len() > 0 && e.Values().Get(0).Number() != 0 ***REMOVED***
		return errors.New("map enum value must have zero number for the first value")
	***REMOVED***
	return nil
***REMOVED***
