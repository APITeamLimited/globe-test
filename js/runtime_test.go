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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		return
	***REMOVED***

	r, err := New()
	assert.NoError(t, err)

	t.Run("Polyfill: Symbol", func(t *testing.T) ***REMOVED***
		v, err := r.VM.Get("Symbol")
		assert.NoError(t, err)
		assert.False(t, v.IsUndefined())
	***REMOVED***)
***REMOVED***

func TestLoad(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		return
	***REMOVED***

	r, err := New()
	assert.NoError(t, err)
	assert.NoError(t, r.VM.Set("__initapi__", InitAPI***REMOVED***r: r***REMOVED***))

	t.Run("Importing Libraries", func(t *testing.T) ***REMOVED***
		_, err := r.load("test.js", []byte(`
			import "k6";
		`))
		assert.NoError(t, err)
		assert.Contains(t, r.lib, "k6.js")
	***REMOVED***)
***REMOVED***

func TestExtractOptions(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		return
	***REMOVED***

	r, err := New()
	assert.NoError(t, err)

	t.Run("nothing", func(t *testing.T) ***REMOVED***
		_, err := r.load("test.js", []byte(``))
		assert.NoError(t, err)
	***REMOVED***)

	t.Run("vus", func(t *testing.T) ***REMOVED***
		_, err := r.load("test.js", []byte(`
			export let options = ***REMOVED*** vus: 12345 ***REMOVED***;
		`))
		assert.NoError(t, err)

		assert.True(t, r.Options.VUs.Valid)
		assert.Equal(t, int64(12345), r.Options.VUs.Int64)
	***REMOVED***)
	t.Run("vusMax", func(t *testing.T) ***REMOVED***
		_, err := r.load("test.js", []byte(`
			export let options = ***REMOVED*** vusMax: 12345 ***REMOVED***;
		`))
		assert.NoError(t, err)

		assert.True(t, r.Options.VUsMax.Valid)
		assert.Equal(t, int64(12345), r.Options.VUsMax.Int64)
	***REMOVED***)
	t.Run("duration", func(t *testing.T) ***REMOVED***
		_, err := r.load("test.js", []byte(`
			export let options = ***REMOVED*** duration: "2m" ***REMOVED***;
		`))
		assert.NoError(t, err)

		assert.True(t, r.Options.Duration.Valid)
		assert.Equal(t, "2m", r.Options.Duration.String)
	***REMOVED***)
	t.Run("maxRedirects", func(t *testing.T) ***REMOVED***
		_, err := r.load("test.js", []byte(`
			export let options = ***REMOVED*** maxRedirects: 12345 ***REMOVED***;
		`))
		assert.NoError(t, err)

		assert.True(t, r.Options.MaxRedirects.Valid)
		assert.Equal(t, int64(12345), r.Options.MaxRedirects.Int64)
	***REMOVED***)
	t.Run("thresholds", func(t *testing.T) ***REMOVED***
		_, err := r.load("test.js", []byte(`
			export let options = ***REMOVED***
				thresholds: ***REMOVED***
					my_metric: ["value<=1000"],
				***REMOVED***
			***REMOVED***
		`))
		assert.NoError(t, err)

		assert.Contains(t, r.Options.Thresholds, "my_metric")
		if assert.Len(t, r.Options.Thresholds["my_metric"].Thresholds, 1) ***REMOVED***
			assert.Equal(t, "value<=1000", r.Options.Thresholds["my_metric"].Thresholds[0].Source)
		***REMOVED***
	***REMOVED***)
***REMOVED***
