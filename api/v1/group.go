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

package v1

import (
	"fmt"

	"github.com/manyminds/api2go/jsonapi"

	"go.k6.io/k6/lib"
)

type Check struct ***REMOVED***
	ID     string `json:"id" yaml:"id"`
	Path   string `json:"path" yaml:"path"`
	Name   string `json:"name" yaml:"name"`
	Passes int64  `json:"passes" yaml:"passes"`
	Fails  int64  `json:"fails" yaml:"fails"`
***REMOVED***

func NewCheck(c *lib.Check) Check ***REMOVED***
	return Check***REMOVED***
		ID:     c.ID,
		Path:   c.Path,
		Name:   c.Name,
		Passes: c.Passes,
		Fails:  c.Fails,
	***REMOVED***
***REMOVED***

type Group struct ***REMOVED***
	ID     string  `json:"-" yaml:"id"`
	Path   string  `json:"path" yaml:"path"`
	Name   string  `json:"name" yaml:"name"`
	Checks []Check `json:"checks" yaml:"checks"`

	Parent   *Group   `json:"-" yaml:"-"`
	ParentID string   `json:"-" yaml:"parent-id"`
	Groups   []*Group `json:"-" yaml:"-"`
	GroupIDs []string `json:"-" yaml:"group-ids"`
***REMOVED***

func NewGroup(g *lib.Group, parent *Group) *Group ***REMOVED***
	group := &Group***REMOVED***
		ID:   g.ID,
		Path: g.Path,
		Name: g.Name,
	***REMOVED***

	if parent != nil ***REMOVED***
		group.Parent = parent
		group.ParentID = parent.ID
	***REMOVED*** else if g.Parent != nil ***REMOVED***
		group.Parent = NewGroup(g.Parent, nil)
		group.ParentID = g.Parent.ID
	***REMOVED***

	for _, gp := range g.Groups ***REMOVED***
		group.Groups = append(group.Groups, NewGroup(gp, group))
		group.GroupIDs = append(group.GroupIDs, gp.ID)
	***REMOVED***
	for _, c := range g.Checks ***REMOVED***
		group.Checks = append(group.Checks, NewCheck(c))
	***REMOVED***

	return group
***REMOVED***

func (g Group) GetID() string ***REMOVED***
	return g.ID
***REMOVED***

func (g *Group) SetID(v string) error ***REMOVED***
	g.ID = v
	return nil
***REMOVED***

func (g Group) GetReferences() []jsonapi.Reference ***REMOVED***
	return []jsonapi.Reference***REMOVED***
		***REMOVED***
			Type:         "groups",
			Name:         "parent",
			Relationship: jsonapi.ToOneRelationship,
		***REMOVED***,
		***REMOVED***
			Type:         "groups",
			Name:         "groups",
			Relationship: jsonapi.ToManyRelationship,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (g Group) GetReferencedIDs() []jsonapi.ReferenceID ***REMOVED***
	refs := []jsonapi.ReferenceID***REMOVED******REMOVED***
	if g.Parent != nil ***REMOVED***
		refs = append(refs, jsonapi.ReferenceID***REMOVED***
			ID:           g.Parent.GetID(),
			Type:         "groups",
			Name:         "parent",
			Relationship: jsonapi.ToOneRelationship,
		***REMOVED***)
	***REMOVED***
	for _, gp := range g.Groups ***REMOVED***
		refs = append(refs, jsonapi.ReferenceID***REMOVED***
			ID:           gp.GetID(),
			Type:         "groups",
			Name:         "groups",
			Relationship: jsonapi.ToManyRelationship,
		***REMOVED***)
	***REMOVED***
	return refs
***REMOVED***

func (g *Group) SetToManyReferenceIDs(name string, ids []string) error ***REMOVED***
	switch name ***REMOVED***
	case "groups":
		g.Groups = nil
		g.GroupIDs = ids
		return nil
	default:
		return fmt.Errorf("unknown to many relation: %s", name)
	***REMOVED***
***REMOVED***

func (g *Group) SetToOneReferenceID(name, id string) error ***REMOVED***
	switch name ***REMOVED***
	case "parent":
		g.Parent = nil
		g.ParentID = id
		return nil
	default:
		return fmt.Errorf("unknown to one relation: %s", name)
	***REMOVED***
***REMOVED***

func FlattenGroup(g *Group) []*Group ***REMOVED***
	groups := []*Group***REMOVED***g***REMOVED***
	for _, gp := range g.Groups ***REMOVED***
		groups = append(groups, FlattenGroup(gp)...)
	***REMOVED***
	return groups
***REMOVED***
