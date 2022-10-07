package collection

type (
	variables struct ***REMOVED***
		set   func(string, string) bool
		get   func(string) string
		has   func(string) bool
		unset func(string) bool
		clear func() bool
		list  func() []string
	***REMOVED***
)

func (mi *CollectionInstance) getVariables() variables ***REMOVED***
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
func (mi *CollectionInstance) set(key string, value string) bool ***REMOVED***
	mi.module.sharedCollection.mu.Lock()
	defer mi.module.sharedCollection.mu.Unlock()

	// Overwrite existing value if key already exists
	mi.module.sharedCollection.data.Variables[key] = value

	return true
***REMOVED***

// get gets a value from the collection.
func (mi *CollectionInstance) get(key string) string ***REMOVED***
	mi.module.sharedCollection.mu.RLock()
	defer mi.module.sharedCollection.mu.RUnlock()

	if value, ok := mi.module.sharedCollection.data.Variables[key]; ok ***REMOVED***
		return value
	***REMOVED***

	return ""
***REMOVED***

// has checks if a key exists in the collection.
func (mi *CollectionInstance) has(key string) bool ***REMOVED***
	mi.module.sharedCollection.mu.RLock()
	defer mi.module.sharedCollection.mu.RUnlock()

	_, ok := mi.module.sharedCollection.data.Variables[key]
	return ok
***REMOVED***

// unset removes a key-value pair from the collection.
func (mi *CollectionInstance) unset(key string) bool ***REMOVED***
	mi.module.sharedCollection.mu.Lock()
	defer mi.module.sharedCollection.mu.Unlock()

	if _, ok := mi.module.sharedCollection.data.Variables[key]; ok ***REMOVED***
		delete(mi.module.sharedCollection.data.Variables, key)
		return true
	***REMOVED***

	return false
***REMOVED***

// clear removes all key-value pairs from the collection.
func (mi *CollectionInstance) clear() bool ***REMOVED***
	mi.module.sharedCollection.mu.Lock()
	defer mi.module.sharedCollection.mu.Unlock()

	mi.module.sharedCollection.data.Variables = make(map[string]string)
	return true
***REMOVED***

// list returns a list of all key-value pairs in the collection.
func (mi *CollectionInstance) list() []string ***REMOVED***
	mi.module.sharedCollection.mu.RLock()
	defer mi.module.sharedCollection.mu.RUnlock()

	list := make([]string, 0, len(mi.module.sharedCollection.data.Variables))
	for _, item := range mi.module.sharedCollection.data.Variables ***REMOVED***
		list = append(list, item)
	***REMOVED***

	return list
***REMOVED***
