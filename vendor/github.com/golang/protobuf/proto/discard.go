// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

// DiscardUnknown recursively discards all unknown fields from this message
// and all embedded messages.
//
// When unmarshaling a message with unrecognized fields, the tags and values
// of such fields are preserved in the Message. This allows a later call to
// marshal to be able to produce a message that continues to have those
// unrecognized fields. To avoid this, DiscardUnknown is used to
// explicitly clear the unknown fields after unmarshaling.
func DiscardUnknown(m Message) ***REMOVED***
	if m != nil ***REMOVED***
		discardUnknown(MessageReflect(m))
	***REMOVED***
***REMOVED***

func discardUnknown(m protoreflect.Message) ***REMOVED***
	m.Range(func(fd protoreflect.FieldDescriptor, val protoreflect.Value) bool ***REMOVED***
		switch ***REMOVED***
		// Handle singular message.
		case fd.Cardinality() != protoreflect.Repeated:
			if fd.Message() != nil ***REMOVED***
				discardUnknown(m.Get(fd).Message())
			***REMOVED***
		// Handle list of messages.
		case fd.IsList():
			if fd.Message() != nil ***REMOVED***
				ls := m.Get(fd).List()
				for i := 0; i < ls.Len(); i++ ***REMOVED***
					discardUnknown(ls.Get(i).Message())
				***REMOVED***
			***REMOVED***
		// Handle map of messages.
		case fd.IsMap():
			if fd.MapValue().Message() != nil ***REMOVED***
				ms := m.Get(fd).Map()
				ms.Range(func(_ protoreflect.MapKey, v protoreflect.Value) bool ***REMOVED***
					discardUnknown(v.Message())
					return true
				***REMOVED***)
			***REMOVED***
		***REMOVED***
		return true
	***REMOVED***)

	// Discard unknown fields.
	if len(m.GetUnknown()) > 0 ***REMOVED***
		m.SetUnknown(nil)
	***REMOVED***
***REMOVED***
