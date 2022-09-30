package executor

import (
	"context"
	"strconv"
	"sync"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/pb"
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
	"github.com/sirupsen/logrus"
)

// BaseExecutor is a helper struct that contains common properties and methods
// between most executors. It is intended to be used as an anonymous struct
// inside of most of the executors, for the purpose of reducing boilerplate
// code.
type BaseExecutor struct ***REMOVED***
	config         libWorker.ExecutorConfig
	executionState *libWorker.ExecutionState
	iterSegIndexMx *sync.Mutex
	iterSegIndex   *libWorker.SegmentedIndex
	logger         *logrus.Entry
	progress       *pb.ProgressBar
***REMOVED***

// NewBaseExecutor returns an initialized BaseExecutor
func NewBaseExecutor(config libWorker.ExecutorConfig, es *libWorker.ExecutionState, logger *logrus.Entry) *BaseExecutor ***REMOVED***
	segIdx := libWorker.NewSegmentedIndex(es.ExecutionTuple)
	return &BaseExecutor***REMOVED***
		config:         config,
		executionState: es,
		logger:         logger,
		iterSegIndexMx: new(sync.Mutex),
		iterSegIndex:   segIdx,
		progress: pb.New(
			pb.WithLeft(config.GetName),
			pb.WithLogger(logger),
		),
	***REMOVED***
***REMOVED***

// nextIterationCounters next scaled(local) and unscaled(global) iteration counters
func (bs *BaseExecutor) nextIterationCounters() (uint64, uint64) ***REMOVED***
	bs.iterSegIndexMx.Lock()
	defer bs.iterSegIndexMx.Unlock()
	scaled, unscaled := bs.iterSegIndex.Next()
	return uint64(scaled - 1), uint64(unscaled - 1)
***REMOVED***

// Init doesn't do anything for most executors, since initialization of all
// planned VUs is handled by the executor.
func (bs *BaseExecutor) Init(_ context.Context) error ***REMOVED***
	return nil
***REMOVED***

// GetConfig returns the configuration with which this executor was launched.
func (bs *BaseExecutor) GetConfig() libWorker.ExecutorConfig ***REMOVED***
	return bs.config
***REMOVED***

// GetLogger returns the executor logger entry.
func (bs *BaseExecutor) GetLogger() *logrus.Entry ***REMOVED***
	return bs.logger
***REMOVED***

// GetProgress just returns the progressbar pointer.
func (bs *BaseExecutor) GetProgress() *pb.ProgressBar ***REMOVED***
	return bs.progress
***REMOVED***

// getMetricTags returns a tag set that can be used to emit metrics by the
// executor. The VU ID is optional.
func (bs *BaseExecutor) getMetricTags(vuID *uint64) *workerMetrics.SampleTags ***REMOVED***
	tags := make(map[string]string, len(bs.executionState.Test.Options.RunTags))
	for k, v := range bs.executionState.Test.Options.RunTags ***REMOVED***
		tags[k] = v
	***REMOVED***
	if bs.executionState.Test.Options.SystemTags.Has(workerMetrics.TagScenario) ***REMOVED***
		tags["scenario"] = bs.config.GetName()
	***REMOVED***
	if vuID != nil && bs.executionState.Test.Options.SystemTags.Has(workerMetrics.TagVU) ***REMOVED***
		tags["vu"] = strconv.FormatUint(*vuID, 10)
	***REMOVED***
	return workerMetrics.IntoSampleTags(&tags)
***REMOVED***
