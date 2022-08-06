//go:build !appengine
// +build !appengine

package util

import (
	"unsafe"
)

// BytesToString converts byte slice to string.
func BytesToString(b []byte) string ***REMOVED***
	return *(*string)(unsafe.Pointer(&b))
***REMOVED***

// StringToBytes converts string to byte slice.
func StringToBytes(s string) []byte ***REMOVED***
	return *(*[]byte)(unsafe.Pointer(
		&struct ***REMOVED***
			string
			Cap int
		***REMOVED******REMOVED***s, len(s)***REMOVED***,
	))
***REMOVED***
