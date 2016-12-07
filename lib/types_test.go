package lib

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
	"testing"
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
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***Thresholds: map[string][]*Threshold***REMOVED***
			"metric": []*Threshold***REMOVED***&Threshold***REMOVED***Source: "1+1==2"***REMOVED******REMOVED***,
		***REMOVED******REMOVED***)
		assert.NotNil(t, opts.Thresholds)
		assert.NotEmpty(t, opts.Thresholds)
	***REMOVED***)
***REMOVED***
