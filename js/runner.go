package js

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat"
	"github.com/loadimpact/speedboat/sampler"
	"github.com/robertkrimen/otto"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"math"
	"os"
)

type Runner struct ***REMOVED***
	Test speedboat.Test

	filename string
	source   string

	logger *log.Logger

	mDuration *sampler.Metric
	mErrors   *sampler.Metric
***REMOVED***

type VU struct ***REMOVED***
	Runner *Runner
	VM     *otto.Otto
	Script *otto.Script

	Client fasthttp.Client

	ID        int64
	Iteration int64
***REMOVED***

func New(t speedboat.Test, filename, source string) *Runner ***REMOVED***
	return &Runner***REMOVED***
		Test:     t,
		filename: filename,
		source:   source,
		logger: &log.Logger***REMOVED***
			Out:       os.Stderr,
			Formatter: &log.TextFormatter***REMOVED******REMOVED***,
			Hooks:     make(log.LevelHooks),
			Level:     log.DebugLevel,
		***REMOVED***,
		mDuration: sampler.Stats("request.duration"),
		mErrors:   sampler.Counter("request.error"),
	***REMOVED***
***REMOVED***

func (r *Runner) NewVU() (speedboat.VU, error) ***REMOVED***
	vm := otto.New()

	script, err := vm.Compile(r.filename, r.source)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	vu := VU***REMOVED***
		Runner: r,
		VM:     vm,
		Script: script,
	***REMOVED***

	vm.Set("print", func(call otto.FunctionCall) otto.Value ***REMOVED***
		fmt.Fprintln(os.Stderr, call.Argument(0))
		return otto.UndefinedValue()
	***REMOVED***)

	vm.Set("$http", map[string]interface***REMOVED******REMOVED******REMOVED***
		"request": func(call otto.FunctionCall) otto.Value ***REMOVED***
			method, err := call.Argument(0).ToString()
			if err != nil ***REMOVED***
				panic(vm.MakeTypeError("method must be a string"))
			***REMOVED***

			url, err := call.Argument(1).ToString()
			if err != nil ***REMOVED***
				panic(vm.MakeTypeError("url must be a string"))
			***REMOVED***

			body, err := bodyFromValue(call.Argument(2))
			if err != nil ***REMOVED***
				panic(vm.MakeTypeError("invalid body"))
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
				panic(vm.MakeCustomError("HTTPError", err.Error()))
			***REMOVED***

			val, err := res.ToValue(vm)
			if err != nil ***REMOVED***
				panic(vm.MakeCustomError("Error", err.Error()))
			***REMOVED***

			return val
		***REMOVED***,
		"setMaxConnsPerHost": func(call otto.FunctionCall) otto.Value ***REMOVED***
			num, err := call.Argument(0).ToInteger()
			if err != nil ***REMOVED***
				panic(vm.MakeTypeError("argument must be an integer"))
			***REMOVED***
			if num <= 0 ***REMOVED***
				panic(vm.MakeRangeError("argument must be >= 1"))
			***REMOVED***
			if num > math.MaxInt32 ***REMOVED***
				num = math.MaxInt32
			***REMOVED***

			vu.Client.MaxConnsPerHost = int(num)

			return otto.UndefinedValue()
		***REMOVED***,
	***REMOVED***)
	vm.Set("$vu", map[string]interface***REMOVED******REMOVED******REMOVED***
		"sleep": func(call otto.FunctionCall) otto.Value ***REMOVED***
			t, err := call.Argument(0).ToFloat()
			if err != nil ***REMOVED***
				panic(vm.MakeTypeError("time must be a number"))
			***REMOVED***

			vu.Sleep(t)

			return otto.UndefinedValue()
		***REMOVED***,
		"id": func(call otto.FunctionCall) otto.Value ***REMOVED***
			val, err := call.Otto.ToValue(vu.ID)
			if err != nil ***REMOVED***
				panic(jsError(vm, err))
			***REMOVED***
			return val
		***REMOVED***,
		"iteration": func(call otto.FunctionCall) otto.Value ***REMOVED***
			val, err := call.Otto.ToValue(vu.Iteration)
			if err != nil ***REMOVED***
				panic(jsError(vm, err))
			***REMOVED***
			return val
		***REMOVED***,
	***REMOVED***)
	vm.Set("$log", map[string]interface***REMOVED******REMOVED******REMOVED***
		"log": func(call otto.FunctionCall) otto.Value ***REMOVED***
			level, err := call.Argument(0).ToString()
			if err != nil ***REMOVED***
				panic(vm.MakeTypeError("level must be a string"))
			***REMOVED***

			msg, err := call.Argument(1).ToString()
			if err != nil ***REMOVED***
				panic(vm.MakeTypeError("message must be a string"))
			***REMOVED***

			fields := make(map[string]interface***REMOVED******REMOVED***)
			fieldsObj := call.Argument(2).Object()
			if fieldsObj != nil ***REMOVED***
				for _, key := range fieldsObj.Keys() ***REMOVED***
					valObj, _ := fieldsObj.Get(key)
					val, err := valObj.Export()
					if err != nil ***REMOVED***
						panic(err)
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
	if _, err := vm.Eval(init); err != nil ***REMOVED***
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
