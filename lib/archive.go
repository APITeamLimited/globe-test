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
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/spf13/afero"

	"go.k6.io/k6/lib/fsext"
	"go.k6.io/k6/loader"
)

//nolint: gochecknoglobals, lll
var (
	volumeRE  = regexp.MustCompile(`^[/\\]?([a-zA-Z]):(.*)`)
	sharedRE  = regexp.MustCompile(`^\\\\([^\\]+)`) // matches a shared folder in Windows before backslack replacement. i.e \\VMBOXSVR\k6\script.js
	homeDirRE = regexp.MustCompile(`(?i)^(/[a-zA-Z])?/(Users|home|Documents and Settings)/(?:[^/]+)`)
)

// NormalizeAndAnonymizePath Normalizes (to use a / path separator) and anonymizes a file path, by scrubbing usernames from home directories.
func NormalizeAndAnonymizePath(path string) string ***REMOVED***
	path = filepath.Clean(path)

	p := volumeRE.ReplaceAllString(path, `/$1$2`)
	p = sharedRE.ReplaceAllString(p, `/nobody`)
	p = strings.Replace(p, "\\", "/", -1)
	return homeDirRE.ReplaceAllString(p, `$1/$2/nobody`)
***REMOVED***

func newNormalizedFs(fs afero.Fs) afero.Fs ***REMOVED***
	return fsext.NewChangePathFs(fs, fsext.ChangePathFunc(func(name string) (string, error) ***REMOVED***
		return NormalizeAndAnonymizePath(name), nil
	***REMOVED***))
***REMOVED***

// An Archive is a rollup of all resources and options needed to reproduce a test identically elsewhere.
type Archive struct ***REMOVED***
	// The runner to use, eg. "js".
	Type string `json:"type"`

	// Options to use.
	Options Options `json:"options"`

	// TODO: rewrite the encoding, decoding of json to use another type with only the fields it
	// needs in order to remove Filename and Pwd from this
	// Filename and contents of the main file being executed.
	Filename    string   `json:"filename"` // only for json
	FilenameURL *url.URL `json:"-"`
	Data        []byte   `json:"-"`

	// Working directory for resolving relative paths.
	Pwd    string   `json:"pwd"` // only for json
	PwdURL *url.URL `json:"-"`

	Filesystems map[string]afero.Fs `json:"-"`

	// Environment variables
	Env map[string]string `json:"env"`

	CompatibilityMode string `json:"compatibilityMode"`

	K6Version string `json:"k6version"`
	Goos      string `json:"goos"`
***REMOVED***

func (arc *Archive) getFs(name string) afero.Fs ***REMOVED***
	fs, ok := arc.Filesystems[name]
	if !ok ***REMOVED***
		fs = afero.NewMemMapFs()
		if name == "file" ***REMOVED***
			fs = newNormalizedFs(fs)
		***REMOVED***
		arc.Filesystems[name] = fs
	***REMOVED***

	return fs
***REMOVED***

func (arc *Archive) loadMetadataJSON(data []byte) (err error) ***REMOVED***
	if err = json.Unmarshal(data, &arc); err != nil ***REMOVED***
		return err
	***REMOVED***
	// Path separator normalization for older archives (<=0.20.0)
	if arc.K6Version == "" ***REMOVED***
		arc.Filename = NormalizeAndAnonymizePath(arc.Filename)
		arc.Pwd = NormalizeAndAnonymizePath(arc.Pwd)
	***REMOVED***
	arc.PwdURL, err = loader.Resolve(&url.URL***REMOVED***Scheme: "file", Path: "/"***REMOVED***, arc.Pwd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	arc.FilenameURL, err = loader.Resolve(&url.URL***REMOVED***Scheme: "file", Path: "/"***REMOVED***, arc.Filename)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

// ReadArchive reads an archive created by Archive.Write from a reader.
func ReadArchive(in io.Reader) (*Archive, error) ***REMOVED***
	r := tar.NewReader(in)
	arc := &Archive***REMOVED***Filesystems: make(map[string]afero.Fs, 2)***REMOVED***
	// initialize both fses
	_ = arc.getFs("https")
	_ = arc.getFs("file")
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
			if err = arc.loadMetadataJSON(data); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			continue
		case "data":
			arc.Data = data
			continue
		***REMOVED***

		// Path separator normalization for older archives (<=0.20.0)
		normPath := NormalizeAndAnonymizePath(hdr.Name)
		idx := strings.IndexRune(normPath, '/')
		if idx == -1 ***REMOVED***
			continue
		***REMOVED***
		pfx := normPath[:idx]
		name := normPath[idx:]

		switch pfx ***REMOVED***
		case "files", "scripts": // old archives
			// in old archives (pre 0.25.0) names without "_" at the beginning were  https, the ones with "_" are local files
			pfx = "https"
			if len(name) >= 2 && name[0:2] == "/_" ***REMOVED***
				pfx = "file"
				name = name[2:]
			***REMOVED***
			fallthrough
		case "https", "file":
			fs := arc.getFs(pfx)
			name = filepath.FromSlash(name)
			err = afero.WriteFile(fs, name, data, os.FileMode(hdr.Mode))
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			err = fs.Chtimes(name, hdr.AccessTime, hdr.ModTime)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		default:
			return nil, fmt.Errorf("unknown file prefix `%s` for file `%s`", pfx, normPath)
		***REMOVED***
	***REMOVED***
	scheme, pathOnFs := getURLPathOnFs(arc.FilenameURL)
	var err error
	pathOnFs, err = url.PathUnescape(pathOnFs)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	err = afero.WriteFile(arc.getFs(scheme), pathOnFs, arc.Data, 0o644) // TODO fix the mode ?
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return arc, nil
***REMOVED***

func normalizeAndAnonymizeURL(u *url.URL) ***REMOVED***
	if u.Scheme == "file" ***REMOVED***
		u.Path = NormalizeAndAnonymizePath(u.Path)
	***REMOVED***
***REMOVED***

func getURLPathOnFs(u *url.URL) (scheme string, pathOnFs string) ***REMOVED***
	scheme = "https"
	switch ***REMOVED***
	case u.Opaque != "":
		return scheme, "/" + u.Opaque
	case u.Scheme == "":
		return scheme, path.Clean(u.String()[len("//"):])
	default:
		scheme = u.Scheme
	***REMOVED***
	return scheme, path.Clean(u.String()[len(u.Scheme)+len(":/"):])
***REMOVED***

func getURLtoString(u *url.URL) string ***REMOVED***
	if u.Opaque == "" && u.Scheme == "" ***REMOVED***
		return u.String()[len("//"):] // https url without a scheme
	***REMOVED***
	return u.String()
***REMOVED***

// Write serialises the archive to a writer.
//
// The format should be treated as opaque; currently it is simply a TAR rollup, but this may
// change. If it does change, ReadArchive must be able to handle all previous formats as well as
// the current one.
func (arc *Archive) Write(out io.Writer) error ***REMOVED***
	w := tar.NewWriter(out)

	now := time.Now()
	metaArc := *arc
	normalizeAndAnonymizeURL(metaArc.FilenameURL)
	normalizeAndAnonymizeURL(metaArc.PwdURL)
	metaArc.Filename = getURLtoString(metaArc.FilenameURL)
	metaArc.Pwd = getURLtoString(metaArc.PwdURL)
	actualDataPath, err := url.PathUnescape(path.Join(getURLPathOnFs(metaArc.FilenameURL)))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	var madeLinkToData bool
	metadata, err := metaArc.json()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_ = w.WriteHeader(&tar.Header***REMOVED***
		Name:     "metadata.json",
		Mode:     0o644,
		Size:     int64(len(metadata)),
		ModTime:  now,
		Typeflag: tar.TypeReg,
	***REMOVED***)
	if _, err = w.Write(metadata); err != nil ***REMOVED***
		return err
	***REMOVED***

	_ = w.WriteHeader(&tar.Header***REMOVED***
		Name:     "data",
		Mode:     0o644,
		Size:     int64(len(arc.Data)),
		ModTime:  now,
		Typeflag: tar.TypeReg,
	***REMOVED***)
	if _, err = w.Write(arc.Data); err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, name := range [...]string***REMOVED***"file", "https"***REMOVED*** ***REMOVED***
		filesystem, ok := arc.Filesystems[name]
		if !ok ***REMOVED***
			continue
		***REMOVED***
		if cachedfs, ok := filesystem.(fsext.CacheLayerGetter); ok ***REMOVED***
			filesystem = cachedfs.GetCachingFs()
		***REMOVED***

		// A couple of things going on here:
		// - You can't just create file entries, you need to create directory entries too.
		//   Figure out which directories are in use here.
		// - We want archives to be comparable by hash, which means the entries need to be written
		//   in the same order every time. Go maps are shuffled, so we need to sort lists of keys.
		// - We don't want to leak private information (eg. usernames) in archives, so make sure to
		//   anonymize paths before stuffing them in a shareable archive.
		foundDirs := make(map[string]bool)
		paths := make([]string, 0, 10)
		infos := make(map[string]os.FileInfo) // ... fix this ?
		files := make(map[string][]byte)

		walkFunc := filepath.WalkFunc(func(filePath string, info os.FileInfo, err error) error ***REMOVED***
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			normalizedPath := NormalizeAndAnonymizePath(filePath)

			infos[normalizedPath] = info
			if info.IsDir() ***REMOVED***
				foundDirs[normalizedPath] = true
				return nil
			***REMOVED***

			paths = append(paths, normalizedPath)
			files[normalizedPath], err = afero.ReadFile(filesystem, filePath)
			return err
		***REMOVED***)

		if err = fsext.Walk(filesystem, afero.FilePathSeparator, walkFunc); err != nil ***REMOVED***
			return err
		***REMOVED***
		if len(files) == 0 ***REMOVED***
			continue // we don't need to write anything for this fs, if this is not done the root will be written
		***REMOVED***
		dirs := make([]string, 0, len(foundDirs))
		for dirpath := range foundDirs ***REMOVED***
			dirs = append(dirs, dirpath)
		***REMOVED***
		sort.Strings(paths)
		sort.Strings(dirs)

		for _, dirPath := range dirs ***REMOVED***
			_ = w.WriteHeader(&tar.Header***REMOVED***
				Name:       path.Clean(path.Join(name, dirPath)),
				Mode:       0o755, // MemMapFs is buggy
				AccessTime: now,   // MemMapFs is buggy
				ChangeTime: now,   // MemMapFs is buggy
				ModTime:    now,   // MemMapFs is buggy
				Typeflag:   tar.TypeDir,
			***REMOVED***)
		***REMOVED***

		for _, filePath := range paths ***REMOVED***
			fullFilePath := path.Clean(path.Join(name, filePath))
			// we either have opaque
			if fullFilePath == actualDataPath ***REMOVED***
				madeLinkToData = true
				err = w.WriteHeader(&tar.Header***REMOVED***
					Name:     fullFilePath,
					Size:     0,
					Typeflag: tar.TypeLink,
					Linkname: "data",
				***REMOVED***)
			***REMOVED*** else ***REMOVED***
				err = w.WriteHeader(&tar.Header***REMOVED***
					Name:       fullFilePath,
					Mode:       0o644, // MemMapFs is buggy
					Size:       int64(len(files[filePath])),
					AccessTime: infos[filePath].ModTime(),
					ChangeTime: infos[filePath].ModTime(),
					ModTime:    infos[filePath].ModTime(),
					Typeflag:   tar.TypeReg,
				***REMOVED***)
				if err == nil ***REMOVED***
					_, err = w.Write(files[filePath])
				***REMOVED***
			***REMOVED***
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if !madeLinkToData ***REMOVED***
		// This should never happen we should always link to `data` from inside the file/https directories
		return fmt.Errorf("archive creation failed because the main script wasn't present in the cached filesystem")
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
