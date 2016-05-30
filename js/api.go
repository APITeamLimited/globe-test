package js

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/olebedev/go-duktape.v2"
)

type apiFunc func(r *Runner, c *duktape.Context) int

func apiHTTPDo(r *Runner, c *duktape.Context) int ***REMOVED***
	method := argString(c, 0)
	if method == "" ***REMOVED***
		log.Error("Missing method in http call")
		return 0
	***REMOVED***

	url := argString(c, 1)
	if url == "" ***REMOVED***
		log.Error("Missing URL in http call")
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
		log.Error("Unknown type for request body")
		return 0
	***REMOVED***

	args := httpArgs***REMOVED******REMOVED***
	if err := argJSON(c, 3, &args); err != nil ***REMOVED***
		log.Error("Invalid arguments to http call")
		return 0
	***REMOVED***

	res, duration, err := httpDo(r.Client, method, url, body, args)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Request error")
	***REMOVED***
	if !args.Quiet ***REMOVED***
		r.mDuration.Update(duration.Nanoseconds())
	***REMOVED***

	pushInstance(c, res, "HTTPResponse")

	return 1
***REMOVED***

func apiHTTPSetMaxConnectionsPerHost(r *Runner, c *duktape.Context) int ***REMOVED***
	num := int(argNumber(c, 0))
	if num < 1 ***REMOVED***
		log.Error("Max connections per host must be at least 1")
		return 0
	***REMOVED***
	r.Client.MaxConnsPerHost = num
	return 0
***REMOVED***

func apiLogType(r *Runner, c *duktape.Context) int ***REMOVED***
	kind := argString(c, 0)
	text := argString(c, 1)
	extra := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	if err := argJSON(c, 2, &extra); err != nil ***REMOVED***
		log.Error("Log context is not an object")
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

func apiTestAbort(r *Runner, c *duktape.Context) int ***REMOVED***
	// TODO: Do this some better way.
	log.Fatal("Test aborted")
	return 0
***REMOVED***
