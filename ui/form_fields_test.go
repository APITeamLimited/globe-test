package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringField(t *testing.T) ***REMOVED***
	t.Parallel()
	t.Run("Creation", func(t *testing.T) ***REMOVED***
		t.Parallel()
		f := StringField***REMOVED***Key: "key", Label: "label"***REMOVED***
		assert.Equal(t, "key", f.GetKey())
		assert.Equal(t, "label", f.GetLabel())
	***REMOVED***)

	t.Run("Valid", func(t *testing.T) ***REMOVED***
		t.Parallel()
		f := StringField***REMOVED***Key: "key", Label: "label"***REMOVED***
		v, err := f.Clean("uwu")
		assert.NoError(t, err)
		assert.Equal(t, "uwu", v)
	***REMOVED***)
	t.Run("Whitespace", func(t *testing.T) ***REMOVED***
		t.Parallel()
		f := StringField***REMOVED***Key: "key", Label: "label"***REMOVED***
		v, err := f.Clean("\r\n\t ")
		assert.NoError(t, err)
		assert.Equal(t, "", v)
	***REMOVED***)
	t.Run("Min", func(t *testing.T) ***REMOVED***
		t.Parallel()
		f := StringField***REMOVED***Key: "key", Label: "label"***REMOVED***
		f.Min = 10
		_, err := f.Clean("short")
		assert.EqualError(t, err, "invalid input, min length is 10")
	***REMOVED***)
	t.Run("Max", func(t *testing.T) ***REMOVED***
		t.Parallel()
		f := StringField***REMOVED***Key: "key", Label: "label"***REMOVED***
		f.Max = 10
		_, err := f.Clean("too dang long")
		assert.EqualError(t, err, "invalid input, max length is 10")
	***REMOVED***)
	t.Run("Default", func(t *testing.T) ***REMOVED***
		t.Parallel()
		f := StringField***REMOVED***Key: "key", Label: "label"***REMOVED***
		f.Default = "default"
		v, err := f.Clean("")
		assert.NoError(t, err)
		assert.Equal(t, "default", v)
	***REMOVED***)
***REMOVED***
