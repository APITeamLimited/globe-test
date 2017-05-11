package cloud

import (
	"os"
	"testing"

	"github.com/loadimpact/k6/lib"
	"github.com/stretchr/testify/assert"
)

func TestGetName(t *testing.T) ***REMOVED***
	nameTests := []struct ***REMOVED***
		lib      *lib.SourceData
		conf     loadimpactConfig
		expected string
	***REMOVED******REMOVED***
		***REMOVED***&lib.SourceData***REMOVED***Filename: ""***REMOVED***, loadimpactConfig***REMOVED******REMOVED***, TestName***REMOVED***,
		***REMOVED***&lib.SourceData***REMOVED***Filename: "-"***REMOVED***, loadimpactConfig***REMOVED******REMOVED***, TestName***REMOVED***,
		***REMOVED***&lib.SourceData***REMOVED***Filename: "script.js"***REMOVED***, loadimpactConfig***REMOVED******REMOVED***, "script.js"***REMOVED***,
		***REMOVED***&lib.SourceData***REMOVED***Filename: "/file/name.js"***REMOVED***, loadimpactConfig***REMOVED******REMOVED***, "name.js"***REMOVED***,
		***REMOVED***&lib.SourceData***REMOVED***Filename: "/file/name"***REMOVED***, loadimpactConfig***REMOVED******REMOVED***, "name"***REMOVED***,
		***REMOVED***&lib.SourceData***REMOVED***Filename: "/file/name"***REMOVED***, loadimpactConfig***REMOVED***Name: "confName"***REMOVED***, "confName"***REMOVED***,
	***REMOVED***

	for _, test := range nameTests ***REMOVED***
		actual := getName(test.lib, test.conf)
		assert.Equal(t, actual, test.expected)
	***REMOVED***

	err := os.Setenv("K6CLOUD_NAME", "envname")
	assert.Nil(t, err)

	for _, test := range nameTests ***REMOVED***
		actual := getName(test.lib, test.conf)
		assert.Equal(t, actual, "envname")

	***REMOVED***
***REMOVED***
