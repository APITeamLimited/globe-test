package lib

import (
	"strings"
	"testing"
	"time"
)

func TestTimeoutError(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		stage, expectedStrContain string
		d                         time.Duration
	***REMOVED******REMOVED***
		***REMOVED***"setup", "1 seconds", time.Second***REMOVED***,
		***REMOVED***"teardown", "2 seconds", time.Second * 2***REMOVED***,
		***REMOVED***"", "0 seconds", time.Duration(0)***REMOVED***,
	***REMOVED***

	for _, tc := range tests ***REMOVED***
		te := NewTimeoutError(tc.stage, tc.d)
		if !strings.Contains(te.String(), tc.expectedStrContain) ***REMOVED***
			t.Errorf("Expected error contains %s, but got: %s", tc.expectedStrContain, te.String())
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestTimeoutErrorHint(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		stage string
		empty bool
	***REMOVED******REMOVED***
		***REMOVED***"setup", false***REMOVED***,
		***REMOVED***"teardown", false***REMOVED***,
		***REMOVED***"not handle", true***REMOVED***,
	***REMOVED***

	for _, tc := range tests ***REMOVED***
		te := NewTimeoutError(tc.stage, time.Second)
		if tc.empty && te.Hint() != "" ***REMOVED***
			t.Errorf("Expected empty hint, got: %s", te.Hint())
		***REMOVED***
		if !tc.empty && te.Hint() == "" ***REMOVED***
			t.Errorf("Expected non-empty hint, got empty")
		***REMOVED***
	***REMOVED***
***REMOVED***
