// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import "net/http"

const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEMSGPACK           = "application/x-msgpack"
	MIMEMSGPACK2          = "application/msgpack"
)

type Binding interface ***REMOVED***
	Name() string
	Bind(*http.Request, interface***REMOVED******REMOVED***) error
***REMOVED***

type StructValidator interface ***REMOVED***
	// ValidateStruct can receive any kind of type and it should never panic, even if the configuration is not right.
	// If the received type is not a struct, any validation should be skipped and nil must be returned.
	// If the received type is a struct or pointer to a struct, the validation should be performed.
	// If the struct is not valid or the validation itself fails, a descriptive error should be returned.
	// Otherwise nil must be returned.
	ValidateStruct(interface***REMOVED******REMOVED***) error
***REMOVED***

var Validator StructValidator = &defaultValidator***REMOVED******REMOVED***

var (
	JSON          = jsonBinding***REMOVED******REMOVED***
	XML           = xmlBinding***REMOVED******REMOVED***
	Form          = formBinding***REMOVED******REMOVED***
	FormPost      = formPostBinding***REMOVED******REMOVED***
	FormMultipart = formMultipartBinding***REMOVED******REMOVED***
	ProtoBuf      = protobufBinding***REMOVED******REMOVED***
	MsgPack       = msgpackBinding***REMOVED******REMOVED***
)

func Default(method, contentType string) Binding ***REMOVED***
	if method == "GET" ***REMOVED***
		return Form
	***REMOVED***

	switch contentType ***REMOVED***
	case MIMEJSON:
		return JSON
	case MIMEXML, MIMEXML2:
		return XML
	case MIMEPROTOBUF:
		return ProtoBuf
	case MIMEMSGPACK, MIMEMSGPACK2:
		return MsgPack
	default: //case MIMEPOSTForm, MIMEMultipartPOSTForm:
		return Form
	***REMOVED***
***REMOVED***

func validate(obj interface***REMOVED******REMOVED***) error ***REMOVED***
	if Validator == nil ***REMOVED***
		return nil
	***REMOVED***
	return Validator.ValidateStruct(obj)
***REMOVED***
