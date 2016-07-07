package accumulate

import (
	"fmt"
	"github.com/loadimpact/speedboat/stats"
	"sort"
	"strings"
	"sync"
)

type StatTree map[StatTreeKey]*StatTreeNode

type StatTreeKey struct ***REMOVED***
	Tag   string
	Value interface***REMOVED******REMOVED***
***REMOVED***

type StatTreeNode struct ***REMOVED***
	Stat     *stats.Stat
	Substats *StatTree
***REMOVED***

type Backend struct ***REMOVED***
	Data    map[*stats.Stat]map[*string]*Dimension
	Only    map[string]bool
	Exclude map[string]bool
	GroupBy []string

	interned    map[string]*string
	vstats      map[*stats.Stat]*StatTree
	submitMutex sync.Mutex
***REMOVED***

func New() *Backend ***REMOVED***
	return &Backend***REMOVED***
		Data:     make(map[*stats.Stat]map[*string]*Dimension),
		Exclude:  make(map[string]bool),
		Only:     make(map[string]bool),
		interned: make(map[string]*string),
		vstats:   make(map[*stats.Stat]*StatTree),
	***REMOVED***
***REMOVED***

func (b *Backend) getVStat(stat *stats.Stat, tags stats.Tags) *stats.Stat ***REMOVED***
	tree := b.vstats[stat]
	if tree == nil ***REMOVED***
		tmp := make(StatTree)
		tree = &tmp
		b.vstats[stat] = tree
	***REMOVED***

	ret := stat
	for n, tag := range b.GroupBy ***REMOVED***
		val, ok := tags[tag]
		if !ok ***REMOVED***
			continue
		***REMOVED***

		key := StatTreeKey***REMOVED***Tag: tag, Value: val***REMOVED***
		node := (*tree)[key]
		if node == nil ***REMOVED***
			tagStrings := make([]string, 0, n)
			for i := 0; i <= n; i++ ***REMOVED***
				t := b.GroupBy[i]
				v, ok := tags[t]
				if !ok ***REMOVED***
					continue
				***REMOVED***
				tagStrings = append(tagStrings, fmt.Sprintf("%s: %v", t, v))
			***REMOVED***

			name := stat.Name
			if len(tagStrings) > 0 ***REMOVED***
				name = fmt.Sprintf("%s***REMOVED***%s***REMOVED***", name, strings.Join(tagStrings, ", "))
			***REMOVED***

			substats := make(StatTree)
			node = &StatTreeNode***REMOVED***
				Stat: &stats.Stat***REMOVED***
					Name:   name,
					Type:   stat.Type,
					Intent: stat.Intent,
				***REMOVED***,
				Substats: &substats,
			***REMOVED***
			(*tree)[key] = node
		***REMOVED***

		ret = node.Stat
	***REMOVED***

	return ret
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

			stat := b.getVStat(p.Stat, p.Tags)
			dimensions, ok := b.Data[stat]
			if !ok ***REMOVED***
				dimensions = make(map[*string]*Dimension)
				b.Data[stat] = dimensions
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
