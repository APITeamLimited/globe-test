package collection

import (
	"go.k6.io/k6/lib"
)

type (
	variables struct {
		set   func(string, string) bool
		get   func(string) string
		has   func(string) bool
		unset func(string) bool
		clear func() bool
		list  func() []lib.KeyValueItem
	}
)

func (mi *Collection) getVariables() variables {
	return variables{
		set:   mi.set,
		get:   mi.get,
		has:   mi.has,
		unset: mi.unset,
		clear: mi.clear,
		list:  mi.list,
	}
}

// set sets a key-value pair in the collection.
func (mi *Collection) set(key string, value string) bool {
	mi.collection.mu.Lock()
	defer mi.collection.mu.Unlock()

	// Overwrite existing value if key already exists
	mi.collection.variables[key] = lib.KeyValueItem{
		Key:   key,
		Value: value,
	}

	return true
}

// get gets a value from the collection.
func (mi *Collection) get(key string) string {
	mi.collection.mu.RLock()
	defer mi.collection.mu.RUnlock()

	if value, ok := mi.collection.variables[key]; ok {
		return value.Value
	}

	return ""
}

// has checks if a key exists in the collection.
func (mi *Collection) has(key string) bool {
	mi.collection.mu.RLock()
	defer mi.collection.mu.RUnlock()

	_, ok := mi.collection.variables[key]
	return ok
}

// unset removes a key-value pair from the collection.
func (mi *Collection) unset(key string) bool {
	mi.collection.mu.Lock()
	defer mi.collection.mu.Unlock()

	if _, ok := mi.collection.variables[key]; ok {
		delete(mi.collection.variables, key)
		return true
	}

	return false
}

// clear removes all key-value pairs from the collection.
func (mi *Collection) clear() bool {
	mi.collection.mu.Lock()
	defer mi.collection.mu.Unlock()

	mi.collection.variables = make(map[string]lib.KeyValueItem)
	return true
}

// list returns a list of all key-value pairs in the collection.
func (mi *Collection) list() []lib.KeyValueItem {
	mi.collection.mu.RLock()
	defer mi.collection.mu.RUnlock()

	list := make([]lib.KeyValueItem, 0, len(mi.collection.variables))
	for _, item := range mi.collection.variables {
		list = append(list, item)
	}

	return list
}
