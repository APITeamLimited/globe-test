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

package js

import (
	"context"
	"errors"
	"github.com/robertkrimen/otto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newSnippetRunner(src string) (*Runner, error) ***REMOVED***
	rt, err := New()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	_ = rt.VM.Set("__initapi__", &InitAPI***REMOVED***r: rt***REMOVED***)
	exp, err := rt.load("__snippet__", []byte(src))
	_ = rt.VM.Set("__initapi__", nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return NewRunner(rt, exp)
***REMOVED***

func runSnippet(src string) error ***REMOVED***
	r, err := newSnippetRunner(src)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	vu, err := r.NewVU()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = vu.RunOnce(context.Background())
	return err
***REMOVED***

func TestCheck(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		return
	***REMOVED***

	vm := otto.New()

	t.Run("String", func(t *testing.T) ***REMOVED***
		t.Run("Something", func(t *testing.T) ***REMOVED***
			v, err := vm.Eval(`"test"`)
			assert.NoError(t, err)
			b, err := Check(v, otto.UndefinedValue())
			assert.NoError(t, err)
			assert.True(t, b)
		***REMOVED***)

		t.Run("Empty", func(t *testing.T) ***REMOVED***
			v, err := vm.Eval(`""`)
			assert.NoError(t, err)
			b, err := Check(v, otto.UndefinedValue())
			assert.NoError(t, err)
			assert.False(t, b)
		***REMOVED***)
	***REMOVED***)

	t.Run("Number", func(t *testing.T) ***REMOVED***
		t.Run("Positive", func(t *testing.T) ***REMOVED***
			v, err := vm.Eval(`1`)
			assert.NoError(t, err)
			b, err := Check(v, otto.UndefinedValue())
			assert.NoError(t, err)
			assert.True(t, b)
		***REMOVED***)
		t.Run("Negative", func(t *testing.T) ***REMOVED***
			v, err := vm.Eval(`-1`)
			assert.NoError(t, err)
			b, err := Check(v, otto.UndefinedValue())
			assert.NoError(t, err)
			assert.True(t, b)
		***REMOVED***)
		t.Run("Zero", func(t *testing.T) ***REMOVED***
			v, err := vm.Eval(`0`)
			assert.NoError(t, err)
			b, err := Check(v, otto.UndefinedValue())
			assert.NoError(t, err)
			assert.False(t, b)
		***REMOVED***)
	***REMOVED***)

	t.Run("Boolean", func(t *testing.T) ***REMOVED***
		t.Run("True", func(t *testing.T) ***REMOVED***
			v, err := vm.Eval(`true`)
			assert.NoError(t, err)
			b, err := Check(v, otto.UndefinedValue())
			assert.NoError(t, err)
			assert.True(t, b)
		***REMOVED***)
		t.Run("False", func(t *testing.T) ***REMOVED***
			v, err := vm.Eval(`false`)
			assert.NoError(t, err)
			b, err := Check(v, otto.UndefinedValue())
			assert.NoError(t, err)
			assert.False(t, b)
		***REMOVED***)
	***REMOVED***)

	t.Run("Function", func(t *testing.T) ***REMOVED***
		fn, err := vm.Eval(`(function(v) ***REMOVED*** return v === true; ***REMOVED***)`)
		assert.NoError(t, err)

		t.Run("True", func(t *testing.T) ***REMOVED***
			b, err := Check(fn, otto.TrueValue())
			assert.NoError(t, err)
			assert.True(t, b)
		***REMOVED***)
		t.Run("False", func(t *testing.T) ***REMOVED***
			b, err := Check(fn, otto.FalseValue())
			assert.NoError(t, err)
			assert.False(t, b)
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestThrow(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		return
	***REMOVED***

	vm := otto.New()
	assert.NoError(t, vm.Set("fn", func() ***REMOVED***
		throw(vm, errors.New("This is a test error"))
	***REMOVED***))
	_, err := vm.Eval(`fn()`)
	assert.EqualError(t, err, "Error: This is a test error")
***REMOVED***
