package node

import (
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/spf13/afero"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/loader"
	"gopkg.in/guregu/null.v3"
)

type consoleWriter struct ***REMOVED***
	rawOut *os.File
	writer io.Writer
	isTTY  bool
	mutex  *sync.Mutex

	// Used for flicker-free persistent objects like the progressbars
	persistentText func()
***REMOVED***

type nodeLoadedTest struct ***REMOVED***
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
	Linger        null.Bool `json:"linger" envconfig:"K6_LINGER"`
	NoUsageReport null.Bool `json:"noUsageReport" envconfig:"K6_NO_USAGE_REPORT"`

	// TODO: deprecate
	Collectors map[string]json.RawMessage `json:"collectors"`
***REMOVED***

// loadedAndConfiguredTest contains the whole loadedTest, as well as the
// consolidated test config and the full test run state.
type nodeLoadedAndConfiguredTest struct ***REMOVED***
	*nodeLoadedTest
	consolidatedConfig Config
	derivedConfig      Config
***REMOVED***

type syncWriter struct ***REMOVED***
	w io.Writer
	m sync.Mutex
***REMOVED***

func (cw *syncWriter) Write(b []byte) (int, error) ***REMOVED***
	cw.m.Lock()
	defer cw.m.Unlock()
	return cw.w.Write(b)
***REMOVED***

const (
	testTypeJS      = "js"
	testTypeArchive = "archive"
)
