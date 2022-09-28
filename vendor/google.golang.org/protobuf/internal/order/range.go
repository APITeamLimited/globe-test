// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package order provides ordered access to messages and maps.
package order

import (
	"sort"
	"sync"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type messageField struct ***REMOVED***
	fd protoreflect.FieldDescriptor
	v  protoreflect.Value
***REMOVED***

var messageFieldPool = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return new([]messageField) ***REMOVED***,
***REMOVED***

type (
	// FieldRnger is an interface for visiting all fields in a message.
	// The protoreflect.Message type implements this interface.
	FieldRanger interface***REMOVED*** Range(VisitField) ***REMOVED***
	// VisitField is called every time a message field is visited.
	VisitField = func(protoreflect.FieldDescriptor, protoreflect.Value) bool
)

// RangeFields iterates over the fields of fs according to the specified order.
func RangeFields(fs FieldRanger, less FieldOrder, fn VisitField) ***REMOVED***
	if less == nil ***REMOVED***
		fs.Range(fn)
		return
	***REMOVED***

	// Obtain a pre-allocated scratch buffer.
	p := messageFieldPool.Get().(*[]messageField)
	fields := (*p)[:0]
	defer func() ***REMOVED***
		if cap(fields) < 1024 ***REMOVED***
			*p = fields
			messageFieldPool.Put(p)
		***REMOVED***
	***REMOVED***()

	// Collect all fields in the message and sort them.
	fs.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool ***REMOVED***
		fields = append(fields, messageField***REMOVED***fd, v***REMOVED***)
		return true
	***REMOVED***)
	sort.Slice(fields, func(i, j int) bool ***REMOVED***
		return less(fields[i].fd, fields[j].fd)
	***REMOVED***)

	// Visit the fields in the specified ordering.
	for _, f := range fields ***REMOVED***
		if !fn(f.fd, f.v) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

type mapEntry struct ***REMOVED***
	k protoreflect.MapKey
	v protoreflect.Value
***REMOVED***

var mapEntryPool = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return new([]mapEntry) ***REMOVED***,
***REMOVED***

type (
	// EntryRanger is an interface for visiting all fields in a message.
	// The protoreflect.Map type implements this interface.
	EntryRanger interface***REMOVED*** Range(VisitEntry) ***REMOVED***
	// VisitEntry is called every time a map entry is visited.
	VisitEntry = func(protoreflect.MapKey, protoreflect.Value) bool
)

// RangeEntries iterates over the entries of es according to the specified order.
func RangeEntries(es EntryRanger, less KeyOrder, fn VisitEntry) ***REMOVED***
	if less == nil ***REMOVED***
		es.Range(fn)
		return
	***REMOVED***

	// Obtain a pre-allocated scratch buffer.
	p := mapEntryPool.Get().(*[]mapEntry)
	entries := (*p)[:0]
	defer func() ***REMOVED***
		if cap(entries) < 1024 ***REMOVED***
			*p = entries
			mapEntryPool.Put(p)
		***REMOVED***
	***REMOVED***()

	// Collect all entries in the map and sort them.
	es.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool ***REMOVED***
		entries = append(entries, mapEntry***REMOVED***k, v***REMOVED***)
		return true
	***REMOVED***)
	sort.Slice(entries, func(i, j int) bool ***REMOVED***
		return less(entries[i].k, entries[j].k)
	***REMOVED***)

	// Visit the entries in the specified ordering.
	for _, e := range entries ***REMOVED***
		if !fn(e.k, e.v) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***
