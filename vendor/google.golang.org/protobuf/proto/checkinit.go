// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
)

// CheckInitialized returns an error if any required fields in m are not set.
func CheckInitialized(m Message) error ***REMOVED***
	// Treat a nil message interface as an "untyped" empty message,
	// which we assume to have no required fields.
	if m == nil ***REMOVED***
		return nil
	***REMOVED***

	return checkInitialized(m.ProtoReflect())
***REMOVED***

// CheckInitialized returns an error if any required fields in m are not set.
func checkInitialized(m protoreflect.Message) error ***REMOVED***
	if methods := protoMethods(m); methods != nil && methods.CheckInitialized != nil ***REMOVED***
		_, err := methods.CheckInitialized(protoiface.CheckInitializedInput***REMOVED***
			Message: m,
		***REMOVED***)
		return err
	***REMOVED***
	return checkInitializedSlow(m)
***REMOVED***

func checkInitializedSlow(m protoreflect.Message) error ***REMOVED***
	md := m.Descriptor()
	fds := md.Fields()
	for i, nums := 0, md.RequiredNumbers(); i < nums.Len(); i++ ***REMOVED***
		fd := fds.ByNumber(nums.Get(i))
		if !m.Has(fd) ***REMOVED***
			return errors.RequiredNotSet(string(fd.FullName()))
		***REMOVED***
	***REMOVED***
	var err error
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool ***REMOVED***
		switch ***REMOVED***
		case fd.IsList():
			if fd.Message() == nil ***REMOVED***
				return true
			***REMOVED***
			for i, list := 0, v.List(); i < list.Len() && err == nil; i++ ***REMOVED***
				err = checkInitialized(list.Get(i).Message())
			***REMOVED***
		case fd.IsMap():
			if fd.MapValue().Message() == nil ***REMOVED***
				return true
			***REMOVED***
			v.Map().Range(func(key protoreflect.MapKey, v protoreflect.Value) bool ***REMOVED***
				err = checkInitialized(v.Message())
				return err == nil
			***REMOVED***)
		default:
			if fd.Message() == nil ***REMOVED***
				return true
			***REMOVED***
			err = checkInitialized(v.Message())
		***REMOVED***
		return err == nil
	***REMOVED***)
	return err
***REMOVED***
