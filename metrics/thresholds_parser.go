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

package metrics

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/guregu/null.v3"
)

// thresholdExpression holds the parsed result of a threshold expression,
// as described in: https://k6.io/docs/using-k6/thresholds/#threshold-syntax
type thresholdExpression struct ***REMOVED***
	// AggregationMethod holds the aggregation method parsed
	// from the threshold expression. Possible values are described
	// by `aggregationMethodTokens`.
	AggregationMethod string

	// AggregationValue will hold the aggregation method's pivot value
	// in the event it is a percentile. For instance: an expression of the form p(99.9) < 200,
	// would result in AggregationValue to be set to 99.9.
	AggregationValue null.Float

	// Operator holds the operator parsed from the threshold expression.
	// Possible values are described by `operatorTokens`.
	Operator string

	// Value holds the value parsed from the threshold expression.
	Value float64
***REMOVED***

// SinkKey computes the key used to index a thresholdExpression in the engine's sinks.
//
// During execution, the engine "sinks" metrics into a internal mapping, so that
// thresholds can be asserted against them. This method is a helper to normalize the
// sink the threshold expression should be applied to.
//
// Because a threshold expression's aggregation method can either be
// a static keyword ("count", "rate", etc...), or a parametric
// expression ("p(somefloatingpointvalue)"), we need to handle this
// case specifically. If we encounter the percentile aggregation method token,
// we recompute the whole "p(value)" expression in order to look for it in the
// sinks.
func (te *thresholdExpression) SinkKey() string ***REMOVED***
	if te.AggregationMethod == tokenPercentile ***REMOVED***
		return fmt.Sprintf("%s(%g)", tokenPercentile, te.AggregationValue.Float64)
	***REMOVED***

	return te.AggregationMethod
***REMOVED***

// parseThresholdAssertion parses a threshold condition expression,
// as defined in a JS script (for instance p(95)<1000), into a thresholdExpression
// instance.
//
// It is expected to be of the form: `aggregation_method operator value`.
// As defined by the following BNF:
// ```
// assertion           -> aggregation_method whitespace* operator whitespace* float
// aggregation_method  -> trend | rate | gauge | counter
// counter             -> "count" | "rate"
// gauge               -> "value"
// rate                -> "rate"
// trend               -> "avg" | "min" | "max" | "med" | percentile
// percentile          -> "p(" float ")"
// operator            -> ">" | ">=" | "<=" | "<" | "==" | "===" | "!="
// float               -> digit+ ("." digit+)?
// digit               -> "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
// whitespace          -> " "
// ```
func parseThresholdExpression(input string) (*thresholdExpression, error) ***REMOVED***
	// Scanning makes no assumption on the underlying values, and only
	// checks that the expression has the right format.
	method, operator, value, err := scanThresholdExpression(input)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed parsing threshold expression %q; reason: %w", input, err)
	***REMOVED***

	parsedMethod, parsedMethodValue, err := parseThresholdAggregationMethod(method)
	if err != nil ***REMOVED***
		err = fmt.Errorf("failed parsing threshold expression's %q left hand side; "+
			"reason: %w", input, err,
		)
		return nil, err
	***REMOVED***

	parsedValue, err := strconv.ParseFloat(value, 64)
	if err != nil ***REMOVED***
		err = fmt.Errorf("failed parsing threshold expresion's %q right hand side; "+
			"reason: %w", input, err,
		)
		return nil, err
	***REMOVED***

	condition := &thresholdExpression***REMOVED***
		AggregationMethod: parsedMethod,
		AggregationValue:  parsedMethodValue,
		Operator:          operator,
		Value:             parsedValue,
	***REMOVED***

	return condition, nil
***REMOVED***

// Define accepted threshold expression operators tokens
const (
	tokenLessEqual     = "<="
	tokenLess          = "<"
	tokenGreaterEqual  = ">="
	tokenGreater       = ">"
	tokenStrictlyEqual = "==="
	tokenLooselyEqual  = "=="
	tokenBangEqual     = "!="
)

// operatorTokens defines the list of operator-related tokens
// used in threshold expressions parsing.
//
// It is meant to be used during the scan of threshold expressions.
// Although declared as a `var`, being an array, it is effectively
// immutable and can be considered constant.
//
// Note that because scanning uses a substring parser, and will match
// the smallest common substring, the actual slice order matters.
// Longer tokens with symbols in common with shorter ones must appear
// first in the list in order to be effectively matched.
var operatorTokens = [7]string***REMOVED*** //nolint:gochecknoglobals
	tokenLessEqual,
	tokenLess,
	tokenGreaterEqual,
	tokenGreater,
	tokenStrictlyEqual,
	tokenLooselyEqual,
	tokenBangEqual,
***REMOVED***

// scanThresholdExpression scans a threshold condition expression of the
// form: `aggregation_method operator value`. An invalid or unknown operator
// will produce an error. However, no assertions regarding
// either the left-hand side aggregation method nor the right-hand
// side value will be made: they will be returned as is, only trimmed from
// their spaces.
func scanThresholdExpression(input string) (string, string, string, error) ***REMOVED***
	for _, op := range operatorTokens ***REMOVED***
		substrings := strings.SplitN(input, op, 2)
		if len(substrings) == 2 ***REMOVED***
			return strings.TrimSpace(substrings[0]), op, strings.TrimSpace(substrings[1]), nil
		***REMOVED***
	***REMOVED***

	return "", "", "", fmt.Errorf("malformed threshold expression")
***REMOVED***

// Define accepted threshold expression aggregation tokens
// Percentile token `p(..)` is accepted too but handled separately.
const (
	tokenValue      = "value"
	tokenCount      = "count"
	tokenRate       = "rate"
	tokenAvg        = "avg"
	tokenMin        = "min"
	tokenMed        = "med"
	tokenMax        = "max"
	tokenPercentile = "p"
)

// aggregationMethodTokens defines the list of aggregation method
// used in the parsing of threshold expressions.
//
// It is meant to be used during the parsing of threshold expressions.
// Although declared as a `var`, being an array, it is effectively
// immutable and can be considered constant.
var aggregationMethodTokens = [8]string***REMOVED*** //nolint:gochecknoglobals
	tokenValue,
	tokenCount,
	tokenRate,
	tokenAvg,
	tokenMin,
	tokenMed,
	tokenMax,
	tokenPercentile,
***REMOVED***

// parseThresholdMethod will parse a threshold condition expression's method.
// It assumes the provided input argument is already trimmed and cleaned up.
// If it encounters a percentile method, it will parse it and verify it
// boils down to an expression of the form: `p(float64)`, but will return
// it verbatim, as a string.
func parseThresholdAggregationMethod(input string) (string, null.Float, error) ***REMOVED***
	// Is the input one of the methods keywords?
	for _, m := range aggregationMethodTokens ***REMOVED***
		// Percentile expressions being of the form p(value),
		// they won't be matched here.
		if m == input ***REMOVED***
			return m, null.Float***REMOVED******REMOVED***, nil
		***REMOVED***
	***REMOVED***

	// Otherwise, attempt to parse a percentile expression
	if strings.HasPrefix(input, tokenPercentile+"(") && strings.HasSuffix(input, ")") ***REMOVED***
		aggregationValue, err := strconv.ParseFloat(trimDelimited("p(", input, ")"), 64)
		if err != nil ***REMOVED***
			return "", null.Float***REMOVED******REMOVED***, fmt.Errorf("malformed percentile value; reason: %w", err)
		***REMOVED***

		return tokenPercentile, null.FloatFrom(aggregationValue), nil
	***REMOVED***

	return "", null.Float***REMOVED******REMOVED***, fmt.Errorf("failed parsing method from expression")
***REMOVED***

func trimDelimited(prefix, input, suffix string) string ***REMOVED***
	return strings.TrimSuffix(strings.TrimPrefix(input, prefix), suffix)
***REMOVED***
