package redis

import (
	"context"
	"sync"
)

// ScanIterator is used to incrementally iterate over a collection of elements.
// It's safe for concurrent use by multiple goroutines.
type ScanIterator struct ***REMOVED***
	mu  sync.Mutex // protects Scanner and pos
	cmd *ScanCmd
	pos int
***REMOVED***

// Err returns the last iterator error, if any.
func (it *ScanIterator) Err() error ***REMOVED***
	it.mu.Lock()
	err := it.cmd.Err()
	it.mu.Unlock()
	return err
***REMOVED***

// Next advances the cursor and returns true if more values can be read.
func (it *ScanIterator) Next(ctx context.Context) bool ***REMOVED***
	it.mu.Lock()
	defer it.mu.Unlock()

	// Instantly return on errors.
	if it.cmd.Err() != nil ***REMOVED***
		return false
	***REMOVED***

	// Advance cursor, check if we are still within range.
	if it.pos < len(it.cmd.page) ***REMOVED***
		it.pos++
		return true
	***REMOVED***

	for ***REMOVED***
		// Return if there is no more data to fetch.
		if it.cmd.cursor == 0 ***REMOVED***
			return false
		***REMOVED***

		// Fetch next page.
		switch it.cmd.args[0] ***REMOVED***
		case "scan", "qscan":
			it.cmd.args[1] = it.cmd.cursor
		default:
			it.cmd.args[2] = it.cmd.cursor
		***REMOVED***

		err := it.cmd.process(ctx, it.cmd)
		if err != nil ***REMOVED***
			return false
		***REMOVED***

		it.pos = 1

		// Redis can occasionally return empty page.
		if len(it.cmd.page) > 0 ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
***REMOVED***

// Val returns the key/field at the current cursor position.
func (it *ScanIterator) Val() string ***REMOVED***
	var v string
	it.mu.Lock()
	if it.cmd.Err() == nil && it.pos > 0 && it.pos <= len(it.cmd.page) ***REMOVED***
		v = it.cmd.page[it.pos-1]
	***REMOVED***
	it.mu.Unlock()
	return v
***REMOVED***
