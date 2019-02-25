package scheduler

import (
	"fmt"
	"time"

	"github.com/loadimpact/k6/lib/types"
	null "gopkg.in/guregu/null.v3"
)

const variableArrivalRateType = "variable-arrival-rate"

func init() ***REMOVED***
	RegisterConfigType(variableArrivalRateType, func(name string, rawJSON []byte) (Config, error) ***REMOVED***
		config := NewVariableArrivalRateConfig(name)
		err := strictJSONUnmarshal(rawJSON, &config)
		return config, err
	***REMOVED***)
***REMOVED***

// VariableArrivalRateConfig stores config for the variable arrival-rate scheduler
type VariableArrivalRateConfig struct ***REMOVED***
	BaseConfig
	StartRate null.Int           `json:"startRate"`
	TimeUnit  types.NullDuration `json:"timeUnit"`
	Stages    []Stage            `json:"stages"`

	// Initialize `PreAllocatedVUs` number of VUs, and if more than that are needed,
	// they will be dynamically allocated, until `MaxVUs` is reached, which is an
	// absolutely hard limit on the number of VUs the scheduler will use
	PreAllocatedVUs null.Int `json:"preAllocatedVUs"`
	MaxVUs          null.Int `json:"maxVUs"`
***REMOVED***

// NewVariableArrivalRateConfig returns a VariableArrivalRateConfig with default values
func NewVariableArrivalRateConfig(name string) VariableArrivalRateConfig ***REMOVED***
	return VariableArrivalRateConfig***REMOVED***
		BaseConfig: NewBaseConfig(name, variableArrivalRateType, false),
		TimeUnit:   types.NewNullDuration(1*time.Second, false),
	***REMOVED***
***REMOVED***

// Make sure we implement the Config interface
var _ Config = &VariableArrivalRateConfig***REMOVED******REMOVED***

// Validate makes sure all options are configured and valid
func (varc VariableArrivalRateConfig) Validate() []error ***REMOVED***
	errors := varc.BaseConfig.Validate()

	if varc.StartRate.Int64 < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the startRate value shouldn't be negative"))
	***REMOVED***

	if time.Duration(varc.TimeUnit.Duration) < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the timeUnit should be more than 0"))
	***REMOVED***

	errors = append(errors, validateStages(varc.Stages)...)

	if !varc.PreAllocatedVUs.Valid ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of preAllocatedVUs isn't specified"))
	***REMOVED*** else if varc.PreAllocatedVUs.Int64 < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of preAllocatedVUs shouldn't be negative"))
	***REMOVED***

	if !varc.MaxVUs.Valid ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of maxVUs isn't specified"))
	***REMOVED*** else if varc.MaxVUs.Int64 < varc.PreAllocatedVUs.Int64 ***REMOVED***
		errors = append(errors, fmt.Errorf("maxVUs shouldn't be less than preAllocatedVUs"))
	***REMOVED***

	return errors
***REMOVED***

// GetMaxVUs returns the absolute maximum number of possible concurrently running VUs
func (varc VariableArrivalRateConfig) GetMaxVUs() int64 ***REMOVED***
	return varc.MaxVUs.Int64
***REMOVED***

// GetMaxDuration returns the maximum duration time for this scheduler, including
// the specified iterationTimeout, if the iterations are uninterruptible
func (varc VariableArrivalRateConfig) GetMaxDuration() time.Duration ***REMOVED***
	var maxDuration types.Duration
	for _, s := range varc.Stages ***REMOVED***
		maxDuration += s.Duration.Duration
	***REMOVED***
	if !varc.Interruptible.Bool ***REMOVED***
		maxDuration += varc.IterationTimeout.Duration
	***REMOVED***
	return time.Duration(maxDuration)
***REMOVED***
