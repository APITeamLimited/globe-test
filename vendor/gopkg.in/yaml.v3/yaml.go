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
	"unicode/utf8"
)

// The Unmarshaler interface may be implemented by types to customize their
// behavior when being unmarshaled from a YAML document.
type Unmarshaler interface ***REMOVED***
	UnmarshalYAML(value *Node) error
***REMOVED***

type obsoleteUnmarshaler interface ***REMOVED***
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

// A Decorder reads and decodes YAML values from an input stream.
type Decoder struct ***REMOVED***
	parser      *parser
	knownFields bool
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

// KnownFields ensures that the keys in decoded mappings to
// exist as fields in the struct being decoded into.
func (dec *Decoder) KnownFields(enable bool) ***REMOVED***
	dec.knownFields = enable
***REMOVED***

// Decode reads the next YAML-encoded value from its input
// and stores it in the value pointed to by v.
//
// See the documentation for Unmarshal for details about the
// conversion of YAML into a Go value.
func (dec *Decoder) Decode(v interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	d := newDecoder()
	d.knownFields = dec.knownFields
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

// Decode decodes the node and stores its data into the value pointed to by v.
//
// See the documentation for Unmarshal for details about the
// conversion of YAML into a Go value.
func (n *Node) Decode(v interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	d := newDecoder()
	defer handleErr(&err)
	out := reflect.ValueOf(v)
	if out.Kind() == reflect.Ptr && !out.IsNil() ***REMOVED***
		out = out.Elem()
	***REMOVED***
	d.unmarshal(n, out)
	if len(d.terrors) > 0 ***REMOVED***
		return &TypeError***REMOVED***d.terrors***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func unmarshal(in []byte, out interface***REMOVED******REMOVED***, strict bool) (err error) ***REMOVED***
	defer handleErr(&err)
	d := newDecoder()
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
//                  case the field will be included if that method returns true.
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

// SetIndent changes the used indentation used when encoding.
func (e *Encoder) SetIndent(spaces int) ***REMOVED***
	if spaces < 0 ***REMOVED***
		panic("yaml: cannot indent to a negative number of spaces")
	***REMOVED***
	e.encoder.indent = spaces
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

type Kind uint32

const (
	DocumentNode Kind = 1 << iota
	SequenceNode
	MappingNode
	ScalarNode
	AliasNode
)

type Style uint32

const (
	TaggedStyle Style = 1 << iota
	DoubleQuotedStyle
	SingleQuotedStyle
	LiteralStyle
	FoldedStyle
	FlowStyle
)

// Node represents an element in the YAML document hierarchy. While documents
// are typically encoded and decoded into higher level types, such as structs
// and maps, Node is an intermediate representation that allows detailed
// control over the content being decoded or encoded.
//
// Values that make use of the Node type interact with the yaml package in the
// same way any other type would do, by encoding and decoding yaml data
// directly or indirectly into them.
//
// For example:
//
//     var person struct ***REMOVED***
//             Name    string
//             Address yaml.Node
//     ***REMOVED***
//     err := yaml.Unmarshal(data, &person)
// 
// Or by itself:
//
//     var person Node
//     err := yaml.Unmarshal(data, &person)
//
type Node struct ***REMOVED***
	// Kind defines whether the node is a document, a mapping, a sequence,
	// a scalar value, or an alias to another node. The specific data type of
	// scalar nodes may be obtained via the ShortTag and LongTag methods.
	Kind  Kind

	// Style allows customizing the apperance of the node in the tree.
	Style Style

	// Tag holds the YAML tag defining the data type for the value.
	// When decoding, this field will always be set to the resolved tag,
	// even when it wasn't explicitly provided in the YAML content.
	// When encoding, if this field is unset the value type will be
	// implied from the node properties, and if it is set, it will only
	// be serialized into the representation if TaggedStyle is used or
	// the implicit tag diverges from the provided one.
	Tag string

	// Value holds the unescaped and unquoted represenation of the value.
	Value string

	// Anchor holds the anchor name for this node, which allows aliases to point to it.
	Anchor string

	// Alias holds the node that this alias points to. Only valid when Kind is AliasNode.
	Alias *Node

	// Content holds contained nodes for documents, mappings, and sequences.
	Content []*Node

	// HeadComment holds any comments in the lines preceding the node and
	// not separated by an empty line.
	HeadComment string

	// LineComment holds any comments at the end of the line where the node is in.
	LineComment string

	// FootComment holds any comments following the node and before empty lines.
	FootComment string

	// Line and Column hold the node position in the decoded YAML text.
	// These fields are not respected when encoding the node.
	Line   int
	Column int
***REMOVED***

// LongTag returns the long form of the tag that indicates the data type for
// the node. If the Tag field isn't explicitly defined, one will be computed
// based on the node properties.
func (n *Node) LongTag() string ***REMOVED***
	return longTag(n.ShortTag())
***REMOVED***

// ShortTag returns the short form of the YAML tag that indicates data type for
// the node. If the Tag field isn't explicitly defined, one will be computed
// based on the node properties.
func (n *Node) ShortTag() string ***REMOVED***
	if n.indicatedString() ***REMOVED***
		return strTag
	***REMOVED***
	if n.Tag == "" || n.Tag == "!" ***REMOVED***
		switch n.Kind ***REMOVED***
		case MappingNode:
			return mapTag
		case SequenceNode:
			return seqTag
		case AliasNode:
			if n.Alias != nil ***REMOVED***
				return n.Alias.ShortTag()
			***REMOVED***
		case ScalarNode:
			tag, _ := resolve("", n.Value)
			return tag
		***REMOVED***
		return ""
	***REMOVED***
	return shortTag(n.Tag)
***REMOVED***

func (n *Node) indicatedString() bool ***REMOVED***
	return n.Kind == ScalarNode &&
		(shortTag(n.Tag) == strTag ||
			(n.Tag == "" || n.Tag == "!") && n.Style&(SingleQuotedStyle|DoubleQuotedStyle|LiteralStyle|FoldedStyle) != 0)
***REMOVED***

// SetString is a convenience function that sets the node to a string value
// and defines its style in a pleasant way depending on its content.
func (n *Node) SetString(s string) ***REMOVED***
	n.Kind = ScalarNode
	if utf8.ValidString(s) ***REMOVED***
		n.Value = s
		n.Tag = strTag
	***REMOVED*** else ***REMOVED***
		n.Value = encodeBase64(s)
		n.Tag = binaryTag
	***REMOVED***
	if strings.Contains(n.Value, "\n") ***REMOVED***
		n.Style = LiteralStyle
	***REMOVED***
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

	// InlineUnmarshalers holds indexes to inlined fields that
	// contain unmarshaler values.
	InlineUnmarshalers [][]int
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
var unmarshalerType reflect.Type

func init() ***REMOVED***
	var v Unmarshaler
	unmarshalerType = reflect.ValueOf(&v).Elem().Type()
***REMOVED***

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
	inlineUnmarshalers := [][]int(nil)
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
					return nil, errors.New(fmt.Sprintf("unsupported flag %q in tag %q of type %s", flag, tag, st))
				***REMOVED***
			***REMOVED***
			tag = fields[0]
		***REMOVED***

		if inline ***REMOVED***
			switch field.Type.Kind() ***REMOVED***
			case reflect.Map:
				if inlineMap >= 0 ***REMOVED***
					return nil, errors.New("multiple ,inline maps in struct " + st.String())
				***REMOVED***
				if field.Type.Key() != reflect.TypeOf("") ***REMOVED***
					return nil, errors.New("option ,inline needs a map with string keys in struct " + st.String())
				***REMOVED***
				inlineMap = info.Num
			case reflect.Struct, reflect.Ptr:
				ftype := field.Type
				for ftype.Kind() == reflect.Ptr ***REMOVED***
					ftype = ftype.Elem()
				***REMOVED***
				if ftype.Kind() != reflect.Struct ***REMOVED***
					return nil, errors.New("option ,inline may only be used on a struct or map field")
				***REMOVED***
				if reflect.PtrTo(ftype).Implements(unmarshalerType) ***REMOVED***
					inlineUnmarshalers = append(inlineUnmarshalers, []int***REMOVED***i***REMOVED***)
				***REMOVED*** else ***REMOVED***
					sinfo, err := getStructInfo(ftype)
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***
					for _, index := range sinfo.InlineUnmarshalers ***REMOVED***
						inlineUnmarshalers = append(inlineUnmarshalers, append([]int***REMOVED***i***REMOVED***, index...))
					***REMOVED***
					for _, finfo := range sinfo.FieldsList ***REMOVED***
						if _, found := fieldsMap[finfo.Key]; found ***REMOVED***
							msg := "duplicated key '" + finfo.Key + "' in struct " + st.String()
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
				***REMOVED***
			default:
				return nil, errors.New("option ,inline may only be used on a struct or map field")
			***REMOVED***
			continue
		***REMOVED***

		if tag != "" ***REMOVED***
			info.Key = tag
		***REMOVED*** else ***REMOVED***
			info.Key = strings.ToLower(field.Name)
		***REMOVED***

		if _, found = fieldsMap[info.Key]; found ***REMOVED***
			msg := "duplicated key '" + info.Key + "' in struct " + st.String()
			return nil, errors.New(msg)
		***REMOVED***

		info.Id = len(fieldsList)
		fieldsList = append(fieldsList, info)
		fieldsMap[info.Key] = info
	***REMOVED***

	sinfo = &structInfo***REMOVED***
		FieldsMap:          fieldsMap,
		FieldsList:         fieldsList,
		InlineMap:          inlineMap,
		InlineUnmarshalers: inlineUnmarshalers,
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
