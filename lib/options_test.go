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

package lib

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
	"testing"
	"time"
)

func TestOptionsApply(t *testing.T) ***REMOVED***
	t.Run("Paused", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***Paused: null.BoolFrom(true)***REMOVED***)
		assert.True(t, opts.Paused.Valid)
		assert.True(t, opts.Paused.Bool)
	***REMOVED***)
	t.Run("VUs", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***VUs: null.IntFrom(12345)***REMOVED***)
		assert.True(t, opts.VUs.Valid)
		assert.Equal(t, int64(12345), opts.VUs.Int64)
	***REMOVED***)
	t.Run("VUsMax", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***VUsMax: null.IntFrom(12345)***REMOVED***)
		assert.True(t, opts.VUsMax.Valid)
		assert.Equal(t, int64(12345), opts.VUsMax.Int64)
	***REMOVED***)
	t.Run("Duration", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***Duration: null.StringFrom("2m")***REMOVED***)
		assert.True(t, opts.Duration.Valid)
		assert.Equal(t, "2m", opts.Duration.String)
	***REMOVED***)
	t.Run("Stages", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***Stages: []Stage***REMOVED******REMOVED***Duration: 1 * time.Second***REMOVED******REMOVED******REMOVED***)
		assert.NotNil(t, opts.Stages)
		assert.Len(t, opts.Stages, 1)
		assert.Equal(t, 1*time.Second, opts.Stages[0].Duration)
	***REMOVED***)
	t.Run("Linger", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***Linger: null.BoolFrom(true)***REMOVED***)
		assert.True(t, opts.Linger.Valid)
		assert.True(t, opts.Linger.Bool)
	***REMOVED***)
	t.Run("AbortOnTaint", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***AbortOnTaint: null.BoolFrom(true)***REMOVED***)
		assert.True(t, opts.AbortOnTaint.Valid)
		assert.True(t, opts.AbortOnTaint.Bool)
	***REMOVED***)
	t.Run("Acceptance", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***Acceptance: null.FloatFrom(12345.0)***REMOVED***)
		assert.True(t, opts.Acceptance.Valid)
		assert.Equal(t, float64(12345.0), opts.Acceptance.Float64)
	***REMOVED***)
	t.Run("MaxRedirects", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***MaxRedirects: null.IntFrom(12345)***REMOVED***)
		assert.True(t, opts.MaxRedirects.Valid)
		assert.Equal(t, int64(12345), opts.MaxRedirects.Int64)
	***REMOVED***)
	t.Run("Thresholds", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***Thresholds: map[string][]string***REMOVED***
			"metric": ***REMOVED***"1+1==2"***REMOVED***,
		***REMOVED******REMOVED***)
		assert.NotNil(t, opts.Thresholds)
		assert.NotEmpty(t, opts.Thresholds)
	***REMOVED***)
***REMOVED***
