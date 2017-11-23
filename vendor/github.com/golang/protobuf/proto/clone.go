// Go support for Protocol Buffers - Google's data interchange format
//
// Copyright 2011 The Go Authors.  All rights reserved.
// https://github.com/golang/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Protocol buffer deep copy and merge.
// TODO: RawMessage.

package proto

import (
	"log"
	"reflect"
	"strings"
)

// Clone returns a deep copy of a protocol buffer.
func Clone(pb Message) Message ***REMOVED***
	in := reflect.ValueOf(pb)
	if in.IsNil() ***REMOVED***
		return pb
	***REMOVED***

	out := reflect.New(in.Type().Elem())
	// out is empty so a merge is a deep copy.
	mergeStruct(out.Elem(), in.Elem())
	return out.Interface().(Message)
***REMOVED***

// Merge merges src into dst.
// Required and optional fields that are set in src will be set to that value in dst.
// Elements of repeated fields will be appended.
// Merge panics if src and dst are not the same type, or if dst is nil.
func Merge(dst, src Message) ***REMOVED***
	in := reflect.ValueOf(src)
	out := reflect.ValueOf(dst)
	if out.IsNil() ***REMOVED***
		panic("proto: nil destination")
	***REMOVED***
	if in.Type() != out.Type() ***REMOVED***
		// Explicit test prior to mergeStruct so that mistyped nils will fail
		panic("proto: type mismatch")
	***REMOVED***
	if in.IsNil() ***REMOVED***
		// Merging nil into non-nil is a quiet no-op
		return
	***REMOVED***
	mergeStruct(out.Elem(), in.Elem())
***REMOVED***

func mergeStruct(out, in reflect.Value) ***REMOVED***
	sprop := GetProperties(in.Type())
	for i := 0; i < in.NumField(); i++ ***REMOVED***
		f := in.Type().Field(i)
		if strings.HasPrefix(f.Name, "XXX_") ***REMOVED***
			continue
		***REMOVED***
		mergeAny(out.Field(i), in.Field(i), false, sprop.Prop[i])
	***REMOVED***

	if emIn, ok := extendable(in.Addr().Interface()); ok ***REMOVED***
		emOut, _ := extendable(out.Addr().Interface())
		mIn, muIn := emIn.extensionsRead()
		if mIn != nil ***REMOVED***
			mOut := emOut.extensionsWrite()
			muIn.Lock()
			mergeExtension(mOut, mIn)
			muIn.Unlock()
		***REMOVED***
	***REMOVED***

	uf := in.FieldByName("XXX_unrecognized")
	if !uf.IsValid() ***REMOVED***
		return
	***REMOVED***
	uin := uf.Bytes()
	if len(uin) > 0 ***REMOVED***
		out.FieldByName("XXX_unrecognized").SetBytes(append([]byte(nil), uin...))
	***REMOVED***
***REMOVED***

// mergeAny performs a merge between two values of the same type.
// viaPtr indicates whether the values were indirected through a pointer (implying proto2).
// prop is set if this is a struct field (it may be nil).
func mergeAny(out, in reflect.Value, viaPtr bool, prop *Properties) ***REMOVED***
	if in.Type() == protoMessageType ***REMOVED***
		if !in.IsNil() ***REMOVED***
			if out.IsNil() ***REMOVED***
				out.Set(reflect.ValueOf(Clone(in.Interface().(Message))))
			***REMOVED*** else ***REMOVED***
				Merge(out.Interface().(Message), in.Interface().(Message))
			***REMOVED***
		***REMOVED***
		return
	***REMOVED***
	switch in.Kind() ***REMOVED***
	case reflect.Bool, reflect.Float32, reflect.Float64, reflect.Int32, reflect.Int64,
		reflect.String, reflect.Uint32, reflect.Uint64:
		if !viaPtr && isProto3Zero(in) ***REMOVED***
			return
		***REMOVED***
		out.Set(in)
	case reflect.Interface:
		// Probably a oneof field; copy non-nil values.
		if in.IsNil() ***REMOVED***
			return
		***REMOVED***
		// Allocate destination if it is not set, or set to a different type.
		// Otherwise we will merge as normal.
		if out.IsNil() || out.Elem().Type() != in.Elem().Type() ***REMOVED***
			out.Set(reflect.New(in.Elem().Elem().Type())) // interface -> *T -> T -> new(T)
		***REMOVED***
		mergeAny(out.Elem(), in.Elem(), false, nil)
	case reflect.Map:
		if in.Len() == 0 ***REMOVED***
			return
		***REMOVED***
		if out.IsNil() ***REMOVED***
			out.Set(reflect.MakeMap(in.Type()))
		***REMOVED***
		// For maps with value types of *T or []byte we need to deep copy each value.
		elemKind := in.Type().Elem().Kind()
		for _, key := range in.MapKeys() ***REMOVED***
			var val reflect.Value
			switch elemKind ***REMOVED***
			case reflect.Ptr:
				val = reflect.New(in.Type().Elem().Elem())
				mergeAny(val, in.MapIndex(key), false, nil)
			case reflect.Slice:
				val = in.MapIndex(key)
				val = reflect.ValueOf(append([]byte***REMOVED******REMOVED***, val.Bytes()...))
			default:
				val = in.MapIndex(key)
			***REMOVED***
			out.SetMapIndex(key, val)
		***REMOVED***
	case reflect.Ptr:
		if in.IsNil() ***REMOVED***
			return
		***REMOVED***
		if out.IsNil() ***REMOVED***
			out.Set(reflect.New(in.Elem().Type()))
		***REMOVED***
		mergeAny(out.Elem(), in.Elem(), true, nil)
	case reflect.Slice:
		if in.IsNil() ***REMOVED***
			return
		***REMOVED***
		if in.Type().Elem().Kind() == reflect.Uint8 ***REMOVED***
			// []byte is a scalar bytes field, not a repeated field.

			// Edge case: if this is in a proto3 message, a zero length
			// bytes field is considered the zero value, and should not
			// be merged.
			if prop != nil && prop.proto3 && in.Len() == 0 ***REMOVED***
				return
			***REMOVED***

			// Make a deep copy.
			// Append to []byte***REMOVED******REMOVED*** instead of []byte(nil) so that we never end up
			// with a nil result.
			out.SetBytes(append([]byte***REMOVED******REMOVED***, in.Bytes()...))
			return
		***REMOVED***
		n := in.Len()
		if out.IsNil() ***REMOVED***
			out.Set(reflect.MakeSlice(in.Type(), 0, n))
		***REMOVED***
		switch in.Type().Elem().Kind() ***REMOVED***
		case reflect.Bool, reflect.Float32, reflect.Float64, reflect.Int32, reflect.Int64,
			reflect.String, reflect.Uint32, reflect.Uint64:
			out.Set(reflect.AppendSlice(out, in))
		default:
			for i := 0; i < n; i++ ***REMOVED***
				x := reflect.Indirect(reflect.New(in.Type().Elem()))
				mergeAny(x, in.Index(i), false, nil)
				out.Set(reflect.Append(out, x))
			***REMOVED***
		***REMOVED***
	case reflect.Struct:
		mergeStruct(out, in)
	default:
		// unknown type, so not a protocol buffer
		log.Printf("proto: don't know how to copy %v", in)
	***REMOVED***
***REMOVED***

func mergeExtension(out, in map[int32]Extension) ***REMOVED***
	for extNum, eIn := range in ***REMOVED***
		eOut := Extension***REMOVED***desc: eIn.desc***REMOVED***
		if eIn.value != nil ***REMOVED***
			v := reflect.New(reflect.TypeOf(eIn.value)).Elem()
			mergeAny(v, reflect.ValueOf(eIn.value), false, nil)
			eOut.value = v.Interface()
		***REMOVED***
		if eIn.enc != nil ***REMOVED***
			eOut.enc = make([]byte, len(eIn.enc))
			copy(eOut.enc, eIn.enc)
		***REMOVED***

		out[extNum] = eOut
	***REMOVED***
***REMOVED***
