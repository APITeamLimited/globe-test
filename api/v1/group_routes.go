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

	"go.k6.io/k6/api/common"
)

func handleGetGroups(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	engine := common.GetEngine(r.Context())

	root := NewGroup(engine.ExecutionScheduler.GetRunner().GetDefaultGroup(), nil)
	groups := FlattenGroup(root)

	data, err := json.Marshal(newGroupsJSONAPI(groups))
	if err != nil ***REMOVED***
		apiError(rw, "Encoding error", err.Error(), http.StatusInternalServerError)
		return
	***REMOVED***
	_, _ = rw.Write(data)
***REMOVED***

func handleGetGroup(rw http.ResponseWriter, r *http.Request, id string) ***REMOVED***
	engine := common.GetEngine(r.Context())

	root := NewGroup(engine.ExecutionScheduler.GetRunner().GetDefaultGroup(), nil)
	groups := FlattenGroup(root)

	var group *Group
	for _, g := range groups ***REMOVED***
		if g.ID == id ***REMOVED***
			group = g
			break
		***REMOVED***
	***REMOVED***
	if group == nil ***REMOVED***
		apiError(rw, "Not Found", "No group with that ID was found", http.StatusNotFound)
		return
	***REMOVED***

	data, err := json.Marshal(newGroupJSONAPI(group))
	if err != nil ***REMOVED***
		apiError(rw, "Encoding error", err.Error(), http.StatusInternalServerError)
		return
	***REMOVED***
	_, _ = rw.Write(data)
***REMOVED***
