package dynamic

import (
	"errors"
	"reflect"

	"github.com/golang/protobuf/proto"

	"github.com/jhump/protoreflect/desc"
)

// Merge merges the given source message into the given destination message. Use
// use this instead of proto.Merge when one or both of the messages might be a
// a dynamic message. If there is a problem merging the messages, such as the
// two messages having different types, then this method will panic (just as
// proto.Merges does).
func Merge(dst, src proto.Message) ***REMOVED***
	if dm, ok := dst.(*Message); ok ***REMOVED***
		if err := dm.MergeFrom(src); err != nil ***REMOVED***
			panic(err.Error())
		***REMOVED***
	***REMOVED*** else if dm, ok := src.(*Message); ok ***REMOVED***
		if err := dm.MergeInto(dst); err != nil ***REMOVED***
			panic(err.Error())
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		proto.Merge(dst, src)
	***REMOVED***
***REMOVED***

// TryMerge merges the given source message into the given destination message.
// You can use this instead of proto.Merge when one or both of the messages
// might be a dynamic message. Unlike proto.Merge, this method will return an
// error on failure instead of panic'ing.
func TryMerge(dst, src proto.Message) error ***REMOVED***
	if dm, ok := dst.(*Message); ok ***REMOVED***
		if err := dm.MergeFrom(src); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else if dm, ok := src.(*Message); ok ***REMOVED***
		if err := dm.MergeInto(dst); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// proto.Merge panics on bad input, so we first verify
		// inputs and return error instead of panic
		out := reflect.ValueOf(dst)
		if out.IsNil() ***REMOVED***
			return errors.New("proto: nil destination")
		***REMOVED***
		in := reflect.ValueOf(src)
		if in.Type() != out.Type() ***REMOVED***
			return errors.New("proto: type mismatch")
		***REMOVED***
		proto.Merge(dst, src)
	***REMOVED***
	return nil
***REMOVED***

func mergeField(m *Message, fd *desc.FieldDescriptor, val interface***REMOVED******REMOVED***) error ***REMOVED***
	rv := reflect.ValueOf(val)

	if fd.IsMap() && rv.Kind() == reflect.Map ***REMOVED***
		return mergeMapField(m, fd, rv)
	***REMOVED***

	if fd.IsRepeated() && rv.Kind() == reflect.Slice && rv.Type() != typeOfBytes ***REMOVED***
		for i := 0; i < rv.Len(); i++ ***REMOVED***
			e := rv.Index(i)
			if e.Kind() == reflect.Interface && !e.IsNil() ***REMOVED***
				e = e.Elem()
			***REMOVED***
			if err := m.addRepeatedField(fd, e.Interface()); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***

	if fd.IsRepeated() ***REMOVED***
		return m.addRepeatedField(fd, val)
	***REMOVED*** else if fd.GetMessageType() == nil ***REMOVED***
		return m.setField(fd, val)
	***REMOVED***

	// it's a message type, so we want to merge contents
	var err error
	if val, err = validFieldValue(fd, val); err != nil ***REMOVED***
		return err
	***REMOVED***

	existing, _ := m.doGetField(fd, true)
	if existing != nil && !reflect.ValueOf(existing).IsNil() ***REMOVED***
		return TryMerge(existing.(proto.Message), val.(proto.Message))
	***REMOVED***

	// no existing message, so just set field
	m.internalSetField(fd, val)
	return nil
***REMOVED***
