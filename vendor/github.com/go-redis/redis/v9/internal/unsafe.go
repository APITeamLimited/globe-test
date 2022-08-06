//go:build !appengine
// +build !appengine

package internal

import "unsafe"

// String converts byte slice to string.
func String(b []byte) string ***REMOVED***
	return *(*string)(unsafe.Pointer(&b))
***REMOVED***

// Bytes converts string to byte slice.
func Bytes(s string) []byte ***REMOVED***
	return *(*[]byte)(unsafe.Pointer(
		&struct ***REMOVED***
			string
			Cap int
		***REMOVED******REMOVED***s, len(s)***REMOVED***,
	))
***REMOVED***
