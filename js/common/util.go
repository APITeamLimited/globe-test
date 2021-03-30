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
	"bytes"
	"fmt"
	"io"

	"github.com/dop251/goja"
)

// Throw a JS error; avoids re-wrapping GoErrors.
func Throw(rt *goja.Runtime, err error) ***REMOVED***
	if e, ok := err.(*goja.Exception); ok ***REMOVED***
		panic(e)
	***REMOVED***
	panic(rt.NewGoError(err))
***REMOVED***

// GetReader tries to return an io.Reader value from an exported goja value.
func GetReader(data interface***REMOVED******REMOVED***) (io.Reader, error) ***REMOVED***
	switch r := data.(type) ***REMOVED***
	case string:
		return bytes.NewBufferString(r), nil
	case []byte:
		return bytes.NewBuffer(r), nil
	case io.Reader:
		return r, nil
	case goja.ArrayBuffer:
		return bytes.NewBuffer(r.Bytes()), nil
	default:
		return nil, fmt.Errorf("invalid type %T, it needs to be a string, byte array or an ArrayBuffer", data)
	***REMOVED***
***REMOVED***

// ToBytes tries to return a byte slice from compatible types.
func ToBytes(data interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	switch dt := data.(type) ***REMOVED***
	case []byte:
		return dt, nil
	case string:
		return []byte(dt), nil
	case goja.ArrayBuffer:
		return dt.Bytes(), nil
	default:
		return nil, fmt.Errorf("invalid type %T, expected string, []byte or ArrayBuffer", data)
	***REMOVED***
***REMOVED***

// ToString tries to return a string from compatible types.
func ToString(data interface***REMOVED******REMOVED***) (string, error) ***REMOVED***
	switch dt := data.(type) ***REMOVED***
	case []byte:
		return string(dt), nil
	case string:
		return dt, nil
	case goja.ArrayBuffer:
		return string(dt.Bytes()), nil
	default:
		return "", fmt.Errorf("invalid type %T, expected string, []byte or ArrayBuffer", data)
	***REMOVED***
***REMOVED***
