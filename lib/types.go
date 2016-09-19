package lib

import (
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/guregu/null.v3"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Status struct ***REMOVED***
	Running null.Bool `json:"running"`
	VUs     null.Int  `json:"vus"`
	VUsMax  null.Int  `json:"vus-max"`
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

type Info struct ***REMOVED***
	Version string `json:"version"`
***REMOVED***

func (i Info) GetName() string ***REMOVED***
	return "info"
***REMOVED***

func (i Info) GetID() string ***REMOVED***
	return "default"
***REMOVED***

type Options struct ***REMOVED***
	VUs      int64         `json:"vus"`
	VUsMax   int64         `json:"vus-max"`
	Duration time.Duration `json:"duration"`
***REMOVED***

func (o Options) GetName() string ***REMOVED***
	return "options"
***REMOVED***

func (o Options) GetID() string ***REMOVED***
	return "default"
***REMOVED***

type Group struct ***REMOVED***
	ID int64 `json:"-"`

	Name   string            `json:"name"`
	Parent *Group            `json:"-"`
	Groups map[string]*Group `json:"-"`
	Tests  map[string]*Test  `json:"-"`

	groupMutex sync.Mutex `json:"-"`
	testMutex  sync.Mutex `json:"-"`
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
		Tests:  make(map[string]*Test),
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

func (g *Group) Test(name string, idCounter *int64) (*Test, bool) ***REMOVED***
	snapshot := g.Tests
	test, ok := snapshot[name]
	if !ok ***REMOVED***
		g.testMutex.Lock()
		test, ok = g.Tests[name]
		if !ok ***REMOVED***
			test = NewTest(name, g, idCounter)
			g.Tests[name] = test
		***REMOVED***
		g.testMutex.Unlock()
	***REMOVED***
	return test, ok
***REMOVED***

func (g Group) GetID() string ***REMOVED***
	return strconv.FormatInt(g.ID, 10)
***REMOVED***

func (g Group) GetReferences() []jsonapi.Reference ***REMOVED***
	return []jsonapi.Reference***REMOVED***
		jsonapi.Reference***REMOVED***
			Name:         "parent",
			Type:         "groups",
			Relationship: jsonapi.ToOneRelationship,
		***REMOVED***,
		jsonapi.Reference***REMOVED***
			Name:         "tests",
			Type:         "tests",
			Relationship: jsonapi.ToManyRelationship,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (g Group) GetReferencedIDs() []jsonapi.ReferenceID ***REMOVED***
	ids := make([]jsonapi.ReferenceID, 0, len(g.Tests)+len(g.Groups))
	for _, test := range g.Tests ***REMOVED***
		ids = append(ids, jsonapi.ReferenceID***REMOVED***
			ID:           test.GetID(),
			Type:         "tests",
			Name:         "tests",
			Relationship: jsonapi.ToManyRelationship,
		***REMOVED***)
	***REMOVED***
	for _, group := range g.Groups ***REMOVED***
		ids = append(ids, jsonapi.ReferenceID***REMOVED***
			ID:           group.GetID(),
			Type:         "groups",
			Name:         "groups",
			Relationship: jsonapi.ToManyRelationship,
		***REMOVED***)
	***REMOVED***
	if g.Parent != nil ***REMOVED***
		ids = append(ids, jsonapi.ReferenceID***REMOVED***
			ID:           g.Parent.GetID(),
			Type:         "groups",
			Name:         "parent",
			Relationship: jsonapi.ToOneRelationship,
		***REMOVED***)
	***REMOVED***
	return ids
***REMOVED***

func (g Group) GetReferencedStructs() []jsonapi.MarshalIdentifier ***REMOVED***
	// Note: we're not sideloading the parent, that snowballs into making requests for a single
	// group return *every single known group* thanks to the common root group.
	refs := make([]jsonapi.MarshalIdentifier, 0, len(g.Tests)+len(g.Groups))
	for _, test := range g.Tests ***REMOVED***
		refs = append(refs, test)
	***REMOVED***
	for _, group := range g.Groups ***REMOVED***
		refs = append(refs, group)
	***REMOVED***
	return refs
***REMOVED***

type Test struct ***REMOVED***
	ID int64 `json:"-"`

	Group *Group `json:"-"`
	Name  string `json:"name"`

	Passes int64 `json:"passes"`
	Fails  int64 `json:"fails"`
***REMOVED***

func NewTest(name string, group *Group, idCounter *int64) *Test ***REMOVED***
	var id int64
	if idCounter != nil ***REMOVED***
		id = atomic.AddInt64(idCounter, 1)
	***REMOVED***
	return &Test***REMOVED***ID: id, Name: name, Group: group***REMOVED***
***REMOVED***

func (t Test) GetID() string ***REMOVED***
	return strconv.FormatInt(t.ID, 10)
***REMOVED***

func (t Test) GetReferences() []jsonapi.Reference ***REMOVED***
	return []jsonapi.Reference***REMOVED***
		jsonapi.Reference***REMOVED***
			Name:         "group",
			Type:         "groups",
			Relationship: jsonapi.ToOneRelationship,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (t Test) GetReferencedIDs() []jsonapi.ReferenceID ***REMOVED***
	return []jsonapi.ReferenceID***REMOVED***
		jsonapi.ReferenceID***REMOVED***
			ID:   t.Group.GetID(),
			Type: "groups",
			Name: "group",
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (t Test) GetReferencedStructs() []jsonapi.MarshalIdentifier ***REMOVED***
	return []jsonapi.MarshalIdentifier***REMOVED***t.Group***REMOVED***
***REMOVED***
