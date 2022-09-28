package afero

import (
	"os"
	"syscall"
	"time"
)

var _ Lstater = (*ReadOnlyFs)(nil)

type ReadOnlyFs struct ***REMOVED***
	source Fs
***REMOVED***

func NewReadOnlyFs(source Fs) Fs ***REMOVED***
	return &ReadOnlyFs***REMOVED***source: source***REMOVED***
***REMOVED***

func (r *ReadOnlyFs) ReadDir(name string) ([]os.FileInfo, error) ***REMOVED***
	return ReadDir(r.source, name)
***REMOVED***

func (r *ReadOnlyFs) Chtimes(n string, a, m time.Time) error ***REMOVED***
	return syscall.EPERM
***REMOVED***

func (r *ReadOnlyFs) Chmod(n string, m os.FileMode) error ***REMOVED***
	return syscall.EPERM
***REMOVED***

func (r *ReadOnlyFs) Name() string ***REMOVED***
	return "ReadOnlyFilter"
***REMOVED***

func (r *ReadOnlyFs) Stat(name string) (os.FileInfo, error) ***REMOVED***
	return r.source.Stat(name)
***REMOVED***

func (r *ReadOnlyFs) LstatIfPossible(name string) (os.FileInfo, bool, error) ***REMOVED***
	if lsf, ok := r.source.(Lstater); ok ***REMOVED***
		return lsf.LstatIfPossible(name)
	***REMOVED***
	fi, err := r.Stat(name)
	return fi, false, err
***REMOVED***

func (r *ReadOnlyFs) Rename(o, n string) error ***REMOVED***
	return syscall.EPERM
***REMOVED***

func (r *ReadOnlyFs) RemoveAll(p string) error ***REMOVED***
	return syscall.EPERM
***REMOVED***

func (r *ReadOnlyFs) Remove(n string) error ***REMOVED***
	return syscall.EPERM
***REMOVED***

func (r *ReadOnlyFs) OpenFile(name string, flag int, perm os.FileMode) (File, error) ***REMOVED***
	if flag&(os.O_WRONLY|syscall.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_TRUNC) != 0 ***REMOVED***
		return nil, syscall.EPERM
	***REMOVED***
	return r.source.OpenFile(name, flag, perm)
***REMOVED***

func (r *ReadOnlyFs) Open(n string) (File, error) ***REMOVED***
	return r.source.Open(n)
***REMOVED***

func (r *ReadOnlyFs) Mkdir(n string, p os.FileMode) error ***REMOVED***
	return syscall.EPERM
***REMOVED***

func (r *ReadOnlyFs) MkdirAll(n string, p os.FileMode) error ***REMOVED***
	return syscall.EPERM
***REMOVED***

func (r *ReadOnlyFs) Create(n string) (File, error) ***REMOVED***
	return nil, syscall.EPERM
***REMOVED***
