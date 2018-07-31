// Copyright © 2015 Steve Francia <spf@spf13.com>.
// Copyright 2013 tsuru authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mem

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

import "time"

const FilePathSeparator = string(filepath.Separator)

type File struct ***REMOVED***
	// atomic requires 64-bit alignment for struct field access
	at           int64
	readDirCount int64
	closed       bool
	readOnly     bool
	fileData     *FileData
***REMOVED***

func NewFileHandle(data *FileData) *File ***REMOVED***
	return &File***REMOVED***fileData: data***REMOVED***
***REMOVED***

func NewReadOnlyFileHandle(data *FileData) *File ***REMOVED***
	return &File***REMOVED***fileData: data, readOnly: true***REMOVED***
***REMOVED***

func (f File) Data() *FileData ***REMOVED***
	return f.fileData
***REMOVED***

type FileData struct ***REMOVED***
	sync.Mutex
	name    string
	data    []byte
	memDir  Dir
	dir     bool
	mode    os.FileMode
	modtime time.Time
***REMOVED***

func (d *FileData) Name() string ***REMOVED***
	d.Lock()
	defer d.Unlock()
	return d.name
***REMOVED***

func CreateFile(name string) *FileData ***REMOVED***
	return &FileData***REMOVED***name: name, mode: os.ModeTemporary, modtime: time.Now()***REMOVED***
***REMOVED***

func CreateDir(name string) *FileData ***REMOVED***
	return &FileData***REMOVED***name: name, memDir: &DirMap***REMOVED******REMOVED***, dir: true***REMOVED***
***REMOVED***

func ChangeFileName(f *FileData, newname string) ***REMOVED***
	f.Lock()
	f.name = newname
	f.Unlock()
***REMOVED***

func SetMode(f *FileData, mode os.FileMode) ***REMOVED***
	f.Lock()
	f.mode = mode
	f.Unlock()
***REMOVED***

func SetModTime(f *FileData, mtime time.Time) ***REMOVED***
	f.Lock()
	setModTime(f, mtime)
	f.Unlock()
***REMOVED***

func setModTime(f *FileData, mtime time.Time) ***REMOVED***
	f.modtime = mtime
***REMOVED***

func GetFileInfo(f *FileData) *FileInfo ***REMOVED***
	return &FileInfo***REMOVED***f***REMOVED***
***REMOVED***

func (f *File) Open() error ***REMOVED***
	atomic.StoreInt64(&f.at, 0)
	atomic.StoreInt64(&f.readDirCount, 0)
	f.fileData.Lock()
	f.closed = false
	f.fileData.Unlock()
	return nil
***REMOVED***

func (f *File) Close() error ***REMOVED***
	f.fileData.Lock()
	f.closed = true
	if !f.readOnly ***REMOVED***
		setModTime(f.fileData, time.Now())
	***REMOVED***
	f.fileData.Unlock()
	return nil
***REMOVED***

func (f *File) Name() string ***REMOVED***
	return f.fileData.Name()
***REMOVED***

func (f *File) Stat() (os.FileInfo, error) ***REMOVED***
	return &FileInfo***REMOVED***f.fileData***REMOVED***, nil
***REMOVED***

func (f *File) Sync() error ***REMOVED***
	return nil
***REMOVED***

func (f *File) Readdir(count int) (res []os.FileInfo, err error) ***REMOVED***
	if !f.fileData.dir ***REMOVED***
		return nil, &os.PathError***REMOVED***Op: "readdir", Path: f.fileData.name, Err: errors.New("not a dir")***REMOVED***
	***REMOVED***
	var outLength int64

	f.fileData.Lock()
	files := f.fileData.memDir.Files()[f.readDirCount:]
	if count > 0 ***REMOVED***
		if len(files) < count ***REMOVED***
			outLength = int64(len(files))
		***REMOVED*** else ***REMOVED***
			outLength = int64(count)
		***REMOVED***
		if len(files) == 0 ***REMOVED***
			err = io.EOF
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		outLength = int64(len(files))
	***REMOVED***
	f.readDirCount += outLength
	f.fileData.Unlock()

	res = make([]os.FileInfo, outLength)
	for i := range res ***REMOVED***
		res[i] = &FileInfo***REMOVED***files[i]***REMOVED***
	***REMOVED***

	return res, err
***REMOVED***

func (f *File) Readdirnames(n int) (names []string, err error) ***REMOVED***
	fi, err := f.Readdir(n)
	names = make([]string, len(fi))
	for i, f := range fi ***REMOVED***
		_, names[i] = filepath.Split(f.Name())
	***REMOVED***
	return names, err
***REMOVED***

func (f *File) Read(b []byte) (n int, err error) ***REMOVED***
	f.fileData.Lock()
	defer f.fileData.Unlock()
	if f.closed == true ***REMOVED***
		return 0, ErrFileClosed
	***REMOVED***
	if len(b) > 0 && int(f.at) == len(f.fileData.data) ***REMOVED***
		return 0, io.EOF
	***REMOVED***
	if int(f.at) > len(f.fileData.data) ***REMOVED***
		return 0, io.ErrUnexpectedEOF
	***REMOVED***
	if len(f.fileData.data)-int(f.at) >= len(b) ***REMOVED***
		n = len(b)
	***REMOVED*** else ***REMOVED***
		n = len(f.fileData.data) - int(f.at)
	***REMOVED***
	copy(b, f.fileData.data[f.at:f.at+int64(n)])
	atomic.AddInt64(&f.at, int64(n))
	return
***REMOVED***

func (f *File) ReadAt(b []byte, off int64) (n int, err error) ***REMOVED***
	atomic.StoreInt64(&f.at, off)
	return f.Read(b)
***REMOVED***

func (f *File) Truncate(size int64) error ***REMOVED***
	if f.closed == true ***REMOVED***
		return ErrFileClosed
	***REMOVED***
	if f.readOnly ***REMOVED***
		return &os.PathError***REMOVED***Op: "truncate", Path: f.fileData.name, Err: errors.New("file handle is read only")***REMOVED***
	***REMOVED***
	if size < 0 ***REMOVED***
		return ErrOutOfRange
	***REMOVED***
	if size > int64(len(f.fileData.data)) ***REMOVED***
		diff := size - int64(len(f.fileData.data))
		f.fileData.data = append(f.fileData.data, bytes.Repeat([]byte***REMOVED***00***REMOVED***, int(diff))...)
	***REMOVED*** else ***REMOVED***
		f.fileData.data = f.fileData.data[0:size]
	***REMOVED***
	setModTime(f.fileData, time.Now())
	return nil
***REMOVED***

func (f *File) Seek(offset int64, whence int) (int64, error) ***REMOVED***
	if f.closed == true ***REMOVED***
		return 0, ErrFileClosed
	***REMOVED***
	switch whence ***REMOVED***
	case 0:
		atomic.StoreInt64(&f.at, offset)
	case 1:
		atomic.AddInt64(&f.at, int64(offset))
	case 2:
		atomic.StoreInt64(&f.at, int64(len(f.fileData.data))+offset)
	***REMOVED***
	return f.at, nil
***REMOVED***

func (f *File) Write(b []byte) (n int, err error) ***REMOVED***
	if f.readOnly ***REMOVED***
		return 0, &os.PathError***REMOVED***Op: "write", Path: f.fileData.name, Err: errors.New("file handle is read only")***REMOVED***
	***REMOVED***
	n = len(b)
	cur := atomic.LoadInt64(&f.at)
	f.fileData.Lock()
	defer f.fileData.Unlock()
	diff := cur - int64(len(f.fileData.data))
	var tail []byte
	if n+int(cur) < len(f.fileData.data) ***REMOVED***
		tail = f.fileData.data[n+int(cur):]
	***REMOVED***
	if diff > 0 ***REMOVED***
		f.fileData.data = append(bytes.Repeat([]byte***REMOVED***00***REMOVED***, int(diff)), b...)
		f.fileData.data = append(f.fileData.data, tail...)
	***REMOVED*** else ***REMOVED***
		f.fileData.data = append(f.fileData.data[:cur], b...)
		f.fileData.data = append(f.fileData.data, tail...)
	***REMOVED***
	setModTime(f.fileData, time.Now())

	atomic.StoreInt64(&f.at, int64(len(f.fileData.data)))
	return
***REMOVED***

func (f *File) WriteAt(b []byte, off int64) (n int, err error) ***REMOVED***
	atomic.StoreInt64(&f.at, off)
	return f.Write(b)
***REMOVED***

func (f *File) WriteString(s string) (ret int, err error) ***REMOVED***
	return f.Write([]byte(s))
***REMOVED***

func (f *File) Info() *FileInfo ***REMOVED***
	return &FileInfo***REMOVED***f.fileData***REMOVED***
***REMOVED***

type FileInfo struct ***REMOVED***
	*FileData
***REMOVED***

// Implements os.FileInfo
func (s *FileInfo) Name() string ***REMOVED***
	s.Lock()
	_, name := filepath.Split(s.name)
	s.Unlock()
	return name
***REMOVED***
func (s *FileInfo) Mode() os.FileMode ***REMOVED***
	s.Lock()
	defer s.Unlock()
	return s.mode
***REMOVED***
func (s *FileInfo) ModTime() time.Time ***REMOVED***
	s.Lock()
	defer s.Unlock()
	return s.modtime
***REMOVED***
func (s *FileInfo) IsDir() bool ***REMOVED***
	s.Lock()
	defer s.Unlock()
	return s.dir
***REMOVED***
func (s *FileInfo) Sys() interface***REMOVED******REMOVED*** ***REMOVED*** return nil ***REMOVED***
func (s *FileInfo) Size() int64 ***REMOVED***
	if s.IsDir() ***REMOVED***
		return int64(42)
	***REMOVED***
	s.Lock()
	defer s.Unlock()
	return int64(len(s.data))
***REMOVED***

var (
	ErrFileClosed        = errors.New("File is closed")
	ErrOutOfRange        = errors.New("Out of range")
	ErrTooLarge          = errors.New("Too large")
	ErrFileNotFound      = os.ErrNotExist
	ErrFileExists        = os.ErrExist
	ErrDestinationExists = os.ErrExist
)
