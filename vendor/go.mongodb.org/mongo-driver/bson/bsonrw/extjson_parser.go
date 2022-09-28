// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonrw

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

const maxNestingDepth = 200

// ErrInvalidJSON indicates the JSON input is invalid
var ErrInvalidJSON = errors.New("invalid JSON input")

type jsonParseState byte

const (
	jpsStartState jsonParseState = iota
	jpsSawBeginObject
	jpsSawEndObject
	jpsSawBeginArray
	jpsSawEndArray
	jpsSawColon
	jpsSawComma
	jpsSawKey
	jpsSawValue
	jpsDoneState
	jpsInvalidState
)

type jsonParseMode byte

const (
	jpmInvalidMode jsonParseMode = iota
	jpmObjectMode
	jpmArrayMode
)

type extJSONValue struct ***REMOVED***
	t bsontype.Type
	v interface***REMOVED******REMOVED***
***REMOVED***

type extJSONObject struct ***REMOVED***
	keys   []string
	values []*extJSONValue
***REMOVED***

type extJSONParser struct ***REMOVED***
	js *jsonScanner
	s  jsonParseState
	m  []jsonParseMode
	k  string
	v  *extJSONValue

	err       error
	canonical bool
	depth     int
	maxDepth  int

	emptyObject bool
	relaxedUUID bool
***REMOVED***

// newExtJSONParser returns a new extended JSON parser, ready to to begin
// parsing from the first character of the argued json input. It will not
// perform any read-ahead and will therefore not report any errors about
// malformed JSON at this point.
func newExtJSONParser(r io.Reader, canonical bool) *extJSONParser ***REMOVED***
	return &extJSONParser***REMOVED***
		js:        &jsonScanner***REMOVED***r: r***REMOVED***,
		s:         jpsStartState,
		m:         []jsonParseMode***REMOVED******REMOVED***,
		canonical: canonical,
		maxDepth:  maxNestingDepth,
	***REMOVED***
***REMOVED***

// peekType examines the next value and returns its BSON Type
func (ejp *extJSONParser) peekType() (bsontype.Type, error) ***REMOVED***
	var t bsontype.Type
	var err error
	initialState := ejp.s

	ejp.advanceState()
	switch ejp.s ***REMOVED***
	case jpsSawValue:
		t = ejp.v.t
	case jpsSawBeginArray:
		t = bsontype.Array
	case jpsInvalidState:
		err = ejp.err
	case jpsSawComma:
		// in array mode, seeing a comma means we need to progress again to actually observe a type
		if ejp.peekMode() == jpmArrayMode ***REMOVED***
			return ejp.peekType()
		***REMOVED***
	case jpsSawEndArray:
		// this would only be a valid state if we were in array mode, so return end-of-array error
		err = ErrEOA
	case jpsSawBeginObject:
		// peek key to determine type
		ejp.advanceState()
		switch ejp.s ***REMOVED***
		case jpsSawEndObject: // empty embedded document
			t = bsontype.EmbeddedDocument
			ejp.emptyObject = true
		case jpsInvalidState:
			err = ejp.err
		case jpsSawKey:
			if initialState == jpsStartState ***REMOVED***
				return bsontype.EmbeddedDocument, nil
			***REMOVED***
			t = wrapperKeyBSONType(ejp.k)

			// if $uuid is encountered, parse as binary subtype 4
			if ejp.k == "$uuid" ***REMOVED***
				ejp.relaxedUUID = true
				t = bsontype.Binary
			***REMOVED***

			switch t ***REMOVED***
			case bsontype.JavaScript:
				// just saw $code, need to check for $scope at same level
				_, err = ejp.readValue(bsontype.JavaScript)
				if err != nil ***REMOVED***
					break
				***REMOVED***

				switch ejp.s ***REMOVED***
				case jpsSawEndObject: // type is TypeJavaScript
				case jpsSawComma:
					ejp.advanceState()

					if ejp.s == jpsSawKey && ejp.k == "$scope" ***REMOVED***
						t = bsontype.CodeWithScope
					***REMOVED*** else ***REMOVED***
						err = fmt.Errorf("invalid extended JSON: unexpected key %s in CodeWithScope object", ejp.k)
					***REMOVED***
				case jpsInvalidState:
					err = ejp.err
				default:
					err = ErrInvalidJSON
				***REMOVED***
			case bsontype.CodeWithScope:
				err = errors.New("invalid extended JSON: code with $scope must contain $code before $scope")
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return t, err
***REMOVED***

// readKey parses the next key and its type and returns them
func (ejp *extJSONParser) readKey() (string, bsontype.Type, error) ***REMOVED***
	if ejp.emptyObject ***REMOVED***
		ejp.emptyObject = false
		return "", 0, ErrEOD
	***REMOVED***

	// advance to key (or return with error)
	switch ejp.s ***REMOVED***
	case jpsStartState:
		ejp.advanceState()
		if ejp.s == jpsSawBeginObject ***REMOVED***
			ejp.advanceState()
		***REMOVED***
	case jpsSawBeginObject:
		ejp.advanceState()
	case jpsSawValue, jpsSawEndObject, jpsSawEndArray:
		ejp.advanceState()
		switch ejp.s ***REMOVED***
		case jpsSawBeginObject, jpsSawComma:
			ejp.advanceState()
		case jpsSawEndObject:
			return "", 0, ErrEOD
		case jpsDoneState:
			return "", 0, io.EOF
		case jpsInvalidState:
			return "", 0, ejp.err
		default:
			return "", 0, ErrInvalidJSON
		***REMOVED***
	case jpsSawKey: // do nothing (key was peeked before)
	default:
		return "", 0, invalidRequestError("key")
	***REMOVED***

	// read key
	var key string

	switch ejp.s ***REMOVED***
	case jpsSawKey:
		key = ejp.k
	case jpsSawEndObject:
		return "", 0, ErrEOD
	case jpsInvalidState:
		return "", 0, ejp.err
	default:
		return "", 0, invalidRequestError("key")
	***REMOVED***

	// check for colon
	ejp.advanceState()
	if err := ensureColon(ejp.s, key); err != nil ***REMOVED***
		return "", 0, err
	***REMOVED***

	// peek at the value to determine type
	t, err := ejp.peekType()
	if err != nil ***REMOVED***
		return "", 0, err
	***REMOVED***

	return key, t, nil
***REMOVED***

// readValue returns the value corresponding to the Type returned by peekType
func (ejp *extJSONParser) readValue(t bsontype.Type) (*extJSONValue, error) ***REMOVED***
	if ejp.s == jpsInvalidState ***REMOVED***
		return nil, ejp.err
	***REMOVED***

	var v *extJSONValue

	switch t ***REMOVED***
	case bsontype.Null, bsontype.Boolean, bsontype.String:
		if ejp.s != jpsSawValue ***REMOVED***
			return nil, invalidRequestError(t.String())
		***REMOVED***
		v = ejp.v
	case bsontype.Int32, bsontype.Int64, bsontype.Double:
		// relaxed version allows these to be literal number values
		if ejp.s == jpsSawValue ***REMOVED***
			v = ejp.v
			break
		***REMOVED***
		fallthrough
	case bsontype.Decimal128, bsontype.Symbol, bsontype.ObjectID, bsontype.MinKey, bsontype.MaxKey, bsontype.Undefined:
		switch ejp.s ***REMOVED***
		case jpsSawKey:
			// read colon
			ejp.advanceState()
			if err := ensureColon(ejp.s, ejp.k); err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			// read value
			ejp.advanceState()
			if ejp.s != jpsSawValue || !ejp.ensureExtValueType(t) ***REMOVED***
				return nil, invalidJSONErrorForType("value", t)
			***REMOVED***

			v = ejp.v

			// read end object
			ejp.advanceState()
			if ejp.s != jpsSawEndObject ***REMOVED***
				return nil, invalidJSONErrorForType("***REMOVED*** after value", t)
			***REMOVED***
		default:
			return nil, invalidRequestError(t.String())
		***REMOVED***
	case bsontype.Binary, bsontype.Regex, bsontype.Timestamp, bsontype.DBPointer:
		if ejp.s != jpsSawKey ***REMOVED***
			return nil, invalidRequestError(t.String())
		***REMOVED***
		// read colon
		ejp.advanceState()
		if err := ensureColon(ejp.s, ejp.k); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		ejp.advanceState()
		if t == bsontype.Binary && ejp.s == jpsSawValue ***REMOVED***
			// convert relaxed $uuid format
			if ejp.relaxedUUID ***REMOVED***
				defer func() ***REMOVED*** ejp.relaxedUUID = false ***REMOVED***()
				uuid, err := ejp.v.parseSymbol()
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***

				// RFC 4122 defines the length of a UUID as 36 and the hyphens in a UUID as appearing
				// in the 8th, 13th, 18th, and 23rd characters.
				//
				// See https://tools.ietf.org/html/rfc4122#section-3
				valid := len(uuid) == 36 &&
					string(uuid[8]) == "-" &&
					string(uuid[13]) == "-" &&
					string(uuid[18]) == "-" &&
					string(uuid[23]) == "-"
				if !valid ***REMOVED***
					return nil, fmt.Errorf("$uuid value does not follow RFC 4122 format regarding length and hyphens")
				***REMOVED***

				// remove hyphens
				uuidNoHyphens := strings.Replace(uuid, "-", "", -1)
				if len(uuidNoHyphens) != 32 ***REMOVED***
					return nil, fmt.Errorf("$uuid value does not follow RFC 4122 format regarding length and hyphens")
				***REMOVED***

				// convert hex to bytes
				bytes, err := hex.DecodeString(uuidNoHyphens)
				if err != nil ***REMOVED***
					return nil, fmt.Errorf("$uuid value does not follow RFC 4122 format regarding hex bytes: %v", err)
				***REMOVED***

				ejp.advanceState()
				if ejp.s != jpsSawEndObject ***REMOVED***
					return nil, invalidJSONErrorForType("$uuid and value and then ***REMOVED***", bsontype.Binary)
				***REMOVED***

				base64 := &extJSONValue***REMOVED***
					t: bsontype.String,
					v: base64.StdEncoding.EncodeToString(bytes),
				***REMOVED***
				subType := &extJSONValue***REMOVED***
					t: bsontype.String,
					v: "04",
				***REMOVED***

				v = &extJSONValue***REMOVED***
					t: bsontype.EmbeddedDocument,
					v: &extJSONObject***REMOVED***
						keys:   []string***REMOVED***"base64", "subType"***REMOVED***,
						values: []*extJSONValue***REMOVED***base64, subType***REMOVED***,
					***REMOVED***,
				***REMOVED***

				break
			***REMOVED***

			// convert legacy $binary format
			base64 := ejp.v

			ejp.advanceState()
			if ejp.s != jpsSawComma ***REMOVED***
				return nil, invalidJSONErrorForType(",", bsontype.Binary)
			***REMOVED***

			ejp.advanceState()
			key, t, err := ejp.readKey()
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if key != "$type" ***REMOVED***
				return nil, invalidJSONErrorForType("$type", bsontype.Binary)
			***REMOVED***

			subType, err := ejp.readValue(t)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			ejp.advanceState()
			if ejp.s != jpsSawEndObject ***REMOVED***
				return nil, invalidJSONErrorForType("2 key-value pairs and then ***REMOVED***", bsontype.Binary)
			***REMOVED***

			v = &extJSONValue***REMOVED***
				t: bsontype.EmbeddedDocument,
				v: &extJSONObject***REMOVED***
					keys:   []string***REMOVED***"base64", "subType"***REMOVED***,
					values: []*extJSONValue***REMOVED***base64, subType***REMOVED***,
				***REMOVED***,
			***REMOVED***
			break
		***REMOVED***

		// read KV pairs
		if ejp.s != jpsSawBeginObject ***REMOVED***
			return nil, invalidJSONErrorForType("***REMOVED***", t)
		***REMOVED***

		keys, vals, err := ejp.readObject(2, true)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		ejp.advanceState()
		if ejp.s != jpsSawEndObject ***REMOVED***
			return nil, invalidJSONErrorForType("2 key-value pairs and then ***REMOVED***", t)
		***REMOVED***

		v = &extJSONValue***REMOVED***t: bsontype.EmbeddedDocument, v: &extJSONObject***REMOVED***keys: keys, values: vals***REMOVED******REMOVED***

	case bsontype.DateTime:
		switch ejp.s ***REMOVED***
		case jpsSawValue:
			v = ejp.v
		case jpsSawKey:
			// read colon
			ejp.advanceState()
			if err := ensureColon(ejp.s, ejp.k); err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			ejp.advanceState()
			switch ejp.s ***REMOVED***
			case jpsSawBeginObject:
				keys, vals, err := ejp.readObject(1, true)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				v = &extJSONValue***REMOVED***t: bsontype.EmbeddedDocument, v: &extJSONObject***REMOVED***keys: keys, values: vals***REMOVED******REMOVED***
			case jpsSawValue:
				if ejp.canonical ***REMOVED***
					return nil, invalidJSONError("***REMOVED***")
				***REMOVED***
				v = ejp.v
			default:
				if ejp.canonical ***REMOVED***
					return nil, invalidJSONErrorForType("object", t)
				***REMOVED***
				return nil, invalidJSONErrorForType("ISO-8601 Internet Date/Time Format as described in RFC-3339", t)
			***REMOVED***

			ejp.advanceState()
			if ejp.s != jpsSawEndObject ***REMOVED***
				return nil, invalidJSONErrorForType("value and then ***REMOVED***", t)
			***REMOVED***
		default:
			return nil, invalidRequestError(t.String())
		***REMOVED***
	case bsontype.JavaScript:
		switch ejp.s ***REMOVED***
		case jpsSawKey:
			// read colon
			ejp.advanceState()
			if err := ensureColon(ejp.s, ejp.k); err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			// read value
			ejp.advanceState()
			if ejp.s != jpsSawValue ***REMOVED***
				return nil, invalidJSONErrorForType("value", t)
			***REMOVED***
			v = ejp.v

			// read end object or comma and just return
			ejp.advanceState()
		case jpsSawEndObject:
			v = ejp.v
		default:
			return nil, invalidRequestError(t.String())
		***REMOVED***
	case bsontype.CodeWithScope:
		if ejp.s == jpsSawKey && ejp.k == "$scope" ***REMOVED***
			v = ejp.v // this is the $code string from earlier

			// read colon
			ejp.advanceState()
			if err := ensureColon(ejp.s, ejp.k); err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			// read ***REMOVED***
			ejp.advanceState()
			if ejp.s != jpsSawBeginObject ***REMOVED***
				return nil, invalidJSONError("$scope to be embedded document")
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			return nil, invalidRequestError(t.String())
		***REMOVED***
	case bsontype.EmbeddedDocument, bsontype.Array:
		return nil, invalidRequestError(t.String())
	***REMOVED***

	return v, nil
***REMOVED***

// readObject is a utility method for reading full objects of known (or expected) size
// it is useful for extended JSON types such as binary, datetime, regex, and timestamp
func (ejp *extJSONParser) readObject(numKeys int, started bool) ([]string, []*extJSONValue, error) ***REMOVED***
	keys := make([]string, numKeys)
	vals := make([]*extJSONValue, numKeys)

	if !started ***REMOVED***
		ejp.advanceState()
		if ejp.s != jpsSawBeginObject ***REMOVED***
			return nil, nil, invalidJSONError("***REMOVED***")
		***REMOVED***
	***REMOVED***

	for i := 0; i < numKeys; i++ ***REMOVED***
		key, t, err := ejp.readKey()
		if err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***

		switch ejp.s ***REMOVED***
		case jpsSawKey:
			v, err := ejp.readValue(t)
			if err != nil ***REMOVED***
				return nil, nil, err
			***REMOVED***

			keys[i] = key
			vals[i] = v
		case jpsSawValue:
			keys[i] = key
			vals[i] = ejp.v
		default:
			return nil, nil, invalidJSONError("value")
		***REMOVED***
	***REMOVED***

	ejp.advanceState()
	if ejp.s != jpsSawEndObject ***REMOVED***
		return nil, nil, invalidJSONError("***REMOVED***")
	***REMOVED***

	return keys, vals, nil
***REMOVED***

// advanceState reads the next JSON token from the scanner and transitions
// from the current state based on that token's type
func (ejp *extJSONParser) advanceState() ***REMOVED***
	if ejp.s == jpsDoneState || ejp.s == jpsInvalidState ***REMOVED***
		return
	***REMOVED***

	jt, err := ejp.js.nextToken()

	if err != nil ***REMOVED***
		ejp.err = err
		ejp.s = jpsInvalidState
		return
	***REMOVED***

	valid := ejp.validateToken(jt.t)
	if !valid ***REMOVED***
		ejp.err = unexpectedTokenError(jt)
		ejp.s = jpsInvalidState
		return
	***REMOVED***

	switch jt.t ***REMOVED***
	case jttBeginObject:
		ejp.s = jpsSawBeginObject
		ejp.pushMode(jpmObjectMode)
		ejp.depth++

		if ejp.depth > ejp.maxDepth ***REMOVED***
			ejp.err = nestingDepthError(jt.p, ejp.depth)
			ejp.s = jpsInvalidState
		***REMOVED***
	case jttEndObject:
		ejp.s = jpsSawEndObject
		ejp.depth--

		if ejp.popMode() != jpmObjectMode ***REMOVED***
			ejp.err = unexpectedTokenError(jt)
			ejp.s = jpsInvalidState
		***REMOVED***
	case jttBeginArray:
		ejp.s = jpsSawBeginArray
		ejp.pushMode(jpmArrayMode)
	case jttEndArray:
		ejp.s = jpsSawEndArray

		if ejp.popMode() != jpmArrayMode ***REMOVED***
			ejp.err = unexpectedTokenError(jt)
			ejp.s = jpsInvalidState
		***REMOVED***
	case jttColon:
		ejp.s = jpsSawColon
	case jttComma:
		ejp.s = jpsSawComma
	case jttEOF:
		ejp.s = jpsDoneState
		if len(ejp.m) != 0 ***REMOVED***
			ejp.err = unexpectedTokenError(jt)
			ejp.s = jpsInvalidState
		***REMOVED***
	case jttString:
		switch ejp.s ***REMOVED***
		case jpsSawComma:
			if ejp.peekMode() == jpmArrayMode ***REMOVED***
				ejp.s = jpsSawValue
				ejp.v = extendJSONToken(jt)
				return
			***REMOVED***
			fallthrough
		case jpsSawBeginObject:
			ejp.s = jpsSawKey
			ejp.k = jt.v.(string)
			return
		***REMOVED***
		fallthrough
	default:
		ejp.s = jpsSawValue
		ejp.v = extendJSONToken(jt)
	***REMOVED***
***REMOVED***

var jpsValidTransitionTokens = map[jsonParseState]map[jsonTokenType]bool***REMOVED***
	jpsStartState: ***REMOVED***
		jttBeginObject: true,
		jttBeginArray:  true,
		jttInt32:       true,
		jttInt64:       true,
		jttDouble:      true,
		jttString:      true,
		jttBool:        true,
		jttNull:        true,
		jttEOF:         true,
	***REMOVED***,
	jpsSawBeginObject: ***REMOVED***
		jttEndObject: true,
		jttString:    true,
	***REMOVED***,
	jpsSawEndObject: ***REMOVED***
		jttEndObject: true,
		jttEndArray:  true,
		jttComma:     true,
		jttEOF:       true,
	***REMOVED***,
	jpsSawBeginArray: ***REMOVED***
		jttBeginObject: true,
		jttBeginArray:  true,
		jttEndArray:    true,
		jttInt32:       true,
		jttInt64:       true,
		jttDouble:      true,
		jttString:      true,
		jttBool:        true,
		jttNull:        true,
	***REMOVED***,
	jpsSawEndArray: ***REMOVED***
		jttEndObject: true,
		jttEndArray:  true,
		jttComma:     true,
		jttEOF:       true,
	***REMOVED***,
	jpsSawColon: ***REMOVED***
		jttBeginObject: true,
		jttBeginArray:  true,
		jttInt32:       true,
		jttInt64:       true,
		jttDouble:      true,
		jttString:      true,
		jttBool:        true,
		jttNull:        true,
	***REMOVED***,
	jpsSawComma: ***REMOVED***
		jttBeginObject: true,
		jttBeginArray:  true,
		jttInt32:       true,
		jttInt64:       true,
		jttDouble:      true,
		jttString:      true,
		jttBool:        true,
		jttNull:        true,
	***REMOVED***,
	jpsSawKey: ***REMOVED***
		jttColon: true,
	***REMOVED***,
	jpsSawValue: ***REMOVED***
		jttEndObject: true,
		jttEndArray:  true,
		jttComma:     true,
		jttEOF:       true,
	***REMOVED***,
	jpsDoneState:    ***REMOVED******REMOVED***,
	jpsInvalidState: ***REMOVED******REMOVED***,
***REMOVED***

func (ejp *extJSONParser) validateToken(jtt jsonTokenType) bool ***REMOVED***
	switch ejp.s ***REMOVED***
	case jpsSawEndObject:
		// if we are at depth zero and the next token is a '***REMOVED***',
		// we can consider it valid only if we are not in array mode.
		if jtt == jttBeginObject && ejp.depth == 0 ***REMOVED***
			return ejp.peekMode() != jpmArrayMode
		***REMOVED***
	case jpsSawComma:
		switch ejp.peekMode() ***REMOVED***
		// the only valid next token after a comma inside a document is a string (a key)
		case jpmObjectMode:
			return jtt == jttString
		case jpmInvalidMode:
			return false
		***REMOVED***
	***REMOVED***

	_, ok := jpsValidTransitionTokens[ejp.s][jtt]
	return ok
***REMOVED***

// ensureExtValueType returns true if the current value has the expected
// value type for single-key extended JSON types. For example,
// ***REMOVED***"$numberInt": v***REMOVED*** v must be TypeString
func (ejp *extJSONParser) ensureExtValueType(t bsontype.Type) bool ***REMOVED***
	switch t ***REMOVED***
	case bsontype.MinKey, bsontype.MaxKey:
		return ejp.v.t == bsontype.Int32
	case bsontype.Undefined:
		return ejp.v.t == bsontype.Boolean
	case bsontype.Int32, bsontype.Int64, bsontype.Double, bsontype.Decimal128, bsontype.Symbol, bsontype.ObjectID:
		return ejp.v.t == bsontype.String
	default:
		return false
	***REMOVED***
***REMOVED***

func (ejp *extJSONParser) pushMode(m jsonParseMode) ***REMOVED***
	ejp.m = append(ejp.m, m)
***REMOVED***

func (ejp *extJSONParser) popMode() jsonParseMode ***REMOVED***
	l := len(ejp.m)
	if l == 0 ***REMOVED***
		return jpmInvalidMode
	***REMOVED***

	m := ejp.m[l-1]
	ejp.m = ejp.m[:l-1]

	return m
***REMOVED***

func (ejp *extJSONParser) peekMode() jsonParseMode ***REMOVED***
	l := len(ejp.m)
	if l == 0 ***REMOVED***
		return jpmInvalidMode
	***REMOVED***

	return ejp.m[l-1]
***REMOVED***

func extendJSONToken(jt *jsonToken) *extJSONValue ***REMOVED***
	var t bsontype.Type

	switch jt.t ***REMOVED***
	case jttInt32:
		t = bsontype.Int32
	case jttInt64:
		t = bsontype.Int64
	case jttDouble:
		t = bsontype.Double
	case jttString:
		t = bsontype.String
	case jttBool:
		t = bsontype.Boolean
	case jttNull:
		t = bsontype.Null
	default:
		return nil
	***REMOVED***

	return &extJSONValue***REMOVED***t: t, v: jt.v***REMOVED***
***REMOVED***

func ensureColon(s jsonParseState, key string) error ***REMOVED***
	if s != jpsSawColon ***REMOVED***
		return fmt.Errorf("invalid JSON input: missing colon after key \"%s\"", key)
	***REMOVED***

	return nil
***REMOVED***

func invalidRequestError(s string) error ***REMOVED***
	return fmt.Errorf("invalid request to read %s", s)
***REMOVED***

func invalidJSONError(expected string) error ***REMOVED***
	return fmt.Errorf("invalid JSON input; expected %s", expected)
***REMOVED***

func invalidJSONErrorForType(expected string, t bsontype.Type) error ***REMOVED***
	return fmt.Errorf("invalid JSON input; expected %s for %s", expected, t)
***REMOVED***

func unexpectedTokenError(jt *jsonToken) error ***REMOVED***
	switch jt.t ***REMOVED***
	case jttInt32, jttInt64, jttDouble:
		return fmt.Errorf("invalid JSON input; unexpected number (%v) at position %d", jt.v, jt.p)
	case jttString:
		return fmt.Errorf("invalid JSON input; unexpected string (\"%v\") at position %d", jt.v, jt.p)
	case jttBool:
		return fmt.Errorf("invalid JSON input; unexpected boolean literal (%v) at position %d", jt.v, jt.p)
	case jttNull:
		return fmt.Errorf("invalid JSON input; unexpected null literal at position %d", jt.p)
	case jttEOF:
		return fmt.Errorf("invalid JSON input; unexpected end of input at position %d", jt.p)
	default:
		return fmt.Errorf("invalid JSON input; unexpected %c at position %d", jt.v.(byte), jt.p)
	***REMOVED***
***REMOVED***

func nestingDepthError(p, depth int) error ***REMOVED***
	return fmt.Errorf("invalid JSON input; nesting too deep (%d levels) at position %d", depth, p)
***REMOVED***
