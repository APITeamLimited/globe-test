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
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
	"time"
)

// ----------------------------------------------------------------------------
// Parser, produces a node tree out of a libyaml event stream.

type parser struct ***REMOVED***
	parser   yaml_parser_t
	event    yaml_event_t
	doc      *Node
	anchors  map[string]*Node
	doneInit bool
***REMOVED***

func newParser(b []byte) *parser ***REMOVED***
	p := parser***REMOVED******REMOVED***
	if !yaml_parser_initialize(&p.parser) ***REMOVED***
		panic("failed to initialize YAML emitter")
	***REMOVED***
	if len(b) == 0 ***REMOVED***
		b = []byte***REMOVED***'\n'***REMOVED***
	***REMOVED***
	yaml_parser_set_input_string(&p.parser, b)
	return &p
***REMOVED***

func newParserFromReader(r io.Reader) *parser ***REMOVED***
	p := parser***REMOVED******REMOVED***
	if !yaml_parser_initialize(&p.parser) ***REMOVED***
		panic("failed to initialize YAML emitter")
	***REMOVED***
	yaml_parser_set_input_reader(&p.parser, r)
	return &p
***REMOVED***

func (p *parser) init() ***REMOVED***
	if p.doneInit ***REMOVED***
		return
	***REMOVED***
	p.anchors = make(map[string]*Node)
	p.expect(yaml_STREAM_START_EVENT)
	p.doneInit = true
***REMOVED***

func (p *parser) destroy() ***REMOVED***
	if p.event.typ != yaml_NO_EVENT ***REMOVED***
		yaml_event_delete(&p.event)
	***REMOVED***
	yaml_parser_delete(&p.parser)
***REMOVED***

// expect consumes an event from the event stream and
// checks that it's of the expected type.
func (p *parser) expect(e yaml_event_type_t) ***REMOVED***
	if p.event.typ == yaml_NO_EVENT ***REMOVED***
		if !yaml_parser_parse(&p.parser, &p.event) ***REMOVED***
			p.fail()
		***REMOVED***
	***REMOVED***
	if p.event.typ == yaml_STREAM_END_EVENT ***REMOVED***
		failf("attempted to go past the end of stream; corrupted value?")
	***REMOVED***
	if p.event.typ != e ***REMOVED***
		p.parser.problem = fmt.Sprintf("expected %s event but got %s", e, p.event.typ)
		p.fail()
	***REMOVED***
	yaml_event_delete(&p.event)
	p.event.typ = yaml_NO_EVENT
***REMOVED***

// peek peeks at the next event in the event stream,
// puts the results into p.event and returns the event type.
func (p *parser) peek() yaml_event_type_t ***REMOVED***
	if p.event.typ != yaml_NO_EVENT ***REMOVED***
		return p.event.typ
	***REMOVED***
	if !yaml_parser_parse(&p.parser, &p.event) ***REMOVED***
		p.fail()
	***REMOVED***
	return p.event.typ
***REMOVED***

func (p *parser) fail() ***REMOVED***
	var where string
	var line int
	if p.parser.problem_mark.line != 0 ***REMOVED***
		line = p.parser.problem_mark.line
		// Scanner errors don't iterate line before returning error
		if p.parser.error == yaml_SCANNER_ERROR ***REMOVED***
			line++
		***REMOVED***
	***REMOVED*** else if p.parser.context_mark.line != 0 ***REMOVED***
		line = p.parser.context_mark.line
	***REMOVED***
	if line != 0 ***REMOVED***
		where = "line " + strconv.Itoa(line) + ": "
	***REMOVED***
	var msg string
	if len(p.parser.problem) > 0 ***REMOVED***
		msg = p.parser.problem
	***REMOVED*** else ***REMOVED***
		msg = "unknown problem parsing YAML content"
	***REMOVED***
	failf("%s%s", where, msg)
***REMOVED***

func (p *parser) anchor(n *Node, anchor []byte) ***REMOVED***
	if anchor != nil ***REMOVED***
		n.Anchor = string(anchor)
		p.anchors[n.Anchor] = n
	***REMOVED***
***REMOVED***

func (p *parser) parse() *Node ***REMOVED***
	p.init()
	switch p.peek() ***REMOVED***
	case yaml_SCALAR_EVENT:
		return p.scalar()
	case yaml_ALIAS_EVENT:
		return p.alias()
	case yaml_MAPPING_START_EVENT:
		return p.mapping()
	case yaml_SEQUENCE_START_EVENT:
		return p.sequence()
	case yaml_DOCUMENT_START_EVENT:
		return p.document()
	case yaml_STREAM_END_EVENT:
		// Happens when attempting to decode an empty buffer.
		return nil
	case yaml_TAIL_COMMENT_EVENT:
		panic("internal error: unexpected tail comment event (please report)")
	default:
		panic("internal error: attempted to parse unknown event (please report): " + p.event.typ.String())
	***REMOVED***
***REMOVED***

func (p *parser) node(kind Kind, defaultTag, tag, value string) *Node ***REMOVED***
	var style Style
	if tag != "" && tag != "!" ***REMOVED***
		tag = shortTag(tag)
		style = TaggedStyle
	***REMOVED*** else if defaultTag != "" ***REMOVED***
		tag = defaultTag
	***REMOVED*** else if kind == ScalarNode ***REMOVED***
		tag, _ = resolve("", value)
	***REMOVED***
	return &Node***REMOVED***
		Kind:        kind,
		Tag:         tag,
		Value:       value,
		Style:       style,
		Line:        p.event.start_mark.line + 1,
		Column:      p.event.start_mark.column + 1,
		HeadComment: string(p.event.head_comment),
		LineComment: string(p.event.line_comment),
		FootComment: string(p.event.foot_comment),
	***REMOVED***
***REMOVED***

func (p *parser) parseChild(parent *Node) *Node ***REMOVED***
	child := p.parse()
	parent.Content = append(parent.Content, child)
	return child
***REMOVED***

func (p *parser) document() *Node ***REMOVED***
	n := p.node(DocumentNode, "", "", "")
	p.doc = n
	p.expect(yaml_DOCUMENT_START_EVENT)
	p.parseChild(n)
	if p.peek() == yaml_DOCUMENT_END_EVENT ***REMOVED***
		n.FootComment = string(p.event.foot_comment)
	***REMOVED***
	p.expect(yaml_DOCUMENT_END_EVENT)
	return n
***REMOVED***

func (p *parser) alias() *Node ***REMOVED***
	n := p.node(AliasNode, "", "", string(p.event.anchor))
	n.Alias = p.anchors[n.Value]
	if n.Alias == nil ***REMOVED***
		failf("unknown anchor '%s' referenced", n.Value)
	***REMOVED***
	p.expect(yaml_ALIAS_EVENT)
	return n
***REMOVED***

func (p *parser) scalar() *Node ***REMOVED***
	var parsedStyle = p.event.scalar_style()
	var nodeStyle Style
	switch ***REMOVED***
	case parsedStyle&yaml_DOUBLE_QUOTED_SCALAR_STYLE != 0:
		nodeStyle = DoubleQuotedStyle
	case parsedStyle&yaml_SINGLE_QUOTED_SCALAR_STYLE != 0:
		nodeStyle = SingleQuotedStyle
	case parsedStyle&yaml_LITERAL_SCALAR_STYLE != 0:
		nodeStyle = LiteralStyle
	case parsedStyle&yaml_FOLDED_SCALAR_STYLE != 0:
		nodeStyle = FoldedStyle
	***REMOVED***
	var nodeValue = string(p.event.value)
	var nodeTag = string(p.event.tag)
	var defaultTag string
	if nodeStyle == 0 ***REMOVED***
		if nodeValue == "<<" ***REMOVED***
			defaultTag = mergeTag
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		defaultTag = strTag
	***REMOVED***
	n := p.node(ScalarNode, defaultTag, nodeTag, nodeValue)
	n.Style |= nodeStyle
	p.anchor(n, p.event.anchor)
	p.expect(yaml_SCALAR_EVENT)
	return n
***REMOVED***

func (p *parser) sequence() *Node ***REMOVED***
	n := p.node(SequenceNode, seqTag, string(p.event.tag), "")
	if p.event.sequence_style()&yaml_FLOW_SEQUENCE_STYLE != 0 ***REMOVED***
		n.Style |= FlowStyle
	***REMOVED***
	p.anchor(n, p.event.anchor)
	p.expect(yaml_SEQUENCE_START_EVENT)
	for p.peek() != yaml_SEQUENCE_END_EVENT ***REMOVED***
		p.parseChild(n)
	***REMOVED***
	n.LineComment = string(p.event.line_comment)
	n.FootComment = string(p.event.foot_comment)
	p.expect(yaml_SEQUENCE_END_EVENT)
	return n
***REMOVED***

func (p *parser) mapping() *Node ***REMOVED***
	n := p.node(MappingNode, mapTag, string(p.event.tag), "")
	block := true
	if p.event.mapping_style()&yaml_FLOW_MAPPING_STYLE != 0 ***REMOVED***
		block = false
		n.Style |= FlowStyle
	***REMOVED***
	p.anchor(n, p.event.anchor)
	p.expect(yaml_MAPPING_START_EVENT)
	for p.peek() != yaml_MAPPING_END_EVENT ***REMOVED***
		k := p.parseChild(n)
		if block && k.FootComment != "" ***REMOVED***
			// Must be a foot comment for the prior value when being dedented.
			if len(n.Content) > 2 ***REMOVED***
				n.Content[len(n.Content)-3].FootComment = k.FootComment
				k.FootComment = ""
			***REMOVED***
		***REMOVED***
		v := p.parseChild(n)
		if k.FootComment == "" && v.FootComment != "" ***REMOVED***
			k.FootComment = v.FootComment
			v.FootComment = ""
		***REMOVED***
		if p.peek() == yaml_TAIL_COMMENT_EVENT ***REMOVED***
			if k.FootComment == "" ***REMOVED***
				k.FootComment = string(p.event.foot_comment)
			***REMOVED***
			p.expect(yaml_TAIL_COMMENT_EVENT)
		***REMOVED***
	***REMOVED***
	n.LineComment = string(p.event.line_comment)
	n.FootComment = string(p.event.foot_comment)
	if n.Style&FlowStyle == 0 && n.FootComment != "" && len(n.Content) > 1 ***REMOVED***
		n.Content[len(n.Content)-2].FootComment = n.FootComment
		n.FootComment = ""
	***REMOVED***
	p.expect(yaml_MAPPING_END_EVENT)
	return n
***REMOVED***

// ----------------------------------------------------------------------------
// Decoder, unmarshals a node into a provided value.

type decoder struct ***REMOVED***
	doc     *Node
	aliases map[*Node]bool
	terrors []string

	stringMapType  reflect.Type
	generalMapType reflect.Type

	knownFields bool
	uniqueKeys  bool
	decodeCount int
	aliasCount  int
	aliasDepth  int
***REMOVED***

var (
	nodeType       = reflect.TypeOf(Node***REMOVED******REMOVED***)
	durationType   = reflect.TypeOf(time.Duration(0))
	stringMapType  = reflect.TypeOf(map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***)
	generalMapType = reflect.TypeOf(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED******REMOVED***)
	ifaceType      = generalMapType.Elem()
	timeType       = reflect.TypeOf(time.Time***REMOVED******REMOVED***)
	ptrTimeType    = reflect.TypeOf(&time.Time***REMOVED******REMOVED***)
)

func newDecoder() *decoder ***REMOVED***
	d := &decoder***REMOVED***
		stringMapType:  stringMapType,
		generalMapType: generalMapType,
		uniqueKeys:     true,
	***REMOVED***
	d.aliases = make(map[*Node]bool)
	return d
***REMOVED***

func (d *decoder) terror(n *Node, tag string, out reflect.Value) ***REMOVED***
	if n.Tag != "" ***REMOVED***
		tag = n.Tag
	***REMOVED***
	value := n.Value
	if tag != seqTag && tag != mapTag ***REMOVED***
		if len(value) > 10 ***REMOVED***
			value = " `" + value[:7] + "...`"
		***REMOVED*** else ***REMOVED***
			value = " `" + value + "`"
		***REMOVED***
	***REMOVED***
	d.terrors = append(d.terrors, fmt.Sprintf("line %d: cannot unmarshal %s%s into %s", n.Line, shortTag(tag), value, out.Type()))
***REMOVED***

func (d *decoder) callUnmarshaler(n *Node, u Unmarshaler) (good bool) ***REMOVED***
	err := u.UnmarshalYAML(n)
	if e, ok := err.(*TypeError); ok ***REMOVED***
		d.terrors = append(d.terrors, e.Errors...)
		return false
	***REMOVED***
	if err != nil ***REMOVED***
		fail(err)
	***REMOVED***
	return true
***REMOVED***

func (d *decoder) callObsoleteUnmarshaler(n *Node, u obsoleteUnmarshaler) (good bool) ***REMOVED***
	terrlen := len(d.terrors)
	err := u.UnmarshalYAML(func(v interface***REMOVED******REMOVED***) (err error) ***REMOVED***
		defer handleErr(&err)
		d.unmarshal(n, reflect.ValueOf(v))
		if len(d.terrors) > terrlen ***REMOVED***
			issues := d.terrors[terrlen:]
			d.terrors = d.terrors[:terrlen]
			return &TypeError***REMOVED***issues***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***)
	if e, ok := err.(*TypeError); ok ***REMOVED***
		d.terrors = append(d.terrors, e.Errors...)
		return false
	***REMOVED***
	if err != nil ***REMOVED***
		fail(err)
	***REMOVED***
	return true
***REMOVED***

// d.prepare initializes and dereferences pointers and calls UnmarshalYAML
// if a value is found to implement it.
// It returns the initialized and dereferenced out value, whether
// unmarshalling was already done by UnmarshalYAML, and if so whether
// its types unmarshalled appropriately.
//
// If n holds a null value, prepare returns before doing anything.
func (d *decoder) prepare(n *Node, out reflect.Value) (newout reflect.Value, unmarshaled, good bool) ***REMOVED***
	if n.ShortTag() == nullTag ***REMOVED***
		return out, false, false
	***REMOVED***
	again := true
	for again ***REMOVED***
		again = false
		if out.Kind() == reflect.Ptr ***REMOVED***
			if out.IsNil() ***REMOVED***
				out.Set(reflect.New(out.Type().Elem()))
			***REMOVED***
			out = out.Elem()
			again = true
		***REMOVED***
		if out.CanAddr() ***REMOVED***
			outi := out.Addr().Interface()
			if u, ok := outi.(Unmarshaler); ok ***REMOVED***
				good = d.callUnmarshaler(n, u)
				return out, true, good
			***REMOVED***
			if u, ok := outi.(obsoleteUnmarshaler); ok ***REMOVED***
				good = d.callObsoleteUnmarshaler(n, u)
				return out, true, good
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return out, false, false
***REMOVED***

func (d *decoder) fieldByIndex(n *Node, v reflect.Value, index []int) (field reflect.Value) ***REMOVED***
	if n.ShortTag() == nullTag ***REMOVED***
		return reflect.Value***REMOVED******REMOVED***
	***REMOVED***
	for _, num := range index ***REMOVED***
		for ***REMOVED***
			if v.Kind() == reflect.Ptr ***REMOVED***
				if v.IsNil() ***REMOVED***
					v.Set(reflect.New(v.Type().Elem()))
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

const (
	// 400,000 decode operations is ~500kb of dense object declarations, or
	// ~5kb of dense object declarations with 10000% alias expansion
	alias_ratio_range_low = 400000

	// 4,000,000 decode operations is ~5MB of dense object declarations, or
	// ~4.5MB of dense object declarations with 10% alias expansion
	alias_ratio_range_high = 4000000

	// alias_ratio_range is the range over which we scale allowed alias ratios
	alias_ratio_range = float64(alias_ratio_range_high - alias_ratio_range_low)
)

func allowedAliasRatio(decodeCount int) float64 ***REMOVED***
	switch ***REMOVED***
	case decodeCount <= alias_ratio_range_low:
		// allow 99% to come from alias expansion for small-to-medium documents
		return 0.99
	case decodeCount >= alias_ratio_range_high:
		// allow 10% to come from alias expansion for very large documents
		return 0.10
	default:
		// scale smoothly from 99% down to 10% over the range.
		// this maps to 396,000 - 400,000 allowed alias-driven decodes over the range.
		// 400,000 decode operations is ~100MB of allocations in worst-case scenarios (single-item maps).
		return 0.99 - 0.89*(float64(decodeCount-alias_ratio_range_low)/alias_ratio_range)
	***REMOVED***
***REMOVED***

func (d *decoder) unmarshal(n *Node, out reflect.Value) (good bool) ***REMOVED***
	d.decodeCount++
	if d.aliasDepth > 0 ***REMOVED***
		d.aliasCount++
	***REMOVED***
	if d.aliasCount > 100 && d.decodeCount > 1000 && float64(d.aliasCount)/float64(d.decodeCount) > allowedAliasRatio(d.decodeCount) ***REMOVED***
		failf("document contains excessive aliasing")
	***REMOVED***
	if out.Type() == nodeType ***REMOVED***
		out.Set(reflect.ValueOf(n).Elem())
		return true
	***REMOVED***
	switch n.Kind ***REMOVED***
	case DocumentNode:
		return d.document(n, out)
	case AliasNode:
		return d.alias(n, out)
	***REMOVED***
	out, unmarshaled, good := d.prepare(n, out)
	if unmarshaled ***REMOVED***
		return good
	***REMOVED***
	switch n.Kind ***REMOVED***
	case ScalarNode:
		good = d.scalar(n, out)
	case MappingNode:
		good = d.mapping(n, out)
	case SequenceNode:
		good = d.sequence(n, out)
	default:
		panic("internal error: unknown node kind: " + strconv.Itoa(int(n.Kind)))
	***REMOVED***
	return good
***REMOVED***

func (d *decoder) document(n *Node, out reflect.Value) (good bool) ***REMOVED***
	if len(n.Content) == 1 ***REMOVED***
		d.doc = n
		d.unmarshal(n.Content[0], out)
		return true
	***REMOVED***
	return false
***REMOVED***

func (d *decoder) alias(n *Node, out reflect.Value) (good bool) ***REMOVED***
	if d.aliases[n] ***REMOVED***
		// TODO this could actually be allowed in some circumstances.
		failf("anchor '%s' value contains itself", n.Value)
	***REMOVED***
	d.aliases[n] = true
	d.aliasDepth++
	good = d.unmarshal(n.Alias, out)
	d.aliasDepth--
	delete(d.aliases, n)
	return good
***REMOVED***

var zeroValue reflect.Value

func resetMap(out reflect.Value) ***REMOVED***
	for _, k := range out.MapKeys() ***REMOVED***
		out.SetMapIndex(k, zeroValue)
	***REMOVED***
***REMOVED***

func (d *decoder) scalar(n *Node, out reflect.Value) bool ***REMOVED***
	var tag string
	var resolved interface***REMOVED******REMOVED***
	if n.indicatedString() ***REMOVED***
		tag = strTag
		resolved = n.Value
	***REMOVED*** else ***REMOVED***
		tag, resolved = resolve(n.Tag, n.Value)
		if tag == binaryTag ***REMOVED***
			data, err := base64.StdEncoding.DecodeString(resolved.(string))
			if err != nil ***REMOVED***
				failf("!!binary value contains invalid base64 data")
			***REMOVED***
			resolved = string(data)
		***REMOVED***
	***REMOVED***
	if resolved == nil ***REMOVED***
		if out.CanAddr() ***REMOVED***
			switch out.Kind() ***REMOVED***
			case reflect.Interface, reflect.Ptr, reflect.Map, reflect.Slice:
				out.Set(reflect.Zero(out.Type()))
				return true
			***REMOVED***
		***REMOVED***
		return false
	***REMOVED***
	if resolvedv := reflect.ValueOf(resolved); out.Type() == resolvedv.Type() ***REMOVED***
		// We've resolved to exactly the type we want, so use that.
		out.Set(resolvedv)
		return true
	***REMOVED***
	// Perhaps we can use the value as a TextUnmarshaler to
	// set its value.
	if out.CanAddr() ***REMOVED***
		u, ok := out.Addr().Interface().(encoding.TextUnmarshaler)
		if ok ***REMOVED***
			var text []byte
			if tag == binaryTag ***REMOVED***
				text = []byte(resolved.(string))
			***REMOVED*** else ***REMOVED***
				// We let any value be unmarshaled into TextUnmarshaler.
				// That might be more lax than we'd like, but the
				// TextUnmarshaler itself should bowl out any dubious values.
				text = []byte(n.Value)
			***REMOVED***
			err := u.UnmarshalText(text)
			if err != nil ***REMOVED***
				fail(err)
			***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	switch out.Kind() ***REMOVED***
	case reflect.String:
		if tag == binaryTag ***REMOVED***
			out.SetString(resolved.(string))
			return true
		***REMOVED***
		out.SetString(n.Value)
		return true
	case reflect.Interface:
		out.Set(reflect.ValueOf(resolved))
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// This used to work in v2, but it's very unfriendly.
		isDuration := out.Type() == durationType

		switch resolved := resolved.(type) ***REMOVED***
		case int:
			if !isDuration && !out.OverflowInt(int64(resolved)) ***REMOVED***
				out.SetInt(int64(resolved))
				return true
			***REMOVED***
		case int64:
			if !isDuration && !out.OverflowInt(resolved) ***REMOVED***
				out.SetInt(resolved)
				return true
			***REMOVED***
		case uint64:
			if !isDuration && resolved <= math.MaxInt64 && !out.OverflowInt(int64(resolved)) ***REMOVED***
				out.SetInt(int64(resolved))
				return true
			***REMOVED***
		case float64:
			if !isDuration && resolved <= math.MaxInt64 && !out.OverflowInt(int64(resolved)) ***REMOVED***
				out.SetInt(int64(resolved))
				return true
			***REMOVED***
		case string:
			if out.Type() == durationType ***REMOVED***
				d, err := time.ParseDuration(resolved)
				if err == nil ***REMOVED***
					out.SetInt(int64(d))
					return true
				***REMOVED***
			***REMOVED***
		***REMOVED***
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch resolved := resolved.(type) ***REMOVED***
		case int:
			if resolved >= 0 && !out.OverflowUint(uint64(resolved)) ***REMOVED***
				out.SetUint(uint64(resolved))
				return true
			***REMOVED***
		case int64:
			if resolved >= 0 && !out.OverflowUint(uint64(resolved)) ***REMOVED***
				out.SetUint(uint64(resolved))
				return true
			***REMOVED***
		case uint64:
			if !out.OverflowUint(uint64(resolved)) ***REMOVED***
				out.SetUint(uint64(resolved))
				return true
			***REMOVED***
		case float64:
			if resolved <= math.MaxUint64 && !out.OverflowUint(uint64(resolved)) ***REMOVED***
				out.SetUint(uint64(resolved))
				return true
			***REMOVED***
		***REMOVED***
	case reflect.Bool:
		switch resolved := resolved.(type) ***REMOVED***
		case bool:
			out.SetBool(resolved)
			return true
		case string:
			// This offers some compatibility with the 1.1 spec (https://yaml.org/type/bool.html).
			// It only works if explicitly attempting to unmarshal into a typed bool value.
			switch resolved ***REMOVED***
			case "y", "Y", "yes", "Yes", "YES", "on", "On", "ON":
				out.SetBool(true)
				return true
			case "n", "N", "no", "No", "NO", "off", "Off", "OFF":
				out.SetBool(false)
				return true
			***REMOVED***
		***REMOVED***
	case reflect.Float32, reflect.Float64:
		switch resolved := resolved.(type) ***REMOVED***
		case int:
			out.SetFloat(float64(resolved))
			return true
		case int64:
			out.SetFloat(float64(resolved))
			return true
		case uint64:
			out.SetFloat(float64(resolved))
			return true
		case float64:
			out.SetFloat(resolved)
			return true
		***REMOVED***
	case reflect.Struct:
		if resolvedv := reflect.ValueOf(resolved); out.Type() == resolvedv.Type() ***REMOVED***
			out.Set(resolvedv)
			return true
		***REMOVED***
	case reflect.Ptr:
		panic("yaml internal error: please report the issue")
	***REMOVED***
	d.terror(n, tag, out)
	return false
***REMOVED***

func settableValueOf(i interface***REMOVED******REMOVED***) reflect.Value ***REMOVED***
	v := reflect.ValueOf(i)
	sv := reflect.New(v.Type()).Elem()
	sv.Set(v)
	return sv
***REMOVED***

func (d *decoder) sequence(n *Node, out reflect.Value) (good bool) ***REMOVED***
	l := len(n.Content)

	var iface reflect.Value
	switch out.Kind() ***REMOVED***
	case reflect.Slice:
		out.Set(reflect.MakeSlice(out.Type(), l, l))
	case reflect.Array:
		if l != out.Len() ***REMOVED***
			failf("invalid array: want %d elements but got %d", out.Len(), l)
		***REMOVED***
	case reflect.Interface:
		// No type hints. Will have to use a generic sequence.
		iface = out
		out = settableValueOf(make([]interface***REMOVED******REMOVED***, l))
	default:
		d.terror(n, seqTag, out)
		return false
	***REMOVED***
	et := out.Type().Elem()

	j := 0
	for i := 0; i < l; i++ ***REMOVED***
		e := reflect.New(et).Elem()
		if ok := d.unmarshal(n.Content[i], e); ok ***REMOVED***
			out.Index(j).Set(e)
			j++
		***REMOVED***
	***REMOVED***
	if out.Kind() != reflect.Array ***REMOVED***
		out.Set(out.Slice(0, j))
	***REMOVED***
	if iface.IsValid() ***REMOVED***
		iface.Set(out)
	***REMOVED***
	return true
***REMOVED***

func (d *decoder) mapping(n *Node, out reflect.Value) (good bool) ***REMOVED***
	l := len(n.Content)
	if d.uniqueKeys ***REMOVED***
		nerrs := len(d.terrors)
		for i := 0; i < l; i += 2 ***REMOVED***
			ni := n.Content[i]
			for j := i + 2; j < l; j += 2 ***REMOVED***
				nj := n.Content[j]
				if ni.Kind == nj.Kind && ni.Value == nj.Value ***REMOVED***
					d.terrors = append(d.terrors, fmt.Sprintf("line %d: mapping key %#v already defined at line %d", nj.Line, nj.Value, ni.Line))
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if len(d.terrors) > nerrs ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	switch out.Kind() ***REMOVED***
	case reflect.Struct:
		return d.mappingStruct(n, out)
	case reflect.Map:
		// okay
	case reflect.Interface:
		iface := out
		if isStringMap(n) ***REMOVED***
			out = reflect.MakeMap(d.stringMapType)
		***REMOVED*** else ***REMOVED***
			out = reflect.MakeMap(d.generalMapType)
		***REMOVED***
		iface.Set(out)
	default:
		d.terror(n, mapTag, out)
		return false
	***REMOVED***

	outt := out.Type()
	kt := outt.Key()
	et := outt.Elem()

	stringMapType := d.stringMapType
	generalMapType := d.generalMapType
	if outt.Elem() == ifaceType ***REMOVED***
		if outt.Key().Kind() == reflect.String ***REMOVED***
			d.stringMapType = outt
		***REMOVED*** else if outt.Key() == ifaceType ***REMOVED***
			d.generalMapType = outt
		***REMOVED***
	***REMOVED***

	if out.IsNil() ***REMOVED***
		out.Set(reflect.MakeMap(outt))
	***REMOVED***
	for i := 0; i < l; i += 2 ***REMOVED***
		if isMerge(n.Content[i]) ***REMOVED***
			d.merge(n.Content[i+1], out)
			continue
		***REMOVED***
		k := reflect.New(kt).Elem()
		if d.unmarshal(n.Content[i], k) ***REMOVED***
			kkind := k.Kind()
			if kkind == reflect.Interface ***REMOVED***
				kkind = k.Elem().Kind()
			***REMOVED***
			if kkind == reflect.Map || kkind == reflect.Slice ***REMOVED***
				failf("invalid map key: %#v", k.Interface())
			***REMOVED***
			e := reflect.New(et).Elem()
			if d.unmarshal(n.Content[i+1], e) ***REMOVED***
				out.SetMapIndex(k, e)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	d.stringMapType = stringMapType
	d.generalMapType = generalMapType
	return true
***REMOVED***

func isStringMap(n *Node) bool ***REMOVED***
	if n.Kind != MappingNode ***REMOVED***
		return false
	***REMOVED***
	l := len(n.Content)
	for i := 0; i < l; i += 2 ***REMOVED***
		if n.Content[i].ShortTag() != strTag ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (d *decoder) mappingStruct(n *Node, out reflect.Value) (good bool) ***REMOVED***
	sinfo, err := getStructInfo(out.Type())
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	var inlineMap reflect.Value
	var elemType reflect.Type
	if sinfo.InlineMap != -1 ***REMOVED***
		inlineMap = out.Field(sinfo.InlineMap)
		inlineMap.Set(reflect.New(inlineMap.Type()).Elem())
		elemType = inlineMap.Type().Elem()
	***REMOVED***

	for _, index := range sinfo.InlineUnmarshalers ***REMOVED***
		field := d.fieldByIndex(n, out, index)
		d.prepare(n, field)
	***REMOVED***

	var doneFields []bool
	if d.uniqueKeys ***REMOVED***
		doneFields = make([]bool, len(sinfo.FieldsList))
	***REMOVED***
	name := settableValueOf("")
	l := len(n.Content)
	for i := 0; i < l; i += 2 ***REMOVED***
		ni := n.Content[i]
		if isMerge(ni) ***REMOVED***
			d.merge(n.Content[i+1], out)
			continue
		***REMOVED***
		if !d.unmarshal(ni, name) ***REMOVED***
			continue
		***REMOVED***
		if info, ok := sinfo.FieldsMap[name.String()]; ok ***REMOVED***
			if d.uniqueKeys ***REMOVED***
				if doneFields[info.Id] ***REMOVED***
					d.terrors = append(d.terrors, fmt.Sprintf("line %d: field %s already set in type %s", ni.Line, name.String(), out.Type()))
					continue
				***REMOVED***
				doneFields[info.Id] = true
			***REMOVED***
			var field reflect.Value
			if info.Inline == nil ***REMOVED***
				field = out.Field(info.Num)
			***REMOVED*** else ***REMOVED***
				field = d.fieldByIndex(n, out, info.Inline)
			***REMOVED***
			d.unmarshal(n.Content[i+1], field)
		***REMOVED*** else if sinfo.InlineMap != -1 ***REMOVED***
			if inlineMap.IsNil() ***REMOVED***
				inlineMap.Set(reflect.MakeMap(inlineMap.Type()))
			***REMOVED***
			value := reflect.New(elemType).Elem()
			d.unmarshal(n.Content[i+1], value)
			inlineMap.SetMapIndex(name, value)
		***REMOVED*** else if d.knownFields ***REMOVED***
			d.terrors = append(d.terrors, fmt.Sprintf("line %d: field %s not found in type %s", ni.Line, name.String(), out.Type()))
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func failWantMap() ***REMOVED***
	failf("map merge requires map or sequence of maps as the value")
***REMOVED***

func (d *decoder) merge(n *Node, out reflect.Value) ***REMOVED***
	switch n.Kind ***REMOVED***
	case MappingNode:
		d.unmarshal(n, out)
	case AliasNode:
		if n.Alias != nil && n.Alias.Kind != MappingNode ***REMOVED***
			failWantMap()
		***REMOVED***
		d.unmarshal(n, out)
	case SequenceNode:
		// Step backwards as earlier nodes take precedence.
		for i := len(n.Content) - 1; i >= 0; i-- ***REMOVED***
			ni := n.Content[i]
			if ni.Kind == AliasNode ***REMOVED***
				if ni.Alias != nil && ni.Alias.Kind != MappingNode ***REMOVED***
					failWantMap()
				***REMOVED***
			***REMOVED*** else if ni.Kind != MappingNode ***REMOVED***
				failWantMap()
			***REMOVED***
			d.unmarshal(ni, out)
		***REMOVED***
	default:
		failWantMap()
	***REMOVED***
***REMOVED***

func isMerge(n *Node) bool ***REMOVED***
	return n.Kind == ScalarNode && n.Value == "<<" && (n.Tag == "" || n.Tag == "!" || shortTag(n.Tag) == mergeTag)
***REMOVED***
