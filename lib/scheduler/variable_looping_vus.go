package scheduler

import (
	"fmt"
	"time"

	"github.com/loadimpact/k6/lib/types"
	null "gopkg.in/guregu/null.v3"
)

const variableLoopingVUsType = "variable-looping-vus"

func init() ***REMOVED***
	RegisterConfigType(variableLoopingVUsType, func(name string, rawJSON []byte) (Config, error) ***REMOVED***
		config := NewVariableLoopingVUsConfig(name)
		err := strictJSONUnmarshal(rawJSON, &config)
		return config, err
	***REMOVED***)
***REMOVED***

// Stage contains
type Stage struct ***REMOVED***
	Duration types.NullDuration `json:"duration"`
	Target   null.Int           `json:"target"` // TODO: maybe rename this to endVUs? something else?
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
func (vlvc VariableLoopingVUsConfig) Validate() []error ***REMOVED***
	errors := vlvc.BaseConfig.Validate()
	if vlvc.StartVUs.Int64 < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of start VUs shouldn't be negative"))
	***REMOVED***

	return append(errors, validateStages(vlvc.Stages)...)
***REMOVED***

// GetMaxVUs returns the absolute maximum number of possible concurrently running VUs
func (vlvc VariableLoopingVUsConfig) GetMaxVUs() int64 ***REMOVED***
	maxVUs := vlvc.StartVUs.Int64
	for _, s := range vlvc.Stages ***REMOVED***
		if s.Target.Int64 > maxVUs ***REMOVED***
			maxVUs = s.Target.Int64
		***REMOVED***
	***REMOVED***
	return maxVUs
***REMOVED***

// GetMaxDuration returns the maximum duration time for this scheduler, including
// the specified iterationTimeout, if the iterations are uninterruptible
func (vlvc VariableLoopingVUsConfig) GetMaxDuration() time.Duration ***REMOVED***
	var maxDuration types.Duration
	for _, s := range vlvc.Stages ***REMOVED***
		maxDuration += s.Duration.Duration
	***REMOVED***
	if !vlvc.Interruptible.Bool ***REMOVED***
		maxDuration += vlvc.IterationTimeout.Duration
	***REMOVED***
	return time.Duration(maxDuration)
***REMOVED***
