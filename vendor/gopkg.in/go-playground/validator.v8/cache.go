package validator

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
)

type tagType uint8

const (
	typeDefault tagType = iota
	typeOmitEmpty
	typeNoStructLevel
	typeStructOnly
	typeDive
	typeOr
	typeExists
)

type structCache struct ***REMOVED***
	lock sync.Mutex
	m    atomic.Value // map[reflect.Type]*cStruct
***REMOVED***

func (sc *structCache) Get(key reflect.Type) (c *cStruct, found bool) ***REMOVED***
	c, found = sc.m.Load().(map[reflect.Type]*cStruct)[key]
	return
***REMOVED***

func (sc *structCache) Set(key reflect.Type, value *cStruct) ***REMOVED***

	m := sc.m.Load().(map[reflect.Type]*cStruct)

	nm := make(map[reflect.Type]*cStruct, len(m)+1)
	for k, v := range m ***REMOVED***
		nm[k] = v
	***REMOVED***
	nm[key] = value
	sc.m.Store(nm)
***REMOVED***

type tagCache struct ***REMOVED***
	lock sync.Mutex
	m    atomic.Value // map[string]*cTag
***REMOVED***

func (tc *tagCache) Get(key string) (c *cTag, found bool) ***REMOVED***
	c, found = tc.m.Load().(map[string]*cTag)[key]
	return
***REMOVED***

func (tc *tagCache) Set(key string, value *cTag) ***REMOVED***

	m := tc.m.Load().(map[string]*cTag)

	nm := make(map[string]*cTag, len(m)+1)
	for k, v := range m ***REMOVED***
		nm[k] = v
	***REMOVED***
	nm[key] = value
	tc.m.Store(nm)
***REMOVED***

type cStruct struct ***REMOVED***
	Name   string
	fields map[int]*cField
	fn     StructLevelFunc
***REMOVED***

type cField struct ***REMOVED***
	Idx     int
	Name    string
	AltName string
	cTags   *cTag
***REMOVED***

type cTag struct ***REMOVED***
	tag            string
	aliasTag       string
	actualAliasTag string
	param          string
	hasAlias       bool
	typeof         tagType
	hasTag         bool
	fn             Func
	next           *cTag
***REMOVED***

func (v *Validate) extractStructCache(current reflect.Value, sName string) *cStruct ***REMOVED***

	v.structCache.lock.Lock()
	defer v.structCache.lock.Unlock() // leave as defer! because if inner panics, it will never get unlocked otherwise!

	typ := current.Type()

	// could have been multiple trying to access, but once first is done this ensures struct
	// isn't parsed again.
	cs, ok := v.structCache.Get(typ)
	if ok ***REMOVED***
		return cs
	***REMOVED***

	cs = &cStruct***REMOVED***Name: sName, fields: make(map[int]*cField), fn: v.structLevelFuncs[typ]***REMOVED***

	numFields := current.NumField()

	var ctag *cTag
	var fld reflect.StructField
	var tag string
	var customName string

	for i := 0; i < numFields; i++ ***REMOVED***

		fld = typ.Field(i)

		if !fld.Anonymous && fld.PkgPath != blank ***REMOVED***
			continue
		***REMOVED***

		tag = fld.Tag.Get(v.tagName)

		if tag == skipValidationTag ***REMOVED***
			continue
		***REMOVED***

		customName = fld.Name

		if v.fieldNameTag != blank ***REMOVED***

			name := strings.SplitN(fld.Tag.Get(v.fieldNameTag), ",", 2)[0]

			// dash check is for json "-" (aka skipValidationTag) means don't output in json
			if name != "" && name != skipValidationTag ***REMOVED***
				customName = name
			***REMOVED***
		***REMOVED***

		// NOTE: cannot use shared tag cache, because tags may be equal, but things like alias may be different
		// and so only struct level caching can be used instead of combined with Field tag caching

		if len(tag) > 0 ***REMOVED***
			ctag, _ = v.parseFieldTagsRecursive(tag, fld.Name, blank, false)
		***REMOVED*** else ***REMOVED***
			// even if field doesn't have validations need cTag for traversing to potential inner/nested
			// elements of the field.
			ctag = new(cTag)
		***REMOVED***

		cs.fields[i] = &cField***REMOVED***Idx: i, Name: fld.Name, AltName: customName, cTags: ctag***REMOVED***
	***REMOVED***

	v.structCache.Set(typ, cs)

	return cs
***REMOVED***

func (v *Validate) parseFieldTagsRecursive(tag string, fieldName string, alias string, hasAlias bool) (firstCtag *cTag, current *cTag) ***REMOVED***

	var t string
	var ok bool
	noAlias := len(alias) == 0
	tags := strings.Split(tag, tagSeparator)

	for i := 0; i < len(tags); i++ ***REMOVED***

		t = tags[i]

		if noAlias ***REMOVED***
			alias = t
		***REMOVED***

		if v.hasAliasValidators ***REMOVED***
			// check map for alias and process new tags, otherwise process as usual
			if tagsVal, found := v.aliasValidators[t]; found ***REMOVED***

				if i == 0 ***REMOVED***
					firstCtag, current = v.parseFieldTagsRecursive(tagsVal, fieldName, t, true)
				***REMOVED*** else ***REMOVED***
					next, curr := v.parseFieldTagsRecursive(tagsVal, fieldName, t, true)
					current.next, current = next, curr

				***REMOVED***

				continue
			***REMOVED***
		***REMOVED***

		if i == 0 ***REMOVED***
			current = &cTag***REMOVED***aliasTag: alias, hasAlias: hasAlias, hasTag: true***REMOVED***
			firstCtag = current
		***REMOVED*** else ***REMOVED***
			current.next = &cTag***REMOVED***aliasTag: alias, hasAlias: hasAlias, hasTag: true***REMOVED***
			current = current.next
		***REMOVED***

		switch t ***REMOVED***

		case diveTag:
			current.typeof = typeDive
			continue

		case omitempty:
			current.typeof = typeOmitEmpty
			continue

		case structOnlyTag:
			current.typeof = typeStructOnly
			continue

		case noStructLevelTag:
			current.typeof = typeNoStructLevel
			continue

		case existsTag:
			current.typeof = typeExists
			continue

		default:

			// if a pipe character is needed within the param you must use the utf8Pipe representation "0x7C"
			orVals := strings.Split(t, orSeparator)

			for j := 0; j < len(orVals); j++ ***REMOVED***

				vals := strings.SplitN(orVals[j], tagKeySeparator, 2)

				if noAlias ***REMOVED***
					alias = vals[0]
					current.aliasTag = alias
				***REMOVED*** else ***REMOVED***
					current.actualAliasTag = t
				***REMOVED***

				if j > 0 ***REMOVED***
					current.next = &cTag***REMOVED***aliasTag: alias, actualAliasTag: current.actualAliasTag, hasAlias: hasAlias, hasTag: true***REMOVED***
					current = current.next
				***REMOVED***

				current.tag = vals[0]
				if len(current.tag) == 0 ***REMOVED***
					panic(strings.TrimSpace(fmt.Sprintf(invalidValidation, fieldName)))
				***REMOVED***

				if current.fn, ok = v.validationFuncs[current.tag]; !ok ***REMOVED***
					panic(strings.TrimSpace(fmt.Sprintf(undefinedValidation, fieldName)))
				***REMOVED***

				if len(orVals) > 1 ***REMOVED***
					current.typeof = typeOr
				***REMOVED***

				if len(vals) > 1 ***REMOVED***
					current.param = strings.Replace(strings.Replace(vals[1], utf8HexComma, ",", -1), utf8Pipe, "|", -1)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return
***REMOVED***
