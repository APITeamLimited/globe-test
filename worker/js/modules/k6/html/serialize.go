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
