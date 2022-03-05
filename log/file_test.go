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

package log

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type nopCloser struct ***REMOVED***
	io.Writer
	closed chan struct***REMOVED******REMOVED***
***REMOVED***

func (nc *nopCloser) Close() error ***REMOVED***
	nc.closed <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
	return nil
***REMOVED***

func TestFileHookFromConfigLine(t *testing.T) ***REMOVED***
	t.Parallel()

	tests := [...]struct ***REMOVED***
		line       string
		err        bool
		errMessage string
		res        fileHook
	***REMOVED******REMOVED***
		***REMOVED***
			line: "file",
			err:  true,
			res: fileHook***REMOVED***
				levels: logrus.AllLevels,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			line: "file=/k6.log,level=info",
			err:  false,
			res: fileHook***REMOVED***
				path:   "/k6.log",
				levels: logrus.AllLevels[:5],
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			line: "file=/a/c/",
			err:  true,
		***REMOVED***,
		***REMOVED***
			line:       "file=,level=info",
			err:        true,
			errMessage: "filepath must not be empty",
		***REMOVED***,
		***REMOVED***
			line: "file=/tmp/k6.log,level=tea",
			err:  true,
		***REMOVED***,
		***REMOVED***
			line: "file=/tmp/k6.log,unknown",
			err:  true,
		***REMOVED***,
		***REMOVED***
			line: "file=/tmp/k6.log,level=",
			err:  true,
		***REMOVED***,
		***REMOVED***
			line: "file=/tmp/k6.log,level=,",
			err:  true,
		***REMOVED***,
		***REMOVED***
			line:       "file=/tmp/k6.log,unknown=something",
			err:        true,
			errMessage: "unknown logfile config key unknown",
		***REMOVED***,
		***REMOVED***
			line:       "unknown=something",
			err:        true,
			errMessage: "logfile configuration should be in the form `file=path-to-local-file` but is `unknown=something`",
		***REMOVED***,
	***REMOVED***

	for _, test := range tests ***REMOVED***
		test := test
		t.Run(test.line, func(t *testing.T) ***REMOVED***
			t.Parallel()

			res, err := FileHookFromConfigLine(
				context.Background(), afero.NewMemMapFs(), logrus.New(), test.line, make(chan struct***REMOVED******REMOVED***),
			)

			if test.err ***REMOVED***
				require.Error(t, err)

				if test.errMessage != "" ***REMOVED***
					require.Equal(t, test.errMessage, err.Error())
				***REMOVED***

				return
			***REMOVED***

			require.NoError(t, err)
			assert.NotNil(t, res.(*fileHook).w)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestFileHookFire(t *testing.T) ***REMOVED***
	t.Parallel()

	var buffer bytes.Buffer
	nc := &nopCloser***REMOVED***
		Writer: &buffer,
		closed: make(chan struct***REMOVED******REMOVED***),
	***REMOVED***

	hook := &fileHook***REMOVED***
		loglines: make(chan []byte),
		w:        nc,
		bw:       bufio.NewWriter(nc),
		levels:   logrus.AllLevels,
		done:     make(chan struct***REMOVED******REMOVED***),
	***REMOVED***

	ctx, cancel := context.WithCancel(context.Background())

	hook.loglines = hook.loop(ctx)

	logger := logrus.New()
	logger.AddHook(hook)
	logger.SetOutput(io.Discard)

	logger.Info("example log line")

	time.Sleep(10 * time.Millisecond)

	cancel()
	<-nc.closed

	assert.Contains(t, buffer.String(), "example log line")
***REMOVED***
