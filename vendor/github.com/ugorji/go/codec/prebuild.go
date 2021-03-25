// +build prebuild

package main

// prebuild.go generates sort implementations for
// various slice types and combination slice+reflect.Value types.
//
// The combination slice+reflect.Value types are used
// during canonical encode, and the others are used during fast-path
// encoding of map keys.

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

// genInternalSortableTypes returns the types
// that are used for fast-path canonical's encoding of maps.
//
// For now, we only support the highest sizes for
// int64, uint64, float64, bool, string, bytes.
func genInternalSortableTypes() []string ***REMOVED***
	return []string***REMOVED***
		"string",
		// "float32",
		"float64",
		// "uint",
		// "uint8",
		// "uint16",
		// "uint32",
		"uint64",
		"uintptr",
		// "int",
		// "int8",
		// "int16",
		// "int32",
		"int64",
		"bool",
		"time",
		"bytes",
	***REMOVED***
***REMOVED***

// genInternalSortablePlusTypes returns the types
// that are used for reflection-based canonical's encoding of maps.
//
// For now, we only support the highest sizes for
// int64, uint64, float64, bool, string, bytes.
func genInternalSortablePlusTypes() []string ***REMOVED***
	return []string***REMOVED***
		"string",
		"float64",
		"uint64",
		"uintptr",
		"int64",
		"bool",
		"time",
		"bytes",
	***REMOVED***
***REMOVED***

func genTypeForShortName(s string) string ***REMOVED***
	switch s ***REMOVED***
	case "time":
		return "time.Time"
	case "bytes":
		return "[]byte"
	***REMOVED***
	return s
***REMOVED***

func genArgs(args ...interface***REMOVED******REMOVED***) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	m := make(map[string]interface***REMOVED******REMOVED***, len(args)/2)
	for i := 0; i < len(args); ***REMOVED***
		m[args[i].(string)] = args[i+1]
		i += 2
	***REMOVED***
	return m
***REMOVED***

func genEndsWith(s0 string, sn ...string) bool ***REMOVED***
	for _, s := range sn ***REMOVED***
		if strings.HasSuffix(s0, s) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func chkerr(err error) ***REMOVED***
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func run(fnameIn, fnameOut string) ***REMOVED***
	var err error

	funcs := make(template.FuncMap)
	funcs["sortables"] = genInternalSortableTypes
	funcs["sortablesplus"] = genInternalSortablePlusTypes
	funcs["tshort"] = genTypeForShortName
	funcs["endswith"] = genEndsWith
	funcs["args"] = genArgs

	t := template.New("").Funcs(funcs)
	fin, err := os.Open(fnameIn)
	chkerr(err)
	defer fin.Close()
	fout, err := os.Create(fnameOut)
	chkerr(err)
	defer fout.Close()
	tmplstr, err := ioutil.ReadAll(fin)
	chkerr(err)
	t, err = t.Parse(string(tmplstr))
	chkerr(err)
	var out bytes.Buffer
	err = t.Execute(&out, 0)
	chkerr(err)
	bout, err := format.Source(out.Bytes())
	if err != nil ***REMOVED***
		fout.Write(out.Bytes()) // write out if error, so we can still see.
	***REMOVED***
	chkerr(err)
	// write out if error, as much as possible, so we can still see.
	_, err = fout.Write(bout)
	chkerr(err)
***REMOVED***

func main() ***REMOVED***
	run("sort-slice.go.tmpl", "sort-slice.generated.go")
***REMOVED***
