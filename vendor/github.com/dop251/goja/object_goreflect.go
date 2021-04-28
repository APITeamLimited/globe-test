package goja

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"

	"github.com/dop251/goja/parser"
	"github.com/dop251/goja/unistring"
)

// JsonEncodable allows custom JSON encoding by JSON.stringify()
// Note that if the returned value itself also implements JsonEncodable, it won't have any effect.
type JsonEncodable interface ***REMOVED***
	JsonEncodable() interface***REMOVED******REMOVED***
***REMOVED***

// FieldNameMapper provides custom mapping between Go and JavaScript property names.
type FieldNameMapper interface ***REMOVED***
	// FieldName returns a JavaScript name for the given struct field in the given type.
	// If this method returns "" the field becomes hidden.
	FieldName(t reflect.Type, f reflect.StructField) string

	// MethodName returns a JavaScript name for the given method in the given type.
	// If this method returns "" the method becomes hidden.
	MethodName(t reflect.Type, m reflect.Method) string
***REMOVED***

type tagFieldNameMapper struct ***REMOVED***
	tagName      string
	uncapMethods bool
***REMOVED***

func (tfm tagFieldNameMapper) FieldName(_ reflect.Type, f reflect.StructField) string ***REMOVED***
	tag := f.Tag.Get(tfm.tagName)
	if idx := strings.IndexByte(tag, ','); idx != -1 ***REMOVED***
		tag = tag[:idx]
	***REMOVED***
	if parser.IsIdentifier(tag) ***REMOVED***
		return tag
	***REMOVED***
	return ""
***REMOVED***

func uncapitalize(s string) string ***REMOVED***
	return strings.ToLower(s[0:1]) + s[1:]
***REMOVED***

func (tfm tagFieldNameMapper) MethodName(_ reflect.Type, m reflect.Method) string ***REMOVED***
	if tfm.uncapMethods ***REMOVED***
		return uncapitalize(m.Name)
	***REMOVED***
	return m.Name
***REMOVED***

type uncapFieldNameMapper struct ***REMOVED***
***REMOVED***

func (u uncapFieldNameMapper) FieldName(_ reflect.Type, f reflect.StructField) string ***REMOVED***
	return uncapitalize(f.Name)
***REMOVED***

func (u uncapFieldNameMapper) MethodName(_ reflect.Type, m reflect.Method) string ***REMOVED***
	return uncapitalize(m.Name)
***REMOVED***

type reflectFieldInfo struct ***REMOVED***
	Index     []int
	Anonymous bool
***REMOVED***

type reflectTypeInfo struct ***REMOVED***
	Fields                  map[string]reflectFieldInfo
	Methods                 map[string]int
	FieldNames, MethodNames []string
***REMOVED***

type objectGoReflect struct ***REMOVED***
	baseObject
	origValue, value reflect.Value

	valueTypeInfo, origValueTypeInfo *reflectTypeInfo

	toJson func() interface***REMOVED******REMOVED***
***REMOVED***

func (o *objectGoReflect) init() ***REMOVED***
	o.baseObject.init()
	switch o.value.Kind() ***REMOVED***
	case reflect.Bool:
		o.class = classBoolean
		o.prototype = o.val.runtime.global.BooleanPrototype
	case reflect.String:
		o.class = classString
		o.prototype = o.val.runtime.global.StringPrototype
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:

		o.class = classNumber
		o.prototype = o.val.runtime.global.NumberPrototype
	default:
		o.class = classObject
		o.prototype = o.val.runtime.global.ObjectPrototype
	***REMOVED***
	o.extensible = true

	o.baseObject._putProp("toString", o.val.runtime.newNativeFunc(o.toStringFunc, nil, "toString", nil, 0), true, false, true)
	o.baseObject._putProp("valueOf", o.val.runtime.newNativeFunc(o.valueOfFunc, nil, "valueOf", nil, 0), true, false, true)

	o.valueTypeInfo = o.val.runtime.typeInfo(o.value.Type())
	o.origValueTypeInfo = o.val.runtime.typeInfo(o.origValue.Type())

	if j, ok := o.origValue.Interface().(JsonEncodable); ok ***REMOVED***
		o.toJson = j.JsonEncodable
	***REMOVED***
***REMOVED***

func (o *objectGoReflect) toStringFunc(FunctionCall) Value ***REMOVED***
	return o.toPrimitiveString()
***REMOVED***

func (o *objectGoReflect) valueOfFunc(FunctionCall) Value ***REMOVED***
	return o.toPrimitive()
***REMOVED***

func (o *objectGoReflect) getStr(name unistring.String, receiver Value) Value ***REMOVED***
	if v := o._get(name.String()); v != nil ***REMOVED***
		return v
	***REMOVED***
	return o.baseObject.getStr(name, receiver)
***REMOVED***

func (o *objectGoReflect) _getField(jsName string) reflect.Value ***REMOVED***
	if info, exists := o.valueTypeInfo.Fields[jsName]; exists ***REMOVED***
		v := o.value.FieldByIndex(info.Index)
		return v
	***REMOVED***

	return reflect.Value***REMOVED******REMOVED***
***REMOVED***

func (o *objectGoReflect) _getMethod(jsName string) reflect.Value ***REMOVED***
	if idx, exists := o.origValueTypeInfo.Methods[jsName]; exists ***REMOVED***
		return o.origValue.Method(idx)
	***REMOVED***

	return reflect.Value***REMOVED******REMOVED***
***REMOVED***

func (o *objectGoReflect) getAddr(v reflect.Value) reflect.Value ***REMOVED***
	if (v.Kind() == reflect.Struct || v.Kind() == reflect.Slice) && v.CanAddr() ***REMOVED***
		return v.Addr()
	***REMOVED***
	return v
***REMOVED***

func (o *objectGoReflect) _get(name string) Value ***REMOVED***
	if o.value.Kind() == reflect.Struct ***REMOVED***
		if v := o._getField(name); v.IsValid() ***REMOVED***
			return o.val.runtime.ToValue(o.getAddr(v).Interface())
		***REMOVED***
	***REMOVED***

	if v := o._getMethod(name); v.IsValid() ***REMOVED***
		return o.val.runtime.ToValue(v.Interface())
	***REMOVED***

	return nil
***REMOVED***

func (o *objectGoReflect) getOwnPropStr(name unistring.String) Value ***REMOVED***
	n := name.String()
	if o.value.Kind() == reflect.Struct ***REMOVED***
		if v := o._getField(n); v.IsValid() ***REMOVED***
			return &valueProperty***REMOVED***
				value:      o.val.runtime.ToValue(o.getAddr(v).Interface()),
				writable:   v.CanSet(),
				enumerable: true,
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if v := o._getMethod(n); v.IsValid() ***REMOVED***
		return &valueProperty***REMOVED***
			value:      o.val.runtime.ToValue(v.Interface()),
			enumerable: true,
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (o *objectGoReflect) setOwnStr(name unistring.String, val Value, throw bool) bool ***REMOVED***
	has, ok := o._put(name.String(), val, throw)
	if !has ***REMOVED***
		if res, ok := o._setForeignStr(name, nil, val, o.val, throw); !ok ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Cannot assign to property %s of a host object", name)
			return false
		***REMOVED*** else ***REMOVED***
			return res
		***REMOVED***
	***REMOVED***
	return ok
***REMOVED***

func (o *objectGoReflect) setForeignStr(name unistring.String, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return o._setForeignStr(name, trueValIfPresent(o._has(name.String())), val, receiver, throw)
***REMOVED***

func (o *objectGoReflect) setForeignIdx(idx valueInt, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return o._setForeignIdx(idx, nil, val, receiver, throw)
***REMOVED***

func (o *objectGoReflect) _put(name string, val Value, throw bool) (has, ok bool) ***REMOVED***
	if o.value.Kind() == reflect.Struct ***REMOVED***
		if v := o._getField(name); v.IsValid() ***REMOVED***
			if !v.CanSet() ***REMOVED***
				o.val.runtime.typeErrorResult(throw, "Cannot assign to a non-addressable or read-only property %s of a host object", name)
				return true, false
			***REMOVED***
			err := o.val.runtime.toReflectValue(val, v, &objectExportCtx***REMOVED******REMOVED***)
			if err != nil ***REMOVED***
				o.val.runtime.typeErrorResult(throw, "Go struct conversion error: %v", err)
				return true, false
			***REMOVED***
			return true, true
		***REMOVED***
	***REMOVED***
	return false, false
***REMOVED***

func (o *objectGoReflect) _putProp(name unistring.String, value Value, writable, enumerable, configurable bool) Value ***REMOVED***
	if _, ok := o._put(name.String(), value, false); ok ***REMOVED***
		return value
	***REMOVED***
	return o.baseObject._putProp(name, value, writable, enumerable, configurable)
***REMOVED***

func (r *Runtime) checkHostObjectPropertyDescr(name unistring.String, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	if descr.Getter != nil || descr.Setter != nil ***REMOVED***
		r.typeErrorResult(throw, "Host objects do not support accessor properties")
		return false
	***REMOVED***
	if descr.Writable == FLAG_FALSE ***REMOVED***
		r.typeErrorResult(throw, "Host object field %s cannot be made read-only", name)
		return false
	***REMOVED***
	if descr.Configurable == FLAG_TRUE ***REMOVED***
		r.typeErrorResult(throw, "Host object field %s cannot be made configurable", name)
		return false
	***REMOVED***
	return true
***REMOVED***

func (o *objectGoReflect) defineOwnPropertyStr(name unistring.String, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	if o.val.runtime.checkHostObjectPropertyDescr(name, descr, throw) ***REMOVED***
		n := name.String()
		if has, ok := o._put(n, descr.Value, throw); !has ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Cannot define property '%s' on a host object", n)
			return false
		***REMOVED*** else ***REMOVED***
			return ok
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (o *objectGoReflect) _has(name string) bool ***REMOVED***
	if o.value.Kind() == reflect.Struct ***REMOVED***
		if v := o._getField(name); v.IsValid() ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	if v := o._getMethod(name); v.IsValid() ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (o *objectGoReflect) hasOwnPropertyStr(name unistring.String) bool ***REMOVED***
	return o._has(name.String())
***REMOVED***

func (o *objectGoReflect) _toNumber() Value ***REMOVED***
	switch o.value.Kind() ***REMOVED***
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intToValue(o.value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return intToValue(int64(o.value.Uint()))
	case reflect.Bool:
		if o.value.Bool() ***REMOVED***
			return intToValue(1)
		***REMOVED*** else ***REMOVED***
			return intToValue(0)
		***REMOVED***
	case reflect.Float32, reflect.Float64:
		return floatToValue(o.value.Float())
	***REMOVED***
	return nil
***REMOVED***

func (o *objectGoReflect) _toString() Value ***REMOVED***
	switch o.value.Kind() ***REMOVED***
	case reflect.String:
		return newStringValue(o.value.String())
	case reflect.Bool:
		if o.value.Interface().(bool) ***REMOVED***
			return stringTrue
		***REMOVED*** else ***REMOVED***
			return stringFalse
		***REMOVED***
	***REMOVED***
	switch v := o.origValue.Interface().(type) ***REMOVED***
	case fmt.Stringer:
		return newStringValue(v.String())
	case error:
		return newStringValue(v.Error())
	***REMOVED***

	return stringObjectObject
***REMOVED***

func (o *objectGoReflect) toPrimitiveNumber() Value ***REMOVED***
	if v := o._toNumber(); v != nil ***REMOVED***
		return v
	***REMOVED***
	return o._toString()
***REMOVED***

func (o *objectGoReflect) toPrimitiveString() Value ***REMOVED***
	if v := o._toNumber(); v != nil ***REMOVED***
		return v.toString()
	***REMOVED***
	return o._toString()
***REMOVED***

func (o *objectGoReflect) toPrimitive() Value ***REMOVED***
	if o.prototype == o.val.runtime.global.NumberPrototype ***REMOVED***
		return o.toPrimitiveNumber()
	***REMOVED***
	return o.toPrimitiveString()
***REMOVED***

func (o *objectGoReflect) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	n := name.String()
	if o._has(n) ***REMOVED***
		o.val.runtime.typeErrorResult(throw, "Cannot delete property %s from a Go type", n)
		return false
	***REMOVED***
	return o.baseObject.deleteStr(name, throw)
***REMOVED***

type goreflectPropIter struct ***REMOVED***
	o   *objectGoReflect
	idx int
***REMOVED***

func (i *goreflectPropIter) nextField() (propIterItem, iterNextFunc) ***REMOVED***
	names := i.o.valueTypeInfo.FieldNames
	if i.idx < len(names) ***REMOVED***
		name := names[i.idx]
		i.idx++
		return propIterItem***REMOVED***name: unistring.NewFromString(name), enumerable: _ENUM_TRUE***REMOVED***, i.nextField
	***REMOVED***

	i.idx = 0
	return i.nextMethod()
***REMOVED***

func (i *goreflectPropIter) nextMethod() (propIterItem, iterNextFunc) ***REMOVED***
	names := i.o.origValueTypeInfo.MethodNames
	if i.idx < len(names) ***REMOVED***
		name := names[i.idx]
		i.idx++
		return propIterItem***REMOVED***name: unistring.NewFromString(name), enumerable: _ENUM_TRUE***REMOVED***, i.nextMethod
	***REMOVED***

	return propIterItem***REMOVED******REMOVED***, nil
***REMOVED***

func (o *objectGoReflect) enumerateOwnKeys() iterNextFunc ***REMOVED***
	r := &goreflectPropIter***REMOVED***
		o: o,
	***REMOVED***
	if o.value.Kind() == reflect.Struct ***REMOVED***
		return r.nextField
	***REMOVED***

	return r.nextMethod
***REMOVED***

func (o *objectGoReflect) ownKeys(_ bool, accum []Value) []Value ***REMOVED***
	// all own keys are enumerable
	for _, name := range o.valueTypeInfo.FieldNames ***REMOVED***
		accum = append(accum, newStringValue(name))
	***REMOVED***

	for _, name := range o.valueTypeInfo.MethodNames ***REMOVED***
		accum = append(accum, newStringValue(name))
	***REMOVED***

	return accum
***REMOVED***

func (o *objectGoReflect) export(*objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
	return o.origValue.Interface()
***REMOVED***

func (o *objectGoReflect) exportType() reflect.Type ***REMOVED***
	return o.origValue.Type()
***REMOVED***

func (o *objectGoReflect) equal(other objectImpl) bool ***REMOVED***
	if other, ok := other.(*objectGoReflect); ok ***REMOVED***
		return o.value.Interface() == other.value.Interface()
	***REMOVED***
	return false
***REMOVED***

func (r *Runtime) buildFieldInfo(t reflect.Type, index []int, info *reflectTypeInfo) ***REMOVED***
	n := t.NumField()
	for i := 0; i < n; i++ ***REMOVED***
		field := t.Field(i)
		name := field.Name
		if !ast.IsExported(name) ***REMOVED***
			continue
		***REMOVED***
		if r.fieldNameMapper != nil ***REMOVED***
			name = r.fieldNameMapper.FieldName(t, field)
		***REMOVED***

		if name != "" ***REMOVED***
			if inf, exists := info.Fields[name]; !exists ***REMOVED***
				info.FieldNames = append(info.FieldNames, name)
			***REMOVED*** else ***REMOVED***
				if len(inf.Index) <= len(index) ***REMOVED***
					continue
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if name != "" || field.Anonymous ***REMOVED***
			idx := make([]int, len(index)+1)
			copy(idx, index)
			idx[len(idx)-1] = i

			if name != "" ***REMOVED***
				info.Fields[name] = reflectFieldInfo***REMOVED***
					Index:     idx,
					Anonymous: field.Anonymous,
				***REMOVED***
			***REMOVED***
			if field.Anonymous ***REMOVED***
				typ := field.Type
				for typ.Kind() == reflect.Ptr ***REMOVED***
					typ = typ.Elem()
				***REMOVED***
				if typ.Kind() == reflect.Struct ***REMOVED***
					r.buildFieldInfo(typ, idx, info)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *Runtime) buildTypeInfo(t reflect.Type) (info *reflectTypeInfo) ***REMOVED***
	info = new(reflectTypeInfo)
	if t.Kind() == reflect.Struct ***REMOVED***
		info.Fields = make(map[string]reflectFieldInfo)
		n := t.NumField()
		info.FieldNames = make([]string, 0, n)
		r.buildFieldInfo(t, nil, info)
	***REMOVED***

	info.Methods = make(map[string]int)
	n := t.NumMethod()
	info.MethodNames = make([]string, 0, n)
	for i := 0; i < n; i++ ***REMOVED***
		method := t.Method(i)
		name := method.Name
		if !ast.IsExported(name) ***REMOVED***
			continue
		***REMOVED***
		if r.fieldNameMapper != nil ***REMOVED***
			name = r.fieldNameMapper.MethodName(t, method)
			if name == "" ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***

		if _, exists := info.Methods[name]; !exists ***REMOVED***
			info.MethodNames = append(info.MethodNames, name)
		***REMOVED***

		info.Methods[name] = i
	***REMOVED***
	return
***REMOVED***

func (r *Runtime) typeInfo(t reflect.Type) (info *reflectTypeInfo) ***REMOVED***
	var exists bool
	if info, exists = r.typeInfoCache[t]; !exists ***REMOVED***
		info = r.buildTypeInfo(t)
		if r.typeInfoCache == nil ***REMOVED***
			r.typeInfoCache = make(map[reflect.Type]*reflectTypeInfo)
		***REMOVED***
		r.typeInfoCache[t] = info
	***REMOVED***

	return
***REMOVED***

// SetFieldNameMapper sets a custom field name mapper for Go types. It can be called at any time, however
// the mapping for any given value is fixed at the point of creation.
// Setting this to nil restores the default behaviour which is all exported fields and methods are mapped to their
// original unchanged names.
func (r *Runtime) SetFieldNameMapper(mapper FieldNameMapper) ***REMOVED***
	r.fieldNameMapper = mapper
	r.typeInfoCache = nil
***REMOVED***

// TagFieldNameMapper returns a FieldNameMapper that uses the given tagName for struct fields and optionally
// uncapitalises (making the first letter lower case) method names.
// The common tag value syntax is supported (name[,options]), however options are ignored.
// Setting name to anything other than a valid ECMAScript identifier makes the field hidden.
func TagFieldNameMapper(tagName string, uncapMethods bool) FieldNameMapper ***REMOVED***
	return tagFieldNameMapper***REMOVED***
		tagName:      tagName,
		uncapMethods: uncapMethods,
	***REMOVED***
***REMOVED***

// UncapFieldNameMapper returns a FieldNameMapper that uncapitalises struct field and method names
// making the first letter lower case.
func UncapFieldNameMapper() FieldNameMapper ***REMOVED***
	return uncapFieldNameMapper***REMOVED******REMOVED***
***REMOVED***
