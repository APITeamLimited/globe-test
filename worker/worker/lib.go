package worker

import (
	"context"
	"encoding/json"
	"io"

	"github.com/APITeamLimited/k6-worker/lib"
	"github.com/APITeamLimited/k6-worker/loader"
	"github.com/APITeamLimited/redis/v9"
	"github.com/spf13/afero"
	"gopkg.in/guregu/null.v3"
)

type consoleWriter struct ***REMOVED***
	ctx      context.Context
	client   *redis.Client
	jobId    string
	workerId string
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
			go lib.HandleStringError(w.ctx, w.client, w.jobId, w.workerId, parsed["error"].(string))
		***REMOVED*** else ***REMOVED***
			go lib.HandleStringError(w.ctx, w.client, w.jobId, w.workerId, parsed["msg"].(string))
		***REMOVED***
		return
	***REMOVED***

	go lib.DispatchMessage(w.ctx, w.client, w.jobId, w.workerId, string(p), "CONSOLE")

	return origLen, err
***REMOVED***

var _ io.Writer = &consoleWriter***REMOVED******REMOVED***

type workerLoadedTest struct ***REMOVED***
	sourceRootPath string
	source         *loader.SourceData
	fs             afero.Fs
	pwd            string
	fileSystems    map[string]afero.Fs
	preInitState   *lib.TestPreInitState
	initRunner     lib.Runner // TODO: rename to something more appropriate
	keyLogger      io.Closer
***REMOVED***

// Config ...
type Config struct ***REMOVED***
	lib.Options

	Out           []string  `json:"out" envconfig:"K6_OUT"`
	Linger        null.Bool `json:"linger" envconfig:"K6_INGER"`
	NoUsageReport null.Bool `json:"noUsageReport" envconfig:"K6_NO_USAGE_REPORT"`

	// TODO: deprecate
	Collectors map[string]json.RawMessage `json:"collectors"`
***REMOVED***

// loadedAndConfiguredTest contains the whole loadedTest, as well as the
// consolidated test config and the full test run state.
type workerLoadedAndConfiguredTest struct ***REMOVED***
	*workerLoadedTest
	consolidatedConfig Config
	derivedConfig      Config
***REMOVED***

const testTypeJS = "js"
