package internal

import (
	"reflect"

	"github.com/golang/protobuf/proto"
)

var typeOfBytes = reflect.TypeOf([]byte(nil))

// GetUnrecognized fetches the bytes of unrecognized fields for the given message.
func GetUnrecognized(msg proto.Message) []byte ***REMOVED***
	val := reflect.Indirect(reflect.ValueOf(msg))
	u := val.FieldByName("XXX_unrecognized")
	if u.IsValid() && u.Type() == typeOfBytes ***REMOVED***
		return u.Interface().([]byte)
	***REMOVED***

	// Fallback to reflection for API v2 messages
	get, _, _, ok := unrecognizedGetSetMethods(val)
	if !ok ***REMOVED***
		return nil
	***REMOVED***

	return get.Call([]reflect.Value(nil))[0].Convert(typeOfBytes).Interface().([]byte)
***REMOVED***

// SetUnrecognized adds the given bytes to the unrecognized fields for the given message.
func SetUnrecognized(msg proto.Message, data []byte) ***REMOVED***
	val := reflect.Indirect(reflect.ValueOf(msg))
	u := val.FieldByName("XXX_unrecognized")
	if u.IsValid() && u.Type() == typeOfBytes ***REMOVED***
		// Just store the bytes in the unrecognized field
		ub := u.Interface().([]byte)
		ub = append(ub, data...)
		u.Set(reflect.ValueOf(ub))
		return
	***REMOVED***

	// Fallback to reflection for API v2 messages
	get, set, argType, ok := unrecognizedGetSetMethods(val)
	if !ok ***REMOVED***
		return
	***REMOVED***

	existing := get.Call([]reflect.Value(nil))[0].Convert(typeOfBytes).Interface().([]byte)
	if len(existing) > 0 ***REMOVED***
		data = append(existing, data...)
	***REMOVED***
	set.Call([]reflect.Value***REMOVED***reflect.ValueOf(data).Convert(argType)***REMOVED***)
***REMOVED***

func unrecognizedGetSetMethods(val reflect.Value) (get reflect.Value, set reflect.Value, argType reflect.Type, ok bool) ***REMOVED***
	// val could be an APIv2 message. We use reflection to interact with
	// this message so that we don't have a hard dependency on the new
	// version of the protobuf package.
	refMethod := val.MethodByName("ProtoReflect")
	if !refMethod.IsValid() ***REMOVED***
		if val.CanAddr() ***REMOVED***
			refMethod = val.Addr().MethodByName("ProtoReflect")
		***REMOVED***
		if !refMethod.IsValid() ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	refType := refMethod.Type()
	if refType.NumIn() != 0 || refType.NumOut() != 1 ***REMOVED***
		return
	***REMOVED***
	ref := refMethod.Call([]reflect.Value(nil))
	getMethod, setMethod := ref[0].MethodByName("GetUnknown"), ref[0].MethodByName("SetUnknown")
	if !getMethod.IsValid() || !setMethod.IsValid() ***REMOVED***
		return
	***REMOVED***
	getType := getMethod.Type()
	setType := setMethod.Type()
	if getType.NumIn() != 0 || getType.NumOut() != 1 || setType.NumIn() != 1 || setType.NumOut() != 0 ***REMOVED***
		return
	***REMOVED***
	arg := setType.In(0)
	if !arg.ConvertibleTo(typeOfBytes) || getType.Out(0) != arg ***REMOVED***
		return
	***REMOVED***

	return getMethod, setMethod, arg, true
***REMOVED***
