// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import "net/http"

type Render interface ***REMOVED***
	Render(http.ResponseWriter) error
	WriteContentType(w http.ResponseWriter)
***REMOVED***

var (
	_ Render     = JSON***REMOVED******REMOVED***
	_ Render     = IndentedJSON***REMOVED******REMOVED***
	_ Render     = XML***REMOVED******REMOVED***
	_ Render     = String***REMOVED******REMOVED***
	_ Render     = Redirect***REMOVED******REMOVED***
	_ Render     = Data***REMOVED******REMOVED***
	_ Render     = HTML***REMOVED******REMOVED***
	_ HTMLRender = HTMLDebug***REMOVED******REMOVED***
	_ HTMLRender = HTMLProduction***REMOVED******REMOVED***
	_ Render     = YAML***REMOVED******REMOVED***
	_ Render     = MsgPack***REMOVED******REMOVED***
	_ Render     = MsgPack***REMOVED******REMOVED***
)

func writeContentType(w http.ResponseWriter, value []string) ***REMOVED***
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 ***REMOVED***
		header["Content-Type"] = value
	***REMOVED***
***REMOVED***
