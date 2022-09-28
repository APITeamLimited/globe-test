package stats

import "math"

// Correlation describes the degree of relationship between two sets of data
func Correlation(data1, data2 Float64Data) (float64, error) ***REMOVED***

	l1 := data1.Len()
	l2 := data2.Len()

	if l1 == 0 || l2 == 0 ***REMOVED***
		return math.NaN(), EmptyInput
	***REMOVED***

	if l1 != l2 ***REMOVED***
		return math.NaN(), SizeErr
	***REMOVED***

	sdev1, _ := StandardDeviationPopulation(data1)
	sdev2, _ := StandardDeviationPopulation(data2)

	if sdev1 == 0 || sdev2 == 0 ***REMOVED***
		return 0, nil
	***REMOVED***

	covp, _ := CovariancePopulation(data1, data2)
	return covp / (sdev1 * sdev2), nil
***REMOVED***

// Pearson calculates the Pearson product-moment correlation coefficient between two variables.
func Pearson(data1, data2 Float64Data) (float64, error) ***REMOVED***
	return Correlation(data1, data2)
***REMOVED***
