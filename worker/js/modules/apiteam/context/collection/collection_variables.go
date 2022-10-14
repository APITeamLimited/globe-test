package collection

import (
	"github.com/APITeamLimited/globe-test/worker/js/common"
	"github.com/dop251/goja"
)

func (mi *CollectionInstance) getVariablesObject() *goja.Object {
	// Make sure variables are enabled before calling

	rt := mi.vu.Runtime()
	variablesObject := rt.NewObject()

	mustExport := func(name string, value interface{}) {
		if err := variablesObject.Set(name, value); err != nil {
			common.Throw(rt, err)
		}
	}

	mustExport("set", mi.set)
	mustExport("get", mi.get)
	mustExport("has", mi.has)
	mustExport("unset", mi.unset)
	mustExport("clear", mi.clear)
	mustExport("list", mi.list)

	return variablesObject
}

// set sets a key-value pair in the collection.
func (mi *CollectionInstance) set(key string, value string) bool {
	mi.module.sharedCollection.mu.Lock()
	defer mi.module.sharedCollection.mu.Unlock()

	// Overwrite existing value if key already exists
	mi.module.sharedCollection.data.Variables[key] = value

	return true
}

// get gets a value from the collection.
func (mi *CollectionInstance) get(key string) string {
	mi.module.sharedCollection.mu.RLock()
	defer mi.module.sharedCollection.mu.RUnlock()

	if value, ok := mi.module.sharedCollection.data.Variables[key]; ok {
		return value
	}

	return ""
}

// has checks if a key exists in the collection.
func (mi *CollectionInstance) has(key string) bool {
	mi.module.sharedCollection.mu.RLock()
	defer mi.module.sharedCollection.mu.RUnlock()

	_, ok := mi.module.sharedCollection.data.Variables[key]

	return ok
}

// unset removes a key-value pair from the collection.
func (mi *CollectionInstance) unset(key string) bool {
	mi.module.sharedCollection.mu.Lock()
	defer mi.module.sharedCollection.mu.Unlock()

	if _, ok := mi.module.sharedCollection.data.Variables[key]; ok {
		delete(mi.module.sharedCollection.data.Variables, key)
		return true
	}

	return false
}

// clear removes all key-value pairs from the collection.
func (mi *CollectionInstance) clear() bool {
	mi.module.sharedCollection.mu.Lock()
	defer mi.module.sharedCollection.mu.Unlock()

	mi.module.sharedCollection.data.Variables = make(map[string]string)
	return true
}

// list returns a dictionary of all key-value pairs in the collection.
func (mi *CollectionInstance) list() map[string]string {
	mi.module.sharedCollection.mu.RLock()
	defer mi.module.sharedCollection.mu.RUnlock()

	return mi.module.sharedCollection.data.Variables
}
