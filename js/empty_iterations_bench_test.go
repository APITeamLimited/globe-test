package js

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/require"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/metrics"
)

func BenchmarkEmptyIteration(b *testing.B) {
	b.StopTimer()

	r, err := getSimpleRunner(b, "/script.js", `exports.default = function() { }`)
	require.NoError(b, err)

	ch := make(chan metrics.SampleContainer, 100)
	defer close(ch)
	go func() { // read the channel so it doesn't block
		for range ch {
		}
	}()
	initVU, err := r.NewVU(1, 1, ch, redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}))
	require.NoError(b, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	vu := initVU.Activate(&lib.VUActivationParams{RunContext: ctx})
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err = vu.RunOnce()
		require.NoError(b, err)
	}
}
