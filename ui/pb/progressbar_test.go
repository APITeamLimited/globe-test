/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
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

package pb

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProgressBarRender(t *testing.T) ***REMOVED***
	t.Parallel()

	testCases := []struct ***REMOVED***
		options  []ProgressBarOption
		expected string
	***REMOVED******REMOVED***
		***REMOVED***[]ProgressBarOption***REMOVED***WithLeft(func() string ***REMOVED*** return "left" ***REMOVED***)***REMOVED***,
			"left [--------------------------------------]"***REMOVED***,
		***REMOVED***[]ProgressBarOption***REMOVED***WithConstLeft("constLeft")***REMOVED***,
			"constLeft [--------------------------------------]"***REMOVED***,
		***REMOVED***[]ProgressBarOption***REMOVED***
			WithLeft(func() string ***REMOVED*** return "left" ***REMOVED***),
			WithProgress(func() (float64, string) ***REMOVED*** return 0, "right" ***REMOVED***),
		***REMOVED***,
			"left [--------------------------------------] right"***REMOVED***,
		***REMOVED***[]ProgressBarOption***REMOVED***
			WithLeft(func() string ***REMOVED*** return "left" ***REMOVED***),
			WithProgress(func() (float64, string) ***REMOVED*** return 0.5, "right" ***REMOVED***),
		***REMOVED***,
			"left [==================>-------------------] right"***REMOVED***,
		***REMOVED***[]ProgressBarOption***REMOVED***
			WithLeft(func() string ***REMOVED*** return "left" ***REMOVED***),
			WithProgress(func() (float64, string) ***REMOVED*** return 1.0, "right" ***REMOVED***),
		***REMOVED***,
			"left [======================================] right"***REMOVED***,
		***REMOVED***[]ProgressBarOption***REMOVED***
			WithLeft(func() string ***REMOVED*** return "left" ***REMOVED***),
			WithProgress(func() (float64, string) ***REMOVED*** return -1, "right" ***REMOVED***),
		***REMOVED***,
			"left [" + strings.Repeat("-", 76) + "] right"***REMOVED***,
		***REMOVED***[]ProgressBarOption***REMOVED***
			WithLeft(func() string ***REMOVED*** return "left" ***REMOVED***),
			WithProgress(func() (float64, string) ***REMOVED*** return 2, "right" ***REMOVED***),
		***REMOVED***,
			"left [" + strings.Repeat("=", 76) + "] right"***REMOVED***,
		***REMOVED***[]ProgressBarOption***REMOVED***
			WithLeft(func() string ***REMOVED*** return "left" ***REMOVED***),
			WithConstProgress(0.2, "constProgress"),
		***REMOVED***,
			"left [======>-------------------------------] constProgress"***REMOVED***,
		***REMOVED***[]ProgressBarOption***REMOVED***
			WithHijack(func() string ***REMOVED*** return "progressbar hijack!" ***REMOVED***),
		***REMOVED***,
			"progressbar hijack!"***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(tc.expected, func(t *testing.T) ***REMOVED***
			pbar := New(tc.options...)
			assert.NotNil(t, pbar)
			assert.Equal(t, tc.expected, pbar.Render(0))
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestProgressBarRenderPaddingMaxLeft(t *testing.T) ***REMOVED***
	t.Parallel()
	testCases := []struct ***REMOVED***
		maxLen   int
		left     string
		expected string
	***REMOVED******REMOVED***
		***REMOVED***-1, "left", "left [--------------------------------------]"***REMOVED***,
		***REMOVED***0, "left", "left [--------------------------------------]"***REMOVED***,
		***REMOVED***10, "left", "left       [--------------------------------------]"***REMOVED***,
		***REMOVED***10, "left_truncated",
			"left_tr... [--------------------------------------]"***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(tc.left, func(t *testing.T) ***REMOVED***
			pbar := New(WithLeft(func() string ***REMOVED*** return tc.left ***REMOVED***))
			assert.NotNil(t, pbar)
			assert.Equal(t, tc.expected, pbar.Render(tc.maxLen))
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestProgressBarLeft(t *testing.T) ***REMOVED***
	t.Parallel()

	testCases := []struct ***REMOVED***
		left     func() string
		expected string
	***REMOVED******REMOVED***
		***REMOVED***nil, ""***REMOVED***,
		***REMOVED***func() string ***REMOVED*** return " left " ***REMOVED***, " left "***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(tc.expected, func(t *testing.T) ***REMOVED***
			pbar := New(WithLeft(tc.left))
			assert.NotNil(t, pbar)
			assert.Equal(t, tc.expected, pbar.Left())
		***REMOVED***)
	***REMOVED***
***REMOVED***
