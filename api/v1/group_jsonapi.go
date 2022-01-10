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

import "encoding/json"

type groupJSONAPI struct ***REMOVED***
	Data groupData `json:"data"`
***REMOVED***

type groupsJSONAPI struct ***REMOVED***
	Data []groupData `json:"data"`
***REMOVED***

type groupData struct ***REMOVED***
	Type          string         `json:"type"`
	ID            string         `json:"id"`
	Attributes    Group          `json:"attributes"`
	Relationships groupRelations `json:"relationships"`
***REMOVED***

type groupRelations struct ***REMOVED***
	Groups struct ***REMOVED***
		Data []groupRelation `json:"data"`
	***REMOVED*** `json:"groups"`
	Parent struct ***REMOVED***
		Data *groupRelation `json:"data"`
	***REMOVED*** `json:"parent"`
***REMOVED***

type groupRelation struct ***REMOVED***
	Type string `json:"type"`
	ID   string `json:"id"`
***REMOVED***

// UnmarshalJSON unmarshal group data properly (extract the ID)
func (g *groupData) UnmarshalJSON(data []byte) error ***REMOVED***
	var raw struct ***REMOVED***
		Type          string         `json:"type"`
		ID            string         `json:"id"`
		Attributes    Group          `json:"attributes"`
		Relationships groupRelations `json:"relationships"`
	***REMOVED***

	if err := json.Unmarshal(data, &raw); err != nil ***REMOVED***
		return err
	***REMOVED***

	g.ID = raw.ID
	g.Type = raw.Type
	g.Relationships = raw.Relationships
	g.Attributes = raw.Attributes

	if g.Attributes.ID == "" ***REMOVED***
		g.Attributes.ID = raw.ID
	***REMOVED***

	if g.Relationships.Parent.Data != nil ***REMOVED***
		g.Attributes.ParentID = g.Relationships.Parent.Data.ID
	***REMOVED***

	g.Attributes.GroupIDs = make([]string, 0, len(g.Relationships.Groups.Data))
	for _, rel := range g.Relationships.Groups.Data ***REMOVED***
		g.Attributes.GroupIDs = append(g.Attributes.GroupIDs, rel.ID)
	***REMOVED***

	return nil
***REMOVED***

func newGroupJSONAPI(g *Group) groupJSONAPI ***REMOVED***
	return groupJSONAPI***REMOVED***
		Data: newGroupData(g),
	***REMOVED***
***REMOVED***

func newGroupsJSONAPI(groups []*Group) groupsJSONAPI ***REMOVED***
	envelop := groupsJSONAPI***REMOVED***
		Data: make([]groupData, 0, len(groups)),
	***REMOVED***

	for _, g := range groups ***REMOVED***
		envelop.Data = append(envelop.Data, newGroupData(g))
	***REMOVED***

	return envelop
***REMOVED***

func newGroupData(group *Group) groupData ***REMOVED***
	data := groupData***REMOVED***
		Type:       "groups",
		ID:         group.ID,
		Attributes: *group,
		Relationships: groupRelations***REMOVED***
			Groups: struct ***REMOVED***
				Data []groupRelation `json:"data"`
			***REMOVED******REMOVED***
				Data: make([]groupRelation, 0, len(group.Groups)),
			***REMOVED***,
			Parent: struct ***REMOVED***
				Data *groupRelation `json:"data"`
			***REMOVED******REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***

	if group.Parent != nil ***REMOVED***
		data.Relationships.Parent.Data = &groupRelation***REMOVED***
			Type: "groups",
			ID:   group.Parent.ID,
		***REMOVED***
	***REMOVED***

	for _, gp := range group.Groups ***REMOVED***
		data.Relationships.Groups.Data = append(data.Relationships.Groups.Data, groupRelation***REMOVED***
			ID:   gp.ID,
			Type: "groups",
		***REMOVED***)
	***REMOVED***

	return data
***REMOVED***
