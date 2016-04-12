package runner

import (
	"testing"
)

func TestScaleNoChange(t *testing.T) ***REMOVED***
	i := 10
	start := func() ***REMOVED*** i++ ***REMOVED***
	stop := func() ***REMOVED*** i-- ***REMOVED***
	scale(i, 10, start, stop)
	if i != 10 ***REMOVED***
		t.Fail()
	***REMOVED***
***REMOVED***

func TestScaleAdd(t *testing.T) ***REMOVED***
	i := 10
	start := func() ***REMOVED*** i++ ***REMOVED***
	stop := func() ***REMOVED*** i-- ***REMOVED***
	scale(i, 15, start, stop)
	if i != 15 ***REMOVED***
		t.Fail()
	***REMOVED***
***REMOVED***

func TestScaleRemove(t *testing.T) ***REMOVED***
	i := 10
	start := func() ***REMOVED*** i++ ***REMOVED***
	stop := func() ***REMOVED*** i-- ***REMOVED***
	scale(i, 5, start, stop)
	if i != 5 ***REMOVED***
		t.Fail()
	***REMOVED***
***REMOVED***
