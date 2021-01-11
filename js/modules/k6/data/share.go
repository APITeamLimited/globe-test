/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2020 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package data

import (
	"context"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
)

// TODO fix it not working really well with setupData or just make it more broken
// TODO fix it working with console.log
type sharedArray struct ***REMOVED***
	arr []string
***REMOVED***

func (s sharedArray) wrap(ctxPtr *context.Context, rt *goja.Runtime) goja.Value ***REMOVED***
	cal, err := rt.RunString(arrayWrapperCode)
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***
	call, _ := goja.AssertFunction(cal)
	wrapped, err := call(goja.Undefined(), rt.ToValue(common.Bind(rt, s, ctxPtr)))
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***

	return wrapped
***REMOVED***

func (s sharedArray) Get(index int) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if index < 0 || index >= len(s.arr) ***REMOVED***
		return goja.Undefined(), nil
	***REMOVED***

	// we specifically use JSON.parse to get the json to an object inside as otherwise we won't be
	// able to freeze it as goja doesn't let us unless it is a pure goja object and this is the
	// easiest way to get one.
	return s.arr[index], nil
***REMOVED***

func (s sharedArray) Length() int ***REMOVED***
	return len(s.arr)
***REMOVED***

/* This implementation is commented as with it - it is harder to deepFreeze it with this implementation.
type sharedArrayIterator struct ***REMOVED***
	a     *sharedArray
	index int
***REMOVED***

func (sai *sharedArrayIterator) Next() (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if sai.index == len(sai.a.arr)-1 ***REMOVED***
		return map[string]bool***REMOVED***"done": true***REMOVED***, nil
	***REMOVED***
	sai.index++
	var tmp interface***REMOVED******REMOVED***
	if err := json.Unmarshal(sai.a.arr[sai.index], &tmp); err != nil ***REMOVED***
		return goja.Undefined(), err
	***REMOVED***
	return map[string]interface***REMOVED******REMOVED******REMOVED***"value": tmp***REMOVED***, nil
***REMOVED***

func (s sharedArray) Iterator() *sharedArrayIterator ***REMOVED***
	return &sharedArrayIterator***REMOVED***a: &s, index: -1***REMOVED***
***REMOVED***
*/

const arrayWrapperCode = `(function(val) ***REMOVED***
	function deepFreeze(o) ***REMOVED***
		Object.freeze(o);
		if (o === undefined) ***REMOVED***
			return o;
		***REMOVED***

		Object.getOwnPropertyNames(o).forEach(function (prop) ***REMOVED***
			if (o[prop] !== null
				&& (typeof o[prop] === "object" || typeof o[prop] === "function")
				&& !Object.isFrozen(o[prop])) ***REMOVED***
				deepFreeze(o[prop]);
			***REMOVED***
		***REMOVED***);

		return o;
	***REMOVED***;

	var arrayHandler = ***REMOVED***
		get: function(target, property, receiver) ***REMOVED***
			switch (property)***REMOVED***
			case "length":
				return target.length();
			case Symbol.iterator:
				return function()***REMOVED***
					var index = 0;
					return ***REMOVED***
						"next": function() ***REMOVED***
							if (index >= target.length()) ***REMOVED***
								return ***REMOVED***done: true***REMOVED***
							***REMOVED***
							var result = ***REMOVED***value: deepFreeze(JSON.parse(target.get(index)))***REMOVED***;
							index++;
							return result;
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
			var i = parseInt(property);
			if (isNaN(i)) ***REMOVED***
				return undefined;
			***REMOVED***

			return deepFreeze(JSON.parse(target.get(i)));
		***REMOVED***
	***REMOVED***;
	return new Proxy(val, arrayHandler);
***REMOVED***)`
