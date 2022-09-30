package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/APITeamLimited/redis/v9"
)

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

func fetchScope(ctx context.Context, scopesClient *redis.Client, scopeId string) (map[string]string, error) ***REMOVED***
	scope, err := scopesClient.Get(ctx, fmt.Sprintf("scope__id:%s", scopeId)).Result()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Check scope not empty
	if len(scope) == 0 ***REMOVED***
		return nil, fmt.Errorf("scope %s is empty", scopeId)
	***REMOVED***

	// Parse scope as map[string]string
	parsedScope := make(map[string]string)
	err = json.Unmarshal([]byte(scope), &parsedScope)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("error unmarshalling scope %s", scopeId)
	***REMOVED***

	return parsedScope, nil
***REMOVED***
