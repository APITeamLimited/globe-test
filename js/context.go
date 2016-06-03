package js

import (
	"fmt"
	// "github.com/loadimpact/speedboat"
	"gopkg.in/olebedev/go-duktape.v2"
)

const (
	scriptProp = "__script__"
)

type JSError struct ***REMOVED***
	Message  string
	Filename string
	Line     int
***REMOVED***

func (e JSError) Error() string ***REMOVED***
	return fmt.Sprintf("%s:%d: %s", e.Filename, e.Line, e.Message)
***REMOVED***

func getJSError(js *duktape.Context) JSError ***REMOVED***
	js.GetPropString(-1, "fileName")
	filename := js.SafeToString(-1)
	js.Pop()

	js.GetPropString(-1, "lineNumber")
	line := js.ToInt(-1)
	js.Pop()

	msg := js.SafeToString(-1)
	return JSError***REMOVED***Message: msg, Filename: filename, Line: line***REMOVED***
***REMOVED***

func setupGlobalObject(js *duktape.Context, id int) ***REMOVED***
	js.PushGlobalObject()
	defer js.Pop()

	js.PushObject()
	js.PutPropString(-2, "__modules__")

	js.PushObject()
	***REMOVED***
		js.PushInt(id)
		js.PutPropString(-2, "id")
	***REMOVED***
	js.PutPropString(-2, "__data__")
***REMOVED***

func putScript(js *duktape.Context, filename, src string) error ***REMOVED***
	js.PushGlobalObject()
	defer js.Pop()

	js.PushString(filename)
	if err := js.PcompileStringFilename(0, src); err != nil ***REMOVED***
		return err
	***REMOVED***
	js.PutPropString(-2, scriptProp)

	return nil
***REMOVED***

func loadScript(js *duktape.Context, filename, src string) error ***REMOVED***
	js.PushString(filename)
	if err := js.PcompileStringFilename(0, src); err != nil ***REMOVED***
		return err
	***REMOVED***

	if js.Pcall(0) != duktape.ErrNone ***REMOVED***
		err := getJSError(js)
		js.Pop()
		return err
	***REMOVED***
	js.Pop()
	return nil
***REMOVED***
