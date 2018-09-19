//+build !appengine

package gjson

import (
	"reflect"
	"unsafe"
)

// getBytes casts the input json bytes to a string and safely returns the
// results as uniquely allocated data. This operation is intended to minimize
// copies and allocations for the large json string->[]byte.
func getBytes(json []byte, path string) Result ***REMOVED***
	var result Result
	if json != nil ***REMOVED***
		// unsafe cast to string
		result = Get(*(*string)(unsafe.Pointer(&json)), path)
		result = fromBytesGet(result)
	***REMOVED***
	return result
***REMOVED***

func fromBytesGet(result Result) Result ***REMOVED***
	// safely get the string headers
	rawhi := *(*reflect.StringHeader)(unsafe.Pointer(&result.Raw))
	strhi := *(*reflect.StringHeader)(unsafe.Pointer(&result.Str))
	// create byte slice headers
	rawh := reflect.SliceHeader***REMOVED***Data: rawhi.Data, Len: rawhi.Len***REMOVED***
	strh := reflect.SliceHeader***REMOVED***Data: strhi.Data, Len: strhi.Len***REMOVED***
	if strh.Data == 0 ***REMOVED***
		// str is nil
		if rawh.Data == 0 ***REMOVED***
			// raw is nil
			result.Raw = ""
		***REMOVED*** else ***REMOVED***
			// raw has data, safely copy the slice header to a string
			result.Raw = string(*(*[]byte)(unsafe.Pointer(&rawh)))
		***REMOVED***
		result.Str = ""
	***REMOVED*** else if rawh.Data == 0 ***REMOVED***
		// raw is nil
		result.Raw = ""
		// str has data, safely copy the slice header to a string
		result.Str = string(*(*[]byte)(unsafe.Pointer(&strh)))
	***REMOVED*** else if strh.Data >= rawh.Data &&
		int(strh.Data)+strh.Len <= int(rawh.Data)+rawh.Len ***REMOVED***
		// Str is a substring of Raw.
		start := int(strh.Data - rawh.Data)
		// safely copy the raw slice header
		result.Raw = string(*(*[]byte)(unsafe.Pointer(&rawh)))
		// substring the raw
		result.Str = result.Raw[start : start+strh.Len]
	***REMOVED*** else ***REMOVED***
		// safely copy both the raw and str slice headers to strings
		result.Raw = string(*(*[]byte)(unsafe.Pointer(&rawh)))
		result.Str = string(*(*[]byte)(unsafe.Pointer(&strh)))
	***REMOVED***
	return result
***REMOVED***

// fillIndex finds the position of Raw data and assigns it to the Index field
// of the resulting value. If the position cannot be found then Index zero is
// used instead.
func fillIndex(json string, c *parseContext) ***REMOVED***
	if len(c.value.Raw) > 0 && !c.calcd ***REMOVED***
		jhdr := *(*reflect.StringHeader)(unsafe.Pointer(&json))
		rhdr := *(*reflect.StringHeader)(unsafe.Pointer(&(c.value.Raw)))
		c.value.Index = int(rhdr.Data - jhdr.Data)
		if c.value.Index < 0 || c.value.Index >= len(json) ***REMOVED***
			c.value.Index = 0
		***REMOVED***
	***REMOVED***
***REMOVED***
