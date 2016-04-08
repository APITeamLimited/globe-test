package testutil

import (
	"github.com/codegangsta/cli"
	"testing"
)

func TestWithAppContext(t *testing.T) ***REMOVED***
	cmd := cli.Command***REMOVED***
		Flags: []cli.Flag***REMOVED***
			cli.StringFlag***REMOVED***Name: "arg"***REMOVED***,
		***REMOVED***,
	***REMOVED***
	err := WithAppContext("--arg value", cmd, func(c *cli.Context) ***REMOVED***
		v := c.String("arg")
		if v != "value" ***REMOVED***
			t.Error("Wrong value:", v)
		***REMOVED***
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***
