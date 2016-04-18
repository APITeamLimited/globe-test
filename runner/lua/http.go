package lua

import (
	"github.com/loadimpact/speedboat/runner"
	"github.com/valyala/fasthttp"
	"github.com/yuin/gopher-lua"
	"time"
)

func (vu *VUContext) HTTPLoader(L *lua.LState) int ***REMOVED***
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction***REMOVED***
		"get": vu.HTTPGet,
	***REMOVED***)
	L.SetField(mod, "name", lua.LString("http"))
	L.Push(mod)
	return 1
***REMOVED***

func (vu *VUContext) HTTPGet(L *lua.LState) int ***REMOVED***
	result := make(chan runner.Result, 1)
	go func() ***REMOVED***
		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)

		res := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(res)

		url := L.ToString(1)
		req.SetRequestURI(url)

		startTime := time.Now()
		err := vu.r.Client.Do(req, res)
		duration := time.Since(startTime)

		result <- runner.Result***REMOVED***Error: err, Time: duration***REMOVED***
	***REMOVED***()

	select ***REMOVED***
	case <-vu.ctx.Done():
	case res := <-result:
		vu.ch <- res
	***REMOVED***

	return 0
***REMOVED***
