// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"google.golang.org/protobuf/internal/filedesc"
	"google.golang.org/protobuf/internal/strs"
	"google.golang.org/protobuf/reflect/protoreflect"
	pref "google.golang.org/protobuf/reflect/protoreflect"
)

// legacyEnumName returns the name of enums used in legacy code.
// It is neither the protobuf full name nor the qualified Go name,
// but rather an odd hybrid of both.
func legacyEnumName(ed pref.EnumDescriptor) string ***REMOVED***
	var protoPkg string
	enumName := string(ed.FullName())
	if fd := ed.ParentFile(); fd != nil ***REMOVED***
		protoPkg = string(fd.Package())
		enumName = strings.TrimPrefix(enumName, protoPkg+".")
	***REMOVED***
	if protoPkg == "" ***REMOVED***
		return strs.GoCamelCase(enumName)
	***REMOVED***
	return protoPkg + "." + strs.GoCamelCase(enumName)
***REMOVED***

// legacyWrapEnum wraps v as a protoreflect.Enum,
// where v must be a int32 kind and not implement the v2 API already.
func legacyWrapEnum(v reflect.Value) pref.Enum ***REMOVED***
	et := legacyLoadEnumType(v.Type())
	return et.New(pref.EnumNumber(v.Int()))
***REMOVED***

var legacyEnumTypeCache sync.Map // map[reflect.Type]protoreflect.EnumType

// legacyLoadEnumType dynamically loads a protoreflect.EnumType for t,
// where t must be an int32 kind and not implement the v2 API already.
func legacyLoadEnumType(t reflect.Type) pref.EnumType ***REMOVED***
	// Fast-path: check if a EnumType is cached for this concrete type.
	if et, ok := legacyEnumTypeCache.Load(t); ok ***REMOVED***
		return et.(pref.EnumType)
	***REMOVED***

	// Slow-path: derive enum descriptor and initialize EnumType.
	var et pref.EnumType
	ed := LegacyLoadEnumDesc(t)
	et = &legacyEnumType***REMOVED***
		desc:   ed,
		goType: t,
	***REMOVED***
	if et, ok := legacyEnumTypeCache.LoadOrStore(t, et); ok ***REMOVED***
		return et.(pref.EnumType)
	***REMOVED***
	return et
***REMOVED***

type legacyEnumType struct ***REMOVED***
	desc   pref.EnumDescriptor
	goType reflect.Type
	m      sync.Map // map[protoreflect.EnumNumber]proto.Enum
***REMOVED***

func (t *legacyEnumType) New(n pref.EnumNumber) pref.Enum ***REMOVED***
	if e, ok := t.m.Load(n); ok ***REMOVED***
		return e.(pref.Enum)
	***REMOVED***
	e := &legacyEnumWrapper***REMOVED***num: n, pbTyp: t, goTyp: t.goType***REMOVED***
	t.m.Store(n, e)
	return e
***REMOVED***
func (t *legacyEnumType) Descriptor() pref.EnumDescriptor ***REMOVED***
	return t.desc
***REMOVED***

type legacyEnumWrapper struct ***REMOVED***
	num   pref.EnumNumber
	pbTyp pref.EnumType
	goTyp reflect.Type
***REMOVED***

func (e *legacyEnumWrapper) Descriptor() pref.EnumDescriptor ***REMOVED***
	return e.pbTyp.Descriptor()
***REMOVED***
func (e *legacyEnumWrapper) Type() pref.EnumType ***REMOVED***
	return e.pbTyp
***REMOVED***
func (e *legacyEnumWrapper) Number() pref.EnumNumber ***REMOVED***
	return e.num
***REMOVED***
func (e *legacyEnumWrapper) ProtoReflect() pref.Enum ***REMOVED***
	return e
***REMOVED***
func (e *legacyEnumWrapper) protoUnwrap() interface***REMOVED******REMOVED*** ***REMOVED***
	v := reflect.New(e.goTyp).Elem()
	v.SetInt(int64(e.num))
	return v.Interface()
***REMOVED***

var (
	_ pref.Enum = (*legacyEnumWrapper)(nil)
	_ unwrapper = (*legacyEnumWrapper)(nil)
)

var legacyEnumDescCache sync.Map // map[reflect.Type]protoreflect.EnumDescriptor

// LegacyLoadEnumDesc returns an EnumDescriptor derived from the Go type,
// which must be an int32 kind and not implement the v2 API already.
//
// This is exported for testing purposes.
func LegacyLoadEnumDesc(t reflect.Type) pref.EnumDescriptor ***REMOVED***
	// Fast-path: check if an EnumDescriptor is cached for this concrete type.
	if ed, ok := legacyEnumDescCache.Load(t); ok ***REMOVED***
		return ed.(pref.EnumDescriptor)
	***REMOVED***

	// Slow-path: initialize EnumDescriptor from the raw descriptor.
	ev := reflect.Zero(t).Interface()
	if _, ok := ev.(pref.Enum); ok ***REMOVED***
		panic(fmt.Sprintf("%v already implements proto.Enum", t))
	***REMOVED***
	edV1, ok := ev.(enumV1)
	if !ok ***REMOVED***
		return aberrantLoadEnumDesc(t)
	***REMOVED***
	b, idxs := edV1.EnumDescriptor()

	var ed pref.EnumDescriptor
	if len(idxs) == 1 ***REMOVED***
		ed = legacyLoadFileDesc(b).Enums().Get(idxs[0])
	***REMOVED*** else ***REMOVED***
		md := legacyLoadFileDesc(b).Messages().Get(idxs[0])
		for _, i := range idxs[1 : len(idxs)-1] ***REMOVED***
			md = md.Messages().Get(i)
		***REMOVED***
		ed = md.Enums().Get(idxs[len(idxs)-1])
	***REMOVED***
	if ed, ok := legacyEnumDescCache.LoadOrStore(t, ed); ok ***REMOVED***
		return ed.(protoreflect.EnumDescriptor)
	***REMOVED***
	return ed
***REMOVED***

var aberrantEnumDescCache sync.Map // map[reflect.Type]protoreflect.EnumDescriptor

// aberrantLoadEnumDesc returns an EnumDescriptor derived from the Go type,
// which must not implement protoreflect.Enum or enumV1.
//
// If the type does not implement enumV1, then there is no reliable
// way to derive the original protobuf type information.
// We are unable to use the global enum registry since it is
// unfortunately keyed by the protobuf full name, which we also do not know.
// Thus, this produces some bogus enum descriptor based on the Go type name.
func aberrantLoadEnumDesc(t reflect.Type) pref.EnumDescriptor ***REMOVED***
	// Fast-path: check if an EnumDescriptor is cached for this concrete type.
	if ed, ok := aberrantEnumDescCache.Load(t); ok ***REMOVED***
		return ed.(pref.EnumDescriptor)
	***REMOVED***

	// Slow-path: construct a bogus, but unique EnumDescriptor.
	ed := &filedesc.Enum***REMOVED***L2: new(filedesc.EnumL2)***REMOVED***
	ed.L0.FullName = AberrantDeriveFullName(t) // e.g., github_com.user.repo.MyEnum
	ed.L0.ParentFile = filedesc.SurrogateProto3
	ed.L2.Values.List = append(ed.L2.Values.List, filedesc.EnumValue***REMOVED******REMOVED***)

	// TODO: Use the presence of a UnmarshalJSON method to determine proto2?

	vd := &ed.L2.Values.List[0]
	vd.L0.FullName = ed.L0.FullName + "_UNKNOWN" // e.g., github_com.user.repo.MyEnum_UNKNOWN
	vd.L0.ParentFile = ed.L0.ParentFile
	vd.L0.Parent = ed

	// TODO: We could use the String method to obtain some enum value names by
	// starting at 0 and print the enum until it produces invalid identifiers.
	// An exhaustive query is clearly impractical, but can be best-effort.

	if ed, ok := aberrantEnumDescCache.LoadOrStore(t, ed); ok ***REMOVED***
		return ed.(pref.EnumDescriptor)
	***REMOVED***
	return ed
***REMOVED***

// AberrantDeriveFullName derives a fully qualified protobuf name for the given Go type
// The provided name is not guaranteed to be stable nor universally unique.
// It should be sufficiently unique within a program.
//
// This is exported for testing purposes.
func AberrantDeriveFullName(t reflect.Type) pref.FullName ***REMOVED***
	sanitize := func(r rune) rune ***REMOVED***
		switch ***REMOVED***
		case r == '/':
			return '.'
		case 'a' <= r && r <= 'z', 'A' <= r && r <= 'Z', '0' <= r && r <= '9':
			return r
		default:
			return '_'
		***REMOVED***
	***REMOVED***
	prefix := strings.Map(sanitize, t.PkgPath())
	suffix := strings.Map(sanitize, t.Name())
	if suffix == "" ***REMOVED***
		suffix = fmt.Sprintf("UnknownX%X", reflect.ValueOf(t).Pointer())
	***REMOVED***

	ss := append(strings.Split(prefix, "."), suffix)
	for i, s := range ss ***REMOVED***
		if s == "" || ('0' <= s[0] && s[0] <= '9') ***REMOVED***
			ss[i] = "x" + s
		***REMOVED***
	***REMOVED***
	return pref.FullName(strings.Join(ss, "."))
***REMOVED***
