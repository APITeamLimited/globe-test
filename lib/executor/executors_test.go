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
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v3"
)

type exp struct ***REMOVED***
	parseError      bool
	validationError bool
	custom          func(t *testing.T, cm lib.ExecutorConfigMap)
***REMOVED***

type configMapTestCase struct ***REMOVED***
	rawJSON  string
	expected exp
***REMOVED***

//nolint:lll,gochecknoglobals
var configMapTestCases = []configMapTestCase***REMOVED***
	***REMOVED***"", exp***REMOVED***parseError: true***REMOVED******REMOVED***,
	***REMOVED***"1234", exp***REMOVED***parseError: true***REMOVED******REMOVED***,
	***REMOVED***"asdf", exp***REMOVED***parseError: true***REMOVED******REMOVED***,
	***REMOVED***"'adsf'", exp***REMOVED***parseError: true***REMOVED******REMOVED***,
	***REMOVED***"[]", exp***REMOVED***parseError: true***REMOVED******REMOVED***,
	***REMOVED***"***REMOVED******REMOVED***", exp***REMOVED***custom: func(t *testing.T, cm lib.ExecutorConfigMap) ***REMOVED***
		assert.Equal(t, cm, lib.ExecutorConfigMap***REMOVED******REMOVED***)
	***REMOVED******REMOVED******REMOVED***,
	***REMOVED***"***REMOVED******REMOVED***asdf", exp***REMOVED***parseError: true***REMOVED******REMOVED***,
	***REMOVED***"null", exp***REMOVED***custom: func(t *testing.T, cm lib.ExecutorConfigMap) ***REMOVED***
		assert.Nil(t, cm)
	***REMOVED******REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"someKey": ***REMOVED******REMOVED******REMOVED***`, exp***REMOVED***parseError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"someKey": ***REMOVED***"type": "constant-blah-blah", "vus": 10, "duration": "60s"***REMOVED******REMOVED***`, exp***REMOVED***parseError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"someKey": ***REMOVED***"type": "constant-looping-vus", "uknownField": "should_error"***REMOVED******REMOVED***`, exp***REMOVED***parseError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"someKey": ***REMOVED***"type": "constant-looping-vus", "vus": 10, "duration": "60s", "env": 123***REMOVED******REMOVED***`, exp***REMOVED***parseError: true***REMOVED******REMOVED***,

	// Validation errors for constant-looping-vus and the base config
	***REMOVED***`***REMOVED***"someKey": ***REMOVED***"type": "constant-looping-vus", "vus": 10, "duration": "60s",
		"gracefulStop": "10s", "startTime": "70s", "env": ***REMOVED***"test": "mest"***REMOVED***, "exec": "someFunc"***REMOVED******REMOVED***`,
		exp***REMOVED***custom: func(t *testing.T, cm lib.ExecutorConfigMap) ***REMOVED***
			sched := NewConstantLoopingVUsConfig("someKey")
			sched.VUs = null.IntFrom(10)
			sched.Duration = types.NullDurationFrom(1 * time.Minute)
			sched.GracefulStop = types.NullDurationFrom(10 * time.Second)
			sched.StartTime = types.NullDurationFrom(70 * time.Second)
			sched.Exec = null.StringFrom("someFunc")
			sched.Env = map[string]string***REMOVED***"test": "mest"***REMOVED***
			require.Equal(t, cm, lib.ExecutorConfigMap***REMOVED***"someKey": sched***REMOVED***)
			require.Equal(t, sched.BaseConfig.Name, cm["someKey"].GetName())
			require.Equal(t, sched.BaseConfig.Type, cm["someKey"].GetType())
			require.Equal(t, sched.BaseConfig.GetGracefulStop(), cm["someKey"].GetGracefulStop())
			require.Equal(t,
				sched.BaseConfig.StartTime.Duration,
				types.Duration(cm["someKey"].GetStartTime()),
			)
			require.Equal(t, sched.BaseConfig.Env, cm["someKey"].GetEnv())

			assert.Empty(t, cm["someKey"].Validate())
			assert.Empty(t, cm.Validate())

			assert.Equal(t, "10 looping VUs for 1m0s (exec: someFunc, startTime: 1m10s, gracefulStop: 10s)", cm["someKey"].GetDescription(nil))

			schedReqs := cm["someKey"].GetExecutionRequirements(nil)
			endOffset, isFinal := lib.GetEndOffset(schedReqs)
			assert.Equal(t, 70*time.Second, endOffset)
			assert.Equal(t, true, isFinal)
			assert.Equal(t, uint64(10), lib.GetMaxPlannedVUs(schedReqs))
			assert.Equal(t, uint64(10), lib.GetMaxPossibleVUs(schedReqs))

			totalReqs := cm.GetFullExecutionRequirements(nil)
			endOffset, isFinal = lib.GetEndOffset(totalReqs)
			assert.Equal(t, 140*time.Second, endOffset)
			assert.Equal(t, true, isFinal)
			assert.Equal(t, uint64(10), lib.GetMaxPlannedVUs(schedReqs))
			assert.Equal(t, uint64(10), lib.GetMaxPossibleVUs(schedReqs))

		***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***`***REMOVED***"aname": ***REMOVED***"type": "constant-looping-vus", "duration": "60s"***REMOVED******REMOVED***`, exp***REMOVED******REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"": ***REMOVED***"type": "constant-looping-vus", "vus": 10, "duration": "60s"***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"aname": ***REMOVED***"type": "constant-looping-vus"***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"aname": ***REMOVED***"type": "constant-looping-vus", "vus": 0.5***REMOVED******REMOVED***`, exp***REMOVED***parseError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"aname": ***REMOVED***"type": "constant-looping-vus", "vus": 10***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"aname": ***REMOVED***"type": "constant-looping-vus", "vus": 0, "duration": "60s"***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"aname": ***REMOVED***"type": "constant-looping-vus", "vus": -1, "duration": "60s"***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"aname": ***REMOVED***"type": "constant-looping-vus", "vus": 10, "duration": "0s"***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"aname": ***REMOVED***"type": "constant-looping-vus", "vus": 10, "duration": "10s", "startTime": "-10s"***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"aname": ***REMOVED***"type": "constant-looping-vus", "vus": 10, "duration": "10s", "exec": ""***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"aname": ***REMOVED***"type": "constant-looping-vus", "vus": 10, "duration": "10s", "gracefulStop": "-2s"***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	// variable-looping-vus
	***REMOVED***`***REMOVED***"varloops": ***REMOVED***"type": "variable-looping-vus", "startVUs": 20, "gracefulStop": "15s", "gracefulRampDown": "10s",
		    "startTime": "23s", "stages": [***REMOVED***"duration": "60s", "target": 30***REMOVED***, ***REMOVED***"duration": "130s", "target": 10***REMOVED***]***REMOVED******REMOVED***`,
		exp***REMOVED***custom: func(t *testing.T, cm lib.ExecutorConfigMap) ***REMOVED***
			sched := NewVariableLoopingVUsConfig("varloops")
			sched.GracefulStop = types.NullDurationFrom(15 * time.Second)
			sched.GracefulRampDown = types.NullDurationFrom(10 * time.Second)
			sched.StartVUs = null.IntFrom(20)
			sched.StartTime = types.NullDurationFrom(23 * time.Second)
			sched.Stages = []Stage***REMOVED***
				***REMOVED***Target: null.IntFrom(30), Duration: types.NullDurationFrom(60 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(10), Duration: types.NullDurationFrom(130 * time.Second)***REMOVED***,
			***REMOVED***
			require.Equal(t, cm, lib.ExecutorConfigMap***REMOVED***"varloops": sched***REMOVED***)

			assert.Empty(t, cm["varloops"].Validate())
			assert.Empty(t, cm.Validate())

			assert.Equal(t, "Up to 30 looping VUs for 3m10s over 2 stages (gracefulRampDown: 10s, startTime: 23s, gracefulStop: 15s)", cm["varloops"].GetDescription(nil))

			schedReqs := cm["varloops"].GetExecutionRequirements(nil)
			endOffset, isFinal := lib.GetEndOffset(schedReqs)
			assert.Equal(t, 205*time.Second, endOffset)
			assert.Equal(t, true, isFinal)
			assert.Equal(t, uint64(30), lib.GetMaxPlannedVUs(schedReqs))
			assert.Equal(t, uint64(30), lib.GetMaxPossibleVUs(schedReqs))

			totalReqs := cm.GetFullExecutionRequirements(nil)
			endOffset, isFinal = lib.GetEndOffset(totalReqs)
			assert.Equal(t, 228*time.Second, endOffset)
			assert.Equal(t, true, isFinal)
			assert.Equal(t, uint64(30), lib.GetMaxPlannedVUs(schedReqs))
			assert.Equal(t, uint64(30), lib.GetMaxPossibleVUs(schedReqs))
		***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***`***REMOVED***"varloops": ***REMOVED***"type": "variable-looping-vus", "startVUs": 0, "stages": [***REMOVED***"duration": "60s", "target": 0***REMOVED***]***REMOVED******REMOVED***`, exp***REMOVED******REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"varloops": ***REMOVED***"type": "variable-looping-vus", "startVUs": -1, "stages": [***REMOVED***"duration": "60s", "target": 30***REMOVED***]***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"varloops": ***REMOVED***"type": "variable-looping-vus", "startVUs": 2, "stages": [***REMOVED***"duration": "-60s", "target": 30***REMOVED***]***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"varloops": ***REMOVED***"type": "variable-looping-vus", "startVUs": 2, "stages": [***REMOVED***"duration": "60s", "target": -30***REMOVED***]***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"varloops": ***REMOVED***"type": "variable-looping-vus", "stages": [***REMOVED***"duration": "60s"***REMOVED***]***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"varloops": ***REMOVED***"type": "variable-looping-vus", "stages": [***REMOVED***"target": 30***REMOVED***]***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"varloops": ***REMOVED***"type": "variable-looping-vus", "stages": []***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"varloops": ***REMOVED***"type": "variable-looping-vus"***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	// shared-iterations
	***REMOVED***`***REMOVED***"ishared": ***REMOVED***"type": "shared-iterations", "iterations": 22, "vus": 12, "maxDuration": "100s"***REMOVED******REMOVED***`,
		exp***REMOVED***custom: func(t *testing.T, cm lib.ExecutorConfigMap) ***REMOVED***
			sched := NewSharedIterationsConfig("ishared")
			sched.Iterations = null.IntFrom(22)
			sched.MaxDuration = types.NullDurationFrom(100 * time.Second)
			sched.VUs = null.IntFrom(12)

			assert.Empty(t, cm["ishared"].Validate())
			assert.Empty(t, cm.Validate())

			assert.Equal(t, "22 iterations shared among 12 VUs (maxDuration: 1m40s, gracefulStop: 30s)", cm["ishared"].GetDescription(nil))

			schedReqs := cm["ishared"].GetExecutionRequirements(nil)
			endOffset, isFinal := lib.GetEndOffset(schedReqs)
			assert.Equal(t, 130*time.Second, endOffset)
			assert.Equal(t, true, isFinal)
			assert.Equal(t, uint64(12), lib.GetMaxPlannedVUs(schedReqs))
			assert.Equal(t, uint64(12), lib.GetMaxPossibleVUs(schedReqs))

			totalReqs := cm.GetFullExecutionRequirements(nil)
			assert.Equal(t, schedReqs, totalReqs)
		***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***`***REMOVED***"ishared": ***REMOVED***"type": "shared-iterations"***REMOVED******REMOVED***`, exp***REMOVED******REMOVED******REMOVED***, // Has 1 VU & 1 iter default values
	***REMOVED***`***REMOVED***"ishared": ***REMOVED***"type": "shared-iterations", "iterations": 20***REMOVED******REMOVED***`, exp***REMOVED******REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"ishared": ***REMOVED***"type": "shared-iterations", "vus": 10***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***, // error because VUs are more than iters
	***REMOVED***`***REMOVED***"ishared": ***REMOVED***"type": "shared-iterations", "iterations": 20, "vus": 10, "maxDuration": "30m"***REMOVED******REMOVED***`, exp***REMOVED******REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"ishared": ***REMOVED***"type": "shared-iterations", "iterations": 20, "vus": 10, "maxDuration": "-3m"***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"ishared": ***REMOVED***"type": "shared-iterations", "iterations": 20, "vus": 10, "maxDuration": "0s"***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"ishared": ***REMOVED***"type": "shared-iterations", "iterations": 20, "vus": -10***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"ishared": ***REMOVED***"type": "shared-iterations", "iterations": -1, "vus": 1***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"ishared": ***REMOVED***"type": "shared-iterations", "iterations": 20, "vus": 30***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	// per-vu-iterations
	***REMOVED***`***REMOVED***"ipervu": ***REMOVED***"type": "per-vu-iterations", "iterations": 23, "vus": 13, "gracefulStop": 0***REMOVED******REMOVED***`,
		exp***REMOVED***custom: func(t *testing.T, cm lib.ExecutorConfigMap) ***REMOVED***
			sched := NewPerVUIterationsConfig("ipervu")
			sched.Iterations = null.IntFrom(23)
			sched.GracefulStop = types.NullDurationFrom(0)
			sched.VUs = null.IntFrom(13)

			assert.Empty(t, cm["ipervu"].Validate())
			assert.Empty(t, cm.Validate())

			assert.Equal(t, "23 iterations for each of 13 VUs (maxDuration: 10m0s)", cm["ipervu"].GetDescription(nil))

			schedReqs := cm["ipervu"].GetExecutionRequirements(nil)
			endOffset, isFinal := lib.GetEndOffset(schedReqs)
			assert.Equal(t, 600*time.Second, endOffset)
			assert.Equal(t, true, isFinal)
			assert.Equal(t, uint64(13), lib.GetMaxPlannedVUs(schedReqs))
			assert.Equal(t, uint64(13), lib.GetMaxPossibleVUs(schedReqs))

			totalReqs := cm.GetFullExecutionRequirements(nil)
			assert.Equal(t, schedReqs, totalReqs)
		***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***`***REMOVED***"ipervu": ***REMOVED***"type": "per-vu-iterations"***REMOVED******REMOVED***`, exp***REMOVED******REMOVED******REMOVED***, // Has 1 VU & 1 iter default values
	***REMOVED***`***REMOVED***"ipervu": ***REMOVED***"type": "per-vu-iterations", "iterations": 20***REMOVED******REMOVED***`, exp***REMOVED******REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"ipervu": ***REMOVED***"type": "per-vu-iterations", "vus": 10***REMOVED******REMOVED***`, exp***REMOVED******REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"ipervu": ***REMOVED***"type": "per-vu-iterations", "iterations": 20, "vus": 10***REMOVED******REMOVED***`, exp***REMOVED******REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"ipervu": ***REMOVED***"type": "per-vu-iterations", "iterations": 20, "vus": 10, "maxDuration": "-3m"***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"ipervu": ***REMOVED***"type": "per-vu-iterations", "iterations": 20, "vus": 10, "maxDuration": "0s"***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"ipervu": ***REMOVED***"type": "per-vu-iterations", "iterations": 20, "vus": -10***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"ipervu": ***REMOVED***"type": "per-vu-iterations", "iterations": -1, "vus": 1***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,

	// constant-arrival-rate
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 30, "timeUnit": "1m", "duration": "10m", "preAllocatedVUs": 20, "maxVUs": 30***REMOVED******REMOVED***`,
		exp***REMOVED***custom: func(t *testing.T, cm lib.ExecutorConfigMap) ***REMOVED***
			sched := NewConstantArrivalRateConfig("carrival")
			sched.Rate = null.IntFrom(30)
			sched.Duration = types.NullDurationFrom(10 * time.Minute)
			sched.TimeUnit = types.NullDurationFrom(1 * time.Minute)
			sched.PreAllocatedVUs = null.IntFrom(20)
			sched.MaxVUs = null.IntFrom(30)

			assert.Empty(t, cm["carrival"].Validate())
			assert.Empty(t, cm.Validate())

			assert.Equal(t, "0.50 iterations/s for 10m0s (maxVUs: 20-30, gracefulStop: 30s)", cm["carrival"].GetDescription(nil))

			schedReqs := cm["carrival"].GetExecutionRequirements(nil)
			endOffset, isFinal := lib.GetEndOffset(schedReqs)
			assert.Equal(t, 630*time.Second, endOffset)
			assert.Equal(t, true, isFinal)
			assert.Equal(t, uint64(20), lib.GetMaxPlannedVUs(schedReqs))
			assert.Equal(t, uint64(30), lib.GetMaxPossibleVUs(schedReqs))

			totalReqs := cm.GetFullExecutionRequirements(nil)
			assert.Equal(t, schedReqs, totalReqs)
		***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 10, "duration": "10m", "preAllocatedVUs": 20, "maxVUs": 30***REMOVED******REMOVED***`, exp***REMOVED******REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 10, "duration": "10m", "preAllocatedVUs": 20, "maxVUs": 30, "timeUnit": "-1s"***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 10, "duration": "10m", "preAllocatedVUs": 20***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 10, "duration": "10m", "maxVUs": 30***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 10, "preAllocatedVUs": 20, "maxVUs": 30***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "duration": "10m", "preAllocatedVUs": 20, "maxVUs": 30***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 10, "duration": "0m", "preAllocatedVUs": 20, "maxVUs": 30***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 0, "duration": "10m", "preAllocatedVUs": 20, "maxVUs": 30***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 10, "duration": "10m", "preAllocatedVUs": 20, "maxVUs": 15***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 10, "duration": "0s", "preAllocatedVUs": 20, "maxVUs": 25***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 10, "duration": "10m", "preAllocatedVUs": -2, "maxVUs": 25***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	// variable-arrival-rate
	***REMOVED***`***REMOVED***"varrival": ***REMOVED***"type": "variable-arrival-rate", "startRate": 10, "timeUnit": "30s", "preAllocatedVUs": 20,
		"maxVUs": 50, "stages": [***REMOVED***"duration": "3m", "target": 30***REMOVED***, ***REMOVED***"duration": "5m", "target": 10***REMOVED***]***REMOVED******REMOVED***`,
		exp***REMOVED***custom: func(t *testing.T, cm lib.ExecutorConfigMap) ***REMOVED***
			sched := NewVariableArrivalRateConfig("varrival")
			sched.StartRate = null.IntFrom(10)
			sched.Stages = []Stage***REMOVED***
				***REMOVED***Target: null.IntFrom(30), Duration: types.NullDurationFrom(180 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(10), Duration: types.NullDurationFrom(300 * time.Second)***REMOVED***,
			***REMOVED***
			sched.TimeUnit = types.NullDurationFrom(30 * time.Second)
			sched.PreAllocatedVUs = null.IntFrom(20)
			sched.MaxVUs = null.IntFrom(50)
			require.Equal(t, cm, lib.ExecutorConfigMap***REMOVED***"varrival": sched***REMOVED***)

			assert.Empty(t, cm["varrival"].Validate())
			assert.Empty(t, cm.Validate())

			assert.Equal(t, "Up to 1.00 iterations/s for 8m0s over 2 stages (maxVUs: 20-50, gracefulStop: 30s)", cm["varrival"].GetDescription(nil))

			schedReqs := cm["varrival"].GetExecutionRequirements(nil)
			endOffset, isFinal := lib.GetEndOffset(schedReqs)
			assert.Equal(t, 510*time.Second, endOffset)
			assert.Equal(t, true, isFinal)
			assert.Equal(t, uint64(20), lib.GetMaxPlannedVUs(schedReqs))
			assert.Equal(t, uint64(50), lib.GetMaxPossibleVUs(schedReqs))

			totalReqs := cm.GetFullExecutionRequirements(nil)
			assert.Equal(t, schedReqs, totalReqs)
		***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***`***REMOVED***"varrival": ***REMOVED***"type": "variable-arrival-rate", "preAllocatedVUs": 20, "maxVUs": 50, "stages": [***REMOVED***"duration": "5m", "target": 10***REMOVED***]***REMOVED******REMOVED***`, exp***REMOVED******REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"varrival": ***REMOVED***"type": "variable-arrival-rate", "preAllocatedVUs": -20, "maxVUs": 50, "stages": [***REMOVED***"duration": "5m", "target": 10***REMOVED***]***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"varrival": ***REMOVED***"type": "variable-arrival-rate", "startRate": -1, "preAllocatedVUs": 20, "maxVUs": 50, "stages": [***REMOVED***"duration": "5m", "target": 10***REMOVED***]***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"varrival": ***REMOVED***"type": "variable-arrival-rate", "preAllocatedVUs": 20, "stages": [***REMOVED***"duration": "5m", "target": 10***REMOVED***]***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"varrival": ***REMOVED***"type": "variable-arrival-rate", "maxVUs": 50, "stages": [***REMOVED***"duration": "5m", "target": 10***REMOVED***]***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"varrival": ***REMOVED***"type": "variable-arrival-rate", "preAllocatedVUs": 20, "maxVUs": 50***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"varrival": ***REMOVED***"type": "variable-arrival-rate", "preAllocatedVUs": 20, "maxVUs": 50, "stages": []***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"varrival": ***REMOVED***"type": "variable-arrival-rate", "preAllocatedVUs": 20, "maxVUs": 50, "stages": [***REMOVED***"duration": "5m", "target": 10***REMOVED***], "timeUnit": "-1s"***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"varrival": ***REMOVED***"type": "variable-arrival-rate", "preAllocatedVUs": 30, "maxVUs": 20, "stages": [***REMOVED***"duration": "5m", "target": 10***REMOVED***]***REMOVED******REMOVED***`, exp***REMOVED***validationError: true***REMOVED******REMOVED***,
	//TODO: more tests of mixed executors and execution plans
***REMOVED***

func TestConfigMapParsingAndValidation(t *testing.T) ***REMOVED***
	t.Parallel()
	for i, tc := range configMapTestCases ***REMOVED***
		tc := tc
		t.Run(fmt.Sprintf("TestCase#%d", i), func(t *testing.T) ***REMOVED***
			t.Logf(tc.rawJSON)
			var result lib.ExecutorConfigMap
			err := json.Unmarshal([]byte(tc.rawJSON), &result)
			if tc.expected.parseError ***REMOVED***
				require.Error(t, err)
				return
			***REMOVED***
			require.NoError(t, err)

			parseErrors := result.Validate()
			if tc.expected.validationError ***REMOVED***
				assert.NotEmpty(t, parseErrors)
			***REMOVED*** else ***REMOVED***
				assert.Empty(t, parseErrors)
			***REMOVED***
			if tc.expected.custom != nil ***REMOVED***
				tc.expected.custom(t, result)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVariableLoopingVUsConfigExecutionPlanExample(t *testing.T) ***REMOVED***
	t.Parallel()
	conf := NewVariableLoopingVUsConfig("test")
	conf.StartVUs = null.IntFrom(4)
	conf.Stages = []Stage***REMOVED***
		***REMOVED***Target: null.IntFrom(6), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(1), Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(5), Duration: types.NullDurationFrom(4 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(1), Duration: types.NullDurationFrom(4 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(4), Duration: types.NullDurationFrom(3 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(4), Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(1), Duration: types.NullDurationFrom(0 * time.Second)***REMOVED***,
		***REMOVED***Target: null.IntFrom(1), Duration: types.NullDurationFrom(3 * time.Second)***REMOVED***,
	***REMOVED***

	expRawStepsNoZeroEnd := []lib.ExecutionStep***REMOVED***
		***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 6***REMOVED***,
		***REMOVED***TimeOffset: 3 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 4 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 5 * time.Second, PlannedVUs: 3***REMOVED***,
		***REMOVED***TimeOffset: 6 * time.Second, PlannedVUs: 2***REMOVED***,
		***REMOVED***TimeOffset: 7 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 8 * time.Second, PlannedVUs: 2***REMOVED***,
		***REMOVED***TimeOffset: 9 * time.Second, PlannedVUs: 3***REMOVED***,
		***REMOVED***TimeOffset: 10 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 11 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 12 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 13 * time.Second, PlannedVUs: 3***REMOVED***,
		***REMOVED***TimeOffset: 14 * time.Second, PlannedVUs: 2***REMOVED***,
		***REMOVED***TimeOffset: 15 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 16 * time.Second, PlannedVUs: 2***REMOVED***,
		***REMOVED***TimeOffset: 17 * time.Second, PlannedVUs: 3***REMOVED***,
		***REMOVED***TimeOffset: 18 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 20 * time.Second, PlannedVUs: 1***REMOVED***,
	***REMOVED***
	rawStepsNoZeroEnd := conf.getRawExecutionSteps(nil, false)
	assert.Equal(t, expRawStepsNoZeroEnd, rawStepsNoZeroEnd)
	endOffset, isFinal := lib.GetEndOffset(rawStepsNoZeroEnd)
	assert.Equal(t, 20*time.Second, endOffset)
	assert.Equal(t, false, isFinal)

	rawStepsZeroEnd := conf.getRawExecutionSteps(nil, true)
	assert.Equal(t,
		append(expRawStepsNoZeroEnd, lib.ExecutionStep***REMOVED***TimeOffset: 23 * time.Second, PlannedVUs: 0***REMOVED***),
		rawStepsZeroEnd,
	)
	endOffset, isFinal = lib.GetEndOffset(rawStepsZeroEnd)
	assert.Equal(t, 23*time.Second, endOffset)
	assert.Equal(t, true, isFinal)

	// GracefulStop and GracefulRampDown equal to the default 30 sec
	assert.Equal(t, []lib.ExecutionStep***REMOVED***
		***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 6***REMOVED***,
		***REMOVED***TimeOffset: 33 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 42 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 50 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 53 * time.Second, PlannedVUs: 0***REMOVED***,
	***REMOVED***, conf.GetExecutionRequirements(nil))

	// Try a longer GracefulStop than the GracefulRampDown
	conf.GracefulStop = types.NullDurationFrom(80 * time.Second)
	assert.Equal(t, []lib.ExecutionStep***REMOVED***
		***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 6***REMOVED***,
		***REMOVED***TimeOffset: 33 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 42 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 50 * time.Second, PlannedVUs: 1***REMOVED***,
		***REMOVED***TimeOffset: 103 * time.Second, PlannedVUs: 0***REMOVED***,
	***REMOVED***, conf.GetExecutionRequirements(nil))

	// Try a much shorter GracefulStop than the GracefulRampDown
	conf.GracefulStop = types.NullDurationFrom(3 * time.Second)
	assert.Equal(t, []lib.ExecutionStep***REMOVED***
		***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 6***REMOVED***,
		***REMOVED***TimeOffset: 26 * time.Second, PlannedVUs: 0***REMOVED***,
	***REMOVED***, conf.GetExecutionRequirements(nil))

	// Try a zero GracefulStop
	conf.GracefulStop = types.NullDurationFrom(0 * time.Second)
	assert.Equal(t, []lib.ExecutionStep***REMOVED***
		***REMOVED***TimeOffset: 0 * time.Second, PlannedVUs: 4***REMOVED***,
		***REMOVED***TimeOffset: 1 * time.Second, PlannedVUs: 5***REMOVED***,
		***REMOVED***TimeOffset: 2 * time.Second, PlannedVUs: 6***REMOVED***,
		***REMOVED***TimeOffset: 23 * time.Second, PlannedVUs: 0***REMOVED***,
	***REMOVED***, conf.GetExecutionRequirements(nil))

	// Try a zero GracefulStop and GracefulRampDown, i.e. raw steps with 0 end cap
	conf.GracefulRampDown = types.NullDurationFrom(0 * time.Second)
	assert.Equal(t, rawStepsZeroEnd, conf.GetExecutionRequirements(nil))
***REMOVED***
