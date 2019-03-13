/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package scheduler

import (
	"fmt"
	"time"

	"github.com/loadimpact/k6/lib/types"
	null "gopkg.in/guregu/null.v3"
)

const constantArrivalRateType = "constant-arrival-rate"

func init() ***REMOVED***
	RegisterConfigType(constantArrivalRateType, func(name string, rawJSON []byte) (Config, error) ***REMOVED***
		config := NewConstantArrivalRateConfig(name)
		err := strictJSONUnmarshal(rawJSON, &config)
		return config, err
	***REMOVED***)
***REMOVED***

// ConstantArrivalRateConfig stores config for the constant arrival-rate scheduler
type ConstantArrivalRateConfig struct ***REMOVED***
	BaseConfig
	Rate     null.Int           `json:"rate"`
	TimeUnit types.NullDuration `json:"timeUnit"`
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
		errors = append(errors, fmt.Errorf("the iteration rate should be more than 0"))
	***REMOVED***

	if time.Duration(carc.TimeUnit.Duration) <= 0 ***REMOVED***
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

// GetMaxVUs returns the absolute maximum number of possible concurrently running VUs
func (carc ConstantArrivalRateConfig) GetMaxVUs() int64 ***REMOVED***
	return carc.MaxVUs.Int64
***REMOVED***

// GetMaxDuration returns the maximum duration time for this scheduler, including
// the specified iterationTimeout, if the iterations are uninterruptible
func (carc ConstantArrivalRateConfig) GetMaxDuration() time.Duration ***REMOVED***
	maxDuration := carc.Duration.Duration
	if !carc.Interruptible.Bool ***REMOVED***
		maxDuration += carc.IterationTimeout.Duration
	***REMOVED***
	return time.Duration(maxDuration)
***REMOVED***
