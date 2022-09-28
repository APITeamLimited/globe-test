// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
)

// StructProperties represents protocol buffer type information for a
// generated protobuf message in the open-struct API.
//
// Deprecated: Do not use.
type StructProperties struct ***REMOVED***
	// Prop are the properties for each field.
	//
	// Fields belonging to a oneof are stored in OneofTypes instead, with a
	// single Properties representing the parent oneof held here.
	//
	// The order of Prop matches the order of fields in the Go struct.
	// Struct fields that are not related to protobufs have a "XXX_" prefix
	// in the Properties.Name and must be ignored by the user.
	Prop []*Properties

	// OneofTypes contains information about the oneof fields in this message.
	// It is keyed by the protobuf field name.
	OneofTypes map[string]*OneofProperties
***REMOVED***

// Properties represents the type information for a protobuf message field.
//
// Deprecated: Do not use.
type Properties struct ***REMOVED***
	// Name is a placeholder name with little meaningful semantic value.
	// If the name has an "XXX_" prefix, the entire Properties must be ignored.
	Name string
	// OrigName is the protobuf field name or oneof name.
	OrigName string
	// JSONName is the JSON name for the protobuf field.
	JSONName string
	// Enum is a placeholder name for enums.
	// For historical reasons, this is neither the Go name for the enum,
	// nor the protobuf name for the enum.
	Enum string // Deprecated: Do not use.
	// Weak contains the full name of the weakly referenced message.
	Weak string
	// Wire is a string representation of the wire type.
	Wire string
	// WireType is the protobuf wire type for the field.
	WireType int
	// Tag is the protobuf field number.
	Tag int
	// Required reports whether this is a required field.
	Required bool
	// Optional reports whether this is a optional field.
	Optional bool
	// Repeated reports whether this is a repeated field.
	Repeated bool
	// Packed reports whether this is a packed repeated field of scalars.
	Packed bool
	// Proto3 reports whether this field operates under the proto3 syntax.
	Proto3 bool
	// Oneof reports whether this field belongs within a oneof.
	Oneof bool

	// Default is the default value in string form.
	Default string
	// HasDefault reports whether the field has a default value.
	HasDefault bool

	// MapKeyProp is the properties for the key field for a map field.
	MapKeyProp *Properties
	// MapValProp is the properties for the value field for a map field.
	MapValProp *Properties
***REMOVED***

// OneofProperties represents the type information for a protobuf oneof.
//
// Deprecated: Do not use.
type OneofProperties struct ***REMOVED***
	// Type is a pointer to the generated wrapper type for the field value.
	// This is nil for messages that are not in the open-struct API.
	Type reflect.Type
	// Field is the index into StructProperties.Prop for the containing oneof.
	Field int
	// Prop is the properties for the field.
	Prop *Properties
***REMOVED***

// String formats the properties in the protobuf struct field tag style.
func (p *Properties) String() string ***REMOVED***
	s := p.Wire
	s += "," + strconv.Itoa(p.Tag)
	if p.Required ***REMOVED***
		s += ",req"
	***REMOVED***
	if p.Optional ***REMOVED***
		s += ",opt"
	***REMOVED***
	if p.Repeated ***REMOVED***
		s += ",rep"
	***REMOVED***
	if p.Packed ***REMOVED***
		s += ",packed"
	***REMOVED***
	s += ",name=" + p.OrigName
	if p.JSONName != "" ***REMOVED***
		s += ",json=" + p.JSONName
	***REMOVED***
	if len(p.Enum) > 0 ***REMOVED***
		s += ",enum=" + p.Enum
	***REMOVED***
	if len(p.Weak) > 0 ***REMOVED***
		s += ",weak=" + p.Weak
	***REMOVED***
	if p.Proto3 ***REMOVED***
		s += ",proto3"
	***REMOVED***
	if p.Oneof ***REMOVED***
		s += ",oneof"
	***REMOVED***
	if p.HasDefault ***REMOVED***
		s += ",def=" + p.Default
	***REMOVED***
	return s
***REMOVED***

// Parse populates p by parsing a string in the protobuf struct field tag style.
func (p *Properties) Parse(tag string) ***REMOVED***
	// For example: "bytes,49,opt,name=foo,def=hello!"
	for len(tag) > 0 ***REMOVED***
		i := strings.IndexByte(tag, ',')
		if i < 0 ***REMOVED***
			i = len(tag)
		***REMOVED***
		switch s := tag[:i]; ***REMOVED***
		case strings.HasPrefix(s, "name="):
			p.OrigName = s[len("name="):]
		case strings.HasPrefix(s, "json="):
			p.JSONName = s[len("json="):]
		case strings.HasPrefix(s, "enum="):
			p.Enum = s[len("enum="):]
		case strings.HasPrefix(s, "weak="):
			p.Weak = s[len("weak="):]
		case strings.Trim(s, "0123456789") == "":
			n, _ := strconv.ParseUint(s, 10, 32)
			p.Tag = int(n)
		case s == "opt":
			p.Optional = true
		case s == "req":
			p.Required = true
		case s == "rep":
			p.Repeated = true
		case s == "varint" || s == "zigzag32" || s == "zigzag64":
			p.Wire = s
			p.WireType = WireVarint
		case s == "fixed32":
			p.Wire = s
			p.WireType = WireFixed32
		case s == "fixed64":
			p.Wire = s
			p.WireType = WireFixed64
		case s == "bytes":
			p.Wire = s
			p.WireType = WireBytes
		case s == "group":
			p.Wire = s
			p.WireType = WireStartGroup
		case s == "packed":
			p.Packed = true
		case s == "proto3":
			p.Proto3 = true
		case s == "oneof":
			p.Oneof = true
		case strings.HasPrefix(s, "def="):
			// The default tag is special in that everything afterwards is the
			// default regardless of the presence of commas.
			p.HasDefault = true
			p.Default, i = tag[len("def="):], len(tag)
		***REMOVED***
		tag = strings.TrimPrefix(tag[i:], ",")
	***REMOVED***
***REMOVED***

// Init populates the properties from a protocol buffer struct tag.
//
// Deprecated: Do not use.
func (p *Properties) Init(typ reflect.Type, name, tag string, f *reflect.StructField) ***REMOVED***
	p.Name = name
	p.OrigName = name
	if tag == "" ***REMOVED***
		return
	***REMOVED***
	p.Parse(tag)

	if typ != nil && typ.Kind() == reflect.Map ***REMOVED***
		p.MapKeyProp = new(Properties)
		p.MapKeyProp.Init(nil, "Key", f.Tag.Get("protobuf_key"), nil)
		p.MapValProp = new(Properties)
		p.MapValProp.Init(nil, "Value", f.Tag.Get("protobuf_val"), nil)
	***REMOVED***
***REMOVED***

var propertiesCache sync.Map // map[reflect.Type]*StructProperties

// GetProperties returns the list of properties for the type represented by t,
// which must be a generated protocol buffer message in the open-struct API,
// where protobuf message fields are represented by exported Go struct fields.
//
// Deprecated: Use protobuf reflection instead.
func GetProperties(t reflect.Type) *StructProperties ***REMOVED***
	if p, ok := propertiesCache.Load(t); ok ***REMOVED***
		return p.(*StructProperties)
	***REMOVED***
	p, _ := propertiesCache.LoadOrStore(t, newProperties(t))
	return p.(*StructProperties)
***REMOVED***

func newProperties(t reflect.Type) *StructProperties ***REMOVED***
	if t.Kind() != reflect.Struct ***REMOVED***
		panic(fmt.Sprintf("%v is not a generated message in the open-struct API", t))
	***REMOVED***

	var hasOneof bool
	prop := new(StructProperties)

	// Construct a list of properties for each field in the struct.
	for i := 0; i < t.NumField(); i++ ***REMOVED***
		p := new(Properties)
		f := t.Field(i)
		tagField := f.Tag.Get("protobuf")
		p.Init(f.Type, f.Name, tagField, &f)

		tagOneof := f.Tag.Get("protobuf_oneof")
		if tagOneof != "" ***REMOVED***
			hasOneof = true
			p.OrigName = tagOneof
		***REMOVED***

		// Rename unrelated struct fields with the "XXX_" prefix since so much
		// user code simply checks for this to exclude special fields.
		if tagField == "" && tagOneof == "" && !strings.HasPrefix(p.Name, "XXX_") ***REMOVED***
			p.Name = "XXX_" + p.Name
			p.OrigName = "XXX_" + p.OrigName
		***REMOVED*** else if p.Weak != "" ***REMOVED***
			p.Name = p.OrigName // avoid possible "XXX_" prefix on weak field
		***REMOVED***

		prop.Prop = append(prop.Prop, p)
	***REMOVED***

	// Construct a mapping of oneof field names to properties.
	if hasOneof ***REMOVED***
		var oneofWrappers []interface***REMOVED******REMOVED***
		if fn, ok := reflect.PtrTo(t).MethodByName("XXX_OneofFuncs"); ok ***REMOVED***
			oneofWrappers = fn.Func.Call([]reflect.Value***REMOVED***reflect.Zero(fn.Type.In(0))***REMOVED***)[3].Interface().([]interface***REMOVED******REMOVED***)
		***REMOVED***
		if fn, ok := reflect.PtrTo(t).MethodByName("XXX_OneofWrappers"); ok ***REMOVED***
			oneofWrappers = fn.Func.Call([]reflect.Value***REMOVED***reflect.Zero(fn.Type.In(0))***REMOVED***)[0].Interface().([]interface***REMOVED******REMOVED***)
		***REMOVED***
		if m, ok := reflect.Zero(reflect.PtrTo(t)).Interface().(protoreflect.ProtoMessage); ok ***REMOVED***
			if m, ok := m.ProtoReflect().(interface***REMOVED*** ProtoMessageInfo() *protoimpl.MessageInfo ***REMOVED***); ok ***REMOVED***
				oneofWrappers = m.ProtoMessageInfo().OneofWrappers
			***REMOVED***
		***REMOVED***

		prop.OneofTypes = make(map[string]*OneofProperties)
		for _, wrapper := range oneofWrappers ***REMOVED***
			p := &OneofProperties***REMOVED***
				Type: reflect.ValueOf(wrapper).Type(), // *T
				Prop: new(Properties),
			***REMOVED***
			f := p.Type.Elem().Field(0)
			p.Prop.Name = f.Name
			p.Prop.Parse(f.Tag.Get("protobuf"))

			// Determine the struct field that contains this oneof.
			// Each wrapper is assignable to exactly one parent field.
			var foundOneof bool
			for i := 0; i < t.NumField() && !foundOneof; i++ ***REMOVED***
				if p.Type.AssignableTo(t.Field(i).Type) ***REMOVED***
					p.Field = i
					foundOneof = true
				***REMOVED***
			***REMOVED***
			if !foundOneof ***REMOVED***
				panic(fmt.Sprintf("%v is not a generated message in the open-struct API", t))
			***REMOVED***
			prop.OneofTypes[p.Prop.OrigName] = p
		***REMOVED***
	***REMOVED***

	return prop
***REMOVED***

func (sp *StructProperties) Len() int           ***REMOVED*** return len(sp.Prop) ***REMOVED***
func (sp *StructProperties) Less(i, j int) bool ***REMOVED*** return false ***REMOVED***
func (sp *StructProperties) Swap(i, j int)      ***REMOVED*** return ***REMOVED***
