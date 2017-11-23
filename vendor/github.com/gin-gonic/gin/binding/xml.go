// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"encoding/xml"
	"net/http"
)

type xmlBinding struct***REMOVED******REMOVED***

func (xmlBinding) Name() string ***REMOVED***
	return "xml"
***REMOVED***

func (xmlBinding) Bind(req *http.Request, obj interface***REMOVED******REMOVED***) error ***REMOVED***
	decoder := xml.NewDecoder(req.Body)
	if err := decoder.Decode(obj); err != nil ***REMOVED***
		return err
	***REMOVED***
	return validate(obj)
***REMOVED***
