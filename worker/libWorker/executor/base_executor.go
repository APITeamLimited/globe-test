package executor

import (
	"context"
	"strconv"
	"sync"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/metrics"
	"github.com/sirupsen/logrus"
)

// BaseExecutor is a helper struct that contains common properties and methods
// between most executors. It is intended to be used as an anonymous struct
// inside of most of the executors, for the purpose of reducing boilerplate
// code.
type BaseExecutor struct {
	config         libWorker.ExecutorConfig
	executionState *libWorker.ExecutionState
	iterSegIndexMx *sync.Mutex
	iterSegIndex   *libWorker.SegmentedIndex
	logger         *logrus.Entry
}

// NewBaseExecutor returns an initialized BaseExecutor
func NewBaseExecutor(config libWorker.ExecutorConfig, es *libWorker.ExecutionState, logger *logrus.Entry) *BaseExecutor {
	segIdx := libWorker.NewSegmentedIndex(es.ExecutionTuple)
	return &BaseExecutor{
		config:         config,
		executionState: es,
		logger:         logger,
		iterSegIndexMx: new(sync.Mutex),
		iterSegIndex:   segIdx,
	}
}

// nextIterationCounters next scaled(local) and unscaled(global) iteration counters
func (bs *BaseExecutor) nextIterationCounters() (uint64, uint64) {
	bs.iterSegIndexMx.Lock()
	defer bs.iterSegIndexMx.Unlock()
	scaled, unscaled := bs.iterSegIndex.Next()
	return uint64(scaled - 1), uint64(unscaled - 1)
}

// Init doesn't do anything for most executors, since initialization of all
// planned VUs is handled by the executor.
func (bs *BaseExecutor) Init(_ context.Context) error {
	return nil
}

// GetConfig returns the configuration with which this executor was launched.
func (bs *BaseExecutor) GetConfig() libWorker.ExecutorConfig {
	return bs.config
}

// GetLogger returns the executor logger entry.
func (bs *BaseExecutor) GetLogger() *logrus.Entry {
	return bs.logger
}

// getMetricTags returns a tag set that can be used to emit metrics by the
// executor. The VU ID is optional.
func (bs *BaseExecutor) getMetricTags(vuID *uint64) *metrics.SampleTags {
	tags := make(map[string]string, len(bs.executionState.Test.Options.RunTags))
	for k, v := range bs.executionState.Test.Options.RunTags {
		tags[k] = v
	}
	if bs.executionState.Test.Options.SystemTags.Has(metrics.TagScenario) {
		tags["scenario"] = bs.config.GetName()
	}
	if vuID != nil && bs.executionState.Test.Options.SystemTags.Has(metrics.TagVU) {
		tags["vu"] = strconv.FormatUint(*vuID, 10)
	}
	return metrics.IntoSampleTags(&tags)
}
