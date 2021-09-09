package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/stats"
)

func TestRegistryNewMetric(t *testing.T) ***REMOVED***
	t.Parallel()
	r := NewRegistry()

	somethingCounter, err := r.NewMetric("something", stats.Counter)
	require.NoError(t, err)
	require.Equal(t, "something", somethingCounter.Name)

	somethingCounterAgain, err := r.NewMetric("something", stats.Counter)
	require.NoError(t, err)
	require.Equal(t, "something", somethingCounterAgain.Name)
	require.Same(t, somethingCounter, somethingCounterAgain)

	_, err = r.NewMetric("something", stats.Gauge)
	require.Error(t, err)

	_, err = r.NewMetric("something", stats.Counter, stats.Time)
	require.Error(t, err)
***REMOVED***

func TestMetricNames(t *testing.T) ***REMOVED***
	t.Parallel()
	testMap := map[string]bool***REMOVED***
		"simple":       true,
		"still_simple": true,
		"":             false,
		"@":            false,
		"a":            true,
		"special\n\t":  false,
		// this has both hangul and japanese numerals .
		"hello.World_in_한글一안녕一세상": true,
		// too long
		"tooolooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooog": false,
	***REMOVED***

	for key, value := range testMap ***REMOVED***
		key, value := key, value
		t.Run(key, func(t *testing.T) ***REMOVED***
			t.Parallel()
			assert.Equal(t, value, checkName(key), key)
		***REMOVED***)
	***REMOVED***
***REMOVED***
