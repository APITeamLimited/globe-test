package stats

import (
	"math"
)

// Validate data for distance calculation
func validateData(dataPointX, dataPointY []float64) error ***REMOVED***
	if len(dataPointX) == 0 || len(dataPointY) == 0 ***REMOVED***
		return EmptyInput
	***REMOVED***

	if len(dataPointX) != len(dataPointY) ***REMOVED***
		return SizeErr
	***REMOVED***
	return nil
***REMOVED***

// Computes Chebyshev distance between two data sets
func ChebyshevDistance(dataPointX, dataPointY []float64) (distance float64, err error) ***REMOVED***
	err = validateData(dataPointX, dataPointY)
	if err != nil ***REMOVED***
		return math.NaN(), err
	***REMOVED***
	var tempDistance float64
	for i := 0; i < len(dataPointY); i++ ***REMOVED***
		tempDistance = math.Abs(dataPointX[i] - dataPointY[i])
		if distance < tempDistance ***REMOVED***
			distance = tempDistance
		***REMOVED***
	***REMOVED***
	return distance, nil
***REMOVED***

//
// Computes Euclidean distance between two data sets
//
func EuclideanDistance(dataPointX, dataPointY []float64) (distance float64, err error) ***REMOVED***

	err = validateData(dataPointX, dataPointY)
	if err != nil ***REMOVED***
		return math.NaN(), err
	***REMOVED***
	distance = 0
	for i := 0; i < len(dataPointX); i++ ***REMOVED***
		distance = distance + ((dataPointX[i] - dataPointY[i]) * (dataPointX[i] - dataPointY[i]))
	***REMOVED***
	return math.Sqrt(distance), nil
***REMOVED***

//
// Computes Manhattan distance between two data sets
//
func ManhattanDistance(dataPointX, dataPointY []float64) (distance float64, err error) ***REMOVED***
	err = validateData(dataPointX, dataPointY)
	if err != nil ***REMOVED***
		return math.NaN(), err
	***REMOVED***
	distance = 0
	for i := 0; i < len(dataPointX); i++ ***REMOVED***
		distance = distance + math.Abs(dataPointX[i]-dataPointY[i])
	***REMOVED***
	return distance, nil
***REMOVED***

//
// Computes minkowski distance between two data sets.
//
// Input:
//    dataPointX: First set of data points
//    dataPointY: Second set of data points. Length of both data
//                sets must be equal.
//    lambda:     aka p or city blocks; With lambda = 1
//                returned distance is manhattan distance and
//                lambda = 2; it is euclidean distance. Lambda
//                reaching to infinite - distance would be chebysev
//                distance.
// Output:
//     Distance or error
//
func MinkowskiDistance(dataPointX, dataPointY []float64, lambda float64) (distance float64, err error) ***REMOVED***
	err = validateData(dataPointX, dataPointY)
	if err != nil ***REMOVED***
		return math.NaN(), err
	***REMOVED***
	for i := 0; i < len(dataPointY); i++ ***REMOVED***
		distance = distance + math.Pow(math.Abs(dataPointX[i]-dataPointY[i]), lambda)
	***REMOVED***
	distance = math.Pow(distance, float64(1/lambda))
	if math.IsInf(distance, 1) == true ***REMOVED***
		return math.NaN(), InfValue
	***REMOVED***
	return distance, nil
***REMOVED***
