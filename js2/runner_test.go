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

package js2

import (
	"context"
	"testing"

	"github.com/loadimpact/k6/lib"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
)

func TestRunnerNew(t *testing.T) ***REMOVED***
	r, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data: []byte(`
			let counter = 0;
			export default function() ***REMOVED*** counter++; ***REMOVED***
		`),
	***REMOVED***, afero.NewMemMapFs())
	assert.NoError(t, err)

	t.Run("NewVU", func(t *testing.T) ***REMOVED***
		vu, err := r.newVU()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), vu.Runtime.Get("counter").Export())

		t.Run("RunOnce", func(t *testing.T) ***REMOVED***
			_, err = vu.RunOnce(context.Background())
			assert.NoError(t, err)
			assert.Equal(t, int64(1), vu.Runtime.Get("counter").Export())
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestRunnerGetDefaultGroup(t *testing.T) ***REMOVED***
	r, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data:     []byte(`export default function() ***REMOVED******REMOVED***;`),
	***REMOVED***, afero.NewMemMapFs())
	assert.NoError(t, err)
	assert.NotNil(t, r.GetDefaultGroup())
***REMOVED***

func TestRunnerOptions(t *testing.T) ***REMOVED***
	r, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data:     []byte(`export default function() ***REMOVED******REMOVED***;`),
	***REMOVED***, afero.NewMemMapFs())
	assert.NoError(t, err)

	assert.Equal(t, r.Bundle.Options, r.GetOptions())
	assert.Equal(t, null.NewBool(false, false), r.Bundle.Options.Paused)
	r.ApplyOptions(lib.Options***REMOVED***Paused: null.BoolFrom(true)***REMOVED***)
	assert.Equal(t, r.Bundle.Options, r.GetOptions())
	assert.Equal(t, null.NewBool(true, true), r.Bundle.Options.Paused)
	r.ApplyOptions(lib.Options***REMOVED***Paused: null.BoolFrom(false)***REMOVED***)
	assert.Equal(t, r.Bundle.Options, r.GetOptions())
	assert.Equal(t, null.NewBool(false, true), r.Bundle.Options.Paused)
***REMOVED***
