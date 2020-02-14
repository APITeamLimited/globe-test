package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/loadimpact/k6/ui/pb"
)

// Return progressbars with different content lengths, to test for
// padding.
func createTestProgressBars(num, padding, colIdx int) []*pb.ProgressBar ***REMOVED***
	pbs := make([]*pb.ProgressBar, num)
	for i := 0; i < num; i++ ***REMOVED***
		left := fmt.Sprintf("left %d", i)
		rightCol1 := fmt.Sprintf("right %d", i)
		progress := 0.0
		status := pb.Running
		if i == colIdx ***REMOVED***
			pad := strings.Repeat("+", padding)
			left += pad
			rightCol1 += pad
			progress = 1.0
			status = pb.Done
		***REMOVED***
		pbs[i] = pb.New(
			pb.WithLeft(func() string ***REMOVED*** return left ***REMOVED***),
			pb.WithStatus(status),
			pb.WithProgress(func() (float64, []string) ***REMOVED***
				return progress, []string***REMOVED***rightCol1, "000"***REMOVED***
			***REMOVED***),
		)
	***REMOVED***
	return pbs
***REMOVED***

func TestRenderMultipleBars(t *testing.T) ***REMOVED***
	t.Parallel()

	testCases := []struct ***REMOVED***
		name     string
		padding  int
		expected string
	***REMOVED******REMOVED***
		***REMOVED***"nopad", 0, `
left 0   [--------------------------------------] right 0  000
left 1 ✓ [======================================] right 1  000
left 2   [--------------------------------------] right 2  000

`***REMOVED***,
		***REMOVED***"pad2", 2, `
left 0     [--------------------------------------] right 0    000
left 1++ ✓ [======================================] right 1++  000
left 2     [--------------------------------------] right 2    000

`***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(tc.name, func(t *testing.T) ***REMOVED***
			pbs := createTestProgressBars(3, tc.padding, 1)
			out := renderMultipleBars(false, false, 6+tc.padding, pbs)
			assert.Equal(t, tc.expected, out)
		***REMOVED***)
	***REMOVED***
***REMOVED***
