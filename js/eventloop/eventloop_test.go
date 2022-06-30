package eventloop

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/js/modulestest"
)

func TestBasicEventLoop(t *testing.T) ***REMOVED***
	t.Parallel()
	loop := New(&modulestest.VU***REMOVED***RuntimeField: goja.New()***REMOVED***)
	var ran int
	f := func() error ***REMOVED*** //nolint:unparam
		ran++
		return nil
	***REMOVED***
	require.NoError(t, loop.Start(f))
	require.Equal(t, 1, ran)
	require.NoError(t, loop.Start(f))
	require.Equal(t, 2, ran)
	require.Error(t, loop.Start(func() error ***REMOVED***
		_ = f()
		loop.RegisterCallback()(f)
		return errors.New("something")
	***REMOVED***))
	require.Equal(t, 3, ran)
***REMOVED***

func TestEventLoopRegistered(t *testing.T) ***REMOVED***
	t.Parallel()
	loop := New(&modulestest.VU***REMOVED***RuntimeField: goja.New()***REMOVED***)
	var ran int
	f := func() error ***REMOVED***
		ran++
		r := loop.RegisterCallback()
		go func() ***REMOVED***
			time.Sleep(time.Second)
			r(func() error ***REMOVED***
				ran++
				return nil
			***REMOVED***)
		***REMOVED***()
		return nil
	***REMOVED***
	start := time.Now()
	require.NoError(t, loop.Start(f))
	took := time.Since(start)
	require.Equal(t, 2, ran)
	require.Less(t, time.Second, took)
	require.Greater(t, time.Second+time.Millisecond*100, took)
***REMOVED***

func TestEventLoopWaitOnRegistered(t *testing.T) ***REMOVED***
	t.Parallel()
	var ran int
	loop := New(&modulestest.VU***REMOVED***RuntimeField: goja.New()***REMOVED***)
	f := func() error ***REMOVED***
		ran++
		r := loop.RegisterCallback()
		go func() ***REMOVED***
			time.Sleep(time.Second)
			r(func() error ***REMOVED***
				ran++
				return nil
			***REMOVED***)
		***REMOVED***()
		return fmt.Errorf("expected")
	***REMOVED***
	start := time.Now()
	require.Error(t, loop.Start(f))
	took := time.Since(start)
	loop.WaitOnRegistered()
	took2 := time.Since(start)
	require.Equal(t, 1, ran)
	require.Greater(t, time.Millisecond*50, took)
	require.Less(t, time.Second, took2)
	require.Greater(t, time.Second+time.Millisecond*100, took2)
***REMOVED***

func TestEventLoopReuse(t *testing.T) ***REMOVED***
	t.Parallel()
	sleepTime := time.Millisecond * 500
	loop := New(&modulestest.VU***REMOVED***RuntimeField: goja.New()***REMOVED***)
	f := func() error ***REMOVED***
		for i := 0; i < 100; i++ ***REMOVED***
			bad := i == 17
			r := loop.RegisterCallback()

			go func() ***REMOVED***
				if !bad ***REMOVED***
					time.Sleep(sleepTime)
				***REMOVED***
				r(func() error ***REMOVED***
					if bad ***REMOVED***
						return errors.New("something")
					***REMOVED***
					panic("this should never execute")
				***REMOVED***)
			***REMOVED***()
		***REMOVED***
		return fmt.Errorf("expected")
	***REMOVED***
	for i := 0; i < 3; i++ ***REMOVED***
		start := time.Now()
		require.Error(t, loop.Start(f))
		took := time.Since(start)
		loop.WaitOnRegistered()
		took2 := time.Since(start)
		require.Greater(t, time.Millisecond*50, took)
		require.Less(t, sleepTime, took2)
		require.Greater(t, sleepTime+time.Millisecond*100, took2)
	***REMOVED***
***REMOVED***

func TestEventLoopPanicOnDoubleCallback(t *testing.T) ***REMOVED***
	t.Parallel()
	loop := New(&modulestest.VU***REMOVED***RuntimeField: goja.New()***REMOVED***)
	var ran int
	f := func() error ***REMOVED***
		ran++
		r := loop.RegisterCallback()
		go func() ***REMOVED***
			time.Sleep(time.Second)
			r(func() error ***REMOVED***
				ran++
				return nil
			***REMOVED***)

			require.Panics(t, func() ***REMOVED*** r(func() error ***REMOVED*** return nil ***REMOVED***) ***REMOVED***)
		***REMOVED***()
		return nil
	***REMOVED***
	start := time.Now()
	require.NoError(t, loop.Start(f))
	took := time.Since(start)
	require.Equal(t, 2, ran)
	require.Less(t, time.Second, took)
	require.Greater(t, time.Second+time.Millisecond*100, took)
***REMOVED***
