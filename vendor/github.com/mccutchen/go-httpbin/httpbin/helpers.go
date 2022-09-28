package httpbin

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// requestHeaders takes in incoming request and returns an http.Header map
// suitable for inclusion in our response data structures.
//
// This is necessary to ensure that the incoming Host header is included,
// because golang only exposes that header on the http.Request struct itself.
func getRequestHeaders(r *http.Request) http.Header ***REMOVED***
	h := r.Header
	h.Set("Host", r.Host)
	return h
***REMOVED***

func getOrigin(r *http.Request) string ***REMOVED***
	origin := r.Header.Get("X-Forwarded-For")
	if origin == "" ***REMOVED***
		origin = r.RemoteAddr
	***REMOVED***
	return origin
***REMOVED***

func getURL(r *http.Request) *url.URL ***REMOVED***
	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" ***REMOVED***
		scheme = r.Header.Get("X-Forwarded-Protocol")
	***REMOVED***
	if scheme == "" && r.Header.Get("X-Forwarded-Ssl") == "on" ***REMOVED***
		scheme = "https"
	***REMOVED***
	if scheme == "" ***REMOVED***
		scheme = "http"
	***REMOVED***

	host := r.URL.Host
	if host == "" ***REMOVED***
		host = r.Host
	***REMOVED***

	return &url.URL***REMOVED***
		Scheme:     scheme,
		Opaque:     r.URL.Opaque,
		User:       r.URL.User,
		Host:       host,
		Path:       r.URL.Path,
		RawPath:    r.URL.RawPath,
		ForceQuery: r.URL.ForceQuery,
		RawQuery:   r.URL.RawQuery,
		Fragment:   r.URL.Fragment,
	***REMOVED***
***REMOVED***

func writeResponse(w http.ResponseWriter, status int, contentType string, body []byte) ***REMOVED***
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
	w.WriteHeader(status)
	w.Write(body)
***REMOVED***

func writeJSON(w http.ResponseWriter, body []byte, status int) ***REMOVED***
	writeResponse(w, status, jsonContentType, body)
***REMOVED***

func writeHTML(w http.ResponseWriter, body []byte, status int) ***REMOVED***
	writeResponse(w, status, htmlContentType, body)
***REMOVED***

// parseBody handles parsing a request body into our standard API response,
// taking care to only consume the request body once based on the Content-Type
// of the request. The given bodyResponse will be modified.
//
// Note: this function expects callers to limit the the maximum size of the
// request body. See, e.g., the limitRequestSize middleware.
func parseBody(w http.ResponseWriter, r *http.Request, resp *bodyResponse) error ***REMOVED***
	if r.Body == nil ***REMOVED***
		return nil
	***REMOVED***

	// Always set resp.Data to the incoming request body, in case we don't know
	// how to handle the content type
	body, err := ioutil.ReadAll(r.Body)
	if err != nil ***REMOVED***
		r.Body.Close()
		return err
	***REMOVED***
	resp.Data = string(body)

	// After reading the body to populate resp.Data, we need to re-wrap it in
	// an io.Reader for further processing below
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	ct := r.Header.Get("Content-Type")
	switch ***REMOVED***
	case strings.HasPrefix(ct, "application/x-www-form-urlencoded"):
		if err := r.ParseForm(); err != nil ***REMOVED***
			return err
		***REMOVED***
		resp.Form = r.PostForm
	case strings.HasPrefix(ct, "multipart/form-data"):
		// The memory limit here only restricts how many parts will be kept in
		// memory before overflowing to disk:
		// http://localhost:8080/pkg/net/http/#Request.ParseMultipartForm
		if err := r.ParseMultipartForm(1024); err != nil ***REMOVED***
			return err
		***REMOVED***
		resp.Form = r.PostForm
	case strings.HasPrefix(ct, "application/json"):
		err := json.NewDecoder(r.Body).Decode(&resp.JSON)
		if err != nil && err != io.EOF ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// parseDuration takes a user's input as a string and attempts to convert it
// into a time.Duration. If not given as a go-style duration string, the input
// is assumed to be seconds as a float.
func parseDuration(input string) (time.Duration, error) ***REMOVED***
	d, err := time.ParseDuration(input)
	if err != nil ***REMOVED***
		n, err := strconv.ParseFloat(input, 64)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		d = time.Duration(n*1000) * time.Millisecond
	***REMOVED***
	return d, nil
***REMOVED***

// parseBoundedDuration parses a time.Duration from user input and ensures that
// it is within a given maximum and minimum time
func parseBoundedDuration(input string, min, max time.Duration) (time.Duration, error) ***REMOVED***
	d, err := parseDuration(input)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	if d > max ***REMOVED***
		err = fmt.Errorf("duration %s longer than %s", d, max)
	***REMOVED*** else if d < min ***REMOVED***
		err = fmt.Errorf("duration %s shorter than %s", d, min)
	***REMOVED***
	return d, err
***REMOVED***

// syntheticByteStream implements the ReadSeeker interface to allow reading
// arbitrary subsets of bytes up to a maximum size given a function for
// generating the byte at a given offset.
type syntheticByteStream struct ***REMOVED***
	mu sync.Mutex

	size    int64
	offset  int64
	factory func(int64) byte
***REMOVED***

// newSyntheticByteStream returns a new stream of bytes of a specific size,
// given a factory function for generating the byte at a given offset.
func newSyntheticByteStream(size int64, factory func(int64) byte) io.ReadSeeker ***REMOVED***
	return &syntheticByteStream***REMOVED***
		size:    size,
		factory: factory,
	***REMOVED***
***REMOVED***

// Read implements the Reader interface for syntheticByteStream
func (s *syntheticByteStream) Read(p []byte) (int, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	start := s.offset
	end := start + int64(len(p))
	var err error
	if end >= s.size ***REMOVED***
		err = io.EOF
		end = s.size
	***REMOVED***

	for idx := start; idx < end; idx++ ***REMOVED***
		p[idx-start] = s.factory(idx)
	***REMOVED***

	s.offset = end

	return int(end - start), err
***REMOVED***

// Seek implements the Seeker interface for syntheticByteStream
func (s *syntheticByteStream) Seek(offset int64, whence int) (int64, error) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	switch whence ***REMOVED***
	case io.SeekStart:
		s.offset = offset
	case io.SeekCurrent:
		s.offset += offset
	case io.SeekEnd:
		s.offset = s.size - offset
	default:
		return 0, errors.New("Seek: invalid whence")
	***REMOVED***

	if s.offset < 0 ***REMOVED***
		return 0, errors.New("Seek: invalid offset")
	***REMOVED***

	return s.offset, nil
***REMOVED***

func sha1hash(input string) string ***REMOVED***
	h := sha1.New()
	return fmt.Sprintf("%x", h.Sum([]byte(input)))
***REMOVED***
