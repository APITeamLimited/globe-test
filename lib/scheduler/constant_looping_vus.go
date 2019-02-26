package scheduler

import (
	"fmt"
	"time"

	"github.com/loadimpact/k6/lib/types"
	null "gopkg.in/guregu/null.v3"
)

const constantLoopingVUsType = "constant-looping-vus"

func init() ***REMOVED***
	RegisterConfigType(constantLoopingVUsType, func(name string, rawJSON []byte) (Config, error) ***REMOVED***
		config := NewConstantLoopingVUsConfig(name)
		err := strictJSONUnmarshal(rawJSON, &config)
		return config, err
	***REMOVED***)
***REMOVED***

// The minimum duration we'll allow users to schedule. This doesn't affect the stages
// configuration, where 0-duration virtual stages are allowed for instantaneous VU jumps
const minDuration = 1 * time.Second

// ConstantLoopingVUsConfig stores VUs and duration
type ConstantLoopingVUsConfig struct ***REMOVED***
	BaseConfig
	VUs      null.Int           `json:"vus"`
	Duration types.NullDuration `json:"duration"`
***REMOVED***

// NewConstantLoopingVUsConfig returns a ConstantLoopingVUsConfig with default values
func NewConstantLoopingVUsConfig(name string) ConstantLoopingVUsConfig ***REMOVED***
	return ConstantLoopingVUsConfig***REMOVED***BaseConfig: NewBaseConfig(name, constantLoopingVUsType, false)***REMOVED***
***REMOVED***

// Make sure we implement the Config interface
var _ Config = &ConstantLoopingVUsConfig***REMOVED******REMOVED***

// Validate makes sure all options are configured and valid
func (lcv ConstantLoopingVUsConfig) Validate() []error ***REMOVED***
	errors := lcv.BaseConfig.Validate()
	if !lcv.VUs.Valid ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of VUs isn't specified"))
	***REMOVED*** else if lcv.VUs.Int64 < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of VUs shouldn't be negative"))
	***REMOVED***

	if !lcv.Duration.Valid ***REMOVED***
		errors = append(errors, fmt.Errorf("the duration is unspecified"))
	***REMOVED*** else if time.Duration(lcv.Duration.Duration) < minDuration ***REMOVED***
		errors = append(errors, fmt.Errorf(
			"the duration should be at least %s, but is %s", minDuration, lcv.Duration,
		))
	***REMOVED***

	return errors
***REMOVED***

// GetMaxVUs returns the absolute maximum number of possible concurrently running VUs
func (lcv ConstantLoopingVUsConfig) GetMaxVUs() int64 ***REMOVED***
	return lcv.VUs.Int64
***REMOVED***

// GetMaxDuration returns the maximum duration time for this scheduler, including
// the specified iterationTimeout, if the iterations are uninterruptible
func (lcv ConstantLoopingVUsConfig) GetMaxDuration() time.Duration ***REMOVED***
	maxDuration := lcv.Duration.Duration
	if !lcv.Interruptible.Bool ***REMOVED***
		maxDuration += lcv.IterationTimeout.Duration
	***REMOVED***
	return time.Duration(maxDuration)
***REMOVED***

//TODO: figure out the most accurate algorithm for the proportional distribution
// of the VUs (which are indivisible items)... Some sort of "pick closest match
// to percentage and adjust remaining"?
