package tagset

import (
	"bytes"
	"encoding/json"
	"strings"
)

// TagSet is a bitmask that is used to keep track
// which system tags should be included with which metrics.
//go:generate enumer -type=TagSet -transform=snake -output tagset_gen.go
type TagSet uint32

//nolint: golint
const (
	// Default system tags includes all of the system tags emitted with metrics by default.
	Proto TagSet = 1 << iota
	SubProto
	Status
	Method
	URL
	Name
	Group
	Check
	Error
	ErrorCode
	TLSVersion

	// System tags not enabled by default.
	Iter
	VU
	OCSPStatus
	IP
)

// Add adds a tag to tag set.
func (ts *TagSet) Add(tag TagSet) ***REMOVED***
	*ts |= tag
***REMOVED***

// Has checks a tag included in tag set.
func (ts *TagSet) Has(tag TagSet) bool ***REMOVED***
	return *ts&tag != 0
***REMOVED***

// FromList converts list of tags to TagSet
func FromList(tags []string) *TagSet ***REMOVED***
	ts := TagSet(0)
	for _, tag := range tags ***REMOVED***
		if v, err := TagSetString(tag); err == nil ***REMOVED***
			ts.Add(v)
		***REMOVED***
	***REMOVED***
	return &ts
***REMOVED***

// MarshalJSON converts the TagSet to a list (JS array).
func (ts *TagSet) MarshalJSON() ([]byte, error) ***REMOVED***
	var tags []string
	for _, tag := range TagSetValues() ***REMOVED***
		if ts.Has(tag) ***REMOVED***
			tags = append(tags, tag.String())
		***REMOVED***
	***REMOVED***
	return json.Marshal(tags)
***REMOVED***

// UnmarshalJSON converts the tag list back to expected tag set.
func (ts *TagSet) UnmarshalJSON(data []byte) error ***REMOVED***
	var tags []string
	if err := json.Unmarshal(data, &tags); err != nil ***REMOVED***
		return err
	***REMOVED***
	if len(tags) != 0 ***REMOVED***
		*ts = *FromList(tags)
	***REMOVED***
	return nil
***REMOVED***

// UnmarshalText converts the tag list to TagSet.
func (ts *TagSet) UnmarshalText(data []byte) error ***REMOVED***
	var list = bytes.Split(data, []byte(","))

	for _, key := range list ***REMOVED***
		key := strings.TrimSpace(string(key))
		if key == "" ***REMOVED***
			continue
		***REMOVED***
		if v, err := TagSetString(key); err == nil ***REMOVED***
			ts.Add(v)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
