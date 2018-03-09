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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/loadimpact/k6/core"
	"github.com/loadimpact/k6/core/local"
	"github.com/loadimpact/k6/lib"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
)

type groupDummyRunner struct ***REMOVED***
	Group *lib.Group
***REMOVED***

func (r groupDummyRunner) MakeArchive() *lib.Archive ***REMOVED*** return nil ***REMOVED***

func (r groupDummyRunner) NewVU() (lib.VU, error) ***REMOVED*** return nil, nil ***REMOVED***

func (r groupDummyRunner) GetDefaultGroup() *lib.Group ***REMOVED*** return r.Group ***REMOVED***

func (r groupDummyRunner) GetOptions() lib.Options ***REMOVED*** return lib.Options***REMOVED******REMOVED*** ***REMOVED***

func (r groupDummyRunner) SetOptions(opts lib.Options) ***REMOVED******REMOVED***

func (r groupDummyRunner) SetCollectorOptions(opts lib.CollectorOptions) ***REMOVED******REMOVED***

func TestGetGroups(t *testing.T) ***REMOVED***
	g0, err := lib.NewGroup("", nil)
	assert.NoError(t, err)
	g1, err := g0.Group("group 1")
	assert.NoError(t, err)
	g2, err := g1.Group("group 2")
	assert.NoError(t, err)

	engine, err := core.NewEngine(local.New(groupDummyRunner***REMOVED***g0***REMOVED***), lib.Options***REMOVED******REMOVED***)
	assert.NoError(t, err)

	t.Run("list", func(t *testing.T) ***REMOVED***
		rw := httptest.NewRecorder()
		NewHandler().ServeHTTP(rw, newRequestWithEngine(engine, "GET", "/v1/groups", nil))
		res := rw.Result()
		body := rw.Body.Bytes()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.NotEmpty(t, body)

		t.Run("document", func(t *testing.T) ***REMOVED***
			var doc jsonapi.Document
			assert.NoError(t, json.Unmarshal(body, &doc))
			assert.Nil(t, doc.Data.DataObject)
			if assert.NotEmpty(t, doc.Data.DataArray) ***REMOVED***
				assert.Equal(t, "groups", doc.Data.DataArray[0].Type)
			***REMOVED***
		***REMOVED***)

		// t.Run("groups", func(t *testing.T) ***REMOVED***
		// 	var groups []Group
		// 	assert.NoError(t, jsonapi.Unmarshal(body, &groups))
		// 	if assert.Len(t, groups, 3) ***REMOVED***
		// 		for _, g := range groups ***REMOVED***
		// 			switch g.ID ***REMOVED***
		// 			case g0.ID:
		// 				assert.Equal(t, "", g.Name)
		// 				assert.Nil(t, g.Parent)
		// 				assert.Equal(t, "", g.ParentID)
		// 				assert.Len(t, g.GroupIDs, 1)
		// 				assert.EqualValues(t, []string***REMOVED***g1.ID***REMOVED***, g.GroupIDs)
		// 			case g1.ID:
		// 				assert.Equal(t, "group 1", g.Name)
		// 				assert.Nil(t, g.Parent)
		// 				assert.Equal(t, g0.ID, g.ParentID)
		// 				assert.EqualValues(t, []string***REMOVED***g2.ID***REMOVED***, g.GroupIDs)
		// 			case g2.ID:
		// 				assert.Equal(t, "group 2", g.Name)
		// 				assert.Nil(t, g.Parent)
		// 				assert.Equal(t, g1.ID, g.ParentID)
		// 				assert.EqualValues(t, []string***REMOVED******REMOVED***, g.GroupIDs)
		// 			default:
		// 				assert.Fail(t, "Unknown ID: "+g.ID)
		// 			***REMOVED***
		// 		***REMOVED***
		// 	***REMOVED***
		// ***REMOVED***)
	***REMOVED***)
	for _, gp := range []*lib.Group***REMOVED***g0, g1, g2***REMOVED*** ***REMOVED***
		t.Run(gp.Name, func(t *testing.T) ***REMOVED***
			rw := httptest.NewRecorder()
			NewHandler().ServeHTTP(rw, newRequestWithEngine(engine, "GET", "/v1/groups/"+gp.ID, nil))
			res := rw.Result()
			assert.Equal(t, http.StatusOK, res.StatusCode)
		***REMOVED***)
	***REMOVED***
***REMOVED***
