// +build go1.6

package humanize

import (
	"bytes"
	"math/big"
	"strings"
)

// BigCommaf produces a string form of the given big.Float in base 10
// with commas after every three orders of magnitude.
func BigCommaf(v *big.Float) string ***REMOVED***
	buf := &bytes.Buffer***REMOVED******REMOVED***
	if v.Sign() < 0 ***REMOVED***
		buf.Write([]byte***REMOVED***'-'***REMOVED***)
		v.Abs(v)
	***REMOVED***

	comma := []byte***REMOVED***','***REMOVED***

	parts := strings.Split(v.Text('f', -1), ".")
	pos := 0
	if len(parts[0])%3 != 0 ***REMOVED***
		pos += len(parts[0]) % 3
		buf.WriteString(parts[0][:pos])
		buf.Write(comma)
	***REMOVED***
	for ; pos < len(parts[0]); pos += 3 ***REMOVED***
		buf.WriteString(parts[0][pos : pos+3])
		buf.Write(comma)
	***REMOVED***
	buf.Truncate(buf.Len() - 1)

	if len(parts) > 1 ***REMOVED***
		buf.Write([]byte***REMOVED***'.'***REMOVED***)
		buf.WriteString(parts[1])
	***REMOVED***
	return buf.String()
***REMOVED***
