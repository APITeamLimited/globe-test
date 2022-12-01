package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/APITeamLimited/globe-test/agent/libAgent"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"github.com/gobwas/ws/wsutil"
)

func handleNewJob(rawMessage []byte, conn *net.Conn, runningJobs *map[string]libOrch.Job,
	setJobCount func(int), orchestratorClient *redis.Client) {
	if len(*runningJobs) >= 5 {
		displayErrorMessage(conn, "You can run a maximum of 5 localhost jobs at once")
		return
	}

	createNewJob(rawMessage, conn, runningJobs, setJobCount, orchestratorClient)
}

func createNewJob(rawMessage []byte, conn *net.Conn, runningJobs *map[string]libOrch.Job,
	setJobCount func(int), orchestratorClient *redis.Client) {
	parsedMessage := libAgent.ClientNewJobMessage{}

	err := json.Unmarshal(rawMessage, &parsedMessage)
	if err != nil {
		displayErrorMessage(conn, "Error parsing job arguments")
		return
	}

	// Create a new job
	(*runningJobs)[parsedMessage.Message.Id] = parsedMessage.Message

	// Send the job to the orchestrator
	marshalledJob, _ := json.Marshal(parsedMessage.Message)

	// Set job info first to prevent race condition
	orchestratorClient.HSet(context.Background(), parsedMessage.Message.Id, "job", string(marshalledJob))
	orchestratorClient.SAdd(context.Background(), "orchestrator:executionHistory", parsedMessage.Message.Id)
	orchestratorClient.Publish(context.Background(), "orchestrator:execution", string(parsedMessage.Message.Id))

	serverNewJobMessage := libAgent.ServerNewJobMessage{
		Type:    "newJob",
		Message: parsedMessage.Message,
	}

	marshalledServerNewJob, _ := json.Marshal(serverNewJobMessage)
	wsutil.WriteServerText(*conn, marshalledServerNewJob)

	go streamGlobeTestMessages(parsedMessage, orchestratorClient, conn, runningJobs, setJobCount)
}

func streamGlobeTestMessages(parsedMessage libAgent.ClientNewJobMessage, orchestratorClient *redis.Client,
	conn *net.Conn, runningJobs *map[string]libOrch.Job, setJobCount func(int)) {
	subscriptionKey := fmt.Sprintf("orchestrator:executionUpdates:%s", parsedMessage.Message.Id)
	subscription := orchestratorClient.Subscribe(context.Background(), subscriptionKey)

	defer subscription.Close()

	for msg := range subscription.Channel() {
		parsedMessage := libOrch.OrchestratorOrWorkerMessage{}
		err := json.Unmarshal([]byte(msg.Payload), &parsedMessage)
		if err != nil {
			displayErrorMessage(conn, "Error parsing job update message")
			continue
		}

		if parsedMessage.MessageType == "STATUS" {
			if parsedMessage.Message == "COMPLETED_SUCCESS" || parsedMessage.Message == "COMPLETED_FAILURE" {
				fmt.Printf("Agent completed job %s with status %s\n", parsedMessage.JobId, parsedMessage.Message)

				writeGlobeTestMessage(conn, msg)

				// Delete the job from the running jobs
				delete(*runningJobs, parsedMessage.JobId)
				notifyJobDeleted(conn, parsedMessage.JobId)
				setJobCount(len(*runningJobs))

				return
			}
		}

		writeGlobeTestMessage(conn, msg)
	}
}

func writeGlobeTestMessage(conn *net.Conn, msg *redis.Message) {
	serverGlobeTestMessage := libAgent.ServerGlobeTestMessage{
		Type:    "globeTestMessage",
		Message: msg.Payload,
	}

	marshalledServerGlobeTestMessage, _ := json.Marshal(serverGlobeTestMessage)
	wsutil.WriteServerText(*conn, marshalledServerGlobeTestMessage)
}
