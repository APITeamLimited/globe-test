package models // import "github.com/influxdata/influxdb1-client/models"

import (
	"reflect"
	"strconv"
	"unsafe"
)

// parseIntBytes is a zero-alloc wrapper around strconv.ParseInt.
func parseIntBytes(b []byte, base int, bitSize int) (i int64, err error) ***REMOVED***
	s := unsafeBytesToString(b)
	return strconv.ParseInt(s, base, bitSize)
***REMOVED***

// parseUintBytes is a zero-alloc wrapper around strconv.ParseUint.
func parseUintBytes(b []byte, base int, bitSize int) (i uint64, err error) ***REMOVED***
	s := unsafeBytesToString(b)
	return strconv.ParseUint(s, base, bitSize)
***REMOVED***

// parseFloatBytes is a zero-alloc wrapper around strconv.ParseFloat.
func parseFloatBytes(b []byte, bitSize int) (float64, error) ***REMOVED***
	s := unsafeBytesToString(b)
	return strconv.ParseFloat(s, bitSize)
***REMOVED***

// parseBoolBytes is a zero-alloc wrapper around strconv.ParseBool.
func parseBoolBytes(b []byte) (bool, error) ***REMOVED***
	return strconv.ParseBool(unsafeBytesToString(b))
***REMOVED***

// unsafeBytesToString converts a []byte to a string without a heap allocation.
//
// It is unsafe, and is intended to prepare input to short-lived functions
// that require strings.
func unsafeBytesToString(in []byte) string ***REMOVED***
	src := *(*reflect.SliceHeader)(unsafe.Pointer(&in))
	dst := reflect.StringHeader***REMOVED***
		Data: src.Data,
		Len:  src.Len,
	***REMOVED***
	s := *(*string)(unsafe.Pointer(&dst))
	return s
***REMOVED***
