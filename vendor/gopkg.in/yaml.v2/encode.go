package yaml

import (
	"encoding"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type encoder struct ***REMOVED***
	emitter yaml_emitter_t
	event   yaml_event_t
	out     []byte
	flow    bool
***REMOVED***

func newEncoder() (e *encoder) ***REMOVED***
	e = &encoder***REMOVED******REMOVED***
	e.must(yaml_emitter_initialize(&e.emitter))
	yaml_emitter_set_output_string(&e.emitter, &e.out)
	yaml_emitter_set_unicode(&e.emitter, true)
	e.must(yaml_stream_start_event_initialize(&e.event, yaml_UTF8_ENCODING))
	e.emit()
	e.must(yaml_document_start_event_initialize(&e.event, nil, nil, true))
	e.emit()
	return e
***REMOVED***

func (e *encoder) finish() ***REMOVED***
	e.must(yaml_document_end_event_initialize(&e.event, true))
	e.emit()
	e.emitter.open_ended = false
	e.must(yaml_stream_end_event_initialize(&e.event))
	e.emit()
***REMOVED***

func (e *encoder) destroy() ***REMOVED***
	yaml_emitter_delete(&e.emitter)
***REMOVED***

func (e *encoder) emit() ***REMOVED***
	// This will internally delete the e.event value.
	if !yaml_emitter_emit(&e.emitter, &e.event) && e.event.typ != yaml_DOCUMENT_END_EVENT && e.event.typ != yaml_STREAM_END_EVENT ***REMOVED***
		e.must(false)
	***REMOVED***
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

func (e *encoder) marshal(tag string, in reflect.Value) ***REMOVED***
	if !in.IsValid() ***REMOVED***
		e.nilv()
		return
	***REMOVED***
	iface := in.Interface()
	if m, ok := iface.(Marshaler); ok ***REMOVED***
		v, err := m.MarshalYAML()
		if err != nil ***REMOVED***
			fail(err)
		***REMOVED***
		if v == nil ***REMOVED***
			e.nilv()
			return
		***REMOVED***
		in = reflect.ValueOf(v)
	***REMOVED*** else if m, ok := iface.(encoding.TextMarshaler); ok ***REMOVED***
		text, err := m.MarshalText()
		if err != nil ***REMOVED***
			fail(err)
		***REMOVED***
		in = reflect.ValueOf(string(text))
	***REMOVED***
	switch in.Kind() ***REMOVED***
	case reflect.Interface:
		if in.IsNil() ***REMOVED***
			e.nilv()
		***REMOVED*** else ***REMOVED***
			e.marshal(tag, in.Elem())
		***REMOVED***
	case reflect.Map:
		e.mapv(tag, in)
	case reflect.Ptr:
		if in.IsNil() ***REMOVED***
			e.nilv()
		***REMOVED*** else ***REMOVED***
			e.marshal(tag, in.Elem())
		***REMOVED***
	case reflect.Struct:
		e.structv(tag, in)
	case reflect.Slice:
		if in.Type().Elem() == mapItemType ***REMOVED***
			e.itemsv(tag, in)
		***REMOVED*** else ***REMOVED***
			e.slicev(tag, in)
		***REMOVED***
	case reflect.String:
		e.stringv(tag, in)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if in.Type() == durationType ***REMOVED***
			e.stringv(tag, reflect.ValueOf(iface.(time.Duration).String()))
		***REMOVED*** else ***REMOVED***
			e.intv(tag, in)
		***REMOVED***
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

func (e *encoder) itemsv(tag string, in reflect.Value) ***REMOVED***
	e.mappingv(tag, func() ***REMOVED***
		slice := in.Convert(reflect.TypeOf([]MapItem***REMOVED******REMOVED***)).Interface().([]MapItem)
		for _, item := range slice ***REMOVED***
			e.marshal("", reflect.ValueOf(item.Key))
			e.marshal("", reflect.ValueOf(item.Value))
		***REMOVED***
	***REMOVED***)
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
				value = in.FieldByIndex(info.Inline)
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
						panic(fmt.Sprintf("Can't have key %q in inlined map; conflicts with struct field", k.String()))
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
	e.must(yaml_mapping_start_event_initialize(&e.event, nil, []byte(tag), implicit, style))
	e.emit()
	f()
	e.must(yaml_mapping_end_event_initialize(&e.event))
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

func (e *encoder) stringv(tag string, in reflect.Value) ***REMOVED***
	var style yaml_scalar_style_t
	s := in.String()
	rtag, rs := resolve("", s)
	if rtag == yaml_BINARY_TAG ***REMOVED***
		if tag == "" || tag == yaml_STR_TAG ***REMOVED***
			tag = rtag
			s = rs.(string)
		***REMOVED*** else if tag == yaml_BINARY_TAG ***REMOVED***
			failf("explicitly tagged !!binary data must be base64-encoded")
		***REMOVED*** else ***REMOVED***
			failf("cannot marshal invalid UTF-8 data as %s", shortTag(tag))
		***REMOVED***
	***REMOVED***
	if tag == "" && (rtag != yaml_STR_TAG || isBase60Float(s)) ***REMOVED***
		style = yaml_DOUBLE_QUOTED_SCALAR_STYLE
	***REMOVED*** else if strings.Contains(s, "\n") ***REMOVED***
		style = yaml_LITERAL_SCALAR_STYLE
	***REMOVED*** else ***REMOVED***
		style = yaml_PLAIN_SCALAR_STYLE
	***REMOVED***
	e.emitScalar(s, "", tag, style)
***REMOVED***

func (e *encoder) boolv(tag string, in reflect.Value) ***REMOVED***
	var s string
	if in.Bool() ***REMOVED***
		s = "true"
	***REMOVED*** else ***REMOVED***
		s = "false"
	***REMOVED***
	e.emitScalar(s, "", tag, yaml_PLAIN_SCALAR_STYLE)
***REMOVED***

func (e *encoder) intv(tag string, in reflect.Value) ***REMOVED***
	s := strconv.FormatInt(in.Int(), 10)
	e.emitScalar(s, "", tag, yaml_PLAIN_SCALAR_STYLE)
***REMOVED***

func (e *encoder) uintv(tag string, in reflect.Value) ***REMOVED***
	s := strconv.FormatUint(in.Uint(), 10)
	e.emitScalar(s, "", tag, yaml_PLAIN_SCALAR_STYLE)
***REMOVED***

func (e *encoder) floatv(tag string, in reflect.Value) ***REMOVED***
	// FIXME: Handle 64 bits here.
	s := strconv.FormatFloat(float64(in.Float()), 'g', -1, 32)
	switch s ***REMOVED***
	case "+Inf":
		s = ".inf"
	case "-Inf":
		s = "-.inf"
	case "NaN":
		s = ".nan"
	***REMOVED***
	e.emitScalar(s, "", tag, yaml_PLAIN_SCALAR_STYLE)
***REMOVED***

func (e *encoder) nilv() ***REMOVED***
	e.emitScalar("null", "", "", yaml_PLAIN_SCALAR_STYLE)
***REMOVED***

func (e *encoder) emitScalar(value, anchor, tag string, style yaml_scalar_style_t) ***REMOVED***
	implicit := tag == ""
	e.must(yaml_scalar_event_initialize(&e.event, []byte(anchor), []byte(tag), []byte(value), implicit, implicit, style))
	e.emit()
***REMOVED***
