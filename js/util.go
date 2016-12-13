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

package js

import (
	"github.com/robertkrimen/otto"
)

func Check(val, arg0 otto.Value) (bool, error) ***REMOVED***
	switch ***REMOVED***
	case val.IsFunction():
		val, err := val.Call(otto.UndefinedValue(), arg0)
		if err != nil ***REMOVED***
			return false, err
		***REMOVED***
		return Check(val, arg0)
	case val.IsBoolean():
		b, err := val.ToBoolean()
		if err != nil ***REMOVED***
			return false, err
		***REMOVED***
		return b, nil
	case val.IsNumber():
		f, err := val.ToFloat()
		if err != nil ***REMOVED***
			return false, err
		***REMOVED***
		return f != 0, nil
	case val.IsString():
		s, err := val.ToString()
		if err != nil ***REMOVED***
			return false, err
		***REMOVED***
		return s != "", nil
	default:
		return false, nil
	***REMOVED***
***REMOVED***

func throw(vm *otto.Otto, v interface***REMOVED******REMOVED***) ***REMOVED***
	if err, ok := v.(error); ok ***REMOVED***
		panic(vm.MakeCustomError("Error", err.Error()))
	***REMOVED***
	panic(v)
***REMOVED***
