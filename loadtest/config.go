package loadtest

import (
	"gopkg.in/yaml.v2"
)

// Configuration type for a state.
type ConfigStage struct ***REMOVED***
	Duration string `yaml:"duration"`
	VUs      []int  `yaml:"vus"`
***REMOVED***

type Config struct ***REMOVED***
	Duration string        `yaml:"duration"`
	Script   string        `yaml:"script"`
	Stages   []ConfigStage `yaml:"stages"`
***REMOVED***

func NewConfig() Config ***REMOVED***
	return Config***REMOVED******REMOVED***
***REMOVED***

func ParseConfig(data []byte, conf *Config) (err error) ***REMOVED***
	return yaml.Unmarshal(data, conf)
***REMOVED***
