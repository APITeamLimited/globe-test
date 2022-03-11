package json

import (
	"bytes"
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/output"
)

func BenchmarkFlushMetrics(b *testing.B) ***REMOVED***
	stdout := new(bytes.Buffer)
	dir := b.TempDir()
	out, err := New(output.Params***REMOVED***
		Logger:         testutils.NewLogger(b),
		StdOut:         stdout,
		FS:             afero.NewOsFs(),
		ConfigArgument: path.Join(dir, "test.gz"),
	***REMOVED***)
	require.NoError(b, err)

	samples, _ := generateTestMetricSamples(b)
	size := 10000
	for len(samples) < size ***REMOVED***
		more, _ := generateTestMetricSamples(b)
		samples = append(samples, more...)
	***REMOVED***
	samples = samples[:size]
	o, _ := out.(*Output)
	require.NoError(b, o.Start())
	o.periodicFlusher.Stop()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ ***REMOVED***
		o.AddMetricSamples(samples)
		o.flushMetrics()
	***REMOVED***
***REMOVED***
