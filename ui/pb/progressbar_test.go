package pb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO(imiric): Consider adding logging tests for 100% pb coverage.
// Unfortunately the following introduces an import cycle: pb -> lib -> pb
// func getTestLogger() *logger.Entry ***REMOVED***
// 	logHook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.WarnLevel***REMOVED******REMOVED***
// 	testLog := logrus.New()
// 	testLog.AddHook(logHook)
// 	testLog.SetOutput(ioutil.Discard)
// 	return logrus.NewEntry(testLog)
// ***REMOVED***

func TestProgressBarRender(t *testing.T) ***REMOVED***
	t.Parallel()

	testCases := []struct ***REMOVED***
		options      []ProgressBarOption
		pbWidthDelta int
		expected     string
	***REMOVED******REMOVED***
		***REMOVED***
			[]ProgressBarOption***REMOVED***WithLeft(func() string ***REMOVED*** return "left" ***REMOVED***)***REMOVED***,
			0, "left   [--------------------------------------]",
		***REMOVED***,
		***REMOVED***
			[]ProgressBarOption***REMOVED***WithConstLeft("constLeft")***REMOVED***,
			0, "constLeft   [--------------------------------------]",
		***REMOVED***,
		***REMOVED***[]ProgressBarOption***REMOVED***
			WithLeft(func() string ***REMOVED*** return "left" ***REMOVED***),
			WithStatus(Done),
		***REMOVED***, 0, "left ✓ [--------------------------------------]"***REMOVED***,
		***REMOVED***[]ProgressBarOption***REMOVED***
			WithLeft(func() string ***REMOVED*** return "left" ***REMOVED***),
			WithProgress(func() (float64, []string) ***REMOVED*** return 0, []string***REMOVED***"right"***REMOVED*** ***REMOVED***),
		***REMOVED***, 0, "left   [--------------------------------------] right"***REMOVED***,
		***REMOVED***[]ProgressBarOption***REMOVED***
			WithLeft(func() string ***REMOVED*** return "left" ***REMOVED***),
			WithProgress(func() (float64, []string) ***REMOVED*** return 0.5, []string***REMOVED***"right"***REMOVED*** ***REMOVED***),
		***REMOVED***, 0, "left   [==================>-------------------] right"***REMOVED***,
		***REMOVED***[]ProgressBarOption***REMOVED***
			WithLeft(func() string ***REMOVED*** return "left" ***REMOVED***),
			WithProgress(func() (float64, []string) ***REMOVED*** return 1.0, []string***REMOVED***"right"***REMOVED*** ***REMOVED***),
		***REMOVED***, 0, "left   [======================================] right"***REMOVED***,
		***REMOVED***[]ProgressBarOption***REMOVED***
			WithLeft(func() string ***REMOVED*** return "left" ***REMOVED***),
			WithProgress(func() (float64, []string) ***REMOVED*** return -1, []string***REMOVED***"right"***REMOVED*** ***REMOVED***),
		***REMOVED***, 0, "left   [--------------------------------------] right"***REMOVED***,
		***REMOVED***[]ProgressBarOption***REMOVED***
			WithLeft(func() string ***REMOVED*** return "left" ***REMOVED***),
			WithProgress(func() (float64, []string) ***REMOVED*** return 2, []string***REMOVED***"right"***REMOVED*** ***REMOVED***),
		***REMOVED***, 0, "left   [======================================] right"***REMOVED***,
		***REMOVED***[]ProgressBarOption***REMOVED***
			WithLeft(func() string ***REMOVED*** return "left" ***REMOVED***),
			WithConstProgress(0.2, "constProgress"),
		***REMOVED***, 0, "left   [======>-------------------------------] constProgress"***REMOVED***,
		***REMOVED***[]ProgressBarOption***REMOVED***
			WithHijack(func() string ***REMOVED*** return "progressbar hijack!" ***REMOVED***),
		***REMOVED***, 0, "progressbar hijack!"***REMOVED***,
		***REMOVED***
			[]ProgressBarOption***REMOVED***WithConstProgress(0.25, "")***REMOVED***,
			-DefaultWidth, "   [  25% ] ",
		***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(tc.expected, func(t *testing.T) ***REMOVED***
			t.Parallel()
			pbar := New(tc.options...)
			assert.NotNil(t, pbar)
			assert.Equal(t, tc.expected, pbar.Render(0, tc.pbWidthDelta).String())
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestProgressBarRenderPaddingMaxLeft(t *testing.T) ***REMOVED***
	t.Parallel()

	testCases := []struct ***REMOVED***
		maxLen   int
		left     string
		expected string
	***REMOVED******REMOVED***
		***REMOVED***-1, "left", "left   [--------------------------------------]"***REMOVED***,
		***REMOVED***0, "left", "left   [--------------------------------------]"***REMOVED***,
		***REMOVED***10, "left_truncated", "left_tr...   [--------------------------------------]"***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(tc.left, func(t *testing.T) ***REMOVED***
			t.Parallel()
			pbar := New(WithLeft(func() string ***REMOVED*** return tc.left ***REMOVED***))
			assert.NotNil(t, pbar)
			assert.Equal(t, tc.expected, pbar.Render(tc.maxLen, 0).String())
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestProgressBarLeft(t *testing.T) ***REMOVED***
	t.Parallel()

	testCases := []struct ***REMOVED***
		left     func() string
		expected string
	***REMOVED******REMOVED***
		***REMOVED***nil, ""***REMOVED***,
		***REMOVED***func() string ***REMOVED*** return " left " ***REMOVED***, " left "***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(tc.expected, func(t *testing.T) ***REMOVED***
			t.Parallel()
			pbar := New(WithLeft(tc.left))
			assert.NotNil(t, pbar)
			assert.Equal(t, tc.expected, pbar.Left())
		***REMOVED***)
	***REMOVED***
***REMOVED***
