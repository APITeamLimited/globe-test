package metrics

import (
	"encoding/json"
	"io"
	"time"
)

// MarshalJSON returns a byte slice containing a JSON representation of all
// the metrics in the Registry.
func (r *StandardRegistry) MarshalJSON() ([]byte, error) ***REMOVED***
	return json.Marshal(r.GetAll())
***REMOVED***

// WriteJSON writes metrics from the given registry  periodically to the
// specified io.Writer as JSON.
func WriteJSON(r Registry, d time.Duration, w io.Writer) ***REMOVED***
	for _ = range time.Tick(d) ***REMOVED***
		WriteJSONOnce(r, w)
	***REMOVED***
***REMOVED***

// WriteJSONOnce writes metrics from the given registry to the specified
// io.Writer as JSON.
func WriteJSONOnce(r Registry, w io.Writer) ***REMOVED***
	json.NewEncoder(w).Encode(r)
***REMOVED***

func (p *PrefixedRegistry) MarshalJSON() ([]byte, error) ***REMOVED***
	return json.Marshal(p.GetAll())
***REMOVED***
