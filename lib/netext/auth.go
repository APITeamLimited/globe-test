package netext

import (
	"sync"
)

type AuthCache struct ***REMOVED***
	sync.Mutex

	entries map[string]string
***REMOVED***

func NewAuthCache() *AuthCache ***REMOVED***
	return &AuthCache***REMOVED***
		entries: make(map[string]string),
	***REMOVED***
***REMOVED***

func (a *AuthCache) Set(key, value string) ***REMOVED***
	a.Lock()
	defer a.Unlock()

	a.entries[key] = value
***REMOVED***

func (a AuthCache) Get(key string) (string, bool) ***REMOVED***
	a.Lock()
	defer a.Unlock()

	value, ok := a.entries[key]
	return value, ok
***REMOVED***
