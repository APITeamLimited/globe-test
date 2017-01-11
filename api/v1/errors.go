package v1

import (
	"encoding/json"
	"github.com/manyminds/api2go"
	"net/http"
	"strconv"
)

func apiError(rw http.ResponseWriter, title, detail string, status int) ***REMOVED***
	doc := map[string][]api2go.Error***REMOVED***
		"errors": []api2go.Error***REMOVED***
			api2go.Error***REMOVED***
				Status: strconv.Itoa(status),
				Title:  title,
				Detail: detail,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	data, err := json.Marshal(doc)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	rw.WriteHeader(status)
	_, _ = rw.Write(data)
***REMOVED***
