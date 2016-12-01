package js

import (
	"context"
	"github.com/loadimpact/k6/lib"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSleep(t *testing.T) ***REMOVED***
	start := time.Now()
	JSAPI***REMOVED******REMOVED***.Sleep(0.2)
	assert.True(t, time.Since(start) > 200*time.Millisecond)
	assert.True(t, time.Since(start) < 1*time.Second)
***REMOVED***

func TestDoGroup(t *testing.T) ***REMOVED***
	r, err := newSnippetRunner(`
	import ***REMOVED*** group ***REMOVED*** from "k6";
	export default function() ***REMOVED***
		group("test", fn);
	***REMOVED***`)
	assert.NoError(t, err)

	vu_, err := r.NewVU()
	assert.NoError(t, err)
	vu := vu_.(*VU)

	vu.vm.Set("fn", func() ***REMOVED***
		assert.Equal(t, "test", vu.group.Name)
	***REMOVED***)

	_, err = vu.RunOnce(context.Background())
	assert.NoError(t, err)
***REMOVED***

func TestDoGroupNested(t *testing.T) ***REMOVED***
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

	vu.vm.Set("fn", func() ***REMOVED***
		assert.Equal(t, "inner", vu.group.Name)
		assert.Equal(t, "outer", vu.group.Parent.Name)
	***REMOVED***)

	_, err = vu.RunOnce(context.Background())
	assert.NoError(t, err)
***REMOVED***

func TestDoGroupReturn(t *testing.T) ***REMOVED***
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

	if !assert.Len(t, r.Checks, 1) ***REMOVED***
		return
	***REMOVED***
	c := r.Checks[0]
	assert.Equal(t, "v === 3", c.Name)
	assert.Equal(t, r.DefaultGroup, c.Group)
	assert.Equal(t, int64(1), c.Passes)
	assert.Equal(t, int64(0), c.Fails)
***REMOVED***

func TestCheckInGroup(t *testing.T) ***REMOVED***
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

	assert.Len(t, r.Groups, 2)
	g := r.Groups[1]
	assert.Equal(t, "group", g.Name)

	assert.Len(t, r.Checks, 1)
	c := r.Checks[0]
	assert.Equal(t, "v === 3", c.Name)
	assert.Equal(t, g, c.Group)
	assert.Equal(t, int64(1), c.Passes)
	assert.Equal(t, int64(0), c.Fails)
***REMOVED***

func TestCheckReturnTrueOnSuccess(t *testing.T) ***REMOVED***
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

func TestCheckReturnFalseAndTaintsOnFailure(t *testing.T) ***REMOVED***
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
	assert.Equal(t, lib.ErrVUWantsTaint, err)
***REMOVED***

func TestTaint(t *testing.T) ***REMOVED***
	r, err := newSnippetRunner(`
	import ***REMOVED*** taint ***REMOVED*** from "k6";
	export default function() ***REMOVED***
		taint();
	***REMOVED***`)
	assert.NoError(t, err)

	vu, err := r.NewVU()
	assert.NoError(t, err)

	_, err = vu.RunOnce(context.Background())
	assert.Equal(t, lib.ErrVUWantsTaint, err)
***REMOVED***
