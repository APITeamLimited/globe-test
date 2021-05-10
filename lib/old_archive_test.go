/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
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
	"net/url"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"go.k6.io/k6/lib/fsext"
)

func dumpMemMapFsToBuf(fs afero.Fs) (*bytes.Buffer, error) ***REMOVED***
	var b = bytes.NewBuffer(nil)
	var w = tar.NewWriter(b)
	err := fsext.Walk(fs, afero.FilePathSeparator,
		filepath.WalkFunc(func(filePath string, info os.FileInfo, err error) error ***REMOVED***
			if filePath == afero.FilePathSeparator ***REMOVED***
				return nil // skip the root
			***REMOVED***
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if info.IsDir() ***REMOVED***
				return w.WriteHeader(&tar.Header***REMOVED***
					Name:     path.Clean(filepath.ToSlash(filePath)[1:]),
					Mode:     0555,
					Typeflag: tar.TypeDir,
				***REMOVED***)
			***REMOVED***
			var data []byte
			data, err = afero.ReadFile(fs, filePath)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			err = w.WriteHeader(&tar.Header***REMOVED***
				Name:     path.Clean(filepath.ToSlash(filePath)[1:]),
				Mode:     0644,
				Size:     int64(len(data)),
				Typeflag: tar.TypeReg,
			***REMOVED***)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			_, err = w.Write(data)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			return nil
		***REMOVED***))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return b, w.Close()
***REMOVED***

func TestOldArchive(t *testing.T) ***REMOVED***
	var testCases = map[string]string***REMOVED***
		// map of filename to data for each main file tested
		"github.com/k6io/k6/samples/example.js": `github file`,
		"cdnjs.com/packages/Faker":              `faker file`,
		"C:/something/path2":                    `windows script`,
		"/absolulte/path2":                      `unix script`,
	***REMOVED***
	for filename, data := range testCases ***REMOVED***
		filename, data := filename, data
		t.Run(filename, func(t *testing.T) ***REMOVED***
			metadata := `***REMOVED***"filename": "` + filename + `", "options": ***REMOVED******REMOVED******REMOVED***`
			fs := makeMemMapFs(t, map[string][]byte***REMOVED***
				// files
				"/files/github.com/k6io/k6/samples/example.js": []byte(`github file`),
				"/files/cdnjs.com/packages/Faker":              []byte(`faker file`),
				"/files/example.com/path/to.js":                []byte(`example.com file`),
				"/files/_/C/something/path":                    []byte(`windows file`),
				"/files/_/absolulte/path":                      []byte(`unix file`),

				// scripts
				"/scripts/github.com/k6io/k6/samples/example.js2": []byte(`github script`),
				"/scripts/cdnjs.com/packages/Faker2":              []byte(`faker script`),
				"/scripts/example.com/path/too.js":                []byte(`example.com script`),
				"/scripts/_/C/something/path2":                    []byte(`windows script`),
				"/scripts/_/absolulte/path2":                      []byte(`unix script`),
				"/data":                                           []byte(data),
				"/metadata.json":                                  []byte(metadata),
			***REMOVED***)

			buf, err := dumpMemMapFsToBuf(fs)
			require.NoError(t, err)

			var (
				expectedFilesystems = map[string]afero.Fs***REMOVED***
					"file": makeMemMapFs(t, map[string][]byte***REMOVED***
						"/C:/something/path":  []byte(`windows file`),
						"/absolulte/path":     []byte(`unix file`),
						"/C:/something/path2": []byte(`windows script`),
						"/absolulte/path2":    []byte(`unix script`),
					***REMOVED***),
					"https": makeMemMapFs(t, map[string][]byte***REMOVED***
						"/example.com/path/to.js":                 []byte(`example.com file`),
						"/example.com/path/too.js":                []byte(`example.com script`),
						"/github.com/k6io/k6/samples/example.js":  []byte(`github file`),
						"/cdnjs.com/packages/Faker":               []byte(`faker file`),
						"/github.com/k6io/k6/samples/example.js2": []byte(`github script`),
						"/cdnjs.com/packages/Faker2":              []byte(`faker script`),
					***REMOVED***),
				***REMOVED***
			)

			arc, err := ReadArchive(buf)
			require.NoError(t, err)

			diffMapFilesystems(t, expectedFilesystems, arc.Filesystems)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestUnknownPrefix(t *testing.T) ***REMOVED***
	fs := makeMemMapFs(t, map[string][]byte***REMOVED***
		"/strange/something": []byte(`github file`),
	***REMOVED***)
	buf, err := dumpMemMapFsToBuf(fs)
	require.NoError(t, err)

	_, err = ReadArchive(buf)
	require.Error(t, err)
	require.Equal(t, err.Error(),
		"unknown file prefix `strange` for file `strange/something`")
***REMOVED***

func TestFilenamePwdResolve(t *testing.T) ***REMOVED***
	var tests = []struct ***REMOVED***
		Filename, Pwd, version              string
		expectedFilenameURL, expectedPwdURL *url.URL
		expectedError                       string
	***REMOVED******REMOVED***
		***REMOVED***
			Filename:            "/home/nobody/something.js",
			Pwd:                 "/home/nobody",
			expectedFilenameURL: &url.URL***REMOVED***Scheme: "file", Path: "/home/nobody/something.js"***REMOVED***,
			expectedPwdURL:      &url.URL***REMOVED***Scheme: "file", Path: "/home/nobody"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Filename:            "github.com/k6io/k6/samples/http2.js",
			Pwd:                 "github.com/k6io/k6/samples",
			expectedFilenameURL: &url.URL***REMOVED***Opaque: "github.com/k6io/k6/samples/http2.js"***REMOVED***,
			expectedPwdURL:      &url.URL***REMOVED***Opaque: "github.com/k6io/k6/samples"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Filename:            "cdnjs.com/libraries/Faker",
			Pwd:                 "/home/nobody",
			expectedFilenameURL: &url.URL***REMOVED***Opaque: "cdnjs.com/libraries/Faker"***REMOVED***,
			expectedPwdURL:      &url.URL***REMOVED***Scheme: "file", Path: "/home/nobody"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Filename:            "example.com/something/dot.js",
			Pwd:                 "example.com/something/",
			expectedFilenameURL: &url.URL***REMOVED***Host: "example.com", Scheme: "", Path: "/something/dot.js"***REMOVED***,
			expectedPwdURL:      &url.URL***REMOVED***Host: "example.com", Scheme: "", Path: "/something"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Filename:            "https://example.com/something/dot.js",
			Pwd:                 "https://example.com/something",
			expectedFilenameURL: &url.URL***REMOVED***Host: "example.com", Scheme: "https", Path: "/something/dot.js"***REMOVED***,
			expectedPwdURL:      &url.URL***REMOVED***Host: "example.com", Scheme: "https", Path: "/something"***REMOVED***,
			version:             "0.25.0",
		***REMOVED***,
		***REMOVED***
			Filename:      "ftps://example.com/something/dot.js",
			Pwd:           "https://example.com/something",
			expectedError: "only supported schemes for imports are file and https",
			version:       "0.25.0",
		***REMOVED***,
		***REMOVED***
			Filename:      "https://example.com/something/dot.js",
			Pwd:           "ftps://example.com/something",
			expectedError: "only supported schemes for imports are file and https",
			version:       "0.25.0",
		***REMOVED***,
	***REMOVED***

	for _, test := range tests ***REMOVED***
		metadata := `***REMOVED***
		"filename": "` + test.Filename + `",
		"pwd": "` + test.Pwd + `",
		"k6version": "` + test.version + `",
		"options": ***REMOVED******REMOVED***
	***REMOVED***`

		buf, err := dumpMemMapFsToBuf(makeMemMapFs(t, map[string][]byte***REMOVED***
			"/metadata.json": []byte(metadata),
		***REMOVED***))
		require.NoError(t, err)

		arc, err := ReadArchive(buf)
		if test.expectedError != "" ***REMOVED***
			require.Error(t, err)
			require.Contains(t, err.Error(), test.expectedError)
		***REMOVED*** else ***REMOVED***
			require.NoError(t, err)
			require.Equal(t, test.expectedFilenameURL, arc.FilenameURL)
			require.Equal(t, test.expectedPwdURL, arc.PwdURL)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestDerivedExecutionDiscarding(t *testing.T) ***REMOVED***
	var emptyConfigMap ScenarioConfigs
	var tests = []struct ***REMOVED***
		metadata     string
		expScenarios interface***REMOVED******REMOVED***
		expError     string
	***REMOVED******REMOVED***
		// Tests to make sure that "execution" in the options, the old name for
		// "scenarios" before #1007 was merged, doesn't mess up the options...
		***REMOVED***
			metadata: `***REMOVED***
				"filename": "/test.js", "pwd": "/",
				"options": ***REMOVED*** "execution": ***REMOVED*** "something": "invalid" ***REMOVED*** ***REMOVED***
			***REMOVED***`,
			expScenarios: emptyConfigMap,
		***REMOVED***,
		***REMOVED***
			metadata: `***REMOVED***
				"filename": "/test.js", "pwd": "/",
				"k6version": "0.24.0",
				"options": ***REMOVED*** "execution": ***REMOVED*** "something": "invalid" ***REMOVED*** ***REMOVED***
			***REMOVED***`,
			expScenarios: emptyConfigMap,
		***REMOVED***,
		***REMOVED***
			metadata: `blah`,
			expError: "invalid character",
		***REMOVED***,
		***REMOVED***
			metadata: `***REMOVED***
				"filename": "/test.js", "pwd": "/",
				"k6version": "0.24.0",
				"options": "something invalid"
			***REMOVED***`,
			expError: "cannot unmarshal string into Go struct field",
		***REMOVED***,
		***REMOVED***
			metadata: `***REMOVED***
				"filename": "/test.js", "pwd": "/",
				"k6version": "0.25.0",
				"options": ***REMOVED*** "scenarios": ***REMOVED*** "something": "invalid" ***REMOVED*** ***REMOVED***
			***REMOVED***`,
			expError: "cannot unmarshal string",
		***REMOVED***,
		// TODO: test an actual scenarios unmarshalling, which is currently
		// impossible due to import cycles...
	***REMOVED***

	for _, test := range tests ***REMOVED***
		buf, err := dumpMemMapFsToBuf(makeMemMapFs(t, map[string][]byte***REMOVED***
			"/metadata.json": []byte(test.metadata),
		***REMOVED***))
		require.NoError(t, err)

		arc, err := ReadArchive(buf)
		if test.expError != "" ***REMOVED***
			require.Errorf(t, err, "expected error '%s' but got nil", test.expError)
			require.Contains(t, err.Error(), test.expError)
		***REMOVED*** else ***REMOVED***
			require.NoError(t, err)
			require.Equal(t, test.expScenarios, arc.Options.Scenarios)
		***REMOVED***
	***REMOVED***
***REMOVED***
