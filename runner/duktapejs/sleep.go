package duktapejs

import (
	"gopkg.in/olebedev/go-duktape.v2"
	"time"
)

func (vu *VUContext) Sleep(c *duktape.Context) int ***REMOVED***
	t := c.GetNumber(-1)
	time.Sleep(time.Duration(t) * time.Second)
	return 0
***REMOVED***
