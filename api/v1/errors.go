package v1

import (
	"encoding/json"
	"fmt"
	"github.com/manyminds/api2go"
	"net/http"
	"strconv"
)

type Error api2go.Error

func (e Error) Error() string ***REMOVED***
	return fmt.Sprintf("%s: %s", e.Title, e.Detail)
***REMOVED***

type ErrorResponse struct ***REMOVED***
	Errors []Error `json:"errors"`
***REMOVED***

func apiError(rw http.ResponseWriter, title, detail string, status int) ***REMOVED***
	doc := ErrorResponse***REMOVED***
		Errors: []Error***REMOVED***
			Error***REMOVED***
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
