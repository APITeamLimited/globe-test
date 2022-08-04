package har

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsAllowedURL(t *testing.T) ***REMOVED***
	var allowed = []struct ***REMOVED***
		url      string
		only     []string
		skip     []string
		expected bool
	***REMOVED******REMOVED***
		***REMOVED***"http://www.google.com/", []string***REMOVED******REMOVED***, []string***REMOVED******REMOVED***, true***REMOVED***,
		***REMOVED***"http://www.google.com/", []string***REMOVED***"google.com"***REMOVED***, []string***REMOVED******REMOVED***, true***REMOVED***,
		***REMOVED***"https://www.google.com/", []string***REMOVED***"google.com"***REMOVED***, []string***REMOVED******REMOVED***, true***REMOVED***,
		***REMOVED***"https://www.google.com/", []string***REMOVED***"http://"***REMOVED***, []string***REMOVED******REMOVED***, false***REMOVED***,
		***REMOVED***"http://www.google.com/?hl=en", []string***REMOVED***"http://www.google.com"***REMOVED***, []string***REMOVED******REMOVED***, true***REMOVED***,
		***REMOVED***"http://www.google.com/?hl=en", []string***REMOVED***"google.com", "google.co.uk"***REMOVED***, []string***REMOVED******REMOVED***, true***REMOVED***,
		***REMOVED***"http://www.google.com/?hl=en", []string***REMOVED******REMOVED***, []string***REMOVED***"google.com"***REMOVED***, false***REMOVED***,
		***REMOVED***"http://www.google.com/?hl=en", []string***REMOVED******REMOVED***, []string***REMOVED***"google.co.uk"***REMOVED***, true***REMOVED***,
	***REMOVED***

	for _, s := range allowed ***REMOVED***
		v := IsAllowedURL(s.url, s.only, s.skip)
		assert.Equal(t, v, s.expected, fmt.Sprintf("params: %v, %v, %v", s.url, s.only, s.skip))
	***REMOVED***
***REMOVED***

func TestSplitEntriesInBatches(t *testing.T) ***REMOVED***
	t1 := time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)

	entries := []*Entry***REMOVED******REMOVED***

	// 10 time entries with increments of 100ms or 200ms
	for i := 1; i <= 10; i++ ***REMOVED***

		period := 100
		if i%2 == 0 ***REMOVED***
			period = 200
		***REMOVED***
		t1 = t1.Add(time.Duration(period) * time.Millisecond)
		entries = append(entries, &Entry***REMOVED***StartedDateTime: t1***REMOVED***)
	***REMOVED***

	splitValues := []struct ***REMOVED***
		diff, groups uint
	***REMOVED******REMOVED***
		***REMOVED***0, 1***REMOVED***,
		***REMOVED***100, 10***REMOVED***,
		***REMOVED***150, 6***REMOVED***,
		***REMOVED***200, 6***REMOVED***,
		***REMOVED***201, 1***REMOVED***,
		***REMOVED***500, 1***REMOVED***,
	***REMOVED***

	sort.Sort(EntryByStarted(entries))

	for _, v := range splitValues ***REMOVED***
		result := SplitEntriesInBatches(entries, v.diff)
		assert.Equal(t, len(result), int(v.groups), fmt.Sprintf("params: entries, %v", v.diff))
	***REMOVED***
***REMOVED***
