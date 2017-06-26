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
	"crypto/tls"
	"testing"
	"time"

	"github.com/loadimpact/k6/stats"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
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
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***Duration: NullDurationFrom(2 * time.Minute)***REMOVED***)
		assert.True(t, opts.Duration.Valid)
		assert.Equal(t, "2m0s", opts.Duration.String())
	***REMOVED***)
	t.Run("Iterations", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***Iterations: null.IntFrom(1234)***REMOVED***)
		assert.True(t, opts.Iterations.Valid)
		assert.Equal(t, int64(1234), opts.Iterations.Int64)
	***REMOVED***)
	t.Run("Stages", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***Stages: []Stage***REMOVED******REMOVED***Duration: NullDurationFrom(1 * time.Second)***REMOVED******REMOVED******REMOVED***)
		assert.NotNil(t, opts.Stages)
		assert.Len(t, opts.Stages, 1)
		assert.Equal(t, 1*time.Second, time.Duration(opts.Stages[0].Duration.Duration))
	***REMOVED***)
	t.Run("Linger", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***Linger: null.BoolFrom(true)***REMOVED***)
		assert.True(t, opts.Linger.Valid)
		assert.True(t, opts.Linger.Bool)
	***REMOVED***)
	t.Run("MaxRedirects", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***MaxRedirects: null.IntFrom(12345)***REMOVED***)
		assert.True(t, opts.MaxRedirects.Valid)
		assert.Equal(t, int64(12345), opts.MaxRedirects.Int64)
	***REMOVED***)
	t.Run("InsecureSkipTLSVerify", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***InsecureSkipTLSVerify: null.BoolFrom(true)***REMOVED***)
		assert.True(t, opts.InsecureSkipTLSVerify.Valid)
		assert.True(t, opts.InsecureSkipTLSVerify.Bool)
	***REMOVED***)
	t.Run("TLSCipherSuites", func(t *testing.T) ***REMOVED***
		for suiteName, suiteID := range SupportedTLSCipherSuites ***REMOVED***
			t.Run(suiteName, func(t *testing.T) ***REMOVED***
				opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***TLSCipherSuites: &TLSCipherSuites***REMOVED***suiteID***REMOVED******REMOVED***)

				assert.NotNil(t, opts.TLSCipherSuites)
				assert.Len(t, *(opts.TLSCipherSuites), 1)
				assert.Equal(t, suiteID, (*opts.TLSCipherSuites)[0])
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
	t.Run("TLSVersion", func(t *testing.T) ***REMOVED***
		version := TLSVersion***REMOVED***Min: tls.VersionSSL30, Max: tls.VersionTLS12***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***TLSVersion: &version***REMOVED***)

		assert.NotNil(t, opts.TLSVersion)
		assert.Equal(t, opts.TLSVersion.Min, tls.VersionSSL30)
		assert.Equal(t, opts.TLSVersion.Max, tls.VersionTLS12)
	***REMOVED***)
	t.Run("NoConnectionReuse", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***NoConnectionReuse: null.BoolFrom(true)***REMOVED***)
		assert.True(t, opts.NoConnectionReuse.Valid)
		assert.True(t, opts.NoConnectionReuse.Bool)
	***REMOVED***)
	t.Run("Thresholds", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***Thresholds: map[string]stats.Thresholds***REMOVED***
			"metric": ***REMOVED***
				Thresholds: []*stats.Threshold***REMOVED******REMOVED******REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED******REMOVED***)
		assert.NotNil(t, opts.Thresholds)
		assert.NotEmpty(t, opts.Thresholds)
	***REMOVED***)
	t.Run("External", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***External: map[string]interface***REMOVED******REMOVED******REMOVED***"a": 1***REMOVED******REMOVED***)
		assert.Equal(t, map[string]interface***REMOVED******REMOVED******REMOVED***"a": 1***REMOVED***, opts.External)
	***REMOVED***)
	t.Run("NoUsageReport", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***NoUsageReport: null.BoolFrom(true)***REMOVED***)
		assert.True(t, opts.NoUsageReport.Valid)
		assert.True(t, opts.NoUsageReport.Bool)
	***REMOVED***)
***REMOVED***
