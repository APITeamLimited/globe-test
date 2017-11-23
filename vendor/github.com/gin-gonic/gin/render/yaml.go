// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

type YAML struct ***REMOVED***
	Data interface***REMOVED******REMOVED***
***REMOVED***

var yamlContentType = []string***REMOVED***"application/x-yaml; charset=utf-8"***REMOVED***

func (r YAML) Render(w http.ResponseWriter) error ***REMOVED***
	r.WriteContentType(w)

	bytes, err := yaml.Marshal(r.Data)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	w.Write(bytes)
	return nil
***REMOVED***

func (r YAML) WriteContentType(w http.ResponseWriter) ***REMOVED***
	writeContentType(w, yamlContentType)
***REMOVED***
