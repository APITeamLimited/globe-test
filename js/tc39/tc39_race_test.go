// +build race
// Heavily influenced by the fantastic work by @dop251 for https://github.com/dop251/goja

package tc39

import (
	"testing"
)

const (
	tc39MaxTestGroupSize = 1000 // to prevent race detector complaining about too many goroutines
)

func (ctx *tc39TestCtx) runTest(name string, f func(t *testing.T)) ***REMOVED***
	ctx.testQueue = append(ctx.testQueue, tc39Test***REMOVED***name: name, f: f***REMOVED***)
	if len(ctx.testQueue) >= tc39MaxTestGroupSize ***REMOVED***
		ctx.flush()
	***REMOVED***
***REMOVED***

func (ctx *tc39TestCtx) flush() ***REMOVED***
	ctx.t.Run("tc39", func(t *testing.T) ***REMOVED***
		for _, tc := range ctx.testQueue ***REMOVED***
			tc := tc
			t.Run(tc.name, func(t *testing.T) ***REMOVED***
				t.Parallel()
				tc.f(t)
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
	ctx.testQueue = ctx.testQueue[:0]
***REMOVED***
