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

package executor

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/metrics"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/stats"
	"go.k6.io/k6/ui/pb"
)

const externallyControlledType = "externally-controlled"

func init() ***REMOVED***
	lib.RegisterExecutorConfigType(
		externallyControlledType,
		func(name string, rawJSON []byte) (lib.ExecutorConfig, error) ***REMOVED***
			config := ExternallyControlledConfig***REMOVED***BaseConfig: NewBaseConfig(name, externallyControlledType)***REMOVED***
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

// ExternallyControlledConfigParams contains all of the options that actually
// determine the scheduling of VUs in the externally controlled executor.
type ExternallyControlledConfigParams struct ***REMOVED***
	VUs      null.Int           `json:"vus"`
	Duration types.NullDuration `json:"duration"` // 0 is a valid value, meaning infinite duration
	MaxVUs   null.Int           `json:"maxVUs"`
***REMOVED***

// Validate just checks the control options in isolation.
func (mecc ExternallyControlledConfigParams) Validate() (errors []error) ***REMOVED***
	if mecc.VUs.Int64 < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of VUs shouldn't be negative"))
	***REMOVED***

	if mecc.MaxVUs.Int64 < mecc.VUs.Int64 ***REMOVED***
		errors = append(errors, fmt.Errorf(
			"the number of active VUs (%d) must be less than or equal to the number of maxVUs (%d)",
			mecc.VUs.Int64, mecc.MaxVUs.Int64,
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

// ExternallyControlledConfig stores the number of currently active VUs, the max
// number of VUs and the executor duration. The duration can be 0, which means
// "infinite duration", i.e. the user has to manually abort the script.
type ExternallyControlledConfig struct ***REMOVED***
	BaseConfig
	ExternallyControlledConfigParams
***REMOVED***

// Make sure we implement the lib.ExecutorConfig interface
var _ lib.ExecutorConfig = &ExternallyControlledConfig***REMOVED******REMOVED***

// GetDescription returns a human-readable description of the executor options
func (mec ExternallyControlledConfig) GetDescription(_ *lib.ExecutionTuple) string ***REMOVED***
	duration := "infinite"
	if mec.Duration.Duration != 0 ***REMOVED***
		duration = mec.Duration.String()
	***REMOVED***
	return fmt.Sprintf(
		"Externally controlled execution with %d VUs, %d max VUs, %s duration",
		mec.VUs.Int64, mec.MaxVUs.Int64, duration,
	)
***REMOVED***

// Validate makes sure all options are configured and valid
func (mec ExternallyControlledConfig) Validate() []error ***REMOVED***
	errors := append(mec.BaseConfig.Validate(), mec.ExternallyControlledConfigParams.Validate()...)
	if mec.GracefulStop.Valid ***REMOVED***
		errors = append(errors, fmt.Errorf(
			"gracefulStop is not supported by the externally controlled executor",
		))
	***REMOVED***
	return errors
***REMOVED***

// GetExecutionRequirements reserves the configured number of max VUs for the
// whole duration of the executor, so these VUs can be externally initialized in
// the beginning of the test.
//
// Importantly, if 0 (i.e. infinite) duration is configured, this executor
// doesn't emit the last step to relinquish these VUs.
//
// Also, the externally controlled executor doesn't set MaxUnplannedVUs in the
// returned steps, since their initialization and usage is directly controlled
// by the user, can be changed during the test runtime, and is effectively
// bounded only by the resources of the machine k6 is running on.
//
// This is not a problem, because the MaxUnplannedVUs are mostly meant to be
// used for calculating the maximum possible number of initialized VUs at any
// point during a test run. That's used for sizing purposes and for user qouta
// checking in the cloud execution, where the externally controlled executor
// isn't supported.
func (mec ExternallyControlledConfig) GetExecutionRequirements(et *lib.ExecutionTuple) []lib.ExecutionStep ***REMOVED***
	startVUs := lib.ExecutionStep***REMOVED***
		TimeOffset:      0,
		PlannedVUs:      uint64(et.Segment.Scale(mec.MaxVUs.Int64)), // user-configured, VUs to be pre-initialized
		MaxUnplannedVUs: 0,                                          // intentional, see function comment
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
// distribute the externally controlled executor.
func (ExternallyControlledConfig) IsDistributable() bool ***REMOVED***
	return false
***REMOVED***

// NewExecutor creates a new ExternallyControlled executor
func (mec ExternallyControlledConfig) NewExecutor(es *lib.ExecutionState, logger *logrus.Entry) (lib.Executor, error) ***REMOVED***
	return &ExternallyControlled***REMOVED***
		BaseExecutor:         NewBaseExecutor(mec, es, logger),
		config:               mec,
		currentControlConfig: mec.ExternallyControlledConfigParams,
		configLock:           &sync.RWMutex***REMOVED******REMOVED***,
		newControlConfigs:    make(chan updateConfigEvent),
		pauseEvents:          make(chan pauseEvent),
		hasStarted:           make(chan struct***REMOVED******REMOVED***),
	***REMOVED***, nil
***REMOVED***

// HasWork reports whether there is any work to be done for the given execution segment.
func (mec ExternallyControlledConfig) HasWork(_ *lib.ExecutionTuple) bool ***REMOVED***
	// We can always initialize new VUs via the REST API, so return true.
	return true
***REMOVED***

type pauseEvent struct ***REMOVED***
	isPaused bool
	err      chan error
***REMOVED***

type updateConfigEvent struct ***REMOVED***
	newConfig ExternallyControlledConfigParams
	err       chan error
***REMOVED***

// ExternallyControlled is an implementation of the old k6 executor that could be
// controlled externally, via the k6 REST API. It implements both the
// lib.PausableExecutor and the lib.LiveUpdatableExecutor interfaces.
type ExternallyControlled struct ***REMOVED***
	*BaseExecutor
	config               ExternallyControlledConfig
	currentControlConfig ExternallyControlledConfigParams
	configLock           *sync.RWMutex
	newControlConfigs    chan updateConfigEvent
	pauseEvents          chan pauseEvent
	hasStarted           chan struct***REMOVED******REMOVED***
***REMOVED***

// Make sure we implement all the interfaces
var (
	_ lib.Executor              = &ExternallyControlled***REMOVED******REMOVED***
	_ lib.PausableExecutor      = &ExternallyControlled***REMOVED******REMOVED***
	_ lib.LiveUpdatableExecutor = &ExternallyControlled***REMOVED******REMOVED***
)

// GetCurrentConfig just returns the executor's current configuration.
func (mex *ExternallyControlled) GetCurrentConfig() ExternallyControlledConfig ***REMOVED***
	mex.configLock.RLock()
	defer mex.configLock.RUnlock()
	return ExternallyControlledConfig***REMOVED***
		BaseConfig:                       mex.config.BaseConfig,
		ExternallyControlledConfigParams: mex.currentControlConfig,
	***REMOVED***
***REMOVED***

// GetConfig just returns the executor's current configuration, it's basically
// an alias of GetCurrentConfig that implements the more generic interface.
func (mex *ExternallyControlled) GetConfig() lib.ExecutorConfig ***REMOVED***
	return mex.GetCurrentConfig()
***REMOVED***

// GetProgress just returns the executor's progress bar instance.
func (mex ExternallyControlled) GetProgress() *pb.ProgressBar ***REMOVED***
	mex.configLock.RLock()
	defer mex.configLock.RUnlock()
	return mex.progress
***REMOVED***

// GetLogger just returns the executor's logger instance.
func (mex ExternallyControlled) GetLogger() *logrus.Entry ***REMOVED***
	mex.configLock.RLock()
	defer mex.configLock.RUnlock()
	return mex.logger
***REMOVED***

// Init doesn't do anything...
func (mex ExternallyControlled) Init(ctx context.Context) error ***REMOVED***
	return nil
***REMOVED***

// SetPaused pauses or resumes the executor.
func (mex *ExternallyControlled) SetPaused(paused bool) error ***REMOVED***
	select ***REMOVED***
	case <-mex.hasStarted:
		event := pauseEvent***REMOVED***isPaused: paused, err: make(chan error)***REMOVED***
		mex.pauseEvents <- event
		return <-event.err
	default:
		return fmt.Errorf("cannot pause the externally controlled executor before it has started")
	***REMOVED***
***REMOVED***

// UpdateConfig validates the supplied config and updates it in real time. It is
// possible to update the configuration even when k6 is paused, either in the
// beginning (i.e. when running k6 with --paused) or in the middle of the script
// execution.
func (mex *ExternallyControlled) UpdateConfig(ctx context.Context, newConf interface***REMOVED******REMOVED***) error ***REMOVED***
	newConfigParams, ok := newConf.(ExternallyControlledConfigParams)
	if !ok ***REMOVED***
		return errors.New("invalid config type")
	***REMOVED***
	if errs := newConfigParams.Validate(); len(errs) != 0 ***REMOVED***
		return fmt.Errorf("invalid configuration supplied: %s", lib.ConcatErrors(errs, ", "))
	***REMOVED***

	if newConfigParams.Duration.Valid && newConfigParams.Duration != mex.config.Duration ***REMOVED***
		return fmt.Errorf("the externally controlled executor duration cannot be changed")
	***REMOVED***
	if newConfigParams.MaxVUs.Valid && newConfigParams.MaxVUs.Int64 < mex.config.MaxVUs.Int64 ***REMOVED***
		// This limitation is because the externally controlled executor is
		// still an executor that participates in the overall k6 scheduling.
		// Thus, any VUs that were explicitly specified by the user in the
		// config may be reused from or by other executors.
		return fmt.Errorf(
			"the new number of max VUs cannot be lower than the starting number of max VUs (%d)",
			mex.config.MaxVUs.Int64,
		)
	***REMOVED***

	mex.configLock.Lock() // guard against a simultaneous start of the test (which will close hasStarted)
	select ***REMOVED***
	case <-mex.hasStarted:
		mex.configLock.Unlock()
		event := updateConfigEvent***REMOVED***newConfig: newConfigParams, err: make(chan error)***REMOVED***
		mex.newControlConfigs <- event
		return <-event.err
	case <-ctx.Done():
		mex.configLock.Unlock()
		return ctx.Err()
	default:
		mex.currentControlConfig = newConfigParams
		mex.configLock.Unlock()
		return nil
	***REMOVED***
***REMOVED***

// This is a helper function that is used in run for non-infinite durations.
func (mex *ExternallyControlled) stopWhenDurationIsReached(ctx context.Context, duration time.Duration, cancel func()) ***REMOVED***
	ctxDone := ctx.Done()
	checkInterval := time.NewTicker(100 * time.Millisecond)
	for ***REMOVED***
		select ***REMOVED***
		case <-ctxDone:
			checkInterval.Stop()
			return

		// TODO: something saner and more optimized that sleeps for pauses and
		// doesn't depend on the global execution state?
		case <-checkInterval.C:
			elapsed := mex.executionState.GetCurrentTestRunDuration() - time.Duration(mex.config.StartTime.Duration)
			if elapsed >= duration ***REMOVED***
				cancel()
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// manualVUHandle is a wrapper around the vuHandle helper, used in the
// ramping-vus executor. Here, instead of using its getVU and returnVU
// methods to retrieve and return a VU from the global buffer, we use them to
// accurately update the local and global active VU counters and to ensure that
// the pausing and reducing VUs operations wait for VUs to fully finish
// executing their current iterations before returning.
type manualVUHandle struct ***REMOVED***
	*vuHandle
	initVU lib.InitializedVU
	wg     *sync.WaitGroup

	// This is the cancel of the local context, used to kill its goroutine when
	// we reduce the number of MaxVUs, so that the Go GC can clean up the VU.
	cancelVU func()
***REMOVED***

func (rs *externallyControlledRunState) newManualVUHandle(
	initVU lib.InitializedVU, logger *logrus.Entry,
) *manualVUHandle ***REMOVED***
	wg := sync.WaitGroup***REMOVED******REMOVED***
	state := rs.executor.executionState
	getVU := func() (lib.InitializedVU, error) ***REMOVED***
		wg.Add(1)
		state.ModCurrentlyActiveVUsCount(+1)
		atomic.AddInt64(rs.activeVUsCount, +1)
		return initVU, nil
	***REMOVED***
	returnVU := func(_ lib.InitializedVU) ***REMOVED***
		state.ModCurrentlyActiveVUsCount(-1)
		atomic.AddInt64(rs.activeVUsCount, -1)
		wg.Done()
	***REMOVED***
	ctx, cancel := context.WithCancel(rs.ctx)
	return &manualVUHandle***REMOVED***
		vuHandle: newStoppedVUHandle(ctx, getVU, returnVU,
			rs.executor.nextIterationCounters,
			&rs.executor.config.BaseConfig, logger),
		initVU:   initVU,
		wg:       &wg,
		cancelVU: cancel,
	***REMOVED***
***REMOVED***

// externallyControlledRunState is created and initialized by the Run() method
// of the externally controlled executor. It is used to track and modify various
// details of the execution, including handling of live config changes.
type externallyControlledRunState struct ***REMOVED***
	ctx             context.Context
	executor        *ExternallyControlled
	startMaxVUs     int64             // the scaled number of initially configured MaxVUs
	duration        time.Duration     // the total duration of the executor, could be 0 for infinite
	activeVUsCount  *int64            // the current number of active VUs, used only for the progress display
	maxVUs          *int64            // the current number of initialized VUs
	vuHandles       []*manualVUHandle // handles for manipulating and tracking all of the VUs
	currentlyPaused bool              // whether the executor is currently paused

	runIteration func(context.Context, lib.ActiveVU) bool // a helper closure function that runs a single iteration
***REMOVED***

// retrieveStartMaxVUs gets and initializes the (scaled) number of MaxVUs
// from the global VU buffer. These are the VUs that the user originally
// specified in the JS config, and that the ExecutionScheduler pre-initialized
// for us.
func (rs *externallyControlledRunState) retrieveStartMaxVUs() error ***REMOVED***
	for i := int64(0); i < rs.startMaxVUs; i++ ***REMOVED*** // get the initial planned VUs from the common buffer
		initVU, vuGetErr := rs.executor.executionState.GetPlannedVU(rs.executor.logger, false)
		if vuGetErr != nil ***REMOVED***
			return vuGetErr
		***REMOVED***
		vuHandle := rs.newManualVUHandle(initVU, rs.executor.logger.WithField("vuNum", i))
		go vuHandle.runLoopsIfPossible(rs.runIteration)
		rs.vuHandles[i] = vuHandle
	***REMOVED***
	return nil
***REMOVED***

func (rs *externallyControlledRunState) progressFn() (float64, []string) ***REMOVED***
	// TODO: simulate spinner for the other case or cycle 0-100?
	currentActiveVUs := atomic.LoadInt64(rs.activeVUsCount)
	currentMaxVUs := atomic.LoadInt64(rs.maxVUs)
	vusFmt := pb.GetFixedLengthIntFormat(currentMaxVUs)
	progVUs := fmt.Sprintf(vusFmt+"/"+vusFmt+" VUs", currentActiveVUs, currentMaxVUs)

	right := []string***REMOVED***progVUs, rs.duration.String(), ""***REMOVED***

	// TODO: use a saner way to calculate the elapsed time, without relying on
	// the global execution state...
	elapsed := rs.executor.executionState.GetCurrentTestRunDuration() - time.Duration(
		rs.executor.config.StartTime.Duration)
	if elapsed > rs.duration ***REMOVED***
		return 1, right
	***REMOVED***

	progress := 0.0
	if rs.duration > 0 ***REMOVED***
		progress = math.Min(1, float64(elapsed)/float64(rs.duration))
	***REMOVED***

	spentDuration := pb.GetFixedLengthDuration(elapsed, rs.duration)
	progDur := fmt.Sprintf("%s/%s", spentDuration, rs.duration)
	right[1] = progDur

	return progress, right
***REMOVED***

func (rs *externallyControlledRunState) handleConfigChange(oldCfg, newCfg ExternallyControlledConfigParams) error ***REMOVED***
	executionState := rs.executor.executionState
	segment := executionState.Options.ExecutionSegment
	oldActiveVUs := segment.Scale(oldCfg.VUs.Int64)
	oldMaxVUs := segment.Scale(oldCfg.MaxVUs.Int64)
	newActiveVUs := segment.Scale(newCfg.VUs.Int64)
	newMaxVUs := segment.Scale(newCfg.MaxVUs.Int64)

	rs.executor.logger.WithFields(logrus.Fields***REMOVED***
		"oldActiveVUs": oldActiveVUs, "oldMaxVUs": oldMaxVUs,
		"newActiveVUs": newActiveVUs, "newMaxVUs": newMaxVUs,
	***REMOVED***).Debug("Updating execution configuration...")

	for i := oldMaxVUs; i < newMaxVUs; i++ ***REMOVED***
		select ***REMOVED*** // check if the user didn't try to abort k6 while we're scaling up the VUs
		case <-rs.ctx.Done():
			return rs.ctx.Err()
		default: // do nothing
		***REMOVED***
		initVU, vuInitErr := executionState.InitializeNewVU(rs.ctx, rs.executor.logger)
		if vuInitErr != nil ***REMOVED***
			return vuInitErr
		***REMOVED***
		vuHandle := rs.newManualVUHandle(initVU, rs.executor.logger.WithField("vuNum", i))
		go vuHandle.runLoopsIfPossible(rs.runIteration)
		rs.vuHandles = append(rs.vuHandles, vuHandle)
	***REMOVED***

	if oldActiveVUs < newActiveVUs ***REMOVED***
		for i := oldActiveVUs; i < newActiveVUs; i++ ***REMOVED***
			if !rs.currentlyPaused ***REMOVED***
				if err := rs.vuHandles[i].start(); err != nil ***REMOVED***
					// TODO: maybe just log it ?
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for i := newActiveVUs; i < oldActiveVUs; i++ ***REMOVED***
			rs.vuHandles[i].hardStop()
		***REMOVED***
		for i := newActiveVUs; i < oldActiveVUs; i++ ***REMOVED***
			rs.vuHandles[i].wg.Wait()
		***REMOVED***
	***REMOVED***

	if oldMaxVUs > newMaxVUs ***REMOVED***
		for i := newMaxVUs; i < oldMaxVUs; i++ ***REMOVED***
			rs.vuHandles[i].cancelVU()
			if i < rs.startMaxVUs ***REMOVED***
				// return the initial planned VUs to the common buffer
				executionState.ReturnVU(rs.vuHandles[i].initVU, false)
			***REMOVED*** else ***REMOVED***
				executionState.ModInitializedVUsCount(-1)
			***REMOVED***
			rs.vuHandles[i] = nil
		***REMOVED***
		rs.vuHandles = rs.vuHandles[:newMaxVUs]
	***REMOVED***

	atomic.StoreInt64(rs.maxVUs, newMaxVUs)
	return nil
***REMOVED***

// Run constantly loops through as many iterations as possible on a variable
// dynamically controlled number of VUs either for the specified duration, or
// until the test is manually stopped.
// nolint:funlen,gocognit,cyclop
func (mex *ExternallyControlled) Run(
	parentCtx context.Context, out chan<- stats.SampleContainer, _ *metrics.BuiltinMetrics,
) (err error) ***REMOVED***
	mex.configLock.RLock()
	// Safely get the current config - it's important that the close of the
	// hasStarted channel is inside of the lock, so that there are no data races
	// between it and the UpdateConfig() method.
	currentControlConfig := mex.currentControlConfig
	close(mex.hasStarted)
	mex.configLock.RUnlock()

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	duration := time.Duration(currentControlConfig.Duration.Duration)
	if duration > 0 ***REMOVED*** // Only keep track of duration if it's not infinite
		go mex.stopWhenDurationIsReached(ctx, duration, cancel)
	***REMOVED***

	mex.logger.WithFields(
		logrus.Fields***REMOVED***"type": externallyControlledType, "duration": duration***REMOVED***,
	).Debug("Starting executor run...")

	startMaxVUs := mex.executionState.Options.ExecutionSegment.Scale(mex.config.MaxVUs.Int64)

	ss := &lib.ScenarioState***REMOVED***
		Name:      mex.config.Name,
		Executor:  mex.config.Type,
		StartTime: time.Now(),
	***REMOVED***
	ctx = lib.WithScenarioState(ctx, ss)

	runState := &externallyControlledRunState***REMOVED***
		ctx:             ctx,
		executor:        mex,
		startMaxVUs:     startMaxVUs,
		duration:        duration,
		vuHandles:       make([]*manualVUHandle, startMaxVUs),
		currentlyPaused: false,
		activeVUsCount:  new(int64),
		maxVUs:          new(int64),
		runIteration:    getIterationRunner(mex.executionState, mex.logger),
	***REMOVED***
	ss.ProgressFn = runState.progressFn

	*runState.maxVUs = startMaxVUs
	if err = runState.retrieveStartMaxVUs(); err != nil ***REMOVED***
		return err
	***REMOVED***

	mex.progress.Modify(pb.WithProgress(runState.progressFn)) // Keep track of the progress
	go trackProgress(parentCtx, ctx, ctx, mex, runState.progressFn)

	err = runState.handleConfigChange( // Start by setting MaxVUs to the starting MaxVUs
		ExternallyControlledConfigParams***REMOVED***MaxVUs: mex.config.MaxVUs***REMOVED***, currentControlConfig,
	)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer func() ***REMOVED*** // Make sure we release the VUs at the end
		err = runState.handleConfigChange(currentControlConfig, ExternallyControlledConfigParams***REMOVED******REMOVED***)
	***REMOVED***()

	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return nil
		case updateConfigEvent := <-mex.newControlConfigs:
			err := runState.handleConfigChange(currentControlConfig, updateConfigEvent.newConfig)
			if err != nil ***REMOVED***
				updateConfigEvent.err <- err
				if ctx.Err() == err ***REMOVED***
					return nil // we've already returned an error to the API client, but k6 should stop normally
				***REMOVED***
				return err
			***REMOVED***
			currentControlConfig = updateConfigEvent.newConfig
			mex.configLock.Lock()
			mex.currentControlConfig = updateConfigEvent.newConfig
			mex.configLock.Unlock()
			updateConfigEvent.err <- nil

		case pauseEvent := <-mex.pauseEvents:
			if pauseEvent.isPaused == runState.currentlyPaused ***REMOVED***
				pauseEvent.err <- nil
				continue
			***REMOVED***
			activeVUs := currentControlConfig.VUs.Int64
			if pauseEvent.isPaused ***REMOVED***
				for i := int64(0); i < activeVUs; i++ ***REMOVED***
					runState.vuHandles[i].gracefulStop()
				***REMOVED***
				for i := int64(0); i < activeVUs; i++ ***REMOVED***
					runState.vuHandles[i].wg.Wait()
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				for i := int64(0); i < activeVUs; i++ ***REMOVED***
					if err := runState.vuHandles[i].start(); err != nil ***REMOVED***
						// TODO again ... just log it?
						pauseEvent.err <- err
						return err
					***REMOVED***
				***REMOVED***
			***REMOVED***
			runState.currentlyPaused = pauseEvent.isPaused
			pauseEvent.err <- nil
		***REMOVED***
	***REMOVED***
***REMOVED***
