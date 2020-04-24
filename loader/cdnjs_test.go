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

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCDNJS(t *testing.T) ***REMOVED***
	t.Skip("skipped to avoid inconsistent API responses")

	paths := map[string]struct ***REMOVED***
		parts []string
		src   string
	***REMOVED******REMOVED***
		"cdnjs.com/libraries/Faker": ***REMOVED***
			[]string***REMOVED***"Faker", "", ""***REMOVED***,
			`^https://cdnjs.cloudflare.com/ajax/libs/Faker/[\d\.]+/faker.min.js$`,
		***REMOVED***,
		"cdnjs.com/libraries/Faker/faker.js": ***REMOVED***
			[]string***REMOVED***"Faker", "", "faker.js"***REMOVED***,
			`^https://cdnjs.cloudflare.com/ajax/libs/Faker/[\d\.]+/faker.js$`,
		***REMOVED***,
		"cdnjs.com/libraries/Faker/locales/en_AU/faker.en_AU.min.js": ***REMOVED***
			[]string***REMOVED***"Faker", "", "locales/en_AU/faker.en_AU.min.js"***REMOVED***,
			`^https://cdnjs.cloudflare.com/ajax/libs/Faker/[\d\.]+/locales/en_AU/faker.en_AU.min.js$`,
		***REMOVED***,
		"cdnjs.com/libraries/Faker/3.1.0": ***REMOVED***
			[]string***REMOVED***"Faker", "3.1.0", ""***REMOVED***,
			`^https://cdnjs.cloudflare.com/ajax/libs/Faker/3.1.0/faker.min.js$`,
		***REMOVED***,
		"cdnjs.com/libraries/Faker/3.1.0/faker.js": ***REMOVED***
			[]string***REMOVED***"Faker", "3.1.0", "faker.js"***REMOVED***,
			`^https://cdnjs.cloudflare.com/ajax/libs/Faker/3.1.0/faker.js$`,
		***REMOVED***,
		"cdnjs.com/libraries/Faker/3.1.0/locales/en_AU/faker.en_AU.min.js": ***REMOVED***
			[]string***REMOVED***"Faker", "3.1.0", "locales/en_AU/faker.en_AU.min.js"***REMOVED***,
			`^https://cdnjs.cloudflare.com/ajax/libs/Faker/3.1.0/locales/en_AU/faker.en_AU.min.js$`,
		***REMOVED***,
		"cdnjs.com/libraries/Faker/0.7.2": ***REMOVED***
			[]string***REMOVED***"Faker", "0.7.2", ""***REMOVED***,
			`^https://cdnjs.cloudflare.com/ajax/libs/Faker/0.7.2/MinFaker.js$`,
		***REMOVED***,
	***REMOVED***

	var root = &url.URL***REMOVED***Scheme: "https", Host: "example.com", Path: "/something/"***REMOVED***
	for path, expected := range paths ***REMOVED***
		path, expected := path, expected
		t.Run(path, func(t *testing.T) ***REMOVED***
			name, loader, parts := pickLoader(path)
			assert.Equal(t, "cdnjs", name)
			assert.Equal(t, expected.parts, parts)

			src, err := loader(path, parts)
			require.NoError(t, err)
			assert.Regexp(t, expected.src, src)

			resolvedURL, err := Resolve(root, path)
			require.NoError(t, err)
			require.Empty(t, resolvedURL.Scheme)
			require.Equal(t, path, resolvedURL.Opaque)

			data, err := Load(map[string]afero.Fs***REMOVED***"https": afero.NewMemMapFs()***REMOVED***, resolvedURL, path)
			require.NoError(t, err)
			assert.Equal(t, resolvedURL, data.URL)
			assert.NotEmpty(t, data.Data)
		***REMOVED***)
	***REMOVED***

	t.Run("cdnjs.com/libraries/nonexistent", func(t *testing.T) ***REMOVED***
		path := "cdnjs.com/libraries/nonexistent"
		name, loader, parts := pickLoader(path)
		assert.Equal(t, "cdnjs", name)
		assert.Equal(t, []string***REMOVED***"nonexistent", "", ""***REMOVED***, parts)
		_, err := loader(path, parts)
		assert.EqualError(t, err, "cdnjs: no such library: nonexistent")
	***REMOVED***)

	t.Run("cdnjs.com/libraries/Faker/3.1.0/nonexistent.js", func(t *testing.T) ***REMOVED***
		path := "cdnjs.com/libraries/Faker/3.1.0/nonexistent.js"
		name, loader, parts := pickLoader(path)
		assert.Equal(t, "cdnjs", name)
		assert.Equal(t, []string***REMOVED***"Faker", "3.1.0", "nonexistent.js"***REMOVED***, parts)
		src, err := loader(path, parts)
		require.NoError(t, err)
		assert.Equal(t, "https://cdnjs.cloudflare.com/ajax/libs/Faker/3.1.0/nonexistent.js", src)

		pathURL, err := url.Parse(src)
		require.NoError(t, err)

		_, err = Load(map[string]afero.Fs***REMOVED***"https": afero.NewMemMapFs()***REMOVED***, pathURL, path)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found: https://cdnjs.cloudflare.com/ajax/libs/Faker/3.1.0/nonexistent.js")
	***REMOVED***)
***REMOVED***
