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

package loader

import (
	"bytes"
	"io"
	"net/url"
	"testing"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/loadimpact/k6/lib/fsext"
)

type errorReader string

func (e errorReader) Read(_ []byte) (int, error) ***REMOVED***
	return 0, errors.New((string)(e))
***REMOVED***

var _ io.Reader = errorReader("")

func TestReadSourceSTDINError(t *testing.T) ***REMOVED***
	_, err := ReadSource("-", "", nil, errorReader("1234"))
	require.Error(t, err)
	require.Equal(t, "1234", err.Error())
***REMOVED***

func TestReadSourceSTDINCache(t *testing.T) ***REMOVED***
	var data = []byte(`test contents`)
	var r = bytes.NewReader(data)
	var fs = afero.NewMemMapFs()
	sourceData, err := ReadSource("-", "/path/to/pwd",
		map[string]afero.Fs***REMOVED***"file": fsext.NewCacheOnReadFs(nil, fs, 0)***REMOVED***, r)
	require.NoError(t, err)
	require.Equal(t, &SourceData***REMOVED***
		URL:  &url.URL***REMOVED***Scheme: "file", Path: "/-"***REMOVED***,
		Data: data***REMOVED***, sourceData)
	fileData, err := afero.ReadFile(fs, "/-")
	require.NoError(t, err)
	require.Equal(t, data, fileData)
***REMOVED***

func TestReadSourceRelative(t *testing.T) ***REMOVED***
	var data = []byte(`test contents`)
	var fs = afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "/path/to/somewhere/script.js", data, 0644))
	sourceData, err := ReadSource("../somewhere/script.js", "/path/to/pwd", map[string]afero.Fs***REMOVED***"file": fs***REMOVED***, nil)
	require.NoError(t, err)
	require.Equal(t, &SourceData***REMOVED***
		URL:  &url.URL***REMOVED***Scheme: "file", Path: "/path/to/somewhere/script.js"***REMOVED***,
		Data: data***REMOVED***, sourceData)
***REMOVED***

func TestReadSourceAbsolute(t *testing.T) ***REMOVED***
	var data = []byte(`test contents`)
	var r = bytes.NewReader(data)
	var fs = afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "/a/b", data, 0644))
	require.NoError(t, afero.WriteFile(fs, "/c/a/b", []byte("wrong"), 0644))
	sourceData, err := ReadSource("/a/b", "/c", map[string]afero.Fs***REMOVED***"file": fs***REMOVED***, r)
	require.NoError(t, err)
	require.Equal(t, &SourceData***REMOVED***
		URL:  &url.URL***REMOVED***Scheme: "file", Path: "/a/b"***REMOVED***,
		Data: data***REMOVED***, sourceData)
***REMOVED***

func TestReadSourceHttps(t *testing.T) ***REMOVED***
	var data = []byte(`test contents`)
	var fs = afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "/github.com/something", data, 0644))
	sourceData, err := ReadSource("https://github.com/something", "/c",
		map[string]afero.Fs***REMOVED***"file": afero.NewMemMapFs(), "https": fs***REMOVED***, nil)
	require.NoError(t, err)
	require.Equal(t, &SourceData***REMOVED***
		URL:  &url.URL***REMOVED***Scheme: "https", Host: "github.com", Path: "/something"***REMOVED***,
		Data: data***REMOVED***, sourceData)
***REMOVED***

func TestReadSourceHttpError(t *testing.T) ***REMOVED***
	var data = []byte(`test contents`)
	var fs = afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "/github.com/something", data, 0644))
	_, err := ReadSource("http://github.com/something", "/c",
		map[string]afero.Fs***REMOVED***"file": afero.NewMemMapFs(), "https": fs***REMOVED***, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), `only supported schemes for imports are file and https`)
***REMOVED***
