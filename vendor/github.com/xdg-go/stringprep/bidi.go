// Copyright 2018 by David A. Golden. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package stringprep

var errHasLCat = "BiDi string can't have runes from category L"
var errFirstRune = "BiDi string first rune must have category R or AL"
var errLastRune = "BiDi string last rune must have category R or AL"

// Check for prohibited characters from table C.8
func checkBiDiProhibitedRune(s string) error ***REMOVED***
	for _, r := range s ***REMOVED***
		if TableC8.Contains(r) ***REMOVED***
			return Error***REMOVED***Msg: errProhibited, Rune: r***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Check for LCat characters from table D.2
func checkBiDiLCat(s string) error ***REMOVED***
	for _, r := range s ***REMOVED***
		if TableD2.Contains(r) ***REMOVED***
			return Error***REMOVED***Msg: errHasLCat, Rune: r***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Check first and last characters are in table D.1; requires non-empty string
func checkBadFirstAndLastRandALCat(s string) error ***REMOVED***
	rs := []rune(s)
	if !TableD1.Contains(rs[0]) ***REMOVED***
		return Error***REMOVED***Msg: errFirstRune, Rune: rs[0]***REMOVED***
	***REMOVED***
	n := len(rs) - 1
	if !TableD1.Contains(rs[n]) ***REMOVED***
		return Error***REMOVED***Msg: errLastRune, Rune: rs[n]***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Look for RandALCat characters from table D.1
func hasBiDiRandALCat(s string) bool ***REMOVED***
	for _, r := range s ***REMOVED***
		if TableD1.Contains(r) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// Check that BiDi rules are satisfied ; let empty string pass this rule
func passesBiDiRules(s string) error ***REMOVED***
	if len(s) == 0 ***REMOVED***
		return nil
	***REMOVED***
	if err := checkBiDiProhibitedRune(s); err != nil ***REMOVED***
		return err
	***REMOVED***
	if hasBiDiRandALCat(s) ***REMOVED***
		if err := checkBiDiLCat(s); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := checkBadFirstAndLastRandALCat(s); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
