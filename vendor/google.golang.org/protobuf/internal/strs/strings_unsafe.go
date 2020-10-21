// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !purego,!appengine

package strs

import (
	"unsafe"

	pref "google.golang.org/protobuf/reflect/protoreflect"
)

type (
	stringHeader struct ***REMOVED***
		Data unsafe.Pointer
		Len  int
	***REMOVED***
	sliceHeader struct ***REMOVED***
		Data unsafe.Pointer
		Len  int
		Cap  int
	***REMOVED***
)

// UnsafeString returns an unsafe string reference of b.
// The caller must treat the input slice as immutable.
//
// WARNING: Use carefully. The returned result must not leak to the end user
// unless the input slice is provably immutable.
func UnsafeString(b []byte) (s string) ***REMOVED***
	src := (*sliceHeader)(unsafe.Pointer(&b))
	dst := (*stringHeader)(unsafe.Pointer(&s))
	dst.Data = src.Data
	dst.Len = src.Len
	return s
***REMOVED***

// UnsafeBytes returns an unsafe bytes slice reference of s.
// The caller must treat returned slice as immutable.
//
// WARNING: Use carefully. The returned result must not leak to the end user.
func UnsafeBytes(s string) (b []byte) ***REMOVED***
	src := (*stringHeader)(unsafe.Pointer(&s))
	dst := (*sliceHeader)(unsafe.Pointer(&b))
	dst.Data = src.Data
	dst.Len = src.Len
	dst.Cap = src.Len
	return b
***REMOVED***

// Builder builds a set of strings with shared lifetime.
// This differs from strings.Builder, which is for building a single string.
type Builder struct ***REMOVED***
	buf []byte
***REMOVED***

// AppendFullName is equivalent to protoreflect.FullName.Append,
// but optimized for large batches where each name has a shared lifetime.
func (sb *Builder) AppendFullName(prefix pref.FullName, name pref.Name) pref.FullName ***REMOVED***
	n := len(prefix) + len(".") + len(name)
	if len(prefix) == 0 ***REMOVED***
		n -= len(".")
	***REMOVED***
	sb.grow(n)
	sb.buf = append(sb.buf, prefix...)
	sb.buf = append(sb.buf, '.')
	sb.buf = append(sb.buf, name...)
	return pref.FullName(sb.last(n))
***REMOVED***

// MakeString is equivalent to string(b), but optimized for large batches
// with a shared lifetime.
func (sb *Builder) MakeString(b []byte) string ***REMOVED***
	sb.grow(len(b))
	sb.buf = append(sb.buf, b...)
	return sb.last(len(b))
***REMOVED***

func (sb *Builder) grow(n int) ***REMOVED***
	if cap(sb.buf)-len(sb.buf) >= n ***REMOVED***
		return
	***REMOVED***

	// Unlike strings.Builder, we do not need to copy over the contents
	// of the old buffer since our builder provides no API for
	// retrieving previously created strings.
	sb.buf = make([]byte, 2*(cap(sb.buf)+n))
***REMOVED***

func (sb *Builder) last(n int) string ***REMOVED***
	return UnsafeString(sb.buf[len(sb.buf)-n:])
***REMOVED***
