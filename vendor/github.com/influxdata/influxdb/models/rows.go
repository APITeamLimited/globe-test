package models

import (
	"sort"
)

// Row represents a single row returned from the execution of a statement.
type Row struct ***REMOVED***
	Name    string            `json:"name,omitempty"`
	Tags    map[string]string `json:"tags,omitempty"`
	Columns []string          `json:"columns,omitempty"`
	Values  [][]interface***REMOVED******REMOVED***   `json:"values,omitempty"`
	Partial bool              `json:"partial,omitempty"`
***REMOVED***

// SameSeries returns true if r contains values for the same series as o.
func (r *Row) SameSeries(o *Row) bool ***REMOVED***
	return r.tagsHash() == o.tagsHash() && r.Name == o.Name
***REMOVED***

// tagsHash returns a hash of tag key/value pairs.
func (r *Row) tagsHash() uint64 ***REMOVED***
	h := NewInlineFNV64a()
	keys := r.tagsKeys()
	for _, k := range keys ***REMOVED***
		h.Write([]byte(k))
		h.Write([]byte(r.Tags[k]))
	***REMOVED***
	return h.Sum64()
***REMOVED***

// tagKeys returns a sorted list of tag keys.
func (r *Row) tagsKeys() []string ***REMOVED***
	a := make([]string, 0, len(r.Tags))
	for k := range r.Tags ***REMOVED***
		a = append(a, k)
	***REMOVED***
	sort.Strings(a)
	return a
***REMOVED***

// Rows represents a collection of rows. Rows implements sort.Interface.
type Rows []*Row

// Len implements sort.Interface.
func (p Rows) Len() int ***REMOVED*** return len(p) ***REMOVED***

// Less implements sort.Interface.
func (p Rows) Less(i, j int) bool ***REMOVED***
	// Sort by name first.
	if p[i].Name != p[j].Name ***REMOVED***
		return p[i].Name < p[j].Name
	***REMOVED***

	// Sort by tag set hash. Tags don't have a meaningful sort order so we
	// just compute a hash and sort by that instead. This allows the tests
	// to receive rows in a predictable order every time.
	return p[i].tagsHash() < p[j].tagsHash()
***REMOVED***

// Swap implements sort.Interface.
func (p Rows) Swap(i, j int) ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***
