// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package descfmt provides functionality to format descriptors.
package descfmt

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"google.golang.org/protobuf/internal/detrand"
	"google.golang.org/protobuf/internal/pragma"
	pref "google.golang.org/protobuf/reflect/protoreflect"
)

type list interface ***REMOVED***
	Len() int
	pragma.DoNotImplement
***REMOVED***

func FormatList(s fmt.State, r rune, vs list) ***REMOVED***
	io.WriteString(s, formatListOpt(vs, true, r == 'v' && (s.Flag('+') || s.Flag('#'))))
***REMOVED***
func formatListOpt(vs list, isRoot, allowMulti bool) string ***REMOVED***
	start, end := "[", "]"
	if isRoot ***REMOVED***
		var name string
		switch vs.(type) ***REMOVED***
		case pref.Names:
			name = "Names"
		case pref.FieldNumbers:
			name = "FieldNumbers"
		case pref.FieldRanges:
			name = "FieldRanges"
		case pref.EnumRanges:
			name = "EnumRanges"
		case pref.FileImports:
			name = "FileImports"
		case pref.Descriptor:
			name = reflect.ValueOf(vs).MethodByName("Get").Type().Out(0).Name() + "s"
		***REMOVED***
		start, end = name+"***REMOVED***", "***REMOVED***"
	***REMOVED***

	var ss []string
	switch vs := vs.(type) ***REMOVED***
	case pref.Names:
		for i := 0; i < vs.Len(); i++ ***REMOVED***
			ss = append(ss, fmt.Sprint(vs.Get(i)))
		***REMOVED***
		return start + joinStrings(ss, false) + end
	case pref.FieldNumbers:
		for i := 0; i < vs.Len(); i++ ***REMOVED***
			ss = append(ss, fmt.Sprint(vs.Get(i)))
		***REMOVED***
		return start + joinStrings(ss, false) + end
	case pref.FieldRanges:
		for i := 0; i < vs.Len(); i++ ***REMOVED***
			r := vs.Get(i)
			if r[0]+1 == r[1] ***REMOVED***
				ss = append(ss, fmt.Sprintf("%d", r[0]))
			***REMOVED*** else ***REMOVED***
				ss = append(ss, fmt.Sprintf("%d:%d", r[0], r[1])) // enum ranges are end exclusive
			***REMOVED***
		***REMOVED***
		return start + joinStrings(ss, false) + end
	case pref.EnumRanges:
		for i := 0; i < vs.Len(); i++ ***REMOVED***
			r := vs.Get(i)
			if r[0] == r[1] ***REMOVED***
				ss = append(ss, fmt.Sprintf("%d", r[0]))
			***REMOVED*** else ***REMOVED***
				ss = append(ss, fmt.Sprintf("%d:%d", r[0], int64(r[1])+1)) // enum ranges are end inclusive
			***REMOVED***
		***REMOVED***
		return start + joinStrings(ss, false) + end
	case pref.FileImports:
		for i := 0; i < vs.Len(); i++ ***REMOVED***
			var rs records
			rs.Append(reflect.ValueOf(vs.Get(i)), "Path", "Package", "IsPublic", "IsWeak")
			ss = append(ss, "***REMOVED***"+rs.Join()+"***REMOVED***")
		***REMOVED***
		return start + joinStrings(ss, allowMulti) + end
	default:
		_, isEnumValue := vs.(pref.EnumValueDescriptors)
		for i := 0; i < vs.Len(); i++ ***REMOVED***
			m := reflect.ValueOf(vs).MethodByName("Get")
			v := m.Call([]reflect.Value***REMOVED***reflect.ValueOf(i)***REMOVED***)[0].Interface()
			ss = append(ss, formatDescOpt(v.(pref.Descriptor), false, allowMulti && !isEnumValue))
		***REMOVED***
		return start + joinStrings(ss, allowMulti && isEnumValue) + end
	***REMOVED***
***REMOVED***

// descriptorAccessors is a list of accessors to print for each descriptor.
//
// Do not print all accessors since some contain redundant information,
// while others are pointers that we do not want to follow since the descriptor
// is actually a cyclic graph.
//
// Using a list allows us to print the accessors in a sensible order.
var descriptorAccessors = map[reflect.Type][]string***REMOVED***
	reflect.TypeOf((*pref.FileDescriptor)(nil)).Elem():      ***REMOVED***"Path", "Package", "Imports", "Messages", "Enums", "Extensions", "Services"***REMOVED***,
	reflect.TypeOf((*pref.MessageDescriptor)(nil)).Elem():   ***REMOVED***"IsMapEntry", "Fields", "Oneofs", "ReservedNames", "ReservedRanges", "RequiredNumbers", "ExtensionRanges", "Messages", "Enums", "Extensions"***REMOVED***,
	reflect.TypeOf((*pref.FieldDescriptor)(nil)).Elem():     ***REMOVED***"Number", "Cardinality", "Kind", "HasJSONName", "JSONName", "HasPresence", "IsExtension", "IsPacked", "IsWeak", "IsList", "IsMap", "MapKey", "MapValue", "HasDefault", "Default", "ContainingOneof", "ContainingMessage", "Message", "Enum"***REMOVED***,
	reflect.TypeOf((*pref.OneofDescriptor)(nil)).Elem():     ***REMOVED***"Fields"***REMOVED***, // not directly used; must keep in sync with formatDescOpt
	reflect.TypeOf((*pref.EnumDescriptor)(nil)).Elem():      ***REMOVED***"Values", "ReservedNames", "ReservedRanges"***REMOVED***,
	reflect.TypeOf((*pref.EnumValueDescriptor)(nil)).Elem(): ***REMOVED***"Number"***REMOVED***,
	reflect.TypeOf((*pref.ServiceDescriptor)(nil)).Elem():   ***REMOVED***"Methods"***REMOVED***,
	reflect.TypeOf((*pref.MethodDescriptor)(nil)).Elem():    ***REMOVED***"Input", "Output", "IsStreamingClient", "IsStreamingServer"***REMOVED***,
***REMOVED***

func FormatDesc(s fmt.State, r rune, t pref.Descriptor) ***REMOVED***
	io.WriteString(s, formatDescOpt(t, true, r == 'v' && (s.Flag('+') || s.Flag('#'))))
***REMOVED***
func formatDescOpt(t pref.Descriptor, isRoot, allowMulti bool) string ***REMOVED***
	rv := reflect.ValueOf(t)
	rt := rv.MethodByName("ProtoType").Type().In(0)

	start, end := "***REMOVED***", "***REMOVED***"
	if isRoot ***REMOVED***
		start = rt.Name() + "***REMOVED***"
	***REMOVED***

	_, isFile := t.(pref.FileDescriptor)
	rs := records***REMOVED***allowMulti: allowMulti***REMOVED***
	if t.IsPlaceholder() ***REMOVED***
		if isFile ***REMOVED***
			rs.Append(rv, "Path", "Package", "IsPlaceholder")
		***REMOVED*** else ***REMOVED***
			rs.Append(rv, "FullName", "IsPlaceholder")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		switch ***REMOVED***
		case isFile:
			rs.Append(rv, "Syntax")
		case isRoot:
			rs.Append(rv, "Syntax", "FullName")
		default:
			rs.Append(rv, "Name")
		***REMOVED***
		switch t := t.(type) ***REMOVED***
		case pref.FieldDescriptor:
			for _, s := range descriptorAccessors[rt] ***REMOVED***
				switch s ***REMOVED***
				case "MapKey":
					if k := t.MapKey(); k != nil ***REMOVED***
						rs.recs = append(rs.recs, [2]string***REMOVED***"MapKey", k.Kind().String()***REMOVED***)
					***REMOVED***
				case "MapValue":
					if v := t.MapValue(); v != nil ***REMOVED***
						switch v.Kind() ***REMOVED***
						case pref.EnumKind:
							rs.recs = append(rs.recs, [2]string***REMOVED***"MapValue", string(v.Enum().FullName())***REMOVED***)
						case pref.MessageKind, pref.GroupKind:
							rs.recs = append(rs.recs, [2]string***REMOVED***"MapValue", string(v.Message().FullName())***REMOVED***)
						default:
							rs.recs = append(rs.recs, [2]string***REMOVED***"MapValue", v.Kind().String()***REMOVED***)
						***REMOVED***
					***REMOVED***
				case "ContainingOneof":
					if od := t.ContainingOneof(); od != nil ***REMOVED***
						rs.recs = append(rs.recs, [2]string***REMOVED***"Oneof", string(od.Name())***REMOVED***)
					***REMOVED***
				case "ContainingMessage":
					if t.IsExtension() ***REMOVED***
						rs.recs = append(rs.recs, [2]string***REMOVED***"Extendee", string(t.ContainingMessage().FullName())***REMOVED***)
					***REMOVED***
				case "Message":
					if !t.IsMap() ***REMOVED***
						rs.Append(rv, s)
					***REMOVED***
				default:
					rs.Append(rv, s)
				***REMOVED***
			***REMOVED***
		case pref.OneofDescriptor:
			var ss []string
			fs := t.Fields()
			for i := 0; i < fs.Len(); i++ ***REMOVED***
				ss = append(ss, string(fs.Get(i).Name()))
			***REMOVED***
			if len(ss) > 0 ***REMOVED***
				rs.recs = append(rs.recs, [2]string***REMOVED***"Fields", "[" + joinStrings(ss, false) + "]"***REMOVED***)
			***REMOVED***
		default:
			rs.Append(rv, descriptorAccessors[rt]...)
		***REMOVED***
		if rv.MethodByName("GoType").IsValid() ***REMOVED***
			rs.Append(rv, "GoType")
		***REMOVED***
	***REMOVED***
	return start + rs.Join() + end
***REMOVED***

type records struct ***REMOVED***
	recs       [][2]string
	allowMulti bool
***REMOVED***

func (rs *records) Append(v reflect.Value, accessors ...string) ***REMOVED***
	for _, a := range accessors ***REMOVED***
		var rv reflect.Value
		if m := v.MethodByName(a); m.IsValid() ***REMOVED***
			rv = m.Call(nil)[0]
		***REMOVED***
		if v.Kind() == reflect.Struct && !rv.IsValid() ***REMOVED***
			rv = v.FieldByName(a)
		***REMOVED***
		if !rv.IsValid() ***REMOVED***
			panic(fmt.Sprintf("unknown accessor: %v.%s", v.Type(), a))
		***REMOVED***
		if _, ok := rv.Interface().(pref.Value); ok ***REMOVED***
			rv = rv.MethodByName("Interface").Call(nil)[0]
			if !rv.IsNil() ***REMOVED***
				rv = rv.Elem()
			***REMOVED***
		***REMOVED***

		// Ignore zero values.
		var isZero bool
		switch rv.Kind() ***REMOVED***
		case reflect.Interface, reflect.Slice:
			isZero = rv.IsNil()
		case reflect.Bool:
			isZero = rv.Bool() == false
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			isZero = rv.Int() == 0
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			isZero = rv.Uint() == 0
		case reflect.String:
			isZero = rv.String() == ""
		***REMOVED***
		if n, ok := rv.Interface().(list); ok ***REMOVED***
			isZero = n.Len() == 0
		***REMOVED***
		if isZero ***REMOVED***
			continue
		***REMOVED***

		// Format the value.
		var s string
		v := rv.Interface()
		switch v := v.(type) ***REMOVED***
		case list:
			s = formatListOpt(v, false, rs.allowMulti)
		case pref.FieldDescriptor, pref.OneofDescriptor, pref.EnumValueDescriptor, pref.MethodDescriptor:
			s = string(v.(pref.Descriptor).Name())
		case pref.Descriptor:
			s = string(v.FullName())
		case string:
			s = strconv.Quote(v)
		case []byte:
			s = fmt.Sprintf("%q", v)
		default:
			s = fmt.Sprint(v)
		***REMOVED***
		rs.recs = append(rs.recs, [2]string***REMOVED***a, s***REMOVED***)
	***REMOVED***
***REMOVED***

func (rs *records) Join() string ***REMOVED***
	var ss []string

	// In single line mode, simply join all records with commas.
	if !rs.allowMulti ***REMOVED***
		for _, r := range rs.recs ***REMOVED***
			ss = append(ss, r[0]+formatColon(0)+r[1])
		***REMOVED***
		return joinStrings(ss, false)
	***REMOVED***

	// In allowMulti line mode, align single line records for more readable output.
	var maxLen int
	flush := func(i int) ***REMOVED***
		for _, r := range rs.recs[len(ss):i] ***REMOVED***
			ss = append(ss, r[0]+formatColon(maxLen-len(r[0]))+r[1])
		***REMOVED***
		maxLen = 0
	***REMOVED***
	for i, r := range rs.recs ***REMOVED***
		if isMulti := strings.Contains(r[1], "\n"); isMulti ***REMOVED***
			flush(i)
			ss = append(ss, r[0]+formatColon(0)+strings.Join(strings.Split(r[1], "\n"), "\n\t"))
		***REMOVED*** else if maxLen < len(r[0]) ***REMOVED***
			maxLen = len(r[0])
		***REMOVED***
	***REMOVED***
	flush(len(rs.recs))
	return joinStrings(ss, true)
***REMOVED***

func formatColon(padding int) string ***REMOVED***
	// Deliberately introduce instability into the debug output to
	// discourage users from performing string comparisons.
	// This provides us flexibility to change the output in the future.
	if detrand.Bool() ***REMOVED***
		return ":" + strings.Repeat(" ", 1+padding) // use non-breaking spaces (U+00a0)
	***REMOVED*** else ***REMOVED***
		return ":" + strings.Repeat(" ", 1+padding) // use regular spaces (U+0020)
	***REMOVED***
***REMOVED***

func joinStrings(ss []string, isMulti bool) string ***REMOVED***
	if len(ss) == 0 ***REMOVED***
		return ""
	***REMOVED***
	if isMulti ***REMOVED***
		return "\n\t" + strings.Join(ss, "\n\t") + "\n"
	***REMOVED***
	return strings.Join(ss, ", ")
***REMOVED***
