package stats

import "math"

// MedianAbsoluteDeviation finds the median of the absolute deviations from the dataset median
func MedianAbsoluteDeviation(input Float64Data) (mad float64, err error) ***REMOVED***
	return MedianAbsoluteDeviationPopulation(input)
***REMOVED***

// MedianAbsoluteDeviationPopulation finds the median of the absolute deviations from the population median
func MedianAbsoluteDeviationPopulation(input Float64Data) (mad float64, err error) ***REMOVED***
	if input.Len() == 0 ***REMOVED***
		return math.NaN(), EmptyInput
	***REMOVED***

	i := copyslice(input)
	m, _ := Median(i)

	for key, value := range i ***REMOVED***
		i[key] = math.Abs(value - m)
	***REMOVED***

	return Median(i)
***REMOVED***

// StandardDeviation the amount of variation in the dataset
func StandardDeviation(input Float64Data) (sdev float64, err error) ***REMOVED***
	return StandardDeviationPopulation(input)
***REMOVED***

// StandardDeviationPopulation finds the amount of variation from the population
func StandardDeviationPopulation(input Float64Data) (sdev float64, err error) ***REMOVED***

	if input.Len() == 0 ***REMOVED***
		return math.NaN(), EmptyInput
	***REMOVED***

	// Get the population variance
	vp, _ := PopulationVariance(input)

	// Return the population standard deviation
	return math.Pow(vp, 0.5), nil
***REMOVED***

// StandardDeviationSample finds the amount of variation from a sample
func StandardDeviationSample(input Float64Data) (sdev float64, err error) ***REMOVED***

	if input.Len() == 0 ***REMOVED***
		return math.NaN(), EmptyInput
	***REMOVED***

	// Get the sample variance
	vs, _ := SampleVariance(input)

	// Return the sample standard deviation
	return math.Pow(vs, 0.5), nil
***REMOVED***
