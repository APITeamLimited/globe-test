package v8js

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ry/v8worker"
	"reflect"
)

type jsCallEnvelope struct ***REMOVED***
	Mod  string        `json:"m"`
	Fn   string        `json:"f"`
	Args []interface***REMOVED******REMOVED*** `json:"a"`
***REMOVED***

func (vu *VUContext) RegisterModules(w *v8worker.Worker) error ***REMOVED***
	vu.mods = map[string]Module***REMOVED***
		"global": Module***REMOVED***
			"sleep": Member***REMOVED***Func: vu.Sleep***REMOVED***,
		***REMOVED***,
		"console": Module***REMOVED***
			"log":   Member***REMOVED***Func: vu.ConsoleLog, Async: true***REMOVED***,
			"warn":  Member***REMOVED***Func: vu.ConsoleWarn, Async: true***REMOVED***,
			"error": Member***REMOVED***Func: vu.ConsoleError, Async: true***REMOVED***,
		***REMOVED***,
		"http": Module***REMOVED***
			"get": Member***REMOVED***Func: vu.HTTPGet***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for modname, mod := range vu.mods ***REMOVED***
		jsMod := fmt.Sprintf(`
		speedboat._modules["%s"] = ***REMOVED******REMOVED***;
		`, modname)
		for name, mem := range mod ***REMOVED***
			t := reflect.TypeOf(mem.Func)

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
				jsFn += fmt.Sprintf("args.push(speedboat._require.%s(arguments[%d]));", aT.Kind().String(), i)
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
			***REMOVED***`, modname, name, mem.Async)
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
	mod, ok := vu.mods[call.Mod]
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

	defer func() ***REMOVED***
		if err := recover(); err != nil ***REMOVED***
			log.WithField("error", err).Error("Go call panicked")
		***REMOVED***
	***REMOVED***()
	fn := reflect.ValueOf(mem.Func)
	log.WithField("T", fn.Type().String()).Debug("Function")
	fn.Call(args)

	return nil
***REMOVED***
