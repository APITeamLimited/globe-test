package accumulate

import (
	"math"
)

type Dimension struct ***REMOVED***
	Values []float64
	Last   float64

	dirty bool
***REMOVED***

func (d Dimension) Sum() float64 ***REMOVED***
	var sum float64
	for _, v := range d.Values ***REMOVED***
		sum += v
	***REMOVED***
	return sum
***REMOVED***

func (d Dimension) Min() float64 ***REMOVED***
	if len(d.Values) == 0 ***REMOVED***
		return 0
	***REMOVED***

	var min float64 = math.MaxFloat64
	for _, v := range d.Values ***REMOVED***
		if v < min ***REMOVED***
			min = v
		***REMOVED***
	***REMOVED***
	return min
***REMOVED***

func (d Dimension) Max() float64 ***REMOVED***
	var max float64
	for _, v := range d.Values ***REMOVED***
		if v > max ***REMOVED***
			max = v
		***REMOVED***
	***REMOVED***
	return max
***REMOVED***

func (d Dimension) Avg() float64 ***REMOVED***
	l := len(d.Values)
	switch l ***REMOVED***
	case 0:
		return 0
	case 1:
		return d.Values[0]
	default:
		return d.Sum() / float64(l)
	***REMOVED***
***REMOVED***

func (d Dimension) Med() float64 ***REMOVED***
	l := len(d.Values)
	switch ***REMOVED***
	case l == 0:
		// No items: median is 0
		return 0
	case l == 1:
		// One item: median is that one item
		return d.Values[0]
	case (l & 0x01) == 0:
		// Even number of items: median is the mean of the middle values
		return (d.Values[l/2] + d.Values[(l/2)-1]) / 2
	default:
		// Odd number of items: median is the middle value
		return d.Values[l/2]
	***REMOVED***
***REMOVED***
