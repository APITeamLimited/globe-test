package echo

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type (
	// Binder is the interface that wraps the Bind method.
	Binder interface ***REMOVED***
		Bind(i interface***REMOVED******REMOVED***, c Context) error
	***REMOVED***

	// DefaultBinder is the default implementation of the Binder interface.
	DefaultBinder struct***REMOVED******REMOVED***

	// BindUnmarshaler is the interface used to wrap the UnmarshalParam method.
	BindUnmarshaler interface ***REMOVED***
		// UnmarshalParam decodes and assigns a value from an form or query param.
		UnmarshalParam(param string) error
	***REMOVED***
)

// Bind implements the `Binder#Bind` function.
func (b *DefaultBinder) Bind(i interface***REMOVED******REMOVED***, c Context) (err error) ***REMOVED***
	req := c.Request()
	if req.ContentLength == 0 ***REMOVED***
		if req.Method == GET || req.Method == DELETE ***REMOVED***
			if err = b.bindData(i, c.QueryParams(), "query"); err != nil ***REMOVED***
				return NewHTTPError(http.StatusBadRequest, err.Error())
			***REMOVED***
			return
		***REMOVED***
		return NewHTTPError(http.StatusBadRequest, "Request body can't be empty")
	***REMOVED***
	ctype := req.Header.Get(HeaderContentType)
	switch ***REMOVED***
	case strings.HasPrefix(ctype, MIMEApplicationJSON):
		if err = json.NewDecoder(req.Body).Decode(i); err != nil ***REMOVED***
			if ute, ok := err.(*json.UnmarshalTypeError); ok ***REMOVED***
				return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, offset=%v", ute.Type, ute.Value, ute.Offset))
			***REMOVED*** else if se, ok := err.(*json.SyntaxError); ok ***REMOVED***
				return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error()))
			***REMOVED*** else ***REMOVED***
				return NewHTTPError(http.StatusBadRequest, err.Error())
			***REMOVED***
		***REMOVED***
	case strings.HasPrefix(ctype, MIMEApplicationXML), strings.HasPrefix(ctype, MIMETextXML):
		if err = xml.NewDecoder(req.Body).Decode(i); err != nil ***REMOVED***
			if ute, ok := err.(*xml.UnsupportedTypeError); ok ***REMOVED***
				return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unsupported type error: type=%v, error=%v", ute.Type, ute.Error()))
			***REMOVED*** else if se, ok := err.(*xml.SyntaxError); ok ***REMOVED***
				return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: line=%v, error=%v", se.Line, se.Error()))
			***REMOVED*** else ***REMOVED***
				return NewHTTPError(http.StatusBadRequest, err.Error())
			***REMOVED***
		***REMOVED***
	case strings.HasPrefix(ctype, MIMEApplicationForm), strings.HasPrefix(ctype, MIMEMultipartForm):
		params, err := c.FormParams()
		if err != nil ***REMOVED***
			return NewHTTPError(http.StatusBadRequest, err.Error())
		***REMOVED***
		if err = b.bindData(i, params, "form"); err != nil ***REMOVED***
			return NewHTTPError(http.StatusBadRequest, err.Error())
		***REMOVED***
	default:
		return ErrUnsupportedMediaType
	***REMOVED***
	return
***REMOVED***

func (b *DefaultBinder) bindData(ptr interface***REMOVED******REMOVED***, data map[string][]string, tag string) error ***REMOVED***
	typ := reflect.TypeOf(ptr).Elem()
	val := reflect.ValueOf(ptr).Elem()

	if typ.Kind() != reflect.Struct ***REMOVED***
		return errors.New("Binding element must be a struct")
	***REMOVED***

	for i := 0; i < typ.NumField(); i++ ***REMOVED***
		typeField := typ.Field(i)
		structField := val.Field(i)
		if !structField.CanSet() ***REMOVED***
			continue
		***REMOVED***
		structFieldKind := structField.Kind()
		inputFieldName := typeField.Tag.Get(tag)

		if inputFieldName == "" ***REMOVED***
			inputFieldName = typeField.Name
			// If tag is nil, we inspect if the field is a struct.
			if _, ok := bindUnmarshaler(structField); !ok && structFieldKind == reflect.Struct ***REMOVED***
				err := b.bindData(structField.Addr().Interface(), data, tag)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
		inputValue, exists := data[inputFieldName]
		if !exists ***REMOVED***
			continue
		***REMOVED***

		// Call this first, in case we're dealing with an alias to an array type
		if ok, err := unmarshalField(typeField.Type.Kind(), inputValue[0], structField); ok ***REMOVED***
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			continue
		***REMOVED***

		numElems := len(inputValue)
		if structFieldKind == reflect.Slice && numElems > 0 ***REMOVED***
			sliceOf := structField.Type().Elem().Kind()
			slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
			for j := 0; j < numElems; j++ ***REMOVED***
				if err := setWithProperType(sliceOf, inputValue[j], slice.Index(j)); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
			val.Field(i).Set(slice)
		***REMOVED*** else ***REMOVED***
			if err := setWithProperType(typeField.Type.Kind(), inputValue[0], structField); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error ***REMOVED***
	// But also call it here, in case we're dealing with an array of BindUnmarshalers
	if ok, err := unmarshalField(valueKind, val, structField); ok ***REMOVED***
		return err
	***REMOVED***

	switch valueKind ***REMOVED***
	case reflect.Ptr:
		return setWithProperType(structField.Elem().Kind(), val, structField.Elem())
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
		return errors.New("unknown type")
	***REMOVED***
	return nil
***REMOVED***

func unmarshalField(valueKind reflect.Kind, val string, field reflect.Value) (bool, error) ***REMOVED***
	switch valueKind ***REMOVED***
	case reflect.Ptr:
		return unmarshalFieldPtr(val, field)
	default:
		return unmarshalFieldNonPtr(val, field)
	***REMOVED***
***REMOVED***

// bindUnmarshaler attempts to unmarshal a reflect.Value into a BindUnmarshaler
func bindUnmarshaler(field reflect.Value) (BindUnmarshaler, bool) ***REMOVED***
	ptr := reflect.New(field.Type())
	if ptr.CanInterface() ***REMOVED***
		iface := ptr.Interface()
		if unmarshaler, ok := iface.(BindUnmarshaler); ok ***REMOVED***
			return unmarshaler, ok
		***REMOVED***
	***REMOVED***
	return nil, false
***REMOVED***

func unmarshalFieldNonPtr(value string, field reflect.Value) (bool, error) ***REMOVED***
	if unmarshaler, ok := bindUnmarshaler(field); ok ***REMOVED***
		err := unmarshaler.UnmarshalParam(value)
		field.Set(reflect.ValueOf(unmarshaler).Elem())
		return true, err
	***REMOVED***
	return false, nil
***REMOVED***

func unmarshalFieldPtr(value string, field reflect.Value) (bool, error) ***REMOVED***
	if field.IsNil() ***REMOVED***
		// Initialize the pointer to a nil value
		field.Set(reflect.New(field.Type().Elem()))
	***REMOVED***
	return unmarshalFieldNonPtr(value, field.Elem())
***REMOVED***

func setIntField(value string, bitSize int, field reflect.Value) error ***REMOVED***
	if value == "" ***REMOVED***
		value = "0"
	***REMOVED***
	intVal, err := strconv.ParseInt(value, 10, bitSize)
	if err == nil ***REMOVED***
		field.SetInt(intVal)
	***REMOVED***
	return err
***REMOVED***

func setUintField(value string, bitSize int, field reflect.Value) error ***REMOVED***
	if value == "" ***REMOVED***
		value = "0"
	***REMOVED***
	uintVal, err := strconv.ParseUint(value, 10, bitSize)
	if err == nil ***REMOVED***
		field.SetUint(uintVal)
	***REMOVED***
	return err
***REMOVED***

func setBoolField(value string, field reflect.Value) error ***REMOVED***
	if value == "" ***REMOVED***
		value = "false"
	***REMOVED***
	boolVal, err := strconv.ParseBool(value)
	if err == nil ***REMOVED***
		field.SetBool(boolVal)
	***REMOVED***
	return err
***REMOVED***

func setFloatField(value string, bitSize int, field reflect.Value) error ***REMOVED***
	if value == "" ***REMOVED***
		value = "0.0"
	***REMOVED***
	floatVal, err := strconv.ParseFloat(value, bitSize)
	if err == nil ***REMOVED***
		field.SetFloat(floatVal)
	***REMOVED***
	return err
***REMOVED***
