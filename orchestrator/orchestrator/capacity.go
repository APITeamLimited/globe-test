package orchestrator

import (
	"errors"
	"sync"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

type ExecutionList struct ***REMOVED***
	currentJobs            map[string]libOrch.Job
	mutex                  sync.Mutex
	maxJobs                int
	maxManagedVUs          int64
	currentJobsCount       int
	currentManagedVUsCount int64
***REMOVED***

// addJob assumes that the execution list is already locked
func (executionList *ExecutionList) addJob(job *libOrch.Job) error ***REMOVED***
	if job.Options == nil ***REMOVED***
		return errors.New("job options should not be nil")
	***REMOVED***

	executionList.currentJobs[job.Id] = *job

	executionList.currentJobsCount++

	if job.Options != nil ***REMOVED***
		executionList.currentManagedVUsCount += job.Options.MaxPossibleVUs.ValueOrZero()
	***REMOVED***

	return nil
***REMOVED***

func (executionList *ExecutionList) removeJob(jobId string) ***REMOVED***
	executionList.mutex.Lock()
	defer executionList.mutex.Unlock()

	job := executionList.currentJobs[jobId]

	if job.Options != nil ***REMOVED***
		managedVUsFreed := job.Options.MaxPossibleVUs.ValueOrZero()
		executionList.currentManagedVUsCount -= managedVUsFreed
	***REMOVED***

	delete(executionList.currentJobs, jobId)
	executionList.currentJobsCount--
***REMOVED***

// Checks if the exectutor has the physical capacity to execute this job, this does
// not concern whether the user has the required credits to execute the job.
func (executionList *ExecutionList) checkExecutionCapacity(options *libWorker.Options) bool ***REMOVED***
	// If more than max permissible jobs, return false
	if executionList.maxJobs >= 0 && executionList.currentJobsCount >= executionList.maxJobs ***REMOVED***
		return false
	***REMOVED***

	if options == nil ***REMOVED***
		return true
	***REMOVED***

	// If more than max permissible managed VUs, return false
	if executionList.maxManagedVUs >= 0 && executionList.currentManagedVUsCount+options.MaxPossibleVUs.ValueOrZero() > executionList.maxManagedVUs ***REMOVED***
		return false
	***REMOVED***

	return true
***REMOVED***
