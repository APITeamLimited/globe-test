//+build !go1.12

package dynamic

import (
	"github.com/jhump/protoreflect/desc"
	"reflect"
)

// Pre-Go-1.12, we must use reflect.Value.MapKeys to reflectively
// iterate a map. (We can be more efficient in Go 1.12 and up...)

func mapsEqual(a, b reflect.Value) bool ***REMOVED***
	if a.Len() != b.Len() ***REMOVED***
		return false
	***REMOVED***
	if a.Len() == 0 && b.Len() == 0 ***REMOVED***
		// Optimize the case where maps are frequently empty because MapKeys()
		// function allocates heavily.
		return true
	***REMOVED***

	for _, k := range a.MapKeys() ***REMOVED***
		av := a.MapIndex(k)
		bv := b.MapIndex(k)
		if !bv.IsValid() ***REMOVED***
			return false
		***REMOVED***
		if !fieldsEqual(av.Interface(), bv.Interface()) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func validFieldValueForMapField(fd *desc.FieldDescriptor, val reflect.Value) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	// make a defensive copy while we check the contents
	// (also converts to map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED*** if it's some other type)
	keyField := fd.GetMessageType().GetFields()[0]
	valField := fd.GetMessageType().GetFields()[1]
	m := map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	for _, k := range val.MapKeys() ***REMOVED***
		if k.Kind() == reflect.Interface ***REMOVED***
			// unwrap it
			k = reflect.ValueOf(k.Interface())
		***REMOVED***
		kk, err := validElementFieldValueForRv(keyField, k, false)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		v := val.MapIndex(k)
		if v.Kind() == reflect.Interface ***REMOVED***
			// unwrap it
			v = reflect.ValueOf(v.Interface())
		***REMOVED***
		vv, err := validElementFieldValueForRv(valField, v, true)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		m[kk] = vv
	***REMOVED***
	return m, nil
***REMOVED***

func canConvertMap(src reflect.Value, target reflect.Type) bool ***REMOVED***
	kt := target.Key()
	vt := target.Elem()
	for _, k := range src.MapKeys() ***REMOVED***
		if !canConvert(k, kt) ***REMOVED***
			return false
		***REMOVED***
		if !canConvert(src.MapIndex(k), vt) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func mergeMapVal(src, target reflect.Value, targetType reflect.Type, deterministic bool) error ***REMOVED***
	tkt := targetType.Key()
	tvt := targetType.Elem()
	for _, k := range src.MapKeys() ***REMOVED***
		v := src.MapIndex(k)
		skt := k.Type()
		svt := v.Type()
		var nk, nv reflect.Value
		if tkt == skt ***REMOVED***
			nk = k
		***REMOVED*** else if tkt.Kind() == reflect.Ptr && tkt.Elem() == skt ***REMOVED***
			nk = k.Addr()
		***REMOVED*** else ***REMOVED***
			nk = reflect.New(tkt).Elem()
			if err := mergeVal(k, nk, deterministic); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if tvt == svt ***REMOVED***
			nv = v
		***REMOVED*** else if tvt.Kind() == reflect.Ptr && tvt.Elem() == svt ***REMOVED***
			nv = v.Addr()
		***REMOVED*** else ***REMOVED***
			nv = reflect.New(tvt).Elem()
			if err := mergeVal(v, nv, deterministic); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if target.IsNil() ***REMOVED***
			target.Set(reflect.MakeMap(targetType))
		***REMOVED***
		target.SetMapIndex(nk, nv)
	***REMOVED***
	return nil
***REMOVED***

func mergeMapField(m *Message, fd *desc.FieldDescriptor, rv reflect.Value) error ***REMOVED***
	for _, k := range rv.MapKeys() ***REMOVED***
		if k.Kind() == reflect.Interface && !k.IsNil() ***REMOVED***
			k = k.Elem()
		***REMOVED***
		v := rv.MapIndex(k)
		if v.Kind() == reflect.Interface && !v.IsNil() ***REMOVED***
			v = v.Elem()
		***REMOVED***
		if err := m.putMapField(fd, k.Interface(), v.Interface()); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
