package compiler

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNew(t *testing.T) ***REMOVED***
	c, err := New()
	assert.NotNil(t, c)
	assert.NoError(t, err)
***REMOVED***

func TestTransform(t *testing.T) ***REMOVED***
	c, err := New()
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	t.Run("blank", func(t *testing.T) ***REMOVED***
		src, srcmap, err := c.Transform("", "test.js")
		assert.NoError(t, err)
		assert.Equal(t, `"use strict";`, src)
		assert.Equal(t, 3, srcmap.Version)
		assert.Equal(t, "test.js", srcmap.File)
		assert.Equal(t, "", srcmap.Mappings)
	***REMOVED***)
	t.Run("double-arrow", func(t *testing.T) ***REMOVED***
		src, srcmap, err := c.Transform("()=> true", "test.js")
		assert.NoError(t, err)
		assert.Equal(t, `"use strict";(function () ***REMOVED***return true;***REMOVED***);`, src)
		assert.Equal(t, 3, srcmap.Version)
		assert.Equal(t, "test.js", srcmap.File)
		assert.Equal(t, "aAAA,qBAAK,IAAL", srcmap.Mappings)
	***REMOVED***)
	t.Run("longer", func(t *testing.T) ***REMOVED***
		src, srcmap, err := c.Transform(strings.Join([]string***REMOVED***
			`function add(a, b) ***REMOVED***`,
			`    return a + b;`,
			`***REMOVED***;`,
			``,
			`let res = add(1, 2);`,
		***REMOVED***, "\n"), "test.js")
		println(src)
		assert.NoError(t, err)
		assert.Equal(t, strings.Join([]string***REMOVED***
			`"use strict";function add(a, b) ***REMOVED***`,
			`    return a + b;`,
			`***REMOVED***;`,
			``,
			`var res = add(1, 2);`,
		***REMOVED***, "\n"), src)
		assert.Equal(t, 3, srcmap.Version)
		assert.Equal(t, "test.js", srcmap.File)
		assert.Equal(t, "aAAA,SAASA,GAAT,CAAaC,CAAb,EAAgBC,CAAhB,EAAmB;AACf,WAAOD,IAAIC,CAAX;AACH;;AAED,IAAIC,MAAMH,IAAI,CAAJ,EAAO,CAAP,CAAV", srcmap.Mappings)
	***REMOVED***)
***REMOVED***
