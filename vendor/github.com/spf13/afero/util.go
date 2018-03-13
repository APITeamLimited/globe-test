// Copyright ©2015 Steve Francia <spf@spf13.com>
// Portions Copyright ©2015 The Hugo Authors
// Portions Copyright 2016-present Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>
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
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Filepath separator defined by os.Separator.
const FilePathSeparator = string(filepath.Separator)

// Takes a reader and a path and writes the content
func (a Afero) WriteReader(path string, r io.Reader) (err error) ***REMOVED***
	return WriteReader(a.Fs, path, r)
***REMOVED***

func WriteReader(fs Fs, path string, r io.Reader) (err error) ***REMOVED***
	dir, _ := filepath.Split(path)
	ospath := filepath.FromSlash(dir)

	if ospath != "" ***REMOVED***
		err = fs.MkdirAll(ospath, 0777) // rwx, rw, r
		if err != nil ***REMOVED***
			if err != os.ErrExist ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	file, err := fs.Create(path)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	defer file.Close()

	_, err = io.Copy(file, r)
	return
***REMOVED***

// Same as WriteReader but checks to see if file/directory already exists.
func (a Afero) SafeWriteReader(path string, r io.Reader) (err error) ***REMOVED***
	return SafeWriteReader(a.Fs, path, r)
***REMOVED***

func SafeWriteReader(fs Fs, path string, r io.Reader) (err error) ***REMOVED***
	dir, _ := filepath.Split(path)
	ospath := filepath.FromSlash(dir)

	if ospath != "" ***REMOVED***
		err = fs.MkdirAll(ospath, 0777) // rwx, rw, r
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	exists, err := Exists(fs, path)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if exists ***REMOVED***
		return fmt.Errorf("%v already exists", path)
	***REMOVED***

	file, err := fs.Create(path)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	defer file.Close()

	_, err = io.Copy(file, r)
	return
***REMOVED***

func (a Afero) GetTempDir(subPath string) string ***REMOVED***
	return GetTempDir(a.Fs, subPath)
***REMOVED***

// GetTempDir returns the default temp directory with trailing slash
// if subPath is not empty then it will be created recursively with mode 777 rwx rwx rwx
func GetTempDir(fs Fs, subPath string) string ***REMOVED***
	addSlash := func(p string) string ***REMOVED***
		if FilePathSeparator != p[len(p)-1:] ***REMOVED***
			p = p + FilePathSeparator
		***REMOVED***
		return p
	***REMOVED***
	dir := addSlash(os.TempDir())

	if subPath != "" ***REMOVED***
		// preserve windows backslash :-(
		if FilePathSeparator == "\\" ***REMOVED***
			subPath = strings.Replace(subPath, "\\", "____", -1)
		***REMOVED***
		dir = dir + UnicodeSanitize((subPath))
		if FilePathSeparator == "\\" ***REMOVED***
			dir = strings.Replace(dir, "____", "\\", -1)
		***REMOVED***

		if exists, _ := Exists(fs, dir); exists ***REMOVED***
			return addSlash(dir)
		***REMOVED***

		err := fs.MkdirAll(dir, 0777)
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		dir = addSlash(dir)
	***REMOVED***
	return dir
***REMOVED***

// Rewrite string to remove non-standard path characters
func UnicodeSanitize(s string) string ***REMOVED***
	source := []rune(s)
	target := make([]rune, 0, len(source))

	for _, r := range source ***REMOVED***
		if unicode.IsLetter(r) ||
			unicode.IsDigit(r) ||
			unicode.IsMark(r) ||
			r == '.' ||
			r == '/' ||
			r == '\\' ||
			r == '_' ||
			r == '-' ||
			r == '%' ||
			r == ' ' ||
			r == '#' ***REMOVED***
			target = append(target, r)
		***REMOVED***
	***REMOVED***

	return string(target)
***REMOVED***

// Transform characters with accents into plain forms.
func NeuterAccents(s string) string ***REMOVED***
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, string(s))

	return result
***REMOVED***

func isMn(r rune) bool ***REMOVED***
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
***REMOVED***

func (a Afero) FileContainsBytes(filename string, subslice []byte) (bool, error) ***REMOVED***
	return FileContainsBytes(a.Fs, filename, subslice)
***REMOVED***

// Check if a file contains a specified byte slice.
func FileContainsBytes(fs Fs, filename string, subslice []byte) (bool, error) ***REMOVED***
	f, err := fs.Open(filename)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	defer f.Close()

	return readerContainsAny(f, subslice), nil
***REMOVED***

func (a Afero) FileContainsAnyBytes(filename string, subslices [][]byte) (bool, error) ***REMOVED***
	return FileContainsAnyBytes(a.Fs, filename, subslices)
***REMOVED***

// Check if a file contains any of the specified byte slices.
func FileContainsAnyBytes(fs Fs, filename string, subslices [][]byte) (bool, error) ***REMOVED***
	f, err := fs.Open(filename)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	defer f.Close()

	return readerContainsAny(f, subslices...), nil
***REMOVED***

// readerContains reports whether any of the subslices is within r.
func readerContainsAny(r io.Reader, subslices ...[]byte) bool ***REMOVED***

	if r == nil || len(subslices) == 0 ***REMOVED***
		return false
	***REMOVED***

	largestSlice := 0

	for _, sl := range subslices ***REMOVED***
		if len(sl) > largestSlice ***REMOVED***
			largestSlice = len(sl)
		***REMOVED***
	***REMOVED***

	if largestSlice == 0 ***REMOVED***
		return false
	***REMOVED***

	bufflen := largestSlice * 4
	halflen := bufflen / 2
	buff := make([]byte, bufflen)
	var err error
	var n, i int

	for ***REMOVED***
		i++
		if i == 1 ***REMOVED***
			n, err = io.ReadAtLeast(r, buff[:halflen], halflen)
		***REMOVED*** else ***REMOVED***
			if i != 2 ***REMOVED***
				// shift left to catch overlapping matches
				copy(buff[:], buff[halflen:])
			***REMOVED***
			n, err = io.ReadAtLeast(r, buff[halflen:], halflen)
		***REMOVED***

		if n > 0 ***REMOVED***
			for _, sl := range subslices ***REMOVED***
				if bytes.Contains(buff, sl) ***REMOVED***
					return true
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if err != nil ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (a Afero) DirExists(path string) (bool, error) ***REMOVED***
	return DirExists(a.Fs, path)
***REMOVED***

// DirExists checks if a path exists and is a directory.
func DirExists(fs Fs, path string) (bool, error) ***REMOVED***
	fi, err := fs.Stat(path)
	if err == nil && fi.IsDir() ***REMOVED***
		return true, nil
	***REMOVED***
	if os.IsNotExist(err) ***REMOVED***
		return false, nil
	***REMOVED***
	return false, err
***REMOVED***

func (a Afero) IsDir(path string) (bool, error) ***REMOVED***
	return IsDir(a.Fs, path)
***REMOVED***

// IsDir checks if a given path is a directory.
func IsDir(fs Fs, path string) (bool, error) ***REMOVED***
	fi, err := fs.Stat(path)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	return fi.IsDir(), nil
***REMOVED***

func (a Afero) IsEmpty(path string) (bool, error) ***REMOVED***
	return IsEmpty(a.Fs, path)
***REMOVED***

// IsEmpty checks if a given file or directory is empty.
func IsEmpty(fs Fs, path string) (bool, error) ***REMOVED***
	if b, _ := Exists(fs, path); !b ***REMOVED***
		return false, fmt.Errorf("%q path does not exist", path)
	***REMOVED***
	fi, err := fs.Stat(path)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	if fi.IsDir() ***REMOVED***
		f, err := fs.Open(path)
		if err != nil ***REMOVED***
			return false, err
		***REMOVED***
		defer f.Close()
		list, err := f.Readdir(-1)
		return len(list) == 0, nil
	***REMOVED***
	return fi.Size() == 0, nil
***REMOVED***

func (a Afero) Exists(path string) (bool, error) ***REMOVED***
	return Exists(a.Fs, path)
***REMOVED***

// Check if a file or directory exists.
func Exists(fs Fs, path string) (bool, error) ***REMOVED***
	_, err := fs.Stat(path)
	if err == nil ***REMOVED***
		return true, nil
	***REMOVED***
	if os.IsNotExist(err) ***REMOVED***
		return false, nil
	***REMOVED***
	return false, err
***REMOVED***

func FullBaseFsPath(basePathFs *BasePathFs, relativePath string) string ***REMOVED***
	combinedPath := filepath.Join(basePathFs.path, relativePath)
	if parent, ok := basePathFs.source.(*BasePathFs); ok ***REMOVED***
		return FullBaseFsPath(parent, combinedPath)
	***REMOVED***

	return combinedPath
***REMOVED***
