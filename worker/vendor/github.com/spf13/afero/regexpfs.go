package afero

import (
	"os"
	"regexp"
	"syscall"
	"time"
)

// The RegexpFs filters files (not directories) by regular expression. Only
// files matching the given regexp will be allowed, all others get a ENOENT error (
// "No such file or directory").
//
type RegexpFs struct ***REMOVED***
	re     *regexp.Regexp
	source Fs
***REMOVED***

func NewRegexpFs(source Fs, re *regexp.Regexp) Fs ***REMOVED***
	return &RegexpFs***REMOVED***source: source, re: re***REMOVED***
***REMOVED***

type RegexpFile struct ***REMOVED***
	f  File
	re *regexp.Regexp
***REMOVED***

func (r *RegexpFs) matchesName(name string) error ***REMOVED***
	if r.re == nil ***REMOVED***
		return nil
	***REMOVED***
	if r.re.MatchString(name) ***REMOVED***
		return nil
	***REMOVED***
	return syscall.ENOENT
***REMOVED***

func (r *RegexpFs) dirOrMatches(name string) error ***REMOVED***
	dir, err := IsDir(r.source, name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if dir ***REMOVED***
		return nil
	***REMOVED***
	return r.matchesName(name)
***REMOVED***

func (r *RegexpFs) Chtimes(name string, a, m time.Time) error ***REMOVED***
	if err := r.dirOrMatches(name); err != nil ***REMOVED***
		return err
	***REMOVED***
	return r.source.Chtimes(name, a, m)
***REMOVED***

func (r *RegexpFs) Chmod(name string, mode os.FileMode) error ***REMOVED***
	if err := r.dirOrMatches(name); err != nil ***REMOVED***
		return err
	***REMOVED***
	return r.source.Chmod(name, mode)
***REMOVED***

func (r *RegexpFs) Name() string ***REMOVED***
	return "RegexpFs"
***REMOVED***

func (r *RegexpFs) Stat(name string) (os.FileInfo, error) ***REMOVED***
	if err := r.dirOrMatches(name); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return r.source.Stat(name)
***REMOVED***

func (r *RegexpFs) Rename(oldname, newname string) error ***REMOVED***
	dir, err := IsDir(r.source, oldname)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if dir ***REMOVED***
		return nil
	***REMOVED***
	if err := r.matchesName(oldname); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := r.matchesName(newname); err != nil ***REMOVED***
		return err
	***REMOVED***
	return r.source.Rename(oldname, newname)
***REMOVED***

func (r *RegexpFs) RemoveAll(p string) error ***REMOVED***
	dir, err := IsDir(r.source, p)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if !dir ***REMOVED***
		if err := r.matchesName(p); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return r.source.RemoveAll(p)
***REMOVED***

func (r *RegexpFs) Remove(name string) error ***REMOVED***
	if err := r.dirOrMatches(name); err != nil ***REMOVED***
		return err
	***REMOVED***
	return r.source.Remove(name)
***REMOVED***

func (r *RegexpFs) OpenFile(name string, flag int, perm os.FileMode) (File, error) ***REMOVED***
	if err := r.dirOrMatches(name); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return r.source.OpenFile(name, flag, perm)
***REMOVED***

func (r *RegexpFs) Open(name string) (File, error) ***REMOVED***
	dir, err := IsDir(r.source, name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !dir ***REMOVED***
		if err := r.matchesName(name); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	f, err := r.source.Open(name)
	return &RegexpFile***REMOVED***f: f, re: r.re***REMOVED***, nil
***REMOVED***

func (r *RegexpFs) Mkdir(n string, p os.FileMode) error ***REMOVED***
	return r.source.Mkdir(n, p)
***REMOVED***

func (r *RegexpFs) MkdirAll(n string, p os.FileMode) error ***REMOVED***
	return r.source.MkdirAll(n, p)
***REMOVED***

func (r *RegexpFs) Create(name string) (File, error) ***REMOVED***
	if err := r.matchesName(name); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return r.source.Create(name)
***REMOVED***

func (f *RegexpFile) Close() error ***REMOVED***
	return f.f.Close()
***REMOVED***

func (f *RegexpFile) Read(s []byte) (int, error) ***REMOVED***
	return f.f.Read(s)
***REMOVED***

func (f *RegexpFile) ReadAt(s []byte, o int64) (int, error) ***REMOVED***
	return f.f.ReadAt(s, o)
***REMOVED***

func (f *RegexpFile) Seek(o int64, w int) (int64, error) ***REMOVED***
	return f.f.Seek(o, w)
***REMOVED***

func (f *RegexpFile) Write(s []byte) (int, error) ***REMOVED***
	return f.f.Write(s)
***REMOVED***

func (f *RegexpFile) WriteAt(s []byte, o int64) (int, error) ***REMOVED***
	return f.f.WriteAt(s, o)
***REMOVED***

func (f *RegexpFile) Name() string ***REMOVED***
	return f.f.Name()
***REMOVED***

func (f *RegexpFile) Readdir(c int) (fi []os.FileInfo, err error) ***REMOVED***
	var rfi []os.FileInfo
	rfi, err = f.f.Readdir(c)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	for _, i := range rfi ***REMOVED***
		if i.IsDir() || f.re.MatchString(i.Name()) ***REMOVED***
			fi = append(fi, i)
		***REMOVED***
	***REMOVED***
	return fi, nil
***REMOVED***

func (f *RegexpFile) Readdirnames(c int) (n []string, err error) ***REMOVED***
	fi, err := f.Readdir(c)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	for _, s := range fi ***REMOVED***
		n = append(n, s.Name())
	***REMOVED***
	return n, nil
***REMOVED***

func (f *RegexpFile) Stat() (os.FileInfo, error) ***REMOVED***
	return f.f.Stat()
***REMOVED***

func (f *RegexpFile) Sync() error ***REMOVED***
	return f.f.Sync()
***REMOVED***

func (f *RegexpFile) Truncate(s int64) error ***REMOVED***
	return f.f.Truncate(s)
***REMOVED***

func (f *RegexpFile) WriteString(s string) (int, error) ***REMOVED***
	return f.f.WriteString(s)
***REMOVED***
