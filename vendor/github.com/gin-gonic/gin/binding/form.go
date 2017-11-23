// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import "net/http"

type formBinding struct***REMOVED******REMOVED***
type formPostBinding struct***REMOVED******REMOVED***
type formMultipartBinding struct***REMOVED******REMOVED***

func (formBinding) Name() string ***REMOVED***
	return "form"
***REMOVED***

func (formBinding) Bind(req *http.Request, obj interface***REMOVED******REMOVED***) error ***REMOVED***
	if err := req.ParseForm(); err != nil ***REMOVED***
		return err
	***REMOVED***
	req.ParseMultipartForm(32 << 10) // 32 MB
	if err := mapForm(obj, req.Form); err != nil ***REMOVED***
		return err
	***REMOVED***
	return validate(obj)
***REMOVED***

func (formPostBinding) Name() string ***REMOVED***
	return "form-urlencoded"
***REMOVED***

func (formPostBinding) Bind(req *http.Request, obj interface***REMOVED******REMOVED***) error ***REMOVED***
	if err := req.ParseForm(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := mapForm(obj, req.PostForm); err != nil ***REMOVED***
		return err
	***REMOVED***
	return validate(obj)
***REMOVED***

func (formMultipartBinding) Name() string ***REMOVED***
	return "multipart/form-data"
***REMOVED***

func (formMultipartBinding) Bind(req *http.Request, obj interface***REMOVED******REMOVED***) error ***REMOVED***
	if err := req.ParseMultipartForm(32 << 10); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := mapForm(obj, req.MultipartForm.Value); err != nil ***REMOVED***
		return err
	***REMOVED***
	return validate(obj)
***REMOVED***
