package mapstructure

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// typedDecodeHook takes a raw DecodeHookFunc (an interface***REMOVED******REMOVED***) and turns
// it into the proper DecodeHookFunc type, such as DecodeHookFuncType.
func typedDecodeHook(h DecodeHookFunc) DecodeHookFunc ***REMOVED***
	// Create variables here so we can reference them with the reflect pkg
	var f1 DecodeHookFuncType
	var f2 DecodeHookFuncKind

	// Fill in the variables into this interface and the rest is done
	// automatically using the reflect package.
	potential := []interface***REMOVED******REMOVED******REMOVED***f1, f2***REMOVED***

	v := reflect.ValueOf(h)
	vt := v.Type()
	for _, raw := range potential ***REMOVED***
		pt := reflect.ValueOf(raw).Type()
		if vt.ConvertibleTo(pt) ***REMOVED***
			return v.Convert(pt).Interface()
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// DecodeHookExec executes the given decode hook. This should be used
// since it'll naturally degrade to the older backwards compatible DecodeHookFunc
// that took reflect.Kind instead of reflect.Type.
func DecodeHookExec(
	raw DecodeHookFunc,
	from reflect.Type, to reflect.Type,
	data interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	switch f := typedDecodeHook(raw).(type) ***REMOVED***
	case DecodeHookFuncType:
		return f(from, to, data)
	case DecodeHookFuncKind:
		return f(from.Kind(), to.Kind(), data)
	default:
		return nil, errors.New("invalid decode hook signature")
	***REMOVED***
***REMOVED***

// ComposeDecodeHookFunc creates a single DecodeHookFunc that
// automatically composes multiple DecodeHookFuncs.
//
// The composed funcs are called in order, with the result of the
// previous transformation.
func ComposeDecodeHookFunc(fs ...DecodeHookFunc) DecodeHookFunc ***REMOVED***
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
		var err error
		for _, f1 := range fs ***REMOVED***
			data, err = DecodeHookExec(f1, f, t, data)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			// Modify the from kind to be correct with the new data
			f = nil
			if val := reflect.ValueOf(data); val.IsValid() ***REMOVED***
				f = val.Type()
			***REMOVED***
		***REMOVED***

		return data, nil
	***REMOVED***
***REMOVED***

// StringToSliceHookFunc returns a DecodeHookFunc that converts
// string to []string by splitting on the given sep.
func StringToSliceHookFunc(sep string) DecodeHookFunc ***REMOVED***
	return func(
		f reflect.Kind,
		t reflect.Kind,
		data interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
		if f != reflect.String || t != reflect.Slice ***REMOVED***
			return data, nil
		***REMOVED***

		raw := data.(string)
		if raw == "" ***REMOVED***
			return []string***REMOVED******REMOVED***, nil
		***REMOVED***

		return strings.Split(raw, sep), nil
	***REMOVED***
***REMOVED***

// StringToTimeDurationHookFunc returns a DecodeHookFunc that converts
// strings to time.Duration.
func StringToTimeDurationHookFunc() DecodeHookFunc ***REMOVED***
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
		if f.Kind() != reflect.String ***REMOVED***
			return data, nil
		***REMOVED***
		if t != reflect.TypeOf(time.Duration(5)) ***REMOVED***
			return data, nil
		***REMOVED***

		// Convert it by parsing
		return time.ParseDuration(data.(string))
	***REMOVED***
***REMOVED***

// StringToIPHookFunc returns a DecodeHookFunc that converts
// strings to net.IP
func StringToIPHookFunc() DecodeHookFunc ***REMOVED***
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
		if f.Kind() != reflect.String ***REMOVED***
			return data, nil
		***REMOVED***
		if t != reflect.TypeOf(net.IP***REMOVED******REMOVED***) ***REMOVED***
			return data, nil
		***REMOVED***

		// Convert it by parsing
		ip := net.ParseIP(data.(string))
		if ip == nil ***REMOVED***
			return net.IP***REMOVED******REMOVED***, fmt.Errorf("failed parsing ip %v", data)
		***REMOVED***

		return ip, nil
	***REMOVED***
***REMOVED***

// StringToIPNetHookFunc returns a DecodeHookFunc that converts
// strings to net.IPNet
func StringToIPNetHookFunc() DecodeHookFunc ***REMOVED***
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
		if f.Kind() != reflect.String ***REMOVED***
			return data, nil
		***REMOVED***
		if t != reflect.TypeOf(net.IPNet***REMOVED******REMOVED***) ***REMOVED***
			return data, nil
		***REMOVED***

		// Convert it by parsing
		_, net, err := net.ParseCIDR(data.(string))
		return net, err
	***REMOVED***
***REMOVED***

// StringToTimeHookFunc returns a DecodeHookFunc that converts
// strings to time.Time.
func StringToTimeHookFunc(layout string) DecodeHookFunc ***REMOVED***
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
		if f.Kind() != reflect.String ***REMOVED***
			return data, nil
		***REMOVED***
		if t != reflect.TypeOf(time.Time***REMOVED******REMOVED***) ***REMOVED***
			return data, nil
		***REMOVED***

		// Convert it by parsing
		return time.Parse(layout, data.(string))
	***REMOVED***
***REMOVED***

// WeaklyTypedHook is a DecodeHookFunc which adds support for weak typing to
// the decoder.
//
// Note that this is significantly different from the WeaklyTypedInput option
// of the DecoderConfig.
func WeaklyTypedHook(
	f reflect.Kind,
	t reflect.Kind,
	data interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	dataVal := reflect.ValueOf(data)
	switch t ***REMOVED***
	case reflect.String:
		switch f ***REMOVED***
		case reflect.Bool:
			if dataVal.Bool() ***REMOVED***
				return "1", nil
			***REMOVED***
			return "0", nil
		case reflect.Float32:
			return strconv.FormatFloat(dataVal.Float(), 'f', -1, 64), nil
		case reflect.Int:
			return strconv.FormatInt(dataVal.Int(), 10), nil
		case reflect.Slice:
			dataType := dataVal.Type()
			elemKind := dataType.Elem().Kind()
			if elemKind == reflect.Uint8 ***REMOVED***
				return string(dataVal.Interface().([]uint8)), nil
			***REMOVED***
		case reflect.Uint:
			return strconv.FormatUint(dataVal.Uint(), 10), nil
		***REMOVED***
	***REMOVED***

	return data, nil
***REMOVED***
