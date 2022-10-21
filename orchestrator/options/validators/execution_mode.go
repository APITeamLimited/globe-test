package validators

import (
	"errors"
	"fmt"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

func ExecutionMode(options *libWorker.Options) error {
	if !options.ExecutionMode.Valid {
		return errors.New("execution mode is not valid")
	}

	if options.ExecutionMode.Value == types.HTTPSingleExecutionMode {
		if (options.VUs.Valid && (options.VUs.ValueOrZero() > 1)) || (options.Iterations.Valid && options.Iterations.ValueOrZero() > 1) {
			return errors.New("cannot use executionMode 'http_single' with more than 1 VU or iteration")
		} else if options.Stages != nil {
			return errors.New("cannot use executionMode 'http_single' with stages")
		} else if options.Scenarios != nil {
			return errors.New("cannot use executionMode 'http_single' with scenarios")
		}

		return nil
	} else if options.ExecutionMode.Value == types.HTTPMultipleExecutionMode {
		return nil
	}

	return fmt.Errorf("execution mode '%s' is not valid", options.ExecutionMode.Value)
}