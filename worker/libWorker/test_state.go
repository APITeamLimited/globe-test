package libWorker

import (
	"github.com/APITeamLimited/globe-test/metrics"
	"github.com/sirupsen/logrus"
)

// TestPreInitState contains all of the state that can be gathered and built
// before the test run is initialized.
type TestPreInitState struct {
	BuiltinMetrics *metrics.BuiltinMetrics

	Logger *logrus.Logger
}

// TestRunState contains the pre-init state as well as all of the state and
// options that are necessary for actually running the test.
type TestRunState struct {
	*TestPreInitState

	Options Options
	Runner  Runner
}
