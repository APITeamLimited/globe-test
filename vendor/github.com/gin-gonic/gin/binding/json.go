// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"encoding/json"
	"net/http"
)

type jsonBinding struct***REMOVED******REMOVED***

func (jsonBinding) Name() string ***REMOVED***
	return "json"
***REMOVED***

func (jsonBinding) Bind(req *http.Request, obj interface***REMOVED******REMOVED***) error ***REMOVED***
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(obj); err != nil ***REMOVED***
		return err
	***REMOVED***
	return validate(obj)
***REMOVED***
