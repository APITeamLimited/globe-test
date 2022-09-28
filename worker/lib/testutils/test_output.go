package testutils

import (
	"io"
	"testing"

	"github.com/sirupsen/logrus"
)

// Something that makes the test also be a valid io.Writer, useful for passing it
// as an output for logs and CLI flag help messages...
type testOutput struct***REMOVED*** testing.TB ***REMOVED***

func (to testOutput) Write(p []byte) (n int, err error) ***REMOVED***
	to.Logf("%s", p)

	return len(p), nil
***REMOVED***

// NewTestOutput returns a simple io.Writer implementation that uses the test's
// logger as an output.
func NewTestOutput(t testing.TB) io.Writer ***REMOVED***
	return testOutput***REMOVED***t***REMOVED***
***REMOVED***

// NewLogger Returns new logger that will log to the testing.TB.Logf
func NewLogger(t testing.TB) *logrus.Logger ***REMOVED***
	l := logrus.New()
	logrus.SetOutput(NewTestOutput(t))

	return l
***REMOVED***
