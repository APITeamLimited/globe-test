package stats

import (
	"sync"
	"time"
)

type Collector struct ***REMOVED***
	batch []Point
	mutex sync.Mutex
***REMOVED***

func (c *Collector) Add(p Point) ***REMOVED***
	if p.Stat == nil || len(p.Values) == 0 ***REMOVED***
		return
	***REMOVED***
	if p.Time.IsZero() ***REMOVED***
		p.Time = time.Now()
	***REMOVED***

	c.mutex.Lock()
	c.batch = append(c.batch, p)
	c.mutex.Unlock()
***REMOVED***

func (c *Collector) drain() []Point ***REMOVED***
	c.mutex.Lock()
	oldBatch := c.batch
	c.batch = nil
	c.mutex.Unlock()

	return oldBatch
***REMOVED***
