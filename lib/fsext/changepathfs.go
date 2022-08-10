package fsext

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/afero"
)

var _ afero.Lstater = (*ChangePathFs)(nil)

// ChangePathFs is a filesystem that wraps another afero.Fs and changes all given paths from all
// file and directory names, with a function, before calling the same method on the wrapped afero.Fs.
// Heavily based on afero.BasePathFs
type ChangePathFs struct ***REMOVED***
	source afero.Fs
	fn     ChangePathFunc
***REMOVED***

// ChangePathFile is a file from ChangePathFs
type ChangePathFile struct ***REMOVED***
	afero.File
	originalName string
***REMOVED***

// NewChangePathFs return a ChangePathFs where all paths will be change with the provided funcs
func NewChangePathFs(source afero.Fs, fn ChangePathFunc) *ChangePathFs ***REMOVED***
	return &ChangePathFs***REMOVED***source: source, fn: fn***REMOVED***
***REMOVED***

// ChangePathFunc is the function that will be called by ChangePathFs to change the path
type ChangePathFunc func(name string) (path string, err error)

// NewTrimFilePathSeparatorFs is ChangePathFs that trims a Afero.FilePathSeparator from all paths
// Heavily based on afero.BasePathFs
func NewTrimFilePathSeparatorFs(source afero.Fs) *ChangePathFs ***REMOVED***
	return &ChangePathFs***REMOVED***source: source, fn: ChangePathFunc(func(name string) (path string, err error) ***REMOVED***
		if !strings.HasPrefix(name, afero.FilePathSeparator) ***REMOVED***
			return name, os.ErrNotExist
		***REMOVED***

		return filepath.Clean(strings.TrimPrefix(name, afero.FilePathSeparator)), nil
	***REMOVED***)***REMOVED***
***REMOVED***

// Name Returns the name of the file
func (f *ChangePathFile) Name() string ***REMOVED***
	return f.originalName
***REMOVED***

// Chown changes the uid and gid of the named file.
func (b *ChangePathFs) Chown(name string, uid, gid int) error ***REMOVED***
	return errors.New("unimplemented Chown")
***REMOVED***

// Chtimes changes the access and modification times of the named file
func (b *ChangePathFs) Chtimes(name string, atime, mtime time.Time) (err error) ***REMOVED***
	var newName string
	if newName, err = b.fn(name); err != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "chtimes", Path: name, Err: err***REMOVED***
	***REMOVED***
	return b.source.Chtimes(newName, atime, mtime)
***REMOVED***

// Chmod changes the mode of the named file to mode.
func (b *ChangePathFs) Chmod(name string, mode os.FileMode) (err error) ***REMOVED***
	var newName string
	if newName, err = b.fn(name); err != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "chmod", Path: name, Err: err***REMOVED***
	***REMOVED***
	return b.source.Chmod(newName, mode)
***REMOVED***

// Name return the name of this FileSystem
func (b *ChangePathFs) Name() string ***REMOVED***
	return "ChangePathFs"
***REMOVED***

// Stat returns a FileInfo describing the named file, or an error, if any
// happens.
func (b *ChangePathFs) Stat(name string) (fi os.FileInfo, err error) ***REMOVED***
	var newName string
	if newName, err = b.fn(name); err != nil ***REMOVED***
		return nil, &os.PathError***REMOVED***Op: "stat", Path: name, Err: err***REMOVED***
	***REMOVED***
	return b.source.Stat(newName)
***REMOVED***

// Rename renames a file.
func (b *ChangePathFs) Rename(oldName, newName string) (err error) ***REMOVED***
	var newOldName, newNewName string
	if newOldName, err = b.fn(oldName); err != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "rename", Path: oldName, Err: err***REMOVED***
	***REMOVED***
	if newNewName, err = b.fn(newName); err != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "rename", Path: newName, Err: err***REMOVED***
	***REMOVED***
	return b.source.Rename(newOldName, newNewName)
***REMOVED***

// RemoveAll removes a directory path and any children it contains. It
// does not fail if the path does not exist (return nil).
func (b *ChangePathFs) RemoveAll(name string) (err error) ***REMOVED***
	var newName string
	if newName, err = b.fn(name); err != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "remove_all", Path: name, Err: err***REMOVED***
	***REMOVED***
	return b.source.RemoveAll(newName)
***REMOVED***

// Remove removes a file identified by name, returning an error, if any
// happens.
func (b *ChangePathFs) Remove(name string) (err error) ***REMOVED***
	var newName string
	if newName, err = b.fn(name); err != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "remove", Path: name, Err: err***REMOVED***
	***REMOVED***
	return b.source.Remove(newName)
***REMOVED***

// OpenFile opens a file using the given flags and the given mode.
func (b *ChangePathFs) OpenFile(name string, flag int, mode os.FileMode) (f afero.File, err error) ***REMOVED***
	var newName string
	if newName, err = b.fn(name); err != nil ***REMOVED***
		return nil, &os.PathError***REMOVED***Op: "openfile", Path: name, Err: err***REMOVED***
	***REMOVED***
	sourcef, err := b.source.OpenFile(newName, flag, mode)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &ChangePathFile***REMOVED***File: sourcef, originalName: name***REMOVED***, nil
***REMOVED***

// Open opens a file, returning it or an error, if any happens.
func (b *ChangePathFs) Open(name string) (f afero.File, err error) ***REMOVED***
	var newName string
	if newName, err = b.fn(name); err != nil ***REMOVED***
		return nil, &os.PathError***REMOVED***Op: "open", Path: name, Err: err***REMOVED***
	***REMOVED***
	sourcef, err := b.source.Open(newName)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &ChangePathFile***REMOVED***File: sourcef, originalName: name***REMOVED***, nil
***REMOVED***

// Mkdir creates a directory in the filesystem, return an error if any
// happens.
func (b *ChangePathFs) Mkdir(name string, mode os.FileMode) (err error) ***REMOVED***
	var newName string
	if newName, err = b.fn(name); err != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "mkdir", Path: name, Err: err***REMOVED***
	***REMOVED***
	return b.source.Mkdir(newName, mode)
***REMOVED***

// MkdirAll creates a directory path and all parents that does not exist
// yet.
func (b *ChangePathFs) MkdirAll(name string, mode os.FileMode) (err error) ***REMOVED***
	var newName string
	if newName, err = b.fn(name); err != nil ***REMOVED***
		return &os.PathError***REMOVED***Op: "mkdir", Path: name, Err: err***REMOVED***
	***REMOVED***
	return b.source.MkdirAll(newName, mode)
***REMOVED***

// Create creates a file in the filesystem, returning the file and an
// error, if any happens
func (b *ChangePathFs) Create(name string) (f afero.File, err error) ***REMOVED***
	var newName string
	if newName, err = b.fn(name); err != nil ***REMOVED***
		return nil, &os.PathError***REMOVED***Op: "create", Path: name, Err: err***REMOVED***
	***REMOVED***
	sourcef, err := b.source.Create(newName)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &ChangePathFile***REMOVED***File: sourcef, originalName: name***REMOVED***, nil
***REMOVED***

// LstatIfPossible implements the afero.Lstater interface
func (b *ChangePathFs) LstatIfPossible(name string) (os.FileInfo, bool, error) ***REMOVED***
	var newName string
	newName, err := b.fn(name)
	if err != nil ***REMOVED***
		return nil, false, &os.PathError***REMOVED***Op: "lstat", Path: name, Err: err***REMOVED***
	***REMOVED***
	if lstater, ok := b.source.(afero.Lstater); ok ***REMOVED***
		return lstater.LstatIfPossible(newName)
	***REMOVED***
	fi, err := b.source.Stat(newName)
	return fi, false, err
***REMOVED***
