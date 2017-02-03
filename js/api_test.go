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
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSleep(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		return
	***REMOVED***

	testdata := map[string]struct ***REMOVED***
		src string
		min time.Duration
	***REMOVED******REMOVED***
		"float,sub-1s": ***REMOVED***`0.2`, 200 * time.Millisecond***REMOVED***,
		"float":        ***REMOVED***`1.0`, 1 * time.Second***REMOVED***,
		"int":          ***REMOVED***`1`, 1 * time.Second***REMOVED***,
		"exceeding":    ***REMOVED***`5`, 2 * time.Second***REMOVED***,
	***REMOVED***
	for name, data := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			r, err := newSnippetRunner(fmt.Sprintf(`
			import ***REMOVED*** sleep ***REMOVED*** from "k6";
			export default function() ***REMOVED***
				sleep(%s);
			***REMOVED***`, data.src))
			assert.NoError(t, err)

			vu, err := r.NewVU()
			assert.NoError(t, err)

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			start := time.Now()

			_, err = vu.RunOnce(ctx)
			assert.NoError(t, err)
			assert.True(t, time.Since(start) > data.min, "ran too short")
			assert.True(t, time.Since(start) < data.min+1*time.Second, "ran too long")
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestDoGroup(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		return
	***REMOVED***

	r, err := newSnippetRunner(`
	import ***REMOVED*** group ***REMOVED*** from "k6";
	export default function() ***REMOVED***
		group("test", fn);
	***REMOVED***`)
	assert.NoError(t, err)

	vu_, err := r.NewVU()
	assert.NoError(t, err)
	vu := vu_.(*VU)

	assert.NoError(t, vu.vm.Set("fn", func() ***REMOVED***
		assert.Equal(t, "test", vu.group.Name)
	***REMOVED***))

	_, err = vu.RunOnce(context.Background())
	assert.NoError(t, err)
***REMOVED***

func TestDoGroupNested(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		return
	***REMOVED***

	r, err := newSnippetRunner(`
	import ***REMOVED*** group ***REMOVED*** from "k6";
	export default function() ***REMOVED***
		group("outer", function() ***REMOVED***
			group("inner", fn);
		***REMOVED***);
	***REMOVED***`)
	assert.NoError(t, err)

	vu_, err := r.NewVU()
	assert.NoError(t, err)
	vu := vu_.(*VU)

	assert.NoError(t, vu.vm.Set("fn", func() ***REMOVED***
		assert.Equal(t, "inner", vu.group.Name)
		assert.Equal(t, "outer", vu.group.Parent.Name)
	***REMOVED***))

	_, err = vu.RunOnce(context.Background())
	assert.NoError(t, err)
***REMOVED***

func TestDoGroupReturn(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		return
	***REMOVED***

	r, err := newSnippetRunner(`
	import ***REMOVED*** group, _assert ***REMOVED*** from "k6";
	export default function() ***REMOVED***
		let v = group("group", function() ***REMOVED***
			return 12345;
		***REMOVED***);
		_assert(v === 12345);
	***REMOVED***`)
	assert.NoError(t, err)

	vu, err := r.NewVU()
	assert.NoError(t, err)
	_, err = vu.RunOnce(context.Background())
	assert.NoError(t, err)
***REMOVED***

func TestDoGroupReturnTrueByDefault(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		return
	***REMOVED***

	r, err := newSnippetRunner(`
	import ***REMOVED*** group, _assert ***REMOVED*** from "k6";
	export default function() ***REMOVED***
		let v = group("group", function() ***REMOVED***
			// no return
		***REMOVED***);
		_assert(v === true);
	***REMOVED***`)
	assert.NoError(t, err)

	vu, err := r.NewVU()
	assert.NoError(t, err)
	_, err = vu.RunOnce(context.Background())
	assert.NoError(t, err)
***REMOVED***

func TestDoCheck(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		return
	***REMOVED***

	r, err := newSnippetRunner(`
	import ***REMOVED*** check ***REMOVED*** from "k6";
	export default function() ***REMOVED***
		check(3, ***REMOVED*** "v === 3": (v) => v === 3 ***REMOVED***);
	***REMOVED***`)
	assert.NoError(t, err)

	vu_, err := r.NewVU()
	assert.NoError(t, err)
	vu := vu_.(*VU)

	_, err = vu.RunOnce(context.Background())
	assert.NoError(t, err)

	c := r.DefaultGroup.Checks["v === 3"]
	assert.NotNil(t, c)
	assert.Equal(t, "v === 3", c.Name)
	assert.Equal(t, r.DefaultGroup, c.Group)
	assert.Equal(t, int64(1), c.Passes)
	assert.Equal(t, int64(0), c.Fails)
***REMOVED***

func TestCheckInGroup(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		return
	***REMOVED***

	r, err := newSnippetRunner(`
	import ***REMOVED*** group, check ***REMOVED*** from "k6";
	export default function() ***REMOVED***
		group("group", function() ***REMOVED***
			check(3, ***REMOVED*** "v === 3": (v) => v === 3 ***REMOVED***);
		***REMOVED***);
	***REMOVED***`)
	assert.NoError(t, err)

	vu_, err := r.NewVU()
	assert.NoError(t, err)
	vu := vu_.(*VU)

	_, err = vu.RunOnce(context.Background())
	assert.NoError(t, err)

	g := r.DefaultGroup.Groups["group"]
	assert.NotNil(t, g)
	assert.Equal(t, "group", g.Name)

	c := g.Checks["v === 3"]
	assert.NotNil(t, c)
	assert.Equal(t, "v === 3", c.Name)
	assert.Equal(t, g, c.Group)
	assert.Equal(t, int64(1), c.Passes)
	assert.Equal(t, int64(0), c.Fails)
***REMOVED***

func TestCheckReturnTrueOnSuccess(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		return
	***REMOVED***

	r, err := newSnippetRunner(`
	import ***REMOVED*** check, _assert ***REMOVED*** from "k6";
	export default function() ***REMOVED***
		let succ = check(null, ***REMOVED*** "true": true ***REMOVED***);
		_assert(succ === true);
	***REMOVED***`)
	assert.NoError(t, err)

	vu, err := r.NewVU()
	assert.NoError(t, err)
	_, err = vu.RunOnce(context.Background())
	assert.NoError(t, err)
***REMOVED***

func TestCheckReturnFalseOnFailure(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		return
	***REMOVED***

	r, err := newSnippetRunner(`
	import ***REMOVED*** check, _assert ***REMOVED*** from "k6";
	export default function() ***REMOVED***
		let succ = check(null, ***REMOVED*** "false": false ***REMOVED***);
		_assert(succ === false);
	***REMOVED***`)
	assert.NoError(t, err)

	vu, err := r.NewVU()
	assert.NoError(t, err)
	_, err = vu.RunOnce(context.Background())
	assert.NoError(t, err)
***REMOVED***
