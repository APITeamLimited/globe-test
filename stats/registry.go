package stats

import (
	"sync"
)

type Registry struct ***REMOVED***
	Backends []Backend

	collectors []*Collector
	mutex      sync.Mutex
***REMOVED***

func (r *Registry) NewCollector() *Collector ***REMOVED***
	collector := &Collector***REMOVED******REMOVED***

	r.mutex.Lock()
	r.collectors = append(r.collectors, collector)
	r.mutex.Unlock()

	return collector
***REMOVED***

func (r *Registry) Submit() error ***REMOVED***
	batches := make([][]Point, 0, len(r.collectors))
	for _, collector := range r.collectors ***REMOVED***
		batch := collector.drain()
		batches = append(batches, batch)
	***REMOVED***

	for _, backend := range r.Backends ***REMOVED***
		if err := backend.Submit(batches); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
