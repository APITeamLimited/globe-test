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
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/loadimpact/k6/lib/types"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
)

// Separator for group IDs.
const GroupSeparator = "::"

// Error emitted if you attempt to instantiate a Group or Check that contains the separator.
var ErrNameContainsGroupSeparator = errors.New("group and check names may not contain '::'")

// Wraps a source file; data and filename.
type SourceData struct ***REMOVED***
	Data     []byte
	Filename string
***REMOVED***

// StageFields defines the fields used for a Stage; this is a dumb hack to make the JSON code
// cleaner. pls fix.
type StageFields struct ***REMOVED***
	// Duration of the stage.
	Duration types.NullDuration `json:"duration"`

	// If Valid, the VU count will be linearly interpolated towards this value.
	Target null.Int `json:"target"`
***REMOVED***

// A Stage defines a step in a test's timeline.
type Stage StageFields

// For some reason, implementing UnmarshalText makes encoding/json treat the type as a string.
func (s *Stage) UnmarshalJSON(b []byte) error ***REMOVED***
	var fields StageFields
	if err := json.Unmarshal(b, &fields); err != nil ***REMOVED***
		return err
	***REMOVED***
	*s = Stage(fields)
	return nil
***REMOVED***

func (s Stage) MarshalJSON() ([]byte, error) ***REMOVED***
	return json.Marshal(StageFields(s))
***REMOVED***

func (s *Stage) UnmarshalText(b []byte) error ***REMOVED***
	var stage Stage
	parts := strings.SplitN(string(b), ":", 2)
	if len(parts) > 0 && parts[0] != "" ***REMOVED***
		d, err := time.ParseDuration(parts[0])
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		stage.Duration = types.NullDurationFrom(d)
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

// A Group is an organisational block, that samples and checks may be tagged with.
//
// For more information, refer to the js/modules/k6.K6.Group() function.
type Group struct ***REMOVED***
	// Arbitrary name of the group.
	Name string `json:"name"`

	// A group may belong to another group, which may belong to another group, etc. The Path
	// describes the hierarchy leading down to this group, with the segments delimited by '::'.
	// As an example: a group "Inner" inside a group named "Outer" would have a path of
	// "::Outer::Inner". The empty first item is the root group, which is always named "".
	Parent *Group `json:"-"`
	Path   string `json:"path"`

	// A group's ID is a hash of the Path. It is deterministic between different k6
	// instances of the same version, but should be treated as opaque - the hash function
	// or length may change.
	ID string `json:"id"`

	// Groups and checks that are children of this group.
	Groups map[string]*Group `json:"groups"`
	Checks map[string]*Check `json:"checks"`

	groupMutex sync.Mutex
	checkMutex sync.Mutex
***REMOVED***

// Creates a new group with the given name and parent group.
//
// The root group must be created with the name "" and parent set to nil; this is the only case
// where a nil parent or empty name is allowed.
func NewGroup(name string, parent *Group) (*Group, error) ***REMOVED***
	if strings.Contains(name, GroupSeparator) ***REMOVED***
		return nil, ErrNameContainsGroupSeparator
	***REMOVED***

	path := name
	if parent != nil ***REMOVED***
		path = parent.Path + GroupSeparator + path
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

// Create a child group belonging to this group.
// This is safe to call from multiple goroutines simultaneously.
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

// Create a check belonging to this group.
// This is safe to call from multiple goroutines simultaneously.
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

// A Check stores a series of successful or failing tests against a value.
//
// For more information, refer to the js/modules/k6.K6.Check() function.
type Check struct ***REMOVED***
	// Arbitrary name of the check.
	Name string `json:"name"`

	// A Check belongs to a Group, which may belong to other groups. The Path describes
	// the hierarchy of these groups, with the segments delimited by '::'.
	// As an example: a check "My Check" within a group "Inner" within a group "Outer"
	// would have a Path of "::Outer::Inner::My Check". The empty first item is the root group,
	// which is always named "".
	Group *Group `json:"-"`
	Path  string `json:"path"`

	// A check's ID is a hash of the Path. It is deterministic between different k6
	// instances of the same version, but should be treated as opaque - the hash function
	// or length may change.
	ID string `json:"id"`

	// Counters for how many times this check has passed and failed respectively.
	Passes int64 `json:"passes"`
	Fails  int64 `json:"fails"`
***REMOVED***

// Creates a new check with the given name and parent group. The group may not be nil.
func NewCheck(name string, group *Group) (*Check, error) ***REMOVED***
	if strings.Contains(name, GroupSeparator) ***REMOVED***
		return nil, ErrNameContainsGroupSeparator
	***REMOVED***

	path := group.Path + GroupSeparator + name
	hash := md5.Sum([]byte(path))
	id := hex.EncodeToString(hash[:])

	return &Check***REMOVED***
		ID:    id,
		Path:  path,
		Group: group,
		Name:  name,
	***REMOVED***, nil
***REMOVED***
