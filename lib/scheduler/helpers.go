package scheduler

import (
	"bytes"
	"encoding/json"
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

// Decode a JSON in a strict manner, emitting an error if there are unknown fields
func strictJSONUnmarshal(data []byte, v interface***REMOVED******REMOVED***) error ***REMOVED***
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	dec.UseNumber()

	if err := dec.Decode(&v); err != nil ***REMOVED***
		return err
	***REMOVED***
	if dec.More() ***REMOVED***
		return fmt.Errorf("unexpected data after the JSON object")
	***REMOVED***
	return nil
***REMOVED***

// A helper function to avoid code duplication
func validateStages(stages []Stage) []error ***REMOVED***
	var errors []error
	if len(stages) == 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("at least one stage has to be specified"))
	***REMOVED*** else ***REMOVED***
		for i, s := range stages ***REMOVED***
			stageNum := i + 1
			if !s.Duration.Valid ***REMOVED***
				errors = append(errors, fmt.Errorf("stage %d doesn't have a duration", stageNum))
			***REMOVED*** else if s.Duration.Duration < 0 ***REMOVED***
				errors = append(errors, fmt.Errorf("the duration for stage %d shouldn't be negative", stageNum))
			***REMOVED***
			if !s.Target.Valid ***REMOVED***
				errors = append(errors, fmt.Errorf("stage %d doesn't have a target", stageNum))
			***REMOVED*** else if s.Target.Int64 < 0 ***REMOVED***
				errors = append(errors, fmt.Errorf("the target for stage %d shouldn't be negative", stageNum))
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return errors
***REMOVED***
