// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"fmt"
	"net/http"
)

type Redirect struct ***REMOVED***
	Code     int
	Request  *http.Request
	Location string
***REMOVED***

func (r Redirect) Render(w http.ResponseWriter) error ***REMOVED***
	if (r.Code < 300 || r.Code > 308) && r.Code != 201 ***REMOVED***
		panic(fmt.Sprintf("Cannot redirect with status code %d", r.Code))
	***REMOVED***
	http.Redirect(w, r.Request, r.Location, r.Code)
	return nil
***REMOVED***

func (r Redirect) WriteContentType(http.ResponseWriter) ***REMOVED******REMOVED***
