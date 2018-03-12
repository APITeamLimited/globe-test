// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

func mapForm(ptr interface***REMOVED******REMOVED***, form map[string][]string) error ***REMOVED***
	typ := reflect.TypeOf(ptr).Elem()
	val := reflect.ValueOf(ptr).Elem()
	for i := 0; i < typ.NumField(); i++ ***REMOVED***
		typeField := typ.Field(i)
		structField := val.Field(i)
		if !structField.CanSet() ***REMOVED***
			continue
		***REMOVED***

		structFieldKind := structField.Kind()
		inputFieldName := typeField.Tag.Get("form")
		if inputFieldName == "" ***REMOVED***
			inputFieldName = typeField.Name

			// if "form" tag is nil, we inspect if the field is a struct.
			// this would not make sense for JSON parsing but it does for a form
			// since data is flatten
			if structFieldKind == reflect.Struct ***REMOVED***
				err := mapForm(structField.Addr().Interface(), form)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
		inputValue, exists := form[inputFieldName]
		if !exists ***REMOVED***
			continue
		***REMOVED***

		numElems := len(inputValue)
		if structFieldKind == reflect.Slice && numElems > 0 ***REMOVED***
			sliceOf := structField.Type().Elem().Kind()
			slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
			for i := 0; i < numElems; i++ ***REMOVED***
				if err := setWithProperType(sliceOf, inputValue[i], slice.Index(i)); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
			val.Field(i).Set(slice)
		***REMOVED*** else ***REMOVED***
			if _, isTime := structField.Interface().(time.Time); isTime ***REMOVED***
				if err := setTimeField(inputValue[0], typeField, structField); err != nil ***REMOVED***
					return err
				***REMOVED***
				continue
			***REMOVED***
			if err := setWithProperType(typeField.Type.Kind(), inputValue[0], structField); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error ***REMOVED***
	switch valueKind ***REMOVED***
	case reflect.Int:
		return setIntField(val, 0, structField)
	case reflect.Int8:
		return setIntField(val, 8, structField)
	case reflect.Int16:
		return setIntField(val, 16, structField)
	case reflect.Int32:
		return setIntField(val, 32, structField)
	case reflect.Int64:
		return setIntField(val, 64, structField)
	case reflect.Uint:
		return setUintField(val, 0, structField)
	case reflect.Uint8:
		return setUintField(val, 8, structField)
	case reflect.Uint16:
		return setUintField(val, 16, structField)
	case reflect.Uint32:
		return setUintField(val, 32, structField)
	case reflect.Uint64:
		return setUintField(val, 64, structField)
	case reflect.Bool:
		return setBoolField(val, structField)
	case reflect.Float32:
		return setFloatField(val, 32, structField)
	case reflect.Float64:
		return setFloatField(val, 64, structField)
	case reflect.String:
		structField.SetString(val)
	default:
		return errors.New("Unknown type")
	***REMOVED***
	return nil
***REMOVED***

func setIntField(val string, bitSize int, field reflect.Value) error ***REMOVED***
	if val == "" ***REMOVED***
		val = "0"
	***REMOVED***
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil ***REMOVED***
		field.SetInt(intVal)
	***REMOVED***
	return err
***REMOVED***

func setUintField(val string, bitSize int, field reflect.Value) error ***REMOVED***
	if val == "" ***REMOVED***
		val = "0"
	***REMOVED***
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil ***REMOVED***
		field.SetUint(uintVal)
	***REMOVED***
	return err
***REMOVED***

func setBoolField(val string, field reflect.Value) error ***REMOVED***
	if val == "" ***REMOVED***
		val = "false"
	***REMOVED***
	boolVal, err := strconv.ParseBool(val)
	if err == nil ***REMOVED***
		field.SetBool(boolVal)
	***REMOVED***
	return nil
***REMOVED***

func setFloatField(val string, bitSize int, field reflect.Value) error ***REMOVED***
	if val == "" ***REMOVED***
		val = "0.0"
	***REMOVED***
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil ***REMOVED***
		field.SetFloat(floatVal)
	***REMOVED***
	return err
***REMOVED***

func setTimeField(val string, structField reflect.StructField, value reflect.Value) error ***REMOVED***
	timeFormat := structField.Tag.Get("time_format")
	if timeFormat == "" ***REMOVED***
		return errors.New("Blank time format")
	***REMOVED***

	if val == "" ***REMOVED***
		value.Set(reflect.ValueOf(time.Time***REMOVED******REMOVED***))
		return nil
	***REMOVED***

	l := time.Local
	if isUTC, _ := strconv.ParseBool(structField.Tag.Get("time_utc")); isUTC ***REMOVED***
		l = time.UTC
	***REMOVED***

	t, err := time.ParseInLocation(timeFormat, val, l)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	value.Set(reflect.ValueOf(t))
	return nil
***REMOVED***

// Don't pass in pointers to bind to. Can lead to bugs. See:
// https://github.com/codegangsta/martini-contrib/issues/40
// https://github.com/codegangsta/martini-contrib/pull/34#issuecomment-29683659
func ensureNotPointer(obj interface***REMOVED******REMOVED***) ***REMOVED***
	if reflect.TypeOf(obj).Kind() == reflect.Ptr ***REMOVED***
		panic("Pointers are not accepted as binding models")
	***REMOVED***
***REMOVED***
