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
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
)

const groupSeparator = "::"

var ErrNameContainsGroupSeparator = errors.Errorf("group and check names may not contain '%s'", groupSeparator)

type SourceData struct ***REMOVED***
	Data     []byte
	Filename string
***REMOVED***

type Stage struct ***REMOVED***
	Duration NullDuration `json:"duration"`
	Target   null.Int     `json:"target"`
***REMOVED***

func (s *Stage) UnmarshalText(b []byte) error ***REMOVED***
	var stage Stage
	parts := strings.SplitN(string(b), ":", 2)
	if len(parts) > 0 && parts[0] != "" ***REMOVED***
		d, err := time.ParseDuration(parts[0])
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		stage.Duration = NullDurationFrom(d)
	***REMOVED***
	if len(parts) > 1 && parts[1] != "" ***REMOVED***
		t, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		stage.Target = null.IntFrom(t)
	***REMOVED***
	*s = stage
	return nil
***REMOVED***

type Group struct ***REMOVED***
	ID     string            `json:"id"`
	Path   string            `json:"path"`
	Name   string            `json:"name"`
	Parent *Group            `json:"-"`
	Groups map[string]*Group `json:"groups"`
	Checks map[string]*Check `json:"checks"`

	groupMutex sync.Mutex `json:"-"`
	checkMutex sync.Mutex `json:"-"`
***REMOVED***

func NewGroup(name string, parent *Group) (*Group, error) ***REMOVED***
	if strings.Contains(name, groupSeparator) ***REMOVED***
		return nil, ErrNameContainsGroupSeparator
	***REMOVED***

	path := name
	if parent != nil ***REMOVED***
		path = parent.Path + groupSeparator + path
	***REMOVED***

	hash := md5.Sum([]byte(path))
	id := hex.EncodeToString(hash[:])

	return &Group***REMOVED***
		ID:     id,
		Path:   path,
		Name:   name,
		Parent: parent,
		Groups: make(map[string]*Group),
		Checks: make(map[string]*Check),
	***REMOVED***, nil
***REMOVED***

func (g *Group) Group(name string) (*Group, error) ***REMOVED***
	snapshot := g.Groups
	group, ok := snapshot[name]
	if !ok ***REMOVED***
		g.groupMutex.Lock()
		defer g.groupMutex.Unlock()

		group, ok := g.Groups[name]
		if !ok ***REMOVED***
			group, err := NewGroup(name, g)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			g.Groups[name] = group
			return group, nil
		***REMOVED***
		return group, nil
	***REMOVED***
	return group, nil
***REMOVED***

func (g *Group) Check(name string) (*Check, error) ***REMOVED***
	snapshot := g.Checks
	check, ok := snapshot[name]
	if !ok ***REMOVED***
		g.checkMutex.Lock()
		defer g.checkMutex.Unlock()
		check, ok := g.Checks[name]
		if !ok ***REMOVED***
			check, err := NewCheck(name, g)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			g.Checks[name] = check
			return check, nil
		***REMOVED***
		return check, nil
	***REMOVED***
	return check, nil
***REMOVED***

type Check struct ***REMOVED***
	ID    string `json:"id"`
	Path  string `json:"path"`
	Group *Group `json:"-"`
	Name  string `json:"name"`

	Passes int64 `json:"passes"`
	Fails  int64 `json:"fails"`
***REMOVED***

func NewCheck(name string, group *Group) (*Check, error) ***REMOVED***
	if strings.Contains(name, groupSeparator) ***REMOVED***
		return nil, ErrNameContainsGroupSeparator
	***REMOVED***

	path := group.Path + groupSeparator + name
	hash := md5.Sum([]byte(path))
	id := hex.EncodeToString(hash[:])

	return &Check***REMOVED***
		ID:    id,
		Path:  path,
		Group: group,
		Name:  name,
	***REMOVED***, nil
***REMOVED***
