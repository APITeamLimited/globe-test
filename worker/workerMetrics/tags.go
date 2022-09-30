package workerMetrics

import (
	"bytes"
	"encoding/json"
	"sort"
	"strings"
)

// EnabledTags is a string to bool map (for lookup efficiency) that is used to keep track
// of which system tags and non-system tags to include.
type EnabledTags map[string]bool

// UnmarshalText converts the tag list to EnabledTags.
func (i *EnabledTags) UnmarshalText(data []byte) error ***REMOVED***
	list := bytes.Split(data, []byte(","))
	if *i == nil ***REMOVED***
		*i = make(EnabledTags, len(list))
	***REMOVED***

	for _, key := range list ***REMOVED***
		key := strings.TrimSpace(string(key))
		if key == "" ***REMOVED***
			continue
		***REMOVED***
		(*i)[key] = true
	***REMOVED***

	return nil
***REMOVED***

// MarshalJSON converts the EnabledTags to a list (JS array).
func (i *EnabledTags) MarshalJSON() ([]byte, error) ***REMOVED***
	var tags []string
	if *i != nil ***REMOVED***
		tags = make([]string, 0, len(*i))
		for tag := range *i ***REMOVED***
			tags = append(tags, tag)
		***REMOVED***
		sort.Strings(tags)
	***REMOVED***

	return json.Marshal(tags)
***REMOVED***

// UnmarshalJSON converts the tag list back to expected tag set.
func (i *EnabledTags) UnmarshalJSON(data []byte) error ***REMOVED***
	var tags []string
	if err := json.Unmarshal(data, &tags); err != nil ***REMOVED***
		return err
	***REMOVED***
	*i = make(EnabledTags, len(tags))
	for _, tag := range tags ***REMOVED***
		(*i)[tag] = true
	***REMOVED***

	return nil
***REMOVED***
