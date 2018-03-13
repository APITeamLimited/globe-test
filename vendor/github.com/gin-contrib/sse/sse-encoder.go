// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sse

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// Server-Sent Events
// W3C Working Draft 29 October 2009
// http://www.w3.org/TR/2009/WD-eventsource-20091029/

const ContentType = "text/event-stream"

var contentType = []string***REMOVED***ContentType***REMOVED***
var noCache = []string***REMOVED***"no-cache"***REMOVED***

var fieldReplacer = strings.NewReplacer(
	"\n", "\\n",
	"\r", "\\r")

var dataReplacer = strings.NewReplacer(
	"\n", "\ndata:",
	"\r", "\\r")

type Event struct ***REMOVED***
	Event string
	Id    string
	Retry uint
	Data  interface***REMOVED******REMOVED***
***REMOVED***

func Encode(writer io.Writer, event Event) error ***REMOVED***
	w := checkWriter(writer)
	writeId(w, event.Id)
	writeEvent(w, event.Event)
	writeRetry(w, event.Retry)
	return writeData(w, event.Data)
***REMOVED***

func writeId(w stringWriter, id string) ***REMOVED***
	if len(id) > 0 ***REMOVED***
		w.WriteString("id:")
		fieldReplacer.WriteString(w, id)
		w.WriteString("\n")
	***REMOVED***
***REMOVED***

func writeEvent(w stringWriter, event string) ***REMOVED***
	if len(event) > 0 ***REMOVED***
		w.WriteString("event:")
		fieldReplacer.WriteString(w, event)
		w.WriteString("\n")
	***REMOVED***
***REMOVED***

func writeRetry(w stringWriter, retry uint) ***REMOVED***
	if retry > 0 ***REMOVED***
		w.WriteString("retry:")
		w.WriteString(strconv.FormatUint(uint64(retry), 10))
		w.WriteString("\n")
	***REMOVED***
***REMOVED***

func writeData(w stringWriter, data interface***REMOVED******REMOVED***) error ***REMOVED***
	w.WriteString("data:")
	switch kindOfData(data) ***REMOVED***
	case reflect.Struct, reflect.Slice, reflect.Map:
		err := json.NewEncoder(w).Encode(data)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		w.WriteString("\n")
	default:
		dataReplacer.WriteString(w, fmt.Sprint(data))
		w.WriteString("\n\n")
	***REMOVED***
	return nil
***REMOVED***

func (r Event) Render(w http.ResponseWriter) error ***REMOVED***
	r.WriteContentType(w)
	return Encode(w, r)
***REMOVED***

func (r Event) WriteContentType(w http.ResponseWriter) ***REMOVED***
	header := w.Header()
	header["Content-Type"] = contentType

	if _, exist := header["Cache-Control"]; !exist ***REMOVED***
		header["Cache-Control"] = noCache
	***REMOVED***
***REMOVED***

func kindOfData(data interface***REMOVED******REMOVED***) reflect.Kind ***REMOVED***
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr ***REMOVED***
		valueType = value.Elem().Kind()
	***REMOVED***
	return valueType
***REMOVED***
