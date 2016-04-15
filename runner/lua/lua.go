package lua

import (
	"github.com/loadimpact/speedboat/runner"
	"golang.org/x/net/context"
	// "github.com/Shopify/go-lua"
)

type LuaRunner struct ***REMOVED***
	Script string
***REMOVED***

func New(filename, src string) *LuaRunner ***REMOVED***
	return &LuaRunner***REMOVED******REMOVED***
***REMOVED***

func (r *LuaRunner) Run(ctx context.Context) <-chan runner.Result ***REMOVED***
	ch := make(chan runner.Result)

	go func() ***REMOVED***
		defer close(ch)
	***REMOVED***()

	return ch
***REMOVED***
