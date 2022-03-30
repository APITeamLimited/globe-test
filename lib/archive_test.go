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
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/lib/consts"
	"go.k6.io/k6/lib/fsext"
	"go.k6.io/k6/metrics"
)

func TestNormalizeAndAnonymizePath(t *testing.T) ***REMOVED***
	testdata := map[string]string***REMOVED***
		"/tmp":                            "/tmp",
		"/tmp/myfile.txt":                 "/tmp/myfile.txt",
		"/home/myname":                    "/home/nobody",
		"/home/myname/foo/bar/myfile.txt": "/home/nobody/foo/bar/myfile.txt",
		"/Users/myname/myfile.txt":        "/Users/nobody/myfile.txt",
		"/Documents and Settings/myname/myfile.txt":           "/Documents and Settings/nobody/myfile.txt",
		"\\\\MYSHARED\\dir\\dir\\myfile.txt":                  "/nobody/dir/dir/myfile.txt",
		"\\NOTSHARED\\dir\\dir\\myfile.txt":                   "/NOTSHARED/dir/dir/myfile.txt",
		"C:\\Users\\myname\\dir\\myfile.txt":                  "/C/Users/nobody/dir/myfile.txt",
		"D:\\Documents and Settings\\myname\\dir\\myfile.txt": "/D/Documents and Settings/nobody/dir/myfile.txt",
		"C:\\uSers\\myname\\dir\\myfile.txt":                  "/C/uSers/nobody/dir/myfile.txt",
		"D:\\doCUMENts aND Settings\\myname\\dir\\myfile.txt": "/D/doCUMENts aND Settings/nobody/dir/myfile.txt",
	***REMOVED***
	// TODO: fix this - the issue is that filepath.Clean replaces `/` with whatever the path
	// separator is on the current OS and as such this gets confused for shared folder on
	// windows :( https://github.com/golang/go/issues/16111
	if runtime.GOOS != "windows" ***REMOVED***
		testdata["//etc/hosts"] = "/etc/hosts"
	***REMOVED***
	for from, to := range testdata ***REMOVED***
		from, to := from, to
		t.Run("path="+from, func(t *testing.T) ***REMOVED***
			res := NormalizeAndAnonymizePath(from)
			assert.Equal(t, to, res)
			assert.Equal(t, res, NormalizeAndAnonymizePath(res))
		***REMOVED***)
	***REMOVED***
***REMOVED***

func makeMemMapFs(t *testing.T, input map[string][]byte) afero.Fs ***REMOVED***
	fs := afero.NewMemMapFs()
	for path, data := range input ***REMOVED***
		require.NoError(t, afero.WriteFile(fs, path, data, 0644))
	***REMOVED***
	return fs
***REMOVED***

func getMapKeys(m map[string]afero.Fs) []string ***REMOVED***
	keys := make([]string, 0, len(m))
	for key := range m ***REMOVED***
		keys = append(keys, key)
	***REMOVED***

	return keys
***REMOVED***

func diffMapFilesystems(t *testing.T, first, second map[string]afero.Fs) bool ***REMOVED***
	require.ElementsMatch(t, getMapKeys(first), getMapKeys(second),
		"fs map keys don't match %s, %s", getMapKeys(first), getMapKeys(second))
	for key, fs := range first ***REMOVED***
		secondFs := second[key]
		diffFilesystems(t, fs, secondFs)
	***REMOVED***

	return true
***REMOVED***

func diffFilesystems(t *testing.T, first, second afero.Fs) ***REMOVED***
	diffFilesystemsDir(t, first, second, "/")
***REMOVED***

func getInfoNames(infos []os.FileInfo) []string ***REMOVED***
	var names = make([]string, len(infos))
	for i, info := range infos ***REMOVED***
		names[i] = info.Name()
	***REMOVED***
	return names
***REMOVED***

func diffFilesystemsDir(t *testing.T, first, second afero.Fs, dirname string) ***REMOVED***
	firstInfos, err := afero.ReadDir(first, dirname)
	require.NoError(t, err, dirname)

	secondInfos, err := afero.ReadDir(first, dirname)
	require.NoError(t, err, dirname)

	require.ElementsMatch(t, getInfoNames(firstInfos), getInfoNames(secondInfos), "directory: "+dirname)
	for _, info := range firstInfos ***REMOVED***
		path := filepath.Join(dirname, info.Name())
		if info.IsDir() ***REMOVED***
			diffFilesystemsDir(t, first, second, path)
			continue
		***REMOVED***
		firstData, err := afero.ReadFile(first, path)
		require.NoError(t, err, path)

		secondData, err := afero.ReadFile(second, path)
		require.NoError(t, err, path)

		assert.Equal(t, firstData, secondData, path)
	***REMOVED***
***REMOVED***

func TestArchiveReadWrite(t *testing.T) ***REMOVED***
	t.Run("Roundtrip", func(t *testing.T) ***REMOVED***
		arc1 := &Archive***REMOVED***
			Type:      "js",
			K6Version: consts.Version,
			Options: Options***REMOVED***
				VUs:        null.IntFrom(12345),
				SystemTags: &metrics.DefaultSystemTagSet,
			***REMOVED***,
			FilenameURL: &url.URL***REMOVED***Scheme: "file", Path: "/path/to/a.js"***REMOVED***,
			Data:        []byte(`// a contents`),
			PwdURL:      &url.URL***REMOVED***Scheme: "file", Path: "/path/to"***REMOVED***,
			Filesystems: map[string]afero.Fs***REMOVED***
				"file": makeMemMapFs(t, map[string][]byte***REMOVED***
					"/path/to/a.js":      []byte(`// a contents`),
					"/path/to/b.js":      []byte(`// b contents`),
					"/path/to/file1.txt": []byte(`hi!`),
					"/path/to/file2.txt": []byte(`bye!`),
				***REMOVED***),
				"https": makeMemMapFs(t, map[string][]byte***REMOVED***
					"/cdnjs.com/libraries/Faker":          []byte(`// faker contents`),
					"/github.com/loadimpact/k6/README.md": []byte(`README`),
				***REMOVED***),
			***REMOVED***,
		***REMOVED***

		buf := bytes.NewBuffer(nil)
		require.NoError(t, arc1.Write(buf))

		arc1Filesystems := arc1.Filesystems
		arc1.Filesystems = nil

		arc2, err := ReadArchive(buf)
		require.NoError(t, err)

		arc2Filesystems := arc2.Filesystems
		arc2.Filesystems = nil
		arc2.Filename = ""
		arc2.Pwd = ""

		assert.Equal(t, arc1, arc2)

		diffMapFilesystems(t, arc1Filesystems, arc2Filesystems)
	***REMOVED***)

	t.Run("Anonymized", func(t *testing.T) ***REMOVED***
		testdata := []struct ***REMOVED***
			Pwd, PwdNormAnon string
		***REMOVED******REMOVED***
			***REMOVED***"/home/myname", "/home/nobody"***REMOVED***,
			***REMOVED***filepath.FromSlash("/C:/Users/Administrator"), "/C/Users/nobody"***REMOVED***,
		***REMOVED***
		for _, entry := range testdata ***REMOVED***
			arc1 := &Archive***REMOVED***
				Type: "js",
				Options: Options***REMOVED***
					VUs:        null.IntFrom(12345),
					SystemTags: &metrics.DefaultSystemTagSet,
				***REMOVED***,
				FilenameURL: &url.URL***REMOVED***Scheme: "file", Path: fmt.Sprintf("%s/a.js", entry.Pwd)***REMOVED***,
				K6Version:   consts.Version,
				Data:        []byte(`// a contents`),
				PwdURL:      &url.URL***REMOVED***Scheme: "file", Path: entry.Pwd***REMOVED***,
				Filesystems: map[string]afero.Fs***REMOVED***
					"file": makeMemMapFs(t, map[string][]byte***REMOVED***
						fmt.Sprintf("%s/a.js", entry.Pwd):      []byte(`// a contents`),
						fmt.Sprintf("%s/b.js", entry.Pwd):      []byte(`// b contents`),
						fmt.Sprintf("%s/file1.txt", entry.Pwd): []byte(`hi!`),
						fmt.Sprintf("%s/file2.txt", entry.Pwd): []byte(`bye!`),
					***REMOVED***),
					"https": makeMemMapFs(t, map[string][]byte***REMOVED***
						"/cdnjs.com/libraries/Faker":          []byte(`// faker contents`),
						"/github.com/loadimpact/k6/README.md": []byte(`README`),
					***REMOVED***),
				***REMOVED***,
			***REMOVED***
			arc1Anon := &Archive***REMOVED***
				Type: "js",
				Options: Options***REMOVED***
					VUs:        null.IntFrom(12345),
					SystemTags: &metrics.DefaultSystemTagSet,
				***REMOVED***,
				FilenameURL: &url.URL***REMOVED***Scheme: "file", Path: fmt.Sprintf("%s/a.js", entry.PwdNormAnon)***REMOVED***,
				K6Version:   consts.Version,
				Data:        []byte(`// a contents`),
				PwdURL:      &url.URL***REMOVED***Scheme: "file", Path: entry.PwdNormAnon***REMOVED***,

				Filesystems: map[string]afero.Fs***REMOVED***
					"file": makeMemMapFs(t, map[string][]byte***REMOVED***
						fmt.Sprintf("%s/a.js", entry.PwdNormAnon):      []byte(`// a contents`),
						fmt.Sprintf("%s/b.js", entry.PwdNormAnon):      []byte(`// b contents`),
						fmt.Sprintf("%s/file1.txt", entry.PwdNormAnon): []byte(`hi!`),
						fmt.Sprintf("%s/file2.txt", entry.PwdNormAnon): []byte(`bye!`),
					***REMOVED***),
					"https": makeMemMapFs(t, map[string][]byte***REMOVED***
						"/cdnjs.com/libraries/Faker":          []byte(`// faker contents`),
						"/github.com/loadimpact/k6/README.md": []byte(`README`),
					***REMOVED***),
				***REMOVED***,
			***REMOVED***

			buf := bytes.NewBuffer(nil)
			require.NoError(t, arc1.Write(buf))

			arc1Filesystems := arc1Anon.Filesystems
			arc1Anon.Filesystems = nil

			arc2, err := ReadArchive(buf)
			assert.NoError(t, err)
			arc2.Filename = ""
			arc2.Pwd = ""

			arc2Filesystems := arc2.Filesystems
			arc2.Filesystems = nil

			assert.Equal(t, arc1Anon, arc2)
			diffMapFilesystems(t, arc1Filesystems, arc2Filesystems)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestArchiveJSONEscape(t *testing.T) ***REMOVED***
	t.Parallel()

	arc := &Archive***REMOVED******REMOVED***
	arc.Filename = "test<.js"
	b, err := arc.json()
	assert.NoError(t, err)
	assert.Contains(t, string(b), "test<.js")
***REMOVED***

func TestUsingCacheFromCacheOnReadFs(t *testing.T) ***REMOVED***
	var base = afero.NewMemMapFs()
	var cached = afero.NewMemMapFs()
	// we specifically have different contents in both places
	require.NoError(t, afero.WriteFile(base, "/wrong", []byte(`ooops`), 0644))
	require.NoError(t, afero.WriteFile(cached, "/correct", []byte(`test`), 0644))

	arc := &Archive***REMOVED***
		Type:        "js",
		FilenameURL: &url.URL***REMOVED***Scheme: "file", Path: "/correct"***REMOVED***,
		K6Version:   consts.Version,
		Data:        []byte(`test`),
		PwdURL:      &url.URL***REMOVED***Scheme: "file", Path: "/"***REMOVED***,
		Filesystems: map[string]afero.Fs***REMOVED***
			"file": fsext.NewCacheOnReadFs(base, cached, 0),
		***REMOVED***,
	***REMOVED***

	buf := bytes.NewBuffer(nil)
	require.NoError(t, arc.Write(buf))

	newArc, err := ReadArchive(buf)
	require.NoError(t, err)

	data, err := afero.ReadFile(newArc.Filesystems["file"], "/correct")
	require.NoError(t, err)
	require.Equal(t, string(data), "test")

	data, err = afero.ReadFile(newArc.Filesystems["file"], "/wrong")
	require.Error(t, err)
	require.Nil(t, data)
***REMOVED***

func TestArchiveWithDataNotInFS(t *testing.T) ***REMOVED***
	t.Parallel()

	arc := &Archive***REMOVED***
		Type:        "js",
		FilenameURL: &url.URL***REMOVED***Scheme: "file", Path: "/script"***REMOVED***,
		K6Version:   consts.Version,
		Data:        []byte(`test`),
		PwdURL:      &url.URL***REMOVED***Scheme: "file", Path: "/"***REMOVED***,
		Filesystems: nil,
	***REMOVED***

	buf := bytes.NewBuffer(nil)
	err := arc.Write(buf)
	require.Error(t, err)
	require.Contains(t, err.Error(), "the main script wasn't present in the cached filesystem")
***REMOVED***

func TestMalformedMetadata(t *testing.T) ***REMOVED***
	var fs = afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "/metadata.json", []byte("***REMOVED***,***REMOVED***"), 0644))
	var b, err = dumpMemMapFsToBuf(fs)
	require.NoError(t, err)
	_, err = ReadArchive(b)
	require.Error(t, err)
	require.Equal(t, err.Error(), `invalid character ',' looking for beginning of object key string`)
***REMOVED***

func TestStrangePaths(t *testing.T) ***REMOVED***
	var pathsToChange = []string***REMOVED***
		`/path/with spaces/a.js`,
		`/path/with spaces/a.js`,
		`/path/with日本語/b.js`,
		`/path/with spaces and 日本語/file1.txt`,
	***REMOVED***
	for _, pathToChange := range pathsToChange ***REMOVED***
		otherMap := make(map[string][]byte, len(pathsToChange))
		for _, other := range pathsToChange ***REMOVED***
			otherMap[other] = []byte(`// ` + other + ` contents`)
		***REMOVED***
		arc1 := &Archive***REMOVED***
			Type:      "js",
			K6Version: consts.Version,
			Options: Options***REMOVED***
				VUs:        null.IntFrom(12345),
				SystemTags: &metrics.DefaultSystemTagSet,
			***REMOVED***,
			FilenameURL: &url.URL***REMOVED***Scheme: "file", Path: pathToChange***REMOVED***,
			Data:        []byte(`// ` + pathToChange + ` contents`),
			PwdURL:      &url.URL***REMOVED***Scheme: "file", Path: path.Dir(pathToChange)***REMOVED***,
			Filesystems: map[string]afero.Fs***REMOVED***
				"file": makeMemMapFs(t, otherMap),
			***REMOVED***,
		***REMOVED***

		buf := bytes.NewBuffer(nil)
		require.NoError(t, arc1.Write(buf), pathToChange)

		arc1Filesystems := arc1.Filesystems
		arc1.Filesystems = nil

		arc2, err := ReadArchive(buf)
		require.NoError(t, err, pathToChange)

		arc2Filesystems := arc2.Filesystems
		arc2.Filesystems = nil
		arc2.Filename = ""
		arc2.Pwd = ""

		assert.Equal(t, arc1, arc2, pathToChange)

		arc1Filesystems["https"] = afero.NewMemMapFs()
		diffMapFilesystems(t, arc1Filesystems, arc2Filesystems)
	***REMOVED***
***REMOVED***

func TestStdinArchive(t *testing.T) ***REMOVED***
	var fs = afero.NewMemMapFs()
	// we specifically have different contents in both places
	require.NoError(t, afero.WriteFile(fs, "/-", []byte(`test`), 0644))

	arc := &Archive***REMOVED***
		Type:        "js",
		FilenameURL: &url.URL***REMOVED***Scheme: "file", Path: "/-"***REMOVED***,
		K6Version:   consts.Version,
		Data:        []byte(`test`),
		PwdURL:      &url.URL***REMOVED***Scheme: "file", Path: "/"***REMOVED***,
		Filesystems: map[string]afero.Fs***REMOVED***
			"file": fs,
		***REMOVED***,
	***REMOVED***

	buf := bytes.NewBuffer(nil)
	require.NoError(t, arc.Write(buf))

	newArc, err := ReadArchive(buf)
	require.NoError(t, err)

	data, err := afero.ReadFile(newArc.Filesystems["file"], "/-")
	require.NoError(t, err)
	require.Equal(t, string(data), "test")

***REMOVED***
