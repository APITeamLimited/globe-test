package test

import (
	"github.com/loadimpact/speedboat/runner"
)

var Members = map[string]interface***REMOVED******REMOVED******REMOVED***
	"abort": Abort,
***REMOVED***

func Abort() <-chan runner.Result ***REMOVED***
	ch := make(chan runner.Result, 1)
	ch <- runner.Result***REMOVED***Abort: true***REMOVED***
	return ch
***REMOVED***
