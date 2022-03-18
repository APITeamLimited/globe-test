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
func Condition(t TestingT, comp assert.Comparison, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Condition(t, comp, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Conditionf uses a Comparison to assert a complex condition.
func Conditionf(t TestingT, comp assert.Comparison, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Conditionf(t, comp, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Contains asserts that the specified string, list(array, slice...) or map contains the
// specified substring or element.
//
//    assert.Contains(t, "Hello World", "World")
//    assert.Contains(t, ["Hello", "World"], "World")
//    assert.Contains(t, ***REMOVED***"Hello": "World"***REMOVED***, "Hello")
func Contains(t TestingT, s interface***REMOVED******REMOVED***, contains interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Contains(t, s, contains, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Containsf asserts that the specified string, list(array, slice...) or map contains the
// specified substring or element.
//
//    assert.Containsf(t, "Hello World", "World", "error message %s", "formatted")
//    assert.Containsf(t, ["Hello", "World"], "World", "error message %s", "formatted")
//    assert.Containsf(t, ***REMOVED***"Hello": "World"***REMOVED***, "Hello", "error message %s", "formatted")
func Containsf(t TestingT, s interface***REMOVED******REMOVED***, contains interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Containsf(t, s, contains, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// DirExists checks whether a directory exists in the given path. It also fails
// if the path is a file rather a directory or there is an error checking whether it exists.
func DirExists(t TestingT, path string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.DirExists(t, path, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// DirExistsf checks whether a directory exists in the given path. It also fails
// if the path is a file rather a directory or there is an error checking whether it exists.
func DirExistsf(t TestingT, path string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.DirExistsf(t, path, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// ElementsMatch asserts that the specified listA(array, slice...) is equal to specified
// listB(array, slice...) ignoring the order of the elements. If there are duplicate elements,
// the number of appearances of each of them in both lists should match.
//
// assert.ElementsMatch(t, [1, 3, 2, 3], [1, 3, 3, 2])
func ElementsMatch(t TestingT, listA interface***REMOVED******REMOVED***, listB interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.ElementsMatch(t, listA, listB, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// ElementsMatchf asserts that the specified listA(array, slice...) is equal to specified
// listB(array, slice...) ignoring the order of the elements. If there are duplicate elements,
// the number of appearances of each of them in both lists should match.
//
// assert.ElementsMatchf(t, [1, 3, 2, 3], [1, 3, 3, 2], "error message %s", "formatted")
func ElementsMatchf(t TestingT, listA interface***REMOVED******REMOVED***, listB interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.ElementsMatchf(t, listA, listB, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Empty asserts that the specified object is empty.  I.e. nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  assert.Empty(t, obj)
func Empty(t TestingT, object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Empty(t, object, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Emptyf asserts that the specified object is empty.  I.e. nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  assert.Emptyf(t, obj, "error message %s", "formatted")
func Emptyf(t TestingT, object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Emptyf(t, object, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Equal asserts that two objects are equal.
//
//    assert.Equal(t, 123, 123)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses). Function equality
// cannot be determined and will always fail.
func Equal(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Equal(t, expected, actual, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// EqualError asserts that a function returned an error (i.e. not `nil`)
// and that it is equal to the provided error.
//
//   actualObj, err := SomeFunction()
//   assert.EqualError(t, err,  expectedErrorString)
func EqualError(t TestingT, theError error, errString string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.EqualError(t, theError, errString, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// EqualErrorf asserts that a function returned an error (i.e. not `nil`)
// and that it is equal to the provided error.
//
//   actualObj, err := SomeFunction()
//   assert.EqualErrorf(t, err,  expectedErrorString, "error message %s", "formatted")
func EqualErrorf(t TestingT, theError error, errString string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.EqualErrorf(t, theError, errString, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// EqualValues asserts that two objects are equal or convertable to the same types
// and equal.
//
//    assert.EqualValues(t, uint32(123), int32(123))
func EqualValues(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.EqualValues(t, expected, actual, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// EqualValuesf asserts that two objects are equal or convertable to the same types
// and equal.
//
//    assert.EqualValuesf(t, uint32(123), int32(123), "error message %s", "formatted")
func EqualValuesf(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.EqualValuesf(t, expected, actual, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Equalf asserts that two objects are equal.
//
//    assert.Equalf(t, 123, 123, "error message %s", "formatted")
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses). Function equality
// cannot be determined and will always fail.
func Equalf(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Equalf(t, expected, actual, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Error asserts that a function returned an error (i.e. not `nil`).
//
//   actualObj, err := SomeFunction()
//   if assert.Error(t, err) ***REMOVED***
// 	   assert.Equal(t, expectedError, err)
//   ***REMOVED***
func Error(t TestingT, err error, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Error(t, err, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// ErrorAs asserts that at least one of the errors in err's chain matches target, and if so, sets target to that error value.
// This is a wrapper for errors.As.
func ErrorAs(t TestingT, err error, target interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.ErrorAs(t, err, target, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// ErrorAsf asserts that at least one of the errors in err's chain matches target, and if so, sets target to that error value.
// This is a wrapper for errors.As.
func ErrorAsf(t TestingT, err error, target interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.ErrorAsf(t, err, target, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// ErrorContains asserts that a function returned an error (i.e. not `nil`)
// and that the error contains the specified substring.
//
//   actualObj, err := SomeFunction()
//   assert.ErrorContains(t, err,  expectedErrorSubString)
func ErrorContains(t TestingT, theError error, contains string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.ErrorContains(t, theError, contains, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// ErrorContainsf asserts that a function returned an error (i.e. not `nil`)
// and that the error contains the specified substring.
//
//   actualObj, err := SomeFunction()
//   assert.ErrorContainsf(t, err,  expectedErrorSubString, "error message %s", "formatted")
func ErrorContainsf(t TestingT, theError error, contains string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.ErrorContainsf(t, theError, contains, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// ErrorIs asserts that at least one of the errors in err's chain matches target.
// This is a wrapper for errors.Is.
func ErrorIs(t TestingT, err error, target error, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.ErrorIs(t, err, target, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// ErrorIsf asserts that at least one of the errors in err's chain matches target.
// This is a wrapper for errors.Is.
func ErrorIsf(t TestingT, err error, target error, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.ErrorIsf(t, err, target, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Errorf asserts that a function returned an error (i.e. not `nil`).
//
//   actualObj, err := SomeFunction()
//   if assert.Errorf(t, err, "error message %s", "formatted") ***REMOVED***
// 	   assert.Equal(t, expectedErrorf, err)
//   ***REMOVED***
func Errorf(t TestingT, err error, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Errorf(t, err, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Eventually asserts that given condition will be met in waitFor time,
// periodically checking target function each tick.
//
//    assert.Eventually(t, func() bool ***REMOVED*** return true; ***REMOVED***, time.Second, 10*time.Millisecond)
func Eventually(t TestingT, condition func() bool, waitFor time.Duration, tick time.Duration, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Eventually(t, condition, waitFor, tick, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Eventuallyf asserts that given condition will be met in waitFor time,
// periodically checking target function each tick.
//
//    assert.Eventuallyf(t, func() bool ***REMOVED*** return true; ***REMOVED***, time.Second, 10*time.Millisecond, "error message %s", "formatted")
func Eventuallyf(t TestingT, condition func() bool, waitFor time.Duration, tick time.Duration, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Eventuallyf(t, condition, waitFor, tick, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Exactly asserts that two objects are equal in value and type.
//
//    assert.Exactly(t, int32(123), int64(123))
func Exactly(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Exactly(t, expected, actual, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Exactlyf asserts that two objects are equal in value and type.
//
//    assert.Exactlyf(t, int32(123), int64(123), "error message %s", "formatted")
func Exactlyf(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Exactlyf(t, expected, actual, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Fail reports a failure through
func Fail(t TestingT, failureMessage string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Fail(t, failureMessage, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// FailNow fails test
func FailNow(t TestingT, failureMessage string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.FailNow(t, failureMessage, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// FailNowf fails test
func FailNowf(t TestingT, failureMessage string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.FailNowf(t, failureMessage, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Failf reports a failure through
func Failf(t TestingT, failureMessage string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Failf(t, failureMessage, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// False asserts that the specified value is false.
//
//    assert.False(t, myBool)
func False(t TestingT, value bool, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.False(t, value, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Falsef asserts that the specified value is false.
//
//    assert.Falsef(t, myBool, "error message %s", "formatted")
func Falsef(t TestingT, value bool, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Falsef(t, value, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// FileExists checks whether a file exists in the given path. It also fails if
// the path points to a directory or there is an error when trying to check the file.
func FileExists(t TestingT, path string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.FileExists(t, path, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// FileExistsf checks whether a file exists in the given path. It also fails if
// the path points to a directory or there is an error when trying to check the file.
func FileExistsf(t TestingT, path string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.FileExistsf(t, path, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Greater asserts that the first element is greater than the second
//
//    assert.Greater(t, 2, 1)
//    assert.Greater(t, float64(2), float64(1))
//    assert.Greater(t, "b", "a")
func Greater(t TestingT, e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Greater(t, e1, e2, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// GreaterOrEqual asserts that the first element is greater than or equal to the second
//
//    assert.GreaterOrEqual(t, 2, 1)
//    assert.GreaterOrEqual(t, 2, 2)
//    assert.GreaterOrEqual(t, "b", "a")
//    assert.GreaterOrEqual(t, "b", "b")
func GreaterOrEqual(t TestingT, e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.GreaterOrEqual(t, e1, e2, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// GreaterOrEqualf asserts that the first element is greater than or equal to the second
//
//    assert.GreaterOrEqualf(t, 2, 1, "error message %s", "formatted")
//    assert.GreaterOrEqualf(t, 2, 2, "error message %s", "formatted")
//    assert.GreaterOrEqualf(t, "b", "a", "error message %s", "formatted")
//    assert.GreaterOrEqualf(t, "b", "b", "error message %s", "formatted")
func GreaterOrEqualf(t TestingT, e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.GreaterOrEqualf(t, e1, e2, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Greaterf asserts that the first element is greater than the second
//
//    assert.Greaterf(t, 2, 1, "error message %s", "formatted")
//    assert.Greaterf(t, float64(2), float64(1), "error message %s", "formatted")
//    assert.Greaterf(t, "b", "a", "error message %s", "formatted")
func Greaterf(t TestingT, e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Greaterf(t, e1, e2, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// HTTPBodyContains asserts that a specified handler returns a
// body that contains a string.
//
//  assert.HTTPBodyContains(t, myHandler, "GET", "www.google.com", nil, "I'm Feeling Lucky")
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPBodyContains(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, str interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.HTTPBodyContains(t, handler, method, url, values, str, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// HTTPBodyContainsf asserts that a specified handler returns a
// body that contains a string.
//
//  assert.HTTPBodyContainsf(t, myHandler, "GET", "www.google.com", nil, "I'm Feeling Lucky", "error message %s", "formatted")
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPBodyContainsf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, str interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.HTTPBodyContainsf(t, handler, method, url, values, str, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// HTTPBodyNotContains asserts that a specified handler returns a
// body that does not contain a string.
//
//  assert.HTTPBodyNotContains(t, myHandler, "GET", "www.google.com", nil, "I'm Feeling Lucky")
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPBodyNotContains(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, str interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.HTTPBodyNotContains(t, handler, method, url, values, str, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// HTTPBodyNotContainsf asserts that a specified handler returns a
// body that does not contain a string.
//
//  assert.HTTPBodyNotContainsf(t, myHandler, "GET", "www.google.com", nil, "I'm Feeling Lucky", "error message %s", "formatted")
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPBodyNotContainsf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, str interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.HTTPBodyNotContainsf(t, handler, method, url, values, str, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// HTTPError asserts that a specified handler returns an error status code.
//
//  assert.HTTPError(t, myHandler, "POST", "/a/b/c", url.Values***REMOVED***"a": []string***REMOVED***"b", "c"***REMOVED******REMOVED***
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPError(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.HTTPError(t, handler, method, url, values, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// HTTPErrorf asserts that a specified handler returns an error status code.
//
//  assert.HTTPErrorf(t, myHandler, "POST", "/a/b/c", url.Values***REMOVED***"a": []string***REMOVED***"b", "c"***REMOVED******REMOVED***
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPErrorf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.HTTPErrorf(t, handler, method, url, values, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// HTTPRedirect asserts that a specified handler returns a redirect status code.
//
//  assert.HTTPRedirect(t, myHandler, "GET", "/a/b/c", url.Values***REMOVED***"a": []string***REMOVED***"b", "c"***REMOVED******REMOVED***
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPRedirect(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.HTTPRedirect(t, handler, method, url, values, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// HTTPRedirectf asserts that a specified handler returns a redirect status code.
//
//  assert.HTTPRedirectf(t, myHandler, "GET", "/a/b/c", url.Values***REMOVED***"a": []string***REMOVED***"b", "c"***REMOVED******REMOVED***
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPRedirectf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.HTTPRedirectf(t, handler, method, url, values, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// HTTPStatusCode asserts that a specified handler returns a specified status code.
//
//  assert.HTTPStatusCode(t, myHandler, "GET", "/notImplemented", nil, 501)
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPStatusCode(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, statuscode int, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.HTTPStatusCode(t, handler, method, url, values, statuscode, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// HTTPStatusCodef asserts that a specified handler returns a specified status code.
//
//  assert.HTTPStatusCodef(t, myHandler, "GET", "/notImplemented", nil, 501, "error message %s", "formatted")
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPStatusCodef(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, statuscode int, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.HTTPStatusCodef(t, handler, method, url, values, statuscode, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// HTTPSuccess asserts that a specified handler returns a success status code.
//
//  assert.HTTPSuccess(t, myHandler, "POST", "http://www.google.com", nil)
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPSuccess(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.HTTPSuccess(t, handler, method, url, values, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// HTTPSuccessf asserts that a specified handler returns a success status code.
//
//  assert.HTTPSuccessf(t, myHandler, "POST", "http://www.google.com", nil, "error message %s", "formatted")
//
// Returns whether the assertion was successful (true) or not (false).
func HTTPSuccessf(t TestingT, handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.HTTPSuccessf(t, handler, method, url, values, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Implements asserts that an object is implemented by the specified interface.
//
//    assert.Implements(t, (*MyInterface)(nil), new(MyObject))
func Implements(t TestingT, interfaceObject interface***REMOVED******REMOVED***, object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Implements(t, interfaceObject, object, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Implementsf asserts that an object is implemented by the specified interface.
//
//    assert.Implementsf(t, (*MyInterface)(nil), new(MyObject), "error message %s", "formatted")
func Implementsf(t TestingT, interfaceObject interface***REMOVED******REMOVED***, object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Implementsf(t, interfaceObject, object, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// InDelta asserts that the two numerals are within delta of each other.
//
// 	 assert.InDelta(t, math.Pi, 22/7.0, 0.01)
func InDelta(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, delta float64, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.InDelta(t, expected, actual, delta, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// InDeltaMapValues is the same as InDelta, but it compares all values between two maps. Both maps must have exactly the same keys.
func InDeltaMapValues(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, delta float64, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.InDeltaMapValues(t, expected, actual, delta, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// InDeltaMapValuesf is the same as InDelta, but it compares all values between two maps. Both maps must have exactly the same keys.
func InDeltaMapValuesf(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, delta float64, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.InDeltaMapValuesf(t, expected, actual, delta, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// InDeltaSlice is the same as InDelta, except it compares two slices.
func InDeltaSlice(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, delta float64, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.InDeltaSlice(t, expected, actual, delta, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// InDeltaSlicef is the same as InDelta, except it compares two slices.
func InDeltaSlicef(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, delta float64, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.InDeltaSlicef(t, expected, actual, delta, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// InDeltaf asserts that the two numerals are within delta of each other.
//
// 	 assert.InDeltaf(t, math.Pi, 22/7.0, 0.01, "error message %s", "formatted")
func InDeltaf(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, delta float64, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.InDeltaf(t, expected, actual, delta, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// InEpsilon asserts that expected and actual have a relative error less than epsilon
func InEpsilon(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, epsilon float64, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.InEpsilon(t, expected, actual, epsilon, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// InEpsilonSlice is the same as InEpsilon, except it compares each value from two slices.
func InEpsilonSlice(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, epsilon float64, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.InEpsilonSlice(t, expected, actual, epsilon, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// InEpsilonSlicef is the same as InEpsilon, except it compares each value from two slices.
func InEpsilonSlicef(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, epsilon float64, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.InEpsilonSlicef(t, expected, actual, epsilon, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// InEpsilonf asserts that expected and actual have a relative error less than epsilon
func InEpsilonf(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, epsilon float64, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.InEpsilonf(t, expected, actual, epsilon, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// IsDecreasing asserts that the collection is decreasing
//
//    assert.IsDecreasing(t, []int***REMOVED***2, 1, 0***REMOVED***)
//    assert.IsDecreasing(t, []float***REMOVED***2, 1***REMOVED***)
//    assert.IsDecreasing(t, []string***REMOVED***"b", "a"***REMOVED***)
func IsDecreasing(t TestingT, object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.IsDecreasing(t, object, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// IsDecreasingf asserts that the collection is decreasing
//
//    assert.IsDecreasingf(t, []int***REMOVED***2, 1, 0***REMOVED***, "error message %s", "formatted")
//    assert.IsDecreasingf(t, []float***REMOVED***2, 1***REMOVED***, "error message %s", "formatted")
//    assert.IsDecreasingf(t, []string***REMOVED***"b", "a"***REMOVED***, "error message %s", "formatted")
func IsDecreasingf(t TestingT, object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.IsDecreasingf(t, object, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// IsIncreasing asserts that the collection is increasing
//
//    assert.IsIncreasing(t, []int***REMOVED***1, 2, 3***REMOVED***)
//    assert.IsIncreasing(t, []float***REMOVED***1, 2***REMOVED***)
//    assert.IsIncreasing(t, []string***REMOVED***"a", "b"***REMOVED***)
func IsIncreasing(t TestingT, object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.IsIncreasing(t, object, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// IsIncreasingf asserts that the collection is increasing
//
//    assert.IsIncreasingf(t, []int***REMOVED***1, 2, 3***REMOVED***, "error message %s", "formatted")
//    assert.IsIncreasingf(t, []float***REMOVED***1, 2***REMOVED***, "error message %s", "formatted")
//    assert.IsIncreasingf(t, []string***REMOVED***"a", "b"***REMOVED***, "error message %s", "formatted")
func IsIncreasingf(t TestingT, object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.IsIncreasingf(t, object, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// IsNonDecreasing asserts that the collection is not decreasing
//
//    assert.IsNonDecreasing(t, []int***REMOVED***1, 1, 2***REMOVED***)
//    assert.IsNonDecreasing(t, []float***REMOVED***1, 2***REMOVED***)
//    assert.IsNonDecreasing(t, []string***REMOVED***"a", "b"***REMOVED***)
func IsNonDecreasing(t TestingT, object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.IsNonDecreasing(t, object, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// IsNonDecreasingf asserts that the collection is not decreasing
//
//    assert.IsNonDecreasingf(t, []int***REMOVED***1, 1, 2***REMOVED***, "error message %s", "formatted")
//    assert.IsNonDecreasingf(t, []float***REMOVED***1, 2***REMOVED***, "error message %s", "formatted")
//    assert.IsNonDecreasingf(t, []string***REMOVED***"a", "b"***REMOVED***, "error message %s", "formatted")
func IsNonDecreasingf(t TestingT, object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.IsNonDecreasingf(t, object, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// IsNonIncreasing asserts that the collection is not increasing
//
//    assert.IsNonIncreasing(t, []int***REMOVED***2, 1, 1***REMOVED***)
//    assert.IsNonIncreasing(t, []float***REMOVED***2, 1***REMOVED***)
//    assert.IsNonIncreasing(t, []string***REMOVED***"b", "a"***REMOVED***)
func IsNonIncreasing(t TestingT, object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.IsNonIncreasing(t, object, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// IsNonIncreasingf asserts that the collection is not increasing
//
//    assert.IsNonIncreasingf(t, []int***REMOVED***2, 1, 1***REMOVED***, "error message %s", "formatted")
//    assert.IsNonIncreasingf(t, []float***REMOVED***2, 1***REMOVED***, "error message %s", "formatted")
//    assert.IsNonIncreasingf(t, []string***REMOVED***"b", "a"***REMOVED***, "error message %s", "formatted")
func IsNonIncreasingf(t TestingT, object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.IsNonIncreasingf(t, object, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// IsType asserts that the specified objects are of the same type.
func IsType(t TestingT, expectedType interface***REMOVED******REMOVED***, object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.IsType(t, expectedType, object, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// IsTypef asserts that the specified objects are of the same type.
func IsTypef(t TestingT, expectedType interface***REMOVED******REMOVED***, object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.IsTypef(t, expectedType, object, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// JSONEq asserts that two JSON strings are equivalent.
//
//  assert.JSONEq(t, `***REMOVED***"hello": "world", "foo": "bar"***REMOVED***`, `***REMOVED***"foo": "bar", "hello": "world"***REMOVED***`)
func JSONEq(t TestingT, expected string, actual string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.JSONEq(t, expected, actual, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// JSONEqf asserts that two JSON strings are equivalent.
//
//  assert.JSONEqf(t, `***REMOVED***"hello": "world", "foo": "bar"***REMOVED***`, `***REMOVED***"foo": "bar", "hello": "world"***REMOVED***`, "error message %s", "formatted")
func JSONEqf(t TestingT, expected string, actual string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.JSONEqf(t, expected, actual, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Len asserts that the specified object has specific length.
// Len also fails if the object has a type that len() not accept.
//
//    assert.Len(t, mySlice, 3)
func Len(t TestingT, object interface***REMOVED******REMOVED***, length int, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Len(t, object, length, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Lenf asserts that the specified object has specific length.
// Lenf also fails if the object has a type that len() not accept.
//
//    assert.Lenf(t, mySlice, 3, "error message %s", "formatted")
func Lenf(t TestingT, object interface***REMOVED******REMOVED***, length int, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Lenf(t, object, length, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Less asserts that the first element is less than the second
//
//    assert.Less(t, 1, 2)
//    assert.Less(t, float64(1), float64(2))
//    assert.Less(t, "a", "b")
func Less(t TestingT, e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Less(t, e1, e2, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// LessOrEqual asserts that the first element is less than or equal to the second
//
//    assert.LessOrEqual(t, 1, 2)
//    assert.LessOrEqual(t, 2, 2)
//    assert.LessOrEqual(t, "a", "b")
//    assert.LessOrEqual(t, "b", "b")
func LessOrEqual(t TestingT, e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.LessOrEqual(t, e1, e2, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// LessOrEqualf asserts that the first element is less than or equal to the second
//
//    assert.LessOrEqualf(t, 1, 2, "error message %s", "formatted")
//    assert.LessOrEqualf(t, 2, 2, "error message %s", "formatted")
//    assert.LessOrEqualf(t, "a", "b", "error message %s", "formatted")
//    assert.LessOrEqualf(t, "b", "b", "error message %s", "formatted")
func LessOrEqualf(t TestingT, e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.LessOrEqualf(t, e1, e2, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Lessf asserts that the first element is less than the second
//
//    assert.Lessf(t, 1, 2, "error message %s", "formatted")
//    assert.Lessf(t, float64(1), float64(2), "error message %s", "formatted")
//    assert.Lessf(t, "a", "b", "error message %s", "formatted")
func Lessf(t TestingT, e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Lessf(t, e1, e2, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Negative asserts that the specified element is negative
//
//    assert.Negative(t, -1)
//    assert.Negative(t, -1.23)
func Negative(t TestingT, e interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Negative(t, e, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Negativef asserts that the specified element is negative
//
//    assert.Negativef(t, -1, "error message %s", "formatted")
//    assert.Negativef(t, -1.23, "error message %s", "formatted")
func Negativef(t TestingT, e interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Negativef(t, e, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Never asserts that the given condition doesn't satisfy in waitFor time,
// periodically checking the target function each tick.
//
//    assert.Never(t, func() bool ***REMOVED*** return false; ***REMOVED***, time.Second, 10*time.Millisecond)
func Never(t TestingT, condition func() bool, waitFor time.Duration, tick time.Duration, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Never(t, condition, waitFor, tick, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Neverf asserts that the given condition doesn't satisfy in waitFor time,
// periodically checking the target function each tick.
//
//    assert.Neverf(t, func() bool ***REMOVED*** return false; ***REMOVED***, time.Second, 10*time.Millisecond, "error message %s", "formatted")
func Neverf(t TestingT, condition func() bool, waitFor time.Duration, tick time.Duration, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Neverf(t, condition, waitFor, tick, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Nil asserts that the specified object is nil.
//
//    assert.Nil(t, err)
func Nil(t TestingT, object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Nil(t, object, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Nilf asserts that the specified object is nil.
//
//    assert.Nilf(t, err, "error message %s", "formatted")
func Nilf(t TestingT, object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Nilf(t, object, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NoDirExists checks whether a directory does not exist in the given path.
// It fails if the path points to an existing _directory_ only.
func NoDirExists(t TestingT, path string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NoDirExists(t, path, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NoDirExistsf checks whether a directory does not exist in the given path.
// It fails if the path points to an existing _directory_ only.
func NoDirExistsf(t TestingT, path string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NoDirExistsf(t, path, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NoError asserts that a function returned no error (i.e. `nil`).
//
//   actualObj, err := SomeFunction()
//   if assert.NoError(t, err) ***REMOVED***
// 	   assert.Equal(t, expectedObj, actualObj)
//   ***REMOVED***
func NoError(t TestingT, err error, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NoError(t, err, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NoErrorf asserts that a function returned no error (i.e. `nil`).
//
//   actualObj, err := SomeFunction()
//   if assert.NoErrorf(t, err, "error message %s", "formatted") ***REMOVED***
// 	   assert.Equal(t, expectedObj, actualObj)
//   ***REMOVED***
func NoErrorf(t TestingT, err error, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NoErrorf(t, err, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NoFileExists checks whether a file does not exist in a given path. It fails
// if the path points to an existing _file_ only.
func NoFileExists(t TestingT, path string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NoFileExists(t, path, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NoFileExistsf checks whether a file does not exist in a given path. It fails
// if the path points to an existing _file_ only.
func NoFileExistsf(t TestingT, path string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NoFileExistsf(t, path, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotContains asserts that the specified string, list(array, slice...) or map does NOT contain the
// specified substring or element.
//
//    assert.NotContains(t, "Hello World", "Earth")
//    assert.NotContains(t, ["Hello", "World"], "Earth")
//    assert.NotContains(t, ***REMOVED***"Hello": "World"***REMOVED***, "Earth")
func NotContains(t TestingT, s interface***REMOVED******REMOVED***, contains interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotContains(t, s, contains, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotContainsf asserts that the specified string, list(array, slice...) or map does NOT contain the
// specified substring or element.
//
//    assert.NotContainsf(t, "Hello World", "Earth", "error message %s", "formatted")
//    assert.NotContainsf(t, ["Hello", "World"], "Earth", "error message %s", "formatted")
//    assert.NotContainsf(t, ***REMOVED***"Hello": "World"***REMOVED***, "Earth", "error message %s", "formatted")
func NotContainsf(t TestingT, s interface***REMOVED******REMOVED***, contains interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotContainsf(t, s, contains, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotEmpty asserts that the specified object is NOT empty.  I.e. not nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  if assert.NotEmpty(t, obj) ***REMOVED***
//    assert.Equal(t, "two", obj[1])
//  ***REMOVED***
func NotEmpty(t TestingT, object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotEmpty(t, object, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotEmptyf asserts that the specified object is NOT empty.  I.e. not nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  if assert.NotEmptyf(t, obj, "error message %s", "formatted") ***REMOVED***
//    assert.Equal(t, "two", obj[1])
//  ***REMOVED***
func NotEmptyf(t TestingT, object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotEmptyf(t, object, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotEqual asserts that the specified values are NOT equal.
//
//    assert.NotEqual(t, obj1, obj2)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
func NotEqual(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotEqual(t, expected, actual, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotEqualValues asserts that two objects are not equal even when converted to the same type
//
//    assert.NotEqualValues(t, obj1, obj2)
func NotEqualValues(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotEqualValues(t, expected, actual, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotEqualValuesf asserts that two objects are not equal even when converted to the same type
//
//    assert.NotEqualValuesf(t, obj1, obj2, "error message %s", "formatted")
func NotEqualValuesf(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotEqualValuesf(t, expected, actual, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotEqualf asserts that the specified values are NOT equal.
//
//    assert.NotEqualf(t, obj1, obj2, "error message %s", "formatted")
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
func NotEqualf(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotEqualf(t, expected, actual, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotErrorIs asserts that at none of the errors in err's chain matches target.
// This is a wrapper for errors.Is.
func NotErrorIs(t TestingT, err error, target error, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotErrorIs(t, err, target, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotErrorIsf asserts that at none of the errors in err's chain matches target.
// This is a wrapper for errors.Is.
func NotErrorIsf(t TestingT, err error, target error, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotErrorIsf(t, err, target, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotNil asserts that the specified object is not nil.
//
//    assert.NotNil(t, err)
func NotNil(t TestingT, object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotNil(t, object, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotNilf asserts that the specified object is not nil.
//
//    assert.NotNilf(t, err, "error message %s", "formatted")
func NotNilf(t TestingT, object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotNilf(t, object, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotPanics asserts that the code inside the specified PanicTestFunc does NOT panic.
//
//   assert.NotPanics(t, func()***REMOVED*** RemainCalm() ***REMOVED***)
func NotPanics(t TestingT, f assert.PanicTestFunc, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotPanics(t, f, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotPanicsf asserts that the code inside the specified PanicTestFunc does NOT panic.
//
//   assert.NotPanicsf(t, func()***REMOVED*** RemainCalm() ***REMOVED***, "error message %s", "formatted")
func NotPanicsf(t TestingT, f assert.PanicTestFunc, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotPanicsf(t, f, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotRegexp asserts that a specified regexp does not match a string.
//
//  assert.NotRegexp(t, regexp.MustCompile("starts"), "it's starting")
//  assert.NotRegexp(t, "^start", "it's not starting")
func NotRegexp(t TestingT, rx interface***REMOVED******REMOVED***, str interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotRegexp(t, rx, str, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotRegexpf asserts that a specified regexp does not match a string.
//
//  assert.NotRegexpf(t, regexp.MustCompile("starts"), "it's starting", "error message %s", "formatted")
//  assert.NotRegexpf(t, "^start", "it's not starting", "error message %s", "formatted")
func NotRegexpf(t TestingT, rx interface***REMOVED******REMOVED***, str interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotRegexpf(t, rx, str, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotSame asserts that two pointers do not reference the same object.
//
//    assert.NotSame(t, ptr1, ptr2)
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func NotSame(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotSame(t, expected, actual, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotSamef asserts that two pointers do not reference the same object.
//
//    assert.NotSamef(t, ptr1, ptr2, "error message %s", "formatted")
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func NotSamef(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotSamef(t, expected, actual, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotSubset asserts that the specified list(array, slice...) contains not all
// elements given in the specified subset(array, slice...).
//
//    assert.NotSubset(t, [1, 3, 4], [1, 2], "But [1, 3, 4] does not contain [1, 2]")
func NotSubset(t TestingT, list interface***REMOVED******REMOVED***, subset interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotSubset(t, list, subset, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotSubsetf asserts that the specified list(array, slice...) contains not all
// elements given in the specified subset(array, slice...).
//
//    assert.NotSubsetf(t, [1, 3, 4], [1, 2], "But [1, 3, 4] does not contain [1, 2]", "error message %s", "formatted")
func NotSubsetf(t TestingT, list interface***REMOVED******REMOVED***, subset interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotSubsetf(t, list, subset, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotZero asserts that i is not the zero value for its type.
func NotZero(t TestingT, i interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotZero(t, i, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// NotZerof asserts that i is not the zero value for its type.
func NotZerof(t TestingT, i interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.NotZerof(t, i, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Panics asserts that the code inside the specified PanicTestFunc panics.
//
//   assert.Panics(t, func()***REMOVED*** GoCrazy() ***REMOVED***)
func Panics(t TestingT, f assert.PanicTestFunc, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Panics(t, f, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// PanicsWithError asserts that the code inside the specified PanicTestFunc
// panics, and that the recovered panic value is an error that satisfies the
// EqualError comparison.
//
//   assert.PanicsWithError(t, "crazy error", func()***REMOVED*** GoCrazy() ***REMOVED***)
func PanicsWithError(t TestingT, errString string, f assert.PanicTestFunc, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.PanicsWithError(t, errString, f, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// PanicsWithErrorf asserts that the code inside the specified PanicTestFunc
// panics, and that the recovered panic value is an error that satisfies the
// EqualError comparison.
//
//   assert.PanicsWithErrorf(t, "crazy error", func()***REMOVED*** GoCrazy() ***REMOVED***, "error message %s", "formatted")
func PanicsWithErrorf(t TestingT, errString string, f assert.PanicTestFunc, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.PanicsWithErrorf(t, errString, f, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// PanicsWithValue asserts that the code inside the specified PanicTestFunc panics, and that
// the recovered panic value equals the expected panic value.
//
//   assert.PanicsWithValue(t, "crazy error", func()***REMOVED*** GoCrazy() ***REMOVED***)
func PanicsWithValue(t TestingT, expected interface***REMOVED******REMOVED***, f assert.PanicTestFunc, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.PanicsWithValue(t, expected, f, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// PanicsWithValuef asserts that the code inside the specified PanicTestFunc panics, and that
// the recovered panic value equals the expected panic value.
//
//   assert.PanicsWithValuef(t, "crazy error", func()***REMOVED*** GoCrazy() ***REMOVED***, "error message %s", "formatted")
func PanicsWithValuef(t TestingT, expected interface***REMOVED******REMOVED***, f assert.PanicTestFunc, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.PanicsWithValuef(t, expected, f, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Panicsf asserts that the code inside the specified PanicTestFunc panics.
//
//   assert.Panicsf(t, func()***REMOVED*** GoCrazy() ***REMOVED***, "error message %s", "formatted")
func Panicsf(t TestingT, f assert.PanicTestFunc, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Panicsf(t, f, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Positive asserts that the specified element is positive
//
//    assert.Positive(t, 1)
//    assert.Positive(t, 1.23)
func Positive(t TestingT, e interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Positive(t, e, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Positivef asserts that the specified element is positive
//
//    assert.Positivef(t, 1, "error message %s", "formatted")
//    assert.Positivef(t, 1.23, "error message %s", "formatted")
func Positivef(t TestingT, e interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Positivef(t, e, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Regexp asserts that a specified regexp matches a string.
//
//  assert.Regexp(t, regexp.MustCompile("start"), "it's starting")
//  assert.Regexp(t, "start...$", "it's not starting")
func Regexp(t TestingT, rx interface***REMOVED******REMOVED***, str interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Regexp(t, rx, str, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Regexpf asserts that a specified regexp matches a string.
//
//  assert.Regexpf(t, regexp.MustCompile("start"), "it's starting", "error message %s", "formatted")
//  assert.Regexpf(t, "start...$", "it's not starting", "error message %s", "formatted")
func Regexpf(t TestingT, rx interface***REMOVED******REMOVED***, str interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Regexpf(t, rx, str, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Same asserts that two pointers reference the same object.
//
//    assert.Same(t, ptr1, ptr2)
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func Same(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Same(t, expected, actual, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Samef asserts that two pointers reference the same object.
//
//    assert.Samef(t, ptr1, ptr2, "error message %s", "formatted")
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func Samef(t TestingT, expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Samef(t, expected, actual, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Subset asserts that the specified list(array, slice...) contains all
// elements given in the specified subset(array, slice...).
//
//    assert.Subset(t, [1, 2, 3], [1, 2], "But [1, 2, 3] does contain [1, 2]")
func Subset(t TestingT, list interface***REMOVED******REMOVED***, subset interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Subset(t, list, subset, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Subsetf asserts that the specified list(array, slice...) contains all
// elements given in the specified subset(array, slice...).
//
//    assert.Subsetf(t, [1, 2, 3], [1, 2], "But [1, 2, 3] does contain [1, 2]", "error message %s", "formatted")
func Subsetf(t TestingT, list interface***REMOVED******REMOVED***, subset interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Subsetf(t, list, subset, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// True asserts that the specified value is true.
//
//    assert.True(t, myBool)
func True(t TestingT, value bool, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.True(t, value, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Truef asserts that the specified value is true.
//
//    assert.Truef(t, myBool, "error message %s", "formatted")
func Truef(t TestingT, value bool, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Truef(t, value, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// WithinDuration asserts that the two times are within duration delta of each other.
//
//   assert.WithinDuration(t, time.Now(), time.Now(), 10*time.Second)
func WithinDuration(t TestingT, expected time.Time, actual time.Time, delta time.Duration, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.WithinDuration(t, expected, actual, delta, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// WithinDurationf asserts that the two times are within duration delta of each other.
//
//   assert.WithinDurationf(t, time.Now(), time.Now(), 10*time.Second, "error message %s", "formatted")
func WithinDurationf(t TestingT, expected time.Time, actual time.Time, delta time.Duration, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.WithinDurationf(t, expected, actual, delta, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// YAMLEq asserts that two YAML strings are equivalent.
func YAMLEq(t TestingT, expected string, actual string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.YAMLEq(t, expected, actual, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// YAMLEqf asserts that two YAML strings are equivalent.
func YAMLEqf(t TestingT, expected string, actual string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.YAMLEqf(t, expected, actual, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Zero asserts that i is the zero value for its type.
func Zero(t TestingT, i interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Zero(t, i, msgAndArgs...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***

// Zerof asserts that i is the zero value for its type.
func Zerof(t TestingT, i interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	if assert.Zerof(t, i, msg, args...) ***REMOVED***
		return
	***REMOVED***
	t.FailNow()
***REMOVED***
