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

			getCwd := func() (string, error) ***REMOVED***
				return "/", nil
			***REMOVED***

			res, err := FileHookFromConfigLine(
				context.Background(), afero.NewMemMapFs(), getCwd, logrus.New(), test.line, make(chan struct***REMOVED******REMOVED***),
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
