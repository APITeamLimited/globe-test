package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/core"
	"go.k6.io/k6/core/local"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/testutils/minirunner"
)

func TestGetStatus(t *testing.T) ***REMOVED***
	t.Parallel()

	testState := getTestRunState(t, lib.Options***REMOVED******REMOVED***, &minirunner.MiniRunner***REMOVED******REMOVED***)
	execScheduler, err := local.NewExecutionScheduler(testState)
	require.NoError(t, err)
	engine, err := core.NewEngine(testState, execScheduler, nil)
	require.NoError(t, err)

	rw := httptest.NewRecorder()
	NewHandler().ServeHTTP(rw, newRequestWithEngine(engine, "GET", "/v1/status", nil))
	res := rw.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	t.Run("document", func(t *testing.T) ***REMOVED***
		t.Parallel()

		var doc StatusJSONAPI
		assert.NoError(t, json.Unmarshal(rw.Body.Bytes(), &doc))
		assert.Equal(t, "status", doc.Data.Type)
	***REMOVED***)

	t.Run("status", func(t *testing.T) ***REMOVED***
		t.Parallel()

		var statusEnvelop StatusJSONAPI

		err := json.Unmarshal(rw.Body.Bytes(), &statusEnvelop)
		assert.NoError(t, err)

		status := statusEnvelop.Status()

		assert.True(t, status.Paused.Valid)
		assert.True(t, status.VUs.Valid)
		assert.True(t, status.VUsMax.Valid)
		assert.False(t, status.Stopped)
		assert.False(t, status.Tainted)
	***REMOVED***)
***REMOVED***

func TestPatchStatus(t *testing.T) ***REMOVED***
	t.Parallel()

	testData := map[string]struct ***REMOVED***
		ExpectedStatusCode int
		ExpectedStatus     Status
		Payload            []byte
	***REMOVED******REMOVED***
		"nothing": ***REMOVED***
			ExpectedStatusCode: 200,
			ExpectedStatus:     Status***REMOVED******REMOVED***,
			Payload:            []byte(`***REMOVED***"data":***REMOVED***"type":"status","id":"default","attributes":***REMOVED***"status":0,"paused":null,"vus":null,"vus-max":null,"stopped":false,"running":false,"tainted":false***REMOVED******REMOVED******REMOVED***`),
		***REMOVED***,
		"paused": ***REMOVED***
			ExpectedStatusCode: 200,
			ExpectedStatus:     Status***REMOVED***Paused: null.BoolFrom(true)***REMOVED***,
			Payload:            []byte(`***REMOVED***"data":***REMOVED***"type":"status","id":"default","attributes":***REMOVED***"status":0,"paused":true,"vus":null,"vus-max":null,"stopped":false,"running":false,"tainted":false***REMOVED******REMOVED******REMOVED***`),
		***REMOVED***,
		"max vus": ***REMOVED***
			ExpectedStatusCode: 200,
			ExpectedStatus:     Status***REMOVED***VUsMax: null.IntFrom(20)***REMOVED***,
			Payload:            []byte(`***REMOVED***"data":***REMOVED***"type":"status","id":"default","attributes":***REMOVED***"status":0,"paused":null,"vus":null,"vus-max":20,"stopped":false,"running":false,"tainted":false***REMOVED******REMOVED******REMOVED***`),
		***REMOVED***,
		"max vus below initial": ***REMOVED***
			ExpectedStatusCode: 400,
			ExpectedStatus:     Status***REMOVED***VUsMax: null.IntFrom(5)***REMOVED***,
			Payload:            []byte(`***REMOVED***"data":***REMOVED***"type":"status","id":"default","attributes":***REMOVED***"status":0,"paused":null,"vus":null,"vus-max":5,"stopped":false,"running":false,"tainted":false***REMOVED******REMOVED******REMOVED***`),
		***REMOVED***,
		"too many vus": ***REMOVED***
			ExpectedStatusCode: 400,
			ExpectedStatus:     Status***REMOVED***VUs: null.IntFrom(10), VUsMax: null.IntFrom(0)***REMOVED***,
			Payload:            []byte(`***REMOVED***"data":***REMOVED***"type":"status","id":"default","attributes":***REMOVED***"status":0,"paused":null,"vus":10,"vus-max":0,"stopped":false,"running":false,"tainted":false***REMOVED******REMOVED******REMOVED***`),
		***REMOVED***,
		"vus": ***REMOVED***
			ExpectedStatusCode: 200,
			ExpectedStatus:     Status***REMOVED***VUs: null.IntFrom(10), VUsMax: null.IntFrom(10)***REMOVED***,
			Payload:            []byte(`***REMOVED***"data":***REMOVED***"type":"status","id":"default","attributes":***REMOVED***"status":0,"paused":null,"vus":10,"vus-max":10,"stopped":false,"running":false,"tainted":false***REMOVED******REMOVED******REMOVED***`),
		***REMOVED***,
	***REMOVED***

	for name, testCase := range testData ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			scenarios := lib.ScenarioConfigs***REMOVED******REMOVED***
			err := json.Unmarshal([]byte(`
			***REMOVED***"external": ***REMOVED***"executor": "externally-controlled",
			"vus": 0, "maxVUs": 10, "duration": "0"***REMOVED******REMOVED***`), &scenarios)
			require.NoError(t, err)

			testState := getTestRunState(t, lib.Options***REMOVED***Scenarios: scenarios***REMOVED***, &minirunner.MiniRunner***REMOVED******REMOVED***)
			execScheduler, err := local.NewExecutionScheduler(testState)
			require.NoError(t, err)
			engine, err := core.NewEngine(testState, execScheduler, nil)
			require.NoError(t, err)

			require.NoError(t, engine.OutputManager.StartOutputs())
			defer engine.OutputManager.StopOutputs()

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			run, wait, err := engine.Init(ctx, ctx)
			require.NoError(t, err)

			defer func() ***REMOVED***
				cancel()
				wait()
			***REMOVED***()

			go func() ***REMOVED***
				assert.NoError(t, run())
			***REMOVED***()
			// wait for the executor to initialize to avoid a potential data race below
			time.Sleep(200 * time.Millisecond)

			rw := httptest.NewRecorder()
			NewHandler().ServeHTTP(rw, newRequestWithEngine(engine, "PATCH", "/v1/status", bytes.NewReader(testCase.Payload)))
			res := rw.Result()

			require.Equal(t, testCase.ExpectedStatusCode, res.StatusCode)

			if testCase.ExpectedStatusCode != 200 ***REMOVED***
				return
			***REMOVED***

			status := NewStatus(engine)
			if testCase.ExpectedStatus.Paused.Valid ***REMOVED***
				assert.Equal(t, testCase.ExpectedStatus.Paused, status.Paused)
			***REMOVED***
			if testCase.ExpectedStatus.VUs.Valid ***REMOVED***
				assert.Equal(t, testCase.ExpectedStatus.VUs, status.VUs)
			***REMOVED***
			if testCase.ExpectedStatus.VUsMax.Valid ***REMOVED***
				assert.Equal(t, testCase.ExpectedStatus.VUsMax, status.VUsMax)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
