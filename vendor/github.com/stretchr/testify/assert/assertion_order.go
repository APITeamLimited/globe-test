package assert

import (
	"fmt"
	"reflect"
)

func compare(obj1, obj2 interface***REMOVED******REMOVED***, kind reflect.Kind) (int, bool) ***REMOVED***
	switch kind ***REMOVED***
	case reflect.Int:
		***REMOVED***
			intobj1 := obj1.(int)
			intobj2 := obj2.(int)
			if intobj1 > intobj2 ***REMOVED***
				return -1, true
			***REMOVED***
			if intobj1 == intobj2 ***REMOVED***
				return 0, true
			***REMOVED***
			if intobj1 < intobj2 ***REMOVED***
				return 1, true
			***REMOVED***
		***REMOVED***
	case reflect.Int8:
		***REMOVED***
			int8obj1 := obj1.(int8)
			int8obj2 := obj2.(int8)
			if int8obj1 > int8obj2 ***REMOVED***
				return -1, true
			***REMOVED***
			if int8obj1 == int8obj2 ***REMOVED***
				return 0, true
			***REMOVED***
			if int8obj1 < int8obj2 ***REMOVED***
				return 1, true
			***REMOVED***
		***REMOVED***
	case reflect.Int16:
		***REMOVED***
			int16obj1 := obj1.(int16)
			int16obj2 := obj2.(int16)
			if int16obj1 > int16obj2 ***REMOVED***
				return -1, true
			***REMOVED***
			if int16obj1 == int16obj2 ***REMOVED***
				return 0, true
			***REMOVED***
			if int16obj1 < int16obj2 ***REMOVED***
				return 1, true
			***REMOVED***
		***REMOVED***
	case reflect.Int32:
		***REMOVED***
			int32obj1 := obj1.(int32)
			int32obj2 := obj2.(int32)
			if int32obj1 > int32obj2 ***REMOVED***
				return -1, true
			***REMOVED***
			if int32obj1 == int32obj2 ***REMOVED***
				return 0, true
			***REMOVED***
			if int32obj1 < int32obj2 ***REMOVED***
				return 1, true
			***REMOVED***
		***REMOVED***
	case reflect.Int64:
		***REMOVED***
			int64obj1 := obj1.(int64)
			int64obj2 := obj2.(int64)
			if int64obj1 > int64obj2 ***REMOVED***
				return -1, true
			***REMOVED***
			if int64obj1 == int64obj2 ***REMOVED***
				return 0, true
			***REMOVED***
			if int64obj1 < int64obj2 ***REMOVED***
				return 1, true
			***REMOVED***
		***REMOVED***
	case reflect.Uint:
		***REMOVED***
			uintobj1 := obj1.(uint)
			uintobj2 := obj2.(uint)
			if uintobj1 > uintobj2 ***REMOVED***
				return -1, true
			***REMOVED***
			if uintobj1 == uintobj2 ***REMOVED***
				return 0, true
			***REMOVED***
			if uintobj1 < uintobj2 ***REMOVED***
				return 1, true
			***REMOVED***
		***REMOVED***
	case reflect.Uint8:
		***REMOVED***
			uint8obj1 := obj1.(uint8)
			uint8obj2 := obj2.(uint8)
			if uint8obj1 > uint8obj2 ***REMOVED***
				return -1, true
			***REMOVED***
			if uint8obj1 == uint8obj2 ***REMOVED***
				return 0, true
			***REMOVED***
			if uint8obj1 < uint8obj2 ***REMOVED***
				return 1, true
			***REMOVED***
		***REMOVED***
	case reflect.Uint16:
		***REMOVED***
			uint16obj1 := obj1.(uint16)
			uint16obj2 := obj2.(uint16)
			if uint16obj1 > uint16obj2 ***REMOVED***
				return -1, true
			***REMOVED***
			if uint16obj1 == uint16obj2 ***REMOVED***
				return 0, true
			***REMOVED***
			if uint16obj1 < uint16obj2 ***REMOVED***
				return 1, true
			***REMOVED***
		***REMOVED***
	case reflect.Uint32:
		***REMOVED***
			uint32obj1 := obj1.(uint32)
			uint32obj2 := obj2.(uint32)
			if uint32obj1 > uint32obj2 ***REMOVED***
				return -1, true
			***REMOVED***
			if uint32obj1 == uint32obj2 ***REMOVED***
				return 0, true
			***REMOVED***
			if uint32obj1 < uint32obj2 ***REMOVED***
				return 1, true
			***REMOVED***
		***REMOVED***
	case reflect.Uint64:
		***REMOVED***
			uint64obj1 := obj1.(uint64)
			uint64obj2 := obj2.(uint64)
			if uint64obj1 > uint64obj2 ***REMOVED***
				return -1, true
			***REMOVED***
			if uint64obj1 == uint64obj2 ***REMOVED***
				return 0, true
			***REMOVED***
			if uint64obj1 < uint64obj2 ***REMOVED***
				return 1, true
			***REMOVED***
		***REMOVED***
	case reflect.Float32:
		***REMOVED***
			float32obj1 := obj1.(float32)
			float32obj2 := obj2.(float32)
			if float32obj1 > float32obj2 ***REMOVED***
				return -1, true
			***REMOVED***
			if float32obj1 == float32obj2 ***REMOVED***
				return 0, true
			***REMOVED***
			if float32obj1 < float32obj2 ***REMOVED***
				return 1, true
			***REMOVED***
		***REMOVED***
	case reflect.Float64:
		***REMOVED***
			float64obj1 := obj1.(float64)
			float64obj2 := obj2.(float64)
			if float64obj1 > float64obj2 ***REMOVED***
				return -1, true
			***REMOVED***
			if float64obj1 == float64obj2 ***REMOVED***
				return 0, true
			***REMOVED***
			if float64obj1 < float64obj2 ***REMOVED***
				return 1, true
			***REMOVED***
		***REMOVED***
	case reflect.String:
		***REMOVED***
			stringobj1 := obj1.(string)
			stringobj2 := obj2.(string)
			if stringobj1 > stringobj2 ***REMOVED***
				return -1, true
			***REMOVED***
			if stringobj1 == stringobj2 ***REMOVED***
				return 0, true
			***REMOVED***
			if stringobj1 < stringobj2 ***REMOVED***
				return 1, true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return 0, false
***REMOVED***

// Greater asserts that the first element is greater than the second
//
//    assert.Greater(t, 2, 1)
//    assert.Greater(t, float64(2), float64(1))
//    assert.Greater(t, "b", "a")
func Greater(t TestingT, e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***

	e1Kind := reflect.ValueOf(e1).Kind()
	e2Kind := reflect.ValueOf(e2).Kind()
	if e1Kind != e2Kind ***REMOVED***
		return Fail(t, "Elements should be the same type", msgAndArgs...)
	***REMOVED***

	res, isComparable := compare(e1, e2, e1Kind)
	if !isComparable ***REMOVED***
		return Fail(t, fmt.Sprintf("Can not compare type \"%s\"", reflect.TypeOf(e1)), msgAndArgs...)
	***REMOVED***

	if res != -1 ***REMOVED***
		return Fail(t, fmt.Sprintf("\"%v\" is not greater than \"%v\"", e1, e2), msgAndArgs...)
	***REMOVED***

	return true
***REMOVED***

// GreaterOrEqual asserts that the first element is greater than or equal to the second
//
//    assert.GreaterOrEqual(t, 2, 1)
//    assert.GreaterOrEqual(t, 2, 2)
//    assert.GreaterOrEqual(t, "b", "a")
//    assert.GreaterOrEqual(t, "b", "b")
func GreaterOrEqual(t TestingT, e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***

	e1Kind := reflect.ValueOf(e1).Kind()
	e2Kind := reflect.ValueOf(e2).Kind()
	if e1Kind != e2Kind ***REMOVED***
		return Fail(t, "Elements should be the same type", msgAndArgs...)
	***REMOVED***

	res, isComparable := compare(e1, e2, e1Kind)
	if !isComparable ***REMOVED***
		return Fail(t, fmt.Sprintf("Can not compare type \"%s\"", reflect.TypeOf(e1)), msgAndArgs...)
	***REMOVED***

	if res != -1 && res != 0 ***REMOVED***
		return Fail(t, fmt.Sprintf("\"%v\" is not greater than or equal to \"%v\"", e1, e2), msgAndArgs...)
	***REMOVED***

	return true
***REMOVED***

// Less asserts that the first element is less than the second
//
//    assert.Less(t, 1, 2)
//    assert.Less(t, float64(1), float64(2))
//    assert.Less(t, "a", "b")
func Less(t TestingT, e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***

	e1Kind := reflect.ValueOf(e1).Kind()
	e2Kind := reflect.ValueOf(e2).Kind()
	if e1Kind != e2Kind ***REMOVED***
		return Fail(t, "Elements should be the same type", msgAndArgs...)
	***REMOVED***

	res, isComparable := compare(e1, e2, e1Kind)
	if !isComparable ***REMOVED***
		return Fail(t, fmt.Sprintf("Can not compare type \"%s\"", reflect.TypeOf(e1)), msgAndArgs...)
	***REMOVED***

	if res != 1 ***REMOVED***
		return Fail(t, fmt.Sprintf("\"%v\" is not less than \"%v\"", e1, e2), msgAndArgs...)
	***REMOVED***

	return true
***REMOVED***

// LessOrEqual asserts that the first element is less than or equal to the second
//
//    assert.LessOrEqual(t, 1, 2)
//    assert.LessOrEqual(t, 2, 2)
//    assert.LessOrEqual(t, "a", "b")
//    assert.LessOrEqual(t, "b", "b")
func LessOrEqual(t TestingT, e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***

	e1Kind := reflect.ValueOf(e1).Kind()
	e2Kind := reflect.ValueOf(e2).Kind()
	if e1Kind != e2Kind ***REMOVED***
		return Fail(t, "Elements should be the same type", msgAndArgs...)
	***REMOVED***

	res, isComparable := compare(e1, e2, e1Kind)
	if !isComparable ***REMOVED***
		return Fail(t, fmt.Sprintf("Can not compare type \"%s\"", reflect.TypeOf(e1)), msgAndArgs...)
	***REMOVED***

	if res != 1 && res != 0 ***REMOVED***
		return Fail(t, fmt.Sprintf("\"%v\" is not less than or equal to \"%v\"", e1, e2), msgAndArgs...)
	***REMOVED***

	return true
***REMOVED***
