package js

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
)

func BenchmarkEmptyIteration(b *testing.B) ***REMOVED***
	b.StopTimer()

	r, err := getSimpleRunner(b, "/script.js", `exports.default = function() ***REMOVED*** ***REMOVED***`)
	require.NoError(b, err)

	ch := make(chan workerMetrics.SampleContainer, 100)
	defer close(ch)
	go func() ***REMOVED*** // read the channel so it doesn't block
		for range ch ***REMOVED***
		***REMOVED***
	***REMOVED***()
	initVU, err := r.NewVU(1, 1, ch, libWorker.GetTestWorkerInfo())
	require.NoError(b, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	vu := initVU.Activate(&libWorker.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		err = vu.RunOnce()
		require.NoError(b, err)
	***REMOVED***
***REMOVED***
