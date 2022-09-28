// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Based on github.com/aws/aws-sdk-go by Amazon.com, Inc. with code from:
// - github.com/aws/aws-sdk-go/blob/v1.34.28/aws/signer/v4/header_rules.go
// - github.com/aws/aws-sdk-go/blob/v1.34.28/internal/strings/strings.go
// See THIRD-PARTY-NOTICES for original license terms

package awsv4

import (
	"strings"
)

// validator houses a set of rule needed for validation of a
// string value
type rules []rule

// rule interface allows for more flexible rules and just simply
// checks whether or not a value adheres to that rule
type rule interface ***REMOVED***
	IsValid(value string) bool
***REMOVED***

// IsValid will iterate through all rules and see if any rules
// apply to the value and supports nested rules
func (r rules) IsValid(value string) bool ***REMOVED***
	for _, rule := range r ***REMOVED***
		if rule.IsValid(value) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// mapRule generic rule for maps
type mapRule map[string]struct***REMOVED******REMOVED***

// IsValid for the map rule satisfies whether it exists in the map
func (m mapRule) IsValid(value string) bool ***REMOVED***
	_, ok := m[value]
	return ok
***REMOVED***

// allowlist is a generic rule for allowlisting
type allowlist struct ***REMOVED***
	rule
***REMOVED***

// IsValid for allowlist checks if the value is within the allowlist
func (a allowlist) IsValid(value string) bool ***REMOVED***
	return a.rule.IsValid(value)
***REMOVED***

// denylist is a generic rule for denylisting
type denylist struct ***REMOVED***
	rule
***REMOVED***

// IsValid for allowlist checks if the value is within the allowlist
func (d denylist) IsValid(value string) bool ***REMOVED***
	return !d.rule.IsValid(value)
***REMOVED***

type patterns []string

// hasPrefixFold tests whether the string s begins with prefix, interpreted as UTF-8 strings,
// under Unicode case-folding.
func hasPrefixFold(s, prefix string) bool ***REMOVED***
	return len(s) >= len(prefix) && strings.EqualFold(s[0:len(prefix)], prefix)
***REMOVED***

// IsValid for patterns checks each pattern and returns if a match has
// been found
func (p patterns) IsValid(value string) bool ***REMOVED***
	for _, pattern := range p ***REMOVED***
		if hasPrefixFold(value, pattern) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// inclusiveRules rules allow for rules to depend on one another
type inclusiveRules []rule

// IsValid will return true if all rules are true
func (r inclusiveRules) IsValid(value string) bool ***REMOVED***
	for _, rule := range r ***REMOVED***
		if !rule.IsValid(value) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***
