package validators

import (
	"errors"
	"fmt"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

func ExecutionMode(options *libWorker.Options) error ***REMOVED***
	if !options.ExecutionMode.Valid ***REMOVED***
		return errors.New("execution mode is not valid")
	***REMOVED***

	if options.ExecutionMode.Value == types.HTTPSingleExecutionMode ***REMOVED***
		if (options.VUs.Valid && (options.VUs.ValueOrZero() > 1)) || (options.Iterations.Valid && options.Iterations.ValueOrZero() > 1) ***REMOVED***
			return errors.New("cannot use executionMode 'http_single' with more than 1 VU or iteration")
		***REMOVED*** else if options.Stages != nil ***REMOVED***
			return errors.New("cannot use executionMode 'http_single' with stages")
		***REMOVED*** else if options.Scenarios != nil ***REMOVED***
			return errors.New("cannot use executionMode 'http_single' with scenarios")
		***REMOVED***

		return nil
	***REMOVED*** else if options.ExecutionMode.Value == types.HTTPMultipleExecutionMode ***REMOVED***
		return nil
	***REMOVED***

	return fmt.Errorf("execution mode '%s' is not valid", options.ExecutionMode.Value)
***REMOVED***
