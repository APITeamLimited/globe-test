/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2021 Load Impact
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

package http

import (
	"context"
	"fmt"
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.k6.io/k6/js/common"
)

func TestHTTPFile(t *testing.T) ***REMOVED***
	t.Parallel()
	rt := goja.New()
	input := []byte***REMOVED***104, 101, 108, 108, 111***REMOVED***

	testCases := []struct ***REMOVED***
		input    interface***REMOVED******REMOVED***
		args     []string
		expected FileData
		expErr   string
	***REMOVED******REMOVED***
		// We can't really test without specifying a filename argument,
		// as File() calls time.Now(), so we'd need some time freezing/mocking
		// or refactoring, or to exclude the field from the assertion.
		***REMOVED***
			input,
			[]string***REMOVED***"test.bin"***REMOVED***,
			FileData***REMOVED***Data: input, Filename: "test.bin", ContentType: "application/octet-stream"***REMOVED***,
			"",
		***REMOVED***,
		***REMOVED***
			string(input),
			[]string***REMOVED***"test.txt", "text/plain"***REMOVED***,
			FileData***REMOVED***Data: input, Filename: "test.txt", ContentType: "text/plain"***REMOVED***,
			"",
		***REMOVED***,
		***REMOVED***
			rt.NewArrayBuffer(input),
			[]string***REMOVED***"test-ab.bin"***REMOVED***,
			FileData***REMOVED***Data: input, Filename: "test-ab.bin", ContentType: "application/octet-stream"***REMOVED***,
			"",
		***REMOVED***,
		***REMOVED***struct***REMOVED******REMOVED******REMOVED******REMOVED***, []string***REMOVED******REMOVED***, FileData***REMOVED******REMOVED***, "invalid type struct ***REMOVED******REMOVED***, expected string, []byte or ArrayBuffer"***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(fmt.Sprintf("%T", tc.input), func(t *testing.T) ***REMOVED***
			if tc.expErr != "" ***REMOVED***
				defer func() ***REMOVED***
					err := recover()
					require.NotNil(t, err)
					require.IsType(t, &goja.Object***REMOVED******REMOVED***, err)
					val := err.(*goja.Object).Export()
					require.EqualError(t, val.(error), tc.expErr)
				***REMOVED***()
			***REMOVED***
			h := new(GlobalHTTP).NewModuleInstancePerVU().(*HTTP)
			ctx := common.WithRuntime(context.Background(), rt)
			out := h.File(ctx, tc.input, tc.args...)
			assert.Equal(t, tc.expected, out)
		***REMOVED***)
	***REMOVED***
***REMOVED***
