package runner

import (
	"errors"
	"github.com/loadimpact/speedboat/runner/js"
	"path"
)

func Get(filename string) (Runner, error) ***REMOVED***
	switch path.Ext(filename) ***REMOVED***
	case "js":
		return js.New()
	default:
		return nil, errors.New("No runner found")
	***REMOVED***
***REMOVED***
