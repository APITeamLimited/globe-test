package lib

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestLerp(t *testing.T) ***REMOVED***
	// data[x][y][t] = v
	data := map[int64]map[int64]map[float64]int64***REMOVED***
		0: map[int64]map[float64]int64***REMOVED***
			0:   map[float64]int64***REMOVED***0.0: 0, 0.10: 0, 0.5: 0, 1.0: 0***REMOVED***,
			100: map[float64]int64***REMOVED***0.0: 0, 0.10: 10, 0.5: 50, 1.0: 100***REMOVED***,
			500: map[float64]int64***REMOVED***0.0: 0, 0.10: 50, 0.5: 250, 1.0: 500***REMOVED***,
		***REMOVED***,
		100: map[int64]map[float64]int64***REMOVED***
			200: map[float64]int64***REMOVED***0.0: 100, 0.1: 110, 0.5: 150, 1.0: 200***REMOVED***,
			0:   map[float64]int64***REMOVED***0.0: 100, 0.1: 90, 0.5: 50, 1.0: 0***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for x, data := range data ***REMOVED***
		t.Run("x="+strconv.FormatInt(x, 10), func(t *testing.T) ***REMOVED***
			for y, data := range data ***REMOVED***
				t.Run("y="+strconv.FormatInt(y, 10), func(t *testing.T) ***REMOVED***
					for t_, x1 := range data ***REMOVED***
						t.Run("t="+strconv.FormatFloat(t_, 'f', 2, 64), func(t *testing.T) ***REMOVED***
							assert.Equal(t, x1, Lerp(x, y, t_))
						***REMOVED***)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
