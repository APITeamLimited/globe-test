package common

import (
	"errors"
	"github.com/loadimpact/speedboat/runner"
	"github.com/loadimpact/speedboat/runner/js"
	"path"
)

func GetRunner(filename string) (runner.Runner, error) ***REMOVED***
	switch path.Ext(filename) ***REMOVED***
	case ".js":
		return js.New()
	default:
		return nil, errors.New("No runner found")
	***REMOVED***
***REMOVED***
