// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"encoding/json"
	"net/http"
)

type (
	JSON struct ***REMOVED***
		Data interface***REMOVED******REMOVED***
	***REMOVED***

	IndentedJSON struct ***REMOVED***
		Data interface***REMOVED******REMOVED***
	***REMOVED***
)

var jsonContentType = []string***REMOVED***"application/json; charset=utf-8"***REMOVED***

func (r JSON) Render(w http.ResponseWriter) (err error) ***REMOVED***
	if err = WriteJSON(w, r.Data); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return
***REMOVED***

func (r JSON) WriteContentType(w http.ResponseWriter) ***REMOVED***
	writeContentType(w, jsonContentType)
***REMOVED***

func WriteJSON(w http.ResponseWriter, obj interface***REMOVED******REMOVED***) error ***REMOVED***
	writeContentType(w, jsonContentType)
	jsonBytes, err := json.Marshal(obj)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	w.Write(jsonBytes)
	return nil
***REMOVED***

func (r IndentedJSON) Render(w http.ResponseWriter) error ***REMOVED***
	r.WriteContentType(w)
	jsonBytes, err := json.MarshalIndent(r.Data, "", "    ")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	w.Write(jsonBytes)
	return nil
***REMOVED***

func (r IndentedJSON) WriteContentType(w http.ResponseWriter) ***REMOVED***
	writeContentType(w, jsonContentType)
***REMOVED***
