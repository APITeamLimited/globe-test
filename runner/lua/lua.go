package lua

import (
	"github.com/loadimpact/speedboat/runner"
	"github.com/yuin/gopher-lua"
	"golang.org/x/net/context"
)

type LuaRunner struct ***REMOVED***
	Filename, Source string
***REMOVED***

func New(filename, src string) *LuaRunner ***REMOVED***
	return &LuaRunner***REMOVED***
		Filename: filename,
		Source:   src,
	***REMOVED***
***REMOVED***

func (r *LuaRunner) Run(ctx context.Context) <-chan runner.Result ***REMOVED***
	ch := make(chan runner.Result)

	go func() ***REMOVED***
		defer close(ch)

		L := lua.NewState()
		defer L.Close()

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
