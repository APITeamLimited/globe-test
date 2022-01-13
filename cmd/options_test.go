/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2018 Load Impact
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

package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTagKeyValue(t *testing.T) ***REMOVED***
	t.Parallel()
	testData := []struct ***REMOVED***
		input string
		name  string
		value string
		err   error
	***REMOVED******REMOVED***
		***REMOVED***
			"",
			"",
			"",
			errTagEmptyString,
		***REMOVED***,
		***REMOVED***
			"=",
			"",
			"",
			errTagEmptyName,
		***REMOVED***,
		***REMOVED***
			"=test",
			"",
			"",
			errTagEmptyName,
		***REMOVED***,
		***REMOVED***
			"test",
			"",
			"",
			errTagEmptyValue,
		***REMOVED***,
		***REMOVED***
			"test=",
			"",
			"",
			errTagEmptyValue,
		***REMOVED***,
		***REMOVED***
			"myTag=foo",
			"myTag",
			"foo",
			nil,
		***REMOVED***,
	***REMOVED***

	for _, data := range testData ***REMOVED***
		data := data
		t.Run(data.input, func(t *testing.T) ***REMOVED***
			t.Parallel()
			name, value, err := parseTagNameValue(data.input)
			assert.Equal(t, name, data.name)
			assert.Equal(t, value, data.value)
			assert.Equal(t, err, data.err)
		***REMOVED***)
	***REMOVED***
***REMOVED***
