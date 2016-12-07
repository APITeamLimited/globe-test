package js

import (
	"github.com/loadimpact/k6/lib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) ***REMOVED***
	r, err := New()
	assert.NoError(t, err)

	t.Run("Polyfill: Symbol", func(t *testing.T) ***REMOVED***
		v, err := r.VM.Get("Symbol")
		assert.NoError(t, err)
		assert.False(t, v.IsUndefined())
	***REMOVED***)
***REMOVED***

func TestLoad(t *testing.T) ***REMOVED***
	r, err := New()
	assert.NoError(t, err)
	assert.NoError(t, r.VM.Set("__initapi__", InitAPI***REMOVED***r: r***REMOVED***))

	t.Run("Importing Libraries", func(t *testing.T) ***REMOVED***
		_, err := r.load("test.js", []byte(`
			import "k6";
		`))
		assert.NoError(t, err)
		assert.Contains(t, r.lib, "k6.js")
	***REMOVED***)
***REMOVED***

func TestExtractOptions(t *testing.T) ***REMOVED***
	r, err := New()
	assert.NoError(t, err)

	t.Run("nothing", func(t *testing.T) ***REMOVED***
		_, err := r.load("test.js", []byte(``))
		assert.NoError(t, err)
	***REMOVED***)

	t.Run("vus", func(t *testing.T) ***REMOVED***
		_, err := r.load("test.js", []byte(`
			export let options = ***REMOVED*** vus: 12345 ***REMOVED***;
		`))
		assert.NoError(t, err)

		assert.True(t, r.Options.VUs.Valid)
		assert.Equal(t, int64(12345), r.Options.VUs.Int64)
	***REMOVED***)
	t.Run("vus-max", func(t *testing.T) ***REMOVED***
		_, err := r.load("test.js", []byte(`
			export let options = ***REMOVED*** "vus-max": 12345 ***REMOVED***;
		`))
		assert.NoError(t, err)

		assert.True(t, r.Options.VUsMax.Valid)
		assert.Equal(t, int64(12345), r.Options.VUsMax.Int64)
	***REMOVED***)
	t.Run("duration", func(t *testing.T) ***REMOVED***
		_, err := r.load("test.js", []byte(`
			export let options = ***REMOVED*** duration: "2m" ***REMOVED***;
		`))
		assert.NoError(t, err)

		assert.True(t, r.Options.Duration.Valid)
		assert.Equal(t, "2m", r.Options.Duration.String)
	***REMOVED***)
	t.Run("max-redirects", func(t *testing.T) ***REMOVED***
		_, err := r.load("test.js", []byte(`
			export let options = ***REMOVED*** "max-redirects": 12345 ***REMOVED***;
		`))
		assert.NoError(t, err)

		assert.True(t, r.Options.MaxRedirects.Valid)
		assert.Equal(t, int64(12345), r.Options.MaxRedirects.Int64)
	***REMOVED***)
	t.Run("thresholds", func(t *testing.T) ***REMOVED***
		_, err := r.load("test.js", []byte(`
			export let options = ***REMOVED***
				thresholds: ***REMOVED***
					my_metric: ["value<=1000"],
				***REMOVED***
			***REMOVED***
		`))
		assert.NoError(t, err)

		assert.Contains(t, r.Options.Thresholds, "my_metric")
		if assert.Len(t, r.Options.Thresholds["my_metric"], 1) ***REMOVED***
			assert.Equal(t, &lib.Threshold***REMOVED***Source: "value<=1000"***REMOVED***, r.Options.Thresholds["my_metric"][0])
		***REMOVED***
	***REMOVED***)
***REMOVED***
