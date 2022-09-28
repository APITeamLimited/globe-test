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
	"encoding"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type encoder struct ***REMOVED***
	emitter  yaml_emitter_t
	event    yaml_event_t
	out      []byte
	flow     bool
	indent   int
	doneInit bool
***REMOVED***

func newEncoder() *encoder ***REMOVED***
	e := &encoder***REMOVED******REMOVED***
	yaml_emitter_initialize(&e.emitter)
	yaml_emitter_set_output_string(&e.emitter, &e.out)
	yaml_emitter_set_unicode(&e.emitter, true)
	return e
***REMOVED***

func newEncoderWithWriter(w io.Writer) *encoder ***REMOVED***
	e := &encoder***REMOVED******REMOVED***
	yaml_emitter_initialize(&e.emitter)
	yaml_emitter_set_output_writer(&e.emitter, w)
	yaml_emitter_set_unicode(&e.emitter, true)
	return e
***REMOVED***

func (e *encoder) init() ***REMOVED***
	if e.doneInit ***REMOVED***
		return
	***REMOVED***
	if e.indent == 0 ***REMOVED***
		e.indent = 4
	***REMOVED***
	e.emitter.best_indent = e.indent
	yaml_stream_start_event_initialize(&e.event, yaml_UTF8_ENCODING)
	e.emit()
	e.doneInit = true
***REMOVED***

func (e *encoder) finish() ***REMOVED***
	e.emitter.open_ended = false
	yaml_stream_end_event_initialize(&e.event)
	e.emit()
***REMOVED***

func (e *encoder) destroy() ***REMOVED***
	yaml_emitter_delete(&e.emitter)
***REMOVED***

func (e *encoder) emit() ***REMOVED***
	// This will internally delete the e.event value.
	e.must(yaml_emitter_emit(&e.emitter, &e.event))
***REMOVED***

func (e *encoder) must(ok bool) ***REMOVED***
	if !ok ***REMOVED***
		msg := e.emitter.problem
		if msg == "" ***REMOVED***
			msg = "unknown problem generating YAML content"
		***REMOVED***
		failf("%s", msg)
	***REMOVED***
***REMOVED***

func (e *encoder) marshalDoc(tag string, in reflect.Value) ***REMOVED***
	e.init()
	var node *Node
	if in.IsValid() ***REMOVED***
		node, _ = in.Interface().(*Node)
	***REMOVED***
	if node != nil && node.Kind == DocumentNode ***REMOVED***
		e.nodev(in)
	***REMOVED*** else ***REMOVED***
		yaml_document_start_event_initialize(&e.event, nil, nil, true)
		e.emit()
		e.marshal(tag, in)
		yaml_document_end_event_initialize(&e.event, true)
		e.emit()
	***REMOVED***
***REMOVED***

func (e *encoder) marshal(tag string, in reflect.Value) ***REMOVED***
	tag = shortTag(tag)
	if !in.IsValid() || in.Kind() == reflect.Ptr && in.IsNil() ***REMOVED***
		e.nilv()
		return
	***REMOVED***
	iface := in.Interface()
	switch value := iface.(type) ***REMOVED***
	case *Node:
		e.nodev(in)
		return
	case Node:
		if !in.CanAddr() ***REMOVED***
			var n = reflect.New(in.Type()).Elem()
			n.Set(in)
			in = n
		***REMOVED***
		e.nodev(in.Addr())
		return
	case time.Time:
		e.timev(tag, in)
		return
	case *time.Time:
		e.timev(tag, in.Elem())
		return
	case time.Duration:
		e.stringv(tag, reflect.ValueOf(value.String()))
		return
	case Marshaler:
		v, err := value.MarshalYAML()
		if err != nil ***REMOVED***
			fail(err)
		***REMOVED***
		if v == nil ***REMOVED***
			e.nilv()
			return
		***REMOVED***
		e.marshal(tag, reflect.ValueOf(v))
		return
	case encoding.TextMarshaler:
		text, err := value.MarshalText()
		if err != nil ***REMOVED***
			fail(err)
		***REMOVED***
		in = reflect.ValueOf(string(text))
	case nil:
		e.nilv()
		return
	***REMOVED***
	switch in.Kind() ***REMOVED***
	case reflect.Interface:
		e.marshal(tag, in.Elem())
	case reflect.Map:
		e.mapv(tag, in)
	case reflect.Ptr:
		e.marshal(tag, in.Elem())
	case reflect.Struct:
		e.structv(tag, in)
	case reflect.Slice, reflect.Array:
		e.slicev(tag, in)
	case reflect.String:
		e.stringv(tag, in)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		e.intv(tag, in)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		e.uintv(tag, in)
	case reflect.Float32, reflect.Float64:
		e.floatv(tag, in)
	case reflect.Bool:
		e.boolv(tag, in)
	default:
		panic("cannot marshal type: " + in.Type().String())
	***REMOVED***
***REMOVED***

func (e *encoder) mapv(tag string, in reflect.Value) ***REMOVED***
	e.mappingv(tag, func() ***REMOVED***
		keys := keyList(in.MapKeys())
		sort.Sort(keys)
		for _, k := range keys ***REMOVED***
			e.marshal("", k)
			e.marshal("", in.MapIndex(k))
		***REMOVED***
	***REMOVED***)
***REMOVED***

func (e *encoder) fieldByIndex(v reflect.Value, index []int) (field reflect.Value) ***REMOVED***
	for _, num := range index ***REMOVED***
		for ***REMOVED***
			if v.Kind() == reflect.Ptr ***REMOVED***
				if v.IsNil() ***REMOVED***
					return reflect.Value***REMOVED******REMOVED***
				***REMOVED***
				v = v.Elem()
				continue
			***REMOVED***
			break
		***REMOVED***
		v = v.Field(num)
	***REMOVED***
	return v
***REMOVED***

func (e *encoder) structv(tag string, in reflect.Value) ***REMOVED***
	sinfo, err := getStructInfo(in.Type())
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	e.mappingv(tag, func() ***REMOVED***
		for _, info := range sinfo.FieldsList ***REMOVED***
			var value reflect.Value
			if info.Inline == nil ***REMOVED***
				value = in.Field(info.Num)
			***REMOVED*** else ***REMOVED***
				value = e.fieldByIndex(in, info.Inline)
				if !value.IsValid() ***REMOVED***
					continue
				***REMOVED***
			***REMOVED***
			if info.OmitEmpty && isZero(value) ***REMOVED***
				continue
			***REMOVED***
			e.marshal("", reflect.ValueOf(info.Key))
			e.flow = info.Flow
			e.marshal("", value)
		***REMOVED***
		if sinfo.InlineMap >= 0 ***REMOVED***
			m := in.Field(sinfo.InlineMap)
			if m.Len() > 0 ***REMOVED***
				e.flow = false
				keys := keyList(m.MapKeys())
				sort.Sort(keys)
				for _, k := range keys ***REMOVED***
					if _, found := sinfo.FieldsMap[k.String()]; found ***REMOVED***
						panic(fmt.Sprintf("cannot have key %q in inlined map: conflicts with struct field", k.String()))
					***REMOVED***
					e.marshal("", k)
					e.flow = false
					e.marshal("", m.MapIndex(k))
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***

func (e *encoder) mappingv(tag string, f func()) ***REMOVED***
	implicit := tag == ""
	style := yaml_BLOCK_MAPPING_STYLE
	if e.flow ***REMOVED***
		e.flow = false
		style = yaml_FLOW_MAPPING_STYLE
	***REMOVED***
	yaml_mapping_start_event_initialize(&e.event, nil, []byte(tag), implicit, style)
	e.emit()
	f()
	yaml_mapping_end_event_initialize(&e.event)
	e.emit()
***REMOVED***

func (e *encoder) slicev(tag string, in reflect.Value) ***REMOVED***
	implicit := tag == ""
	style := yaml_BLOCK_SEQUENCE_STYLE
	if e.flow ***REMOVED***
		e.flow = false
		style = yaml_FLOW_SEQUENCE_STYLE
	***REMOVED***
	e.must(yaml_sequence_start_event_initialize(&e.event, nil, []byte(tag), implicit, style))
	e.emit()
	n := in.Len()
	for i := 0; i < n; i++ ***REMOVED***
		e.marshal("", in.Index(i))
	***REMOVED***
	e.must(yaml_sequence_end_event_initialize(&e.event))
	e.emit()
***REMOVED***

// isBase60 returns whether s is in base 60 notation as defined in YAML 1.1.
//
// The base 60 float notation in YAML 1.1 is a terrible idea and is unsupported
// in YAML 1.2 and by this package, but these should be marshalled quoted for
// the time being for compatibility with other parsers.
func isBase60Float(s string) (result bool) ***REMOVED***
	// Fast path.
	if s == "" ***REMOVED***
		return false
	***REMOVED***
	c := s[0]
	if !(c == '+' || c == '-' || c >= '0' && c <= '9') || strings.IndexByte(s, ':') < 0 ***REMOVED***
		return false
	***REMOVED***
	// Do the full match.
	return base60float.MatchString(s)
***REMOVED***

// From http://yaml.org/type/float.html, except the regular expression there
// is bogus. In practice parsers do not enforce the "\.[0-9_]*" suffix.
var base60float = regexp.MustCompile(`^[-+]?[0-9][0-9_]*(?::[0-5]?[0-9])+(?:\.[0-9_]*)?$`)

// isOldBool returns whether s is bool notation as defined in YAML 1.1.
//
// We continue to force strings that YAML 1.1 would interpret as booleans to be
// rendered as quotes strings so that the marshalled output valid for YAML 1.1
// parsing.
func isOldBool(s string) (result bool) ***REMOVED***
	switch s ***REMOVED***
	case "y", "Y", "yes", "Yes", "YES", "on", "On", "ON",
		"n", "N", "no", "No", "NO", "off", "Off", "OFF":
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

func (e *encoder) stringv(tag string, in reflect.Value) ***REMOVED***
	var style yaml_scalar_style_t
	s := in.String()
	canUsePlain := true
	switch ***REMOVED***
	case !utf8.ValidString(s):
		if tag == binaryTag ***REMOVED***
			failf("explicitly tagged !!binary data must be base64-encoded")
		***REMOVED***
		if tag != "" ***REMOVED***
			failf("cannot marshal invalid UTF-8 data as %s", shortTag(tag))
		***REMOVED***
		// It can't be encoded directly as YAML so use a binary tag
		// and encode it as base64.
		tag = binaryTag
		s = encodeBase64(s)
	case tag == "":
		// Check to see if it would resolve to a specific
		// tag when encoded unquoted. If it doesn't,
		// there's no need to quote it.
		rtag, _ := resolve("", s)
		canUsePlain = rtag == strTag && !(isBase60Float(s) || isOldBool(s))
	***REMOVED***
	// Note: it's possible for user code to emit invalid YAML
	// if they explicitly specify a tag and a string containing
	// text that's incompatible with that tag.
	switch ***REMOVED***
	case strings.Contains(s, "\n"):
		if e.flow ***REMOVED***
			style = yaml_DOUBLE_QUOTED_SCALAR_STYLE
		***REMOVED*** else ***REMOVED***
			style = yaml_LITERAL_SCALAR_STYLE
		***REMOVED***
	case canUsePlain:
		style = yaml_PLAIN_SCALAR_STYLE
	default:
		style = yaml_DOUBLE_QUOTED_SCALAR_STYLE
	***REMOVED***
	e.emitScalar(s, "", tag, style, nil, nil, nil, nil)
***REMOVED***

func (e *encoder) boolv(tag string, in reflect.Value) ***REMOVED***
	var s string
	if in.Bool() ***REMOVED***
		s = "true"
	***REMOVED*** else ***REMOVED***
		s = "false"
	***REMOVED***
	e.emitScalar(s, "", tag, yaml_PLAIN_SCALAR_STYLE, nil, nil, nil, nil)
***REMOVED***

func (e *encoder) intv(tag string, in reflect.Value) ***REMOVED***
	s := strconv.FormatInt(in.Int(), 10)
	e.emitScalar(s, "", tag, yaml_PLAIN_SCALAR_STYLE, nil, nil, nil, nil)
***REMOVED***

func (e *encoder) uintv(tag string, in reflect.Value) ***REMOVED***
	s := strconv.FormatUint(in.Uint(), 10)
	e.emitScalar(s, "", tag, yaml_PLAIN_SCALAR_STYLE, nil, nil, nil, nil)
***REMOVED***

func (e *encoder) timev(tag string, in reflect.Value) ***REMOVED***
	t := in.Interface().(time.Time)
	s := t.Format(time.RFC3339Nano)
	e.emitScalar(s, "", tag, yaml_PLAIN_SCALAR_STYLE, nil, nil, nil, nil)
***REMOVED***

func (e *encoder) floatv(tag string, in reflect.Value) ***REMOVED***
	// Issue #352: When formatting, use the precision of the underlying value
	precision := 64
	if in.Kind() == reflect.Float32 ***REMOVED***
		precision = 32
	***REMOVED***

	s := strconv.FormatFloat(in.Float(), 'g', -1, precision)
	switch s ***REMOVED***
	case "+Inf":
		s = ".inf"
	case "-Inf":
		s = "-.inf"
	case "NaN":
		s = ".nan"
	***REMOVED***
	e.emitScalar(s, "", tag, yaml_PLAIN_SCALAR_STYLE, nil, nil, nil, nil)
***REMOVED***

func (e *encoder) nilv() ***REMOVED***
	e.emitScalar("null", "", "", yaml_PLAIN_SCALAR_STYLE, nil, nil, nil, nil)
***REMOVED***

func (e *encoder) emitScalar(value, anchor, tag string, style yaml_scalar_style_t, head, line, foot, tail []byte) ***REMOVED***
	// TODO Kill this function. Replace all initialize calls by their underlining Go literals.
	implicit := tag == ""
	if !implicit ***REMOVED***
		tag = longTag(tag)
	***REMOVED***
	e.must(yaml_scalar_event_initialize(&e.event, []byte(anchor), []byte(tag), []byte(value), implicit, implicit, style))
	e.event.head_comment = head
	e.event.line_comment = line
	e.event.foot_comment = foot
	e.event.tail_comment = tail
	e.emit()
***REMOVED***

func (e *encoder) nodev(in reflect.Value) ***REMOVED***
	e.node(in.Interface().(*Node), "")
***REMOVED***

func (e *encoder) node(node *Node, tail string) ***REMOVED***
	// Zero nodes behave as nil.
	if node.Kind == 0 && node.IsZero() ***REMOVED***
		e.nilv()
		return
	***REMOVED***

	// If the tag was not explicitly requested, and dropping it won't change the
	// implicit tag of the value, don't include it in the presentation.
	var tag = node.Tag
	var stag = shortTag(tag)
	var forceQuoting bool
	if tag != "" && node.Style&TaggedStyle == 0 ***REMOVED***
		if node.Kind == ScalarNode ***REMOVED***
			if stag == strTag && node.Style&(SingleQuotedStyle|DoubleQuotedStyle|LiteralStyle|FoldedStyle) != 0 ***REMOVED***
				tag = ""
			***REMOVED*** else ***REMOVED***
				rtag, _ := resolve("", node.Value)
				if rtag == stag ***REMOVED***
					tag = ""
				***REMOVED*** else if stag == strTag ***REMOVED***
					tag = ""
					forceQuoting = true
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			var rtag string
			switch node.Kind ***REMOVED***
			case MappingNode:
				rtag = mapTag
			case SequenceNode:
				rtag = seqTag
			***REMOVED***
			if rtag == stag ***REMOVED***
				tag = ""
			***REMOVED***
		***REMOVED***
	***REMOVED***

	switch node.Kind ***REMOVED***
	case DocumentNode:
		yaml_document_start_event_initialize(&e.event, nil, nil, true)
		e.event.head_comment = []byte(node.HeadComment)
		e.emit()
		for _, node := range node.Content ***REMOVED***
			e.node(node, "")
		***REMOVED***
		yaml_document_end_event_initialize(&e.event, true)
		e.event.foot_comment = []byte(node.FootComment)
		e.emit()

	case SequenceNode:
		style := yaml_BLOCK_SEQUENCE_STYLE
		if node.Style&FlowStyle != 0 ***REMOVED***
			style = yaml_FLOW_SEQUENCE_STYLE
		***REMOVED***
		e.must(yaml_sequence_start_event_initialize(&e.event, []byte(node.Anchor), []byte(longTag(tag)), tag == "", style))
		e.event.head_comment = []byte(node.HeadComment)
		e.emit()
		for _, node := range node.Content ***REMOVED***
			e.node(node, "")
		***REMOVED***
		e.must(yaml_sequence_end_event_initialize(&e.event))
		e.event.line_comment = []byte(node.LineComment)
		e.event.foot_comment = []byte(node.FootComment)
		e.emit()

	case MappingNode:
		style := yaml_BLOCK_MAPPING_STYLE
		if node.Style&FlowStyle != 0 ***REMOVED***
			style = yaml_FLOW_MAPPING_STYLE
		***REMOVED***
		yaml_mapping_start_event_initialize(&e.event, []byte(node.Anchor), []byte(longTag(tag)), tag == "", style)
		e.event.tail_comment = []byte(tail)
		e.event.head_comment = []byte(node.HeadComment)
		e.emit()

		// The tail logic below moves the foot comment of prior keys to the following key,
		// since the value for each key may be a nested structure and the foot needs to be
		// processed only the entirety of the value is streamed. The last tail is processed
		// with the mapping end event.
		var tail string
		for i := 0; i+1 < len(node.Content); i += 2 ***REMOVED***
			k := node.Content[i]
			foot := k.FootComment
			if foot != "" ***REMOVED***
				kopy := *k
				kopy.FootComment = ""
				k = &kopy
			***REMOVED***
			e.node(k, tail)
			tail = foot

			v := node.Content[i+1]
			e.node(v, "")
		***REMOVED***

		yaml_mapping_end_event_initialize(&e.event)
		e.event.tail_comment = []byte(tail)
		e.event.line_comment = []byte(node.LineComment)
		e.event.foot_comment = []byte(node.FootComment)
		e.emit()

	case AliasNode:
		yaml_alias_event_initialize(&e.event, []byte(node.Value))
		e.event.head_comment = []byte(node.HeadComment)
		e.event.line_comment = []byte(node.LineComment)
		e.event.foot_comment = []byte(node.FootComment)
		e.emit()

	case ScalarNode:
		value := node.Value
		if !utf8.ValidString(value) ***REMOVED***
			if stag == binaryTag ***REMOVED***
				failf("explicitly tagged !!binary data must be base64-encoded")
			***REMOVED***
			if stag != "" ***REMOVED***
				failf("cannot marshal invalid UTF-8 data as %s", stag)
			***REMOVED***
			// It can't be encoded directly as YAML so use a binary tag
			// and encode it as base64.
			tag = binaryTag
			value = encodeBase64(value)
		***REMOVED***

		style := yaml_PLAIN_SCALAR_STYLE
		switch ***REMOVED***
		case node.Style&DoubleQuotedStyle != 0:
			style = yaml_DOUBLE_QUOTED_SCALAR_STYLE
		case node.Style&SingleQuotedStyle != 0:
			style = yaml_SINGLE_QUOTED_SCALAR_STYLE
		case node.Style&LiteralStyle != 0:
			style = yaml_LITERAL_SCALAR_STYLE
		case node.Style&FoldedStyle != 0:
			style = yaml_FOLDED_SCALAR_STYLE
		case strings.Contains(value, "\n"):
			style = yaml_LITERAL_SCALAR_STYLE
		case forceQuoting:
			style = yaml_DOUBLE_QUOTED_SCALAR_STYLE
		***REMOVED***

		e.emitScalar(value, node.Anchor, tag, style, []byte(node.HeadComment), []byte(node.LineComment), []byte(node.FootComment), []byte(tail))
	default:
		failf("cannot encode node with unknown kind %d", node.Kind)
	***REMOVED***
***REMOVED***
