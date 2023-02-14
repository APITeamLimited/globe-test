package validators

import (
	"errors"
	"fmt"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

func ExecutionMode(options *libWorker.Options, rootNodeVariant string) error {
	if !options.ExecutionMode.Valid {
		// If the execution mode is not set, we default to httpSingle if the root node variant is httpSingle
		if rootNodeVariant == libWorker.HTTPRequestVariant {
			options.ExecutionMode = types.NewNullExecutionMode(types.HTTPSingleExecutionMode, true)
		} else {
			options.ExecutionMode = types.NewNullExecutionMode(types.HTTPMultipleExecutionMode, true)
		}
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
