/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package http

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/dop251/goja"

	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/js/modules/k6/html"
	"github.com/loadimpact/k6/lib/netext/httpext"
)

// Response is a representation of an HTTP response to be returned to the goja VM
type Response struct ***REMOVED***
	*httpext.Response `js:"-"`
***REMOVED***

func responseFromHttpext(resp *httpext.Response) *Response ***REMOVED***
	res := Response***REMOVED***resp***REMOVED***
	return &res
***REMOVED***

// JSON parses the body of a response as json and returns it to the goja VM
func (res *Response) JSON(selector ...string) goja.Value ***REMOVED***
	v, err := res.Response.JSON(selector...)
	if err != nil ***REMOVED***
		common.Throw(common.GetRuntime(res.GetCtx()), err)
	***REMOVED***
	if v == nil ***REMOVED***
		return goja.Undefined()
	***REMOVED***
	return common.GetRuntime(res.GetCtx()).ToValue(v)
***REMOVED***

// HTML returns the body as an html.Selection
func (res *Response) HTML(selector ...string) html.Selection ***REMOVED***
	var body string
	switch b := res.Body.(type) ***REMOVED***
	case []byte:
		body = string(b)
	case string:
		body = b
	default:
		common.Throw(common.GetRuntime(res.GetCtx()), errors.New("invalid response type"))
	***REMOVED***

	sel, err := html.HTML***REMOVED******REMOVED***.ParseHTML(res.GetCtx(), body)
	if err != nil ***REMOVED***
		common.Throw(common.GetRuntime(res.GetCtx()), err)
	***REMOVED***
	sel.URL = res.URL
	if len(selector) > 0 ***REMOVED***
		sel = sel.Find(selector[0])
	***REMOVED***
	return sel
***REMOVED***

// SubmitForm parses the body as an html looking for a from and then submitting it
// TODO: document the actual arguments that can be provided
func (res *Response) SubmitForm(args ...goja.Value) (*Response, error) ***REMOVED***
	rt := common.GetRuntime(res.GetCtx())

	formSelector := "form"
	submitSelector := "[type=\"submit\"]"
	var fields map[string]goja.Value
	requestParams := goja.Null()
	if len(args) > 0 ***REMOVED***
		params := args[0].ToObject(rt)
		for _, k := range params.Keys() ***REMOVED***
			switch k ***REMOVED***
			case "formSelector":
				formSelector = params.Get(k).String()
			case "submitSelector":
				submitSelector = params.Get(k).String()
			case "fields":
				if rt.ExportTo(params.Get(k), &fields) != nil ***REMOVED***
					fields = nil
				***REMOVED***
			case "params":
				requestParams = params.Get(k)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	form := res.HTML(formSelector)
	if form.Size() == 0 ***REMOVED***
		common.Throw(rt, fmt.Errorf("no form found for selector '%s' in response '%s'", formSelector, res.URL))
	***REMOVED***

	methodAttr := form.Attr("method")
	var requestMethod string
	if methodAttr == goja.Undefined() ***REMOVED***
		// Use GET by default
		requestMethod = HTTP_METHOD_GET
	***REMOVED*** else ***REMOVED***
		requestMethod = strings.ToUpper(methodAttr.String())
	***REMOVED***

	responseURL, err := url.Parse(res.URL)
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***

	actionAttr := form.Attr("action")
	var requestURL *url.URL
	if actionAttr == goja.Undefined() ***REMOVED***
		// Use the url of the response if no action is set
		requestURL = responseURL
	***REMOVED*** else ***REMOVED***
		actionURL, err := url.Parse(actionAttr.String())
		if err != nil ***REMOVED***
			common.Throw(rt, err)
		***REMOVED***
		requestURL = responseURL.ResolveReference(actionURL)
	***REMOVED***

	// Set the body based on the form values
	values := form.SerializeObject()

	// Set the name + value of the submit button
	submit := form.Find(submitSelector)
	submitName := submit.Attr("name")
	submitValue := submit.Val()
	if submitName != goja.Undefined() && submitValue != goja.Undefined() ***REMOVED***
		values[submitName.String()] = submitValue
	***REMOVED***

	// Set the values supplied in the arguments, overriding automatically set values
	for k, v := range fields ***REMOVED***
		values[k] = v
	***REMOVED***

	if requestMethod == HTTP_METHOD_GET ***REMOVED***
		q := url.Values***REMOVED******REMOVED***
		for k, v := range values ***REMOVED***
			q.Add(k, v.String())
		***REMOVED***
		requestURL.RawQuery = q.Encode()
		return New().Request(res.GetCtx(), requestMethod, rt.ToValue(requestURL.String()), goja.Null(), requestParams)
	***REMOVED***
	return New().Request(res.GetCtx(), requestMethod, rt.ToValue(requestURL.String()), rt.ToValue(values), requestParams)
***REMOVED***

// ClickLink parses the body as an html, looks for a link and than makes a request as if the link was
// clicked
func (res *Response) ClickLink(args ...goja.Value) (*Response, error) ***REMOVED***
	rt := common.GetRuntime(res.GetCtx())

	selector := "a[href]"
	requestParams := goja.Null()
	if len(args) > 0 ***REMOVED***
		params := args[0].ToObject(rt)
		for _, k := range params.Keys() ***REMOVED***
			switch k ***REMOVED***
			case "selector":
				selector = params.Get(k).String()
			case "params":
				requestParams = params.Get(k)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	responseURL, err := url.Parse(res.URL)
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***

	link := res.HTML(selector)
	if link.Size() == 0 ***REMOVED***
		common.Throw(rt, fmt.Errorf("no element found for selector '%s' in response '%s'", selector, res.URL))
	***REMOVED***
	hrefAttr := link.Attr("href")
	if hrefAttr == goja.Undefined() ***REMOVED***
		common.Throw(rt, fmt.Errorf("no valid href attribute value found on element '%s' in response '%s'", selector, res.URL))
	***REMOVED***
	hrefURL, err := url.Parse(hrefAttr.String())
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***
	requestURL := responseURL.ResolveReference(hrefURL)

	return New().Get(res.GetCtx(), rt.ToValue(requestURL.String()), requestParams)
***REMOVED***
