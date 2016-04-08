package common

import (
	"github.com/loadimpact/speedboat/runner/js"
	"testing"
)

func GetRunnerJS(t *testing.T) ***REMOVED***
	r, err := GetRunner("script.js")
	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
	if _, ok := r.(*js.JSRunner); !ok ***REMOVED***
		t.Error("Not a JS runner")
	***REMOVED***
***REMOVED***

func GetRunnerUnknown(t *testing.T) ***REMOVED***
	r, err := GetRunner("test.doc")
	if err == nil ***REMOVED***
		t.Error("No error")
	***REMOVED***
	if r != nil ***REMOVED***
		t.Error("Something returned")
	***REMOVED***
***REMOVED***
