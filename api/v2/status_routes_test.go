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
	"bytes"
	"encoding/json"
	"github.com/loadimpact/k6/lib"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStatus(t *testing.T) ***REMOVED***
	engine, err := lib.NewEngine(nil)
	assert.NoError(t, err)

	rw := httptest.NewRecorder()
	NewHandler().ServeHTTP(rw, newRequestWithEngine(engine, "GET", "/v2/status", nil))
	res := rw.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	t.Run("document", func(t *testing.T) ***REMOVED***
		var doc jsonapi.Document
		assert.NoError(t, json.Unmarshal(rw.Body.Bytes(), &doc))
		if !assert.NotNil(t, doc.Data.DataObject) ***REMOVED***
			return
		***REMOVED***
		assert.Equal(t, "status", doc.Data.DataObject.Type)
	***REMOVED***)

	t.Run("status", func(t *testing.T) ***REMOVED***
		var status Status
		assert.NoError(t, jsonapi.Unmarshal(rw.Body.Bytes(), &status))
		assert.True(t, status.Running.Valid)
		assert.True(t, status.VUs.Valid)
		assert.True(t, status.VUsMax.Valid)
		assert.False(t, status.Tainted)
	***REMOVED***)
***REMOVED***

func TestPatchStatus(t *testing.T) ***REMOVED***
	testdata := map[string]struct ***REMOVED***
		StatusCode int
		Status     Status
	***REMOVED******REMOVED***
		"nothing":      ***REMOVED***200, Status***REMOVED******REMOVED******REMOVED***,
		"running":      ***REMOVED***200, Status***REMOVED***Running: null.BoolFrom(true)***REMOVED******REMOVED***,
		"max vus":      ***REMOVED***200, Status***REMOVED***VUsMax: null.IntFrom(10)***REMOVED******REMOVED***,
		"too many vus": ***REMOVED***400, Status***REMOVED***VUs: null.IntFrom(10), VUsMax: null.IntFrom(0)***REMOVED******REMOVED***,

		// PANICS DUE TO ENGINE BUG!
		// "vus":          ***REMOVED***200, Status***REMOVED***VUs: null.IntFrom(10), VUsMax: null.IntFrom(10)***REMOVED******REMOVED***,
	***REMOVED***

	for name, indata := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			engine, err := lib.NewEngine(nil)
			assert.NoError(t, err)

			body, err := jsonapi.Marshal(indata.Status)
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***

			rw := httptest.NewRecorder()
			NewHandler().ServeHTTP(rw, newRequestWithEngine(engine, "PATCH", "/v2/status", bytes.NewReader(body)))
			res := rw.Result()

			if !assert.Equal(t, indata.StatusCode, res.StatusCode) ***REMOVED***
				return
			***REMOVED***
			if indata.StatusCode != 200 ***REMOVED***
				return
			***REMOVED***

			status := NewStatus(engine)
			if indata.Status.Running.Valid ***REMOVED***
				assert.Equal(t, indata.Status.Running, status.Running)
			***REMOVED***
			if indata.Status.VUs.Valid ***REMOVED***
				assert.Equal(t, indata.Status.VUs, status.VUs)
			***REMOVED***
			if indata.Status.VUsMax.Valid ***REMOVED***
				assert.Equal(t, indata.Status.VUsMax, status.VUsMax)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
