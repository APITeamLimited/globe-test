// Package escape contains utilities for escaping parts of InfluxQL
// and InfluxDB line protocol.
package escape // import "github.com/influxdata/influxdb1-client/pkg/escape"

import (
	"bytes"
	"strings"
)

// Codes is a map of bytes to be escaped.
var Codes = map[byte][]byte***REMOVED***
	',': []byte(`\,`),
	'"': []byte(`\"`),
	' ': []byte(`\ `),
	'=': []byte(`\=`),
***REMOVED***

// Bytes escapes characters on the input slice, as defined by Codes.
func Bytes(in []byte) []byte ***REMOVED***
	for b, esc := range Codes ***REMOVED***
		in = bytes.Replace(in, []byte***REMOVED***b***REMOVED***, esc, -1)
	***REMOVED***
	return in
***REMOVED***

const escapeChars = `," =`

// IsEscaped returns whether b has any escaped characters,
// i.e. whether b seems to have been processed by Bytes.
func IsEscaped(b []byte) bool ***REMOVED***
	for len(b) > 0 ***REMOVED***
		i := bytes.IndexByte(b, '\\')
		if i < 0 ***REMOVED***
			return false
		***REMOVED***

		if i+1 < len(b) && strings.IndexByte(escapeChars, b[i+1]) >= 0 ***REMOVED***
			return true
		***REMOVED***
		b = b[i+1:]
	***REMOVED***
	return false
***REMOVED***

// AppendUnescaped appends the unescaped version of src to dst
// and returns the resulting slice.
func AppendUnescaped(dst, src []byte) []byte ***REMOVED***
	var pos int
	for len(src) > 0 ***REMOVED***
		next := bytes.IndexByte(src[pos:], '\\')
		if next < 0 || pos+next+1 >= len(src) ***REMOVED***
			return append(dst, src...)
		***REMOVED***

		if pos+next+1 < len(src) && strings.IndexByte(escapeChars, src[pos+next+1]) >= 0 ***REMOVED***
			if pos+next > 0 ***REMOVED***
				dst = append(dst, src[:pos+next]...)
			***REMOVED***
			src = src[pos+next+1:]
			pos = 0
		***REMOVED*** else ***REMOVED***
			pos += next + 1
		***REMOVED***
	***REMOVED***

	return dst
***REMOVED***

// Unescape returns a new slice containing the unescaped version of in.
func Unescape(in []byte) []byte ***REMOVED***
	if len(in) == 0 ***REMOVED***
		return nil
	***REMOVED***

	if bytes.IndexByte(in, '\\') == -1 ***REMOVED***
		return in
	***REMOVED***

	i := 0
	inLen := len(in)

	// The output size will be no more than inLen. Preallocating the
	// capacity of the output is faster and uses less memory than
	// letting append() do its own (over)allocation.
	out := make([]byte, 0, inLen)

	for ***REMOVED***
		if i >= inLen ***REMOVED***
			break
		***REMOVED***
		if in[i] == '\\' && i+1 < inLen ***REMOVED***
			switch in[i+1] ***REMOVED***
			case ',':
				out = append(out, ',')
				i += 2
				continue
			case '"':
				out = append(out, '"')
				i += 2
				continue
			case ' ':
				out = append(out, ' ')
				i += 2
				continue
			case '=':
				out = append(out, '=')
				i += 2
				continue
			***REMOVED***
		***REMOVED***
		out = append(out, in[i])
		i += 1
	***REMOVED***
	return out
***REMOVED***
