package http

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/runner"
	"github.com/valyala/fasthttp"
	"math"
	"time"
)

type context struct ***REMOVED***
	client *fasthttp.Client
***REMOVED***

type RequestArgs struct ***REMOVED***
	Follow bool `json:"follow"`
	Report bool `json:"report"`
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
		"get":     ctx.Get,
		"head":    ctx.Head,
		"post":    ctx.Post,
		"put":     ctx.Put,
		"delete":  ctx.Delete,
		"request": ctx.Request,
	***REMOVED***
***REMOVED***

func (ctx *context) Get(url string, args RequestArgs) <-chan runner.Result ***REMOVED***
	return ctx.Request("GET", url, "", args)
***REMOVED***

func (ctx *context) Head(url string, args RequestArgs) <-chan runner.Result ***REMOVED***
	return ctx.Request("HEAD", url, "", args)
***REMOVED***

func (ctx *context) Post(url, body string, args RequestArgs) <-chan runner.Result ***REMOVED***
	return ctx.Request("POST", url, body, args)
***REMOVED***

func (ctx *context) Put(url, body string, args RequestArgs) <-chan runner.Result ***REMOVED***
	return ctx.Request("PUT", url, body, args)
***REMOVED***

func (ctx *context) Delete(url, body string, args RequestArgs) <-chan runner.Result ***REMOVED***
	return ctx.Request("DELETE", url, body, args)
***REMOVED***

func (ctx *context) Request(method, url, body string, args RequestArgs) <-chan runner.Result ***REMOVED***
	log.WithFields(log.Fields***REMOVED***
		"method": method,
		"url":    url,
		"follow": args.Follow,
		"report": args.Report,
	***REMOVED***).Debug("Request")
	ch := make(chan runner.Result, 1)
	go func() ***REMOVED***
		defer close(ch)

		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)

		res := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(res)

		req.SetRequestURI(url)
		req.Header.SetMethod(method)
		req.SetBodyString(body)

		startTime := time.Now()
		err := ctx.client.Do(req, res)
		duration := time.Since(startTime)

		ch <- runner.Result***REMOVED***Error: err, Time: duration***REMOVED***
	***REMOVED***()
	return ch
***REMOVED***
