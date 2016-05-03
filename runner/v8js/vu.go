package v8js

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/runner"
	"github.com/ry/v8worker"
	"reflect"
	"strings"
)

type jsCallEnvelope struct ***REMOVED***
	Mod  string        `json:"m"`
	Fn   string        `json:"f"`
	Args []interface***REMOVED******REMOVED*** `json:"a"`
***REMOVED***

// Aaaaaa, this is awful, it needs restructuring BADLY x_x
func (vu *VUContext) BridgeAPI(w *v8worker.Worker) error ***REMOVED***
	for modname, mod := range vu.api ***REMOVED***
		jsMod := fmt.Sprintf(`
		speedboat._modules["%s"] = ***REMOVED******REMOVED***;
		`, modname)
		for name, mem := range mod ***REMOVED***
			t := reflect.TypeOf(mem)

			if t.Kind() != reflect.Func ***REMOVED***
				return errors.New("Not a function: " + modname + "." + name)
			***REMOVED***

			jsFn := fmt.Sprintf(`speedboat._modules["%s"]["%s"] = function() ***REMOVED***
				var args = [];
			`, modname, name)

			numArgs := t.NumIn()
			if !t.IsVariadic() ***REMOVED***
				jsFn += fmt.Sprintf(`
					if (arguments.length != %d) ***REMOVED***
						throw new Error("wrong number of arguments");
					***REMOVED***
				`, t.NumIn())
			***REMOVED*** else ***REMOVED***
				numArgs--
			***REMOVED***

			for i := 0; i < numArgs; i++ ***REMOVED***
				aT := t.In(i)
				switch aT.Kind() ***REMOVED***
				case reflect.Struct:
					types := make([]string, 0, aT.NumField())
					for i := 0; i < aT.NumField(); i++ ***REMOVED***
						field := aT.Field(i)
						if field.Anonymous ***REMOVED***
							continue
						***REMOVED***
						key := field.Tag.Get("json") // Does not handle comma params yet!
						if key == "" ***REMOVED***
							key = field.Name
						***REMOVED***
						val := aT.Kind().String()
						types = append(types, fmt.Sprintf(`"%s":"%s"`, key, val))
					***REMOVED***
					jsFn += fmt.Sprintf(`args.push(speedboat._require.struct(***REMOVED***%s***REMOVED***, arguments[%d]));`, strings.Join(types, ","), i)
				default:
					jsFn += fmt.Sprintf("args.push(speedboat._require.%s(arguments[%d]));", aT.Kind().String(), i)
				***REMOVED***
			***REMOVED***
			if t.IsVariadic() ***REMOVED***
				varArg := t.In(numArgs)
				eT := varArg.Elem()
				jsFn += fmt.Sprintf(`
					for (var i = %d; i < arguments.length; i++) ***REMOVED***
						args.push(speedboat._require.%s(arguments[i]));
					***REMOVED***
				`, numArgs, eT.Kind().String())
			***REMOVED***

			jsFn += fmt.Sprintf(`
				return speedboat._invoke('%s', '%s', args, %v);
			***REMOVED***`, modname, name, false)
			jsMod += "\n\n" + jsFn
		***REMOVED***

		if err := w.Load("module:"+modname, jsMod); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Make functions in the "global" module global, preimport console
	makeGlobals := `
	for (key in speedboat._modules['global']) ***REMOVED***
		eval(key + " = speedboat._modules['global']['" + key + "'];");
	***REMOVED***
	var console = speedboat._modules['console'];
	`
	if err := w.Load("internal:preload", makeGlobals); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (vu *VUContext) Recv(raw string) ***REMOVED***
	call := jsCallEnvelope***REMOVED******REMOVED***
	if err := json.Unmarshal([]byte(raw), &call); err != nil ***REMOVED***
		log.WithError(err).Error("Malformed host call")
		return
	***REMOVED***
	log.WithFields(log.Fields***REMOVED***
		"mod":  call.Mod,
		"fn":   call.Fn,
		"args": call.Args,
	***REMOVED***).Debug("Async call")

	if err := vu.invoke(call); err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't invoke")
	***REMOVED***
***REMOVED***

func (vu *VUContext) RecvSync(raw string) string ***REMOVED***
	call := jsCallEnvelope***REMOVED******REMOVED***
	if err := json.Unmarshal([]byte(raw), &call); err != nil ***REMOVED***
		return jsThrow(fmt.Sprintf("malformed host call: %s", err))
	***REMOVED***
	log.WithFields(log.Fields***REMOVED***
		"mod":  call.Mod,
		"fn":   call.Fn,
		"args": call.Args,
	***REMOVED***).Debug("Sync call")

	if err := vu.invoke(call); err != nil ***REMOVED***
		return jsThrow(err.Error())
	***REMOVED***
	return ""
***REMOVED***

func (vu *VUContext) invoke(call jsCallEnvelope) error ***REMOVED***
	mod, ok := vu.api[call.Mod]
	if !ok ***REMOVED***
		return errors.New(fmt.Sprintf("unknown module '%s'", call.Mod))
	***REMOVED***

	mem, ok := mod[call.Fn]
	if !ok ***REMOVED***
		return errors.New(fmt.Sprintf("unrecognized function call: '%s'.'%s'", call.Mod, call.Fn))
	***REMOVED***

	args := make([]reflect.Value, len(call.Args))
	for i, arg := range call.Args ***REMOVED***
		args[i] = reflect.ValueOf(arg)
	***REMOVED***

	fn := reflect.ValueOf(mem)
	fnT := fn.Type()

	for i := 0; i < fnT.NumIn(); i++ ***REMOVED***
		argT := fnT.In(i)
		switch argT.Kind() ***REMOVED***
		case reflect.Struct:
			mapv, ok := args[i].Interface().(map[string]interface***REMOVED******REMOVED***)
			if !ok ***REMOVED***
				return errors.New("argument is not a dictionary")
			***REMOVED***

			v := reflect.New(argT)
			for i := 0; i < argT.NumField(); i++ ***REMOVED***
				f := argT.Field(i)

				key := f.Tag.Get("json")
				if key == "" ***REMOVED***
					key = f.Name
				***REMOVED***
				val, ok := mapv[key]
				if ok ***REMOVED***
					v.Elem().Field(i).Set(reflect.ValueOf(val))
				***REMOVED***
			***REMOVED***

			args[i] = v.Elem()
		default:
		***REMOVED***
	***REMOVED***

	defer func() ***REMOVED***
		if err := recover(); err != nil ***REMOVED***
			log.WithField("error", err).Error("Go call panicked")
		***REMOVED***
	***REMOVED***()
	ret := fn.Call(args)

	for _, val := range ret ***REMOVED***
		switch v := val.Interface().(type) ***REMOVED***
		case <-chan runner.Result:
		readLoop:
			for ***REMOVED***
				select ***REMOVED***
				case <-vu.ctx.Done():
					break readLoop
				case r, ok := <-v:
					if !ok ***REMOVED***
						break readLoop
					***REMOVED***
					vu.ch <- r
				***REMOVED***
			***REMOVED***
		default:
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
