/**
 * Package validator
 *
 * MISC:
 * - anonymous structs - they don't have names so expect the Struct name within StructErrors to be blank
 *
 */

package validator

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
)

const (
	utf8HexComma            = "0x2C"
	utf8Pipe                = "0x7C"
	tagSeparator            = ","
	orSeparator             = "|"
	tagKeySeparator         = "="
	structOnlyTag           = "structonly"
	noStructLevelTag        = "nostructlevel"
	omitempty               = "omitempty"
	skipValidationTag       = "-"
	diveTag                 = "dive"
	existsTag               = "exists"
	fieldErrMsg             = "Key: '%s' Error:Field validation for '%s' failed on the '%s' tag"
	arrayIndexFieldName     = "%s" + leftBracket + "%d" + rightBracket
	mapIndexFieldName       = "%s" + leftBracket + "%v" + rightBracket
	invalidValidation       = "Invalid validation tag on field %s"
	undefinedValidation     = "Undefined validation function on field %s"
	validatorNotInitialized = "Validator instance not initialized"
	fieldNameRequired       = "Field Name Required"
	tagRequired             = "Tag Required"
)

var (
	timeType      = reflect.TypeOf(time.Time***REMOVED******REMOVED***)
	timePtrType   = reflect.TypeOf(&time.Time***REMOVED******REMOVED***)
	defaultCField = new(cField)
)

// StructLevel contains all of the information and helper methods
// for reporting errors during struct level validation
type StructLevel struct ***REMOVED***
	TopStruct     reflect.Value
	CurrentStruct reflect.Value
	errPrefix     string
	nsPrefix      string
	errs          ValidationErrors
	v             *Validate
***REMOVED***

// ReportValidationErrors accepts the key relative to the top level struct and validatin errors.
// Example: had a triple nested struct User, ContactInfo, Country and ran errs := validate.Struct(country)
// from within a User struct level validation would call this method like so:
// ReportValidationErrors("ContactInfo.", errs)
// NOTE: relativeKey can contain both the Field Relative and Custom name relative paths
// i.e. ReportValidationErrors("ContactInfo.|cInfo", errs) where cInfo represents say the JSON name of
// the relative path; this will be split into 2 variables in the next valiator version.
func (sl *StructLevel) ReportValidationErrors(relativeKey string, errs ValidationErrors) ***REMOVED***
	for _, e := range errs ***REMOVED***

		idx := strings.Index(relativeKey, "|")
		var rel string
		var cRel string

		if idx != -1 ***REMOVED***
			rel = relativeKey[:idx]
			cRel = relativeKey[idx+1:]
		***REMOVED*** else ***REMOVED***
			rel = relativeKey
		***REMOVED***

		key := sl.errPrefix + rel + e.Field

		e.FieldNamespace = key
		e.NameNamespace = sl.nsPrefix + cRel + e.Name

		sl.errs[key] = e
	***REMOVED***
***REMOVED***

// ReportError reports an error just by passing the field and tag information
// NOTE: tag can be an existing validation tag or just something you make up
// and precess on the flip side it's up to you.
func (sl *StructLevel) ReportError(field reflect.Value, fieldName string, customName string, tag string) ***REMOVED***

	field, kind := sl.v.ExtractType(field)

	if fieldName == blank ***REMOVED***
		panic(fieldNameRequired)
	***REMOVED***

	if customName == blank ***REMOVED***
		customName = fieldName
	***REMOVED***

	if tag == blank ***REMOVED***
		panic(tagRequired)
	***REMOVED***

	ns := sl.errPrefix + fieldName

	switch kind ***REMOVED***
	case reflect.Invalid:
		sl.errs[ns] = &FieldError***REMOVED***
			FieldNamespace: ns,
			NameNamespace:  sl.nsPrefix + customName,
			Name:           customName,
			Field:          fieldName,
			Tag:            tag,
			ActualTag:      tag,
			Param:          blank,
			Kind:           kind,
		***REMOVED***
	default:
		sl.errs[ns] = &FieldError***REMOVED***
			FieldNamespace: ns,
			NameNamespace:  sl.nsPrefix + customName,
			Name:           customName,
			Field:          fieldName,
			Tag:            tag,
			ActualTag:      tag,
			Param:          blank,
			Value:          field.Interface(),
			Kind:           kind,
			Type:           field.Type(),
		***REMOVED***
	***REMOVED***
***REMOVED***

// Validate contains the validator settings passed in using the Config struct
type Validate struct ***REMOVED***
	tagName             string
	fieldNameTag        string
	validationFuncs     map[string]Func
	structLevelFuncs    map[reflect.Type]StructLevelFunc
	customTypeFuncs     map[reflect.Type]CustomTypeFunc
	aliasValidators     map[string]string
	hasCustomFuncs      bool
	hasAliasValidators  bool
	hasStructLevelFuncs bool
	tagCache            *tagCache
	structCache         *structCache
	errsPool            *sync.Pool
***REMOVED***

func (v *Validate) initCheck() ***REMOVED***
	if v == nil ***REMOVED***
		panic(validatorNotInitialized)
	***REMOVED***
***REMOVED***

// Config contains the options that a Validator instance will use.
// It is passed to the New() function
type Config struct ***REMOVED***
	TagName      string
	FieldNameTag string
***REMOVED***

// CustomTypeFunc allows for overriding or adding custom field type handler functions
// field = field value of the type to return a value to be validated
// example Valuer from sql drive see https://golang.org/src/database/sql/driver/types.go?s=1210:1293#L29
type CustomTypeFunc func(field reflect.Value) interface***REMOVED******REMOVED***

// Func accepts all values needed for file and cross field validation
// v             = validator instance, needed but some built in functions for it's custom types
// topStruct     = top level struct when validating by struct otherwise nil
// currentStruct = current level struct when validating by struct otherwise optional comparison value
// field         = field value for validation
// param         = parameter used in validation i.e. gt=0 param would be 0
type Func func(v *Validate, topStruct reflect.Value, currentStruct reflect.Value, field reflect.Value, fieldtype reflect.Type, fieldKind reflect.Kind, param string) bool

// StructLevelFunc accepts all values needed for struct level validation
type StructLevelFunc func(v *Validate, structLevel *StructLevel)

// ValidationErrors is a type of map[string]*FieldError
// it exists to allow for multiple errors to be passed from this library
// and yet still subscribe to the error interface
type ValidationErrors map[string]*FieldError

// Error is intended for use in development + debugging and not intended to be a production error message.
// It allows ValidationErrors to subscribe to the Error interface.
// All information to create an error message specific to your application is contained within
// the FieldError found within the ValidationErrors map
func (ve ValidationErrors) Error() string ***REMOVED***

	buff := bytes.NewBufferString(blank)

	for key, err := range ve ***REMOVED***
		buff.WriteString(fmt.Sprintf(fieldErrMsg, key, err.Field, err.Tag))
		buff.WriteString("\n")
	***REMOVED***

	return strings.TrimSpace(buff.String())
***REMOVED***

// FieldError contains a single field's validation error along
// with other properties that may be needed for error message creation
type FieldError struct ***REMOVED***
	FieldNamespace string
	NameNamespace  string
	Field          string
	Name           string
	Tag            string
	ActualTag      string
	Kind           reflect.Kind
	Type           reflect.Type
	Param          string
	Value          interface***REMOVED******REMOVED***
***REMOVED***

// New creates a new Validate instance for use.
func New(config *Config) *Validate ***REMOVED***

	tc := new(tagCache)
	tc.m.Store(make(map[string]*cTag))

	sc := new(structCache)
	sc.m.Store(make(map[reflect.Type]*cStruct))

	v := &Validate***REMOVED***
		tagName:      config.TagName,
		fieldNameTag: config.FieldNameTag,
		tagCache:     tc,
		structCache:  sc,
		errsPool: &sync.Pool***REMOVED***New: func() interface***REMOVED******REMOVED*** ***REMOVED***
			return ValidationErrors***REMOVED******REMOVED***
		***REMOVED******REMOVED******REMOVED***

	if len(v.aliasValidators) == 0 ***REMOVED***
		// must copy alias validators for separate validations to be used in each validator instance
		v.aliasValidators = map[string]string***REMOVED******REMOVED***
		for k, val := range bakedInAliasValidators ***REMOVED***
			v.RegisterAliasValidation(k, val)
		***REMOVED***
	***REMOVED***

	if len(v.validationFuncs) == 0 ***REMOVED***
		// must copy validators for separate validations to be used in each instance
		v.validationFuncs = map[string]Func***REMOVED******REMOVED***
		for k, val := range bakedInValidators ***REMOVED***
			v.RegisterValidation(k, val)
		***REMOVED***
	***REMOVED***

	return v
***REMOVED***

// RegisterStructValidation registers a StructLevelFunc against a number of types
// NOTE: this method is not thread-safe it is intended that these all be registered prior to any validation
func (v *Validate) RegisterStructValidation(fn StructLevelFunc, types ...interface***REMOVED******REMOVED***) ***REMOVED***
	v.initCheck()

	if v.structLevelFuncs == nil ***REMOVED***
		v.structLevelFuncs = map[reflect.Type]StructLevelFunc***REMOVED******REMOVED***
	***REMOVED***

	for _, t := range types ***REMOVED***
		v.structLevelFuncs[reflect.TypeOf(t)] = fn
	***REMOVED***

	v.hasStructLevelFuncs = true
***REMOVED***

// RegisterValidation adds a validation Func to a Validate's map of validators denoted by the key
// NOTE: if the key already exists, the previous validation function will be replaced.
// NOTE: this method is not thread-safe it is intended that these all be registered prior to any validation
func (v *Validate) RegisterValidation(key string, fn Func) error ***REMOVED***
	v.initCheck()

	if key == blank ***REMOVED***
		return errors.New("Function Key cannot be empty")
	***REMOVED***

	if fn == nil ***REMOVED***
		return errors.New("Function cannot be empty")
	***REMOVED***

	_, ok := restrictedTags[key]

	if ok || strings.ContainsAny(key, restrictedTagChars) ***REMOVED***
		panic(fmt.Sprintf(restrictedTagErr, key))
	***REMOVED***

	v.validationFuncs[key] = fn

	return nil
***REMOVED***

// RegisterCustomTypeFunc registers a CustomTypeFunc against a number of types
// NOTE: this method is not thread-safe it is intended that these all be registered prior to any validation
func (v *Validate) RegisterCustomTypeFunc(fn CustomTypeFunc, types ...interface***REMOVED******REMOVED***) ***REMOVED***
	v.initCheck()

	if v.customTypeFuncs == nil ***REMOVED***
		v.customTypeFuncs = map[reflect.Type]CustomTypeFunc***REMOVED******REMOVED***
	***REMOVED***

	for _, t := range types ***REMOVED***
		v.customTypeFuncs[reflect.TypeOf(t)] = fn
	***REMOVED***

	v.hasCustomFuncs = true
***REMOVED***

// RegisterAliasValidation registers a mapping of a single validationstag that
// defines a common or complex set of validation(s) to simplify adding validation
// to structs. NOTE: when returning an error the tag returned in FieldError will be
// the alias tag unless the dive tag is part of the alias; everything after the
// dive tag is not reported as the alias tag. Also the ActualTag in the before case
// will be the actual tag within the alias that failed.
// NOTE: this method is not thread-safe it is intended that these all be registered prior to any validation
func (v *Validate) RegisterAliasValidation(alias, tags string) ***REMOVED***
	v.initCheck()

	_, ok := restrictedTags[alias]

	if ok || strings.ContainsAny(alias, restrictedTagChars) ***REMOVED***
		panic(fmt.Sprintf(restrictedAliasErr, alias))
	***REMOVED***

	v.aliasValidators[alias] = tags
	v.hasAliasValidators = true
***REMOVED***

// Field validates a single field using tag style validation and returns nil or ValidationErrors as type error.
// You will need to assert the error if it's not nil i.e. err.(validator.ValidationErrors) to access the map of errors.
// NOTE: it returns ValidationErrors instead of a single FieldError because this can also
// validate Array, Slice and maps fields which may contain more than one error
func (v *Validate) Field(field interface***REMOVED******REMOVED***, tag string) error ***REMOVED***
	v.initCheck()

	if len(tag) == 0 || tag == skipValidationTag ***REMOVED***
		return nil
	***REMOVED***

	errs := v.errsPool.Get().(ValidationErrors)
	fieldVal := reflect.ValueOf(field)

	ctag, ok := v.tagCache.Get(tag)
	if !ok ***REMOVED***
		v.tagCache.lock.Lock()
		defer v.tagCache.lock.Unlock()

		// could have been multiple trying to access, but once first is done this ensures tag
		// isn't parsed again.
		ctag, ok = v.tagCache.Get(tag)
		if !ok ***REMOVED***
			ctag, _ = v.parseFieldTagsRecursive(tag, blank, blank, false)
			v.tagCache.Set(tag, ctag)
		***REMOVED***
	***REMOVED***

	v.traverseField(fieldVal, fieldVal, fieldVal, blank, blank, errs, false, false, nil, nil, defaultCField, ctag)

	if len(errs) == 0 ***REMOVED***
		v.errsPool.Put(errs)
		return nil
	***REMOVED***

	return errs
***REMOVED***

// FieldWithValue validates a single field, against another fields value using tag style validation and returns nil or ValidationErrors.
// You will need to assert the error if it's not nil i.e. err.(validator.ValidationErrors) to access the map of errors.
// NOTE: it returns ValidationErrors instead of a single FieldError because this can also
// validate Array, Slice and maps fields which may contain more than one error
func (v *Validate) FieldWithValue(val interface***REMOVED******REMOVED***, field interface***REMOVED******REMOVED***, tag string) error ***REMOVED***
	v.initCheck()

	if len(tag) == 0 || tag == skipValidationTag ***REMOVED***
		return nil
	***REMOVED***

	errs := v.errsPool.Get().(ValidationErrors)
	topVal := reflect.ValueOf(val)

	ctag, ok := v.tagCache.Get(tag)
	if !ok ***REMOVED***
		v.tagCache.lock.Lock()
		defer v.tagCache.lock.Unlock()

		// could have been multiple trying to access, but once first is done this ensures tag
		// isn't parsed again.
		ctag, ok = v.tagCache.Get(tag)
		if !ok ***REMOVED***
			ctag, _ = v.parseFieldTagsRecursive(tag, blank, blank, false)
			v.tagCache.Set(tag, ctag)
		***REMOVED***
	***REMOVED***

	v.traverseField(topVal, topVal, reflect.ValueOf(field), blank, blank, errs, false, false, nil, nil, defaultCField, ctag)

	if len(errs) == 0 ***REMOVED***
		v.errsPool.Put(errs)
		return nil
	***REMOVED***

	return errs
***REMOVED***

// StructPartial validates the fields passed in only, ignoring all others.
// Fields may be provided in a namespaced fashion relative to the  struct provided
// i.e. NestedStruct.Field or NestedArrayField[0].Struct.Name and returns nil or ValidationErrors as error
// You will need to assert the error if it's not nil i.e. err.(validator.ValidationErrors) to access the map of errors.
func (v *Validate) StructPartial(current interface***REMOVED******REMOVED***, fields ...string) error ***REMOVED***
	v.initCheck()

	sv, _ := v.ExtractType(reflect.ValueOf(current))
	name := sv.Type().Name()
	m := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***

	if fields != nil ***REMOVED***
		for _, k := range fields ***REMOVED***

			flds := strings.Split(k, namespaceSeparator)
			if len(flds) > 0 ***REMOVED***

				key := name + namespaceSeparator
				for _, s := range flds ***REMOVED***

					idx := strings.Index(s, leftBracket)

					if idx != -1 ***REMOVED***
						for idx != -1 ***REMOVED***
							key += s[:idx]
							m[key] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

							idx2 := strings.Index(s, rightBracket)
							idx2++
							key += s[idx:idx2]
							m[key] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
							s = s[idx2:]
							idx = strings.Index(s, leftBracket)
						***REMOVED***
					***REMOVED*** else ***REMOVED***

						key += s
						m[key] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
					***REMOVED***

					key += namespaceSeparator
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	errs := v.errsPool.Get().(ValidationErrors)

	v.ensureValidStruct(sv, sv, sv, blank, blank, errs, true, len(m) != 0, false, m, false)

	if len(errs) == 0 ***REMOVED***
		v.errsPool.Put(errs)
		return nil
	***REMOVED***

	return errs
***REMOVED***

// StructExcept validates all fields except the ones passed in.
// Fields may be provided in a namespaced fashion relative to the  struct provided
// i.e. NestedStruct.Field or NestedArrayField[0].Struct.Name and returns nil or ValidationErrors as error
// You will need to assert the error if it's not nil i.e. err.(validator.ValidationErrors) to access the map of errors.
func (v *Validate) StructExcept(current interface***REMOVED******REMOVED***, fields ...string) error ***REMOVED***
	v.initCheck()

	sv, _ := v.ExtractType(reflect.ValueOf(current))
	name := sv.Type().Name()
	m := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***

	for _, key := range fields ***REMOVED***
		m[name+namespaceSeparator+key] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	errs := v.errsPool.Get().(ValidationErrors)

	v.ensureValidStruct(sv, sv, sv, blank, blank, errs, true, len(m) != 0, true, m, false)

	if len(errs) == 0 ***REMOVED***
		v.errsPool.Put(errs)
		return nil
	***REMOVED***

	return errs
***REMOVED***

// Struct validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified.
// it returns nil or ValidationErrors as error.
// You will need to assert the error if it's not nil i.e. err.(validator.ValidationErrors) to access the map of errors.
func (v *Validate) Struct(current interface***REMOVED******REMOVED***) error ***REMOVED***
	v.initCheck()

	errs := v.errsPool.Get().(ValidationErrors)
	sv := reflect.ValueOf(current)

	v.ensureValidStruct(sv, sv, sv, blank, blank, errs, true, false, false, nil, false)

	if len(errs) == 0 ***REMOVED***
		v.errsPool.Put(errs)
		return nil
	***REMOVED***

	return errs
***REMOVED***

func (v *Validate) ensureValidStruct(topStruct reflect.Value, currentStruct reflect.Value, current reflect.Value, errPrefix string, nsPrefix string, errs ValidationErrors, useStructName bool, partial bool, exclude bool, includeExclude map[string]struct***REMOVED******REMOVED***, isStructOnly bool) ***REMOVED***

	if current.Kind() == reflect.Ptr && !current.IsNil() ***REMOVED***
		current = current.Elem()
	***REMOVED***

	if current.Kind() != reflect.Struct && current.Kind() != reflect.Interface ***REMOVED***
		panic("value passed for validation is not a struct")
	***REMOVED***

	v.tranverseStruct(topStruct, currentStruct, current, errPrefix, nsPrefix, errs, useStructName, partial, exclude, includeExclude, nil, nil)
***REMOVED***

// tranverseStruct traverses a structs fields and then passes them to be validated by traverseField
func (v *Validate) tranverseStruct(topStruct reflect.Value, currentStruct reflect.Value, current reflect.Value, errPrefix string, nsPrefix string, errs ValidationErrors, useStructName bool, partial bool, exclude bool, includeExclude map[string]struct***REMOVED******REMOVED***, cs *cStruct, ct *cTag) ***REMOVED***

	var ok bool
	first := len(nsPrefix) == 0
	typ := current.Type()

	cs, ok = v.structCache.Get(typ)
	if !ok ***REMOVED***
		cs = v.extractStructCache(current, typ.Name())
	***REMOVED***

	if useStructName ***REMOVED***
		errPrefix += cs.Name + namespaceSeparator

		if len(v.fieldNameTag) != 0 ***REMOVED***
			nsPrefix += cs.Name + namespaceSeparator
		***REMOVED***
	***REMOVED***

	// structonly tag present don't tranverseFields
	// but must still check and run below struct level validation
	// if present
	if first || ct == nil || ct.typeof != typeStructOnly ***REMOVED***

		for _, f := range cs.fields ***REMOVED***

			if partial ***REMOVED***

				_, ok = includeExclude[errPrefix+f.Name]

				if (ok && exclude) || (!ok && !exclude) ***REMOVED***
					continue
				***REMOVED***
			***REMOVED***

			v.traverseField(topStruct, currentStruct, current.Field(f.Idx), errPrefix, nsPrefix, errs, partial, exclude, includeExclude, cs, f, f.cTags)
		***REMOVED***
	***REMOVED***

	// check if any struct level validations, after all field validations already checked.
	if cs.fn != nil ***REMOVED***
		cs.fn(v, &StructLevel***REMOVED***v: v, TopStruct: topStruct, CurrentStruct: current, errPrefix: errPrefix, nsPrefix: nsPrefix, errs: errs***REMOVED***)
	***REMOVED***
***REMOVED***

// traverseField validates any field, be it a struct or single field, ensures it's validity and passes it along to be validated via it's tag options
func (v *Validate) traverseField(topStruct reflect.Value, currentStruct reflect.Value, current reflect.Value, errPrefix string, nsPrefix string, errs ValidationErrors, partial bool, exclude bool, includeExclude map[string]struct***REMOVED******REMOVED***, cs *cStruct, cf *cField, ct *cTag) ***REMOVED***

	current, kind, nullable := v.extractTypeInternal(current, false)
	var typ reflect.Type

	switch kind ***REMOVED***
	case reflect.Ptr, reflect.Interface, reflect.Invalid:

		if ct == nil ***REMOVED***
			return
		***REMOVED***

		if ct.typeof == typeOmitEmpty ***REMOVED***
			return
		***REMOVED***

		if ct.hasTag ***REMOVED***

			ns := errPrefix + cf.Name

			if kind == reflect.Invalid ***REMOVED***
				errs[ns] = &FieldError***REMOVED***
					FieldNamespace: ns,
					NameNamespace:  nsPrefix + cf.AltName,
					Name:           cf.AltName,
					Field:          cf.Name,
					Tag:            ct.aliasTag,
					ActualTag:      ct.tag,
					Param:          ct.param,
					Kind:           kind,
				***REMOVED***
				return
			***REMOVED***

			errs[ns] = &FieldError***REMOVED***
				FieldNamespace: ns,
				NameNamespace:  nsPrefix + cf.AltName,
				Name:           cf.AltName,
				Field:          cf.Name,
				Tag:            ct.aliasTag,
				ActualTag:      ct.tag,
				Param:          ct.param,
				Value:          current.Interface(),
				Kind:           kind,
				Type:           current.Type(),
			***REMOVED***

			return
		***REMOVED***

	case reflect.Struct:
		typ = current.Type()

		if typ != timeType ***REMOVED***

			if ct != nil ***REMOVED***
				ct = ct.next
			***REMOVED***

			if ct != nil && ct.typeof == typeNoStructLevel ***REMOVED***
				return
			***REMOVED***

			v.tranverseStruct(topStruct, current, current, errPrefix+cf.Name+namespaceSeparator, nsPrefix+cf.AltName+namespaceSeparator, errs, false, partial, exclude, includeExclude, cs, ct)
			return
		***REMOVED***
	***REMOVED***

	if !ct.hasTag ***REMOVED***
		return
	***REMOVED***

	typ = current.Type()

OUTER:
	for ***REMOVED***
		if ct == nil ***REMOVED***
			return
		***REMOVED***

		switch ct.typeof ***REMOVED***

		case typeExists:
			ct = ct.next
			continue

		case typeOmitEmpty:

			if !nullable && !HasValue(v, topStruct, currentStruct, current, typ, kind, blank) ***REMOVED***
				return
			***REMOVED***

			ct = ct.next
			continue

		case typeDive:

			ct = ct.next

			// traverse slice or map here
			// or panic ;)
			switch kind ***REMOVED***
			case reflect.Slice, reflect.Array:

				for i := 0; i < current.Len(); i++ ***REMOVED***
					v.traverseField(topStruct, currentStruct, current.Index(i), errPrefix, nsPrefix, errs, partial, exclude, includeExclude, cs, &cField***REMOVED***Name: fmt.Sprintf(arrayIndexFieldName, cf.Name, i), AltName: fmt.Sprintf(arrayIndexFieldName, cf.AltName, i)***REMOVED***, ct)
				***REMOVED***

			case reflect.Map:
				for _, key := range current.MapKeys() ***REMOVED***
					v.traverseField(topStruct, currentStruct, current.MapIndex(key), errPrefix, nsPrefix, errs, partial, exclude, includeExclude, cs, &cField***REMOVED***Name: fmt.Sprintf(mapIndexFieldName, cf.Name, key.Interface()), AltName: fmt.Sprintf(mapIndexFieldName, cf.AltName, key.Interface())***REMOVED***, ct)
				***REMOVED***

			default:
				// throw error, if not a slice or map then should not have gotten here
				// bad dive tag
				panic("dive error! can't dive on a non slice or map")
			***REMOVED***

			return

		case typeOr:

			errTag := blank

			for ***REMOVED***

				if ct.fn(v, topStruct, currentStruct, current, typ, kind, ct.param) ***REMOVED***

					// drain rest of the 'or' values, then continue or leave
					for ***REMOVED***

						ct = ct.next

						if ct == nil ***REMOVED***
							return
						***REMOVED***

						if ct.typeof != typeOr ***REMOVED***
							continue OUTER
						***REMOVED***
					***REMOVED***
				***REMOVED***

				errTag += orSeparator + ct.tag

				if ct.next == nil ***REMOVED***
					// if we get here, no valid 'or' value and no more tags

					ns := errPrefix + cf.Name

					if ct.hasAlias ***REMOVED***
						errs[ns] = &FieldError***REMOVED***
							FieldNamespace: ns,
							NameNamespace:  nsPrefix + cf.AltName,
							Name:           cf.AltName,
							Field:          cf.Name,
							Tag:            ct.aliasTag,
							ActualTag:      ct.actualAliasTag,
							Value:          current.Interface(),
							Type:           typ,
							Kind:           kind,
						***REMOVED***
					***REMOVED*** else ***REMOVED***
						errs[errPrefix+cf.Name] = &FieldError***REMOVED***
							FieldNamespace: ns,
							NameNamespace:  nsPrefix + cf.AltName,
							Name:           cf.AltName,
							Field:          cf.Name,
							Tag:            errTag[1:],
							ActualTag:      errTag[1:],
							Value:          current.Interface(),
							Type:           typ,
							Kind:           kind,
						***REMOVED***
					***REMOVED***

					return
				***REMOVED***

				ct = ct.next
			***REMOVED***

		default:
			if !ct.fn(v, topStruct, currentStruct, current, typ, kind, ct.param) ***REMOVED***

				ns := errPrefix + cf.Name

				errs[ns] = &FieldError***REMOVED***
					FieldNamespace: ns,
					NameNamespace:  nsPrefix + cf.AltName,
					Name:           cf.AltName,
					Field:          cf.Name,
					Tag:            ct.aliasTag,
					ActualTag:      ct.tag,
					Value:          current.Interface(),
					Param:          ct.param,
					Type:           typ,
					Kind:           kind,
				***REMOVED***

				return

			***REMOVED***

			ct = ct.next
		***REMOVED***
	***REMOVED***
***REMOVED***
