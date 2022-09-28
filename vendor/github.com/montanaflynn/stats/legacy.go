package stats

// VarP is a shortcut to PopulationVariance
func VarP(input Float64Data) (sdev float64, err error) ***REMOVED***
	return PopulationVariance(input)
***REMOVED***

// VarS is a shortcut to SampleVariance
func VarS(input Float64Data) (sdev float64, err error) ***REMOVED***
	return SampleVariance(input)
***REMOVED***

// StdDevP is a shortcut to StandardDeviationPopulation
func StdDevP(input Float64Data) (sdev float64, err error) ***REMOVED***
	return StandardDeviationPopulation(input)
***REMOVED***

// StdDevS is a shortcut to StandardDeviationSample
func StdDevS(input Float64Data) (sdev float64, err error) ***REMOVED***
	return StandardDeviationSample(input)
***REMOVED***

// LinReg is a shortcut to LinearRegression
func LinReg(s []Coordinate) (regressions []Coordinate, err error) ***REMOVED***
	return LinearRegression(s)
***REMOVED***

// ExpReg is a shortcut to ExponentialRegression
func ExpReg(s []Coordinate) (regressions []Coordinate, err error) ***REMOVED***
	return ExponentialRegression(s)
***REMOVED***

// LogReg is a shortcut to LogarithmicRegression
func LogReg(s []Coordinate) (regressions []Coordinate, err error) ***REMOVED***
	return LogarithmicRegression(s)
***REMOVED***
