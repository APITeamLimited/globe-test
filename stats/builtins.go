package stats

import (
	"encoding/json"
	"io"
)

type JSONBackend struct ***REMOVED***
	encoder *json.Encoder
***REMOVED***

func NewJSONBackend(w io.Writer) Backend ***REMOVED***
	return &JSONBackend***REMOVED***encoder: json.NewEncoder(w)***REMOVED***
***REMOVED***

func (b *JSONBackend) Submit(batches [][]Sample) error ***REMOVED***
	for _, batch := range batches ***REMOVED***
		for _, s := range batch ***REMOVED***
			if err := b.encoder.Encode(b.format(&s)); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (JSONBackend) format(s *Sample) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	data := map[string]interface***REMOVED******REMOVED******REMOVED***
		"time":   s.Time,
		"stat":   s.Stat.Name,
		"tags":   s.Tags,
		"values": s.Values,
	***REMOVED***
	if s.Tags == nil ***REMOVED***
		data["tags"] = Tags***REMOVED******REMOVED***
	***REMOVED***
	return data
***REMOVED***
