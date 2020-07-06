/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
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

package common

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/stats"
)

func TestInitWithoutAddressErrors(t *testing.T) ***REMOVED***
	var c = &Collector***REMOVED***
		Config: Config***REMOVED******REMOVED***,
		Type:   "testtype",
	***REMOVED***
	err := c.Init()
	require.Error(t, err)
***REMOVED***

func TestInitWithBogusAddressErrors(t *testing.T) ***REMOVED***
	var c = &Collector***REMOVED***
		Config: Config***REMOVED***
			Addr: null.StringFrom("localhost:90000"),
		***REMOVED***,
		Type: "testtype",
	***REMOVED***
	err := c.Init()
	require.Error(t, err)
***REMOVED***

func TestLinkReturnAddress(t *testing.T) ***REMOVED***
	var bogusValue = "bogus value"
	var c = &Collector***REMOVED***
		Config: Config***REMOVED***
			Addr: null.StringFrom(bogusValue),
		***REMOVED***,
	***REMOVED***
	require.Equal(t, bogusValue, c.Link())
***REMOVED***

func TestGetRequiredSystemTags(t *testing.T) ***REMOVED***
	var c = &Collector***REMOVED******REMOVED***
	require.Equal(t, stats.SystemTagSet(0), c.GetRequiredSystemTags())
***REMOVED***
