// Copyright © 2014 Steve Francia <spf@spf13.com>.
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

package afero

import (
	"os"
	"time"
)

var _ Lstater = (*OsFs)(nil)

// OsFs is a Fs implementation that uses functions provided by the os package.
//
// For details in any method, check the documentation of the os package
// (http://golang.org/pkg/os/).
type OsFs struct***REMOVED******REMOVED***

func NewOsFs() Fs ***REMOVED***
	return &OsFs***REMOVED******REMOVED***
***REMOVED***

func (OsFs) Name() string ***REMOVED*** return "OsFs" ***REMOVED***

func (OsFs) Create(name string) (File, error) ***REMOVED***
	f, e := os.Create(name)
	if f == nil ***REMOVED***
		// while this looks strange, we need to return a bare nil (of type nil) not
		// a nil value of type *os.File or nil won't be nil
		return nil, e
	***REMOVED***
	return f, e
***REMOVED***

func (OsFs) Mkdir(name string, perm os.FileMode) error ***REMOVED***
	return os.Mkdir(name, perm)
***REMOVED***

func (OsFs) MkdirAll(path string, perm os.FileMode) error ***REMOVED***
	return os.MkdirAll(path, perm)
***REMOVED***

func (OsFs) Open(name string) (File, error) ***REMOVED***
	f, e := os.Open(name)
	if f == nil ***REMOVED***
		// while this looks strange, we need to return a bare nil (of type nil) not
		// a nil value of type *os.File or nil won't be nil
		return nil, e
	***REMOVED***
	return f, e
***REMOVED***

func (OsFs) OpenFile(name string, flag int, perm os.FileMode) (File, error) ***REMOVED***
	f, e := os.OpenFile(name, flag, perm)
	if f == nil ***REMOVED***
		// while this looks strange, we need to return a bare nil (of type nil) not
		// a nil value of type *os.File or nil won't be nil
		return nil, e
	***REMOVED***
	return f, e
***REMOVED***

func (OsFs) Remove(name string) error ***REMOVED***
	return os.Remove(name)
***REMOVED***

func (OsFs) RemoveAll(path string) error ***REMOVED***
	return os.RemoveAll(path)
***REMOVED***

func (OsFs) Rename(oldname, newname string) error ***REMOVED***
	return os.Rename(oldname, newname)
***REMOVED***

func (OsFs) Stat(name string) (os.FileInfo, error) ***REMOVED***
	return os.Stat(name)
***REMOVED***

func (OsFs) Chmod(name string, mode os.FileMode) error ***REMOVED***
	return os.Chmod(name, mode)
***REMOVED***

func (OsFs) Chtimes(name string, atime time.Time, mtime time.Time) error ***REMOVED***
	return os.Chtimes(name, atime, mtime)
***REMOVED***

func (OsFs) LstatIfPossible(name string) (os.FileInfo, bool, error) ***REMOVED***
	fi, err := os.Lstat(name)
	return fi, true, err
***REMOVED***
