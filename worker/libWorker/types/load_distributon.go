package types

import (
	"bytes"
	"encoding/json"
)

const (
	HTTPSingleLoadZone   = "httpSingle"
	HTTPMultipleLoadZone = "httpMultiple"
)

type LoadZone struct ***REMOVED***
	Location string `json:"location"`
	Fraction int    `json:"fraction"`
***REMOVED***

type NullLoadDistribution struct ***REMOVED***
	Value []LoadZone
	Valid bool
***REMOVED***

func NewNullLoadDistribution(loadZone []LoadZone, valid bool) NullLoadDistribution ***REMOVED***
	return NullLoadDistribution***REMOVED***loadZone, valid***REMOVED***
***REMOVED***

func NullLoadDistributionFrom(loadZone []LoadZone) NullLoadDistribution ***REMOVED***
	return NullLoadDistribution***REMOVED***loadZone, true***REMOVED***
***REMOVED***

func (em NullLoadDistribution) MarshalJSON() ([]byte, error) ***REMOVED***
	if !em.Valid ***REMOVED***
		return []byte(`null`), nil
	***REMOVED***
	return json.Marshal(em.Value)
***REMOVED***

func (em *NullLoadDistribution) UnmarshalJSON(data []byte) error ***REMOVED***
	if bytes.Equal(data, []byte(`null`)) ***REMOVED***
		em.Valid = false
		return nil
	***REMOVED***
	if err := json.Unmarshal(data, &em.Value); err != nil ***REMOVED***
		return err
	***REMOVED***
	em.Valid = true
	return nil
***REMOVED***
