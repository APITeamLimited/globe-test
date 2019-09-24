package lib

import (
	"testing"
	"time"
)

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
