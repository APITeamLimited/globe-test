package stats

import (
	"bytes"
	"encoding/json"
	"strings"
)

// SystemTagSet is a bitmask that is used to keep track
// which system tags should be included with which metrics.
//go:generate enumer -type=SystemTagSet -transform=snake -output system_tag_set_gen.go
type SystemTagSet uint32

//nolint: golint
const (
	// Default system tags includes all of the system tags emitted with metrics by default.
	Proto SystemTagSet = 1 << iota
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
func (ts *SystemTagSet) Add(tag SystemTagSet) ***REMOVED***
	*ts |= tag
***REMOVED***

// Has checks a tag included in tag set.
func (ts *SystemTagSet) Has(tag SystemTagSet) bool ***REMOVED***
	return *ts&tag != 0
***REMOVED***

// ToSystemTagSet converts list of tags to SystemTagSet
func ToSystemTagSet(tags []string) *SystemTagSet ***REMOVED***
	ts := SystemTagSet(0)
	for _, tag := range tags ***REMOVED***
		if v, err := SystemTagSetString(tag); err == nil ***REMOVED***
			ts.Add(v)
		***REMOVED***
	***REMOVED***
	return &ts
***REMOVED***

// MarshalJSON converts the SystemTagSet to a list (JS array).
func (ts *SystemTagSet) MarshalJSON() ([]byte, error) ***REMOVED***
	var tags []string
	for _, tag := range SystemTagSetValues() ***REMOVED***
		if ts.Has(tag) ***REMOVED***
			tags = append(tags, tag.String())
		***REMOVED***
	***REMOVED***
	return json.Marshal(tags)
***REMOVED***

// UnmarshalJSON converts the tag list back to expected tag set.
func (ts *SystemTagSet) UnmarshalJSON(data []byte) error ***REMOVED***
	var tags []string
	if err := json.Unmarshal(data, &tags); err != nil ***REMOVED***
		return err
	***REMOVED***
	if len(tags) != 0 ***REMOVED***
		*ts = *ToSystemTagSet(tags)
	***REMOVED***
	return nil
***REMOVED***

// UnmarshalText converts the tag list to SystemTagSet.
func (ts *SystemTagSet) UnmarshalText(data []byte) error ***REMOVED***
	var list = bytes.Split(data, []byte(","))

	for _, key := range list ***REMOVED***
		key := strings.TrimSpace(string(key))
		if key == "" ***REMOVED***
			continue
		***REMOVED***
		if v, err := SystemTagSetString(key); err == nil ***REMOVED***
			ts.Add(v)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
