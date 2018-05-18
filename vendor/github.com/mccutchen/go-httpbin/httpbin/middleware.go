package httpbin

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"
)

func metaRequests(h http.Handler) http.Handler ***REMOVED***
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		origin := r.Header.Get("Origin")
		if origin == "" ***REMOVED***
			origin = "*"
		***REMOVED***
		respHeader := w.Header()
		respHeader.Set("Access-Control-Allow-Origin", origin)
		respHeader.Set("Access-Control-Allow-Credentials", "true")

		switch r.Method ***REMOVED***
		case "OPTIONS":
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Max-Age", "3600")
			if r.Header.Get("Access-Control-Request-Headers") != "" ***REMOVED***
				w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
			***REMOVED***
			w.WriteHeader(200)
		case "HEAD":
			rwRec := httptest.NewRecorder()
			r.Method = "GET"
			h.ServeHTTP(rwRec, r)

			copyHeader(w.Header(), rwRec.Header())
			w.WriteHeader(rwRec.Code)
		default:
			h.ServeHTTP(w, r)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func methods(h http.HandlerFunc, methods ...string) http.HandlerFunc ***REMOVED***
	methodMap := make(map[string]struct***REMOVED******REMOVED***, len(methods))
	for _, m := range methods ***REMOVED***
		methodMap[m] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
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

// metaResponseWriter implements http.ResponseWriter and http.Flusher in order
// to record a response's status code and body size for logging purposes.
type metaResponseWriter struct ***REMOVED***
	w      http.ResponseWriter
	status int
	size   int
***REMOVED***

func (mw *metaResponseWriter) Write(b []byte) (int, error) ***REMOVED***
	size, err := mw.w.Write(b)
	mw.size += size
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

func (mw *metaResponseWriter) Size() int ***REMOVED***
	return mw.size
***REMOVED***

func logger(h http.Handler) http.Handler ***REMOVED***
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		reqMethod, reqURI := r.Method, r.URL.RequestURI()
		mw := &metaResponseWriter***REMOVED***w: w***REMOVED***
		t := time.Now()
		h.ServeHTTP(mw, r)
		duration := time.Now().Sub(t)
		log.Printf("status=%d method=%s uri=%q size=%d duration=%s", mw.Status(), reqMethod, reqURI, mw.Size(), duration)
	***REMOVED***)
***REMOVED***