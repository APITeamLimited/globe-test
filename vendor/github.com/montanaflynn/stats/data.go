package stats

// Float64Data is a named type for []float64 with helper methods
type Float64Data []float64

// Get item in slice
func (f Float64Data) Get(i int) float64 ***REMOVED*** return f[i] ***REMOVED***

// Len returns length of slice
func (f Float64Data) Len() int ***REMOVED*** return len(f) ***REMOVED***

// Less returns if one number is less than another
func (f Float64Data) Less(i, j int) bool ***REMOVED*** return f[i] < f[j] ***REMOVED***

// Swap switches out two numbers in slice
func (f Float64Data) Swap(i, j int) ***REMOVED*** f[i], f[j] = f[j], f[i] ***REMOVED***

// Min returns the minimum number in the data
func (f Float64Data) Min() (float64, error) ***REMOVED*** return Min(f) ***REMOVED***

// Max returns the maximum number in the data
func (f Float64Data) Max() (float64, error) ***REMOVED*** return Max(f) ***REMOVED***

// Sum returns the total of all the numbers in the data
func (f Float64Data) Sum() (float64, error) ***REMOVED*** return Sum(f) ***REMOVED***

// Mean returns the mean of the data
func (f Float64Data) Mean() (float64, error) ***REMOVED*** return Mean(f) ***REMOVED***

// Median returns the median of the data
func (f Float64Data) Median() (float64, error) ***REMOVED*** return Median(f) ***REMOVED***

// Mode returns the mode of the data
func (f Float64Data) Mode() ([]float64, error) ***REMOVED*** return Mode(f) ***REMOVED***

// GeometricMean returns the median of the data
func (f Float64Data) GeometricMean() (float64, error) ***REMOVED*** return GeometricMean(f) ***REMOVED***

// HarmonicMean returns the mode of the data
func (f Float64Data) HarmonicMean() (float64, error) ***REMOVED*** return HarmonicMean(f) ***REMOVED***

// MedianAbsoluteDeviation the median of the absolute deviations from the dataset median
func (f Float64Data) MedianAbsoluteDeviation() (float64, error) ***REMOVED***
	return MedianAbsoluteDeviation(f)
***REMOVED***

// MedianAbsoluteDeviationPopulation finds the median of the absolute deviations from the population median
func (f Float64Data) MedianAbsoluteDeviationPopulation() (float64, error) ***REMOVED***
	return MedianAbsoluteDeviationPopulation(f)
***REMOVED***

// StandardDeviation the amount of variation in the dataset
func (f Float64Data) StandardDeviation() (float64, error) ***REMOVED***
	return StandardDeviation(f)
***REMOVED***

// StandardDeviationPopulation finds the amount of variation from the population
func (f Float64Data) StandardDeviationPopulation() (float64, error) ***REMOVED***
	return StandardDeviationPopulation(f)
***REMOVED***

// StandardDeviationSample finds the amount of variation from a sample
func (f Float64Data) StandardDeviationSample() (float64, error) ***REMOVED***
	return StandardDeviationSample(f)
***REMOVED***

// QuartileOutliers finds the mild and extreme outliers
func (f Float64Data) QuartileOutliers() (Outliers, error) ***REMOVED***
	return QuartileOutliers(f)
***REMOVED***

// Percentile finds the relative standing in a slice of floats
func (f Float64Data) Percentile(p float64) (float64, error) ***REMOVED***
	return Percentile(f, p)
***REMOVED***

// PercentileNearestRank finds the relative standing using the Nearest Rank method
func (f Float64Data) PercentileNearestRank(p float64) (float64, error) ***REMOVED***
	return PercentileNearestRank(f, p)
***REMOVED***

// Correlation describes the degree of relationship between two sets of data
func (f Float64Data) Correlation(d Float64Data) (float64, error) ***REMOVED***
	return Correlation(f, d)
***REMOVED***

// Pearson calculates the Pearson product-moment correlation coefficient between two variables.
func (f Float64Data) Pearson(d Float64Data) (float64, error) ***REMOVED***
	return Pearson(f, d)
***REMOVED***

// Quartile returns the three quartile points from a slice of data
func (f Float64Data) Quartile(d Float64Data) (Quartiles, error) ***REMOVED***
	return Quartile(d)
***REMOVED***

// InterQuartileRange finds the range between Q1 and Q3
func (f Float64Data) InterQuartileRange() (float64, error) ***REMOVED***
	return InterQuartileRange(f)
***REMOVED***

// Midhinge finds the average of the first and third quartiles
func (f Float64Data) Midhinge(d Float64Data) (float64, error) ***REMOVED***
	return Midhinge(d)
***REMOVED***

// Trimean finds the average of the median and the midhinge
func (f Float64Data) Trimean(d Float64Data) (float64, error) ***REMOVED***
	return Trimean(d)
***REMOVED***

// Sample returns sample from input with replacement or without
func (f Float64Data) Sample(n int, r bool) ([]float64, error) ***REMOVED***
	return Sample(f, n, r)
***REMOVED***

// Variance the amount of variation in the dataset
func (f Float64Data) Variance() (float64, error) ***REMOVED***
	return Variance(f)
***REMOVED***

// PopulationVariance finds the amount of variance within a population
func (f Float64Data) PopulationVariance() (float64, error) ***REMOVED***
	return PopulationVariance(f)
***REMOVED***

// SampleVariance finds the amount of variance within a sample
func (f Float64Data) SampleVariance() (float64, error) ***REMOVED***
	return SampleVariance(f)
***REMOVED***

// Covariance is a measure of how much two sets of data change
func (f Float64Data) Covariance(d Float64Data) (float64, error) ***REMOVED***
	return Covariance(f, d)
***REMOVED***

// CovariancePopulation computes covariance for entire population between two variables.
func (f Float64Data) CovariancePopulation(d Float64Data) (float64, error) ***REMOVED***
	return CovariancePopulation(f, d)
***REMOVED***
