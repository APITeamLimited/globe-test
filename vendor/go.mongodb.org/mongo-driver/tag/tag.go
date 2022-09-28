// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// Package tag provides a way to define filters for tagged servers.
package tag // import "go.mongodb.org/mongo-driver/tag"

import (
	"bytes"
	"fmt"
)

// Tag is a name/vlaue pair.
type Tag struct ***REMOVED***
	Name  string
	Value string
***REMOVED***

// String returns a human-readable human-readable description of the tag.
func (tag Tag) String() string ***REMOVED***
	return fmt.Sprintf("%s=%s", tag.Name, tag.Value)
***REMOVED***

// NewTagSetFromMap creates a new tag set from a map.
func NewTagSetFromMap(m map[string]string) Set ***REMOVED***
	var set Set
	for k, v := range m ***REMOVED***
		set = append(set, Tag***REMOVED***Name: k, Value: v***REMOVED***)
	***REMOVED***

	return set
***REMOVED***

// NewTagSetsFromMaps creates new tag sets from maps.
func NewTagSetsFromMaps(maps []map[string]string) []Set ***REMOVED***
	sets := make([]Set, 0, len(maps))
	for _, m := range maps ***REMOVED***
		sets = append(sets, NewTagSetFromMap(m))
	***REMOVED***
	return sets
***REMOVED***

// Set is an ordered list of Tags.
type Set []Tag

// Contains indicates whether the name/value pair exists in the tagset.
func (ts Set) Contains(name, value string) bool ***REMOVED***
	for _, t := range ts ***REMOVED***
		if t.Name == name && t.Value == value ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

// ContainsAll indicates whether all the name/value pairs exist in the tagset.
func (ts Set) ContainsAll(other []Tag) bool ***REMOVED***
	for _, ot := range other ***REMOVED***
		if !ts.Contains(ot.Name, ot.Value) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

// String returns a human-readable human-readable description of the tagset.
func (ts Set) String() string ***REMOVED***
	var b bytes.Buffer
	for i, tag := range ts ***REMOVED***
		if i > 0 ***REMOVED***
			b.WriteString(",")
		***REMOVED***
		b.WriteString(tag.String())
	***REMOVED***
	return b.String()
***REMOVED***
