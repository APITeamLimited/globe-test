/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
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
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/loadimpact/k6/api/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/scheduler"
	"github.com/manyminds/api2go/jsonapi"
)

func HandleGetStatus(rw http.ResponseWriter, r *http.Request, p httprouter.Params) ***REMOVED***
	engine := common.GetEngine(r.Context())

	status := NewStatus(engine)
	data, err := jsonapi.Marshal(status)
	if err != nil ***REMOVED***
		apiError(rw, "Encoding error", err.Error(), http.StatusInternalServerError)
		return
	***REMOVED***
	_, _ = rw.Write(data)
***REMOVED***

func getFirstManualExecutionScheduler(executor lib.Executor) (*scheduler.ManualExecution, error) ***REMOVED***
	schedulers := executor.GetSchedulers()
	for _, s := range schedulers ***REMOVED***
		if mex, ok := s.(*scheduler.ManualExecution); ok ***REMOVED***
			return mex, nil
		***REMOVED***
	***REMOVED***
	return nil, fmt.Errorf("a manual-execution scheduler needs to be configured for live configuration updates")
***REMOVED***

func HandlePatchStatus(rw http.ResponseWriter, r *http.Request, p httprouter.Params) ***REMOVED***
	engine := common.GetEngine(r.Context())

	body, err := ioutil.ReadAll(r.Body)
	if err != nil ***REMOVED***
		apiError(rw, "Couldn't read request", err.Error(), http.StatusBadRequest)
		return
	***REMOVED***

	var status Status
	if err := jsonapi.Unmarshal(body, &status); err != nil ***REMOVED***
		apiError(rw, "Invalid data", err.Error(), http.StatusBadRequest)
		return
	***REMOVED***

	if status.Paused.Valid ***REMOVED***
		if err = engine.Executor.SetPaused(status.Paused.Bool); err != nil ***REMOVED***
			apiError(rw, "Pause error", err.Error(), http.StatusInternalServerError)
			return
		***REMOVED***
	***REMOVED***

	if status.VUsMax.Valid || status.VUs.Valid ***REMOVED***
		//TODO: add ability to specify the actual scheduler id? though thus should
		//likely be in the v2 REST API, where we could implement it in a way that
		//may allow us to eventually support other scheduler types
		scheduler, uptateErr := getFirstManualExecutionScheduler(engine.Executor)
		if uptateErr != nil ***REMOVED***
			apiError(rw, "Execution config error", uptateErr.Error(), http.StatusInternalServerError)
			return
		***REMOVED***
		newConfig := scheduler.GetCurrentConfig().ManualExecutionControlConfig
		if status.VUsMax.Valid ***REMOVED***
			newConfig.MaxVUs = status.VUsMax
		***REMOVED***
		if status.VUs.Valid ***REMOVED***
			newConfig.VUs = status.VUs
		***REMOVED***
		if uptateErr := scheduler.UpdateConfig(r.Context(), newConfig); err != nil ***REMOVED***
			apiError(rw, "Config update error", uptateErr.Error(), http.StatusInternalServerError)
			return
		***REMOVED***
	***REMOVED***
	data, err := jsonapi.Marshal(NewStatus(engine))
	if err != nil ***REMOVED***
		apiError(rw, "Encoding error", err.Error(), http.StatusInternalServerError)
		return
	***REMOVED***
	_, _ = rw.Write(data)
***REMOVED***
