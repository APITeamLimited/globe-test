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

package v2

import (
	"gopkg.in/guregu/null.v3"
)

type Status struct ***REMOVED***
	Running null.Bool `json:"running"`
	VUs     null.Int  `json:"vus"`
	VUsMax  null.Int  `json:"vus-max"`

	// Readonly.
	Tainted bool `json:"tainted"`
***REMOVED***

func (s Status) GetName() string ***REMOVED***
	return "status"
***REMOVED***

func (s Status) GetID() string ***REMOVED***
	return "default"
***REMOVED***

func (s Status) SetID(id string) error ***REMOVED***
	return nil
***REMOVED***
