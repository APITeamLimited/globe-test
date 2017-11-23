package validator

import (
	"reflect"
	"strconv"
	"strings"
)

const (
	blank              = ""
	namespaceSeparator = "."
	leftBracket        = "["
	rightBracket       = "]"
	restrictedTagChars = ".[],|=+()`~!@#$%^&*\\\"/?<>***REMOVED******REMOVED***"
	restrictedAliasErr = "Alias '%s' either contains restricted characters or is the same as a restricted tag needed for normal operation"
	restrictedTagErr   = "Tag '%s' either contains restricted characters or is the same as a restricted tag needed for normal operation"
)

var (
	restrictedTags = map[string]struct***REMOVED******REMOVED******REMOVED***
		diveTag:           ***REMOVED******REMOVED***,
		existsTag:         ***REMOVED******REMOVED***,
		structOnlyTag:     ***REMOVED******REMOVED***,
		omitempty:         ***REMOVED******REMOVED***,
		skipValidationTag: ***REMOVED******REMOVED***,
		utf8HexComma:      ***REMOVED******REMOVED***,
		utf8Pipe:          ***REMOVED******REMOVED***,
		noStructLevelTag:  ***REMOVED******REMOVED***,
	***REMOVED***
)

// ExtractType gets the actual underlying type of field value.
// It will dive into pointers, customTypes and return you the
// underlying value and it's kind.
// it is exposed for use within you Custom Functions
func (v *Validate) ExtractType(current reflect.Value) (reflect.Value, reflect.Kind) ***REMOVED***

	val, k, _ := v.extractTypeInternal(current, false)
	return val, k
***REMOVED***

// only exists to not break backward compatibility, needed to return the third param for a bug fix internally
func (v *Validate) extractTypeInternal(current reflect.Value, nullable bool) (reflect.Value, reflect.Kind, bool) ***REMOVED***

	switch current.Kind() ***REMOVED***
	case reflect.Ptr:

		nullable = true

		if current.IsNil() ***REMOVED***
			return current, reflect.Ptr, nullable
		***REMOVED***

		return v.extractTypeInternal(current.Elem(), nullable)

	case reflect.Interface:

		nullable = true

		if current.IsNil() ***REMOVED***
			return current, reflect.Interface, nullable
		***REMOVED***

		return v.extractTypeInternal(current.Elem(), nullable)

	case reflect.Invalid:
		return current, reflect.Invalid, nullable

	default:

		if v.hasCustomFuncs ***REMOVED***

			if fn, ok := v.customTypeFuncs[current.Type()]; ok ***REMOVED***
				return v.extractTypeInternal(reflect.ValueOf(fn(current)), nullable)
			***REMOVED***
		***REMOVED***

		return current, current.Kind(), nullable
	***REMOVED***
***REMOVED***

// GetStructFieldOK traverses a struct to retrieve a specific field denoted by the provided namespace and
// returns the field, field kind and whether is was successful in retrieving the field at all.
// NOTE: when not successful ok will be false, this can happen when a nested struct is nil and so the field
// could not be retrieved because it didn't exist.
func (v *Validate) GetStructFieldOK(current reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool) ***REMOVED***

	current, kind := v.ExtractType(current)

	if kind == reflect.Invalid ***REMOVED***
		return current, kind, false
	***REMOVED***

	if namespace == blank ***REMOVED***
		return current, kind, true
	***REMOVED***

	switch kind ***REMOVED***

	case reflect.Ptr, reflect.Interface:

		return current, kind, false

	case reflect.Struct:

		typ := current.Type()
		fld := namespace
		ns := namespace

		if typ != timeType && typ != timePtrType ***REMOVED***

			idx := strings.Index(namespace, namespaceSeparator)

			if idx != -1 ***REMOVED***
				fld = namespace[:idx]
				ns = namespace[idx+1:]
			***REMOVED*** else ***REMOVED***
				ns = blank
			***REMOVED***

			bracketIdx := strings.Index(fld, leftBracket)
			if bracketIdx != -1 ***REMOVED***
				fld = fld[:bracketIdx]

				ns = namespace[bracketIdx:]
			***REMOVED***

			current = current.FieldByName(fld)

			return v.GetStructFieldOK(current, ns)
		***REMOVED***

	case reflect.Array, reflect.Slice:
		idx := strings.Index(namespace, leftBracket)
		idx2 := strings.Index(namespace, rightBracket)

		arrIdx, _ := strconv.Atoi(namespace[idx+1 : idx2])

		if arrIdx >= current.Len() ***REMOVED***
			return current, kind, false
		***REMOVED***

		startIdx := idx2 + 1

		if startIdx < len(namespace) ***REMOVED***
			if namespace[startIdx:startIdx+1] == namespaceSeparator ***REMOVED***
				startIdx++
			***REMOVED***
		***REMOVED***

		return v.GetStructFieldOK(current.Index(arrIdx), namespace[startIdx:])

	case reflect.Map:
		idx := strings.Index(namespace, leftBracket) + 1
		idx2 := strings.Index(namespace, rightBracket)

		endIdx := idx2

		if endIdx+1 < len(namespace) ***REMOVED***
			if namespace[endIdx+1:endIdx+2] == namespaceSeparator ***REMOVED***
				endIdx++
			***REMOVED***
		***REMOVED***

		key := namespace[idx:idx2]

		switch current.Type().Key().Kind() ***REMOVED***
		case reflect.Int:
			i, _ := strconv.Atoi(key)
			return v.GetStructFieldOK(current.MapIndex(reflect.ValueOf(i)), namespace[endIdx+1:])
		case reflect.Int8:
			i, _ := strconv.ParseInt(key, 10, 8)
			return v.GetStructFieldOK(current.MapIndex(reflect.ValueOf(int8(i))), namespace[endIdx+1:])
		case reflect.Int16:
			i, _ := strconv.ParseInt(key, 10, 16)
			return v.GetStructFieldOK(current.MapIndex(reflect.ValueOf(int16(i))), namespace[endIdx+1:])
		case reflect.Int32:
			i, _ := strconv.ParseInt(key, 10, 32)
			return v.GetStructFieldOK(current.MapIndex(reflect.ValueOf(int32(i))), namespace[endIdx+1:])
		case reflect.Int64:
			i, _ := strconv.ParseInt(key, 10, 64)
			return v.GetStructFieldOK(current.MapIndex(reflect.ValueOf(i)), namespace[endIdx+1:])
		case reflect.Uint:
			i, _ := strconv.ParseUint(key, 10, 0)
			return v.GetStructFieldOK(current.MapIndex(reflect.ValueOf(uint(i))), namespace[endIdx+1:])
		case reflect.Uint8:
			i, _ := strconv.ParseUint(key, 10, 8)
			return v.GetStructFieldOK(current.MapIndex(reflect.ValueOf(uint8(i))), namespace[endIdx+1:])
		case reflect.Uint16:
			i, _ := strconv.ParseUint(key, 10, 16)
			return v.GetStructFieldOK(current.MapIndex(reflect.ValueOf(uint16(i))), namespace[endIdx+1:])
		case reflect.Uint32:
			i, _ := strconv.ParseUint(key, 10, 32)
			return v.GetStructFieldOK(current.MapIndex(reflect.ValueOf(uint32(i))), namespace[endIdx+1:])
		case reflect.Uint64:
			i, _ := strconv.ParseUint(key, 10, 64)
			return v.GetStructFieldOK(current.MapIndex(reflect.ValueOf(i)), namespace[endIdx+1:])
		case reflect.Float32:
			f, _ := strconv.ParseFloat(key, 32)
			return v.GetStructFieldOK(current.MapIndex(reflect.ValueOf(float32(f))), namespace[endIdx+1:])
		case reflect.Float64:
			f, _ := strconv.ParseFloat(key, 64)
			return v.GetStructFieldOK(current.MapIndex(reflect.ValueOf(f)), namespace[endIdx+1:])
		case reflect.Bool:
			b, _ := strconv.ParseBool(key)
			return v.GetStructFieldOK(current.MapIndex(reflect.ValueOf(b)), namespace[endIdx+1:])

		// reflect.Type = string
		default:
			return v.GetStructFieldOK(current.MapIndex(reflect.ValueOf(key)), namespace[endIdx+1:])
		***REMOVED***
	***REMOVED***

	// if got here there was more namespace, cannot go any deeper
	panic("Invalid field namespace")
***REMOVED***

// asInt returns the parameter as a int64
// or panics if it can't convert
func asInt(param string) int64 ***REMOVED***

	i, err := strconv.ParseInt(param, 0, 64)
	panicIf(err)

	return i
***REMOVED***

// asUint returns the parameter as a uint64
// or panics if it can't convert
func asUint(param string) uint64 ***REMOVED***

	i, err := strconv.ParseUint(param, 0, 64)
	panicIf(err)

	return i
***REMOVED***

// asFloat returns the parameter as a float64
// or panics if it can't convert
func asFloat(param string) float64 ***REMOVED***

	i, err := strconv.ParseFloat(param, 64)
	panicIf(err)

	return i
***REMOVED***

func panicIf(err error) ***REMOVED***
	if err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***
