package api2go

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

// HTTPError is used for errors
type HTTPError struct ***REMOVED***
	err    error
	msg    string
	status int
	Errors []Error `json:"errors,omitempty"`
***REMOVED***

// Error can be used for all kind of application errors
// e.g. you would use it to define form errors or any
// other semantical application problems
// for more information see http://jsonapi.org/format/#errors
type Error struct ***REMOVED***
	ID     string       `json:"id,omitempty"`
	Links  *ErrorLinks  `json:"links,omitempty"`
	Status string       `json:"status,omitempty"`
	Code   string       `json:"code,omitempty"`
	Title  string       `json:"title,omitempty"`
	Detail string       `json:"detail,omitempty"`
	Source *ErrorSource `json:"source,omitempty"`
	Meta   interface***REMOVED******REMOVED***  `json:"meta,omitempty"`
***REMOVED***

// ErrorLinks is used to provide an About URL that leads to
// further details about the particular occurrence of the problem.
//
// for more information see http://jsonapi.org/format/#error-objects
type ErrorLinks struct ***REMOVED***
	About string `json:"about,omitempty"`
***REMOVED***

// ErrorSource is used to provide references to the source of an error.
//
// The Pointer is a JSON Pointer to the associated entity in the request
// document.
// The Paramter is a string indicating which query parameter caused the error.
//
// for more information see http://jsonapi.org/format/#error-objects
type ErrorSource struct ***REMOVED***
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
***REMOVED***

// marshalHTTPError marshals an internal httpError
func marshalHTTPError(input HTTPError) string ***REMOVED***
	if len(input.Errors) == 0 ***REMOVED***
		input.Errors = []Error***REMOVED******REMOVED***Title: input.msg, Status: strconv.Itoa(input.status)***REMOVED******REMOVED***
	***REMOVED***

	data, err := json.Marshal(input)

	if err != nil ***REMOVED***
		log.Println(err)
		return "***REMOVED******REMOVED***"
	***REMOVED***

	return string(data)
***REMOVED***

// NewHTTPError creates a new error with message and status code.
// `err` will be logged (but never sent to a client), `msg` will be sent and `status` is the http status code.
// `err` can be nil.
func NewHTTPError(err error, msg string, status int) HTTPError ***REMOVED***
	return HTTPError***REMOVED***err: err, msg: msg, status: status***REMOVED***
***REMOVED***

// Error returns a nice string represenation including the status
func (e HTTPError) Error() string ***REMOVED***
	msg := fmt.Sprintf("http error (%d) %s and %d more errors", e.status, e.msg, len(e.Errors))
	if e.err != nil ***REMOVED***
		msg += ", " + e.err.Error()
	***REMOVED***

	return msg
***REMOVED***
