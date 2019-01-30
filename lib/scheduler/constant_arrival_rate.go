package scheduler

import (
	"fmt"
	"time"

	"github.com/loadimpact/k6/lib/types"
	null "gopkg.in/guregu/null.v3"
)

const constantArrivalRateType = "constant-arrival-rate"

// ConstantArrivalRateConfig stores config for the constant arrival-rate scheduler
type ConstantArrivalRateConfig struct ***REMOVED***
	BaseConfig
	Rate     null.Int           `json:"rate"`
	TimeUnit types.NullDuration `json:"timeUnit"` //TODO: rename to something else?
	Duration types.NullDuration `json:"duration"`

	// Initialize `PreAllocatedVUs` number of VUs, and if more than that are needed,
	// they will be dynamically allocated, until `MaxVUs` is reached, which is an
	// absolutely hard limit on the number of VUs the scheduler will use
	PreAllocatedVUs null.Int `json:"preAllocatedVUs"`
	MaxVUs          null.Int `json:"maxVUs"`
***REMOVED***

// NewConstantArrivalRateConfig returns a ConstantArrivalRateConfig with default values
func NewConstantArrivalRateConfig(name string) ConstantArrivalRateConfig ***REMOVED***
	return ConstantArrivalRateConfig***REMOVED***
		BaseConfig: NewBaseConfig(name, constantArrivalRateType, false),
		TimeUnit:   types.NewNullDuration(1*time.Second, false),
		//TODO: set some default values for PreAllocatedVUs and MaxVUs?
	***REMOVED***
***REMOVED***

// Make sure we implement the Config interface
var _ Config = &ConstantArrivalRateConfig***REMOVED******REMOVED***

// Validate makes sure all options are configured and valid
func (carc ConstantArrivalRateConfig) Validate() []error ***REMOVED***
	errors := carc.BaseConfig.Validate()
	if !carc.Rate.Valid ***REMOVED***
		errors = append(errors, fmt.Errorf("the iteration rate isn't specified"))
	***REMOVED*** else if carc.Rate.Int64 <= 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the iteration rate should be positive"))
	***REMOVED***

	if time.Duration(carc.TimeUnit.Duration) < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the timeUnit should be more than 0"))
	***REMOVED***

	if !carc.Duration.Valid ***REMOVED***
		errors = append(errors, fmt.Errorf("the duration is unspecified"))
	***REMOVED*** else if time.Duration(carc.Duration.Duration) < minDuration ***REMOVED***
		errors = append(errors, fmt.Errorf(
			"the duration should be at least %s, but is %s", minDuration, carc.Duration,
		))
	***REMOVED***

	if !carc.PreAllocatedVUs.Valid ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of preAllocatedVUs isn't specified"))
	***REMOVED*** else if carc.PreAllocatedVUs.Int64 < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of preAllocatedVUs shouldn't be negative"))
	***REMOVED***

	if !carc.MaxVUs.Valid ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of maxVUs isn't specified"))
	***REMOVED*** else if carc.MaxVUs.Int64 < carc.PreAllocatedVUs.Int64 ***REMOVED***
		errors = append(errors, fmt.Errorf("maxVUs shouldn't be less than preAllocatedVUs"))
	***REMOVED***

	return errors
***REMOVED***
