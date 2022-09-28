// Copyright © 2014 Steve Francia <spf@spf13.com>.
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

package afero

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/afero/mem"
)

type MemMapFs struct ***REMOVED***
	mu   sync.RWMutex
	data map[string]*mem.FileData
	init sync.Once
***REMOVED***

func NewMemMapFs() Fs ***REMOVED***
	return &MemMapFs***REMOVED******REMOVED***
***REMOVED***

func (m *MemMapFs) getData() map[string]*mem.FileData ***REMOVED***
	m.init.Do(func() ***REMOVED***
		m.data = make(map[string]*mem.FileData)
		// Root should always exist, right?
		// TODO: what about windows?
		m.data[FilePathSeparator] = mem.CreateDir(FilePathSeparator)
	***REMOVED***)
	return m.data
***REMOVED***

func (*MemMapFs) Name() string ***REMOVED*** return "MemMapFS" ***REMOVED***

func (m *MemMapFs) Create(name string) (File, error) ***REMOVED***
	name = normalizePath(name)
	m.mu.Lock()
	file := mem.CreateFile(name)
	m.getData()[name] = file
	m.registerWithParent(file)
	m.mu.Unlock()
	return mem.NewFileHandle(file), nil
***REMOVED***

func (m *MemMapFs) unRegisterWithParent(fileName string) error ***REMOVED***
	f, err := m.lockfreeOpen(fileName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	parent := m.findParent(f)
	if parent == nil ***REMOVED***
		log.Panic("parent of ", f.Name(), " is nil")
	***REMOVED***

	parent.Lock()
	mem.RemoveFromMemDir(parent, f)
	parent.Unlock()
	return nil
***REMOVED***

func (m *MemMapFs) findParent(f *mem.FileData) *mem.FileData ***REMOVED***
	pdir, _ := filepath.Split(f.Name())
	pdir = filepath.Clean(pdir)
	pfile, err := m.lockfreeOpen(pdir)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	return pfile
***REMOVED***

func (m *MemMapFs) registerWithParent(f *mem.FileData) ***REMOVED***
	if f == nil ***REMOVED***
		return
	***REMOVED***
	parent := m.findParent(f)
	if parent == nil ***REMOVED***
		pdir := filepath.Dir(filepath.Clean(f.Name()))
		err := m.lockfreeMkdir(pdir, 0777)
		if err != nil ***REMOVED***
			//log.Println("Mkdir error:", err)
			return
		***REMOVED***
		parent, err = m.lockfreeOpen(pdir)
		if err != nil ***REMOVED***
			//log.Println("Open after Mkdir error:", err)
			return
		***REMOVED***
	***REMOVED***

	parent.Lock()
	mem.InitializeDir(parent)
	mem.AddToMemDir(parent, f)
	parent.Unlock()
***REMOVED***

func (m *MemMapFs) lockfreeMkdir(name string, perm os.FileMode) error ***REMOVED***
	name = normalizePath(name)
	x, ok := m.getData()[name]
	if ok ***REMOVED***
		// Only return ErrFileExists if it's a file, not a directory.
		i := mem.FileInfo***REMOVED***FileData: x***REMOVED***
		if !i.IsDir() ***REMOVED***
			return ErrFileExists
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		item := mem.CreateDir(name)
		m.getData()[name] = item
		m.registerWithParent(item)
	***REMOVED***
	return nil
***REMOVED***

func (m *MemMapFs) Mkdir(name string, perm os.FileMode) error ***REMOVED***
	name = normalizePath(name)

	m.mu.RLock()
	_, ok := m.getData()[name]
	m.mu.RUnlock()
	if ok ***REMOVED***
		return &os.PathError***REMOVED***Op: "mkdir", Path: name, Err: ErrFileExists***REMOVED***
	***REMOVED***

	m.mu.Lock()
	item := mem.CreateDir(name)
	m.getData()[name] = item
	m.registerWithParent(item)
	m.mu.Unlock()

	m.Chmod(name, perm|os.ModeDir)

	return nil
***REMOVED***

func (m *MemMapFs) MkdirAll(path string, perm os.FileMode) error ***REMOVED***
	err := m.Mkdir(path, perm)
	if err != nil ***REMOVED***
		if err.(*os.PathError).Err == ErrFileExists ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// Handle some relative paths
func normalizePath(path string) string ***REMOVED***
	path = filepath.Clean(path)

	switch path ***REMOVED***
	case ".":
		return FilePathSeparator
	case "..":
		return FilePathSeparator
	default:
		return path
	***REMOVED***
***REMOVED***

func (m *MemMapFs) Open(name string) (File, error) ***REMOVED***
	f, err := m.open(name)
	if f != nil ***REMOVED***
		return mem.NewReadOnlyFileHandle(f), err
	***REMOVED***
	return nil, err
***REMOVED***

func (m *MemMapFs) openWrite(name string) (File, error) ***REMOVED***
	f, err := m.open(name)
	if f != nil ***REMOVED***
		return mem.NewFileHandle(f), err
	***REMOVED***
	return nil, err
***REMOVED***

func (m *MemMapFs) open(name string) (*mem.FileData, error) ***REMOVED***
	name = normalizePath(name)

	m.mu.RLock()
	f, ok := m.getData()[name]
	m.mu.RUnlock()
	if !ok ***REMOVED***
		return nil, &os.PathError***REMOVED***Op: "open", Path: name, Err: ErrFileNotFound***REMOVED***
	***REMOVED***
	return f, nil
***REMOVED***

func (m *MemMapFs) lockfreeOpen(name string) (*mem.FileData, error) ***REMOVED***
	name = normalizePath(name)
	f, ok := m.getData()[name]
	if ok ***REMOVED***
		return f, nil
	***REMOVED*** else ***REMOVED***
		return nil, ErrFileNotFound
	***REMOVED***
***REMOVED***

func (m *MemMapFs) OpenFile(name string, flag int, perm os.FileMode) (File, error) ***REMOVED***
	chmod := false
	file, err := m.openWrite(name)
	if os.IsNotExist(err) && (flag&os.O_CREATE > 0) ***REMOVED***
		file, err = m.Create(name)
		chmod = true
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if flag == os.O_RDONLY ***REMOVED***
		file = mem.NewReadOnlyFileHandle(file.(*mem.File).Data())
	***REMOVED***
	if flag&os.O_APPEND > 0 ***REMOVED***
		_, err = file.Seek(0, os.SEEK_END)
		if err != nil ***REMOVED***
			file.Close()
			return nil, err
		***REMOVED***
	***REMOVED***
	if flag&os.O_TRUNC > 0 && flag&(os.O_RDWR|os.O_WRONLY) > 0 ***REMOVED***
		err = file.Truncate(0)
		if err != nil ***REMOVED***
			file.Close()
			return nil, err
		***REMOVED***
	***REMOVED***
	if chmod ***REMOVED***
		m.Chmod(name, perm)
	***REMOVED***
	return file, nil
***REMOVED***

func (m *MemMapFs) Remove(name string) error ***REMOVED***
	name = normalizePath(name)

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.getData()[name]; ok ***REMOVED***
		err := m.unRegisterWithParent(name)
		if err != nil ***REMOVED***
			return &os.PathError***REMOVED***Op: "remove", Path: name, Err: err***REMOVED***
		***REMOVED***
		delete(m.getData(), name)
	***REMOVED*** else ***REMOVED***
		return &os.PathError***REMOVED***Op: "remove", Path: name, Err: os.ErrNotExist***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (m *MemMapFs) RemoveAll(path string) error ***REMOVED***
	path = normalizePath(path)
	m.mu.Lock()
	m.unRegisterWithParent(path)
	m.mu.Unlock()

	m.mu.RLock()
	defer m.mu.RUnlock()

	for p, _ := range m.getData() ***REMOVED***
		if strings.HasPrefix(p, path) ***REMOVED***
			m.mu.RUnlock()
			m.mu.Lock()
			delete(m.getData(), p)
			m.mu.Unlock()
			m.mu.RLock()
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (m *MemMapFs) Rename(oldname, newname string) error ***REMOVED***
	oldname = normalizePath(oldname)
	newname = normalizePath(newname)

	if oldname == newname ***REMOVED***
		return nil
	***REMOVED***

	m.mu.RLock()
	defer m.mu.RUnlock()
	if _, ok := m.getData()[oldname]; ok ***REMOVED***
		m.mu.RUnlock()
		m.mu.Lock()
		m.unRegisterWithParent(oldname)
		fileData := m.getData()[oldname]
		delete(m.getData(), oldname)
		mem.ChangeFileName(fileData, newname)
		m.getData()[newname] = fileData
		m.registerWithParent(fileData)
		m.mu.Unlock()
		m.mu.RLock()
	***REMOVED*** else ***REMOVED***
		return &os.PathError***REMOVED***Op: "rename", Path: oldname, Err: ErrFileNotFound***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (m *MemMapFs) Stat(name string) (os.FileInfo, error) ***REMOVED***
	f, err := m.Open(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	fi := mem.GetFileInfo(f.(*mem.File).Data())
	return fi, nil
***REMOVED***

func (m *MemMapFs) Chmod(name string, mode os.FileMode) error ***REMOVED***
	name = normalizePath(name)

	m.mu.RLock()
	f, ok := m.getData()[name]
	m.mu.RUnlock()
	if !ok ***REMOVED***
		return &os.PathError***REMOVED***Op: "chmod", Path: name, Err: ErrFileNotFound***REMOVED***
	***REMOVED***

	m.mu.Lock()
	mem.SetMode(f, mode)
	m.mu.Unlock()

	return nil
***REMOVED***

func (m *MemMapFs) Chtimes(name string, atime time.Time, mtime time.Time) error ***REMOVED***
	name = normalizePath(name)

	m.mu.RLock()
	f, ok := m.getData()[name]
	m.mu.RUnlock()
	if !ok ***REMOVED***
		return &os.PathError***REMOVED***Op: "chtimes", Path: name, Err: ErrFileNotFound***REMOVED***
	***REMOVED***

	m.mu.Lock()
	mem.SetModTime(f, mtime)
	m.mu.Unlock()

	return nil
***REMOVED***

func (m *MemMapFs) List() ***REMOVED***
	for _, x := range m.data ***REMOVED***
		y := mem.FileInfo***REMOVED***FileData: x***REMOVED***
		fmt.Println(x.Name(), y.Size())
	***REMOVED***
***REMOVED***

// func debugMemMapList(fs Fs) ***REMOVED***
// 	if x, ok := fs.(*MemMapFs); ok ***REMOVED***
// 		x.List()
// 	***REMOVED***
// ***REMOVED***
