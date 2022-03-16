/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2020 Load Impact
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

// Package log implements various logrus hooks.
package log

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// fileHookBufferSize is a default size for the fileHook's loglines channel.
const fileHookBufferSize = 100

// fileHook is a hook to handle writing to local files.
type fileHook struct ***REMOVED***
	fs             afero.Fs
	fallbackLogger logrus.FieldLogger
	loglines       chan []byte
	path           string
	w              io.WriteCloser
	bw             *bufio.Writer
	levels         []logrus.Level
	done           chan struct***REMOVED******REMOVED***
***REMOVED***

// FileHookFromConfigLine returns new fileHook hook.
func FileHookFromConfigLine(
	ctx context.Context, fs afero.Fs, fallbackLogger logrus.FieldLogger, line string, done chan struct***REMOVED******REMOVED***,
) (logrus.Hook, error) ***REMOVED***
	// TODO: fix this so it works correctly with relative paths from the CWD

	hook := &fileHook***REMOVED***
		fs:             fs,
		fallbackLogger: fallbackLogger,
		levels:         logrus.AllLevels,
		done:           done,
	***REMOVED***

	parts := strings.SplitN(line, "=", 2)
	if parts[0] != "file" ***REMOVED***
		return nil, fmt.Errorf("logfile configuration should be in the form `file=path-to-local-file` but is `%s`", line)
	***REMOVED***

	if err := hook.parseArgs(line); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := hook.openFile(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	hook.loglines = hook.loop(ctx)

	return hook, nil
***REMOVED***

func (h *fileHook) parseArgs(line string) error ***REMOVED***
	tokens, err := tokenize(line)
	if err != nil ***REMOVED***
		return fmt.Errorf("error while parsing logfile configuration %w", err)
	***REMOVED***

	for _, token := range tokens ***REMOVED***
		switch token.key ***REMOVED***
		case "file":
			if token.value == "" ***REMOVED***
				return fmt.Errorf("filepath must not be empty")
			***REMOVED***
			h.path = token.value
		case "level":
			h.levels, err = parseLevels(token.value)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		default:
			return fmt.Errorf("unknown logfile config key %s", token.key)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// openFile opens logfile and initializes writers.
func (h *fileHook) openFile() error ***REMOVED***
	if _, err := h.fs.Stat(filepath.Dir(h.path)); os.IsNotExist(err) ***REMOVED***
		return fmt.Errorf("provided directory '%s' does not exist", filepath.Dir(h.path))
	***REMOVED***

	file, err := h.fs.OpenFile(h.path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o600)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to open logfile %s: %w", h.path, err)
	***REMOVED***

	h.w = file
	h.bw = bufio.NewWriter(file)

	return nil
***REMOVED***

func (h *fileHook) loop(ctx context.Context) chan []byte ***REMOVED***
	loglines := make(chan []byte, fileHookBufferSize)

	go func() ***REMOVED***
		defer close(h.done)
		for ***REMOVED***
			select ***REMOVED***
			case entry := <-loglines:
				if _, err := h.bw.Write(entry); err != nil ***REMOVED***
					h.fallbackLogger.Errorf("failed to write a log message to a logfile: %w", err)
				***REMOVED***
			case <-ctx.Done():
				if err := h.bw.Flush(); err != nil ***REMOVED***
					h.fallbackLogger.Errorf("failed to flush buffer: %w", err)
				***REMOVED***

				if err := h.w.Close(); err != nil ***REMOVED***
					h.fallbackLogger.Errorf("failed to close logfile: %w", err)
				***REMOVED***

				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return loglines
***REMOVED***

// Fire writes the log file to defined path.
func (h *fileHook) Fire(entry *logrus.Entry) error ***REMOVED***
	message, err := entry.Bytes()
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to get a log entry bytes: %w", err)
	***REMOVED***

	h.loglines <- message
	return nil
***REMOVED***

// Levels returns configured log levels.
func (h *fileHook) Levels() []logrus.Level ***REMOVED***
	return h.levels
***REMOVED***
