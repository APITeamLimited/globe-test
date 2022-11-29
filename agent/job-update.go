package agent

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/APITeamLimited/globe-test/agent/libAgent"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"github.com/gobwas/ws/wsutil"
)

func handleJobUpdate(rawMessage []byte, conn *net.Conn, runningJobs *map[string]libOrch.Job,
	orchestratorClient *redis.Client) ***REMOVED***

	parsedMessage := libAgent.ClientJobUpdateMessage***REMOVED******REMOVED***

	err := json.Unmarshal(rawMessage, &parsedMessage)
	if err != nil ***REMOVED***
		wsutil.WriteServerText(*conn, []byte("Error parsing jobId"))
		return
	***REMOVED***

	// Ensure job exists in running jobs
	/*job*/
	_, ok := (*runningJobs)[parsedMessage.Message.JobId]
	if !ok ***REMOVED***
		wsutil.WriteServerText(*conn, []byte(fmt.Sprintf("Job does not exist with id %s", parsedMessage.Message.JobId)))
		return
	***REMOVED***

	// TODO: Implement this once updates are implemented

	panic("Not implemented")
***REMOVED***
