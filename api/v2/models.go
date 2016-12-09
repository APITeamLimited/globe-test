package v2

import (
	"gopkg.in/guregu/null.v3"
)

type Status struct ***REMOVED***
	Running null.Bool `json:"running"`
	Tainted null.Bool `json:"tainted"`
	VUs     null.Int  `json:"vus"`
	VUsMax  null.Int  `json:"vus-max"`
***REMOVED***

func (s Status) GetName() string ***REMOVED***
	return "status"
***REMOVED***

func (s Status) GetID() string ***REMOVED***
	return "default"
***REMOVED***

func (s Status) SetID(id string) error ***REMOVED***
	return nil
***REMOVED***
