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
	"github.com/loadimpact/k6/lib"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBundle(t *testing.T) ***REMOVED***
	t.Run("Blank", func(t *testing.T) ***REMOVED***
		_, err := NewBundle(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data:     []byte(``),
		***REMOVED***, afero.NewMemMapFs())
		assert.EqualError(t, err, "script must export a default function")
	***REMOVED***)
	t.Run("DefaultUndefined", func(t *testing.T) ***REMOVED***
		_, err := NewBundle(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data: []byte(`
				export default undefined;
			`),
		***REMOVED***, afero.NewMemMapFs())
		assert.EqualError(t, err, "script must export a default function")
	***REMOVED***)
	t.Run("DefaultNull", func(t *testing.T) ***REMOVED***
		_, err := NewBundle(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data: []byte(`
				export default null;
			`),
		***REMOVED***, afero.NewMemMapFs())
		assert.EqualError(t, err, "script must export a default function")
	***REMOVED***)
	t.Run("DefaultWrongType", func(t *testing.T) ***REMOVED***
		_, err := NewBundle(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data: []byte(`
				export default 12345;
			`),
		***REMOVED***, afero.NewMemMapFs())
		assert.EqualError(t, err, "default export must be a function")
	***REMOVED***)
	t.Run("Minimal", func(t *testing.T) ***REMOVED***
		_, err := NewBundle(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data: []byte(`
				export default function() ***REMOVED******REMOVED***;
			`),
		***REMOVED***, afero.NewMemMapFs())
		assert.NoError(t, err)
	***REMOVED***)
***REMOVED***
