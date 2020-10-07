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
	"net/url"
	"testing"

	"github.com/loadimpact/k6/lib/testutils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGithub(t *testing.T) ***REMOVED***
	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))
	path := "github.com/github/gitignore/Go.gitignore"
	expectedEndSrc := "https://raw.githubusercontent.com/github/gitignore/master/Go.gitignore"
	name, loader, parts := pickLoader(path)
	assert.Equal(t, "github", name)
	assert.Equal(t, []string***REMOVED***"github", "gitignore", "Go.gitignore"***REMOVED***, parts)
	src, err := loader(logger, path, parts)
	assert.NoError(t, err)
	assert.Equal(t, expectedEndSrc, src)

	var root = &url.URL***REMOVED***Scheme: "https", Host: "example.com", Path: "/something/"***REMOVED***
	resolvedURL, err := Resolve(root, path)
	require.NoError(t, err)
	require.Empty(t, resolvedURL.Scheme)
	require.Equal(t, path, resolvedURL.Opaque)
	t.Run("not cached", func(t *testing.T) ***REMOVED***
		data, err := Load(logger, map[string]afero.Fs***REMOVED***"https": afero.NewMemMapFs()***REMOVED***, resolvedURL, path)
		require.NoError(t, err)
		assert.Equal(t, data.URL, resolvedURL)
		assert.Equal(t, path, data.URL.String())
		assert.NotEmpty(t, data.Data)
	***REMOVED***)

	t.Run("cached", func(t *testing.T) ***REMOVED***
		fs := afero.NewMemMapFs()
		testData := []byte("test data")

		err := afero.WriteFile(fs, "/github.com/github/gitignore/Go.gitignore", testData, 0644)
		require.NoError(t, err)

		data, err := Load(logger, map[string]afero.Fs***REMOVED***"https": fs***REMOVED***, resolvedURL, path)
		require.NoError(t, err)
		assert.Equal(t, path, data.URL.String())
		assert.Equal(t, data.Data, testData)
	***REMOVED***)

	t.Run("relative", func(t *testing.T) ***REMOVED***
		var tests = map[string]string***REMOVED***
			"./something.else":  "github.com/github/gitignore/something.else",
			"../something.else": "github.com/github/something.else",
			"/something.else":   "github.com/something.else",
		***REMOVED***
		for relative, expected := range tests ***REMOVED***
			relativeURL, err := Resolve(Dir(resolvedURL), relative)
			require.NoError(t, err)
			assert.Equal(t, expected, relativeURL.String())
		***REMOVED***
	***REMOVED***)

	t.Run("dir", func(t *testing.T) ***REMOVED***
		require.Equal(t, &url.URL***REMOVED***Opaque: "github.com/github/gitignore"***REMOVED***, Dir(resolvedURL))
	***REMOVED***)
***REMOVED***
