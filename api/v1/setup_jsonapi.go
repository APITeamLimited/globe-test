package v1

type setUpJSONAPI struct ***REMOVED***
	Data setUpData `json:"data"`
***REMOVED***

type setUpData struct ***REMOVED***
	Type       string      `json:"type"`
	ID         string      `json:"id"`
	Attributes interface***REMOVED******REMOVED*** `json:"attributes"`
***REMOVED***

func newSetUpJSONAPI(setup interface***REMOVED******REMOVED***) setUpJSONAPI ***REMOVED***
	return setUpJSONAPI***REMOVED***
		Data: setUpData***REMOVED***
			Type:       "setupData",
			ID:         "default",
			Attributes: setup,
		***REMOVED***,
	***REMOVED***
***REMOVED***
