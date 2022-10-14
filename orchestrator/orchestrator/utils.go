package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
)

func (e *ExecutionList) addJob(job libOrch.Job) ***REMOVED***
	e.mutex.Lock()
	e.currentJobs[job.Id] = job
	e.mutex.Unlock()
***REMOVED***

func (e *ExecutionList) removeJob(jobId string) ***REMOVED***
	e.mutex.Lock()
	delete(e.currentJobs, jobId)
	e.mutex.Unlock()
***REMOVED***

func fetchJob(ctx context.Context, orchestratorClient *redis.Client, jobId string) (*libOrch.Job, error) ***REMOVED***
	jobRaw, err := orchestratorClient.HGet(ctx, jobId, "job").Result()

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Check job not empty
	if jobRaw == "" ***REMOVED***
		return nil, fmt.Errorf("job %s is empty", jobId)
	***REMOVED***

	job := libOrch.Job***REMOVED******REMOVED***
	// Parse job as libOrch.Job
	err = json.Unmarshal([]byte(jobRaw), &job)
	if err != nil ***REMOVED***
		fmt.Println("error unmarshalling job", err)
		return nil, fmt.Errorf("error unmarshalling job %s", jobId)
	***REMOVED***

	return &job, nil
***REMOVED***
