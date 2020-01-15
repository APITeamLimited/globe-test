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

package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/loadimpact/k6/core/local"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/ui/pb"
)

// TODO: Make configurable
const maxLeftLength = 30

// A writer that syncs writes with a mutex and, if the output is a TTY, clears before newlines.
type consoleWriter struct ***REMOVED***
	Writer io.Writer
	IsTTY  bool
	Mutex  *sync.Mutex

	// Used for flicker-free persistent objects like the progressbars
	PersistentText func()
***REMOVED***

func (w *consoleWriter) Write(p []byte) (n int, err error) ***REMOVED***
	origLen := len(p)
	if w.IsTTY ***REMOVED***
		// Add a TTY code to erase till the end of line with each new line
		// TODO: check how cross-platform this is...
		p = bytes.Replace(p, []byte***REMOVED***'\n'***REMOVED***, []byte***REMOVED***'\x1b', '[', '0', 'K', '\n'***REMOVED***, -1)
	***REMOVED***

	w.Mutex.Lock()
	n, err = w.Writer.Write(p)
	if w.PersistentText != nil ***REMOVED***
		w.PersistentText()
	***REMOVED***
	w.Mutex.Unlock()

	if err != nil && n < origLen ***REMOVED***
		return n, err
	***REMOVED***
	return origLen, err
***REMOVED***

func printBar(bar *pb.ProgressBar, rightText string) ***REMOVED***
	end := "\n"
	if stdout.IsTTY ***REMOVED***
		// If we're in a TTY, instead of printing the bar and going to the next
		// line, erase everything till the end of the line and return to the
		// start, so that the next print will overwrite the same line.
		//
		// TODO: check for cross platform support
		end = "\x1b[0K\r"
	***REMOVED***
	rendered := bar.Render(0)
	// Only output the left and middle part of the progress bar
	fprintf(stdout, "%s %s %s%s", rendered.Left, rendered.Progress, rightText, end)
***REMOVED***

func renderMultipleBars(isTTY, goBack bool, leftMax int, pbs []*pb.ProgressBar) string ***REMOVED***
	lineEnd := "\n"
	if isTTY ***REMOVED***
		//TODO: check for cross platform support
		lineEnd = "\x1b[K\n" // erase till end of line
	***REMOVED***

	var (
		// Maximum length of each right side column except last,
		// used to calculate the padding between columns.
		maxRColumnLen = make([]int, 1)
		pbsCount      = len(pbs)
		rendered      = make([]pb.ProgressBarRender, pbsCount)
		result        = make([]string, pbsCount+2)
	)

	result[0] = lineEnd // start with an empty line

	// First pass to render all progressbars and get the maximum
	// lengths of right-side columns.
	for i, pb := range pbs ***REMOVED***
		rend := pb.Render(leftMax)
		for i := range rend.Right ***REMOVED***
			// Don't calculate for last column, since there's nothing to align
			// after it (yet?).
			if i == len(rend.Right)-1 ***REMOVED***
				break
			***REMOVED***
			if len(rend.Right[i]) > maxRColumnLen[i] ***REMOVED***
				maxRColumnLen[i] = len(rend.Right[i])
			***REMOVED***
		***REMOVED***
		rendered[i] = rend
	***REMOVED***

	// Second pass to render final output, applying padding where needed
	for i := range rendered ***REMOVED***
		rend := rendered[i]
		if rend.Hijack != "" ***REMOVED***
			result[i+1] = rend.Hijack + lineEnd
			continue
		***REMOVED***
		var leftText, rightText string
		leftPadFmt := fmt.Sprintf("%%-%ds %%s ", leftMax)
		leftText = fmt.Sprintf(leftPadFmt, rend.Left, rend.Status)
		for i := range rend.Right ***REMOVED***
			rpad := 0
			if len(maxRColumnLen) > i ***REMOVED***
				rpad = maxRColumnLen[i]
			***REMOVED***
			rightPadFmt := fmt.Sprintf(" %%-%ds", rpad+1)
			rightText += fmt.Sprintf(rightPadFmt, rend.Right[i])
		***REMOVED***
		result[i+1] = leftText + rend.Progress + rightText + lineEnd
	***REMOVED***

	if isTTY && goBack ***REMOVED***
		// Go back to the beginning
		//TODO: check for cross platform support
		result[pbsCount+1] = fmt.Sprintf("\r\x1b[%dA", pbsCount+1)
	***REMOVED*** else ***REMOVED***
		result[pbsCount+1] = "\n"
	***REMOVED***
	return strings.Join(result, "")
***REMOVED***

//TODO: show other information here?
//TODO: add a no-progress option that will disable these
//TODO: don't use global variables...
func showProgress(ctx context.Context, conf Config, execScheduler *local.ExecutionScheduler) ***REMOVED***
	if quiet || conf.HTTPDebug.Valid && conf.HTTPDebug.String != "" ***REMOVED***
		return
	***REMOVED***

	pbs := []*pb.ProgressBar***REMOVED***execScheduler.GetInitProgressBar()***REMOVED***
	for _, s := range execScheduler.GetExecutors() ***REMOVED***
		pbs = append(pbs, s.GetProgress())
	***REMOVED***

	// Get the longest left side string length, to align progress bars
	// horizontally and trim excess text.
	var leftLen int64
	for _, pb := range pbs ***REMOVED***
		l := pb.Left()
		leftLen = lib.Max(int64(len(l)), leftLen)
	***REMOVED***

	// Limit to maximum left text length
	leftMax := int(lib.Min(leftLen, maxLeftLength))

	// For flicker-free progressbars!
	progressBarsLastRender := []byte(renderMultipleBars(stdoutTTY, true, leftMax, pbs))
	progressBarsPrint := func() ***REMOVED***
		_, _ = stdout.Writer.Write(progressBarsLastRender)
	***REMOVED***

	//TODO: make configurable?
	updateFreq := 1 * time.Second
	//TODO: remove !noColor after we fix how we handle colors (see the related
	//description in the TODO message in cmd/root.go)
	if stdoutTTY && !noColor ***REMOVED***
		updateFreq = 100 * time.Millisecond
		outMutex.Lock()
		stdout.PersistentText = progressBarsPrint
		stderr.PersistentText = progressBarsPrint
		outMutex.Unlock()
		defer func() ***REMOVED***
			outMutex.Lock()
			stdout.PersistentText = nil
			stderr.PersistentText = nil
			if ctx.Err() != nil ***REMOVED***
				// Render a last plain-text progressbar in an error
				progressBarsLastRender = []byte(renderMultipleBars(stdoutTTY, false, leftMax, pbs))
				progressBarsPrint()
			***REMOVED***
			outMutex.Unlock()
		***REMOVED***()
	***REMOVED***

	ctxDone := ctx.Done()
	ticker := time.NewTicker(updateFreq)
	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			barText := renderMultipleBars(stdoutTTY, true, leftMax, pbs)
			outMutex.Lock()
			progressBarsLastRender = []byte(barText)
			progressBarsPrint()
			outMutex.Unlock()
		case <-ctxDone:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***
