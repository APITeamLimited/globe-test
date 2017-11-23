package negroni

import (
	"bytes"

	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

// LoggerEntry is the structure
// passed to the template.
type LoggerEntry struct ***REMOVED***
	StartTime string
	Status    int
	Duration  time.Duration
	Hostname  string
	Method    string
	Path      string
	Request   *http.Request
***REMOVED***

// LoggerDefaultFormat is the format
// logged used by the default Logger instance.
var LoggerDefaultFormat = "***REMOVED******REMOVED***.StartTime***REMOVED******REMOVED*** | ***REMOVED******REMOVED***.Status***REMOVED******REMOVED*** | \t ***REMOVED******REMOVED***.Duration***REMOVED******REMOVED*** | ***REMOVED******REMOVED***.Hostname***REMOVED******REMOVED*** | ***REMOVED******REMOVED***.Method***REMOVED******REMOVED*** ***REMOVED******REMOVED***.Path***REMOVED******REMOVED*** \n"

// LoggerDefaultDateFormat is the
// format used for date by the
// default Logger instance.
var LoggerDefaultDateFormat = time.RFC3339

// ALogger interface
type ALogger interface ***REMOVED***
	Println(v ...interface***REMOVED******REMOVED***)
	Printf(format string, v ...interface***REMOVED******REMOVED***)
***REMOVED***

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct ***REMOVED***
	// ALogger implements just enough log.Logger interface to be compatible with other implementations
	ALogger
	dateFormat string
	template   *template.Template
***REMOVED***

// NewLogger returns a new Logger instance
func NewLogger() *Logger ***REMOVED***
	logger := &Logger***REMOVED***ALogger: log.New(os.Stdout, "[negroni] ", 0), dateFormat: LoggerDefaultDateFormat***REMOVED***
	logger.SetFormat(LoggerDefaultFormat)
	return logger
***REMOVED***

func (l *Logger) SetFormat(format string) ***REMOVED***
	l.template = template.Must(template.New("negroni_parser").Parse(format))
***REMOVED***

func (l *Logger) SetDateFormat(format string) ***REMOVED***
	l.dateFormat = format
***REMOVED***

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) ***REMOVED***
	start := time.Now()

	next(rw, r)

	res := rw.(ResponseWriter)
	log := LoggerEntry***REMOVED***
		StartTime: start.Format(l.dateFormat),
		Status:    res.Status(),
		Duration:  time.Since(start),
		Hostname:  r.Host,
		Method:    r.Method,
		Path:      r.URL.Path,
		Request:   r,
	***REMOVED***

	buff := &bytes.Buffer***REMOVED******REMOVED***
	l.template.Execute(buff, log)
	l.Printf(buff.String())
***REMOVED***
