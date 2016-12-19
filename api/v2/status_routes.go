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

package v2

import (
	"github.com/julienschmidt/httprouter"
	"github.com/loadimpact/k6/api/common"
	"github.com/manyminds/api2go/jsonapi"
	"io/ioutil"
	"net/http"
)

func HandleGetStatus(rw http.ResponseWriter, r *http.Request, p httprouter.Params) ***REMOVED***
	engine := common.GetEngine(r.Context())

	status := NewStatus(engine)
	data, err := jsonapi.Marshal(status)
	if err != nil ***REMOVED***
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	***REMOVED***
	_, _ = rw.Write(data)
***REMOVED***

func HandlePatchStatus(rw http.ResponseWriter, r *http.Request, p httprouter.Params) ***REMOVED***
	engine := common.GetEngine(r.Context())

	body, err := ioutil.ReadAll(r.Body)
	if err != nil ***REMOVED***
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	***REMOVED***

	var status Status
	if err := jsonapi.Unmarshal(body, &status); err != nil ***REMOVED***
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	***REMOVED***

	if status.VUsMax.Valid ***REMOVED***
		if err := engine.SetVUsMax(status.VUsMax.Int64); err != nil ***REMOVED***
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		***REMOVED***
	***REMOVED***
	if status.VUs.Valid ***REMOVED***
		if err := engine.SetVUs(status.VUs.Int64); err != nil ***REMOVED***
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		***REMOVED***
	***REMOVED***
	if status.Paused.Valid ***REMOVED***
		engine.SetPaused(status.Paused.Bool)
	***REMOVED***
***REMOVED***
