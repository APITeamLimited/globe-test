package afero

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

var _ Lstater = (*CopyOnWriteFs)(nil)

// The CopyOnWriteFs is a union filesystem: a read only base file system with
// a possibly writeable layer on top. Changes to the file system will only
// be made in the overlay: Changing an existing file in the base layer which
// is not present in the overlay will copy the file to the overlay ("changing"
// includes also calls to e.g. Chtimes() and Chmod()).
//
// Reading directories is currently only supported via Open(), not OpenFile().
type CopyOnWriteFs struct ***REMOVED***
	base  Fs
	layer Fs
***REMOVED***

func NewCopyOnWriteFs(base Fs, layer Fs) Fs ***REMOVED***
	return &CopyOnWriteFs***REMOVED***base: base, layer: layer***REMOVED***
***REMOVED***

// Returns true if the file is not in the overlay
func (u *CopyOnWriteFs) isBaseFile(name string) (bool, error) ***REMOVED***
	if _, err := u.layer.Stat(name); err == nil ***REMOVED***
		return false, nil
	***REMOVED***
	_, err := u.base.Stat(name)
	if err != nil ***REMOVED***
		if oerr, ok := err.(*os.PathError); ok ***REMOVED***
			if oerr.Err == os.ErrNotExist || oerr.Err == syscall.ENOENT || oerr.Err == syscall.ENOTDIR ***REMOVED***
				return false, nil
			***REMOVED***
		***REMOVED***
		if err == syscall.ENOENT ***REMOVED***
			return false, nil
		***REMOVED***
	***REMOVED***
	return true, err
***REMOVED***

func (u *CopyOnWriteFs) copyToLayer(name string) error ***REMOVED***
	return copyToLayer(u.base, u.layer, name)
***REMOVED***

func (u *CopyOnWriteFs) Chtimes(name string, atime, mtime time.Time) error ***REMOVED***
	b, err := u.isBaseFile(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if b ***REMOVED***
		if err := u.copyToLayer(name); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return u.layer.Chtimes(name, atime, mtime)
***REMOVED***

func (u *CopyOnWriteFs) Chmod(name string, mode os.FileMode) error ***REMOVED***
	b, err := u.isBaseFile(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if b ***REMOVED***
		if err := u.copyToLayer(name); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return u.layer.Chmod(name, mode)
***REMOVED***

func (u *CopyOnWriteFs) Stat(name string) (os.FileInfo, error) ***REMOVED***
	fi, err := u.layer.Stat(name)
	if err != nil ***REMOVED***
		isNotExist := u.isNotExist(err)
		if isNotExist ***REMOVED***
			return u.base.Stat(name)
		***REMOVED***
		return nil, err
	***REMOVED***
	return fi, nil
***REMOVED***

func (u *CopyOnWriteFs) LstatIfPossible(name string) (os.FileInfo, bool, error) ***REMOVED***
	llayer, ok1 := u.layer.(Lstater)
	lbase, ok2 := u.base.(Lstater)

	if ok1 ***REMOVED***
		fi, b, err := llayer.LstatIfPossible(name)
		if err == nil ***REMOVED***
			return fi, b, nil
		***REMOVED***

		if !u.isNotExist(err) ***REMOVED***
			return nil, b, err
		***REMOVED***
	***REMOVED***

	if ok2 ***REMOVED***
		fi, b, err := lbase.LstatIfPossible(name)
		if err == nil ***REMOVED***
			return fi, b, nil
		***REMOVED***
		if !u.isNotExist(err) ***REMOVED***
			return nil, b, err
		***REMOVED***
	***REMOVED***

	fi, err := u.Stat(name)

	return fi, false, err
***REMOVED***

func (u *CopyOnWriteFs) isNotExist(err error) bool ***REMOVED***
	if e, ok := err.(*os.PathError); ok ***REMOVED***
		err = e.Err
	***REMOVED***
	if err == os.ErrNotExist || err == syscall.ENOENT || err == syscall.ENOTDIR ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

// Renaming files present only in the base layer is not permitted
func (u *CopyOnWriteFs) Rename(oldname, newname string) error ***REMOVED***
	b, err := u.isBaseFile(oldname)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if b ***REMOVED***
		return syscall.EPERM
	***REMOVED***
	return u.layer.Rename(oldname, newname)
***REMOVED***

// Removing files present only in the base layer is not permitted. If
// a file is present in the base layer and the overlay, only the overlay
// will be removed.
func (u *CopyOnWriteFs) Remove(name string) error ***REMOVED***
	err := u.layer.Remove(name)
	switch err ***REMOVED***
	case syscall.ENOENT:
		_, err = u.base.Stat(name)
		if err == nil ***REMOVED***
			return syscall.EPERM
		***REMOVED***
		return syscall.ENOENT
	default:
		return err
	***REMOVED***
***REMOVED***

func (u *CopyOnWriteFs) RemoveAll(name string) error ***REMOVED***
	err := u.layer.RemoveAll(name)
	switch err ***REMOVED***
	case syscall.ENOENT:
		_, err = u.base.Stat(name)
		if err == nil ***REMOVED***
			return syscall.EPERM
		***REMOVED***
		return syscall.ENOENT
	default:
		return err
	***REMOVED***
***REMOVED***

func (u *CopyOnWriteFs) OpenFile(name string, flag int, perm os.FileMode) (File, error) ***REMOVED***
	b, err := u.isBaseFile(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if flag&(os.O_WRONLY|os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_TRUNC) != 0 ***REMOVED***
		if b ***REMOVED***
			if err = u.copyToLayer(name); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return u.layer.OpenFile(name, flag, perm)
		***REMOVED***

		dir := filepath.Dir(name)
		isaDir, err := IsDir(u.base, dir)
		if err != nil && !os.IsNotExist(err) ***REMOVED***
			return nil, err
		***REMOVED***
		if isaDir ***REMOVED***
			if err = u.layer.MkdirAll(dir, 0777); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return u.layer.OpenFile(name, flag, perm)
		***REMOVED***

		isaDir, err = IsDir(u.layer, dir)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if isaDir ***REMOVED***
			return u.layer.OpenFile(name, flag, perm)
		***REMOVED***

		return nil, &os.PathError***REMOVED***Op: "open", Path: name, Err: syscall.ENOTDIR***REMOVED*** // ...or os.ErrNotExist?
	***REMOVED***
	if b ***REMOVED***
		return u.base.OpenFile(name, flag, perm)
	***REMOVED***
	return u.layer.OpenFile(name, flag, perm)
***REMOVED***

// This function handles the 9 different possibilities caused
// by the union which are the intersection of the following...
//  layer: doesn't exist, exists as a file, and exists as a directory
//  base:  doesn't exist, exists as a file, and exists as a directory
func (u *CopyOnWriteFs) Open(name string) (File, error) ***REMOVED***
	// Since the overlay overrides the base we check that first
	b, err := u.isBaseFile(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// If overlay doesn't exist, return the base (base state irrelevant)
	if b ***REMOVED***
		return u.base.Open(name)
	***REMOVED***

	// If overlay is a file, return it (base state irrelevant)
	dir, err := IsDir(u.layer, name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !dir ***REMOVED***
		return u.layer.Open(name)
	***REMOVED***

	// Overlay is a directory, base state now matters.
	// Base state has 3 states to check but 2 outcomes:
	// A. It's a file or non-readable in the base (return just the overlay)
	// B. It's an accessible directory in the base (return a UnionFile)

	// If base is file or nonreadable, return overlay
	dir, err = IsDir(u.base, name)
	if !dir || err != nil ***REMOVED***
		return u.layer.Open(name)
	***REMOVED***

	// Both base & layer are directories
	// Return union file (if opens are without error)
	bfile, bErr := u.base.Open(name)
	lfile, lErr := u.layer.Open(name)

	// If either have errors at this point something is very wrong. Return nil and the errors
	if bErr != nil || lErr != nil ***REMOVED***
		return nil, fmt.Errorf("BaseErr: %v\nOverlayErr: %v", bErr, lErr)
	***REMOVED***

	return &UnionFile***REMOVED***Base: bfile, Layer: lfile***REMOVED***, nil
***REMOVED***

func (u *CopyOnWriteFs) Mkdir(name string, perm os.FileMode) error ***REMOVED***
	dir, err := IsDir(u.base, name)
	if err != nil ***REMOVED***
		return u.layer.MkdirAll(name, perm)
	***REMOVED***
	if dir ***REMOVED***
		return syscall.EEXIST
	***REMOVED***
	return u.layer.MkdirAll(name, perm)
***REMOVED***

func (u *CopyOnWriteFs) Name() string ***REMOVED***
	return "CopyOnWriteFs"
***REMOVED***

func (u *CopyOnWriteFs) MkdirAll(name string, perm os.FileMode) error ***REMOVED***
	dir, err := IsDir(u.base, name)
	if err != nil ***REMOVED***
		return u.layer.MkdirAll(name, perm)
	***REMOVED***
	if dir ***REMOVED***
		return syscall.EEXIST
	***REMOVED***
	return u.layer.MkdirAll(name, perm)
***REMOVED***

func (u *CopyOnWriteFs) Create(name string) (File, error) ***REMOVED***
	return u.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
***REMOVED***
