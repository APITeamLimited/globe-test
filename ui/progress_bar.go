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

package ui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	faint = color.New(color.Faint)
)

type ProgressBar struct ***REMOVED***
	Width    int
	Progress float64
***REMOVED***

func (b ProgressBar) String() string ***REMOVED***
	space := b.Width - 2
	filled := int(float64(space) * b.Progress)

	filling := ""
	caret := ""
	if filled > 0 ***REMOVED***
		if filled < space ***REMOVED***
			filling = strings.Repeat("=", filled-1)
			caret = ">"
		***REMOVED*** else ***REMOVED***
			filling = strings.Repeat("=", filled)
		***REMOVED***
	***REMOVED***

	padding := ""
	filler := "="
	if color.NoColor ***REMOVED***
		filler = " "
	***REMOVED***
	if space > filled ***REMOVED***
		padding = faint.Sprint(strings.Repeat(filler, space-filled))
	***REMOVED***

	return fmt.Sprintf("[%s%s%s]", filling, caret, padding)
***REMOVED***
