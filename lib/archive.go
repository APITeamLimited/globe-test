/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package lib

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/spf13/afero"
)

var volumeRE = regexp.MustCompile(`^([a-zA-Z]):(.*)`)
var sharedRE = regexp.MustCompile(`^\\\\([^\\]+)`) // matches a shared folder in Windows before backslack replacement. i.e \\VMBOXSVR\k6\script.js
var homeDirRE = regexp.MustCompile(`^(/[a-zA-Z])?/(Users|home|Documents and Settings)/(?:[^/]+)`)

// NormalizeAndAnonymizePath Normalizes (to use a / path separator) and anonymizes a file path, by scrubbing usernames from home directories.
func NormalizeAndAnonymizePath(path string) string ***REMOVED***
	path = filepath.Clean(path)

	p := volumeRE.ReplaceAllString(path, `/$1$2`)
	p = sharedRE.ReplaceAllString(p, `/nobody`)
	p = strings.Replace(p, "\\", "/", -1)
	return homeDirRE.ReplaceAllString(p, `$1/$2/nobody`)
***REMOVED***

type normalizedFS struct ***REMOVED***
	afero.Fs
***REMOVED***

func (m *normalizedFS) Open(name string) (afero.File, error) ***REMOVED***
	return m.Fs.Open(NormalizeAndAnonymizePath(name))
***REMOVED***

func (m *normalizedFS) OpenFile(name string, flag int, mode os.FileMode) (afero.File, error) ***REMOVED***
	return m.Fs.OpenFile(NormalizeAndAnonymizePath(name), flag, mode)
***REMOVED***

func (m *normalizedFS) Stat(name string) (os.FileInfo, error) ***REMOVED***
	return m.Fs.Stat(NormalizeAndAnonymizePath(name))
***REMOVED***

// An Archive is a rollup of all resources and options needed to reproduce a test identically elsewhere.
type Archive struct ***REMOVED***
	// The runner to use, eg. "js".
	Type string `json:"type"`

	// Options to use.
	Options Options `json:"options"`

	// Filename and contents of the main file being executed.
	Filename string `json:"filename"`
	Data     []byte `json:"-"`

	// Working directory for resolving relative paths.
	Pwd string `json:"pwd"`

	// Archived filesystem.
	Scripts map[string][]byte `json:"-"` // included scripts
	Files   map[string][]byte `json:"-"` // non-script resources

	FS afero.Fs `json:"-"`

	// Environment variables
	Env map[string]string `json:"env"`
***REMOVED***

// Reads an archive created by Archive.Write from a reader.
func ReadArchive(in io.Reader) (*Archive, error) ***REMOVED***
	r := tar.NewReader(in)
	arc := &Archive***REMOVED***
		Scripts: make(map[string][]byte),
		Files:   make(map[string][]byte),
		FS:      &normalizedFS***REMOVED***Fs: afero.NewMemMapFs()***REMOVED***,
	***REMOVED***

	for ***REMOVED***
		hdr, err := r.Next()
		if err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				break
			***REMOVED***
			return nil, err
		***REMOVED***
		if hdr.Typeflag != tar.TypeReg && hdr.Typeflag != tar.TypeRegA ***REMOVED***
			continue
		***REMOVED***

		data, err := ioutil.ReadAll(r)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		switch hdr.Name ***REMOVED***
		case "metadata.json":
			if err := json.Unmarshal(data, &arc); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			// Path separator normalization for older archives (<=0.20.0)
			arc.Filename = NormalizeAndAnonymizePath(arc.Filename)
			arc.Pwd = NormalizeAndAnonymizePath(arc.Pwd)
			continue
		case "data":
			arc.Data = data
		***REMOVED***

		// Path separator normalization for older archives (<=0.20.0)
		normPath := NormalizeAndAnonymizePath(hdr.Name)
		idx := strings.IndexRune(normPath, '/')
		if idx == -1 ***REMOVED***
			continue
		***REMOVED***
		pfx := normPath[:idx]
		name := normPath[idx+1:]
		if name != "" && name[0] == '_' ***REMOVED***
			name = name[1:]
		***REMOVED***

		var dst map[string][]byte
		switch pfx ***REMOVED***
		case "files":
			dst = arc.Files
		case "scripts":
			dst = arc.Scripts
		default:
			continue
		***REMOVED***

		dst[name] = data

		err = afero.WriteFile(arc.FS, name, data, os.ModePerm)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return arc, nil
***REMOVED***

// Write serialises the archive to a writer.
//
// The format should be treated as opaque; currently it is simply a TAR rollup, but this may
// change. If it does change, ReadArchive must be able to handle all previous formats as well as
// the current one.
func (arc *Archive) Write(out io.Writer) error ***REMOVED***
	w := tar.NewWriter(out)
	t := time.Now()

	metaArc := *arc
	metaArc.Filename = NormalizeAndAnonymizePath(metaArc.Filename)
	metaArc.Pwd = NormalizeAndAnonymizePath(metaArc.Pwd)
	metadata, err := metaArc.json()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_ = w.WriteHeader(&tar.Header***REMOVED***
		Name:     "metadata.json",
		Mode:     0644,
		Size:     int64(len(metadata)),
		ModTime:  t,
		Typeflag: tar.TypeReg,
	***REMOVED***)
	if _, err := w.Write(metadata); err != nil ***REMOVED***
		return err
	***REMOVED***

	_ = w.WriteHeader(&tar.Header***REMOVED***
		Name:     "data",
		Mode:     0644,
		Size:     int64(len(arc.Data)),
		ModTime:  t,
		Typeflag: tar.TypeReg,
	***REMOVED***)
	if _, err := w.Write(arc.Data); err != nil ***REMOVED***
		return err
	***REMOVED***

	arcfs := []struct ***REMOVED***
		name  string
		files map[string][]byte
	***REMOVED******REMOVED***
		***REMOVED***"scripts", arc.Scripts***REMOVED***,
		***REMOVED***"files", arc.Files***REMOVED***,
	***REMOVED***
	for _, entry := range arcfs ***REMOVED***
		_ = w.WriteHeader(&tar.Header***REMOVED***
			Name:     entry.name,
			Mode:     0755,
			ModTime:  t,
			Typeflag: tar.TypeDir,
		***REMOVED***)

		// A couple of things going on here:
		// - You can't just create file entries, you need to create directory entries too.
		//   Figure out which directories are in use here.
		// - We want archives to be comparable by hash, which means the entries need to be written
		//   in the same order every time. Go maps are shuffled, so we need to sort lists of keys.
		// - We don't want to leak private information (eg. usernames) in archives, so make sure to
		//   anonymize paths before stuffing them in a shareable archive.
		foundDirs := make(map[string]bool)
		paths := make([]string, 0, len(entry.files))
		files := make(map[string][]byte, len(entry.files))
		for filePath, data := range entry.files ***REMOVED***
			filePath = NormalizeAndAnonymizePath(filePath)
			files[filePath] = data
			paths = append(paths, filePath)
			dir := path.Dir(filePath)
			for ***REMOVED***
				foundDirs[dir] = true
				idx := strings.LastIndexByte(dir, os.PathSeparator)
				if idx == -1 ***REMOVED***
					break
				***REMOVED***
				dir = dir[:idx]
			***REMOVED***
		***REMOVED***
		dirs := make([]string, 0, len(foundDirs))
		for dirpath := range foundDirs ***REMOVED***
			dirs = append(dirs, dirpath)
		***REMOVED***
		sort.Strings(paths)
		sort.Strings(dirs)

		for _, dirpath := range dirs ***REMOVED***
			if dirpath == "" || dirpath[0] == '/' ***REMOVED***
				dirpath = "_" + dirpath
			***REMOVED***
			_ = w.WriteHeader(&tar.Header***REMOVED***
				Name:     path.Clean(entry.name + "/" + dirpath),
				Mode:     0755,
				ModTime:  t,
				Typeflag: tar.TypeDir,
			***REMOVED***)
		***REMOVED***

		for _, filePath := range paths ***REMOVED***
			data := files[filePath]
			if filePath[0] == '/' ***REMOVED***
				filePath = "_" + filePath
			***REMOVED***
			_ = w.WriteHeader(&tar.Header***REMOVED***
				Name:     path.Clean(entry.name + "/" + filePath),
				Mode:     0644,
				Size:     int64(len(data)),
				ModTime:  t,
				Typeflag: tar.TypeReg,
			***REMOVED***)
			if _, err := w.Write(data); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return w.Close()
***REMOVED***

func (arc *Archive) json() ([]byte, error) ***REMOVED***
	buffer := &bytes.Buffer***REMOVED******REMOVED***
	encoder := json.NewEncoder(buffer)
	// this prevents <, >, and & from being escaped in JSON strings
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(arc); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return buffer.Bytes(), nil
***REMOVED***
