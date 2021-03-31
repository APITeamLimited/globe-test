// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"reflect"

	"google.golang.org/protobuf/internal/descopts"
	"google.golang.org/protobuf/internal/encoding/messageset"
	ptag "google.golang.org/protobuf/internal/encoding/tag"
	"google.golang.org/protobuf/internal/filedesc"
	"google.golang.org/protobuf/internal/pragma"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	preg "google.golang.org/protobuf/reflect/protoregistry"
	piface "google.golang.org/protobuf/runtime/protoiface"
)

func (xi *ExtensionInfo) initToLegacy() ***REMOVED***
	xd := xi.desc
	var parent piface.MessageV1
	messageName := xd.ContainingMessage().FullName()
	if mt, _ := preg.GlobalTypes.FindMessageByName(messageName); mt != nil ***REMOVED***
		// Create a new parent message and unwrap it if possible.
		mv := mt.New().Interface()
		t := reflect.TypeOf(mv)
		if mv, ok := mv.(unwrapper); ok ***REMOVED***
			t = reflect.TypeOf(mv.protoUnwrap())
		***REMOVED***

		// Check whether the message implements the legacy v1 Message interface.
		mz := reflect.Zero(t).Interface()
		if mz, ok := mz.(piface.MessageV1); ok ***REMOVED***
			parent = mz
		***REMOVED***
	***REMOVED***

	// Determine the v1 extension type, which is unfortunately not the same as
	// the v2 ExtensionType.GoType.
	extType := xi.goType
	switch extType.Kind() ***REMOVED***
	case reflect.Bool, reflect.Int32, reflect.Int64, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:
		extType = reflect.PtrTo(extType) // T -> *T for singular scalar fields
	***REMOVED***

	// Reconstruct the legacy enum full name.
	var enumName string
	if xd.Kind() == pref.EnumKind ***REMOVED***
		enumName = legacyEnumName(xd.Enum())
	***REMOVED***

	// Derive the proto file that the extension was declared within.
	var filename string
	if fd := xd.ParentFile(); fd != nil ***REMOVED***
		filename = fd.Path()
	***REMOVED***

	// For MessageSet extensions, the name used is the parent message.
	name := xd.FullName()
	if messageset.IsMessageSetExtension(xd) ***REMOVED***
		name = name.Parent()
	***REMOVED***

	xi.ExtendedType = parent
	xi.ExtensionType = reflect.Zero(extType).Interface()
	xi.Field = int32(xd.Number())
	xi.Name = string(name)
	xi.Tag = ptag.Marshal(xd, enumName)
	xi.Filename = filename
***REMOVED***

// initFromLegacy initializes an ExtensionInfo from
// the contents of the deprecated exported fields of the type.
func (xi *ExtensionInfo) initFromLegacy() ***REMOVED***
	// The v1 API returns "type incomplete" descriptors where only the
	// field number is specified. In such a case, use a placeholder.
	if xi.ExtendedType == nil || xi.ExtensionType == nil ***REMOVED***
		xd := placeholderExtension***REMOVED***
			name:   pref.FullName(xi.Name),
			number: pref.FieldNumber(xi.Field),
		***REMOVED***
		xi.desc = extensionTypeDescriptor***REMOVED***xd, xi***REMOVED***
		return
	***REMOVED***

	// Resolve enum or message dependencies.
	var ed pref.EnumDescriptor
	var md pref.MessageDescriptor
	t := reflect.TypeOf(xi.ExtensionType)
	isOptional := t.Kind() == reflect.Ptr && t.Elem().Kind() != reflect.Struct
	isRepeated := t.Kind() == reflect.Slice && t.Elem().Kind() != reflect.Uint8
	if isOptional || isRepeated ***REMOVED***
		t = t.Elem()
	***REMOVED***
	switch v := reflect.Zero(t).Interface().(type) ***REMOVED***
	case pref.Enum:
		ed = v.Descriptor()
	case enumV1:
		ed = LegacyLoadEnumDesc(t)
	case pref.ProtoMessage:
		md = v.ProtoReflect().Descriptor()
	case messageV1:
		md = LegacyLoadMessageDesc(t)
	***REMOVED***

	// Derive basic field information from the struct tag.
	var evs pref.EnumValueDescriptors
	if ed != nil ***REMOVED***
		evs = ed.Values()
	***REMOVED***
	fd := ptag.Unmarshal(xi.Tag, t, evs).(*filedesc.Field)

	// Construct a v2 ExtensionType.
	xd := &filedesc.Extension***REMOVED***L2: new(filedesc.ExtensionL2)***REMOVED***
	xd.L0.ParentFile = filedesc.SurrogateProto2
	xd.L0.FullName = pref.FullName(xi.Name)
	xd.L1.Number = pref.FieldNumber(xi.Field)
	xd.L1.Cardinality = fd.L1.Cardinality
	xd.L1.Kind = fd.L1.Kind
	xd.L2.IsPacked = fd.L1.IsPacked
	xd.L2.Default = fd.L1.Default
	xd.L1.Extendee = Export***REMOVED******REMOVED***.MessageDescriptorOf(xi.ExtendedType)
	xd.L2.Enum = ed
	xd.L2.Message = md

	// Derive real extension field name for MessageSets.
	if messageset.IsMessageSet(xd.L1.Extendee) && md.FullName() == xd.L0.FullName ***REMOVED***
		xd.L0.FullName = xd.L0.FullName.Append(messageset.ExtensionName)
	***REMOVED***

	tt := reflect.TypeOf(xi.ExtensionType)
	if isOptional ***REMOVED***
		tt = tt.Elem()
	***REMOVED***
	xi.goType = tt
	xi.desc = extensionTypeDescriptor***REMOVED***xd, xi***REMOVED***
***REMOVED***

type placeholderExtension struct ***REMOVED***
	name   pref.FullName
	number pref.FieldNumber
***REMOVED***

func (x placeholderExtension) ParentFile() pref.FileDescriptor            ***REMOVED*** return nil ***REMOVED***
func (x placeholderExtension) Parent() pref.Descriptor                    ***REMOVED*** return nil ***REMOVED***
func (x placeholderExtension) Index() int                                 ***REMOVED*** return 0 ***REMOVED***
func (x placeholderExtension) Syntax() pref.Syntax                        ***REMOVED*** return 0 ***REMOVED***
func (x placeholderExtension) Name() pref.Name                            ***REMOVED*** return x.name.Name() ***REMOVED***
func (x placeholderExtension) FullName() pref.FullName                    ***REMOVED*** return x.name ***REMOVED***
func (x placeholderExtension) IsPlaceholder() bool                        ***REMOVED*** return true ***REMOVED***
func (x placeholderExtension) Options() pref.ProtoMessage                 ***REMOVED*** return descopts.Field ***REMOVED***
func (x placeholderExtension) Number() pref.FieldNumber                   ***REMOVED*** return x.number ***REMOVED***
func (x placeholderExtension) Cardinality() pref.Cardinality              ***REMOVED*** return 0 ***REMOVED***
func (x placeholderExtension) Kind() pref.Kind                            ***REMOVED*** return 0 ***REMOVED***
func (x placeholderExtension) HasJSONName() bool                          ***REMOVED*** return false ***REMOVED***
func (x placeholderExtension) JSONName() string                           ***REMOVED*** return "[" + string(x.name) + "]" ***REMOVED***
func (x placeholderExtension) TextName() string                           ***REMOVED*** return "[" + string(x.name) + "]" ***REMOVED***
func (x placeholderExtension) HasPresence() bool                          ***REMOVED*** return false ***REMOVED***
func (x placeholderExtension) HasOptionalKeyword() bool                   ***REMOVED*** return false ***REMOVED***
func (x placeholderExtension) IsExtension() bool                          ***REMOVED*** return true ***REMOVED***
func (x placeholderExtension) IsWeak() bool                               ***REMOVED*** return false ***REMOVED***
func (x placeholderExtension) IsPacked() bool                             ***REMOVED*** return false ***REMOVED***
func (x placeholderExtension) IsList() bool                               ***REMOVED*** return false ***REMOVED***
func (x placeholderExtension) IsMap() bool                                ***REMOVED*** return false ***REMOVED***
func (x placeholderExtension) MapKey() pref.FieldDescriptor               ***REMOVED*** return nil ***REMOVED***
func (x placeholderExtension) MapValue() pref.FieldDescriptor             ***REMOVED*** return nil ***REMOVED***
func (x placeholderExtension) HasDefault() bool                           ***REMOVED*** return false ***REMOVED***
func (x placeholderExtension) Default() pref.Value                        ***REMOVED*** return pref.Value***REMOVED******REMOVED*** ***REMOVED***
func (x placeholderExtension) DefaultEnumValue() pref.EnumValueDescriptor ***REMOVED*** return nil ***REMOVED***
func (x placeholderExtension) ContainingOneof() pref.OneofDescriptor      ***REMOVED*** return nil ***REMOVED***
func (x placeholderExtension) ContainingMessage() pref.MessageDescriptor  ***REMOVED*** return nil ***REMOVED***
func (x placeholderExtension) Enum() pref.EnumDescriptor                  ***REMOVED*** return nil ***REMOVED***
func (x placeholderExtension) Message() pref.MessageDescriptor            ***REMOVED*** return nil ***REMOVED***
func (x placeholderExtension) ProtoType(pref.FieldDescriptor)             ***REMOVED*** return ***REMOVED***
func (x placeholderExtension) ProtoInternal(pragma.DoNotImplement)        ***REMOVED*** return ***REMOVED***
