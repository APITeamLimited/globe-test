// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"errors"
	"io"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"
)

// Content-Type MIME of the most common data formats
const (
	MIMEJSON              = binding.MIMEJSON
	MIMEHTML              = binding.MIMEHTML
	MIMEXML               = binding.MIMEXML
	MIMEXML2              = binding.MIMEXML2
	MIMEPlain             = binding.MIMEPlain
	MIMEPOSTForm          = binding.MIMEPOSTForm
	MIMEMultipartPOSTForm = binding.MIMEMultipartPOSTForm
)

const (
	defaultMemory      = 32 << 20 // 32 MB
	abortIndex    int8 = math.MaxInt8 / 2
)

// Context is the most important part of gin. It allows us to pass variables between middleware,
// manage the flow, validate the JSON of a request and render a JSON response for example.
type Context struct ***REMOVED***
	writermem responseWriter
	Request   *http.Request
	Writer    ResponseWriter

	Params   Params
	handlers HandlersChain
	index    int8

	engine   *Engine
	Keys     map[string]interface***REMOVED******REMOVED***
	Errors   errorMsgs
	Accepted []string
***REMOVED***

/************************************/
/********** CONTEXT CREATION ********/
/************************************/

func (c *Context) reset() ***REMOVED***
	c.Writer = &c.writermem
	c.Params = c.Params[0:0]
	c.handlers = nil
	c.index = -1
	c.Keys = nil
	c.Errors = c.Errors[0:0]
	c.Accepted = nil
***REMOVED***

// Copy returns a copy of the current context that can be safely used outside the request's scope.
// This has to be used when the context has to be passed to a goroutine.
func (c *Context) Copy() *Context ***REMOVED***
	var cp = *c
	cp.writermem.ResponseWriter = nil
	cp.Writer = &cp.writermem
	cp.index = abortIndex
	cp.handlers = nil
	return &cp
***REMOVED***

// HandlerName returns the main handler's name. For example if the handler is "handleGetUsers()", this
// function will return "main.handleGetUsers"
func (c *Context) HandlerName() string ***REMOVED***
	return nameOfFunction(c.handlers.Last())
***REMOVED***

// Handler returns the main handler.
func (c *Context) Handler() HandlerFunc ***REMOVED***
	return c.handlers.Last()
***REMOVED***

/************************************/
/*********** FLOW CONTROL ***********/
/************************************/

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in GitHub.
func (c *Context) Next() ***REMOVED***
	c.index++
	s := int8(len(c.handlers))
	for ; c.index < s; c.index++ ***REMOVED***
		c.handlers[c.index](c)
	***REMOVED***
***REMOVED***

// IsAborted returns true if the current context was aborted.
func (c *Context) IsAborted() bool ***REMOVED***
	return c.index >= abortIndex
***REMOVED***

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
// Let's say you have an authorization middleware that validates that the current request is authorized. If the
// authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
// for this request are not called.
func (c *Context) Abort() ***REMOVED***
	c.index = abortIndex
***REMOVED***

// AbortWithStatus calls `Abort()` and writes the headers with the specified status code.
// For example, a failed attempt to authenticate a request could use: context.AbortWithStatus(401).
func (c *Context) AbortWithStatus(code int) ***REMOVED***
	c.Status(code)
	c.Writer.WriteHeaderNow()
	c.Abort()
***REMOVED***

// AbortWithStatusJSON calls `Abort()` and then `JSON` internally. This method stops the chain, writes the status code and return a JSON body
// It also sets the Content-Type as "application/json".
func (c *Context) AbortWithStatusJSON(code int, jsonObj interface***REMOVED******REMOVED***) ***REMOVED***
	c.Abort()
	c.JSON(code, jsonObj)
***REMOVED***

// AbortWithError calls `AbortWithStatus()` and `Error()` internally. This method stops the chain, writes the status code and
// pushes the specified error to `c.Errors`.
// See Context.Error() for more details.
func (c *Context) AbortWithError(code int, err error) *Error ***REMOVED***
	c.AbortWithStatus(code)
	return c.Error(err)
***REMOVED***

/************************************/
/********* ERROR MANAGEMENT *********/
/************************************/

// Attaches an error to the current context. The error is pushed to a list of errors.
// It's a good idea to call Error for each error that occurred during the resolution of a request.
// A middleware can be used to collect all the errors
// and push them to a database together, print a log, or append it in the HTTP response.
func (c *Context) Error(err error) *Error ***REMOVED***
	var parsedError *Error
	switch err.(type) ***REMOVED***
	case *Error:
		parsedError = err.(*Error)
	default:
		parsedError = &Error***REMOVED***
			Err:  err,
			Type: ErrorTypePrivate,
		***REMOVED***
	***REMOVED***
	c.Errors = append(c.Errors, parsedError)
	return parsedError
***REMOVED***

/************************************/
/******** METADATA MANAGEMENT********/
/************************************/

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *Context) Set(key string, value interface***REMOVED******REMOVED***) ***REMOVED***
	if c.Keys == nil ***REMOVED***
		c.Keys = make(map[string]interface***REMOVED******REMOVED***)
	***REMOVED***
	c.Keys[key] = value
***REMOVED***

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *Context) Get(key string) (value interface***REMOVED******REMOVED***, exists bool) ***REMOVED***
	value, exists = c.Keys[key]
	return
***REMOVED***

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (c *Context) MustGet(key string) interface***REMOVED******REMOVED*** ***REMOVED***
	if value, exists := c.Get(key); exists ***REMOVED***
		return value
	***REMOVED***
	panic("Key \"" + key + "\" does not exist")
***REMOVED***

// GetString returns the value associated with the key as a string.
func (c *Context) GetString(key string) (s string) ***REMOVED***
	if val, ok := c.Get(key); ok && val != nil ***REMOVED***
		s, _ = val.(string)
	***REMOVED***
	return
***REMOVED***

// GetBool returns the value associated with the key as a boolean.
func (c *Context) GetBool(key string) (b bool) ***REMOVED***
	if val, ok := c.Get(key); ok && val != nil ***REMOVED***
		b, _ = val.(bool)
	***REMOVED***
	return
***REMOVED***

// GetInt returns the value associated with the key as an integer.
func (c *Context) GetInt(key string) (i int) ***REMOVED***
	if val, ok := c.Get(key); ok && val != nil ***REMOVED***
		i, _ = val.(int)
	***REMOVED***
	return
***REMOVED***

// GetInt64 returns the value associated with the key as an integer.
func (c *Context) GetInt64(key string) (i64 int64) ***REMOVED***
	if val, ok := c.Get(key); ok && val != nil ***REMOVED***
		i64, _ = val.(int64)
	***REMOVED***
	return
***REMOVED***

// GetFloat64 returns the value associated with the key as a float64.
func (c *Context) GetFloat64(key string) (f64 float64) ***REMOVED***
	if val, ok := c.Get(key); ok && val != nil ***REMOVED***
		f64, _ = val.(float64)
	***REMOVED***
	return
***REMOVED***

// GetTime returns the value associated with the key as time.
func (c *Context) GetTime(key string) (t time.Time) ***REMOVED***
	if val, ok := c.Get(key); ok && val != nil ***REMOVED***
		t, _ = val.(time.Time)
	***REMOVED***
	return
***REMOVED***

// GetDuration returns the value associated with the key as a duration.
func (c *Context) GetDuration(key string) (d time.Duration) ***REMOVED***
	if val, ok := c.Get(key); ok && val != nil ***REMOVED***
		d, _ = val.(time.Duration)
	***REMOVED***
	return
***REMOVED***

// GetStringSlice returns the value associated with the key as a slice of strings.
func (c *Context) GetStringSlice(key string) (ss []string) ***REMOVED***
	if val, ok := c.Get(key); ok && val != nil ***REMOVED***
		ss, _ = val.([]string)
	***REMOVED***
	return
***REMOVED***

// GetStringMap returns the value associated with the key as a map of interfaces.
func (c *Context) GetStringMap(key string) (sm map[string]interface***REMOVED******REMOVED***) ***REMOVED***
	if val, ok := c.Get(key); ok && val != nil ***REMOVED***
		sm, _ = val.(map[string]interface***REMOVED******REMOVED***)
	***REMOVED***
	return
***REMOVED***

// GetStringMapString returns the value associated with the key as a map of strings.
func (c *Context) GetStringMapString(key string) (sms map[string]string) ***REMOVED***
	if val, ok := c.Get(key); ok && val != nil ***REMOVED***
		sms, _ = val.(map[string]string)
	***REMOVED***
	return
***REMOVED***

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (c *Context) GetStringMapStringSlice(key string) (smss map[string][]string) ***REMOVED***
	if val, ok := c.Get(key); ok && val != nil ***REMOVED***
		smss, _ = val.(map[string][]string)
	***REMOVED***
	return
***REMOVED***

/************************************/
/************ INPUT DATA ************/
/************************************/

// Param returns the value of the URL param.
// It is a shortcut for c.Params.ByName(key)
//		router.GET("/user/:id", func(c *gin.Context) ***REMOVED***
//			// a GET request to /user/john
//			id := c.Param("id") // id == "john"
//		***REMOVED***)
func (c *Context) Param(key string) string ***REMOVED***
	return c.Params.ByName(key)
***REMOVED***

// Query returns the keyed url query value if it exists,
// otherwise it returns an empty string `("")`.
// It is shortcut for `c.Request.URL.Query().Get(key)`
// 		GET /path?id=1234&name=Manu&value=
// 		c.Query("id") == "1234"
// 		c.Query("name") == "Manu"
// 		c.Query("value") == ""
// 		c.Query("wtf") == ""
func (c *Context) Query(key string) string ***REMOVED***
	value, _ := c.GetQuery(key)
	return value
***REMOVED***

// DefaultQuery returns the keyed url query value if it exists,
// otherwise it returns the specified defaultValue string.
// See: Query() and GetQuery() for further information.
// 		GET /?name=Manu&lastname=
// 		c.DefaultQuery("name", "unknown") == "Manu"
// 		c.DefaultQuery("id", "none") == "none"
// 		c.DefaultQuery("lastname", "none") == ""
func (c *Context) DefaultQuery(key, defaultValue string) string ***REMOVED***
	if value, ok := c.GetQuery(key); ok ***REMOVED***
		return value
	***REMOVED***
	return defaultValue
***REMOVED***

// GetQuery is like Query(), it returns the keyed url query value
// if it exists `(value, true)` (even when the value is an empty string),
// otherwise it returns `("", false)`.
// It is shortcut for `c.Request.URL.Query().Get(key)`
// 		GET /?name=Manu&lastname=
// 		("Manu", true) == c.GetQuery("name")
// 		("", false) == c.GetQuery("id")
// 		("", true) == c.GetQuery("lastname")
func (c *Context) GetQuery(key string) (string, bool) ***REMOVED***
	if values, ok := c.GetQueryArray(key); ok ***REMOVED***
		return values[0], ok
	***REMOVED***
	return "", false
***REMOVED***

// QueryArray returns a slice of strings for a given query key.
// The length of the slice depends on the number of params with the given key.
func (c *Context) QueryArray(key string) []string ***REMOVED***
	values, _ := c.GetQueryArray(key)
	return values
***REMOVED***

// GetQueryArray returns a slice of strings for a given query key, plus
// a boolean value whether at least one value exists for the given key.
func (c *Context) GetQueryArray(key string) ([]string, bool) ***REMOVED***
	req := c.Request
	if values, ok := req.URL.Query()[key]; ok && len(values) > 0 ***REMOVED***
		return values, true
	***REMOVED***
	return []string***REMOVED******REMOVED***, false
***REMOVED***

// PostForm returns the specified key from a POST urlencoded form or multipart form
// when it exists, otherwise it returns an empty string `("")`.
func (c *Context) PostForm(key string) string ***REMOVED***
	value, _ := c.GetPostForm(key)
	return value
***REMOVED***

// DefaultPostForm returns the specified key from a POST urlencoded form or multipart form
// when it exists, otherwise it returns the specified defaultValue string.
// See: PostForm() and GetPostForm() for further information.
func (c *Context) DefaultPostForm(key, defaultValue string) string ***REMOVED***
	if value, ok := c.GetPostForm(key); ok ***REMOVED***
		return value
	***REMOVED***
	return defaultValue
***REMOVED***

// GetPostForm is like PostForm(key). It returns the specified key from a POST urlencoded
// form or multipart form when it exists `(value, true)` (even when the value is an empty string),
// otherwise it returns ("", false).
// For example, during a PATCH request to update the user's email:
// 		email=mail@example.com  -->  ("mail@example.com", true) := GetPostForm("email") // set email to "mail@example.com"
// 		email=  			  	-->  ("", true) := GetPostForm("email") // set email to ""
//							 	-->  ("", false) := GetPostForm("email") // do nothing with email
func (c *Context) GetPostForm(key string) (string, bool) ***REMOVED***
	if values, ok := c.GetPostFormArray(key); ok ***REMOVED***
		return values[0], ok
	***REMOVED***
	return "", false
***REMOVED***

// PostFormArray returns a slice of strings for a given form key.
// The length of the slice depends on the number of params with the given key.
func (c *Context) PostFormArray(key string) []string ***REMOVED***
	values, _ := c.GetPostFormArray(key)
	return values
***REMOVED***

// GetPostFormArray returns a slice of strings for a given form key, plus
// a boolean value whether at least one value exists for the given key.
func (c *Context) GetPostFormArray(key string) ([]string, bool) ***REMOVED***
	req := c.Request
	req.ParseForm()
	req.ParseMultipartForm(defaultMemory)
	if values := req.PostForm[key]; len(values) > 0 ***REMOVED***
		return values, true
	***REMOVED***
	if req.MultipartForm != nil && req.MultipartForm.File != nil ***REMOVED***
		if values := req.MultipartForm.Value[key]; len(values) > 0 ***REMOVED***
			return values, true
		***REMOVED***
	***REMOVED***
	return []string***REMOVED******REMOVED***, false
***REMOVED***

// FormFile returns the first file for the provided form key.
func (c *Context) FormFile(name string) (*multipart.FileHeader, error) ***REMOVED***
	_, fh, err := c.Request.FormFile(name)
	return fh, err
***REMOVED***

// MultipartForm is the parsed multipart form, including file uploads.
func (c *Context) MultipartForm() (*multipart.Form, error) ***REMOVED***
	err := c.Request.ParseMultipartForm(defaultMemory)
	return c.Request.MultipartForm, err
***REMOVED***

// Bind checks the Content-Type to select a binding engine automatically,
// Depending the "Content-Type" header different bindings are used:
// 		"application/json" --> JSON binding
// 		"application/xml"  --> XML binding
// otherwise --> returns an error
// It parses the request's body as JSON if Content-Type == "application/json" using JSON or XML as a JSON input.
// It decodes the json payload into the struct specified as a pointer.
// Like ParseBody() but this method also writes a 400 error if the json is not valid.
func (c *Context) Bind(obj interface***REMOVED******REMOVED***) error ***REMOVED***
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.MustBindWith(obj, b)
***REMOVED***

// BindJSON is a shortcut for c.MustBindWith(obj, binding.JSON)
func (c *Context) BindJSON(obj interface***REMOVED******REMOVED***) error ***REMOVED***
	return c.MustBindWith(obj, binding.JSON)
***REMOVED***

// MustBindWith binds the passed struct pointer using the specified binding
// engine. It will abort the request with HTTP 400 if any error ocurrs.
// See the binding package.
func (c *Context) MustBindWith(obj interface***REMOVED******REMOVED***, b binding.Binding) (err error) ***REMOVED***
	if err = c.ShouldBindWith(obj, b); err != nil ***REMOVED***
		c.AbortWithError(400, err).SetType(ErrorTypeBind)
	***REMOVED***

	return
***REMOVED***

// ShouldBindWith binds the passed struct pointer using the specified binding
// engine.
// See the binding package.
func (c *Context) ShouldBindWith(obj interface***REMOVED******REMOVED***, b binding.Binding) error ***REMOVED***
	return b.Bind(c.Request, obj)
***REMOVED***

// ClientIP implements a best effort algorithm to return the real client IP, it parses
// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
// Use X-Forwarded-For before X-Real-Ip as nginx uses X-Real-Ip with the proxy's IP.
func (c *Context) ClientIP() string ***REMOVED***
	if c.engine.ForwardedByClientIP ***REMOVED***
		clientIP := c.requestHeader("X-Forwarded-For")
		if index := strings.IndexByte(clientIP, ','); index >= 0 ***REMOVED***
			clientIP = clientIP[0:index]
		***REMOVED***
		clientIP = strings.TrimSpace(clientIP)
		if len(clientIP) > 0 ***REMOVED***
			return clientIP
		***REMOVED***
		clientIP = strings.TrimSpace(c.requestHeader("X-Real-Ip"))
		if len(clientIP) > 0 ***REMOVED***
			return clientIP
		***REMOVED***
	***REMOVED***

	if c.engine.AppEngine ***REMOVED***
		if addr := c.Request.Header.Get("X-Appengine-Remote-Addr"); addr != "" ***REMOVED***
			return addr
		***REMOVED***
	***REMOVED***

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr)); err == nil ***REMOVED***
		return ip
	***REMOVED***

	return ""
***REMOVED***

// ContentType returns the Content-Type header of the request.
func (c *Context) ContentType() string ***REMOVED***
	return filterFlags(c.requestHeader("Content-Type"))
***REMOVED***

// IsWebsocket returns true if the request headers indicate that a websocket
// handshake is being initiated by the client.
func (c *Context) IsWebsocket() bool ***REMOVED***
	if strings.Contains(strings.ToLower(c.requestHeader("Connection")), "upgrade") &&
		strings.ToLower(c.requestHeader("Upgrade")) == "websocket" ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (c *Context) requestHeader(key string) string ***REMOVED***
	if values, _ := c.Request.Header[key]; len(values) > 0 ***REMOVED***
		return values[0]
	***REMOVED***
	return ""
***REMOVED***

/************************************/
/******** RESPONSE RENDERING ********/
/************************************/

// bodyAllowedForStatus is a copy of http.bodyAllowedForStatus non-exported function
func bodyAllowedForStatus(status int) bool ***REMOVED***
	switch ***REMOVED***
	case status >= 100 && status <= 199:
		return false
	case status == 204:
		return false
	case status == 304:
		return false
	***REMOVED***
	return true
***REMOVED***

func (c *Context) Status(code int) ***REMOVED***
	c.writermem.WriteHeader(code)
***REMOVED***

// Header is a intelligent shortcut for c.Writer.Header().Set(key, value)
// It writes a header in the response.
// If value == "", this method removes the header `c.Writer.Header().Del(key)`
func (c *Context) Header(key, value string) ***REMOVED***
	if len(value) == 0 ***REMOVED***
		c.Writer.Header().Del(key)
	***REMOVED*** else ***REMOVED***
		c.Writer.Header().Set(key, value)
	***REMOVED***
***REMOVED***

// GetHeader returns value from request headers
func (c *Context) GetHeader(key string) string ***REMOVED***
	return c.requestHeader(key)
***REMOVED***

// GetRawData return stream data
func (c *Context) GetRawData() ([]byte, error) ***REMOVED***
	return ioutil.ReadAll(c.Request.Body)
***REMOVED***

func (c *Context) SetCookie(
	name string,
	value string,
	maxAge int,
	path string,
	domain string,
	secure bool,
	httpOnly bool,
) ***REMOVED***
	if path == "" ***REMOVED***
		path = "/"
	***REMOVED***
	http.SetCookie(c.Writer, &http.Cookie***REMOVED***
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	***REMOVED***)
***REMOVED***

func (c *Context) Cookie(name string) (string, error) ***REMOVED***
	cookie, err := c.Request.Cookie(name)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
***REMOVED***

func (c *Context) Render(code int, r render.Render) ***REMOVED***
	c.Status(code)

	if !bodyAllowedForStatus(code) ***REMOVED***
		r.WriteContentType(c.Writer)
		c.Writer.WriteHeaderNow()
		return
	***REMOVED***

	if err := r.Render(c.Writer); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

// HTML renders the HTTP template specified by its file name.
// It also updates the HTTP code and sets the Content-Type as "text/html".
// See http://golang.org/doc/articles/wiki/
func (c *Context) HTML(code int, name string, obj interface***REMOVED******REMOVED***) ***REMOVED***
	instance := c.engine.HTMLRender.Instance(name, obj)
	c.Render(code, instance)
***REMOVED***

// IndentedJSON serializes the given struct as pretty JSON (indented + endlines) into the response body.
// It also sets the Content-Type as "application/json".
// WARNING: we recommend to use this only for development purposes since printing pretty JSON is
// more CPU and bandwidth consuming. Use Context.JSON() instead.
func (c *Context) IndentedJSON(code int, obj interface***REMOVED******REMOVED***) ***REMOVED***
	c.Render(code, render.IndentedJSON***REMOVED***Data: obj***REMOVED***)
***REMOVED***

// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json".
func (c *Context) JSON(code int, obj interface***REMOVED******REMOVED***) ***REMOVED***
	c.Render(code, render.JSON***REMOVED***Data: obj***REMOVED***)
***REMOVED***

// XML serializes the given struct as XML into the response body.
// It also sets the Content-Type as "application/xml".
func (c *Context) XML(code int, obj interface***REMOVED******REMOVED***) ***REMOVED***
	c.Render(code, render.XML***REMOVED***Data: obj***REMOVED***)
***REMOVED***

// YAML serializes the given struct as YAML into the response body.
func (c *Context) YAML(code int, obj interface***REMOVED******REMOVED***) ***REMOVED***
	c.Render(code, render.YAML***REMOVED***Data: obj***REMOVED***)
***REMOVED***

// String writes the given string into the response body.
func (c *Context) String(code int, format string, values ...interface***REMOVED******REMOVED***) ***REMOVED***
	c.Render(code, render.String***REMOVED***Format: format, Data: values***REMOVED***)
***REMOVED***

// Redirect returns a HTTP redirect to the specific location.
func (c *Context) Redirect(code int, location string) ***REMOVED***
	c.Render(-1, render.Redirect***REMOVED***
		Code:     code,
		Location: location,
		Request:  c.Request,
	***REMOVED***)
***REMOVED***

// Data writes some data into the body stream and updates the HTTP code.
func (c *Context) Data(code int, contentType string, data []byte) ***REMOVED***
	c.Render(code, render.Data***REMOVED***
		ContentType: contentType,
		Data:        data,
	***REMOVED***)
***REMOVED***

// File writes the specified file into the body stream in a efficient way.
func (c *Context) File(filepath string) ***REMOVED***
	http.ServeFile(c.Writer, c.Request, filepath)
***REMOVED***

// SSEvent writes a Server-Sent Event into the body stream.
func (c *Context) SSEvent(name string, message interface***REMOVED******REMOVED***) ***REMOVED***
	c.Render(-1, sse.Event***REMOVED***
		Event: name,
		Data:  message,
	***REMOVED***)
***REMOVED***

func (c *Context) Stream(step func(w io.Writer) bool) ***REMOVED***
	w := c.Writer
	clientGone := w.CloseNotify()
	for ***REMOVED***
		select ***REMOVED***
		case <-clientGone:
			return
		default:
			keepOpen := step(w)
			w.Flush()
			if !keepOpen ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

/************************************/
/******** CONTENT NEGOTIATION *******/
/************************************/

type Negotiate struct ***REMOVED***
	Offered  []string
	HTMLName string
	HTMLData interface***REMOVED******REMOVED***
	JSONData interface***REMOVED******REMOVED***
	XMLData  interface***REMOVED******REMOVED***
	Data     interface***REMOVED******REMOVED***
***REMOVED***

func (c *Context) Negotiate(code int, config Negotiate) ***REMOVED***
	switch c.NegotiateFormat(config.Offered...) ***REMOVED***
	case binding.MIMEJSON:
		data := chooseData(config.JSONData, config.Data)
		c.JSON(code, data)

	case binding.MIMEHTML:
		data := chooseData(config.HTMLData, config.Data)
		c.HTML(code, config.HTMLName, data)

	case binding.MIMEXML:
		data := chooseData(config.XMLData, config.Data)
		c.XML(code, data)

	default:
		c.AbortWithError(http.StatusNotAcceptable, errors.New("the accepted formats are not offered by the server"))
	***REMOVED***
***REMOVED***

func (c *Context) NegotiateFormat(offered ...string) string ***REMOVED***
	assert1(len(offered) > 0, "you must provide at least one offer")

	if c.Accepted == nil ***REMOVED***
		c.Accepted = parseAccept(c.requestHeader("Accept"))
	***REMOVED***
	if len(c.Accepted) == 0 ***REMOVED***
		return offered[0]
	***REMOVED***
	for _, accepted := range c.Accepted ***REMOVED***
		for _, offert := range offered ***REMOVED***
			if accepted == offert ***REMOVED***
				return offert
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

func (c *Context) SetAccepted(formats ...string) ***REMOVED***
	c.Accepted = formats
***REMOVED***

/************************************/
/***** GOLANG.ORG/X/NET/CONTEXT *****/
/************************************/

func (c *Context) Deadline() (deadline time.Time, ok bool) ***REMOVED***
	return
***REMOVED***

func (c *Context) Done() <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	return nil
***REMOVED***

func (c *Context) Err() error ***REMOVED***
	return nil
***REMOVED***

func (c *Context) Value(key interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	if key == 0 ***REMOVED***
		return c.Request
	***REMOVED***
	if keyAsString, ok := key.(string); ok ***REMOVED***
		val, _ := c.Get(keyAsString)
		return val
	***REMOVED***
	return nil
***REMOVED***
