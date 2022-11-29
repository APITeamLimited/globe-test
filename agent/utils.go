package agent

import (
	"encoding/json"
	"net"

	"github.com/APITeamLimited/globe-test/agent/libAgent"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/gobwas/ws/wsutil"
)

func broadcastMessageToAll(message []byte, connections *map[string]*net.Conn) {
	for _, conn := range *connections {
		wsutil.WriteServerText(*conn, message)
	}
}

func sendRunningJobsToClient(conn *net.Conn, runningJobs *map[string]libOrch.Job) {
	// Make array of running jobs
	runningJobsArray := []libOrch.Job{}

	for _, job := range *runningJobs {
		runningJobsArray = append(runningJobsArray, job)
	}

	message := libAgent.ServerRunningJobsMessage{
		Type:    "runningJobs",
		Message: runningJobsArray,
	}

	marshalledMessage, err := json.Marshal(message)
	if err != nil {
		wsutil.WriteServerText(*conn, []byte("Error marshalling running jobs"))
		return
	}

	wsutil.WriteServerText(*conn, marshalledMessage)
}
