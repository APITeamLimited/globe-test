package afero

import (
	"os"
	"syscall"
	"time"
)

// If the cache duration is 0, cache time will be unlimited, i.e. once
// a file is in the layer, the base will never be read again for this file.
//
// For cache times greater than 0, the modification time of a file is
// checked. Note that a lot of file system implementations only allow a
// resolution of a second for timestamps... or as the godoc for os.Chtimes()
// states: "The underlying filesystem may truncate or round the values to a
// less precise time unit."
//
// This caching union will forward all write calls also to the base file
// system first. To prevent writing to the base Fs, wrap it in a read-only
// filter - Note: this will also make the overlay read-only, for writing files
// in the overlay, use the overlay Fs directly, not via the union Fs.
type CacheOnReadFs struct ***REMOVED***
	base      Fs
	layer     Fs
	cacheTime time.Duration
***REMOVED***

func NewCacheOnReadFs(base Fs, layer Fs, cacheTime time.Duration) Fs ***REMOVED***
	return &CacheOnReadFs***REMOVED***base: base, layer: layer, cacheTime: cacheTime***REMOVED***
***REMOVED***

type cacheState int

const (
	// not present in the overlay, unknown if it exists in the base:
	cacheMiss cacheState = iota
	// present in the overlay and in base, base file is newer:
	cacheStale
	// present in the overlay - with cache time == 0 it may exist in the base,
	// with cacheTime > 0 it exists in the base and is same age or newer in the
	// overlay
	cacheHit
	// happens if someone writes directly to the overlay without
	// going through this union
	cacheLocal
)

func (u *CacheOnReadFs) cacheStatus(name string) (state cacheState, fi os.FileInfo, err error) ***REMOVED***
	var lfi, bfi os.FileInfo
	lfi, err = u.layer.Stat(name)
	if err == nil ***REMOVED***
		if u.cacheTime == 0 ***REMOVED***
			return cacheHit, lfi, nil
		***REMOVED***
		if lfi.ModTime().Add(u.cacheTime).Before(time.Now()) ***REMOVED***
			bfi, err = u.base.Stat(name)
			if err != nil ***REMOVED***
				return cacheLocal, lfi, nil
			***REMOVED***
			if bfi.ModTime().After(lfi.ModTime()) ***REMOVED***
				return cacheStale, bfi, nil
			***REMOVED***
		***REMOVED***
		return cacheHit, lfi, nil
	***REMOVED***

	if err == syscall.ENOENT || os.IsNotExist(err) ***REMOVED***
		return cacheMiss, nil, nil
	***REMOVED***

	return cacheMiss, nil, err
***REMOVED***

func (u *CacheOnReadFs) copyToLayer(name string) error ***REMOVED***
	return copyToLayer(u.base, u.layer, name)
***REMOVED***

func (u *CacheOnReadFs) Chtimes(name string, atime, mtime time.Time) error ***REMOVED***
	st, _, err := u.cacheStatus(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	switch st ***REMOVED***
	case cacheLocal:
	case cacheHit:
		err = u.base.Chtimes(name, atime, mtime)
	case cacheStale, cacheMiss:
		if err := u.copyToLayer(name); err != nil ***REMOVED***
			return err
		***REMOVED***
		err = u.base.Chtimes(name, atime, mtime)
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return u.layer.Chtimes(name, atime, mtime)
***REMOVED***

func (u *CacheOnReadFs) Chmod(name string, mode os.FileMode) error ***REMOVED***
	st, _, err := u.cacheStatus(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	switch st ***REMOVED***
	case cacheLocal:
	case cacheHit:
		err = u.base.Chmod(name, mode)
	case cacheStale, cacheMiss:
		if err := u.copyToLayer(name); err != nil ***REMOVED***
			return err
		***REMOVED***
		err = u.base.Chmod(name, mode)
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return u.layer.Chmod(name, mode)
***REMOVED***

func (u *CacheOnReadFs) Stat(name string) (os.FileInfo, error) ***REMOVED***
	st, fi, err := u.cacheStatus(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	switch st ***REMOVED***
	case cacheMiss:
		return u.base.Stat(name)
	default: // cacheStale has base, cacheHit and cacheLocal the layer os.FileInfo
		return fi, nil
	***REMOVED***
***REMOVED***

func (u *CacheOnReadFs) Rename(oldname, newname string) error ***REMOVED***
	st, _, err := u.cacheStatus(oldname)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	switch st ***REMOVED***
	case cacheLocal:
	case cacheHit:
		err = u.base.Rename(oldname, newname)
	case cacheStale, cacheMiss:
		if err := u.copyToLayer(oldname); err != nil ***REMOVED***
			return err
		***REMOVED***
		err = u.base.Rename(oldname, newname)
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return u.layer.Rename(oldname, newname)
***REMOVED***

func (u *CacheOnReadFs) Remove(name string) error ***REMOVED***
	st, _, err := u.cacheStatus(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	switch st ***REMOVED***
	case cacheLocal:
	case cacheHit, cacheStale, cacheMiss:
		err = u.base.Remove(name)
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return u.layer.Remove(name)
***REMOVED***

func (u *CacheOnReadFs) RemoveAll(name string) error ***REMOVED***
	st, _, err := u.cacheStatus(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	switch st ***REMOVED***
	case cacheLocal:
	case cacheHit, cacheStale, cacheMiss:
		err = u.base.RemoveAll(name)
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return u.layer.RemoveAll(name)
***REMOVED***

func (u *CacheOnReadFs) OpenFile(name string, flag int, perm os.FileMode) (File, error) ***REMOVED***
	st, _, err := u.cacheStatus(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	switch st ***REMOVED***
	case cacheLocal, cacheHit:
	default:
		if err := u.copyToLayer(name); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	if flag&(os.O_WRONLY|syscall.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_TRUNC) != 0 ***REMOVED***
		bfi, err := u.base.OpenFile(name, flag, perm)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		lfi, err := u.layer.OpenFile(name, flag, perm)
		if err != nil ***REMOVED***
			bfi.Close() // oops, what if O_TRUNC was set and file opening in the layer failed...?
			return nil, err
		***REMOVED***
		return &UnionFile***REMOVED***Base: bfi, Layer: lfi***REMOVED***, nil
	***REMOVED***
	return u.layer.OpenFile(name, flag, perm)
***REMOVED***

func (u *CacheOnReadFs) Open(name string) (File, error) ***REMOVED***
	st, fi, err := u.cacheStatus(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	switch st ***REMOVED***
	case cacheLocal:
		return u.layer.Open(name)

	case cacheMiss:
		bfi, err := u.base.Stat(name)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if bfi.IsDir() ***REMOVED***
			return u.base.Open(name)
		***REMOVED***
		if err := u.copyToLayer(name); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return u.layer.Open(name)

	case cacheStale:
		if !fi.IsDir() ***REMOVED***
			if err := u.copyToLayer(name); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return u.layer.Open(name)
		***REMOVED***
	case cacheHit:
		if !fi.IsDir() ***REMOVED***
			return u.layer.Open(name)
		***REMOVED***
	***REMOVED***
	// the dirs from cacheHit, cacheStale fall down here:
	bfile, _ := u.base.Open(name)
	lfile, err := u.layer.Open(name)
	if err != nil && bfile == nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &UnionFile***REMOVED***Base: bfile, Layer: lfile***REMOVED***, nil
***REMOVED***

func (u *CacheOnReadFs) Mkdir(name string, perm os.FileMode) error ***REMOVED***
	err := u.base.Mkdir(name, perm)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return u.layer.MkdirAll(name, perm) // yes, MkdirAll... we cannot assume it exists in the cache
***REMOVED***

func (u *CacheOnReadFs) Name() string ***REMOVED***
	return "CacheOnReadFs"
***REMOVED***

func (u *CacheOnReadFs) MkdirAll(name string, perm os.FileMode) error ***REMOVED***
	err := u.base.MkdirAll(name, perm)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return u.layer.MkdirAll(name, perm)
***REMOVED***

func (u *CacheOnReadFs) Create(name string) (File, error) ***REMOVED***
	bfh, err := u.base.Create(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	lfh, err := u.layer.Create(name)
	if err != nil ***REMOVED***
		// oops, see comment about OS_TRUNC above, should we remove? then we have to
		// remember if the file did not exist before
		bfh.Close()
		return nil, err
	***REMOVED***
	return &UnionFile***REMOVED***Base: bfh, Layer: lfh***REMOVED***, nil
***REMOVED***
