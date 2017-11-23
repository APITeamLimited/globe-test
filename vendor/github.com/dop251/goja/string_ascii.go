package goja

import (
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
	"strings"
)

type asciiString string

type asciiRuneReader struct ***REMOVED***
	s   asciiString
	pos int
***REMOVED***

func (rr *asciiRuneReader) ReadRune() (r rune, size int, err error) ***REMOVED***
	if rr.pos < len(rr.s) ***REMOVED***
		r = rune(rr.s[rr.pos])
		size = 1
		rr.pos++
	***REMOVED*** else ***REMOVED***
		err = io.EOF
	***REMOVED***
	return
***REMOVED***

func (s asciiString) reader(start int) io.RuneReader ***REMOVED***
	return &asciiRuneReader***REMOVED***
		s: s[start:],
	***REMOVED***
***REMOVED***

// ss must be trimmed
func strToInt(ss string) (int64, error) ***REMOVED***
	if ss == "" ***REMOVED***
		return 0, nil
	***REMOVED***
	if ss == "-0" ***REMOVED***
		return 0, strconv.ErrSyntax
	***REMOVED***
	if len(ss) > 2 ***REMOVED***
		switch ss[:2] ***REMOVED***
		case "0x", "0X":
			i, _ := strconv.ParseInt(ss[2:], 16, 64)
			return i, nil
		case "0b", "0B":
			i, _ := strconv.ParseInt(ss[2:], 2, 64)
			return i, nil
		case "0o", "0O":
			i, _ := strconv.ParseInt(ss[2:], 8, 64)
			return i, nil
		***REMOVED***
	***REMOVED***
	return strconv.ParseInt(ss, 10, 64)
***REMOVED***

func (s asciiString) _toInt() (int64, error) ***REMOVED***
	return strToInt(strings.TrimSpace(string(s)))
***REMOVED***

func isRangeErr(err error) bool ***REMOVED***
	if err, ok := err.(*strconv.NumError); ok ***REMOVED***
		return err.Err == strconv.ErrRange
	***REMOVED***
	return false
***REMOVED***

func (s asciiString) _toFloat() (float64, error) ***REMOVED***
	ss := strings.TrimSpace(string(s))
	if ss == "" ***REMOVED***
		return 0, nil
	***REMOVED***
	if ss == "-0" ***REMOVED***
		var f float64
		return -f, nil
	***REMOVED***
	f, err := strconv.ParseFloat(ss, 64)
	if isRangeErr(err) ***REMOVED***
		err = nil
	***REMOVED***
	return f, err
***REMOVED***

func (s asciiString) ToInteger() int64 ***REMOVED***
	if s == "" ***REMOVED***
		return 0
	***REMOVED***
	if s == "Infinity" || s == "+Infinity" ***REMOVED***
		return math.MaxInt64
	***REMOVED***
	if s == "-Infinity" ***REMOVED***
		return math.MinInt64
	***REMOVED***
	i, err := s._toInt()
	if err != nil ***REMOVED***
		f, err := s._toFloat()
		if err == nil ***REMOVED***
			return int64(f)
		***REMOVED***
	***REMOVED***
	return i
***REMOVED***

func (s asciiString) ToString() valueString ***REMOVED***
	return s
***REMOVED***

func (s asciiString) String() string ***REMOVED***
	return string(s)
***REMOVED***

func (s asciiString) ToFloat() float64 ***REMOVED***
	if s == "" ***REMOVED***
		return 0
	***REMOVED***
	if s == "Infinity" || s == "+Infinity" ***REMOVED***
		return math.Inf(1)
	***REMOVED***
	if s == "-Infinity" ***REMOVED***
		return math.Inf(-1)
	***REMOVED***
	f, err := s._toFloat()
	if err != nil ***REMOVED***
		i, err := s._toInt()
		if err == nil ***REMOVED***
			return float64(i)
		***REMOVED***
		f = math.NaN()
	***REMOVED***
	return f
***REMOVED***

func (s asciiString) ToBoolean() bool ***REMOVED***
	return s != ""
***REMOVED***

func (s asciiString) ToNumber() Value ***REMOVED***
	if s == "" ***REMOVED***
		return intToValue(0)
	***REMOVED***
	if s == "Infinity" || s == "+Infinity" ***REMOVED***
		return _positiveInf
	***REMOVED***
	if s == "-Infinity" ***REMOVED***
		return _negativeInf
	***REMOVED***

	if i, err := s._toInt(); err == nil ***REMOVED***
		return intToValue(i)
	***REMOVED***

	if f, err := s._toFloat(); err == nil ***REMOVED***
		return floatToValue(f)
	***REMOVED***

	return _NaN
***REMOVED***

func (s asciiString) ToObject(r *Runtime) *Object ***REMOVED***
	return r._newString(s)
***REMOVED***

func (s asciiString) SameAs(other Value) bool ***REMOVED***
	if otherStr, ok := other.(asciiString); ok ***REMOVED***
		return s == otherStr
	***REMOVED***
	return false
***REMOVED***

func (s asciiString) Equals(other Value) bool ***REMOVED***
	if o, ok := other.(asciiString); ok ***REMOVED***
		return s == o
	***REMOVED***

	if o, ok := other.assertInt(); ok ***REMOVED***
		if o1, e := s._toInt(); e == nil ***REMOVED***
			return o1 == o
		***REMOVED***
		return false
	***REMOVED***

	if o, ok := other.assertFloat(); ok ***REMOVED***
		return s.ToFloat() == o
	***REMOVED***

	if o, ok := other.(valueBool); ok ***REMOVED***
		if o1, e := s._toFloat(); e == nil ***REMOVED***
			return o1 == o.ToFloat()
		***REMOVED***
		return false
	***REMOVED***

	if o, ok := other.(*Object); ok ***REMOVED***
		return s.Equals(o.self.toPrimitive())
	***REMOVED***
	return false
***REMOVED***

func (s asciiString) StrictEquals(other Value) bool ***REMOVED***
	if otherStr, ok := other.(asciiString); ok ***REMOVED***
		return s == otherStr
	***REMOVED***
	return false
***REMOVED***

func (s asciiString) assertInt() (int64, bool) ***REMOVED***
	return 0, false
***REMOVED***

func (s asciiString) assertFloat() (float64, bool) ***REMOVED***
	return 0, false
***REMOVED***

func (s asciiString) assertString() (valueString, bool) ***REMOVED***
	return s, true
***REMOVED***

func (s asciiString) baseObject(r *Runtime) *Object ***REMOVED***
	ss := r.stringSingleton
	ss.value = s
	ss.setLength()
	return ss.val
***REMOVED***

func (s asciiString) charAt(idx int64) rune ***REMOVED***
	return rune(s[idx])
***REMOVED***

func (s asciiString) length() int64 ***REMOVED***
	return int64(len(s))
***REMOVED***

func (s asciiString) concat(other valueString) valueString ***REMOVED***
	switch other := other.(type) ***REMOVED***
	case asciiString:
		b := make([]byte, len(s)+len(other))
		copy(b, s)
		copy(b[len(s):], other)
		return asciiString(b)
		//return asciiString(string(s) + string(other))
	case unicodeString:
		b := make([]uint16, len(s)+len(other))
		for i := 0; i < len(s); i++ ***REMOVED***
			b[i] = uint16(s[i])
		***REMOVED***
		copy(b[len(s):], other)
		return unicodeString(b)
	default:
		panic(fmt.Errorf("Unknown string type: %T", other))
	***REMOVED***
***REMOVED***

func (s asciiString) substring(start, end int64) valueString ***REMOVED***
	return asciiString(s[start:end])
***REMOVED***

func (s asciiString) compareTo(other valueString) int ***REMOVED***
	switch other := other.(type) ***REMOVED***
	case asciiString:
		return strings.Compare(string(s), string(other))
	case unicodeString:
		return strings.Compare(string(s), other.String())
	default:
		panic(fmt.Errorf("Unknown string type: %T", other))
	***REMOVED***
***REMOVED***

func (s asciiString) index(substr valueString, start int64) int64 ***REMOVED***
	if substr, ok := substr.(asciiString); ok ***REMOVED***
		p := int64(strings.Index(string(s[start:]), string(substr)))
		if p >= 0 ***REMOVED***
			return p + start
		***REMOVED***
	***REMOVED***
	return -1
***REMOVED***

func (s asciiString) lastIndex(substr valueString, pos int64) int64 ***REMOVED***
	if substr, ok := substr.(asciiString); ok ***REMOVED***
		end := pos + int64(len(substr))
		var ss string
		if end > int64(len(s)) ***REMOVED***
			ss = string(s)
		***REMOVED*** else ***REMOVED***
			ss = string(s[:end])
		***REMOVED***
		return int64(strings.LastIndex(ss, string(substr)))
	***REMOVED***
	return -1
***REMOVED***

func (s asciiString) toLower() valueString ***REMOVED***
	return asciiString(strings.ToLower(string(s)))
***REMOVED***

func (s asciiString) toUpper() valueString ***REMOVED***
	return asciiString(strings.ToUpper(string(s)))
***REMOVED***

func (s asciiString) toTrimmedUTF8() string ***REMOVED***
	return strings.TrimSpace(string(s))
***REMOVED***

func (s asciiString) Export() interface***REMOVED******REMOVED*** ***REMOVED***
	return string(s)
***REMOVED***

func (s asciiString) ExportType() reflect.Type ***REMOVED***
	return reflectTypeString
***REMOVED***
