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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRunner(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		return
	***REMOVED***

	rt, err := New()
	assert.NoError(t, err)
	exp, err := rt.load("test.js", []byte(`export default function() ***REMOVED******REMOVED***`))
	assert.NoError(t, err)
	r, err := NewRunner(rt, exp)
	assert.NoError(t, err)
	if !assert.NotNil(t, r) ***REMOVED***
		return
	***REMOVED***

	t.Run("GetDefaultGroup", func(t *testing.T) ***REMOVED***
		assert.Equal(t, r.DefaultGroup, r.GetDefaultGroup())
	***REMOVED***)

	t.Run("VU", func(t *testing.T) ***REMOVED***
		vu_, err := r.NewVU()
		assert.NoError(t, err)
		vu := vu_.(*VU)

		t.Run("Reconfigure", func(t *testing.T) ***REMOVED***
			assert.NoError(t, vu.Reconfigure(12345))
			assert.Equal(t, int64(12345), vu.ID)
		***REMOVED***)

		t.Run("RunOnce", func(t *testing.T) ***REMOVED***
			_, err := vu.RunOnce(context.Background())
			assert.NoError(t, err)
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestVUSelfIdentity(t *testing.T) ***REMOVED***
	r, err := newSnippetRunner(`
	import ***REMOVED*** _assert ***REMOVED*** from "k6";
	export default function() ***REMOVED******REMOVED***
	`)
	assert.NoError(t, err)

	vu_, err := r.NewVU()
	assert.NoError(t, err)
	vu := vu_.(*VU)

	assert.NoError(t, vu.Reconfigure(1234))
	_, err = vu.vm.Eval(`_assert(__VU == 1234)`)
	_, err = vu.vm.Eval(`_assert(__ITER == 0)`)

	_, err = vu.RunOnce(context.Background())
	assert.NoError(t, err)
	_, err = vu.vm.Eval(`_assert(__VU == 1234)`)
	_, err = vu.vm.Eval(`_assert(__ITER == 1)`)

	_, err = vu.RunOnce(context.Background())
	assert.NoError(t, err)
	_, err = vu.vm.Eval(`_assert(__VU == 1234)`)
	_, err = vu.vm.Eval(`_assert(__ITER == 2)`)

	assert.NoError(t, vu.Reconfigure(1234))
	_, err = vu.vm.Eval(`_assert(__VU == 1234)`)
	_, err = vu.vm.Eval(`_assert(__ITER == 0)`)
***REMOVED***
