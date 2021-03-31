//
// Copyright (c) 2011-2019 Canonical Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package yaml

import (
	"encoding/base64"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type resolveMapItem struct ***REMOVED***
	value interface***REMOVED******REMOVED***
	tag   string
***REMOVED***

var resolveTable = make([]byte, 256)
var resolveMap = make(map[string]resolveMapItem)

func init() ***REMOVED***
	t := resolveTable
	t[int('+')] = 'S' // Sign
	t[int('-')] = 'S'
	for _, c := range "0123456789" ***REMOVED***
		t[int(c)] = 'D' // Digit
	***REMOVED***
	for _, c := range "yYnNtTfFoO~" ***REMOVED***
		t[int(c)] = 'M' // In map
	***REMOVED***
	t[int('.')] = '.' // Float (potentially in map)

	var resolveMapList = []struct ***REMOVED***
		v   interface***REMOVED******REMOVED***
		tag string
		l   []string
	***REMOVED******REMOVED***
		***REMOVED***true, boolTag, []string***REMOVED***"true", "True", "TRUE"***REMOVED******REMOVED***,
		***REMOVED***false, boolTag, []string***REMOVED***"false", "False", "FALSE"***REMOVED******REMOVED***,
		***REMOVED***nil, nullTag, []string***REMOVED***"", "~", "null", "Null", "NULL"***REMOVED******REMOVED***,
		***REMOVED***math.NaN(), floatTag, []string***REMOVED***".nan", ".NaN", ".NAN"***REMOVED******REMOVED***,
		***REMOVED***math.Inf(+1), floatTag, []string***REMOVED***".inf", ".Inf", ".INF"***REMOVED******REMOVED***,
		***REMOVED***math.Inf(+1), floatTag, []string***REMOVED***"+.inf", "+.Inf", "+.INF"***REMOVED******REMOVED***,
		***REMOVED***math.Inf(-1), floatTag, []string***REMOVED***"-.inf", "-.Inf", "-.INF"***REMOVED******REMOVED***,
		***REMOVED***"<<", mergeTag, []string***REMOVED***"<<"***REMOVED******REMOVED***,
	***REMOVED***

	m := resolveMap
	for _, item := range resolveMapList ***REMOVED***
		for _, s := range item.l ***REMOVED***
			m[s] = resolveMapItem***REMOVED***item.v, item.tag***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

const (
	nullTag      = "!!null"
	boolTag      = "!!bool"
	strTag       = "!!str"
	intTag       = "!!int"
	floatTag     = "!!float"
	timestampTag = "!!timestamp"
	seqTag       = "!!seq"
	mapTag       = "!!map"
	binaryTag    = "!!binary"
	mergeTag     = "!!merge"
)

var longTags = make(map[string]string)
var shortTags = make(map[string]string)

func init() ***REMOVED***
	for _, stag := range []string***REMOVED***nullTag, boolTag, strTag, intTag, floatTag, timestampTag, seqTag, mapTag, binaryTag, mergeTag***REMOVED*** ***REMOVED***
		ltag := longTag(stag)
		longTags[stag] = ltag
		shortTags[ltag] = stag
	***REMOVED***
***REMOVED***

const longTagPrefix = "tag:yaml.org,2002:"

func shortTag(tag string) string ***REMOVED***
	if strings.HasPrefix(tag, longTagPrefix) ***REMOVED***
		if stag, ok := shortTags[tag]; ok ***REMOVED***
			return stag
		***REMOVED***
		return "!!" + tag[len(longTagPrefix):]
	***REMOVED***
	return tag
***REMOVED***

func longTag(tag string) string ***REMOVED***
	if strings.HasPrefix(tag, "!!") ***REMOVED***
		if ltag, ok := longTags[tag]; ok ***REMOVED***
			return ltag
		***REMOVED***
		return longTagPrefix + tag[2:]
	***REMOVED***
	return tag
***REMOVED***

func resolvableTag(tag string) bool ***REMOVED***
	switch tag ***REMOVED***
	case "", strTag, boolTag, intTag, floatTag, nullTag, timestampTag:
		return true
	***REMOVED***
	return false
***REMOVED***

var yamlStyleFloat = regexp.MustCompile(`^[-+]?(\.[0-9]+|[0-9]+(\.[0-9]*)?)([eE][-+]?[0-9]+)?$`)

func resolve(tag string, in string) (rtag string, out interface***REMOVED******REMOVED***) ***REMOVED***
	tag = shortTag(tag)
	if !resolvableTag(tag) ***REMOVED***
		return tag, in
	***REMOVED***

	defer func() ***REMOVED***
		switch tag ***REMOVED***
		case "", rtag, strTag, binaryTag:
			return
		case floatTag:
			if rtag == intTag ***REMOVED***
				switch v := out.(type) ***REMOVED***
				case int64:
					rtag = floatTag
					out = float64(v)
					return
				case int:
					rtag = floatTag
					out = float64(v)
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***
		failf("cannot decode %s `%s` as a %s", shortTag(rtag), in, shortTag(tag))
	***REMOVED***()

	// Any data is accepted as a !!str or !!binary.
	// Otherwise, the prefix is enough of a hint about what it might be.
	hint := byte('N')
	if in != "" ***REMOVED***
		hint = resolveTable[in[0]]
	***REMOVED***
	if hint != 0 && tag != strTag && tag != binaryTag ***REMOVED***
		// Handle things we can lookup in a map.
		if item, ok := resolveMap[in]; ok ***REMOVED***
			return item.tag, item.value
		***REMOVED***

		// Base 60 floats are a bad idea, were dropped in YAML 1.2, and
		// are purposefully unsupported here. They're still quoted on
		// the way out for compatibility with other parser, though.

		switch hint ***REMOVED***
		case 'M':
			// We've already checked the map above.

		case '.':
			// Not in the map, so maybe a normal float.
			floatv, err := strconv.ParseFloat(in, 64)
			if err == nil ***REMOVED***
				return floatTag, floatv
			***REMOVED***

		case 'D', 'S':
			// Int, float, or timestamp.
			// Only try values as a timestamp if the value is unquoted or there's an explicit
			// !!timestamp tag.
			if tag == "" || tag == timestampTag ***REMOVED***
				t, ok := parseTimestamp(in)
				if ok ***REMOVED***
					return timestampTag, t
				***REMOVED***
			***REMOVED***

			plain := strings.Replace(in, "_", "", -1)
			intv, err := strconv.ParseInt(plain, 0, 64)
			if err == nil ***REMOVED***
				if intv == int64(int(intv)) ***REMOVED***
					return intTag, int(intv)
				***REMOVED*** else ***REMOVED***
					return intTag, intv
				***REMOVED***
			***REMOVED***
			uintv, err := strconv.ParseUint(plain, 0, 64)
			if err == nil ***REMOVED***
				return intTag, uintv
			***REMOVED***
			if yamlStyleFloat.MatchString(plain) ***REMOVED***
				floatv, err := strconv.ParseFloat(plain, 64)
				if err == nil ***REMOVED***
					return floatTag, floatv
				***REMOVED***
			***REMOVED***
			if strings.HasPrefix(plain, "0b") ***REMOVED***
				intv, err := strconv.ParseInt(plain[2:], 2, 64)
				if err == nil ***REMOVED***
					if intv == int64(int(intv)) ***REMOVED***
						return intTag, int(intv)
					***REMOVED*** else ***REMOVED***
						return intTag, intv
					***REMOVED***
				***REMOVED***
				uintv, err := strconv.ParseUint(plain[2:], 2, 64)
				if err == nil ***REMOVED***
					return intTag, uintv
				***REMOVED***
			***REMOVED*** else if strings.HasPrefix(plain, "-0b") ***REMOVED***
				intv, err := strconv.ParseInt("-"+plain[3:], 2, 64)
				if err == nil ***REMOVED***
					if true || intv == int64(int(intv)) ***REMOVED***
						return intTag, int(intv)
					***REMOVED*** else ***REMOVED***
						return intTag, intv
					***REMOVED***
				***REMOVED***
			***REMOVED***
			// Octals as introduced in version 1.2 of the spec.
			// Octals from the 1.1 spec, spelled as 0777, are still
			// decoded by default in v3 as well for compatibility.
			// May be dropped in v4 depending on how usage evolves.
			if strings.HasPrefix(plain, "0o") ***REMOVED***
				intv, err := strconv.ParseInt(plain[2:], 8, 64)
				if err == nil ***REMOVED***
					if intv == int64(int(intv)) ***REMOVED***
						return intTag, int(intv)
					***REMOVED*** else ***REMOVED***
						return intTag, intv
					***REMOVED***
				***REMOVED***
				uintv, err := strconv.ParseUint(plain[2:], 8, 64)
				if err == nil ***REMOVED***
					return intTag, uintv
				***REMOVED***
			***REMOVED*** else if strings.HasPrefix(plain, "-0o") ***REMOVED***
				intv, err := strconv.ParseInt("-"+plain[3:], 8, 64)
				if err == nil ***REMOVED***
					if true || intv == int64(int(intv)) ***REMOVED***
						return intTag, int(intv)
					***REMOVED*** else ***REMOVED***
						return intTag, intv
					***REMOVED***
				***REMOVED***
			***REMOVED***
		default:
			panic("internal error: missing handler for resolver table: " + string(rune(hint)) + " (with " + in + ")")
		***REMOVED***
	***REMOVED***
	return strTag, in
***REMOVED***

// encodeBase64 encodes s as base64 that is broken up into multiple lines
// as appropriate for the resulting length.
func encodeBase64(s string) string ***REMOVED***
	const lineLen = 70
	encLen := base64.StdEncoding.EncodedLen(len(s))
	lines := encLen/lineLen + 1
	buf := make([]byte, encLen*2+lines)
	in := buf[0:encLen]
	out := buf[encLen:]
	base64.StdEncoding.Encode(in, []byte(s))
	k := 0
	for i := 0; i < len(in); i += lineLen ***REMOVED***
		j := i + lineLen
		if j > len(in) ***REMOVED***
			j = len(in)
		***REMOVED***
		k += copy(out[k:], in[i:j])
		if lines > 1 ***REMOVED***
			out[k] = '\n'
			k++
		***REMOVED***
	***REMOVED***
	return string(out[:k])
***REMOVED***

// This is a subset of the formats allowed by the regular expression
// defined at http://yaml.org/type/timestamp.html.
var allowedTimestampFormats = []string***REMOVED***
	"2006-1-2T15:4:5.999999999Z07:00", // RCF3339Nano with short date fields.
	"2006-1-2t15:4:5.999999999Z07:00", // RFC3339Nano with short date fields and lower-case "t".
	"2006-1-2 15:4:5.999999999",       // space separated with no time zone
	"2006-1-2",                        // date only
	// Notable exception: time.Parse cannot handle: "2001-12-14 21:59:43.10 -5"
	// from the set of examples.
***REMOVED***

// parseTimestamp parses s as a timestamp string and
// returns the timestamp and reports whether it succeeded.
// Timestamp formats are defined at http://yaml.org/type/timestamp.html
func parseTimestamp(s string) (time.Time, bool) ***REMOVED***
	// TODO write code to check all the formats supported by
	// http://yaml.org/type/timestamp.html instead of using time.Parse.

	// Quick check: all date formats start with YYYY-.
	i := 0
	for ; i < len(s); i++ ***REMOVED***
		if c := s[i]; c < '0' || c > '9' ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	if i != 4 || i == len(s) || s[i] != '-' ***REMOVED***
		return time.Time***REMOVED******REMOVED***, false
	***REMOVED***
	for _, format := range allowedTimestampFormats ***REMOVED***
		if t, err := time.Parse(format, s); err == nil ***REMOVED***
			return t, true
		***REMOVED***
	***REMOVED***
	return time.Time***REMOVED******REMOVED***, false
***REMOVED***
