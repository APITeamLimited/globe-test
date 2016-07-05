package js

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/lib"
	"github.com/loadimpact/speedboat/stats"
	"github.com/robertkrimen/otto"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"math"
	"os"
)

type Runner struct ***REMOVED***
	filename string
	source   string

	logger *log.Logger
***REMOVED***

type VU struct ***REMOVED***
	Runner *Runner
	VM     *otto.Otto
	Script *otto.Script

	Collector *stats.Collector

	Client fasthttp.Client

	ID        int64
	Iteration int64
***REMOVED***

func New(filename, source string) *Runner ***REMOVED***
	return &Runner***REMOVED***
		filename: filename,
		source:   source,
		logger: &log.Logger***REMOVED***
			Out:       os.Stderr,
			Formatter: &log.TextFormatter***REMOVED******REMOVED***,
			Hooks:     make(log.LevelHooks),
			Level:     log.DebugLevel,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (r *Runner) NewVU() (lib.VU, error) ***REMOVED***
	vuVM := otto.New()

	script, err := vuVM.Compile(r.filename, r.source)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	vu := VU***REMOVED***
		Runner: r,
		VM:     vuVM,
		Script: script,

		Collector: stats.NewCollector(),
	***REMOVED***

	vu.VM.Set("print", func(call otto.FunctionCall) otto.Value ***REMOVED***
		fmt.Fprintln(os.Stderr, call.Argument(0))
		return otto.UndefinedValue()
	***REMOVED***)

	vu.VM.Set("$http", map[string]interface***REMOVED******REMOVED******REMOVED***
		"request": func(call otto.FunctionCall) otto.Value ***REMOVED***
			method, _ := call.Argument(0).ToString()
			url, _ := call.Argument(1).ToString()

			body, err := bodyFromValue(call.Argument(2))
			if err != nil ***REMOVED***
				panic(call.Otto.MakeTypeError("invalid body"))
			***REMOVED***

			params, err := paramsFromObject(call.Argument(3).Object())
			if err != nil ***REMOVED***
				panic(err)
			***REMOVED***

			log.WithFields(log.Fields***REMOVED***
				"method": method,
				"url":    url,
				"body":   body,
				"params": params,
			***REMOVED***).Debug("Request")
			res, err := vu.HTTPRequest(method, url, body, params)
			if err != nil ***REMOVED***
				panic(jsCustomError(call.Otto, "HTTPError", err))
			***REMOVED***

			val, err := res.ToValue(call.Otto)
			if err != nil ***REMOVED***
				panic(jsError(call.Otto, err))
			***REMOVED***

			return val
		***REMOVED***,
		"setMaxConnsPerHost": func(call otto.FunctionCall) otto.Value ***REMOVED***
			num, err := call.Argument(0).ToInteger()
			if err != nil ***REMOVED***
				panic(call.Otto.MakeTypeError("argument must be an integer"))
			***REMOVED***
			if num <= 0 ***REMOVED***
				panic(call.Otto.MakeRangeError("argument must be >= 1"))
			***REMOVED***
			if num > math.MaxInt32 ***REMOVED***
				num = math.MaxInt32
			***REMOVED***

			vu.Client.MaxConnsPerHost = int(num)

			return otto.UndefinedValue()
		***REMOVED***,
	***REMOVED***)
	vu.VM.Set("$vu", map[string]interface***REMOVED******REMOVED******REMOVED***
		"sleep": func(call otto.FunctionCall) otto.Value ***REMOVED***
			t, err := call.Argument(0).ToFloat()
			if err != nil ***REMOVED***
				panic(call.Otto.MakeTypeError("time must be a number"))
			***REMOVED***

			vu.Sleep(t)

			return otto.UndefinedValue()
		***REMOVED***,
		"id": func(call otto.FunctionCall) otto.Value ***REMOVED***
			val, err := call.Otto.ToValue(vu.ID)
			if err != nil ***REMOVED***
				panic(jsError(call.Otto, err))
			***REMOVED***
			return val
		***REMOVED***,
		"iteration": func(call otto.FunctionCall) otto.Value ***REMOVED***
			val, err := call.Otto.ToValue(vu.Iteration)
			if err != nil ***REMOVED***
				panic(jsError(call.Otto, err))
			***REMOVED***
			return val
		***REMOVED***,
	***REMOVED***)
	vu.VM.Set("$test", map[string]interface***REMOVED******REMOVED******REMOVED***
		"env": func(call otto.FunctionCall) otto.Value ***REMOVED***
			key, _ := call.Argument(0).ToString()

			value, ok := os.LookupEnv(key)
			if !ok ***REMOVED***
				return otto.UndefinedValue()
			***REMOVED***

			val, err := call.Otto.ToValue(value)
			if err != nil ***REMOVED***
				panic(jsError(call.Otto, err))
			***REMOVED***
			return val
		***REMOVED***,
		"abort": func(call otto.FunctionCall) otto.Value ***REMOVED***
			panic(lib.AbortTest)
			return otto.UndefinedValue()
		***REMOVED***,
	***REMOVED***)
	vu.VM.Set("$log", map[string]interface***REMOVED******REMOVED******REMOVED***
		"log": func(call otto.FunctionCall) otto.Value ***REMOVED***
			level, _ := call.Argument(0).ToString()
			msg, _ := call.Argument(1).ToString()

			fields := make(map[string]interface***REMOVED******REMOVED***)
			fieldsObj := call.Argument(2).Object()
			if fieldsObj != nil ***REMOVED***
				for _, key := range fieldsObj.Keys() ***REMOVED***
					valObj, _ := fieldsObj.Get(key)
					val, err := valObj.Export()
					if err != nil ***REMOVED***
						panic(jsError(call.Otto, err))
					***REMOVED***
					fields[key] = val
				***REMOVED***
			***REMOVED***

			vu.Log(level, msg, fields)

			return otto.UndefinedValue()
		***REMOVED***,
	***REMOVED***)

	init := `
	function HTTPResponse() ***REMOVED***
		this.json = function() ***REMOVED***
			return JSON.parse(this.body);
		***REMOVED***;
	***REMOVED***
	
	$http.get = function(url, data, params) ***REMOVED*** return $http.request('GET', url, data, params); ***REMOVED***;
	$http.post = function(url, data, params) ***REMOVED*** return $http.request('POST', url, data, params); ***REMOVED***;
	$http.put = function(url, data, params) ***REMOVED*** return $http.request('PUT', url, data, params); ***REMOVED***;
	$http.delete = function(url, data, params) ***REMOVED*** return $http.request('DELETE', url, data, params); ***REMOVED***;
	$http.patch = function(url, data, params) ***REMOVED*** return $http.request('PATCH', url, data, params); ***REMOVED***;
	$http.options = function(url, data, params) ***REMOVED*** return $http.request('OPTIONS', url, data, params); ***REMOVED***;
	$http.head = function(url, data, params) ***REMOVED*** return $http.request('HEAD', url, data, params); ***REMOVED***;
	
	$log.debug = function(msg, fields) ***REMOVED*** $log.log('debug', msg, fields); ***REMOVED***;
	$log.info = function(msg, fields) ***REMOVED*** $log.log('info', msg, fields); ***REMOVED***;
	$log.warn = function(msg, fields) ***REMOVED*** $log.log('warn', msg, fields); ***REMOVED***;
	$log.error = function(msg, fields) ***REMOVED*** $log.log('error', msg, fields); ***REMOVED***;
	`
	if _, err := vu.VM.Eval(init); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &vu, nil
***REMOVED***

func (u *VU) Reconfigure(id int64) error ***REMOVED***
	u.ID = id
	u.Iteration = 0
	return nil
***REMOVED***

func (u *VU) RunOnce(ctx context.Context) error ***REMOVED***
	u.Iteration++
	if _, err := u.VM.Run(u.Script); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***
