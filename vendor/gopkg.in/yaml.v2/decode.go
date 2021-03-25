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

const (
	documentNode = 1 << iota
	mappingNode
	sequenceNode
	scalarNode
	aliasNode
)

type node struct ***REMOVED***
	kind         int
	line, column int
	tag          string
	// For an alias node, alias holds the resolved alias.
	alias    *node
	value    string
	implicit bool
	children []*node
	anchors  map[string]*node
***REMOVED***

// ----------------------------------------------------------------------------
// Parser, produces a node tree out of a libyaml event stream.

type parser struct ***REMOVED***
	parser   yaml_parser_t
	event    yaml_event_t
	doc      *node
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

func (p *parser) anchor(n *node, anchor []byte) ***REMOVED***
	if anchor != nil ***REMOVED***
		p.doc.anchors[string(anchor)] = n
	***REMOVED***
***REMOVED***

func (p *parser) parse() *node ***REMOVED***
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
	default:
		panic("attempted to parse unknown event: " + p.event.typ.String())
	***REMOVED***
***REMOVED***

func (p *parser) node(kind int) *node ***REMOVED***
	return &node***REMOVED***
		kind:   kind,
		line:   p.event.start_mark.line,
		column: p.event.start_mark.column,
	***REMOVED***
***REMOVED***

func (p *parser) document() *node ***REMOVED***
	n := p.node(documentNode)
	n.anchors = make(map[string]*node)
	p.doc = n
	p.expect(yaml_DOCUMENT_START_EVENT)
	n.children = append(n.children, p.parse())
	p.expect(yaml_DOCUMENT_END_EVENT)
	return n
***REMOVED***

func (p *parser) alias() *node ***REMOVED***
	n := p.node(aliasNode)
	n.value = string(p.event.anchor)
	n.alias = p.doc.anchors[n.value]
	if n.alias == nil ***REMOVED***
		failf("unknown anchor '%s' referenced", n.value)
	***REMOVED***
	p.expect(yaml_ALIAS_EVENT)
	return n
***REMOVED***

func (p *parser) scalar() *node ***REMOVED***
	n := p.node(scalarNode)
	n.value = string(p.event.value)
	n.tag = string(p.event.tag)
	n.implicit = p.event.implicit
	p.anchor(n, p.event.anchor)
	p.expect(yaml_SCALAR_EVENT)
	return n
***REMOVED***

func (p *parser) sequence() *node ***REMOVED***
	n := p.node(sequenceNode)
	p.anchor(n, p.event.anchor)
	p.expect(yaml_SEQUENCE_START_EVENT)
	for p.peek() != yaml_SEQUENCE_END_EVENT ***REMOVED***
		n.children = append(n.children, p.parse())
	***REMOVED***
	p.expect(yaml_SEQUENCE_END_EVENT)
	return n
***REMOVED***

func (p *parser) mapping() *node ***REMOVED***
	n := p.node(mappingNode)
	p.anchor(n, p.event.anchor)
	p.expect(yaml_MAPPING_START_EVENT)
	for p.peek() != yaml_MAPPING_END_EVENT ***REMOVED***
		n.children = append(n.children, p.parse(), p.parse())
	***REMOVED***
	p.expect(yaml_MAPPING_END_EVENT)
	return n
***REMOVED***

// ----------------------------------------------------------------------------
// Decoder, unmarshals a node into a provided value.

type decoder struct ***REMOVED***
	doc     *node
	aliases map[*node]bool
	mapType reflect.Type
	terrors []string
	strict  bool

	decodeCount int
	aliasCount  int
	aliasDepth  int
***REMOVED***

var (
	mapItemType    = reflect.TypeOf(MapItem***REMOVED******REMOVED***)
	durationType   = reflect.TypeOf(time.Duration(0))
	defaultMapType = reflect.TypeOf(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED******REMOVED***)
	ifaceType      = defaultMapType.Elem()
	timeType       = reflect.TypeOf(time.Time***REMOVED******REMOVED***)
	ptrTimeType    = reflect.TypeOf(&time.Time***REMOVED******REMOVED***)
)

func newDecoder(strict bool) *decoder ***REMOVED***
	d := &decoder***REMOVED***mapType: defaultMapType, strict: strict***REMOVED***
	d.aliases = make(map[*node]bool)
	return d
***REMOVED***

func (d *decoder) terror(n *node, tag string, out reflect.Value) ***REMOVED***
	if n.tag != "" ***REMOVED***
		tag = n.tag
	***REMOVED***
	value := n.value
	if tag != yaml_SEQ_TAG && tag != yaml_MAP_TAG ***REMOVED***
		if len(value) > 10 ***REMOVED***
			value = " `" + value[:7] + "...`"
		***REMOVED*** else ***REMOVED***
			value = " `" + value + "`"
		***REMOVED***
	***REMOVED***
	d.terrors = append(d.terrors, fmt.Sprintf("line %d: cannot unmarshal %s%s into %s", n.line+1, shortTag(tag), value, out.Type()))
***REMOVED***

func (d *decoder) callUnmarshaler(n *node, u Unmarshaler) (good bool) ***REMOVED***
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
func (d *decoder) prepare(n *node, out reflect.Value) (newout reflect.Value, unmarshaled, good bool) ***REMOVED***
	if n.tag == yaml_NULL_TAG || n.kind == scalarNode && n.tag == "" && (n.value == "null" || n.value == "~" || n.value == "" && n.implicit) ***REMOVED***
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
			if u, ok := out.Addr().Interface().(Unmarshaler); ok ***REMOVED***
				good = d.callUnmarshaler(n, u)
				return out, true, good
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return out, false, false
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

func (d *decoder) unmarshal(n *node, out reflect.Value) (good bool) ***REMOVED***
	d.decodeCount++
	if d.aliasDepth > 0 ***REMOVED***
		d.aliasCount++
	***REMOVED***
	if d.aliasCount > 100 && d.decodeCount > 1000 && float64(d.aliasCount)/float64(d.decodeCount) > allowedAliasRatio(d.decodeCount) ***REMOVED***
		failf("document contains excessive aliasing")
	***REMOVED***
	switch n.kind ***REMOVED***
	case documentNode:
		return d.document(n, out)
	case aliasNode:
		return d.alias(n, out)
	***REMOVED***
	out, unmarshaled, good := d.prepare(n, out)
	if unmarshaled ***REMOVED***
		return good
	***REMOVED***
	switch n.kind ***REMOVED***
	case scalarNode:
		good = d.scalar(n, out)
	case mappingNode:
		good = d.mapping(n, out)
	case sequenceNode:
		good = d.sequence(n, out)
	default:
		panic("internal error: unknown node kind: " + strconv.Itoa(n.kind))
	***REMOVED***
	return good
***REMOVED***

func (d *decoder) document(n *node, out reflect.Value) (good bool) ***REMOVED***
	if len(n.children) == 1 ***REMOVED***
		d.doc = n
		d.unmarshal(n.children[0], out)
		return true
	***REMOVED***
	return false
***REMOVED***

func (d *decoder) alias(n *node, out reflect.Value) (good bool) ***REMOVED***
	if d.aliases[n] ***REMOVED***
		// TODO this could actually be allowed in some circumstances.
		failf("anchor '%s' value contains itself", n.value)
	***REMOVED***
	d.aliases[n] = true
	d.aliasDepth++
	good = d.unmarshal(n.alias, out)
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

func (d *decoder) scalar(n *node, out reflect.Value) bool ***REMOVED***
	var tag string
	var resolved interface***REMOVED******REMOVED***
	if n.tag == "" && !n.implicit ***REMOVED***
		tag = yaml_STR_TAG
		resolved = n.value
	***REMOVED*** else ***REMOVED***
		tag, resolved = resolve(n.tag, n.value)
		if tag == yaml_BINARY_TAG ***REMOVED***
			data, err := base64.StdEncoding.DecodeString(resolved.(string))
			if err != nil ***REMOVED***
				failf("!!binary value contains invalid base64 data")
			***REMOVED***
			resolved = string(data)
		***REMOVED***
	***REMOVED***
	if resolved == nil ***REMOVED***
		if out.Kind() == reflect.Map && !out.CanAddr() ***REMOVED***
			resetMap(out)
		***REMOVED*** else ***REMOVED***
			out.Set(reflect.Zero(out.Type()))
		***REMOVED***
		return true
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
			if tag == yaml_BINARY_TAG ***REMOVED***
				text = []byte(resolved.(string))
			***REMOVED*** else ***REMOVED***
				// We let any value be unmarshaled into TextUnmarshaler.
				// That might be more lax than we'd like, but the
				// TextUnmarshaler itself should bowl out any dubious values.
				text = []byte(n.value)
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
		if tag == yaml_BINARY_TAG ***REMOVED***
			out.SetString(resolved.(string))
			return true
		***REMOVED***
		if resolved != nil ***REMOVED***
			out.SetString(n.value)
			return true
		***REMOVED***
	case reflect.Interface:
		if resolved == nil ***REMOVED***
			out.Set(reflect.Zero(out.Type()))
		***REMOVED*** else if tag == yaml_TIMESTAMP_TAG ***REMOVED***
			// It looks like a timestamp but for backward compatibility
			// reasons we set it as a string, so that code that unmarshals
			// timestamp-like values into interface***REMOVED******REMOVED*** will continue to
			// see a string and not a time.Time.
			// TODO(v3) Drop this.
			out.Set(reflect.ValueOf(n.value))
		***REMOVED*** else ***REMOVED***
			out.Set(reflect.ValueOf(resolved))
		***REMOVED***
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch resolved := resolved.(type) ***REMOVED***
		case int:
			if !out.OverflowInt(int64(resolved)) ***REMOVED***
				out.SetInt(int64(resolved))
				return true
			***REMOVED***
		case int64:
			if !out.OverflowInt(resolved) ***REMOVED***
				out.SetInt(resolved)
				return true
			***REMOVED***
		case uint64:
			if resolved <= math.MaxInt64 && !out.OverflowInt(int64(resolved)) ***REMOVED***
				out.SetInt(int64(resolved))
				return true
			***REMOVED***
		case float64:
			if resolved <= math.MaxInt64 && !out.OverflowInt(int64(resolved)) ***REMOVED***
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
		if out.Type().Elem() == reflect.TypeOf(resolved) ***REMOVED***
			// TODO DOes this make sense? When is out a Ptr except when decoding a nil value?
			elem := reflect.New(out.Type().Elem())
			elem.Elem().Set(reflect.ValueOf(resolved))
			out.Set(elem)
			return true
		***REMOVED***
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

func (d *decoder) sequence(n *node, out reflect.Value) (good bool) ***REMOVED***
	l := len(n.children)

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
		d.terror(n, yaml_SEQ_TAG, out)
		return false
	***REMOVED***
	et := out.Type().Elem()

	j := 0
	for i := 0; i < l; i++ ***REMOVED***
		e := reflect.New(et).Elem()
		if ok := d.unmarshal(n.children[i], e); ok ***REMOVED***
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

func (d *decoder) mapping(n *node, out reflect.Value) (good bool) ***REMOVED***
	switch out.Kind() ***REMOVED***
	case reflect.Struct:
		return d.mappingStruct(n, out)
	case reflect.Slice:
		return d.mappingSlice(n, out)
	case reflect.Map:
		// okay
	case reflect.Interface:
		if d.mapType.Kind() == reflect.Map ***REMOVED***
			iface := out
			out = reflect.MakeMap(d.mapType)
			iface.Set(out)
		***REMOVED*** else ***REMOVED***
			slicev := reflect.New(d.mapType).Elem()
			if !d.mappingSlice(n, slicev) ***REMOVED***
				return false
			***REMOVED***
			out.Set(slicev)
			return true
		***REMOVED***
	default:
		d.terror(n, yaml_MAP_TAG, out)
		return false
	***REMOVED***
	outt := out.Type()
	kt := outt.Key()
	et := outt.Elem()

	mapType := d.mapType
	if outt.Key() == ifaceType && outt.Elem() == ifaceType ***REMOVED***
		d.mapType = outt
	***REMOVED***

	if out.IsNil() ***REMOVED***
		out.Set(reflect.MakeMap(outt))
	***REMOVED***
	l := len(n.children)
	for i := 0; i < l; i += 2 ***REMOVED***
		if isMerge(n.children[i]) ***REMOVED***
			d.merge(n.children[i+1], out)
			continue
		***REMOVED***
		k := reflect.New(kt).Elem()
		if d.unmarshal(n.children[i], k) ***REMOVED***
			kkind := k.Kind()
			if kkind == reflect.Interface ***REMOVED***
				kkind = k.Elem().Kind()
			***REMOVED***
			if kkind == reflect.Map || kkind == reflect.Slice ***REMOVED***
				failf("invalid map key: %#v", k.Interface())
			***REMOVED***
			e := reflect.New(et).Elem()
			if d.unmarshal(n.children[i+1], e) ***REMOVED***
				d.setMapIndex(n.children[i+1], out, k, e)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	d.mapType = mapType
	return true
***REMOVED***

func (d *decoder) setMapIndex(n *node, out, k, v reflect.Value) ***REMOVED***
	if d.strict && out.MapIndex(k) != zeroValue ***REMOVED***
		d.terrors = append(d.terrors, fmt.Sprintf("line %d: key %#v already set in map", n.line+1, k.Interface()))
		return
	***REMOVED***
	out.SetMapIndex(k, v)
***REMOVED***

func (d *decoder) mappingSlice(n *node, out reflect.Value) (good bool) ***REMOVED***
	outt := out.Type()
	if outt.Elem() != mapItemType ***REMOVED***
		d.terror(n, yaml_MAP_TAG, out)
		return false
	***REMOVED***

	mapType := d.mapType
	d.mapType = outt

	var slice []MapItem
	var l = len(n.children)
	for i := 0; i < l; i += 2 ***REMOVED***
		if isMerge(n.children[i]) ***REMOVED***
			d.merge(n.children[i+1], out)
			continue
		***REMOVED***
		item := MapItem***REMOVED******REMOVED***
		k := reflect.ValueOf(&item.Key).Elem()
		if d.unmarshal(n.children[i], k) ***REMOVED***
			v := reflect.ValueOf(&item.Value).Elem()
			if d.unmarshal(n.children[i+1], v) ***REMOVED***
				slice = append(slice, item)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	out.Set(reflect.ValueOf(slice))
	d.mapType = mapType
	return true
***REMOVED***

func (d *decoder) mappingStruct(n *node, out reflect.Value) (good bool) ***REMOVED***
	sinfo, err := getStructInfo(out.Type())
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	name := settableValueOf("")
	l := len(n.children)

	var inlineMap reflect.Value
	var elemType reflect.Type
	if sinfo.InlineMap != -1 ***REMOVED***
		inlineMap = out.Field(sinfo.InlineMap)
		inlineMap.Set(reflect.New(inlineMap.Type()).Elem())
		elemType = inlineMap.Type().Elem()
	***REMOVED***

	var doneFields []bool
	if d.strict ***REMOVED***
		doneFields = make([]bool, len(sinfo.FieldsList))
	***REMOVED***
	for i := 0; i < l; i += 2 ***REMOVED***
		ni := n.children[i]
		if isMerge(ni) ***REMOVED***
			d.merge(n.children[i+1], out)
			continue
		***REMOVED***
		if !d.unmarshal(ni, name) ***REMOVED***
			continue
		***REMOVED***
		if info, ok := sinfo.FieldsMap[name.String()]; ok ***REMOVED***
			if d.strict ***REMOVED***
				if doneFields[info.Id] ***REMOVED***
					d.terrors = append(d.terrors, fmt.Sprintf("line %d: field %s already set in type %s", ni.line+1, name.String(), out.Type()))
					continue
				***REMOVED***
				doneFields[info.Id] = true
			***REMOVED***
			var field reflect.Value
			if info.Inline == nil ***REMOVED***
				field = out.Field(info.Num)
			***REMOVED*** else ***REMOVED***
				field = out.FieldByIndex(info.Inline)
			***REMOVED***
			d.unmarshal(n.children[i+1], field)
		***REMOVED*** else if sinfo.InlineMap != -1 ***REMOVED***
			if inlineMap.IsNil() ***REMOVED***
				inlineMap.Set(reflect.MakeMap(inlineMap.Type()))
			***REMOVED***
			value := reflect.New(elemType).Elem()
			d.unmarshal(n.children[i+1], value)
			d.setMapIndex(n.children[i+1], inlineMap, name, value)
		***REMOVED*** else if d.strict ***REMOVED***
			d.terrors = append(d.terrors, fmt.Sprintf("line %d: field %s not found in type %s", ni.line+1, name.String(), out.Type()))
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func failWantMap() ***REMOVED***
	failf("map merge requires map or sequence of maps as the value")
***REMOVED***

func (d *decoder) merge(n *node, out reflect.Value) ***REMOVED***
	switch n.kind ***REMOVED***
	case mappingNode:
		d.unmarshal(n, out)
	case aliasNode:
		if n.alias != nil && n.alias.kind != mappingNode ***REMOVED***
			failWantMap()
		***REMOVED***
		d.unmarshal(n, out)
	case sequenceNode:
		// Step backwards as earlier nodes take precedence.
		for i := len(n.children) - 1; i >= 0; i-- ***REMOVED***
			ni := n.children[i]
			if ni.kind == aliasNode ***REMOVED***
				if ni.alias != nil && ni.alias.kind != mappingNode ***REMOVED***
					failWantMap()
				***REMOVED***
			***REMOVED*** else if ni.kind != mappingNode ***REMOVED***
				failWantMap()
			***REMOVED***
			d.unmarshal(ni, out)
		***REMOVED***
	default:
		failWantMap()
	***REMOVED***
***REMOVED***

func isMerge(n *node) bool ***REMOVED***
	return n.kind == scalarNode && n.value == "<<" && (n.implicit == true || n.tag == yaml_MERGE_TAG)
***REMOVED***
