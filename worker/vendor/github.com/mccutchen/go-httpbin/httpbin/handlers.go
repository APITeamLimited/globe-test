package httpbin

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mccutchen/go-httpbin/httpbin/assets"
	"github.com/mccutchen/go-httpbin/httpbin/digest"
)

var acceptedMediaTypes = []string***REMOVED***
	"image/webp",
	"image/svg+xml",
	"image/jpeg",
	"image/png",
	"image/",
***REMOVED***

func notImplementedHandler(w http.ResponseWriter, r *http.Request) ***REMOVED***
	http.Error(w, "Not implemented", http.StatusNotImplemented)
***REMOVED***

// Index renders an HTML index page
func (h *HTTPBin) Index(w http.ResponseWriter, r *http.Request) ***REMOVED***
	if r.URL.Path != "/" ***REMOVED***
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	***REMOVED***
	w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' camo.githubusercontent.com")
	writeHTML(w, assets.MustAsset("index.html"), http.StatusOK)
***REMOVED***

// FormsPost renders an HTML form that submits a request to the /post endpoint
func (h *HTTPBin) FormsPost(w http.ResponseWriter, r *http.Request) ***REMOVED***
	writeHTML(w, assets.MustAsset("forms-post.html"), http.StatusOK)
***REMOVED***

// UTF8 renders an HTML encoding stress test
func (h *HTTPBin) UTF8(w http.ResponseWriter, r *http.Request) ***REMOVED***
	writeHTML(w, assets.MustAsset("utf8.html"), http.StatusOK)
***REMOVED***

// Get handles HTTP GET requests
func (h *HTTPBin) Get(w http.ResponseWriter, r *http.Request) ***REMOVED***
	resp := &getResponse***REMOVED***
		Args:    r.URL.Query(),
		Headers: getRequestHeaders(r),
		Origin:  getOrigin(r),
		URL:     getURL(r).String(),
	***REMOVED***
	body, _ := json.Marshal(resp)
	writeJSON(w, body, http.StatusOK)
***REMOVED***

// RequestWithBody handles POST, PUT, and PATCH requests
func (h *HTTPBin) RequestWithBody(w http.ResponseWriter, r *http.Request) ***REMOVED***
	resp := &bodyResponse***REMOVED***
		Args:    r.URL.Query(),
		Headers: getRequestHeaders(r),
		Origin:  getOrigin(r),
		URL:     getURL(r).String(),
	***REMOVED***

	err := parseBody(w, r, resp)
	if err != nil ***REMOVED***
		http.Error(w, fmt.Sprintf("error parsing request body: %s", err), http.StatusBadRequest)
		return
	***REMOVED***

	body, _ := json.Marshal(resp)
	writeJSON(w, body, http.StatusOK)
***REMOVED***

// Gzip returns a gzipped response
func (h *HTTPBin) Gzip(w http.ResponseWriter, r *http.Request) ***REMOVED***
	resp := &gzipResponse***REMOVED***
		Headers: getRequestHeaders(r),
		Origin:  getOrigin(r),
		Gzipped: true,
	***REMOVED***
	body, _ := json.Marshal(resp)

	buf := &bytes.Buffer***REMOVED******REMOVED***
	gzw := gzip.NewWriter(buf)
	gzw.Write(body)
	gzw.Close()

	gzBody := buf.Bytes()

	w.Header().Set("Content-Encoding", "gzip")
	writeJSON(w, gzBody, http.StatusOK)
***REMOVED***

// Deflate returns a gzipped response
func (h *HTTPBin) Deflate(w http.ResponseWriter, r *http.Request) ***REMOVED***
	resp := &deflateResponse***REMOVED***
		Headers:  getRequestHeaders(r),
		Origin:   getOrigin(r),
		Deflated: true,
	***REMOVED***
	body, _ := json.Marshal(resp)

	buf := &bytes.Buffer***REMOVED******REMOVED***
	w2 := zlibWorker.NewWriter(buf)
	w2.Write(body)
	w2.Close()

	compressedBody := buf.Bytes()

	w.Header().Set("Content-Encoding", "deflate")
	writeJSON(w, compressedBody, http.StatusOK)
***REMOVED***

// IP echoes the IP address of the incoming request
func (h *HTTPBin) IP(w http.ResponseWriter, r *http.Request) ***REMOVED***
	body, _ := json.Marshal(&ipResponse***REMOVED***
		Origin: getOrigin(r),
	***REMOVED***)
	writeJSON(w, body, http.StatusOK)
***REMOVED***

// UserAgent echoes the incoming User-Agent header
func (h *HTTPBin) UserAgent(w http.ResponseWriter, r *http.Request) ***REMOVED***
	body, _ := json.Marshal(&userAgentResponse***REMOVED***
		UserAgent: r.Header.Get("User-Agent"),
	***REMOVED***)
	writeJSON(w, body, http.StatusOK)
***REMOVED***

// Headers echoes the incoming request headers
func (h *HTTPBin) Headers(w http.ResponseWriter, r *http.Request) ***REMOVED***
	body, _ := json.Marshal(&headersResponse***REMOVED***
		Headers: getRequestHeaders(r),
	***REMOVED***)
	writeJSON(w, body, http.StatusOK)
***REMOVED***

// Status responds with the specified status code. TODO: support random choice
// from multiple, optionally weighted status codes.
func (h *HTTPBin) Status(w http.ResponseWriter, r *http.Request) ***REMOVED***
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 ***REMOVED***
		http.Error(w, "Not found", http.StatusNotFound)
		return
	***REMOVED***
	code, err := strconv.Atoi(parts[2])
	if err != nil ***REMOVED***
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	***REMOVED***

	type statusCase struct ***REMOVED***
		headers map[string]string
		body    []byte
	***REMOVED***

	redirectHeaders := &statusCase***REMOVED***
		headers: map[string]string***REMOVED***
			"Location": "/redirect/1",
		***REMOVED***,
	***REMOVED***
	notAcceptableBody, _ := json.Marshal(map[string]interface***REMOVED******REMOVED******REMOVED***
		"message": "Client did not request a supported media type",
		"accept":  acceptedMediaTypes,
	***REMOVED***)

	specialCases := map[int]*statusCase***REMOVED***
		301: redirectHeaders,
		302: redirectHeaders,
		303: redirectHeaders,
		305: redirectHeaders,
		307: redirectHeaders,
		401: ***REMOVED***
			headers: map[string]string***REMOVED***
				"WWW-Authenticate": `Basic realm="Fake Realm"`,
			***REMOVED***,
		***REMOVED***,
		402: ***REMOVED***
			body: []byte("Fuck you, pay me!"),
			headers: map[string]string***REMOVED***
				"X-More-Info": "http://vimeo.com/22053820",
			***REMOVED***,
		***REMOVED***,
		406: ***REMOVED***
			body: notAcceptableBody,
			headers: map[string]string***REMOVED***
				"Content-Type": jsonContentType,
			***REMOVED***,
		***REMOVED***,
		407: ***REMOVED***
			headers: map[string]string***REMOVED***
				"Proxy-Authenticate": `Basic realm="Fake Realm"`,
			***REMOVED***,
		***REMOVED***,
		418: ***REMOVED***
			body: []byte("I'm a teapot!"),
			headers: map[string]string***REMOVED***
				"X-More-Info": "http://tools.ietf.org/html/rfc2324",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	if specialCase, ok := specialCases[code]; ok ***REMOVED***
		if specialCase.headers != nil ***REMOVED***
			for key, val := range specialCase.headers ***REMOVED***
				w.Header().Set(key, val)
			***REMOVED***
		***REMOVED***
		w.WriteHeader(code)
		if specialCase.body != nil ***REMOVED***
			w.Write(specialCase.body)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		w.WriteHeader(code)
	***REMOVED***
***REMOVED***

// ResponseHeaders responds with a map of header values
func (h *HTTPBin) ResponseHeaders(w http.ResponseWriter, r *http.Request) ***REMOVED***
	args := r.URL.Query()
	for k, vs := range args ***REMOVED***
		for _, v := range vs ***REMOVED***
			w.Header().Add(http.CanonicalHeaderKey(k), v)
		***REMOVED***
	***REMOVED***
	body, _ := json.Marshal(args)
	if contentType := w.Header().Get("Content-Type"); contentType == "" ***REMOVED***
		w.Header().Set("Content-Type", jsonContentType)
	***REMOVED***
	w.Write(body)
***REMOVED***

func redirectLocation(r *http.Request, relative bool, n int) string ***REMOVED***
	var location string
	var path string

	if n < 1 ***REMOVED***
		path = "/get"
	***REMOVED*** else if relative ***REMOVED***
		path = fmt.Sprintf("/relative-redirect/%d", n)
	***REMOVED*** else ***REMOVED***
		path = fmt.Sprintf("/absolute-redirect/%d", n)
	***REMOVED***

	if relative ***REMOVED***
		location = path
	***REMOVED*** else ***REMOVED***
		u := getURL(r)
		u.Path = path
		u.RawQuery = ""
		location = u.String()
	***REMOVED***

	return location
***REMOVED***

func doRedirect(w http.ResponseWriter, r *http.Request, relative bool) ***REMOVED***
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 ***REMOVED***
		http.Error(w, "Not found", http.StatusNotFound)
		return
	***REMOVED***
	n, err := strconv.Atoi(parts[2])
	if err != nil || n < 1 ***REMOVED***
		http.Error(w, "Invalid redirect", http.StatusBadRequest)
		return
	***REMOVED***

	w.Header().Set("Location", redirectLocation(r, relative, n-1))
	w.WriteHeader(http.StatusFound)
***REMOVED***

// Redirect responds with 302 redirect a given number of times. Defaults to a
// relative redirect, but an ?absolute=true query param will trigger an
// absolute redirect.
func (h *HTTPBin) Redirect(w http.ResponseWriter, r *http.Request) ***REMOVED***
	params := r.URL.Query()
	relative := strings.ToLower(params.Get("absolute")) != "true"
	doRedirect(w, r, relative)
***REMOVED***

// RelativeRedirect responds with an HTTP 302 redirect a given number of times
func (h *HTTPBin) RelativeRedirect(w http.ResponseWriter, r *http.Request) ***REMOVED***
	doRedirect(w, r, true)
***REMOVED***

// AbsoluteRedirect responds with an HTTP 302 redirect a given number of times
func (h *HTTPBin) AbsoluteRedirect(w http.ResponseWriter, r *http.Request) ***REMOVED***
	doRedirect(w, r, false)
***REMOVED***

// RedirectTo responds with a redirect to a specific URL with an optional
// status code, which defaults to 302
func (h *HTTPBin) RedirectTo(w http.ResponseWriter, r *http.Request) ***REMOVED***
	q := r.URL.Query()

	url := q.Get("url")
	if url == "" ***REMOVED***
		http.Error(w, "Missing URL", http.StatusBadRequest)
		return
	***REMOVED***

	var err error
	statusCode := http.StatusFound
	rawStatusCode := q.Get("status_code")
	if rawStatusCode != "" ***REMOVED***
		statusCode, err = strconv.Atoi(q.Get("status_code"))
		if err != nil || statusCode < 300 || statusCode > 399 ***REMOVED***
			http.Error(w, "Invalid status code", http.StatusBadRequest)
			return
		***REMOVED***
	***REMOVED***

	w.Header().Set("Location", url)
	w.WriteHeader(statusCode)
***REMOVED***

// Cookies responds with the cookies in the incoming request
func (h *HTTPBin) Cookies(w http.ResponseWriter, r *http.Request) ***REMOVED***
	resp := cookiesResponse***REMOVED******REMOVED***
	for _, c := range r.Cookies() ***REMOVED***
		resp[c.Name] = c.Value
	***REMOVED***
	body, _ := json.Marshal(resp)
	writeJSON(w, body, http.StatusOK)
***REMOVED***

// SetCookies sets cookies as specified in query params and redirects to
// Cookies endpoint
func (h *HTTPBin) SetCookies(w http.ResponseWriter, r *http.Request) ***REMOVED***
	params := r.URL.Query()
	for k := range params ***REMOVED***
		http.SetCookie(w, &http.Cookie***REMOVED***
			Name:     k,
			Value:    params.Get(k),
			HttpOnly: true,
		***REMOVED***)
	***REMOVED***
	w.Header().Set("Location", "/cookies")
	w.WriteHeader(http.StatusFound)
***REMOVED***

// DeleteCookies deletes cookies specified in query params and redirects to
// Cookies endpoint
func (h *HTTPBin) DeleteCookies(w http.ResponseWriter, r *http.Request) ***REMOVED***
	params := r.URL.Query()
	for k := range params ***REMOVED***
		http.SetCookie(w, &http.Cookie***REMOVED***
			Name:     k,
			Value:    params.Get(k),
			HttpOnly: true,
			MaxAge:   -1,
			Expires:  time.Now().Add(-1 * 24 * 365 * time.Hour),
		***REMOVED***)
	***REMOVED***
	w.Header().Set("Location", "/cookies")
	w.WriteHeader(http.StatusFound)
***REMOVED***

// BasicAuth requires basic authentication
func (h *HTTPBin) BasicAuth(w http.ResponseWriter, r *http.Request) ***REMOVED***
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 ***REMOVED***
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	***REMOVED***
	expectedUser := parts[2]
	expectedPass := parts[3]

	givenUser, givenPass, _ := r.BasicAuth()

	status := http.StatusOK
	authorized := givenUser == expectedUser && givenPass == expectedPass
	if !authorized ***REMOVED***
		status = http.StatusUnauthorized
		w.Header().Set("WWW-Authenticate", `Basic realm="Fake Realm"`)
	***REMOVED***

	body, _ := json.Marshal(&authResponse***REMOVED***
		Authorized: authorized,
		User:       givenUser,
	***REMOVED***)
	writeJSON(w, body, status)
***REMOVED***

// HiddenBasicAuth requires HTTP Basic authentication but returns a status of
// 404 if the request is unauthorized
func (h *HTTPBin) HiddenBasicAuth(w http.ResponseWriter, r *http.Request) ***REMOVED***
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 ***REMOVED***
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	***REMOVED***
	expectedUser := parts[2]
	expectedPass := parts[3]

	givenUser, givenPass, _ := r.BasicAuth()

	authorized := givenUser == expectedUser && givenPass == expectedPass
	if !authorized ***REMOVED***
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	***REMOVED***

	body, _ := json.Marshal(&authResponse***REMOVED***
		Authorized: authorized,
		User:       givenUser,
	***REMOVED***)
	writeJSON(w, body, http.StatusOK)
***REMOVED***

// Stream responds with max(n, 100) lines of JSON-encoded request data.
func (h *HTTPBin) Stream(w http.ResponseWriter, r *http.Request) ***REMOVED***
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 ***REMOVED***
		http.Error(w, "Not found", http.StatusNotFound)
		return
	***REMOVED***
	n, err := strconv.Atoi(parts[2])
	if err != nil ***REMOVED***
		http.Error(w, "Invalid integer", http.StatusBadRequest)
		return
	***REMOVED***

	if n > 100 ***REMOVED***
		n = 100
	***REMOVED*** else if n < 1 ***REMOVED***
		n = 1
	***REMOVED***

	resp := &streamResponse***REMOVED***
		Args:    r.URL.Query(),
		Headers: getRequestHeaders(r),
		Origin:  getOrigin(r),
		URL:     getURL(r).String(),
	***REMOVED***

	f := w.(http.Flusher)
	for i := 0; i < n; i++ ***REMOVED***
		resp.ID = i
		line, _ := json.Marshal(resp)
		w.Write(line)
		w.Write([]byte("\n"))
		f.Flush()
	***REMOVED***
***REMOVED***

// Delay waits for a given amount of time before responding, where the time may
// be specified as a golang-style duration or seconds in floating point.
func (h *HTTPBin) Delay(w http.ResponseWriter, r *http.Request) ***REMOVED***
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 ***REMOVED***
		http.Error(w, "Not found", http.StatusNotFound)
		return
	***REMOVED***

	delay, err := parseBoundedDuration(parts[2], 0, h.MaxDuration)
	if err != nil ***REMOVED***
		http.Error(w, "Invalid duration", http.StatusBadRequest)
		return
	***REMOVED***

	select ***REMOVED***
	case <-r.Context().Done():
		return
	case <-time.After(delay):
	***REMOVED***
	h.RequestWithBody(w, r)
***REMOVED***

// Drip returns data over a duration after an optional initial delay, then
// (optionally) returns with the given status code.
func (h *HTTPBin) Drip(w http.ResponseWriter, r *http.Request) ***REMOVED***
	q := r.URL.Query()

	duration := time.Duration(0)
	delay := time.Duration(0)
	numbytes := int64(10)
	code := http.StatusOK

	var err error

	userDuration := q.Get("duration")
	if userDuration != "" ***REMOVED***
		duration, err = parseBoundedDuration(userDuration, 0, h.MaxDuration)
		if err != nil ***REMOVED***
			http.Error(w, "Invalid duration", http.StatusBadRequest)
			return
		***REMOVED***
	***REMOVED***

	userDelay := q.Get("delay")
	if userDelay != "" ***REMOVED***
		delay, err = parseBoundedDuration(userDelay, 0, h.MaxDuration)
		if err != nil ***REMOVED***
			http.Error(w, "Invalid delay", http.StatusBadRequest)
			return
		***REMOVED***
	***REMOVED***

	userNumBytes := q.Get("numbytes")
	if userNumBytes != "" ***REMOVED***
		numbytes, err = strconv.ParseInt(userNumBytes, 10, 64)
		if err != nil || numbytes <= 0 || numbytes > h.MaxBodySize ***REMOVED***
			http.Error(w, "Invalid numbytes", http.StatusBadRequest)
			return
		***REMOVED***
	***REMOVED***

	userCode := q.Get("code")
	if userCode != "" ***REMOVED***
		code, err = strconv.Atoi(userCode)
		if err != nil || code < 100 || code >= 600 ***REMOVED***
			http.Error(w, "Invalid code", http.StatusBadRequest)
			return
		***REMOVED***
	***REMOVED***

	if duration+delay > h.MaxDuration ***REMOVED***
		http.Error(w, "Too much time", http.StatusBadRequest)
		return
	***REMOVED***

	pause := duration / time.Duration(numbytes)

	select ***REMOVED***
	case <-r.Context().Done():
		return
	case <-time.After(delay):
	***REMOVED***

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/octet-stream")

	f := w.(http.Flusher)
	for i := int64(0); i < numbytes; i++ ***REMOVED***
		w.Write([]byte("*"))
		f.Flush()

		select ***REMOVED***
		case <-r.Context().Done():
			return
		case <-time.After(pause):
		***REMOVED***
	***REMOVED***
***REMOVED***

// Range returns up to N bytes, with support for HTTP Range requests.
//
// This departs from httpbin by not supporting the chunk_size or duration
// parameters.
func (h *HTTPBin) Range(w http.ResponseWriter, r *http.Request) ***REMOVED***
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 ***REMOVED***
		http.Error(w, "Not found", http.StatusNotFound)
		return
	***REMOVED***

	numBytes, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil ***REMOVED***
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	***REMOVED***

	w.Header().Add("ETag", fmt.Sprintf("range%d", numBytes))
	w.Header().Add("Accept-Ranges", "bytes")

	if numBytes <= 0 || numBytes > h.MaxBodySize ***REMOVED***
		http.Error(w, "Invalid number of bytes", http.StatusBadRequest)
		return
	***REMOVED***

	content := newSyntheticByteStream(numBytes, func(offset int64) byte ***REMOVED***
		return byte(97 + (offset % 26))
	***REMOVED***)
	var modtime time.Time
	http.ServeContent(w, r, "", modtime, content)
***REMOVED***

// HTML renders a basic HTML page
func (h *HTTPBin) HTML(w http.ResponseWriter, r *http.Request) ***REMOVED***
	writeHTML(w, assets.MustAsset("moby.html"), http.StatusOK)
***REMOVED***

// Robots renders a basic robots.txt file
func (h *HTTPBin) Robots(w http.ResponseWriter, r *http.Request) ***REMOVED***
	robotsTxt := []byte(`User-agent: *
Disallow: /deny
`)
	writeResponse(w, http.StatusOK, "text/plain", robotsTxt)
***REMOVED***

// Deny renders a basic page that robots should never access
func (h *HTTPBin) Deny(w http.ResponseWriter, r *http.Request) ***REMOVED***
	writeResponse(w, http.StatusOK, "text/plain", []byte(`YOU SHOULDN'T BE HERE`))
***REMOVED***

// Cache returns a 304 if an If-Modified-Since or an If-None-Match header is
// present, otherwise returns the same response as Get.
func (h *HTTPBin) Cache(w http.ResponseWriter, r *http.Request) ***REMOVED***
	if r.Header.Get("If-Modified-Since") != "" || r.Header.Get("If-None-Match") != "" ***REMOVED***
		w.WriteHeader(http.StatusNotModified)
		return
	***REMOVED***

	lastModified := time.Now().Format(time.RFC1123)
	w.Header().Add("Last-Modified", lastModified)
	w.Header().Add("ETag", sha1hash(lastModified))
	h.Get(w, r)
***REMOVED***

// CacheControl sets a Cache-Control header for N seconds for /cache/N requests
func (h *HTTPBin) CacheControl(w http.ResponseWriter, r *http.Request) ***REMOVED***
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 ***REMOVED***
		http.Error(w, "Not found", http.StatusNotFound)
		return
	***REMOVED***

	seconds, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil ***REMOVED***
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	***REMOVED***

	w.Header().Add("Cache-Control", fmt.Sprintf("public, max-age=%d", seconds))
	h.Get(w, r)
***REMOVED***

// ETag assumes the resource has the given etag and response to If-None-Match
// and If-Match headers appropriately.
func (h *HTTPBin) ETag(w http.ResponseWriter, r *http.Request) ***REMOVED***
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 ***REMOVED***
		http.Error(w, "Not found", http.StatusNotFound)
		return
	***REMOVED***

	etag := parts[2]
	w.Header().Set("ETag", fmt.Sprintf(`"%s"`, etag))

	// TODO: This mostly duplicates the work of Get() above, should this be
	// pulled into a little helper?
	resp := &getResponse***REMOVED***
		Args:    r.URL.Query(),
		Headers: getRequestHeaders(r),
		Origin:  getOrigin(r),
		URL:     getURL(r).String(),
	***REMOVED***
	body, _ := json.Marshal(resp)

	// Let http.ServeContent deal with If-None-Match and If-Match headers:
	// https://golang.org/pkg/net/http/#ServeContent
	http.ServeContent(w, r, "response.json", time.Now(), bytes.NewReader(body))
***REMOVED***

// Bytes returns N random bytes generated with an optional seed
func (h *HTTPBin) Bytes(w http.ResponseWriter, r *http.Request) ***REMOVED***
	handleBytes(w, r, false)
***REMOVED***

// StreamBytes streams N random bytes generated with an optional seed in chunks
// of a given size.
func (h *HTTPBin) StreamBytes(w http.ResponseWriter, r *http.Request) ***REMOVED***
	handleBytes(w, r, true)
***REMOVED***

// handleBytes consolidates the logic for validating input params of the Bytes
// and StreamBytes endpoints and knows how to write the response in chunks if
// streaming is true.
func handleBytes(w http.ResponseWriter, r *http.Request, streaming bool) ***REMOVED***
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 ***REMOVED***
		http.Error(w, "Not found", http.StatusNotFound)
		return
	***REMOVED***

	numBytes, err := strconv.Atoi(parts[2])
	if err != nil ***REMOVED***
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	***REMOVED***

	if numBytes < 1 ***REMOVED***
		numBytes = 1
	***REMOVED*** else if numBytes > 100*1024 ***REMOVED***
		numBytes = 100 * 1024
	***REMOVED***

	var chunkSize int
	var write func([]byte)

	if streaming ***REMOVED***
		if r.URL.Query().Get("chunk_size") != "" ***REMOVED***
			chunkSize, err = strconv.Atoi(r.URL.Query().Get("chunk_size"))
			if err != nil ***REMOVED***
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			chunkSize = 10 * 1024
		***REMOVED***

		write = func() func(chunk []byte) ***REMOVED***
			f := w.(http.Flusher)
			return func(chunk []byte) ***REMOVED***
				w.Write(chunk)
				f.Flush()
			***REMOVED***
		***REMOVED***()
	***REMOVED*** else ***REMOVED***
		chunkSize = numBytes
		write = func(chunk []byte) ***REMOVED***
			w.Header().Set("Content-Length", strconv.Itoa(len(chunk)))
			w.Write(chunk)
		***REMOVED***
	***REMOVED***

	var seed int64
	rawSeed := r.URL.Query().Get("seed")
	if rawSeed != "" ***REMOVED***
		seed, err = strconv.ParseInt(rawSeed, 10, 64)
		if err != nil ***REMOVED***
			http.Error(w, "invalid seed", http.StatusBadRequest)
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		seed = time.Now().Unix()
	***REMOVED***

	src := rand.NewSource(seed)
	rng := rand.New(src)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)

	var chunk []byte
	for i := 0; i < numBytes; i++ ***REMOVED***
		chunk = append(chunk, byte(rng.Intn(256)))
		if len(chunk) == chunkSize ***REMOVED***
			write(chunk)
			chunk = nil
		***REMOVED***
	***REMOVED***
	if len(chunk) > 0 ***REMOVED***
		write(chunk)
	***REMOVED***
***REMOVED***

// Links redirects to the first page in a series of N links
func (h *HTTPBin) Links(w http.ResponseWriter, r *http.Request) ***REMOVED***
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 && len(parts) != 4 ***REMOVED***
		http.Error(w, "Not found", http.StatusNotFound)
		return
	***REMOVED***

	n, err := strconv.Atoi(parts[2])
	if err != nil || n < 0 || n > 256 ***REMOVED***
		http.Error(w, "Invalid link count", http.StatusBadRequest)
		return
	***REMOVED***

	// Are we handling /links/<n>/<offset>? If so, render an HTML page
	if len(parts) == 4 ***REMOVED***
		offset, err := strconv.Atoi(parts[3])
		if err != nil ***REMOVED***
			http.Error(w, "Invalid offset", http.StatusBadRequest)
		***REMOVED***
		doLinksPage(w, r, n, offset)
		return
	***REMOVED***

	// Otherwise, redirect from /links/<n> to /links/<n>/0
	r.URL.Path = r.URL.Path + "/0"
	w.Header().Set("Location", r.URL.String())
	w.WriteHeader(http.StatusFound)
***REMOVED***

// doLinksPage renders a page with a series of N links
func doLinksPage(w http.ResponseWriter, r *http.Request, n int, offset int) ***REMOVED***
	w.Header().Add("Content-Type", htmlContentType)
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("<html><head><title>Links</title></head><body>"))
	for i := 0; i < n; i++ ***REMOVED***
		if i == offset ***REMOVED***
			fmt.Fprintf(w, "%d ", i)
		***REMOVED*** else ***REMOVED***
			fmt.Fprintf(w, `<a href="/links/%d/%d">%d</a> `, n, i, i)
		***REMOVED***
	***REMOVED***
	w.Write([]byte("</body></html>"))
***REMOVED***

// ImageAccept responds with an appropriate image based on the Accept header
func (h *HTTPBin) ImageAccept(w http.ResponseWriter, r *http.Request) ***REMOVED***
	accept := r.Header.Get("Accept")
	if accept == "" || strings.Contains(accept, "image/png") || strings.Contains(accept, "image/*") ***REMOVED***
		doImage(w, "png")
	***REMOVED*** else if strings.Contains(accept, "image/webp") ***REMOVED***
		doImage(w, "webp")
	***REMOVED*** else if strings.Contains(accept, "image/svg+xml") ***REMOVED***
		doImage(w, "svg")
	***REMOVED*** else if strings.Contains(accept, "image/jpeg") ***REMOVED***
		doImage(w, "jpeg")
	***REMOVED*** else ***REMOVED***
		http.Error(w, "Unsupported media type", http.StatusUnsupportedMediaType)
	***REMOVED***
***REMOVED***

// Image responds with an image of a specific kind, from /image/<kind>
func (h *HTTPBin) Image(w http.ResponseWriter, r *http.Request) ***REMOVED***
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 ***REMOVED***
		http.Error(w, "Not found", http.StatusNotFound)
		return
	***REMOVED***
	doImage(w, parts[2])
***REMOVED***

// doImage responds with a specific kind of image, if there is an image asset
// of the given kind.
func doImage(w http.ResponseWriter, kind string) ***REMOVED***
	img, err := assets.Asset("image." + kind)
	if err != nil ***REMOVED***
		http.Error(w, "Not Found", http.StatusNotFound)
	***REMOVED***
	contentType := "image/" + kind
	if kind == "svg" ***REMOVED***
		contentType = "image/svg+xml"
	***REMOVED***
	writeResponse(w, http.StatusOK, contentType, img)
***REMOVED***

// XML responds with an XML document
func (h *HTTPBin) XML(w http.ResponseWriter, r *http.Request) ***REMOVED***
	writeResponse(w, http.StatusOK, "application/xml", assets.MustAsset("sample.xml"))
***REMOVED***

// DigestAuth handles a simple implementation of HTTP Digest Authentication,
// which supports the "auth" QOP and the MD5 and SHA-256 crypto algorithms.
//
// /digest-auth/<qop>/<user>/<passwd>
// /digest-auth/<qop>/<user>/<passwd>/<algorithm>
func (h *HTTPBin) DigestAuth(w http.ResponseWriter, r *http.Request) ***REMOVED***
	parts := strings.Split(r.URL.Path, "/")
	count := len(parts)

	if count != 5 && count != 6 ***REMOVED***
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	***REMOVED***

	qop := strings.ToLower(parts[2])
	user := parts[3]
	password := parts[4]

	algoName := "MD5"
	if count == 6 ***REMOVED***
		algoName = strings.ToUpper(parts[5])
	***REMOVED***

	if qop != "auth" ***REMOVED***
		http.Error(w, "Invalid QOP directive", http.StatusBadRequest)
		return
	***REMOVED***
	if algoName != "MD5" && algoName != "SHA-256" ***REMOVED***
		http.Error(w, "Invalid algorithm", http.StatusBadRequest)
		return
	***REMOVED***

	algorithm := digest.MD5
	if algoName == "SHA-256" ***REMOVED***
		algorithm = digest.SHA256
	***REMOVED***

	if !digest.Check(r, user, password) ***REMOVED***
		w.Header().Set("WWW-Authenticate", digest.Challenge("go-httpbin", algorithm))
		w.WriteHeader(http.StatusUnauthorized)
		return
	***REMOVED***

	resp, _ := json.Marshal(&authResponse***REMOVED***
		Authorized: true,
		User:       user,
	***REMOVED***)
	writeJSON(w, resp, http.StatusOK)
***REMOVED***
