package stats

import "math"

// _variance finds the variance for both population and sample data
func _variance(input Float64Data, sample int) (variance float64, err error) ***REMOVED***

	if input.Len() == 0 ***REMOVED***
		return math.NaN(), EmptyInput
	***REMOVED***

	// Sum the square of the mean subtracted from each number
	m, _ := Mean(input)

	for _, n := range input ***REMOVED***
		variance += (float64(n) - m) * (float64(n) - m)
	***REMOVED***

	// When getting the mean of the squared differences
	// "sample" will allow us to know if it's a sample
	// or population and wether to subtract by one or not
	return variance / float64((input.Len() - (1 * sample))), nil
***REMOVED***

// Variance the amount of variation in the dataset
func Variance(input Float64Data) (sdev float64, err error) ***REMOVED***
	return PopulationVariance(input)
***REMOVED***

// PopulationVariance finds the amount of variance within a population
func PopulationVariance(input Float64Data) (pvar float64, err error) ***REMOVED***

	v, err := _variance(input, 0)
	if err != nil ***REMOVED***
		return math.NaN(), err
	***REMOVED***

	return v, nil
***REMOVED***

// SampleVariance finds the amount of variance within a sample
func SampleVariance(input Float64Data) (svar float64, err error) ***REMOVED***

	v, err := _variance(input, 1)
	if err != nil ***REMOVED***
		return math.NaN(), err
	***REMOVED***

	return v, nil
***REMOVED***

// Covariance is a measure of how much two sets of data change
func Covariance(data1, data2 Float64Data) (float64, error) ***REMOVED***

	l1 := data1.Len()
	l2 := data2.Len()

	if l1 == 0 || l2 == 0 ***REMOVED***
		return math.NaN(), EmptyInput
	***REMOVED***

	if l1 != l2 ***REMOVED***
		return math.NaN(), SizeErr
	***REMOVED***

	m1, _ := Mean(data1)
	m2, _ := Mean(data2)

	// Calculate sum of squares
	var ss float64
	for i := 0; i < l1; i++ ***REMOVED***
		delta1 := (data1.Get(i) - m1)
		delta2 := (data2.Get(i) - m2)
		ss += (delta1*delta2 - ss) / float64(i+1)
	***REMOVED***

	return ss * float64(l1) / float64(l1-1), nil
***REMOVED***

// CovariancePopulation computes covariance for entire population between two variables.
func CovariancePopulation(data1, data2 Float64Data) (float64, error) ***REMOVED***

	l1 := data1.Len()
	l2 := data2.Len()

	if l1 == 0 || l2 == 0 ***REMOVED***
		return math.NaN(), EmptyInput
	***REMOVED***

	if l1 != l2 ***REMOVED***
		return math.NaN(), SizeErr
	***REMOVED***

	m1, _ := Mean(data1)
	m2, _ := Mean(data2)

	var s float64
	for i := 0; i < l1; i++ ***REMOVED***
		delta1 := (data1.Get(i) - m1)
		delta2 := (data2.Get(i) - m2)
		s += delta1 * delta2
	***REMOVED***

	return s / float64(l1), nil
***REMOVED***
