package accumulate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSumEmpty(t *testing.T) ***REMOVED***
	assert.Equal(t, 0.0, Dimension***REMOVED******REMOVED***.Sum())
***REMOVED***

func TestSum(t *testing.T) ***REMOVED***
	assert.Equal(t, 20.0, Dimension***REMOVED***Values: []float64***REMOVED***5, 3, 4, 7, 1***REMOVED******REMOVED***.Sum())
***REMOVED***

func TestMinEmpty(t *testing.T) ***REMOVED***
	assert.Equal(t, 0.0, Dimension***REMOVED******REMOVED***.Min())
***REMOVED***

func TestMin(t *testing.T) ***REMOVED***
	assert.Equal(t, 1.0, Dimension***REMOVED***Values: []float64***REMOVED***5, 3, 4, 7, 1***REMOVED******REMOVED***.Min())
***REMOVED***

func TestMaxEmpty(t *testing.T) ***REMOVED***
	assert.Equal(t, 0.0, Dimension***REMOVED******REMOVED***.Max())
***REMOVED***

func TestMax(t *testing.T) ***REMOVED***
	assert.Equal(t, 7.0, Dimension***REMOVED***Values: []float64***REMOVED***5, 3, 4, 7, 1***REMOVED******REMOVED***.Max())
***REMOVED***

func TestAvgEmpty(t *testing.T) ***REMOVED***
	assert.Equal(t, 0.0, Dimension***REMOVED******REMOVED***.Avg())
***REMOVED***

func TestAvgOne(t *testing.T) ***REMOVED***
	assert.Equal(t, 5.0, Dimension***REMOVED***Values: []float64***REMOVED***5***REMOVED******REMOVED***.Avg())
***REMOVED***

func TestAvgTwo(t *testing.T) ***REMOVED***
	assert.Equal(t, 4.0, Dimension***REMOVED***Values: []float64***REMOVED***5, 3***REMOVED******REMOVED***.Avg())
***REMOVED***

func TestAvgThree(t *testing.T) ***REMOVED***
	assert.Equal(t, 4.0, Dimension***REMOVED***Values: []float64***REMOVED***5, 3, 4***REMOVED******REMOVED***.Avg())
***REMOVED***

func TestMedEmpty(t *testing.T) ***REMOVED***
	assert.Equal(t, 0.0, Dimension***REMOVED******REMOVED***.Med())
***REMOVED***

func TestMedOne(t *testing.T) ***REMOVED***
	assert.Equal(t, 5.0, Dimension***REMOVED***Values: []float64***REMOVED***5***REMOVED******REMOVED***.Med())
***REMOVED***

func TestMedTwo(t *testing.T) ***REMOVED***
	assert.Equal(t, 4.0, Dimension***REMOVED***Values: []float64***REMOVED***5, 3***REMOVED******REMOVED***.Med())
***REMOVED***

func TestMedThree(t *testing.T) ***REMOVED***
	assert.Equal(t, 3.0, Dimension***REMOVED***Values: []float64***REMOVED***5, 3, 4***REMOVED******REMOVED***.Med())
***REMOVED***

func TestMedFour(t *testing.T) ***REMOVED***
	assert.Equal(t, 3.5, Dimension***REMOVED***Values: []float64***REMOVED***5, 3, 4, 7***REMOVED******REMOVED***.Med())
***REMOVED***

func TestMedFive(t *testing.T) ***REMOVED***
	assert.Equal(t, 4.0, Dimension***REMOVED***Values: []float64***REMOVED***5, 3, 4, 7, 1***REMOVED******REMOVED***.Med())
***REMOVED***
