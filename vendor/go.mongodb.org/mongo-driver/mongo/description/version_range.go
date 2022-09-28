// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package description

import "fmt"

// VersionRange represents a range of versions.
type VersionRange struct ***REMOVED***
	Min int32
	Max int32
***REMOVED***

// NewVersionRange creates a new VersionRange given a min and a max.
func NewVersionRange(min, max int32) VersionRange ***REMOVED***
	return VersionRange***REMOVED***Min: min, Max: max***REMOVED***
***REMOVED***

// Includes returns a bool indicating whether the supplied integer is included
// in the range.
func (vr VersionRange) Includes(v int32) bool ***REMOVED***
	return v >= vr.Min && v <= vr.Max
***REMOVED***

// Equals returns a bool indicating whether the supplied VersionRange is equal.
func (vr *VersionRange) Equals(other *VersionRange) bool ***REMOVED***
	if vr == nil && other == nil ***REMOVED***
		return true
	***REMOVED***
	if vr == nil || other == nil ***REMOVED***
		return false
	***REMOVED***
	return vr.Min == other.Min && vr.Max == other.Max
***REMOVED***

// String implements the fmt.Stringer interface.
func (vr VersionRange) String() string ***REMOVED***
	return fmt.Sprintf("[%d, %d]", vr.Min, vr.Max)
***REMOVED***
