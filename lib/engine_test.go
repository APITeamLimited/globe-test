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

package lib

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
	"runtime"
	"testing"
	"time"
)

func TestNewEngine(t *testing.T) ***REMOVED***
	_, err := NewEngine(nil, Options***REMOVED******REMOVED***)
	assert.NoError(t, err)
***REMOVED***

func TestNewEngineOptions(t *testing.T) ***REMOVED***
	t.Run("VUsMax", func(t *testing.T) ***REMOVED***
		t.Run("not set", func(t *testing.T) ***REMOVED***
			e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(0), e.GetVUsMax())
			assert.Equal(t, int64(0), e.GetVUs())
		***REMOVED***)
		t.Run("set", func(t *testing.T) ***REMOVED***
			e, err := NewEngine(nil, Options***REMOVED***
				VUsMax: null.IntFrom(10),
			***REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(10), e.GetVUsMax())
			assert.Equal(t, int64(0), e.GetVUs())
		***REMOVED***)
	***REMOVED***)
	t.Run("VUs", func(t *testing.T) ***REMOVED***
		t.Run("no max", func(t *testing.T) ***REMOVED***
			_, err := NewEngine(nil, Options***REMOVED***
				VUs: null.IntFrom(10),
			***REMOVED***)
			assert.EqualError(t, err, "more vus than allocated requested")
		***REMOVED***)
		t.Run("max too low", func(t *testing.T) ***REMOVED***
			_, err := NewEngine(nil, Options***REMOVED***
				VUsMax: null.IntFrom(1),
				VUs:    null.IntFrom(10),
			***REMOVED***)
			assert.EqualError(t, err, "more vus than allocated requested")
		***REMOVED***)
		t.Run("max higher", func(t *testing.T) ***REMOVED***
			e, err := NewEngine(nil, Options***REMOVED***
				VUsMax: null.IntFrom(10),
				VUs:    null.IntFrom(1),
			***REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(10), e.GetVUsMax())
			assert.Equal(t, int64(1), e.GetVUs())
		***REMOVED***)
		t.Run("max just right", func(t *testing.T) ***REMOVED***
			e, err := NewEngine(nil, Options***REMOVED***
				VUsMax: null.IntFrom(10),
				VUs:    null.IntFrom(10),
			***REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(10), e.GetVUsMax())
			assert.Equal(t, int64(10), e.GetVUs())
		***REMOVED***)
	***REMOVED***)
	t.Run("Paused", func(t *testing.T) ***REMOVED***
		t.Run("not set", func(t *testing.T) ***REMOVED***
			e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
			assert.NoError(t, err)
			assert.False(t, e.IsPaused())
		***REMOVED***)
		t.Run("false", func(t *testing.T) ***REMOVED***
			e, err := NewEngine(nil, Options***REMOVED***
				Paused: null.BoolFrom(false),
			***REMOVED***)
			assert.NoError(t, err)
			assert.False(t, e.IsPaused())
		***REMOVED***)
		t.Run("true", func(t *testing.T) ***REMOVED***
			e, err := NewEngine(nil, Options***REMOVED***
				Paused: null.BoolFrom(true),
			***REMOVED***)
			assert.NoError(t, err)
			assert.True(t, e.IsPaused())
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestEngineRun(t *testing.T) ***REMOVED***
	t.Run("exits with context", func(t *testing.T) ***REMOVED***
		startTime := time.Now()
		duration := 100 * time.Millisecond
		e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)

		ctx, _ := context.WithTimeout(context.Background(), duration)
		assert.NoError(t, e.Run(ctx))
		assert.WithinDuration(t, startTime.Add(duration), time.Now(), 100*time.Millisecond)
	***REMOVED***)
	t.Run("terminates subctx", func(t *testing.T) ***REMOVED***
		e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)

		subctx := e.subctx
		select ***REMOVED***
		case <-subctx.Done():
			assert.Fail(t, "context is already terminated")
		default:
		***REMOVED***

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		e.Run(ctx)

		assert.NotEqual(t, subctx, e.subctx, "subcontext not changed")
		select ***REMOVED***
		case <-subctx.Done():
		default:
			assert.Fail(t, "context was not terminated")
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestEngineIsRunning(t *testing.T) ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
	assert.NoError(t, err)

	go e.Run(ctx)
	runtime.Gosched()
	assert.True(t, e.IsRunning())

	cancel()
	runtime.Gosched()
	assert.False(t, e.IsRunning())
***REMOVED***

func TestEngineSetPaused(t *testing.T) ***REMOVED***
	e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
	assert.NoError(t, err)
	assert.False(t, e.IsPaused())

	e.SetPaused(true)
	assert.True(t, e.IsPaused())

	e.SetPaused(false)
	assert.False(t, e.IsPaused())
***REMOVED***

func TestEngineSetVUsMax(t *testing.T) ***REMOVED***
	t.Run("not set", func(t *testing.T) ***REMOVED***
		e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), e.GetVUsMax())
	***REMOVED***)
	t.Run("set", func(t *testing.T) ***REMOVED***
		e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)
		assert.NoError(t, e.SetVUsMax(10))
		assert.Equal(t, int64(10), e.GetVUsMax())
	***REMOVED***)
	t.Run("set negative", func(t *testing.T) ***REMOVED***
		e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)
		assert.EqualError(t, e.SetVUsMax(-1), "vus-max can't be negative")
	***REMOVED***)
	t.Run("set too low", func(t *testing.T) ***REMOVED***
		e, err := NewEngine(nil, Options***REMOVED***
			VUsMax: null.IntFrom(10),
			VUs:    null.IntFrom(10),
		***REMOVED***)
		assert.NoError(t, err)
		assert.EqualError(t, e.SetVUsMax(5), "can't reduce vus-max below vus")
	***REMOVED***)
***REMOVED***

func TestEngineSetVUs(t *testing.T) ***REMOVED***
	t.Run("not set", func(t *testing.T) ***REMOVED***
		e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), e.GetVUsMax())
		assert.Equal(t, int64(0), e.GetVUs())
	***REMOVED***)
	t.Run("set", func(t *testing.T) ***REMOVED***
		e, err := NewEngine(nil, Options***REMOVED***VUsMax: null.IntFrom(10)***REMOVED***)
		assert.NoError(t, err)
		assert.NoError(t, e.SetVUs(10))
		assert.Equal(t, int64(10), e.GetVUs())
	***REMOVED***)
	t.Run("set too high", func(t *testing.T) ***REMOVED***
		e, err := NewEngine(nil, Options***REMOVED***VUsMax: null.IntFrom(10)***REMOVED***)
		assert.NoError(t, err)
		assert.EqualError(t, e.SetVUs(20), "more vus than allocated requested")
	***REMOVED***)
***REMOVED***
