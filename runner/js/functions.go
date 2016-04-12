package js

import (
	"fmt"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"net/http"
	"time"
)

type JSError string

func (e JSError) Error() string ***REMOVED*** return string(e) ***REMOVED***

func jsSleepFactory(impl func(time.Duration)) func(otto.FunctionCall) otto.Value ***REMOVED***
	return func(call otto.FunctionCall) otto.Value ***REMOVED***
		seconds, err := call.Argument(0).ToFloat()
		if err != nil ***REMOVED***
			seconds = 0.0
		***REMOVED***
		impl(time.Duration(seconds * float64(time.Second)))
		return otto.UndefinedValue()
	***REMOVED***
***REMOVED***

func jsLogFactory(impl func(string)) func(otto.FunctionCall) otto.Value ***REMOVED***
	return func(call otto.FunctionCall) otto.Value ***REMOVED***
		text, err := call.Argument(0).ToString()
		if err != nil ***REMOVED***
			text = "[ERROR]"
		***REMOVED***
		impl(text)
		return otto.UndefinedValue()
	***REMOVED***
***REMOVED***

func jsHTTPGetFactory(vm *otto.Otto, impl func(url string) (*http.Response, error)) func(otto.FunctionCall) otto.Value ***REMOVED***
	return func(call otto.FunctionCall) otto.Value ***REMOVED***
		url, err := call.Argument(0).ToString()
		if err != nil ***REMOVED***
			panic(JSError(fmt.Sprintf("Couldn't call function: %s", err)))
		***REMOVED***

		res, err := impl(url)
		if err != nil ***REMOVED***
			panic(JSError(fmt.Sprintf("HTTP GET impl error: %s", err)))
		***REMOVED***
		defer res.Body.Close()

		obj, err := vm.Object("new Object()")
		if err != nil ***REMOVED***
			panic(JSError(fmt.Sprintf("Couldn't create an Object(): %s", err)))
		***REMOVED***
		body, _ := ioutil.ReadAll(res.Body)
		obj.Set("body", string(body))
		obj.Set("statusCode", res.StatusCode)
		obj.Set("header", res.Header)

		return obj.Value()
	***REMOVED***
***REMOVED***
