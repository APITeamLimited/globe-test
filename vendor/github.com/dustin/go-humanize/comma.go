package humanize

import (
	"bytes"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// Comma produces a string form of the given number in base 10 with
// commas after every three orders of magnitude.
//
// e.g. Comma(834142) -> 834,142
func Comma(v int64) string ***REMOVED***
	sign := ""

	// Min int64 can't be negated to a usable value, so it has to be special cased.
	if v == math.MinInt64 ***REMOVED***
		return "-9,223,372,036,854,775,808"
	***REMOVED***

	if v < 0 ***REMOVED***
		sign = "-"
		v = 0 - v
	***REMOVED***

	parts := []string***REMOVED***"", "", "", "", "", "", ""***REMOVED***
	j := len(parts) - 1

	for v > 999 ***REMOVED***
		parts[j] = strconv.FormatInt(v%1000, 10)
		switch len(parts[j]) ***REMOVED***
		case 2:
			parts[j] = "0" + parts[j]
		case 1:
			parts[j] = "00" + parts[j]
		***REMOVED***
		v = v / 1000
		j--
	***REMOVED***
	parts[j] = strconv.Itoa(int(v))
	return sign + strings.Join(parts[j:], ",")
***REMOVED***

// Commaf produces a string form of the given number in base 10 with
// commas after every three orders of magnitude.
//
// e.g. Commaf(834142.32) -> 834,142.32
func Commaf(v float64) string ***REMOVED***
	buf := &bytes.Buffer***REMOVED******REMOVED***
	if v < 0 ***REMOVED***
		buf.Write([]byte***REMOVED***'-'***REMOVED***)
		v = 0 - v
	***REMOVED***

	comma := []byte***REMOVED***','***REMOVED***

	parts := strings.Split(strconv.FormatFloat(v, 'f', -1, 64), ".")
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

// BigComma produces a string form of the given big.Int in base 10
// with commas after every three orders of magnitude.
func BigComma(b *big.Int) string ***REMOVED***
	sign := ""
	if b.Sign() < 0 ***REMOVED***
		sign = "-"
		b.Abs(b)
	***REMOVED***

	athousand := big.NewInt(1000)
	c := (&big.Int***REMOVED******REMOVED***).Set(b)
	_, m := oom(c, athousand)
	parts := make([]string, m+1)
	j := len(parts) - 1

	mod := &big.Int***REMOVED******REMOVED***
	for b.Cmp(athousand) >= 0 ***REMOVED***
		b.DivMod(b, athousand, mod)
		parts[j] = strconv.FormatInt(mod.Int64(), 10)
		switch len(parts[j]) ***REMOVED***
		case 2:
			parts[j] = "0" + parts[j]
		case 1:
			parts[j] = "00" + parts[j]
		***REMOVED***
		j--
	***REMOVED***
	parts[j] = strconv.Itoa(int(b.Int64()))
	return sign + strings.Join(parts[j:], ",")
***REMOVED***
