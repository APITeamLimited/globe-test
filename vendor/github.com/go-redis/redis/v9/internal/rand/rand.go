package rand

import (
	"math/rand"
	"sync"
)

// Int returns a non-negative pseudo-random int.
func Int() int ***REMOVED*** return pseudo.Int() ***REMOVED***

// Intn returns, as an int, a non-negative pseudo-random number in [0,n).
// It panics if n <= 0.
func Intn(n int) int ***REMOVED*** return pseudo.Intn(n) ***REMOVED***

// Int63n returns, as an int64, a non-negative pseudo-random number in [0,n).
// It panics if n <= 0.
func Int63n(n int64) int64 ***REMOVED*** return pseudo.Int63n(n) ***REMOVED***

// Perm returns, as a slice of n ints, a pseudo-random permutation of the integers [0,n).
func Perm(n int) []int ***REMOVED*** return pseudo.Perm(n) ***REMOVED***

// Seed uses the provided seed value to initialize the default Source to a
// deterministic state. If Seed is not called, the generator behaves as if
// seeded by Seed(1).
func Seed(n int64) ***REMOVED*** pseudo.Seed(n) ***REMOVED***

var pseudo = rand.New(&source***REMOVED***src: rand.NewSource(1)***REMOVED***)

type source struct ***REMOVED***
	src rand.Source
	mu  sync.Mutex
***REMOVED***

func (s *source) Int63() int64 ***REMOVED***
	s.mu.Lock()
	n := s.src.Int63()
	s.mu.Unlock()
	return n
***REMOVED***

func (s *source) Seed(seed int64) ***REMOVED***
	s.mu.Lock()
	s.src.Seed(seed)
	s.mu.Unlock()
***REMOVED***

// Shuffle pseudo-randomizes the order of elements.
// n is the number of elements.
// swap swaps the elements with indexes i and j.
func Shuffle(n int, swap func(i, j int)) ***REMOVED*** pseudo.Shuffle(n, swap) ***REMOVED***
