package libWorker

import (
	"io"

	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
	"github.com/sirupsen/logrus"
)

// TestPreInitState contains all of the state that can be gathered and built
// before the test run is initialized.
type TestPreInitState struct {
	RuntimeOptions RuntimeOptions
	Registry       *workerMetrics.Registry
	BuiltinMetrics *workerMetrics.BuiltinMetrics
	KeyLogger      io.Writer

	Logger *logrus.Logger
}

// TestRunState contains the pre-init state as well as all of the state and
// options that are necessary for actually running the test.
type TestRunState struct {
	*TestPreInitState

	Options Options
	Runner  Runner // TODO: rename to something better, see type comment

	// TODO: add atlas root node

	// TODO: add other properties that are computed or derived after init, e.g.
	// thresholds?
}
