// Copyright (c) 2016 Kelsey Hightower and others. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package envconfig

import (
	"encoding"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/tabwriter"
	"text/template"
)

const (
	// DefaultListFormat constant to use to display usage in a list format
	DefaultListFormat = `This application is configured via the environment. The following environment
variables can be used:
***REMOVED******REMOVED***range .***REMOVED******REMOVED***
***REMOVED******REMOVED***usage_key .***REMOVED******REMOVED***
  [description] ***REMOVED******REMOVED***usage_description .***REMOVED******REMOVED***
  [type]        ***REMOVED******REMOVED***usage_type .***REMOVED******REMOVED***
  [default]     ***REMOVED******REMOVED***usage_default .***REMOVED******REMOVED***
  [required]    ***REMOVED******REMOVED***usage_required .***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***
`
	// DefaultTableFormat constant to use to display usage in a tabular format
	DefaultTableFormat = `This application is configured via the environment. The following environment
variables can be used:

KEY	TYPE	DEFAULT	REQUIRED	DESCRIPTION
***REMOVED******REMOVED***range .***REMOVED******REMOVED******REMOVED******REMOVED***usage_key .***REMOVED******REMOVED***	***REMOVED******REMOVED***usage_type .***REMOVED******REMOVED***	***REMOVED******REMOVED***usage_default .***REMOVED******REMOVED***	***REMOVED******REMOVED***usage_required .***REMOVED******REMOVED***	***REMOVED******REMOVED***usage_description .***REMOVED******REMOVED***
***REMOVED******REMOVED***end***REMOVED******REMOVED***`
)

var (
	decoderType           = reflect.TypeOf((*Decoder)(nil)).Elem()
	setterType            = reflect.TypeOf((*Setter)(nil)).Elem()
	textUnmarshalerType   = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
	binaryUnmarshalerType = reflect.TypeOf((*encoding.BinaryUnmarshaler)(nil)).Elem()
)

func implementsInterface(t reflect.Type) bool ***REMOVED***
	return t.Implements(decoderType) ||
		reflect.PtrTo(t).Implements(decoderType) ||
		t.Implements(setterType) ||
		reflect.PtrTo(t).Implements(setterType) ||
		t.Implements(textUnmarshalerType) ||
		reflect.PtrTo(t).Implements(textUnmarshalerType) ||
		t.Implements(binaryUnmarshalerType) ||
		reflect.PtrTo(t).Implements(binaryUnmarshalerType)
***REMOVED***

// toTypeDescription converts Go types into a human readable description
func toTypeDescription(t reflect.Type) string ***REMOVED***
	switch t.Kind() ***REMOVED***
	case reflect.Array, reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 ***REMOVED***
			return "String"
		***REMOVED***
		return fmt.Sprintf("Comma-separated list of %s", toTypeDescription(t.Elem()))
	case reflect.Map:
		return fmt.Sprintf(
			"Comma-separated list of %s:%s pairs",
			toTypeDescription(t.Key()),
			toTypeDescription(t.Elem()),
		)
	case reflect.Ptr:
		return toTypeDescription(t.Elem())
	case reflect.Struct:
		if implementsInterface(t) && t.Name() != "" ***REMOVED***
			return t.Name()
		***REMOVED***
		return ""
	case reflect.String:
		name := t.Name()
		if name != "" && name != "string" ***REMOVED***
			return name
		***REMOVED***
		return "String"
	case reflect.Bool:
		name := t.Name()
		if name != "" && name != "bool" ***REMOVED***
			return name
		***REMOVED***
		return "True or False"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		name := t.Name()
		if name != "" && !strings.HasPrefix(name, "int") ***REMOVED***
			return name
		***REMOVED***
		return "Integer"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		name := t.Name()
		if name != "" && !strings.HasPrefix(name, "uint") ***REMOVED***
			return name
		***REMOVED***
		return "Unsigned Integer"
	case reflect.Float32, reflect.Float64:
		name := t.Name()
		if name != "" && !strings.HasPrefix(name, "float") ***REMOVED***
			return name
		***REMOVED***
		return "Float"
	***REMOVED***
	return fmt.Sprintf("%+v", t)
***REMOVED***

// Usage writes usage information to stdout using the default header and table format
func Usage(prefix string, spec interface***REMOVED******REMOVED***) error ***REMOVED***
	// The default is to output the usage information as a table
	// Create tabwriter instance to support table output
	tabs := tabwriter.NewWriter(os.Stdout, 1, 0, 4, ' ', 0)

	err := Usagef(prefix, spec, tabs, DefaultTableFormat)
	tabs.Flush()
	return err
***REMOVED***

// Usagef writes usage information to the specified io.Writer using the specifed template specification
func Usagef(prefix string, spec interface***REMOVED******REMOVED***, out io.Writer, format string) error ***REMOVED***

	// Specify the default usage template functions
	functions := template.FuncMap***REMOVED***
		"usage_key":         func(v varInfo) string ***REMOVED*** return v.Key ***REMOVED***,
		"usage_description": func(v varInfo) string ***REMOVED*** return v.Tags.Get("desc") ***REMOVED***,
		"usage_type":        func(v varInfo) string ***REMOVED*** return toTypeDescription(v.Field.Type()) ***REMOVED***,
		"usage_default":     func(v varInfo) string ***REMOVED*** return v.Tags.Get("default") ***REMOVED***,
		"usage_required": func(v varInfo) (string, error) ***REMOVED***
			req := v.Tags.Get("required")
			if req != "" ***REMOVED***
				reqB, err := strconv.ParseBool(req)
				if err != nil ***REMOVED***
					return "", err
				***REMOVED***
				if reqB ***REMOVED***
					req = "true"
				***REMOVED***
			***REMOVED***
			return req, nil
		***REMOVED***,
	***REMOVED***

	tmpl, err := template.New("envconfig").Funcs(functions).Parse(format)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return Usaget(prefix, spec, out, tmpl)
***REMOVED***

// Usaget writes usage information to the specified io.Writer using the specified template
func Usaget(prefix string, spec interface***REMOVED******REMOVED***, out io.Writer, tmpl *template.Template) error ***REMOVED***
	// gather first
	infos, err := gatherInfo(prefix, spec)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return tmpl.Execute(out, infos)
***REMOVED***
