// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ptypes

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	anypb "github.com/golang/protobuf/ptypes/any"
)

const urlPrefix = "type.googleapis.com/"

// AnyMessageName returns the message name contained in an anypb.Any message.
// Most type assertions should use the Is function instead.
func AnyMessageName(any *anypb.Any) (string, error) ***REMOVED***
	name, err := anyMessageName(any)
	return string(name), err
***REMOVED***
func anyMessageName(any *anypb.Any) (protoreflect.FullName, error) ***REMOVED***
	if any == nil ***REMOVED***
		return "", fmt.Errorf("message is nil")
	***REMOVED***
	name := protoreflect.FullName(any.TypeUrl)
	if i := strings.LastIndex(any.TypeUrl, "/"); i >= 0 ***REMOVED***
		name = name[i+len("/"):]
	***REMOVED***
	if !name.IsValid() ***REMOVED***
		return "", fmt.Errorf("message type url %q is invalid", any.TypeUrl)
	***REMOVED***
	return name, nil
***REMOVED***

// MarshalAny marshals the given message m into an anypb.Any message.
func MarshalAny(m proto.Message) (*anypb.Any, error) ***REMOVED***
	switch dm := m.(type) ***REMOVED***
	case DynamicAny:
		m = dm.Message
	case *DynamicAny:
		if dm == nil ***REMOVED***
			return nil, proto.ErrNil
		***REMOVED***
		m = dm.Message
	***REMOVED***
	b, err := proto.Marshal(m)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &anypb.Any***REMOVED***TypeUrl: urlPrefix + proto.MessageName(m), Value: b***REMOVED***, nil
***REMOVED***

// Empty returns a new message of the type specified in an anypb.Any message.
// It returns protoregistry.NotFound if the corresponding message type could not
// be resolved in the global registry.
func Empty(any *anypb.Any) (proto.Message, error) ***REMOVED***
	name, err := anyMessageName(any)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	mt, err := protoregistry.GlobalTypes.FindMessageByName(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return proto.MessageV1(mt.New().Interface()), nil
***REMOVED***

// UnmarshalAny unmarshals the encoded value contained in the anypb.Any message
// into the provided message m. It returns an error if the target message
// does not match the type in the Any message or if an unmarshal error occurs.
//
// The target message m may be a *DynamicAny message. If the underlying message
// type could not be resolved, then this returns protoregistry.NotFound.
func UnmarshalAny(any *anypb.Any, m proto.Message) error ***REMOVED***
	if dm, ok := m.(*DynamicAny); ok ***REMOVED***
		if dm.Message == nil ***REMOVED***
			var err error
			dm.Message, err = Empty(any)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		m = dm.Message
	***REMOVED***

	anyName, err := AnyMessageName(any)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	msgName := proto.MessageName(m)
	if anyName != msgName ***REMOVED***
		return fmt.Errorf("mismatched message type: got %q want %q", anyName, msgName)
	***REMOVED***
	return proto.Unmarshal(any.Value, m)
***REMOVED***

// Is reports whether the Any message contains a message of the specified type.
func Is(any *anypb.Any, m proto.Message) bool ***REMOVED***
	if any == nil || m == nil ***REMOVED***
		return false
	***REMOVED***
	name := proto.MessageName(m)
	if !strings.HasSuffix(any.TypeUrl, name) ***REMOVED***
		return false
	***REMOVED***
	return len(any.TypeUrl) == len(name) || any.TypeUrl[len(any.TypeUrl)-len(name)-1] == '/'
***REMOVED***

// DynamicAny is a value that can be passed to UnmarshalAny to automatically
// allocate a proto.Message for the type specified in an anypb.Any message.
// The allocated message is stored in the embedded proto.Message.
//
// Example:
//   var x ptypes.DynamicAny
//   if err := ptypes.UnmarshalAny(a, &x); err != nil ***REMOVED*** ... ***REMOVED***
//   fmt.Printf("unmarshaled message: %v", x.Message)
type DynamicAny struct***REMOVED*** proto.Message ***REMOVED***

func (m DynamicAny) String() string ***REMOVED***
	if m.Message == nil ***REMOVED***
		return "<nil>"
	***REMOVED***
	return m.Message.String()
***REMOVED***
func (m DynamicAny) Reset() ***REMOVED***
	if m.Message == nil ***REMOVED***
		return
	***REMOVED***
	m.Message.Reset()
***REMOVED***
func (m DynamicAny) ProtoMessage() ***REMOVED***
	return
***REMOVED***
func (m DynamicAny) ProtoReflect() protoreflect.Message ***REMOVED***
	if m.Message == nil ***REMOVED***
		return nil
	***REMOVED***
	return dynamicAny***REMOVED***proto.MessageReflect(m.Message)***REMOVED***
***REMOVED***

type dynamicAny struct***REMOVED*** protoreflect.Message ***REMOVED***

func (m dynamicAny) Type() protoreflect.MessageType ***REMOVED***
	return dynamicAnyType***REMOVED***m.Message.Type()***REMOVED***
***REMOVED***
func (m dynamicAny) New() protoreflect.Message ***REMOVED***
	return dynamicAnyType***REMOVED***m.Message.Type()***REMOVED***.New()
***REMOVED***
func (m dynamicAny) Interface() protoreflect.ProtoMessage ***REMOVED***
	return DynamicAny***REMOVED***proto.MessageV1(m.Message.Interface())***REMOVED***
***REMOVED***

type dynamicAnyType struct***REMOVED*** protoreflect.MessageType ***REMOVED***

func (t dynamicAnyType) New() protoreflect.Message ***REMOVED***
	return dynamicAny***REMOVED***t.MessageType.New()***REMOVED***
***REMOVED***
func (t dynamicAnyType) Zero() protoreflect.Message ***REMOVED***
	return dynamicAny***REMOVED***t.MessageType.Zero()***REMOVED***
***REMOVED***
