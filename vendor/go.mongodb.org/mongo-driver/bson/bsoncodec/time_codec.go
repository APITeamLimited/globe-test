// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsoncodec

import (
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsonoptions"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	timeFormatString = "2006-01-02T15:04:05.999Z07:00"
)

// TimeCodec is the Codec used for time.Time values.
type TimeCodec struct ***REMOVED***
	UseLocalTimeZone bool
***REMOVED***

var (
	defaultTimeCodec = NewTimeCodec()

	_ ValueCodec  = defaultTimeCodec
	_ typeDecoder = defaultTimeCodec
)

// NewTimeCodec returns a TimeCodec with options opts.
func NewTimeCodec(opts ...*bsonoptions.TimeCodecOptions) *TimeCodec ***REMOVED***
	timeOpt := bsonoptions.MergeTimeCodecOptions(opts...)

	codec := TimeCodec***REMOVED******REMOVED***
	if timeOpt.UseLocalTimeZone != nil ***REMOVED***
		codec.UseLocalTimeZone = *timeOpt.UseLocalTimeZone
	***REMOVED***
	return &codec
***REMOVED***

func (tc *TimeCodec) decodeType(dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	if t != tTime ***REMOVED***
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "TimeDecodeValue",
			Types:    []reflect.Type***REMOVED***tTime***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***

	var timeVal time.Time
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.DateTime:
		dt, err := vr.ReadDateTime()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		timeVal = time.Unix(dt/1000, dt%1000*1000000)
	case bsontype.String:
		// assume strings are in the isoTimeFormat
		timeStr, err := vr.ReadString()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		timeVal, err = time.Parse(timeFormatString, timeStr)
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
	case bsontype.Int64:
		i64, err := vr.ReadInt64()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		timeVal = time.Unix(i64/1000, i64%1000*1000000)
	case bsontype.Timestamp:
		t, _, err := vr.ReadTimestamp()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		timeVal = time.Unix(int64(t), 0)
	case bsontype.Null:
		if err := vr.ReadNull(); err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
	case bsontype.Undefined:
		if err := vr.ReadUndefined(); err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a time.Time", vrType)
	***REMOVED***

	if !tc.UseLocalTimeZone ***REMOVED***
		timeVal = timeVal.UTC()
	***REMOVED***
	return reflect.ValueOf(timeVal), nil
***REMOVED***

// DecodeValue is the ValueDecoderFunc for time.Time.
func (tc *TimeCodec) DecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tTime ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "TimeDecodeValue", Types: []reflect.Type***REMOVED***tTime***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	elem, err := tc.decodeType(dc, vr, tTime)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.Set(elem)
	return nil
***REMOVED***

// EncodeValue is the ValueEncoderFunc for time.TIme.
func (tc *TimeCodec) EncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tTime ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "TimeEncodeValue", Types: []reflect.Type***REMOVED***tTime***REMOVED***, Received: val***REMOVED***
	***REMOVED***
	tt := val.Interface().(time.Time)
	dt := primitive.NewDateTimeFromTime(tt)
	return vw.WriteDateTime(int64(dt))
***REMOVED***
