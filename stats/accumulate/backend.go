package accumulate

import (
	"github.com/loadimpact/speedboat/stats"
	"sort"
	"sync"
)

type Backend struct ***REMOVED***
	Data    map[*stats.Stat]map[*string]*Dimension
	Only    map[string]bool
	Exclude map[string]bool

	interned    map[string]*string
	submitMutex sync.Mutex
***REMOVED***

func New() *Backend ***REMOVED***
	return &Backend***REMOVED***
		Data:     make(map[*stats.Stat]map[*string]*Dimension),
		Exclude:  make(map[string]bool),
		Only:     make(map[string]bool),
		interned: make(map[string]*string),
	***REMOVED***
***REMOVED***

func (b *Backend) Get(stat *stats.Stat, dname string) *Dimension ***REMOVED***
	dimensions, ok := b.Data[stat]
	if !ok ***REMOVED***
		return nil
	***REMOVED***

	return dimensions[b.interned[dname]]
***REMOVED***

func (b *Backend) Submit(batches [][]stats.Point) error ***REMOVED***
	b.submitMutex.Lock()

	hasOnly := len(b.Only) > 0

	for _, batch := range batches ***REMOVED***
		for _, p := range batch ***REMOVED***
			if hasOnly && !b.Only[p.Stat.Name] ***REMOVED***
				continue
			***REMOVED***

			if b.Exclude[p.Stat.Name] ***REMOVED***
				continue
			***REMOVED***

			dimensions, ok := b.Data[p.Stat]
			if !ok ***REMOVED***
				dimensions = make(map[*string]*Dimension)
				b.Data[p.Stat] = dimensions
			***REMOVED***

			for dname, val := range p.Values ***REMOVED***
				interned, ok := b.interned[dname]
				if !ok ***REMOVED***
					interned = &dname
					b.interned[dname] = interned
				***REMOVED***

				dim, ok := dimensions[interned]
				if !ok ***REMOVED***
					dim = &Dimension***REMOVED******REMOVED***
					dimensions[interned] = dim
				***REMOVED***

				dim.Values = append(dim.Values, val)
				dim.Last = val
				dim.dirty = true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	for _, dimensions := range b.Data ***REMOVED***
		for _, dim := range dimensions ***REMOVED***
			if dim.dirty ***REMOVED***
				sort.Float64s(dim.Values)
				dim.dirty = false
			***REMOVED***
		***REMOVED***
	***REMOVED***

	b.submitMutex.Unlock()

	return nil
***REMOVED***
