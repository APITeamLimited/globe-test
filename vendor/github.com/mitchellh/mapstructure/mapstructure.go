// Package mapstructure exposes functionality to convert an arbitrary
// map[string]interface***REMOVED******REMOVED*** into a native Go structure.
//
// The Go structure can be arbitrarily complex, containing slices,
// other structs, etc. and the decoder will properly decode nested
// maps and so on into the proper structures in the native Go struct.
// See the examples to see what the decoder is capable of.
package mapstructure

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// DecodeHookFunc is the callback function that can be used for
// data transformations. See "DecodeHook" in the DecoderConfig
// struct.
//
// The type should be DecodeHookFuncType or DecodeHookFuncKind.
// Either is accepted. Types are a superset of Kinds (Types can return
// Kinds) and are generally a richer thing to use, but Kinds are simpler
// if you only need those.
//
// The reason DecodeHookFunc is multi-typed is for backwards compatibility:
// we started with Kinds and then realized Types were the better solution,
// but have a promise to not break backwards compat so we now support
// both.
type DecodeHookFunc interface***REMOVED******REMOVED***

// DecodeHookFuncType is a DecodeHookFunc which has complete information about
// the source and target types.
type DecodeHookFuncType func(reflect.Type, reflect.Type, interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error)

// DecodeHookFuncKind is a DecodeHookFunc which knows only the Kinds of the
// source and target types.
type DecodeHookFuncKind func(reflect.Kind, reflect.Kind, interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error)

// DecoderConfig is the configuration that is used to create a new decoder
// and allows customization of various aspects of decoding.
type DecoderConfig struct ***REMOVED***
	// DecodeHook, if set, will be called before any decoding and any
	// type conversion (if WeaklyTypedInput is on). This lets you modify
	// the values before they're set down onto the resulting struct.
	//
	// If an error is returned, the entire decode will fail with that
	// error.
	DecodeHook DecodeHookFunc

	// If ErrorUnused is true, then it is an error for there to exist
	// keys in the original map that were unused in the decoding process
	// (extra keys).
	ErrorUnused bool

	// ZeroFields, if set to true, will zero fields before writing them.
	// For example, a map will be emptied before decoded values are put in
	// it. If this is false, a map will be merged.
	ZeroFields bool

	// If WeaklyTypedInput is true, the decoder will make the following
	// "weak" conversions:
	//
	//   - bools to string (true = "1", false = "0")
	//   - numbers to string (base 10)
	//   - bools to int/uint (true = 1, false = 0)
	//   - strings to int/uint (base implied by prefix)
	//   - int to bool (true if value != 0)
	//   - string to bool (accepts: 1, t, T, TRUE, true, True, 0, f, F,
	//     FALSE, false, False. Anything else is an error)
	//   - empty array = empty map and vice versa
	//   - negative numbers to overflowed uint values (base 10)
	//   - slice of maps to a merged map
	//   - single values are converted to slices if required. Each
	//     element is weakly decoded. For example: "4" can become []int***REMOVED***4***REMOVED***
	//     if the target type is an int slice.
	//
	WeaklyTypedInput bool

	// Metadata is the struct that will contain extra metadata about
	// the decoding. If this is nil, then no metadata will be tracked.
	Metadata *Metadata

	// Result is a pointer to the struct that will contain the decoded
	// value.
	Result interface***REMOVED******REMOVED***

	// The tag name that mapstructure reads for field names. This
	// defaults to "mapstructure"
	TagName string
***REMOVED***

// A Decoder takes a raw interface value and turns it into structured
// data, keeping track of rich error information along the way in case
// anything goes wrong. Unlike the basic top-level Decode method, you can
// more finely control how the Decoder behaves using the DecoderConfig
// structure. The top-level Decode method is just a convenience that sets
// up the most basic Decoder.
type Decoder struct ***REMOVED***
	config *DecoderConfig
***REMOVED***

// Metadata contains information about decoding a structure that
// is tedious or difficult to get otherwise.
type Metadata struct ***REMOVED***
	// Keys are the keys of the structure which were successfully decoded
	Keys []string

	// Unused is a slice of keys that were found in the raw value but
	// weren't decoded since there was no matching field in the result interface
	Unused []string
***REMOVED***

// Decode takes an input structure and uses reflection to translate it to
// the output structure. output must be a pointer to a map or struct.
func Decode(input interface***REMOVED******REMOVED***, output interface***REMOVED******REMOVED***) error ***REMOVED***
	config := &DecoderConfig***REMOVED***
		Metadata: nil,
		Result:   output,
	***REMOVED***

	decoder, err := NewDecoder(config)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return decoder.Decode(input)
***REMOVED***

// WeakDecode is the same as Decode but is shorthand to enable
// WeaklyTypedInput. See DecoderConfig for more info.
func WeakDecode(input, output interface***REMOVED******REMOVED***) error ***REMOVED***
	config := &DecoderConfig***REMOVED***
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
	***REMOVED***

	decoder, err := NewDecoder(config)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return decoder.Decode(input)
***REMOVED***

// DecodeMetadata is the same as Decode, but is shorthand to
// enable metadata collection. See DecoderConfig for more info.
func DecodeMetadata(input interface***REMOVED******REMOVED***, output interface***REMOVED******REMOVED***, metadata *Metadata) error ***REMOVED***
	config := &DecoderConfig***REMOVED***
		Metadata: metadata,
		Result:   output,
	***REMOVED***

	decoder, err := NewDecoder(config)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return decoder.Decode(input)
***REMOVED***

// WeakDecodeMetadata is the same as Decode, but is shorthand to
// enable both WeaklyTypedInput and metadata collection. See
// DecoderConfig for more info.
func WeakDecodeMetadata(input interface***REMOVED******REMOVED***, output interface***REMOVED******REMOVED***, metadata *Metadata) error ***REMOVED***
	config := &DecoderConfig***REMOVED***
		Metadata:         metadata,
		Result:           output,
		WeaklyTypedInput: true,
	***REMOVED***

	decoder, err := NewDecoder(config)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return decoder.Decode(input)
***REMOVED***

// NewDecoder returns a new decoder for the given configuration. Once
// a decoder has been returned, the same configuration must not be used
// again.
func NewDecoder(config *DecoderConfig) (*Decoder, error) ***REMOVED***
	val := reflect.ValueOf(config.Result)
	if val.Kind() != reflect.Ptr ***REMOVED***
		return nil, errors.New("result must be a pointer")
	***REMOVED***

	val = val.Elem()
	if !val.CanAddr() ***REMOVED***
		return nil, errors.New("result must be addressable (a pointer)")
	***REMOVED***

	if config.Metadata != nil ***REMOVED***
		if config.Metadata.Keys == nil ***REMOVED***
			config.Metadata.Keys = make([]string, 0)
		***REMOVED***

		if config.Metadata.Unused == nil ***REMOVED***
			config.Metadata.Unused = make([]string, 0)
		***REMOVED***
	***REMOVED***

	if config.TagName == "" ***REMOVED***
		config.TagName = "mapstructure"
	***REMOVED***

	result := &Decoder***REMOVED***
		config: config,
	***REMOVED***

	return result, nil
***REMOVED***

// Decode decodes the given raw interface to the target pointer specified
// by the configuration.
func (d *Decoder) Decode(input interface***REMOVED******REMOVED***) error ***REMOVED***
	return d.decode("", input, reflect.ValueOf(d.config.Result).Elem())
***REMOVED***

// Decodes an unknown data type into a specific reflection value.
func (d *Decoder) decode(name string, input interface***REMOVED******REMOVED***, outVal reflect.Value) error ***REMOVED***
	var inputVal reflect.Value
	if input != nil ***REMOVED***
		inputVal = reflect.ValueOf(input)

		// We need to check here if input is a typed nil. Typed nils won't
		// match the "input == nil" below so we check that here.
		if inputVal.Kind() == reflect.Ptr && inputVal.IsNil() ***REMOVED***
			input = nil
		***REMOVED***
	***REMOVED***

	if input == nil ***REMOVED***
		// If the data is nil, then we don't set anything, unless ZeroFields is set
		// to true.
		if d.config.ZeroFields ***REMOVED***
			outVal.Set(reflect.Zero(outVal.Type()))

			if d.config.Metadata != nil && name != "" ***REMOVED***
				d.config.Metadata.Keys = append(d.config.Metadata.Keys, name)
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***

	if !inputVal.IsValid() ***REMOVED***
		// If the input value is invalid, then we just set the value
		// to be the zero value.
		outVal.Set(reflect.Zero(outVal.Type()))
		if d.config.Metadata != nil && name != "" ***REMOVED***
			d.config.Metadata.Keys = append(d.config.Metadata.Keys, name)
		***REMOVED***
		return nil
	***REMOVED***

	if d.config.DecodeHook != nil ***REMOVED***
		// We have a DecodeHook, so let's pre-process the input.
		var err error
		input, err = DecodeHookExec(
			d.config.DecodeHook,
			inputVal.Type(), outVal.Type(), input)
		if err != nil ***REMOVED***
			return fmt.Errorf("error decoding '%s': %s", name, err)
		***REMOVED***
	***REMOVED***

	var err error
	outputKind := getKind(outVal)
	switch outputKind ***REMOVED***
	case reflect.Bool:
		err = d.decodeBool(name, input, outVal)
	case reflect.Interface:
		err = d.decodeBasic(name, input, outVal)
	case reflect.String:
		err = d.decodeString(name, input, outVal)
	case reflect.Int:
		err = d.decodeInt(name, input, outVal)
	case reflect.Uint:
		err = d.decodeUint(name, input, outVal)
	case reflect.Float32:
		err = d.decodeFloat(name, input, outVal)
	case reflect.Struct:
		err = d.decodeStruct(name, input, outVal)
	case reflect.Map:
		err = d.decodeMap(name, input, outVal)
	case reflect.Ptr:
		err = d.decodePtr(name, input, outVal)
	case reflect.Slice:
		err = d.decodeSlice(name, input, outVal)
	case reflect.Array:
		err = d.decodeArray(name, input, outVal)
	case reflect.Func:
		err = d.decodeFunc(name, input, outVal)
	default:
		// If we reached this point then we weren't able to decode it
		return fmt.Errorf("%s: unsupported type: %s", name, outputKind)
	***REMOVED***

	// If we reached here, then we successfully decoded SOMETHING, so
	// mark the key as used if we're tracking metainput.
	if d.config.Metadata != nil && name != "" ***REMOVED***
		d.config.Metadata.Keys = append(d.config.Metadata.Keys, name)
	***REMOVED***

	return err
***REMOVED***

// This decodes a basic type (bool, int, string, etc.) and sets the
// value to "data" of that type.
func (d *Decoder) decodeBasic(name string, data interface***REMOVED******REMOVED***, val reflect.Value) error ***REMOVED***
	if val.IsValid() && val.Elem().IsValid() ***REMOVED***
		return d.decode(name, data, val.Elem())
	***REMOVED***

	dataVal := reflect.ValueOf(data)

	// If the input data is a pointer, and the assigned type is the dereference
	// of that exact pointer, then indirect it so that we can assign it.
	// Example: *string to string
	if dataVal.Kind() == reflect.Ptr && dataVal.Type().Elem() == val.Type() ***REMOVED***
		dataVal = reflect.Indirect(dataVal)
	***REMOVED***

	if !dataVal.IsValid() ***REMOVED***
		dataVal = reflect.Zero(val.Type())
	***REMOVED***

	dataValType := dataVal.Type()
	if !dataValType.AssignableTo(val.Type()) ***REMOVED***
		return fmt.Errorf(
			"'%s' expected type '%s', got '%s'",
			name, val.Type(), dataValType)
	***REMOVED***

	val.Set(dataVal)
	return nil
***REMOVED***

func (d *Decoder) decodeString(name string, data interface***REMOVED******REMOVED***, val reflect.Value) error ***REMOVED***
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)

	converted := true
	switch ***REMOVED***
	case dataKind == reflect.String:
		val.SetString(dataVal.String())
	case dataKind == reflect.Bool && d.config.WeaklyTypedInput:
		if dataVal.Bool() ***REMOVED***
			val.SetString("1")
		***REMOVED*** else ***REMOVED***
			val.SetString("0")
		***REMOVED***
	case dataKind == reflect.Int && d.config.WeaklyTypedInput:
		val.SetString(strconv.FormatInt(dataVal.Int(), 10))
	case dataKind == reflect.Uint && d.config.WeaklyTypedInput:
		val.SetString(strconv.FormatUint(dataVal.Uint(), 10))
	case dataKind == reflect.Float32 && d.config.WeaklyTypedInput:
		val.SetString(strconv.FormatFloat(dataVal.Float(), 'f', -1, 64))
	case dataKind == reflect.Slice && d.config.WeaklyTypedInput,
		dataKind == reflect.Array && d.config.WeaklyTypedInput:
		dataType := dataVal.Type()
		elemKind := dataType.Elem().Kind()
		switch elemKind ***REMOVED***
		case reflect.Uint8:
			var uints []uint8
			if dataKind == reflect.Array ***REMOVED***
				uints = make([]uint8, dataVal.Len(), dataVal.Len())
				for i := range uints ***REMOVED***
					uints[i] = dataVal.Index(i).Interface().(uint8)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				uints = dataVal.Interface().([]uint8)
			***REMOVED***
			val.SetString(string(uints))
		default:
			converted = false
		***REMOVED***
	default:
		converted = false
	***REMOVED***

	if !converted ***REMOVED***
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s'",
			name, val.Type(), dataVal.Type())
	***REMOVED***

	return nil
***REMOVED***

func (d *Decoder) decodeInt(name string, data interface***REMOVED******REMOVED***, val reflect.Value) error ***REMOVED***
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)
	dataType := dataVal.Type()

	switch ***REMOVED***
	case dataKind == reflect.Int:
		val.SetInt(dataVal.Int())
	case dataKind == reflect.Uint:
		val.SetInt(int64(dataVal.Uint()))
	case dataKind == reflect.Float32:
		val.SetInt(int64(dataVal.Float()))
	case dataKind == reflect.Bool && d.config.WeaklyTypedInput:
		if dataVal.Bool() ***REMOVED***
			val.SetInt(1)
		***REMOVED*** else ***REMOVED***
			val.SetInt(0)
		***REMOVED***
	case dataKind == reflect.String && d.config.WeaklyTypedInput:
		i, err := strconv.ParseInt(dataVal.String(), 0, val.Type().Bits())
		if err == nil ***REMOVED***
			val.SetInt(i)
		***REMOVED*** else ***REMOVED***
			return fmt.Errorf("cannot parse '%s' as int: %s", name, err)
		***REMOVED***
	case dataType.PkgPath() == "encoding/json" && dataType.Name() == "Number":
		jn := data.(json.Number)
		i, err := jn.Int64()
		if err != nil ***REMOVED***
			return fmt.Errorf(
				"error decoding json.Number into %s: %s", name, err)
		***REMOVED***
		val.SetInt(i)
	default:
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s'",
			name, val.Type(), dataVal.Type())
	***REMOVED***

	return nil
***REMOVED***

func (d *Decoder) decodeUint(name string, data interface***REMOVED******REMOVED***, val reflect.Value) error ***REMOVED***
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)

	switch ***REMOVED***
	case dataKind == reflect.Int:
		i := dataVal.Int()
		if i < 0 && !d.config.WeaklyTypedInput ***REMOVED***
			return fmt.Errorf("cannot parse '%s', %d overflows uint",
				name, i)
		***REMOVED***
		val.SetUint(uint64(i))
	case dataKind == reflect.Uint:
		val.SetUint(dataVal.Uint())
	case dataKind == reflect.Float32:
		f := dataVal.Float()
		if f < 0 && !d.config.WeaklyTypedInput ***REMOVED***
			return fmt.Errorf("cannot parse '%s', %f overflows uint",
				name, f)
		***REMOVED***
		val.SetUint(uint64(f))
	case dataKind == reflect.Bool && d.config.WeaklyTypedInput:
		if dataVal.Bool() ***REMOVED***
			val.SetUint(1)
		***REMOVED*** else ***REMOVED***
			val.SetUint(0)
		***REMOVED***
	case dataKind == reflect.String && d.config.WeaklyTypedInput:
		i, err := strconv.ParseUint(dataVal.String(), 0, val.Type().Bits())
		if err == nil ***REMOVED***
			val.SetUint(i)
		***REMOVED*** else ***REMOVED***
			return fmt.Errorf("cannot parse '%s' as uint: %s", name, err)
		***REMOVED***
	default:
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s'",
			name, val.Type(), dataVal.Type())
	***REMOVED***

	return nil
***REMOVED***

func (d *Decoder) decodeBool(name string, data interface***REMOVED******REMOVED***, val reflect.Value) error ***REMOVED***
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)

	switch ***REMOVED***
	case dataKind == reflect.Bool:
		val.SetBool(dataVal.Bool())
	case dataKind == reflect.Int && d.config.WeaklyTypedInput:
		val.SetBool(dataVal.Int() != 0)
	case dataKind == reflect.Uint && d.config.WeaklyTypedInput:
		val.SetBool(dataVal.Uint() != 0)
	case dataKind == reflect.Float32 && d.config.WeaklyTypedInput:
		val.SetBool(dataVal.Float() != 0)
	case dataKind == reflect.String && d.config.WeaklyTypedInput:
		b, err := strconv.ParseBool(dataVal.String())
		if err == nil ***REMOVED***
			val.SetBool(b)
		***REMOVED*** else if dataVal.String() == "" ***REMOVED***
			val.SetBool(false)
		***REMOVED*** else ***REMOVED***
			return fmt.Errorf("cannot parse '%s' as bool: %s", name, err)
		***REMOVED***
	default:
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s'",
			name, val.Type(), dataVal.Type())
	***REMOVED***

	return nil
***REMOVED***

func (d *Decoder) decodeFloat(name string, data interface***REMOVED******REMOVED***, val reflect.Value) error ***REMOVED***
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)
	dataType := dataVal.Type()

	switch ***REMOVED***
	case dataKind == reflect.Int:
		val.SetFloat(float64(dataVal.Int()))
	case dataKind == reflect.Uint:
		val.SetFloat(float64(dataVal.Uint()))
	case dataKind == reflect.Float32:
		val.SetFloat(dataVal.Float())
	case dataKind == reflect.Bool && d.config.WeaklyTypedInput:
		if dataVal.Bool() ***REMOVED***
			val.SetFloat(1)
		***REMOVED*** else ***REMOVED***
			val.SetFloat(0)
		***REMOVED***
	case dataKind == reflect.String && d.config.WeaklyTypedInput:
		f, err := strconv.ParseFloat(dataVal.String(), val.Type().Bits())
		if err == nil ***REMOVED***
			val.SetFloat(f)
		***REMOVED*** else ***REMOVED***
			return fmt.Errorf("cannot parse '%s' as float: %s", name, err)
		***REMOVED***
	case dataType.PkgPath() == "encoding/json" && dataType.Name() == "Number":
		jn := data.(json.Number)
		i, err := jn.Float64()
		if err != nil ***REMOVED***
			return fmt.Errorf(
				"error decoding json.Number into %s: %s", name, err)
		***REMOVED***
		val.SetFloat(i)
	default:
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s'",
			name, val.Type(), dataVal.Type())
	***REMOVED***

	return nil
***REMOVED***

func (d *Decoder) decodeMap(name string, data interface***REMOVED******REMOVED***, val reflect.Value) error ***REMOVED***
	valType := val.Type()
	valKeyType := valType.Key()
	valElemType := valType.Elem()

	// By default we overwrite keys in the current map
	valMap := val

	// If the map is nil or we're purposely zeroing fields, make a new map
	if valMap.IsNil() || d.config.ZeroFields ***REMOVED***
		// Make a new map to hold our result
		mapType := reflect.MapOf(valKeyType, valElemType)
		valMap = reflect.MakeMap(mapType)
	***REMOVED***

	// Check input type and based on the input type jump to the proper func
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	switch dataVal.Kind() ***REMOVED***
	case reflect.Map:
		return d.decodeMapFromMap(name, dataVal, val, valMap)

	case reflect.Struct:
		return d.decodeMapFromStruct(name, dataVal, val, valMap)

	case reflect.Array, reflect.Slice:
		if d.config.WeaklyTypedInput ***REMOVED***
			return d.decodeMapFromSlice(name, dataVal, val, valMap)
		***REMOVED***

		fallthrough

	default:
		return fmt.Errorf("'%s' expected a map, got '%s'", name, dataVal.Kind())
	***REMOVED***
***REMOVED***

func (d *Decoder) decodeMapFromSlice(name string, dataVal reflect.Value, val reflect.Value, valMap reflect.Value) error ***REMOVED***
	// Special case for BC reasons (covered by tests)
	if dataVal.Len() == 0 ***REMOVED***
		val.Set(valMap)
		return nil
	***REMOVED***

	for i := 0; i < dataVal.Len(); i++ ***REMOVED***
		err := d.decode(
			fmt.Sprintf("%s[%d]", name, i),
			dataVal.Index(i).Interface(), val)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (d *Decoder) decodeMapFromMap(name string, dataVal reflect.Value, val reflect.Value, valMap reflect.Value) error ***REMOVED***
	valType := val.Type()
	valKeyType := valType.Key()
	valElemType := valType.Elem()

	// Accumulate errors
	errors := make([]string, 0)

	// If the input data is empty, then we just match what the input data is.
	if dataVal.Len() == 0 ***REMOVED***
		if dataVal.IsNil() ***REMOVED***
			if !val.IsNil() ***REMOVED***
				val.Set(dataVal)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// Set to empty allocated value
			val.Set(valMap)
		***REMOVED***

		return nil
	***REMOVED***

	for _, k := range dataVal.MapKeys() ***REMOVED***
		fieldName := fmt.Sprintf("%s[%s]", name, k)

		// First decode the key into the proper type
		currentKey := reflect.Indirect(reflect.New(valKeyType))
		if err := d.decode(fieldName, k.Interface(), currentKey); err != nil ***REMOVED***
			errors = appendErrors(errors, err)
			continue
		***REMOVED***

		// Next decode the data into the proper type
		v := dataVal.MapIndex(k).Interface()
		currentVal := reflect.Indirect(reflect.New(valElemType))
		if err := d.decode(fieldName, v, currentVal); err != nil ***REMOVED***
			errors = appendErrors(errors, err)
			continue
		***REMOVED***

		valMap.SetMapIndex(currentKey, currentVal)
	***REMOVED***

	// Set the built up map to the value
	val.Set(valMap)

	// If we had errors, return those
	if len(errors) > 0 ***REMOVED***
		return &Error***REMOVED***errors***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (d *Decoder) decodeMapFromStruct(name string, dataVal reflect.Value, val reflect.Value, valMap reflect.Value) error ***REMOVED***
	typ := dataVal.Type()
	for i := 0; i < typ.NumField(); i++ ***REMOVED***
		// Get the StructField first since this is a cheap operation. If the
		// field is unexported, then ignore it.
		f := typ.Field(i)
		if f.PkgPath != "" ***REMOVED***
			continue
		***REMOVED***

		// Next get the actual value of this field and verify it is assignable
		// to the map value.
		v := dataVal.Field(i)
		if !v.Type().AssignableTo(valMap.Type().Elem()) ***REMOVED***
			return fmt.Errorf("cannot assign type '%s' to map value field of type '%s'", v.Type(), valMap.Type().Elem())
		***REMOVED***

		tagValue := f.Tag.Get(d.config.TagName)
		tagParts := strings.Split(tagValue, ",")

		// Determine the name of the key in the map
		keyName := f.Name
		if tagParts[0] != "" ***REMOVED***
			if tagParts[0] == "-" ***REMOVED***
				continue
			***REMOVED***
			keyName = tagParts[0]
		***REMOVED***

		// If "squash" is specified in the tag, we squash the field down.
		squash := false
		for _, tag := range tagParts[1:] ***REMOVED***
			if tag == "squash" ***REMOVED***
				squash = true
				break
			***REMOVED***
		***REMOVED***
		if squash && v.Kind() != reflect.Struct ***REMOVED***
			return fmt.Errorf("cannot squash non-struct type '%s'", v.Type())
		***REMOVED***

		switch v.Kind() ***REMOVED***
		// this is an embedded struct, so handle it differently
		case reflect.Struct:
			x := reflect.New(v.Type())
			x.Elem().Set(v)

			vType := valMap.Type()
			vKeyType := vType.Key()
			vElemType := vType.Elem()
			mType := reflect.MapOf(vKeyType, vElemType)
			vMap := reflect.MakeMap(mType)

			err := d.decode(keyName, x.Interface(), vMap)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			if squash ***REMOVED***
				for _, k := range vMap.MapKeys() ***REMOVED***
					valMap.SetMapIndex(k, vMap.MapIndex(k))
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				valMap.SetMapIndex(reflect.ValueOf(keyName), vMap)
			***REMOVED***

		default:
			valMap.SetMapIndex(reflect.ValueOf(keyName), v)
		***REMOVED***
	***REMOVED***

	if val.CanAddr() ***REMOVED***
		val.Set(valMap)
	***REMOVED***

	return nil
***REMOVED***

func (d *Decoder) decodePtr(name string, data interface***REMOVED******REMOVED***, val reflect.Value) error ***REMOVED***
	// If the input data is nil, then we want to just set the output
	// pointer to be nil as well.
	isNil := data == nil
	if !isNil ***REMOVED***
		switch v := reflect.Indirect(reflect.ValueOf(data)); v.Kind() ***REMOVED***
		case reflect.Chan,
			reflect.Func,
			reflect.Interface,
			reflect.Map,
			reflect.Ptr,
			reflect.Slice:
			isNil = v.IsNil()
		***REMOVED***
	***REMOVED***
	if isNil ***REMOVED***
		if !val.IsNil() && val.CanSet() ***REMOVED***
			nilValue := reflect.New(val.Type()).Elem()
			val.Set(nilValue)
		***REMOVED***

		return nil
	***REMOVED***

	// Create an element of the concrete (non pointer) type and decode
	// into that. Then set the value of the pointer to this type.
	valType := val.Type()
	valElemType := valType.Elem()
	if val.CanSet() ***REMOVED***
		realVal := val
		if realVal.IsNil() || d.config.ZeroFields ***REMOVED***
			realVal = reflect.New(valElemType)
		***REMOVED***

		if err := d.decode(name, data, reflect.Indirect(realVal)); err != nil ***REMOVED***
			return err
		***REMOVED***

		val.Set(realVal)
	***REMOVED*** else ***REMOVED***
		if err := d.decode(name, data, reflect.Indirect(val)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (d *Decoder) decodeFunc(name string, data interface***REMOVED******REMOVED***, val reflect.Value) error ***REMOVED***
	// Create an element of the concrete (non pointer) type and decode
	// into that. Then set the value of the pointer to this type.
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	if val.Type() != dataVal.Type() ***REMOVED***
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s'",
			name, val.Type(), dataVal.Type())
	***REMOVED***
	val.Set(dataVal)
	return nil
***REMOVED***

func (d *Decoder) decodeSlice(name string, data interface***REMOVED******REMOVED***, val reflect.Value) error ***REMOVED***
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataValKind := dataVal.Kind()
	valType := val.Type()
	valElemType := valType.Elem()
	sliceType := reflect.SliceOf(valElemType)

	valSlice := val
	if valSlice.IsNil() || d.config.ZeroFields ***REMOVED***
		if d.config.WeaklyTypedInput ***REMOVED***
			switch ***REMOVED***
			// Slice and array we use the normal logic
			case dataValKind == reflect.Slice, dataValKind == reflect.Array:
				break

			// Empty maps turn into empty slices
			case dataValKind == reflect.Map:
				if dataVal.Len() == 0 ***REMOVED***
					val.Set(reflect.MakeSlice(sliceType, 0, 0))
					return nil
				***REMOVED***
				// Create slice of maps of other sizes
				return d.decodeSlice(name, []interface***REMOVED******REMOVED******REMOVED***data***REMOVED***, val)

			case dataValKind == reflect.String && valElemType.Kind() == reflect.Uint8:
				return d.decodeSlice(name, []byte(dataVal.String()), val)

			// All other types we try to convert to the slice type
			// and "lift" it into it. i.e. a string becomes a string slice.
			default:
				// Just re-try this function with data as a slice.
				return d.decodeSlice(name, []interface***REMOVED******REMOVED******REMOVED***data***REMOVED***, val)
			***REMOVED***
		***REMOVED***

		// Check input type
		if dataValKind != reflect.Array && dataValKind != reflect.Slice ***REMOVED***
			return fmt.Errorf(
				"'%s': source data must be an array or slice, got %s", name, dataValKind)

		***REMOVED***

		// If the input value is empty, then don't allocate since non-nil != nil
		if dataVal.Len() == 0 ***REMOVED***
			return nil
		***REMOVED***

		// Make a new slice to hold our result, same size as the original data.
		valSlice = reflect.MakeSlice(sliceType, dataVal.Len(), dataVal.Len())
	***REMOVED***

	// Accumulate any errors
	errors := make([]string, 0)

	for i := 0; i < dataVal.Len(); i++ ***REMOVED***
		currentData := dataVal.Index(i).Interface()
		for valSlice.Len() <= i ***REMOVED***
			valSlice = reflect.Append(valSlice, reflect.Zero(valElemType))
		***REMOVED***
		currentField := valSlice.Index(i)

		fieldName := fmt.Sprintf("%s[%d]", name, i)
		if err := d.decode(fieldName, currentData, currentField); err != nil ***REMOVED***
			errors = appendErrors(errors, err)
		***REMOVED***
	***REMOVED***

	// Finally, set the value to the slice we built up
	val.Set(valSlice)

	// If there were errors, we return those
	if len(errors) > 0 ***REMOVED***
		return &Error***REMOVED***errors***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (d *Decoder) decodeArray(name string, data interface***REMOVED******REMOVED***, val reflect.Value) error ***REMOVED***
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataValKind := dataVal.Kind()
	valType := val.Type()
	valElemType := valType.Elem()
	arrayType := reflect.ArrayOf(valType.Len(), valElemType)

	valArray := val

	if valArray.Interface() == reflect.Zero(valArray.Type()).Interface() || d.config.ZeroFields ***REMOVED***
		// Check input type
		if dataValKind != reflect.Array && dataValKind != reflect.Slice ***REMOVED***
			if d.config.WeaklyTypedInput ***REMOVED***
				switch ***REMOVED***
				// Empty maps turn into empty arrays
				case dataValKind == reflect.Map:
					if dataVal.Len() == 0 ***REMOVED***
						val.Set(reflect.Zero(arrayType))
						return nil
					***REMOVED***

				// All other types we try to convert to the array type
				// and "lift" it into it. i.e. a string becomes a string array.
				default:
					// Just re-try this function with data as a slice.
					return d.decodeArray(name, []interface***REMOVED******REMOVED******REMOVED***data***REMOVED***, val)
				***REMOVED***
			***REMOVED***

			return fmt.Errorf(
				"'%s': source data must be an array or slice, got %s", name, dataValKind)

		***REMOVED***
		if dataVal.Len() > arrayType.Len() ***REMOVED***
			return fmt.Errorf(
				"'%s': expected source data to have length less or equal to %d, got %d", name, arrayType.Len(), dataVal.Len())

		***REMOVED***

		// Make a new array to hold our result, same size as the original data.
		valArray = reflect.New(arrayType).Elem()
	***REMOVED***

	// Accumulate any errors
	errors := make([]string, 0)

	for i := 0; i < dataVal.Len(); i++ ***REMOVED***
		currentData := dataVal.Index(i).Interface()
		currentField := valArray.Index(i)

		fieldName := fmt.Sprintf("%s[%d]", name, i)
		if err := d.decode(fieldName, currentData, currentField); err != nil ***REMOVED***
			errors = appendErrors(errors, err)
		***REMOVED***
	***REMOVED***

	// Finally, set the value to the array we built up
	val.Set(valArray)

	// If there were errors, we return those
	if len(errors) > 0 ***REMOVED***
		return &Error***REMOVED***errors***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (d *Decoder) decodeStruct(name string, data interface***REMOVED******REMOVED***, val reflect.Value) error ***REMOVED***
	dataVal := reflect.Indirect(reflect.ValueOf(data))

	// If the type of the value to write to and the data match directly,
	// then we just set it directly instead of recursing into the structure.
	if dataVal.Type() == val.Type() ***REMOVED***
		val.Set(dataVal)
		return nil
	***REMOVED***

	dataValKind := dataVal.Kind()
	switch dataValKind ***REMOVED***
	case reflect.Map:
		return d.decodeStructFromMap(name, dataVal, val)

	case reflect.Struct:
		// Not the most efficient way to do this but we can optimize later if
		// we want to. To convert from struct to struct we go to map first
		// as an intermediary.
		m := make(map[string]interface***REMOVED******REMOVED***)
		mval := reflect.Indirect(reflect.ValueOf(&m))
		if err := d.decodeMapFromStruct(name, dataVal, mval, mval); err != nil ***REMOVED***
			return err
		***REMOVED***

		result := d.decodeStructFromMap(name, mval, val)
		return result

	default:
		return fmt.Errorf("'%s' expected a map, got '%s'", name, dataVal.Kind())
	***REMOVED***
***REMOVED***

func (d *Decoder) decodeStructFromMap(name string, dataVal, val reflect.Value) error ***REMOVED***
	dataValType := dataVal.Type()
	if kind := dataValType.Key().Kind(); kind != reflect.String && kind != reflect.Interface ***REMOVED***
		return fmt.Errorf(
			"'%s' needs a map with string keys, has '%s' keys",
			name, dataValType.Key().Kind())
	***REMOVED***

	dataValKeys := make(map[reflect.Value]struct***REMOVED******REMOVED***)
	dataValKeysUnused := make(map[interface***REMOVED******REMOVED***]struct***REMOVED******REMOVED***)
	for _, dataValKey := range dataVal.MapKeys() ***REMOVED***
		dataValKeys[dataValKey] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		dataValKeysUnused[dataValKey.Interface()] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	errors := make([]string, 0)

	// This slice will keep track of all the structs we'll be decoding.
	// There can be more than one struct if there are embedded structs
	// that are squashed.
	structs := make([]reflect.Value, 1, 5)
	structs[0] = val

	// Compile the list of all the fields that we're going to be decoding
	// from all the structs.
	type field struct ***REMOVED***
		field reflect.StructField
		val   reflect.Value
	***REMOVED***
	fields := []field***REMOVED******REMOVED***
	for len(structs) > 0 ***REMOVED***
		structVal := structs[0]
		structs = structs[1:]

		structType := structVal.Type()

		for i := 0; i < structType.NumField(); i++ ***REMOVED***
			fieldType := structType.Field(i)
			fieldKind := fieldType.Type.Kind()

			// If "squash" is specified in the tag, we squash the field down.
			squash := false
			tagParts := strings.Split(fieldType.Tag.Get(d.config.TagName), ",")
			for _, tag := range tagParts[1:] ***REMOVED***
				if tag == "squash" ***REMOVED***
					squash = true
					break
				***REMOVED***
			***REMOVED***

			if squash ***REMOVED***
				if fieldKind != reflect.Struct ***REMOVED***
					errors = appendErrors(errors,
						fmt.Errorf("%s: unsupported type for squash: %s", fieldType.Name, fieldKind))
				***REMOVED*** else ***REMOVED***
					structs = append(structs, structVal.FieldByName(fieldType.Name))
				***REMOVED***
				continue
			***REMOVED***

			// Normal struct field, store it away
			fields = append(fields, field***REMOVED***fieldType, structVal.Field(i)***REMOVED***)
		***REMOVED***
	***REMOVED***

	// for fieldType, field := range fields ***REMOVED***
	for _, f := range fields ***REMOVED***
		field, fieldValue := f.field, f.val
		fieldName := field.Name

		tagValue := field.Tag.Get(d.config.TagName)
		tagValue = strings.SplitN(tagValue, ",", 2)[0]
		if tagValue != "" ***REMOVED***
			fieldName = tagValue
		***REMOVED***

		rawMapKey := reflect.ValueOf(fieldName)
		rawMapVal := dataVal.MapIndex(rawMapKey)
		if !rawMapVal.IsValid() ***REMOVED***
			// Do a slower search by iterating over each key and
			// doing case-insensitive search.
			for dataValKey := range dataValKeys ***REMOVED***
				mK, ok := dataValKey.Interface().(string)
				if !ok ***REMOVED***
					// Not a string key
					continue
				***REMOVED***

				if strings.EqualFold(mK, fieldName) ***REMOVED***
					rawMapKey = dataValKey
					rawMapVal = dataVal.MapIndex(dataValKey)
					break
				***REMOVED***
			***REMOVED***

			if !rawMapVal.IsValid() ***REMOVED***
				// There was no matching key in the map for the value in
				// the struct. Just ignore.
				continue
			***REMOVED***
		***REMOVED***

		// Delete the key we're using from the unused map so we stop tracking
		delete(dataValKeysUnused, rawMapKey.Interface())

		if !fieldValue.IsValid() ***REMOVED***
			// This should never happen
			panic("field is not valid")
		***REMOVED***

		// If we can't set the field, then it is unexported or something,
		// and we just continue onwards.
		if !fieldValue.CanSet() ***REMOVED***
			continue
		***REMOVED***

		// If the name is empty string, then we're at the root, and we
		// don't dot-join the fields.
		if name != "" ***REMOVED***
			fieldName = fmt.Sprintf("%s.%s", name, fieldName)
		***REMOVED***

		if err := d.decode(fieldName, rawMapVal.Interface(), fieldValue); err != nil ***REMOVED***
			errors = appendErrors(errors, err)
		***REMOVED***
	***REMOVED***

	if d.config.ErrorUnused && len(dataValKeysUnused) > 0 ***REMOVED***
		keys := make([]string, 0, len(dataValKeysUnused))
		for rawKey := range dataValKeysUnused ***REMOVED***
			keys = append(keys, rawKey.(string))
		***REMOVED***
		sort.Strings(keys)

		err := fmt.Errorf("'%s' has invalid keys: %s", name, strings.Join(keys, ", "))
		errors = appendErrors(errors, err)
	***REMOVED***

	if len(errors) > 0 ***REMOVED***
		return &Error***REMOVED***errors***REMOVED***
	***REMOVED***

	// Add the unused keys to the list of unused keys if we're tracking metadata
	if d.config.Metadata != nil ***REMOVED***
		for rawKey := range dataValKeysUnused ***REMOVED***
			key := rawKey.(string)
			if name != "" ***REMOVED***
				key = fmt.Sprintf("%s.%s", name, key)
			***REMOVED***

			d.config.Metadata.Unused = append(d.config.Metadata.Unused, key)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func getKind(val reflect.Value) reflect.Kind ***REMOVED***
	kind := val.Kind()

	switch ***REMOVED***
	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float32
	default:
		return kind
	***REMOVED***
***REMOVED***
