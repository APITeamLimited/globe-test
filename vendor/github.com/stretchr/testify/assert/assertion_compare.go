package assert

import (
	"fmt"
	"reflect"
)

type CompareType int

const (
	compareLess CompareType = iota - 1
	compareEqual
	compareGreater
)

var (
	intType   = reflect.TypeOf(int(1))
	int8Type  = reflect.TypeOf(int8(1))
	int16Type = reflect.TypeOf(int16(1))
	int32Type = reflect.TypeOf(int32(1))
	int64Type = reflect.TypeOf(int64(1))

	uintType   = reflect.TypeOf(uint(1))
	uint8Type  = reflect.TypeOf(uint8(1))
	uint16Type = reflect.TypeOf(uint16(1))
	uint32Type = reflect.TypeOf(uint32(1))
	uint64Type = reflect.TypeOf(uint64(1))

	float32Type = reflect.TypeOf(float32(1))
	float64Type = reflect.TypeOf(float64(1))

	stringType = reflect.TypeOf("")
)

func compare(obj1, obj2 interface***REMOVED******REMOVED***, kind reflect.Kind) (CompareType, bool) ***REMOVED***
	obj1Value := reflect.ValueOf(obj1)
	obj2Value := reflect.ValueOf(obj2)

	// throughout this switch we try and avoid calling .Convert() if possible,
	// as this has a pretty big performance impact
	switch kind ***REMOVED***
	case reflect.Int:
		***REMOVED***
			intobj1, ok := obj1.(int)
			if !ok ***REMOVED***
				intobj1 = obj1Value.Convert(intType).Interface().(int)
			***REMOVED***
			intobj2, ok := obj2.(int)
			if !ok ***REMOVED***
				intobj2 = obj2Value.Convert(intType).Interface().(int)
			***REMOVED***
			if intobj1 > intobj2 ***REMOVED***
				return compareGreater, true
			***REMOVED***
			if intobj1 == intobj2 ***REMOVED***
				return compareEqual, true
			***REMOVED***
			if intobj1 < intobj2 ***REMOVED***
				return compareLess, true
			***REMOVED***
		***REMOVED***
	case reflect.Int8:
		***REMOVED***
			int8obj1, ok := obj1.(int8)
			if !ok ***REMOVED***
				int8obj1 = obj1Value.Convert(int8Type).Interface().(int8)
			***REMOVED***
			int8obj2, ok := obj2.(int8)
			if !ok ***REMOVED***
				int8obj2 = obj2Value.Convert(int8Type).Interface().(int8)
			***REMOVED***
			if int8obj1 > int8obj2 ***REMOVED***
				return compareGreater, true
			***REMOVED***
			if int8obj1 == int8obj2 ***REMOVED***
				return compareEqual, true
			***REMOVED***
			if int8obj1 < int8obj2 ***REMOVED***
				return compareLess, true
			***REMOVED***
		***REMOVED***
	case reflect.Int16:
		***REMOVED***
			int16obj1, ok := obj1.(int16)
			if !ok ***REMOVED***
				int16obj1 = obj1Value.Convert(int16Type).Interface().(int16)
			***REMOVED***
			int16obj2, ok := obj2.(int16)
			if !ok ***REMOVED***
				int16obj2 = obj2Value.Convert(int16Type).Interface().(int16)
			***REMOVED***
			if int16obj1 > int16obj2 ***REMOVED***
				return compareGreater, true
			***REMOVED***
			if int16obj1 == int16obj2 ***REMOVED***
				return compareEqual, true
			***REMOVED***
			if int16obj1 < int16obj2 ***REMOVED***
				return compareLess, true
			***REMOVED***
		***REMOVED***
	case reflect.Int32:
		***REMOVED***
			int32obj1, ok := obj1.(int32)
			if !ok ***REMOVED***
				int32obj1 = obj1Value.Convert(int32Type).Interface().(int32)
			***REMOVED***
			int32obj2, ok := obj2.(int32)
			if !ok ***REMOVED***
				int32obj2 = obj2Value.Convert(int32Type).Interface().(int32)
			***REMOVED***
			if int32obj1 > int32obj2 ***REMOVED***
				return compareGreater, true
			***REMOVED***
			if int32obj1 == int32obj2 ***REMOVED***
				return compareEqual, true
			***REMOVED***
			if int32obj1 < int32obj2 ***REMOVED***
				return compareLess, true
			***REMOVED***
		***REMOVED***
	case reflect.Int64:
		***REMOVED***
			int64obj1, ok := obj1.(int64)
			if !ok ***REMOVED***
				int64obj1 = obj1Value.Convert(int64Type).Interface().(int64)
			***REMOVED***
			int64obj2, ok := obj2.(int64)
			if !ok ***REMOVED***
				int64obj2 = obj2Value.Convert(int64Type).Interface().(int64)
			***REMOVED***
			if int64obj1 > int64obj2 ***REMOVED***
				return compareGreater, true
			***REMOVED***
			if int64obj1 == int64obj2 ***REMOVED***
				return compareEqual, true
			***REMOVED***
			if int64obj1 < int64obj2 ***REMOVED***
				return compareLess, true
			***REMOVED***
		***REMOVED***
	case reflect.Uint:
		***REMOVED***
			uintobj1, ok := obj1.(uint)
			if !ok ***REMOVED***
				uintobj1 = obj1Value.Convert(uintType).Interface().(uint)
			***REMOVED***
			uintobj2, ok := obj2.(uint)
			if !ok ***REMOVED***
				uintobj2 = obj2Value.Convert(uintType).Interface().(uint)
			***REMOVED***
			if uintobj1 > uintobj2 ***REMOVED***
				return compareGreater, true
			***REMOVED***
			if uintobj1 == uintobj2 ***REMOVED***
				return compareEqual, true
			***REMOVED***
			if uintobj1 < uintobj2 ***REMOVED***
				return compareLess, true
			***REMOVED***
		***REMOVED***
	case reflect.Uint8:
		***REMOVED***
			uint8obj1, ok := obj1.(uint8)
			if !ok ***REMOVED***
				uint8obj1 = obj1Value.Convert(uint8Type).Interface().(uint8)
			***REMOVED***
			uint8obj2, ok := obj2.(uint8)
			if !ok ***REMOVED***
				uint8obj2 = obj2Value.Convert(uint8Type).Interface().(uint8)
			***REMOVED***
			if uint8obj1 > uint8obj2 ***REMOVED***
				return compareGreater, true
			***REMOVED***
			if uint8obj1 == uint8obj2 ***REMOVED***
				return compareEqual, true
			***REMOVED***
			if uint8obj1 < uint8obj2 ***REMOVED***
				return compareLess, true
			***REMOVED***
		***REMOVED***
	case reflect.Uint16:
		***REMOVED***
			uint16obj1, ok := obj1.(uint16)
			if !ok ***REMOVED***
				uint16obj1 = obj1Value.Convert(uint16Type).Interface().(uint16)
			***REMOVED***
			uint16obj2, ok := obj2.(uint16)
			if !ok ***REMOVED***
				uint16obj2 = obj2Value.Convert(uint16Type).Interface().(uint16)
			***REMOVED***
			if uint16obj1 > uint16obj2 ***REMOVED***
				return compareGreater, true
			***REMOVED***
			if uint16obj1 == uint16obj2 ***REMOVED***
				return compareEqual, true
			***REMOVED***
			if uint16obj1 < uint16obj2 ***REMOVED***
				return compareLess, true
			***REMOVED***
		***REMOVED***
	case reflect.Uint32:
		***REMOVED***
			uint32obj1, ok := obj1.(uint32)
			if !ok ***REMOVED***
				uint32obj1 = obj1Value.Convert(uint32Type).Interface().(uint32)
			***REMOVED***
			uint32obj2, ok := obj2.(uint32)
			if !ok ***REMOVED***
				uint32obj2 = obj2Value.Convert(uint32Type).Interface().(uint32)
			***REMOVED***
			if uint32obj1 > uint32obj2 ***REMOVED***
				return compareGreater, true
			***REMOVED***
			if uint32obj1 == uint32obj2 ***REMOVED***
				return compareEqual, true
			***REMOVED***
			if uint32obj1 < uint32obj2 ***REMOVED***
				return compareLess, true
			***REMOVED***
		***REMOVED***
	case reflect.Uint64:
		***REMOVED***
			uint64obj1, ok := obj1.(uint64)
			if !ok ***REMOVED***
				uint64obj1 = obj1Value.Convert(uint64Type).Interface().(uint64)
			***REMOVED***
			uint64obj2, ok := obj2.(uint64)
			if !ok ***REMOVED***
				uint64obj2 = obj2Value.Convert(uint64Type).Interface().(uint64)
			***REMOVED***
			if uint64obj1 > uint64obj2 ***REMOVED***
				return compareGreater, true
			***REMOVED***
			if uint64obj1 == uint64obj2 ***REMOVED***
				return compareEqual, true
			***REMOVED***
			if uint64obj1 < uint64obj2 ***REMOVED***
				return compareLess, true
			***REMOVED***
		***REMOVED***
	case reflect.Float32:
		***REMOVED***
			float32obj1, ok := obj1.(float32)
			if !ok ***REMOVED***
				float32obj1 = obj1Value.Convert(float32Type).Interface().(float32)
			***REMOVED***
			float32obj2, ok := obj2.(float32)
			if !ok ***REMOVED***
				float32obj2 = obj2Value.Convert(float32Type).Interface().(float32)
			***REMOVED***
			if float32obj1 > float32obj2 ***REMOVED***
				return compareGreater, true
			***REMOVED***
			if float32obj1 == float32obj2 ***REMOVED***
				return compareEqual, true
			***REMOVED***
			if float32obj1 < float32obj2 ***REMOVED***
				return compareLess, true
			***REMOVED***
		***REMOVED***
	case reflect.Float64:
		***REMOVED***
			float64obj1, ok := obj1.(float64)
			if !ok ***REMOVED***
				float64obj1 = obj1Value.Convert(float64Type).Interface().(float64)
			***REMOVED***
			float64obj2, ok := obj2.(float64)
			if !ok ***REMOVED***
				float64obj2 = obj2Value.Convert(float64Type).Interface().(float64)
			***REMOVED***
			if float64obj1 > float64obj2 ***REMOVED***
				return compareGreater, true
			***REMOVED***
			if float64obj1 == float64obj2 ***REMOVED***
				return compareEqual, true
			***REMOVED***
			if float64obj1 < float64obj2 ***REMOVED***
				return compareLess, true
			***REMOVED***
		***REMOVED***
	case reflect.String:
		***REMOVED***
			stringobj1, ok := obj1.(string)
			if !ok ***REMOVED***
				stringobj1 = obj1Value.Convert(stringType).Interface().(string)
			***REMOVED***
			stringobj2, ok := obj2.(string)
			if !ok ***REMOVED***
				stringobj2 = obj2Value.Convert(stringType).Interface().(string)
			***REMOVED***
			if stringobj1 > stringobj2 ***REMOVED***
				return compareGreater, true
			***REMOVED***
			if stringobj1 == stringobj2 ***REMOVED***
				return compareEqual, true
			***REMOVED***
			if stringobj1 < stringobj2 ***REMOVED***
				return compareLess, true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return compareEqual, false
***REMOVED***

// Greater asserts that the first element is greater than the second
//
//    assert.Greater(t, 2, 1)
//    assert.Greater(t, float64(2), float64(1))
//    assert.Greater(t, "b", "a")
func Greater(t TestingT, e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return compareTwoValues(t, e1, e2, []CompareType***REMOVED***compareGreater***REMOVED***, "\"%v\" is not greater than \"%v\"", msgAndArgs)
***REMOVED***

// GreaterOrEqual asserts that the first element is greater than or equal to the second
//
//    assert.GreaterOrEqual(t, 2, 1)
//    assert.GreaterOrEqual(t, 2, 2)
//    assert.GreaterOrEqual(t, "b", "a")
//    assert.GreaterOrEqual(t, "b", "b")
func GreaterOrEqual(t TestingT, e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return compareTwoValues(t, e1, e2, []CompareType***REMOVED***compareGreater, compareEqual***REMOVED***, "\"%v\" is not greater than or equal to \"%v\"", msgAndArgs)
***REMOVED***

// Less asserts that the first element is less than the second
//
//    assert.Less(t, 1, 2)
//    assert.Less(t, float64(1), float64(2))
//    assert.Less(t, "a", "b")
func Less(t TestingT, e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return compareTwoValues(t, e1, e2, []CompareType***REMOVED***compareLess***REMOVED***, "\"%v\" is not less than \"%v\"", msgAndArgs)
***REMOVED***

// LessOrEqual asserts that the first element is less than or equal to the second
//
//    assert.LessOrEqual(t, 1, 2)
//    assert.LessOrEqual(t, 2, 2)
//    assert.LessOrEqual(t, "a", "b")
//    assert.LessOrEqual(t, "b", "b")
func LessOrEqual(t TestingT, e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	return compareTwoValues(t, e1, e2, []CompareType***REMOVED***compareLess, compareEqual***REMOVED***, "\"%v\" is not less than or equal to \"%v\"", msgAndArgs)
***REMOVED***

// Positive asserts that the specified element is positive
//
//    assert.Positive(t, 1)
//    assert.Positive(t, 1.23)
func Positive(t TestingT, e interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	zero := reflect.Zero(reflect.TypeOf(e))
	return compareTwoValues(t, e, zero.Interface(), []CompareType***REMOVED***compareGreater***REMOVED***, "\"%v\" is not positive", msgAndArgs)
***REMOVED***

// Negative asserts that the specified element is negative
//
//    assert.Negative(t, -1)
//    assert.Negative(t, -1.23)
func Negative(t TestingT, e interface***REMOVED******REMOVED***, msgAndArgs ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	zero := reflect.Zero(reflect.TypeOf(e))
	return compareTwoValues(t, e, zero.Interface(), []CompareType***REMOVED***compareLess***REMOVED***, "\"%v\" is not negative", msgAndArgs)
***REMOVED***

func compareTwoValues(t TestingT, e1 interface***REMOVED******REMOVED***, e2 interface***REMOVED******REMOVED***, allowedComparesResults []CompareType, failMessage string, msgAndArgs ...interface***REMOVED******REMOVED***) bool ***REMOVED***
	if h, ok := t.(tHelper); ok ***REMOVED***
		h.Helper()
	***REMOVED***

	e1Kind := reflect.ValueOf(e1).Kind()
	e2Kind := reflect.ValueOf(e2).Kind()
	if e1Kind != e2Kind ***REMOVED***
		return Fail(t, "Elements should be the same type", msgAndArgs...)
	***REMOVED***

	compareResult, isComparable := compare(e1, e2, e1Kind)
	if !isComparable ***REMOVED***
		return Fail(t, fmt.Sprintf("Can not compare type \"%s\"", reflect.TypeOf(e1)), msgAndArgs...)
	***REMOVED***

	if !containsValue(allowedComparesResults, compareResult) ***REMOVED***
		return Fail(t, fmt.Sprintf(failMessage, e1, e2), msgAndArgs...)
	***REMOVED***

	return true
***REMOVED***

func containsValue(values []CompareType, value CompareType) bool ***REMOVED***
	for _, v := range values ***REMOVED***
		if v == value ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***
