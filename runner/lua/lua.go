package lua

import (
	"github.com/loadimpact/speedboat/runner"
	"github.com/valyala/fasthttp"
	"github.com/yuin/gopher-lua"
	"golang.org/x/net/context"
	"time"
)

type LuaRunner struct ***REMOVED***
	Filename string
	Source   string
	Client   *fasthttp.Client
***REMOVED***

type VUContext struct ***REMOVED***
	r   *LuaRunner
	ctx context.Context
	ch  chan runner.Result
***REMOVED***

func New(filename, src string) *LuaRunner ***REMOVED***
	return &LuaRunner***REMOVED***
		Filename: filename,
		Source:   src,
		Client: &fasthttp.Client***REMOVED***
			MaxIdleConnDuration: time.Duration(0),
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (r *LuaRunner) Run(ctx context.Context) <-chan runner.Result ***REMOVED***
	ch := make(chan runner.Result)

	go func() ***REMOVED***
		defer close(ch)

		vu := VUContext***REMOVED***r: r, ctx: ctx, ch: ch***REMOVED***

		L := lua.NewState()
		defer L.Close()

		L.SetGlobal("sleep", L.NewFunction(vu.Sleep))

		L.PreloadModule("http", vu.HTTPLoader)

		// Try to load the script, abort execution if it fails
		lfn, err := L.LoadString(r.Source)
		if err != nil ***REMOVED***
			ch <- runner.Result***REMOVED***Error: err***REMOVED***
			return
		***REMOVED***

		for ***REMOVED***
			L.Push(lfn)
			if err := L.PCall(0, 0, nil); err != nil ***REMOVED***
				ch <- runner.Result***REMOVED***Error: err***REMOVED***
			***REMOVED***

			select ***REMOVED***
			case <-ctx.Done():
				return
			default:
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***
