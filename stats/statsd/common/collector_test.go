package common

import (
	"testing"

	"github.com/loadimpact/k6/stats"

	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v3"
)

func TestInitWithoutAddressErrors(t *testing.T) ***REMOVED***
	var c = &Collector***REMOVED***
		Config: Config***REMOVED******REMOVED***,
		Type:   "testtype",
	***REMOVED***
	err := c.Init()
	require.Error(t, err)
***REMOVED***

func TestInitWithBogusAddressErrors(t *testing.T) ***REMOVED***
	var c = &Collector***REMOVED***
		Config: Config***REMOVED***
			Addr: null.StringFrom("localhost:90000"),
		***REMOVED***,
		Type: "testtype",
	***REMOVED***
	err := c.Init()
	require.Error(t, err)
***REMOVED***

func TestLinkReturnAddress(t *testing.T) ***REMOVED***
	var bogusValue = "bogus value"
	var c = &Collector***REMOVED***
		Config: Config***REMOVED***
			Addr: null.StringFrom(bogusValue),
		***REMOVED***,
	***REMOVED***
	require.Equal(t, bogusValue, c.Link())
***REMOVED***

func TestGetRequiredSystemTags(t *testing.T) ***REMOVED***
	var c = &Collector***REMOVED******REMOVED***
	require.Equal(t, stats.SystemTagSet(0), c.GetRequiredSystemTags())
***REMOVED***
