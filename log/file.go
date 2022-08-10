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
	ctx context.Context, fs afero.Fs, getCwd func() (string, error),
	fallbackLogger logrus.FieldLogger, line string, done chan struct***REMOVED******REMOVED***,
) (logrus.Hook, error) ***REMOVED***
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

	if err := hook.openFile(getCwd); err != nil ***REMOVED***
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
func (h *fileHook) openFile(getCwd func() (string, error)) error ***REMOVED***
	path := h.path
	if !filepath.IsAbs(path) ***REMOVED***
		cwd, err := getCwd()
		if err != nil ***REMOVED***
			return fmt.Errorf("'%s' is a relative path but could not determine CWD: %w", path, err)
		***REMOVED***
		path = filepath.Join(cwd, path)
	***REMOVED***

	if _, err := h.fs.Stat(filepath.Dir(path)); os.IsNotExist(err) ***REMOVED***
		return fmt.Errorf("provided directory '%s' does not exist", filepath.Dir(path))
	***REMOVED***

	file, err := h.fs.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o600)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to open logfile %s: %w", path, err)
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
