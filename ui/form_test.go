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

package ui

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForm(t *testing.T) ***REMOVED***
	t.Parallel()
	t.Run("Blank", func(t *testing.T) ***REMOVED***
		t.Parallel()
		data, err := Form***REMOVED******REMOVED***.Run(strings.NewReader(""), bytes.NewBuffer(nil))
		assert.NoError(t, err)
		assert.Equal(t, map[string]string***REMOVED******REMOVED***, data)
	***REMOVED***)
	t.Run("Banner", func(t *testing.T) ***REMOVED***
		t.Parallel()
		out := bytes.NewBuffer(nil)
		data, err := Form***REMOVED***Banner: "Hi!"***REMOVED***.Run(strings.NewReader(""), out)
		assert.NoError(t, err)
		assert.Equal(t, map[string]string***REMOVED******REMOVED***, data)
		assert.Equal(t, "Hi!\n\n", out.String())
	***REMOVED***)
	t.Run("Field", func(t *testing.T) ***REMOVED***
		t.Parallel()
		f := Form***REMOVED***
			Fields: []Field***REMOVED***
				StringField***REMOVED***Key: "key", Label: "label"***REMOVED***,
			***REMOVED***,
		***REMOVED***
		in := "Value\n"
		out := bytes.NewBuffer(nil)
		data, err := f.Run(strings.NewReader(in), out)
		assert.NoError(t, err)
		assert.Equal(t, map[string]string***REMOVED***"key": "Value"***REMOVED***, data)
		assert.Equal(t, "  label: ", out.String())
	***REMOVED***)
	t.Run("Fields", func(t *testing.T) ***REMOVED***
		t.Parallel()
		f := Form***REMOVED***
			Fields: []Field***REMOVED***
				StringField***REMOVED***Key: "a", Label: "label a"***REMOVED***,
				StringField***REMOVED***Key: "b", Label: "label b"***REMOVED***,
			***REMOVED***,
		***REMOVED***
		in := "1\n2\n"
		out := bytes.NewBuffer(nil)
		data, err := f.Run(strings.NewReader(in), out)
		assert.NoError(t, err)
		assert.Equal(t, map[string]string***REMOVED***"a": "1", "b": "2"***REMOVED***, data)
		assert.Equal(t, "  label a:   label b: ", out.String())
	***REMOVED***)
	t.Run("Defaults", func(t *testing.T) ***REMOVED***
		t.Parallel()
		f := Form***REMOVED***
			Fields: []Field***REMOVED***
				StringField***REMOVED***Key: "a", Label: "label a", Default: "default a"***REMOVED***,
				StringField***REMOVED***Key: "b", Label: "label b", Default: "default b"***REMOVED***,
			***REMOVED***,
		***REMOVED***
		in := "\n2\n"
		out := bytes.NewBuffer(nil)
		data, err := f.Run(strings.NewReader(in), out)
		assert.NoError(t, err)
		assert.Equal(t, map[string]string***REMOVED***"a": "default a", "b": "2"***REMOVED***, data)
		assert.Equal(t, "  label a [default a]:   label b [default b]: ", out.String())
	***REMOVED***)
	t.Run("Errors", func(t *testing.T) ***REMOVED***
		t.Parallel()
		f := Form***REMOVED***
			Fields: []Field***REMOVED***
				StringField***REMOVED***Key: "key", Label: "label", Min: 6, Max: 10***REMOVED***,
			***REMOVED***,
		***REMOVED***
		in := "short\ntoo damn long\nperfect\n"
		out := bytes.NewBuffer(nil)
		data, err := f.Run(strings.NewReader(in), out)
		assert.NoError(t, err)
		assert.Equal(t, map[string]string***REMOVED***"key": "perfect"***REMOVED***, data)
		assert.Equal(t, "  label: - invalid input, min length is 6\n  label: - invalid input, max length is 10\n  label: ", out.String())
	***REMOVED***)
***REMOVED***
