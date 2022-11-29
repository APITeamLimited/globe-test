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
	setJobCount func(int), orchestratorClient *redis.Client, connections *map[string]*net.Conn) {
	parsedMessage := libAgent.ClientNewJobMessage{}

	err := json.Unmarshal(rawMessage, &parsedMessage)
	if err != nil {
		wsutil.WriteServerText(*conn, []byte("Error parsing job arguments"))
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
	broadcastMessage(marshalledServerNewJob, connections)

	go streamGlobeTestMessages(parsedMessage, orchestratorClient, connections, runningJobs, setJobCount)

}

func streamGlobeTestMessages(parsedMessage libAgent.ClientNewJobMessage, orchestratorClient *redis.Client,
	connections *map[string]*net.Conn, runningJobs *map[string]libOrch.Job, setJobCount func(int)) {
	subscriptionKey := fmt.Sprintf("jobUserUpdates:%s:%s:%s", parsedMessage.Message.Scope.Variant, parsedMessage.Message.Scope.VariantTargetId, parsedMessage.Message.Id)
	subscription := orchestratorClient.Subscribe(context.Background(), subscriptionKey)
	defer subscription.Close()

	for msg := range subscription.Channel() {
		parsedMessage := libOrch.OrchestratorOrWorkerMessage{}
		err := json.Unmarshal([]byte(msg.Payload), &parsedMessage)
		if err != nil {
			fmt.Println("Error parsing job update message")
			continue
		}

		if parsedMessage.MessageType == "STATUS" {
			if parsedMessage.Message == "COMPLETED_SUCCESS" || parsedMessage.Message == "COMPLETED_FAILURE" {
				fmt.Printf("Job %s completed with status %s", parsedMessage.JobId, parsedMessage.Message)

				// Delete the job from the running jobs
				delete(*runningJobs, parsedMessage.JobId)
				setJobCount(len(*runningJobs))
				return
			}
		}

		broadcastMessage([]byte(msg.Payload), connections)
	}
}
