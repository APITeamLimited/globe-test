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

package loader

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/loadimpact/k6/lib/testutils"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestDir(t *testing.T) ***REMOVED***
	testdata := map[string]string***REMOVED***
		"/path/to/file.txt": "/path/to",
		"-":                 "/",
	***REMOVED***
	for name, dir := range testdata ***REMOVED***
		t.Run("path="+name, func(t *testing.T) ***REMOVED***
			assert.Equal(t, dir, Dir(name))
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestLoad(t *testing.T) ***REMOVED***
	tb := testutils.NewHTTPMultiBin(t)
	sr := tb.Replacer.Replace

	oldHTTPTransport := http.DefaultTransport
	http.DefaultTransport = tb.HTTPTransport

	defer func() ***REMOVED***
		tb.Cleanup()
		http.DefaultTransport = oldHTTPTransport
	***REMOVED***()

	t.Run("Blank", func(t *testing.T) ***REMOVED***
		_, err := Load(nil, "/", "")
		assert.EqualError(t, err, "local or remote path required")
	***REMOVED***)

	t.Run("Protocol", func(t *testing.T) ***REMOVED***
		_, err := Load(nil, "/", sr("HTTPSBIN_URL/html"))
		assert.EqualError(t, err, "imports should not contain a protocol")
	***REMOVED***)

	t.Run("Local", func(t *testing.T) ***REMOVED***
		fs := afero.NewMemMapFs()
		assert.NoError(t, fs.MkdirAll("/path/to", 0755))
		assert.NoError(t, afero.WriteFile(fs, "/path/to/file.txt", []byte("hi"), 0644))

		testdata := map[string]struct***REMOVED*** pwd, path string ***REMOVED******REMOVED***
			"Absolute": ***REMOVED***"/path", "/path/to/file.txt"***REMOVED***,
			"Relative": ***REMOVED***"/path", "./to/file.txt"***REMOVED***,
			"Adjacent": ***REMOVED***"/path/to", "./file.txt"***REMOVED***,
		***REMOVED***
		for name, data := range testdata ***REMOVED***
			t.Run(name, func(t *testing.T) ***REMOVED***
				src, err := Load(fs, data.pwd, data.path)
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, "/path/to/file.txt", src.Filename)
					assert.Equal(t, "hi", string(src.Data))
				***REMOVED***
			***REMOVED***)
		***REMOVED***

		t.Run("Nonexistent", func(t *testing.T) ***REMOVED***
			_, err := Load(fs, "/", "/nonexistent")
			assert.EqualError(t, err, "open /nonexistent: file does not exist")
		***REMOVED***)

		t.Run("Remote Lifting Denied", func(t *testing.T) ***REMOVED***
			_, err := Load(fs, "example.com", "/etc/shadow")
			assert.EqualError(t, err, "origin (example.com) not allowed to load local file: /etc/shadow")
		***REMOVED***)
	***REMOVED***)

	t.Run("Remote", func(t *testing.T) ***REMOVED***
		src, err := Load(nil, "/", sr("HTTPSBIN_DOMAIN:HTTPSBIN_PORT/html"))
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, src.Filename, sr("HTTPSBIN_DOMAIN:HTTPSBIN_PORT/html"))
			assert.Contains(t, string(src.Data), "Herman Melville - Moby-Dick")
		***REMOVED***

		t.Run("Absolute", func(t *testing.T) ***REMOVED***
			src, err := Load(nil, sr("HTTPSBIN_DOMAIN:HTTPSBIN_PORT"), sr("HTTPSBIN_DOMAIN:HTTPSBIN_PORT/robots.txt"))
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, src.Filename, sr("HTTPSBIN_DOMAIN:HTTPSBIN_PORT/robots.txt"))
				assert.Equal(t, string(src.Data), "User-agent: *\nDisallow: /deny\n")
			***REMOVED***
		***REMOVED***)

		t.Run("Relative", func(t *testing.T) ***REMOVED***
			src, err := Load(nil, sr("HTTPSBIN_DOMAIN:HTTPSBIN_PORT"), "./robots.txt")
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, src.Filename, sr("HTTPSBIN_DOMAIN:HTTPSBIN_PORT/robots.txt"))
				assert.Equal(t, string(src.Data), "User-agent: *\nDisallow: /deny\n")
			***REMOVED***
		***REMOVED***)
	***REMOVED***)

	const responseStr = "export function fn() ***REMOVED***\r\n    return 1234;\r\n***REMOVED***"
	tb.Mux.HandleFunc("/raw/something", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if _, ok := r.URL.Query()["_k6"]; ok ***REMOVED***
			http.Error(w, "Internal server error", 500)
			return
		***REMOVED***
		_, err := fmt.Fprint(w, responseStr)
		assert.NoError(t, err)
	***REMOVED***)

	t.Run("No _k6=1 Fallback", func(t *testing.T) ***REMOVED***
		src, err := Load(nil, "/", sr("HTTPSBIN_DOMAIN:HTTPSBIN_PORT/raw/something"))
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, src.Filename, sr("HTTPSBIN_DOMAIN:HTTPSBIN_PORT/raw/something"))
			assert.Equal(t, responseStr, string(src.Data))
		***REMOVED***
	***REMOVED***)

	tb.Mux.HandleFunc("/invalid", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		http.Error(w, "Internal server error", 500)
	***REMOVED***)

	t.Run("Invalid", func(t *testing.T) ***REMOVED***
		src, err := Load(nil, "/", sr("HTTPSBIN_DOMAIN:HTTPSBIN_PORT/invalid"))
		assert.Nil(t, src)
		assert.Error(t, err)

		t.Run("Host", func(t *testing.T) ***REMOVED***
			src, err := Load(nil, "/", "some-path-that-doesnt-exist.js")
			assert.Nil(t, src)
			assert.Error(t, err)
		***REMOVED***)
		t.Run("URL", func(t *testing.T) ***REMOVED***
			src, err := Load(nil, "/", "192.168.0.%31")
			assert.Nil(t, src)
			assert.Error(t, err)
		***REMOVED***)
	***REMOVED***)
***REMOVED***
