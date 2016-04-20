package v8js

import (
	"strconv"
	"time"
)

func (vu *VUContext) Sleep(ts string) ***REMOVED***
	t, err := strconv.ParseFloat(ts, 64)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	time.Sleep(time.Duration(t) * time.Second)
***REMOVED***
