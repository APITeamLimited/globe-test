package orchestrator

import (
	"errors"
	"sync"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

type ExecutionList struct {
	currentJobs            map[string]libOrch.Job
	mutex                  sync.Mutex
	maxJobs                int
	maxManagedVUs          int64
	currentJobsCount       int
	currentManagedVUsCount int64
}

// addJob assumes that the execution list is already locked
func (executionList *ExecutionList) addJob(job *libOrch.Job) error {
	if job.Options == nil {
		return errors.New("job options should not be nil")
	}

	executionList.currentJobs[job.Id] = *job

	executionList.currentJobsCount++

	if job.Options != nil {
		executionList.currentManagedVUsCount += job.Options.MaxPossibleVUs.ValueOrZero()
	}

	return nil
}

func (executionList *ExecutionList) removeJob(jobId string) {
	executionList.mutex.Lock()
	defer executionList.mutex.Unlock()

	job := executionList.currentJobs[jobId]

	if job.Options != nil {
		managedVUsFreed := job.Options.MaxPossibleVUs.ValueOrZero()
		executionList.currentManagedVUsCount -= managedVUsFreed
	}

	delete(executionList.currentJobs, jobId)
	executionList.currentJobsCount--
}

// Checks if the exectutor has the physical capacity to execute this job, this does
// not concern whether the user has the required credits to execute the job.
func (executionList *ExecutionList) checkExecutionCapacity(options *libWorker.Options) bool {
	// If more than max permissible jobs, return false
	if executionList.maxJobs >= 0 && executionList.currentJobsCount >= executionList.maxJobs {
		return false
	}

	if options == nil {
		return true
	}

	// If more than max permissible managed VUs, return false
	if executionList.maxManagedVUs >= 0 && executionList.currentManagedVUsCount+options.MaxPossibleVUs.ValueOrZero() > executionList.maxManagedVUs {
		return false
	}

	return true
}
