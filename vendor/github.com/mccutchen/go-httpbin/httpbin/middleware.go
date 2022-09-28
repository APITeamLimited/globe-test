package httpbin

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func preflight(h http.Handler) http.Handler ***REMOVED***
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		origin := r.Header.Get("Origin")
		if origin == "" ***REMOVED***
			origin = "*"
		***REMOVED***
		respHeader := w.Header()
		respHeader.Set("Access-Control-Allow-Origin", origin)
		respHeader.Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" ***REMOVED***
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Max-Age", "3600")
			if r.Header.Get("Access-Control-Request-Headers") != "" ***REMOVED***
				w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
			***REMOVED***
			w.WriteHeader(200)
			return
		***REMOVED***

		h.ServeHTTP(w, r)
	***REMOVED***)
***REMOVED***

func methods(h http.HandlerFunc, methods ...string) http.HandlerFunc ***REMOVED***
	methodMap := make(map[string]struct***REMOVED******REMOVED***, len(methods))
	for _, m := range methods ***REMOVED***
		methodMap[m] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		// GET implies support for HEAD
		if m == "GET" ***REMOVED***
			methodMap["HEAD"] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
	return func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if _, ok := methodMap[r.Method]; !ok ***REMOVED***
			http.Error(w, fmt.Sprintf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
			return
		***REMOVED***
		h.ServeHTTP(w, r)
	***REMOVED***
***REMOVED***

func limitRequestSize(maxSize int64, h http.Handler) http.Handler ***REMOVED***
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Body != nil ***REMOVED***
			r.Body = http.MaxBytesReader(w, r.Body, maxSize)
		***REMOVED***
		h.ServeHTTP(w, r)
	***REMOVED***)
***REMOVED***

// headResponseWriter implements http.ResponseWriter in order to discard the
// body of the response
type headResponseWriter struct ***REMOVED***
	http.ResponseWriter
***REMOVED***

func (hw *headResponseWriter) Write(b []byte) (int, error) ***REMOVED***
	return 0, nil
***REMOVED***

// autohead automatically discards the body of responses to HEAD requests
func autohead(h http.Handler) http.Handler ***REMOVED***
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method == "HEAD" ***REMOVED***
			w = &headResponseWriter***REMOVED***w***REMOVED***
		***REMOVED***
		h.ServeHTTP(w, r)
	***REMOVED***)
***REMOVED***

// metaResponseWriter implements http.ResponseWriter and http.Flusher in order
// to record a response's status code and body size for logging purposes.
type metaResponseWriter struct ***REMOVED***
	w      http.ResponseWriter
	status int
	size   int64
***REMOVED***

func (mw *metaResponseWriter) Write(b []byte) (int, error) ***REMOVED***
	size, err := mw.w.Write(b)
	mw.size += int64(size)
	return size, err
***REMOVED***

func (mw *metaResponseWriter) WriteHeader(s int) ***REMOVED***
	mw.w.WriteHeader(s)
	mw.status = s
***REMOVED***

func (mw *metaResponseWriter) Flush() ***REMOVED***
	f := mw.w.(http.Flusher)
	f.Flush()
***REMOVED***

func (mw *metaResponseWriter) Header() http.Header ***REMOVED***
	return mw.w.Header()
***REMOVED***

func (mw *metaResponseWriter) Status() int ***REMOVED***
	if mw.status == 0 ***REMOVED***
		return http.StatusOK
	***REMOVED***
	return mw.status
***REMOVED***

func (mw *metaResponseWriter) Size() int64 ***REMOVED***
	return mw.size
***REMOVED***

func observe(o Observer, h http.Handler) http.Handler ***REMOVED***
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		mw := &metaResponseWriter***REMOVED***w: w***REMOVED***
		t := time.Now()
		h.ServeHTTP(mw, r)
		o(Result***REMOVED***
			Status:   mw.Status(),
			Method:   r.Method,
			URI:      r.URL.RequestURI(),
			Size:     mw.Size(),
			Duration: time.Now().Sub(t),
		***REMOVED***)
	***REMOVED***)
***REMOVED***

// Result is the result of handling a request, used for instrumentation
type Result struct ***REMOVED***
	Status   int
	Method   string
	URI      string
	Size     int64
	Duration time.Duration
***REMOVED***

// Observer is a function that will be called with the details of a handled
// request, which can be used for logging, instrumentation, etc
type Observer func(result Result)

// StdLogObserver creates an Observer that will log each request in structured
// format using the given stdlib logger
func StdLogObserver(l *log.Logger) Observer ***REMOVED***
	const (
		logFmt  = "time=%q status=%d method=%q uri=%q size_bytes=%d duration_ms=%0.02f"
		dateFmt = "2006-01-02T15:04:05.9999"
	)
	return func(result Result) ***REMOVED***
		l.Printf(
			logFmt,
			time.Now().Format(dateFmt),
			result.Status,
			result.Method,
			result.URI,
			result.Size,
			result.Duration.Seconds()*1e3, // https://github.com/golang/go/issues/5491#issuecomment-66079585
		)
	***REMOVED***
***REMOVED***
