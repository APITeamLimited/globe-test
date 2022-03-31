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

package html

import (
	neturl "net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
)

type FormValue struct ***REMOVED***
	Name  string
	Value goja.Value
***REMOVED***

func (s Selection) SerializeArray() []FormValue ***REMOVED***
	submittableSelector := "input,select,textarea,keygen"
	var formElements *goquery.Selection
	if s.sel.Is("form") ***REMOVED***
		formElements = s.sel.Find(submittableSelector)
	***REMOVED*** else ***REMOVED***
		formElements = s.sel.Filter(submittableSelector)
	***REMOVED***

	formElements = formElements.FilterFunction(func(_ int, sel *goquery.Selection) bool ***REMOVED***
		name := sel.AttrOr("name", "")
		inputType := sel.AttrOr("type", "")
		disabled := sel.AttrOr("disabled", "")
		checked := sel.AttrOr("checked", "")

		return name != "" && // Must have a non-empty name
			disabled != "disabled" && // Must not be disabled
			inputType != "submit" && // Must not be a button
			inputType != "button" &&
			inputType != "reset" &&
			inputType != "image" && // Must not be an image or file
			inputType != "file" &&
			(checked == "checked" ||
				(inputType != "checkbox" && inputType != "radio")) // Must be checked if it is an checkbox or radio
	***REMOVED***)

	result := make([]FormValue, len(formElements.Nodes))
	formElements.Each(func(i int, sel *goquery.Selection) ***REMOVED***
		element := Selection***REMOVED***s.rt, sel, s.URL***REMOVED***
		name, _ := sel.Attr("name")
		result[i] = FormValue***REMOVED***Name: name, Value: element.Val()***REMOVED***
	***REMOVED***)
	return result
***REMOVED***

func (s Selection) SerializeObject() map[string]goja.Value ***REMOVED***
	formValues := s.SerializeArray()
	result := make(map[string]goja.Value)
	for i := range formValues ***REMOVED***
		formValue := formValues[i]
		result[formValue.Name] = formValue.Value
	***REMOVED***

	return result
***REMOVED***

func (s Selection) Serialize() string ***REMOVED***
	formValues := s.SerializeArray()
	urlValues := make(neturl.Values, len(formValues))
	for i := range formValues ***REMOVED***
		formValue := formValues[i]
		value := formValue.Value.Export()
		switch v := value.(type) ***REMOVED***
		case string:
			urlValues.Set(formValue.Name, v)
		case []string:
			urlValues[formValue.Name] = v
		***REMOVED***
	***REMOVED***
	return urlValues.Encode()
***REMOVED***
