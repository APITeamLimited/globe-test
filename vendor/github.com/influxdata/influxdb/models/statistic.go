package models

// Statistic is the representation of a statistic used by the monitoring service.
type Statistic struct ***REMOVED***
	Name   string                 `json:"name"`
	Tags   map[string]string      `json:"tags"`
	Values map[string]interface***REMOVED******REMOVED*** `json:"values"`
***REMOVED***

// NewStatistic returns an initialized Statistic.
func NewStatistic(name string) Statistic ***REMOVED***
	return Statistic***REMOVED***
		Name:   name,
		Tags:   make(map[string]string),
		Values: make(map[string]interface***REMOVED******REMOVED***),
	***REMOVED***
***REMOVED***

// StatisticTags is a map that can be merged with others without causing
// mutations to either map.
type StatisticTags map[string]string

// Merge creates a new map containing the merged contents of tags and t.
// If both tags and the receiver map contain the same key, the value in tags
// is used in the resulting map.
//
// Merge always returns a usable map.
func (t StatisticTags) Merge(tags map[string]string) map[string]string ***REMOVED***
	// Add everything in tags to the result.
	out := make(map[string]string, len(tags))
	for k, v := range tags ***REMOVED***
		out[k] = v
	***REMOVED***

	// Only add values from t that don't appear in tags.
	for k, v := range t ***REMOVED***
		if _, ok := tags[k]; !ok ***REMOVED***
			out[k] = v
		***REMOVED***
	***REMOVED***
	return out
***REMOVED***
