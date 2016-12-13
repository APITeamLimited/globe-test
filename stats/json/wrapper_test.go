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

package json

import (
	"github.com/loadimpact/k6/stats"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWrapWithNilArg(t *testing.T) ***REMOVED***
	out := Wrap(nil)
	assert.Equal(t, out, (*Envelope)(nil))
***REMOVED***

func TestWrapWithUnusedType(t *testing.T) ***REMOVED***
	out := Wrap(JSONSample***REMOVED******REMOVED***)
	assert.Equal(t, out, (*Envelope)(nil))
***REMOVED***

func TestWrapWithSample(t *testing.T) ***REMOVED***
	out := Wrap(stats.Sample***REMOVED***
		Metric: &stats.Metric***REMOVED******REMOVED***,
	***REMOVED***)
	assert.NotEqual(t, out, (*Envelope)(nil))
***REMOVED***

func TestWrapWithMetricPointer(t *testing.T) ***REMOVED***
	out := Wrap(&stats.Metric***REMOVED******REMOVED***)
	assert.NotEqual(t, out, (*Envelope)(nil))
***REMOVED***

func TestWrapWithMetric(t *testing.T) ***REMOVED***
	out := Wrap(stats.Metric***REMOVED******REMOVED***)
	assert.Equal(t, out, (*Envelope)(nil))
***REMOVED***
