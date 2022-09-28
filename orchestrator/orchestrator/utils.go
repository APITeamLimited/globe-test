package orchestrator

func (e *ExecutionList) addJob(job map[string]string) ***REMOVED***
	e.mutex.Lock()
	e.currentJobs[job["id"]] = job
	e.mutex.Unlock()
***REMOVED***

func (e *ExecutionList) removeJob(jobId string) ***REMOVED***
	e.mutex.Lock()
	delete(e.currentJobs, jobId)
	e.mutex.Unlock()
***REMOVED***
