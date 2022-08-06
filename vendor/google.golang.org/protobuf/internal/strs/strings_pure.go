// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build purego || appengine
// +build purego appengine

package strs

import pref "google.golang.org/protobuf/reflect/protoreflect"

func UnsafeString(b []byte) string ***REMOVED***
	return string(b)
***REMOVED***

func UnsafeBytes(s string) []byte ***REMOVED***
	return []byte(s)
***REMOVED***

type Builder struct***REMOVED******REMOVED***

func (*Builder) AppendFullName(prefix pref.FullName, name pref.Name) pref.FullName ***REMOVED***
	return prefix.Append(name)
***REMOVED***

func (*Builder) MakeString(b []byte) string ***REMOVED***
	return string(b)
***REMOVED***
