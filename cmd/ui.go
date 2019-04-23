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
	"github.com/loadimpact/k6/ui/pb"
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
		//TODO: check how cross-platform this is...
		p = bytes.ReplaceAll(p, []byte***REMOVED***'\n'***REMOVED***, []byte***REMOVED***'\x1b', '[', '0', 'K', '\n'***REMOVED***)
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
		//TODO: check for cross platform support
		end = "\x1b[0K\r"
	***REMOVED***
	fprintf(stdout, "%s %s%s", bar.String(), rightText, end)
***REMOVED***

func renderMultipleBars(isTTY, goBack bool, pbs []*pb.ProgressBar) string ***REMOVED***
	lineEnd := "\n"
	if isTTY ***REMOVED***
		//TODO: check for cross platform support
		lineEnd = "\x1b[K\n" // erase till end of line
	***REMOVED***

	pbsCount := len(pbs)
	result := make([]string, pbsCount+2)
	result[0] = lineEnd // start with an empty line
	for i, pb := range pbs ***REMOVED***
		result[i+1] = pb.String() + lineEnd
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
func showProgress(ctx context.Context, wg *sync.WaitGroup, conf Config, executor *local.Executor) ***REMOVED***
	defer wg.Done()
	if quiet || conf.HttpDebug.Valid && conf.HttpDebug.String != "" ***REMOVED***
		return
	***REMOVED***

	pbs := []*pb.ProgressBar***REMOVED***executor.GetInitProgressBar()***REMOVED***
	for _, s := range executor.GetSchedulers() ***REMOVED***
		pbs = append(pbs, s.GetProgress())
	***REMOVED***

	// For flicker-free progressbars!
	progressBarsLastRender := []byte(renderMultipleBars(stdoutTTY, true, pbs))
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
				progressBarsLastRender = []byte(renderMultipleBars(stdoutTTY, false, pbs))
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
			barText := renderMultipleBars(stdoutTTY, true, pbs)
			outMutex.Lock()
			progressBarsLastRender = []byte(barText)
			progressBarsPrint()
			outMutex.Unlock()
		case <-ctxDone:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***
