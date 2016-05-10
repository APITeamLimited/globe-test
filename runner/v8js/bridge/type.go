package bridge

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"reflect"
)

type Type struct ***REMOVED***
	Type    reflect.Type
	Spec    map[string]Type
	JSONKey string
***REMOVED***

func (t *Type) Cast(v *interface***REMOVED******REMOVED***) error ***REMOVED***
	rV := reflect.ValueOf(*v)
	vT := rV.Type()
	if vT == t.Type ***REMOVED***
		return nil
	***REMOVED***

	switch t.Type.Kind() ***REMOVED***
	case reflect.Struct:
		if vT.Kind() != reflect.Map ***REMOVED***
			return errors.New("Invalid argument")
		***REMOVED***
	default:
		if !vT.ConvertibleTo(t.Type) ***REMOVED***
			log.WithFields(log.Fields***REMOVED***
				"expected": t.Type,
				"actual":   vT,
			***REMOVED***).Debug("Invalid argument")
			return errors.New("Invalid argument")
		***REMOVED***
		rV = rV.Convert(t.Type)
		*v = rV.Interface()
	***REMOVED***

	return nil
***REMOVED***

// Creates a bridged type.
// Panics if raw is a function.
func BridgeType(raw reflect.Type) Type ***REMOVED***
	tp := Type***REMOVED***Type: raw***REMOVED***
	kind := tp.Type.Kind()

	if kind == reflect.Func ***REMOVED***
		panic(errors.New("That's a function, bridge it as such"))
	***REMOVED***

	if kind == reflect.Struct ***REMOVED***
		tp.Spec = make(map[string]Type)
		for i := 0; i < raw.NumField(); i++ ***REMOVED***
			f := raw.Field(i)
			tag := f.Tag.Get("json")
			if tag == "" || tag == "-" ***REMOVED***
				continue
			***REMOVED***

			ftp := BridgeType(f.Type)
			ftp.JSONKey = tag
			tp.Spec[f.Name] = ftp
		***REMOVED***
	***REMOVED***

	return tp
***REMOVED***
