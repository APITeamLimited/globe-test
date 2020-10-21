package dynamic

import (
	"bytes"
	"reflect"

	"github.com/golang/protobuf/proto"

	"github.com/jhump/protoreflect/desc"
)

// Equal returns true if the given two dynamic messages are equal. Two messages are equal when they
// have the same message type and same fields set to equal values. For proto3 messages, fields set
// to their zero value are considered unset.
func Equal(a, b *Message) bool ***REMOVED***
	if a == b ***REMOVED***
		return true
	***REMOVED***
	if (a == nil) != (b == nil) ***REMOVED***
		return false
	***REMOVED***
	if a.md.GetFullyQualifiedName() != b.md.GetFullyQualifiedName() ***REMOVED***
		return false
	***REMOVED***
	if len(a.values) != len(b.values) ***REMOVED***
		return false
	***REMOVED***
	if len(a.unknownFields) != len(b.unknownFields) ***REMOVED***
		return false
	***REMOVED***
	for tag, aval := range a.values ***REMOVED***
		bval, ok := b.values[tag]
		if !ok ***REMOVED***
			return false
		***REMOVED***
		if !fieldsEqual(aval, bval) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	for tag, au := range a.unknownFields ***REMOVED***
		bu, ok := b.unknownFields[tag]
		if !ok ***REMOVED***
			return false
		***REMOVED***
		if len(au) != len(bu) ***REMOVED***
			return false
		***REMOVED***
		for i, aval := range au ***REMOVED***
			bval := bu[i]
			if aval.Encoding != bval.Encoding ***REMOVED***
				return false
			***REMOVED***
			if aval.Encoding == proto.WireBytes || aval.Encoding == proto.WireStartGroup ***REMOVED***
				if !bytes.Equal(aval.Contents, bval.Contents) ***REMOVED***
					return false
				***REMOVED***
			***REMOVED*** else if aval.Value != bval.Value ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// all checks pass!
	return true
***REMOVED***

func fieldsEqual(aval, bval interface***REMOVED******REMOVED***) bool ***REMOVED***
	arv := reflect.ValueOf(aval)
	brv := reflect.ValueOf(bval)
	if arv.Type() != brv.Type() ***REMOVED***
		// it is possible that one is a dynamic message and one is not
		apm, ok := aval.(proto.Message)
		if !ok ***REMOVED***
			return false
		***REMOVED***
		bpm, ok := bval.(proto.Message)
		if !ok ***REMOVED***
			return false
		***REMOVED***
		return MessagesEqual(apm, bpm)

	***REMOVED*** else ***REMOVED***
		switch arv.Kind() ***REMOVED***
		case reflect.Ptr:
			apm, ok := aval.(proto.Message)
			if !ok ***REMOVED***
				// Don't know how to compare pointer values that aren't messages!
				// Maybe this should panic?
				return false
			***REMOVED***
			bpm := bval.(proto.Message) // we know it will succeed because we know a and b have same type
			return MessagesEqual(apm, bpm)

		case reflect.Map:
			return mapsEqual(arv, brv)

		case reflect.Slice:
			if arv.Type() == typeOfBytes ***REMOVED***
				return bytes.Equal(aval.([]byte), bval.([]byte))
			***REMOVED*** else ***REMOVED***
				return slicesEqual(arv, brv)
			***REMOVED***

		default:
			return aval == bval
		***REMOVED***
	***REMOVED***
***REMOVED***

func slicesEqual(a, b reflect.Value) bool ***REMOVED***
	if a.Len() != b.Len() ***REMOVED***
		return false
	***REMOVED***
	for i := 0; i < a.Len(); i++ ***REMOVED***
		ai := a.Index(i)
		bi := b.Index(i)
		if !fieldsEqual(ai.Interface(), bi.Interface()) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// MessagesEqual returns true if the given two messages are equal. Use this instead of proto.Equal
// when one or both of the messages might be a dynamic message.
func MessagesEqual(a, b proto.Message) bool ***REMOVED***
	da, aok := a.(*Message)
	db, bok := b.(*Message)
	// Both dynamic messages
	if aok && bok ***REMOVED***
		return Equal(da, db)
	***REMOVED***
	// Neither dynamic messages
	if !aok && !bok ***REMOVED***
		return proto.Equal(a, b)
	***REMOVED***
	// Mixed
	if bok ***REMOVED***
		// we want a to be the dynamic one
		b, da = a, db
	***REMOVED***

	// Instead of panic'ing below if we have a nil dynamic message, check
	// now and return false if the input message is not also nil.
	if da == nil ***REMOVED***
		return isNil(b)
	***REMOVED***

	md, err := desc.LoadMessageDescriptorForMessage(b)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	db = NewMessageWithMessageFactory(md, da.mf)
	if db.ConvertFrom(b) != nil ***REMOVED***
		return false
	***REMOVED***
	return Equal(da, db)
***REMOVED***
