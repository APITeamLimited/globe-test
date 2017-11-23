package negroni

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"text/template"
)

const (
	panicText = "PANIC: %s\n%s"
	panicHTML = `<html>
<head><title>PANIC: ***REMOVED******REMOVED***.RecoveredPanic***REMOVED******REMOVED***</title></head>
<style type="text/css">
html, body ***REMOVED***
	font-family: Helvetica, Arial, Sans;
	color: #333333;
	background-color: #ffffff;
	margin: 0px;
***REMOVED***
h1 ***REMOVED***
	color: #ffffff;
	background-color: #f14c4c;
	padding: 20px;
	border-bottom: 1px solid #2b3848;
***REMOVED***
.block ***REMOVED***
	margin: 2em;
***REMOVED***
.panic-interface ***REMOVED***
***REMOVED***

.panic-stack-raw pre ***REMOVED***
	padding: 1em;
	background: #f6f8fa;
	border: dashed 1px;
***REMOVED***
.panic-interface-title ***REMOVED***
	font-weight: bold;
***REMOVED***
</style>
<body>
<h1>Negroni - PANIC</h1>

<div class="panic-interface block">
	<h3>***REMOVED******REMOVED***.RequestDescription***REMOVED******REMOVED***</h3>
	<span class="panic-interface-title">Runtime error:</span> <span class="panic-interface-element">***REMOVED******REMOVED***.RecoveredPanic***REMOVED******REMOVED***</span>
</div>

***REMOVED******REMOVED*** if .Stack ***REMOVED******REMOVED***
<div class="panic-stack-raw block">
	<h3>Runtime Stack</h3>
	<pre>***REMOVED******REMOVED***.StackAsString***REMOVED******REMOVED***</pre>
</div>
***REMOVED******REMOVED*** end ***REMOVED******REMOVED***

</body>
</html>`
	nilRequestMessage = "Request is nil"
)

var panicHTMLTemplate = template.Must(template.New("PanicPage").Parse(panicHTML))

// PanicInformation contains all
// elements for printing stack informations.
type PanicInformation struct ***REMOVED***
	RecoveredPanic interface***REMOVED******REMOVED***
	Stack          []byte
	Request        *http.Request
***REMOVED***

// StackAsString returns a printable version of the stack
func (p *PanicInformation) StackAsString() string ***REMOVED***
	return string(p.Stack)
***REMOVED***

// RequestDescription returns a printable description of the url
func (p *PanicInformation) RequestDescription() string ***REMOVED***

	if p.Request == nil ***REMOVED***
		return nilRequestMessage
	***REMOVED***

	var queryOutput string
	if p.Request.URL.RawQuery != "" ***REMOVED***
		queryOutput = "?" + p.Request.URL.RawQuery
	***REMOVED***
	return fmt.Sprintf("%s %s%s", p.Request.Method, p.Request.URL.Path, queryOutput)
***REMOVED***

// PanicFormatter is an interface on object can implement
// to be able to output the stack trace
type PanicFormatter interface ***REMOVED***
	// FormatPanicError output the stack for a given answer/response.
	// In case the the middleware should not output the stack trace,
	// the field `Stack` of the passed `PanicInformation` instance equals `[]byte***REMOVED******REMOVED***`.
	FormatPanicError(rw http.ResponseWriter, r *http.Request, infos *PanicInformation)
***REMOVED***

// TextPanicFormatter output the stack
// as simple text on os.Stdout. If no `Content-Type` is set,
// it will output the data as `text/plain; charset=utf-8`.
// Otherwise, the origin `Content-Type` is kept.
type TextPanicFormatter struct***REMOVED******REMOVED***

func (t *TextPanicFormatter) FormatPanicError(rw http.ResponseWriter, r *http.Request, infos *PanicInformation) ***REMOVED***
	if rw.Header().Get("Content-Type") == "" ***REMOVED***
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	***REMOVED***
	fmt.Fprintf(rw, panicText, infos.RecoveredPanic, infos.Stack)
***REMOVED***

// HTMLPanicFormatter output the stack inside
// an HTML page. This has been largely inspired by
// https://github.com/go-martini/martini/pull/156/commits.
type HTMLPanicFormatter struct***REMOVED******REMOVED***

func (t *HTMLPanicFormatter) FormatPanicError(rw http.ResponseWriter, r *http.Request, infos *PanicInformation) ***REMOVED***
	if rw.Header().Get("Content-Type") == "" ***REMOVED***
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	***REMOVED***
	panicHTMLTemplate.Execute(rw, infos)
***REMOVED***

// Recovery is a Negroni middleware that recovers from any panics and writes a 500 if there was one.
type Recovery struct ***REMOVED***
	Logger           ALogger
	PrintStack       bool
	PanicHandlerFunc func(*PanicInformation)
	StackAll         bool
	StackSize        int
	Formatter        PanicFormatter

	// Deprecated: Use PanicHandlerFunc instead to receive panic
	// error with additional information (see PanicInformation)
	ErrorHandlerFunc func(interface***REMOVED******REMOVED***)
***REMOVED***

// NewRecovery returns a new instance of Recovery
func NewRecovery() *Recovery ***REMOVED***
	return &Recovery***REMOVED***
		Logger:     log.New(os.Stdout, "[negroni] ", 0),
		PrintStack: true,
		StackAll:   false,
		StackSize:  1024 * 8,
		Formatter:  &TextPanicFormatter***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

func (rec *Recovery) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) ***REMOVED***
	defer func() ***REMOVED***
		if err := recover(); err != nil ***REMOVED***
			rw.WriteHeader(http.StatusInternalServerError)

			stack := make([]byte, rec.StackSize)
			stack = stack[:runtime.Stack(stack, rec.StackAll)]
			infos := &PanicInformation***REMOVED***RecoveredPanic: err, Request: r***REMOVED***

			if rec.PrintStack ***REMOVED***
				infos.Stack = stack
			***REMOVED***
			rec.Logger.Printf(panicText, err, stack)
			rec.Formatter.FormatPanicError(rw, r, infos)

			if rec.ErrorHandlerFunc != nil ***REMOVED***
				func() ***REMOVED***
					defer func() ***REMOVED***
						if err := recover(); err != nil ***REMOVED***
							rec.Logger.Printf("provided ErrorHandlerFunc panic'd: %s, trace:\n%s", err, debug.Stack())
							rec.Logger.Printf("%s\n", debug.Stack())
						***REMOVED***
					***REMOVED***()
					rec.ErrorHandlerFunc(err)
				***REMOVED***()
			***REMOVED***
			if rec.PanicHandlerFunc != nil ***REMOVED***
				func() ***REMOVED***
					defer func() ***REMOVED***
						if err := recover(); err != nil ***REMOVED***
							rec.Logger.Printf("provided PanicHandlerFunc panic'd: %s, trace:\n%s", err, debug.Stack())
							rec.Logger.Printf("%s\n", debug.Stack())
						***REMOVED***
					***REMOVED***()
					rec.PanicHandlerFunc(infos)
				***REMOVED***()
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	next(rw, r)
***REMOVED***
