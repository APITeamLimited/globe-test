/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2021 Load Impact
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

package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
)

func TestParseThresholdExpression(t *testing.T) ***REMOVED***
	t.Parallel()

	tests := []struct ***REMOVED***
		name           string
		input          string
		wantExpression *thresholdExpression
		wantErr        bool
	***REMOVED******REMOVED***
		***REMOVED***
			name:           "unknown expression's operator fails",
			input:          "count!20",
			wantExpression: nil,
			wantErr:        true,
		***REMOVED***,
		***REMOVED***
			name:           "unknown expression's method fails",
			input:          "foo>20",
			wantExpression: nil,
			wantErr:        true,
		***REMOVED***,
		***REMOVED***
			name:           "non numerical expression's value fails",
			input:          "count>abc",
			wantExpression: nil,
			wantErr:        true,
		***REMOVED***,
		***REMOVED***
			name:           "valid threshold expression syntax",
			input:          "count>20",
			wantExpression: &thresholdExpression***REMOVED***AggregationMethod: "count", Operator: ">", Value: 20***REMOVED***,
			wantErr:        false,
		***REMOVED***,
	***REMOVED***
	for _, testCase := range tests ***REMOVED***
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			gotExpression, gotErr := parseThresholdExpression(testCase.input)

			assert.Equal(t,
				testCase.wantErr,
				gotErr != nil,
				"parseThresholdExpression() error = %v, wantErr %v", gotErr, testCase.wantErr,
			)

			assert.Equal(t,
				testCase.wantExpression,
				gotExpression,
				"parseThresholdExpression() gotExpression = %v, want %v", gotExpression, testCase.wantExpression,
			)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func BenchmarkParseThresholdExpression(b *testing.B) ***REMOVED***
	for i := 0; i < b.N; i++ ***REMOVED***
		parseThresholdExpression("count>20") // nolint
	***REMOVED***
***REMOVED***

func TestParseThresholdAggregationMethod(t *testing.T) ***REMOVED***
	t.Parallel()

	tests := []struct ***REMOVED***
		name            string
		input           string
		wantMethod      string
		wantMethodValue null.Float
		wantErr         bool
	***REMOVED******REMOVED***
		***REMOVED***
			name:            "count method is parsed",
			input:           "count",
			wantMethod:      "count",
			wantMethodValue: null.Float***REMOVED******REMOVED***,
			wantErr:         false,
		***REMOVED***,
		***REMOVED***
			name:            "rate method is parsed",
			input:           "rate",
			wantMethod:      "rate",
			wantMethodValue: null.Float***REMOVED******REMOVED***,
			wantErr:         false,
		***REMOVED***,
		***REMOVED***
			name:            "value method is parsed",
			input:           "value",
			wantMethod:      "value",
			wantMethodValue: null.Float***REMOVED******REMOVED***,
			wantErr:         false,
		***REMOVED***,
		***REMOVED***
			name:            "avg method is parsed",
			input:           "avg",
			wantMethod:      "avg",
			wantMethodValue: null.Float***REMOVED******REMOVED***,
			wantErr:         false,
		***REMOVED***,
		***REMOVED***
			name:            "min method is parsed",
			input:           "min",
			wantMethod:      "min",
			wantMethodValue: null.Float***REMOVED******REMOVED***,
			wantErr:         false,
		***REMOVED***,
		***REMOVED***
			name:            "max method is parsed",
			input:           "max",
			wantMethod:      "max",
			wantMethodValue: null.Float***REMOVED******REMOVED***,
			wantErr:         false,
		***REMOVED***,
		***REMOVED***
			name:            "med method is parsed",
			input:           "med",
			wantMethod:      "med",
			wantMethodValue: null.Float***REMOVED******REMOVED***,
			wantErr:         false,
		***REMOVED***,
		***REMOVED***
			name:            "percentile method with integer value is parsed",
			input:           "p(99)",
			wantMethod:      "p(99)",
			wantMethodValue: null.FloatFrom(99),
			wantErr:         false,
		***REMOVED***,
		***REMOVED***
			name:            "percentile method with floating point value is parsed",
			input:           "p(99.9)",
			wantMethod:      "p(99.9)",
			wantMethodValue: null.FloatFrom(99.9),
			wantErr:         false,
		***REMOVED***,
		***REMOVED***
			name:            "parsing invalid method fails",
			input:           "foo",
			wantMethod:      "",
			wantMethodValue: null.Float***REMOVED******REMOVED***,
			wantErr:         true,
		***REMOVED***,
		***REMOVED***
			name:            "parsing empty percentile expression fails",
			input:           "p()",
			wantMethod:      "",
			wantMethodValue: null.Float***REMOVED******REMOVED***,
			wantErr:         true,
		***REMOVED***,
		***REMOVED***
			name:            "parsing incomplete percentile expression fails",
			input:           "p(99",
			wantMethod:      "",
			wantMethodValue: null.Float***REMOVED******REMOVED***,
			wantErr:         true,
		***REMOVED***,
		***REMOVED***
			name:            "parsing non-numerical percentile value fails",
			input:           "p(foo)",
			wantMethod:      "",
			wantMethodValue: null.Float***REMOVED******REMOVED***,
			wantErr:         true,
		***REMOVED***,
	***REMOVED***
	for _, testCase := range tests ***REMOVED***
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			gotMethod, gotMethodValue, gotErr := parseThresholdAggregationMethod(testCase.input)

			assert.Equal(t,
				testCase.wantErr,
				gotErr != nil,
				"parseThresholdAggregationMethod() error = %v, wantErr %v", gotErr, testCase.wantErr,
			)

			assert.Equal(t,
				testCase.wantMethod,
				gotMethod,
				"parseThresholdAggregationMethod() gotMethod = %v, want %v", gotMethod, testCase.wantMethod,
			)

			assert.Equal(t,
				testCase.wantMethodValue,
				gotMethodValue,
				"parseThresholdAggregationMethod() gotMethodValue = %v, want %v", gotMethodValue, testCase.wantMethodValue,
			)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func BenchmarkParseThresholdAggregationMethod(b *testing.B) ***REMOVED***
	for i := 0; i < b.N; i++ ***REMOVED***
		parseThresholdAggregationMethod("p(99.9)") // nolint
	***REMOVED***
***REMOVED***

func TestScanThresholdExpression(t *testing.T) ***REMOVED***
	t.Parallel()

	tests := []struct ***REMOVED***
		name         string
		input        string
		wantMethod   string
		wantOperator string
		wantValue    string
		wantErr      bool
	***REMOVED******REMOVED***
		***REMOVED***
			name:         "expression with <= operator is scanned",
			input:        "foo<=bar",
			wantMethod:   "foo",
			wantOperator: "<=",
			wantValue:    "bar",
			wantErr:      false,
		***REMOVED***,
		***REMOVED***
			name:         "expression with < operator is scanned",
			input:        "foo<bar",
			wantMethod:   "foo",
			wantOperator: "<",
			wantValue:    "bar",
			wantErr:      false,
		***REMOVED***,
		***REMOVED***
			name:         "expression with >= operator is scanned",
			input:        "foo>=bar",
			wantMethod:   "foo",
			wantOperator: ">=",
			wantValue:    "bar",
			wantErr:      false,
		***REMOVED***,
		***REMOVED***
			name:         "expression with > operator is scanned",
			input:        "foo>bar",
			wantMethod:   "foo",
			wantOperator: ">",
			wantValue:    "bar",
			wantErr:      false,
		***REMOVED***,
		***REMOVED***
			name:         "expression with === operator is scanned",
			input:        "foo===bar",
			wantMethod:   "foo",
			wantOperator: "===",
			wantValue:    "bar",
			wantErr:      false,
		***REMOVED***,
		***REMOVED***
			name:         "expression with == operator is scanned",
			input:        "foo==bar",
			wantMethod:   "foo",
			wantOperator: "==",
			wantValue:    "bar",
			wantErr:      false,
		***REMOVED***,
		***REMOVED***
			name:         "expression with != operator is scanned",
			input:        "foo!=bar",
			wantMethod:   "foo",
			wantOperator: "!=",
			wantValue:    "bar",
			wantErr:      false,
		***REMOVED***,
		***REMOVED***
			name:         "expression's method is trimmed",
			input:        "  foo  <=bar",
			wantMethod:   "foo",
			wantOperator: "<=",
			wantValue:    "bar",
			wantErr:      false,
		***REMOVED***,
		***REMOVED***
			name:         "expression's value is trimmed",
			input:        "foo<=  bar  ",
			wantMethod:   "foo",
			wantOperator: "<=",
			wantValue:    "bar",
			wantErr:      false,
		***REMOVED***,
		***REMOVED***
			name:         "expression with unknown operator fails",
			input:        "foo!bar",
			wantMethod:   "",
			wantOperator: "",
			wantValue:    "",
			wantErr:      true,
		***REMOVED***,
	***REMOVED***
	for _, testCase := range tests ***REMOVED***
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			gotMethod, gotOperator, gotValue, gotErr := scanThresholdExpression(testCase.input)

			assert.Equal(t,
				testCase.wantErr,
				gotErr != nil,
				"scanThresholdExpression() error = %v, wantErr %v", gotErr, testCase.wantErr,
			)

			assert.Equal(t,
				testCase.wantMethod,
				gotMethod,
				"scanThresholdExpression() gotMethod = %v, want %v", gotMethod, testCase.wantMethod,
			)

			assert.Equal(t,
				testCase.wantOperator,
				gotOperator,
				"scanThresholdExpression() gotOperator = %v, want %v", gotOperator, testCase.wantOperator,
			)

			assert.Equal(t,
				testCase.wantValue,
				gotValue,
				"scanThresholdExpression() gotValue = %v, want %v", gotValue, testCase.wantValue,
			)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func BenchmarkScanThresholdExpression(b *testing.B) ***REMOVED***
	for i := 0; i < b.N; i++ ***REMOVED***
		scanThresholdExpression("foo<=bar") // nolint
	***REMOVED***
***REMOVED***
