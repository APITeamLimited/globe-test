package sse

import "io"

type stringWriter interface ***REMOVED***
	io.Writer
	WriteString(string) (int, error)
***REMOVED***

type stringWrapper struct ***REMOVED***
	io.Writer
***REMOVED***

func (w stringWrapper) WriteString(str string) (int, error) ***REMOVED***
	return w.Writer.Write([]byte(str))
***REMOVED***

func checkWriter(writer io.Writer) stringWriter ***REMOVED***
	if w, ok := writer.(stringWriter); ok ***REMOVED***
		return w
	***REMOVED*** else ***REMOVED***
		return stringWrapper***REMOVED***writer***REMOVED***
	***REMOVED***
***REMOVED***
