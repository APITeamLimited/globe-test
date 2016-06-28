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

func (b *JSONBackend) Submit(batches [][]Point) error ***REMOVED***
	for _, batch := range batches ***REMOVED***
		for _, p := range batch ***REMOVED***
			if err := b.encoder.Encode(b.format(&p)); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (JSONBackend) format(p *Point) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	data := map[string]interface***REMOVED******REMOVED******REMOVED***
		"time":   p.Time,
		"stat":   p.Stat.Name,
		"tags":   p.Tags,
		"values": p.Values,
	***REMOVED***
	if p.Tags == nil ***REMOVED***
		data["tags"] = Tags***REMOVED******REMOVED***
	***REMOVED***
	return data
***REMOVED***
