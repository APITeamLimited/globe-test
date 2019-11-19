// Copyright (c) 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package envconfig

import (
	"encoding"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ErrInvalidSpecification indicates that a specification is of the wrong type.
var ErrInvalidSpecification = errors.New("specification must be a struct pointer")

var gatherRegexp = regexp.MustCompile("([^A-Z]+|[A-Z]+[^A-Z]+|[A-Z]+)")
var acronymRegexp = regexp.MustCompile("([A-Z]+)([A-Z][^A-Z]+)")

// A ParseError occurs when an environment variable cannot be converted to
// the type required by a struct field during assignment.
type ParseError struct ***REMOVED***
	KeyName   string
	FieldName string
	TypeName  string
	Value     string
	Err       error
***REMOVED***

// Decoder has the same semantics as Setter, but takes higher precedence.
// It is provided for historical compatibility.
type Decoder interface ***REMOVED***
	Decode(value string) error
***REMOVED***

// Setter is implemented by types can self-deserialize values.
// Any type that implements flag.Value also implements Setter.
type Setter interface ***REMOVED***
	Set(value string) error
***REMOVED***

func (e *ParseError) Error() string ***REMOVED***
	return fmt.Sprintf("envconfig.Process: assigning %[1]s to %[2]s: converting '%[3]s' to type %[4]s. details: %[5]s", e.KeyName, e.FieldName, e.Value, e.TypeName, e.Err)
***REMOVED***

// varInfo maintains information about the configuration variable
type varInfo struct ***REMOVED***
	Name  string
	Alt   string
	Key   string
	Field reflect.Value
	Tags  reflect.StructTag
***REMOVED***

// GatherInfo gathers information about the specified struct
func gatherInfo(prefix string, spec interface***REMOVED******REMOVED***) ([]varInfo, error) ***REMOVED***
	s := reflect.ValueOf(spec)

	if s.Kind() != reflect.Ptr ***REMOVED***
		return nil, ErrInvalidSpecification
	***REMOVED***
	s = s.Elem()
	if s.Kind() != reflect.Struct ***REMOVED***
		return nil, ErrInvalidSpecification
	***REMOVED***
	typeOfSpec := s.Type()

	// over allocate an info array, we will extend if needed later
	infos := make([]varInfo, 0, s.NumField())
	for i := 0; i < s.NumField(); i++ ***REMOVED***
		f := s.Field(i)
		ftype := typeOfSpec.Field(i)
		if !f.CanSet() || isTrue(ftype.Tag.Get("ignored")) ***REMOVED***
			continue
		***REMOVED***

		for f.Kind() == reflect.Ptr ***REMOVED***
			if f.IsNil() ***REMOVED***
				if f.Type().Elem().Kind() != reflect.Struct ***REMOVED***
					// nil pointer to a non-struct: leave it alone
					break
				***REMOVED***
				// nil pointer to struct: create a zero instance
				f.Set(reflect.New(f.Type().Elem()))
			***REMOVED***
			f = f.Elem()
		***REMOVED***

		// Capture information about the config variable
		info := varInfo***REMOVED***
			Name:  ftype.Name,
			Field: f,
			Tags:  ftype.Tag,
			Alt:   strings.ToUpper(ftype.Tag.Get("envconfig")),
		***REMOVED***

		// Default to the field name as the env var name (will be upcased)
		info.Key = info.Name

		// Best effort to un-pick camel casing as separate words
		if isTrue(ftype.Tag.Get("split_words")) ***REMOVED***
			words := gatherRegexp.FindAllStringSubmatch(ftype.Name, -1)
			if len(words) > 0 ***REMOVED***
				var name []string
				for _, words := range words ***REMOVED***
					if m := acronymRegexp.FindStringSubmatch(words[0]); len(m) == 3 ***REMOVED***
						name = append(name, m[1], m[2])
					***REMOVED*** else ***REMOVED***
						name = append(name, words[0])
					***REMOVED***
				***REMOVED***

				info.Key = strings.Join(name, "_")
			***REMOVED***
		***REMOVED***
		if info.Alt != "" ***REMOVED***
			info.Key = info.Alt
		***REMOVED***
		if prefix != "" ***REMOVED***
			info.Key = fmt.Sprintf("%s_%s", prefix, info.Key)
		***REMOVED***
		info.Key = strings.ToUpper(info.Key)
		infos = append(infos, info)

		if f.Kind() == reflect.Struct ***REMOVED***
			// honor Decode if present
			if decoderFrom(f) == nil && setterFrom(f) == nil && textUnmarshaler(f) == nil && binaryUnmarshaler(f) == nil ***REMOVED***
				innerPrefix := prefix
				if !ftype.Anonymous ***REMOVED***
					innerPrefix = info.Key
				***REMOVED***

				embeddedPtr := f.Addr().Interface()
				embeddedInfos, err := gatherInfo(innerPrefix, embeddedPtr)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				infos = append(infos[:len(infos)-1], embeddedInfos...)

				continue
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return infos, nil
***REMOVED***

// CheckDisallowed checks that no environment variables with the prefix are set
// that we don't know how or want to parse. This is likely only meaningful with
// a non-empty prefix.
func CheckDisallowed(prefix string, spec interface***REMOVED******REMOVED***) error ***REMOVED***
	infos, err := gatherInfo(prefix, spec)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	vars := make(map[string]struct***REMOVED******REMOVED***)
	for _, info := range infos ***REMOVED***
		vars[info.Key] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	if prefix != "" ***REMOVED***
		prefix = strings.ToUpper(prefix) + "_"
	***REMOVED***

	for _, env := range os.Environ() ***REMOVED***
		if !strings.HasPrefix(env, prefix) ***REMOVED***
			continue
		***REMOVED***
		v := strings.SplitN(env, "=", 2)[0]
		if _, found := vars[v]; !found ***REMOVED***
			return fmt.Errorf("unknown environment variable %s", v)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// Process populates the specified struct based on environment variables
func Process(prefix string, spec interface***REMOVED******REMOVED***) error ***REMOVED***
	infos, err := gatherInfo(prefix, spec)

	for _, info := range infos ***REMOVED***

		// `os.Getenv` cannot differentiate between an explicitly set empty value
		// and an unset value. `os.LookupEnv` is preferred to `syscall.Getenv`,
		// but it is only available in go1.5 or newer. We're using Go build tags
		// here to use os.LookupEnv for >=go1.5
		value, ok := lookupEnv(info.Key)
		if !ok && info.Alt != "" ***REMOVED***
			value, ok = lookupEnv(info.Alt)
		***REMOVED***

		def := info.Tags.Get("default")
		if def != "" && !ok ***REMOVED***
			value = def
		***REMOVED***

		req := info.Tags.Get("required")
		if !ok && def == "" ***REMOVED***
			if isTrue(req) ***REMOVED***
				key := info.Key
				if info.Alt != "" ***REMOVED***
					key = info.Alt
				***REMOVED***
				return fmt.Errorf("required key %s missing value", key)
			***REMOVED***
			continue
		***REMOVED***

		err = processField(value, info.Field)
		if err != nil ***REMOVED***
			return &ParseError***REMOVED***
				KeyName:   info.Key,
				FieldName: info.Name,
				TypeName:  info.Field.Type().String(),
				Value:     value,
				Err:       err,
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return err
***REMOVED***

// MustProcess is the same as Process but panics if an error occurs
func MustProcess(prefix string, spec interface***REMOVED******REMOVED***) ***REMOVED***
	if err := Process(prefix, spec); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func processField(value string, field reflect.Value) error ***REMOVED***
	typ := field.Type()

	decoder := decoderFrom(field)
	if decoder != nil ***REMOVED***
		return decoder.Decode(value)
	***REMOVED***
	// look for Set method if Decode not defined
	setter := setterFrom(field)
	if setter != nil ***REMOVED***
		return setter.Set(value)
	***REMOVED***

	if t := textUnmarshaler(field); t != nil ***REMOVED***
		return t.UnmarshalText([]byte(value))
	***REMOVED***

	if b := binaryUnmarshaler(field); b != nil ***REMOVED***
		return b.UnmarshalBinary([]byte(value))
	***REMOVED***

	if typ.Kind() == reflect.Ptr ***REMOVED***
		typ = typ.Elem()
		if field.IsNil() ***REMOVED***
			field.Set(reflect.New(typ))
		***REMOVED***
		field = field.Elem()
	***REMOVED***

	switch typ.Kind() ***REMOVED***
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var (
			val int64
			err error
		)
		if field.Kind() == reflect.Int64 && typ.PkgPath() == "time" && typ.Name() == "Duration" ***REMOVED***
			var d time.Duration
			d, err = time.ParseDuration(value)
			val = int64(d)
		***REMOVED*** else ***REMOVED***
			val, err = strconv.ParseInt(value, 0, typ.Bits())
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		field.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(value, 0, typ.Bits())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		field.SetUint(val)
	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		field.SetBool(val)
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(value, typ.Bits())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		field.SetFloat(val)
	case reflect.Slice:
		sl := reflect.MakeSlice(typ, 0, 0)
		if typ.Elem().Kind() == reflect.Uint8 ***REMOVED***
			sl = reflect.ValueOf([]byte(value))
		***REMOVED*** else if len(strings.TrimSpace(value)) != 0 ***REMOVED***
			vals := strings.Split(value, ",")
			sl = reflect.MakeSlice(typ, len(vals), len(vals))
			for i, val := range vals ***REMOVED***
				err := processField(val, sl.Index(i))
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
		field.Set(sl)
	case reflect.Map:
		mp := reflect.MakeMap(typ)
		if len(strings.TrimSpace(value)) != 0 ***REMOVED***
			pairs := strings.Split(value, ",")
			for _, pair := range pairs ***REMOVED***
				kvpair := strings.Split(pair, ":")
				if len(kvpair) != 2 ***REMOVED***
					return fmt.Errorf("invalid map item: %q", pair)
				***REMOVED***
				k := reflect.New(typ.Key()).Elem()
				err := processField(kvpair[0], k)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				v := reflect.New(typ.Elem()).Elem()
				err = processField(kvpair[1], v)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				mp.SetMapIndex(k, v)
			***REMOVED***
		***REMOVED***
		field.Set(mp)
	***REMOVED***

	return nil
***REMOVED***

func interfaceFrom(field reflect.Value, fn func(interface***REMOVED******REMOVED***, *bool)) ***REMOVED***
	// it may be impossible for a struct field to fail this check
	if !field.CanInterface() ***REMOVED***
		return
	***REMOVED***
	var ok bool
	fn(field.Interface(), &ok)
	if !ok && field.CanAddr() ***REMOVED***
		fn(field.Addr().Interface(), &ok)
	***REMOVED***
***REMOVED***

func decoderFrom(field reflect.Value) (d Decoder) ***REMOVED***
	interfaceFrom(field, func(v interface***REMOVED******REMOVED***, ok *bool) ***REMOVED*** d, *ok = v.(Decoder) ***REMOVED***)
	return d
***REMOVED***

func setterFrom(field reflect.Value) (s Setter) ***REMOVED***
	interfaceFrom(field, func(v interface***REMOVED******REMOVED***, ok *bool) ***REMOVED*** s, *ok = v.(Setter) ***REMOVED***)
	return s
***REMOVED***

func textUnmarshaler(field reflect.Value) (t encoding.TextUnmarshaler) ***REMOVED***
	interfaceFrom(field, func(v interface***REMOVED******REMOVED***, ok *bool) ***REMOVED*** t, *ok = v.(encoding.TextUnmarshaler) ***REMOVED***)
	return t
***REMOVED***

func binaryUnmarshaler(field reflect.Value) (b encoding.BinaryUnmarshaler) ***REMOVED***
	interfaceFrom(field, func(v interface***REMOVED******REMOVED***, ok *bool) ***REMOVED*** b, *ok = v.(encoding.BinaryUnmarshaler) ***REMOVED***)
	return b
***REMOVED***

func isTrue(s string) bool ***REMOVED***
	b, _ := strconv.ParseBool(s)
	return b
***REMOVED***
