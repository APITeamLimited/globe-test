package js

import (
	"testing"
)

func BenchmarkRunIteration(b *testing.B) ***REMOVED***
	r, err := New()
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	err = r.Load("script.js", "")
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		for res := range r.RunIteration() ***REMOVED***
			if err, ok := res.(error); ok ***REMOVED***
				b.Error(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
