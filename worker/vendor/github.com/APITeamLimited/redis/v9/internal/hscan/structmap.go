package hscan

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// structMap contains the map of struct fields for target structs
// indexed by the struct type.
type structMap struct ***REMOVED***
	m sync.Map
***REMOVED***

func newStructMap() *structMap ***REMOVED***
	return new(structMap)
***REMOVED***

func (s *structMap) get(t reflect.Type) *structSpec ***REMOVED***
	if v, ok := s.m.Load(t); ok ***REMOVED***
		return v.(*structSpec)
	***REMOVED***

	spec := newStructSpec(t, "redis")
	s.m.Store(t, spec)
	return spec
***REMOVED***

//------------------------------------------------------------------------------

// structSpec contains the list of all fields in a target struct.
type structSpec struct ***REMOVED***
	m map[string]*structField
***REMOVED***

func (s *structSpec) set(tag string, sf *structField) ***REMOVED***
	s.m[tag] = sf
***REMOVED***

func newStructSpec(t reflect.Type, fieldTag string) *structSpec ***REMOVED***
	numField := t.NumField()
	out := &structSpec***REMOVED***
		m: make(map[string]*structField, numField),
	***REMOVED***

	for i := 0; i < numField; i++ ***REMOVED***
		f := t.Field(i)

		tag := f.Tag.Get(fieldTag)
		if tag == "" || tag == "-" ***REMOVED***
			continue
		***REMOVED***

		tag = strings.Split(tag, ",")[0]
		if tag == "" ***REMOVED***
			continue
		***REMOVED***

		// Use the built-in decoder.
		out.set(tag, &structField***REMOVED***index: i, fn: decoders[f.Type.Kind()]***REMOVED***)
	***REMOVED***

	return out
***REMOVED***

//------------------------------------------------------------------------------

// structField represents a single field in a target struct.
type structField struct ***REMOVED***
	index int
	fn    decoderFunc
***REMOVED***

//------------------------------------------------------------------------------

type StructValue struct ***REMOVED***
	spec  *structSpec
	value reflect.Value
***REMOVED***

func (s StructValue) Scan(key string, value string) error ***REMOVED***
	field, ok := s.spec.m[key]
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	if err := field.fn(s.value.Field(field.index), value); err != nil ***REMOVED***
		t := s.value.Type()
		return fmt.Errorf("cannot scan redis.result %s into struct field %s.%s of type %s, error-%s",
			value, t.Name(), t.Field(field.index).Name, t.Field(field.index).Type, err.Error())
	***REMOVED***
	return nil
***REMOVED***
