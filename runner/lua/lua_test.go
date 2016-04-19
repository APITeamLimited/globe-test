package lua

import (
	"golang.org/x/net/context"
	"testing"
)

func BenchmarkRunEmpty(b *testing.B) ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	r := New("script.lua", "")
	for i := 0; i < b.N; i++ ***REMOVED***
		r.Run(ctx, int64(i))
	***REMOVED***
***REMOVED***
