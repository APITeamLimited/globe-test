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

package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/ui/pb"
	"github.com/sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"
)

//TODO: remove globals and use some type of explicit dependency injection?
//nolint:gochecknoglobals
var (
	schedulerConfigTypesMutex   sync.RWMutex
	schedulerConfigConstructors = make(map[string]SchedulerConfigConstructor)
)

// ExecutionStep is used by different schedulers to specify the planned number
// of VUs they will need at a particular time. The times are relative to their
// StartTime, i.e. they don't take into account the specific starting time of
// the scheduler, as that will be considered by the external executor
// separately.
//
// A slice [***REMOVED***t1, v1***REMOVED***, ***REMOVED***t2, v2***REMOVED***, ***REMOVED***t3, v3***REMOVED***, ..., ***REMOVED***tn, vn***REMOVED***] of execution steps
// means that a scheduler will need 0 VUs until t1, it will need v1 number of
// VUs from time t1 until t2, need v2 number of VUs from time t2 to t3, and so
// on. t1 is usually 0, tn is usually the same as GetMaxDuration() and vn is
// usually 0.
//
// Keep in mind that t(i) may be exactly equal to t(i+i), when there's an abrupt
// transition in the number of VUs required by a scheduler. For example, the
// variable-looping-vus scheduler may have 0-duration stages, or it may scale up
// VUs in its last stage right until the end. These
//
// []ExecutionStep is also used by the SchedulerConfigMap, to represent the
// amount of needed VUs among all schedulers, during the whole execution of a
// test script. In that context, each scheduler's StartTime is accounted for and
// included in the offsets.
type ExecutionStep struct ***REMOVED***
	TimeOffset      time.Duration
	PlannedVUs      uint64
	MaxUnplannedVUs uint64
***REMOVED***

//TODO: make []ExecutionStep or []SchedulerConfig their own type?

// SchedulerConfig is an interface that should be implemented by all scheduler config types
type SchedulerConfig interface ***REMOVED***
	Validate() []error

	GetName() string
	GetType() string
	GetStartTime() time.Duration
	GetGracefulStop() time.Duration

	// This is used to validate whether a particular script can run in the cloud
	// or, in the future, in the native k6 distributed execution. Currently only
	// the manual-execution scheduler should return false.
	IsDistributable() bool

	GetEnv() map[string]string
	GetExec() null.String //TODO: use interface***REMOVED******REMOVED*** so plain http requests can be specified?

	// Calculates the VU requirements in different stages of the scheduler's
	// execution, including any extensions caused by waiting for iterations to
	// finish with graceful stops or ramp-downs.
	GetExecutionRequirements(*ExecutionSegment) []ExecutionStep

	// Return a human-readable description of the scheduler
	GetDescription(es *ExecutionSegment) string

	NewScheduler(*ExecutorState, *logrus.Entry) (Scheduler, error)
***REMOVED***

// InitVUFunc is just a shorthand so we don't have to type the function
// signature every time.
type InitVUFunc func(context.Context, *logrus.Entry) (VU, error)

// Scheduler is the interface all schedulers should implement
type Scheduler interface ***REMOVED***
	GetConfig() SchedulerConfig
	GetProgress() *pb.ProgressBar
	GetLogger() *logrus.Entry

	Init(ctx context.Context) error
	Run(ctx context.Context, engineOut chan<- stats.SampleContainer) error

	// Currently used only in the manual executor
	IsPausable() bool // TODO: make a separate interfaces
	LiveUpdate(pause bool, newConfig SchedulerConfig) error
***REMOVED***

// SchedulerConfigConstructor is a simple function that returns a concrete
// Config instance with the specified name and all default values correctly
// initialized
type SchedulerConfigConstructor func(name string, rawJSON []byte) (SchedulerConfig, error)

// RegisterSchedulerConfigType adds the supplied SchedulerConfigConstructor as
// the constructor for its type in the configConstructors map, in a thread-safe
// manner
func RegisterSchedulerConfigType(configType string, constructor SchedulerConfigConstructor) ***REMOVED***
	schedulerConfigTypesMutex.Lock()
	defer schedulerConfigTypesMutex.Unlock()

	if constructor == nil ***REMOVED***
		panic("scheduler configs: constructor is nil")
	***REMOVED***
	if _, configTypeExists := schedulerConfigConstructors[configType]; configTypeExists ***REMOVED***
		panic("scheduler configs: lib.RegisterSchedulerConfigType called twice for  " + configType)
	***REMOVED***

	schedulerConfigConstructors[configType] = constructor
***REMOVED***

// SchedulerConfigMap can contain mixed scheduler config types
type SchedulerConfigMap map[string]SchedulerConfig

// UnmarshalJSON implements the json.Unmarshaler interface in a two-step manner,
// creating the correct type of configs based on the `type` property.
func (scs *SchedulerConfigMap) UnmarshalJSON(data []byte) error ***REMOVED***
	if len(data) == 0 ***REMOVED***
		return nil
	***REMOVED***

	if len(data) == 4 && string(data) == "null" ***REMOVED***
		return nil
	***REMOVED***

	//TODO: use a more sophisticated combination of dec.Token() and dec.More(),
	// which would allow us to support both arrays and maps for this config?
	var protoConfigs map[string]protoSchedulerConfig
	if err := StrictJSONUnmarshal(data, &protoConfigs); err != nil ***REMOVED***
		return err
	***REMOVED***

	result := make(SchedulerConfigMap, len(protoConfigs))
	for k, v := range protoConfigs ***REMOVED***
		if v.configType == "" ***REMOVED***
			return fmt.Errorf("execution config '%s' doesn't have a type value", k)
		***REMOVED***
		config, err := GetParsedSchedulerConfig(k, v.configType, v.rawJSON)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		result[k] = config
	***REMOVED***

	*scs = result

	return nil
***REMOVED***

// Validate checks if all of the specified scheduler options make sense
func (scs SchedulerConfigMap) Validate() (errors []error) ***REMOVED***
	for name, scheduler := range scs ***REMOVED***
		if schedErr := scheduler.Validate(); len(schedErr) != 0 ***REMOVED***
			errors = append(errors,
				fmt.Errorf("scheduler %s has errors: %s", name, concatErrors(schedErr, ", ")))
		***REMOVED***
	***REMOVED***
	return errors
***REMOVED***

// GetSortedSchedulerConfigs returns a slice with the scheduler configurations,
// sorted in a consistent and predictable manner. It is useful when we want or
// have to avoid using maps with string keys (and tons of string lookups in
// them) and avoid the unpredictable iterations over Go maps. Slices allow us
// constant-time lookups and ordered iterations.
//
// The configs in the returned slice will be sorted by their start times in an
// ascending order, and alphabetically by their names (which are unique) if
// there are ties.
func (scs SchedulerConfigMap) GetSortedSchedulerConfigs() []SchedulerConfig ***REMOVED***
	configs := make([]SchedulerConfig, len(scs))

	// Populate the configs slice with sorted scheduler configs
	i := 0
	for _, config := range scs ***REMOVED***
		configs[i] = config // populate the slice in an unordered manner
		i++
	***REMOVED***
	sort.Slice(configs, func(a, b int) bool ***REMOVED*** // sort by (start time, name)
		switch ***REMOVED***
		case configs[a].GetStartTime() < configs[b].GetStartTime():
			return true
		case configs[a].GetStartTime() == configs[b].GetStartTime():
			return strings.Compare(configs[a].GetName(), configs[b].GetName()) < 0
		default:
			return false
		***REMOVED***
	***REMOVED***)

	return configs
***REMOVED***

// GetFullExecutionRequirements combines the execution requirements from all of
// the configured schedulers. It takes into account their start times and their
// individual VU requirements and calculates the total VU requirements for each
// moment in the test execution.
func (scs SchedulerConfigMap) GetFullExecutionRequirements(executionSegment *ExecutionSegment) []ExecutionStep ***REMOVED***
	sortedConfigs := scs.GetSortedSchedulerConfigs()

	// Combine the steps and requirements from all different schedulers, and
	// sort them by their time offset, counting the schedulers' startTimes as
	// well.
	type trackedStep struct ***REMOVED***
		ExecutionStep
		configID int
	***REMOVED***
	trackedSteps := []trackedStep***REMOVED******REMOVED***
	for configID, config := range sortedConfigs ***REMOVED*** // orderly iteration over a slice
		configStartTime := config.GetStartTime()
		configSteps := config.GetExecutionRequirements(executionSegment)
		for _, cs := range configSteps ***REMOVED***
			cs.TimeOffset += configStartTime // add the scheduler start time to the step time offset
			trackedSteps = append(trackedSteps, trackedStep***REMOVED***cs, configID***REMOVED***)
		***REMOVED***
	***REMOVED***
	// Sort by (time offset, config id). It's important that we use stable
	// sorting algorithm, since there are could be steps with the same time from
	// the same scheduler and their order is important.
	sort.SliceStable(trackedSteps, func(a, b int) bool ***REMOVED***
		switch ***REMOVED***
		case trackedSteps[a].TimeOffset < trackedSteps[b].TimeOffset:
			return true
		case trackedSteps[a].TimeOffset == trackedSteps[b].TimeOffset:
			return trackedSteps[a].configID < trackedSteps[b].configID
		default:
			return false
		***REMOVED***
	***REMOVED***)

	// Go through all of the sorted steps from all of the schedulers, and
	// build a new list of execution steps that consolidates all of their
	// requirements. If multiple schedulers have an execution step at exactly
	// the same time offset, they will be combined into a single new execution
	// step with the sum of the values from the previous ones.
	currentTimeOffset := time.Duration(0)
	currentPlannedVUs := make([]uint64, len(scs))
	currentMaxUnplannedVUs := make([]uint64, len(scs))
	sum := func(data []uint64) (result uint64) ***REMOVED*** // sigh...
		for _, val := range data ***REMOVED***
			result += val
		***REMOVED***
		return result
	***REMOVED***
	consolidatedSteps := []ExecutionStep***REMOVED******REMOVED***
	addCurrentStepIfDifferent := func() ***REMOVED***
		newPlannedVUs := sum(currentPlannedVUs)
		newMaxUnplannedVUs := sum(currentMaxUnplannedVUs)
		stepsLen := len(consolidatedSteps)
		if stepsLen == 0 ||
			consolidatedSteps[stepsLen-1].PlannedVUs != newPlannedVUs ||
			consolidatedSteps[stepsLen-1].MaxUnplannedVUs != newMaxUnplannedVUs ***REMOVED***

			consolidatedSteps = append(consolidatedSteps, ExecutionStep***REMOVED***
				TimeOffset:      currentTimeOffset,
				PlannedVUs:      newPlannedVUs,
				MaxUnplannedVUs: newMaxUnplannedVUs,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	for _, step := range trackedSteps ***REMOVED***
		// If the time offset is different, create a new step with the current values
		if step.TimeOffset != currentTimeOffset ***REMOVED***
			addCurrentStepIfDifferent()
			currentTimeOffset = step.TimeOffset
		***REMOVED***
		currentPlannedVUs[step.configID] = step.PlannedVUs
		currentMaxUnplannedVUs[step.configID] = step.MaxUnplannedVUs
	***REMOVED***
	addCurrentStepIfDifferent() // Add the last step
	return consolidatedSteps
***REMOVED***

// GetParsedSchedulerConfig returns a struct instance corresponding to the supplied
// config type. It will be fully initialized - with both the default values of
// the type, as well as with whatever the user had specified in the JSON
func GetParsedSchedulerConfig(name, configType string, rawJSON []byte) (result SchedulerConfig, err error) ***REMOVED***
	schedulerConfigTypesMutex.Lock()
	defer schedulerConfigTypesMutex.Unlock()

	constructor, exists := schedulerConfigConstructors[configType]
	if !exists ***REMOVED***
		return nil, fmt.Errorf("unknown execution scheduler type '%s'", configType)
	***REMOVED***
	return constructor(name, rawJSON)
***REMOVED***

type protoSchedulerConfig struct ***REMOVED***
	configType string
	rawJSON    json.RawMessage
***REMOVED***

// UnmarshalJSON just reads unmarshals the base config (to get the type), but it also
// stores the unprocessed JSON so we can parse the full config in the next step
func (pc *protoSchedulerConfig) UnmarshalJSON(b []byte) error ***REMOVED***
	var tmp struct ***REMOVED***
		ConfigType string `json:"type"`
	***REMOVED***
	err := json.Unmarshal(b, &tmp)
	*pc = protoSchedulerConfig***REMOVED***tmp.ConfigType, b***REMOVED***
	return err
***REMOVED***
