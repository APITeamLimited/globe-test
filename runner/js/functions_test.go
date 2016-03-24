package js

import (
	"github.com/robertkrimen/otto"
	"testing"
	"time"
)

func makeTestSleepFunc() (func(time.Duration), <-chan time.Duration) ***REMOVED***
	ch := make(chan time.Duration)
	fn := func(d time.Duration) ***REMOVED***
		go func() ***REMOVED***
			ch <- d
		***REMOVED***()
	***REMOVED***

	return fn, ch
***REMOVED***

func TestJSSleep(t *testing.T) ***REMOVED***
	sleep, times := makeTestSleepFunc()

	vm := otto.New()
	vm.Set("sleep", jsSleepFactory(sleep))
	_, err := vm.Run(`sleep(1)`)
	if err != nil ***REMOVED***
		t.Error("JS Error", err)
	***REMOVED***

	d := <-times
	if d != time.Duration(1)*time.Second ***REMOVED***
		t.Error("Wrong amount of sleep", d)
	***REMOVED***
***REMOVED***

func TestJSSleepFraction(t *testing.T) ***REMOVED***
	sleep, times := makeTestSleepFunc()

	vm := otto.New()
	vm.Set("sleep", jsSleepFactory(sleep))
	_, err := vm.Run(`sleep(0.1)`)
	if err != nil ***REMOVED***
		t.Error("JS Error", err)
	***REMOVED***

	d := <-times
	if d != time.Duration(100)*time.Millisecond ***REMOVED***
		t.Error("Wrong amount of sleep", d)
	***REMOVED***
***REMOVED***
