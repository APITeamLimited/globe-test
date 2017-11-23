package yaml

import (
	"encoding"
	"encoding/base64"
	"fmt"
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
	value        string
	implicit     bool
	children     []*node
	anchors      map[string]*node
***REMOVED***

// ----------------------------------------------------------------------------
// Parser, produces a node tree out of a libyaml event stream.

type parser struct ***REMOVED***
	parser yaml_parser_t
	event  yaml_event_t
	doc    *node
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

	p.skip()
	if p.event.typ != yaml_STREAM_START_EVENT ***REMOVED***
		panic("expected stream start event, got " + strconv.Itoa(int(p.event.typ)))
	***REMOVED***
	p.skip()
	return &p
***REMOVED***

func (p *parser) destroy() ***REMOVED***
	if p.event.typ != yaml_NO_EVENT ***REMOVED***
		yaml_event_delete(&p.event)
	***REMOVED***
	yaml_parser_delete(&p.parser)
***REMOVED***

func (p *parser) skip() ***REMOVED***
	if p.event.typ != yaml_NO_EVENT ***REMOVED***
		if p.event.typ == yaml_STREAM_END_EVENT ***REMOVED***
			failf("attempted to go past the end of stream; corrupted value?")
		***REMOVED***
		yaml_event_delete(&p.event)
	***REMOVED***
	if !yaml_parser_parse(&p.parser, &p.event) ***REMOVED***
		p.fail()
	***REMOVED***
***REMOVED***

func (p *parser) fail() ***REMOVED***
	var where string
	var line int
	if p.parser.problem_mark.line != 0 ***REMOVED***
		line = p.parser.problem_mark.line
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
	switch p.event.typ ***REMOVED***
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
		panic("attempted to parse unknown event: " + strconv.Itoa(int(p.event.typ)))
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
	p.skip()
	n.children = append(n.children, p.parse())
	if p.event.typ != yaml_DOCUMENT_END_EVENT ***REMOVED***
		panic("expected end of document event but got " + strconv.Itoa(int(p.event.typ)))
	***REMOVED***
	p.skip()
	return n
***REMOVED***

func (p *parser) alias() *node ***REMOVED***
	n := p.node(aliasNode)
	n.value = string(p.event.anchor)
	p.skip()
	return n
***REMOVED***

func (p *parser) scalar() *node ***REMOVED***
	n := p.node(scalarNode)
	n.value = string(p.event.value)
	n.tag = string(p.event.tag)
	n.implicit = p.event.implicit
	p.anchor(n, p.event.anchor)
	p.skip()
	return n
***REMOVED***

func (p *parser) sequence() *node ***REMOVED***
	n := p.node(sequenceNode)
	p.anchor(n, p.event.anchor)
	p.skip()
	for p.event.typ != yaml_SEQUENCE_END_EVENT ***REMOVED***
		n.children = append(n.children, p.parse())
	***REMOVED***
	p.skip()
	return n
***REMOVED***

func (p *parser) mapping() *node ***REMOVED***
	n := p.node(mappingNode)
	p.anchor(n, p.event.anchor)
	p.skip()
	for p.event.typ != yaml_MAPPING_END_EVENT ***REMOVED***
		n.children = append(n.children, p.parse(), p.parse())
	***REMOVED***
	p.skip()
	return n
***REMOVED***

// ----------------------------------------------------------------------------
// Decoder, unmarshals a node into a provided value.

type decoder struct ***REMOVED***
	doc     *node
	aliases map[string]bool
	mapType reflect.Type
	terrors []string
	strict  bool
***REMOVED***

var (
	mapItemType    = reflect.TypeOf(MapItem***REMOVED******REMOVED***)
	durationType   = reflect.TypeOf(time.Duration(0))
	defaultMapType = reflect.TypeOf(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED******REMOVED***)
	ifaceType      = defaultMapType.Elem()
)

func newDecoder(strict bool) *decoder ***REMOVED***
	d := &decoder***REMOVED***mapType: defaultMapType, strict: strict***REMOVED***
	d.aliases = make(map[string]bool)
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

func (d *decoder) unmarshal(n *node, out reflect.Value) (good bool) ***REMOVED***
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
	an, ok := d.doc.anchors[n.value]
	if !ok ***REMOVED***
		failf("unknown anchor '%s' referenced", n.value)
	***REMOVED***
	if d.aliases[n.value] ***REMOVED***
		failf("anchor '%s' value contains itself", n.value)
	***REMOVED***
	d.aliases[n.value] = true
	good = d.unmarshal(an, out)
	delete(d.aliases, n.value)
	return good
***REMOVED***

var zeroValue reflect.Value

func resetMap(out reflect.Value) ***REMOVED***
	for _, k := range out.MapKeys() ***REMOVED***
		out.SetMapIndex(k, zeroValue)
	***REMOVED***
***REMOVED***

func (d *decoder) scalar(n *node, out reflect.Value) (good bool) ***REMOVED***
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
	if s, ok := resolved.(string); ok && out.CanAddr() ***REMOVED***
		if u, ok := out.Addr().Interface().(encoding.TextUnmarshaler); ok ***REMOVED***
			err := u.UnmarshalText([]byte(s))
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
			good = true
		***REMOVED*** else if resolved != nil ***REMOVED***
			out.SetString(n.value)
			good = true
		***REMOVED***
	case reflect.Interface:
		if resolved == nil ***REMOVED***
			out.Set(reflect.Zero(out.Type()))
		***REMOVED*** else ***REMOVED***
			out.Set(reflect.ValueOf(resolved))
		***REMOVED***
		good = true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch resolved := resolved.(type) ***REMOVED***
		case int:
			if !out.OverflowInt(int64(resolved)) ***REMOVED***
				out.SetInt(int64(resolved))
				good = true
			***REMOVED***
		case int64:
			if !out.OverflowInt(resolved) ***REMOVED***
				out.SetInt(resolved)
				good = true
			***REMOVED***
		case uint64:
			if resolved <= math.MaxInt64 && !out.OverflowInt(int64(resolved)) ***REMOVED***
				out.SetInt(int64(resolved))
				good = true
			***REMOVED***
		case float64:
			if resolved <= math.MaxInt64 && !out.OverflowInt(int64(resolved)) ***REMOVED***
				out.SetInt(int64(resolved))
				good = true
			***REMOVED***
		case string:
			if out.Type() == durationType ***REMOVED***
				d, err := time.ParseDuration(resolved)
				if err == nil ***REMOVED***
					out.SetInt(int64(d))
					good = true
				***REMOVED***
			***REMOVED***
		***REMOVED***
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch resolved := resolved.(type) ***REMOVED***
		case int:
			if resolved >= 0 && !out.OverflowUint(uint64(resolved)) ***REMOVED***
				out.SetUint(uint64(resolved))
				good = true
			***REMOVED***
		case int64:
			if resolved >= 0 && !out.OverflowUint(uint64(resolved)) ***REMOVED***
				out.SetUint(uint64(resolved))
				good = true
			***REMOVED***
		case uint64:
			if !out.OverflowUint(uint64(resolved)) ***REMOVED***
				out.SetUint(uint64(resolved))
				good = true
			***REMOVED***
		case float64:
			if resolved <= math.MaxUint64 && !out.OverflowUint(uint64(resolved)) ***REMOVED***
				out.SetUint(uint64(resolved))
				good = true
			***REMOVED***
		***REMOVED***
	case reflect.Bool:
		switch resolved := resolved.(type) ***REMOVED***
		case bool:
			out.SetBool(resolved)
			good = true
		***REMOVED***
	case reflect.Float32, reflect.Float64:
		switch resolved := resolved.(type) ***REMOVED***
		case int:
			out.SetFloat(float64(resolved))
			good = true
		case int64:
			out.SetFloat(float64(resolved))
			good = true
		case uint64:
			out.SetFloat(float64(resolved))
			good = true
		case float64:
			out.SetFloat(resolved)
			good = true
		***REMOVED***
	case reflect.Ptr:
		if out.Type().Elem() == reflect.TypeOf(resolved) ***REMOVED***
			// TODO DOes this make sense? When is out a Ptr except when decoding a nil value?
			elem := reflect.New(out.Type().Elem())
			elem.Elem().Set(reflect.ValueOf(resolved))
			out.Set(elem)
			good = true
		***REMOVED***
	***REMOVED***
	if !good ***REMOVED***
		d.terror(n, tag, out)
	***REMOVED***
	return good
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
	out.Set(out.Slice(0, j))
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
				out.SetMapIndex(k, e)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	d.mapType = mapType
	return true
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
			inlineMap.SetMapIndex(name, value)
		***REMOVED*** else if d.strict ***REMOVED***
			d.terrors = append(d.terrors, fmt.Sprintf("line %d: field %s not found in struct %s", n.line+1, name.String(), out.Type()))
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
		an, ok := d.doc.anchors[n.value]
		if ok && an.kind != mappingNode ***REMOVED***
			failWantMap()
		***REMOVED***
		d.unmarshal(n, out)
	case sequenceNode:
		// Step backwards as earlier nodes take precedence.
		for i := len(n.children) - 1; i >= 0; i-- ***REMOVED***
			ni := n.children[i]
			if ni.kind == aliasNode ***REMOVED***
				an, ok := d.doc.anchors[ni.value]
				if ok && an.kind != mappingNode ***REMOVED***
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
