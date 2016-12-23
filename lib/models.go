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

package lib

import (
	"gopkg.in/guregu/null.v3"
	"sync"
	"sync/atomic"
	"time"
)

type Stage struct ***REMOVED***
	Duration time.Duration `json:"duration"`
	Target   null.Int      `json:"target"`
***REMOVED***

type Group struct ***REMOVED***
	ID int64

	Name   string
	Parent *Group
	Groups map[string]*Group
	Checks map[string]*Check

	groupMutex sync.Mutex
	checkMutex sync.Mutex
***REMOVED***

func NewGroup(name string, parent *Group, idCounter *int64) *Group ***REMOVED***
	var id int64
	if idCounter != nil ***REMOVED***
		id = atomic.AddInt64(idCounter, 1)
	***REMOVED***

	return &Group***REMOVED***
		ID:     id,
		Name:   name,
		Parent: parent,
		Groups: make(map[string]*Group),
		Checks: make(map[string]*Check),
	***REMOVED***
***REMOVED***

func (g *Group) Group(name string, idCounter *int64) (*Group, bool) ***REMOVED***
	snapshot := g.Groups
	group, ok := snapshot[name]
	if !ok ***REMOVED***
		g.groupMutex.Lock()
		group, ok = g.Groups[name]
		if !ok ***REMOVED***
			group = NewGroup(name, g, idCounter)
			g.Groups[name] = group
		***REMOVED***
		g.groupMutex.Unlock()
	***REMOVED***
	return group, ok
***REMOVED***

func (g *Group) Check(name string, idCounter *int64) (*Check, bool) ***REMOVED***
	snapshot := g.Checks
	check, ok := snapshot[name]
	if !ok ***REMOVED***
		g.checkMutex.Lock()
		check, ok = g.Checks[name]
		if !ok ***REMOVED***
			check = NewCheck(name, g, idCounter)
			g.Checks[name] = check
		***REMOVED***
		g.checkMutex.Unlock()
	***REMOVED***
	return check, ok
***REMOVED***

type Check struct ***REMOVED***
	ID int64

	Group *Group
	Name  string

	Passes int64
	Fails  int64
***REMOVED***

func NewCheck(name string, group *Group, idCounter *int64) *Check ***REMOVED***
	var id int64
	if idCounter != nil ***REMOVED***
		id = atomic.AddInt64(idCounter, 1)
	***REMOVED***
	return &Check***REMOVED***ID: id, Name: name, Group: group***REMOVED***
***REMOVED***
