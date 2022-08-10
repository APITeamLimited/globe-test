package v1

import (
	"go.k6.io/k6/core"
)

// StatusJSONAPI is JSON API envelop for metrics
type StatusJSONAPI struct ***REMOVED***
	Data statusData `json:"data"`
***REMOVED***

// NewStatusJSONAPI creates the JSON API status envelop
func NewStatusJSONAPI(s Status) StatusJSONAPI ***REMOVED***
	return StatusJSONAPI***REMOVED***
		Data: statusData***REMOVED***
			ID:         "default",
			Type:       "status",
			Attributes: s,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// Status extract the v1.Status from the JSON API envelop
func (s StatusJSONAPI) Status() Status ***REMOVED***
	return s.Data.Attributes
***REMOVED***

type statusData struct ***REMOVED***
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes Status `json:"attributes"`
***REMOVED***

func newStatusJSONAPIFromEngine(engine *core.Engine) StatusJSONAPI ***REMOVED***
	return NewStatusJSONAPI(NewStatus(engine))
***REMOVED***
