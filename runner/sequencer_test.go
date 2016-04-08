package runner

import (
	"testing"
	"time"
)

func TestAddCount(t *testing.T) ***REMOVED***
	seq := NewSequencer()
	if seq.Count() != 0 ***REMOVED***
		t.Error("Why does a blank sequencer have metrics!?")
	***REMOVED***
	seq.Add(Metric***REMOVED***Duration: time.Duration(10) * time.Second***REMOVED***)
	if seq.Count() != 1 ***REMOVED***
		t.Error("Add() didn't seem to add anything")
	***REMOVED***
***REMOVED***

func TestStatDuration(t *testing.T) ***REMOVED***
	seq := NewSequencer()
	seq.Metrics = []Metric***REMOVED***
		Metric***REMOVED***Duration: time.Duration(10) * time.Second***REMOVED***,
		Metric***REMOVED***Duration: time.Duration(15) * time.Second***REMOVED***,
		Metric***REMOVED***Duration: time.Duration(20) * time.Second***REMOVED***,
		Metric***REMOVED***Duration: time.Duration(25) * time.Second***REMOVED***,
	***REMOVED***
	s := seq.StatDuration()
	if s.Avg != 17.5 ***REMOVED***
		t.Error("Wrong average", s.Avg)
	***REMOVED***
	if s.Med != 20 ***REMOVED***
		t.Error("Wrong median", s.Med)
	***REMOVED***
	if s.Min != 10 ***REMOVED***
		t.Error("Wrong min", s.Min)
	***REMOVED***
	if s.Max != 25 ***REMOVED***
		t.Error("Wrong max", s.Max)
	***REMOVED***
***REMOVED***

func TestStatDurationNoMetrics(t *testing.T) ***REMOVED***
	seq := NewSequencer()
	s := seq.StatDuration()
	if s.Avg != 0 || s.Med != 0 || s.Min != 0 || s.Max != 0 ***REMOVED***
		t.Error("Nonzero values", s)
	***REMOVED***
***REMOVED***
