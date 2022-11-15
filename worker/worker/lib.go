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
}

func (w *consoleWriter) Write(p []byte) (n int, err error) {
	origLen := len(p)

	// Intercept the write message so can assess log errors parse json
	parsed := make(map[string]interface{})
	if err := json.Unmarshal(p, &parsed); err != nil {

		return origLen, err
	}

	// Check message level, if error then log error
	if parsed["level"] == "error" {
		if parsed["error"] != nil {
			go libWorker.HandleStringError(w.gs, parsed["error"].(string))
		} else {
			go libWorker.HandleStringError(w.gs, parsed["msg"].(string))
		}
		return
	}

	libWorker.DispatchMessage(w.gs, string(p), "CONSOLE")

	return origLen, err
}

var _ io.Writer = &consoleWriter{}

type workerLoadedTest struct {
	sourceRootPath string
	source         *loader.SourceData
	fs             afero.Fs
	pwd            string
	fileSystems    map[string]afero.Fs
	preInitState   *libWorker.TestPreInitState
	initRunner     libWorker.Runner // TODO: rename to something more appropriate
	keyLogger      io.Closer
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

type jobUserUpdate struct {
	UpdateType string `json:"updateType"`
}
