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

package har

import (
	"encoding/json"
	"io"
	"strings"
	"time"
)

// Define new types to sort
type EntryByStarted []*Entry

func (e EntryByStarted) Len() int ***REMOVED*** return len(e) ***REMOVED***

func (e EntryByStarted) Swap(i, j int) ***REMOVED*** e[i], e[j] = e[j], e[i] ***REMOVED***

func (e EntryByStarted) Less(i, j int) bool ***REMOVED***
	return e[i].StartedDateTime.Before(e[j].StartedDateTime)
***REMOVED***

type PageByStarted []Page

func (e PageByStarted) Len() int ***REMOVED*** return len(e) ***REMOVED***

func (e PageByStarted) Swap(i, j int) ***REMOVED*** e[i], e[j] = e[j], e[i] ***REMOVED***

func (e PageByStarted) Less(i, j int) bool ***REMOVED***
	return e[i].StartedDateTime.Before(e[j].StartedDateTime)
***REMOVED***

func Decode(r io.Reader) (HAR, error) ***REMOVED***
	var har HAR
	if err := json.NewDecoder(r).Decode(&har); err != nil ***REMOVED***
		return HAR***REMOVED******REMOVED***, err
	***REMOVED***

	return har, nil
***REMOVED***

// Returns true if the given url is allowed from the only (only domains) and skip (skip domains) values, otherwise false
func IsAllowedURL(url string, only, skip []string) bool ***REMOVED***
	if len(only) != 0 ***REMOVED***
		for _, v := range only ***REMOVED***
			v = strings.Trim(v, " ")
			if v != "" && strings.Contains(url, v) ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
		return false
	***REMOVED***
	if len(skip) != 0 ***REMOVED***
		for _, v := range skip ***REMOVED***
			v = strings.Trim(v, " ")
			if v != "" && strings.Contains(url, v) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func SplitEntriesInBatches(entries []*Entry, interval uint) [][]*Entry ***REMOVED***
	var r [][]*Entry
	r = append(r, []*Entry***REMOVED******REMOVED***)

	if interval > 0 && len(entries) > 1 ***REMOVED***
		j := 0
		d := time.Duration(interval) * time.Millisecond
		for i, e := range entries ***REMOVED***

			if i != 0 ***REMOVED***
				prev := entries[i-1]
				if e.StartedDateTime.Sub(prev.StartedDateTime) >= d ***REMOVED***
					r = append(r, []*Entry***REMOVED******REMOVED***)
					j++
				***REMOVED***
			***REMOVED***
			r[j] = append(r[j], e)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r[0] = entries
	***REMOVED***

	return r
***REMOVED***