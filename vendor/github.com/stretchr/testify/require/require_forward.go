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
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Condition(a.t, comp, msgAndArgs...)
***REMOVED***

// Conditionf uses a Comparison to assert a complex condition.
func (a *Assertions) Conditionf(comp assert.Comparison, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Conditionf(a.t, comp, msg, args...)
***REMOVED***

// Contains asserts that the specified string, list(array, slice...) or map contains the
// specified substring or element.
//
//    a.Contains("Hello World", "World")
//    a.Contains(["Hello", "World"], "World")
//    a.Contains(***REMOVED***"Hello": "World"***REMOVED***, "Hello")
func (a *Assertions) Contains(s interface***REMOVED******REMOVED***, contains interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Contains(a.t, s, contains, msgAndArgs...)
***REMOVED***

// Containsf asserts that the specified string, list(array, slice...) or map contains the
// specified substring or element.
//
//    a.Containsf("Hello World", "World", "error message %s", "formatted")
//    a.Containsf(["Hello", "World"], "World", "error message %s", "formatted")
//    a.Containsf(***REMOVED***"Hello": "World"***REMOVED***, "Hello", "error message %s", "formatted")
func (a *Assertions) Containsf(s interface***REMOVED******REMOVED***, contains interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Containsf(a.t, s, contains, msg, args...)
***REMOVED***

// DirExists checks whether a directory exists in the given path. It also fails
// if the path is a file rather a directory or there is an error checking whether it exists.
func (a *Assertions) DirExists(path string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	DirExists(a.t, path, msgAndArgs...)
***REMOVED***

// DirExistsf checks whether a directory exists in the given path. It also fails
// if the path is a file rather a directory or there is an error checking whether it exists.
func (a *Assertions) DirExistsf(path string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	DirExistsf(a.t, path, msg, args...)
***REMOVED***

// ElementsMatch asserts that the specified listA(array, slice...) is equal to specified
// listB(array, slice...) ignoring the order of the elements. If there are duplicate elements,
// the number of appearances of each of them in both lists should match.
//
// a.ElementsMatch([1, 3, 2, 3], [1, 3, 3, 2])
func (a *Assertions) ElementsMatch(listA interface***REMOVED******REMOVED***, listB interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	ElementsMatch(a.t, listA, listB, msgAndArgs...)
***REMOVED***

// ElementsMatchf asserts that the specified listA(array, slice...) is equal to specified
// listB(array, slice...) ignoring the order of the elements. If there are duplicate elements,
// the number of appearances of each of them in both lists should match.
//
// a.ElementsMatchf([1, 3, 2, 3], [1, 3, 3, 2], "error message %s", "formatted")
func (a *Assertions) ElementsMatchf(listA interface***REMOVED******REMOVED***, listB interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	ElementsMatchf(a.t, listA, listB, msg, args...)
***REMOVED***

// Empty asserts that the specified object is empty.  I.e. nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  a.Empty(obj)
func (a *Assertions) Empty(object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Empty(a.t, object, msgAndArgs...)
***REMOVED***

// Emptyf asserts that the specified object is empty.  I.e. nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  a.Emptyf(obj, "error message %s", "formatted")
func (a *Assertions) Emptyf(object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
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
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Equal(a.t, expected, actual, msgAndArgs...)
***REMOVED***

// EqualError asserts that a function returned an error (i.e. not `nil`)
// and that it is equal to the provided error.
//
//   actualObj, err := SomeFunction()
//   a.EqualError(err,  expectedErrorString)
func (a *Assertions) EqualError(theError error, errString string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	EqualError(a.t, theError, errString, msgAndArgs...)
***REMOVED***

// EqualErrorf asserts that a function returned an error (i.e. not `nil`)
// and that it is equal to the provided error.
//
//   actualObj, err := SomeFunction()
//   a.EqualErrorf(err,  expectedErrorString, "error message %s", "formatted")
func (a *Assertions) EqualErrorf(theError error, errString string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	EqualErrorf(a.t, theError, errString, msg, args...)
***REMOVED***

// EqualValues asserts that two objects are equal or convertable to the same types
// and equal.
//
//    a.EqualValues(uint32(123), int32(123))
func (a *Assertions) EqualValues(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	EqualValues(a.t, expected, actual, msgAndArgs...)
***REMOVED***

// EqualValuesf asserts that two objects are equal or convertable to the same types
// and equal.
//
//    a.EqualValuesf(uint32(123), int32(123), "error message %s", "formatted")
func (a *Assertions) EqualValuesf(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
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
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Equalf(a.t, expected, actual, msg, args...)
***REMOVED***

// Error asserts that a function returned an error (i.e. not `nil`).
//
//   actualObj, err := SomeFunction()
//   if a.Error(err) ***REMOVED***
// 	   assert.Equal(t, expectedError, err)
//   ***REMOVED***
func (a *Assertions) Error(err error, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Error(a.t, err, msgAndArgs...)
***REMOVED***

// ErrorAs asserts that at least one of the errors in err's chain matches target, and if so, sets target to that error value.
// This is a wrapper for errors.As.
func (a *Assertions) ErrorAs(err error, target interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	ErrorAs(a.t, err, target, msgAndArgs...)
***REMOVED***

// ErrorAsf asserts that at least one of the errors in err's chain matches target, and if so, sets target to that error value.
// This is a wrapper for errors.As.
func (a *Assertions) ErrorAsf(err error, target interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	ErrorAsf(a.t, err, target, msg, args...)
***REMOVED***

// ErrorIs asserts that at least one of the errors in err's chain matches target.
// This is a wrapper for errors.Is.
func (a *Assertions) ErrorIs(err error, target error, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	ErrorIs(a.t, err, target, msgAndArgs...)
***REMOVED***

// ErrorIsf asserts that at least one of the errors in err's chain matches target.
// This is a wrapper for errors.Is.
func (a *Assertions) ErrorIsf(err error, target error, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	ErrorIsf(a.t, err, target, msg, args...)
***REMOVED***

// Errorf asserts that a function returned an error (i.e. not `nil`).
//
//   actualObj, err := SomeFunction()
//   if a.Errorf(err, "error message %s", "formatted") ***REMOVED***
// 	   assert.Equal(t, expectedErrorf, err)
//   ***REMOVED***
func (a *Assertions) Errorf(err error, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Errorf(a.t, err, msg, args...)
***REMOVED***

// Eventually asserts that given condition will be met in waitFor time,
// periodically checking target function each tick.
//
//    a.Eventually(func() bool ***REMOVED*** return true; ***REMOVED***, time.Second, 10*time.Millisecond)
func (a *Assertions) Eventually(condition func() bool, waitFor time.Duration, tick time.Duration, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Eventually(a.t, condition, waitFor, tick, msgAndArgs...)
***REMOVED***

// Eventuallyf asserts that given condition will be met in waitFor time,
// periodically checking target function each tick.
//
//    a.Eventuallyf(func() bool ***REMOVED*** return true; ***REMOVED***, time.Second, 10*time.Millisecond, "error message %s", "formatted")
func (a *Assertions) Eventuallyf(condition func() bool, waitFor time.Duration, tick time.Duration, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Eventuallyf(a.t, condition, waitFor, tick, msg, args...)
***REMOVED***

// Exactly asserts that two objects are equal in value and type.
//
//    a.Exactly(int32(123), int64(123))
func (a *Assertions) Exactly(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Exactly(a.t, expected, actual, msgAndArgs...)
***REMOVED***

// Exactlyf asserts that two objects are equal in value and type.
//
//    a.Exactlyf(int32(123), int64(123), "error message %s", "formatted")
func (a *Assertions) Exactlyf(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Exactlyf(a.t, expected, actual, msg, args...)
***REMOVED***

// Fail reports a failure through
func (a *Assertions) Fail(failureMessage string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Fail(a.t, failureMessage, msgAndArgs...)
***REMOVED***

// FailNow fails test
func (a *Assertions) FailNow(failureMessage string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	FailNow(a.t, failureMessage, msgAndArgs...)
***REMOVED***

// FailNowf fails test
func (a *Assertions) FailNowf(failureMessage string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	FailNowf(a.t, failureMessage, msg, args...)
***REMOVED***

// Failf reports a failure through
func (a *Assertions) Failf(failureMessage string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Failf(a.t, failureMessage, msg, args...)
***REMOVED***

// False asserts that the specified value is false.
//
//    a.False(myBool)
func (a *Assertions) False(value bool, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	False(a.t, value, msgAndArgs...)
***REMOVED***

// Falsef asserts that the specified value is false.
//
//    a.Falsef(myBool, "error message %s", "formatted")
func (a *Assertions) Falsef(value bool, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Falsef(a.t, value, msg, args...)
***REMOVED***

// FileExists checks whether a file exists in the given path. It also fails if
// the path points to a directory or there is an error when trying to check the file.
func (a *Assertions) FileExists(path string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	FileExists(a.t, path, msgAndArgs...)
***REMOVED***

// FileExistsf checks whether a file exists in the given path. It also fails if
// the path points to a directory or there is an error when trying to check the file.
func (a *Assertions) FileExistsf(path string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	FileExistsf(a.t, path, msg, args...)
***REMOVED***

// Greater asserts that the first element is greater than the second
//
//    a.Greater(2, 1)
//    a.Greater(float64(2), float64(1))
//    a.Greater("b", "a")
func (a *Assertions) Greater(e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Greater(a.t, e1, e2, msgAndArgs...)
***REMOVED***

// GreaterOrEqual asserts that the first element is greater than or equal to the second
//
//    a.GreaterOrEqual(2, 1)
//    a.GreaterOrEqual(2, 2)
//    a.GreaterOrEqual("b", "a")
//    a.GreaterOrEqual("b", "b")
func (a *Assertions) GreaterOrEqual(e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	GreaterOrEqual(a.t, e1, e2, msgAndArgs...)
***REMOVED***

// GreaterOrEqualf asserts that the first element is greater than or equal to the second
//
//    a.GreaterOrEqualf(2, 1, "error message %s", "formatted")
//    a.GreaterOrEqualf(2, 2, "error message %s", "formatted")
//    a.GreaterOrEqualf("b", "a", "error message %s", "formatted")
//    a.GreaterOrEqualf("b", "b", "error message %s", "formatted")
func (a *Assertions) GreaterOrEqualf(e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	GreaterOrEqualf(a.t, e1, e2, msg, args...)
***REMOVED***

// Greaterf asserts that the first element is greater than the second
//
//    a.Greaterf(2, 1, "error message %s", "formatted")
//    a.Greaterf(float64(2), float64(1), "error message %s", "formatted")
//    a.Greaterf("b", "a", "error message %s", "formatted")
func (a *Assertions) Greaterf(e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Greaterf(a.t, e1, e2, msg, args...)
***REMOVED***

// HTTPBodyContains asserts that a specified handler returns a
// body that contains a string.
//
//  a.HTTPBodyContains(myHandler, "GET", "www.google.com", nil, "I'm Feeling Lucky")
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPBodyContains(handler http.HandlerFunc, method string, url string, values url.Values, str interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	HTTPBodyContains(a.t, handler, method, url, values, str, msgAndArgs...)
***REMOVED***

// HTTPBodyContainsf asserts that a specified handler returns a
// body that contains a string.
//
//  a.HTTPBodyContainsf(myHandler, "GET", "www.google.com", nil, "I'm Feeling Lucky", "error message %s", "formatted")
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPBodyContainsf(handler http.HandlerFunc, method string, url string, values url.Values, str interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	HTTPBodyContainsf(a.t, handler, method, url, values, str, msg, args...)
***REMOVED***

// HTTPBodyNotContains asserts that a specified handler returns a
// body that does not contain a string.
//
//  a.HTTPBodyNotContains(myHandler, "GET", "www.google.com", nil, "I'm Feeling Lucky")
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPBodyNotContains(handler http.HandlerFunc, method string, url string, values url.Values, str interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	HTTPBodyNotContains(a.t, handler, method, url, values, str, msgAndArgs...)
***REMOVED***

// HTTPBodyNotContainsf asserts that a specified handler returns a
// body that does not contain a string.
//
//  a.HTTPBodyNotContainsf(myHandler, "GET", "www.google.com", nil, "I'm Feeling Lucky", "error message %s", "formatted")
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPBodyNotContainsf(handler http.HandlerFunc, method string, url string, values url.Values, str interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	HTTPBodyNotContainsf(a.t, handler, method, url, values, str, msg, args...)
***REMOVED***

// HTTPError asserts that a specified handler returns an error status code.
//
//  a.HTTPError(myHandler, "POST", "/a/b/c", url.Values***REMOVED***"a": []string***REMOVED***"b", "c"***REMOVED******REMOVED***
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPError(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	HTTPError(a.t, handler, method, url, values, msgAndArgs...)
***REMOVED***

// HTTPErrorf asserts that a specified handler returns an error status code.
//
//  a.HTTPErrorf(myHandler, "POST", "/a/b/c", url.Values***REMOVED***"a": []string***REMOVED***"b", "c"***REMOVED******REMOVED***
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPErrorf(handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	HTTPErrorf(a.t, handler, method, url, values, msg, args...)
***REMOVED***

// HTTPRedirect asserts that a specified handler returns a redirect status code.
//
//  a.HTTPRedirect(myHandler, "GET", "/a/b/c", url.Values***REMOVED***"a": []string***REMOVED***"b", "c"***REMOVED******REMOVED***
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPRedirect(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	HTTPRedirect(a.t, handler, method, url, values, msgAndArgs...)
***REMOVED***

// HTTPRedirectf asserts that a specified handler returns a redirect status code.
//
//  a.HTTPRedirectf(myHandler, "GET", "/a/b/c", url.Values***REMOVED***"a": []string***REMOVED***"b", "c"***REMOVED******REMOVED***
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPRedirectf(handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	HTTPRedirectf(a.t, handler, method, url, values, msg, args...)
***REMOVED***

// HTTPStatusCode asserts that a specified handler returns a specified status code.
//
//  a.HTTPStatusCode(myHandler, "GET", "/notImplemented", nil, 501)
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPStatusCode(handler http.HandlerFunc, method string, url string, values url.Values, statuscode int, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	HTTPStatusCode(a.t, handler, method, url, values, statuscode, msgAndArgs...)
***REMOVED***

// HTTPStatusCodef asserts that a specified handler returns a specified status code.
//
//  a.HTTPStatusCodef(myHandler, "GET", "/notImplemented", nil, 501, "error message %s", "formatted")
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPStatusCodef(handler http.HandlerFunc, method string, url string, values url.Values, statuscode int, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	HTTPStatusCodef(a.t, handler, method, url, values, statuscode, msg, args...)
***REMOVED***

// HTTPSuccess asserts that a specified handler returns a success status code.
//
//  a.HTTPSuccess(myHandler, "POST", "http://www.google.com", nil)
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPSuccess(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	HTTPSuccess(a.t, handler, method, url, values, msgAndArgs...)
***REMOVED***

// HTTPSuccessf asserts that a specified handler returns a success status code.
//
//  a.HTTPSuccessf(myHandler, "POST", "http://www.google.com", nil, "error message %s", "formatted")
//
// Returns whether the assertion was successful (true) or not (false).
func (a *Assertions) HTTPSuccessf(handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	HTTPSuccessf(a.t, handler, method, url, values, msg, args...)
***REMOVED***

// Implements asserts that an object is implemented by the specified interface.
//
//    a.Implements((*MyInterface)(nil), new(MyObject))
func (a *Assertions) Implements(interfaceObject interface***REMOVED******REMOVED***, object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Implements(a.t, interfaceObject, object, msgAndArgs...)
***REMOVED***

// Implementsf asserts that an object is implemented by the specified interface.
//
//    a.Implementsf((*MyInterface)(nil), new(MyObject), "error message %s", "formatted")
func (a *Assertions) Implementsf(interfaceObject interface***REMOVED******REMOVED***, object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Implementsf(a.t, interfaceObject, object, msg, args...)
***REMOVED***

// InDelta asserts that the two numerals are within delta of each other.
//
// 	 a.InDelta(math.Pi, 22/7.0, 0.01)
func (a *Assertions) InDelta(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, delta float64, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	InDelta(a.t, expected, actual, delta, msgAndArgs...)
***REMOVED***

// InDeltaMapValues is the same as InDelta, but it compares all values between two maps. Both maps must have exactly the same keys.
func (a *Assertions) InDeltaMapValues(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, delta float64, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	InDeltaMapValues(a.t, expected, actual, delta, msgAndArgs...)
***REMOVED***

// InDeltaMapValuesf is the same as InDelta, but it compares all values between two maps. Both maps must have exactly the same keys.
func (a *Assertions) InDeltaMapValuesf(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, delta float64, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	InDeltaMapValuesf(a.t, expected, actual, delta, msg, args...)
***REMOVED***

// InDeltaSlice is the same as InDelta, except it compares two slices.
func (a *Assertions) InDeltaSlice(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, delta float64, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	InDeltaSlice(a.t, expected, actual, delta, msgAndArgs...)
***REMOVED***

// InDeltaSlicef is the same as InDelta, except it compares two slices.
func (a *Assertions) InDeltaSlicef(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, delta float64, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	InDeltaSlicef(a.t, expected, actual, delta, msg, args...)
***REMOVED***

// InDeltaf asserts that the two numerals are within delta of each other.
//
// 	 a.InDeltaf(math.Pi, 22/7.0, 0.01, "error message %s", "formatted")
func (a *Assertions) InDeltaf(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, delta float64, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	InDeltaf(a.t, expected, actual, delta, msg, args...)
***REMOVED***

// InEpsilon asserts that expected and actual have a relative error less than epsilon
func (a *Assertions) InEpsilon(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, epsilon float64, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	InEpsilon(a.t, expected, actual, epsilon, msgAndArgs...)
***REMOVED***

// InEpsilonSlice is the same as InEpsilon, except it compares each value from two slices.
func (a *Assertions) InEpsilonSlice(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, epsilon float64, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	InEpsilonSlice(a.t, expected, actual, epsilon, msgAndArgs...)
***REMOVED***

// InEpsilonSlicef is the same as InEpsilon, except it compares each value from two slices.
func (a *Assertions) InEpsilonSlicef(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, epsilon float64, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	InEpsilonSlicef(a.t, expected, actual, epsilon, msg, args...)
***REMOVED***

// InEpsilonf asserts that expected and actual have a relative error less than epsilon
func (a *Assertions) InEpsilonf(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, epsilon float64, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	InEpsilonf(a.t, expected, actual, epsilon, msg, args...)
***REMOVED***

// IsDecreasing asserts that the collection is decreasing
//
//    a.IsDecreasing([]int***REMOVED***2, 1, 0***REMOVED***)
//    a.IsDecreasing([]float***REMOVED***2, 1***REMOVED***)
//    a.IsDecreasing([]string***REMOVED***"b", "a"***REMOVED***)
func (a *Assertions) IsDecreasing(object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	IsDecreasing(a.t, object, msgAndArgs...)
***REMOVED***

// IsDecreasingf asserts that the collection is decreasing
//
//    a.IsDecreasingf([]int***REMOVED***2, 1, 0***REMOVED***, "error message %s", "formatted")
//    a.IsDecreasingf([]float***REMOVED***2, 1***REMOVED***, "error message %s", "formatted")
//    a.IsDecreasingf([]string***REMOVED***"b", "a"***REMOVED***, "error message %s", "formatted")
func (a *Assertions) IsDecreasingf(object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	IsDecreasingf(a.t, object, msg, args...)
***REMOVED***

// IsIncreasing asserts that the collection is increasing
//
//    a.IsIncreasing([]int***REMOVED***1, 2, 3***REMOVED***)
//    a.IsIncreasing([]float***REMOVED***1, 2***REMOVED***)
//    a.IsIncreasing([]string***REMOVED***"a", "b"***REMOVED***)
func (a *Assertions) IsIncreasing(object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	IsIncreasing(a.t, object, msgAndArgs...)
***REMOVED***

// IsIncreasingf asserts that the collection is increasing
//
//    a.IsIncreasingf([]int***REMOVED***1, 2, 3***REMOVED***, "error message %s", "formatted")
//    a.IsIncreasingf([]float***REMOVED***1, 2***REMOVED***, "error message %s", "formatted")
//    a.IsIncreasingf([]string***REMOVED***"a", "b"***REMOVED***, "error message %s", "formatted")
func (a *Assertions) IsIncreasingf(object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	IsIncreasingf(a.t, object, msg, args...)
***REMOVED***

// IsNonDecreasing asserts that the collection is not decreasing
//
//    a.IsNonDecreasing([]int***REMOVED***1, 1, 2***REMOVED***)
//    a.IsNonDecreasing([]float***REMOVED***1, 2***REMOVED***)
//    a.IsNonDecreasing([]string***REMOVED***"a", "b"***REMOVED***)
func (a *Assertions) IsNonDecreasing(object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	IsNonDecreasing(a.t, object, msgAndArgs...)
***REMOVED***

// IsNonDecreasingf asserts that the collection is not decreasing
//
//    a.IsNonDecreasingf([]int***REMOVED***1, 1, 2***REMOVED***, "error message %s", "formatted")
//    a.IsNonDecreasingf([]float***REMOVED***1, 2***REMOVED***, "error message %s", "formatted")
//    a.IsNonDecreasingf([]string***REMOVED***"a", "b"***REMOVED***, "error message %s", "formatted")
func (a *Assertions) IsNonDecreasingf(object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	IsNonDecreasingf(a.t, object, msg, args...)
***REMOVED***

// IsNonIncreasing asserts that the collection is not increasing
//
//    a.IsNonIncreasing([]int***REMOVED***2, 1, 1***REMOVED***)
//    a.IsNonIncreasing([]float***REMOVED***2, 1***REMOVED***)
//    a.IsNonIncreasing([]string***REMOVED***"b", "a"***REMOVED***)
func (a *Assertions) IsNonIncreasing(object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	IsNonIncreasing(a.t, object, msgAndArgs...)
***REMOVED***

// IsNonIncreasingf asserts that the collection is not increasing
//
//    a.IsNonIncreasingf([]int***REMOVED***2, 1, 1***REMOVED***, "error message %s", "formatted")
//    a.IsNonIncreasingf([]float***REMOVED***2, 1***REMOVED***, "error message %s", "formatted")
//    a.IsNonIncreasingf([]string***REMOVED***"b", "a"***REMOVED***, "error message %s", "formatted")
func (a *Assertions) IsNonIncreasingf(object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	IsNonIncreasingf(a.t, object, msg, args...)
***REMOVED***

// IsType asserts that the specified objects are of the same type.
func (a *Assertions) IsType(expectedType interface***REMOVED******REMOVED***, object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	IsType(a.t, expectedType, object, msgAndArgs...)
***REMOVED***

// IsTypef asserts that the specified objects are of the same type.
func (a *Assertions) IsTypef(expectedType interface***REMOVED******REMOVED***, object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	IsTypef(a.t, expectedType, object, msg, args...)
***REMOVED***

// JSONEq asserts that two JSON strings are equivalent.
//
//  a.JSONEq(`***REMOVED***"hello": "world", "foo": "bar"***REMOVED***`, `***REMOVED***"foo": "bar", "hello": "world"***REMOVED***`)
func (a *Assertions) JSONEq(expected string, actual string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	JSONEq(a.t, expected, actual, msgAndArgs...)
***REMOVED***

// JSONEqf asserts that two JSON strings are equivalent.
//
//  a.JSONEqf(`***REMOVED***"hello": "world", "foo": "bar"***REMOVED***`, `***REMOVED***"foo": "bar", "hello": "world"***REMOVED***`, "error message %s", "formatted")
func (a *Assertions) JSONEqf(expected string, actual string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	JSONEqf(a.t, expected, actual, msg, args...)
***REMOVED***

// Len asserts that the specified object has specific length.
// Len also fails if the object has a type that len() not accept.
//
//    a.Len(mySlice, 3)
func (a *Assertions) Len(object interface***REMOVED******REMOVED***, length int, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Len(a.t, object, length, msgAndArgs...)
***REMOVED***

// Lenf asserts that the specified object has specific length.
// Lenf also fails if the object has a type that len() not accept.
//
//    a.Lenf(mySlice, 3, "error message %s", "formatted")
func (a *Assertions) Lenf(object interface***REMOVED******REMOVED***, length int, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Lenf(a.t, object, length, msg, args...)
***REMOVED***

// Less asserts that the first element is less than the second
//
//    a.Less(1, 2)
//    a.Less(float64(1), float64(2))
//    a.Less("a", "b")
func (a *Assertions) Less(e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Less(a.t, e1, e2, msgAndArgs...)
***REMOVED***

// LessOrEqual asserts that the first element is less than or equal to the second
//
//    a.LessOrEqual(1, 2)
//    a.LessOrEqual(2, 2)
//    a.LessOrEqual("a", "b")
//    a.LessOrEqual("b", "b")
func (a *Assertions) LessOrEqual(e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	LessOrEqual(a.t, e1, e2, msgAndArgs...)
***REMOVED***

// LessOrEqualf asserts that the first element is less than or equal to the second
//
//    a.LessOrEqualf(1, 2, "error message %s", "formatted")
//    a.LessOrEqualf(2, 2, "error message %s", "formatted")
//    a.LessOrEqualf("a", "b", "error message %s", "formatted")
//    a.LessOrEqualf("b", "b", "error message %s", "formatted")
func (a *Assertions) LessOrEqualf(e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	LessOrEqualf(a.t, e1, e2, msg, args...)
***REMOVED***

// Lessf asserts that the first element is less than the second
//
//    a.Lessf(1, 2, "error message %s", "formatted")
//    a.Lessf(float64(1), float64(2), "error message %s", "formatted")
//    a.Lessf("a", "b", "error message %s", "formatted")
func (a *Assertions) Lessf(e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Lessf(a.t, e1, e2, msg, args...)
***REMOVED***

// Negative asserts that the specified element is negative
//
//    a.Negative(-1)
//    a.Negative(-1.23)
func (a *Assertions) Negative(e interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Negative(a.t, e, msgAndArgs...)
***REMOVED***

// Negativef asserts that the specified element is negative
//
//    a.Negativef(-1, "error message %s", "formatted")
//    a.Negativef(-1.23, "error message %s", "formatted")
func (a *Assertions) Negativef(e interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Negativef(a.t, e, msg, args...)
***REMOVED***

// Never asserts that the given condition doesn't satisfy in waitFor time,
// periodically checking the target function each tick.
//
//    a.Never(func() bool ***REMOVED*** return false; ***REMOVED***, time.Second, 10*time.Millisecond)
func (a *Assertions) Never(condition func() bool, waitFor time.Duration, tick time.Duration, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Never(a.t, condition, waitFor, tick, msgAndArgs...)
***REMOVED***

// Neverf asserts that the given condition doesn't satisfy in waitFor time,
// periodically checking the target function each tick.
//
//    a.Neverf(func() bool ***REMOVED*** return false; ***REMOVED***, time.Second, 10*time.Millisecond, "error message %s", "formatted")
func (a *Assertions) Neverf(condition func() bool, waitFor time.Duration, tick time.Duration, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Neverf(a.t, condition, waitFor, tick, msg, args...)
***REMOVED***

// Nil asserts that the specified object is nil.
//
//    a.Nil(err)
func (a *Assertions) Nil(object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Nil(a.t, object, msgAndArgs...)
***REMOVED***

// Nilf asserts that the specified object is nil.
//
//    a.Nilf(err, "error message %s", "formatted")
func (a *Assertions) Nilf(object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Nilf(a.t, object, msg, args...)
***REMOVED***

// NoDirExists checks whether a directory does not exist in the given path.
// It fails if the path points to an existing _directory_ only.
func (a *Assertions) NoDirExists(path string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NoDirExists(a.t, path, msgAndArgs...)
***REMOVED***

// NoDirExistsf checks whether a directory does not exist in the given path.
// It fails if the path points to an existing _directory_ only.
func (a *Assertions) NoDirExistsf(path string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NoDirExistsf(a.t, path, msg, args...)
***REMOVED***

// NoError asserts that a function returned no error (i.e. `nil`).
//
//   actualObj, err := SomeFunction()
//   if a.NoError(err) ***REMOVED***
// 	   assert.Equal(t, expectedObj, actualObj)
//   ***REMOVED***
func (a *Assertions) NoError(err error, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NoError(a.t, err, msgAndArgs...)
***REMOVED***

// NoErrorf asserts that a function returned no error (i.e. `nil`).
//
//   actualObj, err := SomeFunction()
//   if a.NoErrorf(err, "error message %s", "formatted") ***REMOVED***
// 	   assert.Equal(t, expectedObj, actualObj)
//   ***REMOVED***
func (a *Assertions) NoErrorf(err error, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NoErrorf(a.t, err, msg, args...)
***REMOVED***

// NoFileExists checks whether a file does not exist in a given path. It fails
// if the path points to an existing _file_ only.
func (a *Assertions) NoFileExists(path string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NoFileExists(a.t, path, msgAndArgs...)
***REMOVED***

// NoFileExistsf checks whether a file does not exist in a given path. It fails
// if the path points to an existing _file_ only.
func (a *Assertions) NoFileExistsf(path string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NoFileExistsf(a.t, path, msg, args...)
***REMOVED***

// NotContains asserts that the specified string, list(array, slice...) or map does NOT contain the
// specified substring or element.
//
//    a.NotContains("Hello World", "Earth")
//    a.NotContains(["Hello", "World"], "Earth")
//    a.NotContains(***REMOVED***"Hello": "World"***REMOVED***, "Earth")
func (a *Assertions) NotContains(s interface***REMOVED******REMOVED***, contains interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotContains(a.t, s, contains, msgAndArgs...)
***REMOVED***

// NotContainsf asserts that the specified string, list(array, slice...) or map does NOT contain the
// specified substring or element.
//
//    a.NotContainsf("Hello World", "Earth", "error message %s", "formatted")
//    a.NotContainsf(["Hello", "World"], "Earth", "error message %s", "formatted")
//    a.NotContainsf(***REMOVED***"Hello": "World"***REMOVED***, "Earth", "error message %s", "formatted")
func (a *Assertions) NotContainsf(s interface***REMOVED******REMOVED***, contains interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotContainsf(a.t, s, contains, msg, args...)
***REMOVED***

// NotEmpty asserts that the specified object is NOT empty.  I.e. not nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  if a.NotEmpty(obj) ***REMOVED***
//    assert.Equal(t, "two", obj[1])
//  ***REMOVED***
func (a *Assertions) NotEmpty(object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotEmpty(a.t, object, msgAndArgs...)
***REMOVED***

// NotEmptyf asserts that the specified object is NOT empty.  I.e. not nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  if a.NotEmptyf(obj, "error message %s", "formatted") ***REMOVED***
//    assert.Equal(t, "two", obj[1])
//  ***REMOVED***
func (a *Assertions) NotEmptyf(object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotEmptyf(a.t, object, msg, args...)
***REMOVED***

// NotEqual asserts that the specified values are NOT equal.
//
//    a.NotEqual(obj1, obj2)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
func (a *Assertions) NotEqual(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotEqual(a.t, expected, actual, msgAndArgs...)
***REMOVED***

// NotEqualValues asserts that two objects are not equal even when converted to the same type
//
//    a.NotEqualValues(obj1, obj2)
func (a *Assertions) NotEqualValues(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotEqualValues(a.t, expected, actual, msgAndArgs...)
***REMOVED***

// NotEqualValuesf asserts that two objects are not equal even when converted to the same type
//
//    a.NotEqualValuesf(obj1, obj2, "error message %s", "formatted")
func (a *Assertions) NotEqualValuesf(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotEqualValuesf(a.t, expected, actual, msg, args...)
***REMOVED***

// NotEqualf asserts that the specified values are NOT equal.
//
//    a.NotEqualf(obj1, obj2, "error message %s", "formatted")
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
func (a *Assertions) NotEqualf(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotEqualf(a.t, expected, actual, msg, args...)
***REMOVED***

// NotErrorIs asserts that at none of the errors in err's chain matches target.
// This is a wrapper for errors.Is.
func (a *Assertions) NotErrorIs(err error, target error, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotErrorIs(a.t, err, target, msgAndArgs...)
***REMOVED***

// NotErrorIsf asserts that at none of the errors in err's chain matches target.
// This is a wrapper for errors.Is.
func (a *Assertions) NotErrorIsf(err error, target error, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotErrorIsf(a.t, err, target, msg, args...)
***REMOVED***

// NotNil asserts that the specified object is not nil.
//
//    a.NotNil(err)
func (a *Assertions) NotNil(object interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotNil(a.t, object, msgAndArgs...)
***REMOVED***

// NotNilf asserts that the specified object is not nil.
//
//    a.NotNilf(err, "error message %s", "formatted")
func (a *Assertions) NotNilf(object interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotNilf(a.t, object, msg, args...)
***REMOVED***

// NotPanics asserts that the code inside the specified PanicTestFunc does NOT panic.
//
//   a.NotPanics(func()***REMOVED*** RemainCalm() ***REMOVED***)
func (a *Assertions) NotPanics(f assert.PanicTestFunc, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotPanics(a.t, f, msgAndArgs...)
***REMOVED***

// NotPanicsf asserts that the code inside the specified PanicTestFunc does NOT panic.
//
//   a.NotPanicsf(func()***REMOVED*** RemainCalm() ***REMOVED***, "error message %s", "formatted")
func (a *Assertions) NotPanicsf(f assert.PanicTestFunc, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotPanicsf(a.t, f, msg, args...)
***REMOVED***

// NotRegexp asserts that a specified regexp does not match a string.
//
//  a.NotRegexp(regexp.MustCompile("starts"), "it's starting")
//  a.NotRegexp("^start", "it's not starting")
func (a *Assertions) NotRegexp(rx interface***REMOVED******REMOVED***, str interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotRegexp(a.t, rx, str, msgAndArgs...)
***REMOVED***

// NotRegexpf asserts that a specified regexp does not match a string.
//
//  a.NotRegexpf(regexp.MustCompile("starts"), "it's starting", "error message %s", "formatted")
//  a.NotRegexpf("^start", "it's not starting", "error message %s", "formatted")
func (a *Assertions) NotRegexpf(rx interface***REMOVED******REMOVED***, str interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotRegexpf(a.t, rx, str, msg, args...)
***REMOVED***

// NotSame asserts that two pointers do not reference the same object.
//
//    a.NotSame(ptr1, ptr2)
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func (a *Assertions) NotSame(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotSame(a.t, expected, actual, msgAndArgs...)
***REMOVED***

// NotSamef asserts that two pointers do not reference the same object.
//
//    a.NotSamef(ptr1, ptr2, "error message %s", "formatted")
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func (a *Assertions) NotSamef(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotSamef(a.t, expected, actual, msg, args...)
***REMOVED***

// NotSubset asserts that the specified list(array, slice...) contains not all
// elements given in the specified subset(array, slice...).
//
//    a.NotSubset([1, 3, 4], [1, 2], "But [1, 3, 4] does not contain [1, 2]")
func (a *Assertions) NotSubset(list interface***REMOVED******REMOVED***, subset interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotSubset(a.t, list, subset, msgAndArgs...)
***REMOVED***

// NotSubsetf asserts that the specified list(array, slice...) contains not all
// elements given in the specified subset(array, slice...).
//
//    a.NotSubsetf([1, 3, 4], [1, 2], "But [1, 3, 4] does not contain [1, 2]", "error message %s", "formatted")
func (a *Assertions) NotSubsetf(list interface***REMOVED******REMOVED***, subset interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotSubsetf(a.t, list, subset, msg, args...)
***REMOVED***

// NotZero asserts that i is not the zero value for its type.
func (a *Assertions) NotZero(i interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotZero(a.t, i, msgAndArgs...)
***REMOVED***

// NotZerof asserts that i is not the zero value for its type.
func (a *Assertions) NotZerof(i interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	NotZerof(a.t, i, msg, args...)
***REMOVED***

// Panics asserts that the code inside the specified PanicTestFunc panics.
//
//   a.Panics(func()***REMOVED*** GoCrazy() ***REMOVED***)
func (a *Assertions) Panics(f assert.PanicTestFunc, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Panics(a.t, f, msgAndArgs...)
***REMOVED***

// PanicsWithError asserts that the code inside the specified PanicTestFunc
// panics, and that the recovered panic value is an error that satisfies the
// EqualError comparison.
//
//   a.PanicsWithError("crazy error", func()***REMOVED*** GoCrazy() ***REMOVED***)
func (a *Assertions) PanicsWithError(errString string, f assert.PanicTestFunc, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	PanicsWithError(a.t, errString, f, msgAndArgs...)
***REMOVED***

// PanicsWithErrorf asserts that the code inside the specified PanicTestFunc
// panics, and that the recovered panic value is an error that satisfies the
// EqualError comparison.
//
//   a.PanicsWithErrorf("crazy error", func()***REMOVED*** GoCrazy() ***REMOVED***, "error message %s", "formatted")
func (a *Assertions) PanicsWithErrorf(errString string, f assert.PanicTestFunc, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	PanicsWithErrorf(a.t, errString, f, msg, args...)
***REMOVED***

// PanicsWithValue asserts that the code inside the specified PanicTestFunc panics, and that
// the recovered panic value equals the expected panic value.
//
//   a.PanicsWithValue("crazy error", func()***REMOVED*** GoCrazy() ***REMOVED***)
func (a *Assertions) PanicsWithValue(expected interface***REMOVED******REMOVED***, f assert.PanicTestFunc, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	PanicsWithValue(a.t, expected, f, msgAndArgs...)
***REMOVED***

// PanicsWithValuef asserts that the code inside the specified PanicTestFunc panics, and that
// the recovered panic value equals the expected panic value.
//
//   a.PanicsWithValuef("crazy error", func()***REMOVED*** GoCrazy() ***REMOVED***, "error message %s", "formatted")
func (a *Assertions) PanicsWithValuef(expected interface***REMOVED******REMOVED***, f assert.PanicTestFunc, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	PanicsWithValuef(a.t, expected, f, msg, args...)
***REMOVED***

// Panicsf asserts that the code inside the specified PanicTestFunc panics.
//
//   a.Panicsf(func()***REMOVED*** GoCrazy() ***REMOVED***, "error message %s", "formatted")
func (a *Assertions) Panicsf(f assert.PanicTestFunc, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Panicsf(a.t, f, msg, args...)
***REMOVED***

// Positive asserts that the specified element is positive
//
//    a.Positive(1)
//    a.Positive(1.23)
func (a *Assertions) Positive(e interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Positive(a.t, e, msgAndArgs...)
***REMOVED***

// Positivef asserts that the specified element is positive
//
//    a.Positivef(1, "error message %s", "formatted")
//    a.Positivef(1.23, "error message %s", "formatted")
func (a *Assertions) Positivef(e interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Positivef(a.t, e, msg, args...)
***REMOVED***

// Regexp asserts that a specified regexp matches a string.
//
//  a.Regexp(regexp.MustCompile("start"), "it's starting")
//  a.Regexp("start...$", "it's not starting")
func (a *Assertions) Regexp(rx interface***REMOVED******REMOVED***, str interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Regexp(a.t, rx, str, msgAndArgs...)
***REMOVED***

// Regexpf asserts that a specified regexp matches a string.
//
//  a.Regexpf(regexp.MustCompile("start"), "it's starting", "error message %s", "formatted")
//  a.Regexpf("start...$", "it's not starting", "error message %s", "formatted")
func (a *Assertions) Regexpf(rx interface***REMOVED******REMOVED***, str interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Regexpf(a.t, rx, str, msg, args...)
***REMOVED***

// Same asserts that two pointers reference the same object.
//
//    a.Same(ptr1, ptr2)
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func (a *Assertions) Same(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Same(a.t, expected, actual, msgAndArgs...)
***REMOVED***

// Samef asserts that two pointers reference the same object.
//
//    a.Samef(ptr1, ptr2, "error message %s", "formatted")
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func (a *Assertions) Samef(expected interface***REMOVED******REMOVED***, actual interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Samef(a.t, expected, actual, msg, args...)
***REMOVED***

// Subset asserts that the specified list(array, slice...) contains all
// elements given in the specified subset(array, slice...).
//
//    a.Subset([1, 2, 3], [1, 2], "But [1, 2, 3] does contain [1, 2]")
func (a *Assertions) Subset(list interface***REMOVED******REMOVED***, subset interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Subset(a.t, list, subset, msgAndArgs...)
***REMOVED***

// Subsetf asserts that the specified list(array, slice...) contains all
// elements given in the specified subset(array, slice...).
//
//    a.Subsetf([1, 2, 3], [1, 2], "But [1, 2, 3] does contain [1, 2]", "error message %s", "formatted")
func (a *Assertions) Subsetf(list interface***REMOVED******REMOVED***, subset interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Subsetf(a.t, list, subset, msg, args...)
***REMOVED***

// True asserts that the specified value is true.
//
//    a.True(myBool)
func (a *Assertions) True(value bool, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	True(a.t, value, msgAndArgs...)
***REMOVED***

// Truef asserts that the specified value is true.
//
//    a.Truef(myBool, "error message %s", "formatted")
func (a *Assertions) Truef(value bool, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Truef(a.t, value, msg, args...)
***REMOVED***

// WithinDuration asserts that the two times are within duration delta of each other.
//
//   a.WithinDuration(time.Now(), time.Now(), 10*time.Second)
func (a *Assertions) WithinDuration(expected time.Time, actual time.Time, delta time.Duration, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	WithinDuration(a.t, expected, actual, delta, msgAndArgs...)
***REMOVED***

// WithinDurationf asserts that the two times are within duration delta of each other.
//
//   a.WithinDurationf(time.Now(), time.Now(), 10*time.Second, "error message %s", "formatted")
func (a *Assertions) WithinDurationf(expected time.Time, actual time.Time, delta time.Duration, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	WithinDurationf(a.t, expected, actual, delta, msg, args...)
***REMOVED***

// YAMLEq asserts that two YAML strings are equivalent.
func (a *Assertions) YAMLEq(expected string, actual string, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	YAMLEq(a.t, expected, actual, msgAndArgs...)
***REMOVED***

// YAMLEqf asserts that two YAML strings are equivalent.
func (a *Assertions) YAMLEqf(expected string, actual string, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	YAMLEqf(a.t, expected, actual, msg, args...)
***REMOVED***

// Zero asserts that i is the zero value for its type.
func (a *Assertions) Zero(i interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Zero(a.t, i, msgAndArgs...)
***REMOVED***

// Zerof asserts that i is the zero value for its type.
func (a *Assertions) Zerof(i interface***REMOVED******REMOVED***, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if h, ok := a.t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***
	Zerof(a.t, i, msg, args...)
***REMOVED***
