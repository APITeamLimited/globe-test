// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import "net/http"

type Data struct ***REMOVED***
	ContentType string
	Data        []byte
***REMOVED***

// Render (Data) writes data with custom ContentType
func (r Data) Render(w http.ResponseWriter) (err error) ***REMOVED***
	r.WriteContentType(w)
	_, err = w.Write(r.Data)
	return
***REMOVED***

func (r Data) WriteContentType(w http.ResponseWriter) ***REMOVED***
	writeContentType(w, []string***REMOVED***r.ContentType***REMOVED***)
***REMOVED***
