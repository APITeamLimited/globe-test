package api2go

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/manyminds/api2go/jsonapi"
)

// The Response struct implements api2go.Responder and can be used as a default
// implementation for your responses
// you can fill the field `Meta` with all the metadata your application needs
// like license, tokens, etc
type Response struct ***REMOVED***
	Res        interface***REMOVED******REMOVED***
	Code       int
	Meta       map[string]interface***REMOVED******REMOVED***
	Pagination Pagination
***REMOVED***

// Metadata returns additional meta data
func (r Response) Metadata() map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	return r.Meta
***REMOVED***

// Result returns the actual payload
func (r Response) Result() interface***REMOVED******REMOVED*** ***REMOVED***
	return r.Res
***REMOVED***

// StatusCode sets the return status code
func (r Response) StatusCode() int ***REMOVED***
	return r.Code
***REMOVED***

func buildLink(base string, r *http.Request, pagination map[string]string) jsonapi.Link ***REMOVED***
	params := r.URL.Query()
	for k, v := range pagination ***REMOVED***
		qk := fmt.Sprintf("page[%s]", k)
		params.Set(qk, v)
	***REMOVED***
	if len(params) == 0 ***REMOVED***
		return jsonapi.Link***REMOVED***Href: base***REMOVED***
	***REMOVED***
	query, _ := url.QueryUnescape(params.Encode())
	return jsonapi.Link***REMOVED***Href: fmt.Sprintf("%s?%s", base, query)***REMOVED***
***REMOVED***

// Links returns a jsonapi.Links object to include in the top-level response
func (r Response) Links(req *http.Request, baseURL string) (ret jsonapi.Links) ***REMOVED***
	ret = make(jsonapi.Links)

	if r.Pagination.Next != nil ***REMOVED***
		ret["next"] = buildLink(baseURL, req, r.Pagination.Next)
	***REMOVED***
	if r.Pagination.Prev != nil ***REMOVED***
		ret["prev"] = buildLink(baseURL, req, r.Pagination.Prev)
	***REMOVED***
	if r.Pagination.First != nil ***REMOVED***
		ret["first"] = buildLink(baseURL, req, r.Pagination.First)
	***REMOVED***
	if r.Pagination.Last != nil ***REMOVED***
		ret["last"] = buildLink(baseURL, req, r.Pagination.Last)
	***REMOVED***
	return
***REMOVED***
