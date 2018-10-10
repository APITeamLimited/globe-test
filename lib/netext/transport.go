package netext

import (
	"net"
	"net/http"
	"strconv"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	"github.com/pkg/errors"
)

type Transport struct ***REMOVED***
	roundTripper http.RoundTripper
	options      *lib.Options
	tags         map[string]string
	trail        *Trail
	tlsInfo      TLSInfo
	samplesCh    chan<- stats.SampleContainer
***REMOVED***

func NewTransport(transport http.RoundTripper, samplesCh chan<- stats.SampleContainer, options *lib.Options, tags map[string]string) *Transport ***REMOVED***
	return &Transport***REMOVED***
		roundTripper: transport,
		tags:         tags,
		options:      options,
		samplesCh:    samplesCh,
	***REMOVED***
***REMOVED***

func (t *Transport) SetOptions(options *lib.Options) ***REMOVED***
	t.options = options
***REMOVED***

func (t *Transport) GetTrail() *Trail ***REMOVED***
	return t.trail
***REMOVED***

func (t *Transport) TLSInfo() TLSInfo ***REMOVED***
	return t.tlsInfo
***REMOVED***

func (t *Transport) RoundTrip(req *http.Request) (res *http.Response, err error) ***REMOVED***
	if t.roundTripper == nil ***REMOVED***
		return nil, errors.New("No roundtrip defined")
	***REMOVED***

	tags := map[string]string***REMOVED******REMOVED***
	for k, v := range t.tags ***REMOVED***
		tags[k] = v
	***REMOVED***

	ctx := req.Context()
	tracer := Tracer***REMOVED******REMOVED***
	reqWithTracer := req.WithContext(WithTracer(ctx, &tracer))

	resp, err := t.roundTripper.RoundTrip(reqWithTracer)
	trail := tracer.Done()
	if err != nil ***REMOVED***
		if t.options.SystemTags["error"] ***REMOVED***
			tags["error"] = err.Error()
		***REMOVED***

		//TODO: expand/replace this so we can recognize the different non-HTTP
		// errors, probably by using a type switch for resErr
		if t.options.SystemTags["status"] ***REMOVED***
			tags["status"] = "0"
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if t.options.SystemTags["url"] ***REMOVED***
			tags["url"] = req.URL.String()
		***REMOVED***
		if t.options.SystemTags["status"] ***REMOVED***
			tags["status"] = strconv.Itoa(resp.StatusCode)
		***REMOVED***
		if t.options.SystemTags["proto"] ***REMOVED***
			tags["proto"] = resp.Proto
		***REMOVED***

		if resp.TLS != nil ***REMOVED***
			tlsInfo, oscp := ParseTLSConnState(resp.TLS)
			if t.options.SystemTags["tls_version"] ***REMOVED***
				tags["tls_version"] = tlsInfo.Version
			***REMOVED***
			if t.options.SystemTags["ocsp_status"] ***REMOVED***
				tags["ocsp_status"] = oscp.Status
			***REMOVED***

			t.tlsInfo = tlsInfo
		***REMOVED***
	***REMOVED***
	if t.options.SystemTags["ip"] && trail.ConnRemoteAddr != nil ***REMOVED***
		if ip, _, err := net.SplitHostPort(trail.ConnRemoteAddr.String()); err == nil ***REMOVED***
			tags["ip"] = ip
		***REMOVED***
	***REMOVED***

	t.trail = trail
	trail.SaveSamples(stats.IntoSampleTags(&tags))
	stats.PushIfNotCancelled(ctx, t.samplesCh, trail)

	return resp, err
***REMOVED***
