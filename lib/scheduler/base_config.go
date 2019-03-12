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

const minPercentage = 0.01

// The maximum time k6 will wait after an iteration is supposed to be done
const maxIterationTimeout = 300 * time.Second

// BaseConfig contains the common config fields for all schedulers
type BaseConfig struct ***REMOVED***
	Name             string             `json:"-"` // set via the JS object key
	Type             string             `json:"type"`
	StartTime        types.NullDuration `json:"startTime"`
	Interruptible    null.Bool          `json:"interruptible"`
	IterationTimeout types.NullDuration `json:"iterationTimeout"`
	Env              map[string]string  `json:"env"`
	Exec             null.String        `json:"exec"` // function name, externally validated
	Percentage       float64            `json:"-"`    // 100, unless Split() was called

	//TODO: future extensions like tags, distribution, others?
***REMOVED***

// NewBaseConfig returns a default base config with the default values
func NewBaseConfig(name, configType string, interruptible bool) BaseConfig ***REMOVED***
	return BaseConfig***REMOVED***
		Name:             name,
		Type:             configType,
		Interruptible:    null.NewBool(interruptible, false),
		IterationTimeout: types.NewNullDuration(30*time.Second, false),
		Percentage:       100,
	***REMOVED***
***REMOVED***

// Validate checks some basic things like present name, type, and a positive start time
func (bc BaseConfig) Validate() (errors []error) ***REMOVED***
	// Some just-in-case checks, since those things are likely checked in other places or
	// even assigned by us:
	if bc.Name == "" ***REMOVED***
		errors = append(errors, fmt.Errorf("scheduler name shouldn't be empty"))
	***REMOVED***
	if bc.Type == "" ***REMOVED***
		errors = append(errors, fmt.Errorf("missing or empty type field"))
	***REMOVED***
	if bc.Percentage < minPercentage || bc.Percentage > 100 ***REMOVED***
		errors = append(errors, fmt.Errorf(
			"percentage should be between %f and 100, but is %f", minPercentage, bc.Percentage,
		))
	***REMOVED***
	if bc.Exec.Valid && bc.Exec.String == "" ***REMOVED***
		errors = append(errors, fmt.Errorf("exec value cannot be empty"))
	***REMOVED***
	// The actually reasonable checks:
	if bc.StartTime.Duration < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("scheduler start time can't be negative"))
	***REMOVED***
	iterTimeout := time.Duration(bc.IterationTimeout.Duration)
	if iterTimeout < 0 || iterTimeout > maxIterationTimeout ***REMOVED***
		errors = append(errors, fmt.Errorf(
			"the iteration timeout should be between 0 and %s, but is %s", maxIterationTimeout, iterTimeout,
		))
	***REMOVED***
	return errors
***REMOVED***

// GetBaseConfig just returns itself
func (bc BaseConfig) GetBaseConfig() BaseConfig ***REMOVED***
	return bc
***REMOVED***

// CopyWithPercentage is a helper function that just sets the percentage to
// the specified amount.
func (bc BaseConfig) CopyWithPercentage(percentage float64) *BaseConfig ***REMOVED***
	c := bc
	c.Percentage = percentage
	return &c
***REMOVED***
