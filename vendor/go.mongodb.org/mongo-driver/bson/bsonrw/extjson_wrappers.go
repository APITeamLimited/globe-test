// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonrw

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func wrapperKeyBSONType(key string) bsontype.Type ***REMOVED***
	switch key ***REMOVED***
	case "$numberInt":
		return bsontype.Int32
	case "$numberLong":
		return bsontype.Int64
	case "$oid":
		return bsontype.ObjectID
	case "$symbol":
		return bsontype.Symbol
	case "$numberDouble":
		return bsontype.Double
	case "$numberDecimal":
		return bsontype.Decimal128
	case "$binary":
		return bsontype.Binary
	case "$code":
		return bsontype.JavaScript
	case "$scope":
		return bsontype.CodeWithScope
	case "$timestamp":
		return bsontype.Timestamp
	case "$regularExpression":
		return bsontype.Regex
	case "$dbPointer":
		return bsontype.DBPointer
	case "$date":
		return bsontype.DateTime
	case "$minKey":
		return bsontype.MinKey
	case "$maxKey":
		return bsontype.MaxKey
	case "$undefined":
		return bsontype.Undefined
	***REMOVED***

	return bsontype.EmbeddedDocument
***REMOVED***

func (ejv *extJSONValue) parseBinary() (b []byte, subType byte, err error) ***REMOVED***
	if ejv.t != bsontype.EmbeddedDocument ***REMOVED***
		return nil, 0, fmt.Errorf("$binary value should be object, but instead is %s", ejv.t)
	***REMOVED***

	binObj := ejv.v.(*extJSONObject)
	bFound := false
	stFound := false

	for i, key := range binObj.keys ***REMOVED***
		val := binObj.values[i]

		switch key ***REMOVED***
		case "base64":
			if bFound ***REMOVED***
				return nil, 0, errors.New("duplicate base64 key in $binary")
			***REMOVED***

			if val.t != bsontype.String ***REMOVED***
				return nil, 0, fmt.Errorf("$binary base64 value should be string, but instead is %s", val.t)
			***REMOVED***

			base64Bytes, err := base64.StdEncoding.DecodeString(val.v.(string))
			if err != nil ***REMOVED***
				return nil, 0, fmt.Errorf("invalid $binary base64 string: %s", val.v.(string))
			***REMOVED***

			b = base64Bytes
			bFound = true
		case "subType":
			if stFound ***REMOVED***
				return nil, 0, errors.New("duplicate subType key in $binary")
			***REMOVED***

			if val.t != bsontype.String ***REMOVED***
				return nil, 0, fmt.Errorf("$binary subType value should be string, but instead is %s", val.t)
			***REMOVED***

			i, err := strconv.ParseInt(val.v.(string), 16, 64)
			if err != nil ***REMOVED***
				return nil, 0, fmt.Errorf("invalid $binary subType string: %s", val.v.(string))
			***REMOVED***

			subType = byte(i)
			stFound = true
		default:
			return nil, 0, fmt.Errorf("invalid key in $binary object: %s", key)
		***REMOVED***
	***REMOVED***

	if !bFound ***REMOVED***
		return nil, 0, errors.New("missing base64 field in $binary object")
	***REMOVED***

	if !stFound ***REMOVED***
		return nil, 0, errors.New("missing subType field in $binary object")

	***REMOVED***

	return b, subType, nil
***REMOVED***

func (ejv *extJSONValue) parseDBPointer() (ns string, oid primitive.ObjectID, err error) ***REMOVED***
	if ejv.t != bsontype.EmbeddedDocument ***REMOVED***
		return "", primitive.NilObjectID, fmt.Errorf("$dbPointer value should be object, but instead is %s", ejv.t)
	***REMOVED***

	dbpObj := ejv.v.(*extJSONObject)
	oidFound := false
	nsFound := false

	for i, key := range dbpObj.keys ***REMOVED***
		val := dbpObj.values[i]

		switch key ***REMOVED***
		case "$ref":
			if nsFound ***REMOVED***
				return "", primitive.NilObjectID, errors.New("duplicate $ref key in $dbPointer")
			***REMOVED***

			if val.t != bsontype.String ***REMOVED***
				return "", primitive.NilObjectID, fmt.Errorf("$dbPointer $ref value should be string, but instead is %s", val.t)
			***REMOVED***

			ns = val.v.(string)
			nsFound = true
		case "$id":
			if oidFound ***REMOVED***
				return "", primitive.NilObjectID, errors.New("duplicate $id key in $dbPointer")
			***REMOVED***

			if val.t != bsontype.String ***REMOVED***
				return "", primitive.NilObjectID, fmt.Errorf("$dbPointer $id value should be string, but instead is %s", val.t)
			***REMOVED***

			oid, err = primitive.ObjectIDFromHex(val.v.(string))
			if err != nil ***REMOVED***
				return "", primitive.NilObjectID, err
			***REMOVED***

			oidFound = true
		default:
			return "", primitive.NilObjectID, fmt.Errorf("invalid key in $dbPointer object: %s", key)
		***REMOVED***
	***REMOVED***

	if !nsFound ***REMOVED***
		return "", oid, errors.New("missing $ref field in $dbPointer object")
	***REMOVED***

	if !oidFound ***REMOVED***
		return "", oid, errors.New("missing $id field in $dbPointer object")
	***REMOVED***

	return ns, oid, nil
***REMOVED***

const (
	rfc3339Milli = "2006-01-02T15:04:05.999Z07:00"
)

var (
	timeFormats = []string***REMOVED***rfc3339Milli, "2006-01-02T15:04:05.999Z0700"***REMOVED***
)

func (ejv *extJSONValue) parseDateTime() (int64, error) ***REMOVED***
	switch ejv.t ***REMOVED***
	case bsontype.Int32:
		return int64(ejv.v.(int32)), nil
	case bsontype.Int64:
		return ejv.v.(int64), nil
	case bsontype.String:
		return parseDatetimeString(ejv.v.(string))
	case bsontype.EmbeddedDocument:
		return parseDatetimeObject(ejv.v.(*extJSONObject))
	default:
		return 0, fmt.Errorf("$date value should be string or object, but instead is %s", ejv.t)
	***REMOVED***
***REMOVED***

func parseDatetimeString(data string) (int64, error) ***REMOVED***
	var t time.Time
	var err error
	// try acceptable time formats until one matches
	for _, format := range timeFormats ***REMOVED***
		t, err = time.Parse(format, data)
		if err == nil ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	if err != nil ***REMOVED***
		return 0, fmt.Errorf("invalid $date value string: %s", data)
	***REMOVED***

	return int64(primitive.NewDateTimeFromTime(t)), nil
***REMOVED***

func parseDatetimeObject(data *extJSONObject) (d int64, err error) ***REMOVED***
	dFound := false

	for i, key := range data.keys ***REMOVED***
		val := data.values[i]

		switch key ***REMOVED***
		case "$numberLong":
			if dFound ***REMOVED***
				return 0, errors.New("duplicate $numberLong key in $date")
			***REMOVED***

			if val.t != bsontype.String ***REMOVED***
				return 0, fmt.Errorf("$date $numberLong field should be string, but instead is %s", val.t)
			***REMOVED***

			d, err = val.parseInt64()
			if err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			dFound = true
		default:
			return 0, fmt.Errorf("invalid key in $date object: %s", key)
		***REMOVED***
	***REMOVED***

	if !dFound ***REMOVED***
		return 0, errors.New("missing $numberLong field in $date object")
	***REMOVED***

	return d, nil
***REMOVED***

func (ejv *extJSONValue) parseDecimal128() (primitive.Decimal128, error) ***REMOVED***
	if ejv.t != bsontype.String ***REMOVED***
		return primitive.Decimal128***REMOVED******REMOVED***, fmt.Errorf("$numberDecimal value should be string, but instead is %s", ejv.t)
	***REMOVED***

	d, err := primitive.ParseDecimal128(ejv.v.(string))
	if err != nil ***REMOVED***
		return primitive.Decimal128***REMOVED******REMOVED***, fmt.Errorf("$invalid $numberDecimal string: %s", ejv.v.(string))
	***REMOVED***

	return d, nil
***REMOVED***

func (ejv *extJSONValue) parseDouble() (float64, error) ***REMOVED***
	if ejv.t == bsontype.Double ***REMOVED***
		return ejv.v.(float64), nil
	***REMOVED***

	if ejv.t != bsontype.String ***REMOVED***
		return 0, fmt.Errorf("$numberDouble value should be string, but instead is %s", ejv.t)
	***REMOVED***

	switch ejv.v.(string) ***REMOVED***
	case "Infinity":
		return math.Inf(1), nil
	case "-Infinity":
		return math.Inf(-1), nil
	case "NaN":
		return math.NaN(), nil
	***REMOVED***

	f, err := strconv.ParseFloat(ejv.v.(string), 64)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	return f, nil
***REMOVED***

func (ejv *extJSONValue) parseInt32() (int32, error) ***REMOVED***
	if ejv.t == bsontype.Int32 ***REMOVED***
		return ejv.v.(int32), nil
	***REMOVED***

	if ejv.t != bsontype.String ***REMOVED***
		return 0, fmt.Errorf("$numberInt value should be string, but instead is %s", ejv.t)
	***REMOVED***

	i, err := strconv.ParseInt(ejv.v.(string), 10, 64)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	if i < math.MinInt32 || i > math.MaxInt32 ***REMOVED***
		return 0, fmt.Errorf("$numberInt value should be int32 but instead is int64: %d", i)
	***REMOVED***

	return int32(i), nil
***REMOVED***

func (ejv *extJSONValue) parseInt64() (int64, error) ***REMOVED***
	if ejv.t == bsontype.Int64 ***REMOVED***
		return ejv.v.(int64), nil
	***REMOVED***

	if ejv.t != bsontype.String ***REMOVED***
		return 0, fmt.Errorf("$numberLong value should be string, but instead is %s", ejv.t)
	***REMOVED***

	i, err := strconv.ParseInt(ejv.v.(string), 10, 64)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	return i, nil
***REMOVED***

func (ejv *extJSONValue) parseJavascript() (code string, err error) ***REMOVED***
	if ejv.t != bsontype.String ***REMOVED***
		return "", fmt.Errorf("$code value should be string, but instead is %s", ejv.t)
	***REMOVED***

	return ejv.v.(string), nil
***REMOVED***

func (ejv *extJSONValue) parseMinMaxKey(minmax string) error ***REMOVED***
	if ejv.t != bsontype.Int32 ***REMOVED***
		return fmt.Errorf("$%sKey value should be int32, but instead is %s", minmax, ejv.t)
	***REMOVED***

	if ejv.v.(int32) != 1 ***REMOVED***
		return fmt.Errorf("$%sKey value must be 1, but instead is %d", minmax, ejv.v.(int32))
	***REMOVED***

	return nil
***REMOVED***

func (ejv *extJSONValue) parseObjectID() (primitive.ObjectID, error) ***REMOVED***
	if ejv.t != bsontype.String ***REMOVED***
		return primitive.NilObjectID, fmt.Errorf("$oid value should be string, but instead is %s", ejv.t)
	***REMOVED***

	return primitive.ObjectIDFromHex(ejv.v.(string))
***REMOVED***

func (ejv *extJSONValue) parseRegex() (pattern, options string, err error) ***REMOVED***
	if ejv.t != bsontype.EmbeddedDocument ***REMOVED***
		return "", "", fmt.Errorf("$regularExpression value should be object, but instead is %s", ejv.t)
	***REMOVED***

	regexObj := ejv.v.(*extJSONObject)
	patFound := false
	optFound := false

	for i, key := range regexObj.keys ***REMOVED***
		val := regexObj.values[i]

		switch key ***REMOVED***
		case "pattern":
			if patFound ***REMOVED***
				return "", "", errors.New("duplicate pattern key in $regularExpression")
			***REMOVED***

			if val.t != bsontype.String ***REMOVED***
				return "", "", fmt.Errorf("$regularExpression pattern value should be string, but instead is %s", val.t)
			***REMOVED***

			pattern = val.v.(string)
			patFound = true
		case "options":
			if optFound ***REMOVED***
				return "", "", errors.New("duplicate options key in $regularExpression")
			***REMOVED***

			if val.t != bsontype.String ***REMOVED***
				return "", "", fmt.Errorf("$regularExpression options value should be string, but instead is %s", val.t)
			***REMOVED***

			options = val.v.(string)
			optFound = true
		default:
			return "", "", fmt.Errorf("invalid key in $regularExpression object: %s", key)
		***REMOVED***
	***REMOVED***

	if !patFound ***REMOVED***
		return "", "", errors.New("missing pattern field in $regularExpression object")
	***REMOVED***

	if !optFound ***REMOVED***
		return "", "", errors.New("missing options field in $regularExpression object")

	***REMOVED***

	return pattern, options, nil
***REMOVED***

func (ejv *extJSONValue) parseSymbol() (string, error) ***REMOVED***
	if ejv.t != bsontype.String ***REMOVED***
		return "", fmt.Errorf("$symbol value should be string, but instead is %s", ejv.t)
	***REMOVED***

	return ejv.v.(string), nil
***REMOVED***

func (ejv *extJSONValue) parseTimestamp() (t, i uint32, err error) ***REMOVED***
	if ejv.t != bsontype.EmbeddedDocument ***REMOVED***
		return 0, 0, fmt.Errorf("$timestamp value should be object, but instead is %s", ejv.t)
	***REMOVED***

	handleKey := func(key string, val *extJSONValue, flag bool) (uint32, error) ***REMOVED***
		if flag ***REMOVED***
			return 0, fmt.Errorf("duplicate %s key in $timestamp", key)
		***REMOVED***

		switch val.t ***REMOVED***
		case bsontype.Int32:
			value := val.v.(int32)

			if value < 0 ***REMOVED***
				return 0, fmt.Errorf("$timestamp %s number should be uint32: %d", key, value)
			***REMOVED***

			return uint32(value), nil
		case bsontype.Int64:
			value := val.v.(int64)
			if value < 0 || value > int64(math.MaxUint32) ***REMOVED***
				return 0, fmt.Errorf("$timestamp %s number should be uint32: %d", key, value)
			***REMOVED***

			return uint32(value), nil
		default:
			return 0, fmt.Errorf("$timestamp %s value should be uint32, but instead is %s", key, val.t)
		***REMOVED***
	***REMOVED***

	tsObj := ejv.v.(*extJSONObject)
	tFound := false
	iFound := false

	for j, key := range tsObj.keys ***REMOVED***
		val := tsObj.values[j]

		switch key ***REMOVED***
		case "t":
			if t, err = handleKey(key, val, tFound); err != nil ***REMOVED***
				return 0, 0, err
			***REMOVED***

			tFound = true
		case "i":
			if i, err = handleKey(key, val, iFound); err != nil ***REMOVED***
				return 0, 0, err
			***REMOVED***

			iFound = true
		default:
			return 0, 0, fmt.Errorf("invalid key in $timestamp object: %s", key)
		***REMOVED***
	***REMOVED***

	if !tFound ***REMOVED***
		return 0, 0, errors.New("missing t field in $timestamp object")
	***REMOVED***

	if !iFound ***REMOVED***
		return 0, 0, errors.New("missing i field in $timestamp object")
	***REMOVED***

	return t, i, nil
***REMOVED***

func (ejv *extJSONValue) parseUndefined() error ***REMOVED***
	if ejv.t != bsontype.Boolean ***REMOVED***
		return fmt.Errorf("undefined value should be boolean, but instead is %s", ejv.t)
	***REMOVED***

	if !ejv.v.(bool) ***REMOVED***
		return fmt.Errorf("$undefined balue boolean should be true, but instead is %v", ejv.v.(bool))
	***REMOVED***

	return nil
***REMOVED***
