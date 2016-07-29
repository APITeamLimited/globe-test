package stats

import (
	"sync"
)

type Registry struct ***REMOVED***
	Backends  []Backend
	ExtraTags Tags

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
	r.mutex.Lock()
	defer r.mutex.Unlock()

	batches := make([][]Sample, 0, len(r.collectors))
	for _, collector := range r.collectors ***REMOVED***
		batch := collector.drain()
		batches = append(batches, batch)
	***REMOVED***

	if len(r.ExtraTags) > 0 ***REMOVED***
		for _, batch := range batches ***REMOVED***
			for i, p := range batch ***REMOVED***
				if p.Tags == nil ***REMOVED***
					p.Tags = r.ExtraTags
					batch[i] = p
				***REMOVED*** else ***REMOVED***
					for key, val := range r.ExtraTags ***REMOVED***
						p.Tags[key] = val
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	for _, backend := range r.Backends ***REMOVED***
		if err := backend.Submit(batches); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
