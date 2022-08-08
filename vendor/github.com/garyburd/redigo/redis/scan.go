// Copyright 2012 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package redis

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

func ensureLen(d reflect.Value, n int) ***REMOVED***
	if n > d.Cap() ***REMOVED***
		d.Set(reflect.MakeSlice(d.Type(), n, n))
	***REMOVED*** else ***REMOVED***
		d.SetLen(n)
	***REMOVED***
***REMOVED***

func cannotConvert(d reflect.Value, s interface***REMOVED******REMOVED***) error ***REMOVED***
	var sname string
	switch s.(type) ***REMOVED***
	case string:
		sname = "Redis simple string"
	case Error:
		sname = "Redis error"
	case int64:
		sname = "Redis integer"
	case []byte:
		sname = "Redis bulk string"
	case []interface***REMOVED******REMOVED***:
		sname = "Redis array"
	default:
		sname = reflect.TypeOf(s).String()
	***REMOVED***
	return fmt.Errorf("cannot convert from %s to %s", sname, d.Type())
***REMOVED***

func convertAssignBulkString(d reflect.Value, s []byte) (err error) ***REMOVED***
	switch d.Type().Kind() ***REMOVED***
	case reflect.Float32, reflect.Float64:
		var x float64
		x, err = strconv.ParseFloat(string(s), d.Type().Bits())
		d.SetFloat(x)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var x int64
		x, err = strconv.ParseInt(string(s), 10, d.Type().Bits())
		d.SetInt(x)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var x uint64
		x, err = strconv.ParseUint(string(s), 10, d.Type().Bits())
		d.SetUint(x)
	case reflect.Bool:
		var x bool
		x, err = strconv.ParseBool(string(s))
		d.SetBool(x)
	case reflect.String:
		d.SetString(string(s))
	case reflect.Slice:
		if d.Type().Elem().Kind() != reflect.Uint8 ***REMOVED***
			err = cannotConvert(d, s)
		***REMOVED*** else ***REMOVED***
			d.SetBytes(s)
		***REMOVED***
	default:
		err = cannotConvert(d, s)
	***REMOVED***
	return
***REMOVED***

func convertAssignInt(d reflect.Value, s int64) (err error) ***REMOVED***
	switch d.Type().Kind() ***REMOVED***
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		d.SetInt(s)
		if d.Int() != s ***REMOVED***
			err = strconv.ErrRange
			d.SetInt(0)
		***REMOVED***
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if s < 0 ***REMOVED***
			err = strconv.ErrRange
		***REMOVED*** else ***REMOVED***
			x := uint64(s)
			d.SetUint(x)
			if d.Uint() != x ***REMOVED***
				err = strconv.ErrRange
				d.SetUint(0)
			***REMOVED***
		***REMOVED***
	case reflect.Bool:
		d.SetBool(s != 0)
	default:
		err = cannotConvert(d, s)
	***REMOVED***
	return
***REMOVED***

func convertAssignValue(d reflect.Value, s interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	if d.Kind() != reflect.Ptr ***REMOVED***
		if d.CanAddr() ***REMOVED***
			d2 := d.Addr()
			if d2.CanInterface() ***REMOVED***
				if scanner, ok := d2.Interface().(Scanner); ok ***REMOVED***
					return scanner.RedisScan(s)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED*** else if d.CanInterface() ***REMOVED***
		// Already a reflect.Ptr
		if d.IsNil() ***REMOVED***
			d.Set(reflect.New(d.Type().Elem()))
		***REMOVED***
		if scanner, ok := d.Interface().(Scanner); ok ***REMOVED***
			return scanner.RedisScan(s)
		***REMOVED***
	***REMOVED***

	switch s := s.(type) ***REMOVED***
	case []byte:
		err = convertAssignBulkString(d, s)
	case int64:
		err = convertAssignInt(d, s)
	default:
		err = cannotConvert(d, s)
	***REMOVED***
	return err
***REMOVED***

func convertAssignArray(d reflect.Value, s []interface***REMOVED******REMOVED***) error ***REMOVED***
	if d.Type().Kind() != reflect.Slice ***REMOVED***
		return cannotConvert(d, s)
	***REMOVED***
	ensureLen(d, len(s))
	for i := 0; i < len(s); i++ ***REMOVED***
		if err := convertAssignValue(d.Index(i), s[i]); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func convertAssign(d interface***REMOVED******REMOVED***, s interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	if scanner, ok := d.(Scanner); ok ***REMOVED***
		return scanner.RedisScan(s)
	***REMOVED***

	// Handle the most common destination types using type switches and
	// fall back to reflection for all other types.
	switch s := s.(type) ***REMOVED***
	case nil:
		// ignore
	case []byte:
		switch d := d.(type) ***REMOVED***
		case *string:
			*d = string(s)
		case *int:
			*d, err = strconv.Atoi(string(s))
		case *bool:
			*d, err = strconv.ParseBool(string(s))
		case *[]byte:
			*d = s
		case *interface***REMOVED******REMOVED***:
			*d = s
		case nil:
			// skip value
		default:
			if d := reflect.ValueOf(d); d.Type().Kind() != reflect.Ptr ***REMOVED***
				err = cannotConvert(d, s)
			***REMOVED*** else ***REMOVED***
				err = convertAssignBulkString(d.Elem(), s)
			***REMOVED***
		***REMOVED***
	case int64:
		switch d := d.(type) ***REMOVED***
		case *int:
			x := int(s)
			if int64(x) != s ***REMOVED***
				err = strconv.ErrRange
				x = 0
			***REMOVED***
			*d = x
		case *bool:
			*d = s != 0
		case *interface***REMOVED******REMOVED***:
			*d = s
		case nil:
			// skip value
		default:
			if d := reflect.ValueOf(d); d.Type().Kind() != reflect.Ptr ***REMOVED***
				err = cannotConvert(d, s)
			***REMOVED*** else ***REMOVED***
				err = convertAssignInt(d.Elem(), s)
			***REMOVED***
		***REMOVED***
	case string:
		switch d := d.(type) ***REMOVED***
		case *string:
			*d = s
		case *interface***REMOVED******REMOVED***:
			*d = s
		case nil:
			// skip value
		default:
			err = cannotConvert(reflect.ValueOf(d), s)
		***REMOVED***
	case []interface***REMOVED******REMOVED***:
		switch d := d.(type) ***REMOVED***
		case *[]interface***REMOVED******REMOVED***:
			*d = s
		case *interface***REMOVED******REMOVED***:
			*d = s
		case nil:
			// skip value
		default:
			if d := reflect.ValueOf(d); d.Type().Kind() != reflect.Ptr ***REMOVED***
				err = cannotConvert(d, s)
			***REMOVED*** else ***REMOVED***
				err = convertAssignArray(d.Elem(), s)
			***REMOVED***
		***REMOVED***
	case Error:
		err = s
	default:
		err = cannotConvert(reflect.ValueOf(d), s)
	***REMOVED***
	return
***REMOVED***

// Scan copies from src to the values pointed at by dest.
//
// Scan uses RedisScan if available otherwise:
//
// The values pointed at by dest must be an integer, float, boolean, string,
// []byte, interface***REMOVED******REMOVED*** or slices of these types. Scan uses the standard strconv
// package to convert bulk strings to numeric and boolean types.
//
// If a dest value is nil, then the corresponding src value is skipped.
//
// If a src element is nil, then the corresponding dest value is not modified.
//
// To enable easy use of Scan in a loop, Scan returns the slice of src
// following the copied values.
func Scan(src []interface***REMOVED******REMOVED***, dest ...interface***REMOVED******REMOVED***) ([]interface***REMOVED******REMOVED***, error) ***REMOVED***
	if len(src) < len(dest) ***REMOVED***
		return nil, errors.New("redigo.Scan: array short")
	***REMOVED***
	var err error
	for i, d := range dest ***REMOVED***
		err = convertAssign(d, src[i])
		if err != nil ***REMOVED***
			err = fmt.Errorf("redigo.Scan: cannot assign to dest %d: %v", i, err)
			break
		***REMOVED***
	***REMOVED***
	return src[len(dest):], err
***REMOVED***

type fieldSpec struct ***REMOVED***
	name      string
	index     []int
	omitEmpty bool
***REMOVED***

type structSpec struct ***REMOVED***
	m map[string]*fieldSpec
	l []*fieldSpec
***REMOVED***

func (ss *structSpec) fieldSpec(name []byte) *fieldSpec ***REMOVED***
	return ss.m[string(name)]
***REMOVED***

func compileStructSpec(t reflect.Type, depth map[string]int, index []int, ss *structSpec) ***REMOVED***
	for i := 0; i < t.NumField(); i++ ***REMOVED***
		f := t.Field(i)
		switch ***REMOVED***
		case f.PkgPath != "" && !f.Anonymous:
			// Ignore unexported fields.
		case f.Anonymous:
			// TODO: Handle pointers. Requires change to decoder and
			// protection against infinite recursion.
			if f.Type.Kind() == reflect.Struct ***REMOVED***
				compileStructSpec(f.Type, depth, append(index, i), ss)
			***REMOVED***
		default:
			fs := &fieldSpec***REMOVED***name: f.Name***REMOVED***
			tag := f.Tag.Get("redis")
			p := strings.Split(tag, ",")
			if len(p) > 0 ***REMOVED***
				if p[0] == "-" ***REMOVED***
					continue
				***REMOVED***
				if len(p[0]) > 0 ***REMOVED***
					fs.name = p[0]
				***REMOVED***
				for _, s := range p[1:] ***REMOVED***
					switch s ***REMOVED***
					case "omitempty":
						fs.omitEmpty = true
					default:
						panic(fmt.Errorf("redigo: unknown field tag %s for type %s", s, t.Name()))
					***REMOVED***
				***REMOVED***
			***REMOVED***
			d, found := depth[fs.name]
			if !found ***REMOVED***
				d = 1 << 30
			***REMOVED***
			switch ***REMOVED***
			case len(index) == d:
				// At same depth, remove from result.
				delete(ss.m, fs.name)
				j := 0
				for i := 0; i < len(ss.l); i++ ***REMOVED***
					if fs.name != ss.l[i].name ***REMOVED***
						ss.l[j] = ss.l[i]
						j += 1
					***REMOVED***
				***REMOVED***
				ss.l = ss.l[:j]
			case len(index) < d:
				fs.index = make([]int, len(index)+1)
				copy(fs.index, index)
				fs.index[len(index)] = i
				depth[fs.name] = len(index)
				ss.m[fs.name] = fs
				ss.l = append(ss.l, fs)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

var (
	structSpecMutex  sync.RWMutex
	structSpecCache  = make(map[reflect.Type]*structSpec)
	defaultFieldSpec = &fieldSpec***REMOVED******REMOVED***
)

func structSpecForType(t reflect.Type) *structSpec ***REMOVED***

	structSpecMutex.RLock()
	ss, found := structSpecCache[t]
	structSpecMutex.RUnlock()
	if found ***REMOVED***
		return ss
	***REMOVED***

	structSpecMutex.Lock()
	defer structSpecMutex.Unlock()
	ss, found = structSpecCache[t]
	if found ***REMOVED***
		return ss
	***REMOVED***

	ss = &structSpec***REMOVED***m: make(map[string]*fieldSpec)***REMOVED***
	compileStructSpec(t, make(map[string]int), nil, ss)
	structSpecCache[t] = ss
	return ss
***REMOVED***

var errScanStructValue = errors.New("redigo.ScanStruct: value must be non-nil pointer to a struct")

// ScanStruct scans alternating names and values from src to a struct. The
// HGETALL and CONFIG GET commands return replies in this format.
//
// ScanStruct uses exported field names to match values in the response. Use
// 'redis' field tag to override the name:
//
//      Field int `redis:"myName"`
//
// Fields with the tag redis:"-" are ignored.
//
// Each field uses RedisScan if available otherwise:
// Integer, float, boolean, string and []byte fields are supported. Scan uses the
// standard strconv package to convert bulk string values to numeric and
// boolean types.
//
// If a src element is nil, then the corresponding field is not modified.
func ScanStruct(src []interface***REMOVED******REMOVED***, dest interface***REMOVED******REMOVED***) error ***REMOVED***
	d := reflect.ValueOf(dest)
	if d.Kind() != reflect.Ptr || d.IsNil() ***REMOVED***
		return errScanStructValue
	***REMOVED***
	d = d.Elem()
	if d.Kind() != reflect.Struct ***REMOVED***
		return errScanStructValue
	***REMOVED***
	ss := structSpecForType(d.Type())

	if len(src)%2 != 0 ***REMOVED***
		return errors.New("redigo.ScanStruct: number of values not a multiple of 2")
	***REMOVED***

	for i := 0; i < len(src); i += 2 ***REMOVED***
		s := src[i+1]
		if s == nil ***REMOVED***
			continue
		***REMOVED***
		name, ok := src[i].([]byte)
		if !ok ***REMOVED***
			return fmt.Errorf("redigo.ScanStruct: key %d not a bulk string value", i)
		***REMOVED***
		fs := ss.fieldSpec(name)
		if fs == nil ***REMOVED***
			continue
		***REMOVED***
		if err := convertAssignValue(d.FieldByIndex(fs.index), s); err != nil ***REMOVED***
			return fmt.Errorf("redigo.ScanStruct: cannot assign field %s: %v", fs.name, err)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

var (
	errScanSliceValue = errors.New("redigo.ScanSlice: dest must be non-nil pointer to a struct")
)

// ScanSlice scans src to the slice pointed to by dest. The elements the dest
// slice must be integer, float, boolean, string, struct or pointer to struct
// values.
//
// Struct fields must be integer, float, boolean or string values. All struct
// fields are used unless a subset is specified using fieldNames.
func ScanSlice(src []interface***REMOVED******REMOVED***, dest interface***REMOVED******REMOVED***, fieldNames ...string) error ***REMOVED***
	d := reflect.ValueOf(dest)
	if d.Kind() != reflect.Ptr || d.IsNil() ***REMOVED***
		return errScanSliceValue
	***REMOVED***
	d = d.Elem()
	if d.Kind() != reflect.Slice ***REMOVED***
		return errScanSliceValue
	***REMOVED***

	isPtr := false
	t := d.Type().Elem()
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct ***REMOVED***
		isPtr = true
		t = t.Elem()
	***REMOVED***

	if t.Kind() != reflect.Struct ***REMOVED***
		ensureLen(d, len(src))
		for i, s := range src ***REMOVED***
			if s == nil ***REMOVED***
				continue
			***REMOVED***
			if err := convertAssignValue(d.Index(i), s); err != nil ***REMOVED***
				return fmt.Errorf("redigo.ScanSlice: cannot assign element %d: %v", i, err)
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***

	ss := structSpecForType(t)
	fss := ss.l
	if len(fieldNames) > 0 ***REMOVED***
		fss = make([]*fieldSpec, len(fieldNames))
		for i, name := range fieldNames ***REMOVED***
			fss[i] = ss.m[name]
			if fss[i] == nil ***REMOVED***
				return fmt.Errorf("redigo.ScanSlice: ScanSlice bad field name %s", name)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if len(fss) == 0 ***REMOVED***
		return errors.New("redigo.ScanSlice: no struct fields")
	***REMOVED***

	n := len(src) / len(fss)
	if n*len(fss) != len(src) ***REMOVED***
		return errors.New("redigo.ScanSlice: length not a multiple of struct field count")
	***REMOVED***

	ensureLen(d, n)
	for i := 0; i < n; i++ ***REMOVED***
		d := d.Index(i)
		if isPtr ***REMOVED***
			if d.IsNil() ***REMOVED***
				d.Set(reflect.New(t))
			***REMOVED***
			d = d.Elem()
		***REMOVED***
		for j, fs := range fss ***REMOVED***
			s := src[i*len(fss)+j]
			if s == nil ***REMOVED***
				continue
			***REMOVED***
			if err := convertAssignValue(d.FieldByIndex(fs.index), s); err != nil ***REMOVED***
				return fmt.Errorf("redigo.ScanSlice: cannot assign element %d to field %s: %v", i*len(fss)+j, fs.name, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Args is a helper for constructing command arguments from structured values.
type Args []interface***REMOVED******REMOVED***

// Add returns the result of appending value to args.
func (args Args) Add(value ...interface***REMOVED******REMOVED***) Args ***REMOVED***
	return append(args, value...)
***REMOVED***

// AddFlat returns the result of appending the flattened value of v to args.
//
// Maps are flattened by appending the alternating keys and map values to args.
//
// Slices are flattened by appending the slice elements to args.
//
// Structs are flattened by appending the alternating names and values of
// exported fields to args. If v is a nil struct pointer, then nothing is
// appended. The 'redis' field tag overrides struct field names. See ScanStruct
// for more information on the use of the 'redis' field tag.
//
// Other types are appended to args as is.
func (args Args) AddFlat(v interface***REMOVED******REMOVED***) Args ***REMOVED***
	rv := reflect.ValueOf(v)
	switch rv.Kind() ***REMOVED***
	case reflect.Struct:
		args = flattenStruct(args, rv)
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ ***REMOVED***
			args = append(args, rv.Index(i).Interface())
		***REMOVED***
	case reflect.Map:
		for _, k := range rv.MapKeys() ***REMOVED***
			args = append(args, k.Interface(), rv.MapIndex(k).Interface())
		***REMOVED***
	case reflect.Ptr:
		if rv.Type().Elem().Kind() == reflect.Struct ***REMOVED***
			if !rv.IsNil() ***REMOVED***
				args = flattenStruct(args, rv.Elem())
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			args = append(args, v)
		***REMOVED***
	default:
		args = append(args, v)
	***REMOVED***
	return args
***REMOVED***

func flattenStruct(args Args, v reflect.Value) Args ***REMOVED***
	ss := structSpecForType(v.Type())
	for _, fs := range ss.l ***REMOVED***
		fv := v.FieldByIndex(fs.index)
		if fs.omitEmpty ***REMOVED***
			var empty = false
			switch fv.Kind() ***REMOVED***
			case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
				empty = fv.Len() == 0
			case reflect.Bool:
				empty = !fv.Bool()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				empty = fv.Int() == 0
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				empty = fv.Uint() == 0
			case reflect.Float32, reflect.Float64:
				empty = fv.Float() == 0
			case reflect.Interface, reflect.Ptr:
				empty = fv.IsNil()
			***REMOVED***
			if empty ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
		args = append(args, fs.name, fv.Interface())
	***REMOVED***
	return args
***REMOVED***
