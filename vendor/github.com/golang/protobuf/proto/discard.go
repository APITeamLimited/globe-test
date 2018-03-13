// Go support for Protocol Buffers - Google's data interchange format
//
// Copyright 2017 The Go Authors.  All rights reserved.
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

package proto

import (
	"fmt"
	"reflect"
	"strings"
)

// DiscardUnknown recursively discards all unknown fields from this message
// and all embedded messages.
//
// When unmarshaling a message with unrecognized fields, the tags and values
// of such fields are preserved in the Message. This allows a later call to
// marshal to be able to produce a message that continues to have those
// unrecognized fields. To avoid this, DiscardUnknown is used to
// explicitly clear the unknown fields after unmarshaling.
//
// For proto2 messages, the unknown fields of message extensions are only
// discarded from messages that have been accessed via GetExtension.
func DiscardUnknown(m Message) ***REMOVED***
	discardLegacy(m)
***REMOVED***

func discardLegacy(m Message) ***REMOVED***
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Ptr || v.IsNil() ***REMOVED***
		return
	***REMOVED***
	v = v.Elem()
	if v.Kind() != reflect.Struct ***REMOVED***
		return
	***REMOVED***
	t := v.Type()

	for i := 0; i < v.NumField(); i++ ***REMOVED***
		f := t.Field(i)
		if strings.HasPrefix(f.Name, "XXX_") ***REMOVED***
			continue
		***REMOVED***
		vf := v.Field(i)
		tf := f.Type

		// Unwrap tf to get its most basic type.
		var isPointer, isSlice bool
		if tf.Kind() == reflect.Slice && tf.Elem().Kind() != reflect.Uint8 ***REMOVED***
			isSlice = true
			tf = tf.Elem()
		***REMOVED***
		if tf.Kind() == reflect.Ptr ***REMOVED***
			isPointer = true
			tf = tf.Elem()
		***REMOVED***
		if isPointer && isSlice && tf.Kind() != reflect.Struct ***REMOVED***
			panic(fmt.Sprintf("%T.%s cannot be a slice of pointers to primitive types", m, f.Name))
		***REMOVED***

		switch tf.Kind() ***REMOVED***
		case reflect.Struct:
			switch ***REMOVED***
			case !isPointer:
				panic(fmt.Sprintf("%T.%s cannot be a direct struct value", m, f.Name))
			case isSlice: // E.g., []*pb.T
				for j := 0; j < vf.Len(); j++ ***REMOVED***
					discardLegacy(vf.Index(j).Interface().(Message))
				***REMOVED***
			default: // E.g., *pb.T
				discardLegacy(vf.Interface().(Message))
			***REMOVED***
		case reflect.Map:
			switch ***REMOVED***
			case isPointer || isSlice:
				panic(fmt.Sprintf("%T.%s cannot be a pointer to a map or a slice of map values", m, f.Name))
			default: // E.g., map[K]V
				tv := vf.Type().Elem()
				if tv.Kind() == reflect.Ptr && tv.Implements(protoMessageType) ***REMOVED*** // Proto struct (e.g., *T)
					for _, key := range vf.MapKeys() ***REMOVED***
						val := vf.MapIndex(key)
						discardLegacy(val.Interface().(Message))
					***REMOVED***
				***REMOVED***
			***REMOVED***
		case reflect.Interface:
			// Must be oneof field.
			switch ***REMOVED***
			case isPointer || isSlice:
				panic(fmt.Sprintf("%T.%s cannot be a pointer to a interface or a slice of interface values", m, f.Name))
			default: // E.g., test_proto.isCommunique_Union interface
				if !vf.IsNil() && f.Tag.Get("protobuf_oneof") != "" ***REMOVED***
					vf = vf.Elem() // E.g., *test_proto.Communique_Msg
					if !vf.IsNil() ***REMOVED***
						vf = vf.Elem()   // E.g., test_proto.Communique_Msg
						vf = vf.Field(0) // E.g., Proto struct (e.g., *T) or primitive value
						if vf.Kind() == reflect.Ptr ***REMOVED***
							discardLegacy(vf.Interface().(Message))
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if vf := v.FieldByName("XXX_unrecognized"); vf.IsValid() ***REMOVED***
		if vf.Type() != reflect.TypeOf([]byte***REMOVED******REMOVED***) ***REMOVED***
			panic("expected XXX_unrecognized to be of type []byte")
		***REMOVED***
		vf.Set(reflect.ValueOf([]byte(nil)))
	***REMOVED***

	// For proto2 messages, only discard unknown fields in message extensions
	// that have been accessed via GetExtension.
	if em, ok := extendable(m); ok ***REMOVED***
		// Ignore lock since discardLegacy is not concurrency safe.
		emm, _ := em.extensionsRead()
		for _, mx := range emm ***REMOVED***
			if m, ok := mx.value.(Message); ok ***REMOVED***
				discardLegacy(m)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
