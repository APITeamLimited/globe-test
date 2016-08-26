package stats

import (
	"sync"
	"time"
)

type Collector struct ***REMOVED***
	Batch []Sample
	mutex sync.Mutex
***REMOVED***

func (c *Collector) Add(s Sample) ***REMOVED***
	if s.Stat == nil || len(s.Values) == 0 ***REMOVED***
		return
	***REMOVED***
	if s.Time.IsZero() ***REMOVED***
		s.Time = time.Now()
	***REMOVED***

	c.mutex.Lock()
	c.Batch = append(c.Batch, s)
	c.mutex.Unlock()
***REMOVED***

func (c *Collector) drain() []Sample ***REMOVED***
	c.mutex.Lock()
	oldBatch := c.Batch
	c.Batch = nil
	c.mutex.Unlock()

	return oldBatch
***REMOVED***
