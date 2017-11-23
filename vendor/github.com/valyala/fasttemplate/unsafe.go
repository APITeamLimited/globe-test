// +build !appengine

package fasttemplate

import (
	"reflect"
	"unsafe"
)

func unsafeBytes2String(b []byte) string ***REMOVED***
	return *(*string)(unsafe.Pointer(&b))
***REMOVED***

func unsafeString2Bytes(s string) []byte ***REMOVED***
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader***REMOVED***
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	***REMOVED***
	return *(*[]byte)(unsafe.Pointer(&bh))
***REMOVED***
