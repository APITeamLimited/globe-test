package scheduler

import (
	"fmt"
	"math"
	"strings"
)

// A helper function to verify percentage distributions
func checkPercentagesSum(percentages []float64) error ***REMOVED***
	var sum float64
	for _, v := range percentages ***REMOVED***
		sum += v
	***REMOVED***
	if math.Abs(100-sum) >= minPercentage ***REMOVED***
		return fmt.Errorf("split percentage sum is %.2f while it should be 100", sum)
	***REMOVED***
	return nil
***REMOVED***

// A helper function for joining error messages into a single string
func concatErrors(errors []error, separator string) string ***REMOVED***
	errStrings := make([]string, len(errors))
	for i, e := range errors ***REMOVED***
		errStrings[i] = e.Error()
	***REMOVED***
	return strings.Join(errStrings, separator)
***REMOVED***
