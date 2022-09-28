// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package internal

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// StringSliceFromRawElement decodes the provided BSON element into a []string. This internally calls
// StringSliceFromRawValue on the element's value. The error conditions outlined in that function's documentation
// apply for this function as well.
func StringSliceFromRawElement(element bson.RawElement) ([]string, error) ***REMOVED***
	return StringSliceFromRawValue(element.Key(), element.Value())
***REMOVED***

// StringSliceFromRawValue decodes the provided BSON value into a []string. This function returns an error if the value
// is not an array or any of the elements in the array are not strings. The name parameter is used to add context to
// error messages.
func StringSliceFromRawValue(name string, val bson.RawValue) ([]string, error) ***REMOVED***
	arr, ok := val.ArrayOK()
	if !ok ***REMOVED***
		return nil, fmt.Errorf("expected '%s' to be an array but it's a BSON %s", name, val.Type)
	***REMOVED***

	arrayValues, err := arr.Values()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	strs := make([]string, 0, len(arrayValues))
	for _, arrayVal := range arrayValues ***REMOVED***
		str, ok := arrayVal.StringValueOK()
		if !ok ***REMOVED***
			return nil, fmt.Errorf("expected '%s' to be an array of strings, but found a BSON %s", name, arrayVal.Type)
		***REMOVED***
		strs = append(strs, str)
	***REMOVED***
	return strs, nil
***REMOVED***
