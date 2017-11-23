// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"fmt"
	"io"
	"net/http"
)

type String struct ***REMOVED***
	Format string
	Data   []interface***REMOVED******REMOVED***
***REMOVED***

var plainContentType = []string***REMOVED***"text/plain; charset=utf-8"***REMOVED***

func (r String) Render(w http.ResponseWriter) error ***REMOVED***
	WriteString(w, r.Format, r.Data)
	return nil
***REMOVED***

func (r String) WriteContentType(w http.ResponseWriter) ***REMOVED***
	writeContentType(w, plainContentType)
***REMOVED***

func WriteString(w http.ResponseWriter, format string, data []interface***REMOVED******REMOVED***) ***REMOVED***
	writeContentType(w, plainContentType)
	if len(data) > 0 ***REMOVED***
		fmt.Fprintf(w, format, data...)
	***REMOVED*** else ***REMOVED***
		io.WriteString(w, format)
	***REMOVED***
***REMOVED***
