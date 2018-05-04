// Copyright © 2014 Steve Francia <spf@spf13.com>.
// Copyright 2009 The Go Authors. All rights reserved.

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
	"path/filepath"
	"sort"
	"strings"
)

// Glob returns the names of all files matching pattern or nil
// if there is no matching file. The syntax of patterns is the same
// as in Match. The pattern may describe hierarchical names such as
// /usr/*/bin/ed (assuming the Separator is '/').
//
// Glob ignores file system errors such as I/O errors reading directories.
// The only possible returned error is ErrBadPattern, when pattern
// is malformed.
//
// This was adapted from (http://golang.org/pkg/path/filepath) and uses several
// built-ins from that package.
func Glob(fs Fs, pattern string) (matches []string, err error) ***REMOVED***
	if !hasMeta(pattern) ***REMOVED***
		// Lstat not supported by a ll filesystems.
		if _, err = lstatIfPossible(fs, pattern); err != nil ***REMOVED***
			return nil, nil
		***REMOVED***
		return []string***REMOVED***pattern***REMOVED***, nil
	***REMOVED***

	dir, file := filepath.Split(pattern)
	switch dir ***REMOVED***
	case "":
		dir = "."
	case string(filepath.Separator):
	// nothing
	default:
		dir = dir[0 : len(dir)-1] // chop off trailing separator
	***REMOVED***

	if !hasMeta(dir) ***REMOVED***
		return glob(fs, dir, file, nil)
	***REMOVED***

	var m []string
	m, err = Glob(fs, dir)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	for _, d := range m ***REMOVED***
		matches, err = glob(fs, d, file, matches)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// glob searches for files matching pattern in the directory dir
// and appends them to matches. If the directory cannot be
// opened, it returns the existing matches. New matches are
// added in lexicographical order.
func glob(fs Fs, dir, pattern string, matches []string) (m []string, e error) ***REMOVED***
	m = matches
	fi, err := fs.Stat(dir)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if !fi.IsDir() ***REMOVED***
		return
	***REMOVED***
	d, err := fs.Open(dir)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	defer d.Close()

	names, _ := d.Readdirnames(-1)
	sort.Strings(names)

	for _, n := range names ***REMOVED***
		matched, err := filepath.Match(pattern, n)
		if err != nil ***REMOVED***
			return m, err
		***REMOVED***
		if matched ***REMOVED***
			m = append(m, filepath.Join(dir, n))
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// hasMeta reports whether path contains any of the magic characters
// recognized by Match.
func hasMeta(path string) bool ***REMOVED***
	// TODO(niemeyer): Should other magic characters be added here?
	return strings.IndexAny(path, "*?[") >= 0
***REMOVED***
