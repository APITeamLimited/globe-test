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
	orchestratorClient *redis.Client) {

	parsedMessage := libAgent.ClientJobUpdateMessage{}

	err := json.Unmarshal(rawMessage, &parsedMessage)
	if err != nil {
		wsutil.WriteServerText(*conn, []byte("Error parsing jobId"))
		return
	}

	// Ensure job exists in running jobs
	/*job*/
	_, ok := (*runningJobs)[parsedMessage.Message.JobId]
	if !ok {
		wsutil.WriteServerText(*conn, []byte(fmt.Sprintf("Job does not exist with id %s", parsedMessage.Message.JobId)))
		return
	}

	// TODO: Implement this once updates are implemented

	panic("Not implemented")
}
