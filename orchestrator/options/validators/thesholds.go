package validators

import (
	"github.com/APITeamLimited/globe-test/worker/errext"
	"github.com/APITeamLimited/globe-test/worker/errext/exitcodes"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/metrics"
)

func Thresholds(options *libWorker.Options) error {
	dummyRegistry := metrics.NewRegistry()
	metrics.RegisterBuiltinMetrics(dummyRegistry)

	for metricName, thresholdsDefinition := range options.Thresholds {
		err := thresholdsDefinition.Parse()
		if err != nil {
			return errext.WithExitCodeIfNone(err, exitcodes.InvalidConfig)
		}

		err = thresholdsDefinition.Validate(metricName, dummyRegistry)
		if err != nil {
			return errext.WithExitCodeIfNone(err, exitcodes.InvalidConfig)
		}

		options.Thresholds[metricName] = thresholdsDefinition
	}

	return nil
}
