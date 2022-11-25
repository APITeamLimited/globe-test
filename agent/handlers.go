package agent

import (
	"encoding/json"

	socketio "github.com/googollee/go-socket.io"
)

func abortAllJobs(runningJobs map[string]runningJob, setJobCount func(int)) {
	// Loop through all running jobs and cancel them
	for _, job := range runningJobs {
		abortJob(job.JobId, runningJobs, setJobCount)
	}
	setJobCount(len(runningJobs))
}

func abortJob(jobId string, runningJobs map[string]runningJob, setJobCount func(int)) {
	_, ok := runningJobs[jobId]
	if ok {
		delete(runningJobs, jobId)
	}
	setJobCount(len(runningJobs))
}

func handleJobCreation(msg string, s socketio.Conn, runningJobs map[string]runningJob, setJobCount func(int)) {
	// Parse the message
	var args jobArgs
	err := json.Unmarshal([]byte(msg), &args)
	if err != nil {
		s.Emit("error", err.Error())
		return
	}

	// Create a new job
	runningJobs[args.JobId] = runningJob{
		JobId:      args.JobId,
		Source:     args.Source,
		SourceName: args.SourceName,
		Messages:   []string{},
	}
	5
}
