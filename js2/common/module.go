/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
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

package common

import (
	"context"
	"reflect"

	"github.com/dop251/goja"
	"github.com/pkg/errors"
)

type Module struct ***REMOVED***
	Context context.Context
	Impl    interface***REMOVED******REMOVED***
***REMOVED***

func (m *Module) Export(rt *goja.Runtime) goja.Value ***REMOVED***
	return m.Proxy(rt, m.Impl)
***REMOVED***

func (m *Module) Proxy(rt *goja.Runtime, v interface***REMOVED******REMOVED***) goja.Value ***REMOVED***
	ctxT := reflect.TypeOf((*context.Context)(nil)).Elem()
	errorT := reflect.TypeOf((*error)(nil)).Elem()

	exports := rt.NewObject()
	mapper := FieldNameMapper***REMOVED******REMOVED***

	val := reflect.ValueOf(v)
	typ := val.Type()
	for i := 0; i < typ.NumMethod(); i++ ***REMOVED***
		i := i
		methT := typ.Method(i)
		name := mapper.MethodName(typ, methT)
		if name == "" ***REMOVED***
			continue
		***REMOVED***
		meth := val.Method(i)

		in := make([]reflect.Type, methT.Type.NumIn())
		for i := 0; i < len(in); i++ ***REMOVED***
			in[i] = methT.Type.In(i)
		***REMOVED***
		out := make([]reflect.Type, methT.Type.NumOut())
		for i := 0; i < len(out); i++ ***REMOVED***
			out[i] = methT.Type.Out(i)
		***REMOVED***

		// Skip over the first input arg; it'll be the bound object.
		in = in[1:]

		// If the first argument is a context.Context, inject the given context.
		// The function will error if called outside of a valid context.
		if len(in) > 0 && in[0].Implements(ctxT) ***REMOVED***
			in = in[1:]
			meth = m.injectContext(in, out, methT, meth, rt)
		***REMOVED***

		// If the last return value is an error, turn it into a JS throw.
		if len(out) > 0 && out[len(out)-1] == errorT ***REMOVED***
			out = out[:len(out)-1]
			meth = m.injectErrorHandler(in, out, methT, meth, rt)
		***REMOVED***

		_ = exports.Set(name, meth.Interface())
	***REMOVED***

	return rt.ToValue(exports)
***REMOVED***

func (m *Module) injectContext(in, out []reflect.Type, methT reflect.Method, meth reflect.Value, rt *goja.Runtime) reflect.Value ***REMOVED***
	return reflect.MakeFunc(
		reflect.FuncOf(in, out, methT.Type.IsVariadic()),
		func(args []reflect.Value) []reflect.Value ***REMOVED***
			if m.Context == nil ***REMOVED***
				panic(rt.NewGoError(errors.Errorf("%s needs a valid VU context", methT.Name)))
			***REMOVED***

			select ***REMOVED***
			case <-m.Context.Done():
				panic(rt.NewGoError(errors.Errorf("test has ended")))
			default:
			***REMOVED***

			ctx := reflect.ValueOf(m.Context)
			return meth.Call(append([]reflect.Value***REMOVED***ctx***REMOVED***, args...))
		***REMOVED***,
	)
***REMOVED***

func (m *Module) injectErrorHandler(in, out []reflect.Type, methT reflect.Method, meth reflect.Value, rt *goja.Runtime) reflect.Value ***REMOVED***
	return reflect.MakeFunc(
		reflect.FuncOf(in, out, methT.Type.IsVariadic()),
		func(args []reflect.Value) []reflect.Value ***REMOVED***
			ret := meth.Call(args)
			err := ret[len(ret)-1]
			if !err.IsNil() ***REMOVED***
				panic(rt.NewGoError(err.Interface().(error)))
			***REMOVED***
			return ret[:len(ret)-1]
		***REMOVED***,
	)
***REMOVED***
