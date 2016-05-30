package speedboat

import (
	"testing"
	"time"
)

func TestVUsAtSingleStage(t *testing.T) ***REMOVED***
	test := Test***REMOVED***
		Stages: []TestStage***REMOVED***
			TestStage***REMOVED***Duration: 10 * time.Second, StartVUs: 0, EndVUs: 10***REMOVED***,
		***REMOVED***,
	***REMOVED***
	if n := test.VUsAt(0 * time.Second); n != 0 ***REMOVED***
		t.Errorf("Wrong number at 0s: %d", n)
	***REMOVED***
	if n := test.VUsAt(5 * time.Second); n != 5 ***REMOVED***
		t.Errorf("Wrong number at 5s: %d", n)
	***REMOVED***
	if n := test.VUsAt(10 * time.Second); n != 10 ***REMOVED***
		t.Errorf("Wrong number at 10s: %d", n)
	***REMOVED***
***REMOVED***

func TestVUsAtMultiStage(t *testing.T) ***REMOVED***
	test := Test***REMOVED***
		Stages: []TestStage***REMOVED***
			TestStage***REMOVED***Duration: 5 * time.Second, StartVUs: 0, EndVUs: 10***REMOVED***,
			TestStage***REMOVED***Duration: 10 * time.Second, StartVUs: 10, EndVUs: 20***REMOVED***,
		***REMOVED***,
	***REMOVED***
	if n := test.VUsAt(5 * time.Second); n != 10 ***REMOVED***
		t.Errorf("Wrong number at 5s: %d", n)
	***REMOVED***
	if n := test.VUsAt(10 * time.Second); n != 15 ***REMOVED***
		t.Errorf("Wrong number at 10s: %d", n)
	***REMOVED***
	if n := test.VUsAt(15 * time.Second); n != 20 ***REMOVED***
		t.Errorf("Wrong number at 15s: %d", n)
	***REMOVED***
***REMOVED***
