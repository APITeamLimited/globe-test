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
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/ui/pb"
	"github.com/sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"
)

const manualExecutionType = "manual-execution"

func init() ***REMOVED***
	lib.RegisterSchedulerConfigType(
		manualExecutionType,
		func(name string, rawJSON []byte) (lib.SchedulerConfig, error) ***REMOVED***
			config := ManualExecutionConfig***REMOVED***BaseConfig: NewBaseConfig(name, manualExecutionType)***REMOVED***
			err := lib.StrictJSONUnmarshal(rawJSON, &config)
			if err != nil ***REMOVED***
				return config, err
			***REMOVED***
			if !config.MaxVUs.Valid ***REMOVED***
				config.MaxVUs = config.VUs
			***REMOVED***
			return config, nil
		***REMOVED***,
	)
***REMOVED***

// ManualExecutionControlConfig contains all of the options that actually
// determine the scheduling of VUs in the manual execution scheduler.
type ManualExecutionControlConfig struct ***REMOVED***
	VUs      null.Int           `json:"vus"`
	Duration types.NullDuration `json:"duration"`
	MaxVUs   null.Int           `json:"maxVUs"`
***REMOVED***

// Validate just checks the control options in isolation.
func (mecc ManualExecutionControlConfig) Validate() (errors []error) ***REMOVED***
	if mecc.VUs.Int64 < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of VUs shouldn't be negative"))
	***REMOVED***

	if mecc.MaxVUs.Int64 < mecc.VUs.Int64 ***REMOVED***
		errors = append(errors, fmt.Errorf(
			"the specified maxVUs (%d) should more than or equal to the the number of active VUs (%d)",
			mecc.MaxVUs.Int64, mecc.VUs.Int64,
		))
	***REMOVED***

	if !mecc.Duration.Valid ***REMOVED***
		errors = append(errors, fmt.Errorf("the duration should be specified, for infinite duration use 0"))
	***REMOVED*** else if time.Duration(mecc.Duration.Duration) < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf(
			"the duration shouldn't be negative, for infinite duration use 0",
		))
	***REMOVED***

	return errors
***REMOVED***

// ManualExecutionConfig stores the number of currently active VUs, the max
// number of VUs and the scheduler duration. The duration can be 0, which means
// "infinite duration", i.e. the user has to manually abort the script.
type ManualExecutionConfig struct ***REMOVED***
	BaseConfig
	ManualExecutionControlConfig
***REMOVED***

// Make sure we implement the lib.SchedulerConfig interface
var _ lib.SchedulerConfig = &ManualExecutionConfig***REMOVED******REMOVED***

// GetDescription returns a human-readable description of the scheduler options
func (mec ManualExecutionConfig) GetDescription(_ *lib.ExecutionSegment) string ***REMOVED***
	duration := "infinite"
	if mec.Duration.Duration != 0 ***REMOVED***
		duration = mec.Duration.String()
	***REMOVED***
	return fmt.Sprintf(
		"Manually controlled execution with %d VUs, %d max VUs, %s duration",
		mec.VUs.Int64, mec.MaxVUs.Int64, duration,
	)
***REMOVED***

// Validate makes sure all options are configured and valid
func (mec ManualExecutionConfig) Validate() []error ***REMOVED***
	errors := append(mec.BaseConfig.Validate(), mec.ManualExecutionControlConfig.Validate()...)
	if mec.GracefulStop.Valid ***REMOVED***
		errors = append(errors, fmt.Errorf(
			"gracefulStop is not supported by the manual execution scheduler",
		))
	***REMOVED***
	return errors
***REMOVED***

// GetExecutionRequirements just reserves the specified number of max VUs for
// the whole duration of the scheduler, so these VUs can be initialized in the
// beginning of the test.
//
// Importantly, if 0 (i.e. infinite) duration is configured, this scheduler
// doesn't emit the last step to relinquish these VUs.
//
// Also, the manual execution scheduler doesn't set MaxUnplannedVUs in the
// returned steps, since their initialization and usage is directly controlled
// by the user and is effectively bounded only by the resources of the machine
// k6 is running on.
//
// This is not a problem, because the MaxUnplannedVUs are mostly meant to be
// used for calculating the maximum possble number of initialized VUs at any
// point during a test run. That's used for sizing purposes and for user qouta
// checking in the cloud execution, where the manual scheduler isn't supported.
func (mec ManualExecutionConfig) GetExecutionRequirements(es *lib.ExecutionSegment) []lib.ExecutionStep ***REMOVED***
	startVUs := lib.ExecutionStep***REMOVED***
		TimeOffset:      0,
		PlannedVUs:      uint64(es.Scale(mec.MaxVUs.Int64)), // use
		MaxUnplannedVUs: 0,                                  // intentional, see function comment
	***REMOVED***

	maxDuration := time.Duration(mec.Duration.Duration)
	if maxDuration == 0 ***REMOVED***
		// Infinite duration, don't emit 0 VUs at the end since there's no planned end
		return []lib.ExecutionStep***REMOVED***startVUs***REMOVED***
	***REMOVED***
	return []lib.ExecutionStep***REMOVED***startVUs, ***REMOVED***
		TimeOffset:      maxDuration,
		PlannedVUs:      0,
		MaxUnplannedVUs: 0, // intentional, see function comment
	***REMOVED******REMOVED***
***REMOVED***

// IsDistributable simply returns false because there's no way to reliably
// distribute the manual execution scheduler.
func (ManualExecutionConfig) IsDistributable() bool ***REMOVED***
	return false
***REMOVED***

// NewScheduler creates a new ManualExecution "scheduler"
func (mec ManualExecutionConfig) NewScheduler(
	es *lib.ExecutorState, logger *logrus.Entry) (lib.Scheduler, error) ***REMOVED***

	return &ManualExecution***REMOVED***
		startConfig:          mec,
		currentControlConfig: mec.ManualExecutionControlConfig,
		configLock:           &sync.RWMutex***REMOVED******REMOVED***,
		newControlConfigs:    make(chan updateConfigEvent),
		pauseEvents:          make(chan pauseEvent),
		hasStarted:           make(chan struct***REMOVED******REMOVED***),

		executorState: es,
		logger:        logger,
		progress:      pb.New(pb.WithLeft(mec.GetName)),
	***REMOVED***, nil
***REMOVED***

type pauseEvent struct ***REMOVED***
	isPaused bool
	err      chan error
***REMOVED***

type updateConfigEvent struct ***REMOVED***
	newConfig ManualExecutionControlConfig
	err       chan error
***REMOVED***

// ManualExecution is an implementation of the old k6 scheduler that could be
// controlled externally, via the k6 REST API. It implements both the
// lib.PausableScheduler and the lib.LiveUpdatableScheduler interfaces.
type ManualExecution struct ***REMOVED***
	startConfig          ManualExecutionConfig
	currentControlConfig ManualExecutionControlConfig
	configLock           *sync.RWMutex
	newControlConfigs    chan updateConfigEvent
	pauseEvents          chan pauseEvent
	hasStarted           chan struct***REMOVED******REMOVED***

	executorState *lib.ExecutorState
	logger        *logrus.Entry
	progress      *pb.ProgressBar
***REMOVED***

// Make sure we implement all the interfaces
var _ lib.Scheduler = &ManualExecution***REMOVED******REMOVED***
var _ lib.PausableScheduler = &ManualExecution***REMOVED******REMOVED***
var _ lib.LiveUpdatableScheduler = &ManualExecution***REMOVED******REMOVED***

// GetCurrentConfig just returns the scheduler's current configuration.
func (mex *ManualExecution) GetCurrentConfig() ManualExecutionConfig ***REMOVED***
	mex.configLock.RLock()
	defer mex.configLock.RUnlock()
	return ManualExecutionConfig***REMOVED***
		BaseConfig:                   mex.startConfig.BaseConfig,
		ManualExecutionControlConfig: mex.currentControlConfig,
	***REMOVED***
***REMOVED***

// GetConfig just returns the scheduler's current configuration, it's basically
// an alias of GetCurrentConfig that implements the more generic interface.
func (mex *ManualExecution) GetConfig() lib.SchedulerConfig ***REMOVED***
	return mex.GetCurrentConfig()
***REMOVED***

// GetProgress just returns the scheduler's progress bar instance.
func (mex ManualExecution) GetProgress() *pb.ProgressBar ***REMOVED***
	return mex.progress
***REMOVED***

// GetLogger just returns the scheduler's logger instance.
func (mex ManualExecution) GetLogger() *logrus.Entry ***REMOVED***
	return mex.logger
***REMOVED***

// Init doesn't do anything...
func (mex ManualExecution) Init(ctx context.Context) error ***REMOVED***
	return nil
***REMOVED***

// SetPaused pauses or resumes the scheduler.
func (mex *ManualExecution) SetPaused(paused bool) error ***REMOVED***
	select ***REMOVED***
	case <-mex.hasStarted:
		event := pauseEvent***REMOVED***isPaused: paused, err: make(chan error)***REMOVED***
		mex.pauseEvents <- event
		return <-event.err
	default:
		return fmt.Errorf("cannot pause the manual scheduler before it has started")
	***REMOVED***
***REMOVED***

// UpdateConfig validates the supplied config and updates it in real time. It is
// possible to update the configuration even when k6 is paused, either in the
// beginning (i.e. when running k6 with --paused) or in the middle of the script
// execution.
func (mex *ManualExecution) UpdateConfig(ctx context.Context, newConf interface***REMOVED******REMOVED***) error ***REMOVED***
	newManualConfig, ok := newConf.(ManualExecutionControlConfig)
	if !ok ***REMOVED***
		return errors.New("invalid config type")
	***REMOVED***
	if errs := newManualConfig.Validate(); len(errs) != 0 ***REMOVED***
		return fmt.Errorf("invalid confiuguration supplied: %s", lib.ConcatErrors(errs, ", "))
	***REMOVED***

	if newManualConfig.Duration != mex.startConfig.Duration ***REMOVED***
		return fmt.Errorf("the manual scheduler duration cannot be changed")
	***REMOVED***
	if newManualConfig.MaxVUs.Int64 < mex.startConfig.MaxVUs.Int64 ***REMOVED***
		// This limitation is because the manual execution scheduler is still a
		// scheduler that participates in the overall k6 scheduling. Thus, any
		// VUs that were explicitly specified by the user in the config may be
		// reused from or by other schedulers.
		return fmt.Errorf(
			"the new number of max VUs cannot be lower than the starting number of max VUs (%d)",
			mex.startConfig.MaxVUs.Int64,
		)
	***REMOVED***

	mex.configLock.Lock()
	select ***REMOVED***
	case <-mex.hasStarted:
		mex.configLock.Unlock()
		event := updateConfigEvent***REMOVED***newConfig: newManualConfig, err: make(chan error)***REMOVED***
		mex.newControlConfigs <- event
		return <-event.err
	case <-ctx.Done():
		mex.configLock.Unlock()
		return ctx.Err()
	default:
		mex.currentControlConfig = newManualConfig
		mex.configLock.Unlock()
		return nil
	***REMOVED***
***REMOVED***

// This is a helper function that is used in run for non-infinite durations.
func (mex *ManualExecution) stopWhenDurationIsReached(ctx context.Context, duration time.Duration, cancel func()) ***REMOVED***
	ctxDone := ctx.Done()
	checkInterval := time.NewTicker(100 * time.Millisecond)
	for ***REMOVED***
		select ***REMOVED***
		case <-ctxDone:
			checkInterval.Stop()
			return

		//TODO: something more optimized that sleeps for pauses?
		case <-checkInterval.C:
			if mex.executorState.GetCurrentTestRunDuration() >= duration ***REMOVED***
				cancel()
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// manualVUHandle is a wrapper around the vuHandle helper, used in the
// variable-looping-vus scheduler. Here, instead of using its getVU and returnVU
// methods to retrieve and return a VU from the global buffer, we use them to
// accurately update the local and global active VU counters and to ensure that
// the pausing and reducing VUs operations wait for VUs to fully finish
// executing their current iterations before returning.
type manualVUHandle struct ***REMOVED***
	*vuHandle
	vu lib.VU
	wg *sync.WaitGroup

	// This is the cancel of the local context, used to kill its goroutine when
	// we reduce the number of MaxVUs, so that the Go GC can clean up the VU.
	cancelVU func()
***REMOVED***

func newManualVUHandle(
	parentCtx context.Context, state *lib.ExecutorState, localActiveVUsCount *int64, vu lib.VU, logger *logrus.Entry,
) *manualVUHandle ***REMOVED***

	wg := sync.WaitGroup***REMOVED******REMOVED***
	getVU := func() (lib.VU, error) ***REMOVED***
		wg.Add(1)
		state.ModCurrentlyActiveVUsCount(+1)
		atomic.AddInt64(localActiveVUsCount, +1)
		return vu, nil
	***REMOVED***
	returnVU := func(_ lib.VU) ***REMOVED***
		state.ModCurrentlyActiveVUsCount(-1)
		atomic.AddInt64(localActiveVUsCount, -1)
		wg.Done()
	***REMOVED***
	ctx, cancel := context.WithCancel(parentCtx)
	return &manualVUHandle***REMOVED***
		vuHandle: newStoppedVUHandle(ctx, getVU, returnVU, logger),
		vu:       vu,
		wg:       &wg,
		cancelVU: cancel,
	***REMOVED***
***REMOVED***

// Run constantly loops through as many iterations as possible on a variable
// dynamically controlled number of VUs either for the specified duration, or
// until the test is manually stopped.
//
//TODO: split this up? somehow... :/
func (mex *ManualExecution) Run(parentCtx context.Context, out chan<- stats.SampleContainer) (err error) ***REMOVED***
	mex.configLock.RLock()
	// Safely get the current config - it's important that the close of the
	// hasStarted channel is inside of the lock, so that there are no data races
	// between it and the UpdateConfig() method.
	currentControlConfig := mex.currentControlConfig
	close(mex.hasStarted)
	mex.configLock.RUnlock()

	segment := mex.executorState.Options.ExecutionSegment
	duration := time.Duration(currentControlConfig.Duration.Duration)

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	if duration > 0 ***REMOVED*** // Only keep track of duration if it's not infinite
		go mex.stopWhenDurationIsReached(ctx, duration, cancel)
	***REMOVED***

	mex.logger.WithFields(
		logrus.Fields***REMOVED***"type": manualExecutionType, "duration": duration***REMOVED***,
	).Debug("Starting scheduler run...")

	// Retrieve and initialize the (scaled) number of MaxVUs from the global VU
	// buffer that the user originally specified in the JS config.
	startMaxVUs := segment.Scale(mex.startConfig.MaxVUs.Int64)
	vuHandles := make([]*manualVUHandle, startMaxVUs)
	activeVUsCount := new(int64)
	runIteration := getIterationRunner(mex.executorState, mex.logger, out)
	for i := int64(0); i < startMaxVUs; i++ ***REMOVED*** // get the initial planned VUs from the common buffer
		vu, vuGetErr := mex.executorState.GetPlannedVU(mex.logger, false)
		if vuGetErr != nil ***REMOVED***
			return vuGetErr
		***REMOVED***
		vuHandle := newManualVUHandle(
			parentCtx, mex.executorState, activeVUsCount, vu, mex.logger.WithField("vuNum", i),
		)
		go vuHandle.runLoopsIfPossible(runIteration)
		vuHandles[i] = vuHandle
	***REMOVED***

	// Keep track of the progress
	maxVUs := new(int64)
	*maxVUs = startMaxVUs
	progresFn := func() (float64, string) ***REMOVED***
		spent := mex.executorState.GetCurrentTestRunDuration()
		progress := 0.0
		if duration > 0 ***REMOVED***
			progress = math.Min(1, float64(spent)/float64(duration))
		***REMOVED***
		//TODO: simulate spinner for the other case or cycle 0-100?
		currentActiveVUs := atomic.LoadInt64(activeVUsCount)
		currentMaxVUs := atomic.LoadInt64(maxVUs)
		vusFmt := pb.GetFixedLengthIntFormat(currentMaxVUs)
		return progress, fmt.Sprintf(
			"currently "+vusFmt+" out of "+vusFmt+" active looping VUs, %s/%s", currentActiveVUs, currentMaxVUs,
			pb.GetFixedLengthDuration(spent, duration), duration,
		)
	***REMOVED***
	mex.progress.Modify(pb.WithProgress(progresFn))
	go trackProgress(parentCtx, ctx, ctx, mex, progresFn)

	currentlyPaused := false
	waitVUs := func(from, to int64) ***REMOVED***
		for i := from; i < to; i++ ***REMOVED***
			vuHandles[i].wg.Wait()
		***REMOVED***
	***REMOVED***
	handleConfigChange := func(oldControlConfig, newControlConfig ManualExecutionControlConfig) error ***REMOVED***
		oldActiveVUs := segment.Scale(oldControlConfig.VUs.Int64)
		oldMaxVUs := segment.Scale(oldControlConfig.MaxVUs.Int64)
		newActiveVUs := segment.Scale(newControlConfig.VUs.Int64)
		newMaxVUs := segment.Scale(newControlConfig.MaxVUs.Int64)

		mex.logger.WithFields(logrus.Fields***REMOVED***
			"oldActiveVUs": oldActiveVUs, "oldMaxVUs": oldMaxVUs,
			"newActiveVUs": newActiveVUs, "newMaxVUs": newMaxVUs,
		***REMOVED***).Debug("Updating execution configuration...")

		for i := oldMaxVUs; i < newMaxVUs; i++ ***REMOVED***
			vu, vuInitErr := mex.executorState.InitializeNewVU(ctx, mex.logger)
			if vuInitErr != nil ***REMOVED***
				return vuInitErr
			***REMOVED***
			vuHandle := newManualVUHandle(
				ctx, mex.executorState, activeVUsCount, vu, mex.logger.WithField("vuNum", i),
			)
			go vuHandle.runLoopsIfPossible(runIteration)
			vuHandles = append(vuHandles, vuHandle)
		***REMOVED***

		if oldActiveVUs < newActiveVUs ***REMOVED***
			for i := oldActiveVUs; i < newActiveVUs; i++ ***REMOVED***

				if !currentlyPaused ***REMOVED***
					vuHandles[i].start()
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for i := newActiveVUs; i < oldActiveVUs; i++ ***REMOVED***
				vuHandles[i].hardStop()
			***REMOVED***
			waitVUs(newActiveVUs, oldActiveVUs)
		***REMOVED***

		if oldMaxVUs > newMaxVUs ***REMOVED***
			for i := newMaxVUs; i < oldMaxVUs; i++ ***REMOVED***
				vuHandles[i].cancelVU()
				if i < startMaxVUs ***REMOVED***
					// return the initial planned VUs to the common buffer
					mex.executorState.ReturnVU(vuHandles[i].vu, false)
				***REMOVED*** else ***REMOVED***
					mex.executorState.ModInitializedVUsCount(-1)
				***REMOVED***
				vuHandles[i] = nil
			***REMOVED***
			vuHandles = vuHandles[:newMaxVUs]
		***REMOVED***

		atomic.StoreInt64(maxVUs, newMaxVUs)
		return nil
	***REMOVED***

	err = handleConfigChange(ManualExecutionControlConfig***REMOVED***MaxVUs: mex.startConfig.MaxVUs***REMOVED***, currentControlConfig)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer func() ***REMOVED***
		err = handleConfigChange(currentControlConfig, ManualExecutionControlConfig***REMOVED******REMOVED***)
	***REMOVED***()

	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return nil
		case updateConfigEvent := <-mex.newControlConfigs:
			err := handleConfigChange(currentControlConfig, updateConfigEvent.newConfig)
			if err != nil ***REMOVED***
				updateConfigEvent.err <- err
				return err
			***REMOVED***
			currentControlConfig = updateConfigEvent.newConfig
			mex.configLock.Lock()
			mex.currentControlConfig = updateConfigEvent.newConfig
			mex.configLock.Unlock()
			updateConfigEvent.err <- nil

		case pauseEvent := <-mex.pauseEvents:
			if pauseEvent.isPaused == currentlyPaused ***REMOVED***
				pauseEvent.err <- nil
				continue
			***REMOVED***
			activeVUs := currentControlConfig.VUs.Int64
			if pauseEvent.isPaused ***REMOVED***
				for i := int64(0); i < activeVUs; i++ ***REMOVED***
					vuHandles[i].gracefulStop()
				***REMOVED***
				waitVUs(0, activeVUs)
			***REMOVED*** else ***REMOVED***
				for i := int64(0); i < activeVUs; i++ ***REMOVED***
					vuHandles[i].start()
				***REMOVED***
			***REMOVED***
			currentlyPaused = pauseEvent.isPaused
			pauseEvent.err <- nil
		***REMOVED***
	***REMOVED***
***REMOVED***
