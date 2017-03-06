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
	"github.com/loadimpact/k6/lib"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
	"testing"
	"time"
)

func TestParseStage(t *testing.T) ***REMOVED***
	testdata := map[string]lib.Stage***REMOVED***
		"":        ***REMOVED******REMOVED***,
		":":       ***REMOVED******REMOVED***,
		"10s":     ***REMOVED***Duration: 10 * time.Second***REMOVED***,
		"10s:":    ***REMOVED***Duration: 10 * time.Second***REMOVED***,
		"10s:100": ***REMOVED***Duration: 10 * time.Second, Target: null.IntFrom(100)***REMOVED***,
		":100":    ***REMOVED***Target: null.IntFrom(100)***REMOVED***,
	***REMOVED***
	for s, st := range testdata ***REMOVED***
		t.Run(s, func(t *testing.T) ***REMOVED***
			parsed, err := ParseStage(s)
			assert.NoError(t, err)
			assert.Equal(t, st, parsed)
		***REMOVED***)
	***REMOVED***
***REMOVED***
