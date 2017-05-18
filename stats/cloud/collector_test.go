/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2017 Load Impact
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

package cloud

import (
	"os"
	"testing"

	"github.com/loadimpact/k6/lib"
	"github.com/stretchr/testify/assert"
)

func TestGetName(t *testing.T) ***REMOVED***
	nameTests := []struct ***REMOVED***
		lib      *lib.SourceData
		conf     loadimpactConfig
		expected string
	***REMOVED******REMOVED***
		***REMOVED***&lib.SourceData***REMOVED***Filename: ""***REMOVED***, loadimpactConfig***REMOVED******REMOVED***, TestName***REMOVED***,
		***REMOVED***&lib.SourceData***REMOVED***Filename: "-"***REMOVED***, loadimpactConfig***REMOVED******REMOVED***, TestName***REMOVED***,
		***REMOVED***&lib.SourceData***REMOVED***Filename: "script.js"***REMOVED***, loadimpactConfig***REMOVED******REMOVED***, "script.js"***REMOVED***,
		***REMOVED***&lib.SourceData***REMOVED***Filename: "/file/name.js"***REMOVED***, loadimpactConfig***REMOVED******REMOVED***, "name.js"***REMOVED***,
		***REMOVED***&lib.SourceData***REMOVED***Filename: "/file/name"***REMOVED***, loadimpactConfig***REMOVED******REMOVED***, "name"***REMOVED***,
		***REMOVED***&lib.SourceData***REMOVED***Filename: "/file/name"***REMOVED***, loadimpactConfig***REMOVED***Name: "confName"***REMOVED***, "confName"***REMOVED***,
	***REMOVED***

	for _, test := range nameTests ***REMOVED***
		actual := getName(test.lib, test.conf)
		assert.Equal(t, actual, test.expected)
	***REMOVED***

	err := os.Setenv("K6CLOUD_NAME", "envname")
	assert.Nil(t, err)

	for _, test := range nameTests ***REMOVED***
		actual := getName(test.lib, test.conf)
		assert.Equal(t, actual, "envname")
	***REMOVED***
***REMOVED***
