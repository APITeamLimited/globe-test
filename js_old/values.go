package js

import (
	"encoding/json"
	"gopkg.in/olebedev/go-duktape.v2"
)

func argNumber(js *duktape.Context, index int) float64 ***REMOVED***
	if js.GetTopIndex() < index ***REMOVED***
		return 0
	***REMOVED***

	return js.ToNumber(index)
***REMOVED***

func argString(js *duktape.Context, index int) string ***REMOVED***
	if js.GetTopIndex() < index ***REMOVED***
		return ""
	***REMOVED***

	return js.ToString(index)
***REMOVED***

func argJSON(js *duktape.Context, index int, out interface***REMOVED******REMOVED***) error ***REMOVED***
	if js.GetTopIndex() < index ***REMOVED***
		return nil
	***REMOVED***

	js.JsonEncode(index)
	str := js.GetString(index)
	return json.Unmarshal([]byte(str), out)
***REMOVED***

func pushObject(js *duktape.Context, obj interface***REMOVED******REMOVED***, t string) error ***REMOVED***
	s, err := json.Marshal(obj)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	js.PushString(string(s))
	js.JsonDecode(-1)

	if t != "" ***REMOVED***
		js.PushGlobalObject()
		***REMOVED***
			js.GetPropString(-1, t)
			js.SetPrototype(-3)
		***REMOVED***
		js.Pop()
	***REMOVED***

	return nil
***REMOVED***
