// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http/httputil"
	"runtime"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() HandlerFunc ***REMOVED***
	return RecoveryWithWriter(DefaultErrorWriter)
***REMOVED***

func RecoveryWithWriter(out io.Writer) HandlerFunc ***REMOVED***
	var logger *log.Logger
	if out != nil ***REMOVED***
		logger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
	***REMOVED***
	return func(c *Context) ***REMOVED***
		defer func() ***REMOVED***
			if err := recover(); err != nil ***REMOVED***
				if logger != nil ***REMOVED***
					stack := stack(3)
					httprequest, _ := httputil.DumpRequest(c.Request, false)
					logger.Printf("[Recovery] panic recovered:\n%s\n%s\n%s%s", string(httprequest), err, stack, reset)
				***REMOVED***
				c.AbortWithStatus(500)
			***REMOVED***
		***REMOVED***()
		c.Next()
	***REMOVED***
***REMOVED***

// stack returns a nicely formated stack frame, skipping skip frames
func stack(skip int) []byte ***REMOVED***
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ ***REMOVED*** // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok ***REMOVED***
			break
		***REMOVED***
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile ***REMOVED***
			data, err := ioutil.ReadFile(file)
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			lines = bytes.Split(data, []byte***REMOVED***'\n'***REMOVED***)
			lastFile = file
		***REMOVED***
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	***REMOVED***
	return buf.Bytes()
***REMOVED***

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte ***REMOVED***
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) ***REMOVED***
		return dunno
	***REMOVED***
	return bytes.TrimSpace(lines[n])
***REMOVED***

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte ***REMOVED***
	fn := runtime.FuncForPC(pc)
	if fn == nil ***REMOVED***
		return dunno
	***REMOVED***
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 ***REMOVED***
		name = name[lastslash+1:]
	***REMOVED***
	if period := bytes.Index(name, dot); period >= 0 ***REMOVED***
		name = name[period+1:]
	***REMOVED***
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
***REMOVED***
