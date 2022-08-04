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
