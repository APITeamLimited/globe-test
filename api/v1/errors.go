package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Error is an api error
type Error struct ***REMOVED***
	Status string `json:"status,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
***REMOVED***

func (e Error) Error() string ***REMOVED***
	return fmt.Sprintf("%s: %s", e.Title, e.Detail)
***REMOVED***

// ErrorResponse is a struct wrapper around multiple errors
type ErrorResponse struct ***REMOVED***
	Errors []Error `json:"errors"`
***REMOVED***

func apiError(rw http.ResponseWriter, title, detail string, status int) ***REMOVED***
	doc := ErrorResponse***REMOVED***
		Errors: []Error***REMOVED***
			***REMOVED***
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
