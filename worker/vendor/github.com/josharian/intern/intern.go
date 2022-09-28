// Package intern interns strings.
// Interning is best effort only.
// Interned strings may be removed automatically
// at any time without notification.
// All functions may be called concurrently
// with themselves and each other.
package intern

import "sync"

var (
	pool sync.Pool = sync.Pool***REMOVED***
		New: func() interface***REMOVED******REMOVED*** ***REMOVED***
			return make(map[string]string)
		***REMOVED***,
	***REMOVED***
)

// String returns s, interned.
func String(s string) string ***REMOVED***
	m := pool.Get().(map[string]string)
	c, ok := m[s]
	if ok ***REMOVED***
		pool.Put(m)
		return c
	***REMOVED***
	m[s] = s
	pool.Put(m)
	return s
***REMOVED***

// Bytes returns b converted to a string, interned.
func Bytes(b []byte) string ***REMOVED***
	m := pool.Get().(map[string]string)
	c, ok := m[string(b)]
	if ok ***REMOVED***
		pool.Put(m)
		return c
	***REMOVED***
	s := string(b)
	m[s] = s
	pool.Put(m)
	return s
***REMOVED***
