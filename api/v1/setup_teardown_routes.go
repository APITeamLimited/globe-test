/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2018 Load Impact
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

package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/manyminds/api2go/jsonapi"

	"github.com/loadimpact/k6/api/common"
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

// GetName is a dummy method so we can satisfy the jsonapi.EntityNamer interface
func (sd SetupData) GetName() string ***REMOVED***
	return "setupData"
***REMOVED***

// GetID is a dummy method so we can satisfy the jsonapi.MarshalIdentifier interface
func (sd SetupData) GetID() string ***REMOVED***
	return "default"
***REMOVED***

func handleSetupDataOutput(rw http.ResponseWriter, setupData json.RawMessage) ***REMOVED***
	rw.Header().Set("Content-Type", "application/json")
	var err error
	var data []byte

	if setupData == nil ***REMOVED***
		data, err = jsonapi.Marshal(NullSetupData***REMOVED***Data: nil***REMOVED***)
	***REMOVED*** else ***REMOVED***
		data, err = jsonapi.Marshal(SetupData***REMOVED***setupData***REMOVED***)
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
