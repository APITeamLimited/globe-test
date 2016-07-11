package stats

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testBackend struct ***REMOVED***
	submitted []Sample
***REMOVED***

func (b *testBackend) Submit(batches [][]Sample) error ***REMOVED***
	for _, batch := range batches ***REMOVED***
		for _, s := range batch ***REMOVED***
			b.submitted = append(b.submitted, s)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func TestNewCollector(t *testing.T) ***REMOVED***
	r := Registry***REMOVED******REMOVED***
	c := r.NewCollector()
	assert.Equal(t, 1, len(r.collectors))
	assert.Equal(t, c, r.collectors[0])
***REMOVED***

func TestSubmit(t *testing.T) ***REMOVED***
	backend := &testBackend***REMOVED******REMOVED***
	r := Registry***REMOVED***
		Backends: []Backend***REMOVED***backend***REMOVED***,
	***REMOVED***
	stat := Stat***REMOVED***Name: "test"***REMOVED***

	c1 := r.NewCollector()
	c1.Add(Sample***REMOVED***Stat: &stat, Values: Value(1)***REMOVED***)
	c1.Add(Sample***REMOVED***Stat: &stat, Values: Value(2)***REMOVED***)

	c2 := r.NewCollector()
	c2.Add(Sample***REMOVED***Stat: &stat, Values: Value(3)***REMOVED***)
	c2.Add(Sample***REMOVED***Stat: &stat, Values: Value(4)***REMOVED***)

	err := r.Submit()
	assert.NoError(t, err)
	assert.Len(t, backend.submitted, 4)
***REMOVED***
