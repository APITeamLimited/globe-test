package scheduler

import (
	"fmt"
	"time"

	"github.com/loadimpact/k6/lib/types"
	null "gopkg.in/guregu/null.v3"
)

const sharedIterationsType = "shared-iterations"

func init() ***REMOVED***
	RegisterConfigType(sharedIterationsType, func(name string, rawJSON []byte) (Config, error) ***REMOVED***
		config := NewSharedIterationsConfig(name)
		err := strictJSONUnmarshal(rawJSON, &config)
		return config, err
	***REMOVED***)
***REMOVED***

// SharedIteationsConfig stores the number of VUs iterations, as well as maxDuration settings
type SharedIteationsConfig struct ***REMOVED***
	BaseConfig
	VUs         null.Int           `json:"vus"`
	Iterations  null.Int           `json:"iterations"`
	MaxDuration types.NullDuration `json:"maxDuration"`
***REMOVED***

// NewSharedIterationsConfig returns a SharedIteationsConfig with default values
func NewSharedIterationsConfig(name string) SharedIteationsConfig ***REMOVED***
	return SharedIteationsConfig***REMOVED***
		BaseConfig:  NewBaseConfig(name, sharedIterationsType, false),
		VUs:         null.NewInt(1, false), //TODO: set defaults in another place?
		Iterations:  null.NewInt(1, false),
		MaxDuration: types.NewNullDuration(1*time.Hour, false),
	***REMOVED***
***REMOVED***

// Make sure we implement the Config interface
var _ Config = &SharedIteationsConfig***REMOVED******REMOVED***

// Validate makes sure all options are configured and valid
func (sic SharedIteationsConfig) Validate() []error ***REMOVED***
	errors := sic.BaseConfig.Validate()
	if sic.VUs.Int64 <= 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of VUs should be more than 0"))
	***REMOVED***

	if sic.Iterations.Int64 < sic.VUs.Int64 ***REMOVED***
		errors = append(errors, fmt.Errorf(
			"the number of iterations (%d) shouldn't be less than the number of VUs (%d)",
			sic.Iterations.Int64, sic.VUs.Int64,
		))
	***REMOVED***

	if time.Duration(sic.MaxDuration.Duration) < minDuration ***REMOVED***
		errors = append(errors, fmt.Errorf(
			"the maxDuration should be at least %s, but is %s", minDuration, sic.MaxDuration,
		))
	***REMOVED***

	return errors
***REMOVED***

// GetMaxVUs returns the absolute maximum number of possible concurrently running VUs
func (sic SharedIteationsConfig) GetMaxVUs() int64 ***REMOVED***
	return sic.VUs.Int64
***REMOVED***

// GetMaxDuration returns the maximum duration time for this scheduler, including
// the specified iterationTimeout, if the iterations are uninterruptible
func (sic SharedIteationsConfig) GetMaxDuration() time.Duration ***REMOVED***
	maxDuration := sic.MaxDuration.Duration
	if !sic.Interruptible.Bool ***REMOVED***
		maxDuration += sic.IterationTimeout.Duration
	***REMOVED***
	return time.Duration(maxDuration)
***REMOVED***
