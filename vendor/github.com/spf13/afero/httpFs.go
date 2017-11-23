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
	"errors"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type httpDir struct ***REMOVED***
	basePath string
	fs       HttpFs
***REMOVED***

func (d httpDir) Open(name string) (http.File, error) ***REMOVED***
	if filepath.Separator != '/' && strings.IndexRune(name, filepath.Separator) >= 0 ||
		strings.Contains(name, "\x00") ***REMOVED***
		return nil, errors.New("http: invalid character in file path")
	***REMOVED***
	dir := string(d.basePath)
	if dir == "" ***REMOVED***
		dir = "."
	***REMOVED***

	f, err := d.fs.Open(filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name))))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return f, nil
***REMOVED***

type HttpFs struct ***REMOVED***
	source Fs
***REMOVED***

func NewHttpFs(source Fs) *HttpFs ***REMOVED***
	return &HttpFs***REMOVED***source: source***REMOVED***
***REMOVED***

func (h HttpFs) Dir(s string) *httpDir ***REMOVED***
	return &httpDir***REMOVED***basePath: s, fs: h***REMOVED***
***REMOVED***

func (h HttpFs) Name() string ***REMOVED*** return "h HttpFs" ***REMOVED***

func (h HttpFs) Create(name string) (File, error) ***REMOVED***
	return h.source.Create(name)
***REMOVED***

func (h HttpFs) Chmod(name string, mode os.FileMode) error ***REMOVED***
	return h.source.Chmod(name, mode)
***REMOVED***

func (h HttpFs) Chtimes(name string, atime time.Time, mtime time.Time) error ***REMOVED***
	return h.source.Chtimes(name, atime, mtime)
***REMOVED***

func (h HttpFs) Mkdir(name string, perm os.FileMode) error ***REMOVED***
	return h.source.Mkdir(name, perm)
***REMOVED***

func (h HttpFs) MkdirAll(path string, perm os.FileMode) error ***REMOVED***
	return h.source.MkdirAll(path, perm)
***REMOVED***

func (h HttpFs) Open(name string) (http.File, error) ***REMOVED***
	f, err := h.source.Open(name)
	if err == nil ***REMOVED***
		if httpfile, ok := f.(http.File); ok ***REMOVED***
			return httpfile, nil
		***REMOVED***
	***REMOVED***
	return nil, err
***REMOVED***

func (h HttpFs) OpenFile(name string, flag int, perm os.FileMode) (File, error) ***REMOVED***
	return h.source.OpenFile(name, flag, perm)
***REMOVED***

func (h HttpFs) Remove(name string) error ***REMOVED***
	return h.source.Remove(name)
***REMOVED***

func (h HttpFs) RemoveAll(path string) error ***REMOVED***
	return h.source.RemoveAll(path)
***REMOVED***

func (h HttpFs) Rename(oldname, newname string) error ***REMOVED***
	return h.source.Rename(oldname, newname)
***REMOVED***

func (h HttpFs) Stat(name string) (os.FileInfo, error) ***REMOVED***
	return h.source.Stat(name)
***REMOVED***
