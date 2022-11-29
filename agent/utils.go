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

func displayErrorMessage(conn *net.Conn, message string) ***REMOVED***
	errorMessage := libAgent.ServerDisplayableErrorMessage***REMOVED***
		Type:    "displayableErrorMessage",
		Message: message,
	***REMOVED***

	marshalledMessage, err := json.Marshal(errorMessage)
	if err != nil ***REMOVED***
		wsutil.WriteServerText(*conn, []byte("Error marshalling error message"))
		return
	***REMOVED***

	wsutil.WriteServerText(*conn, marshalledMessage)
***REMOVED***

func displaySuccessMessage(conn *net.Conn, message string) ***REMOVED***
	successMessage := libAgent.ServerDisplayableSuccessMessage***REMOVED***
		Type:    "displayableSuccessMessage",
		Message: message,
	***REMOVED***

	marshalledMessage, err := json.Marshal(successMessage)
	if err != nil ***REMOVED***
		wsutil.WriteServerText(*conn, []byte("Error marshalling success message"))
		return
	***REMOVED***

	wsutil.WriteServerText(*conn, marshalledMessage)
***REMOVED***

func notifyJobDeleted(conn *net.Conn, jobId string) ***REMOVED***
	message := libAgent.ServerJobDeletedMessage***REMOVED***
		Type:    "jobDeleted",
		Message: jobId,
	***REMOVED***

	marshalledMessage, err := json.Marshal(message)
	if err != nil ***REMOVED***
		wsutil.WriteServerText(*conn, []byte("Error marshalling job deleted message"))
		return
	***REMOVED***

	wsutil.WriteServerText(*conn, marshalledMessage)
***REMOVED***
