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
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/loadimpact/k6/lib/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v3"
)

type configMapTestCase struct ***REMOVED***
	rawJSON               string
	expectParseError      bool
	expectValidationError bool
	customValidator       func(t *testing.T, cm ConfigMap)
***REMOVED***

var configMapTestCases = []configMapTestCase***REMOVED***
	***REMOVED***"", true, false, nil***REMOVED***,
	***REMOVED***"1234", true, false, nil***REMOVED***,
	***REMOVED***"asdf", true, false, nil***REMOVED***,
	***REMOVED***"'adsf'", true, false, nil***REMOVED***,
	***REMOVED***"[]", true, false, nil***REMOVED***,
	***REMOVED***"***REMOVED******REMOVED***", false, false, func(t *testing.T, cm ConfigMap) ***REMOVED***
		assert.Equal(t, cm, ConfigMap***REMOVED******REMOVED***)
	***REMOVED******REMOVED***,
	***REMOVED***"***REMOVED******REMOVED***asdf", true, false, nil***REMOVED***,
	***REMOVED***"null", false, false, func(t *testing.T, cm ConfigMap) ***REMOVED***
		assert.Nil(t, cm)
	***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"someKey": ***REMOVED******REMOVED******REMOVED***`, true, false, nil***REMOVED***,
	***REMOVED***`***REMOVED***"someKey": ***REMOVED***"type": "constant-blah-blah", "vus": 10, "duration": "60s"***REMOVED******REMOVED***`, true, false, nil***REMOVED***,
	***REMOVED***`***REMOVED***"someKey": ***REMOVED***"type": "constant-looping-vus", "uknownField": "should_error"***REMOVED******REMOVED***`, true, false, nil***REMOVED***,
	***REMOVED***`***REMOVED***"someKey": ***REMOVED***"type": "constant-looping-vus", "vus": 10, "duration": "60s", "env": 123***REMOVED******REMOVED***`, true, false, nil***REMOVED***,

	// Validation errors for constant-looping-vus and the base config
	***REMOVED***`***REMOVED***"someKey": ***REMOVED***"type": "constant-looping-vus", "vus": 10, "duration": "60s", "interruptible": false,
		"iterationTimeout": "10s", "startTime": "70s", "env": ***REMOVED***"test": "mest"***REMOVED***, "exec": "someFunc"***REMOVED******REMOVED***`,
		false, false, func(t *testing.T, cm ConfigMap) ***REMOVED***
			sched := NewConstantLoopingVUsConfig("someKey")
			sched.VUs = null.IntFrom(10)
			sched.Duration = types.NullDurationFrom(1 * time.Minute)
			sched.Interruptible = null.BoolFrom(false)
			sched.IterationTimeout = types.NullDurationFrom(10 * time.Second)
			sched.StartTime = types.NullDurationFrom(70 * time.Second)
			sched.Exec = null.StringFrom("someFunc")
			sched.Env = map[string]string***REMOVED***"test": "mest"***REMOVED***
			require.Equal(t, cm, ConfigMap***REMOVED***"someKey": sched***REMOVED***)
			require.Equal(t, sched.BaseConfig, cm["someKey"].GetBaseConfig())
			assert.Equal(t, 70*time.Second, cm["someKey"].GetMaxDuration())
			assert.Equal(t, int64(10), cm["someKey"].GetMaxVUs())
			assert.Empty(t, cm["someKey"].Validate())
		***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"": ***REMOVED***"type": "constant-looping-vus", "vus": 10, "duration": "60s"***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"aname": ***REMOVED***"type": "constant-looping-vus"***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"aname": ***REMOVED***"type": "constant-looping-vus", "vus": 10***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"aname": ***REMOVED***"type": "constant-looping-vus", "duration": "60s"***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"aname": ***REMOVED***"type": "constant-looping-vus", "vus": -1, "duration": "60s"***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"aname": ***REMOVED***"type": "constant-looping-vus", "vus": 10, "duration": "0s"***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"aname": ***REMOVED***"type": "constant-looping-vus", "vus": 10, "duration": "10s", "startTime": "-10s"***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"aname": ***REMOVED***"type": "constant-looping-vus", "vus": 10, "duration": "10s", "exec": ""***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"aname": ***REMOVED***"type": "constant-looping-vus", "vus": 10, "duration": "10s", "iterationTimeout": "-2s"***REMOVED******REMOVED***`, false, true, nil***REMOVED***,

	// variable-looping-vus
	***REMOVED***`***REMOVED***"varloops": ***REMOVED***"type": "variable-looping-vus", "startVUs": 20, "iterationTimeout": "15s",
		"stages": [***REMOVED***"duration": "60s", "target": 30***REMOVED***, ***REMOVED***"duration": "120s", "target": 10***REMOVED***]***REMOVED******REMOVED***`,
		false, false, func(t *testing.T, cm ConfigMap) ***REMOVED***
			sched := NewVariableLoopingVUsConfig("varloops")
			sched.IterationTimeout = types.NullDurationFrom(15 * time.Second)
			sched.StartVUs = null.IntFrom(20)
			sched.Stages = []Stage***REMOVED***
				***REMOVED***Target: null.IntFrom(30), Duration: types.NullDurationFrom(60 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(10), Duration: types.NullDurationFrom(120 * time.Second)***REMOVED***,
			***REMOVED***
			require.Equal(t, cm, ConfigMap***REMOVED***"varloops": sched***REMOVED***)
			assert.Equal(t, int64(30), cm["varloops"].GetMaxVUs())
			assert.Equal(t, 195*time.Second, cm["varloops"].GetMaxDuration())
			assert.Empty(t, cm["varloops"].Validate())
		***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"varloops": ***REMOVED***"type": "variable-looping-vus", "startVUs": 0, "stages": [***REMOVED***"duration": "60s", "target": 0***REMOVED***]***REMOVED******REMOVED***`, false, false, nil***REMOVED***,
	***REMOVED***`***REMOVED***"varloops": ***REMOVED***"type": "variable-looping-vus", "startVUs": -1, "stages": [***REMOVED***"duration": "60s", "target": 30***REMOVED***]***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"varloops": ***REMOVED***"type": "variable-looping-vus", "startVUs": 2, "stages": [***REMOVED***"duration": "-60s", "target": 30***REMOVED***]***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"varloops": ***REMOVED***"type": "variable-looping-vus", "startVUs": 2, "stages": [***REMOVED***"duration": "60s", "target": -30***REMOVED***]***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"varloops": ***REMOVED***"type": "variable-looping-vus", "stages": [***REMOVED***"duration": "60s"***REMOVED***]***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"varloops": ***REMOVED***"type": "variable-looping-vus", "stages": [***REMOVED***"target": 30***REMOVED***]***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"varloops": ***REMOVED***"type": "variable-looping-vus", "stages": []***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"varloops": ***REMOVED***"type": "variable-looping-vus"***REMOVED******REMOVED***`, false, true, nil***REMOVED***,

	// shared-iterations
	***REMOVED***`***REMOVED***"ishared": ***REMOVED***"type": "shared-iterations", "iterations": 20, "vus": 10***REMOVED******REMOVED***`,
		false, false, func(t *testing.T, cm ConfigMap) ***REMOVED***
			sched := NewSharedIterationsConfig("ishared")
			sched.Iterations = null.IntFrom(20)
			sched.VUs = null.IntFrom(10)
			require.Equal(t, cm, ConfigMap***REMOVED***"ishared": sched***REMOVED***)
			assert.Equal(t, int64(10), cm["ishared"].GetMaxVUs())
			assert.Equal(t, 3630*time.Second, cm["ishared"].GetMaxDuration())
			assert.Empty(t, cm["ishared"].Validate())
		***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"ishared": ***REMOVED***"type": "shared-iterations"***REMOVED******REMOVED***`, false, false, nil***REMOVED***, // Has 1 VU & 1 iter default values
	***REMOVED***`***REMOVED***"ishared": ***REMOVED***"type": "shared-iterations", "iterations": 20***REMOVED******REMOVED***`, false, false, nil***REMOVED***,
	***REMOVED***`***REMOVED***"ishared": ***REMOVED***"type": "shared-iterations", "vus": 10***REMOVED******REMOVED***`, false, true, nil***REMOVED***, // error because VUs are more than iters
	***REMOVED***`***REMOVED***"ishared": ***REMOVED***"type": "shared-iterations", "iterations": 20, "vus": 10, "maxDuration": "30m"***REMOVED******REMOVED***`, false, false, nil***REMOVED***,
	***REMOVED***`***REMOVED***"ishared": ***REMOVED***"type": "shared-iterations", "iterations": 20, "vus": 10, "maxDuration": "-3m"***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"ishared": ***REMOVED***"type": "shared-iterations", "iterations": 20, "vus": 10, "maxDuration": "0s"***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"ishared": ***REMOVED***"type": "shared-iterations", "iterations": 20, "vus": -10***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"ishared": ***REMOVED***"type": "shared-iterations", "iterations": -1, "vus": 1***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"ishared": ***REMOVED***"type": "shared-iterations", "iterations": 20, "vus": 30***REMOVED******REMOVED***`, false, true, nil***REMOVED***,

	// per-vu-iterations
	***REMOVED***`***REMOVED***"ipervu": ***REMOVED***"type": "per-vu-iterations", "iterations": 20, "vus": 10***REMOVED******REMOVED***`,
		false, false, func(t *testing.T, cm ConfigMap) ***REMOVED***
			sched := NewPerVUIterationsConfig("ipervu")
			sched.Iterations = null.IntFrom(20)
			sched.VUs = null.IntFrom(10)
			require.Equal(t, cm, ConfigMap***REMOVED***"ipervu": sched***REMOVED***)
			assert.Equal(t, int64(10), cm["ipervu"].GetMaxVUs())
			assert.Equal(t, 3630*time.Second, cm["ipervu"].GetMaxDuration())
			assert.Empty(t, cm["ipervu"].Validate())
		***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"ipervu": ***REMOVED***"type": "per-vu-iterations"***REMOVED******REMOVED***`, false, false, nil***REMOVED***, // Has 1 VU & 1 iter default values
	***REMOVED***`***REMOVED***"ipervu": ***REMOVED***"type": "per-vu-iterations", "iterations": 20***REMOVED******REMOVED***`, false, false, nil***REMOVED***,
	***REMOVED***`***REMOVED***"ipervu": ***REMOVED***"type": "per-vu-iterations", "vus": 10***REMOVED******REMOVED***`, false, false, nil***REMOVED***,
	***REMOVED***`***REMOVED***"ipervu": ***REMOVED***"type": "per-vu-iterations", "iterations": 20, "vus": 10***REMOVED******REMOVED***`, false, false, nil***REMOVED***,
	***REMOVED***`***REMOVED***"ipervu": ***REMOVED***"type": "per-vu-iterations", "iterations": 20, "vus": 10, "maxDuration": "-3m"***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"ipervu": ***REMOVED***"type": "per-vu-iterations", "iterations": 20, "vus": 10, "maxDuration": "0s"***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"ipervu": ***REMOVED***"type": "per-vu-iterations", "iterations": 20, "vus": -10***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"ipervu": ***REMOVED***"type": "per-vu-iterations", "iterations": -1, "vus": 1***REMOVED******REMOVED***`, false, true, nil***REMOVED***,

	// constant-arrival-rate
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 10, "timeUnit": "1m", "duration": "10m", "preAllocatedVUs": 20, "maxVUs": 30***REMOVED******REMOVED***`,
		false, false, func(t *testing.T, cm ConfigMap) ***REMOVED***
			sched := NewConstantArrivalRateConfig("carrival")
			sched.Rate = null.IntFrom(10)
			sched.Duration = types.NullDurationFrom(10 * time.Minute)
			sched.TimeUnit = types.NullDurationFrom(1 * time.Minute)
			sched.PreAllocatedVUs = null.IntFrom(20)
			sched.MaxVUs = null.IntFrom(30)
			require.Equal(t, cm, ConfigMap***REMOVED***"carrival": sched***REMOVED***)
			assert.Equal(t, int64(30), cm["carrival"].GetMaxVUs())
			assert.Equal(t, 630*time.Second, cm["carrival"].GetMaxDuration())
			assert.Empty(t, cm["carrival"].Validate())
		***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 10, "duration": "10m", "preAllocatedVUs": 20, "maxVUs": 30***REMOVED******REMOVED***`, false, false, nil***REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 10, "duration": "10m", "preAllocatedVUs": 20, "maxVUs": 30, "timeUnit": "-1s"***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 10, "duration": "10m", "preAllocatedVUs": 20***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 10, "duration": "10m", "maxVUs": 30***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 10, "preAllocatedVUs": 20, "maxVUs": 30***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "duration": "10m", "preAllocatedVUs": 20, "maxVUs": 30***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 10, "duration": "0m", "preAllocatedVUs": 20, "maxVUs": 30***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 0, "duration": "10m", "preAllocatedVUs": 20, "maxVUs": 30***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 10, "duration": "10m", "preAllocatedVUs": 20, "maxVUs": 15***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 10, "duration": "0s", "preAllocatedVUs": 20, "maxVUs": 25***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"carrival": ***REMOVED***"type": "constant-arrival-rate", "rate": 10, "duration": "10m", "preAllocatedVUs": -2, "maxVUs": 25***REMOVED******REMOVED***`, false, true, nil***REMOVED***,

	// variable-arrival-rate
	***REMOVED***`***REMOVED***"varrival": ***REMOVED***"type": "variable-arrival-rate", "startRate": 10, "timeUnit": "30s", "preAllocatedVUs": 20, "maxVUs": 50,
		"stages": [***REMOVED***"duration": "3m", "target": 30***REMOVED***, ***REMOVED***"duration": "5m", "target": 10***REMOVED***]***REMOVED******REMOVED***`,
		false, false, func(t *testing.T, cm ConfigMap) ***REMOVED***
			sched := NewVariableArrivalRateConfig("varrival")
			sched.StartRate = null.IntFrom(10)
			sched.Stages = []Stage***REMOVED***
				***REMOVED***Target: null.IntFrom(30), Duration: types.NullDurationFrom(180 * time.Second)***REMOVED***,
				***REMOVED***Target: null.IntFrom(10), Duration: types.NullDurationFrom(300 * time.Second)***REMOVED***,
			***REMOVED***
			sched.TimeUnit = types.NullDurationFrom(30 * time.Second)
			sched.PreAllocatedVUs = null.IntFrom(20)
			sched.MaxVUs = null.IntFrom(50)
			require.Equal(t, cm, ConfigMap***REMOVED***"varrival": sched***REMOVED***)
			assert.Equal(t, int64(50), cm["varrival"].GetMaxVUs())
			assert.Equal(t, 510*time.Second, cm["varrival"].GetMaxDuration())
			assert.Empty(t, cm["varrival"].Validate())
		***REMOVED******REMOVED***,
	***REMOVED***`***REMOVED***"varrival": ***REMOVED***"type": "variable-arrival-rate", "preAllocatedVUs": 20, "maxVUs": 50, "stages": [***REMOVED***"duration": "5m", "target": 10***REMOVED***]***REMOVED******REMOVED***`, false, false, nil***REMOVED***,
	***REMOVED***`***REMOVED***"varrival": ***REMOVED***"type": "variable-arrival-rate", "preAllocatedVUs": -20, "maxVUs": 50, "stages": [***REMOVED***"duration": "5m", "target": 10***REMOVED***]***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"varrival": ***REMOVED***"type": "variable-arrival-rate", "startRate": -1, "preAllocatedVUs": 20, "maxVUs": 50, "stages": [***REMOVED***"duration": "5m", "target": 10***REMOVED***]***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"varrival": ***REMOVED***"type": "variable-arrival-rate", "preAllocatedVUs": 20, "stages": [***REMOVED***"duration": "5m", "target": 10***REMOVED***]***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"varrival": ***REMOVED***"type": "variable-arrival-rate", "maxVUs": 50, "stages": [***REMOVED***"duration": "5m", "target": 10***REMOVED***]***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"varrival": ***REMOVED***"type": "variable-arrival-rate", "preAllocatedVUs": 20, "maxVUs": 50***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"varrival": ***REMOVED***"type": "variable-arrival-rate", "preAllocatedVUs": 20, "maxVUs": 50, "stages": []***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"varrival": ***REMOVED***"type": "variable-arrival-rate", "preAllocatedVUs": 20, "maxVUs": 50, "stages": [***REMOVED***"duration": "5m", "target": 10***REMOVED***], "timeUnit": "-1s"***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
	***REMOVED***`***REMOVED***"varrival": ***REMOVED***"type": "variable-arrival-rate", "preAllocatedVUs": 30, "maxVUs": 20, "stages": [***REMOVED***"duration": "5m", "target": 10***REMOVED***]***REMOVED******REMOVED***`, false, true, nil***REMOVED***,
***REMOVED***

func TestConfigMapParsingAndValidation(t *testing.T) ***REMOVED***
	t.Parallel()
	for i, tc := range configMapTestCases ***REMOVED***
		tc := tc
		t.Run(fmt.Sprintf("TestCase#%d", i), func(t *testing.T) ***REMOVED***
			t.Logf(tc.rawJSON)
			var result ConfigMap
			err := json.Unmarshal([]byte(tc.rawJSON), &result)
			if tc.expectParseError ***REMOVED***
				require.Error(t, err)
				return
			***REMOVED***
			require.NoError(t, err)

			validationErrors := result.Validate()
			if tc.expectValidationError ***REMOVED***
				assert.NotEmpty(t, validationErrors)
			***REMOVED*** else ***REMOVED***
				assert.Empty(t, validationErrors)
			***REMOVED***
			if tc.customValidator != nil ***REMOVED***
				tc.customValidator(t, result)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

//TODO: check percentage split calculations
