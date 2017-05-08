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

package main

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func Test_getSrcData(t *testing.T) ***REMOVED***
	t.Run("Files", func(t *testing.T) ***REMOVED***
		fs := afero.NewMemMapFs()
		assert.NoError(t, fs.MkdirAll("/path/to", 0755))
		assert.NoError(t, afero.WriteFile(fs, "/path/to/file.js", []byte(`hi!`), 0644))

		testdata := map[string]struct***REMOVED*** filename, pwd string ***REMOVED******REMOVED***
			"Absolute":        ***REMOVED***"/path/to/file.js", "/path"***REMOVED***,
			"Relative":        ***REMOVED***"./to/file.js", "/path"***REMOVED***,
			"Adjacent":        ***REMOVED***"./file.js", "/path/to"***REMOVED***,
			"ImpliedRelative": ***REMOVED***"to/file.js", "/path"***REMOVED***,
			"ImpliedAdjacent": ***REMOVED***"file.js", "/path/to"***REMOVED***,
		***REMOVED***
		for name, data := range testdata ***REMOVED***
			t.Run(name, func(t *testing.T) ***REMOVED***
				src, err := getSrcData(data.filename, data.pwd, nil, fs)
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, "/path/to/file.js", src.Filename)
					assert.Equal(t, "hi!", string(src.Data))
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***
