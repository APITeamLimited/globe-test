package lib

import (
	"gopkg.in/guregu/null.v3"
)

type Status struct ***REMOVED***
	Running     null.Bool `json:"running" yaml:"Running"`
	ActiveVUs   null.Int  `json:"active-vus" yaml:"ActiveVUs"`
	InactiveVUs null.Int  `json:"inactive-vus" yaml:"InactiveVUs"`
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

type Info struct ***REMOVED***
	ID      string `json:"-"`
	Version string `json:"version"`
***REMOVED***

func (i Info) GetName() string ***REMOVED***
	return "info"
***REMOVED***

func (i Info) GetID() string ***REMOVED***
	return "default"
***REMOVED***
