package afero

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var _ Lstater = (*BasePathFs)(nil)

// The BasePathFs restricts all operations to a given path within an Fs.
// The given file name to the operations on this Fs will be prepended with
// the base path before calling the base Fs.
// Any file name (after filepath.Clean()) outside this base path will be
// treated as non existing file.
//
// Note that it does not clean the error messages on return, so you may
// reveal the real path on errors.
type BasePathFs struct ***REMOVED***
	source Fs
	path   string
***REMOVED***

type BasePathFile struct ***REMOVED***
	File
	path string
***REMOVED***

func (f *BasePathFile) Name() string ***REMOVED***
	sourcename := f.File.Name()
	return strings.TrimPrefix(sourcename, filepath.Clean(f.path))
***REMOVED***

func NewBasePathFs(source Fs, path string) Fs ***REMOVED***
	return &BasePathFs***REMOVED***source: source, path: path***REMOVED***
***REMOVED***

// on a file outside the base path it returns the given file name and an error,
// else the given file with the base path prepended
func (b *BasePathFs) RealPath(name string) (path string, err error) ***REMOVED***
	if err := validateBasePathName(name); err != nil ***REMOVED***
		return name, err
	***REMOVED***

	bpath := filepath.Clean(b.path)
	path = filepath.Clean(filepath.Join(bpath, name))
	if !strings.HasPrefix(path, bpath) ***REMOVED***
		return name, os.ErrNotExist
	***REMOVED***

	return path, nil
***REMOVED***

func validateBasePathName(name string) error ***REMOVED***
	if runtime.GOOS != "windows" ***REMOVED***
		// Not much to do here;
		// the virtual file paths all look absolute on *nix.
		return nil
	***REMOVED***

	// On Windows a common mistake would be to provide an absolute OS path
	// We could strip out the base part, but that would not be very portable.
	if filepath.IsAbs(name) ***REMOVED***
		return os.ErrNotExist
	***REMOVED***

	return nil
***REMOVED***

func (b *BasePathFs) Chtimes(name string, atime, mtime time.Time) (err error) ***REMOVED***
	if name, err = b.RealPath(name); err != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "chtimes", Path: name, Err: err***REMOVED***
	***REMOVED***
	return b.source.Chtimes(name, atime, mtime)
***REMOVED***

func (b *BasePathFs) Chmod(name string, mode os.FileMode) (err error) ***REMOVED***
	if name, err = b.RealPath(name); err != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "chmod", Path: name, Err: err***REMOVED***
	***REMOVED***
	return b.source.Chmod(name, mode)
***REMOVED***

func (b *BasePathFs) Name() string ***REMOVED***
	return "BasePathFs"
***REMOVED***

func (b *BasePathFs) Stat(name string) (fi os.FileInfo, err error) ***REMOVED***
	if name, err = b.RealPath(name); err != nil ***REMOVED***
		return nil, &os.PathError***REMOVED***Op: "stat", Path: name, Err: err***REMOVED***
	***REMOVED***
	return b.source.Stat(name)
***REMOVED***

func (b *BasePathFs) Rename(oldname, newname string) (err error) ***REMOVED***
	if oldname, err = b.RealPath(oldname); err != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "rename", Path: oldname, Err: err***REMOVED***
	***REMOVED***
	if newname, err = b.RealPath(newname); err != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "rename", Path: newname, Err: err***REMOVED***
	***REMOVED***
	return b.source.Rename(oldname, newname)
***REMOVED***

func (b *BasePathFs) RemoveAll(name string) (err error) ***REMOVED***
	if name, err = b.RealPath(name); err != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "remove_all", Path: name, Err: err***REMOVED***
	***REMOVED***
	return b.source.RemoveAll(name)
***REMOVED***

func (b *BasePathFs) Remove(name string) (err error) ***REMOVED***
	if name, err = b.RealPath(name); err != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "remove", Path: name, Err: err***REMOVED***
	***REMOVED***
	return b.source.Remove(name)
***REMOVED***

func (b *BasePathFs) OpenFile(name string, flag int, mode os.FileMode) (f File, err error) ***REMOVED***
	if name, err = b.RealPath(name); err != nil ***REMOVED***
		return nil, &os.PathError***REMOVED***Op: "openfile", Path: name, Err: err***REMOVED***
	***REMOVED***
	sourcef, err := b.source.OpenFile(name, flag, mode)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &BasePathFile***REMOVED***sourcef, b.path***REMOVED***, nil
***REMOVED***

func (b *BasePathFs) Open(name string) (f File, err error) ***REMOVED***
	if name, err = b.RealPath(name); err != nil ***REMOVED***
		return nil, &os.PathError***REMOVED***Op: "open", Path: name, Err: err***REMOVED***
	***REMOVED***
	sourcef, err := b.source.Open(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &BasePathFile***REMOVED***File: sourcef, path: b.path***REMOVED***, nil
***REMOVED***

func (b *BasePathFs) Mkdir(name string, mode os.FileMode) (err error) ***REMOVED***
	if name, err = b.RealPath(name); err != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "mkdir", Path: name, Err: err***REMOVED***
	***REMOVED***
	return b.source.Mkdir(name, mode)
***REMOVED***

func (b *BasePathFs) MkdirAll(name string, mode os.FileMode) (err error) ***REMOVED***
	if name, err = b.RealPath(name); err != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "mkdir", Path: name, Err: err***REMOVED***
	***REMOVED***
	return b.source.MkdirAll(name, mode)
***REMOVED***

func (b *BasePathFs) Create(name string) (f File, err error) ***REMOVED***
	if name, err = b.RealPath(name); err != nil ***REMOVED***
		return nil, &os.PathError***REMOVED***Op: "create", Path: name, Err: err***REMOVED***
	***REMOVED***
	sourcef, err := b.source.Create(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &BasePathFile***REMOVED***File: sourcef, path: b.path***REMOVED***, nil
***REMOVED***

func (b *BasePathFs) LstatIfPossible(name string) (os.FileInfo, bool, error) ***REMOVED***
	name, err := b.RealPath(name)
	if err != nil ***REMOVED***
		return nil, false, &os.PathError***REMOVED***Op: "lstat", Path: name, Err: err***REMOVED***
	***REMOVED***
	if lstater, ok := b.source.(Lstater); ok ***REMOVED***
		return lstater.LstatIfPossible(name)
	***REMOVED***
	fi, err := b.source.Stat(name)
	return fi, false, err
***REMOVED***

// vim: ts=4 sw=4 noexpandtab nolist syn=go
