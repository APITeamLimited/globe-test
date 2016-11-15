package lib

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestEase(t *testing.T) ***REMOVED***
	// x[y][t] = x1, tx = 0, ty = 100
	data := map[int64]map[int64]map[int64]int64***REMOVED***
		0: map[int64]map[int64]int64***REMOVED***
			0:   map[int64]int64***REMOVED***0: 0, 10: 0, 50: 0, 100: 0***REMOVED***,
			100: map[int64]int64***REMOVED***0: 0, 10: 10, 50: 50, 100: 100***REMOVED***,
			500: map[int64]int64***REMOVED***0: 0, 10: 50, 50: 250, 100: 500***REMOVED***,
		***REMOVED***,
		100: map[int64]map[int64]int64***REMOVED***
			200: map[int64]int64***REMOVED***0: 100, 10: 110, 50: 150, 100: 200***REMOVED***,
			0:   map[int64]int64***REMOVED***0: 100, 10: 90, 50: 50, 100: 0***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for x, data := range data ***REMOVED***
		t.Run("x="+strconv.FormatInt(x, 10), func(t *testing.T) ***REMOVED***
			for y, data := range data ***REMOVED***
				t.Run("y="+strconv.FormatInt(y, 10), func(t *testing.T) ***REMOVED***
					for t0, x1 := range data ***REMOVED***
						t.Run("t="+strconv.FormatInt(t0, 10), func(t *testing.T) ***REMOVED***
							assert.Equal(t, x1, Ease(t0, 0, 100, x, y))
						***REMOVED***)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
