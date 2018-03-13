// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

type ErrorType uint64

const (
	ErrorTypeBind    ErrorType = 1 << 63 // used when c.Bind() fails
	ErrorTypeRender  ErrorType = 1 << 62 // used when c.Render() fails
	ErrorTypePrivate ErrorType = 1 << 0
	ErrorTypePublic  ErrorType = 1 << 1

	ErrorTypeAny ErrorType = 1<<64 - 1
	ErrorTypeNu            = 2
)

type (
	Error struct ***REMOVED***
		Err  error
		Type ErrorType
		Meta interface***REMOVED******REMOVED***
	***REMOVED***

	errorMsgs []*Error
)

var _ error = &Error***REMOVED******REMOVED***

func (msg *Error) SetType(flags ErrorType) *Error ***REMOVED***
	msg.Type = flags
	return msg
***REMOVED***

func (msg *Error) SetMeta(data interface***REMOVED******REMOVED***) *Error ***REMOVED***
	msg.Meta = data
	return msg
***REMOVED***

func (msg *Error) JSON() interface***REMOVED******REMOVED*** ***REMOVED***
	json := H***REMOVED******REMOVED***
	if msg.Meta != nil ***REMOVED***
		value := reflect.ValueOf(msg.Meta)
		switch value.Kind() ***REMOVED***
		case reflect.Struct:
			return msg.Meta
		case reflect.Map:
			for _, key := range value.MapKeys() ***REMOVED***
				json[key.String()] = value.MapIndex(key).Interface()
			***REMOVED***
		default:
			json["meta"] = msg.Meta
		***REMOVED***
	***REMOVED***
	if _, ok := json["error"]; !ok ***REMOVED***
		json["error"] = msg.Error()
	***REMOVED***
	return json
***REMOVED***

// MarshalJSON implements the json.Marshaller interface
func (msg *Error) MarshalJSON() ([]byte, error) ***REMOVED***
	return json.Marshal(msg.JSON())
***REMOVED***

// Implements the error interface
func (msg Error) Error() string ***REMOVED***
	return msg.Err.Error()
***REMOVED***

func (msg *Error) IsType(flags ErrorType) bool ***REMOVED***
	return (msg.Type & flags) > 0
***REMOVED***

// Returns a readonly copy filtered the byte.
// ie ByType(gin.ErrorTypePublic) returns a slice of errors with type=ErrorTypePublic
func (a errorMsgs) ByType(typ ErrorType) errorMsgs ***REMOVED***
	if len(a) == 0 ***REMOVED***
		return nil
	***REMOVED***
	if typ == ErrorTypeAny ***REMOVED***
		return a
	***REMOVED***
	var result errorMsgs
	for _, msg := range a ***REMOVED***
		if msg.IsType(typ) ***REMOVED***
			result = append(result, msg)
		***REMOVED***
	***REMOVED***
	return result
***REMOVED***

// Returns the last error in the slice. It returns nil if the array is empty.
// Shortcut for errors[len(errors)-1]
func (a errorMsgs) Last() *Error ***REMOVED***
	length := len(a)
	if length > 0 ***REMOVED***
		return a[length-1]
	***REMOVED***
	return nil
***REMOVED***

// Returns an array will all the error messages.
// Example:
// 		c.Error(errors.New("first"))
// 		c.Error(errors.New("second"))
// 		c.Error(errors.New("third"))
// 		c.Errors.Errors() // == []string***REMOVED***"first", "second", "third"***REMOVED***
func (a errorMsgs) Errors() []string ***REMOVED***
	if len(a) == 0 ***REMOVED***
		return nil
	***REMOVED***
	errorStrings := make([]string, len(a))
	for i, err := range a ***REMOVED***
		errorStrings[i] = err.Error()
	***REMOVED***
	return errorStrings
***REMOVED***

func (a errorMsgs) JSON() interface***REMOVED******REMOVED*** ***REMOVED***
	switch len(a) ***REMOVED***
	case 0:
		return nil
	case 1:
		return a.Last().JSON()
	default:
		json := make([]interface***REMOVED******REMOVED***, len(a))
		for i, err := range a ***REMOVED***
			json[i] = err.JSON()
		***REMOVED***
		return json
	***REMOVED***
***REMOVED***

func (a errorMsgs) MarshalJSON() ([]byte, error) ***REMOVED***
	return json.Marshal(a.JSON())
***REMOVED***

func (a errorMsgs) String() string ***REMOVED***
	if len(a) == 0 ***REMOVED***
		return ""
	***REMOVED***
	var buffer bytes.Buffer
	for i, msg := range a ***REMOVED***
		fmt.Fprintf(&buffer, "Error #%02d: %s\n", (i + 1), msg.Err)
		if msg.Meta != nil ***REMOVED***
			fmt.Fprintf(&buffer, "     Meta: %v\n", msg.Meta)
		***REMOVED***
	***REMOVED***
	return buffer.String()
***REMOVED***
