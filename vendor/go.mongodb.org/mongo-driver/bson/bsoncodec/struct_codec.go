// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsoncodec

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsonoptions"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// DecodeError represents an error that occurs when unmarshalling BSON bytes into a native Go type.
type DecodeError struct ***REMOVED***
	keys    []string
	wrapped error
***REMOVED***

// Unwrap returns the underlying error
func (de *DecodeError) Unwrap() error ***REMOVED***
	return de.wrapped
***REMOVED***

// Error implements the error interface.
func (de *DecodeError) Error() string ***REMOVED***
	// The keys are stored in reverse order because the de.keys slice is builtup while propagating the error up the
	// stack of BSON keys, so we call de.Keys(), which reverses them.
	keyPath := strings.Join(de.Keys(), ".")
	return fmt.Sprintf("error decoding key %s: %v", keyPath, de.wrapped)
***REMOVED***

// Keys returns the BSON key path that caused an error as a slice of strings. The keys in the slice are in top-down
// order. For example, if the document being unmarshalled was ***REMOVED***a: ***REMOVED***b: ***REMOVED***c: 1***REMOVED******REMOVED******REMOVED*** and the value for c was supposed to be
// a string, the keys slice will be ["a", "b", "c"].
func (de *DecodeError) Keys() []string ***REMOVED***
	reversedKeys := make([]string, 0, len(de.keys))
	for idx := len(de.keys) - 1; idx >= 0; idx-- ***REMOVED***
		reversedKeys = append(reversedKeys, de.keys[idx])
	***REMOVED***

	return reversedKeys
***REMOVED***

// Zeroer allows custom struct types to implement a report of zero
// state. All struct types that don't implement Zeroer or where IsZero
// returns false are considered to be not zero.
type Zeroer interface ***REMOVED***
	IsZero() bool
***REMOVED***

// StructCodec is the Codec used for struct values.
type StructCodec struct ***REMOVED***
	cache                            map[reflect.Type]*structDescription
	l                                sync.RWMutex
	parser                           StructTagParser
	DecodeZeroStruct                 bool
	DecodeDeepZeroInline             bool
	EncodeOmitDefaultStruct          bool
	AllowUnexportedFields            bool
	OverwriteDuplicatedInlinedFields bool
***REMOVED***

var _ ValueEncoder = &StructCodec***REMOVED******REMOVED***
var _ ValueDecoder = &StructCodec***REMOVED******REMOVED***

// NewStructCodec returns a StructCodec that uses p for struct tag parsing.
func NewStructCodec(p StructTagParser, opts ...*bsonoptions.StructCodecOptions) (*StructCodec, error) ***REMOVED***
	if p == nil ***REMOVED***
		return nil, errors.New("a StructTagParser must be provided to NewStructCodec")
	***REMOVED***

	structOpt := bsonoptions.MergeStructCodecOptions(opts...)

	codec := &StructCodec***REMOVED***
		cache:  make(map[reflect.Type]*structDescription),
		parser: p,
	***REMOVED***

	if structOpt.DecodeZeroStruct != nil ***REMOVED***
		codec.DecodeZeroStruct = *structOpt.DecodeZeroStruct
	***REMOVED***
	if structOpt.DecodeDeepZeroInline != nil ***REMOVED***
		codec.DecodeDeepZeroInline = *structOpt.DecodeDeepZeroInline
	***REMOVED***
	if structOpt.EncodeOmitDefaultStruct != nil ***REMOVED***
		codec.EncodeOmitDefaultStruct = *structOpt.EncodeOmitDefaultStruct
	***REMOVED***
	if structOpt.OverwriteDuplicatedInlinedFields != nil ***REMOVED***
		codec.OverwriteDuplicatedInlinedFields = *structOpt.OverwriteDuplicatedInlinedFields
	***REMOVED***
	if structOpt.AllowUnexportedFields != nil ***REMOVED***
		codec.AllowUnexportedFields = *structOpt.AllowUnexportedFields
	***REMOVED***

	return codec, nil
***REMOVED***

// EncodeValue handles encoding generic struct types.
func (sc *StructCodec) EncodeValue(r EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Kind() != reflect.Struct ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "StructCodec.EncodeValue", Kinds: []reflect.Kind***REMOVED***reflect.Struct***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	sd, err := sc.describeStruct(r.Registry, val.Type())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	dw, err := vw.WriteDocument()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	var rv reflect.Value
	for _, desc := range sd.fl ***REMOVED***
		if desc.inline == nil ***REMOVED***
			rv = val.Field(desc.idx)
		***REMOVED*** else ***REMOVED***
			rv, err = fieldByIndexErr(val, desc.inline)
			if err != nil ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***

		desc.encoder, rv, err = defaultValueEncoders.lookupElementEncoder(r, desc.encoder, rv)

		if err != nil && err != errInvalidValue ***REMOVED***
			return err
		***REMOVED***

		if err == errInvalidValue ***REMOVED***
			if desc.omitEmpty ***REMOVED***
				continue
			***REMOVED***
			vw2, err := dw.WriteDocumentElement(desc.name)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			err = vw2.WriteNull()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			continue
		***REMOVED***

		if desc.encoder == nil ***REMOVED***
			return ErrNoEncoder***REMOVED***Type: rv.Type()***REMOVED***
		***REMOVED***

		encoder := desc.encoder

		var isZero bool
		rvInterface := rv.Interface()
		if cz, ok := encoder.(CodecZeroer); ok ***REMOVED***
			isZero = cz.IsTypeZero(rvInterface)
		***REMOVED*** else if rv.Kind() == reflect.Interface ***REMOVED***
			// sc.isZero will not treat an interface rv as an interface, so we need to check for the zero interface separately.
			isZero = rv.IsNil()
		***REMOVED*** else ***REMOVED***
			isZero = sc.isZero(rvInterface)
		***REMOVED***
		if desc.omitEmpty && isZero ***REMOVED***
			continue
		***REMOVED***

		vw2, err := dw.WriteDocumentElement(desc.name)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		ectx := EncodeContext***REMOVED***Registry: r.Registry, MinSize: desc.minSize***REMOVED***
		err = encoder.EncodeValue(ectx, vw2, rv)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if sd.inlineMap >= 0 ***REMOVED***
		rv := val.Field(sd.inlineMap)
		collisionFn := func(key string) bool ***REMOVED***
			_, exists := sd.fm[key]
			return exists
		***REMOVED***

		return defaultMapCodec.mapEncodeValue(r, dw, rv, collisionFn)
	***REMOVED***

	return dw.WriteDocumentEnd()
***REMOVED***

func newDecodeError(key string, original error) error ***REMOVED***
	de, ok := original.(*DecodeError)
	if !ok ***REMOVED***
		return &DecodeError***REMOVED***
			keys:    []string***REMOVED***key***REMOVED***,
			wrapped: original,
		***REMOVED***
	***REMOVED***

	de.keys = append(de.keys, key)
	return de
***REMOVED***

// DecodeValue implements the Codec interface.
// By default, map types in val will not be cleared. If a map has existing key/value pairs, it will be extended with the new ones from vr.
// For slices, the decoder will set the length of the slice to zero and append all elements. The underlying array will not be cleared.
func (sc *StructCodec) DecodeValue(r DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Kind() != reflect.Struct ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "StructCodec.DecodeValue", Kinds: []reflect.Kind***REMOVED***reflect.Struct***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.Type(0), bsontype.EmbeddedDocument:
	case bsontype.Null:
		if err := vr.ReadNull(); err != nil ***REMOVED***
			return err
		***REMOVED***

		val.Set(reflect.Zero(val.Type()))
		return nil
	case bsontype.Undefined:
		if err := vr.ReadUndefined(); err != nil ***REMOVED***
			return err
		***REMOVED***

		val.Set(reflect.Zero(val.Type()))
		return nil
	default:
		return fmt.Errorf("cannot decode %v into a %s", vrType, val.Type())
	***REMOVED***

	sd, err := sc.describeStruct(r.Registry, val.Type())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if sc.DecodeZeroStruct ***REMOVED***
		val.Set(reflect.Zero(val.Type()))
	***REMOVED***
	if sc.DecodeDeepZeroInline && sd.inline ***REMOVED***
		val.Set(deepZero(val.Type()))
	***REMOVED***

	var decoder ValueDecoder
	var inlineMap reflect.Value
	if sd.inlineMap >= 0 ***REMOVED***
		inlineMap = val.Field(sd.inlineMap)
		decoder, err = r.LookupDecoder(inlineMap.Type().Elem())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	dr, err := vr.ReadDocument()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for ***REMOVED***
		name, vr, err := dr.ReadElement()
		if err == bsonrw.ErrEOD ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		fd, exists := sd.fm[name]
		if !exists ***REMOVED***
			// if the original name isn't found in the struct description, try again with the name in lowercase
			// this could match if a BSON tag isn't specified because by default, describeStruct lowercases all field
			// names
			fd, exists = sd.fm[strings.ToLower(name)]
		***REMOVED***

		if !exists ***REMOVED***
			if sd.inlineMap < 0 ***REMOVED***
				// The encoding/json package requires a flag to return on error for non-existent fields.
				// This functionality seems appropriate for the struct codec.
				err = vr.Skip()
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				continue
			***REMOVED***

			if inlineMap.IsNil() ***REMOVED***
				inlineMap.Set(reflect.MakeMap(inlineMap.Type()))
			***REMOVED***

			elem := reflect.New(inlineMap.Type().Elem()).Elem()
			r.Ancestor = inlineMap.Type()
			err = decoder.DecodeValue(r, vr, elem)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			inlineMap.SetMapIndex(reflect.ValueOf(name), elem)
			continue
		***REMOVED***

		var field reflect.Value
		if fd.inline == nil ***REMOVED***
			field = val.Field(fd.idx)
		***REMOVED*** else ***REMOVED***
			field, err = getInlineField(val, fd.inline)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		if !field.CanSet() ***REMOVED*** // Being settable is a super set of being addressable.
			innerErr := fmt.Errorf("field %v is not settable", field)
			return newDecodeError(fd.name, innerErr)
		***REMOVED***
		if field.Kind() == reflect.Ptr && field.IsNil() ***REMOVED***
			field.Set(reflect.New(field.Type().Elem()))
		***REMOVED***
		field = field.Addr()

		dctx := DecodeContext***REMOVED***Registry: r.Registry, Truncate: fd.truncate || r.Truncate***REMOVED***
		if fd.decoder == nil ***REMOVED***
			return newDecodeError(fd.name, ErrNoDecoder***REMOVED***Type: field.Elem().Type()***REMOVED***)
		***REMOVED***

		err = fd.decoder.DecodeValue(dctx, vr, field.Elem())
		if err != nil ***REMOVED***
			return newDecodeError(fd.name, err)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (sc *StructCodec) isZero(i interface***REMOVED******REMOVED***) bool ***REMOVED***
	v := reflect.ValueOf(i)

	// check the value validity
	if !v.IsValid() ***REMOVED***
		return true
	***REMOVED***

	if z, ok := v.Interface().(Zeroer); ok && (v.Kind() != reflect.Ptr || !v.IsNil()) ***REMOVED***
		return z.IsZero()
	***REMOVED***

	switch v.Kind() ***REMOVED***
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Struct:
		if sc.EncodeOmitDefaultStruct ***REMOVED***
			vt := v.Type()
			if vt == tTime ***REMOVED***
				return v.Interface().(time.Time).IsZero()
			***REMOVED***
			for i := 0; i < v.NumField(); i++ ***REMOVED***
				if vt.Field(i).PkgPath != "" && !vt.Field(i).Anonymous ***REMOVED***
					continue // Private field
				***REMOVED***
				fld := v.Field(i)
				if !sc.isZero(fld.Interface()) ***REMOVED***
					return false
				***REMOVED***
			***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

type structDescription struct ***REMOVED***
	fm        map[string]fieldDescription
	fl        []fieldDescription
	inlineMap int
	inline    bool
***REMOVED***

type fieldDescription struct ***REMOVED***
	name      string // BSON key name
	fieldName string // struct field name
	idx       int
	omitEmpty bool
	minSize   bool
	truncate  bool
	inline    []int
	encoder   ValueEncoder
	decoder   ValueDecoder
***REMOVED***

type byIndex []fieldDescription

func (bi byIndex) Len() int ***REMOVED*** return len(bi) ***REMOVED***

func (bi byIndex) Swap(i, j int) ***REMOVED*** bi[i], bi[j] = bi[j], bi[i] ***REMOVED***

func (bi byIndex) Less(i, j int) bool ***REMOVED***
	// If a field is inlined, its index in the top level struct is stored at inline[0]
	iIdx, jIdx := bi[i].idx, bi[j].idx
	if len(bi[i].inline) > 0 ***REMOVED***
		iIdx = bi[i].inline[0]
	***REMOVED***
	if len(bi[j].inline) > 0 ***REMOVED***
		jIdx = bi[j].inline[0]
	***REMOVED***
	if iIdx != jIdx ***REMOVED***
		return iIdx < jIdx
	***REMOVED***
	for k, biik := range bi[i].inline ***REMOVED***
		if k >= len(bi[j].inline) ***REMOVED***
			return false
		***REMOVED***
		if biik != bi[j].inline[k] ***REMOVED***
			return biik < bi[j].inline[k]
		***REMOVED***
	***REMOVED***
	return len(bi[i].inline) < len(bi[j].inline)
***REMOVED***

func (sc *StructCodec) describeStruct(r *Registry, t reflect.Type) (*structDescription, error) ***REMOVED***
	// We need to analyze the struct, including getting the tags, collecting
	// information about inlining, and create a map of the field name to the field.
	sc.l.RLock()
	ds, exists := sc.cache[t]
	sc.l.RUnlock()
	if exists ***REMOVED***
		return ds, nil
	***REMOVED***

	numFields := t.NumField()
	sd := &structDescription***REMOVED***
		fm:        make(map[string]fieldDescription, numFields),
		fl:        make([]fieldDescription, 0, numFields),
		inlineMap: -1,
	***REMOVED***

	var fields []fieldDescription
	for i := 0; i < numFields; i++ ***REMOVED***
		sf := t.Field(i)
		if sf.PkgPath != "" && (!sc.AllowUnexportedFields || !sf.Anonymous) ***REMOVED***
			// field is private or unexported fields aren't allowed, ignore
			continue
		***REMOVED***

		sfType := sf.Type
		encoder, err := r.LookupEncoder(sfType)
		if err != nil ***REMOVED***
			encoder = nil
		***REMOVED***
		decoder, err := r.LookupDecoder(sfType)
		if err != nil ***REMOVED***
			decoder = nil
		***REMOVED***

		description := fieldDescription***REMOVED***
			fieldName: sf.Name,
			idx:       i,
			encoder:   encoder,
			decoder:   decoder,
		***REMOVED***

		stags, err := sc.parser.ParseStructTags(sf)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if stags.Skip ***REMOVED***
			continue
		***REMOVED***
		description.name = stags.Name
		description.omitEmpty = stags.OmitEmpty
		description.minSize = stags.MinSize
		description.truncate = stags.Truncate

		if stags.Inline ***REMOVED***
			sd.inline = true
			switch sfType.Kind() ***REMOVED***
			case reflect.Map:
				if sd.inlineMap >= 0 ***REMOVED***
					return nil, errors.New("(struct " + t.String() + ") multiple inline maps")
				***REMOVED***
				if sfType.Key() != tString ***REMOVED***
					return nil, errors.New("(struct " + t.String() + ") inline map must have a string keys")
				***REMOVED***
				sd.inlineMap = description.idx
			case reflect.Ptr:
				sfType = sfType.Elem()
				if sfType.Kind() != reflect.Struct ***REMOVED***
					return nil, fmt.Errorf("(struct %s) inline fields must be a struct, a struct pointer, or a map", t.String())
				***REMOVED***
				fallthrough
			case reflect.Struct:
				inlinesf, err := sc.describeStruct(r, sfType)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				for _, fd := range inlinesf.fl ***REMOVED***
					if fd.inline == nil ***REMOVED***
						fd.inline = []int***REMOVED***i, fd.idx***REMOVED***
					***REMOVED*** else ***REMOVED***
						fd.inline = append([]int***REMOVED***i***REMOVED***, fd.inline...)
					***REMOVED***
					fields = append(fields, fd)

				***REMOVED***
			default:
				return nil, fmt.Errorf("(struct %s) inline fields must be a struct, a struct pointer, or a map", t.String())
			***REMOVED***
			continue
		***REMOVED***
		fields = append(fields, description)
	***REMOVED***

	// Sort fieldDescriptions by name and use dominance rules to determine which should be added for each name
	sort.Slice(fields, func(i, j int) bool ***REMOVED***
		x := fields
		// sort field by name, breaking ties with depth, then
		// breaking ties with index sequence.
		if x[i].name != x[j].name ***REMOVED***
			return x[i].name < x[j].name
		***REMOVED***
		if len(x[i].inline) != len(x[j].inline) ***REMOVED***
			return len(x[i].inline) < len(x[j].inline)
		***REMOVED***
		return byIndex(x).Less(i, j)
	***REMOVED***)

	for advance, i := 0, 0; i < len(fields); i += advance ***REMOVED***
		// One iteration per name.
		// Find the sequence of fields with the name of this first field.
		fi := fields[i]
		name := fi.name
		for advance = 1; i+advance < len(fields); advance++ ***REMOVED***
			fj := fields[i+advance]
			if fj.name != name ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		if advance == 1 ***REMOVED*** // Only one field with this name
			sd.fl = append(sd.fl, fi)
			sd.fm[name] = fi
			continue
		***REMOVED***
		dominant, ok := dominantField(fields[i : i+advance])
		if !ok || !sc.OverwriteDuplicatedInlinedFields ***REMOVED***
			return nil, fmt.Errorf("struct %s has duplicated key %s", t.String(), name)
		***REMOVED***
		sd.fl = append(sd.fl, dominant)
		sd.fm[name] = dominant
	***REMOVED***

	sort.Sort(byIndex(sd.fl))

	sc.l.Lock()
	sc.cache[t] = sd
	sc.l.Unlock()

	return sd, nil
***REMOVED***

// dominantField looks through the fields, all of which are known to
// have the same name, to find the single field that dominates the
// others using Go's inlining rules. If there are multiple top-level
// fields, the boolean will be false: This condition is an error in Go
// and we skip all the fields.
func dominantField(fields []fieldDescription) (fieldDescription, bool) ***REMOVED***
	// The fields are sorted in increasing index-length order, then by presence of tag.
	// That means that the first field is the dominant one. We need only check
	// for error cases: two fields at top level.
	if len(fields) > 1 &&
		len(fields[0].inline) == len(fields[1].inline) ***REMOVED***
		return fieldDescription***REMOVED******REMOVED***, false
	***REMOVED***
	return fields[0], true
***REMOVED***

func fieldByIndexErr(v reflect.Value, index []int) (result reflect.Value, err error) ***REMOVED***
	defer func() ***REMOVED***
		if recovered := recover(); recovered != nil ***REMOVED***
			switch r := recovered.(type) ***REMOVED***
			case string:
				err = fmt.Errorf("%s", r)
			case error:
				err = r
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	result = v.FieldByIndex(index)
	return
***REMOVED***

func getInlineField(val reflect.Value, index []int) (reflect.Value, error) ***REMOVED***
	field, err := fieldByIndexErr(val, index)
	if err == nil ***REMOVED***
		return field, nil
	***REMOVED***

	// if parent of this element doesn't exist, fix its parent
	inlineParent := index[:len(index)-1]
	var fParent reflect.Value
	if fParent, err = fieldByIndexErr(val, inlineParent); err != nil ***REMOVED***
		fParent, err = getInlineField(val, inlineParent)
		if err != nil ***REMOVED***
			return fParent, err
		***REMOVED***
	***REMOVED***
	fParent.Set(reflect.New(fParent.Type().Elem()))

	return fieldByIndexErr(val, index)
***REMOVED***

// DeepZero returns recursive zero object
func deepZero(st reflect.Type) (result reflect.Value) ***REMOVED***
	result = reflect.Indirect(reflect.New(st))

	if result.Kind() == reflect.Struct ***REMOVED***
		for i := 0; i < result.NumField(); i++ ***REMOVED***
			if f := result.Field(i); f.Kind() == reflect.Ptr ***REMOVED***
				if f.CanInterface() ***REMOVED***
					if ft := reflect.TypeOf(f.Interface()); ft.Elem().Kind() == reflect.Struct ***REMOVED***
						result.Field(i).Set(recursivePointerTo(deepZero(ft.Elem())))
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

// recursivePointerTo calls reflect.New(v.Type) but recursively for its fields inside
func recursivePointerTo(v reflect.Value) reflect.Value ***REMOVED***
	v = reflect.Indirect(v)
	result := reflect.New(v.Type())
	if v.Kind() == reflect.Struct ***REMOVED***
		for i := 0; i < v.NumField(); i++ ***REMOVED***
			if f := v.Field(i); f.Kind() == reflect.Ptr ***REMOVED***
				if f.Elem().Kind() == reflect.Struct ***REMOVED***
					result.Elem().Field(i).Set(recursivePointerTo(f))
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return result
***REMOVED***
