// Package gjson provides searching for json strings.
package gjson

import (
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"

	"github.com/tidwall/match"
	"github.com/tidwall/pretty"
)

// Type is Result type
type Type int

const (
	// Null is a null json value
	Null Type = iota
	// False is a json false boolean
	False
	// Number is json number
	Number
	// String is a json string
	String
	// True is a json true boolean
	True
	// JSON is a raw block of JSON
	JSON
)

// String returns a string representation of the type.
func (t Type) String() string ***REMOVED***
	switch t ***REMOVED***
	default:
		return ""
	case Null:
		return "Null"
	case False:
		return "False"
	case Number:
		return "Number"
	case String:
		return "String"
	case True:
		return "True"
	case JSON:
		return "JSON"
	***REMOVED***
***REMOVED***

// Result represents a json value that is returned from Get().
type Result struct ***REMOVED***
	// Type is the json type
	Type Type
	// Raw is the raw json
	Raw string
	// Str is the json string
	Str string
	// Num is the json number
	Num float64
	// Index of raw value in original json, zero means index unknown
	Index int
	// Indexes of all the elements that match on a path containing the '#'
	// query character.
	Indexes []int
***REMOVED***

// String returns a string representation of the value.
func (t Result) String() string ***REMOVED***
	switch t.Type ***REMOVED***
	default:
		return ""
	case False:
		return "false"
	case Number:
		if len(t.Raw) == 0 ***REMOVED***
			// calculated result
			return strconv.FormatFloat(t.Num, 'f', -1, 64)
		***REMOVED***
		var i int
		if t.Raw[0] == '-' ***REMOVED***
			i++
		***REMOVED***
		for ; i < len(t.Raw); i++ ***REMOVED***
			if t.Raw[i] < '0' || t.Raw[i] > '9' ***REMOVED***
				return strconv.FormatFloat(t.Num, 'f', -1, 64)
			***REMOVED***
		***REMOVED***
		return t.Raw
	case String:
		return t.Str
	case JSON:
		return t.Raw
	case True:
		return "true"
	***REMOVED***
***REMOVED***

// Bool returns an boolean representation.
func (t Result) Bool() bool ***REMOVED***
	switch t.Type ***REMOVED***
	default:
		return false
	case True:
		return true
	case String:
		b, _ := strconv.ParseBool(strings.ToLower(t.Str))
		return b
	case Number:
		return t.Num != 0
	***REMOVED***
***REMOVED***

// Int returns an integer representation.
func (t Result) Int() int64 ***REMOVED***
	switch t.Type ***REMOVED***
	default:
		return 0
	case True:
		return 1
	case String:
		n, _ := parseInt(t.Str)
		return n
	case Number:
		// try to directly convert the float64 to int64
		i, ok := safeInt(t.Num)
		if ok ***REMOVED***
			return i
		***REMOVED***
		// now try to parse the raw string
		i, ok = parseInt(t.Raw)
		if ok ***REMOVED***
			return i
		***REMOVED***
		// fallback to a standard conversion
		return int64(t.Num)
	***REMOVED***
***REMOVED***

// Uint returns an unsigned integer representation.
func (t Result) Uint() uint64 ***REMOVED***
	switch t.Type ***REMOVED***
	default:
		return 0
	case True:
		return 1
	case String:
		n, _ := parseUint(t.Str)
		return n
	case Number:
		// try to directly convert the float64 to uint64
		i, ok := safeInt(t.Num)
		if ok && i >= 0 ***REMOVED***
			return uint64(i)
		***REMOVED***
		// now try to parse the raw string
		u, ok := parseUint(t.Raw)
		if ok ***REMOVED***
			return u
		***REMOVED***
		// fallback to a standard conversion
		return uint64(t.Num)
	***REMOVED***
***REMOVED***

// Float returns an float64 representation.
func (t Result) Float() float64 ***REMOVED***
	switch t.Type ***REMOVED***
	default:
		return 0
	case True:
		return 1
	case String:
		n, _ := strconv.ParseFloat(t.Str, 64)
		return n
	case Number:
		return t.Num
	***REMOVED***
***REMOVED***

// Time returns a time.Time representation.
func (t Result) Time() time.Time ***REMOVED***
	res, _ := time.Parse(time.RFC3339, t.String())
	return res
***REMOVED***

// Array returns back an array of values.
// If the result represents a null value or is non-existent, then an empty
// array will be returned.
// If the result is not a JSON array, the return value will be an
// array containing one result.
func (t Result) Array() []Result ***REMOVED***
	if t.Type == Null ***REMOVED***
		return []Result***REMOVED******REMOVED***
	***REMOVED***
	if !t.IsArray() ***REMOVED***
		return []Result***REMOVED***t***REMOVED***
	***REMOVED***
	r := t.arrayOrMap('[', false)
	return r.a
***REMOVED***

// IsObject returns true if the result value is a JSON object.
func (t Result) IsObject() bool ***REMOVED***
	return t.Type == JSON && len(t.Raw) > 0 && t.Raw[0] == '***REMOVED***'
***REMOVED***

// IsArray returns true if the result value is a JSON array.
func (t Result) IsArray() bool ***REMOVED***
	return t.Type == JSON && len(t.Raw) > 0 && t.Raw[0] == '['
***REMOVED***

// IsBool returns true if the result value is a JSON boolean.
func (t Result) IsBool() bool ***REMOVED***
	return t.Type == True || t.Type == False
***REMOVED***

// ForEach iterates through values.
// If the result represents a non-existent value, then no values will be
// iterated. If the result is an Object, the iterator will pass the key and
// value of each item. If the result is an Array, the iterator will only pass
// the value of each item. If the result is not a JSON array or object, the
// iterator will pass back one value equal to the result.
func (t Result) ForEach(iterator func(key, value Result) bool) ***REMOVED***
	if !t.Exists() ***REMOVED***
		return
	***REMOVED***
	if t.Type != JSON ***REMOVED***
		iterator(Result***REMOVED******REMOVED***, t)
		return
	***REMOVED***
	json := t.Raw
	var obj bool
	var i int
	var key, value Result
	for ; i < len(json); i++ ***REMOVED***
		if json[i] == '***REMOVED***' ***REMOVED***
			i++
			key.Type = String
			obj = true
			break
		***REMOVED*** else if json[i] == '[' ***REMOVED***
			i++
			key.Type = Number
			key.Num = -1
			break
		***REMOVED***
		if json[i] > ' ' ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	var str string
	var vesc bool
	var ok bool
	var idx int
	for ; i < len(json); i++ ***REMOVED***
		if obj ***REMOVED***
			if json[i] != '"' ***REMOVED***
				continue
			***REMOVED***
			s := i
			i, str, vesc, ok = parseString(json, i+1)
			if !ok ***REMOVED***
				return
			***REMOVED***
			if vesc ***REMOVED***
				key.Str = unescape(str[1 : len(str)-1])
			***REMOVED*** else ***REMOVED***
				key.Str = str[1 : len(str)-1]
			***REMOVED***
			key.Raw = str
			key.Index = s + t.Index
		***REMOVED*** else ***REMOVED***
			key.Num += 1
		***REMOVED***
		for ; i < len(json); i++ ***REMOVED***
			if json[i] <= ' ' || json[i] == ',' || json[i] == ':' ***REMOVED***
				continue
			***REMOVED***
			break
		***REMOVED***
		s := i
		i, value, ok = parseAny(json, i, true)
		if !ok ***REMOVED***
			return
		***REMOVED***
		if t.Indexes != nil ***REMOVED***
			if idx < len(t.Indexes) ***REMOVED***
				value.Index = t.Indexes[idx]
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			value.Index = s + t.Index
		***REMOVED***
		if !iterator(key, value) ***REMOVED***
			return
		***REMOVED***
		idx++
	***REMOVED***
***REMOVED***

// Map returns back a map of values. The result should be a JSON object.
// If the result is not a JSON object, the return value will be an empty map.
func (t Result) Map() map[string]Result ***REMOVED***
	if t.Type != JSON ***REMOVED***
		return map[string]Result***REMOVED******REMOVED***
	***REMOVED***
	r := t.arrayOrMap('***REMOVED***', false)
	return r.o
***REMOVED***

// Get searches result for the specified path.
// The result should be a JSON array or object.
func (t Result) Get(path string) Result ***REMOVED***
	r := Get(t.Raw, path)
	if r.Indexes != nil ***REMOVED***
		for i := 0; i < len(r.Indexes); i++ ***REMOVED***
			r.Indexes[i] += t.Index
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r.Index += t.Index
	***REMOVED***
	return r
***REMOVED***

type arrayOrMapResult struct ***REMOVED***
	a  []Result
	ai []interface***REMOVED******REMOVED***
	o  map[string]Result
	oi map[string]interface***REMOVED******REMOVED***
	vc byte
***REMOVED***

func (t Result) arrayOrMap(vc byte, valueize bool) (r arrayOrMapResult) ***REMOVED***
	var json = t.Raw
	var i int
	var value Result
	var count int
	var key Result
	if vc == 0 ***REMOVED***
		for ; i < len(json); i++ ***REMOVED***
			if json[i] == '***REMOVED***' || json[i] == '[' ***REMOVED***
				r.vc = json[i]
				i++
				break
			***REMOVED***
			if json[i] > ' ' ***REMOVED***
				goto end
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for ; i < len(json); i++ ***REMOVED***
			if json[i] == vc ***REMOVED***
				i++
				break
			***REMOVED***
			if json[i] > ' ' ***REMOVED***
				goto end
			***REMOVED***
		***REMOVED***
		r.vc = vc
	***REMOVED***
	if r.vc == '***REMOVED***' ***REMOVED***
		if valueize ***REMOVED***
			r.oi = make(map[string]interface***REMOVED******REMOVED***)
		***REMOVED*** else ***REMOVED***
			r.o = make(map[string]Result)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if valueize ***REMOVED***
			r.ai = make([]interface***REMOVED******REMOVED***, 0)
		***REMOVED*** else ***REMOVED***
			r.a = make([]Result, 0)
		***REMOVED***
	***REMOVED***
	for ; i < len(json); i++ ***REMOVED***
		if json[i] <= ' ' ***REMOVED***
			continue
		***REMOVED***
		// get next value
		if json[i] == ']' || json[i] == '***REMOVED***' ***REMOVED***
			break
		***REMOVED***
		switch json[i] ***REMOVED***
		default:
			if (json[i] >= '0' && json[i] <= '9') || json[i] == '-' ***REMOVED***
				value.Type = Number
				value.Raw, value.Num = tonum(json[i:])
				value.Str = ""
			***REMOVED*** else ***REMOVED***
				continue
			***REMOVED***
		case '***REMOVED***', '[':
			value.Type = JSON
			value.Raw = squash(json[i:])
			value.Str, value.Num = "", 0
		case 'n':
			value.Type = Null
			value.Raw = tolit(json[i:])
			value.Str, value.Num = "", 0
		case 't':
			value.Type = True
			value.Raw = tolit(json[i:])
			value.Str, value.Num = "", 0
		case 'f':
			value.Type = False
			value.Raw = tolit(json[i:])
			value.Str, value.Num = "", 0
		case '"':
			value.Type = String
			value.Raw, value.Str = tostr(json[i:])
			value.Num = 0
		***REMOVED***
		value.Index = i + t.Index

		i += len(value.Raw) - 1

		if r.vc == '***REMOVED***' ***REMOVED***
			if count%2 == 0 ***REMOVED***
				key = value
			***REMOVED*** else ***REMOVED***
				if valueize ***REMOVED***
					if _, ok := r.oi[key.Str]; !ok ***REMOVED***
						r.oi[key.Str] = value.Value()
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					if _, ok := r.o[key.Str]; !ok ***REMOVED***
						r.o[key.Str] = value
					***REMOVED***
				***REMOVED***
			***REMOVED***
			count++
		***REMOVED*** else ***REMOVED***
			if valueize ***REMOVED***
				r.ai = append(r.ai, value.Value())
			***REMOVED*** else ***REMOVED***
				r.a = append(r.a, value)
			***REMOVED***
		***REMOVED***
	***REMOVED***
end:
	if t.Indexes != nil ***REMOVED***
		if len(t.Indexes) != len(r.a) ***REMOVED***
			for i := 0; i < len(r.a); i++ ***REMOVED***
				r.a[i].Index = 0
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for i := 0; i < len(r.a); i++ ***REMOVED***
				r.a[i].Index = t.Indexes[i]
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// Parse parses the json and returns a result.
//
// This function expects that the json is well-formed, and does not validate.
// Invalid json will not panic, but it may return back unexpected results.
// If you are consuming JSON from an unpredictable source then you may want to
// use the Valid function first.
func Parse(json string) Result ***REMOVED***
	var value Result
	i := 0
	for ; i < len(json); i++ ***REMOVED***
		if json[i] == '***REMOVED***' || json[i] == '[' ***REMOVED***
			value.Type = JSON
			value.Raw = json[i:] // just take the entire raw
			break
		***REMOVED***
		if json[i] <= ' ' ***REMOVED***
			continue
		***REMOVED***
		switch json[i] ***REMOVED***
		case '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'i', 'I', 'N':
			value.Type = Number
			value.Raw, value.Num = tonum(json[i:])
		case 'n':
			if i+1 < len(json) && json[i+1] != 'u' ***REMOVED***
				// nan
				value.Type = Number
				value.Raw, value.Num = tonum(json[i:])
			***REMOVED*** else ***REMOVED***
				// null
				value.Type = Null
				value.Raw = tolit(json[i:])
			***REMOVED***
		case 't':
			value.Type = True
			value.Raw = tolit(json[i:])
		case 'f':
			value.Type = False
			value.Raw = tolit(json[i:])
		case '"':
			value.Type = String
			value.Raw, value.Str = tostr(json[i:])
		default:
			return Result***REMOVED******REMOVED***
		***REMOVED***
		break
	***REMOVED***
	if value.Exists() ***REMOVED***
		value.Index = i
	***REMOVED***
	return value
***REMOVED***

// ParseBytes parses the json and returns a result.
// If working with bytes, this method preferred over Parse(string(data))
func ParseBytes(json []byte) Result ***REMOVED***
	return Parse(string(json))
***REMOVED***

func squash(json string) string ***REMOVED***
	// expects that the lead character is a '[' or '***REMOVED***' or '(' or '"'
	// squash the value, ignoring all nested arrays and objects.
	var i, depth int
	if json[0] != '"' ***REMOVED***
		i, depth = 1, 1
	***REMOVED***
	for ; i < len(json); i++ ***REMOVED***
		if json[i] >= '"' && json[i] <= '***REMOVED***' ***REMOVED***
			switch json[i] ***REMOVED***
			case '"':
				i++
				s2 := i
				for ; i < len(json); i++ ***REMOVED***
					if json[i] > '\\' ***REMOVED***
						continue
					***REMOVED***
					if json[i] == '"' ***REMOVED***
						// look for an escaped slash
						if json[i-1] == '\\' ***REMOVED***
							n := 0
							for j := i - 2; j > s2-1; j-- ***REMOVED***
								if json[j] != '\\' ***REMOVED***
									break
								***REMOVED***
								n++
							***REMOVED***
							if n%2 == 0 ***REMOVED***
								continue
							***REMOVED***
						***REMOVED***
						break
					***REMOVED***
				***REMOVED***
				if depth == 0 ***REMOVED***
					if i >= len(json) ***REMOVED***
						return json
					***REMOVED***
					return json[:i+1]
				***REMOVED***
			case '***REMOVED***', '[', '(':
				depth++
			case '***REMOVED***', ']', ')':
				depth--
				if depth == 0 ***REMOVED***
					return json[:i+1]
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return json
***REMOVED***

func tonum(json string) (raw string, num float64) ***REMOVED***
	for i := 1; i < len(json); i++ ***REMOVED***
		// less than dash might have valid characters
		if json[i] <= '-' ***REMOVED***
			if json[i] <= ' ' || json[i] == ',' ***REMOVED***
				// break on whitespace and comma
				raw = json[:i]
				num, _ = strconv.ParseFloat(raw, 64)
				return
			***REMOVED***
			// could be a '+' or '-'. let's assume so.
		***REMOVED*** else if json[i] == ']' || json[i] == '***REMOVED***' ***REMOVED***
			// break on ']' or '***REMOVED***'
			raw = json[:i]
			num, _ = strconv.ParseFloat(raw, 64)
			return
		***REMOVED***
	***REMOVED***
	raw = json
	num, _ = strconv.ParseFloat(raw, 64)
	return
***REMOVED***

func tolit(json string) (raw string) ***REMOVED***
	for i := 1; i < len(json); i++ ***REMOVED***
		if json[i] < 'a' || json[i] > 'z' ***REMOVED***
			return json[:i]
		***REMOVED***
	***REMOVED***
	return json
***REMOVED***

func tostr(json string) (raw string, str string) ***REMOVED***
	// expects that the lead character is a '"'
	for i := 1; i < len(json); i++ ***REMOVED***
		if json[i] > '\\' ***REMOVED***
			continue
		***REMOVED***
		if json[i] == '"' ***REMOVED***
			return json[:i+1], json[1:i]
		***REMOVED***
		if json[i] == '\\' ***REMOVED***
			i++
			for ; i < len(json); i++ ***REMOVED***
				if json[i] > '\\' ***REMOVED***
					continue
				***REMOVED***
				if json[i] == '"' ***REMOVED***
					// look for an escaped slash
					if json[i-1] == '\\' ***REMOVED***
						n := 0
						for j := i - 2; j > 0; j-- ***REMOVED***
							if json[j] != '\\' ***REMOVED***
								break
							***REMOVED***
							n++
						***REMOVED***
						if n%2 == 0 ***REMOVED***
							continue
						***REMOVED***
					***REMOVED***
					return json[:i+1], unescape(json[1:i])
				***REMOVED***
			***REMOVED***
			var ret string
			if i+1 < len(json) ***REMOVED***
				ret = json[:i+1]
			***REMOVED*** else ***REMOVED***
				ret = json[:i]
			***REMOVED***
			return ret, unescape(json[1:i])
		***REMOVED***
	***REMOVED***
	return json, json[1:]
***REMOVED***

// Exists returns true if value exists.
//
//  if gjson.Get(json, "name.last").Exists()***REMOVED***
//		println("value exists")
//  ***REMOVED***
func (t Result) Exists() bool ***REMOVED***
	return t.Type != Null || len(t.Raw) != 0
***REMOVED***

// Value returns one of these types:
//
//	bool, for JSON booleans
//	float64, for JSON numbers
//	Number, for JSON numbers
//	string, for JSON string literals
//	nil, for JSON null
//	map[string]interface***REMOVED******REMOVED***, for JSON objects
//	[]interface***REMOVED******REMOVED***, for JSON arrays
//
func (t Result) Value() interface***REMOVED******REMOVED*** ***REMOVED***
	if t.Type == String ***REMOVED***
		return t.Str
	***REMOVED***
	switch t.Type ***REMOVED***
	default:
		return nil
	case False:
		return false
	case Number:
		return t.Num
	case JSON:
		r := t.arrayOrMap(0, true)
		if r.vc == '***REMOVED***' ***REMOVED***
			return r.oi
		***REMOVED*** else if r.vc == '[' ***REMOVED***
			return r.ai
		***REMOVED***
		return nil
	case True:
		return true
	***REMOVED***
***REMOVED***

func parseString(json string, i int) (int, string, bool, bool) ***REMOVED***
	var s = i
	for ; i < len(json); i++ ***REMOVED***
		if json[i] > '\\' ***REMOVED***
			continue
		***REMOVED***
		if json[i] == '"' ***REMOVED***
			return i + 1, json[s-1 : i+1], false, true
		***REMOVED***
		if json[i] == '\\' ***REMOVED***
			i++
			for ; i < len(json); i++ ***REMOVED***
				if json[i] > '\\' ***REMOVED***
					continue
				***REMOVED***
				if json[i] == '"' ***REMOVED***
					// look for an escaped slash
					if json[i-1] == '\\' ***REMOVED***
						n := 0
						for j := i - 2; j > 0; j-- ***REMOVED***
							if json[j] != '\\' ***REMOVED***
								break
							***REMOVED***
							n++
						***REMOVED***
						if n%2 == 0 ***REMOVED***
							continue
						***REMOVED***
					***REMOVED***
					return i + 1, json[s-1 : i+1], true, true
				***REMOVED***
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return i, json[s-1:], false, false
***REMOVED***

func parseNumber(json string, i int) (int, string) ***REMOVED***
	var s = i
	i++
	for ; i < len(json); i++ ***REMOVED***
		if json[i] <= ' ' || json[i] == ',' || json[i] == ']' ||
			json[i] == '***REMOVED***' ***REMOVED***
			return i, json[s:i]
		***REMOVED***
	***REMOVED***
	return i, json[s:]
***REMOVED***

func parseLiteral(json string, i int) (int, string) ***REMOVED***
	var s = i
	i++
	for ; i < len(json); i++ ***REMOVED***
		if json[i] < 'a' || json[i] > 'z' ***REMOVED***
			return i, json[s:i]
		***REMOVED***
	***REMOVED***
	return i, json[s:]
***REMOVED***

type arrayPathResult struct ***REMOVED***
	part    string
	path    string
	pipe    string
	piped   bool
	more    bool
	alogok  bool
	arrch   bool
	alogkey string
	query   struct ***REMOVED***
		on    bool
		all   bool
		path  string
		op    string
		value string
	***REMOVED***
***REMOVED***

func parseArrayPath(path string) (r arrayPathResult) ***REMOVED***
	for i := 0; i < len(path); i++ ***REMOVED***
		if path[i] == '|' ***REMOVED***
			r.part = path[:i]
			r.pipe = path[i+1:]
			r.piped = true
			return
		***REMOVED***
		if path[i] == '.' ***REMOVED***
			r.part = path[:i]
			if !r.arrch && i < len(path)-1 && isDotPiperChar(path[i+1]) ***REMOVED***
				r.pipe = path[i+1:]
				r.piped = true
			***REMOVED*** else ***REMOVED***
				r.path = path[i+1:]
				r.more = true
			***REMOVED***
			return
		***REMOVED***
		if path[i] == '#' ***REMOVED***
			r.arrch = true
			if i == 0 && len(path) > 1 ***REMOVED***
				if path[1] == '.' ***REMOVED***
					r.alogok = true
					r.alogkey = path[2:]
					r.path = path[:1]
				***REMOVED*** else if path[1] == '[' || path[1] == '(' ***REMOVED***
					// query
					r.query.on = true
					qpath, op, value, _, fi, vesc, ok :=
						parseQuery(path[i:])
					if !ok ***REMOVED***
						// bad query, end now
						break
					***REMOVED***
					if len(value) >= 2 && value[0] == '"' &&
						value[len(value)-1] == '"' ***REMOVED***
						value = value[1 : len(value)-1]
						if vesc ***REMOVED***
							value = unescape(value)
						***REMOVED***
					***REMOVED***
					r.query.path = qpath
					r.query.op = op
					r.query.value = value

					i = fi - 1
					if i+1 < len(path) && path[i+1] == '#' ***REMOVED***
						r.query.all = true
					***REMOVED***
				***REMOVED***
			***REMOVED***
			continue
		***REMOVED***
	***REMOVED***
	r.part = path
	r.path = ""
	return
***REMOVED***

// splitQuery takes a query and splits it into three parts:
//   path, op, middle, and right.
// So for this query:
//   #(first_name=="Murphy").last
// Becomes
//   first_name   # path
//   =="Murphy"   # middle
//   .last        # right
// Or,
//   #(service_roles.#(=="one")).cap
// Becomes
//   service_roles.#(=="one")   # path
//                              # middle
//   .cap                       # right
func parseQuery(query string) (
	path, op, value, remain string, i int, vesc, ok bool,
) ***REMOVED***
	if len(query) < 2 || query[0] != '#' ||
		(query[1] != '(' && query[1] != '[') ***REMOVED***
		return "", "", "", "", i, false, false
	***REMOVED***
	i = 2
	j := 0 // start of value part
	depth := 1
	for ; i < len(query); i++ ***REMOVED***
		if depth == 1 && j == 0 ***REMOVED***
			switch query[i] ***REMOVED***
			case '!', '=', '<', '>', '%':
				// start of the value part
				j = i
				continue
			***REMOVED***
		***REMOVED***
		if query[i] == '\\' ***REMOVED***
			i++
		***REMOVED*** else if query[i] == '[' || query[i] == '(' ***REMOVED***
			depth++
		***REMOVED*** else if query[i] == ']' || query[i] == ')' ***REMOVED***
			depth--
			if depth == 0 ***REMOVED***
				break
			***REMOVED***
		***REMOVED*** else if query[i] == '"' ***REMOVED***
			// inside selector string, balance quotes
			i++
			for ; i < len(query); i++ ***REMOVED***
				if query[i] == '\\' ***REMOVED***
					vesc = true
					i++
				***REMOVED*** else if query[i] == '"' ***REMOVED***
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if depth > 0 ***REMOVED***
		return "", "", "", "", i, false, false
	***REMOVED***
	if j > 0 ***REMOVED***
		path = trim(query[2:j])
		value = trim(query[j:i])
		remain = query[i+1:]
		// parse the compare op from the value
		var opsz int
		switch ***REMOVED***
		case len(value) == 1:
			opsz = 1
		case value[0] == '!' && value[1] == '=':
			opsz = 2
		case value[0] == '!' && value[1] == '%':
			opsz = 2
		case value[0] == '<' && value[1] == '=':
			opsz = 2
		case value[0] == '>' && value[1] == '=':
			opsz = 2
		case value[0] == '=' && value[1] == '=':
			value = value[1:]
			opsz = 1
		case value[0] == '<':
			opsz = 1
		case value[0] == '>':
			opsz = 1
		case value[0] == '=':
			opsz = 1
		case value[0] == '%':
			opsz = 1
		***REMOVED***
		op = value[:opsz]
		value = trim(value[opsz:])
	***REMOVED*** else ***REMOVED***
		path = trim(query[2:i])
		remain = query[i+1:]
	***REMOVED***
	return path, op, value, remain, i + 1, vesc, true
***REMOVED***

func trim(s string) string ***REMOVED***
left:
	if len(s) > 0 && s[0] <= ' ' ***REMOVED***
		s = s[1:]
		goto left
	***REMOVED***
right:
	if len(s) > 0 && s[len(s)-1] <= ' ' ***REMOVED***
		s = s[:len(s)-1]
		goto right
	***REMOVED***
	return s
***REMOVED***

// peek at the next byte and see if it's a '@', '[', or '***REMOVED***'.
func isDotPiperChar(c byte) bool ***REMOVED***
	return !DisableModifiers && (c == '@' || c == '[' || c == '***REMOVED***')
***REMOVED***

type objectPathResult struct ***REMOVED***
	part  string
	path  string
	pipe  string
	piped bool
	wild  bool
	more  bool
***REMOVED***

func parseObjectPath(path string) (r objectPathResult) ***REMOVED***
	for i := 0; i < len(path); i++ ***REMOVED***
		if path[i] == '|' ***REMOVED***
			r.part = path[:i]
			r.pipe = path[i+1:]
			r.piped = true
			return
		***REMOVED***
		if path[i] == '.' ***REMOVED***
			r.part = path[:i]
			if i < len(path)-1 && isDotPiperChar(path[i+1]) ***REMOVED***
				r.pipe = path[i+1:]
				r.piped = true
			***REMOVED*** else ***REMOVED***
				r.path = path[i+1:]
				r.more = true
			***REMOVED***
			return
		***REMOVED***
		if path[i] == '*' || path[i] == '?' ***REMOVED***
			r.wild = true
			continue
		***REMOVED***
		if path[i] == '\\' ***REMOVED***
			// go into escape mode. this is a slower path that
			// strips off the escape character from the part.
			epart := []byte(path[:i])
			i++
			if i < len(path) ***REMOVED***
				epart = append(epart, path[i])
				i++
				for ; i < len(path); i++ ***REMOVED***
					if path[i] == '\\' ***REMOVED***
						i++
						if i < len(path) ***REMOVED***
							epart = append(epart, path[i])
						***REMOVED***
						continue
					***REMOVED*** else if path[i] == '.' ***REMOVED***
						r.part = string(epart)
						if i < len(path)-1 && isDotPiperChar(path[i+1]) ***REMOVED***
							r.pipe = path[i+1:]
							r.piped = true
						***REMOVED*** else ***REMOVED***
							r.path = path[i+1:]
						***REMOVED***
						r.more = true
						return
					***REMOVED*** else if path[i] == '|' ***REMOVED***
						r.part = string(epart)
						r.pipe = path[i+1:]
						r.piped = true
						return
					***REMOVED*** else if path[i] == '*' || path[i] == '?' ***REMOVED***
						r.wild = true
					***REMOVED***
					epart = append(epart, path[i])
				***REMOVED***
			***REMOVED***
			// append the last part
			r.part = string(epart)
			return
		***REMOVED***
	***REMOVED***
	r.part = path
	return
***REMOVED***

func parseSquash(json string, i int) (int, string) ***REMOVED***
	// expects that the lead character is a '[' or '***REMOVED***' or '('
	// squash the value, ignoring all nested arrays and objects.
	// the first '[' or '***REMOVED***' or '(' has already been read
	s := i
	i++
	depth := 1
	for ; i < len(json); i++ ***REMOVED***
		if json[i] >= '"' && json[i] <= '***REMOVED***' ***REMOVED***
			switch json[i] ***REMOVED***
			case '"':
				i++
				s2 := i
				for ; i < len(json); i++ ***REMOVED***
					if json[i] > '\\' ***REMOVED***
						continue
					***REMOVED***
					if json[i] == '"' ***REMOVED***
						// look for an escaped slash
						if json[i-1] == '\\' ***REMOVED***
							n := 0
							for j := i - 2; j > s2-1; j-- ***REMOVED***
								if json[j] != '\\' ***REMOVED***
									break
								***REMOVED***
								n++
							***REMOVED***
							if n%2 == 0 ***REMOVED***
								continue
							***REMOVED***
						***REMOVED***
						break
					***REMOVED***
				***REMOVED***
			case '***REMOVED***', '[', '(':
				depth++
			case '***REMOVED***', ']', ')':
				depth--
				if depth == 0 ***REMOVED***
					i++
					return i, json[s:i]
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return i, json[s:]
***REMOVED***

func parseObject(c *parseContext, i int, path string) (int, bool) ***REMOVED***
	var pmatch, kesc, vesc, ok, hit bool
	var key, val string
	rp := parseObjectPath(path)
	if !rp.more && rp.piped ***REMOVED***
		c.pipe = rp.pipe
		c.piped = true
	***REMOVED***
	for i < len(c.json) ***REMOVED***
		for ; i < len(c.json); i++ ***REMOVED***
			if c.json[i] == '"' ***REMOVED***
				// parse_key_string
				// this is slightly different from getting s string value
				// because we don't need the outer quotes.
				i++
				var s = i
				for ; i < len(c.json); i++ ***REMOVED***
					if c.json[i] > '\\' ***REMOVED***
						continue
					***REMOVED***
					if c.json[i] == '"' ***REMOVED***
						i, key, kesc, ok = i+1, c.json[s:i], false, true
						goto parse_key_string_done
					***REMOVED***
					if c.json[i] == '\\' ***REMOVED***
						i++
						for ; i < len(c.json); i++ ***REMOVED***
							if c.json[i] > '\\' ***REMOVED***
								continue
							***REMOVED***
							if c.json[i] == '"' ***REMOVED***
								// look for an escaped slash
								if c.json[i-1] == '\\' ***REMOVED***
									n := 0
									for j := i - 2; j > 0; j-- ***REMOVED***
										if c.json[j] != '\\' ***REMOVED***
											break
										***REMOVED***
										n++
									***REMOVED***
									if n%2 == 0 ***REMOVED***
										continue
									***REMOVED***
								***REMOVED***
								i, key, kesc, ok = i+1, c.json[s:i], true, true
								goto parse_key_string_done
							***REMOVED***
						***REMOVED***
						break
					***REMOVED***
				***REMOVED***
				key, kesc, ok = c.json[s:], false, false
			parse_key_string_done:
				break
			***REMOVED***
			if c.json[i] == '***REMOVED***' ***REMOVED***
				return i + 1, false
			***REMOVED***
		***REMOVED***
		if !ok ***REMOVED***
			return i, false
		***REMOVED***
		if rp.wild ***REMOVED***
			if kesc ***REMOVED***
				pmatch = matchLimit(unescape(key), rp.part)
			***REMOVED*** else ***REMOVED***
				pmatch = matchLimit(key, rp.part)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if kesc ***REMOVED***
				pmatch = rp.part == unescape(key)
			***REMOVED*** else ***REMOVED***
				pmatch = rp.part == key
			***REMOVED***
		***REMOVED***
		hit = pmatch && !rp.more
		for ; i < len(c.json); i++ ***REMOVED***
			var num bool
			switch c.json[i] ***REMOVED***
			default:
				continue
			case '"':
				i++
				i, val, vesc, ok = parseString(c.json, i)
				if !ok ***REMOVED***
					return i, false
				***REMOVED***
				if hit ***REMOVED***
					if vesc ***REMOVED***
						c.value.Str = unescape(val[1 : len(val)-1])
					***REMOVED*** else ***REMOVED***
						c.value.Str = val[1 : len(val)-1]
					***REMOVED***
					c.value.Raw = val
					c.value.Type = String
					return i, true
				***REMOVED***
			case '***REMOVED***':
				if pmatch && !hit ***REMOVED***
					i, hit = parseObject(c, i+1, rp.path)
					if hit ***REMOVED***
						return i, true
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					i, val = parseSquash(c.json, i)
					if hit ***REMOVED***
						c.value.Raw = val
						c.value.Type = JSON
						return i, true
					***REMOVED***
				***REMOVED***
			case '[':
				if pmatch && !hit ***REMOVED***
					i, hit = parseArray(c, i+1, rp.path)
					if hit ***REMOVED***
						return i, true
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					i, val = parseSquash(c.json, i)
					if hit ***REMOVED***
						c.value.Raw = val
						c.value.Type = JSON
						return i, true
					***REMOVED***
				***REMOVED***
			case 'n':
				if i+1 < len(c.json) && c.json[i+1] != 'u' ***REMOVED***
					num = true
					break
				***REMOVED***
				fallthrough
			case 't', 'f':
				vc := c.json[i]
				i, val = parseLiteral(c.json, i)
				if hit ***REMOVED***
					c.value.Raw = val
					switch vc ***REMOVED***
					case 't':
						c.value.Type = True
					case 'f':
						c.value.Type = False
					***REMOVED***
					return i, true
				***REMOVED***
			case '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
				'i', 'I', 'N':
				num = true
			***REMOVED***
			if num ***REMOVED***
				i, val = parseNumber(c.json, i)
				if hit ***REMOVED***
					c.value.Raw = val
					c.value.Type = Number
					c.value.Num, _ = strconv.ParseFloat(val, 64)
					return i, true
				***REMOVED***
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return i, false
***REMOVED***

// matchLimit will limit the complexity of the match operation to avoid ReDos
// attacks from arbritary inputs.
// See the github.com/tidwall/match.MatchLimit function for more information.
func matchLimit(str, pattern string) bool ***REMOVED***
	matched, _ := match.MatchLimit(str, pattern, 10000)
	return matched
***REMOVED***

func queryMatches(rp *arrayPathResult, value Result) bool ***REMOVED***
	rpv := rp.query.value
	if len(rpv) > 0 && rpv[0] == '~' ***REMOVED***
		// convert to bool
		rpv = rpv[1:]
		if value.Bool() ***REMOVED***
			value = Result***REMOVED***Type: True***REMOVED***
		***REMOVED*** else ***REMOVED***
			value = Result***REMOVED***Type: False***REMOVED***
		***REMOVED***
	***REMOVED***
	if !value.Exists() ***REMOVED***
		return false
	***REMOVED***
	if rp.query.op == "" ***REMOVED***
		// the query is only looking for existence, such as:
		//   friends.#(name)
		// which makes sure that the array "friends" has an element of
		// "name" that exists
		return true
	***REMOVED***
	switch value.Type ***REMOVED***
	case String:
		switch rp.query.op ***REMOVED***
		case "=":
			return value.Str == rpv
		case "!=":
			return value.Str != rpv
		case "<":
			return value.Str < rpv
		case "<=":
			return value.Str <= rpv
		case ">":
			return value.Str > rpv
		case ">=":
			return value.Str >= rpv
		case "%":
			return matchLimit(value.Str, rpv)
		case "!%":
			return !matchLimit(value.Str, rpv)
		***REMOVED***
	case Number:
		rpvn, _ := strconv.ParseFloat(rpv, 64)
		switch rp.query.op ***REMOVED***
		case "=":
			return value.Num == rpvn
		case "!=":
			return value.Num != rpvn
		case "<":
			return value.Num < rpvn
		case "<=":
			return value.Num <= rpvn
		case ">":
			return value.Num > rpvn
		case ">=":
			return value.Num >= rpvn
		***REMOVED***
	case True:
		switch rp.query.op ***REMOVED***
		case "=":
			return rpv == "true"
		case "!=":
			return rpv != "true"
		case ">":
			return rpv == "false"
		case ">=":
			return true
		***REMOVED***
	case False:
		switch rp.query.op ***REMOVED***
		case "=":
			return rpv == "false"
		case "!=":
			return rpv != "false"
		case "<":
			return rpv == "true"
		case "<=":
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
func parseArray(c *parseContext, i int, path string) (int, bool) ***REMOVED***
	var pmatch, vesc, ok, hit bool
	var val string
	var h int
	var alog []int
	var partidx int
	var multires []byte
	var queryIndexes []int
	rp := parseArrayPath(path)
	if !rp.arrch ***REMOVED***
		n, ok := parseUint(rp.part)
		if !ok ***REMOVED***
			partidx = -1
		***REMOVED*** else ***REMOVED***
			partidx = int(n)
		***REMOVED***
	***REMOVED***
	if !rp.more && rp.piped ***REMOVED***
		c.pipe = rp.pipe
		c.piped = true
	***REMOVED***

	procQuery := func(qval Result) bool ***REMOVED***
		if rp.query.all ***REMOVED***
			if len(multires) == 0 ***REMOVED***
				multires = append(multires, '[')
			***REMOVED***
		***REMOVED***
		var tmp parseContext
		tmp.value = qval
		fillIndex(c.json, &tmp)
		parentIndex := tmp.value.Index
		var res Result
		if qval.Type == JSON ***REMOVED***
			res = qval.Get(rp.query.path)
		***REMOVED*** else ***REMOVED***
			if rp.query.path != "" ***REMOVED***
				return false
			***REMOVED***
			res = qval
		***REMOVED***
		if queryMatches(&rp, res) ***REMOVED***
			if rp.more ***REMOVED***
				left, right, ok := splitPossiblePipe(rp.path)
				if ok ***REMOVED***
					rp.path = left
					c.pipe = right
					c.piped = true
				***REMOVED***
				res = qval.Get(rp.path)
			***REMOVED*** else ***REMOVED***
				res = qval
			***REMOVED***
			if rp.query.all ***REMOVED***
				raw := res.Raw
				if len(raw) == 0 ***REMOVED***
					raw = res.String()
				***REMOVED***
				if raw != "" ***REMOVED***
					if len(multires) > 1 ***REMOVED***
						multires = append(multires, ',')
					***REMOVED***
					multires = append(multires, raw...)
					queryIndexes = append(queryIndexes, res.Index+parentIndex)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				c.value = res
				return true
			***REMOVED***
		***REMOVED***
		return false
	***REMOVED***
	for i < len(c.json)+1 ***REMOVED***
		if !rp.arrch ***REMOVED***
			pmatch = partidx == h
			hit = pmatch && !rp.more
		***REMOVED***
		h++
		if rp.alogok ***REMOVED***
			alog = append(alog, i)
		***REMOVED***
		for ; ; i++ ***REMOVED***
			var ch byte
			if i > len(c.json) ***REMOVED***
				break
			***REMOVED*** else if i == len(c.json) ***REMOVED***
				ch = ']'
			***REMOVED*** else ***REMOVED***
				ch = c.json[i]
			***REMOVED***
			var num bool
			switch ch ***REMOVED***
			default:
				continue
			case '"':
				i++
				i, val, vesc, ok = parseString(c.json, i)
				if !ok ***REMOVED***
					return i, false
				***REMOVED***
				if rp.query.on ***REMOVED***
					var qval Result
					if vesc ***REMOVED***
						qval.Str = unescape(val[1 : len(val)-1])
					***REMOVED*** else ***REMOVED***
						qval.Str = val[1 : len(val)-1]
					***REMOVED***
					qval.Raw = val
					qval.Type = String
					if procQuery(qval) ***REMOVED***
						return i, true
					***REMOVED***
				***REMOVED*** else if hit ***REMOVED***
					if rp.alogok ***REMOVED***
						break
					***REMOVED***
					if vesc ***REMOVED***
						c.value.Str = unescape(val[1 : len(val)-1])
					***REMOVED*** else ***REMOVED***
						c.value.Str = val[1 : len(val)-1]
					***REMOVED***
					c.value.Raw = val
					c.value.Type = String
					return i, true
				***REMOVED***
			case '***REMOVED***':
				if pmatch && !hit ***REMOVED***
					i, hit = parseObject(c, i+1, rp.path)
					if hit ***REMOVED***
						if rp.alogok ***REMOVED***
							break
						***REMOVED***
						return i, true
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					i, val = parseSquash(c.json, i)
					if rp.query.on ***REMOVED***
						if procQuery(Result***REMOVED***Raw: val, Type: JSON***REMOVED***) ***REMOVED***
							return i, true
						***REMOVED***
					***REMOVED*** else if hit ***REMOVED***
						if rp.alogok ***REMOVED***
							break
						***REMOVED***
						c.value.Raw = val
						c.value.Type = JSON
						return i, true
					***REMOVED***
				***REMOVED***
			case '[':
				if pmatch && !hit ***REMOVED***
					i, hit = parseArray(c, i+1, rp.path)
					if hit ***REMOVED***
						if rp.alogok ***REMOVED***
							break
						***REMOVED***
						return i, true
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					i, val = parseSquash(c.json, i)
					if rp.query.on ***REMOVED***
						if procQuery(Result***REMOVED***Raw: val, Type: JSON***REMOVED***) ***REMOVED***
							return i, true
						***REMOVED***
					***REMOVED*** else if hit ***REMOVED***
						if rp.alogok ***REMOVED***
							break
						***REMOVED***
						c.value.Raw = val
						c.value.Type = JSON
						return i, true
					***REMOVED***
				***REMOVED***
			case 'n':
				if i+1 < len(c.json) && c.json[i+1] != 'u' ***REMOVED***
					num = true
					break
				***REMOVED***
				fallthrough
			case 't', 'f':
				vc := c.json[i]
				i, val = parseLiteral(c.json, i)
				if rp.query.on ***REMOVED***
					var qval Result
					qval.Raw = val
					switch vc ***REMOVED***
					case 't':
						qval.Type = True
					case 'f':
						qval.Type = False
					***REMOVED***
					if procQuery(qval) ***REMOVED***
						return i, true
					***REMOVED***
				***REMOVED*** else if hit ***REMOVED***
					if rp.alogok ***REMOVED***
						break
					***REMOVED***
					c.value.Raw = val
					switch vc ***REMOVED***
					case 't':
						c.value.Type = True
					case 'f':
						c.value.Type = False
					***REMOVED***
					return i, true
				***REMOVED***
			case '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
				'i', 'I', 'N':
				num = true
			case ']':
				if rp.arrch && rp.part == "#" ***REMOVED***
					if rp.alogok ***REMOVED***
						left, right, ok := splitPossiblePipe(rp.alogkey)
						if ok ***REMOVED***
							rp.alogkey = left
							c.pipe = right
							c.piped = true
						***REMOVED***
						var indexes = make([]int, 0, 64)
						var jsons = make([]byte, 0, 64)
						jsons = append(jsons, '[')
						for j, k := 0, 0; j < len(alog); j++ ***REMOVED***
							idx := alog[j]
							for idx < len(c.json) ***REMOVED***
								switch c.json[idx] ***REMOVED***
								case ' ', '\t', '\r', '\n':
									idx++
									continue
								***REMOVED***
								break
							***REMOVED***
							if idx < len(c.json) && c.json[idx] != ']' ***REMOVED***
								_, res, ok := parseAny(c.json, idx, true)
								if ok ***REMOVED***
									res := res.Get(rp.alogkey)
									if res.Exists() ***REMOVED***
										if k > 0 ***REMOVED***
											jsons = append(jsons, ',')
										***REMOVED***
										raw := res.Raw
										if len(raw) == 0 ***REMOVED***
											raw = res.String()
										***REMOVED***
										jsons = append(jsons, []byte(raw)...)
										indexes = append(indexes, res.Index)
										k++
									***REMOVED***
								***REMOVED***
							***REMOVED***
						***REMOVED***
						jsons = append(jsons, ']')
						c.value.Type = JSON
						c.value.Raw = string(jsons)
						c.value.Indexes = indexes
						return i + 1, true
					***REMOVED***
					if rp.alogok ***REMOVED***
						break
					***REMOVED***

					c.value.Type = Number
					c.value.Num = float64(h - 1)
					c.value.Raw = strconv.Itoa(h - 1)
					c.calcd = true
					return i + 1, true
				***REMOVED***
				if !c.value.Exists() ***REMOVED***
					if len(multires) > 0 ***REMOVED***
						c.value = Result***REMOVED***
							Raw:     string(append(multires, ']')),
							Type:    JSON,
							Indexes: queryIndexes,
						***REMOVED***
					***REMOVED*** else if rp.query.all ***REMOVED***
						c.value = Result***REMOVED***
							Raw:  "[]",
							Type: JSON,
						***REMOVED***
					***REMOVED***
				***REMOVED***
				return i + 1, false
			***REMOVED***
			if num ***REMOVED***
				i, val = parseNumber(c.json, i)
				if rp.query.on ***REMOVED***
					var qval Result
					qval.Raw = val
					qval.Type = Number
					qval.Num, _ = strconv.ParseFloat(val, 64)
					if procQuery(qval) ***REMOVED***
						return i, true
					***REMOVED***
				***REMOVED*** else if hit ***REMOVED***
					if rp.alogok ***REMOVED***
						break
					***REMOVED***
					c.value.Raw = val
					c.value.Type = Number
					c.value.Num, _ = strconv.ParseFloat(val, 64)
					return i, true
				***REMOVED***
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return i, false
***REMOVED***

func splitPossiblePipe(path string) (left, right string, ok bool) ***REMOVED***
	// take a quick peek for the pipe character. If found we'll split the piped
	// part of the path into the c.pipe field and shorten the rp.
	var possible bool
	for i := 0; i < len(path); i++ ***REMOVED***
		if path[i] == '|' ***REMOVED***
			possible = true
			break
		***REMOVED***
	***REMOVED***
	if !possible ***REMOVED***
		return
	***REMOVED***

	if len(path) > 0 && path[0] == '***REMOVED***' ***REMOVED***
		squashed := squash(path[1:])
		if len(squashed) < len(path)-1 ***REMOVED***
			squashed = path[:len(squashed)+1]
			remain := path[len(squashed):]
			if remain[0] == '|' ***REMOVED***
				return squashed, remain[1:], true
			***REMOVED***
		***REMOVED***
		return
	***REMOVED***

	// split the left and right side of the path with the pipe character as
	// the delimiter. This is a little tricky because we'll need to basically
	// parse the entire path.
	for i := 0; i < len(path); i++ ***REMOVED***
		if path[i] == '\\' ***REMOVED***
			i++
		***REMOVED*** else if path[i] == '.' ***REMOVED***
			if i == len(path)-1 ***REMOVED***
				return
			***REMOVED***
			if path[i+1] == '#' ***REMOVED***
				i += 2
				if i == len(path) ***REMOVED***
					return
				***REMOVED***
				if path[i] == '[' || path[i] == '(' ***REMOVED***
					var start, end byte
					if path[i] == '[' ***REMOVED***
						start, end = '[', ']'
					***REMOVED*** else ***REMOVED***
						start, end = '(', ')'
					***REMOVED***
					// inside selector, balance brackets
					i++
					depth := 1
					for ; i < len(path); i++ ***REMOVED***
						if path[i] == '\\' ***REMOVED***
							i++
						***REMOVED*** else if path[i] == start ***REMOVED***
							depth++
						***REMOVED*** else if path[i] == end ***REMOVED***
							depth--
							if depth == 0 ***REMOVED***
								break
							***REMOVED***
						***REMOVED*** else if path[i] == '"' ***REMOVED***
							// inside selector string, balance quotes
							i++
							for ; i < len(path); i++ ***REMOVED***
								if path[i] == '\\' ***REMOVED***
									i++
								***REMOVED*** else if path[i] == '"' ***REMOVED***
									break
								***REMOVED***
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED*** else if path[i] == '|' ***REMOVED***
			return path[:i], path[i+1:], true
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// ForEachLine iterates through lines of JSON as specified by the JSON Lines
// format (http://jsonlines.org/).
// Each line is returned as a GJSON Result.
func ForEachLine(json string, iterator func(line Result) bool) ***REMOVED***
	var res Result
	var i int
	for ***REMOVED***
		i, res, _ = parseAny(json, i, true)
		if !res.Exists() ***REMOVED***
			break
		***REMOVED***
		if !iterator(res) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

type subSelector struct ***REMOVED***
	name string
	path string
***REMOVED***

// parseSubSelectors returns the subselectors belonging to a '[path1,path2]' or
// '***REMOVED***"field1":path1,"field2":path2***REMOVED***' type subSelection. It's expected that the
// first character in path is either '[' or '***REMOVED***', and has already been checked
// prior to calling this function.
func parseSubSelectors(path string) (sels []subSelector, out string, ok bool) ***REMOVED***
	modifier := 0
	depth := 1
	colon := 0
	start := 1
	i := 1
	pushSel := func() ***REMOVED***
		var sel subSelector
		if colon == 0 ***REMOVED***
			sel.path = path[start:i]
		***REMOVED*** else ***REMOVED***
			sel.name = path[start:colon]
			sel.path = path[colon+1 : i]
		***REMOVED***
		sels = append(sels, sel)
		colon = 0
		modifier = 0
		start = i + 1
	***REMOVED***
	for ; i < len(path); i++ ***REMOVED***
		switch path[i] ***REMOVED***
		case '\\':
			i++
		case '@':
			if modifier == 0 && i > 0 && (path[i-1] == '.' || path[i-1] == '|') ***REMOVED***
				modifier = i
			***REMOVED***
		case ':':
			if modifier == 0 && colon == 0 && depth == 1 ***REMOVED***
				colon = i
			***REMOVED***
		case ',':
			if depth == 1 ***REMOVED***
				pushSel()
			***REMOVED***
		case '"':
			i++
		loop:
			for ; i < len(path); i++ ***REMOVED***
				switch path[i] ***REMOVED***
				case '\\':
					i++
				case '"':
					break loop
				***REMOVED***
			***REMOVED***
		case '[', '(', '***REMOVED***':
			depth++
		case ']', ')', '***REMOVED***':
			depth--
			if depth == 0 ***REMOVED***
				pushSel()
				path = path[i+1:]
				return sels, path, true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// nameOfLast returns the name of the last component
func nameOfLast(path string) string ***REMOVED***
	for i := len(path) - 1; i >= 0; i-- ***REMOVED***
		if path[i] == '|' || path[i] == '.' ***REMOVED***
			if i > 0 ***REMOVED***
				if path[i-1] == '\\' ***REMOVED***
					continue
				***REMOVED***
			***REMOVED***
			return path[i+1:]
		***REMOVED***
	***REMOVED***
	return path
***REMOVED***

func isSimpleName(component string) bool ***REMOVED***
	for i := 0; i < len(component); i++ ***REMOVED***
		if component[i] < ' ' ***REMOVED***
			return false
		***REMOVED***
		switch component[i] ***REMOVED***
		case '[', ']', '***REMOVED***', '***REMOVED***', '(', ')', '#', '|', '!':
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

var hexchars = [...]byte***REMOVED***
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f',
***REMOVED***

func appendHex16(dst []byte, x uint16) []byte ***REMOVED***
	return append(dst,
		hexchars[x>>12&0xF], hexchars[x>>8&0xF],
		hexchars[x>>4&0xF], hexchars[x>>0&0xF],
	)
***REMOVED***

// AppendJSONString is a convenience function that converts the provided string
// to a valid JSON string and appends it to dst.
func AppendJSONString(dst []byte, s string) []byte ***REMOVED***
	dst = append(dst, make([]byte, len(s)+2)...)
	dst = append(dst[:len(dst)-len(s)-2], '"')
	for i := 0; i < len(s); i++ ***REMOVED***
		if s[i] < ' ' ***REMOVED***
			dst = append(dst, '\\')
			switch s[i] ***REMOVED***
			case '\n':
				dst = append(dst, 'n')
			case '\r':
				dst = append(dst, 'r')
			case '\t':
				dst = append(dst, 't')
			default:
				dst = append(dst, 'u')
				dst = appendHex16(dst, uint16(s[i]))
			***REMOVED***
		***REMOVED*** else if s[i] == '>' || s[i] == '<' || s[i] == '&' ***REMOVED***
			dst = append(dst, '\\', 'u')
			dst = appendHex16(dst, uint16(s[i]))
		***REMOVED*** else if s[i] == '\\' ***REMOVED***
			dst = append(dst, '\\', '\\')
		***REMOVED*** else if s[i] == '"' ***REMOVED***
			dst = append(dst, '\\', '"')
		***REMOVED*** else if s[i] > 127 ***REMOVED***
			// read utf8 character
			r, n := utf8.DecodeRuneInString(s[i:])
			if n == 0 ***REMOVED***
				break
			***REMOVED***
			if r == utf8.RuneError && n == 1 ***REMOVED***
				dst = append(dst, `\ufffd`...)
			***REMOVED*** else if r == '\u2028' || r == '\u2029' ***REMOVED***
				dst = append(dst, `\u202`...)
				dst = append(dst, hexchars[r&0xF])
			***REMOVED*** else ***REMOVED***
				dst = append(dst, s[i:i+n]...)
			***REMOVED***
			i = i + n - 1
		***REMOVED*** else ***REMOVED***
			dst = append(dst, s[i])
		***REMOVED***
	***REMOVED***
	return append(dst, '"')
***REMOVED***

type parseContext struct ***REMOVED***
	json  string
	value Result
	pipe  string
	piped bool
	calcd bool
	lines bool
***REMOVED***

// Get searches json for the specified path.
// A path is in dot syntax, such as "name.last" or "age".
// When the value is found it's returned immediately.
//
// A path is a series of keys separated by a dot.
// A key may contain special wildcard characters '*' and '?'.
// To access an array value use the index as the key.
// To get the number of elements in an array or to access a child path, use
// the '#' character.
// The dot and wildcard character can be escaped with '\'.
//
//  ***REMOVED***
//    "name": ***REMOVED***"first": "Tom", "last": "Anderson"***REMOVED***,
//    "age":37,
//    "children": ["Sara","Alex","Jack"],
//    "friends": [
//      ***REMOVED***"first": "James", "last": "Murphy"***REMOVED***,
//      ***REMOVED***"first": "Roger", "last": "Craig"***REMOVED***
//    ]
//  ***REMOVED***
//  "name.last"          >> "Anderson"
//  "age"                >> 37
//  "children"           >> ["Sara","Alex","Jack"]
//  "children.#"         >> 3
//  "children.1"         >> "Alex"
//  "child*.2"           >> "Jack"
//  "c?ildren.0"         >> "Sara"
//  "friends.#.first"    >> ["James","Roger"]
//
// This function expects that the json is well-formed, and does not validate.
// Invalid json will not panic, but it may return back unexpected results.
// If you are consuming JSON from an unpredictable source then you may want to
// use the Valid function first.
func Get(json, path string) Result ***REMOVED***
	if len(path) > 1 ***REMOVED***
		if (path[0] == '@' && !DisableModifiers) || path[0] == '!' ***REMOVED***
			// possible modifier
			var ok bool
			var npath string
			var rjson string
			if path[0] == '@' && !DisableModifiers ***REMOVED***
				npath, rjson, ok = execModifier(json, path)
			***REMOVED*** else if path[0] == '!' ***REMOVED***
				npath, rjson, ok = execStatic(json, path)
			***REMOVED***
			if ok ***REMOVED***
				path = npath
				if len(path) > 0 && (path[0] == '|' || path[0] == '.') ***REMOVED***
					res := Get(rjson, path[1:])
					res.Index = 0
					res.Indexes = nil
					return res
				***REMOVED***
				return Parse(rjson)
			***REMOVED***
		***REMOVED***
		if path[0] == '[' || path[0] == '***REMOVED***' ***REMOVED***
			// using a subselector path
			kind := path[0]
			var ok bool
			var subs []subSelector
			subs, path, ok = parseSubSelectors(path)
			if ok ***REMOVED***
				if len(path) == 0 || (path[0] == '|' || path[0] == '.') ***REMOVED***
					var b []byte
					b = append(b, kind)
					var i int
					for _, sub := range subs ***REMOVED***
						res := Get(json, sub.path)
						if res.Exists() ***REMOVED***
							if i > 0 ***REMOVED***
								b = append(b, ',')
							***REMOVED***
							if kind == '***REMOVED***' ***REMOVED***
								if len(sub.name) > 0 ***REMOVED***
									if sub.name[0] == '"' && Valid(sub.name) ***REMOVED***
										b = append(b, sub.name...)
									***REMOVED*** else ***REMOVED***
										b = AppendJSONString(b, sub.name)
									***REMOVED***
								***REMOVED*** else ***REMOVED***
									last := nameOfLast(sub.path)
									if isSimpleName(last) ***REMOVED***
										b = AppendJSONString(b, last)
									***REMOVED*** else ***REMOVED***
										b = AppendJSONString(b, "_")
									***REMOVED***
								***REMOVED***
								b = append(b, ':')
							***REMOVED***
							var raw string
							if len(res.Raw) == 0 ***REMOVED***
								raw = res.String()
								if len(raw) == 0 ***REMOVED***
									raw = "null"
								***REMOVED***
							***REMOVED*** else ***REMOVED***
								raw = res.Raw
							***REMOVED***
							b = append(b, raw...)
							i++
						***REMOVED***
					***REMOVED***
					b = append(b, kind+2)
					var res Result
					res.Raw = string(b)
					res.Type = JSON
					if len(path) > 0 ***REMOVED***
						res = res.Get(path[1:])
					***REMOVED***
					res.Index = 0
					return res
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	var i int
	var c = &parseContext***REMOVED***json: json***REMOVED***
	if len(path) >= 2 && path[0] == '.' && path[1] == '.' ***REMOVED***
		c.lines = true
		parseArray(c, 0, path[2:])
	***REMOVED*** else ***REMOVED***
		for ; i < len(c.json); i++ ***REMOVED***
			if c.json[i] == '***REMOVED***' ***REMOVED***
				i++
				parseObject(c, i, path)
				break
			***REMOVED***
			if c.json[i] == '[' ***REMOVED***
				i++
				parseArray(c, i, path)
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if c.piped ***REMOVED***
		res := c.value.Get(c.pipe)
		res.Index = 0
		return res
	***REMOVED***
	fillIndex(json, c)
	return c.value
***REMOVED***

// GetBytes searches json for the specified path.
// If working with bytes, this method preferred over Get(string(data), path)
func GetBytes(json []byte, path string) Result ***REMOVED***
	return getBytes(json, path)
***REMOVED***

// runeit returns the rune from the the \uXXXX
func runeit(json string) rune ***REMOVED***
	n, _ := strconv.ParseUint(json[:4], 16, 64)
	return rune(n)
***REMOVED***

// unescape unescapes a string
func unescape(json string) string ***REMOVED***
	var str = make([]byte, 0, len(json))
	for i := 0; i < len(json); i++ ***REMOVED***
		switch ***REMOVED***
		default:
			str = append(str, json[i])
		case json[i] < ' ':
			return string(str)
		case json[i] == '\\':
			i++
			if i >= len(json) ***REMOVED***
				return string(str)
			***REMOVED***
			switch json[i] ***REMOVED***
			default:
				return string(str)
			case '\\':
				str = append(str, '\\')
			case '/':
				str = append(str, '/')
			case 'b':
				str = append(str, '\b')
			case 'f':
				str = append(str, '\f')
			case 'n':
				str = append(str, '\n')
			case 'r':
				str = append(str, '\r')
			case 't':
				str = append(str, '\t')
			case '"':
				str = append(str, '"')
			case 'u':
				if i+5 > len(json) ***REMOVED***
					return string(str)
				***REMOVED***
				r := runeit(json[i+1:])
				i += 5
				if utf16.IsSurrogate(r) ***REMOVED***
					// need another code
					if len(json[i:]) >= 6 && json[i] == '\\' &&
						json[i+1] == 'u' ***REMOVED***
						// we expect it to be correct so just consume it
						r = utf16.DecodeRune(r, runeit(json[i+2:]))
						i += 6
					***REMOVED***
				***REMOVED***
				// provide enough space to encode the largest utf8 possible
				str = append(str, 0, 0, 0, 0, 0, 0, 0, 0)
				n := utf8.EncodeRune(str[len(str)-8:], r)
				str = str[:len(str)-8+n]
				i-- // backtrack index by one
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return string(str)
***REMOVED***

// Less return true if a token is less than another token.
// The caseSensitive paramater is used when the tokens are Strings.
// The order when comparing two different type is:
//
//  Null < False < Number < String < True < JSON
//
func (t Result) Less(token Result, caseSensitive bool) bool ***REMOVED***
	if t.Type < token.Type ***REMOVED***
		return true
	***REMOVED***
	if t.Type > token.Type ***REMOVED***
		return false
	***REMOVED***
	if t.Type == String ***REMOVED***
		if caseSensitive ***REMOVED***
			return t.Str < token.Str
		***REMOVED***
		return stringLessInsensitive(t.Str, token.Str)
	***REMOVED***
	if t.Type == Number ***REMOVED***
		return t.Num < token.Num
	***REMOVED***
	return t.Raw < token.Raw
***REMOVED***

func stringLessInsensitive(a, b string) bool ***REMOVED***
	for i := 0; i < len(a) && i < len(b); i++ ***REMOVED***
		if a[i] >= 'A' && a[i] <= 'Z' ***REMOVED***
			if b[i] >= 'A' && b[i] <= 'Z' ***REMOVED***
				// both are uppercase, do nothing
				if a[i] < b[i] ***REMOVED***
					return true
				***REMOVED*** else if a[i] > b[i] ***REMOVED***
					return false
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				// a is uppercase, convert a to lowercase
				if a[i]+32 < b[i] ***REMOVED***
					return true
				***REMOVED*** else if a[i]+32 > b[i] ***REMOVED***
					return false
				***REMOVED***
			***REMOVED***
		***REMOVED*** else if b[i] >= 'A' && b[i] <= 'Z' ***REMOVED***
			// b is uppercase, convert b to lowercase
			if a[i] < b[i]+32 ***REMOVED***
				return true
			***REMOVED*** else if a[i] > b[i]+32 ***REMOVED***
				return false
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// neither are uppercase
			if a[i] < b[i] ***REMOVED***
				return true
			***REMOVED*** else if a[i] > b[i] ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return len(a) < len(b)
***REMOVED***

// parseAny parses the next value from a json string.
// A Result is returned when the hit param is set.
// The return values are (i int, res Result, ok bool)
func parseAny(json string, i int, hit bool) (int, Result, bool) ***REMOVED***
	var res Result
	var val string
	for ; i < len(json); i++ ***REMOVED***
		if json[i] == '***REMOVED***' || json[i] == '[' ***REMOVED***
			i, val = parseSquash(json, i)
			if hit ***REMOVED***
				res.Raw = val
				res.Type = JSON
			***REMOVED***
			var tmp parseContext
			tmp.value = res
			fillIndex(json, &tmp)
			return i, tmp.value, true
		***REMOVED***
		if json[i] <= ' ' ***REMOVED***
			continue
		***REMOVED***
		var num bool
		switch json[i] ***REMOVED***
		case '"':
			i++
			var vesc bool
			var ok bool
			i, val, vesc, ok = parseString(json, i)
			if !ok ***REMOVED***
				return i, res, false
			***REMOVED***
			if hit ***REMOVED***
				res.Type = String
				res.Raw = val
				if vesc ***REMOVED***
					res.Str = unescape(val[1 : len(val)-1])
				***REMOVED*** else ***REMOVED***
					res.Str = val[1 : len(val)-1]
				***REMOVED***
			***REMOVED***
			return i, res, true
		case 'n':
			if i+1 < len(json) && json[i+1] != 'u' ***REMOVED***
				num = true
				break
			***REMOVED***
			fallthrough
		case 't', 'f':
			vc := json[i]
			i, val = parseLiteral(json, i)
			if hit ***REMOVED***
				res.Raw = val
				switch vc ***REMOVED***
				case 't':
					res.Type = True
				case 'f':
					res.Type = False
				***REMOVED***
				return i, res, true
			***REMOVED***
		case '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'i', 'I', 'N':
			num = true
		***REMOVED***
		if num ***REMOVED***
			i, val = parseNumber(json, i)
			if hit ***REMOVED***
				res.Raw = val
				res.Type = Number
				res.Num, _ = strconv.ParseFloat(val, 64)
			***REMOVED***
			return i, res, true
		***REMOVED***

	***REMOVED***
	return i, res, false
***REMOVED***

// GetMany searches json for the multiple paths.
// The return value is a Result array where the number of items
// will be equal to the number of input paths.
func GetMany(json string, path ...string) []Result ***REMOVED***
	res := make([]Result, len(path))
	for i, path := range path ***REMOVED***
		res[i] = Get(json, path)
	***REMOVED***
	return res
***REMOVED***

// GetManyBytes searches json for the multiple paths.
// The return value is a Result array where the number of items
// will be equal to the number of input paths.
func GetManyBytes(json []byte, path ...string) []Result ***REMOVED***
	res := make([]Result, len(path))
	for i, path := range path ***REMOVED***
		res[i] = GetBytes(json, path)
	***REMOVED***
	return res
***REMOVED***

func validpayload(data []byte, i int) (outi int, ok bool) ***REMOVED***
	for ; i < len(data); i++ ***REMOVED***
		switch data[i] ***REMOVED***
		default:
			i, ok = validany(data, i)
			if !ok ***REMOVED***
				return i, false
			***REMOVED***
			for ; i < len(data); i++ ***REMOVED***
				switch data[i] ***REMOVED***
				default:
					return i, false
				case ' ', '\t', '\n', '\r':
					continue
				***REMOVED***
			***REMOVED***
			return i, true
		case ' ', '\t', '\n', '\r':
			continue
		***REMOVED***
	***REMOVED***
	return i, false
***REMOVED***
func validany(data []byte, i int) (outi int, ok bool) ***REMOVED***
	for ; i < len(data); i++ ***REMOVED***
		switch data[i] ***REMOVED***
		default:
			return i, false
		case ' ', '\t', '\n', '\r':
			continue
		case '***REMOVED***':
			return validobject(data, i+1)
		case '[':
			return validarray(data, i+1)
		case '"':
			return validstring(data, i+1)
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return validnumber(data, i+1)
		case 't':
			return validtrue(data, i+1)
		case 'f':
			return validfalse(data, i+1)
		case 'n':
			return validnull(data, i+1)
		***REMOVED***
	***REMOVED***
	return i, false
***REMOVED***
func validobject(data []byte, i int) (outi int, ok bool) ***REMOVED***
	for ; i < len(data); i++ ***REMOVED***
		switch data[i] ***REMOVED***
		default:
			return i, false
		case ' ', '\t', '\n', '\r':
			continue
		case '***REMOVED***':
			return i + 1, true
		case '"':
		key:
			if i, ok = validstring(data, i+1); !ok ***REMOVED***
				return i, false
			***REMOVED***
			if i, ok = validcolon(data, i); !ok ***REMOVED***
				return i, false
			***REMOVED***
			if i, ok = validany(data, i); !ok ***REMOVED***
				return i, false
			***REMOVED***
			if i, ok = validcomma(data, i, '***REMOVED***'); !ok ***REMOVED***
				return i, false
			***REMOVED***
			if data[i] == '***REMOVED***' ***REMOVED***
				return i + 1, true
			***REMOVED***
			i++
			for ; i < len(data); i++ ***REMOVED***
				switch data[i] ***REMOVED***
				default:
					return i, false
				case ' ', '\t', '\n', '\r':
					continue
				case '"':
					goto key
				***REMOVED***
			***REMOVED***
			return i, false
		***REMOVED***
	***REMOVED***
	return i, false
***REMOVED***
func validcolon(data []byte, i int) (outi int, ok bool) ***REMOVED***
	for ; i < len(data); i++ ***REMOVED***
		switch data[i] ***REMOVED***
		default:
			return i, false
		case ' ', '\t', '\n', '\r':
			continue
		case ':':
			return i + 1, true
		***REMOVED***
	***REMOVED***
	return i, false
***REMOVED***
func validcomma(data []byte, i int, end byte) (outi int, ok bool) ***REMOVED***
	for ; i < len(data); i++ ***REMOVED***
		switch data[i] ***REMOVED***
		default:
			return i, false
		case ' ', '\t', '\n', '\r':
			continue
		case ',':
			return i, true
		case end:
			return i, true
		***REMOVED***
	***REMOVED***
	return i, false
***REMOVED***
func validarray(data []byte, i int) (outi int, ok bool) ***REMOVED***
	for ; i < len(data); i++ ***REMOVED***
		switch data[i] ***REMOVED***
		default:
			for ; i < len(data); i++ ***REMOVED***
				if i, ok = validany(data, i); !ok ***REMOVED***
					return i, false
				***REMOVED***
				if i, ok = validcomma(data, i, ']'); !ok ***REMOVED***
					return i, false
				***REMOVED***
				if data[i] == ']' ***REMOVED***
					return i + 1, true
				***REMOVED***
			***REMOVED***
		case ' ', '\t', '\n', '\r':
			continue
		case ']':
			return i + 1, true
		***REMOVED***
	***REMOVED***
	return i, false
***REMOVED***
func validstring(data []byte, i int) (outi int, ok bool) ***REMOVED***
	for ; i < len(data); i++ ***REMOVED***
		if data[i] < ' ' ***REMOVED***
			return i, false
		***REMOVED*** else if data[i] == '\\' ***REMOVED***
			i++
			if i == len(data) ***REMOVED***
				return i, false
			***REMOVED***
			switch data[i] ***REMOVED***
			default:
				return i, false
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
			case 'u':
				for j := 0; j < 4; j++ ***REMOVED***
					i++
					if i >= len(data) ***REMOVED***
						return i, false
					***REMOVED***
					if !((data[i] >= '0' && data[i] <= '9') ||
						(data[i] >= 'a' && data[i] <= 'f') ||
						(data[i] >= 'A' && data[i] <= 'F')) ***REMOVED***
						return i, false
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED*** else if data[i] == '"' ***REMOVED***
			return i + 1, true
		***REMOVED***
	***REMOVED***
	return i, false
***REMOVED***
func validnumber(data []byte, i int) (outi int, ok bool) ***REMOVED***
	i--
	// sign
	if data[i] == '-' ***REMOVED***
		i++
		if i == len(data) ***REMOVED***
			return i, false
		***REMOVED***
		if data[i] < '0' || data[i] > '9' ***REMOVED***
			return i, false
		***REMOVED***
	***REMOVED***
	// int
	if i == len(data) ***REMOVED***
		return i, false
	***REMOVED***
	if data[i] == '0' ***REMOVED***
		i++
	***REMOVED*** else ***REMOVED***
		for ; i < len(data); i++ ***REMOVED***
			if data[i] >= '0' && data[i] <= '9' ***REMOVED***
				continue
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	// frac
	if i == len(data) ***REMOVED***
		return i, true
	***REMOVED***
	if data[i] == '.' ***REMOVED***
		i++
		if i == len(data) ***REMOVED***
			return i, false
		***REMOVED***
		if data[i] < '0' || data[i] > '9' ***REMOVED***
			return i, false
		***REMOVED***
		i++
		for ; i < len(data); i++ ***REMOVED***
			if data[i] >= '0' && data[i] <= '9' ***REMOVED***
				continue
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	// exp
	if i == len(data) ***REMOVED***
		return i, true
	***REMOVED***
	if data[i] == 'e' || data[i] == 'E' ***REMOVED***
		i++
		if i == len(data) ***REMOVED***
			return i, false
		***REMOVED***
		if data[i] == '+' || data[i] == '-' ***REMOVED***
			i++
		***REMOVED***
		if i == len(data) ***REMOVED***
			return i, false
		***REMOVED***
		if data[i] < '0' || data[i] > '9' ***REMOVED***
			return i, false
		***REMOVED***
		i++
		for ; i < len(data); i++ ***REMOVED***
			if data[i] >= '0' && data[i] <= '9' ***REMOVED***
				continue
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return i, true
***REMOVED***

func validtrue(data []byte, i int) (outi int, ok bool) ***REMOVED***
	if i+3 <= len(data) && data[i] == 'r' && data[i+1] == 'u' &&
		data[i+2] == 'e' ***REMOVED***
		return i + 3, true
	***REMOVED***
	return i, false
***REMOVED***
func validfalse(data []byte, i int) (outi int, ok bool) ***REMOVED***
	if i+4 <= len(data) && data[i] == 'a' && data[i+1] == 'l' &&
		data[i+2] == 's' && data[i+3] == 'e' ***REMOVED***
		return i + 4, true
	***REMOVED***
	return i, false
***REMOVED***
func validnull(data []byte, i int) (outi int, ok bool) ***REMOVED***
	if i+3 <= len(data) && data[i] == 'u' && data[i+1] == 'l' &&
		data[i+2] == 'l' ***REMOVED***
		return i + 3, true
	***REMOVED***
	return i, false
***REMOVED***

// Valid returns true if the input is valid json.
//
//  if !gjson.Valid(json) ***REMOVED***
//  	return errors.New("invalid json")
//  ***REMOVED***
//  value := gjson.Get(json, "name.last")
//
func Valid(json string) bool ***REMOVED***
	_, ok := validpayload(stringBytes(json), 0)
	return ok
***REMOVED***

// ValidBytes returns true if the input is valid json.
//
//  if !gjson.Valid(json) ***REMOVED***
//  	return errors.New("invalid json")
//  ***REMOVED***
//  value := gjson.Get(json, "name.last")
//
// If working with bytes, this method preferred over ValidBytes(string(data))
//
func ValidBytes(json []byte) bool ***REMOVED***
	_, ok := validpayload(json, 0)
	return ok
***REMOVED***

func parseUint(s string) (n uint64, ok bool) ***REMOVED***
	var i int
	if i == len(s) ***REMOVED***
		return 0, false
	***REMOVED***
	for ; i < len(s); i++ ***REMOVED***
		if s[i] >= '0' && s[i] <= '9' ***REMOVED***
			n = n*10 + uint64(s[i]-'0')
		***REMOVED*** else ***REMOVED***
			return 0, false
		***REMOVED***
	***REMOVED***
	return n, true
***REMOVED***

func parseInt(s string) (n int64, ok bool) ***REMOVED***
	var i int
	var sign bool
	if len(s) > 0 && s[0] == '-' ***REMOVED***
		sign = true
		i++
	***REMOVED***
	if i == len(s) ***REMOVED***
		return 0, false
	***REMOVED***
	for ; i < len(s); i++ ***REMOVED***
		if s[i] >= '0' && s[i] <= '9' ***REMOVED***
			n = n*10 + int64(s[i]-'0')
		***REMOVED*** else ***REMOVED***
			return 0, false
		***REMOVED***
	***REMOVED***
	if sign ***REMOVED***
		return n * -1, true
	***REMOVED***
	return n, true
***REMOVED***

// safeInt validates a given JSON number
// ensures it lies within the minimum and maximum representable JSON numbers
func safeInt(f float64) (n int64, ok bool) ***REMOVED***
	// https://tc39.es/ecma262/#sec-number.min_safe_integer
	// https://tc39.es/ecma262/#sec-number.max_safe_integer
	if f < -9007199254740991 || f > 9007199254740991 ***REMOVED***
		return 0, false
	***REMOVED***
	return int64(f), true
***REMOVED***

// execStatic parses the path to find a static value.
// The input expects that the path already starts with a '!'
func execStatic(json, path string) (pathOut, res string, ok bool) ***REMOVED***
	name := path[1:]
	if len(name) > 0 ***REMOVED***
		switch name[0] ***REMOVED***
		case '***REMOVED***', '[', '"', '+', '-', '0', '1', '2', '3', '4', '5', '6', '7',
			'8', '9':
			_, res = parseSquash(name, 0)
			pathOut = name[len(res):]
			return pathOut, res, true
		***REMOVED***
	***REMOVED***
	for i := 1; i < len(path); i++ ***REMOVED***
		if path[i] == '|' ***REMOVED***
			pathOut = path[i:]
			name = path[1:i]
			break
		***REMOVED***
		if path[i] == '.' ***REMOVED***
			pathOut = path[i:]
			name = path[1:i]
			break
		***REMOVED***
	***REMOVED***
	switch strings.ToLower(name) ***REMOVED***
	case "true", "false", "null", "nan", "inf":
		return pathOut, name, true
	***REMOVED***
	return pathOut, res, false
***REMOVED***

// execModifier parses the path to find a matching modifier function.
// The input expects that the path already starts with a '@'
func execModifier(json, path string) (pathOut, res string, ok bool) ***REMOVED***
	name := path[1:]
	var hasArgs bool
	for i := 1; i < len(path); i++ ***REMOVED***
		if path[i] == ':' ***REMOVED***
			pathOut = path[i+1:]
			name = path[1:i]
			hasArgs = len(pathOut) > 0
			break
		***REMOVED***
		if path[i] == '|' ***REMOVED***
			pathOut = path[i:]
			name = path[1:i]
			break
		***REMOVED***
		if path[i] == '.' ***REMOVED***
			pathOut = path[i:]
			name = path[1:i]
			break
		***REMOVED***
	***REMOVED***
	if fn, ok := modifiers[name]; ok ***REMOVED***
		var args string
		if hasArgs ***REMOVED***
			var parsedArgs bool
			switch pathOut[0] ***REMOVED***
			case '***REMOVED***', '[', '"':
				res := Parse(pathOut)
				if res.Exists() ***REMOVED***
					args = squash(pathOut)
					pathOut = pathOut[len(args):]
					parsedArgs = true
				***REMOVED***
			***REMOVED***
			if !parsedArgs ***REMOVED***
				idx := strings.IndexByte(pathOut, '|')
				if idx == -1 ***REMOVED***
					args = pathOut
					pathOut = ""
				***REMOVED*** else ***REMOVED***
					args = pathOut[:idx]
					pathOut = pathOut[idx:]
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return pathOut, fn(json, args), true
	***REMOVED***
	return pathOut, res, false
***REMOVED***

// unwrap removes the '[]' or '***REMOVED******REMOVED***' characters around json
func unwrap(json string) string ***REMOVED***
	json = trim(json)
	if len(json) >= 2 && (json[0] == '[' || json[0] == '***REMOVED***') ***REMOVED***
		json = json[1 : len(json)-1]
	***REMOVED***
	return json
***REMOVED***

// DisableModifiers will disable the modifier syntax
var DisableModifiers = false

var modifiers = map[string]func(json, arg string) string***REMOVED***
	"pretty":  modPretty,
	"ugly":    modUgly,
	"reverse": modReverse,
	"this":    modThis,
	"flatten": modFlatten,
	"join":    modJoin,
	"valid":   modValid,
	"keys":    modKeys,
	"values":  modValues,
	"tostr":   modToStr,
	"fromstr": modFromStr,
	"group":   modGroup,
***REMOVED***

// AddModifier binds a custom modifier command to the GJSON syntax.
// This operation is not thread safe and should be executed prior to
// using all other gjson function.
func AddModifier(name string, fn func(json, arg string) string) ***REMOVED***
	modifiers[name] = fn
***REMOVED***

// ModifierExists returns true when the specified modifier exists.
func ModifierExists(name string, fn func(json, arg string) string) bool ***REMOVED***
	_, ok := modifiers[name]
	return ok
***REMOVED***

// cleanWS remove any non-whitespace from string
func cleanWS(s string) string ***REMOVED***
	for i := 0; i < len(s); i++ ***REMOVED***
		switch s[i] ***REMOVED***
		case ' ', '\t', '\n', '\r':
			continue
		default:
			var s2 []byte
			for i := 0; i < len(s); i++ ***REMOVED***
				switch s[i] ***REMOVED***
				case ' ', '\t', '\n', '\r':
					s2 = append(s2, s[i])
				***REMOVED***
			***REMOVED***
			return string(s2)
		***REMOVED***
	***REMOVED***
	return s
***REMOVED***

// @pretty modifier makes the json look nice.
func modPretty(json, arg string) string ***REMOVED***
	if len(arg) > 0 ***REMOVED***
		opts := *pretty.DefaultOptions
		Parse(arg).ForEach(func(key, value Result) bool ***REMOVED***
			switch key.String() ***REMOVED***
			case "sortKeys":
				opts.SortKeys = value.Bool()
			case "indent":
				opts.Indent = cleanWS(value.String())
			case "prefix":
				opts.Prefix = cleanWS(value.String())
			case "width":
				opts.Width = int(value.Int())
			***REMOVED***
			return true
		***REMOVED***)
		return bytesString(pretty.PrettyOptions(stringBytes(json), &opts))
	***REMOVED***
	return bytesString(pretty.Pretty(stringBytes(json)))
***REMOVED***

// @this returns the current element. Can be used to retrieve the root element.
func modThis(json, arg string) string ***REMOVED***
	return json
***REMOVED***

// @ugly modifier removes all whitespace.
func modUgly(json, arg string) string ***REMOVED***
	return bytesString(pretty.Ugly(stringBytes(json)))
***REMOVED***

// @reverse reverses array elements or root object members.
func modReverse(json, arg string) string ***REMOVED***
	res := Parse(json)
	if res.IsArray() ***REMOVED***
		var values []Result
		res.ForEach(func(_, value Result) bool ***REMOVED***
			values = append(values, value)
			return true
		***REMOVED***)
		out := make([]byte, 0, len(json))
		out = append(out, '[')
		for i, j := len(values)-1, 0; i >= 0; i, j = i-1, j+1 ***REMOVED***
			if j > 0 ***REMOVED***
				out = append(out, ',')
			***REMOVED***
			out = append(out, values[i].Raw...)
		***REMOVED***
		out = append(out, ']')
		return bytesString(out)
	***REMOVED***
	if res.IsObject() ***REMOVED***
		var keyValues []Result
		res.ForEach(func(key, value Result) bool ***REMOVED***
			keyValues = append(keyValues, key, value)
			return true
		***REMOVED***)
		out := make([]byte, 0, len(json))
		out = append(out, '***REMOVED***')
		for i, j := len(keyValues)-2, 0; i >= 0; i, j = i-2, j+1 ***REMOVED***
			if j > 0 ***REMOVED***
				out = append(out, ',')
			***REMOVED***
			out = append(out, keyValues[i+0].Raw...)
			out = append(out, ':')
			out = append(out, keyValues[i+1].Raw...)
		***REMOVED***
		out = append(out, '***REMOVED***')
		return bytesString(out)
	***REMOVED***
	return json
***REMOVED***

// @flatten an array with child arrays.
//   [1,[2],[3,4],[5,[6,7]]] -> [1,2,3,4,5,[6,7]]
// The ***REMOVED***"deep":true***REMOVED*** arg can be provide for deep flattening.
//   [1,[2],[3,4],[5,[6,7]]] -> [1,2,3,4,5,6,7]
// The original json is returned when the json is not an array.
func modFlatten(json, arg string) string ***REMOVED***
	res := Parse(json)
	if !res.IsArray() ***REMOVED***
		return json
	***REMOVED***
	var deep bool
	if arg != "" ***REMOVED***
		Parse(arg).ForEach(func(key, value Result) bool ***REMOVED***
			if key.String() == "deep" ***REMOVED***
				deep = value.Bool()
			***REMOVED***
			return true
		***REMOVED***)
	***REMOVED***
	var out []byte
	out = append(out, '[')
	var idx int
	res.ForEach(func(_, value Result) bool ***REMOVED***
		var raw string
		if value.IsArray() ***REMOVED***
			if deep ***REMOVED***
				raw = unwrap(modFlatten(value.Raw, arg))
			***REMOVED*** else ***REMOVED***
				raw = unwrap(value.Raw)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			raw = value.Raw
		***REMOVED***
		raw = strings.TrimSpace(raw)
		if len(raw) > 0 ***REMOVED***
			if idx > 0 ***REMOVED***
				out = append(out, ',')
			***REMOVED***
			out = append(out, raw...)
			idx++
		***REMOVED***
		return true
	***REMOVED***)
	out = append(out, ']')
	return bytesString(out)
***REMOVED***

// @keys extracts the keys from an object.
//  ***REMOVED***"first":"Tom","last":"Smith"***REMOVED*** -> ["first","last"]
func modKeys(json, arg string) string ***REMOVED***
	v := Parse(json)
	if !v.Exists() ***REMOVED***
		return "[]"
	***REMOVED***
	obj := v.IsObject()
	var out strings.Builder
	out.WriteByte('[')
	var i int
	v.ForEach(func(key, _ Result) bool ***REMOVED***
		if i > 0 ***REMOVED***
			out.WriteByte(',')
		***REMOVED***
		if obj ***REMOVED***
			out.WriteString(key.Raw)
		***REMOVED*** else ***REMOVED***
			out.WriteString("null")
		***REMOVED***
		i++
		return true
	***REMOVED***)
	out.WriteByte(']')
	return out.String()
***REMOVED***

// @values extracts the values from an object.
//   ***REMOVED***"first":"Tom","last":"Smith"***REMOVED*** -> ["Tom","Smith"]
func modValues(json, arg string) string ***REMOVED***
	v := Parse(json)
	if !v.Exists() ***REMOVED***
		return "[]"
	***REMOVED***
	if v.IsArray() ***REMOVED***
		return json
	***REMOVED***
	var out strings.Builder
	out.WriteByte('[')
	var i int
	v.ForEach(func(_, value Result) bool ***REMOVED***
		if i > 0 ***REMOVED***
			out.WriteByte(',')
		***REMOVED***
		out.WriteString(value.Raw)
		i++
		return true
	***REMOVED***)
	out.WriteByte(']')
	return out.String()
***REMOVED***

// @join multiple objects into a single object.
//   [***REMOVED***"first":"Tom"***REMOVED***,***REMOVED***"last":"Smith"***REMOVED***] -> ***REMOVED***"first","Tom","last":"Smith"***REMOVED***
// The arg can be "true" to specify that duplicate keys should be preserved.
//   [***REMOVED***"first":"Tom","age":37***REMOVED***,***REMOVED***"age":41***REMOVED***] -> ***REMOVED***"first","Tom","age":37,"age":41***REMOVED***
// Without preserved keys:
//   [***REMOVED***"first":"Tom","age":37***REMOVED***,***REMOVED***"age":41***REMOVED***] -> ***REMOVED***"first","Tom","age":41***REMOVED***
// The original json is returned when the json is not an object.
func modJoin(json, arg string) string ***REMOVED***
	res := Parse(json)
	if !res.IsArray() ***REMOVED***
		return json
	***REMOVED***
	var preserve bool
	if arg != "" ***REMOVED***
		Parse(arg).ForEach(func(key, value Result) bool ***REMOVED***
			if key.String() == "preserve" ***REMOVED***
				preserve = value.Bool()
			***REMOVED***
			return true
		***REMOVED***)
	***REMOVED***
	var out []byte
	out = append(out, '***REMOVED***')
	if preserve ***REMOVED***
		// Preserve duplicate keys.
		var idx int
		res.ForEach(func(_, value Result) bool ***REMOVED***
			if !value.IsObject() ***REMOVED***
				return true
			***REMOVED***
			if idx > 0 ***REMOVED***
				out = append(out, ',')
			***REMOVED***
			out = append(out, unwrap(value.Raw)...)
			idx++
			return true
		***REMOVED***)
	***REMOVED*** else ***REMOVED***
		// Deduplicate keys and generate an object with stable ordering.
		var keys []Result
		kvals := make(map[string]Result)
		res.ForEach(func(_, value Result) bool ***REMOVED***
			if !value.IsObject() ***REMOVED***
				return true
			***REMOVED***
			value.ForEach(func(key, value Result) bool ***REMOVED***
				k := key.String()
				if _, ok := kvals[k]; !ok ***REMOVED***
					keys = append(keys, key)
				***REMOVED***
				kvals[k] = value
				return true
			***REMOVED***)
			return true
		***REMOVED***)
		for i := 0; i < len(keys); i++ ***REMOVED***
			if i > 0 ***REMOVED***
				out = append(out, ',')
			***REMOVED***
			out = append(out, keys[i].Raw...)
			out = append(out, ':')
			out = append(out, kvals[keys[i].String()].Raw...)
		***REMOVED***
	***REMOVED***
	out = append(out, '***REMOVED***')
	return bytesString(out)
***REMOVED***

// @valid ensures that the json is valid before moving on. An empty string is
// returned when the json is not valid, otherwise it returns the original json.
func modValid(json, arg string) string ***REMOVED***
	if !Valid(json) ***REMOVED***
		return ""
	***REMOVED***
	return json
***REMOVED***

// @fromstr converts a string to json
//   "***REMOVED***\"id\":1023,\"name\":\"alert\"***REMOVED***" -> ***REMOVED***"id":1023,"name":"alert"***REMOVED***
func modFromStr(json, arg string) string ***REMOVED***
	if !Valid(json) ***REMOVED***
		return ""
	***REMOVED***
	return Parse(json).String()
***REMOVED***

// @tostr converts a string to json
//   ***REMOVED***"id":1023,"name":"alert"***REMOVED*** -> "***REMOVED***\"id\":1023,\"name\":\"alert\"***REMOVED***"
func modToStr(str, arg string) string ***REMOVED***
	return string(AppendJSONString(nil, str))
***REMOVED***

func modGroup(json, arg string) string ***REMOVED***
	res := Parse(json)
	if !res.IsObject() ***REMOVED***
		return ""
	***REMOVED***
	var all [][]byte
	res.ForEach(func(key, value Result) bool ***REMOVED***
		if !value.IsArray() ***REMOVED***
			return true
		***REMOVED***
		var idx int
		value.ForEach(func(_, value Result) bool ***REMOVED***
			if idx == len(all) ***REMOVED***
				all = append(all, []byte***REMOVED******REMOVED***)
			***REMOVED***
			all[idx] = append(all[idx], ("," + key.Raw + ":" + value.Raw)...)
			idx++
			return true
		***REMOVED***)
		return true
	***REMOVED***)
	var data []byte
	data = append(data, '[')
	for i, item := range all ***REMOVED***
		if i > 0 ***REMOVED***
			data = append(data, ',')
		***REMOVED***
		data = append(data, '***REMOVED***')
		data = append(data, item[1:]...)
		data = append(data, '***REMOVED***')
	***REMOVED***
	data = append(data, ']')
	return string(data)
***REMOVED***

// stringHeader instead of reflect.StringHeader
type stringHeader struct ***REMOVED***
	data unsafe.Pointer
	len  int
***REMOVED***

// sliceHeader instead of reflect.SliceHeader
type sliceHeader struct ***REMOVED***
	data unsafe.Pointer
	len  int
	cap  int
***REMOVED***

// getBytes casts the input json bytes to a string and safely returns the
// results as uniquely allocated data. This operation is intended to minimize
// copies and allocations for the large json string->[]byte.
func getBytes(json []byte, path string) Result ***REMOVED***
	var result Result
	if json != nil ***REMOVED***
		// unsafe cast to string
		result = Get(*(*string)(unsafe.Pointer(&json)), path)
		// safely get the string headers
		rawhi := *(*stringHeader)(unsafe.Pointer(&result.Raw))
		strhi := *(*stringHeader)(unsafe.Pointer(&result.Str))
		// create byte slice headers
		rawh := sliceHeader***REMOVED***data: rawhi.data, len: rawhi.len, cap: rawhi.len***REMOVED***
		strh := sliceHeader***REMOVED***data: strhi.data, len: strhi.len, cap: rawhi.len***REMOVED***
		if strh.data == nil ***REMOVED***
			// str is nil
			if rawh.data == nil ***REMOVED***
				// raw is nil
				result.Raw = ""
			***REMOVED*** else ***REMOVED***
				// raw has data, safely copy the slice header to a string
				result.Raw = string(*(*[]byte)(unsafe.Pointer(&rawh)))
			***REMOVED***
			result.Str = ""
		***REMOVED*** else if rawh.data == nil ***REMOVED***
			// raw is nil
			result.Raw = ""
			// str has data, safely copy the slice header to a string
			result.Str = string(*(*[]byte)(unsafe.Pointer(&strh)))
		***REMOVED*** else if uintptr(strh.data) >= uintptr(rawh.data) &&
			uintptr(strh.data)+uintptr(strh.len) <=
				uintptr(rawh.data)+uintptr(rawh.len) ***REMOVED***
			// Str is a substring of Raw.
			start := uintptr(strh.data) - uintptr(rawh.data)
			// safely copy the raw slice header
			result.Raw = string(*(*[]byte)(unsafe.Pointer(&rawh)))
			// substring the raw
			result.Str = result.Raw[start : start+uintptr(strh.len)]
		***REMOVED*** else ***REMOVED***
			// safely copy both the raw and str slice headers to strings
			result.Raw = string(*(*[]byte)(unsafe.Pointer(&rawh)))
			result.Str = string(*(*[]byte)(unsafe.Pointer(&strh)))
		***REMOVED***
	***REMOVED***
	return result
***REMOVED***

// fillIndex finds the position of Raw data and assigns it to the Index field
// of the resulting value. If the position cannot be found then Index zero is
// used instead.
func fillIndex(json string, c *parseContext) ***REMOVED***
	if len(c.value.Raw) > 0 && !c.calcd ***REMOVED***
		jhdr := *(*stringHeader)(unsafe.Pointer(&json))
		rhdr := *(*stringHeader)(unsafe.Pointer(&(c.value.Raw)))
		c.value.Index = int(uintptr(rhdr.data) - uintptr(jhdr.data))
		if c.value.Index < 0 || c.value.Index >= len(json) ***REMOVED***
			c.value.Index = 0
		***REMOVED***
	***REMOVED***
***REMOVED***

func stringBytes(s string) []byte ***REMOVED***
	return *(*[]byte)(unsafe.Pointer(&sliceHeader***REMOVED***
		data: (*stringHeader)(unsafe.Pointer(&s)).data,
		len:  len(s),
		cap:  len(s),
	***REMOVED***))
***REMOVED***

func bytesString(b []byte) string ***REMOVED***
	return *(*string)(unsafe.Pointer(&b))
***REMOVED***

func revSquash(json string) string ***REMOVED***
	// reverse squash
	// expects that the tail character is a ']' or '***REMOVED***' or ')' or '"'
	// squash the value, ignoring all nested arrays and objects.
	i := len(json) - 1
	var depth int
	if json[i] != '"' ***REMOVED***
		depth++
	***REMOVED***
	if json[i] == '***REMOVED***' || json[i] == ']' || json[i] == ')' ***REMOVED***
		i--
	***REMOVED***
	for ; i >= 0; i-- ***REMOVED***
		switch json[i] ***REMOVED***
		case '"':
			i--
			for ; i >= 0; i-- ***REMOVED***
				if json[i] == '"' ***REMOVED***
					esc := 0
					for i > 0 && json[i-1] == '\\' ***REMOVED***
						i--
						esc++
					***REMOVED***
					if esc%2 == 1 ***REMOVED***
						continue
					***REMOVED***
					i += esc
					break
				***REMOVED***
			***REMOVED***
			if depth == 0 ***REMOVED***
				if i < 0 ***REMOVED***
					i = 0
				***REMOVED***
				return json[i:]
			***REMOVED***
		case '***REMOVED***', ']', ')':
			depth++
		case '***REMOVED***', '[', '(':
			depth--
			if depth == 0 ***REMOVED***
				return json[i:]
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return json
***REMOVED***

// Paths returns the original GJSON paths for a Result where the Result came
// from a simple query path that returns an array, like:
//
//    gjson.Get(json, "friends.#.first")
//
// The returned value will be in the form of a JSON array:
//
//    ["friends.0.first","friends.1.first","friends.2.first"]
//
// The param 'json' must be the original JSON used when calling Get.
//
// Returns an empty string if the paths cannot be determined, which can happen
// when the Result came from a path that contained a multipath, modifier,
// or a nested query.
func (t Result) Paths(json string) []string ***REMOVED***
	if t.Indexes == nil ***REMOVED***
		return nil
	***REMOVED***
	paths := make([]string, 0, len(t.Indexes))
	t.ForEach(func(_, value Result) bool ***REMOVED***
		paths = append(paths, value.Path(json))
		return true
	***REMOVED***)
	if len(paths) != len(t.Indexes) ***REMOVED***
		return nil
	***REMOVED***
	return paths
***REMOVED***

// Path returns the original GJSON path for a Result where the Result came
// from a simple path that returns a single value, like:
//
//    gjson.Get(json, "friends.#(last=Murphy)")
//
// The returned value will be in the form of a JSON string:
//
//    "friends.0"
//
// The param 'json' must be the original JSON used when calling Get.
//
// Returns an empty string if the paths cannot be determined, which can happen
// when the Result came from a path that contained a multipath, modifier,
// or a nested query.
func (t Result) Path(json string) string ***REMOVED***
	var path []byte
	var comps []string // raw components
	i := t.Index - 1
	if t.Index+len(t.Raw) > len(json) ***REMOVED***
		// JSON cannot safely contain Result.
		goto fail
	***REMOVED***
	if !strings.HasPrefix(json[t.Index:], t.Raw) ***REMOVED***
		// Result is not at the JSON index as exepcted.
		goto fail
	***REMOVED***
	for ; i >= 0; i-- ***REMOVED***
		if json[i] <= ' ' ***REMOVED***
			continue
		***REMOVED***
		if json[i] == ':' ***REMOVED***
			// inside of object, get the key
			for ; i >= 0; i-- ***REMOVED***
				if json[i] != '"' ***REMOVED***
					continue
				***REMOVED***
				break
			***REMOVED***
			raw := revSquash(json[:i+1])
			i = i - len(raw)
			comps = append(comps, raw)
			// key gotten, now squash the rest
			raw = revSquash(json[:i+1])
			i = i - len(raw)
			i++ // increment the index for next loop step
		***REMOVED*** else if json[i] == '***REMOVED***' ***REMOVED***
			// Encountered an open object. The original result was probably an
			// object key.
			goto fail
		***REMOVED*** else if json[i] == ',' || json[i] == '[' ***REMOVED***
			// inside of an array, count the position
			var arrIdx int
			if json[i] == ',' ***REMOVED***
				arrIdx++
				i--
			***REMOVED***
			for ; i >= 0; i-- ***REMOVED***
				if json[i] == ':' ***REMOVED***
					// Encountered an unexpected colon. The original result was
					// probably an object key.
					goto fail
				***REMOVED*** else if json[i] == ',' ***REMOVED***
					arrIdx++
				***REMOVED*** else if json[i] == '[' ***REMOVED***
					comps = append(comps, strconv.Itoa(arrIdx))
					break
				***REMOVED*** else if json[i] == ']' || json[i] == '***REMOVED***' || json[i] == '"' ***REMOVED***
					raw := revSquash(json[:i+1])
					i = i - len(raw) + 1
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if len(comps) == 0 ***REMOVED***
		if DisableModifiers ***REMOVED***
			goto fail
		***REMOVED***
		return "@this"
	***REMOVED***
	for i := len(comps) - 1; i >= 0; i-- ***REMOVED***
		rcomp := Parse(comps[i])
		if !rcomp.Exists() ***REMOVED***
			goto fail
		***REMOVED***
		comp := escapeComp(rcomp.String())
		path = append(path, '.')
		path = append(path, comp...)
	***REMOVED***
	if len(path) > 0 ***REMOVED***
		path = path[1:]
	***REMOVED***
	return string(path)
fail:
	return ""
***REMOVED***

// isSafePathKeyChar returns true if the input character is safe for not
// needing escaping.
func isSafePathKeyChar(c byte) bool ***REMOVED***
	return c <= ' ' || c > '~' || c == '_' || c == '-' || c == ':' ||
		(c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9')
***REMOVED***

// escapeComp escaped a path compontent, making it safe for generating a
// path for later use.
func escapeComp(comp string) string ***REMOVED***
	for i := 0; i < len(comp); i++ ***REMOVED***
		if !isSafePathKeyChar(comp[i]) ***REMOVED***
			ncomp := []byte(comp[:i])
			for ; i < len(comp); i++ ***REMOVED***
				if !isSafePathKeyChar(comp[i]) ***REMOVED***
					ncomp = append(ncomp, '\\')
				***REMOVED***
				ncomp = append(ncomp, comp[i])
			***REMOVED***
			return string(ncomp)
		***REMOVED***
	***REMOVED***
	return comp
***REMOVED***
