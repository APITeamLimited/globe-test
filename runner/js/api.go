package js

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/runner"
	"gopkg.in/olebedev/go-duktape.v2"
)

type apiFunc func(r *Runner, c *duktape.Context, ch chan<- runner.Result) int

func apiHTTPDo(r *Runner, c *duktape.Context, ch chan<- runner.Result) int ***REMOVED***
	method := argString(c, 0)
	if method == "" ***REMOVED***
		ch <- runner.Result***REMOVED***Error: errors.New("Missing method in http call")***REMOVED***
		return 0
	***REMOVED***

	url := argString(c, 1)
	if url == "" ***REMOVED***
		ch <- runner.Result***REMOVED***Error: errors.New("Missing URL in http call")***REMOVED***
		return 0
	***REMOVED***

	body := ""
	switch c.GetType(2) ***REMOVED***
	case duktape.TypeNone, duktape.TypeNull, duktape.TypeUndefined:
	case duktape.TypeString, duktape.TypeNumber, duktape.TypeBoolean:
		body = c.ToString(2)
	case duktape.TypeObject:
		body = c.JsonEncode(2)
	default:
		ch <- runner.Result***REMOVED***Error: errors.New("Unknown type for request body")***REMOVED***
		return 0
	***REMOVED***

	args := httpArgs***REMOVED******REMOVED***
	if err := argJSON(c, 3, &args); err != nil ***REMOVED***
		ch <- runner.Result***REMOVED***Error: errors.New("Invalid arguments to http call")***REMOVED***
		return 0
	***REMOVED***

	res, duration, err := httpDo(r.Client, method, url, body, args)
	if !args.Quiet ***REMOVED***
		ch <- runner.Result***REMOVED***Error: err, Time: duration***REMOVED***
	***REMOVED***

	pushInstance(c, res, "HTTPResponse")

	return 1
***REMOVED***

func apiHTTPSetMaxConnectionsPerHost(r *Runner, c *duktape.Context, ch chan<- runner.Result) int ***REMOVED***
	num := int(argNumber(c, 0))
	if num < 1 ***REMOVED***
		ch <- runner.Result***REMOVED***Error: errors.New("Max connections per host must be at least 1")***REMOVED***
		return 0
	***REMOVED***
	r.Client.MaxConnsPerHost = num
	return 0
***REMOVED***

func apiLogType(r *Runner, c *duktape.Context, ch chan<- runner.Result) int ***REMOVED***
	kind := argString(c, 0)
	text := argString(c, 1)
	extra := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	if err := argJSON(c, 2, &extra); err != nil ***REMOVED***
		ch <- runner.Result***REMOVED***Error: errors.New("Log context is not an object")***REMOVED***
		return 0
	***REMOVED***

	l := log.WithFields(log.Fields(extra))
	switch kind ***REMOVED***
	case "debug":
		l.Debug(text)
	case "info":
		l.Info(text)
	case "warn":
		l.Warn(text)
	case "error":
		l.Error(text)
	***REMOVED***

	return 0
***REMOVED***

func apiTestAbort(r *Runner, c *duktape.Context, ch chan<- runner.Result) int ***REMOVED***
	ch <- runner.Result***REMOVED***Abort: true***REMOVED***
	return 0
***REMOVED***
