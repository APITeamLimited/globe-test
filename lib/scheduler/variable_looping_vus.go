package scheduler

import (
	"encoding/json"
	"fmt"

	"github.com/loadimpact/k6/lib/types"
	null "gopkg.in/guregu/null.v3"
)

const variableLoopingVUsType = "variable-looping-vus"

func init() ***REMOVED***
	RegisterConfigType(variableLoopingVUsType, func(name string, rawJSON []byte) (Config, error) ***REMOVED***
		config := NewVariableLoopingVUsConfig(name)
		err := json.Unmarshal(rawJSON, &config)
		return config, err
	***REMOVED***)
***REMOVED***

// Stage contains
type Stage struct ***REMOVED***
	Duration types.NullDuration `json:"duration"`
	Target   null.Int           `json:"target"` // TODO: maybe rename this to endVUs?
***REMOVED***

// VariableLoopingVUsConfig stores the configuration for the stages scheduler
type VariableLoopingVUsConfig struct ***REMOVED***
	BaseConfig
	StartVUs null.Int `json:"startVUs"`
	Stages   []Stage  `json:"stages"`
***REMOVED***

// NewVariableLoopingVUsConfig returns a VariableLoopingVUsConfig with its default values
func NewVariableLoopingVUsConfig(name string) VariableLoopingVUsConfig ***REMOVED***
	return VariableLoopingVUsConfig***REMOVED***BaseConfig: NewBaseConfig(name, variableLoopingVUsType, false)***REMOVED***
***REMOVED***

// Make sure we implement the Config interface
var _ Config = &VariableLoopingVUsConfig***REMOVED******REMOVED***

// Validate makes sure all options are configured and valid
func (ls VariableLoopingVUsConfig) Validate() []error ***REMOVED***
	errors := ls.BaseConfig.Validate()
	if ls.StartVUs.Int64 < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of start VUs shouldn't be negative"))
	***REMOVED***

	if len(ls.Stages) == 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("at least one stage has to be specified"))
	***REMOVED*** else ***REMOVED***
		for i, s := range ls.Stages ***REMOVED***
			stageNum := i + 1
			if !s.Duration.Valid ***REMOVED***
				errors = append(errors, fmt.Errorf("stage %d doesn't have a duration", stageNum))
			***REMOVED*** else if s.Duration.Duration < 0 ***REMOVED***
				errors = append(errors, fmt.Errorf("the duration for stage %d shouldn't be negative", stageNum))
			***REMOVED***
			if !s.Target.Valid ***REMOVED***
				errors = append(errors, fmt.Errorf("stage %d doesn't have a target", stageNum))
			***REMOVED*** else if s.Target.Int64 < 0 ***REMOVED***
				errors = append(errors, fmt.Errorf("the target for stage %d shouldn't be negative", stageNum))
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return errors
***REMOVED***
