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

package influxdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigText(t *testing.T) ***REMOVED***
	testdata := map[string]Config***REMOVED***
		"":                             ***REMOVED******REMOVED***,
		"dbname":                       ***REMOVED***DB: "dbname"***REMOVED***,
		"/dbname":                      ***REMOVED***DB: "dbname"***REMOVED***,
		"http://localhost:8086":        ***REMOVED***Addr: "http://localhost:8086"***REMOVED***,
		"http://localhost:8086/dbname": ***REMOVED***Addr: "http://localhost:8086", DB: "dbname"***REMOVED***,
	***REMOVED***
	queries := map[string]struct ***REMOVED***
		Config Config
		Err    string
	***REMOVED******REMOVED***
		"?":                ***REMOVED***Config***REMOVED******REMOVED***, ""***REMOVED***,
		"?insecure=false":  ***REMOVED***Config***REMOVED***Insecure: false***REMOVED***, ""***REMOVED***,
		"?insecure=true":   ***REMOVED***Config***REMOVED***Insecure: true***REMOVED***, ""***REMOVED***,
		"?insecure=ture":   ***REMOVED***Config***REMOVED******REMOVED***, "insecure must be true or false, not ture"***REMOVED***,
		"?payload_size=69": ***REMOVED***Config***REMOVED***PayloadSize: 69***REMOVED***, ""***REMOVED***,
		"?payload_size=a":  ***REMOVED***Config***REMOVED******REMOVED***, "strconv.Atoi: parsing \"a\": invalid syntax"***REMOVED***,
	***REMOVED***
	for str, data := range testdata ***REMOVED***
		t.Run(str, func(t *testing.T) ***REMOVED***
			config, err := ParseURL(str)
			assert.NoError(t, err)
			assert.Equal(t, data, config)

			for q, qdata := range queries ***REMOVED***
				t.Run(q, func(t *testing.T) ***REMOVED***
					config, err := ParseURL(str + q)
					if qdata.Err != "" ***REMOVED***
						assert.EqualError(t, err, qdata.Err)
					***REMOVED*** else ***REMOVED***
						expected2 := qdata.Config
						expected2.DB = data.DB
						expected2.Addr = data.Addr
						assert.Equal(t, expected2, config)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
