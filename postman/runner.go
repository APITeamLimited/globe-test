package postman

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GeertJohan/go.rice"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/lib"
	"github.com/loadimpact/speedboat/stats"
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
	"golang.org/x/net/context"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"time"
)

var (
	mRequests = stats.Stat***REMOVED***Name: "requests", Type: stats.HistogramType, Intent: stats.TimeIntent***REMOVED***
	mErrors   = stats.Stat***REMOVED***Name: "errors", Type: stats.CounterType***REMOVED***
)

const vuSetupScript = `
	var globals = ***REMOVED******REMOVED***;
	var environment = ***REMOVED******REMOVED***;
	
	var postman = ***REMOVED******REMOVED***;
	
	postman.setEnvironmentVariable = function(name, value) ***REMOVED***
		environment[name] = value;
	***REMOVED***
	postman.setGlobalVariable = function(name, value) ***REMOVED***
		globals[name] = value;
	***REMOVED***
	postman.clearEnvironmentVariable = function(name) ***REMOVED***
		delete environment[name];
	***REMOVED***
	postman.clearGlobalVariable = function(name) ***REMOVED***
		delete globals[name];
	***REMOVED***
	postman.clearEnvironmentVariables = function() ***REMOVED***
		environment = ***REMOVED******REMOVED***;
	***REMOVED***
	postman.clearGlobalVariables = function() ***REMOVED***
		globals = ***REMOVED******REMOVED***;
	***REMOVED***
	
	postman.getResponseHeader = function(name) ***REMOVED***
		// Normalize captialization; "content-type"/"CONTENT-TYPE" -> "Content-Type"
		return responseHeaders[name.toLowerCase().replace(/(?:^|-)(\w)/g, function(txt) ***REMOVED***
			return txt.toUpperCase();
		***REMOVED***)];
	***REMOVED***
`

var libFiles = []string***REMOVED***
	"sugar/release/sugar.js",
***REMOVED***

var libPatches = map[string]map[string]string***REMOVED***
	"sugar/release/sugar.js": map[string]string***REMOVED***
		// Patch out functions using unsupported regex features.
		`function cleanDateInput(str) ***REMOVED***
      str = str.trim().replace(/^just (?=now)|\.+$/i, '');
      return convertAsianDigits(str);
    ***REMOVED***`: "",
		`function truncateOnWord(str, limit, fromLeft) ***REMOVED***
      if (fromLeft) ***REMOVED***
        return reverseString(truncateOnWord(reverseString(str), limit));
      ***REMOVED***
      var reg = RegExp('(?=[' + getTrimmableCharacters() + '])');
      var words = str.split(reg);
      var count = 0;
      return words.filter(function(word) ***REMOVED***
        count += word.length;
        return count <= limit;
      ***REMOVED***).join('');
    ***REMOVED***`: "",
		// We don't need to fully patch out this one, we just have to drop support for -昨 (last...)
		// This regex is only used to tell whether a character with multiple meanings is used as a
		// number or as a word, which is not something we're expecting people to do here anyways.
		`AsianDigitReg = RegExp('([期週周])?([' + KanjiDigits + FullWidthDigits + ']+)(?!昨)', 'g');`: `AsianDigitReg = RegExp('([期週周])?([' + KanjiDigits + FullWidthDigits + ']+)', 'g');`,
	***REMOVED***,
***REMOVED***

type ErrorWithLineNumber struct ***REMOVED***
	Wrapped error
	Line    int
***REMOVED***

func (e ErrorWithLineNumber) Error() string ***REMOVED***
	return fmt.Sprintf("%s (line %d)", e.Wrapped.Error(), e.Line)
***REMOVED***

type Runner struct ***REMOVED***
	VM         *otto.Otto
	Collection Collection
	Endpoints  []Endpoint
***REMOVED***

type VU struct ***REMOVED***
	Runner    *Runner
	VM        *otto.Otto
	Client    http.Client
	Collector *stats.Collector
	Iteration int64
***REMOVED***

func New(source []byte) (*Runner, error) ***REMOVED***
	var collection Collection
	if err := json.Unmarshal(source, &collection); err != nil ***REMOVED***
		switch e := err.(type) ***REMOVED***
		case *json.SyntaxError:
			src := string(source)
			line := strings.Count(src[:e.Offset], "\n") + 1
			return nil, ErrorWithLineNumber***REMOVED***Wrapped: e, Line: line***REMOVED***
		case *json.UnmarshalTypeError:
			src := string(source)
			line := strings.Count(src[:e.Offset], "\n") + 1
			return nil, ErrorWithLineNumber***REMOVED***Wrapped: e, Line: line***REMOVED***
		***REMOVED***
		return nil, err
	***REMOVED***

	vm := otto.New()
	lib, err := rice.FindBox("lib")
	if err != nil ***REMOVED***
		return nil, errors.New(fmt.Sprintf("couldn't find postman lib files; this can happen if you run from the wrong working directory with a non-boxed binary: %s", err.Error()))
	***REMOVED***
	for _, filename := range libFiles ***REMOVED***
		src, err := lib.String(filename)
		if err != nil ***REMOVED***
			return nil, errors.New(fmt.Sprintf("couldn't load lib file (%s): %s", filename, err.Error()))
		***REMOVED***
		for find, repl := range libPatches[filename] ***REMOVED***
			src = strings.Replace(src, find, repl, 1)
		***REMOVED***
		if _, err := vm.Eval(src); err != nil ***REMOVED***
			return nil, errors.New(fmt.Sprintf("couldn't eval lib file (%s): %s", filename, err.Error()))
		***REMOVED***
	***REMOVED***

	eps, err := MakeEndpoints(collection, vm)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &Runner***REMOVED***
		Collection: collection,
		Endpoints:  eps,
		VM:         vm,
	***REMOVED***, nil
***REMOVED***

func (r *Runner) NewVU() (lib.VU, error) ***REMOVED***
	return &VU***REMOVED***
		Runner: r,
		VM:     r.VM.Copy(),
		Client: http.Client***REMOVED***
			Transport: &http.Transport***REMOVED***
				MaxIdleConnsPerHost: math.MaxInt32,
			***REMOVED***,
		***REMOVED***,
		Collector: stats.NewCollector(),
	***REMOVED***, nil
***REMOVED***

func (u *VU) Reconfigure(id int64) error ***REMOVED***
	u.Iteration = 0

	return nil
***REMOVED***

func (u *VU) RunOnce(ctx context.Context) error ***REMOVED***
	u.Iteration++
	u.VM.Set("iteration", u.Iteration)

	if _, err := u.VM.Run(vuSetupScript); err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, ep := range u.Runner.Endpoints ***REMOVED***
		req := ep.Request()

		startTime := time.Now()
		res, err := u.Client.Do(&req)
		duration := time.Since(startTime)

		var status int
		var body []byte
		if err == nil ***REMOVED***
			status = res.StatusCode
			body, err = ioutil.ReadAll(res.Body)
			if err != nil ***REMOVED***
				res.Body.Close()
				return err
			***REMOVED***
			res.Body.Close()
		***REMOVED***

		tags := stats.Tags***REMOVED***"method": ep.Method, "url": ep.URLString, "status": status***REMOVED***
		u.Collector.Add(stats.Sample***REMOVED***
			Stat:   &mRequests,
			Tags:   tags,
			Values: stats.Values***REMOVED***"duration": float64(duration)***REMOVED***,
		***REMOVED***)

		if err != nil ***REMOVED***
			log.WithError(err).Error("Request error")
			u.Collector.Add(stats.Sample***REMOVED***
				Stat:   &mErrors,
				Tags:   tags,
				Values: stats.Value(1),
			***REMOVED***)
			return err
		***REMOVED***

		if len(ep.Tests) > 0 ***REMOVED***
			u.VM.Set("request", map[string]interface***REMOVED******REMOVED******REMOVED***
				"data":    ep.BodyMap,
				"headers": ep.HeaderMap,
				"method":  ep.Method,
				"url":     ep.URLString,
			***REMOVED***)

			responseHeaders := make(map[string]string)
			for key, values := range res.Header ***REMOVED***
				responseHeaders[key] = strings.Join(values, ", ")
			***REMOVED***
			u.VM.Set("responseHeaders", responseHeaders)

			// JSON seems to be geting automatically decoded by Postman? Is it decided by
			// Content-Type? Always attempted? We don't know, because it's nowhere in the docs!
			var obj interface***REMOVED******REMOVED***
			if err := json.Unmarshal(body, &obj); err != nil ***REMOVED***
				u.VM.Set("responseBody", string(body))
			***REMOVED*** else ***REMOVED***
				log.WithField("body", obj).Info("Body")
				u.VM.Set("responseBody", obj)
			***REMOVED***

			u.VM.Set("responseTime", duration/time.Millisecond)
			u.VM.Set("responseCode", map[string]interface***REMOVED******REMOVED******REMOVED***
				"code":   res.StatusCode,
				"name":   res.Status,
				"detail": res.Status, // The docs are vague on this one
			***REMOVED***)
			u.VM.Set("tests", map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***)

			for _, script := range ep.Tests ***REMOVED***
				if _, err := u.VM.Run(script); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
