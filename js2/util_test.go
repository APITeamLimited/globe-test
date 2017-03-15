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
	"testing"

	"github.com/loadimpact/k6/lib"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestQuickInstance(t *testing.T) ***REMOVED***
	rt, err := QuickInstance(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data:     []byte(`export default function() ***REMOVED******REMOVED***`),
	***REMOVED***, afero.NewMemMapFs())
	assert.NoError(t, err)
	assert.NotNil(t, rt)

	t.Run("invalid", func(t *testing.T) ***REMOVED***
		_, err := QuickInstance(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data:     []byte(``),
		***REMOVED***, afero.NewMemMapFs())
		assert.Error(t, err)
	***REMOVED***)
***REMOVED***
