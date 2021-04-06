/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2017 Load Impact
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

package cloudapi

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrNotAuthorized    = errors.New("Not allowed to upload result to Load Impact cloud")
	ErrNotAuthenticated = errors.New("Failed to authenticate with Load Impact cloud")
	ErrUnknown          = errors.New("An error occurred talking to Load Impact cloud")
)

// ErrorResponse represents an error cause by talking to the API
type ErrorResponse struct ***REMOVED***
	Response *http.Response `json:"-"`

	Code        int                 `json:"code"`
	Message     string              `json:"message"`
	Details     map[string][]string `json:"details"`
	FieldErrors map[string][]string `json:"field_errors"`
	Errors      []string            `json:"errors"`
***REMOVED***

func contains(s []string, e string) bool ***REMOVED***
	for _, a := range s ***REMOVED***
		if a == e ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (e ErrorResponse) Error() string ***REMOVED***
	msg := e.Message

	for _, v := range e.Errors ***REMOVED***
		// atm: `errors` and `message` could be duplicated
		// TODO: remove condition when the API changes
		if v != msg ***REMOVED***
			msg += "\n " + v
		***REMOVED***
	***REMOVED***

	// `e.Details` is the old API version
	// TODO: do not handle `details` when the old API becomes obsolete
	var details []string
	var detail string
	for k, v := range e.Details ***REMOVED***
		detail = k + ": " + strings.Join(v, ", ")
		details = append(details, detail)
	***REMOVED***

	for k, v := range e.FieldErrors ***REMOVED***
		detail = k + ": " + strings.Join(v, ", ")
		// atm: `details` and `field_errors` could be duplicated
		if !contains(details, detail) ***REMOVED***
			details = append(details, detail)
		***REMOVED***
	***REMOVED***

	if len(details) > 0 ***REMOVED***
		msg += "\n " + strings.Join(details, "\n")
	***REMOVED***

	var code string
	if e.Code > 0 && e.Response != nil ***REMOVED***
		code = fmt.Sprintf("%d/E%d", e.Response.StatusCode, e.Code)
	***REMOVED*** else if e.Response != nil ***REMOVED***
		code = fmt.Sprintf("%d", e.Response.StatusCode)
	***REMOVED*** else if e.Code > 0 ***REMOVED***
		code = fmt.Sprintf("E%d", e.Code)
	***REMOVED***

	if len(code) > 0 ***REMOVED***
		msg = fmt.Sprintf("(%s) %s", code, msg)
	***REMOVED***

	return msg
***REMOVED***
