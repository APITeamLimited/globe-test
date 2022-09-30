package validators

import (
	"errors"
	"fmt"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

func ExecutionMode(options *libWorker.Options) error {
	if options.ExecutionMode.Value == types.RESTSingleExecutionMode {
		if (options.VUs.Valid && (options.VUs.ValueOrZero() > 1)) || (options.Iterations.Valid && options.Iterations.ValueOrZero() > 1) {
			return errors.New("cannot use execution mode 'rest_single' with more than 1 VU or iteration")
		} else if options.Stages != nil {
			return errors.New("cannot use execution mode 'rest_single' with stages")
		} else if options.Scenarios != nil {
			return errors.New("cannot use execution mode 'rest_single' with scenarios")
		}
	} else if options.ExecutionMode.Valid && options.ExecutionMode.Value != types.RESTMultipleExecutionMode {
		return fmt.Errorf("invalid execution mode '%s'", options.ExecutionMode.Value)
	} else if !options.ExecutionMode.Valid {
		options.ExecutionMode = types.NullExecutionModeFrom(types.RESTMultipleExecutionMode)
	}
	return nil
}
