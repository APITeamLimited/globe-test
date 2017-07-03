/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package local

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	null "gopkg.in/guregu/null.v3"
)

func TestExecutorRun(t *testing.T) ***REMOVED***
	e := New(nil)
	assert.NoError(t, e.SetVUsMax(10))
	assert.NoError(t, e.SetVUs(10))

	ctx, cancel := context.WithCancel(context.Background())
	err := make(chan error, 1)
	go func() ***REMOVED*** err <- e.Run(ctx, nil) ***REMOVED***()
	cancel()
	assert.NoError(t, <-err)
***REMOVED***

func TestExecutorSetLogger(t *testing.T) ***REMOVED***
	logger, _ := logtest.NewNullLogger()
	e := New(nil)
	e.SetLogger(logger)
	assert.Equal(t, logger, e.GetLogger())
***REMOVED***

func TestExecutorEndTime(t *testing.T) ***REMOVED***
	e := New(nil)
	assert.NoError(t, e.SetVUsMax(10))
	assert.NoError(t, e.SetVUs(10))
	e.SetEndTime(lib.NullDurationFrom(1 * time.Second))
	assert.Equal(t, lib.NullDurationFrom(1*time.Second), e.GetEndTime())

	startTime := time.Now()
	assert.NoError(t, e.Run(context.Background(), nil))
	assert.True(t, time.Now().After(startTime.Add(1*time.Second)), "test did not take 1s")
***REMOVED***

func TestExecutorEndIterations(t *testing.T) ***REMOVED***
	var i int64
	e := New(lib.RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
		atomic.AddInt64(&i, 1)
		return nil, nil
	***REMOVED***))
	assert.NoError(t, e.SetVUsMax(10))
	assert.NoError(t, e.SetVUs(10))
	e.SetEndIterations(null.IntFrom(100))
	assert.Equal(t, null.IntFrom(100), e.GetEndIterations())
	assert.NoError(t, e.Run(context.Background(), nil))
	assert.Equal(t, int64(100), e.GetIterations())
	assert.Equal(t, int64(100), i)
***REMOVED***

func TestExecutorIsRunning(t *testing.T) ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	e := New(nil)

	err := make(chan error)
	go func() ***REMOVED*** err <- e.Run(ctx, nil) ***REMOVED***()
	for !e.IsRunning() ***REMOVED***
	***REMOVED***
	cancel()
	for e.IsRunning() ***REMOVED***
	***REMOVED***
	assert.NoError(t, <-err)
***REMOVED***

func TestExecutorSetVUsMax(t *testing.T) ***REMOVED***
	t.Run("Negative", func(t *testing.T) ***REMOVED***
		assert.EqualError(t, New(nil).SetVUsMax(-1), "vu cap can't be negative")
	***REMOVED***)

	t.Run("Raise", func(t *testing.T) ***REMOVED***
		e := New(nil)

		assert.NoError(t, e.SetVUsMax(50))
		assert.Equal(t, int64(50), e.GetVUsMax())

		assert.NoError(t, e.SetVUsMax(100))
		assert.Equal(t, int64(100), e.GetVUsMax())

		t.Run("Lower", func(t *testing.T) ***REMOVED***
			assert.NoError(t, e.SetVUsMax(50))
			assert.Equal(t, int64(50), e.GetVUsMax())
		***REMOVED***)
	***REMOVED***)

	t.Run("TooLow", func(t *testing.T) ***REMOVED***
		e := New(nil)
		e.ctx = context.Background()

		assert.NoError(t, e.SetVUsMax(100))
		assert.Equal(t, int64(100), e.GetVUsMax())

		assert.NoError(t, e.SetVUs(100))
		assert.Equal(t, int64(100), e.GetVUs())

		assert.EqualError(t, e.SetVUsMax(50), "can't lower vu cap (to 50) below vu count (100)")
	***REMOVED***)
***REMOVED***

func TestExecutorSetVUs(t *testing.T) ***REMOVED***
	t.Run("Negative", func(t *testing.T) ***REMOVED***
		assert.EqualError(t, New(nil).SetVUs(-1), "vu count can't be negative")
	***REMOVED***)

	t.Run("Too High", func(t *testing.T) ***REMOVED***
		assert.EqualError(t, New(nil).SetVUs(100), "can't raise vu count (to 100) above vu cap (0)")
	***REMOVED***)

	t.Run("Raise", func(t *testing.T) ***REMOVED***
		e := New(lib.RunnerFunc(func(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
			return nil, nil
		***REMOVED***))
		e.ctx = context.Background()

		assert.NoError(t, e.SetVUsMax(100))
		assert.Equal(t, int64(100), e.GetVUsMax())
		if assert.Len(t, e.vus, 100) ***REMOVED***
			num := 0
			for i, handle := range e.vus ***REMOVED***
				num++
				assert.NotNil(t, handle.vu, "vu %d lacks impl", i)
				assert.Nil(t, handle.ctx, "vu %d has ctx", i)
				assert.Nil(t, handle.cancel, "vu %d has cancel", i)
			***REMOVED***
			assert.Equal(t, 100, num)
		***REMOVED***

		assert.NoError(t, e.SetVUs(50))
		assert.Equal(t, int64(50), e.GetVUs())
		if assert.Len(t, e.vus, 100) ***REMOVED***
			num := 0
			for i, handle := range e.vus ***REMOVED***
				if i < 50 ***REMOVED***
					assert.NotNil(t, handle.cancel, "vu %d lacks cancel", i)
					num++
				***REMOVED*** else ***REMOVED***
					assert.Nil(t, handle.cancel, "vu %d has cancel", i)
				***REMOVED***
			***REMOVED***
			assert.Equal(t, 50, num)
		***REMOVED***

		assert.NoError(t, e.SetVUs(100))
		assert.Equal(t, int64(100), e.GetVUs())
		if assert.Len(t, e.vus, 100) ***REMOVED***
			num := 0
			for i, handle := range e.vus ***REMOVED***
				assert.NotNil(t, handle.cancel, "vu %d lacks cancel", i)
				num++
			***REMOVED***
			assert.Equal(t, 100, num)
		***REMOVED***

		t.Run("Lower", func(t *testing.T) ***REMOVED***
			assert.NoError(t, e.SetVUs(50))
			assert.Equal(t, int64(50), e.GetVUs())
			if assert.Len(t, e.vus, 100) ***REMOVED***
				num := 0
				for i, handle := range e.vus ***REMOVED***
					if i < 50 ***REMOVED***
						assert.NotNil(t, handle.cancel, "vu %d lacks cancel", i)
						num++
					***REMOVED*** else ***REMOVED***
						assert.Nil(t, handle.cancel, "vu %d has cancel", i)
					***REMOVED***
				***REMOVED***
				assert.Equal(t, 50, num)
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
***REMOVED***
