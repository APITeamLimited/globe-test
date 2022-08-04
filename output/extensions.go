package output

import (
	"fmt"
	"sync"
)

//nolint:gochecknoglobals
var (
	extensions = make(map[string]func(Params) (Output, error))
	mx         sync.RWMutex
)

// GetExtensions returns all registered extensions.
func GetExtensions() map[string]func(Params) (Output, error) ***REMOVED***
	mx.RLock()
	defer mx.RUnlock()
	res := make(map[string]func(Params) (Output, error), len(extensions))
	for k, v := range extensions ***REMOVED***
		res[k] = v
	***REMOVED***
	return res
***REMOVED***

// RegisterExtension registers the given output extension constructor. This
// function panics if a module with the same name is already registered.
func RegisterExtension(name string, mod func(Params) (Output, error)) ***REMOVED***
	mx.Lock()
	defer mx.Unlock()

	if _, ok := extensions[name]; ok ***REMOVED***
		panic(fmt.Sprintf("output extension already registered: %s", name))
	***REMOVED***
	extensions[name] = mod
***REMOVED***
