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

package js

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/consts"
	"github.com/loadimpact/k6/lib/fsext"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/loader"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"
)

const isWindows = runtime.GOOS == "windows"

func getSimpleBundle(filename, data string) (*Bundle, error) ***REMOVED***
	return getSimpleBundleWithFs(filename, data, afero.NewMemMapFs())
***REMOVED***

func getSimpleBundleWithOptions(filename, data string, options lib.RuntimeOptions) (*Bundle, error) ***REMOVED***
	return NewBundle(
		&loader.SourceData***REMOVED***
			URL:  &url.URL***REMOVED***Path: filename, Scheme: "file"***REMOVED***,
			Data: []byte(data),
		***REMOVED***,
		map[string]afero.Fs***REMOVED***"file": afero.NewMemMapFs(), "https": afero.NewMemMapFs()***REMOVED***,
		options,
	)
***REMOVED***

func getSimpleBundleWithFs(filename, data string, fs afero.Fs) (*Bundle, error) ***REMOVED***
	return NewBundle(
		&loader.SourceData***REMOVED***
			URL:  &url.URL***REMOVED***Path: filename, Scheme: "file"***REMOVED***,
			Data: []byte(data),
		***REMOVED***,
		map[string]afero.Fs***REMOVED***"file": fs, "https": afero.NewMemMapFs()***REMOVED***,
		lib.RuntimeOptions***REMOVED******REMOVED***,
	)
***REMOVED***

func TestNewBundle(t *testing.T) ***REMOVED***
	t.Run("Blank", func(t *testing.T) ***REMOVED***
		_, err := getSimpleBundle("/script.js", "")
		assert.EqualError(t, err, "script must export a default function")
	***REMOVED***)
	t.Run("Invalid", func(t *testing.T) ***REMOVED***
		_, err := getSimpleBundle("/script.js", "\x00")
		assert.Contains(t, err.Error(), "SyntaxError: file:///script.js: Unexpected character '\x00' (1:0)\n> 1 | \x00\n")
	***REMOVED***)
	t.Run("Error", func(t *testing.T) ***REMOVED***
		_, err := getSimpleBundle("/script.js", `throw new Error("aaaa");`)
		assert.EqualError(t, err, "Error: aaaa at file:///script.js:1:7(3)")
	***REMOVED***)
	t.Run("InvalidExports", func(t *testing.T) ***REMOVED***
		_, err := getSimpleBundle("/script.js", `exports = null`)
		assert.EqualError(t, err, "exports must be an object")
	***REMOVED***)
	t.Run("DefaultUndefined", func(t *testing.T) ***REMOVED***
		_, err := getSimpleBundle("/script.js", `export default undefined;`)
		assert.EqualError(t, err, "script must export a default function")
	***REMOVED***)
	t.Run("DefaultNull", func(t *testing.T) ***REMOVED***
		_, err := getSimpleBundle("/script.js", `export default null;`)
		assert.EqualError(t, err, "script must export a default function")
	***REMOVED***)
	t.Run("DefaultWrongType", func(t *testing.T) ***REMOVED***
		_, err := getSimpleBundle("/script.js", `export default 12345;`)
		assert.EqualError(t, err, "default export must be a function")
	***REMOVED***)
	t.Run("Minimal", func(t *testing.T) ***REMOVED***
		_, err := getSimpleBundle("/script.js", `export default function() ***REMOVED******REMOVED***;`)
		assert.NoError(t, err)
	***REMOVED***)
	t.Run("stdin", func(t *testing.T) ***REMOVED***
		b, err := getSimpleBundle("-", `export default function() ***REMOVED******REMOVED***;`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "file://-", b.Filename.String())
			assert.Equal(t, "file:///", b.BaseInitContext.pwd.String())
		***REMOVED***
	***REMOVED***)
	t.Run("Options", func(t *testing.T) ***REMOVED***
		t.Run("Empty", func(t *testing.T) ***REMOVED***
			_, err := getSimpleBundle("/script.js", `
				export let options = ***REMOVED******REMOVED***;
				export default function() ***REMOVED******REMOVED***;
			`)
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("Invalid", func(t *testing.T) ***REMOVED***
			invalidOptions := map[string]struct ***REMOVED***
				Expr, Error string
			***REMOVED******REMOVED***
				"Array":    ***REMOVED***`[]`, "json: cannot unmarshal array into Go value of type lib.Options"***REMOVED***,
				"Function": ***REMOVED***`function()***REMOVED******REMOVED***`, "json: unsupported type: func(goja.FunctionCall) goja.Value"***REMOVED***,
			***REMOVED***
			for name, data := range invalidOptions ***REMOVED***
				t.Run(name, func(t *testing.T) ***REMOVED***
					_, err := getSimpleBundle("/script.js", fmt.Sprintf(`
						export let options = %s;
						export default function() ***REMOVED******REMOVED***;
					`, data.Expr))
					assert.EqualError(t, err, data.Error)
				***REMOVED***)
			***REMOVED***
		***REMOVED***)

		t.Run("Paused", func(t *testing.T) ***REMOVED***
			b, err := getSimpleBundle("/script.js", `
				export let options = ***REMOVED***
					paused: true,
				***REMOVED***;
				export default function() ***REMOVED******REMOVED***;
			`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.BoolFrom(true), b.Options.Paused)
			***REMOVED***
		***REMOVED***)
		t.Run("VUs", func(t *testing.T) ***REMOVED***
			b, err := getSimpleBundle("/script.js", `
				export let options = ***REMOVED***
					vus: 100,
				***REMOVED***;
				export default function() ***REMOVED******REMOVED***;
			`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.IntFrom(100), b.Options.VUs)
			***REMOVED***
		***REMOVED***)
		t.Run("Duration", func(t *testing.T) ***REMOVED***
			b, err := getSimpleBundle("/script.js", `
				export let options = ***REMOVED***
					duration: "10s",
				***REMOVED***;
				export default function() ***REMOVED******REMOVED***;
			`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, types.NullDurationFrom(10*time.Second), b.Options.Duration)
			***REMOVED***
		***REMOVED***)
		t.Run("Iterations", func(t *testing.T) ***REMOVED***
			b, err := getSimpleBundle("/script.js", `
				export let options = ***REMOVED***
					iterations: 100,
				***REMOVED***;
				export default function() ***REMOVED******REMOVED***;
			`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.IntFrom(100), b.Options.Iterations)
			***REMOVED***
		***REMOVED***)
		t.Run("Stages", func(t *testing.T) ***REMOVED***
			b, err := getSimpleBundle("/script.js", `
				export let options = ***REMOVED***
					stages: [],
				***REMOVED***;
				export default function() ***REMOVED******REMOVED***;
			`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Len(t, b.Options.Stages, 0)
			***REMOVED***

			t.Run("Empty", func(t *testing.T) ***REMOVED***
				b, err := getSimpleBundle("/script.js", `
					export let options = ***REMOVED***
						stages: [
							***REMOVED******REMOVED***,
						],
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`)
				if assert.NoError(t, err) ***REMOVED***
					if assert.Len(t, b.Options.Stages, 1) ***REMOVED***
						assert.Equal(t, lib.Stage***REMOVED******REMOVED***, b.Options.Stages[0])
					***REMOVED***
				***REMOVED***
			***REMOVED***)
			t.Run("Target", func(t *testing.T) ***REMOVED***
				b, err := getSimpleBundle("/script.js", `
					export let options = ***REMOVED***
						stages: [
							***REMOVED***target: 10***REMOVED***,
						],
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`)
				if assert.NoError(t, err) ***REMOVED***
					if assert.Len(t, b.Options.Stages, 1) ***REMOVED***
						assert.Equal(t, lib.Stage***REMOVED***Target: null.IntFrom(10)***REMOVED***, b.Options.Stages[0])
					***REMOVED***
				***REMOVED***
			***REMOVED***)
			t.Run("Duration", func(t *testing.T) ***REMOVED***
				b, err := getSimpleBundle("/script.js", `
					export let options = ***REMOVED***
						stages: [
							***REMOVED***duration: "10s"***REMOVED***,
						],
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`)
				if assert.NoError(t, err) ***REMOVED***
					if assert.Len(t, b.Options.Stages, 1) ***REMOVED***
						assert.Equal(t, lib.Stage***REMOVED***Duration: types.NullDurationFrom(10 * time.Second)***REMOVED***, b.Options.Stages[0])
					***REMOVED***
				***REMOVED***
			***REMOVED***)
			t.Run("DurationAndTarget", func(t *testing.T) ***REMOVED***
				b, err := getSimpleBundle("/script.js", `
					export let options = ***REMOVED***
						stages: [
							***REMOVED***duration: "10s", target: 10***REMOVED***,
						],
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`)
				if assert.NoError(t, err) ***REMOVED***
					if assert.Len(t, b.Options.Stages, 1) ***REMOVED***
						assert.Equal(t, lib.Stage***REMOVED***Duration: types.NullDurationFrom(10 * time.Second), Target: null.IntFrom(10)***REMOVED***, b.Options.Stages[0])
					***REMOVED***
				***REMOVED***
			***REMOVED***)
			t.Run("RampUpAndPlateau", func(t *testing.T) ***REMOVED***
				b, err := getSimpleBundle("/script.js", `
					export let options = ***REMOVED***
						stages: [
							***REMOVED***duration: "10s", target: 10***REMOVED***,
							***REMOVED***duration: "5s"***REMOVED***,
						],
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`)
				if assert.NoError(t, err) ***REMOVED***
					if assert.Len(t, b.Options.Stages, 2) ***REMOVED***
						assert.Equal(t, lib.Stage***REMOVED***Duration: types.NullDurationFrom(10 * time.Second), Target: null.IntFrom(10)***REMOVED***, b.Options.Stages[0])
						assert.Equal(t, lib.Stage***REMOVED***Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***, b.Options.Stages[1])
					***REMOVED***
				***REMOVED***
			***REMOVED***)
		***REMOVED***)
		t.Run("MaxRedirects", func(t *testing.T) ***REMOVED***
			b, err := getSimpleBundle("/script.js", `
				export let options = ***REMOVED***
					maxRedirects: 10,
				***REMOVED***;
				export default function() ***REMOVED******REMOVED***;
			`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.IntFrom(10), b.Options.MaxRedirects)
			***REMOVED***
		***REMOVED***)
		t.Run("InsecureSkipTLSVerify", func(t *testing.T) ***REMOVED***
			b, err := getSimpleBundle("/script.js", `
				export let options = ***REMOVED***
					insecureSkipTLSVerify: true,
				***REMOVED***;
				export default function() ***REMOVED******REMOVED***;
			`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.BoolFrom(true), b.Options.InsecureSkipTLSVerify)
			***REMOVED***
		***REMOVED***)
		t.Run("TLSCipherSuites", func(t *testing.T) ***REMOVED***
			for suiteName, suiteID := range lib.SupportedTLSCipherSuites ***REMOVED***
				t.Run(suiteName, func(t *testing.T) ***REMOVED***
					script := `
					export let options = ***REMOVED***
						tlsCipherSuites: ["%s"]
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
					`
					script = fmt.Sprintf(script, suiteName)

					b, err := getSimpleBundle("/script.js", script)
					if assert.NoError(t, err) ***REMOVED***
						if assert.Len(t, *b.Options.TLSCipherSuites, 1) ***REMOVED***
							assert.Equal(t, (*b.Options.TLSCipherSuites)[0], suiteID)
						***REMOVED***
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
		t.Run("TLSVersion", func(t *testing.T) ***REMOVED***
			t.Run("Object", func(t *testing.T) ***REMOVED***
				b, err := getSimpleBundle("/script.js", `
					export let options = ***REMOVED***
						tlsVersion: ***REMOVED***
							min: "ssl3.0",
							max: "tls1.2"
						***REMOVED***
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`)
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, b.Options.TLSVersion.Min, lib.TLSVersion(tls.VersionSSL30))
					assert.Equal(t, b.Options.TLSVersion.Max, lib.TLSVersion(tls.VersionTLS12))
				***REMOVED***
			***REMOVED***)
			t.Run("String", func(t *testing.T) ***REMOVED***
				b, err := getSimpleBundle("/script.js", `
					export let options = ***REMOVED***
						tlsVersion: "ssl3.0"
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`)
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, b.Options.TLSVersion.Min, lib.TLSVersion(tls.VersionSSL30))
					assert.Equal(t, b.Options.TLSVersion.Max, lib.TLSVersion(tls.VersionSSL30))
				***REMOVED***

			***REMOVED***)
		***REMOVED***)
		t.Run("Thresholds", func(t *testing.T) ***REMOVED***
			b, err := getSimpleBundle("/script.js", `
				export let options = ***REMOVED***
					thresholds: ***REMOVED***
						http_req_duration: ["avg<100"],
					***REMOVED***,
				***REMOVED***;
				export default function() ***REMOVED******REMOVED***;
			`)
			if assert.NoError(t, err) ***REMOVED***
				if assert.Len(t, b.Options.Thresholds["http_req_duration"].Thresholds, 1) ***REMOVED***
					assert.Equal(t, "avg<100", b.Options.Thresholds["http_req_duration"].Thresholds[0].Source)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestNewBundleFromArchive(t *testing.T) ***REMOVED***
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/path/to", 0755))
	assert.NoError(t, afero.WriteFile(fs, "/path/to/file.txt", []byte(`hi`), 0644))
	assert.NoError(t, afero.WriteFile(fs, "/path/to/exclaim.js", []byte(`export default function(s) ***REMOVED*** return s + "!" ***REMOVED***;`), 0644))

	data := `
			import exclaim from "./exclaim.js";
			export let options = ***REMOVED*** vus: 12345 ***REMOVED***;
			export let file = open("./file.txt");
			export default function() ***REMOVED*** return exclaim(file); ***REMOVED***;
		`
	b, err := getSimpleBundleWithFs("/path/to/script.js", data, fs)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	assert.Equal(t, lib.Options***REMOVED***VUs: null.IntFrom(12345)***REMOVED***, b.Options)

	bi, err := b.Instantiate()
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	v, err := bi.Default(goja.Undefined())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	assert.Equal(t, "hi!", v.Export())

	arc := b.makeArchive()
	assert.Equal(t, "js", arc.Type)
	assert.Equal(t, lib.Options***REMOVED***VUs: null.IntFrom(12345)***REMOVED***, arc.Options)
	assert.Equal(t, "file:///path/to/script.js", arc.FilenameURL.String())
	assert.Equal(t, data, string(arc.Data))
	assert.Equal(t, "file:///path/to/", arc.PwdURL.String())

	exclaimData, err := afero.ReadFile(arc.Filesystems["file"], "/path/to/exclaim.js")
	assert.NoError(t, err)
	assert.Equal(t, `export default function(s) ***REMOVED*** return s + "!" ***REMOVED***;`, string(exclaimData))

	fileData, err := afero.ReadFile(arc.Filesystems["file"], "/path/to/file.txt")
	assert.NoError(t, err)
	assert.Equal(t, `hi`, string(fileData))
	assert.Equal(t, consts.Version, arc.K6Version)

	b2, err := NewBundleFromArchive(arc, lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	assert.Equal(t, lib.Options***REMOVED***VUs: null.IntFrom(12345)***REMOVED***, b2.Options)

	bi2, err := b.Instantiate()
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	v2, err := bi2.Default(goja.Undefined())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	assert.Equal(t, "hi!", v2.Export())
***REMOVED***

func TestOpen(t *testing.T) ***REMOVED***
	var testCases = [...]struct ***REMOVED***
		name           string
		openPath       string
		pwd            string
		isError        bool
		isArchiveError bool
	***REMOVED******REMOVED***
		***REMOVED***
			name:     "notOpeningUrls",
			openPath: "github.com",
			isError:  true,
			pwd:      "/path/to",
		***REMOVED***,
		***REMOVED***
			name:     "simple",
			openPath: "file.txt",
			pwd:      "/path/to",
		***REMOVED***,
		***REMOVED***
			name:     "simple with dot",
			openPath: "./file.txt",
			pwd:      "/path/to",
		***REMOVED***,
		***REMOVED***
			name:     "simple with two dots",
			openPath: "../to/file.txt",
			pwd:      "/path/not",
		***REMOVED***,
		***REMOVED***
			name:     "fullpath",
			openPath: "/path/to/file.txt",
			pwd:      "/path/to",
		***REMOVED***,
		***REMOVED***
			name:     "fullpath2",
			openPath: "/path/to/file.txt",
			pwd:      "/path",
		***REMOVED***,
		***REMOVED***
			name:     "file is dir",
			openPath: "/path/to/",
			pwd:      "/path/to",
			isError:  true,
		***REMOVED***,
		***REMOVED***
			name:     "file is missing",
			openPath: "/path/to/missing.txt",
			isError:  true,
		***REMOVED***,
		***REMOVED***
			name:     "relative1",
			openPath: "to/file.txt",
			pwd:      "/path",
		***REMOVED***,
		***REMOVED***
			name:     "relative2",
			openPath: "./path/to/file.txt",
			pwd:      "/",
		***REMOVED***,
		***REMOVED***
			name:     "relative wonky",
			openPath: "../path/to/file.txt",
			pwd:      "/path",
		***REMOVED***,
		***REMOVED***
			name:     "empty open doesn't panic",
			openPath: "",
			pwd:      "/path",
			isError:  true,
		***REMOVED***,
	***REMOVED***
	fss := map[string]func() (afero.Fs, string, func())***REMOVED***
		"MemMapFS": func() (afero.Fs, string, func()) ***REMOVED***
			fs := afero.NewMemMapFs()
			require.NoError(t, fs.MkdirAll("/path/to", 0755))
			require.NoError(t, afero.WriteFile(fs, "/path/to/file.txt", []byte(`hi`), 0644))
			return fs, "", func() ***REMOVED******REMOVED***
		***REMOVED***,
		"OsFS": func() (afero.Fs, string, func()) ***REMOVED***
			prefix, err := ioutil.TempDir("", "k6_open_test")
			require.NoError(t, err)
			fs := afero.NewOsFs()
			filePath := filepath.Join(prefix, "/path/to/file.txt")
			require.NoError(t, fs.MkdirAll(filepath.Join(prefix, "/path/to"), 0755))
			require.NoError(t, afero.WriteFile(fs, filePath, []byte(`hi`), 0644))
			if isWindows ***REMOVED***
				fs = fsext.NewTrimFilePathSeparatorFs(fs)
			***REMOVED***
			return fs, prefix, func() ***REMOVED*** require.NoError(t, os.RemoveAll(prefix)) ***REMOVED***
		***REMOVED***,
	***REMOVED***

	for name, fsInit := range fss ***REMOVED***
		fs, prefix, cleanUp := fsInit()
		defer cleanUp()
		fs = afero.NewReadOnlyFs(fs)
		t.Run(name, func(t *testing.T) ***REMOVED***
			for _, tCase := range testCases ***REMOVED***
				tCase := tCase

				var testFunc = func(t *testing.T) ***REMOVED***
					var openPath = tCase.openPath
					// if fullpath prepend prefix
					if openPath != "" && (openPath[0] == '/' || openPath[0] == '\\') ***REMOVED***
						openPath = filepath.Join(prefix, openPath)
					***REMOVED***
					if isWindows ***REMOVED***
						openPath = strings.Replace(openPath, `\`, `\\`, -1)
					***REMOVED***
					var pwd = tCase.pwd
					if pwd == "" ***REMOVED***
						pwd = "/path/to/"
					***REMOVED***
					data := `
						export let file = open("` + openPath + `");
						export default function() ***REMOVED*** return file ***REMOVED***;`

					sourceBundle, err := getSimpleBundleWithFs(filepath.ToSlash(filepath.Join(prefix, pwd, "script.js")), data, fs)
					if tCase.isError ***REMOVED***
						assert.Error(t, err)
						return
					***REMOVED***
					require.NoError(t, err)

					arcBundle, err := NewBundleFromArchive(sourceBundle.makeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)

					require.NoError(t, err)

					for source, b := range map[string]*Bundle***REMOVED***"source": sourceBundle, "archive": arcBundle***REMOVED*** ***REMOVED***
						b := b
						t.Run(source, func(t *testing.T) ***REMOVED***
							bi, err := b.Instantiate()
							require.NoError(t, err)
							v, err := bi.Default(goja.Undefined())
							require.NoError(t, err)
							assert.Equal(t, "hi", v.Export())
						***REMOVED***)
					***REMOVED***
				***REMOVED***

				t.Run(tCase.name, testFunc)
				if isWindows ***REMOVED***
					// windowsify the testcase
					tCase.openPath = strings.Replace(tCase.openPath, `/`, `\`, -1)
					tCase.pwd = strings.Replace(tCase.pwd, `/`, `\`, -1)
					t.Run(tCase.name+" with windows slash", testFunc)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestBundleInstantiate(t *testing.T) ***REMOVED***
	b, err := getSimpleBundle("/script.js", `
		export let options = ***REMOVED***
			vus: 5,
			teardownTimeout: '1s',
		***REMOVED***;
		let val = true;
		export default function() ***REMOVED*** return val; ***REMOVED***
	`)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	bi, err := b.Instantiate()
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	t.Run("Run", func(t *testing.T) ***REMOVED***
		v, err := bi.Default(goja.Undefined())
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, true, v.Export())
		***REMOVED***
	***REMOVED***)

	t.Run("SetAndRun", func(t *testing.T) ***REMOVED***
		bi.Runtime.Set("val", false)
		v, err := bi.Default(goja.Undefined())
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, false, v.Export())
		***REMOVED***
	***REMOVED***)

	t.Run("Options", func(t *testing.T) ***REMOVED***
		// Ensure `options` properties are correctly marshalled
		jsOptions := bi.Runtime.Get("options").ToObject(bi.Runtime)
		vus := jsOptions.Get("vus").Export()
		assert.Equal(t, int64(5), vus)
		tdt := jsOptions.Get("teardownTimeout").Export()
		assert.Equal(t, "1s", tdt)

		// Ensure options propagate correctly from outside to the script
		optOrig := b.Options.VUs
		b.Options.VUs = null.IntFrom(10)
		bi2, err := b.Instantiate()
		assert.NoError(t, err)
		jsOptions = bi2.Runtime.Get("options").ToObject(bi2.Runtime)
		vus = jsOptions.Get("vus").Export()
		assert.Equal(t, int64(10), vus)
		b.Options.VUs = optOrig
	***REMOVED***)
***REMOVED***

func TestBundleEnv(t *testing.T) ***REMOVED***
	rtOpts := lib.RuntimeOptions***REMOVED***Env: map[string]string***REMOVED***
		"TEST_A": "1",
		"TEST_B": "",
	***REMOVED******REMOVED***
	data := `
		export default function() ***REMOVED***
			if (__ENV.TEST_A !== "1") ***REMOVED*** throw new Error("Invalid TEST_A: " + __ENV.TEST_A); ***REMOVED***
			if (__ENV.TEST_B !== "") ***REMOVED*** throw new Error("Invalid TEST_B: " + __ENV.TEST_B); ***REMOVED***
		***REMOVED***
	`
	b1, err := getSimpleBundleWithOptions("/script.js", data, rtOpts)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	b2, err := NewBundleFromArchive(b1.makeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	bundles := map[string]*Bundle***REMOVED***"Source": b1, "Archive": b2***REMOVED***
	for name, b := range bundles ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			assert.Equal(t, "1", b.Env["TEST_A"])
			assert.Equal(t, "", b.Env["TEST_B"])

			bi, err := b.Instantiate()
			if assert.NoError(t, err) ***REMOVED***
				_, err := bi.Default(goja.Undefined())
				assert.NoError(t, err)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
