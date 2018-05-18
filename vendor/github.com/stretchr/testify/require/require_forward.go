/*
* CODE GENERATED AUTOMATICALLY WITH github.com/stretchr/testify/_codegen
* THIS FILE MUST NOT BE EDITED BY HAND
 */

package require

import (
	assert "github.com/stretchr/testify/assert"
	http "net/http"
	url "net/url"
	time "time"
)

// Condition uses a Comparison to assert a complex condition.
func (a *Assertions) Condition(comp assert.Comparison, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	Condition(a.t, comp, msgAndArgs...)
***REMOVED***

// Conditionf uses a Comparison to assert a complex condition.
func (a *Assertions) Conditionf(comp assert.Comparison, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	Conditionf(a.t, comp, msg, args...)
***REMOVED***

// Contains asserts that the specified string, list(array, slice...) or map contains the
// specified substring or element.
//
//    a.Contains("Hello World", "World")
//    a.Contains(["Hello", "World"], "World")
//    a.Contains(***REMOVED***"Hello": "World"***REMOVED***, "Hello")
func (a *Assertions) Contains(s interface***REMOVED******REMOVED***, contains interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	Contains(a.t, s, contains, msgAndArgs...)
***REMOVED***

// Containsf asserts that the specified string, list(array, slice...) or map contains the
// specified substring or element.
//
//    a.Containsf("Hello World", "World", "error message %s", "formatted")
//    a.Containsf(["Hello", "World"], "World", "error message %s", "formatted")
//    a.Containsf(***REMOVED***"Hello": "World"***REMOVED***, "Hello", "error message %s", "formatted")
func (a *Assertions) Containsf(s interface***REMOVED******REMOVED***, contains interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	Containsf(a.t, s, contains, msg, args...)
***REMOVED***

// DirExists checks whether a directory exists in the given path. It also fails if the path is a file rather a directory or there is an error checking whether it exists.
func (a *Assertions) DirExists(path string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	DirExists(a.t, path, msgAndArgs...)
***REMOVED***

// DirExistsf checks whether a directory exists in the given path. It also fails if the path is a file rather a directory or there is an error checking whether it exists.
func (a *Assertions) DirExistsf(path string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	DirExistsf(a.t, path, msg, args...)
***REMOVED***

// ElementsMatch asserts that the specified listA(array, slice...) is equal to specified
// listB(array, slice...) ignoring the order of the elements. If there are duplicate elements,
// the number of appearances of each of them in both lists should match.
//
// a.ElementsMatch([1, 3, 2, 3], [1, 3, 3, 2])
func (a *Assertions) ElementsMatch(listA interface***REMOVED******REMOVED***, listB interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	ElementsMatch(a.t, listA, listB, msgAndArgs...)
***REMOVED***

// ElementsMatchf asserts that the specified listA(array, slice...) is equal to specified
// listB(array, slice...) ignoring the order of the elements. If there are duplicate elements,
// the number of appearances of each of them in both lists should match.
//
// a.ElementsMatchf([1, 3, 2, 3], [1, 3, 3, 2], "error message %s", "formatted")
func (a *Assertions) ElementsMatchf(listA interface***REMOVED******REMOVED***, listB interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	ElementsMatchf(a.t, listA, listB, msg, args...)
***REMOVED***

// Empty asserts that the specified object is empty.  I.e. nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  a.Empty(obj)
func (a *Assertions) Empty(object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	Empty(a.t, object, msgAndArgs...)
***REMOVED***

// Emptyf asserts that the specified object is empty.  I.e. nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  a.Emptyf(obj, "error message %s", "formatted")
func (a *Assertions) Emptyf(object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	Emptyf(a.t, object, msg, args...)
***REMOVED***

// Equal asserts that two objects are equal.
//
//    a.Equal(123, 123)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses). Function equality
// cannot be determined and will always fail.
func (a *Assertions) Equal(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	Equal(a.t, expected, actual, msgAndArgs...)
***REMOVED***

// EqualError asserts that a function returned an error (i.e. not `nil`)
// and that it is equal to the provided error.
//
//   actualObj, err := SomeFunction()
//   a.EqualError(err,  expectedErrorString)
func (a *Assertions) EqualError(theError error, errString string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	EqualError(a.t, theError, errString, msgAndArgs...)
***REMOVED***

// EqualErrorf asserts that a function returned an error (i.e. not `nil`)
// and that it is equal to the provided error.
//
//   actualObj, err := SomeFunction()
//   a.EqualErrorf(err,  expectedErrorString, "error message %s", "formatted")
func (a *Assertions) EqualErrorf(theError error, errString string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	EqualErrorf(a.t, theError, errString, msg, args...)
***REMOVED***

// EqualValues asserts that two objects are equal or convertable to the same types
// and equal.
//
//    a.EqualValues(uint32(123), int32(123))
func (a *Assertions) EqualValues(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	EqualValues(a.t, expected, actual, msgAndArgs...)
***REMOVED***

// EqualValuesf asserts that two objects are equal or convertable to the same types
// and equal.
//
//    a.EqualValuesf(uint32(123, "error message %s", "formatted"), int32(123))
func (a *Assertions) EqualValuesf(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	EqualValuesf(a.t, expected, actual, msg, args...)
***REMOVED***

// Equalf asserts that two objects are equal.
//
//    a.Equalf(123, 123, "error message %s", "formatted")
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses). Function equality
// cannot be determined and will always fail.
func (a *Assertions) Equalf(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	Equalf(a.t, expected, actual, msg, args...)
***REMOVED***

// Error asserts that a function returned an error (i.e. not `nil`).
//
//   actualObj, err := SomeFunction()
//   if a.Error(err) ***REMOVED***
// 	   assert.Equal(t, expectedError, err)
//   ***REMOVED***
func (a *Assertions) Error(err error, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	Error(a.t, err, msgAndArgs...)
***REMOVED***

// Errorf asserts that a function returned an error (i.e. not `nil`).
//
//   actualObj, err := SomeFunction()
//   if a.Errorf(err, "error message %s", "formatted") ***REMOVED***
// 	   assert.Equal(t, expectedErrorf, err)
//   ***REMOVED***
func (a *Assertions) Errorf(err error, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	Errorf(a.t, err, msg, args...)
***REMOVED***

// Exactly asserts that two objects are equal in value and type.
//
//    a.Exactly(int32(123), int64(123))
func (a *Assertions) Exactly(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	Exactly(a.t, expected, actual, msgAndArgs...)
***REMOVED***

// Exactlyf asserts that two objects are equal in value and type.
//
//    a.Exactlyf(int32(123, "error message %s", "formatted"), int64(123))
func (a *Assertions) Exactlyf(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	Exactlyf(a.t, expected, actual, msg, args...)
***REMOVED***

// Fail reports a failure through
func (a *Assertions) Fail(failureMessage string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	Fail(a.t, failureMessage, msgAndArgs...)
***REMOVED***

// FailNow fails test
func (a *Assertions) FailNow(failureMessage string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	FailNow(a.t, failureMessage, msgAndArgs...)
***REMOVED***

// FailNowf fails test
func (a *Assertions) FailNowf(failureMessage string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	FailNowf(a.t, failureMessage, msg, args...)
***REMOVED***

// Failf reports a failure through
func (a *Assertions) Failf(failureMessage string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	Failf(a.t, failureMessage, msg, args...)
***REMOVED***

// False asserts that the specified value is false.
//
//    a.False(myBool)
func (a *Assertions) False(value bool, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	False(a.t, value, msgAndArgs...)
***REMOVED***

// Falsef asserts that the specified value is false.
//
//    a.Falsef(myBool, "error message %s", "formatted")
func (a *Assertions) Falsef(value bool, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	Falsef(a.t, value, msg, args...)
***REMOVED***

// FileExists checks whether a file exists in the given path. It also fails if the path points to a directory or there is an error when trying to check the file.
func (a *Assertions) FileExists(path string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	FileExists(a.t, path, msgAndArgs...)
***REMOVED***

// FileExistsf checks whether a file exists in the given path. It also fails if the path points to a directory or there is an error when trying to check the file.
func (a *Assertions) FileExistsf(path string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	FileExistsf(a.t, path, msg, args...)
***REMOVED***

// HTTPBodyContains asserts that a specified handler returns a
// body that contains a string.
//
//  a.HTTPBodyContains(myHandler, "www.google.com", nil, "I'm Feeling Lucky")
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPBodyContains(handler http.HandlerFunc, method string, url string, values url.Values, str interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	HTTPBodyContains(a.t, handler, method, url, values, str, msgAndArgs...)
***REMOVED***

// HTTPBodyContainsf asserts that a specified handler returns a
// body that contains a string.
//
//  a.HTTPBodyContainsf(myHandler, "www.google.com", nil, "I'm Feeling Lucky", "error message %s", "formatted")
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPBodyContainsf(handler http.HandlerFunc, method string, url string, values url.Values, str interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	HTTPBodyContainsf(a.t, handler, method, url, values, str, msg, args...)
***REMOVED***

// HTTPBodyNotContains asserts that a specified handler returns a
// body that does not contain a string.
//
//  a.HTTPBodyNotContains(myHandler, "www.google.com", nil, "I'm Feeling Lucky")
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPBodyNotContains(handler http.HandlerFunc, method string, url string, values url.Values, str interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	HTTPBodyNotContains(a.t, handler, method, url, values, str, msgAndArgs...)
***REMOVED***

// HTTPBodyNotContainsf asserts that a specified handler returns a
// body that does not contain a string.
//
//  a.HTTPBodyNotContainsf(myHandler, "www.google.com", nil, "I'm Feeling Lucky", "error message %s", "formatted")
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPBodyNotContainsf(handler http.HandlerFunc, method string, url string, values url.Values, str interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	HTTPBodyNotContainsf(a.t, handler, method, url, values, str, msg, args...)
***REMOVED***

// HTTPError asserts that a specified handler returns an error status code.
//
//  a.HTTPError(myHandler, "POST", "/a/b/c", url.Values***REMOVED***"a": []string***REMOVED***"b", "c"***REMOVED******REMOVED***
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPError(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	HTTPError(a.t, handler, method, url, values, msgAndArgs...)
***REMOVED***

// HTTPErrorf asserts that a specified handler returns an error status code.
//
//  a.HTTPErrorf(myHandler, "POST", "/a/b/c", url.Values***REMOVED***"a": []string***REMOVED***"b", "c"***REMOVED******REMOVED***
//
// Returns whether the assertion was successful (true, "error message %s", "formatted") or not (false).
func (a *Assertions) HTTPErrorf(handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	HTTPErrorf(a.t, handler, method, url, values, msg, args...)
***REMOVED***

// HTTPRedirect asserts that a specified handler returns a redirect status code.
//
//  a.HTTPRedirect(myHandler, "GET", "/a/b/c", url.Values***REMOVED***"a": []string***REMOVED***"b", "c"***REMOVED******REMOVED***
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPRedirect(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	HTTPRedirect(a.t, handler, method, url, values, msgAndArgs...)
***REMOVED***

// HTTPRedirectf asserts that a specified handler returns a redirect status code.
//
//  a.HTTPRedirectf(myHandler, "GET", "/a/b/c", url.Values***REMOVED***"a": []string***REMOVED***"b", "c"***REMOVED******REMOVED***
//
// Returns whether the assertion was successful (true, "error message %s", "formatted") or not (false).
func (a *Assertions) HTTPRedirectf(handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	HTTPRedirectf(a.t, handler, method, url, values, msg, args...)
***REMOVED***

// HTTPSuccess asserts that a specified handler returns a success status code.
//
//  a.HTTPSuccess(myHandler, "POST", "http://www.google.com", nil)
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPSuccess(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	HTTPSuccess(a.t, handler, method, url, values, msgAndArgs...)
***REMOVED***

// HTTPSuccessf asserts that a specified handler returns a success status code.
//
//  a.HTTPSuccessf(myHandler, "POST", "http://www.google.com", nil, "error message %s", "formatted")
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPSuccessf(handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	HTTPSuccessf(a.t, handler, method, url, values, msg, args...)
***REMOVED***

// Implements asserts that an object is implemented by the specified interface.
//
//    a.Implements((*MyInterface)(nil), new(MyObject))
func (a *Assertions) Implements(interfaceObject interface***REMOVED******REMOVED***, object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	Implements(a.t, interfaceObject, object, msgAndArgs...)
***REMOVED***

// Implementsf asserts that an object is implemented by the specified interface.
//
//    a.Implementsf((*MyInterface, "error message %s", "formatted")(nil), new(MyObject))
func (a *Assertions) Implementsf(interfaceObject interface***REMOVED******REMOVED***, object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	Implementsf(a.t, interfaceObject, object, msg, args...)
***REMOVED***

// InDelta asserts that the two numerals are within delta of each other.
//
// 	 a.InDelta(math.Pi, (22 / 7.0), 0.01)
func (a *Assertions) InDelta(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, delta float64, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	InDelta(a.t, expected, actual, delta, msgAndArgs...)
***REMOVED***

// InDeltaMapValues is the same as InDelta, but it compares all values between two maps. Both maps must have exactly the same keys.
func (a *Assertions) InDeltaMapValues(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, delta float64, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	InDeltaMapValues(a.t, expected, actual, delta, msgAndArgs...)
***REMOVED***

// InDeltaMapValuesf is the same as InDelta, but it compares all values between two maps. Both maps must have exactly the same keys.
func (a *Assertions) InDeltaMapValuesf(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, delta float64, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	InDeltaMapValuesf(a.t, expected, actual, delta, msg, args...)
***REMOVED***

// InDeltaSlice is the same as InDelta, except it compares two slices.
func (a *Assertions) InDeltaSlice(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, delta float64, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	InDeltaSlice(a.t, expected, actual, delta, msgAndArgs...)
***REMOVED***

// InDeltaSlicef is the same as InDelta, except it compares two slices.
func (a *Assertions) InDeltaSlicef(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, delta float64, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	InDeltaSlicef(a.t, expected, actual, delta, msg, args...)
***REMOVED***

// InDeltaf asserts that the two numerals are within delta of each other.
//
// 	 a.InDeltaf(math.Pi, (22 / 7.0, "error message %s", "formatted"), 0.01)
func (a *Assertions) InDeltaf(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, delta float64, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	InDeltaf(a.t, expected, actual, delta, msg, args...)
***REMOVED***

// InEpsilon asserts that expected and actual have a relative error less than epsilon
func (a *Assertions) InEpsilon(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, epsilon float64, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	InEpsilon(a.t, expected, actual, epsilon, msgAndArgs...)
***REMOVED***

// InEpsilonSlice is the same as InEpsilon, except it compares each value from two slices.
func (a *Assertions) InEpsilonSlice(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, epsilon float64, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	InEpsilonSlice(a.t, expected, actual, epsilon, msgAndArgs...)
***REMOVED***

// InEpsilonSlicef is the same as InEpsilon, except it compares each value from two slices.
func (a *Assertions) InEpsilonSlicef(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, epsilon float64, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	InEpsilonSlicef(a.t, expected, actual, epsilon, msg, args...)
***REMOVED***

// InEpsilonf asserts that expected and actual have a relative error less than epsilon
func (a *Assertions) InEpsilonf(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, epsilon float64, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	InEpsilonf(a.t, expected, actual, epsilon, msg, args...)
***REMOVED***

// IsType asserts that the specified objects are of the same type.
func (a *Assertions) IsType(expectedType interface***REMOVED******REMOVED***, object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	IsType(a.t, expectedType, object, msgAndArgs...)
***REMOVED***

// IsTypef asserts that the specified objects are of the same type.
func (a *Assertions) IsTypef(expectedType interface***REMOVED******REMOVED***, object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	IsTypef(a.t, expectedType, object, msg, args...)
***REMOVED***

// JSONEq asserts that two JSON strings are equivalent.
//
//  a.JSONEq(`***REMOVED***"hello": "world", "foo": "bar"***REMOVED***`, `***REMOVED***"foo": "bar", "hello": "world"***REMOVED***`)
func (a *Assertions) JSONEq(expected string, actual string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	JSONEq(a.t, expected, actual, msgAndArgs...)
***REMOVED***

// JSONEqf asserts that two JSON strings are equivalent.
//
//  a.JSONEqf(`***REMOVED***"hello": "world", "foo": "bar"***REMOVED***`, `***REMOVED***"foo": "bar", "hello": "world"***REMOVED***`, "error message %s", "formatted")
func (a *Assertions) JSONEqf(expected string, actual string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	JSONEqf(a.t, expected, actual, msg, args...)
***REMOVED***

// Len asserts that the specified object has specific length.
// Len also fails if the object has a type that len() not accept.
//
//    a.Len(mySlice, 3)
func (a *Assertions) Len(object interface***REMOVED******REMOVED***, length int, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	Len(a.t, object, length, msgAndArgs...)
***REMOVED***

// Lenf asserts that the specified object has specific length.
// Lenf also fails if the object has a type that len() not accept.
//
//    a.Lenf(mySlice, 3, "error message %s", "formatted")
func (a *Assertions) Lenf(object interface***REMOVED******REMOVED***, length int, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	Lenf(a.t, object, length, msg, args...)
***REMOVED***

// Nil asserts that the specified object is nil.
//
//    a.Nil(err)
func (a *Assertions) Nil(object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	Nil(a.t, object, msgAndArgs...)
***REMOVED***

// Nilf asserts that the specified object is nil.
//
//    a.Nilf(err, "error message %s", "formatted")
func (a *Assertions) Nilf(object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	Nilf(a.t, object, msg, args...)
***REMOVED***

// NoError asserts that a function returned no error (i.e. `nil`).
//
//   actualObj, err := SomeFunction()
//   if a.NoError(err) ***REMOVED***
// 	   assert.Equal(t, expectedObj, actualObj)
//   ***REMOVED***
func (a *Assertions) NoError(err error, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	NoError(a.t, err, msgAndArgs...)
***REMOVED***

// NoErrorf asserts that a function returned no error (i.e. `nil`).
//
//   actualObj, err := SomeFunction()
//   if a.NoErrorf(err, "error message %s", "formatted") ***REMOVED***
// 	   assert.Equal(t, expectedObj, actualObj)
//   ***REMOVED***
func (a *Assertions) NoErrorf(err error, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	NoErrorf(a.t, err, msg, args...)
***REMOVED***

// NotContains asserts that the specified string, list(array, slice...) or map does NOT contain the
// specified substring or element.
//
//    a.NotContains("Hello World", "Earth")
//    a.NotContains(["Hello", "World"], "Earth")
//    a.NotContains(***REMOVED***"Hello": "World"***REMOVED***, "Earth")
func (a *Assertions) NotContains(s interface***REMOVED******REMOVED***, contains interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	NotContains(a.t, s, contains, msgAndArgs...)
***REMOVED***

// NotContainsf asserts that the specified string, list(array, slice...) or map does NOT contain the
// specified substring or element.
//
//    a.NotContainsf("Hello World", "Earth", "error message %s", "formatted")
//    a.NotContainsf(["Hello", "World"], "Earth", "error message %s", "formatted")
//    a.NotContainsf(***REMOVED***"Hello": "World"***REMOVED***, "Earth", "error message %s", "formatted")
func (a *Assertions) NotContainsf(s interface***REMOVED******REMOVED***, contains interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	NotContainsf(a.t, s, contains, msg, args...)
***REMOVED***

// NotEmpty asserts that the specified object is NOT empty.  I.e. not nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  if a.NotEmpty(obj) ***REMOVED***
//    assert.Equal(t, "two", obj[1])
//  ***REMOVED***
func (a *Assertions) NotEmpty(object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	NotEmpty(a.t, object, msgAndArgs...)
***REMOVED***

// NotEmptyf asserts that the specified object is NOT empty.  I.e. not nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  if a.NotEmptyf(obj, "error message %s", "formatted") ***REMOVED***
//    assert.Equal(t, "two", obj[1])
//  ***REMOVED***
func (a *Assertions) NotEmptyf(object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	NotEmptyf(a.t, object, msg, args...)
***REMOVED***

// NotEqual asserts that the specified values are NOT equal.
//
//    a.NotEqual(obj1, obj2)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
func (a *Assertions) NotEqual(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	NotEqual(a.t, expected, actual, msgAndArgs...)
***REMOVED***

// NotEqualf asserts that the specified values are NOT equal.
//
//    a.NotEqualf(obj1, obj2, "error message %s", "formatted")
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
func (a *Assertions) NotEqualf(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	NotEqualf(a.t, expected, actual, msg, args...)
***REMOVED***

// NotNil asserts that the specified object is not nil.
//
//    a.NotNil(err)
func (a *Assertions) NotNil(object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	NotNil(a.t, object, msgAndArgs...)
***REMOVED***

// NotNilf asserts that the specified object is not nil.
//
//    a.NotNilf(err, "error message %s", "formatted")
func (a *Assertions) NotNilf(object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	NotNilf(a.t, object, msg, args...)
***REMOVED***

// NotPanics asserts that the code inside the specified PanicTestFunc does NOT panic.
//
//   a.NotPanics(func()***REMOVED*** RemainCalm() ***REMOVED***)
func (a *Assertions) NotPanics(f assert.PanicTestFunc, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	NotPanics(a.t, f, msgAndArgs...)
***REMOVED***

// NotPanicsf asserts that the code inside the specified PanicTestFunc does NOT panic.
//
//   a.NotPanicsf(func()***REMOVED*** RemainCalm() ***REMOVED***, "error message %s", "formatted")
func (a *Assertions) NotPanicsf(f assert.PanicTestFunc, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	NotPanicsf(a.t, f, msg, args...)
***REMOVED***

// NotRegexp asserts that a specified regexp does not match a string.
//
//  a.NotRegexp(regexp.MustCompile("starts"), "it's starting")
//  a.NotRegexp("^start", "it's not starting")
func (a *Assertions) NotRegexp(rx interface***REMOVED******REMOVED***, str interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	NotRegexp(a.t, rx, str, msgAndArgs...)
***REMOVED***

// NotRegexpf asserts that a specified regexp does not match a string.
//
//  a.NotRegexpf(regexp.MustCompile("starts", "error message %s", "formatted"), "it's starting")
//  a.NotRegexpf("^start", "it's not starting", "error message %s", "formatted")
func (a *Assertions) NotRegexpf(rx interface***REMOVED******REMOVED***, str interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	NotRegexpf(a.t, rx, str, msg, args...)
***REMOVED***

// NotSubset asserts that the specified list(array, slice...) contains not all
// elements given in the specified subset(array, slice...).
//
//    a.NotSubset([1, 3, 4], [1, 2], "But [1, 3, 4] does not contain [1, 2]")
func (a *Assertions) NotSubset(list interface***REMOVED******REMOVED***, subset interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	NotSubset(a.t, list, subset, msgAndArgs...)
***REMOVED***

// NotSubsetf asserts that the specified list(array, slice...) contains not all
// elements given in the specified subset(array, slice...).
//
//    a.NotSubsetf([1, 3, 4], [1, 2], "But [1, 3, 4] does not contain [1, 2]", "error message %s", "formatted")
func (a *Assertions) NotSubsetf(list interface***REMOVED******REMOVED***, subset interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	NotSubsetf(a.t, list, subset, msg, args...)
***REMOVED***

// NotZero asserts that i is not the zero value for its type.
func (a *Assertions) NotZero(i interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	NotZero(a.t, i, msgAndArgs...)
***REMOVED***

// NotZerof asserts that i is not the zero value for its type.
func (a *Assertions) NotZerof(i interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	NotZerof(a.t, i, msg, args...)
***REMOVED***

// Panics asserts that the code inside the specified PanicTestFunc panics.
//
//   a.Panics(func()***REMOVED*** GoCrazy() ***REMOVED***)
func (a *Assertions) Panics(f assert.PanicTestFunc, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	Panics(a.t, f, msgAndArgs...)
***REMOVED***

// PanicsWithValue asserts that the code inside the specified PanicTestFunc panics, and that
// the recovered panic value equals the expected panic value.
//
//   a.PanicsWithValue("crazy error", func()***REMOVED*** GoCrazy() ***REMOVED***)
func (a *Assertions) PanicsWithValue(expected interface***REMOVED******REMOVED***, f assert.PanicTestFunc, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	PanicsWithValue(a.t, expected, f, msgAndArgs...)
***REMOVED***

// PanicsWithValuef asserts that the code inside the specified PanicTestFunc panics, and that
// the recovered panic value equals the expected panic value.
//
//   a.PanicsWithValuef("crazy error", func()***REMOVED*** GoCrazy() ***REMOVED***, "error message %s", "formatted")
func (a *Assertions) PanicsWithValuef(expected interface***REMOVED******REMOVED***, f assert.PanicTestFunc, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	PanicsWithValuef(a.t, expected, f, msg, args...)
***REMOVED***

// Panicsf asserts that the code inside the specified PanicTestFunc panics.
//
//   a.Panicsf(func()***REMOVED*** GoCrazy() ***REMOVED***, "error message %s", "formatted")
func (a *Assertions) Panicsf(f assert.PanicTestFunc, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	Panicsf(a.t, f, msg, args...)
***REMOVED***

// Regexp asserts that a specified regexp matches a string.
//
//  a.Regexp(regexp.MustCompile("start"), "it's starting")
//  a.Regexp("start...$", "it's not starting")
func (a *Assertions) Regexp(rx interface***REMOVED******REMOVED***, str interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	Regexp(a.t, rx, str, msgAndArgs...)
***REMOVED***

// Regexpf asserts that a specified regexp matches a string.
//
//  a.Regexpf(regexp.MustCompile("start", "error message %s", "formatted"), "it's starting")
//  a.Regexpf("start...$", "it's not starting", "error message %s", "formatted")
func (a *Assertions) Regexpf(rx interface***REMOVED******REMOVED***, str interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	Regexpf(a.t, rx, str, msg, args...)
***REMOVED***

// Subset asserts that the specified list(array, slice...) contains all
// elements given in the specified subset(array, slice...).
//
//    a.Subset([1, 2, 3], [1, 2], "But [1, 2, 3] does contain [1, 2]")
func (a *Assertions) Subset(list interface***REMOVED******REMOVED***, subset interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	Subset(a.t, list, subset, msgAndArgs...)
***REMOVED***

// Subsetf asserts that the specified list(array, slice...) contains all
// elements given in the specified subset(array, slice...).
//
//    a.Subsetf([1, 2, 3], [1, 2], "But [1, 2, 3] does contain [1, 2]", "error message %s", "formatted")
func (a *Assertions) Subsetf(list interface***REMOVED******REMOVED***, subset interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	Subsetf(a.t, list, subset, msg, args...)
***REMOVED***

// True asserts that the specified value is true.
//
//    a.True(myBool)
func (a *Assertions) True(value bool, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	True(a.t, value, msgAndArgs...)
***REMOVED***

// Truef asserts that the specified value is true.
//
//    a.Truef(myBool, "error message %s", "formatted")
func (a *Assertions) Truef(value bool, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	Truef(a.t, value, msg, args...)
***REMOVED***

// WithinDuration asserts that the two times are within duration delta of each other.
//
//   a.WithinDuration(time.Now(), time.Now(), 10*time.Second)
func (a *Assertions) WithinDuration(expected time.Time, actual time.Time, delta time.Duration, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	WithinDuration(a.t, expected, actual, delta, msgAndArgs...)
***REMOVED***

// WithinDurationf asserts that the two times are within duration delta of each other.
//
//   a.WithinDurationf(time.Now(), time.Now(), 10*time.Second, "error message %s", "formatted")
func (a *Assertions) WithinDurationf(expected time.Time, actual time.Time, delta time.Duration, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	WithinDurationf(a.t, expected, actual, delta, msg, args...)
***REMOVED***

// Zero asserts that i is the zero value for its type.
func (a *Assertions) Zero(i interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	Zero(a.t, i, msgAndArgs...)
***REMOVED***

// Zerof asserts that i is the zero value for its type.
func (a *Assertions) Zerof(i interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	Zerof(a.t, i, msg, args...)
***REMOVED***