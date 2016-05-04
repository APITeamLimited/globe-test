package http

import (
	"github.com/valyala/fasthttp"
	"math"
	"time"
)

type context struct ***REMOVED***
	client   *fasthttp.Client
	defaults RequestArgs
***REMOVED***

type RequestArgs struct ***REMOVED***
	Follow    bool   `json:"follow"`
	Report    bool   `json:"report"`
	UserAgent string `json:"userAgent"`
***REMOVED***

func (args *RequestArgs) ApplyDefaults(def RequestArgs) ***REMOVED***
	if !args.Follow && def.Follow ***REMOVED***
		args.Follow = true
	***REMOVED***
	if !args.Report && def.Follow ***REMOVED***
		args.Report = true
	***REMOVED***
	if args.UserAgent == "" ***REMOVED***
		args.UserAgent = def.UserAgent
	***REMOVED***
***REMOVED***

func New() map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	ctx := &context***REMOVED***
		client: &fasthttp.Client***REMOVED***
			Dial:                fasthttp.Dial,
			MaxIdleConnDuration: time.Duration(0),
			MaxConnsPerHost:     math.MaxInt64,
		***REMOVED***,
	***REMOVED***
	return map[string]interface***REMOVED******REMOVED******REMOVED***
		"get":                ctx.Get,
		"head":               ctx.Head,
		"post":               ctx.Post,
		"put":                ctx.Put,
		"delete":             ctx.Delete,
		"request":            ctx.Request,
		"setMaxConnsPerHost": ctx.SetMaxConnsPerHost,
	***REMOVED***
***REMOVED***
