// Copyright © 2013 Steve Francia <spf@spf13.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Commands similar to git, go tools and other modern CLI tools
// inspired by go, go-Commander, gh and subcommand

package cobra

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode"
)

var templateFuncs = template.FuncMap***REMOVED***
	"trim":                    strings.TrimSpace,
	"trimRightSpace":          trimRightSpace,
	"trimTrailingWhitespaces": trimRightSpace,
	"appendIfNotPresent":      appendIfNotPresent,
	"rpad":                    rpad,
	"gt":                      Gt,
	"eq":                      Eq,
***REMOVED***

var initializers []func()

// EnablePrefixMatching allows to set automatic prefix matching. Automatic prefix matching can be a dangerous thing
// to automatically enable in CLI tools.
// Set this to true to enable it.
var EnablePrefixMatching = false

// EnableCommandSorting controls sorting of the slice of commands, which is turned on by default.
// To disable sorting, set it to false.
var EnableCommandSorting = true

// MousetrapHelpText enables an information splash screen on Windows
// if the CLI is started from explorer.exe.
// To disable the mousetrap, just set this variable to blank string ("").
// Works only on Microsoft Windows.
var MousetrapHelpText = `This is a command line tool.

You need to open cmd.exe and run it from there.
`

// MousetrapDisplayDuration controls how long the MousetrapHelpText message is displayed on Windows
// if the CLI is started from explorer.exe. Set to 0 to wait for the return key to be pressed.
// To disable the mousetrap, just set MousetrapHelpText to blank string ("").
// Works only on Microsoft Windows.
var MousetrapDisplayDuration = 5 * time.Second

// AddTemplateFunc adds a template function that's available to Usage and Help
// template generation.
func AddTemplateFunc(name string, tmplFunc interface***REMOVED******REMOVED***) ***REMOVED***
	templateFuncs[name] = tmplFunc
***REMOVED***

// AddTemplateFuncs adds multiple template functions that are available to Usage and
// Help template generation.
func AddTemplateFuncs(tmplFuncs template.FuncMap) ***REMOVED***
	for k, v := range tmplFuncs ***REMOVED***
		templateFuncs[k] = v
	***REMOVED***
***REMOVED***

// OnInitialize sets the passed functions to be run when each command's
// Execute method is called.
func OnInitialize(y ...func()) ***REMOVED***
	initializers = append(initializers, y...)
***REMOVED***

// FIXME Gt is unused by cobra and should be removed in a version 2. It exists only for compatibility with users of cobra.

// Gt takes two types and checks whether the first type is greater than the second. In case of types Arrays, Chans,
// Maps and Slices, Gt will compare their lengths. Ints are compared directly while strings are first parsed as
// ints and then compared.
func Gt(a interface***REMOVED******REMOVED***, b interface***REMOVED******REMOVED***) bool ***REMOVED***
	var left, right int64
	av := reflect.ValueOf(a)

	switch av.Kind() ***REMOVED***
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		left = int64(av.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		left = av.Int()
	case reflect.String:
		left, _ = strconv.ParseInt(av.String(), 10, 64)
	***REMOVED***

	bv := reflect.ValueOf(b)

	switch bv.Kind() ***REMOVED***
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		right = int64(bv.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		right = bv.Int()
	case reflect.String:
		right, _ = strconv.ParseInt(bv.String(), 10, 64)
	***REMOVED***

	return left > right
***REMOVED***

// FIXME Eq is unused by cobra and should be removed in a version 2. It exists only for compatibility with users of cobra.

// Eq takes two types and checks whether they are equal. Supported types are int and string. Unsupported types will panic.
func Eq(a interface***REMOVED******REMOVED***, b interface***REMOVED******REMOVED***) bool ***REMOVED***
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() ***REMOVED***
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		panic("Eq called on unsupported type")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return av.Int() == bv.Int()
	case reflect.String:
		return av.String() == bv.String()
	***REMOVED***
	return false
***REMOVED***

func trimRightSpace(s string) string ***REMOVED***
	return strings.TrimRightFunc(s, unicode.IsSpace)
***REMOVED***

// FIXME appendIfNotPresent is unused by cobra and should be removed in a version 2. It exists only for compatibility with users of cobra.

// appendIfNotPresent will append stringToAppend to the end of s, but only if it's not yet present in s.
func appendIfNotPresent(s, stringToAppend string) string ***REMOVED***
	if strings.Contains(s, stringToAppend) ***REMOVED***
		return s
	***REMOVED***
	return s + " " + stringToAppend
***REMOVED***

// rpad adds padding to the right of a string.
func rpad(s string, padding int) string ***REMOVED***
	template := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(template, s)
***REMOVED***

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface***REMOVED******REMOVED***) error ***REMOVED***
	t := template.New("top")
	t.Funcs(templateFuncs)
	template.Must(t.Parse(text))
	return t.Execute(w, data)
***REMOVED***

// ld compares two strings and returns the levenshtein distance between them.
func ld(s, t string, ignoreCase bool) int ***REMOVED***
	if ignoreCase ***REMOVED***
		s = strings.ToLower(s)
		t = strings.ToLower(t)
	***REMOVED***
	d := make([][]int, len(s)+1)
	for i := range d ***REMOVED***
		d[i] = make([]int, len(t)+1)
	***REMOVED***
	for i := range d ***REMOVED***
		d[i][0] = i
	***REMOVED***
	for j := range d[0] ***REMOVED***
		d[0][j] = j
	***REMOVED***
	for j := 1; j <= len(t); j++ ***REMOVED***
		for i := 1; i <= len(s); i++ ***REMOVED***
			if s[i-1] == t[j-1] ***REMOVED***
				d[i][j] = d[i-1][j-1]
			***REMOVED*** else ***REMOVED***
				min := d[i-1][j]
				if d[i][j-1] < min ***REMOVED***
					min = d[i][j-1]
				***REMOVED***
				if d[i-1][j-1] < min ***REMOVED***
					min = d[i-1][j-1]
				***REMOVED***
				d[i][j] = min + 1
			***REMOVED***
		***REMOVED***

	***REMOVED***
	return d[len(s)][len(t)]
***REMOVED***

func stringInSlice(a string, list []string) bool ***REMOVED***
	for _, b := range list ***REMOVED***
		if b == a ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// CheckErr prints the msg with the prefix 'Error:' and exits with error code 1. If the msg is nil, it does nothing.
func CheckErr(msg interface***REMOVED******REMOVED***) ***REMOVED***
	if msg != nil ***REMOVED***
		fmt.Fprintln(os.Stderr, "Error:", msg)
		os.Exit(1)
	***REMOVED***
***REMOVED***

// WriteStringAndCheck writes a string into a buffer, and checks if the error is not nil.
func WriteStringAndCheck(b io.StringWriter, s string) ***REMOVED***
	_, err := b.WriteString(s)
	CheckErr(err)
***REMOVED***
