package v8js

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/api"
	"github.com/loadimpact/speedboat/runner/v8js/bridge"
	"github.com/ry/v8worker"
)

type jsCallEnvelope struct ***REMOVED***
	Mod  string        `json:"m"`
	Fn   string        `json:"f"`
	Args []interface***REMOVED******REMOVED*** `json:"a"`
***REMOVED***

func (vu *VUContext) bridgeAPI(w *v8worker.Worker) error ***REMOVED***
	base := api.New()
	vu.api = make(map[string]bridge.Module)

	for modname, mod := range base ***REMOVED***
		bmod := bridge.BridgeModule(modname, mod)
		vu.api[modname] = bmod

		js := bmod.JS()
		println(js)
		if err := w.Load("internal:api:"+modname, js); err != nil ***REMOVED***
			log.WithError(err).WithField("mod", modname).Error("Couldn't bridge module")
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (vu *VUContext) invoke(call jsCallEnvelope) error ***REMOVED***
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
