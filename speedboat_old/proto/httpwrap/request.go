package httpwrap

import (
	"context"
	"github.com/loadimpact/speedboat/stats"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
)

type Params struct ***REMOVED***
	TakeSample bool
	KeepBody   bool
***REMOVED***

func Do(ctx context.Context, client *http.Client, req *http.Request, params Params) (*http.Response, []byte, stats.Sample, error) ***REMOVED***
	var t Tracer
	if params.TakeSample ***REMOVED***
		trace := t.MakeClientTrace()
		ctx = httptrace.WithClientTrace(ctx, &trace)
	***REMOVED***

	res, err := client.Do(req.WithContext(ctx))
	if err != nil ***REMOVED***
		return res, nil, stats.Sample***REMOVED******REMOVED***, err
	***REMOVED***

	var body []byte
	if params.KeepBody ***REMOVED***
		body, err = ioutil.ReadAll(res.Body)
	***REMOVED*** else ***REMOVED***
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	***REMOVED***

	var sample stats.Sample
	if params.TakeSample ***REMOVED***
		t.RequestDone()
		sample.Tags = stats.Tags***REMOVED***
			"proto":  res.Proto,
			"method": req.Method,
			"url":    req.URL.String(),
			"status": res.StatusCode,
		***REMOVED***
		sample.Values = stats.Values***REMOVED***
			"duration":     float64(t.Duration),
			"ttfb":         float64(t.TimeToFirstByte),
			"lookup":       float64(t.TimeForDNS),
			"connect":      float64(t.TimeForConnect),
			"send":         float64(t.TimeForWriteHeaders + t.TimeForWriteBody),
			"send_headers": float64(t.TimeForWriteHeaders),
			"send_body":    float64(t.TimeForWriteBody),
			"wait":         float64(t.TimeWaiting),
			"receive":      float64(t.Duration - t.TimeToFirstByte),
		***REMOVED***
	***REMOVED***

	return res, body, sample, err
***REMOVED***
