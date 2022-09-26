package worker

import (
	"context"
	"encoding/json"
	"io"
	"sync"

	"github.com/APITeamLimited/redis/v9"
	"github.com/spf13/afero"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/loader"
	"gopkg.in/guregu/null.v3"
)

type consoleWriter struct {
	ctx      context.Context
	client   *redis.Client
	jobId    string
	workerId string
}

type workerLoadedTest struct {
	sourceRootPath string
	source         *loader.SourceData
	fs             afero.Fs
	pwd            string
	fileSystems    map[string]afero.Fs
	preInitState   *lib.TestPreInitState
	initRunner     lib.Runner // TODO: rename to something more appropriate
	keyLogger      io.Closer
}

// Config ...
type Config struct {
	lib.Options

	Out           []string  `json:"out" envconfig:"K6_OUT"`
	Linger        null.Bool `json:"linger" envconfig:"K6_LINGER"`
	NoUsageReport null.Bool `json:"noUsageReport" envconfig:"K6_NO_USAGE_REPORT"`

	// TODO: deprecate
	Collectors map[string]json.RawMessage `json:"collectors"`
}

// loadedAndConfiguredTest contains the whole loadedTest, as well as the
// consolidated test config and the full test run state.
type workerLoadedAndConfiguredTest struct {
	*workerLoadedTest
	consolidatedConfig Config
	derivedConfig      Config
}

type syncWriter struct {
	w io.Writer
	m sync.Mutex
}

func (cw *syncWriter) Write(b []byte) (int, error) {
	cw.m.Lock()
	defer cw.m.Unlock()
	return cw.w.Write(b)
}

const (
	testTypeJS = "js"
)
