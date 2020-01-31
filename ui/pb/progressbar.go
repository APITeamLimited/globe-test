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
	"github.com/sirupsen/logrus"
)

const defaultWidth = 40
const defaultBarColor = color.Faint

// Status of the progress bar
type Status rune

// Progress bar status symbols
const (
	Running     Status = ' '
	Waiting     Status = '•'
	Stopping    Status = '↓'
	Interrupted Status = '✗'
	Done        Status = '✓'
)

//nolint:gochecknoglobals
var statusColors = map[Status]*color.Color***REMOVED***
	Interrupted: color.New(color.FgRed),
	Done:        color.New(color.FgGreen),
	Waiting:     color.New(defaultBarColor),
***REMOVED***

// ProgressBar is a simple thread-safe progressbar implementation with
// callbacks.
type ProgressBar struct ***REMOVED***
	mutex  sync.RWMutex
	width  int
	color  *color.Color
	logger *logrus.Entry
	status Status

	left     func() string
	progress func() (progress float64, right []string)
	hijack   func() string
***REMOVED***

// ProgressBarOption is used for helper functions that modify the progressbar
// parameters, either in the constructor or via the Modify() method.
type ProgressBarOption func(*ProgressBar)

// WithLeft modifies the function that returns the left progressbar value.
func WithLeft(left func() string) ProgressBarOption ***REMOVED***
	return func(pb *ProgressBar) ***REMOVED*** pb.left = left ***REMOVED***
***REMOVED***

// WithConstLeft sets the left progressbar value to the supplied const.
func WithConstLeft(left string) ProgressBarOption ***REMOVED***
	return func(pb *ProgressBar) ***REMOVED***
		pb.left = func() string ***REMOVED*** return left ***REMOVED***
	***REMOVED***
***REMOVED***

// WithLogger modifies the logger instance
func WithLogger(logger *logrus.Entry) ProgressBarOption ***REMOVED***
	return func(pb *ProgressBar) ***REMOVED*** pb.logger = logger ***REMOVED***
***REMOVED***

// WithProgress modifies the progress calculation function.
func WithProgress(progress func() (float64, []string)) ProgressBarOption ***REMOVED***
	return func(pb *ProgressBar) ***REMOVED*** pb.progress = progress ***REMOVED***
***REMOVED***

// WithStatus modifies the progressbar status
func WithStatus(status Status) ProgressBarOption ***REMOVED***
	return func(pb *ProgressBar) ***REMOVED*** pb.status = status ***REMOVED***
***REMOVED***

// WithConstProgress sets the progress and right values to the supplied consts.
func WithConstProgress(progress float64, right ...string) ProgressBarOption ***REMOVED***
	return func(pb *ProgressBar) ***REMOVED***
		pb.progress = func() (float64, []string) ***REMOVED*** return progress, right ***REMOVED***
	***REMOVED***
***REMOVED***

// WithHijack replaces the progressbar Render function with the argument.
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

// renderLeft renders the left part of the progressbar, replacing text
// exceeding maxLen with an ellipsis.
func (pb *ProgressBar) renderLeft(maxLen int) string ***REMOVED***
	var left string
	if pb.left != nil ***REMOVED***
		l := pb.left()
		if maxLen > 0 && len(l) > maxLen ***REMOVED***
			l = l[:maxLen-3] + "..."
		***REMOVED***
		left = l
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

// ProgressBarRender stores the different rendered parts of the
// progress bar UI.
type ProgressBarRender struct ***REMOVED***
	Left, Status, Progress, Hijack string
	Right                          []string
***REMOVED***

func (pbr ProgressBarRender) String() string ***REMOVED***
	if pbr.Hijack != "" ***REMOVED***
		return pbr.Hijack
	***REMOVED***
	var right string
	if len(pbr.Right) > 0 ***REMOVED***
		right = " " + strings.Join(pbr.Right, "  ")
	***REMOVED***
	return fmt.Sprintf("%s %-1s %s%s",
		pbr.Left, pbr.Status, pbr.Progress, right)
***REMOVED***

// Render locks the progressbar struct for reading and calls all of
// its methods to return the final output. A struct is returned over a
// plain string to allow dynamic padding and positioning of elements
// depending on other elements on the screen.
// - leftMax defines the maximum character length of the left-side
//   text. Characters exceeding this length will be replaced with a
//   single ellipsis. Passing <=0 disables this.
func (pb *ProgressBar) Render(leftMax int) ProgressBarRender ***REMOVED***
	pb.mutex.RLock()
	defer pb.mutex.RUnlock()

	var out ProgressBarRender
	if pb.hijack != nil ***REMOVED***
		out.Hijack = pb.hijack()
		return out
	***REMOVED***

	var progress float64
	if pb.progress != nil ***REMOVED***
		progress, out.Right = pb.progress()
		progressClamped := Clampf(progress, 0, 1)
		if progress != progressClamped ***REMOVED***
			progress = progressClamped
			if pb.logger != nil ***REMOVED***
				pb.logger.Warnf("progress value %.2f exceeds valid range, clamped between 0 and 1", progress)
			***REMOVED***
		***REMOVED***
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

	out.Left = pb.renderLeft(leftMax)

	switch c, ok := statusColors[pb.status]; ***REMOVED***
	case ok:
		out.Status = c.Sprint(string(pb.status))
	case pb.status > 0:
		out.Status = string(pb.status)
	default:
		out.Status = " "
	***REMOVED***

	out.Progress = fmt.Sprintf("[%s%s%s]", filling, caret, padding)

	return out
***REMOVED***
