/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
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

package ui

import (
	"strings"

	"github.com/pkg/errors"
)

var _ Field = StringField***REMOVED******REMOVED***

type StringField struct ***REMOVED***
	Key     string
	Label   string
	Default string

	// Length constraints.
	Min, Max int
***REMOVED***

func (f StringField) GetKey() string ***REMOVED***
	return f.Key
***REMOVED***

func (f StringField) GetLabel() string ***REMOVED***
	return f.Label
***REMOVED***

func (f StringField) GetLabelExtra() string ***REMOVED***
	return f.Default
***REMOVED***

func (f StringField) Clean(s string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	s = strings.TrimSpace(s)
	if f.Min != 0 && len(s) < f.Min ***REMOVED***
		return nil, errors.Errorf("invalid input, min length is %d", f.Min)
	***REMOVED***
	if f.Max != 0 && len(s) > f.Max ***REMOVED***
		return nil, errors.Errorf("invalid input, max length is %d", f.Max)
	***REMOVED***
	if s == "" ***REMOVED***
		s = f.Default
	***REMOVED***
	return s, nil
***REMOVED***
