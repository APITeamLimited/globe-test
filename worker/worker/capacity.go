package worker

import (
	"sync"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

type ExecutionList struct {
	currentJobs      map[string]libOrch.ChildJob
	mutex            sync.Mutex
	maxJobs          int
	maxVUs           int64
	currentJobsCount int
	currentVUsCount  int64
}

func (executionList *ExecutionList) addJob(job libOrch.ChildJob) {
	executionList.currentJobs[job.Id] = job

	executionList.currentJobsCount++
	executionList.currentVUsCount += job.Options.MaxPossibleVUs.ValueOrZero()
}

func (executionList *ExecutionList) removeJob(childJobId string) {
	executionList.mutex.Lock()
	managedVUsFreed := executionList.currentJobs[childJobId].Options.MaxPossibleVUs.ValueOrZero()

	executionList.currentVUsCount -= managedVUsFreed
	executionList.currentJobsCount--

	delete(executionList.currentJobs, childJobId)

	executionList.mutex.Unlock()
}

// Checks if the exectutor has the physical capacity to execute this job, this does
// not concern whether the user has the required credits to execute the job.
func (executionList *ExecutionList) checkExecutionCapacity(options libWorker.Options) bool {
	if executionList.maxJobs >= 0 && executionList.currentJobsCount >= executionList.maxJobs {
		return false
	}

	// If more than max permissible managed VUs, return false
	if executionList.maxVUs >= 0 && executionList.currentVUsCount+options.MaxPossibleVUs.ValueOrZero() > executionList.maxVUs {
		return false
	}

	return true
}
