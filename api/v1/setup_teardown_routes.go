package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go.k6.io/k6/api/common"
)

// NullSetupData is wrapper around null to satisfy jsonapi
type NullSetupData struct ***REMOVED***
	SetupData
	Data interface***REMOVED******REMOVED*** `json:"data,omitempty" yaml:"data"`
***REMOVED***

// SetupData is just a simple wrapper to satisfy jsonapi
type SetupData struct ***REMOVED***
	Data interface***REMOVED******REMOVED*** `json:"data" yaml:"data"`
***REMOVED***

func handleSetupDataOutput(rw http.ResponseWriter, setupData json.RawMessage) ***REMOVED***
	rw.Header().Set("Content-Type", "application/json")
	var err error
	var data []byte

	if setupData == nil ***REMOVED***
		data, err = json.Marshal(newSetUpJSONAPI(NullSetupData***REMOVED***Data: nil***REMOVED***))
	***REMOVED*** else ***REMOVED***
		data, err = json.Marshal(newSetUpJSONAPI(SetupData***REMOVED***setupData***REMOVED***))
	***REMOVED***
	if err != nil ***REMOVED***
		apiError(rw, "Encoding error", err.Error(), http.StatusInternalServerError)
		return
	***REMOVED***

	_, _ = rw.Write(data)
***REMOVED***

// handleGetSetupData just returns the current JSON-encoded setup data
func handleGetSetupData(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	runner := common.GetEngine(r.Context()).ExecutionScheduler.GetRunner()
	handleSetupDataOutput(rw, runner.GetSetupData())
***REMOVED***

// handleSetSetupData just parses the JSON request body and sets the result as setup data for the runner
func handleSetSetupData(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	body, err := ioutil.ReadAll(r.Body)
	if err != nil ***REMOVED***
		apiError(rw, "Error reading request body", err.Error(), http.StatusBadRequest)
		return
	***REMOVED***

	var data interface***REMOVED******REMOVED***
	if len(body) > 0 ***REMOVED***
		if err := json.Unmarshal(body, &data); err != nil ***REMOVED***
			apiError(rw, "Error parsing request body", err.Error(), http.StatusBadRequest)
			return
		***REMOVED***
	***REMOVED***

	runner := common.GetEngine(r.Context()).ExecutionScheduler.GetRunner()

	if len(body) == 0 ***REMOVED***
		runner.SetSetupData(nil)
	***REMOVED*** else ***REMOVED***
		runner.SetSetupData(body)
	***REMOVED***

	handleSetupDataOutput(rw, runner.GetSetupData())
***REMOVED***

// handleRunSetup executes the runner's Setup() method and returns the result
func handleRunSetup(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	engine := common.GetEngine(r.Context())
	runner := engine.ExecutionScheduler.GetRunner()

	if err := runner.Setup(r.Context(), engine.Samples); err != nil ***REMOVED***
		apiError(rw, "Error executing setup", err.Error(), http.StatusInternalServerError)
		return
	***REMOVED***

	handleSetupDataOutput(rw, runner.GetSetupData())
***REMOVED***

// handleRunTeardown executes the runner's Teardown() method
func handleRunTeardown(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	engine := common.GetEngine(r.Context())
	runner := common.GetEngine(r.Context()).ExecutionScheduler.GetRunner()

	if err := runner.Teardown(r.Context(), engine.Samples); err != nil ***REMOVED***
		apiError(rw, "Error executing teardown", err.Error(), http.StatusInternalServerError)
	***REMOVED***
***REMOVED***
