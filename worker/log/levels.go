package log

import (
	"fmt"
	"sort"

	"github.com/sirupsen/logrus"
)

func parseLevels(level string) ([]logrus.Level, error) ***REMOVED***
	lvl, err := logrus.ParseLevel(level)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("unknown log level %s", level) // specifically use a custom error
	***REMOVED***
	index := sort.Search(len(logrus.AllLevels), func(i int) bool ***REMOVED***
		return logrus.AllLevels[i] > lvl
	***REMOVED***)

	return logrus.AllLevels[:index], nil
***REMOVED***
