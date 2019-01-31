package scheduler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/loadimpact/k6/lib/types"
	null "gopkg.in/guregu/null.v3"
)

const perVUIterationsType = "per-vu-iterations"

func init() ***REMOVED***
	RegisterConfigType(perVUIterationsType, func(name string, rawJSON []byte) (Config, error) ***REMOVED***
		config := NewPerVUIterationsConfig(name)
		err := json.Unmarshal(rawJSON, &config)
		return config, err
	***REMOVED***)
***REMOVED***

// PerVUIteationsConfig stores the number of VUs iterations, as well as maxDuration settings
type PerVUIteationsConfig struct ***REMOVED***
	BaseConfig
	VUs         null.Int           `json:"vus"`
	Iterations  null.Int           `json:"iterations"`
	MaxDuration types.NullDuration `json:"maxDuration"`
***REMOVED***

// NewPerVUIterationsConfig returns a PerVUIteationsConfig with default values
func NewPerVUIterationsConfig(name string) PerVUIteationsConfig ***REMOVED***
	return PerVUIteationsConfig***REMOVED***
		BaseConfig:  NewBaseConfig(name, perVUIterationsType, false),
		MaxDuration: types.NewNullDuration(1*time.Hour, false),
	***REMOVED***
***REMOVED***

// Make sure we implement the Config interface
var _ Config = &PerVUIteationsConfig***REMOVED******REMOVED***

// Validate makes sure all options are configured and valid
func (pvic PerVUIteationsConfig) Validate() []error ***REMOVED***
	errors := pvic.BaseConfig.Validate()
	if !pvic.VUs.Valid ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of VUs isn't specified"))
	***REMOVED*** else if pvic.VUs.Int64 <= 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of VUs should be more than 0"))
	***REMOVED***

	if !pvic.Iterations.Valid ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of iterations isn't specified"))
	***REMOVED*** else if pvic.Iterations.Int64 <= 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of iterations should be more than 0"))
	***REMOVED***

	if time.Duration(pvic.MaxDuration.Duration) < minDuration ***REMOVED***
		errors = append(errors, fmt.Errorf(
			"the maxDuration should be at least %s, but is %s", minDuration, pvic.MaxDuration,
		))
	***REMOVED***

	return errors
***REMOVED***
