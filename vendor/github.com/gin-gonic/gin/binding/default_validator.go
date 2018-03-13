// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"reflect"
	"sync"

	"gopkg.in/go-playground/validator.v8"
)

type defaultValidator struct ***REMOVED***
	once     sync.Once
	validate *validator.Validate
***REMOVED***

var _ StructValidator = &defaultValidator***REMOVED******REMOVED***

func (v *defaultValidator) ValidateStruct(obj interface***REMOVED******REMOVED***) error ***REMOVED***
	if kindOfData(obj) == reflect.Struct ***REMOVED***
		v.lazyinit()
		if err := v.validate.Struct(obj); err != nil ***REMOVED***
			return error(err)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (v *defaultValidator) lazyinit() ***REMOVED***
	v.once.Do(func() ***REMOVED***
		config := &validator.Config***REMOVED***TagName: "binding"***REMOVED***
		v.validate = validator.New(config)
	***REMOVED***)
***REMOVED***

func kindOfData(data interface***REMOVED******REMOVED***) reflect.Kind ***REMOVED***
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr ***REMOVED***
		valueType = value.Elem().Kind()
	***REMOVED***
	return valueType
***REMOVED***
