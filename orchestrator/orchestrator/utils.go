package orchestrator

func (e *ExecutionList) addJob(job map[string]string) {
	e.mutex.Lock()
	e.currentJobs[job["id"]] = job
	e.mutex.Unlock()
}

func (e *ExecutionList) removeJob(jobId string) {
	e.mutex.Lock()
	delete(e.currentJobs, jobId)
	e.mutex.Unlock()
}
