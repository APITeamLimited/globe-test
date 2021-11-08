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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/dop251/goja"
	"github.com/tidwall/gjson"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules/k6/html"
	"go.k6.io/k6/lib/netext/httpext"
)

// Response is a representation of an HTTP response to be returned to the goja VM
type Response struct ***REMOVED***
	*httpext.Response `js:"-"`
	client            *Client

	cachedJSON    interface***REMOVED******REMOVED***
	validatedJSON bool
***REMOVED***

type jsonError struct ***REMOVED***
	line      int
	character int
	err       error
***REMOVED***

func (j jsonError) Error() string ***REMOVED***
	errMessage := "cannot parse json due to an error at line"
	return fmt.Sprintf("%s %d, character %d , error: %v", errMessage, j.line, j.character, j.err)
***REMOVED***

// HTML returns the body as an html.Selection
func (res *Response) HTML(selector ...string) html.Selection ***REMOVED***
	rt := res.client.moduleInstance.vu.Runtime()
	if res.Body == nil ***REMOVED***
		err := fmt.Errorf("the body is null so we can't transform it to HTML" +
			" - this likely was because of a request error getting the response")
		common.Throw(rt, err)
	***REMOVED***

	body, err := common.ToString(res.Body)
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***

	sel, err := html.ParseHTML(rt, body)
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***
	sel.URL = res.URL
	if len(selector) > 0 ***REMOVED***
		sel = sel.Find(selector[0])
	***REMOVED***
	return sel
***REMOVED***

// JSON parses the body of a response as JSON and returns it to the goja VM.
func (res *Response) JSON(selector ...string) goja.Value ***REMOVED***
	rt := res.client.moduleInstance.vu.Runtime()

	if res.Body == nil ***REMOVED***
		err := fmt.Errorf("the body is null so we can't transform it to JSON" +
			" - this likely was because of a request error getting the response")
		common.Throw(rt, err)
	***REMOVED***

	hasSelector := len(selector) > 0
	if res.cachedJSON == nil || hasSelector ***REMOVED*** //nolint:nestif
		var v interface***REMOVED******REMOVED***

		body, err := common.ToBytes(res.Body)
		if err != nil ***REMOVED***
			common.Throw(rt, err)
		***REMOVED***

		if hasSelector ***REMOVED***
			if !res.validatedJSON ***REMOVED***
				if !gjson.ValidBytes(body) ***REMOVED***
					return goja.Undefined()
				***REMOVED***
				res.validatedJSON = true
			***REMOVED***

			result := gjson.GetBytes(body, selector[0])

			if !result.Exists() ***REMOVED***
				return goja.Undefined()
			***REMOVED***
			return rt.ToValue(result.Value())
		***REMOVED***

		if err := json.Unmarshal(body, &v); err != nil ***REMOVED***
			var syntaxError *json.SyntaxError
			if errors.As(err, &syntaxError) ***REMOVED***
				err = checkErrorInJSON(body, int(syntaxError.Offset), err)
			***REMOVED***
			common.Throw(rt, err)
		***REMOVED***
		res.validatedJSON = true
		res.cachedJSON = v
	***REMOVED***

	return rt.ToValue(res.cachedJSON)
***REMOVED***

func checkErrorInJSON(input []byte, offset int, err error) error ***REMOVED***
	lf := '\n'
	str := string(input)

	// Humans tend to count from 1.
	line := 1
	character := 0

	for i, b := range str ***REMOVED***
		if b == lf ***REMOVED***
			line++
			character = 0
		***REMOVED***
		character++
		if i == offset ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	return jsonError***REMOVED***line: line, character: character, err: err***REMOVED***
***REMOVED***

// SubmitForm parses the body as an html looking for a from and then submitting it
// TODO: document the actual arguments that can be provided
func (res *Response) SubmitForm(args ...goja.Value) (*Response, error) ***REMOVED***
	rt := res.client.moduleInstance.vu.Runtime()

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
		requestMethod = http.MethodGet
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

	if requestMethod == http.MethodGet ***REMOVED***
		q := url.Values***REMOVED******REMOVED***
		for k, v := range values ***REMOVED***
			q.Add(k, v.String())
		***REMOVED***
		requestURL.RawQuery = q.Encode()
		return res.client.Request(requestMethod, rt.ToValue(requestURL.String()), goja.Null(), requestParams)
	***REMOVED***
	return res.client.Request(
		requestMethod, rt.ToValue(requestURL.String()),
		rt.ToValue(values), requestParams,
	)
***REMOVED***

// ClickLink parses the body as an html, looks for a link and than makes a request as if the link was
// clicked
func (res *Response) ClickLink(args ...goja.Value) (*Response, error) ***REMOVED***
	rt := res.client.moduleInstance.vu.Runtime()

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

	return res.client.Request(http.MethodGet, rt.ToValue(requestURL.String()), goja.Undefined(), requestParams)
***REMOVED***
