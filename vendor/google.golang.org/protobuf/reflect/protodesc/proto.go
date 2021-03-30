// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protodesc

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/internal/encoding/defval"
	"google.golang.org/protobuf/internal/strs"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"google.golang.org/protobuf/types/descriptorpb"
)

// ToFileDescriptorProto copies a protoreflect.FileDescriptor into a
// google.protobuf.FileDescriptorProto message.
func ToFileDescriptorProto(file protoreflect.FileDescriptor) *descriptorpb.FileDescriptorProto ***REMOVED***
	p := &descriptorpb.FileDescriptorProto***REMOVED***
		Name:    proto.String(file.Path()),
		Package: proto.String(string(file.Package())),
		Options: proto.Clone(file.Options()).(*descriptorpb.FileOptions),
	***REMOVED***
	for i, imports := 0, file.Imports(); i < imports.Len(); i++ ***REMOVED***
		imp := imports.Get(i)
		p.Dependency = append(p.Dependency, imp.Path())
		if imp.IsPublic ***REMOVED***
			p.PublicDependency = append(p.PublicDependency, int32(i))
		***REMOVED***
		if imp.IsWeak ***REMOVED***
			p.WeakDependency = append(p.WeakDependency, int32(i))
		***REMOVED***
	***REMOVED***
	for i, locs := 0, file.SourceLocations(); i < locs.Len(); i++ ***REMOVED***
		loc := locs.Get(i)
		l := &descriptorpb.SourceCodeInfo_Location***REMOVED******REMOVED***
		l.Path = append(l.Path, loc.Path...)
		if loc.StartLine == loc.EndLine ***REMOVED***
			l.Span = []int32***REMOVED***int32(loc.StartLine), int32(loc.StartColumn), int32(loc.EndColumn)***REMOVED***
		***REMOVED*** else ***REMOVED***
			l.Span = []int32***REMOVED***int32(loc.StartLine), int32(loc.StartColumn), int32(loc.EndLine), int32(loc.EndColumn)***REMOVED***
		***REMOVED***
		l.LeadingDetachedComments = append([]string(nil), loc.LeadingDetachedComments...)
		if loc.LeadingComments != "" ***REMOVED***
			l.LeadingComments = proto.String(loc.LeadingComments)
		***REMOVED***
		if loc.TrailingComments != "" ***REMOVED***
			l.TrailingComments = proto.String(loc.TrailingComments)
		***REMOVED***
		if p.SourceCodeInfo == nil ***REMOVED***
			p.SourceCodeInfo = &descriptorpb.SourceCodeInfo***REMOVED******REMOVED***
		***REMOVED***
		p.SourceCodeInfo.Location = append(p.SourceCodeInfo.Location, l)

	***REMOVED***
	for i, messages := 0, file.Messages(); i < messages.Len(); i++ ***REMOVED***
		p.MessageType = append(p.MessageType, ToDescriptorProto(messages.Get(i)))
	***REMOVED***
	for i, enums := 0, file.Enums(); i < enums.Len(); i++ ***REMOVED***
		p.EnumType = append(p.EnumType, ToEnumDescriptorProto(enums.Get(i)))
	***REMOVED***
	for i, services := 0, file.Services(); i < services.Len(); i++ ***REMOVED***
		p.Service = append(p.Service, ToServiceDescriptorProto(services.Get(i)))
	***REMOVED***
	for i, exts := 0, file.Extensions(); i < exts.Len(); i++ ***REMOVED***
		p.Extension = append(p.Extension, ToFieldDescriptorProto(exts.Get(i)))
	***REMOVED***
	if syntax := file.Syntax(); syntax != protoreflect.Proto2 ***REMOVED***
		p.Syntax = proto.String(file.Syntax().String())
	***REMOVED***
	return p
***REMOVED***

// ToDescriptorProto copies a protoreflect.MessageDescriptor into a
// google.protobuf.DescriptorProto message.
func ToDescriptorProto(message protoreflect.MessageDescriptor) *descriptorpb.DescriptorProto ***REMOVED***
	p := &descriptorpb.DescriptorProto***REMOVED***
		Name:    proto.String(string(message.Name())),
		Options: proto.Clone(message.Options()).(*descriptorpb.MessageOptions),
	***REMOVED***
	for i, fields := 0, message.Fields(); i < fields.Len(); i++ ***REMOVED***
		p.Field = append(p.Field, ToFieldDescriptorProto(fields.Get(i)))
	***REMOVED***
	for i, exts := 0, message.Extensions(); i < exts.Len(); i++ ***REMOVED***
		p.Extension = append(p.Extension, ToFieldDescriptorProto(exts.Get(i)))
	***REMOVED***
	for i, messages := 0, message.Messages(); i < messages.Len(); i++ ***REMOVED***
		p.NestedType = append(p.NestedType, ToDescriptorProto(messages.Get(i)))
	***REMOVED***
	for i, enums := 0, message.Enums(); i < enums.Len(); i++ ***REMOVED***
		p.EnumType = append(p.EnumType, ToEnumDescriptorProto(enums.Get(i)))
	***REMOVED***
	for i, xranges := 0, message.ExtensionRanges(); i < xranges.Len(); i++ ***REMOVED***
		xrange := xranges.Get(i)
		p.ExtensionRange = append(p.ExtensionRange, &descriptorpb.DescriptorProto_ExtensionRange***REMOVED***
			Start:   proto.Int32(int32(xrange[0])),
			End:     proto.Int32(int32(xrange[1])),
			Options: proto.Clone(message.ExtensionRangeOptions(i)).(*descriptorpb.ExtensionRangeOptions),
		***REMOVED***)
	***REMOVED***
	for i, oneofs := 0, message.Oneofs(); i < oneofs.Len(); i++ ***REMOVED***
		p.OneofDecl = append(p.OneofDecl, ToOneofDescriptorProto(oneofs.Get(i)))
	***REMOVED***
	for i, ranges := 0, message.ReservedRanges(); i < ranges.Len(); i++ ***REMOVED***
		rrange := ranges.Get(i)
		p.ReservedRange = append(p.ReservedRange, &descriptorpb.DescriptorProto_ReservedRange***REMOVED***
			Start: proto.Int32(int32(rrange[0])),
			End:   proto.Int32(int32(rrange[1])),
		***REMOVED***)
	***REMOVED***
	for i, names := 0, message.ReservedNames(); i < names.Len(); i++ ***REMOVED***
		p.ReservedName = append(p.ReservedName, string(names.Get(i)))
	***REMOVED***
	return p
***REMOVED***

// ToFieldDescriptorProto copies a protoreflect.FieldDescriptor into a
// google.protobuf.FieldDescriptorProto message.
func ToFieldDescriptorProto(field protoreflect.FieldDescriptor) *descriptorpb.FieldDescriptorProto ***REMOVED***
	p := &descriptorpb.FieldDescriptorProto***REMOVED***
		Name:    proto.String(string(field.Name())),
		Number:  proto.Int32(int32(field.Number())),
		Label:   descriptorpb.FieldDescriptorProto_Label(field.Cardinality()).Enum(),
		Options: proto.Clone(field.Options()).(*descriptorpb.FieldOptions),
	***REMOVED***
	if field.IsExtension() ***REMOVED***
		p.Extendee = fullNameOf(field.ContainingMessage())
	***REMOVED***
	if field.Kind().IsValid() ***REMOVED***
		p.Type = descriptorpb.FieldDescriptorProto_Type(field.Kind()).Enum()
	***REMOVED***
	if field.Enum() != nil ***REMOVED***
		p.TypeName = fullNameOf(field.Enum())
	***REMOVED***
	if field.Message() != nil ***REMOVED***
		p.TypeName = fullNameOf(field.Message())
	***REMOVED***
	if field.HasJSONName() ***REMOVED***
		// A bug in older versions of protoc would always populate the
		// "json_name" option for extensions when it is meaningless.
		// When it did so, it would always use the camel-cased field name.
		if field.IsExtension() ***REMOVED***
			p.JsonName = proto.String(strs.JSONCamelCase(string(field.Name())))
		***REMOVED*** else ***REMOVED***
			p.JsonName = proto.String(field.JSONName())
		***REMOVED***
	***REMOVED***
	if field.Syntax() == protoreflect.Proto3 && field.HasOptionalKeyword() ***REMOVED***
		p.Proto3Optional = proto.Bool(true)
	***REMOVED***
	if field.HasDefault() ***REMOVED***
		def, err := defval.Marshal(field.Default(), field.DefaultEnumValue(), field.Kind(), defval.Descriptor)
		if err != nil && field.DefaultEnumValue() != nil ***REMOVED***
			def = string(field.DefaultEnumValue().Name()) // occurs for unresolved enum values
		***REMOVED*** else if err != nil ***REMOVED***
			panic(fmt.Sprintf("%v: %v", field.FullName(), err))
		***REMOVED***
		p.DefaultValue = proto.String(def)
	***REMOVED***
	if oneof := field.ContainingOneof(); oneof != nil ***REMOVED***
		p.OneofIndex = proto.Int32(int32(oneof.Index()))
	***REMOVED***
	return p
***REMOVED***

// ToOneofDescriptorProto copies a protoreflect.OneofDescriptor into a
// google.protobuf.OneofDescriptorProto message.
func ToOneofDescriptorProto(oneof protoreflect.OneofDescriptor) *descriptorpb.OneofDescriptorProto ***REMOVED***
	return &descriptorpb.OneofDescriptorProto***REMOVED***
		Name:    proto.String(string(oneof.Name())),
		Options: proto.Clone(oneof.Options()).(*descriptorpb.OneofOptions),
	***REMOVED***
***REMOVED***

// ToEnumDescriptorProto copies a protoreflect.EnumDescriptor into a
// google.protobuf.EnumDescriptorProto message.
func ToEnumDescriptorProto(enum protoreflect.EnumDescriptor) *descriptorpb.EnumDescriptorProto ***REMOVED***
	p := &descriptorpb.EnumDescriptorProto***REMOVED***
		Name:    proto.String(string(enum.Name())),
		Options: proto.Clone(enum.Options()).(*descriptorpb.EnumOptions),
	***REMOVED***
	for i, values := 0, enum.Values(); i < values.Len(); i++ ***REMOVED***
		p.Value = append(p.Value, ToEnumValueDescriptorProto(values.Get(i)))
	***REMOVED***
	for i, ranges := 0, enum.ReservedRanges(); i < ranges.Len(); i++ ***REMOVED***
		rrange := ranges.Get(i)
		p.ReservedRange = append(p.ReservedRange, &descriptorpb.EnumDescriptorProto_EnumReservedRange***REMOVED***
			Start: proto.Int32(int32(rrange[0])),
			End:   proto.Int32(int32(rrange[1])),
		***REMOVED***)
	***REMOVED***
	for i, names := 0, enum.ReservedNames(); i < names.Len(); i++ ***REMOVED***
		p.ReservedName = append(p.ReservedName, string(names.Get(i)))
	***REMOVED***
	return p
***REMOVED***

// ToEnumValueDescriptorProto copies a protoreflect.EnumValueDescriptor into a
// google.protobuf.EnumValueDescriptorProto message.
func ToEnumValueDescriptorProto(value protoreflect.EnumValueDescriptor) *descriptorpb.EnumValueDescriptorProto ***REMOVED***
	return &descriptorpb.EnumValueDescriptorProto***REMOVED***
		Name:    proto.String(string(value.Name())),
		Number:  proto.Int32(int32(value.Number())),
		Options: proto.Clone(value.Options()).(*descriptorpb.EnumValueOptions),
	***REMOVED***
***REMOVED***

// ToServiceDescriptorProto copies a protoreflect.ServiceDescriptor into a
// google.protobuf.ServiceDescriptorProto message.
func ToServiceDescriptorProto(service protoreflect.ServiceDescriptor) *descriptorpb.ServiceDescriptorProto ***REMOVED***
	p := &descriptorpb.ServiceDescriptorProto***REMOVED***
		Name:    proto.String(string(service.Name())),
		Options: proto.Clone(service.Options()).(*descriptorpb.ServiceOptions),
	***REMOVED***
	for i, methods := 0, service.Methods(); i < methods.Len(); i++ ***REMOVED***
		p.Method = append(p.Method, ToMethodDescriptorProto(methods.Get(i)))
	***REMOVED***
	return p
***REMOVED***

// ToMethodDescriptorProto copies a protoreflect.MethodDescriptor into a
// google.protobuf.MethodDescriptorProto message.
func ToMethodDescriptorProto(method protoreflect.MethodDescriptor) *descriptorpb.MethodDescriptorProto ***REMOVED***
	p := &descriptorpb.MethodDescriptorProto***REMOVED***
		Name:       proto.String(string(method.Name())),
		InputType:  fullNameOf(method.Input()),
		OutputType: fullNameOf(method.Output()),
		Options:    proto.Clone(method.Options()).(*descriptorpb.MethodOptions),
	***REMOVED***
	if method.IsStreamingClient() ***REMOVED***
		p.ClientStreaming = proto.Bool(true)
	***REMOVED***
	if method.IsStreamingServer() ***REMOVED***
		p.ServerStreaming = proto.Bool(true)
	***REMOVED***
	return p
***REMOVED***

func fullNameOf(d protoreflect.Descriptor) *string ***REMOVED***
	if d == nil ***REMOVED***
		return nil
	***REMOVED***
	if strings.HasPrefix(string(d.FullName()), unknownPrefix) ***REMOVED***
		return proto.String(string(d.FullName()[len(unknownPrefix):]))
	***REMOVED***
	return proto.String("." + string(d.FullName()))
***REMOVED***
