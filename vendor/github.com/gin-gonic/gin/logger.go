// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mattn/go-isatty"
)

var (
	green        = string([]byte***REMOVED***27, 91, 57, 55, 59, 52, 50, 109***REMOVED***)
	white        = string([]byte***REMOVED***27, 91, 57, 48, 59, 52, 55, 109***REMOVED***)
	yellow       = string([]byte***REMOVED***27, 91, 57, 55, 59, 52, 51, 109***REMOVED***)
	red          = string([]byte***REMOVED***27, 91, 57, 55, 59, 52, 49, 109***REMOVED***)
	blue         = string([]byte***REMOVED***27, 91, 57, 55, 59, 52, 52, 109***REMOVED***)
	magenta      = string([]byte***REMOVED***27, 91, 57, 55, 59, 52, 53, 109***REMOVED***)
	cyan         = string([]byte***REMOVED***27, 91, 57, 55, 59, 52, 54, 109***REMOVED***)
	reset        = string([]byte***REMOVED***27, 91, 48, 109***REMOVED***)
	disableColor = false
)

func DisableConsoleColor() ***REMOVED***
	disableColor = true
***REMOVED***

func ErrorLogger() HandlerFunc ***REMOVED***
	return ErrorLoggerT(ErrorTypeAny)
***REMOVED***

func ErrorLoggerT(typ ErrorType) HandlerFunc ***REMOVED***
	return func(c *Context) ***REMOVED***
		c.Next()
		errors := c.Errors.ByType(typ)
		if len(errors) > 0 ***REMOVED***
			c.JSON(-1, errors)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Logger instances a Logger middleware that will write the logs to gin.DefaultWriter
// By default gin.DefaultWriter = os.Stdout
func Logger() HandlerFunc ***REMOVED***
	return LoggerWithWriter(DefaultWriter)
***REMOVED***

// LoggerWithWriter instance a Logger middleware with the specified writter buffer.
// Example: os.Stdout, a file opened in write mode, a socket...
func LoggerWithWriter(out io.Writer, notlogged ...string) HandlerFunc ***REMOVED***
	isTerm := true

	if w, ok := out.(*os.File); !ok ||
		(os.Getenv("TERM") == "dumb" || (!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd()))) ||
		disableColor ***REMOVED***
		isTerm = false
	***REMOVED***

	var skip map[string]struct***REMOVED******REMOVED***

	if length := len(notlogged); length > 0 ***REMOVED***
		skip = make(map[string]struct***REMOVED******REMOVED***, length)

		for _, path := range notlogged ***REMOVED***
			skip[path] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	return func(c *Context) ***REMOVED***
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok ***REMOVED***
			// Stop timer
			end := time.Now()
			latency := end.Sub(start)

			clientIP := c.ClientIP()
			method := c.Request.Method
			statusCode := c.Writer.Status()
			var statusColor, methodColor string
			if isTerm ***REMOVED***
				statusColor = colorForStatus(statusCode)
				methodColor = colorForMethod(method)
			***REMOVED***
			comment := c.Errors.ByType(ErrorTypePrivate).String()

			fmt.Fprintf(out, "[GIN] %v |%s %3d %s| %13v | %15s |%s  %s %-7s %s\n%s",
				end.Format("2006/01/02 - 15:04:05"),
				statusColor, statusCode, reset,
				latency,
				clientIP,
				methodColor, method, reset,
				path,
				comment,
			)
		***REMOVED***
	***REMOVED***
***REMOVED***

func colorForStatus(code int) string ***REMOVED***
	switch ***REMOVED***
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	***REMOVED***
***REMOVED***

func colorForMethod(method string) string ***REMOVED***
	switch method ***REMOVED***
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	***REMOVED***
***REMOVED***
