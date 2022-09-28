package collection

import (
	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

type (
	variables struct ***REMOVED***
		set   func(string, string) bool
		get   func(string) string
		has   func(string) bool
		unset func(string) bool
		clear func() bool
		list  func() []libWorker.KeyValueItem
	***REMOVED***
)

func (mi *Collection) getVariables() variables ***REMOVED***
	return variables***REMOVED***
		set:   mi.set,
		get:   mi.get,
		has:   mi.has,
		unset: mi.unset,
		clear: mi.clear,
		list:  mi.list,
	***REMOVED***
***REMOVED***

// set sets a key-value pair in the collection.
func (mi *Collection) set(key string, value string) bool ***REMOVED***
	mi.collection.mu.Lock()
	defer mi.collection.mu.Unlock()

	// Overwrite existing value if key already exists
	mi.collection.variables[key] = libWorker.KeyValueItem***REMOVED***
		Key:   key,
		Value: value,
	***REMOVED***

	return true
***REMOVED***

// get gets a value from the collection.
func (mi *Collection) get(key string) string ***REMOVED***
	mi.collection.mu.RLock()
	defer mi.collection.mu.RUnlock()

	if value, ok := mi.collection.variables[key]; ok ***REMOVED***
		return value.Value
	***REMOVED***

	return ""
***REMOVED***

// has checks if a key exists in the collection.
func (mi *Collection) has(key string) bool ***REMOVED***
	mi.collection.mu.RLock()
	defer mi.collection.mu.RUnlock()

	_, ok := mi.collection.variables[key]
	return ok
***REMOVED***

// unset removes a key-value pair from the collection.
func (mi *Collection) unset(key string) bool ***REMOVED***
	mi.collection.mu.Lock()
	defer mi.collection.mu.Unlock()

	if _, ok := mi.collection.variables[key]; ok ***REMOVED***
		delete(mi.collection.variables, key)
		return true
	***REMOVED***

	return false
***REMOVED***

// clear removes all key-value pairs from the collection.
func (mi *Collection) clear() bool ***REMOVED***
	mi.collection.mu.Lock()
	defer mi.collection.mu.Unlock()

	mi.collection.variables = make(map[string]libWorker.KeyValueItem)
	return true
***REMOVED***

// list returns a list of all key-value pairs in the collection.
func (mi *Collection) list() []libWorker.KeyValueItem ***REMOVED***
	mi.collection.mu.RLock()
	defer mi.collection.mu.RUnlock()

	list := make([]libWorker.KeyValueItem, 0, len(mi.collection.variables))
	for _, item := range mi.collection.variables ***REMOVED***
		list = append(list, item)
	***REMOVED***

	return list
***REMOVED***
