// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"encoding/xml"
	"net/http"
)

type XML struct ***REMOVED***
	Data interface***REMOVED******REMOVED***
***REMOVED***

var xmlContentType = []string***REMOVED***"application/xml; charset=utf-8"***REMOVED***

func (r XML) Render(w http.ResponseWriter) error ***REMOVED***
	r.WriteContentType(w)
	return xml.NewEncoder(w).Encode(r.Data)
***REMOVED***

func (r XML) WriteContentType(w http.ResponseWriter) ***REMOVED***
	writeContentType(w, xmlContentType)
***REMOVED***
