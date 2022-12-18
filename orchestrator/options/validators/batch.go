package validators

import (
	"fmt"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

const maxBatchSize = 20

func Batch(options *libWorker.Options) error {
	// Ensure that the batch size is not too large
	if options.Batch.Int64 > maxBatchSize {
		return fmt.Errorf("batch size %d is too large, the maximum is %d", options.Batch.Int64, maxBatchSize)
	}

	return nil
}
