package function

import (
	"github.com/APITeamLimited/globe-test/worker/worker"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("WorkerCloud", worker.RunGoogleCloud)
}
