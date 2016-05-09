package bridge

import (
	"errors"
	"reflect"
)

type Type struct ***REMOVED***
	Kind    reflect.Kind
	Spec    map[string]Type
	JSONKey string
***REMOVED***

// Creates a bridged type.
// Panics if raw is a function.
func BridgeType(raw reflect.Type) Type ***REMOVED***
	tp := Type***REMOVED***Kind: raw.Kind()***REMOVED***

	if tp.Kind == reflect.Func ***REMOVED***
		panic(errors.New("That's a function, bridge it as such"))
	***REMOVED***

	if tp.Kind == reflect.Struct ***REMOVED***
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
