package echo

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type (
	// Context represents the context of the current HTTP request. It holds request and
	// response objects, path, path parameters, data and registered handler.
	Context interface ***REMOVED***
		// Request returns `*http.Request`.
		Request() *http.Request

		// SetRequest sets `*http.Request`.
		SetRequest(r *http.Request)

		// Response returns `*Response`.
		Response() *Response

		// IsTLS returns true if HTTP connection is TLS otherwise false.
		IsTLS() bool

		// IsWebSocket returns true if HTTP connection is WebSocket otherwise false.
		IsWebSocket() bool

		// Scheme returns the HTTP protocol scheme, `http` or `https`.
		Scheme() string

		// RealIP returns the client's network address based on `X-Forwarded-For`
		// or `X-Real-IP` request header.
		RealIP() string

		// Path returns the registered path for the handler.
		Path() string

		// SetPath sets the registered path for the handler.
		SetPath(p string)

		// Param returns path parameter by name.
		Param(name string) string

		// ParamNames returns path parameter names.
		ParamNames() []string

		// SetParamNames sets path parameter names.
		SetParamNames(names ...string)

		// ParamValues returns path parameter values.
		ParamValues() []string

		// SetParamValues sets path parameter values.
		SetParamValues(values ...string)

		// QueryParam returns the query param for the provided name.
		QueryParam(name string) string

		// QueryParams returns the query parameters as `url.Values`.
		QueryParams() url.Values

		// QueryString returns the URL query string.
		QueryString() string

		// FormValue returns the form field value for the provided name.
		FormValue(name string) string

		// FormParams returns the form parameters as `url.Values`.
		FormParams() (url.Values, error)

		// FormFile returns the multipart form file for the provided name.
		FormFile(name string) (*multipart.FileHeader, error)

		// MultipartForm returns the multipart form.
		MultipartForm() (*multipart.Form, error)

		// Cookie returns the named cookie provided in the request.
		Cookie(name string) (*http.Cookie, error)

		// SetCookie adds a `Set-Cookie` header in HTTP response.
		SetCookie(cookie *http.Cookie)

		// Cookies returns the HTTP cookies sent with the request.
		Cookies() []*http.Cookie

		// Get retrieves data from the context.
		Get(key string) interface***REMOVED******REMOVED***

		// Set saves data in the context.
		Set(key string, val interface***REMOVED******REMOVED***)

		// Bind binds the request body into provided type `i`. The default binder
		// does it based on Content-Type header.
		Bind(i interface***REMOVED******REMOVED***) error

		// Validate validates provided `i`. It is usually called after `Context#Bind()`.
		// Validator must be registered using `Echo#Validator`.
		Validate(i interface***REMOVED******REMOVED***) error

		// Render renders a template with data and sends a text/html response with status
		// code. Renderer must be registered using `Echo.Renderer`.
		Render(code int, name string, data interface***REMOVED******REMOVED***) error

		// HTML sends an HTTP response with status code.
		HTML(code int, html string) error

		// HTMLBlob sends an HTTP blob response with status code.
		HTMLBlob(code int, b []byte) error

		// String sends a string response with status code.
		String(code int, s string) error

		// JSON sends a JSON response with status code.
		JSON(code int, i interface***REMOVED******REMOVED***) error

		// JSONPretty sends a pretty-print JSON with status code.
		JSONPretty(code int, i interface***REMOVED******REMOVED***, indent string) error

		// JSONBlob sends a JSON blob response with status code.
		JSONBlob(code int, b []byte) error

		// JSONP sends a JSONP response with status code. It uses `callback` to construct
		// the JSONP payload.
		JSONP(code int, callback string, i interface***REMOVED******REMOVED***) error

		// JSONPBlob sends a JSONP blob response with status code. It uses `callback`
		// to construct the JSONP payload.
		JSONPBlob(code int, callback string, b []byte) error

		// XML sends an XML response with status code.
		XML(code int, i interface***REMOVED******REMOVED***) error

		// XMLPretty sends a pretty-print XML with status code.
		XMLPretty(code int, i interface***REMOVED******REMOVED***, indent string) error

		// XMLBlob sends an XML blob response with status code.
		XMLBlob(code int, b []byte) error

		// Blob sends a blob response with status code and content type.
		Blob(code int, contentType string, b []byte) error

		// Stream sends a streaming response with status code and content type.
		Stream(code int, contentType string, r io.Reader) error

		// File sends a response with the content of the file.
		File(file string) error

		// Attachment sends a response as attachment, prompting client to save the
		// file.
		Attachment(file string, name string) error

		// Inline sends a response as inline, opening the file in the browser.
		Inline(file string, name string) error

		// NoContent sends a response with no body and a status code.
		NoContent(code int) error

		// Redirect redirects the request to a provided URL with status code.
		Redirect(code int, url string) error

		// Error invokes the registered HTTP error handler. Generally used by middleware.
		Error(err error)

		// Handler returns the matched handler by router.
		Handler() HandlerFunc

		// SetHandler sets the matched handler by router.
		SetHandler(h HandlerFunc)

		// Logger returns the `Logger` instance.
		Logger() Logger

		// Echo returns the `Echo` instance.
		Echo() *Echo

		// Reset resets the context after request completes. It must be called along
		// with `Echo#AcquireContext()` and `Echo#ReleaseContext()`.
		// See `Echo#ServeHTTP()`
		Reset(r *http.Request, w http.ResponseWriter)
	***REMOVED***

	context struct ***REMOVED***
		request  *http.Request
		response *Response
		path     string
		pnames   []string
		pvalues  []string
		query    url.Values
		handler  HandlerFunc
		store    Map
		echo     *Echo
	***REMOVED***
)

const (
	defaultMemory = 32 << 20 // 32 MB
	indexPage     = "index.html"
)

func (c *context) Request() *http.Request ***REMOVED***
	return c.request
***REMOVED***

func (c *context) SetRequest(r *http.Request) ***REMOVED***
	c.request = r
***REMOVED***

func (c *context) Response() *Response ***REMOVED***
	return c.response
***REMOVED***

func (c *context) IsTLS() bool ***REMOVED***
	return c.request.TLS != nil
***REMOVED***

func (c *context) IsWebSocket() bool ***REMOVED***
	upgrade := c.request.Header.Get(HeaderUpgrade)
	return upgrade == "websocket" || upgrade == "Websocket"
***REMOVED***

func (c *context) Scheme() string ***REMOVED***
	// Can't use `r.Request.URL.Scheme`
	// See: https://groups.google.com/forum/#!topic/golang-nuts/pMUkBlQBDF0
	if c.IsTLS() ***REMOVED***
		return "https"
	***REMOVED***
	if scheme := c.request.Header.Get(HeaderXForwardedProto); scheme != "" ***REMOVED***
		return scheme
	***REMOVED***
	if scheme := c.request.Header.Get(HeaderXForwardedProtocol); scheme != "" ***REMOVED***
		return scheme
	***REMOVED***
	if ssl := c.request.Header.Get(HeaderXForwardedSsl); ssl == "on" ***REMOVED***
		return "https"
	***REMOVED***
	if scheme := c.request.Header.Get(HeaderXUrlScheme); scheme != "" ***REMOVED***
		return scheme
	***REMOVED***
	return "http"
***REMOVED***

func (c *context) RealIP() string ***REMOVED***
	ra := c.request.RemoteAddr
	if ip := c.request.Header.Get(HeaderXForwardedFor); ip != "" ***REMOVED***
		ra = strings.Split(ip, ", ")[0]
	***REMOVED*** else if ip := c.request.Header.Get(HeaderXRealIP); ip != "" ***REMOVED***
		ra = ip
	***REMOVED*** else ***REMOVED***
		ra, _, _ = net.SplitHostPort(ra)
	***REMOVED***
	return ra
***REMOVED***

func (c *context) Path() string ***REMOVED***
	return c.path
***REMOVED***

func (c *context) SetPath(p string) ***REMOVED***
	c.path = p
***REMOVED***

func (c *context) Param(name string) string ***REMOVED***
	for i, n := range c.pnames ***REMOVED***
		if i < len(c.pvalues) ***REMOVED***
			if n == name ***REMOVED***
				return c.pvalues[i]
			***REMOVED***

			// Param name with aliases
			for _, p := range strings.Split(n, ",") ***REMOVED***
				if p == name ***REMOVED***
					return c.pvalues[i]
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

func (c *context) ParamNames() []string ***REMOVED***
	return c.pnames
***REMOVED***

func (c *context) SetParamNames(names ...string) ***REMOVED***
	c.pnames = names
***REMOVED***

func (c *context) ParamValues() []string ***REMOVED***
	return c.pvalues[:len(c.pnames)]
***REMOVED***

func (c *context) SetParamValues(values ...string) ***REMOVED***
	c.pvalues = values
***REMOVED***

func (c *context) QueryParam(name string) string ***REMOVED***
	if c.query == nil ***REMOVED***
		c.query = c.request.URL.Query()
	***REMOVED***
	return c.query.Get(name)
***REMOVED***

func (c *context) QueryParams() url.Values ***REMOVED***
	if c.query == nil ***REMOVED***
		c.query = c.request.URL.Query()
	***REMOVED***
	return c.query
***REMOVED***

func (c *context) QueryString() string ***REMOVED***
	return c.request.URL.RawQuery
***REMOVED***

func (c *context) FormValue(name string) string ***REMOVED***
	return c.request.FormValue(name)
***REMOVED***

func (c *context) FormParams() (url.Values, error) ***REMOVED***
	if strings.HasPrefix(c.request.Header.Get(HeaderContentType), MIMEMultipartForm) ***REMOVED***
		if err := c.request.ParseMultipartForm(defaultMemory); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if err := c.request.ParseForm(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return c.request.Form, nil
***REMOVED***

func (c *context) FormFile(name string) (*multipart.FileHeader, error) ***REMOVED***
	_, fh, err := c.request.FormFile(name)
	return fh, err
***REMOVED***

func (c *context) MultipartForm() (*multipart.Form, error) ***REMOVED***
	err := c.request.ParseMultipartForm(defaultMemory)
	return c.request.MultipartForm, err
***REMOVED***

func (c *context) Cookie(name string) (*http.Cookie, error) ***REMOVED***
	return c.request.Cookie(name)
***REMOVED***

func (c *context) SetCookie(cookie *http.Cookie) ***REMOVED***
	http.SetCookie(c.Response(), cookie)
***REMOVED***

func (c *context) Cookies() []*http.Cookie ***REMOVED***
	return c.request.Cookies()
***REMOVED***

func (c *context) Get(key string) interface***REMOVED******REMOVED*** ***REMOVED***
	return c.store[key]
***REMOVED***

func (c *context) Set(key string, val interface***REMOVED******REMOVED***) ***REMOVED***
	if c.store == nil ***REMOVED***
		c.store = make(Map)
	***REMOVED***
	c.store[key] = val
***REMOVED***

func (c *context) Bind(i interface***REMOVED******REMOVED***) error ***REMOVED***
	return c.echo.Binder.Bind(i, c)
***REMOVED***

func (c *context) Validate(i interface***REMOVED******REMOVED***) error ***REMOVED***
	if c.echo.Validator == nil ***REMOVED***
		return ErrValidatorNotRegistered
	***REMOVED***
	return c.echo.Validator.Validate(i)
***REMOVED***

func (c *context) Render(code int, name string, data interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	if c.echo.Renderer == nil ***REMOVED***
		return ErrRendererNotRegistered
	***REMOVED***
	buf := new(bytes.Buffer)
	if err = c.echo.Renderer.Render(buf, name, data, c); err != nil ***REMOVED***
		return
	***REMOVED***
	return c.HTMLBlob(code, buf.Bytes())
***REMOVED***

func (c *context) HTML(code int, html string) (err error) ***REMOVED***
	return c.HTMLBlob(code, []byte(html))
***REMOVED***

func (c *context) HTMLBlob(code int, b []byte) (err error) ***REMOVED***
	return c.Blob(code, MIMETextHTMLCharsetUTF8, b)
***REMOVED***

func (c *context) String(code int, s string) (err error) ***REMOVED***
	return c.Blob(code, MIMETextPlainCharsetUTF8, []byte(s))
***REMOVED***

func (c *context) JSON(code int, i interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	_, pretty := c.QueryParams()["pretty"]
	if c.echo.Debug || pretty ***REMOVED***
		return c.JSONPretty(code, i, "  ")
	***REMOVED***
	b, err := json.Marshal(i)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return c.JSONBlob(code, b)
***REMOVED***

func (c *context) JSONPretty(code int, i interface***REMOVED******REMOVED***, indent string) (err error) ***REMOVED***
	b, err := json.MarshalIndent(i, "", indent)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return c.JSONBlob(code, b)
***REMOVED***

func (c *context) JSONBlob(code int, b []byte) (err error) ***REMOVED***
	return c.Blob(code, MIMEApplicationJSONCharsetUTF8, b)
***REMOVED***

func (c *context) JSONP(code int, callback string, i interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	b, err := json.Marshal(i)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return c.JSONPBlob(code, callback, b)
***REMOVED***

func (c *context) JSONPBlob(code int, callback string, b []byte) (err error) ***REMOVED***
	c.response.Header().Set(HeaderContentType, MIMEApplicationJavaScriptCharsetUTF8)
	c.response.WriteHeader(code)
	if _, err = c.response.Write([]byte(callback + "(")); err != nil ***REMOVED***
		return
	***REMOVED***
	if _, err = c.response.Write(b); err != nil ***REMOVED***
		return
	***REMOVED***
	_, err = c.response.Write([]byte(");"))
	return
***REMOVED***

func (c *context) XML(code int, i interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	_, pretty := c.QueryParams()["pretty"]
	if c.echo.Debug || pretty ***REMOVED***
		return c.XMLPretty(code, i, "  ")
	***REMOVED***
	b, err := xml.Marshal(i)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return c.XMLBlob(code, b)
***REMOVED***

func (c *context) XMLPretty(code int, i interface***REMOVED******REMOVED***, indent string) (err error) ***REMOVED***
	b, err := xml.MarshalIndent(i, "", indent)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return c.XMLBlob(code, b)
***REMOVED***

func (c *context) XMLBlob(code int, b []byte) (err error) ***REMOVED***
	c.response.Header().Set(HeaderContentType, MIMEApplicationXMLCharsetUTF8)
	c.response.WriteHeader(code)
	if _, err = c.response.Write([]byte(xml.Header)); err != nil ***REMOVED***
		return
	***REMOVED***
	_, err = c.response.Write(b)
	return
***REMOVED***

func (c *context) Blob(code int, contentType string, b []byte) (err error) ***REMOVED***
	c.response.Header().Set(HeaderContentType, contentType)
	c.response.WriteHeader(code)
	_, err = c.response.Write(b)
	return
***REMOVED***

func (c *context) Stream(code int, contentType string, r io.Reader) (err error) ***REMOVED***
	c.response.Header().Set(HeaderContentType, contentType)
	c.response.WriteHeader(code)
	_, err = io.Copy(c.response, r)
	return
***REMOVED***

func (c *context) File(file string) (err error) ***REMOVED***
	f, err := os.Open(file)
	if err != nil ***REMOVED***
		return NotFoundHandler(c)
	***REMOVED***
	defer f.Close()

	fi, _ := f.Stat()
	if fi.IsDir() ***REMOVED***
		file = filepath.Join(file, indexPage)
		f, err = os.Open(file)
		if err != nil ***REMOVED***
			return NotFoundHandler(c)
		***REMOVED***
		defer f.Close()
		if fi, err = f.Stat(); err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	http.ServeContent(c.Response(), c.Request(), fi.Name(), fi.ModTime(), f)
	return
***REMOVED***

func (c *context) Attachment(file, name string) (err error) ***REMOVED***
	return c.contentDisposition(file, name, "attachment")
***REMOVED***

func (c *context) Inline(file, name string) (err error) ***REMOVED***
	return c.contentDisposition(file, name, "inline")
***REMOVED***

func (c *context) contentDisposition(file, name, dispositionType string) (err error) ***REMOVED***
	c.response.Header().Set(HeaderContentDisposition, fmt.Sprintf("%s; filename=%q", dispositionType, name))
	c.File(file)
	return
***REMOVED***

func (c *context) NoContent(code int) error ***REMOVED***
	c.response.WriteHeader(code)
	return nil
***REMOVED***

func (c *context) Redirect(code int, url string) error ***REMOVED***
	if code < 300 || code > 308 ***REMOVED***
		return ErrInvalidRedirectCode
	***REMOVED***
	c.response.Header().Set(HeaderLocation, url)
	c.response.WriteHeader(code)
	return nil
***REMOVED***

func (c *context) Error(err error) ***REMOVED***
	c.echo.HTTPErrorHandler(err, c)
***REMOVED***

func (c *context) Echo() *Echo ***REMOVED***
	return c.echo
***REMOVED***

func (c *context) Handler() HandlerFunc ***REMOVED***
	return c.handler
***REMOVED***

func (c *context) SetHandler(h HandlerFunc) ***REMOVED***
	c.handler = h
***REMOVED***

func (c *context) Logger() Logger ***REMOVED***
	return c.echo.Logger
***REMOVED***

func (c *context) Reset(r *http.Request, w http.ResponseWriter) ***REMOVED***
	c.request = r
	c.response.reset(w)
	c.query = nil
	c.handler = NotFoundHandler
	c.store = nil
	c.path = ""
	c.pnames = nil
	// NOTE: Don't reset because it has to have length c.echo.maxParam at all times
	// c.pvalues = nil
***REMOVED***
