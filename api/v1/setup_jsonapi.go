/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2021 Load Impact
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

package v1

type setUpJSONAPI struct ***REMOVED***
	Data setUpData `json:"data"`
***REMOVED***

type setUpData struct ***REMOVED***
	Type       string      `json:"type"`
	ID         string      `json:"id"`
	Attributes interface***REMOVED******REMOVED*** `json:"attributes"`
***REMOVED***

func newSetUpJSONAPI(setup interface***REMOVED******REMOVED***) setUpJSONAPI ***REMOVED***
	return setUpJSONAPI***REMOVED***
		Data: setUpData***REMOVED***
			Type:       "setupData",
			ID:         "default",
			Attributes: setup,
		***REMOVED***,
	***REMOVED***
***REMOVED***
