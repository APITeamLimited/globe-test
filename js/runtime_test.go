package js

import (
	"github.com/loadimpact/speedboat/lib"
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
			import "speedboat";
		`))
		assert.NoError(t, err)
		assert.Contains(t, r.lib, "speedboat.js")
	***REMOVED***)
***REMOVED***

func TestExtractOptions(t *testing.T) ***REMOVED***
	r, err := New()
	assert.NoError(t, err)

	t.Run("nothing", func(t *testing.T) ***REMOVED***
		exp, err := r.load("test.js", []byte(``))
		assert.NoError(t, err)

		var opts lib.Options
		assert.NoError(t, r.ExtractOptions(exp, &opts))
	***REMOVED***)

	t.Run("vus", func(t *testing.T) ***REMOVED***
		exp, err := r.load("test.js", []byte(`
			export let options = ***REMOVED*** vus: 12345 ***REMOVED***;
		`))
		assert.NoError(t, err)

		var opts lib.Options
		assert.NoError(t, r.ExtractOptions(exp, &opts))
		assert.True(t, opts.VUs.Valid)
		assert.Equal(t, int64(12345), opts.VUs.Int64)
	***REMOVED***)
	t.Run("vus-max", func(t *testing.T) ***REMOVED***
		exp, err := r.load("test.js", []byte(`
			export let options = ***REMOVED*** "vus-max": 12345 ***REMOVED***;
		`))
		assert.NoError(t, err)

		var opts lib.Options
		assert.NoError(t, r.ExtractOptions(exp, &opts))
		assert.True(t, opts.VUsMax.Valid)
		assert.Equal(t, int64(12345), opts.VUsMax.Int64)
	***REMOVED***)
	t.Run("duration", func(t *testing.T) ***REMOVED***
		exp, err := r.load("test.js", []byte(`
			export let options = ***REMOVED*** duration: "2m" ***REMOVED***;
		`))
		assert.NoError(t, err)

		var opts lib.Options
		assert.NoError(t, r.ExtractOptions(exp, &opts))
		assert.True(t, opts.Duration.Valid)
		assert.Equal(t, "2m", opts.Duration.String)
	***REMOVED***)
	t.Run("thresholds", func(t *testing.T) ***REMOVED***
		exp, err := r.load("test.js", []byte(`
			export let options = ***REMOVED***
				thresholds: ***REMOVED***
					my_metric: ["value<=1000"],
				***REMOVED***
			***REMOVED***
		`))
		assert.NoError(t, err)

		var opts lib.Options
		assert.NoError(t, r.ExtractOptions(exp, &opts))
		assert.Contains(t, opts.Thresholds, "my_metric")
		if assert.Len(t, opts.Thresholds["my_metric"], 1) ***REMOVED***
			assert.Equal(t, &lib.Threshold***REMOVED***Source: "value<=1000"***REMOVED***, opts.Thresholds["my_metric"][0])
		***REMOVED***
	***REMOVED***)
***REMOVED***
