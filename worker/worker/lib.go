package worker

import (
	"encoding/json"
	"io"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/loader"
	"github.com/spf13/afero"
)

type consoleWriter struct ***REMOVED***
	gs libWorker.BaseGlobalState
***REMOVED***

func (w *consoleWriter) Write(p []byte) (n int, err error) ***REMOVED***
	origLen := len(p)

	// Intercept the write message so can assess log errors parse json
	parsed := make(map[string]interface***REMOVED******REMOVED***)
	if err := json.Unmarshal(p, &parsed); err != nil ***REMOVED***

		return origLen, err
	***REMOVED***

	// Check message level, if error then log error
	if parsed["level"] == "error" ***REMOVED***
		if parsed["error"] != nil ***REMOVED***
			go libWorker.HandleStringError(w.gs, parsed["error"].(string))
		***REMOVED*** else ***REMOVED***
			go libWorker.HandleStringError(w.gs, parsed["msg"].(string))
		***REMOVED***
		return
	***REMOVED***

	go libWorker.DispatchMessage(w.gs, string(p), "CONSOLE")

	return origLen, err
***REMOVED***

var _ io.Writer = &consoleWriter***REMOVED******REMOVED***

type workerLoadedTest struct ***REMOVED***
	sourceRootPath string
	source         *loader.SourceData
	fs             afero.Fs
	pwd            string
	fileSystems    map[string]afero.Fs
	preInitState   *libWorker.TestPreInitState
	initRunner     libWorker.Runner // TODO: rename to something more appropriate
	keyLogger      io.Closer
***REMOVED***

// Config ...
type Config struct ***REMOVED***
	libWorker.Options

	Out []string `json:"out" envconfig:"K6_OUT"`

	// TODO: deprecate
	Collectors map[string]json.RawMessage `json:"collectors"`
***REMOVED***

// loadedAndConfiguredTest contains the whole loadedTest, as well as the
// consolidated test config and the full test run state.
type workerLoadedAndConfiguredTest struct ***REMOVED***
	*workerLoadedTest
	derivedConfig Config
***REMOVED***

const testTypeJS = "js"
