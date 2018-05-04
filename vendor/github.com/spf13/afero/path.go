// Copyright ©2015 The Go Authors
// Copyright ©2015 Steve Francia <spf@spf13.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package afero

import (
	"os"
	"path/filepath"
	"sort"
)

// readDirNames reads the directory named by dirname and returns
// a sorted list of directory entries.
// adapted from https://golang.org/src/path/filepath/path.go
func readDirNames(fs Fs, dirname string) ([]string, error) ***REMOVED***
	f, err := fs.Open(dirname)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	sort.Strings(names)
	return names, nil
***REMOVED***

// walk recursively descends path, calling walkFn
// adapted from https://golang.org/src/path/filepath/path.go
func walk(fs Fs, path string, info os.FileInfo, walkFn filepath.WalkFunc) error ***REMOVED***
	err := walkFn(path, info, nil)
	if err != nil ***REMOVED***
		if info.IsDir() && err == filepath.SkipDir ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***

	if !info.IsDir() ***REMOVED***
		return nil
	***REMOVED***

	names, err := readDirNames(fs, path)
	if err != nil ***REMOVED***
		return walkFn(path, info, err)
	***REMOVED***

	for _, name := range names ***REMOVED***
		filename := filepath.Join(path, name)
		fileInfo, err := lstatIfPossible(fs, filename)
		if err != nil ***REMOVED***
			if err := walkFn(filename, fileInfo, err); err != nil && err != filepath.SkipDir ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			err = walk(fs, filename, fileInfo, walkFn)
			if err != nil ***REMOVED***
				if !fileInfo.IsDir() || err != filepath.SkipDir ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// if the filesystem supports it, use Lstat, else use fs.Stat
func lstatIfPossible(fs Fs, path string) (os.FileInfo, error) ***REMOVED***
	if lfs, ok := fs.(Lstater); ok ***REMOVED***
		fi, _, err := lfs.LstatIfPossible(path)
		return fi, err
	***REMOVED***
	return fs.Stat(path)
***REMOVED***

// Walk walks the file tree rooted at root, calling walkFn for each file or
// directory in the tree, including root. All errors that arise visiting files
// and directories are filtered by walkFn. The files are walked in lexical
// order, which makes the output deterministic but means that for very
// large directories Walk can be inefficient.
// Walk does not follow symbolic links.

func (a Afero) Walk(root string, walkFn filepath.WalkFunc) error ***REMOVED***
	return Walk(a.Fs, root, walkFn)
***REMOVED***

func Walk(fs Fs, root string, walkFn filepath.WalkFunc) error ***REMOVED***
	info, err := lstatIfPossible(fs, root)
	if err != nil ***REMOVED***
		return walkFn(root, nil, err)
	***REMOVED***
	return walk(fs, root, info, walkFn)
***REMOVED***
