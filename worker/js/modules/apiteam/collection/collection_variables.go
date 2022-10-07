package collection

type (
	variables struct {
		set   func(string, string) bool
		get   func(string) string
		has   func(string) bool
		unset func(string) bool
		clear func() bool
		list  func() []string
	}
)

func (mi *CollectionInstance) getVariables() variables {
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

// list returns a list of all key-value pairs in the collection.
func (mi *CollectionInstance) list() []string {
	mi.module.sharedCollection.mu.RLock()
	defer mi.module.sharedCollection.mu.RUnlock()

	list := make([]string, 0, len(mi.module.sharedCollection.data.Variables))
	for _, item := range mi.module.sharedCollection.data.Variables {
		list = append(list, item)
	}

	return list
}
