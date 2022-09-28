// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo // import "go.mongodb.org/mongo-driver/mongo"

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Dialer is used to make network connections.
type Dialer interface ***REMOVED***
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
***REMOVED***

// BSONAppender is an interface implemented by types that can marshal a
// provided type into BSON bytes and append those bytes to the provided []byte.
// The AppendBSON can return a non-nil error and non-nil []byte. The AppendBSON
// method may also write incomplete BSON to the []byte.
type BSONAppender interface ***REMOVED***
	AppendBSON([]byte, interface***REMOVED******REMOVED***) ([]byte, error)
***REMOVED***

// BSONAppenderFunc is an adapter function that allows any function that
// satisfies the AppendBSON method signature to be used where a BSONAppender is
// used.
type BSONAppenderFunc func([]byte, interface***REMOVED******REMOVED***) ([]byte, error)

// AppendBSON implements the BSONAppender interface
func (baf BSONAppenderFunc) AppendBSON(dst []byte, val interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return baf(dst, val)
***REMOVED***

// MarshalError is returned when attempting to transform a value into a document
// results in an error.
type MarshalError struct ***REMOVED***
	Value interface***REMOVED******REMOVED***
	Err   error
***REMOVED***

// Error implements the error interface.
func (me MarshalError) Error() string ***REMOVED***
	return fmt.Sprintf("cannot transform type %s to a BSON Document: %v", reflect.TypeOf(me.Value), me.Err)
***REMOVED***

// Pipeline is a type that makes creating aggregation pipelines easier. It is a
// helper and is intended for serializing to BSON.
//
// Example usage:
//
//		mongo.Pipeline***REMOVED***
//			***REMOVED******REMOVED***"$group", bson.D***REMOVED******REMOVED***"_id", "$state"***REMOVED***, ***REMOVED***"totalPop", bson.D***REMOVED******REMOVED***"$sum", "$pop"***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***,
//			***REMOVED******REMOVED***"$match", bson.D***REMOVED******REMOVED***"totalPop", bson.D***REMOVED******REMOVED***"$gte", 10*1000*1000***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***,
//		***REMOVED***
//
type Pipeline []bson.D

// transformAndEnsureID is a hack that makes it easy to get a RawValue as the _id value.
// It will also add an ObjectID _id as the first key if it not already present in the passed-in val.
func transformAndEnsureID(registry *bsoncodec.Registry, val interface***REMOVED******REMOVED***) (bsoncore.Document, interface***REMOVED******REMOVED***, error) ***REMOVED***
	if registry == nil ***REMOVED***
		registry = bson.NewRegistryBuilder().Build()
	***REMOVED***
	switch tt := val.(type) ***REMOVED***
	case nil:
		return nil, nil, ErrNilDocument
	case bsonx.Doc:
		val = tt.Copy()
	case []byte:
		// Slight optimization so we'll just use MarshalBSON and not go through the codec machinery.
		val = bson.Raw(tt)
	***REMOVED***

	// TODO(skriptble): Use a pool of these instead.
	doc := make(bsoncore.Document, 0, 256)
	doc, err := bson.MarshalAppendWithRegistry(registry, doc, val)
	if err != nil ***REMOVED***
		return nil, nil, MarshalError***REMOVED***Value: val, Err: err***REMOVED***
	***REMOVED***

	var id interface***REMOVED******REMOVED***

	value := doc.Lookup("_id")
	switch value.Type ***REMOVED***
	case bsontype.Type(0):
		value = bsoncore.Value***REMOVED***Type: bsontype.ObjectID, Data: bsoncore.AppendObjectID(nil, primitive.NewObjectID())***REMOVED***
		olddoc := doc
		doc = make(bsoncore.Document, 0, len(olddoc)+17) // type byte + _id + null byte + object ID
		_, doc = bsoncore.ReserveLength(doc)
		doc = bsoncore.AppendValueElement(doc, "_id", value)
		doc = append(doc, olddoc[4:]...) // remove the length
		doc = bsoncore.UpdateLength(doc, 0, int32(len(doc)))
	default:
		// We copy the bytes here to ensure that any bytes returned to the user aren't modified
		// later.
		buf := make([]byte, len(value.Data))
		copy(buf, value.Data)
		value.Data = buf
	***REMOVED***

	err = bson.RawValue***REMOVED***Type: value.Type, Value: value.Data***REMOVED***.UnmarshalWithRegistry(registry, &id)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	return doc, id, nil
***REMOVED***

func transformBsoncoreDocument(registry *bsoncodec.Registry, val interface***REMOVED******REMOVED***, mapAllowed bool, paramName string) (bsoncore.Document, error) ***REMOVED***
	if registry == nil ***REMOVED***
		registry = bson.DefaultRegistry
	***REMOVED***
	if val == nil ***REMOVED***
		return nil, ErrNilDocument
	***REMOVED***
	if bs, ok := val.([]byte); ok ***REMOVED***
		// Slight optimization so we'll just use MarshalBSON and not go through the codec machinery.
		val = bson.Raw(bs)
	***REMOVED***
	if !mapAllowed ***REMOVED***
		refValue := reflect.ValueOf(val)
		if refValue.Kind() == reflect.Map && refValue.Len() > 1 ***REMOVED***
			return nil, ErrMapForOrderedArgument***REMOVED***paramName***REMOVED***
		***REMOVED***
	***REMOVED***

	// TODO(skriptble): Use a pool of these instead.
	buf := make([]byte, 0, 256)
	b, err := bson.MarshalAppendWithRegistry(registry, buf[:0], val)
	if err != nil ***REMOVED***
		return nil, MarshalError***REMOVED***Value: val, Err: err***REMOVED***
	***REMOVED***
	return b, nil
***REMOVED***

func ensureDollarKey(doc bsoncore.Document) error ***REMOVED***
	firstElem, err := doc.IndexErr(0)
	if err != nil ***REMOVED***
		return errors.New("update document must have at least one element")
	***REMOVED***

	if !strings.HasPrefix(firstElem.Key(), "$") ***REMOVED***
		return errors.New("update document must contain key beginning with '$'")
	***REMOVED***
	return nil
***REMOVED***

func ensureNoDollarKey(doc bsoncore.Document) error ***REMOVED***
	if elem, err := doc.IndexErr(0); err == nil && strings.HasPrefix(elem.Key(), "$") ***REMOVED***
		return errors.New("replacement document cannot contain keys beginning with '$'")
	***REMOVED***

	return nil
***REMOVED***

func transformAggregatePipeline(registry *bsoncodec.Registry, pipeline interface***REMOVED******REMOVED***) (bsoncore.Document, bool, error) ***REMOVED***
	switch t := pipeline.(type) ***REMOVED***
	case bsoncodec.ValueMarshaler:
		btype, val, err := t.MarshalBSONValue()
		if err != nil ***REMOVED***
			return nil, false, err
		***REMOVED***
		if btype != bsontype.Array ***REMOVED***
			return nil, false, fmt.Errorf("ValueMarshaler returned a %v, but was expecting %v", btype, bsontype.Array)
		***REMOVED***

		var hasOutputStage bool
		pipelineDoc := bsoncore.Document(val)
		values, _ := pipelineDoc.Values()
		if pipelineLen := len(values); pipelineLen > 0 ***REMOVED***
			if finalDoc, ok := values[pipelineLen-1].DocumentOK(); ok ***REMOVED***
				if elem, err := finalDoc.IndexErr(0); err == nil && (elem.Key() == "$out" || elem.Key() == "$merge") ***REMOVED***
					hasOutputStage = true
				***REMOVED***
			***REMOVED***
		***REMOVED***

		return pipelineDoc, hasOutputStage, nil
	default:
		val := reflect.ValueOf(t)
		if !val.IsValid() || (val.Kind() != reflect.Slice && val.Kind() != reflect.Array) ***REMOVED***
			return nil, false, fmt.Errorf("can only transform slices and arrays into aggregation pipelines, but got %v", val.Kind())
		***REMOVED***

		var hasOutputStage bool
		valLen := val.Len()

		switch t := pipeline.(type) ***REMOVED***
		// Explicitly forbid non-empty pipelines that are semantically single documents
		// and are implemented as slices.
		case bson.D, bson.Raw, bsoncore.Document:
			if valLen > 0 ***REMOVED***
				return nil, false,
					fmt.Errorf("%T is not an allowed pipeline type as it represents a single document. Use bson.A or mongo.Pipeline instead", t)
			***REMOVED***
		// bsoncore.Arrays do not need to be transformed. Only check validity and presence of output stage.
		case bsoncore.Array:
			if err := t.Validate(); err != nil ***REMOVED***
				return nil, false, err
			***REMOVED***

			values, err := t.Values()
			if err != nil ***REMOVED***
				return nil, false, err
			***REMOVED***

			numVals := len(values)
			if numVals == 0 ***REMOVED***
				return bsoncore.Document(t), false, nil
			***REMOVED***

			// If not empty, check if first value of the last stage is $out or $merge.
			if lastStage, ok := values[numVals-1].DocumentOK(); ok ***REMOVED***
				if elem, err := lastStage.IndexErr(0); err == nil && (elem.Key() == "$out" || elem.Key() == "$merge") ***REMOVED***
					hasOutputStage = true
				***REMOVED***
			***REMOVED***
			return bsoncore.Document(t), hasOutputStage, nil
		***REMOVED***

		aidx, arr := bsoncore.AppendArrayStart(nil)
		for idx := 0; idx < valLen; idx++ ***REMOVED***
			doc, err := transformBsoncoreDocument(registry, val.Index(idx).Interface(), true, fmt.Sprintf("pipeline stage :%v", idx))
			if err != nil ***REMOVED***
				return nil, false, err
			***REMOVED***

			if idx == valLen-1 ***REMOVED***
				if elem, err := doc.IndexErr(0); err == nil && (elem.Key() == "$out" || elem.Key() == "$merge") ***REMOVED***
					hasOutputStage = true
				***REMOVED***
			***REMOVED***
			arr = bsoncore.AppendDocumentElement(arr, strconv.Itoa(idx), doc)
		***REMOVED***
		arr, _ = bsoncore.AppendArrayEnd(arr, aidx)
		return arr, hasOutputStage, nil
	***REMOVED***
***REMOVED***

func transformUpdateValue(registry *bsoncodec.Registry, update interface***REMOVED******REMOVED***, dollarKeysAllowed bool) (bsoncore.Value, error) ***REMOVED***
	documentCheckerFunc := ensureDollarKey
	if !dollarKeysAllowed ***REMOVED***
		documentCheckerFunc = ensureNoDollarKey
	***REMOVED***

	var u bsoncore.Value
	var err error
	switch t := update.(type) ***REMOVED***
	case nil:
		return u, ErrNilDocument
	case primitive.D, bsonx.Doc:
		u.Type = bsontype.EmbeddedDocument
		u.Data, err = transformBsoncoreDocument(registry, update, true, "update")
		if err != nil ***REMOVED***
			return u, err
		***REMOVED***

		return u, documentCheckerFunc(u.Data)
	case bson.Raw:
		u.Type = bsontype.EmbeddedDocument
		u.Data = t
		return u, documentCheckerFunc(u.Data)
	case bsoncore.Document:
		u.Type = bsontype.EmbeddedDocument
		u.Data = t
		return u, documentCheckerFunc(u.Data)
	case []byte:
		u.Type = bsontype.EmbeddedDocument
		u.Data = t
		return u, documentCheckerFunc(u.Data)
	case bsoncodec.Marshaler:
		u.Type = bsontype.EmbeddedDocument
		u.Data, err = t.MarshalBSON()
		if err != nil ***REMOVED***
			return u, err
		***REMOVED***

		return u, documentCheckerFunc(u.Data)
	case bsoncodec.ValueMarshaler:
		u.Type, u.Data, err = t.MarshalBSONValue()
		if err != nil ***REMOVED***
			return u, err
		***REMOVED***
		if u.Type != bsontype.Array && u.Type != bsontype.EmbeddedDocument ***REMOVED***
			return u, fmt.Errorf("ValueMarshaler returned a %v, but was expecting %v or %v", u.Type, bsontype.Array, bsontype.EmbeddedDocument)
		***REMOVED***
		return u, err
	default:
		val := reflect.ValueOf(t)
		if !val.IsValid() ***REMOVED***
			return u, fmt.Errorf("can only transform slices and arrays into update pipelines, but got %v", val.Kind())
		***REMOVED***
		if val.Kind() != reflect.Slice && val.Kind() != reflect.Array ***REMOVED***
			u.Type = bsontype.EmbeddedDocument
			u.Data, err = transformBsoncoreDocument(registry, update, true, "update")
			if err != nil ***REMOVED***
				return u, err
			***REMOVED***

			return u, documentCheckerFunc(u.Data)
		***REMOVED***

		u.Type = bsontype.Array
		aidx, arr := bsoncore.AppendArrayStart(nil)
		valLen := val.Len()
		for idx := 0; idx < valLen; idx++ ***REMOVED***
			doc, err := transformBsoncoreDocument(registry, val.Index(idx).Interface(), true, "update")
			if err != nil ***REMOVED***
				return u, err
			***REMOVED***

			if err := documentCheckerFunc(doc); err != nil ***REMOVED***
				return u, err
			***REMOVED***

			arr = bsoncore.AppendDocumentElement(arr, strconv.Itoa(idx), doc)
		***REMOVED***
		u.Data, _ = bsoncore.AppendArrayEnd(arr, aidx)
		return u, err
	***REMOVED***
***REMOVED***

func transformValue(registry *bsoncodec.Registry, val interface***REMOVED******REMOVED***, mapAllowed bool, paramName string) (bsoncore.Value, error) ***REMOVED***
	if registry == nil ***REMOVED***
		registry = bson.DefaultRegistry
	***REMOVED***
	if val == nil ***REMOVED***
		return bsoncore.Value***REMOVED******REMOVED***, ErrNilValue
	***REMOVED***

	if !mapAllowed ***REMOVED***
		refValue := reflect.ValueOf(val)
		if refValue.Kind() == reflect.Map && refValue.Len() > 1 ***REMOVED***
			return bsoncore.Value***REMOVED******REMOVED***, ErrMapForOrderedArgument***REMOVED***paramName***REMOVED***
		***REMOVED***
	***REMOVED***

	buf := make([]byte, 0, 256)
	bsonType, bsonValue, err := bson.MarshalValueAppendWithRegistry(registry, buf[:0], val)
	if err != nil ***REMOVED***
		return bsoncore.Value***REMOVED******REMOVED***, MarshalError***REMOVED***Value: val, Err: err***REMOVED***
	***REMOVED***

	return bsoncore.Value***REMOVED***Type: bsonType, Data: bsonValue***REMOVED***, nil
***REMOVED***

// Build the aggregation pipeline for the CountDocument command.
func countDocumentsAggregatePipeline(registry *bsoncodec.Registry, filter interface***REMOVED******REMOVED***, opts *options.CountOptions) (bsoncore.Document, error) ***REMOVED***
	filterDoc, err := transformBsoncoreDocument(registry, filter, true, "filter")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	aidx, arr := bsoncore.AppendArrayStart(nil)
	didx, arr := bsoncore.AppendDocumentElementStart(arr, strconv.Itoa(0))
	arr = bsoncore.AppendDocumentElement(arr, "$match", filterDoc)
	arr, _ = bsoncore.AppendDocumentEnd(arr, didx)

	index := 1
	if opts != nil ***REMOVED***
		if opts.Skip != nil ***REMOVED***
			didx, arr = bsoncore.AppendDocumentElementStart(arr, strconv.Itoa(index))
			arr = bsoncore.AppendInt64Element(arr, "$skip", *opts.Skip)
			arr, _ = bsoncore.AppendDocumentEnd(arr, didx)
			index++
		***REMOVED***
		if opts.Limit != nil ***REMOVED***
			didx, arr = bsoncore.AppendDocumentElementStart(arr, strconv.Itoa(index))
			arr = bsoncore.AppendInt64Element(arr, "$limit", *opts.Limit)
			arr, _ = bsoncore.AppendDocumentEnd(arr, didx)
			index++
		***REMOVED***
	***REMOVED***

	didx, arr = bsoncore.AppendDocumentElementStart(arr, strconv.Itoa(index))
	iidx, arr := bsoncore.AppendDocumentElementStart(arr, "$group")
	arr = bsoncore.AppendInt32Element(arr, "_id", 1)
	iiidx, arr := bsoncore.AppendDocumentElementStart(arr, "n")
	arr = bsoncore.AppendInt32Element(arr, "$sum", 1)
	arr, _ = bsoncore.AppendDocumentEnd(arr, iiidx)
	arr, _ = bsoncore.AppendDocumentEnd(arr, iidx)
	arr, _ = bsoncore.AppendDocumentEnd(arr, didx)

	return bsoncore.AppendArrayEnd(arr, aidx)
***REMOVED***
