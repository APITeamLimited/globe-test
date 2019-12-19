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
	"fmt"
	"strings"
	"sync"

	"github.com/fatih/color"
)

const defaultWidth = 40
const defaultBarColor = color.Faint

// ProgressBar is just a simple thread-safe progressbar implementation with
// callbacks.
type ProgressBar struct ***REMOVED***
	mutex sync.RWMutex
	width int
	color *color.Color

	left     func() string
	progress func() (progress float64, right string)
	hijack   func() string
***REMOVED***

// ProgressBarOption is used for helper functions that modify the progressbar
// parameters, either in the constructor or via the Modify() method.
type ProgressBarOption func(*ProgressBar)

// WithLeft modifies the function that returns the left progressbar padding.
func WithLeft(left func() string) ProgressBarOption ***REMOVED***
	return func(pb *ProgressBar) ***REMOVED*** pb.left = left ***REMOVED***
***REMOVED***

// WithConstLeft sets the left progressbar padding to the supplied const.
func WithConstLeft(left string) ProgressBarOption ***REMOVED***
	return func(pb *ProgressBar) ***REMOVED***
		pb.left = func() string ***REMOVED*** return left ***REMOVED***
	***REMOVED***
***REMOVED***

// WithProgress modifies the progress calculation function.
func WithProgress(progress func() (float64, string)) ProgressBarOption ***REMOVED***
	return func(pb *ProgressBar) ***REMOVED*** pb.progress = progress ***REMOVED***
***REMOVED***

// WithConstProgress sets the progress and right padding to the supplied consts.
func WithConstProgress(progress float64, right string) ProgressBarOption ***REMOVED***
	return func(pb *ProgressBar) ***REMOVED***
		pb.progress = func() (float64, string) ***REMOVED*** return progress, right ***REMOVED***
	***REMOVED***
***REMOVED***

// WithHijack replaces the progressbar String function with the argument.
func WithHijack(hijack func() string) ProgressBarOption ***REMOVED***
	return func(pb *ProgressBar) ***REMOVED*** pb.hijack = hijack ***REMOVED***
***REMOVED***

// New creates and initializes a new ProgressBar struct, calling all of the
// supplied options
func New(options ...ProgressBarOption) *ProgressBar ***REMOVED***
	pb := &ProgressBar***REMOVED***
		mutex: sync.RWMutex***REMOVED******REMOVED***,
		width: defaultWidth,
		color: color.New(defaultBarColor),
	***REMOVED***
	pb.Modify(options...)
	return pb
***REMOVED***

// Left returns the left part of the progressbar in a thread-safe way.
func (pb *ProgressBar) Left() string ***REMOVED***
	pb.mutex.RLock()
	defer pb.mutex.RUnlock()

	return pb.renderLeft(0)
***REMOVED***

func (pb *ProgressBar) renderLeft(pad int) string ***REMOVED***
	var left string
	if pb.left != nil ***REMOVED***
		padFmt := fmt.Sprintf("%%-%ds", pad)
		left = fmt.Sprintf(padFmt, pb.left())
	***REMOVED***
	return left
***REMOVED***

// Modify changes the progressbar options in a thread-safe way.
func (pb *ProgressBar) Modify(options ...ProgressBarOption) ***REMOVED***
	pb.mutex.Lock()
	defer pb.mutex.Unlock()
	for _, option := range options ***REMOVED***
		option(pb)
	***REMOVED***
***REMOVED***

// String locks the progressbar struct for reading and calls all of its methods
// to assemble the progress bar and return it as a string.
//TODO: something prettier? paddings, right-alignment of the left column, line trimming by terminal size
func (pb *ProgressBar) String(leftPad int) string ***REMOVED***
	pb.mutex.RLock()
	defer pb.mutex.RUnlock()

	if pb.hijack != nil ***REMOVED***
		return pb.hijack()
	***REMOVED***

	var (
		progress float64
		right    string
	)
	if pb.progress != nil ***REMOVED***
		progress, right = pb.progress()
		right = " " + right
	***REMOVED***

	space := pb.width - 2
	filled := int(float64(space) * progress)

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
	if space > filled ***REMOVED***
		padding = pb.color.Sprint(strings.Repeat("-", space-filled))
	***REMOVED***

	return fmt.Sprintf("%s [%s%s%s]%s", pb.renderLeft(leftPad), filling, caret, padding, right)
***REMOVED***
