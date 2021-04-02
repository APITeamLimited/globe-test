// Package yaml implements YAML support for the Go language.
//
// Source code and other details for the project are available at GitHub:
//
//   https://github.com/go-yaml/yaml
//
package yaml

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
)

// MapSlice encodes and decodes as a YAML map.
// The order of keys is preserved when encoding and decoding.
type MapSlice []MapItem

// MapItem is an item in a MapSlice.
type MapItem struct ***REMOVED***
	Key, Value interface***REMOVED******REMOVED***
***REMOVED***

// The Unmarshaler interface may be implemented by types to customize their
// behavior when being unmarshaled from a YAML document. The UnmarshalYAML
// method receives a function that may be called to unmarshal the original
// YAML value into a field or variable. It is safe to call the unmarshal
// function parameter more than once if necessary.
type Unmarshaler interface ***REMOVED***
	UnmarshalYAML(unmarshal func(interface***REMOVED******REMOVED***) error) error
***REMOVED***

// The Marshaler interface may be implemented by types to customize their
// behavior when being marshaled into a YAML document. The returned value
// is marshaled in place of the original value implementing Marshaler.
//
// If an error is returned by MarshalYAML, the marshaling procedure stops
// and returns with the provided error.
type Marshaler interface ***REMOVED***
	MarshalYAML() (interface***REMOVED******REMOVED***, error)
***REMOVED***

// Unmarshal decodes the first document found within the in byte slice
// and assigns decoded values into the out value.
//
// Maps and pointers (to a struct, string, int, etc) are accepted as out
// values. If an internal pointer within a struct is not initialized,
// the yaml package will initialize it if necessary for unmarshalling
// the provided data. The out parameter must not be nil.
//
// The type of the decoded values should be compatible with the respective
// values in out. If one or more values cannot be decoded due to a type
// mismatches, decoding continues partially until the end of the YAML
// content, and a *yaml.TypeError is returned with details for all
// missed values.
//
// Struct fields are only unmarshalled if they are exported (have an
// upper case first letter), and are unmarshalled using the field name
// lowercased as the default key. Custom keys may be defined via the
// "yaml" name in the field tag: the content preceding the first comma
// is used as the key, and the following comma-separated options are
// used to tweak the marshalling process (see Marshal).
// Conflicting names result in a runtime error.
//
// For example:
//
//     type T struct ***REMOVED***
//         F int `yaml:"a,omitempty"`
//         B int
//     ***REMOVED***
//     var t T
//     yaml.Unmarshal([]byte("a: 1\nb: 2"), &t)
//
// See the documentation of Marshal for the format of tags and a list of
// supported tag options.
//
func Unmarshal(in []byte, out interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	return unmarshal(in, out, false)
***REMOVED***

// UnmarshalStrict is like Unmarshal except that any fields that are found
// in the data that do not have corresponding struct members, or mapping
// keys that are duplicates, will result in
// an error.
func UnmarshalStrict(in []byte, out interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	return unmarshal(in, out, true)
***REMOVED***

// A Decoder reads and decodes YAML values from an input stream.
type Decoder struct ***REMOVED***
	strict bool
	parser *parser
***REMOVED***

// NewDecoder returns a new decoder that reads from r.
//
// The decoder introduces its own buffering and may read
// data from r beyond the YAML values requested.
func NewDecoder(r io.Reader) *Decoder ***REMOVED***
	return &Decoder***REMOVED***
		parser: newParserFromReader(r),
	***REMOVED***
***REMOVED***

// SetStrict sets whether strict decoding behaviour is enabled when
// decoding items in the data (see UnmarshalStrict). By default, decoding is not strict.
func (dec *Decoder) SetStrict(strict bool) ***REMOVED***
	dec.strict = strict
***REMOVED***

// Decode reads the next YAML-encoded value from its input
// and stores it in the value pointed to by v.
//
// See the documentation for Unmarshal for details about the
// conversion of YAML into a Go value.
func (dec *Decoder) Decode(v interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	d := newDecoder(dec.strict)
	defer handleErr(&err)
	node := dec.parser.parse()
	if node == nil ***REMOVED***
		return io.EOF
	***REMOVED***
	out := reflect.ValueOf(v)
	if out.Kind() == reflect.Ptr && !out.IsNil() ***REMOVED***
		out = out.Elem()
	***REMOVED***
	d.unmarshal(node, out)
	if len(d.terrors) > 0 ***REMOVED***
		return &TypeError***REMOVED***d.terrors***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func unmarshal(in []byte, out interface***REMOVED******REMOVED***, strict bool) (err error) ***REMOVED***
	defer handleErr(&err)
	d := newDecoder(strict)
	p := newParser(in)
	defer p.destroy()
	node := p.parse()
	if node != nil ***REMOVED***
		v := reflect.ValueOf(out)
		if v.Kind() == reflect.Ptr && !v.IsNil() ***REMOVED***
			v = v.Elem()
		***REMOVED***
		d.unmarshal(node, v)
	***REMOVED***
	if len(d.terrors) > 0 ***REMOVED***
		return &TypeError***REMOVED***d.terrors***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Marshal serializes the value provided into a YAML document. The structure
// of the generated document will reflect the structure of the value itself.
// Maps and pointers (to struct, string, int, etc) are accepted as the in value.
//
// Struct fields are only marshalled if they are exported (have an upper case
// first letter), and are marshalled using the field name lowercased as the
// default key. Custom keys may be defined via the "yaml" name in the field
// tag: the content preceding the first comma is used as the key, and the
// following comma-separated options are used to tweak the marshalling process.
// Conflicting names result in a runtime error.
//
// The field tag format accepted is:
//
//     `(...) yaml:"[<key>][,<flag1>[,<flag2>]]" (...)`
//
// The following flags are currently supported:
//
//     omitempty    Only include the field if it's not set to the zero
//                  value for the type or to empty slices or maps.
//                  Zero valued structs will be omitted if all their public
//                  fields are zero, unless they implement an IsZero
//                  method (see the IsZeroer interface type), in which
//                  case the field will be excluded if IsZero returns true.
//
//     flow         Marshal using a flow style (useful for structs,
//                  sequences and maps).
//
//     inline       Inline the field, which must be a struct or a map,
//                  causing all of its fields or keys to be processed as if
//                  they were part of the outer struct. For maps, keys must
//                  not conflict with the yaml keys of other struct fields.
//
// In addition, if the key is "-", the field is ignored.
//
// For example:
//
//     type T struct ***REMOVED***
//         F int `yaml:"a,omitempty"`
//         B int
//     ***REMOVED***
//     yaml.Marshal(&T***REMOVED***B: 2***REMOVED***) // Returns "b: 2\n"
//     yaml.Marshal(&T***REMOVED***F: 1***REMOVED******REMOVED*** // Returns "a: 1\nb: 0\n"
//
func Marshal(in interface***REMOVED******REMOVED***) (out []byte, err error) ***REMOVED***
	defer handleErr(&err)
	e := newEncoder()
	defer e.destroy()
	e.marshalDoc("", reflect.ValueOf(in))
	e.finish()
	out = e.out
	return
***REMOVED***

// An Encoder writes YAML values to an output stream.
type Encoder struct ***REMOVED***
	encoder *encoder
***REMOVED***

// NewEncoder returns a new encoder that writes to w.
// The Encoder should be closed after use to flush all data
// to w.
func NewEncoder(w io.Writer) *Encoder ***REMOVED***
	return &Encoder***REMOVED***
		encoder: newEncoderWithWriter(w),
	***REMOVED***
***REMOVED***

// Encode writes the YAML encoding of v to the stream.
// If multiple items are encoded to the stream, the
// second and subsequent document will be preceded
// with a "---" document separator, but the first will not.
//
// See the documentation for Marshal for details about the conversion of Go
// values to YAML.
func (e *Encoder) Encode(v interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	defer handleErr(&err)
	e.encoder.marshalDoc("", reflect.ValueOf(v))
	return nil
***REMOVED***

// Close closes the encoder by writing any remaining data.
// It does not write a stream terminating string "...".
func (e *Encoder) Close() (err error) ***REMOVED***
	defer handleErr(&err)
	e.encoder.finish()
	return nil
***REMOVED***

func handleErr(err *error) ***REMOVED***
	if v := recover(); v != nil ***REMOVED***
		if e, ok := v.(yamlError); ok ***REMOVED***
			*err = e.err
		***REMOVED*** else ***REMOVED***
			panic(v)
		***REMOVED***
	***REMOVED***
***REMOVED***

type yamlError struct ***REMOVED***
	err error
***REMOVED***

func fail(err error) ***REMOVED***
	panic(yamlError***REMOVED***err***REMOVED***)
***REMOVED***

func failf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	panic(yamlError***REMOVED***fmt.Errorf("yaml: "+format, args...)***REMOVED***)
***REMOVED***

// A TypeError is returned by Unmarshal when one or more fields in
// the YAML document cannot be properly decoded into the requested
// types. When this error is returned, the value is still
// unmarshaled partially.
type TypeError struct ***REMOVED***
	Errors []string
***REMOVED***

func (e *TypeError) Error() string ***REMOVED***
	return fmt.Sprintf("yaml: unmarshal errors:\n  %s", strings.Join(e.Errors, "\n  "))
***REMOVED***

// --------------------------------------------------------------------------
// Maintain a mapping of keys to structure field indexes

// The code in this section was copied from mgo/bson.

// structInfo holds details for the serialization of fields of
// a given struct.
type structInfo struct ***REMOVED***
	FieldsMap  map[string]fieldInfo
	FieldsList []fieldInfo

	// InlineMap is the number of the field in the struct that
	// contains an ,inline map, or -1 if there's none.
	InlineMap int
***REMOVED***

type fieldInfo struct ***REMOVED***
	Key       string
	Num       int
	OmitEmpty bool
	Flow      bool
	// Id holds the unique field identifier, so we can cheaply
	// check for field duplicates without maintaining an extra map.
	Id int

	// Inline holds the field index if the field is part of an inlined struct.
	Inline []int
***REMOVED***

var structMap = make(map[reflect.Type]*structInfo)
var fieldMapMutex sync.RWMutex

func getStructInfo(st reflect.Type) (*structInfo, error) ***REMOVED***
	fieldMapMutex.RLock()
	sinfo, found := structMap[st]
	fieldMapMutex.RUnlock()
	if found ***REMOVED***
		return sinfo, nil
	***REMOVED***

	n := st.NumField()
	fieldsMap := make(map[string]fieldInfo)
	fieldsList := make([]fieldInfo, 0, n)
	inlineMap := -1
	for i := 0; i != n; i++ ***REMOVED***
		field := st.Field(i)
		if field.PkgPath != "" && !field.Anonymous ***REMOVED***
			continue // Private field
		***REMOVED***

		info := fieldInfo***REMOVED***Num: i***REMOVED***

		tag := field.Tag.Get("yaml")
		if tag == "" && strings.Index(string(field.Tag), ":") < 0 ***REMOVED***
			tag = string(field.Tag)
		***REMOVED***
		if tag == "-" ***REMOVED***
			continue
		***REMOVED***

		inline := false
		fields := strings.Split(tag, ",")
		if len(fields) > 1 ***REMOVED***
			for _, flag := range fields[1:] ***REMOVED***
				switch flag ***REMOVED***
				case "omitempty":
					info.OmitEmpty = true
				case "flow":
					info.Flow = true
				case "inline":
					inline = true
				default:
					return nil, errors.New(fmt.Sprintf("Unsupported flag %q in tag %q of type %s", flag, tag, st))
				***REMOVED***
			***REMOVED***
			tag = fields[0]
		***REMOVED***

		if inline ***REMOVED***
			switch field.Type.Kind() ***REMOVED***
			case reflect.Map:
				if inlineMap >= 0 ***REMOVED***
					return nil, errors.New("Multiple ,inline maps in struct " + st.String())
				***REMOVED***
				if field.Type.Key() != reflect.TypeOf("") ***REMOVED***
					return nil, errors.New("Option ,inline needs a map with string keys in struct " + st.String())
				***REMOVED***
				inlineMap = info.Num
			case reflect.Struct:
				sinfo, err := getStructInfo(field.Type)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				for _, finfo := range sinfo.FieldsList ***REMOVED***
					if _, found := fieldsMap[finfo.Key]; found ***REMOVED***
						msg := "Duplicated key '" + finfo.Key + "' in struct " + st.String()
						return nil, errors.New(msg)
					***REMOVED***
					if finfo.Inline == nil ***REMOVED***
						finfo.Inline = []int***REMOVED***i, finfo.Num***REMOVED***
					***REMOVED*** else ***REMOVED***
						finfo.Inline = append([]int***REMOVED***i***REMOVED***, finfo.Inline...)
					***REMOVED***
					finfo.Id = len(fieldsList)
					fieldsMap[finfo.Key] = finfo
					fieldsList = append(fieldsList, finfo)
				***REMOVED***
			default:
				//return nil, errors.New("Option ,inline needs a struct value or map field")
				return nil, errors.New("Option ,inline needs a struct value field")
			***REMOVED***
			continue
		***REMOVED***

		if tag != "" ***REMOVED***
			info.Key = tag
		***REMOVED*** else ***REMOVED***
			info.Key = strings.ToLower(field.Name)
		***REMOVED***

		if _, found = fieldsMap[info.Key]; found ***REMOVED***
			msg := "Duplicated key '" + info.Key + "' in struct " + st.String()
			return nil, errors.New(msg)
		***REMOVED***

		info.Id = len(fieldsList)
		fieldsList = append(fieldsList, info)
		fieldsMap[info.Key] = info
	***REMOVED***

	sinfo = &structInfo***REMOVED***
		FieldsMap:  fieldsMap,
		FieldsList: fieldsList,
		InlineMap:  inlineMap,
	***REMOVED***

	fieldMapMutex.Lock()
	structMap[st] = sinfo
	fieldMapMutex.Unlock()
	return sinfo, nil
***REMOVED***

// IsZeroer is used to check whether an object is zero to
// determine whether it should be omitted when marshaling
// with the omitempty flag. One notable implementation
// is time.Time.
type IsZeroer interface ***REMOVED***
	IsZero() bool
***REMOVED***

func isZero(v reflect.Value) bool ***REMOVED***
	kind := v.Kind()
	if z, ok := v.Interface().(IsZeroer); ok ***REMOVED***
		if (kind == reflect.Ptr || kind == reflect.Interface) && v.IsNil() ***REMOVED***
			return true
		***REMOVED***
		return z.IsZero()
	***REMOVED***
	switch kind ***REMOVED***
	case reflect.String:
		return len(v.String()) == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Slice:
		return v.Len() == 0
	case reflect.Map:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Struct:
		vt := v.Type()
		for i := v.NumField() - 1; i >= 0; i-- ***REMOVED***
			if vt.Field(i).PkgPath != "" ***REMOVED***
				continue // Private field
			***REMOVED***
			if !isZero(v.Field(i)) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

// FutureLineWrap globally disables line wrapping when encoding long strings.
// This is a temporary and thus deprecated method introduced to faciliate
// migration towards v3, which offers more control of line lengths on
// individual encodings, and has a default matching the behavior introduced
// by this function.
//
// The default formatting of v2 was erroneously changed in v2.3.0 and reverted
// in v2.4.0, at which point this function was introduced to help migration.
func FutureLineWrap() ***REMOVED***
	disableLineWrapping = true
***REMOVED***
