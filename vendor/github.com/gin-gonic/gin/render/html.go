// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"html/template"
	"net/http"
)

type (
	Delims struct ***REMOVED***
		Left  string
		Right string
	***REMOVED***

	HTMLRender interface ***REMOVED***
		Instance(string, interface***REMOVED******REMOVED***) Render
	***REMOVED***

	HTMLProduction struct ***REMOVED***
		Template *template.Template
		Delims   Delims
	***REMOVED***

	HTMLDebug struct ***REMOVED***
		Files   []string
		Glob    string
		Delims  Delims
		FuncMap template.FuncMap
	***REMOVED***

	HTML struct ***REMOVED***
		Template *template.Template
		Name     string
		Data     interface***REMOVED******REMOVED***
	***REMOVED***
)

var htmlContentType = []string***REMOVED***"text/html; charset=utf-8"***REMOVED***

func (r HTMLProduction) Instance(name string, data interface***REMOVED******REMOVED***) Render ***REMOVED***
	return HTML***REMOVED***
		Template: r.Template,
		Name:     name,
		Data:     data,
	***REMOVED***
***REMOVED***

func (r HTMLDebug) Instance(name string, data interface***REMOVED******REMOVED***) Render ***REMOVED***
	return HTML***REMOVED***
		Template: r.loadTemplate(),
		Name:     name,
		Data:     data,
	***REMOVED***
***REMOVED***
func (r HTMLDebug) loadTemplate() *template.Template ***REMOVED***
	if r.FuncMap == nil ***REMOVED***
		r.FuncMap = template.FuncMap***REMOVED******REMOVED***
	***REMOVED***
	if len(r.Files) > 0 ***REMOVED***
		return template.Must(template.New("").Delims(r.Delims.Left, r.Delims.Right).Funcs(r.FuncMap).ParseFiles(r.Files...))
	***REMOVED***
	if len(r.Glob) > 0 ***REMOVED***
		return template.Must(template.New("").Delims(r.Delims.Left, r.Delims.Right).Funcs(r.FuncMap).ParseGlob(r.Glob))
	***REMOVED***
	panic("the HTML debug render was created without files or glob pattern")
***REMOVED***

func (r HTML) Render(w http.ResponseWriter) error ***REMOVED***
	r.WriteContentType(w)

	if len(r.Name) == 0 ***REMOVED***
		return r.Template.Execute(w, r.Data)
	***REMOVED***
	return r.Template.ExecuteTemplate(w, r.Name, r.Data)
***REMOVED***

func (r HTML) WriteContentType(w http.ResponseWriter) ***REMOVED***
	writeContentType(w, htmlContentType)
***REMOVED***
