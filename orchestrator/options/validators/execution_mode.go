package validators

import (
	"errors"
	"fmt"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

func ExecutionMode(options *libWorker.Options) error ***REMOVED***
	if options.ExecutionMode.Value == types.RESTSingleExecutionMode ***REMOVED***
		if (options.VUs.Valid && (options.VUs.ValueOrZero() > 1)) || (options.Iterations.Valid && options.Iterations.ValueOrZero() > 1) ***REMOVED***
			return errors.New("cannot use execution mode 'rest_single' with more than 1 VU or iteration")
		***REMOVED*** else if options.Stages != nil ***REMOVED***
			return errors.New("cannot use execution mode 'rest_single' with stages")
		***REMOVED*** else if options.Scenarios != nil ***REMOVED***
			return errors.New("cannot use execution mode 'rest_single' with scenarios")
		***REMOVED***
	***REMOVED*** else if options.ExecutionMode.Valid && options.ExecutionMode.Value != types.RESTMultipleExecutionMode ***REMOVED***
		return fmt.Errorf("invalid execution mode '%s'", options.ExecutionMode.Value)
	***REMOVED*** else if !options.ExecutionMode.Valid ***REMOVED***
		options.ExecutionMode = types.NullExecutionModeFrom(types.RESTMultipleExecutionMode)
	***REMOVED***
	return nil
***REMOVED***
