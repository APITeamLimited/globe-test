package validators

import (
	"errors"
	"fmt"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

func ExecutionMode(options *libWorker.Options) error {
	if !options.ExecutionMode.Valid {
		return errors.New("the option executionMode must be specified")
	}

	if options.ExecutionMode.Value == types.HTTPSingleExecutionMode {
		if (options.VUs.Valid && (options.VUs.ValueOrZero() > 1)) || (options.Iterations.Valid && options.Iterations.ValueOrZero() > 1) {
			return errors.New("cannot use executionMode 'httpSingle' with more than 1 VU or iteration")
		} else if options.Stages != nil {
			return errors.New("cannot use executionMode 'httpSingle' with stages")
		} else if options.Scenarios != nil {
			return errors.New("cannot use executionMode 'httpSingle' with scenarios")
		}

		return nil
	} else if options.ExecutionMode.Value == types.HTTPMultipleExecutionMode {
		return nil
	}

	return fmt.Errorf("execution mode '%s' is not valid", options.ExecutionMode.Value)
}
