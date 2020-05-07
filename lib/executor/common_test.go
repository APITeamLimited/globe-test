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
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/testutils"
	"github.com/loadimpact/k6/lib/testutils/minirunner"
	"github.com/loadimpact/k6/stats"
)

func simpleRunner(vuFn func(context.Context) error) lib.Runner ***REMOVED***
	return &minirunner.MiniRunner***REMOVED***
		Fn: func(ctx context.Context, _ chan<- stats.SampleContainer) error ***REMOVED***
			return vuFn(ctx)
		***REMOVED***,
	***REMOVED***
***REMOVED***

func setupExecutor(t *testing.T, config lib.ExecutorConfig, es *lib.ExecutionState, runner lib.Runner) (
	context.Context, context.CancelFunc, lib.Executor, *testutils.SimpleLogrusHook,
) ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	engineOut := make(chan stats.SampleContainer, 100) // TODO: return this for more complicated tests?

	logHook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.WarnLevel***REMOVED******REMOVED***
	testLog := logrus.New()
	testLog.AddHook(logHook)
	testLog.SetOutput(ioutil.Discard)
	logEntry := logrus.NewEntry(testLog)

	initVUFunc := func(_ context.Context, logger *logrus.Entry) (lib.InitializedVU, error) ***REMOVED***
		return runner.NewVU(int64(es.GetUniqueVUIdentifier()), engineOut)
	***REMOVED***
	es.SetInitVUFunc(initVUFunc)

	et, err := lib.NewExecutionTuple(es.Options.ExecutionSegment, es.Options.ExecutionSegmentSequence)
	require.NoError(t, err)

	maxPlannedVUs := lib.GetMaxPlannedVUs(config.GetExecutionRequirements(et))
	initializeVUs(ctx, t, logEntry, es, maxPlannedVUs, initVUFunc)

	executor, err := config.NewExecutor(es, logEntry)
	require.NoError(t, err)

	err = executor.Init(ctx)
	require.NoError(t, err)
	return ctx, cancel, executor, logHook
***REMOVED***

func initializeVUs(
	ctx context.Context, t testing.TB, logEntry *logrus.Entry, es *lib.ExecutionState, number uint64, initVU lib.InitVUFunc,
) ***REMOVED***
	// This is not how the local ExecutionScheduler initializes VUs, but should do the same job
	for i := uint64(0); i < number; i++ ***REMOVED***
		// Not calling es.InitializeNewVU() here to avoid a double increment of initializedVUs,
		// which is done in es.AddInitializedVU().
		vu, err := initVU(ctx, logEntry)
		require.NoError(t, err)
		es.ReturnVU(vu, false)
	***REMOVED***
***REMOVED***
