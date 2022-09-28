// Package easyjson contains marshaler/unmarshaler interfaces and helper functions.
package easyjson

import (
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"unsafe"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

// Marshaler is an easyjson-compatible marshaler interface.
type Marshaler interface ***REMOVED***
	MarshalEasyJSON(w *jwriter.Writer)
***REMOVED***

// Marshaler is an easyjson-compatible unmarshaler interface.
type Unmarshaler interface ***REMOVED***
	UnmarshalEasyJSON(w *jlexer.Lexer)
***REMOVED***

// MarshalerUnmarshaler is an easyjson-compatible marshaler/unmarshaler interface.
type MarshalerUnmarshaler interface ***REMOVED***
	Marshaler
	Unmarshaler
***REMOVED***

// Optional defines an undefined-test method for a type to integrate with 'omitempty' logic.
type Optional interface ***REMOVED***
	IsDefined() bool
***REMOVED***

// UnknownsUnmarshaler provides a method to unmarshal unknown struct fileds and save them as you want
type UnknownsUnmarshaler interface ***REMOVED***
	UnmarshalUnknown(in *jlexer.Lexer, key string)
***REMOVED***

// UnknownsMarshaler provides a method to write additional struct fields
type UnknownsMarshaler interface ***REMOVED***
	MarshalUnknowns(w *jwriter.Writer, first bool)
***REMOVED***

func isNilInterface(i interface***REMOVED******REMOVED***) bool ***REMOVED***
	return (*[2]uintptr)(unsafe.Pointer(&i))[1] == 0
***REMOVED***

// Marshal returns data as a single byte slice. Method is suboptimal as the data is likely to be copied
// from a chain of smaller chunks.
func Marshal(v Marshaler) ([]byte, error) ***REMOVED***
	if isNilInterface(v) ***REMOVED***
		return nullBytes, nil
	***REMOVED***

	w := jwriter.Writer***REMOVED******REMOVED***
	v.MarshalEasyJSON(&w)
	return w.BuildBytes()
***REMOVED***

// MarshalToWriter marshals the data to an io.Writer.
func MarshalToWriter(v Marshaler, w io.Writer) (written int, err error) ***REMOVED***
	if isNilInterface(v) ***REMOVED***
		return w.Write(nullBytes)
	***REMOVED***

	jw := jwriter.Writer***REMOVED******REMOVED***
	v.MarshalEasyJSON(&jw)
	return jw.DumpTo(w)
***REMOVED***

// MarshalToHTTPResponseWriter sets Content-Length and Content-Type headers for the
// http.ResponseWriter, and send the data to the writer. started will be equal to
// false if an error occurred before any http.ResponseWriter methods were actually
// invoked (in this case a 500 reply is possible).
func MarshalToHTTPResponseWriter(v Marshaler, w http.ResponseWriter) (started bool, written int, err error) ***REMOVED***
	if isNilInterface(v) ***REMOVED***
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(nullBytes)))
		written, err = w.Write(nullBytes)
		return true, written, err
	***REMOVED***

	jw := jwriter.Writer***REMOVED******REMOVED***
	v.MarshalEasyJSON(&jw)
	if jw.Error != nil ***REMOVED***
		return false, 0, jw.Error
	***REMOVED***
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(jw.Size()))

	started = true
	written, err = jw.DumpTo(w)
	return
***REMOVED***

// Unmarshal decodes the JSON in data into the object.
func Unmarshal(data []byte, v Unmarshaler) error ***REMOVED***
	l := jlexer.Lexer***REMOVED***Data: data***REMOVED***
	v.UnmarshalEasyJSON(&l)
	return l.Error()
***REMOVED***

// UnmarshalFromReader reads all the data in the reader and decodes as JSON into the object.
func UnmarshalFromReader(r io.Reader, v Unmarshaler) error ***REMOVED***
	data, err := ioutil.ReadAll(r)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	l := jlexer.Lexer***REMOVED***Data: data***REMOVED***
	v.UnmarshalEasyJSON(&l)
	return l.Error()
***REMOVED***
