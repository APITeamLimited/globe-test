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
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/loadimpact/k6/core/local"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/ui/pb"
)

const (
	// Max length of left-side progress bar text before trimming is forced
	maxLeftLength = 30
	// Amount of padding in chars between rendered progress
	// bar text and right-side terminal window edge.
	termPadding = 1
)

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

func printBar(bar *pb.ProgressBar) ***REMOVED***
	end := "\n"
	if stdout.IsTTY ***REMOVED***
		// If we're in a TTY, instead of printing the bar and going to the next
		// line, erase everything till the end of the line and return to the
		// start, so that the next print will overwrite the same line.
		//
		// TODO: check for cross platform support
		end = "\x1b[0K\r"
	***REMOVED***
	rendered := bar.Render(0, 0)
	// Only output the left and middle part of the progress bar
	fprintf(stdout, "%s%s", rendered.String(), end)
***REMOVED***

//nolint: funlen
func renderMultipleBars(
	isTTY, goBack bool, maxLeft, termWidth, widthDelta int, pbs []*pb.ProgressBar,
) (string, int) ***REMOVED***
	lineEnd := "\n"
	if isTTY ***REMOVED***
		//TODO: check for cross platform support
		lineEnd = "\x1b[K\n" // erase till end of line
	***REMOVED***

	var (
		// Amount of times line lengths exceed termWidth.
		// Needed to factor into the amount of lines to jump
		// back with [A and avoid scrollback issues.
		lineBreaks  int
		longestLine int
		// Maximum length of each right side column except last,
		// used to calculate the padding between columns.
		maxRColumnLen = make([]int, 2)
		pbsCount      = len(pbs)
		rendered      = make([]pb.ProgressBarRender, pbsCount)
		result        = make([]string, pbsCount+2)
	)

	result[0] = lineEnd // start with an empty line

	// First pass to render all progressbars and get the maximum
	// lengths of right-side columns.
	for i, pb := range pbs ***REMOVED***
		rend := pb.Render(maxLeft, widthDelta)
		for i := range rend.Right ***REMOVED***
			// Skip last column, since there's nothing to align after it (yet?).
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
			runeCount := utf8.RuneCountInString(rend.Hijack)
			lineBreaks += (runeCount - termPadding) / termWidth
			continue
		***REMOVED***
		var leftText, rightText string
		leftPadFmt := fmt.Sprintf("%%-%ds", maxLeft)
		leftText = fmt.Sprintf(leftPadFmt, rend.Left)
		for i := range rend.Right ***REMOVED***
			rpad := 0
			if len(maxRColumnLen) > i ***REMOVED***
				rpad = maxRColumnLen[i]
			***REMOVED***
			rightPadFmt := fmt.Sprintf(" %%-%ds", rpad+1)
			rightText += fmt.Sprintf(rightPadFmt, rend.Right[i])
		***REMOVED***
		// Get visible line length, without ANSI escape sequences (color)
		status := fmt.Sprintf(" %s ", rend.Status())
		line := leftText + status + rend.Progress() + rightText
		lineRuneCount := utf8.RuneCountInString(line)
		if lineRuneCount > longestLine ***REMOVED***
			longestLine = lineRuneCount
		***REMOVED***
		lineBreaks += (lineRuneCount - termPadding) / termWidth
		if !noColor ***REMOVED***
			rend.Color = true
			status = fmt.Sprintf(" %s ", rend.Status())
			line = fmt.Sprintf(leftPadFmt+"%s%s%s",
				rend.Left, status, rend.Progress(), rightText)
		***REMOVED***
		result[i+1] = line + lineEnd
	***REMOVED***

	if isTTY && goBack ***REMOVED***
		// Clear screen and go back to the beginning
		//TODO: check for cross platform support
		result[pbsCount+1] = fmt.Sprintf("\r\x1b[J\x1b[%dA", pbsCount+lineBreaks+1)
	***REMOVED*** else ***REMOVED***
		result[pbsCount+1] = lineEnd
	***REMOVED***

	return strings.Join(result, ""), longestLine
***REMOVED***

//TODO: show other information here?
//TODO: add a no-progress option that will disable these
//TODO: don't use global variables...
// nolint:funlen
func showProgress(
	ctx context.Context, conf Config,
	execScheduler *local.ExecutionScheduler, logger *logrus.Logger,
) ***REMOVED***
	if quiet || conf.HTTPDebug.Valid && conf.HTTPDebug.String != "" ***REMOVED***
		return
	***REMOVED***

	pbs := []*pb.ProgressBar***REMOVED***execScheduler.GetInitProgressBar()***REMOVED***
	for _, s := range execScheduler.GetExecutors() ***REMOVED***
		pbs = append(pbs, s.GetProgress())
	***REMOVED***

	termWidth, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil && stdoutTTY ***REMOVED***
		logger.WithError(err).Warn("error getting terminal size")
		termWidth = 80 // TODO: something safer, return error?
	***REMOVED***

	// Get the longest left side string length, to align progress bars
	// horizontally and trim excess text.
	var leftLen int64
	for _, pb := range pbs ***REMOVED***
		l := pb.Left()
		leftLen = lib.Max(int64(len(l)), leftLen)
	***REMOVED***
	// Limit to maximum left text length
	maxLeft := int(lib.Min(leftLen, maxLeftLength))

	var progressBarsLastRender []byte

	printProgressBars := func() ***REMOVED***
		_, _ = stdout.Writer.Write(progressBarsLastRender)
	***REMOVED***

	var widthDelta int
	// Default to responsive progress bars when in an interactive terminal
	renderProgressBars := func(goBack bool) ***REMOVED***
		barText, longestLine := renderMultipleBars(stdoutTTY, goBack, maxLeft, termWidth, widthDelta, pbs)
		widthDelta = termWidth - longestLine - termPadding
		progressBarsLastRender = []byte(barText)
	***REMOVED***

	// Otherwise fallback to fixed compact progress bars
	if !stdoutTTY ***REMOVED***
		widthDelta = -pb.DefaultWidth
		renderProgressBars = func(goBack bool) ***REMOVED***
			barText, _ := renderMultipleBars(stdoutTTY, goBack, maxLeft, termWidth, widthDelta, pbs)
			progressBarsLastRender = []byte(barText)
		***REMOVED***
	***REMOVED***

	//TODO: make configurable?
	updateFreq := 1 * time.Second
	//TODO: remove !noColor after we fix how we handle colors (see the related
	//description in the TODO message in cmd/root.go)
	if stdoutTTY && !noColor ***REMOVED***
		updateFreq = 100 * time.Millisecond
		outMutex.Lock()
		stdout.PersistentText = printProgressBars
		stderr.PersistentText = printProgressBars
		outMutex.Unlock()
		defer func() ***REMOVED***
			outMutex.Lock()
			stdout.PersistentText = nil
			stderr.PersistentText = nil
			outMutex.Unlock()
		***REMOVED***()
	***REMOVED***

	var (
		fd     = int(os.Stdout.Fd())
		ticker = time.NewTicker(updateFreq)
	)

	var winch chan os.Signal
	if sig := getWinchSignal(); sig != nil ***REMOVED***
		winch = make(chan os.Signal, 1)
		signal.Notify(winch, sig)
	***REMOVED***

	ctxDone := ctx.Done()
	for ***REMOVED***
		select ***REMOVED***
		case <-ctxDone:
			renderProgressBars(false)
			outMutex.Lock()
			printProgressBars()
			outMutex.Unlock()
			return
		case <-winch:
			// More responsive progress bar resizing on platforms with SIGWINCH (*nix)
			termWidth, _, _ = terminal.GetSize(fd)
		case <-ticker.C:
			// Default ticker-based progress bar resizing
			if winch == nil ***REMOVED***
				termWidth, _, _ = terminal.GetSize(fd)
			***REMOVED***
		***REMOVED***
		renderProgressBars(true)
		outMutex.Lock()
		printProgressBars()
		outMutex.Unlock()
	***REMOVED***
***REMOVED***
