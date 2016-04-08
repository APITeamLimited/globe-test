package master

import (
	"testing"
)

func TestMasterAddresses(t *testing.T) ***REMOVED***
	in, out := MasterAddresses("127.0.0.1", 1234)
	if in != "tcp://127.0.0.1:1234" ***REMOVED***
		t.Error("Invalid in address", in)
	***REMOVED***
	if out != "tcp://127.0.0.1:1235" ***REMOVED***
		t.Error("Invalid out address", out)
	***REMOVED***
***REMOVED***
