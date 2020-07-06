package js

import (
	"context"
	"testing"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkEmptyIteration(b *testing.B) ***REMOVED***
	b.StopTimer()

	r, err := getSimpleRunner("/script.js", `exports.default = function() ***REMOVED*** ***REMOVED***`)
	if !assert.NoError(b, err) ***REMOVED***
		return
	***REMOVED***
	require.NoError(b, err)

	var ch = make(chan stats.SampleContainer, 100)
	defer close(ch)
	go func() ***REMOVED*** // read the channel so it doesn't block
		for range ch ***REMOVED***
		***REMOVED***
	***REMOVED***()
	initVU, err := r.NewVU(1, ch)
	if !assert.NoError(b, err) ***REMOVED***
		return
	***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		err = vu.RunOnce()
		assert.NoError(b, err)
	***REMOVED***
***REMOVED***
