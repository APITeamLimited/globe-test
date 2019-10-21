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

package ui

import (
	"testing"

	"github.com/loadimpact/k6/stats"
	"github.com/stretchr/testify/assert"
)

var verifyTests = []struct ***REMOVED***
	in  string
	out error
***REMOVED******REMOVED***
	***REMOVED***"avg", nil***REMOVED***,
	***REMOVED***"min", nil***REMOVED***,
	***REMOVED***"med", nil***REMOVED***,
	***REMOVED***"max", nil***REMOVED***,
	***REMOVED***"p(0)", nil***REMOVED***,
	***REMOVED***"p(90)", nil***REMOVED***,
	***REMOVED***"p(95)", nil***REMOVED***,
	***REMOVED***"p(99)", nil***REMOVED***,
	***REMOVED***"p(99.9)", nil***REMOVED***,
	***REMOVED***"p(99.9999)", nil***REMOVED***,
	***REMOVED***"nil", ErrStatUnknownFormat***REMOVED***,
	***REMOVED***" avg", ErrStatUnknownFormat***REMOVED***,
	***REMOVED***"avg ", ErrStatUnknownFormat***REMOVED***,
	***REMOVED***"", ErrStatEmptyString***REMOVED***,
***REMOVED***

var defaultTrendColumns = TrendColumns

func createTestTrendSink(count int) *stats.TrendSink ***REMOVED***
	sink := stats.TrendSink***REMOVED******REMOVED***

	for i := 0; i < count; i++ ***REMOVED***
		sink.Add(stats.Sample***REMOVED***Value: float64(i)***REMOVED***)
	***REMOVED***

	return &sink
***REMOVED***

func TestVerifyTrendColumnStat(t *testing.T) ***REMOVED***
	for _, testCase := range verifyTests ***REMOVED***
		err := VerifyTrendColumnStat(testCase.in)
		assert.Equal(t, testCase.out, err)
	***REMOVED***
***REMOVED***

func TestUpdateTrendColumns(t *testing.T) ***REMOVED***
	sink := createTestTrendSink(100)

	t.Run("No stats", func(t *testing.T) ***REMOVED***
		TrendColumns = defaultTrendColumns

		UpdateTrendColumns(make([]string, 0))

		assert.Equal(t, defaultTrendColumns, TrendColumns)
	***REMOVED***)

	t.Run("One stat", func(t *testing.T) ***REMOVED***
		TrendColumns = defaultTrendColumns

		UpdateTrendColumns([]string***REMOVED***"avg"***REMOVED***)

		assert.Exactly(t, 1, len(TrendColumns))
		assert.Exactly(t,
			sink.Avg,
			TrendColumns[0].Get(sink))
	***REMOVED***)

	t.Run("Multiple stats", func(t *testing.T) ***REMOVED***
		TrendColumns = defaultTrendColumns

		UpdateTrendColumns([]string***REMOVED***"med", "max"***REMOVED***)

		assert.Exactly(t, 2, len(TrendColumns))
		assert.Exactly(t, sink.Med, TrendColumns[0].Get(sink))
		assert.Exactly(t, sink.Max, TrendColumns[1].Get(sink))
	***REMOVED***)

	t.Run("Ignore invalid stats", func(t *testing.T) ***REMOVED***
		TrendColumns = defaultTrendColumns

		UpdateTrendColumns([]string***REMOVED***"med", "max", "invalid"***REMOVED***)

		assert.Exactly(t, 2, len(TrendColumns))
		assert.Exactly(t, sink.Med, TrendColumns[0].Get(sink))
		assert.Exactly(t, sink.Max, TrendColumns[1].Get(sink))
	***REMOVED***)

	t.Run("Percentile stats", func(t *testing.T) ***REMOVED***
		TrendColumns = defaultTrendColumns

		UpdateTrendColumns([]string***REMOVED***"p(99.9999)"***REMOVED***)

		assert.Exactly(t, 1, len(TrendColumns))
		assert.Exactly(t, sink.P(0.999999), TrendColumns[0].Get(sink))
	***REMOVED***)
***REMOVED***

func TestGeneratePercentileTrendColumn(t *testing.T) ***REMOVED***
	sink := createTestTrendSink(100)

	t.Run("Happy path", func(t *testing.T) ***REMOVED***
		colFunc, err := generatePercentileTrendColumn("p(99)")
		assert.Nil(t, err)
		assert.NotNil(t, colFunc)
		assert.Exactly(t, sink.P(0.99), colFunc(sink))
		assert.NotEqual(t, sink.P(0.98), colFunc(sink))
		assert.Nil(t, err)
	***REMOVED***)

	t.Run("Empty stat", func(t *testing.T) ***REMOVED***
		colFunc, err := generatePercentileTrendColumn("")

		assert.Nil(t, colFunc)
		assert.Exactly(t, err, ErrStatEmptyString)
	***REMOVED***)

	t.Run("Invalid format", func(t *testing.T) ***REMOVED***
		colFunc, err := generatePercentileTrendColumn("p90")

		assert.Nil(t, colFunc)
		assert.Exactly(t, err, ErrStatUnknownFormat)
	***REMOVED***)

	t.Run("Invalid format 2", func(t *testing.T) ***REMOVED***
		colFunc, err := generatePercentileTrendColumn("p(90")

		assert.Nil(t, colFunc)
		assert.Exactly(t, err, ErrStatUnknownFormat)
	***REMOVED***)

	t.Run("Invalid float", func(t *testing.T) ***REMOVED***
		colFunc, err := generatePercentileTrendColumn("p(a)")

		assert.Nil(t, colFunc)
		assert.Exactly(t, err, ErrPercentileStatInvalidValue)
	***REMOVED***)
***REMOVED***
