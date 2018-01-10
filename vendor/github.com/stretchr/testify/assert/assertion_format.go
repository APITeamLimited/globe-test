/*
* CODE GENERATED AUTOMATICALLY WITH github.com/stretchr/testify/_codegen
* THIS FILE MUST NOT BE EDITED BY HAND
 */

package assert

import (
	http "net/http"
	url "net/url"
	time "time"
)

// Conditionf uses a Comparison to assert a complex condition.
func Conditionf(t TestingT, comp Comparison, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return Condition(t, comp, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// Containsf asserts that the specified string, list(array, slice...) or map contains the
// specified substring or element.
//
//    assert.Containsf(t, "Hello World", "World", "error message %s", "formatted")
//    assert.Containsf(t, ["Hello", "World"], "World", "error message %s", "formatted")
//    assert.Containsf(t, ***REMOVED***"Hello": "World"***REMOVED***, "Hello", "error message %s", "formatted")
func Containsf(t TestingT, s interface***REMOVED******REMOVED***, contains interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return Contains(t, s, contains, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// DirExistsf checks whether a directory exists in the given path. It also fails if the path is a file rather a directory or there is an error checking whether it exists.
func DirExistsf(t TestingT, path string, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return DirExists(t, path, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// ElementsMatchf asserts that the specified listA(array, slice...) is equal to specified
// listB(array, slice...) ignoring the order of the elements. If there are duplicate elements,
// the number of appearances of each of them in both lists should match.
//
// assert.ElementsMatchf(t, [1, 3, 2, 3], [1, 3, 3, 2], "error message %s", "formatted"))
func ElementsMatchf(t TestingT, listA interface***REMOVED******REMOVED***, listB interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return ElementsMatch(t, listA, listB, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// Emptyf asserts that the specified object is empty.  I.e. nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  assert.Emptyf(t, obj, "error message %s", "formatted")
func Emptyf(t TestingT, object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return Empty(t, object, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// Equalf asserts that two objects are equal.
//
//    assert.Equalf(t, 123, 123, "error message %s", "formatted")
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses). Function equality
// cannot be determined and will always fail.
func Equalf(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return Equal(t, expected, actual, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// EqualErrorf asserts that a function returned an error (i.e. not `nil`)
// and that it is equal to the provided error.
//
//   actualObj, err := SomeFunction()
//   assert.EqualErrorf(t, err,  expectedErrorString, "error message %s", "formatted")
func EqualErrorf(t TestingT, theError error, errString string, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return EqualError(t, theError, errString, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// EqualValuesf asserts that two objects are equal or convertable to the same types
// and equal.
//
//    assert.EqualValuesf(t, uint32(123, "error message %s", "formatted"), int32(123))
func EqualValuesf(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return EqualValues(t, expected, actual, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// Errorf asserts that a function returned an error (i.e. not `nil`).
//
//   actualObj, err := SomeFunction()
//   if assert.Errorf(t, err, "error message %s", "formatted") ***REMOVED***
// 	   assert.Equal(t, expectedErrorf, err)
//   ***REMOVED***
func Errorf(t TestingT, err error, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return Error(t, err, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// Exactlyf asserts that two objects are equal in value and type.
//
//    assert.Exactlyf(t, int32(123, "error message %s", "formatted"), int64(123))
func Exactlyf(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return Exactly(t, expected, actual, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// Failf reports a failure through
func Failf(t TestingT, failureMessage string, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return Fail(t, failureMessage, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// FailNowf fails test
func FailNowf(t TestingT, failureMessage string, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return FailNow(t, failureMessage, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// Falsef asserts that the specified value is false.
//
//    assert.Falsef(t, myBool, "error message %s", "formatted")
func Falsef(t TestingT, value bool, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return False(t, value, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// FileExistsf checks whether a file exists in the given path. It also fails if the path points to a directory or there is an error when trying to check the file.
func FileExistsf(t TestingT, path string, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return FileExists(t, path, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// HTTPBodyContainsf asserts that a specified handler returns a
// body that contains a string.
//
//  assert.HTTPBodyContainsf(t, myHandler, "www.google.com", nil, "I'm Feeling Lucky", "error message %s", "formatted")
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPBodyContainsf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, str interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return HTTPBodyContains(t, handler, method, url, values, str, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// HTTPBodyNotContainsf asserts that a specified handler returns a
// body that does not contain a string.
//
//  assert.HTTPBodyNotContainsf(t, myHandler, "www.google.com", nil, "I'm Feeling Lucky", "error message %s", "formatted")
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPBodyNotContainsf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, str interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return HTTPBodyNotContains(t, handler, method, url, values, str, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// HTTPErrorf asserts that a specified handler returns an error status code.
//
//  assert.HTTPErrorf(t, myHandler, "POST", "/a/b/c", url.Values***REMOVED***"a": []string***REMOVED***"b", "c"***REMOVED******REMOVED***
//
// Returns whether the assertion was successful (true, "error message %s", "formatted") or not (false).
func HTTPErrorf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return HTTPError(t, handler, method, url, values, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// HTTPRedirectf asserts that a specified handler returns a redirect status code.
//
//  assert.HTTPRedirectf(t, myHandler, "GET", "/a/b/c", url.Values***REMOVED***"a": []string***REMOVED***"b", "c"***REMOVED******REMOVED***
//
// Returns whether the assertion was successful (true, "error message %s", "formatted") or not (false).
func HTTPRedirectf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return HTTPRedirect(t, handler, method, url, values, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// HTTPSuccessf asserts that a specified handler returns a success status code.
//
//  assert.HTTPSuccessf(t, myHandler, "POST", "http://www.google.com", nil, "error message %s", "formatted")
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPSuccessf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return HTTPSuccess(t, handler, method, url, values, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// Implementsf asserts that an object is implemented by the specified interface.
//
//    assert.Implementsf(t, (*MyInterface, "error message %s", "formatted")(nil), new(MyObject))
func Implementsf(t TestingT, interfaceObject interface***REMOVED******REMOVED***, object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return Implements(t, interfaceObject, object, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// InDeltaf asserts that the two numerals are within delta of each other.
//
// 	 assert.InDeltaf(t, math.Pi, (22 / 7.0, "error message %s", "formatted"), 0.01)
func InDeltaf(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, delta float64, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return InDelta(t, expected, actual, delta, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// InDeltaMapValuesf is the same as InDelta, but it compares all values between two maps. Both maps must have exactly the same keys.
func InDeltaMapValuesf(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, delta float64, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return InDeltaMapValues(t, expected, actual, delta, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// InDeltaSlicef is the same as InDelta, except it compares two slices.
func InDeltaSlicef(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, delta float64, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return InDeltaSlice(t, expected, actual, delta, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// InEpsilonf asserts that expected and actual have a relative error less than epsilon
func InEpsilonf(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, epsilon float64, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return InEpsilon(t, expected, actual, epsilon, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// InEpsilonSlicef is the same as InEpsilon, except it compares each value from two slices.
func InEpsilonSlicef(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, epsilon float64, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return InEpsilonSlice(t, expected, actual, epsilon, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// IsTypef asserts that the specified objects are of the same type.
func IsTypef(t TestingT, expectedType interface***REMOVED******REMOVED***, object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return IsType(t, expectedType, object, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// JSONEqf asserts that two JSON strings are equivalent.
//
//  assert.JSONEqf(t, `***REMOVED***"hello": "world", "foo": "bar"***REMOVED***`, `***REMOVED***"foo": "bar", "hello": "world"***REMOVED***`, "error message %s", "formatted")
func JSONEqf(t TestingT, expected string, actual string, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return JSONEq(t, expected, actual, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// Lenf asserts that the specified object has specific length.
// Lenf also fails if the object has a type that len() not accept.
//
//    assert.Lenf(t, mySlice, 3, "error message %s", "formatted")
func Lenf(t TestingT, object interface***REMOVED******REMOVED***, length int, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return Len(t, object, length, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// Nilf asserts that the specified object is nil.
//
//    assert.Nilf(t, err, "error message %s", "formatted")
func Nilf(t TestingT, object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return Nil(t, object, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// NoErrorf asserts that a function returned no error (i.e. `nil`).
//
//   actualObj, err := SomeFunction()
//   if assert.NoErrorf(t, err, "error message %s", "formatted") ***REMOVED***
// 	   assert.Equal(t, expectedObj, actualObj)
//   ***REMOVED***
func NoErrorf(t TestingT, err error, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return NoError(t, err, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// NotContainsf asserts that the specified string, list(array, slice...) or map does NOT contain the
// specified substring or element.
//
//    assert.NotContainsf(t, "Hello World", "Earth", "error message %s", "formatted")
//    assert.NotContainsf(t, ["Hello", "World"], "Earth", "error message %s", "formatted")
//    assert.NotContainsf(t, ***REMOVED***"Hello": "World"***REMOVED***, "Earth", "error message %s", "formatted")
func NotContainsf(t TestingT, s interface***REMOVED******REMOVED***, contains interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return NotContains(t, s, contains, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// NotEmptyf asserts that the specified object is NOT empty.  I.e. not nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  if assert.NotEmptyf(t, obj, "error message %s", "formatted") ***REMOVED***
//    assert.Equal(t, "two", obj[1])
//  ***REMOVED***
func NotEmptyf(t TestingT, object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return NotEmpty(t, object, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// NotEqualf asserts that the specified values are NOT equal.
//
//    assert.NotEqualf(t, obj1, obj2, "error message %s", "formatted")
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
func NotEqualf(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return NotEqual(t, expected, actual, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// NotNilf asserts that the specified object is not nil.
//
//    assert.NotNilf(t, err, "error message %s", "formatted")
func NotNilf(t TestingT, object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return NotNil(t, object, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// NotPanicsf asserts that the code inside the specified PanicTestFunc does NOT panic.
//
//   assert.NotPanicsf(t, func()***REMOVED*** RemainCalm() ***REMOVED***, "error message %s", "formatted")
func NotPanicsf(t TestingT, f PanicTestFunc, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return NotPanics(t, f, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// NotRegexpf asserts that a specified regexp does not match a string.
//
//  assert.NotRegexpf(t, regexp.MustCompile("starts", "error message %s", "formatted"), "it's starting")
//  assert.NotRegexpf(t, "^start", "it's not starting", "error message %s", "formatted")
func NotRegexpf(t TestingT, rx interface***REMOVED******REMOVED***, str interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return NotRegexp(t, rx, str, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// NotSubsetf asserts that the specified list(array, slice...) contains not all
// elements given in the specified subset(array, slice...).
//
//    assert.NotSubsetf(t, [1, 3, 4], [1, 2], "But [1, 3, 4] does not contain [1, 2]", "error message %s", "formatted")
func NotSubsetf(t TestingT, list interface***REMOVED******REMOVED***, subset interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return NotSubset(t, list, subset, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// NotZerof asserts that i is not the zero value for its type.
func NotZerof(t TestingT, i interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return NotZero(t, i, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// Panicsf asserts that the code inside the specified PanicTestFunc panics.
//
//   assert.Panicsf(t, func()***REMOVED*** GoCrazy() ***REMOVED***, "error message %s", "formatted")
func Panicsf(t TestingT, f PanicTestFunc, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return Panics(t, f, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// PanicsWithValuef asserts that the code inside the specified PanicTestFunc panics, and that
// the recovered panic value equals the expected panic value.
//
//   assert.PanicsWithValuef(t, "crazy error", func()***REMOVED*** GoCrazy() ***REMOVED***, "error message %s", "formatted")
func PanicsWithValuef(t TestingT, expected interface***REMOVED******REMOVED***, f PanicTestFunc, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return PanicsWithValue(t, expected, f, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// Regexpf asserts that a specified regexp matches a string.
//
//  assert.Regexpf(t, regexp.MustCompile("start", "error message %s", "formatted"), "it's starting")
//  assert.Regexpf(t, "start...$", "it's not starting", "error message %s", "formatted")
func Regexpf(t TestingT, rx interface***REMOVED******REMOVED***, str interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return Regexp(t, rx, str, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// Subsetf asserts that the specified list(array, slice...) contains all
// elements given in the specified subset(array, slice...).
//
//    assert.Subsetf(t, [1, 2, 3], [1, 2], "But [1, 2, 3] does contain [1, 2]", "error message %s", "formatted")
func Subsetf(t TestingT, list interface***REMOVED******REMOVED***, subset interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return Subset(t, list, subset, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// Truef asserts that the specified value is true.
//
//    assert.Truef(t, myBool, "error message %s", "formatted")
func Truef(t TestingT, value bool, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return True(t, value, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// WithinDurationf asserts that the two times are within duration delta of each other.
//
//   assert.WithinDurationf(t, time.Now(), time.Now(), 10*time.Second, "error message %s", "formatted")
func WithinDurationf(t TestingT, expected time.Time, actual time.Time, delta time.Duration, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return WithinDuration(t, expected, actual, delta, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***

// Zerof asserts that i is the zero value for its type.
func Zerof(t TestingT, i interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return Zero(t, i, append([]interface***REMOVED******REMOVED******REMOVED***msg***REMOVED***, args...)...)
***REMOVED***
