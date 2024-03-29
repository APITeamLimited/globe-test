package worker

import (
	"encoding/json"
	"io"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/loader"
	"github.com/spf13/afero"
)

type consoleWriter struct {
	gs libWorker.BaseGlobalState

	loggerChannel chan map[string]interface{}
}

func (w *consoleWriter) Write(p []byte) (n int, err error) {
	origLen := len(p)

	// Intercept the write message so can assess log errors parse json
	parsed := make(map[string]interface{})
	if err := json.Unmarshal(p, &parsed); err != nil {
		return origLen, err
	}

	// Ignore debug messages
	if parsed["level"] == "debug" {
		return origLen, err
	}

	w.loggerChannel <- parsed

	// Check message level, if error then log error
	if parsed["level"] == "error" {
		if parsed["error"] != nil {
			libWorker.HandleStringError(w.gs, parsed["error"].(string))
		} else {
			libWorker.HandleStringError(w.gs, parsed["msg"].(string))
		}
		return
	}

	return origLen, err
}

var _ io.Writer = &consoleWriter{}

type workerLoadedTest struct {
	fs           afero.Fs
	pwd          string
	fileSystems  map[string]afero.Fs
	preInitState *libWorker.TestPreInitState
	initRunner   libWorker.Runner // TODO: rename to something more appropriate
	keyLogger    io.Closer
	sourceData   *[]*loader.SourceData
}

// Config ...
type Config struct {
	libWorker.Options

	Out []string `json:"out" envconfig:"K6_OUT"`

	// TODO: deprecate
	Collectors map[string]json.RawMessage `json:"collectors"`
}

// loadedAndConfiguredTest contains the whole loadedTest, as well as the
// consolidated test config and the full test run state.
type workerLoadedAndConfiguredTest struct {
	*workerLoadedTest
	derivedConfig Config
}

const testTypeJS = "js"

type JobUserUpdate struct {
	UpdateType string `json:"updateType"`
}
