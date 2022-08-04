// Package v1 implements the v1 of the k6's REST API
package v1

import (
	"net/http"
)

func NewHandler() http.Handler ***REMOVED***
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/status", func(rw http.ResponseWriter, r *http.Request) ***REMOVED***
		switch r.Method ***REMOVED***
		case http.MethodGet:
			handleGetStatus(rw, r)
		case http.MethodPatch:
			handlePatchStatus(rw, r)
		default:
			rw.WriteHeader(http.StatusMethodNotAllowed)
		***REMOVED***
	***REMOVED***)

	mux.HandleFunc("/v1/metrics", func(rw http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method != http.MethodGet ***REMOVED***
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		***REMOVED***
		handleGetMetrics(rw, r)
	***REMOVED***)

	mux.HandleFunc("/v1/metrics/", func(rw http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method != http.MethodGet ***REMOVED***
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		***REMOVED***

		id := r.URL.Path[len("/v1/metrics/"):]
		handleGetMetric(rw, r, id)
	***REMOVED***)

	mux.HandleFunc("/v1/groups", func(rw http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method != http.MethodGet ***REMOVED***
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		***REMOVED***

		handleGetGroups(rw, r)
	***REMOVED***)

	mux.HandleFunc("/v1/groups/", func(rw http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method != http.MethodGet ***REMOVED***
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		***REMOVED***

		id := r.URL.Path[len("/v1/groups/"):]
		handleGetGroup(rw, r, id)
	***REMOVED***)

	mux.HandleFunc("/v1/setup", func(rw http.ResponseWriter, r *http.Request) ***REMOVED***
		switch r.Method ***REMOVED***
		case http.MethodPost:
			handleRunSetup(rw, r)
		case http.MethodPut:
			handleSetSetupData(rw, r)
		case http.MethodGet:
			handleGetSetupData(rw, r)
		default:
			rw.WriteHeader(http.StatusMethodNotAllowed)
		***REMOVED***
	***REMOVED***)

	mux.HandleFunc("/v1/teardown", func(rw http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method != http.MethodPost ***REMOVED***
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		***REMOVED***

		handleRunTeardown(rw, r)
	***REMOVED***)

	return mux
***REMOVED***
