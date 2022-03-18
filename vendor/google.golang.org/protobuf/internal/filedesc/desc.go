// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filedesc

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"

	"google.golang.org/protobuf/internal/descfmt"
	"google.golang.org/protobuf/internal/descopts"
	"google.golang.org/protobuf/internal/encoding/defval"
	"google.golang.org/protobuf/internal/encoding/messageset"
	"google.golang.org/protobuf/internal/genid"
	"google.golang.org/protobuf/internal/pragma"
	"google.golang.org/protobuf/internal/strs"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// The types in this file may have a suffix:
//	• L0: Contains fields common to all descriptors (except File) and
//	must be initialized up front.
//	• L1: Contains fields specific to a descriptor and
//	must be initialized up front.
//	• L2: Contains fields that are lazily initialized when constructing
//	from the raw file descriptor. When constructing as a literal, the L2
//	fields must be initialized up front.
//
// The types are exported so that packages like reflect/protodesc can
// directly construct descriptors.

type (
	File struct ***REMOVED***
		fileRaw
		L1 FileL1

		once uint32     // atomically set if L2 is valid
		mu   sync.Mutex // protects L2
		L2   *FileL2
	***REMOVED***
	FileL1 struct ***REMOVED***
		Syntax  pref.Syntax
		Path    string
		Package pref.FullName

		Enums      Enums
		Messages   Messages
		Extensions Extensions
		Services   Services
	***REMOVED***
	FileL2 struct ***REMOVED***
		Options   func() pref.ProtoMessage
		Imports   FileImports
		Locations SourceLocations
	***REMOVED***
)

func (fd *File) ParentFile() pref.FileDescriptor ***REMOVED*** return fd ***REMOVED***
func (fd *File) Parent() pref.Descriptor         ***REMOVED*** return nil ***REMOVED***
func (fd *File) Index() int                      ***REMOVED*** return 0 ***REMOVED***
func (fd *File) Syntax() pref.Syntax             ***REMOVED*** return fd.L1.Syntax ***REMOVED***
func (fd *File) Name() pref.Name                 ***REMOVED*** return fd.L1.Package.Name() ***REMOVED***
func (fd *File) FullName() pref.FullName         ***REMOVED*** return fd.L1.Package ***REMOVED***
func (fd *File) IsPlaceholder() bool             ***REMOVED*** return false ***REMOVED***
func (fd *File) Options() pref.ProtoMessage ***REMOVED***
	if f := fd.lazyInit().Options; f != nil ***REMOVED***
		return f()
	***REMOVED***
	return descopts.File
***REMOVED***
func (fd *File) Path() string                          ***REMOVED*** return fd.L1.Path ***REMOVED***
func (fd *File) Package() pref.FullName                ***REMOVED*** return fd.L1.Package ***REMOVED***
func (fd *File) Imports() pref.FileImports             ***REMOVED*** return &fd.lazyInit().Imports ***REMOVED***
func (fd *File) Enums() pref.EnumDescriptors           ***REMOVED*** return &fd.L1.Enums ***REMOVED***
func (fd *File) Messages() pref.MessageDescriptors     ***REMOVED*** return &fd.L1.Messages ***REMOVED***
func (fd *File) Extensions() pref.ExtensionDescriptors ***REMOVED*** return &fd.L1.Extensions ***REMOVED***
func (fd *File) Services() pref.ServiceDescriptors     ***REMOVED*** return &fd.L1.Services ***REMOVED***
func (fd *File) SourceLocations() pref.SourceLocations ***REMOVED*** return &fd.lazyInit().Locations ***REMOVED***
func (fd *File) Format(s fmt.State, r rune)            ***REMOVED*** descfmt.FormatDesc(s, r, fd) ***REMOVED***
func (fd *File) ProtoType(pref.FileDescriptor)         ***REMOVED******REMOVED***
func (fd *File) ProtoInternal(pragma.DoNotImplement)   ***REMOVED******REMOVED***

func (fd *File) lazyInit() *FileL2 ***REMOVED***
	if atomic.LoadUint32(&fd.once) == 0 ***REMOVED***
		fd.lazyInitOnce()
	***REMOVED***
	return fd.L2
***REMOVED***

func (fd *File) lazyInitOnce() ***REMOVED***
	fd.mu.Lock()
	if fd.L2 == nil ***REMOVED***
		fd.lazyRawInit() // recursively initializes all L2 structures
	***REMOVED***
	atomic.StoreUint32(&fd.once, 1)
	fd.mu.Unlock()
***REMOVED***

// GoPackagePath is a pseudo-internal API for determining the Go package path
// that this file descriptor is declared in.
//
// WARNING: This method is exempt from the compatibility promise and may be
// removed in the future without warning.
func (fd *File) GoPackagePath() string ***REMOVED***
	return fd.builder.GoPackagePath
***REMOVED***

type (
	Enum struct ***REMOVED***
		Base
		L1 EnumL1
		L2 *EnumL2 // protected by fileDesc.once
	***REMOVED***
	EnumL1 struct ***REMOVED***
		eagerValues bool // controls whether EnumL2.Values is already populated
	***REMOVED***
	EnumL2 struct ***REMOVED***
		Options        func() pref.ProtoMessage
		Values         EnumValues
		ReservedNames  Names
		ReservedRanges EnumRanges
	***REMOVED***

	EnumValue struct ***REMOVED***
		Base
		L1 EnumValueL1
	***REMOVED***
	EnumValueL1 struct ***REMOVED***
		Options func() pref.ProtoMessage
		Number  pref.EnumNumber
	***REMOVED***
)

func (ed *Enum) Options() pref.ProtoMessage ***REMOVED***
	if f := ed.lazyInit().Options; f != nil ***REMOVED***
		return f()
	***REMOVED***
	return descopts.Enum
***REMOVED***
func (ed *Enum) Values() pref.EnumValueDescriptors ***REMOVED***
	if ed.L1.eagerValues ***REMOVED***
		return &ed.L2.Values
	***REMOVED***
	return &ed.lazyInit().Values
***REMOVED***
func (ed *Enum) ReservedNames() pref.Names       ***REMOVED*** return &ed.lazyInit().ReservedNames ***REMOVED***
func (ed *Enum) ReservedRanges() pref.EnumRanges ***REMOVED*** return &ed.lazyInit().ReservedRanges ***REMOVED***
func (ed *Enum) Format(s fmt.State, r rune)      ***REMOVED*** descfmt.FormatDesc(s, r, ed) ***REMOVED***
func (ed *Enum) ProtoType(pref.EnumDescriptor)   ***REMOVED******REMOVED***
func (ed *Enum) lazyInit() *EnumL2 ***REMOVED***
	ed.L0.ParentFile.lazyInit() // implicitly initializes L2
	return ed.L2
***REMOVED***

func (ed *EnumValue) Options() pref.ProtoMessage ***REMOVED***
	if f := ed.L1.Options; f != nil ***REMOVED***
		return f()
	***REMOVED***
	return descopts.EnumValue
***REMOVED***
func (ed *EnumValue) Number() pref.EnumNumber            ***REMOVED*** return ed.L1.Number ***REMOVED***
func (ed *EnumValue) Format(s fmt.State, r rune)         ***REMOVED*** descfmt.FormatDesc(s, r, ed) ***REMOVED***
func (ed *EnumValue) ProtoType(pref.EnumValueDescriptor) ***REMOVED******REMOVED***

type (
	Message struct ***REMOVED***
		Base
		L1 MessageL1
		L2 *MessageL2 // protected by fileDesc.once
	***REMOVED***
	MessageL1 struct ***REMOVED***
		Enums        Enums
		Messages     Messages
		Extensions   Extensions
		IsMapEntry   bool // promoted from google.protobuf.MessageOptions
		IsMessageSet bool // promoted from google.protobuf.MessageOptions
	***REMOVED***
	MessageL2 struct ***REMOVED***
		Options               func() pref.ProtoMessage
		Fields                Fields
		Oneofs                Oneofs
		ReservedNames         Names
		ReservedRanges        FieldRanges
		RequiredNumbers       FieldNumbers // must be consistent with Fields.Cardinality
		ExtensionRanges       FieldRanges
		ExtensionRangeOptions []func() pref.ProtoMessage // must be same length as ExtensionRanges
	***REMOVED***

	Field struct ***REMOVED***
		Base
		L1 FieldL1
	***REMOVED***
	FieldL1 struct ***REMOVED***
		Options          func() pref.ProtoMessage
		Number           pref.FieldNumber
		Cardinality      pref.Cardinality // must be consistent with Message.RequiredNumbers
		Kind             pref.Kind
		StringName       stringName
		IsProto3Optional bool // promoted from google.protobuf.FieldDescriptorProto
		IsWeak           bool // promoted from google.protobuf.FieldOptions
		HasPacked        bool // promoted from google.protobuf.FieldOptions
		IsPacked         bool // promoted from google.protobuf.FieldOptions
		HasEnforceUTF8   bool // promoted from google.protobuf.FieldOptions
		EnforceUTF8      bool // promoted from google.protobuf.FieldOptions
		Default          defaultValue
		ContainingOneof  pref.OneofDescriptor // must be consistent with Message.Oneofs.Fields
		Enum             pref.EnumDescriptor
		Message          pref.MessageDescriptor
	***REMOVED***

	Oneof struct ***REMOVED***
		Base
		L1 OneofL1
	***REMOVED***
	OneofL1 struct ***REMOVED***
		Options func() pref.ProtoMessage
		Fields  OneofFields // must be consistent with Message.Fields.ContainingOneof
	***REMOVED***
)

func (md *Message) Options() pref.ProtoMessage ***REMOVED***
	if f := md.lazyInit().Options; f != nil ***REMOVED***
		return f()
	***REMOVED***
	return descopts.Message
***REMOVED***
func (md *Message) IsMapEntry() bool                   ***REMOVED*** return md.L1.IsMapEntry ***REMOVED***
func (md *Message) Fields() pref.FieldDescriptors      ***REMOVED*** return &md.lazyInit().Fields ***REMOVED***
func (md *Message) Oneofs() pref.OneofDescriptors      ***REMOVED*** return &md.lazyInit().Oneofs ***REMOVED***
func (md *Message) ReservedNames() pref.Names          ***REMOVED*** return &md.lazyInit().ReservedNames ***REMOVED***
func (md *Message) ReservedRanges() pref.FieldRanges   ***REMOVED*** return &md.lazyInit().ReservedRanges ***REMOVED***
func (md *Message) RequiredNumbers() pref.FieldNumbers ***REMOVED*** return &md.lazyInit().RequiredNumbers ***REMOVED***
func (md *Message) ExtensionRanges() pref.FieldRanges  ***REMOVED*** return &md.lazyInit().ExtensionRanges ***REMOVED***
func (md *Message) ExtensionRangeOptions(i int) pref.ProtoMessage ***REMOVED***
	if f := md.lazyInit().ExtensionRangeOptions[i]; f != nil ***REMOVED***
		return f()
	***REMOVED***
	return descopts.ExtensionRange
***REMOVED***
func (md *Message) Enums() pref.EnumDescriptors           ***REMOVED*** return &md.L1.Enums ***REMOVED***
func (md *Message) Messages() pref.MessageDescriptors     ***REMOVED*** return &md.L1.Messages ***REMOVED***
func (md *Message) Extensions() pref.ExtensionDescriptors ***REMOVED*** return &md.L1.Extensions ***REMOVED***
func (md *Message) ProtoType(pref.MessageDescriptor)      ***REMOVED******REMOVED***
func (md *Message) Format(s fmt.State, r rune)            ***REMOVED*** descfmt.FormatDesc(s, r, md) ***REMOVED***
func (md *Message) lazyInit() *MessageL2 ***REMOVED***
	md.L0.ParentFile.lazyInit() // implicitly initializes L2
	return md.L2
***REMOVED***

// IsMessageSet is a pseudo-internal API for checking whether a message
// should serialize in the proto1 message format.
//
// WARNING: This method is exempt from the compatibility promise and may be
// removed in the future without warning.
func (md *Message) IsMessageSet() bool ***REMOVED***
	return md.L1.IsMessageSet
***REMOVED***

func (fd *Field) Options() pref.ProtoMessage ***REMOVED***
	if f := fd.L1.Options; f != nil ***REMOVED***
		return f()
	***REMOVED***
	return descopts.Field
***REMOVED***
func (fd *Field) Number() pref.FieldNumber      ***REMOVED*** return fd.L1.Number ***REMOVED***
func (fd *Field) Cardinality() pref.Cardinality ***REMOVED*** return fd.L1.Cardinality ***REMOVED***
func (fd *Field) Kind() pref.Kind               ***REMOVED*** return fd.L1.Kind ***REMOVED***
func (fd *Field) HasJSONName() bool             ***REMOVED*** return fd.L1.StringName.hasJSON ***REMOVED***
func (fd *Field) JSONName() string              ***REMOVED*** return fd.L1.StringName.getJSON(fd) ***REMOVED***
func (fd *Field) TextName() string              ***REMOVED*** return fd.L1.StringName.getText(fd) ***REMOVED***
func (fd *Field) HasPresence() bool ***REMOVED***
	return fd.L1.Cardinality != pref.Repeated && (fd.L0.ParentFile.L1.Syntax == pref.Proto2 || fd.L1.Message != nil || fd.L1.ContainingOneof != nil)
***REMOVED***
func (fd *Field) HasOptionalKeyword() bool ***REMOVED***
	return (fd.L0.ParentFile.L1.Syntax == pref.Proto2 && fd.L1.Cardinality == pref.Optional && fd.L1.ContainingOneof == nil) || fd.L1.IsProto3Optional
***REMOVED***
func (fd *Field) IsPacked() bool ***REMOVED***
	if !fd.L1.HasPacked && fd.L0.ParentFile.L1.Syntax != pref.Proto2 && fd.L1.Cardinality == pref.Repeated ***REMOVED***
		switch fd.L1.Kind ***REMOVED***
		case pref.StringKind, pref.BytesKind, pref.MessageKind, pref.GroupKind:
		default:
			return true
		***REMOVED***
	***REMOVED***
	return fd.L1.IsPacked
***REMOVED***
func (fd *Field) IsExtension() bool ***REMOVED*** return false ***REMOVED***
func (fd *Field) IsWeak() bool      ***REMOVED*** return fd.L1.IsWeak ***REMOVED***
func (fd *Field) IsList() bool      ***REMOVED*** return fd.Cardinality() == pref.Repeated && !fd.IsMap() ***REMOVED***
func (fd *Field) IsMap() bool       ***REMOVED*** return fd.Message() != nil && fd.Message().IsMapEntry() ***REMOVED***
func (fd *Field) MapKey() pref.FieldDescriptor ***REMOVED***
	if !fd.IsMap() ***REMOVED***
		return nil
	***REMOVED***
	return fd.Message().Fields().ByNumber(genid.MapEntry_Key_field_number)
***REMOVED***
func (fd *Field) MapValue() pref.FieldDescriptor ***REMOVED***
	if !fd.IsMap() ***REMOVED***
		return nil
	***REMOVED***
	return fd.Message().Fields().ByNumber(genid.MapEntry_Value_field_number)
***REMOVED***
func (fd *Field) HasDefault() bool                           ***REMOVED*** return fd.L1.Default.has ***REMOVED***
func (fd *Field) Default() pref.Value                        ***REMOVED*** return fd.L1.Default.get(fd) ***REMOVED***
func (fd *Field) DefaultEnumValue() pref.EnumValueDescriptor ***REMOVED*** return fd.L1.Default.enum ***REMOVED***
func (fd *Field) ContainingOneof() pref.OneofDescriptor      ***REMOVED*** return fd.L1.ContainingOneof ***REMOVED***
func (fd *Field) ContainingMessage() pref.MessageDescriptor ***REMOVED***
	return fd.L0.Parent.(pref.MessageDescriptor)
***REMOVED***
func (fd *Field) Enum() pref.EnumDescriptor ***REMOVED***
	return fd.L1.Enum
***REMOVED***
func (fd *Field) Message() pref.MessageDescriptor ***REMOVED***
	if fd.L1.IsWeak ***REMOVED***
		if d, _ := protoregistry.GlobalFiles.FindDescriptorByName(fd.L1.Message.FullName()); d != nil ***REMOVED***
			return d.(pref.MessageDescriptor)
		***REMOVED***
	***REMOVED***
	return fd.L1.Message
***REMOVED***
func (fd *Field) Format(s fmt.State, r rune)     ***REMOVED*** descfmt.FormatDesc(s, r, fd) ***REMOVED***
func (fd *Field) ProtoType(pref.FieldDescriptor) ***REMOVED******REMOVED***

// EnforceUTF8 is a pseudo-internal API to determine whether to enforce UTF-8
// validation for the string field. This exists for Google-internal use only
// since proto3 did not enforce UTF-8 validity prior to the open-source release.
// If this method does not exist, the default is to enforce valid UTF-8.
//
// WARNING: This method is exempt from the compatibility promise and may be
// removed in the future without warning.
func (fd *Field) EnforceUTF8() bool ***REMOVED***
	if fd.L1.HasEnforceUTF8 ***REMOVED***
		return fd.L1.EnforceUTF8
	***REMOVED***
	return fd.L0.ParentFile.L1.Syntax == pref.Proto3
***REMOVED***

func (od *Oneof) IsSynthetic() bool ***REMOVED***
	return od.L0.ParentFile.L1.Syntax == pref.Proto3 && len(od.L1.Fields.List) == 1 && od.L1.Fields.List[0].HasOptionalKeyword()
***REMOVED***
func (od *Oneof) Options() pref.ProtoMessage ***REMOVED***
	if f := od.L1.Options; f != nil ***REMOVED***
		return f()
	***REMOVED***
	return descopts.Oneof
***REMOVED***
func (od *Oneof) Fields() pref.FieldDescriptors  ***REMOVED*** return &od.L1.Fields ***REMOVED***
func (od *Oneof) Format(s fmt.State, r rune)     ***REMOVED*** descfmt.FormatDesc(s, r, od) ***REMOVED***
func (od *Oneof) ProtoType(pref.OneofDescriptor) ***REMOVED******REMOVED***

type (
	Extension struct ***REMOVED***
		Base
		L1 ExtensionL1
		L2 *ExtensionL2 // protected by fileDesc.once
	***REMOVED***
	ExtensionL1 struct ***REMOVED***
		Number      pref.FieldNumber
		Extendee    pref.MessageDescriptor
		Cardinality pref.Cardinality
		Kind        pref.Kind
	***REMOVED***
	ExtensionL2 struct ***REMOVED***
		Options          func() pref.ProtoMessage
		StringName       stringName
		IsProto3Optional bool // promoted from google.protobuf.FieldDescriptorProto
		IsPacked         bool // promoted from google.protobuf.FieldOptions
		Default          defaultValue
		Enum             pref.EnumDescriptor
		Message          pref.MessageDescriptor
	***REMOVED***
)

func (xd *Extension) Options() pref.ProtoMessage ***REMOVED***
	if f := xd.lazyInit().Options; f != nil ***REMOVED***
		return f()
	***REMOVED***
	return descopts.Field
***REMOVED***
func (xd *Extension) Number() pref.FieldNumber      ***REMOVED*** return xd.L1.Number ***REMOVED***
func (xd *Extension) Cardinality() pref.Cardinality ***REMOVED*** return xd.L1.Cardinality ***REMOVED***
func (xd *Extension) Kind() pref.Kind               ***REMOVED*** return xd.L1.Kind ***REMOVED***
func (xd *Extension) HasJSONName() bool             ***REMOVED*** return xd.lazyInit().StringName.hasJSON ***REMOVED***
func (xd *Extension) JSONName() string              ***REMOVED*** return xd.lazyInit().StringName.getJSON(xd) ***REMOVED***
func (xd *Extension) TextName() string              ***REMOVED*** return xd.lazyInit().StringName.getText(xd) ***REMOVED***
func (xd *Extension) HasPresence() bool             ***REMOVED*** return xd.L1.Cardinality != pref.Repeated ***REMOVED***
func (xd *Extension) HasOptionalKeyword() bool ***REMOVED***
	return (xd.L0.ParentFile.L1.Syntax == pref.Proto2 && xd.L1.Cardinality == pref.Optional) || xd.lazyInit().IsProto3Optional
***REMOVED***
func (xd *Extension) IsPacked() bool                             ***REMOVED*** return xd.lazyInit().IsPacked ***REMOVED***
func (xd *Extension) IsExtension() bool                          ***REMOVED*** return true ***REMOVED***
func (xd *Extension) IsWeak() bool                               ***REMOVED*** return false ***REMOVED***
func (xd *Extension) IsList() bool                               ***REMOVED*** return xd.Cardinality() == pref.Repeated ***REMOVED***
func (xd *Extension) IsMap() bool                                ***REMOVED*** return false ***REMOVED***
func (xd *Extension) MapKey() pref.FieldDescriptor               ***REMOVED*** return nil ***REMOVED***
func (xd *Extension) MapValue() pref.FieldDescriptor             ***REMOVED*** return nil ***REMOVED***
func (xd *Extension) HasDefault() bool                           ***REMOVED*** return xd.lazyInit().Default.has ***REMOVED***
func (xd *Extension) Default() pref.Value                        ***REMOVED*** return xd.lazyInit().Default.get(xd) ***REMOVED***
func (xd *Extension) DefaultEnumValue() pref.EnumValueDescriptor ***REMOVED*** return xd.lazyInit().Default.enum ***REMOVED***
func (xd *Extension) ContainingOneof() pref.OneofDescriptor      ***REMOVED*** return nil ***REMOVED***
func (xd *Extension) ContainingMessage() pref.MessageDescriptor  ***REMOVED*** return xd.L1.Extendee ***REMOVED***
func (xd *Extension) Enum() pref.EnumDescriptor                  ***REMOVED*** return xd.lazyInit().Enum ***REMOVED***
func (xd *Extension) Message() pref.MessageDescriptor            ***REMOVED*** return xd.lazyInit().Message ***REMOVED***
func (xd *Extension) Format(s fmt.State, r rune)                 ***REMOVED*** descfmt.FormatDesc(s, r, xd) ***REMOVED***
func (xd *Extension) ProtoType(pref.FieldDescriptor)             ***REMOVED******REMOVED***
func (xd *Extension) ProtoInternal(pragma.DoNotImplement)        ***REMOVED******REMOVED***
func (xd *Extension) lazyInit() *ExtensionL2 ***REMOVED***
	xd.L0.ParentFile.lazyInit() // implicitly initializes L2
	return xd.L2
***REMOVED***

type (
	Service struct ***REMOVED***
		Base
		L1 ServiceL1
		L2 *ServiceL2 // protected by fileDesc.once
	***REMOVED***
	ServiceL1 struct***REMOVED******REMOVED***
	ServiceL2 struct ***REMOVED***
		Options func() pref.ProtoMessage
		Methods Methods
	***REMOVED***

	Method struct ***REMOVED***
		Base
		L1 MethodL1
	***REMOVED***
	MethodL1 struct ***REMOVED***
		Options           func() pref.ProtoMessage
		Input             pref.MessageDescriptor
		Output            pref.MessageDescriptor
		IsStreamingClient bool
		IsStreamingServer bool
	***REMOVED***
)

func (sd *Service) Options() pref.ProtoMessage ***REMOVED***
	if f := sd.lazyInit().Options; f != nil ***REMOVED***
		return f()
	***REMOVED***
	return descopts.Service
***REMOVED***
func (sd *Service) Methods() pref.MethodDescriptors     ***REMOVED*** return &sd.lazyInit().Methods ***REMOVED***
func (sd *Service) Format(s fmt.State, r rune)          ***REMOVED*** descfmt.FormatDesc(s, r, sd) ***REMOVED***
func (sd *Service) ProtoType(pref.ServiceDescriptor)    ***REMOVED******REMOVED***
func (sd *Service) ProtoInternal(pragma.DoNotImplement) ***REMOVED******REMOVED***
func (sd *Service) lazyInit() *ServiceL2 ***REMOVED***
	sd.L0.ParentFile.lazyInit() // implicitly initializes L2
	return sd.L2
***REMOVED***

func (md *Method) Options() pref.ProtoMessage ***REMOVED***
	if f := md.L1.Options; f != nil ***REMOVED***
		return f()
	***REMOVED***
	return descopts.Method
***REMOVED***
func (md *Method) Input() pref.MessageDescriptor       ***REMOVED*** return md.L1.Input ***REMOVED***
func (md *Method) Output() pref.MessageDescriptor      ***REMOVED*** return md.L1.Output ***REMOVED***
func (md *Method) IsStreamingClient() bool             ***REMOVED*** return md.L1.IsStreamingClient ***REMOVED***
func (md *Method) IsStreamingServer() bool             ***REMOVED*** return md.L1.IsStreamingServer ***REMOVED***
func (md *Method) Format(s fmt.State, r rune)          ***REMOVED*** descfmt.FormatDesc(s, r, md) ***REMOVED***
func (md *Method) ProtoType(pref.MethodDescriptor)     ***REMOVED******REMOVED***
func (md *Method) ProtoInternal(pragma.DoNotImplement) ***REMOVED******REMOVED***

// Surrogate files are can be used to create standalone descriptors
// where the syntax is only information derived from the parent file.
var (
	SurrogateProto2 = &File***REMOVED***L1: FileL1***REMOVED***Syntax: pref.Proto2***REMOVED***, L2: &FileL2***REMOVED******REMOVED******REMOVED***
	SurrogateProto3 = &File***REMOVED***L1: FileL1***REMOVED***Syntax: pref.Proto3***REMOVED***, L2: &FileL2***REMOVED******REMOVED******REMOVED***
)

type (
	Base struct ***REMOVED***
		L0 BaseL0
	***REMOVED***
	BaseL0 struct ***REMOVED***
		FullName   pref.FullName // must be populated
		ParentFile *File         // must be populated
		Parent     pref.Descriptor
		Index      int
	***REMOVED***
)

func (d *Base) Name() pref.Name         ***REMOVED*** return d.L0.FullName.Name() ***REMOVED***
func (d *Base) FullName() pref.FullName ***REMOVED*** return d.L0.FullName ***REMOVED***
func (d *Base) ParentFile() pref.FileDescriptor ***REMOVED***
	if d.L0.ParentFile == SurrogateProto2 || d.L0.ParentFile == SurrogateProto3 ***REMOVED***
		return nil // surrogate files are not real parents
	***REMOVED***
	return d.L0.ParentFile
***REMOVED***
func (d *Base) Parent() pref.Descriptor             ***REMOVED*** return d.L0.Parent ***REMOVED***
func (d *Base) Index() int                          ***REMOVED*** return d.L0.Index ***REMOVED***
func (d *Base) Syntax() pref.Syntax                 ***REMOVED*** return d.L0.ParentFile.Syntax() ***REMOVED***
func (d *Base) IsPlaceholder() bool                 ***REMOVED*** return false ***REMOVED***
func (d *Base) ProtoInternal(pragma.DoNotImplement) ***REMOVED******REMOVED***

type stringName struct ***REMOVED***
	hasJSON  bool
	once     sync.Once
	nameJSON string
	nameText string
***REMOVED***

// InitJSON initializes the name. It is exported for use by other internal packages.
func (s *stringName) InitJSON(name string) ***REMOVED***
	s.hasJSON = true
	s.nameJSON = name
***REMOVED***

func (s *stringName) lazyInit(fd pref.FieldDescriptor) *stringName ***REMOVED***
	s.once.Do(func() ***REMOVED***
		if fd.IsExtension() ***REMOVED***
			// For extensions, JSON and text are formatted the same way.
			var name string
			if messageset.IsMessageSetExtension(fd) ***REMOVED***
				name = string("[" + fd.FullName().Parent() + "]")
			***REMOVED*** else ***REMOVED***
				name = string("[" + fd.FullName() + "]")
			***REMOVED***
			s.nameJSON = name
			s.nameText = name
		***REMOVED*** else ***REMOVED***
			// Format the JSON name.
			if !s.hasJSON ***REMOVED***
				s.nameJSON = strs.JSONCamelCase(string(fd.Name()))
			***REMOVED***

			// Format the text name.
			s.nameText = string(fd.Name())
			if fd.Kind() == pref.GroupKind ***REMOVED***
				s.nameText = string(fd.Message().Name())
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	return s
***REMOVED***

func (s *stringName) getJSON(fd pref.FieldDescriptor) string ***REMOVED*** return s.lazyInit(fd).nameJSON ***REMOVED***
func (s *stringName) getText(fd pref.FieldDescriptor) string ***REMOVED*** return s.lazyInit(fd).nameText ***REMOVED***

func DefaultValue(v pref.Value, ev pref.EnumValueDescriptor) defaultValue ***REMOVED***
	dv := defaultValue***REMOVED***has: v.IsValid(), val: v, enum: ev***REMOVED***
	if b, ok := v.Interface().([]byte); ok ***REMOVED***
		// Store a copy of the default bytes, so that we can detect
		// accidental mutations of the original value.
		dv.bytes = append([]byte(nil), b...)
	***REMOVED***
	return dv
***REMOVED***

func unmarshalDefault(b []byte, k pref.Kind, pf *File, ed pref.EnumDescriptor) defaultValue ***REMOVED***
	var evs pref.EnumValueDescriptors
	if k == pref.EnumKind ***REMOVED***
		// If the enum is declared within the same file, be careful not to
		// blindly call the Values method, lest we bind ourselves in a deadlock.
		if e, ok := ed.(*Enum); ok && e.L0.ParentFile == pf ***REMOVED***
			evs = &e.L2.Values
		***REMOVED*** else ***REMOVED***
			evs = ed.Values()
		***REMOVED***

		// If we are unable to resolve the enum dependency, use a placeholder
		// enum value since we will not be able to parse the default value.
		if ed.IsPlaceholder() && pref.Name(b).IsValid() ***REMOVED***
			v := pref.ValueOfEnum(0)
			ev := PlaceholderEnumValue(ed.FullName().Parent().Append(pref.Name(b)))
			return DefaultValue(v, ev)
		***REMOVED***
	***REMOVED***

	v, ev, err := defval.Unmarshal(string(b), k, evs, defval.Descriptor)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return DefaultValue(v, ev)
***REMOVED***

type defaultValue struct ***REMOVED***
	has   bool
	val   pref.Value
	enum  pref.EnumValueDescriptor
	bytes []byte
***REMOVED***

func (dv *defaultValue) get(fd pref.FieldDescriptor) pref.Value ***REMOVED***
	// Return the zero value as the default if unpopulated.
	if !dv.has ***REMOVED***
		if fd.Cardinality() == pref.Repeated ***REMOVED***
			return pref.Value***REMOVED******REMOVED***
		***REMOVED***
		switch fd.Kind() ***REMOVED***
		case pref.BoolKind:
			return pref.ValueOfBool(false)
		case pref.Int32Kind, pref.Sint32Kind, pref.Sfixed32Kind:
			return pref.ValueOfInt32(0)
		case pref.Int64Kind, pref.Sint64Kind, pref.Sfixed64Kind:
			return pref.ValueOfInt64(0)
		case pref.Uint32Kind, pref.Fixed32Kind:
			return pref.ValueOfUint32(0)
		case pref.Uint64Kind, pref.Fixed64Kind:
			return pref.ValueOfUint64(0)
		case pref.FloatKind:
			return pref.ValueOfFloat32(0)
		case pref.DoubleKind:
			return pref.ValueOfFloat64(0)
		case pref.StringKind:
			return pref.ValueOfString("")
		case pref.BytesKind:
			return pref.ValueOfBytes(nil)
		case pref.EnumKind:
			if evs := fd.Enum().Values(); evs.Len() > 0 ***REMOVED***
				return pref.ValueOfEnum(evs.Get(0).Number())
			***REMOVED***
			return pref.ValueOfEnum(0)
		***REMOVED***
	***REMOVED***

	if len(dv.bytes) > 0 && !bytes.Equal(dv.bytes, dv.val.Bytes()) ***REMOVED***
		// TODO: Avoid panic if we're running with the race detector
		// and instead spawn a goroutine that periodically resets
		// this value back to the original to induce a race.
		panic(fmt.Sprintf("detected mutation on the default bytes for %v", fd.FullName()))
	***REMOVED***
	return dv.val
***REMOVED***
