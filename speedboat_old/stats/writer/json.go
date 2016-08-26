package writer

import (
	"encoding/json"
)

type JSONFormatter struct***REMOVED******REMOVED***

func (JSONFormatter) Format(data interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return json.Marshal(data)
***REMOVED***

type PrettyJSONFormatter struct***REMOVED******REMOVED***

func (PrettyJSONFormatter) Format(data interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return json.MarshalIndent(data, "", "    ")
***REMOVED***
