package postman

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/lib"
	"github.com/loadimpact/speedboat/stats"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	mRequests = stats.Stat***REMOVED***Name: "requests", Type: stats.HistogramType, Intent: stats.TimeIntent***REMOVED***
	mErrors   = stats.Stat***REMOVED***Name: "errors", Type: stats.CounterType***REMOVED***
)

type ErrorWithLineNumber struct ***REMOVED***
	Wrapped error
	Line    int
***REMOVED***

func (e ErrorWithLineNumber) Error() string ***REMOVED***
	return fmt.Sprintf("%s (line %d)", e.Wrapped.Error(), e.Line)
***REMOVED***

type Runner struct ***REMOVED***
	Collection Collection
***REMOVED***

type VU struct ***REMOVED***
	Runner    *Runner
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

	return &Runner***REMOVED***
		Collection: collection,
	***REMOVED***, nil
***REMOVED***

func (r *Runner) NewVU() (lib.VU, error) ***REMOVED***
	return &VU***REMOVED***
		Runner: r,
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
	for _, item := range u.Runner.Collection.Item ***REMOVED***
		if err := u.runItem(item, u.Runner.Collection.Auth); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (u *VU) runItem(i Item, a Auth) error ***REMOVED***
	if i.Auth.Type != "" ***REMOVED***
		a = i.Auth
	***REMOVED***

	if i.Request.URL != "" ***REMOVED***
		var buffer *bytes.Buffer
		switch i.Request.Body.Mode ***REMOVED***
		case "raw":
			buffer = bytes.NewBufferString(i.Request.Body.Raw)
		case "formdata":
			buffer = &bytes.Buffer***REMOVED******REMOVED***
			w := multipart.NewWriter(buffer)
			for _, field := range i.Request.Body.FormData ***REMOVED***
				if !field.Enabled ***REMOVED***
					continue
				***REMOVED***

				if err := w.WriteField(field.Key, field.Value); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		case "urlencoded":
			v := make(url.Values)
			for _, field := range i.Request.Body.URLEncoded ***REMOVED***
				if !field.Enabled ***REMOVED***
					continue
				***REMOVED***
				v[field.Key] = append(v[field.Key], field.Value)
			***REMOVED***
			buffer = bytes.NewBufferString(v.Encode())
		***REMOVED***

		req, err := http.NewRequest(i.Request.Method, i.Request.URL, buffer)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		startTime := time.Now()
		res, err := u.Client.Do(req)
		duration := time.Since(startTime)

		status := 0
		if err == nil ***REMOVED***
			status = res.StatusCode
			io.Copy(ioutil.Discard, res.Body)
			res.Body.Close()
		***REMOVED***

		tags := stats.Tags***REMOVED***"method": i.Request.Method, "url": i.Request.URL, "status": status***REMOVED***
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
	***REMOVED***

	for _, item := range i.Item ***REMOVED***
		if err := u.runItem(item, a); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
