package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/APITeamLimited/globe-test/agent/libAgent"
	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"github.com/gobwas/ws/wsutil"
)

func handleAbortJob(rawMessage []byte, conn *net.Conn, runningJobs *map[string]libOrch.Job,
	setJobCount func(int), orchestratorClient *redis.Client) ***REMOVED***
	// Parse the rawMessage
	parsedMessage := libAgent.ClientAbortJobMessage***REMOVED******REMOVED***
	err := json.Unmarshal(rawMessage, &parsedMessage)
	if err != nil ***REMOVED***
		wsutil.WriteServerText(*conn, []byte("Error parsing jobId"))
		return
	***REMOVED***

	// Ensure job exists in running jobs
	job, ok := (*runningJobs)[parsedMessage.Message]
	if !ok ***REMOVED***
		wsutil.WriteServerText(*conn, []byte(fmt.Sprintf("Job does not exist with id %s", parsedMessage.Message)))
		return
	***REMOVED***

	// Abort the job
	processAbortion(job, runningJobs, setJobCount, orchestratorClient)
	displaySuccessMessage(conn, "Stopping test run")
***REMOVED***

func processAbortion(job libOrch.Job, runningJobs *map[string]libOrch.Job, setJobCount func(int),
	orchestratorClient *redis.Client) ***REMOVED***
	_, ok := (*runningJobs)[job.Id]
	if !ok ***REMOVED***
		fmt.Println("Attempted to abort job that does not exist")
		return
	***REMOVED***

	cancelMessage := lib.JobUserUpdate***REMOVED***
		UpdateType: "CANCEL",
	***REMOVED***

	marshalledCancel, _ := json.Marshal(cancelMessage)
	orchestratorClient.Publish(context.Background(), fmt.Sprintf("jobUserUpdates:%s:%s:%s", job.Scope.Variant, job.Scope.VariantTargetId, job.Id), string(marshalledCancel))
***REMOVED***
