package types

import (
	"bytes"
	"encoding/json"
)

const (
	HTTPSingleExecutionMode   = "httpSingle"
	HTTPMultipleExecutionMode = "httpMultiple"
)

type NullExecutionMode struct ***REMOVED***
	Value string
	Valid bool
***REMOVED***

func NewNullExecutionMode(s string, valid bool) NullExecutionMode ***REMOVED***
	return NullExecutionMode***REMOVED***s, valid***REMOVED***
***REMOVED***

func NullExecutionModeFrom(s string) NullExecutionMode ***REMOVED***
	return NullExecutionMode***REMOVED***s, true***REMOVED***
***REMOVED***

func (em NullExecutionMode) ValueOrZero() string ***REMOVED***
	if !em.Valid ***REMOVED***
		return ""
	***REMOVED***
	return em.Value
***REMOVED***

func (em NullExecutionMode) MarshalJSON() ([]byte, error) ***REMOVED***
	if !em.Valid ***REMOVED***
		return []byte(`null`), nil
	***REMOVED***
	return json.Marshal(em.Value)
***REMOVED***

func (em *NullExecutionMode) UnmarshalJSON(data []byte) error ***REMOVED***
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
