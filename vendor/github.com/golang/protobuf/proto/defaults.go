// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

// SetDefaults sets unpopulated scalar fields to their default values.
// Fields within a oneof are not set even if they have a default value.
// SetDefaults is recursively called upon any populated message fields.
func SetDefaults(m Message) ***REMOVED***
	if m != nil ***REMOVED***
		setDefaults(MessageReflect(m))
	***REMOVED***
***REMOVED***

func setDefaults(m protoreflect.Message) ***REMOVED***
	fds := m.Descriptor().Fields()
	for i := 0; i < fds.Len(); i++ ***REMOVED***
		fd := fds.Get(i)
		if !m.Has(fd) ***REMOVED***
			if fd.HasDefault() && fd.ContainingOneof() == nil ***REMOVED***
				v := fd.Default()
				if fd.Kind() == protoreflect.BytesKind ***REMOVED***
					v = protoreflect.ValueOf(append([]byte(nil), v.Bytes()...)) // copy the default bytes
				***REMOVED***
				m.Set(fd, v)
			***REMOVED***
			continue
		***REMOVED***
	***REMOVED***

	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool ***REMOVED***
		switch ***REMOVED***
		// Handle singular message.
		case fd.Cardinality() != protoreflect.Repeated:
			if fd.Message() != nil ***REMOVED***
				setDefaults(m.Get(fd).Message())
			***REMOVED***
		// Handle list of messages.
		case fd.IsList():
			if fd.Message() != nil ***REMOVED***
				ls := m.Get(fd).List()
				for i := 0; i < ls.Len(); i++ ***REMOVED***
					setDefaults(ls.Get(i).Message())
				***REMOVED***
			***REMOVED***
		// Handle map of messages.
		case fd.IsMap():
			if fd.MapValue().Message() != nil ***REMOVED***
				ms := m.Get(fd).Map()
				ms.Range(func(_ protoreflect.MapKey, v protoreflect.Value) bool ***REMOVED***
					setDefaults(v.Message())
					return true
				***REMOVED***)
			***REMOVED***
		***REMOVED***
		return true
	***REMOVED***)
***REMOVED***
