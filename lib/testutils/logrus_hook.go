package testutils

import (
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// SimpleLogrusHook implements the logrus.Hook interface and could be used to check
// if log messages were outputted
type SimpleLogrusHook struct ***REMOVED***
	HookedLevels []logrus.Level
	mutex        sync.Mutex
	messageCache []logrus.Entry
***REMOVED***

// Levels just returns whatever was stored in the HookedLevels slice
func (smh *SimpleLogrusHook) Levels() []logrus.Level ***REMOVED***
	return smh.HookedLevels
***REMOVED***

// Fire saves whatever message the logrus library passed in the cache
func (smh *SimpleLogrusHook) Fire(e *logrus.Entry) error ***REMOVED***
	smh.mutex.Lock()
	defer smh.mutex.Unlock()
	smh.messageCache = append(smh.messageCache, *e)
	return nil
***REMOVED***

// Drain returns the currently stored messages and deletes them from the cache
func (smh *SimpleLogrusHook) Drain() []logrus.Entry ***REMOVED***
	smh.mutex.Lock()
	defer smh.mutex.Unlock()
	res := smh.messageCache
	smh.messageCache = []logrus.Entry***REMOVED******REMOVED***
	return res
***REMOVED***

var _ logrus.Hook = &SimpleLogrusHook***REMOVED******REMOVED***

// LogContains is a helper function that checks the provided list of log entries
// for a message matching the provided level and contents.
func LogContains(logEntries []logrus.Entry, expLevel logrus.Level, expContents string) bool ***REMOVED***
	for _, entry := range logEntries ***REMOVED***
		if entry.Level == expLevel && strings.Contains(entry.Message, expContents) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
