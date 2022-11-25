package agent

import (
	"encoding/json"

	socketio "github.com/googollee/go-socket.io"
)

func abortAllJobs(runningJobs map[string]runningJob, setJobCount func(int)) ***REMOVED***
	// Loop through all running jobs and cancel them
	for _, job := range runningJobs ***REMOVED***
		abortJob(job.JobId, runningJobs, setJobCount)
	***REMOVED***
	setJobCount(len(runningJobs))
***REMOVED***

func abortJob(jobId string, runningJobs map[string]runningJob, setJobCount func(int)) ***REMOVED***
	_, ok := runningJobs[jobId]
	if ok ***REMOVED***
		delete(runningJobs, jobId)
	***REMOVED***
	setJobCount(len(runningJobs))
***REMOVED***

func handleJobCreation(msg string, s socketio.Conn, runningJobs map[string]runningJob, setJobCount func(int)) ***REMOVED***
	// Parse the message
	var args jobArgs
	err := json.Unmarshal([]byte(msg), &args)
	if err != nil ***REMOVED***
		s.Emit("error", err.Error())
		return
	***REMOVED***

	// Create a new job
	runningJobs[args.JobId] = runningJob***REMOVED***
		JobId:      args.JobId,
		Source:     args.Source,
		SourceName: args.SourceName,
		Messages:   []string***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***
