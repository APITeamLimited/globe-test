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
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/loadimpact/k6/api/common"
	"github.com/manyminds/api2go/jsonapi"
)

func HandleGetMetrics(rw http.ResponseWriter, r *http.Request, p httprouter.Params) ***REMOVED***
	engine := common.GetEngine(r.Context())

	var t time.Duration
	if engine.Executor != nil ***REMOVED***
		t = engine.Executor.GetState().GetCurrentTestRunDuration()
	***REMOVED***

	metrics := make([]Metric, 0)
	for _, m := range engine.Metrics ***REMOVED***
		metrics = append(metrics, NewMetric(m, t))
	***REMOVED***

	data, err := jsonapi.Marshal(metrics)
	if err != nil ***REMOVED***
		apiError(rw, "Encoding error", err.Error(), http.StatusInternalServerError)
		return
	***REMOVED***
	_, _ = rw.Write(data)
***REMOVED***

func HandleGetMetric(rw http.ResponseWriter, r *http.Request, p httprouter.Params) ***REMOVED***
	id := p.ByName("id")
	engine := common.GetEngine(r.Context())

	var t time.Duration
	if engine.Executor != nil ***REMOVED***
		t = engine.Executor.GetState().GetCurrentTestRunDuration()
	***REMOVED***

	var metric Metric
	var found bool
	for _, m := range engine.Metrics ***REMOVED***
		if m.Name == id ***REMOVED***
			metric = NewMetric(m, t)
			found = true
			break
		***REMOVED***
	***REMOVED***

	if !found ***REMOVED***
		apiError(rw, "Not Found", "No metric with that ID was found", http.StatusNotFound)
		return
	***REMOVED***

	data, err := jsonapi.Marshal(metric)
	if err != nil ***REMOVED***
		apiError(rw, "Encoding error", err.Error(), http.StatusInternalServerError)
		return
	***REMOVED***
	_, _ = rw.Write(data)
***REMOVED***
