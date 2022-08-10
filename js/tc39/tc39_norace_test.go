//go:build !race
// +build !race

// Heavily influenced by the fantastic work by @dop251 for https://github.com/dop251/goja

package tc39

import "testing"

func (ctx *tc39TestCtx) runTest(name string, f func(t *testing.T)) ***REMOVED***
	ctx.t.Run(name, func(t *testing.T) ***REMOVED***
		t.Parallel()
		f(t)
	***REMOVED***)
***REMOVED***

func (ctx *tc39TestCtx) flush() ***REMOVED***
***REMOVED***
