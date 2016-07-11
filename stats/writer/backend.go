package writer

import (
	"github.com/loadimpact/speedboat/stats"
	"io"
	"sync"
)

type Formatter interface ***REMOVED***
	Format(data interface***REMOVED******REMOVED***) ([]byte, error)
***REMOVED***

type Backend struct ***REMOVED***
	Only    map[string]bool
	Exclude map[string]bool

	Writer    io.Writer
	Formatter Formatter

	mutex sync.Mutex
***REMOVED***

func (b Backend) Submit(batches [][]stats.Sample) error ***REMOVED***
	b.mutex.Lock()
	defer b.mutex.Unlock()

	hasOnly := len(b.Only) > 0

	for _, batch := range batches ***REMOVED***
		for _, s := range batch ***REMOVED***
			if hasOnly && !b.Only[s.Stat.Name] ***REMOVED***
				continue
			***REMOVED***
			if b.Exclude[s.Stat.Name] ***REMOVED***
				continue
			***REMOVED***

			data := b.Format(&s)
			bytes, err := b.Formatter.Format(data)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			if _, err := b.Writer.Write(bytes); err != nil ***REMOVED***
				return err
			***REMOVED***
			if _, err := b.Writer.Write([]byte***REMOVED***'\n'***REMOVED***); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (b Backend) Format(s *stats.Sample) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	data := map[string]interface***REMOVED******REMOVED******REMOVED***
		"time":   s.Time,
		"stat":   s.Stat.Name,
		"tags":   s.Tags,
		"values": s.Values,
	***REMOVED***
	if s.Tags == nil ***REMOVED***
		data["tags"] = stats.Tags***REMOVED******REMOVED***
	***REMOVED***
	return data
***REMOVED***
