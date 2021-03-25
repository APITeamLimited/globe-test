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

// jsonNumber is the interface of the encoding/json.Number datatype.
// Repeating the interface here avoids a dependency on encoding/json, and also
// supports other libraries like jsoniter, which use a similar datatype with
// the same interface. Detecting this interface is useful when dealing with
// structures containing json.Number, which is a string under the hood. The
// encoder should prefer the use of Int64(), Float64() and string(), in that
// order, when encoding this type.
type jsonNumber interface ***REMOVED***
	Float64() (float64, error)
	Int64() (int64, error)
	String() string
***REMOVED***

type encoder struct ***REMOVED***
	emitter yaml_emitter_t
	event   yaml_event_t
	out     []byte
	flow    bool
	// doneInit holds whether the initial stream_start_event has been
	// emitted.
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
	yaml_document_start_event_initialize(&e.event, nil, nil, true)
	e.emit()
	e.marshal(tag, in)
	yaml_document_end_event_initialize(&e.event, true)
	e.emit()
***REMOVED***

func (e *encoder) marshal(tag string, in reflect.Value) ***REMOVED***
	if !in.IsValid() || in.Kind() == reflect.Ptr && in.IsNil() ***REMOVED***
		e.nilv()
		return
	***REMOVED***
	iface := in.Interface()
	switch m := iface.(type) ***REMOVED***
	case jsonNumber:
		integer, err := m.Int64()
		if err == nil ***REMOVED***
			// In this case the json.Number is a valid int64
			in = reflect.ValueOf(integer)
			break
		***REMOVED***
		float, err := m.Float64()
		if err == nil ***REMOVED***
			// In this case the json.Number is a valid float64
			in = reflect.ValueOf(float)
			break
		***REMOVED***
		// fallback case - no number could be obtained
		in = reflect.ValueOf(m.String())
	case time.Time, *time.Time:
		// Although time.Time implements TextMarshaler,
		// we don't want to treat it as a string for YAML
		// purposes because YAML has special support for
		// timestamps.
	case Marshaler:
		v, err := m.MarshalYAML()
		if err != nil ***REMOVED***
			fail(err)
		***REMOVED***
		if v == nil ***REMOVED***
			e.nilv()
			return
		***REMOVED***
		in = reflect.ValueOf(v)
	case encoding.TextMarshaler:
		text, err := m.MarshalText()
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
		if in.Type() == ptrTimeType ***REMOVED***
			e.timev(tag, in.Elem())
		***REMOVED*** else ***REMOVED***
			e.marshal(tag, in.Elem())
		***REMOVED***
	case reflect.Struct:
		if in.Type() == timeType ***REMOVED***
			e.timev(tag, in)
		***REMOVED*** else ***REMOVED***
			e.structv(tag, in)
		***REMOVED***
	case reflect.Slice, reflect.Array:
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

func (e *encoder) stringv(tag string, in reflect.Value) ***REMOVED***
	var style yaml_scalar_style_t
	s := in.String()
	canUsePlain := true
	switch ***REMOVED***
	case !utf8.ValidString(s):
		if tag == yaml_BINARY_TAG ***REMOVED***
			failf("explicitly tagged !!binary data must be base64-encoded")
		***REMOVED***
		if tag != "" ***REMOVED***
			failf("cannot marshal invalid UTF-8 data as %s", shortTag(tag))
		***REMOVED***
		// It can't be encoded directly as YAML so use a binary tag
		// and encode it as base64.
		tag = yaml_BINARY_TAG
		s = encodeBase64(s)
	case tag == "":
		// Check to see if it would resolve to a specific
		// tag when encoded unquoted. If it doesn't,
		// there's no need to quote it.
		rtag, _ := resolve("", s)
		canUsePlain = rtag == yaml_STR_TAG && !isBase60Float(s)
	***REMOVED***
	// Note: it's possible for user code to emit invalid YAML
	// if they explicitly specify a tag and a string containing
	// text that's incompatible with that tag.
	switch ***REMOVED***
	case strings.Contains(s, "\n"):
		style = yaml_LITERAL_SCALAR_STYLE
	case canUsePlain:
		style = yaml_PLAIN_SCALAR_STYLE
	default:
		style = yaml_DOUBLE_QUOTED_SCALAR_STYLE
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

func (e *encoder) timev(tag string, in reflect.Value) ***REMOVED***
	t := in.Interface().(time.Time)
	s := t.Format(time.RFC3339Nano)
	e.emitScalar(s, "", tag, yaml_PLAIN_SCALAR_STYLE)
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
