package fsext

import (
	"time"

	"github.com/spf13/afero"
)

// CacheOnReadFs is wrapper around afero.CacheOnReadFs with the ability to return the filesystem
// that is used as cache
type CacheOnReadFs struct ***REMOVED***
	afero.Fs
	cache afero.Fs
***REMOVED***

// NewCacheOnReadFs returns a new CacheOnReadFs
func NewCacheOnReadFs(base, layer afero.Fs, cacheTime time.Duration) afero.Fs ***REMOVED***
	return CacheOnReadFs***REMOVED***
		Fs:    afero.NewCacheOnReadFs(base, layer, cacheTime),
		cache: layer,
	***REMOVED***
***REMOVED***

// GetCachingFs returns the afero.Fs being used for cache
func (c CacheOnReadFs) GetCachingFs() afero.Fs ***REMOVED***
	return c.cache
***REMOVED***
