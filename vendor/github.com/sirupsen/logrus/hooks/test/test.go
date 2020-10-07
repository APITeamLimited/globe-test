// The Test package is used for testing logrus.
// It provides a simple hooks which register logged messages.
package test

import (
	"io/ioutil"
	"sync"

	"github.com/sirupsen/logrus"
)

// Hook is a hook designed for dealing with logs in test scenarios.
type Hook struct ***REMOVED***
	// Entries is an array of all entries that have been received by this hook.
	// For safe access, use the AllEntries() method, rather than reading this
	// value directly.
	Entries []logrus.Entry
	mu      sync.RWMutex
***REMOVED***

// NewGlobal installs a test hook for the global logger.
func NewGlobal() *Hook ***REMOVED***

	hook := new(Hook)
	logrus.AddHook(hook)

	return hook

***REMOVED***

// NewLocal installs a test hook for a given local logger.
func NewLocal(logger *logrus.Logger) *Hook ***REMOVED***

	hook := new(Hook)
	logger.Hooks.Add(hook)

	return hook

***REMOVED***

// NewNullLogger creates a discarding logger and installs the test hook.
func NewNullLogger() (*logrus.Logger, *Hook) ***REMOVED***

	logger := logrus.New()
	logger.Out = ioutil.Discard

	return logger, NewLocal(logger)

***REMOVED***

func (t *Hook) Fire(e *logrus.Entry) error ***REMOVED***
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Entries = append(t.Entries, *e)
	return nil
***REMOVED***

func (t *Hook) Levels() []logrus.Level ***REMOVED***
	return logrus.AllLevels
***REMOVED***

// LastEntry returns the last entry that was logged or nil.
func (t *Hook) LastEntry() *logrus.Entry ***REMOVED***
	t.mu.RLock()
	defer t.mu.RUnlock()
	i := len(t.Entries) - 1
	if i < 0 ***REMOVED***
		return nil
	***REMOVED***
	return &t.Entries[i]
***REMOVED***

// AllEntries returns all entries that were logged.
func (t *Hook) AllEntries() []*logrus.Entry ***REMOVED***
	t.mu.RLock()
	defer t.mu.RUnlock()
	// Make a copy so the returned value won't race with future log requests
	entries := make([]*logrus.Entry, len(t.Entries))
	for i := 0; i < len(t.Entries); i++ ***REMOVED***
		// Make a copy, for safety
		entries[i] = &t.Entries[i]
	***REMOVED***
	return entries
***REMOVED***

// Reset removes all Entries from this test hook.
func (t *Hook) Reset() ***REMOVED***
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Entries = make([]logrus.Entry, 0)
***REMOVED***
