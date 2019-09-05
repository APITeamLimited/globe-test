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

package csv

import (
	"testing"
	"time"

	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib/types"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) ***REMOVED***
	config := NewConfig()
	assert.Equal(t, "file.csv", config.FileName.String)
	assert.Equal(t, "1s", config.SaveInterval.String())
***REMOVED***

func TestApply(t *testing.T) ***REMOVED***
	configs := []Config***REMOVED***
		***REMOVED***
			FileName:     null.StringFrom(""),
			SaveInterval: types.NullDurationFrom(2 * time.Second),
		***REMOVED***,
		***REMOVED***
			FileName:     null.StringFrom("newPath"),
			SaveInterval: types.NewNullDuration(time.Duration(1), false),
		***REMOVED***,
	***REMOVED***
	expected := []struct ***REMOVED***
		FileName     string
		SaveInterval string
	***REMOVED******REMOVED***
		***REMOVED***
			FileName:     "",
			SaveInterval: "2s",
		***REMOVED***,
		***REMOVED***
			FileName:     "newPath",
			SaveInterval: "1s",
		***REMOVED***,
	***REMOVED***

	for i := range configs ***REMOVED***
		config := configs[i]
		expected := expected[i]
		t.Run(expected.FileName+"_"+expected.SaveInterval, func(t *testing.T) ***REMOVED***
			baseConfig := NewConfig()
			baseConfig = baseConfig.Apply(config)

			assert.Equal(t, expected.FileName, baseConfig.FileName.String)
			assert.Equal(t, expected.SaveInterval, baseConfig.SaveInterval.String())
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestParseArg(t *testing.T) ***REMOVED***
	args := []string***REMOVED***
		"test_file.csv",
		"file_name=test.csv,save_interval=5s",
	***REMOVED***

	expected := []Config***REMOVED***
		***REMOVED***
			FileName:     null.StringFrom("test_file.csv"),
			SaveInterval: types.NullDurationFrom(1 * time.Second),
		***REMOVED***,
		***REMOVED***
			FileName:     null.StringFrom("test.csv"),
			SaveInterval: types.NullDurationFrom(5 * time.Second),
		***REMOVED***,
	***REMOVED***

	for i := range args ***REMOVED***
		arg := args[i]
		expected := expected[i]

		t.Run(expected.FileName.String+"_"+expected.SaveInterval.String(), func(t *testing.T) ***REMOVED***
			config, err := ParseArg(arg)

			assert.Nil(t, err)
			assert.Equal(t, expected.FileName.String, config.FileName.String)
			assert.Equal(t, expected.SaveInterval.String(), config.SaveInterval.String())
		***REMOVED***)
	***REMOVED***
***REMOVED***
