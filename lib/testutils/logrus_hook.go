package testutils

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

// SimpleLogrusHook implements the logrus.Hook interface and could be used to check
// if log messages were outputted
type SimpleLogrusHook struct ***REMOVED***
	HookedLevels []log.Level
	mutex        sync.Mutex
	messageCache []log.Entry
***REMOVED***

// Levels just returns whatever was stored in the HookedLevels slice
func (smh *SimpleLogrusHook) Levels() []log.Level ***REMOVED***
	return smh.HookedLevels
***REMOVED***

// Fire saves whatever message the logrus library passed in the cache
func (smh *SimpleLogrusHook) Fire(e *log.Entry) error ***REMOVED***
	smh.mutex.Lock()
	defer smh.mutex.Unlock()
	smh.messageCache = append(smh.messageCache, *e)
	return nil
***REMOVED***

// Drain returns the currently stored messages and deletes them from the cache
func (smh *SimpleLogrusHook) Drain() []log.Entry ***REMOVED***
	smh.mutex.Lock()
	defer smh.mutex.Unlock()
	res := smh.messageCache
	smh.messageCache = []log.Entry***REMOVED******REMOVED***
	return res
***REMOVED***

var _ log.Hook = &SimpleLogrusHook***REMOVED******REMOVED***
