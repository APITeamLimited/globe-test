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

	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"
)

func TestMakeBatchConfig(t *testing.T) ***REMOVED***
	t.Run("Empty", func(t *testing.T) ***REMOVED***
		assert.Equal(t,
			client.BatchPointsConfig***REMOVED***Database: "k6"***REMOVED***,
			MakeBatchConfig(Config***REMOVED******REMOVED***),
		)
	***REMOVED***)
	t.Run("DB Set", func(t *testing.T) ***REMOVED***
		assert.Equal(t,
			client.BatchPointsConfig***REMOVED***Database: "dbname"***REMOVED***,
			MakeBatchConfig(Config***REMOVED***DB: null.StringFrom("dbname")***REMOVED***),
		)
	***REMOVED***)
***REMOVED***

func TestFieldKinds(t *testing.T) ***REMOVED***
	var fieldKinds map[string]FieldKind
	var err error

	conf := NewConfig()
	conf.TagsAsFields = []string***REMOVED***"vu", "iter", "url", "boolField", "floatField", "intField"***REMOVED***

	// Error case 1 (duplicated bool fields)
	conf.TagsAsFields = []string***REMOVED***"vu", "iter", "url", "boolField:bool", "boolField:bool"***REMOVED***
	_, err = MakeFieldKinds(*conf)
	require.Error(t, err)

	// Error case 2 (duplicated fields in bool and float ields)
	conf.TagsAsFields = []string***REMOVED***"vu", "iter", "url", "boolField:bool", "boolField:float"***REMOVED***
	_, err = MakeFieldKinds(*conf)
	require.Error(t, err)

	// Error case 3 (duplicated fields in BoolFields and IntFields)
	conf.TagsAsFields = []string***REMOVED***"vu", "iter", "url", "boolField:bool", "floatField:float", "boolField:int"***REMOVED***
	_, err = MakeFieldKinds(*conf)
	require.Error(t, err)

	// Normal case
	conf.TagsAsFields = []string***REMOVED***"vu", "iter", "url", "boolField:bool", "floatField:float", "intField:int"***REMOVED***
	fieldKinds, err = MakeFieldKinds(*conf)
	require.NoError(t, err)

	require.Equal(t, fieldKinds["boolField"], Bool)
	require.Equal(t, fieldKinds["floatField"], Float)
	require.Equal(t, fieldKinds["intField"], Int)
***REMOVED***
