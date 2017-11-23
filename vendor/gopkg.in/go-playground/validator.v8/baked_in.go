package validator

import (
	"fmt"
	"net"
	"net/url"
	"reflect"
	"strings"
	"time"
	"unicode/utf8"
)

// BakedInAliasValidators is a default mapping of a single validationstag that
// defines a common or complex set of validation(s) to simplify
// adding validation to structs. i.e. set key "_ageok" and the tags
// are "gt=0,lte=130" or key "_preferredname" and tags "omitempty,gt=0,lte=60"
var bakedInAliasValidators = map[string]string***REMOVED***
	"iscolor": "hexcolor|rgb|rgba|hsl|hsla",
***REMOVED***

// BakedInValidators is the default map of ValidationFunc
// you can add, remove or even replace items to suite your needs,
// or even disregard and use your own map if so desired.
var bakedInValidators = map[string]Func***REMOVED***
	"required":     HasValue,
	"len":          HasLengthOf,
	"min":          HasMinOf,
	"max":          HasMaxOf,
	"eq":           IsEq,
	"ne":           IsNe,
	"lt":           IsLt,
	"lte":          IsLte,
	"gt":           IsGt,
	"gte":          IsGte,
	"eqfield":      IsEqField,
	"eqcsfield":    IsEqCrossStructField,
	"necsfield":    IsNeCrossStructField,
	"gtcsfield":    IsGtCrossStructField,
	"gtecsfield":   IsGteCrossStructField,
	"ltcsfield":    IsLtCrossStructField,
	"ltecsfield":   IsLteCrossStructField,
	"nefield":      IsNeField,
	"gtefield":     IsGteField,
	"gtfield":      IsGtField,
	"ltefield":     IsLteField,
	"ltfield":      IsLtField,
	"alpha":        IsAlpha,
	"alphanum":     IsAlphanum,
	"numeric":      IsNumeric,
	"number":       IsNumber,
	"hexadecimal":  IsHexadecimal,
	"hexcolor":     IsHEXColor,
	"rgb":          IsRGB,
	"rgba":         IsRGBA,
	"hsl":          IsHSL,
	"hsla":         IsHSLA,
	"email":        IsEmail,
	"url":          IsURL,
	"uri":          IsURI,
	"base64":       IsBase64,
	"contains":     Contains,
	"containsany":  ContainsAny,
	"containsrune": ContainsRune,
	"excludes":     Excludes,
	"excludesall":  ExcludesAll,
	"excludesrune": ExcludesRune,
	"isbn":         IsISBN,
	"isbn10":       IsISBN10,
	"isbn13":       IsISBN13,
	"uuid":         IsUUID,
	"uuid3":        IsUUID3,
	"uuid4":        IsUUID4,
	"uuid5":        IsUUID5,
	"ascii":        IsASCII,
	"printascii":   IsPrintableASCII,
	"multibyte":    HasMultiByteCharacter,
	"datauri":      IsDataURI,
	"latitude":     IsLatitude,
	"longitude":    IsLongitude,
	"ssn":          IsSSN,
	"ipv4":         IsIPv4,
	"ipv6":         IsIPv6,
	"ip":           IsIP,
	"cidrv4":       IsCIDRv4,
	"cidrv6":       IsCIDRv6,
	"cidr":         IsCIDR,
	"tcp4_addr":    IsTCP4AddrResolvable,
	"tcp6_addr":    IsTCP6AddrResolvable,
	"tcp_addr":     IsTCPAddrResolvable,
	"udp4_addr":    IsUDP4AddrResolvable,
	"udp6_addr":    IsUDP6AddrResolvable,
	"udp_addr":     IsUDPAddrResolvable,
	"ip4_addr":     IsIP4AddrResolvable,
	"ip6_addr":     IsIP6AddrResolvable,
	"ip_addr":      IsIPAddrResolvable,
	"unix_addr":    IsUnixAddrResolvable,
	"mac":          IsMAC,
***REMOVED***

// IsMAC is the validation function for validating if the field's value is a valid MAC address.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsMAC(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	_, err := net.ParseMAC(field.String())
	return err == nil
***REMOVED***

// IsCIDRv4 is the validation function for validating if the field's value is a valid v4 CIDR address.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsCIDRv4(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	ip, _, err := net.ParseCIDR(field.String())

	return err == nil && ip.To4() != nil
***REMOVED***

// IsCIDRv6 is the validation function for validating if the field's value is a valid v6 CIDR address.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsCIDRv6(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	ip, _, err := net.ParseCIDR(field.String())

	return err == nil && ip.To4() == nil
***REMOVED***

// IsCIDR is the validation function for validating if the field's value is a valid v4 or v6 CIDR address.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsCIDR(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	_, _, err := net.ParseCIDR(field.String())

	return err == nil
***REMOVED***

// IsIPv4 is the validation function for validating if a value is a valid v4 IP address.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsIPv4(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	ip := net.ParseIP(field.String())

	return ip != nil && ip.To4() != nil
***REMOVED***

// IsIPv6 is the validation function for validating if the field's value is a valid v6 IP address.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsIPv6(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	ip := net.ParseIP(field.String())

	return ip != nil && ip.To4() == nil
***REMOVED***

// IsIP is the validation function for validating if the field's value is a valid v4 or v6 IP address.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsIP(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	ip := net.ParseIP(field.String())

	return ip != nil
***REMOVED***

// IsSSN is the validation function for validating if the field's value is a valid SSN.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsSSN(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	if field.Len() != 11 ***REMOVED***
		return false
	***REMOVED***

	return sSNRegex.MatchString(field.String())
***REMOVED***

// IsLongitude is the validation function for validating if the field's value is a valid longitude coordinate.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsLongitude(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return longitudeRegex.MatchString(field.String())
***REMOVED***

// IsLatitude is the validation function for validating if the field's value is a valid latitude coordinate.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsLatitude(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return latitudeRegex.MatchString(field.String())
***REMOVED***

// IsDataURI is the validation function for validating if the field's value is a valid data URI.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsDataURI(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	uri := strings.SplitN(field.String(), ",", 2)

	if len(uri) != 2 ***REMOVED***
		return false
	***REMOVED***

	if !dataURIRegex.MatchString(uri[0]) ***REMOVED***
		return false
	***REMOVED***

	fld := reflect.ValueOf(uri[1])

	return IsBase64(v, topStruct, currentStructOrField, fld, fld.Type(), fld.Kind(), param)
***REMOVED***

// HasMultiByteCharacter is the validation function for validating if the field's value has a multi byte character.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func HasMultiByteCharacter(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	if field.Len() == 0 ***REMOVED***
		return true
	***REMOVED***

	return multibyteRegex.MatchString(field.String())
***REMOVED***

// IsPrintableASCII is the validation function for validating if the field's value is a valid printable ASCII character.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsPrintableASCII(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return printableASCIIRegex.MatchString(field.String())
***REMOVED***

// IsASCII is the validation function for validating if the field's value is a valid ASCII character.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsASCII(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return aSCIIRegex.MatchString(field.String())
***REMOVED***

// IsUUID5 is the validation function for validating if the field's value is a valid v5 UUID.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsUUID5(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return uUID5Regex.MatchString(field.String())
***REMOVED***

// IsUUID4 is the validation function for validating if the field's value is a valid v4 UUID.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsUUID4(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return uUID4Regex.MatchString(field.String())
***REMOVED***

// IsUUID3 is the validation function for validating if the field's value is a valid v3 UUID.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsUUID3(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return uUID3Regex.MatchString(field.String())
***REMOVED***

// IsUUID is the validation function for validating if the field's value is a valid UUID of any version.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsUUID(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return uUIDRegex.MatchString(field.String())
***REMOVED***

// IsISBN is the validation function for validating if the field's value is a valid v10 or v13 ISBN.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsISBN(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return IsISBN10(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param) || IsISBN13(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param)
***REMOVED***

// IsISBN13 is the validation function for validating if the field's value is a valid v13 ISBN.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsISBN13(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	s := strings.Replace(strings.Replace(field.String(), "-", "", 4), " ", "", 4)

	if !iSBN13Regex.MatchString(s) ***REMOVED***
		return false
	***REMOVED***

	var checksum int32
	var i int32

	factor := []int32***REMOVED***1, 3***REMOVED***

	for i = 0; i < 12; i++ ***REMOVED***
		checksum += factor[i%2] * int32(s[i]-'0')
	***REMOVED***

	return (int32(s[12]-'0'))-((10-(checksum%10))%10) == 0
***REMOVED***

// IsISBN10 is the validation function for validating if the field's value is a valid v10 ISBN.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsISBN10(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	s := strings.Replace(strings.Replace(field.String(), "-", "", 3), " ", "", 3)

	if !iSBN10Regex.MatchString(s) ***REMOVED***
		return false
	***REMOVED***

	var checksum int32
	var i int32

	for i = 0; i < 9; i++ ***REMOVED***
		checksum += (i + 1) * int32(s[i]-'0')
	***REMOVED***

	if s[9] == 'X' ***REMOVED***
		checksum += 10 * 10
	***REMOVED*** else ***REMOVED***
		checksum += 10 * int32(s[9]-'0')
	***REMOVED***

	return checksum%11 == 0
***REMOVED***

// ExcludesRune is the validation function for validating that the field's value does not contain the rune specified within the param.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func ExcludesRune(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return !ContainsRune(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param)
***REMOVED***

// ExcludesAll is the validation function for validating that the field's value does not contain any of the characters specified within the param.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func ExcludesAll(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return !ContainsAny(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param)
***REMOVED***

// Excludes is the validation function for validating that the field's value does not contain the text specified within the param.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func Excludes(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return !Contains(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param)
***REMOVED***

// ContainsRune is the validation function for validating that the field's value contains the rune specified within the param.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func ContainsRune(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	r, _ := utf8.DecodeRuneInString(param)

	return strings.ContainsRune(field.String(), r)
***REMOVED***

// ContainsAny is the validation function for validating that the field's value contains any of the characters specified within the param.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func ContainsAny(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return strings.ContainsAny(field.String(), param)
***REMOVED***

// Contains is the validation function for validating that the field's value contains the text specified within the param.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func Contains(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return strings.Contains(field.String(), param)
***REMOVED***

// IsNeField is the validation function for validating if the current field's value is not equal to the field specified by the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsNeField(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	currentField, currentKind, ok := v.GetStructFieldOK(currentStructOrField, param)

	if !ok || currentKind != fieldKind ***REMOVED***
		return true
	***REMOVED***

	switch fieldKind ***REMOVED***

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() != currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() != currentField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() != currentField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) != int64(currentField.Len())

	case reflect.Struct:

		// Not Same underlying type i.e. struct and time
		if fieldType != currentField.Type() ***REMOVED***
			return true
		***REMOVED***

		if fieldType == timeType ***REMOVED***

			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)

			return !fieldTime.Equal(t)
		***REMOVED***

	***REMOVED***

	// default reflect.String:
	return field.String() != currentField.String()
***REMOVED***

// IsNe is the validation function for validating that the field's value does not equal the provided param value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsNe(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return !IsEq(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param)
***REMOVED***

// IsLteCrossStructField is the validation function for validating if the current field's value is less than or equal to the field, within a separate struct, specified by the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsLteCrossStructField(v *Validate, topStruct reflect.Value, current reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	topField, topKind, ok := v.GetStructFieldOK(topStruct, param)
	if !ok || topKind != fieldKind ***REMOVED***
		return false
	***REMOVED***

	switch fieldKind ***REMOVED***

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() <= topField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() <= topField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() <= topField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) <= int64(topField.Len())

	case reflect.Struct:

		// Not Same underlying type i.e. struct and time
		if fieldType != topField.Type() ***REMOVED***
			return false
		***REMOVED***

		if fieldType == timeType ***REMOVED***

			fieldTime := field.Interface().(time.Time)
			topTime := topField.Interface().(time.Time)

			return fieldTime.Before(topTime) || fieldTime.Equal(topTime)
		***REMOVED***
	***REMOVED***

	// default reflect.String:
	return field.String() <= topField.String()
***REMOVED***

// IsLtCrossStructField is the validation function for validating if the current field's value is less than the field, within a separate struct, specified by the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsLtCrossStructField(v *Validate, topStruct reflect.Value, current reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	topField, topKind, ok := v.GetStructFieldOK(topStruct, param)
	if !ok || topKind != fieldKind ***REMOVED***
		return false
	***REMOVED***

	switch fieldKind ***REMOVED***

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() < topField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() < topField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() < topField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) < int64(topField.Len())

	case reflect.Struct:

		// Not Same underlying type i.e. struct and time
		if fieldType != topField.Type() ***REMOVED***
			return false
		***REMOVED***

		if fieldType == timeType ***REMOVED***

			fieldTime := field.Interface().(time.Time)
			topTime := topField.Interface().(time.Time)

			return fieldTime.Before(topTime)
		***REMOVED***
	***REMOVED***

	// default reflect.String:
	return field.String() < topField.String()
***REMOVED***

// IsGteCrossStructField is the validation function for validating if the current field's value is greater than or equal to the field, within a separate struct, specified by the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsGteCrossStructField(v *Validate, topStruct reflect.Value, current reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	topField, topKind, ok := v.GetStructFieldOK(topStruct, param)
	if !ok || topKind != fieldKind ***REMOVED***
		return false
	***REMOVED***

	switch fieldKind ***REMOVED***

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() >= topField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() >= topField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() >= topField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) >= int64(topField.Len())

	case reflect.Struct:

		// Not Same underlying type i.e. struct and time
		if fieldType != topField.Type() ***REMOVED***
			return false
		***REMOVED***

		if fieldType == timeType ***REMOVED***

			fieldTime := field.Interface().(time.Time)
			topTime := topField.Interface().(time.Time)

			return fieldTime.After(topTime) || fieldTime.Equal(topTime)
		***REMOVED***
	***REMOVED***

	// default reflect.String:
	return field.String() >= topField.String()
***REMOVED***

// IsGtCrossStructField is the validation function for validating if the current field's value is greater than the field, within a separate struct, specified by the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsGtCrossStructField(v *Validate, topStruct reflect.Value, current reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	topField, topKind, ok := v.GetStructFieldOK(topStruct, param)
	if !ok || topKind != fieldKind ***REMOVED***
		return false
	***REMOVED***

	switch fieldKind ***REMOVED***

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() > topField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() > topField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() > topField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) > int64(topField.Len())

	case reflect.Struct:

		// Not Same underlying type i.e. struct and time
		if fieldType != topField.Type() ***REMOVED***
			return false
		***REMOVED***

		if fieldType == timeType ***REMOVED***

			fieldTime := field.Interface().(time.Time)
			topTime := topField.Interface().(time.Time)

			return fieldTime.After(topTime)
		***REMOVED***
	***REMOVED***

	// default reflect.String:
	return field.String() > topField.String()
***REMOVED***

// IsNeCrossStructField is the validation function for validating that the current field's value is not equal to the field, within a separate struct, specified by the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsNeCrossStructField(v *Validate, topStruct reflect.Value, current reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	topField, currentKind, ok := v.GetStructFieldOK(topStruct, param)
	if !ok || currentKind != fieldKind ***REMOVED***
		return true
	***REMOVED***

	switch fieldKind ***REMOVED***

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return topField.Int() != field.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return topField.Uint() != field.Uint()

	case reflect.Float32, reflect.Float64:
		return topField.Float() != field.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(topField.Len()) != int64(field.Len())

	case reflect.Struct:

		// Not Same underlying type i.e. struct and time
		if fieldType != topField.Type() ***REMOVED***
			return true
		***REMOVED***

		if fieldType == timeType ***REMOVED***

			t := field.Interface().(time.Time)
			fieldTime := topField.Interface().(time.Time)

			return !fieldTime.Equal(t)
		***REMOVED***
	***REMOVED***

	// default reflect.String:
	return topField.String() != field.String()
***REMOVED***

// IsEqCrossStructField is the validation function for validating that the current field's value is equal to the field, within a separate struct, specified by the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsEqCrossStructField(v *Validate, topStruct reflect.Value, current reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	topField, topKind, ok := v.GetStructFieldOK(topStruct, param)
	if !ok || topKind != fieldKind ***REMOVED***
		return false
	***REMOVED***

	switch fieldKind ***REMOVED***

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return topField.Int() == field.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return topField.Uint() == field.Uint()

	case reflect.Float32, reflect.Float64:
		return topField.Float() == field.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(topField.Len()) == int64(field.Len())

	case reflect.Struct:

		// Not Same underlying type i.e. struct and time
		if fieldType != topField.Type() ***REMOVED***
			return false
		***REMOVED***

		if fieldType == timeType ***REMOVED***

			t := field.Interface().(time.Time)
			fieldTime := topField.Interface().(time.Time)

			return fieldTime.Equal(t)
		***REMOVED***
	***REMOVED***

	// default reflect.String:
	return topField.String() == field.String()
***REMOVED***

// IsEqField is the validation function for validating if the current field's value is equal to the field specified by the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsEqField(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	currentField, currentKind, ok := v.GetStructFieldOK(currentStructOrField, param)
	if !ok || currentKind != fieldKind ***REMOVED***
		return false
	***REMOVED***

	switch fieldKind ***REMOVED***

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() == currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() == currentField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() == currentField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) == int64(currentField.Len())

	case reflect.Struct:

		// Not Same underlying type i.e. struct and time
		if fieldType != currentField.Type() ***REMOVED***
			return false
		***REMOVED***

		if fieldType == timeType ***REMOVED***

			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)

			return fieldTime.Equal(t)
		***REMOVED***

	***REMOVED***

	// default reflect.String:
	return field.String() == currentField.String()
***REMOVED***

// IsEq is the validation function for validating if the current field's value is equal to the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsEq(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	switch fieldKind ***REMOVED***

	case reflect.String:
		return field.String() == param

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) == p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)

		return field.Int() == p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() == p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)

		return field.Float() == p
	***REMOVED***

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
***REMOVED***

// IsBase64 is the validation function for validating if the current field's value is a valid base 64.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsBase64(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return base64Regex.MatchString(field.String())
***REMOVED***

// IsURI is the validation function for validating if the current field's value is a valid URI.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsURI(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	switch fieldKind ***REMOVED***

	case reflect.String:

		s := field.String()

		// checks needed as of Go 1.6 because of change https://github.com/golang/go/commit/617c93ce740c3c3cc28cdd1a0d712be183d0b328#diff-6c2d018290e298803c0c9419d8739885L195
		// emulate browser and strip the '#' suffix prior to validation. see issue-#237
		if i := strings.Index(s, "#"); i > -1 ***REMOVED***
			s = s[:i]
		***REMOVED***

		if s == blank ***REMOVED***
			return false
		***REMOVED***

		_, err := url.ParseRequestURI(s)

		return err == nil
	***REMOVED***

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
***REMOVED***

// IsURL is the validation function for validating if the current field's value is a valid URL.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsURL(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	switch fieldKind ***REMOVED***

	case reflect.String:

		var i int
		s := field.String()

		// checks needed as of Go 1.6 because of change https://github.com/golang/go/commit/617c93ce740c3c3cc28cdd1a0d712be183d0b328#diff-6c2d018290e298803c0c9419d8739885L195
		// emulate browser and strip the '#' suffix prior to validation. see issue-#237
		if i = strings.Index(s, "#"); i > -1 ***REMOVED***
			s = s[:i]
		***REMOVED***

		if s == blank ***REMOVED***
			return false
		***REMOVED***

		url, err := url.ParseRequestURI(s)

		if err != nil || url.Scheme == blank ***REMOVED***
			return false
		***REMOVED***

		return err == nil
	***REMOVED***

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
***REMOVED***

// IsEmail is the validation function for validating if the current field's value is a valid email address.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsEmail(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return emailRegex.MatchString(field.String())
***REMOVED***

// IsHSLA is the validation function for validating if the current field's value is a valid HSLA color.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsHSLA(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return hslaRegex.MatchString(field.String())
***REMOVED***

// IsHSL is the validation function for validating if the current field's value is a valid HSL color.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsHSL(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return hslRegex.MatchString(field.String())
***REMOVED***

// IsRGBA is the validation function for validating if the current field's value is a valid RGBA color.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsRGBA(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return rgbaRegex.MatchString(field.String())
***REMOVED***

// IsRGB is the validation function for validating if the current field's value is a valid RGB color.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsRGB(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return rgbRegex.MatchString(field.String())
***REMOVED***

// IsHEXColor is the validation function for validating if the current field's value is a valid HEX color.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsHEXColor(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return hexcolorRegex.MatchString(field.String())
***REMOVED***

// IsHexadecimal is the validation function for validating if the current field's value is a valid hexadecimal.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsHexadecimal(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return hexadecimalRegex.MatchString(field.String())
***REMOVED***

// IsNumber is the validation function for validating if the current field's value is a valid number.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsNumber(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return numberRegex.MatchString(field.String())
***REMOVED***

// IsNumeric is the validation function for validating if the current field's value is a valid numeric value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsNumeric(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return numericRegex.MatchString(field.String())
***REMOVED***

// IsAlphanum is the validation function for validating if the current field's value is a valid alphanumeric value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsAlphanum(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return alphaNumericRegex.MatchString(field.String())
***REMOVED***

// IsAlpha is the validation function for validating if the current field's value is a valid alpha value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsAlpha(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return alphaRegex.MatchString(field.String())
***REMOVED***

// HasValue is the validation function for validating if the current field's value is not the default static value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func HasValue(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	switch fieldKind ***REMOVED***
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		return field.IsValid() && field.Interface() != reflect.Zero(fieldType).Interface()
	***REMOVED***
***REMOVED***

// IsGteField is the validation function for validating if the current field's value is greater than or equal to the field specified by the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsGteField(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	currentField, currentKind, ok := v.GetStructFieldOK(currentStructOrField, param)
	if !ok || currentKind != fieldKind ***REMOVED***
		return false
	***REMOVED***

	switch fieldKind ***REMOVED***

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		return field.Int() >= currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:

		return field.Uint() >= currentField.Uint()

	case reflect.Float32, reflect.Float64:

		return field.Float() >= currentField.Float()

	case reflect.Struct:

		// Not Same underlying type i.e. struct and time
		if fieldType != currentField.Type() ***REMOVED***
			return false
		***REMOVED***

		if fieldType == timeType ***REMOVED***

			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)

			return fieldTime.After(t) || fieldTime.Equal(t)
		***REMOVED***
	***REMOVED***

	// default reflect.String
	return len(field.String()) >= len(currentField.String())
***REMOVED***

// IsGtField is the validation function for validating if the current field's value is greater than the field specified by the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsGtField(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	currentField, currentKind, ok := v.GetStructFieldOK(currentStructOrField, param)
	if !ok || currentKind != fieldKind ***REMOVED***
		return false
	***REMOVED***

	switch fieldKind ***REMOVED***

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		return field.Int() > currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:

		return field.Uint() > currentField.Uint()

	case reflect.Float32, reflect.Float64:

		return field.Float() > currentField.Float()

	case reflect.Struct:

		// Not Same underlying type i.e. struct and time
		if fieldType != currentField.Type() ***REMOVED***
			return false
		***REMOVED***

		if fieldType == timeType ***REMOVED***

			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)

			return fieldTime.After(t)
		***REMOVED***
	***REMOVED***

	// default reflect.String
	return len(field.String()) > len(currentField.String())
***REMOVED***

// IsGte is the validation function for validating if the current field's value is greater than or equal to the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsGte(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	switch fieldKind ***REMOVED***

	case reflect.String:
		p := asInt(param)

		return int64(utf8.RuneCountInString(field.String())) >= p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) >= p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)

		return field.Int() >= p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() >= p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)

		return field.Float() >= p

	case reflect.Struct:

		if fieldType == timeType || fieldType == timePtrType ***REMOVED***

			now := time.Now().UTC()
			t := field.Interface().(time.Time)

			return t.After(now) || t.Equal(now)
		***REMOVED***
	***REMOVED***

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
***REMOVED***

// IsGt is the validation function for validating if the current field's value is greater than the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsGt(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	switch fieldKind ***REMOVED***

	case reflect.String:
		p := asInt(param)

		return int64(utf8.RuneCountInString(field.String())) > p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) > p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)

		return field.Int() > p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() > p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)

		return field.Float() > p
	case reflect.Struct:

		if fieldType == timeType || fieldType == timePtrType ***REMOVED***

			return field.Interface().(time.Time).After(time.Now().UTC())
		***REMOVED***
	***REMOVED***

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
***REMOVED***

// HasLengthOf is the validation function for validating if the current field's value is equal to the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func HasLengthOf(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	switch fieldKind ***REMOVED***

	case reflect.String:
		p := asInt(param)

		return int64(utf8.RuneCountInString(field.String())) == p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) == p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)

		return field.Int() == p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() == p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)

		return field.Float() == p
	***REMOVED***

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
***REMOVED***

// HasMinOf is the validation function for validating if the current field's value is greater than or equal to the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func HasMinOf(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	return IsGte(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param)
***REMOVED***

// IsLteField is the validation function for validating if the current field's value is less than or equal to the field specified by the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsLteField(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	currentField, currentKind, ok := v.GetStructFieldOK(currentStructOrField, param)
	if !ok || currentKind != fieldKind ***REMOVED***
		return false
	***REMOVED***

	switch fieldKind ***REMOVED***

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		return field.Int() <= currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:

		return field.Uint() <= currentField.Uint()

	case reflect.Float32, reflect.Float64:

		return field.Float() <= currentField.Float()

	case reflect.Struct:

		// Not Same underlying type i.e. struct and time
		if fieldType != currentField.Type() ***REMOVED***
			return false
		***REMOVED***

		if fieldType == timeType ***REMOVED***

			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)

			return fieldTime.Before(t) || fieldTime.Equal(t)
		***REMOVED***
	***REMOVED***

	// default reflect.String
	return len(field.String()) <= len(currentField.String())
***REMOVED***

// IsLtField is the validation function for validating if the current field's value is less than the field specified by the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsLtField(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	currentField, currentKind, ok := v.GetStructFieldOK(currentStructOrField, param)
	if !ok || currentKind != fieldKind ***REMOVED***
		return false
	***REMOVED***

	switch fieldKind ***REMOVED***

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		return field.Int() < currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:

		return field.Uint() < currentField.Uint()

	case reflect.Float32, reflect.Float64:

		return field.Float() < currentField.Float()

	case reflect.Struct:

		// Not Same underlying type i.e. struct and time
		if fieldType != currentField.Type() ***REMOVED***
			return false
		***REMOVED***

		if fieldType == timeType ***REMOVED***

			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)

			return fieldTime.Before(t)
		***REMOVED***
	***REMOVED***

	// default reflect.String
	return len(field.String()) < len(currentField.String())
***REMOVED***

// IsLte is the validation function for validating if the current field's value is less than or equal to the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsLte(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	switch fieldKind ***REMOVED***

	case reflect.String:
		p := asInt(param)

		return int64(utf8.RuneCountInString(field.String())) <= p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) <= p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)

		return field.Int() <= p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() <= p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)

		return field.Float() <= p

	case reflect.Struct:

		if fieldType == timeType || fieldType == timePtrType ***REMOVED***

			now := time.Now().UTC()
			t := field.Interface().(time.Time)

			return t.Before(now) || t.Equal(now)
		***REMOVED***
	***REMOVED***

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
***REMOVED***

// IsLt is the validation function for validating if the current field's value is less than the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsLt(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	switch fieldKind ***REMOVED***

	case reflect.String:
		p := asInt(param)

		return int64(utf8.RuneCountInString(field.String())) < p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) < p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)

		return field.Int() < p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() < p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)

		return field.Float() < p

	case reflect.Struct:

		if fieldType == timeType || fieldType == timePtrType ***REMOVED***

			return field.Interface().(time.Time).Before(time.Now().UTC())
		***REMOVED***
	***REMOVED***

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
***REMOVED***

// HasMaxOf is the validation function for validating if the current field's value is less than or equal to the param's value.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func HasMaxOf(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	return IsLte(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param)
***REMOVED***

// IsTCP4AddrResolvable is the validation function for validating if the field's value is a resolvable tcp4 address.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsTCP4AddrResolvable(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	if !isIP4Addr(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param) ***REMOVED***
		return false
	***REMOVED***

	_, err := net.ResolveTCPAddr("tcp4", field.String())
	return err == nil
***REMOVED***

// IsTCP6AddrResolvable is the validation function for validating if the field's value is a resolvable tcp6 address.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsTCP6AddrResolvable(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	if !isIP6Addr(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param) ***REMOVED***
		return false
	***REMOVED***

	_, err := net.ResolveTCPAddr("tcp6", field.String())
	return err == nil
***REMOVED***

// IsTCPAddrResolvable is the validation function for validating if the field's value is a resolvable tcp address.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsTCPAddrResolvable(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	if !isIP4Addr(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param) &&
		!isIP6Addr(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param) ***REMOVED***
		return false
	***REMOVED***

	_, err := net.ResolveTCPAddr("tcp", field.String())
	return err == nil
***REMOVED***

// IsUDP4AddrResolvable is the validation function for validating if the field's value is a resolvable udp4 address.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsUDP4AddrResolvable(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	if !isIP4Addr(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param) ***REMOVED***
		return false
	***REMOVED***

	_, err := net.ResolveUDPAddr("udp4", field.String())
	return err == nil
***REMOVED***

// IsUDP6AddrResolvable is the validation function for validating if the field's value is a resolvable udp6 address.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsUDP6AddrResolvable(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	if !isIP6Addr(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param) ***REMOVED***
		return false
	***REMOVED***

	_, err := net.ResolveUDPAddr("udp6", field.String())
	return err == nil
***REMOVED***

// IsUDPAddrResolvable is the validation function for validating if the field's value is a resolvable udp address.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsUDPAddrResolvable(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	if !isIP4Addr(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param) &&
		!isIP6Addr(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param) ***REMOVED***
		return false
	***REMOVED***

	_, err := net.ResolveUDPAddr("udp", field.String())
	return err == nil
***REMOVED***

// IsIP4AddrResolvable is the validation function for validating if the field's value is a resolvable ip4 address.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsIP4AddrResolvable(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	if !IsIPv4(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param) ***REMOVED***
		return false
	***REMOVED***

	_, err := net.ResolveIPAddr("ip4", field.String())
	return err == nil
***REMOVED***

// IsIP6AddrResolvable is the validation function for validating if the field's value is a resolvable ip6 address.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsIP6AddrResolvable(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	if !IsIPv6(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param) ***REMOVED***
		return false
	***REMOVED***

	_, err := net.ResolveIPAddr("ip6", field.String())
	return err == nil
***REMOVED***

// IsIPAddrResolvable is the validation function for validating if the field's value is a resolvable ip address.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsIPAddrResolvable(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***

	if !IsIP(v, topStruct, currentStructOrField, field, fieldType, fieldKind, param) ***REMOVED***
		return false
	***REMOVED***

	_, err := net.ResolveIPAddr("ip", field.String())
	return err == nil
***REMOVED***

// IsUnixAddrResolvable is the validation function for validating if the field's value is a resolvable unix address.
// NOTE: This is exposed for use within your own custom functions and not intended to be called directly.
func IsUnixAddrResolvable(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	_, err := net.ResolveUnixAddr("unix", field.String())
	return err == nil
***REMOVED***

func isIP4Addr(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	val := field.String()

	if idx := strings.LastIndex(val, ":"); idx != -1 ***REMOVED***
		val = val[0:idx]
	***REMOVED***

	if !IsIPv4(v, topStruct, currentStructOrField, reflect.ValueOf(val), fieldType, fieldKind, param) ***REMOVED***
		return false
	***REMOVED***

	return true
***REMOVED***

func isIP6Addr(v *Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool ***REMOVED***
	val := field.String()

	if idx := strings.LastIndex(val, ":"); idx != -1 ***REMOVED***
		if idx != 0 && val[idx-1:idx] == "]" ***REMOVED***
			val = val[1 : idx-1]
		***REMOVED***
	***REMOVED***

	if !IsIPv6(v, topStruct, currentStructOrField, reflect.ValueOf(val), fieldType, fieldKind, param) ***REMOVED***
		return false
	***REMOVED***

	return true
***REMOVED***
