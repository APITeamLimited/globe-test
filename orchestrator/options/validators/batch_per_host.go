package validators

import (
	"fmt"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

const maxBatchPerHostSize = 20

func BatchPerHost(options *libWorker.Options) error {
	// Ensure that the batch size is not too large
	if options.BatchPerHost.Int64 > maxBatchPerHostSize {
		return fmt.Errorf("batchPerHost size %d is too large, the maximum is %d", options.BatchPerHost.Int64, maxBatchPerHostSize)
	}

	return nil
}
