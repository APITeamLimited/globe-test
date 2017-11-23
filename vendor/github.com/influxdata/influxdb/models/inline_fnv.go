package models

// from stdlib hash/fnv/fnv.go
const (
	prime64  = 1099511628211
	offset64 = 14695981039346656037
)

// InlineFNV64a is an alloc-free port of the standard library's fnv64a.
// See https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function.
type InlineFNV64a uint64

// NewInlineFNV64a returns a new instance of InlineFNV64a.
func NewInlineFNV64a() InlineFNV64a ***REMOVED***
	return offset64
***REMOVED***

// Write adds data to the running hash.
func (s *InlineFNV64a) Write(data []byte) (int, error) ***REMOVED***
	hash := uint64(*s)
	for _, c := range data ***REMOVED***
		hash ^= uint64(c)
		hash *= prime64
	***REMOVED***
	*s = InlineFNV64a(hash)
	return len(data), nil
***REMOVED***

// Sum64 returns the uint64 of the current resulting hash.
func (s *InlineFNV64a) Sum64() uint64 ***REMOVED***
	return uint64(*s)
***REMOVED***
