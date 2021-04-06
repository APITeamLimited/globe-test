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
	"errors"
	"strconv"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/js/internal/modules"
	"github.com/loadimpact/k6/lib"
)

type data struct***REMOVED******REMOVED***

func init() ***REMOVED***
	modules.Register("k6/data", new(data))
***REMOVED***

const sharedArrayNamePrefix = "k6/data/SharedArray."

// XSharedArray is a constructor returning a shareable read-only array
// indentified by the name and having their contents be whatever the call returns
func (d *data) XSharedArray(ctx context.Context, name string, call goja.Callable) (goja.Value, error) ***REMOVED***
	if lib.GetState(ctx) != nil ***REMOVED***
		return nil, errors.New("new SharedArray must be called in the init context")
	***REMOVED***

	initEnv := common.GetInitEnv(ctx)
	if initEnv == nil ***REMOVED***
		return nil, errors.New("missing init environment")
	***REMOVED***
	if len(name) == 0 ***REMOVED***
		return nil, errors.New("empty name provided to SharedArray's constructor")
	***REMOVED***

	name = sharedArrayNamePrefix + name
	value := initEnv.SharedObjects.GetOrCreateShare(name, func() interface***REMOVED******REMOVED*** ***REMOVED***
		return getShareArrayFromCall(common.GetRuntime(ctx), call)
	***REMOVED***)
	array, ok := value.(sharedArray)
	if !ok ***REMOVED*** // TODO more info in the error?
		return nil, errors.New("wrong type of shared object")
	***REMOVED***

	return array.wrap(common.GetRuntime(ctx)), nil
***REMOVED***

func getShareArrayFromCall(rt *goja.Runtime, call goja.Callable) sharedArray ***REMOVED***
	gojaValue, err := call(goja.Undefined())
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***
	obj := gojaValue.ToObject(rt)
	if obj.ClassName() != "Array" ***REMOVED***
		common.Throw(rt, errors.New("only arrays can be made into SharedArray")) // TODO better error
	***REMOVED***
	arr := make([]string, obj.Get("length").ToInteger())

	stringify, _ := goja.AssertFunction(rt.GlobalObject().Get("JSON").ToObject(rt).Get("stringify"))
	var val goja.Value
	for i := range arr ***REMOVED***
		val, err = stringify(goja.Undefined(), obj.Get(strconv.Itoa(i)))
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		arr[i] = val.String()
	***REMOVED***

	return sharedArray***REMOVED***arr: arr***REMOVED***
***REMOVED***
