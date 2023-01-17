package function

import (
	"net/http"

	"github.com/APITeamLimited/globe-test/worker/worker"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	runFunction := func(w http.ResponseWriter, r *http.Request) {
		worker.RunWorkerFunction(w, r, false)
	}

	functions.HTTP("WorkerCloud", runFunction)
}
