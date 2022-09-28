// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// Reset clears every field in the message.
// The resulting message shares no observable memory with its previous state
// other than the memory for the message itself.
func Reset(m Message) ***REMOVED***
	if mr, ok := m.(interface***REMOVED*** Reset() ***REMOVED***); ok && hasProtoMethods ***REMOVED***
		mr.Reset()
		return
	***REMOVED***
	resetMessage(m.ProtoReflect())
***REMOVED***

func resetMessage(m protoreflect.Message) ***REMOVED***
	if !m.IsValid() ***REMOVED***
		panic(fmt.Sprintf("cannot reset invalid %v message", m.Descriptor().FullName()))
	***REMOVED***

	// Clear all known fields.
	fds := m.Descriptor().Fields()
	for i := 0; i < fds.Len(); i++ ***REMOVED***
		m.Clear(fds.Get(i))
	***REMOVED***

	// Clear extension fields.
	m.Range(func(fd protoreflect.FieldDescriptor, _ protoreflect.Value) bool ***REMOVED***
		m.Clear(fd)
		return true
	***REMOVED***)

	// Clear unknown fields.
	m.SetUnknown(nil)
***REMOVED***
