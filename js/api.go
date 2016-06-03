package js

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/js/http"
	"golang.org/x/net/context"
	"gopkg.in/olebedev/go-duktape.v2"
	"time"
)

type APIFunc func(js *duktape.Context, ctx context.Context) int

func apiSleep(js *duktape.Context, ctx context.Context) int ***REMOVED***
	time.Sleep(time.Duration(argNumber(js, 0) * float64(time.Second)))
	return 0
***REMOVED***

func apiHTTPDo(js *duktape.Context, ctx context.Context) int ***REMOVED***
	method := argString(js, 0)
	if method == "" ***REMOVED***
		log.Error("Missing method in http call")
		return 0
	***REMOVED***

	url := argString(js, 1)
	if url == "" ***REMOVED***
		log.Error("Missing URL in http call")
		return 0
	***REMOVED***

	body := ""
	switch js.GetType(2) ***REMOVED***
	case duktape.TypeNone, duktape.TypeNull, duktape.TypeUndefined:
	case duktape.TypeString, duktape.TypeNumber, duktape.TypeBoolean:
		body = js.ToString(2)
	case duktape.TypeObject:
		body = js.JsonEncode(2)
	default:
		log.Error("Unknown type for request body")
		return 0
	***REMOVED***

	args := http.Args***REMOVED******REMOVED***
	if err := argJSON(js, 3, &args); err != nil ***REMOVED***
		log.Error("Invalid arguments to http call")
		return 0
	***REMOVED***

	res, err := http.Do(ctx, method, url, body, args)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Request error")
	***REMOVED***

	pushObject(js, res, "HTTPResponse")

	return 1
***REMOVED***
