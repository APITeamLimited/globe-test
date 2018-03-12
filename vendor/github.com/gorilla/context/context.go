// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package context

import (
	"net/http"
	"sync"
	"time"
)

var (
	mutex sync.RWMutex
	data  = make(map[*http.Request]map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***)
	datat = make(map[*http.Request]int64)
)

// Set stores a value for a given key in a given request.
func Set(r *http.Request, key, val interface***REMOVED******REMOVED***) ***REMOVED***
	mutex.Lock()
	if data[r] == nil ***REMOVED***
		data[r] = make(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***)
		datat[r] = time.Now().Unix()
	***REMOVED***
	data[r][key] = val
	mutex.Unlock()
***REMOVED***

// Get returns a value stored for a given key in a given request.
func Get(r *http.Request, key interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	mutex.RLock()
	if ctx := data[r]; ctx != nil ***REMOVED***
		value := ctx[key]
		mutex.RUnlock()
		return value
	***REMOVED***
	mutex.RUnlock()
	return nil
***REMOVED***

// GetOk returns stored value and presence state like multi-value return of map access.
func GetOk(r *http.Request, key interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, bool) ***REMOVED***
	mutex.RLock()
	if _, ok := data[r]; ok ***REMOVED***
		value, ok := data[r][key]
		mutex.RUnlock()
		return value, ok
	***REMOVED***
	mutex.RUnlock()
	return nil, false
***REMOVED***

// GetAll returns all stored values for the request as a map. Nil is returned for invalid requests.
func GetAll(r *http.Request) map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED*** ***REMOVED***
	mutex.RLock()
	if context, ok := data[r]; ok ***REMOVED***
		result := make(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***, len(context))
		for k, v := range context ***REMOVED***
			result[k] = v
		***REMOVED***
		mutex.RUnlock()
		return result
	***REMOVED***
	mutex.RUnlock()
	return nil
***REMOVED***

// GetAllOk returns all stored values for the request as a map and a boolean value that indicates if
// the request was registered.
func GetAllOk(r *http.Request) (map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***, bool) ***REMOVED***
	mutex.RLock()
	context, ok := data[r]
	result := make(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***, len(context))
	for k, v := range context ***REMOVED***
		result[k] = v
	***REMOVED***
	mutex.RUnlock()
	return result, ok
***REMOVED***

// Delete removes a value stored for a given key in a given request.
func Delete(r *http.Request, key interface***REMOVED******REMOVED***) ***REMOVED***
	mutex.Lock()
	if data[r] != nil ***REMOVED***
		delete(data[r], key)
	***REMOVED***
	mutex.Unlock()
***REMOVED***

// Clear removes all values stored for a given request.
//
// This is usually called by a handler wrapper to clean up request
// variables at the end of a request lifetime. See ClearHandler().
func Clear(r *http.Request) ***REMOVED***
	mutex.Lock()
	clear(r)
	mutex.Unlock()
***REMOVED***

// clear is Clear without the lock.
func clear(r *http.Request) ***REMOVED***
	delete(data, r)
	delete(datat, r)
***REMOVED***

// Purge removes request data stored for longer than maxAge, in seconds.
// It returns the amount of requests removed.
//
// If maxAge <= 0, all request data is removed.
//
// This is only used for sanity check: in case context cleaning was not
// properly set some request data can be kept forever, consuming an increasing
// amount of memory. In case this is detected, Purge() must be called
// periodically until the problem is fixed.
func Purge(maxAge int) int ***REMOVED***
	mutex.Lock()
	count := 0
	if maxAge <= 0 ***REMOVED***
		count = len(data)
		data = make(map[*http.Request]map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***)
		datat = make(map[*http.Request]int64)
	***REMOVED*** else ***REMOVED***
		min := time.Now().Unix() - int64(maxAge)
		for r := range data ***REMOVED***
			if datat[r] < min ***REMOVED***
				clear(r)
				count++
			***REMOVED***
		***REMOVED***
	***REMOVED***
	mutex.Unlock()
	return count
***REMOVED***

// ClearHandler wraps an http.Handler and clears request values at the end
// of a request lifetime.
func ClearHandler(h http.Handler) http.Handler ***REMOVED***
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		defer Clear(r)
		h.ServeHTTP(w, r)
	***REMOVED***)
***REMOVED***
