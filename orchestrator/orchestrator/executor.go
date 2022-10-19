package orchestrator

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/redis/v9"
)

type locatedMesaage struct ***REMOVED***
	location string
	msg      *redis.Message
***REMOVED***

func runExecution(gs libOrch.BaseGlobalState, options *libWorker.Options, scope libOrch.Scope, childJobs map[string]jobDistribution, jobId string) (string, error) ***REMOVED***
	libOrch.UpdateStatus(gs.Ctx(), gs.Client(), jobId, gs.OrchestratorId(), "LOADING")

	workerSubscriptions := make(map[string]*redis.PubSub)
	for location, jobDistribution := range childJobs ***REMOVED***
		if jobDistribution.jobs != nil && len(jobDistribution.jobs) > 0 ***REMOVED***
			workerSubscriptions[location] = jobDistribution.workerClient.Subscribe(gs.Ctx(), fmt.Sprintf("worker:executionUpdates:%s", jobId))
		***REMOVED***
	***REMOVED***

	workerChannels := make(map[string]<-chan *redis.Message)
	for location, subscription := range workerSubscriptions ***REMOVED***
		workerChannels[location] = subscription.Channel()
	***REMOVED***

	// Update the status
	libOrch.UpdateStatus(gs.Ctx(), gs.Client(), jobId, gs.OrchestratorId(), "RUNNING")

	// Check if workerSubscriptions is empty
	if len(workerSubscriptions) == 0 ***REMOVED***
		libOrch.DispatchMessage(gs.Ctx(), gs.Client(), gs.JobId(), gs.OrchestratorId(), "No child jobs were created", "INFO")
		return "SUCCESS", nil
	***REMOVED***

	for _, jobDistribution := range childJobs ***REMOVED***
		for _, job := range jobDistribution.jobs ***REMOVED***
			err := dispatchJob(gs, jobDistribution.workerClient, job, options)
			if err != nil ***REMOVED***
				return "", err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	unifiedChannel := make(chan locatedMesaage)

	for location, channel := range workerChannels ***REMOVED***
		go func(channel <-chan *redis.Message) ***REMOVED***
			for msg := range channel ***REMOVED***
				unifiedChannel <- locatedMesaage***REMOVED***
					location,
					msg,
				***REMOVED***
			***REMOVED***
		***REMOVED***(channel)
	***REMOVED***

	chilJobCount := len(childJobs)
	jobsInitialised := 0

	for locatedMessage := range unifiedChannel ***REMOVED***
		workerMessage := libOrch.WorkerMessage***REMOVED******REMOVED***
		err := json.Unmarshal([]byte(locatedMessage.msg.Payload), &workerMessage)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***

		if workerMessage.MessageType == "STATUS" ***REMOVED***
			if workerMessage.Message == "READY" ***REMOVED***
				jobsInitialised++
				if jobsInitialised == chilJobCount ***REMOVED***
					// Broadcast the start message to all child jobs

					for _, jobDistribution := range childJobs ***REMOVED***
						for _, job := range jobDistribution.jobs ***REMOVED***
							jobDistribution.workerClient.Publish(gs.Ctx(), fmt.Sprintf("%s:go", job.ChildJobId), "GO TIME")
						***REMOVED***
					***REMOVED***

				***REMOVED***

				libOrch.UpdateStatus(gs.Ctx(), gs.Client(), jobId, gs.OrchestratorId(), "RUNNING")
			***REMOVED*** else if workerMessage.Message == "FAILURE" ***REMOVED***
				return "FAILURE", nil
			***REMOVED*** else if workerMessage.Message == "SUCCESS" ***REMOVED***
				return "SUCCESS", nil
			***REMOVED***
			// Ignore other kinds of status messages
			//libOrch.DispatchWorkerMessage(gs.Ctx(), gs.Client(), gs.JobId(), workerMessage.WorkerId, workerMessage.Message, "STATUS")

			// Sometimes errors don't stop the execution automatically so stop them here
		***REMOVED*** else if workerMessage.MessageType == "ERROR" ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx(), gs.Client(), gs.JobId(), workerMessage.WorkerId, workerMessage.Message, workerMessage.MessageType)
			return "FAILURE", nil
		***REMOVED*** else if workerMessage.MessageType == "METRICS" ***REMOVED***
			(*gs.MetricsStore()).AddMessage(workerMessage, locatedMessage.location)
		***REMOVED*** else if workerMessage.MessageType == "DEBUG" ***REMOVED***
			// TODO: make this configurable
			continue
		***REMOVED*** else ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx(), gs.Client(), gs.JobId(), workerMessage.WorkerId, workerMessage.Message, workerMessage.MessageType)
		***REMOVED***

		// Could handle these differently, but for now just dispatch them

		/*else if workerMessage.MessageType == "MARK" ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job.Id, workerMessage.WorkerId, workerMessage.Message, "MARK")
		***REMOVED*** else if workerMessage.MessageType == "CONSOLE" ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job.Id, workerMessage.WorkerId, workerMessage.Message, "CONSOLE")
		***REMOVED*** else if workerMessage.MessageType == "METRICS" ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job.Id, workerMessage.WorkerId, workerMessage.Message, "METRICS")
		***REMOVED*** else if workerMessage.MessageType == "SUMMARY_METRICS" ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job.Id, workerMessage.WorkerId, workerMessage.Message, "SUMMARY_METRICS")
			workerSubscription.Close()
		***REMOVED*** else if workerMessage.MessageType == "ERROR" ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job.Id, workerMessage.WorkerId, workerMessage.Message, "ERROR")
			libOrch.HandleStringError(gs.Ctx, orchestratorClient, job.Id, gs.orchestratorId, workerMessage.Message)
			return
		***REMOVED*** else if workerMessage.MessageType == "DEBUG" ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job.Id, workerMessage.WorkerId, workerMessage.Message, "DEBUG")
		***REMOVED****/
	***REMOVED***

	// Shouldn't get here
	return "", errors.New("an unexpected error occurred")
***REMOVED***

func dispatchJob(gs libOrch.BaseGlobalState, workerClient *redis.Client, job libOrch.ChildJob, options *libWorker.Options) error ***REMOVED***
	// Convert options to json
	marshalledChildJob, err := json.Marshal(job)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	workerClient.HSet(gs.Ctx(), job.ChildJobId, "job", marshalledChildJob)

	workerClient.SAdd(gs.Ctx(), "worker:executionHistory", job.ChildJobId)
	workerClient.Publish(gs.Ctx(), "worker:execution", job.ChildJobId)

	return nil
***REMOVED***
