package writer

import (
	"gopkg.in/yaml.v2"
)

type YAMLFormatter struct***REMOVED******REMOVED***

func (YAMLFormatter) Format(data interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return yaml.Marshal(data)
***REMOVED***
