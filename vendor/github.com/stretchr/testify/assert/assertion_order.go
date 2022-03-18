package assert

import (
	"fmt"
	"reflect"
)

// isOrdered checks that collection contains orderable elements.
func isOrdered(t TestingT, object interface***REMOVED******REMOVED***, allowedComparesResults []CompareType, failMessage string, msgAndArgs ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	objKind := reflect.TypeOf(object).Kind()
	if objKind != reflect.Slice && objKind != reflect.Array ***REMOVED***
		return false
	***REMOVED***

	objValue := reflect.ValueOf(object)
	objLen := objValue.Len()

	if objLen <= 1 ***REMOVED***
		return true
	***REMOVED***

	value := objValue.Index(0)
	valueInterface := value.Interface()
	firstValueKind := value.Kind()

	for i := 1; i < objLen; i++ ***REMOVED***
		prevValue := value
		prevValueInterface := valueInterface

		value = objValue.Index(i)
		valueInterface = value.Interface()

		compareResult, isComparable := compare(prevValueInterface, valueInterface, firstValueKind)

		if !isComparable ***REMOVED***
			return Fail(t, fmt.Sprintf("Can not compare type \"%s\" and \"%s\"", reflect.TypeOf(value), reflect.TypeOf(prevValue)), msgAndArgs...)
		***REMOVED***

		if !containsValue(allowedComparesResults, compareResult) ***REMOVED***
			return Fail(t, fmt.Sprintf(failMessage, prevValue, value), msgAndArgs...)
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

// IsIncreasing asserts that the collection is increasing
//
//    assert.IsIncreasing(t, []int***REMOVED***1, 2, 3***REMOVED***)
//    assert.IsIncreasing(t, []float***REMOVED***1, 2***REMOVED***)
//    assert.IsIncreasing(t, []string***REMOVED***"a", "b"***REMOVED***)
func IsIncreasing(t TestingT, object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return isOrdered(t, object, []CompareType***REMOVED***compareLess***REMOVED***, "\"%v\" is not less than \"%v\"", msgAndArgs...)
***REMOVED***

// IsNonIncreasing asserts that the collection is not increasing
//
//    assert.IsNonIncreasing(t, []int***REMOVED***2, 1, 1***REMOVED***)
//    assert.IsNonIncreasing(t, []float***REMOVED***2, 1***REMOVED***)
//    assert.IsNonIncreasing(t, []string***REMOVED***"b", "a"***REMOVED***)
func IsNonIncreasing(t TestingT, object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return isOrdered(t, object, []CompareType***REMOVED***compareEqual, compareGreater***REMOVED***, "\"%v\" is not greater than or equal to \"%v\"", msgAndArgs...)
***REMOVED***

// IsDecreasing asserts that the collection is decreasing
//
//    assert.IsDecreasing(t, []int***REMOVED***2, 1, 0***REMOVED***)
//    assert.IsDecreasing(t, []float***REMOVED***2, 1***REMOVED***)
//    assert.IsDecreasing(t, []string***REMOVED***"b", "a"***REMOVED***)
func IsDecreasing(t TestingT, object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return isOrdered(t, object, []CompareType***REMOVED***compareGreater***REMOVED***, "\"%v\" is not greater than \"%v\"", msgAndArgs...)
***REMOVED***

// IsNonDecreasing asserts that the collection is not decreasing
//
//    assert.IsNonDecreasing(t, []int***REMOVED***1, 1, 2***REMOVED***)
//    assert.IsNonDecreasing(t, []float***REMOVED***1, 2***REMOVED***)
//    assert.IsNonDecreasing(t, []string***REMOVED***"a", "b"***REMOVED***)
func IsNonDecreasing(t TestingT, object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return isOrdered(t, object, []CompareType***REMOVED***compareLess, compareEqual***REMOVED***, "\"%v\" is not less than or equal to \"%v\"", msgAndArgs...)
***REMOVED***
