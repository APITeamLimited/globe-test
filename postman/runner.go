package postman

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/lib"
	"github.com/loadimpact/speedboat/stats"
	"github.com/robertkrimen/otto"
	"golang.org/x/net/context"
	"io"
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

const SETUP_SRC = `
// Scripts populate this with test resultks; keys are strings, values are bools,
// or just something truthy/falsy, because this is javascript.
var tests = ***REMOVED******REMOVED***;
`

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
	if _, err := vm.Eval(SETUP_SRC); err != nil ***REMOVED***
		return nil, err
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
	return nil
***REMOVED***

func (u *VU) RunOnce(ctx context.Context) error ***REMOVED***
	for _, ep := range u.Runner.Endpoints ***REMOVED***
		req := ep.Request()

		startTime := time.Now()
		res, err := u.Client.Do(&req)
		duration := time.Since(startTime)

		status := 0
		if err == nil ***REMOVED***
			status = res.StatusCode
			io.Copy(ioutil.Discard, res.Body)
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

		for _, script := range ep.Tests ***REMOVED***
			if _, err := u.VM.Run(script); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
