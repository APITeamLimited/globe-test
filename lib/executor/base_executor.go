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
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/metrics"
	"go.k6.io/k6/ui/pb"
)

// BaseExecutor is a helper struct that contains common properties and methods
// between most executors. It is intended to be used as an anonymous struct
// inside of most of the executors, for the purpose of reducing boilerplate
// code.
type BaseExecutor struct ***REMOVED***
	config         lib.ExecutorConfig
	executionState *lib.ExecutionState
	iterSegIndexMx *sync.Mutex
	iterSegIndex   *lib.SegmentedIndex
	logger         *logrus.Entry
	progress       *pb.ProgressBar
***REMOVED***

// NewBaseExecutor returns an initialized BaseExecutor
func NewBaseExecutor(config lib.ExecutorConfig, es *lib.ExecutionState, logger *logrus.Entry) *BaseExecutor ***REMOVED***
	segIdx := lib.NewSegmentedIndex(es.ExecutionTuple)
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
func (bs *BaseExecutor) GetConfig() lib.ExecutorConfig ***REMOVED***
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
func (bs *BaseExecutor) getMetricTags(vuID *uint64) *metrics.SampleTags ***REMOVED***
	tags := bs.executionState.Options.RunTags.CloneTags()
	if bs.executionState.Options.SystemTags.Has(metrics.TagScenario) ***REMOVED***
		tags["scenario"] = bs.config.GetName()
	***REMOVED***
	if vuID != nil && bs.executionState.Options.SystemTags.Has(metrics.TagVU) ***REMOVED***
		tags["vu"] = strconv.FormatUint(*vuID, 10)
	***REMOVED***
	return metrics.IntoSampleTags(&tags)
***REMOVED***
