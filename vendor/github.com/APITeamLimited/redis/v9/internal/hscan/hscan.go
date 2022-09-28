package hscan

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// decoderFunc represents decoding functions for default built-in types.
type decoderFunc func(reflect.Value, string) error

var (
	// List of built-in decoders indexed by their numeric constant values (eg: reflect.Bool = 1).
	decoders = []decoderFunc***REMOVED***
		reflect.Bool:          decodeBool,
		reflect.Int:           decodeInt,
		reflect.Int8:          decodeInt8,
		reflect.Int16:         decodeInt16,
		reflect.Int32:         decodeInt32,
		reflect.Int64:         decodeInt64,
		reflect.Uint:          decodeUint,
		reflect.Uint8:         decodeUint8,
		reflect.Uint16:        decodeUint16,
		reflect.Uint32:        decodeUint32,
		reflect.Uint64:        decodeUint64,
		reflect.Float32:       decodeFloat32,
		reflect.Float64:       decodeFloat64,
		reflect.Complex64:     decodeUnsupported,
		reflect.Complex128:    decodeUnsupported,
		reflect.Array:         decodeUnsupported,
		reflect.Chan:          decodeUnsupported,
		reflect.Func:          decodeUnsupported,
		reflect.Interface:     decodeUnsupported,
		reflect.Map:           decodeUnsupported,
		reflect.Ptr:           decodeUnsupported,
		reflect.Slice:         decodeSlice,
		reflect.String:        decodeString,
		reflect.Struct:        decodeUnsupported,
		reflect.UnsafePointer: decodeUnsupported,
	***REMOVED***

	// Global map of struct field specs that is populated once for every new
	// struct type that is scanned. This caches the field types and the corresponding
	// decoder functions to avoid iterating through struct fields on subsequent scans.
	globalStructMap = newStructMap()
)

func Struct(dst interface***REMOVED******REMOVED***) (StructValue, error) ***REMOVED***
	v := reflect.ValueOf(dst)

	// The destination to scan into should be a struct pointer.
	if v.Kind() != reflect.Ptr || v.IsNil() ***REMOVED***
		return StructValue***REMOVED******REMOVED***, fmt.Errorf("redis.Scan(non-pointer %T)", dst)
	***REMOVED***

	v = v.Elem()
	if v.Kind() != reflect.Struct ***REMOVED***
		return StructValue***REMOVED******REMOVED***, fmt.Errorf("redis.Scan(non-struct %T)", dst)
	***REMOVED***

	return StructValue***REMOVED***
		spec:  globalStructMap.get(v.Type()),
		value: v,
	***REMOVED***, nil
***REMOVED***

// Scan scans the results from a key-value Redis map result set to a destination struct.
// The Redis keys are matched to the struct's field with the `redis` tag.
func Scan(dst interface***REMOVED******REMOVED***, keys []interface***REMOVED******REMOVED***, vals []interface***REMOVED******REMOVED***) error ***REMOVED***
	if len(keys) != len(vals) ***REMOVED***
		return errors.New("args should have the same number of keys and vals")
	***REMOVED***

	strct, err := Struct(dst)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Iterate through the (key, value) sequence.
	for i := 0; i < len(vals); i++ ***REMOVED***
		key, ok := keys[i].(string)
		if !ok ***REMOVED***
			continue
		***REMOVED***

		val, ok := vals[i].(string)
		if !ok ***REMOVED***
			continue
		***REMOVED***

		if err := strct.Scan(key, val); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func decodeBool(f reflect.Value, s string) error ***REMOVED***
	b, err := strconv.ParseBool(s)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	f.SetBool(b)
	return nil
***REMOVED***

func decodeInt8(f reflect.Value, s string) error ***REMOVED***
	return decodeNumber(f, s, 8)
***REMOVED***

func decodeInt16(f reflect.Value, s string) error ***REMOVED***
	return decodeNumber(f, s, 16)
***REMOVED***

func decodeInt32(f reflect.Value, s string) error ***REMOVED***
	return decodeNumber(f, s, 32)
***REMOVED***

func decodeInt64(f reflect.Value, s string) error ***REMOVED***
	return decodeNumber(f, s, 64)
***REMOVED***

func decodeInt(f reflect.Value, s string) error ***REMOVED***
	return decodeNumber(f, s, 0)
***REMOVED***

func decodeNumber(f reflect.Value, s string, bitSize int) error ***REMOVED***
	v, err := strconv.ParseInt(s, 10, bitSize)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	f.SetInt(v)
	return nil
***REMOVED***

func decodeUint8(f reflect.Value, s string) error ***REMOVED***
	return decodeUnsignedNumber(f, s, 8)
***REMOVED***

func decodeUint16(f reflect.Value, s string) error ***REMOVED***
	return decodeUnsignedNumber(f, s, 16)
***REMOVED***

func decodeUint32(f reflect.Value, s string) error ***REMOVED***
	return decodeUnsignedNumber(f, s, 32)
***REMOVED***

func decodeUint64(f reflect.Value, s string) error ***REMOVED***
	return decodeUnsignedNumber(f, s, 64)
***REMOVED***

func decodeUint(f reflect.Value, s string) error ***REMOVED***
	return decodeUnsignedNumber(f, s, 0)
***REMOVED***

func decodeUnsignedNumber(f reflect.Value, s string, bitSize int) error ***REMOVED***
	v, err := strconv.ParseUint(s, 10, bitSize)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	f.SetUint(v)
	return nil
***REMOVED***

func decodeFloat32(f reflect.Value, s string) error ***REMOVED***
	v, err := strconv.ParseFloat(s, 32)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	f.SetFloat(v)
	return nil
***REMOVED***

// although the default is float64, but we better define it.
func decodeFloat64(f reflect.Value, s string) error ***REMOVED***
	v, err := strconv.ParseFloat(s, 64)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	f.SetFloat(v)
	return nil
***REMOVED***

func decodeString(f reflect.Value, s string) error ***REMOVED***
	f.SetString(s)
	return nil
***REMOVED***

func decodeSlice(f reflect.Value, s string) error ***REMOVED***
	// []byte slice ([]uint8).
	if f.Type().Elem().Kind() == reflect.Uint8 ***REMOVED***
		f.SetBytes([]byte(s))
	***REMOVED***
	return nil
***REMOVED***

func decodeUnsupported(v reflect.Value, s string) error ***REMOVED***
	return fmt.Errorf("redis.Scan(unsupported %s)", v.Type())
***REMOVED***
