package agent

import (
	"encoding/json"
	"net"

	"github.com/APITeamLimited/globe-test/agent/libAgent"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/gobwas/ws/wsutil"
)

func broadcastMessageToAll(message []byte, connections *map[string]*net.Conn) ***REMOVED***
	for _, conn := range *connections ***REMOVED***
		wsutil.WriteServerText(*conn, message)
	***REMOVED***
***REMOVED***

func sendRunningJobsToClient(conn *net.Conn, runningJobs *map[string]libOrch.Job) ***REMOVED***
	// Make array of running jobs
	runningJobsArray := []libOrch.Job***REMOVED******REMOVED***

	for _, job := range *runningJobs ***REMOVED***
		runningJobsArray = append(runningJobsArray, job)
	***REMOVED***

	message := libAgent.ServerRunningJobsMessage***REMOVED***
		Type:    "runningJobs",
		Message: runningJobsArray,
	***REMOVED***

	marshalledMessage, err := json.Marshal(message)
	if err != nil ***REMOVED***
		wsutil.WriteServerText(*conn, []byte("Error marshalling running jobs"))
		return
	***REMOVED***

	wsutil.WriteServerText(*conn, marshalledMessage)
***REMOVED***
