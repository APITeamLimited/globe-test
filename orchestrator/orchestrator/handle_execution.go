package orchestrator

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/orchestrator/orchMetrics"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/redis/v9"
)

type locatedMesaage struct ***REMOVED***
	location string
	msg      *redis.Message
***REMOVED***

type jobUserUpdate struct ***REMOVED***
	UpdateType string `json:"updateType"`
***REMOVED***

func handleExecution(gs libOrch.BaseGlobalState, options *libWorker.Options, scope libOrch.Scope, childJobs map[string]jobDistribution, jobId string) (string, error) ***REMOVED***
	libOrch.UpdateStatus(gs, "LOADING")

	// Create a handler for aborts
	jobUserUpdatesChannel := gs.Client().Subscribe(gs.Ctx(), fmt.Sprintf("jobUserUpdates:%s:%s:%s", scope.Variant, scope.VariantTargetId, jobId)).Channel()

	workerSubscriptions := make(map[string]*redis.PubSub)
	for location, jobDistribution := range childJobs ***REMOVED***
		if jobDistribution.Jobs != nil && len(jobDistribution.Jobs) > 0 ***REMOVED***
			workerSubscriptions[location] = jobDistribution.workerClient.Subscribe(gs.Ctx(), fmt.Sprintf("worker:executionUpdates:%s", jobId))
		***REMOVED***
	***REMOVED***

	workerChannels := make(map[string]<-chan *redis.Message)
	for location, subscription := range workerSubscriptions ***REMOVED***
		workerChannels[location] = subscription.Channel()
	***REMOVED***

	// Update the status
	libOrch.UpdateStatus(gs, "RUNNING")

	// Check if workerSubscriptions is empty
	if len(workerSubscriptions) == 0 ***REMOVED***
		libOrch.DispatchMessage(gs, "No child jobs were created", "INFO")
		return "SUCCESS", nil
	***REMOVED***

	for _, jobDistribution := range childJobs ***REMOVED***
		for _, job := range jobDistribution.Jobs ***REMOVED***
			err := dispatchJob(gs, jobDistribution.workerClient, job, options)
			if err != nil ***REMOVED***
				return "FAILURE", err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	unifiedChannel := make(chan locatedMesaage)

	for location, channel := range workerChannels ***REMOVED***
		// Variable declaration here for locationLoop seems to be required to avoid
		// capturing the loop variable error
		locationLoop := location
		go func(channel <-chan *redis.Message) ***REMOVED***
			for msg := range channel ***REMOVED***
				newMessage := locatedMesaage***REMOVED***
					location: locationLoop,
					msg:      msg,
				***REMOVED***

				unifiedChannel <- newMessage
			***REMOVED***
		***REMOVED***(channel)
	***REMOVED***

	// Add jobUserUpdatesChannel to the unifiedChannel
	go func() ***REMOVED***
		for msg := range jobUserUpdatesChannel ***REMOVED***
			newMessage := locatedMesaage***REMOVED***
				location: "jobUserUpdates",
				msg:      msg,
			***REMOVED***

			unifiedChannel <- newMessage
		***REMOVED***
	***REMOVED***()

	childJobCount := 0
	for _, jobDistribution := range childJobs ***REMOVED***
		childJobCount += len(jobDistribution.Jobs)
	***REMOVED***

	jobsInitialised := []string***REMOVED******REMOVED***
	jobsMutex := &sync.Mutex***REMOVED******REMOVED***

	successCount := 0
	failureCount := 0
	resolutionMutex := sync.Mutex***REMOVED******REMOVED***

	summaryBank := orchMetrics.NewSummaryBank(gs, options)

	for locatedMessage := range unifiedChannel ***REMOVED***
		// Handle user updates separately
		if locatedMessage.location == "jobUserUpdates" ***REMOVED***
			var jobUserUpdate = jobUserUpdate***REMOVED******REMOVED***

			err := json.Unmarshal([]byte(locatedMessage.msg.Payload), &jobUserUpdate)
			if err != nil ***REMOVED***
				return "FAILURE", err
			***REMOVED***

			if jobUserUpdate.UpdateType == "CANCEL" ***REMOVED***
				fmt.Println("Received job abort signal")
				libOrch.DispatchMessage(gs, "Job cancelled by user", "INFO")

				// Cancel all child jobs
				for _, jobDistribution := range childJobs ***REMOVED***
					for _, job := range jobDistribution.Jobs ***REMOVED***
						err := jobDistribution.workerClient.Publish(gs.Ctx(), fmt.Sprintf("childJobUserUpdates:%s", job.ChildJobId), locatedMessage.msg.Payload).Err()
						if err != nil ***REMOVED***
							return "FAILURE", err
						***REMOVED***
					***REMOVED***
				***REMOVED***

				return "FAILURE", errors.New("job cancelled by user")
			***REMOVED***

			// jobUserUpdates are different to other messages, so we can skip the rest of the loop
			continue
		***REMOVED***

		var workerMessage = libOrch.WorkerMessage***REMOVED******REMOVED***
		err := json.Unmarshal([]byte(locatedMessage.msg.Payload), &workerMessage)
		if err != nil ***REMOVED***
			return "FAILURE", err
		***REMOVED***

		if workerMessage.MessageType == "STATUS" ***REMOVED***
			gs.SetChildJobState(workerMessage.WorkerId, workerMessage.ChildJobId, workerMessage.Message)

			if workerMessage.Message == "READY" ***REMOVED***
				jobsMutex.Lock()

				alreadyInitialised := false

				// Check if the job has already been initialised
				for _, initialisedJob := range jobsInitialised ***REMOVED***
					if initialisedJob == workerMessage.ChildJobId ***REMOVED***
						jobsMutex.Unlock()
						alreadyInitialised = true
						break
					***REMOVED***
				***REMOVED***

				if !alreadyInitialised ***REMOVED***
					jobsInitialised = append(jobsInitialised, workerMessage.ChildJobId)
					jobsMutex.Unlock()

					if len(jobsInitialised) == childJobCount ***REMOVED***
						// Broadcast the start message to all child jobs
						for _, jobDistribution := range childJobs ***REMOVED***
							for _, job := range jobDistribution.Jobs ***REMOVED***
								jobDistribution.workerClient.Publish(gs.Ctx(), fmt.Sprintf("%s:go", job.ChildJobId), "GO TIME")
							***REMOVED***
						***REMOVED***

						libOrch.UpdateStatus(gs, "RUNNING")
					***REMOVED***
				***REMOVED***
			***REMOVED*** else if workerMessage.Message == "SUCCESS" || workerMessage.Message == "FAILURE" ***REMOVED***
				resolutionMutex.Lock()
				if workerMessage.Message == "SUCCESS" ***REMOVED***
					successCount++
				***REMOVED*** else ***REMOVED***
					failureCount++
				***REMOVED***
				resolutionMutex.Unlock()

				if successCount+failureCount == childJobCount ***REMOVED***
					// All jobs have finished
					if failureCount > 0 ***REMOVED***
						libOrch.UpdateStatus(gs, "FAILURE")
						return "FAILURE", nil
					***REMOVED*** else ***REMOVED***
						libOrch.UpdateStatus(gs, "SUCCESS")
						return "SUCCESS", nil
					***REMOVED***
				***REMOVED***
			***REMOVED***
			// Ignore other kinds of status messages
			//libOrch.DispatchWorkerMessage(gs.Ctx(), gs.Client(), gs.JobId(), workerMessage.WorkerId, workerMessage.Message, "STATUS")

			// Sometimes errors don't stop the execution automatically so stop them here
		***REMOVED*** else if workerMessage.MessageType == "ERROR" ***REMOVED***
			libOrch.DispatchWorkerMessage(gs, workerMessage.WorkerId, workerMessage.ChildJobId, workerMessage.Message, workerMessage.MessageType)

			resolutionMutex.Lock()
			failureCount++
			resolutionMutex.Unlock()

			if successCount+failureCount == childJobCount ***REMOVED***
				// All jobs have finished
				libOrch.UpdateStatus(gs, "FAILURE")
				return "FAILURE", nil
			***REMOVED***
		***REMOVED*** else if workerMessage.MessageType == "METRICS" ***REMOVED***
			childJob := findChildJob(&childJobs, locatedMessage.location, workerMessage.ChildJobId)
			if childJob == nil ***REMOVED***
				return "FAILURE", fmt.Errorf("could not find child job with id %s to add summary metrics to", workerMessage.ChildJobId)
			***REMOVED***

			(*gs.MetricsStore()).AddMessage(workerMessage, locatedMessage.location, childJob.SubFraction)
		***REMOVED*** else if workerMessage.MessageType == "SUMMARY_METRICS" ***REMOVED***
			childJob := findChildJob(&childJobs, locatedMessage.location, workerMessage.ChildJobId)
			if childJob == nil ***REMOVED***
				return "FAILURE", fmt.Errorf("could not find child job with id %s to add summary metrics to", workerMessage.ChildJobId)
			***REMOVED***

			summaryBank.AddMessage(workerMessage, locatedMessage.location, childJob.SubFraction)

			if summaryBank.Size() == childJobCount ***REMOVED***
				err := summaryBank.CalculateAndDispatchSummaryMetrics()
				if err != nil ***REMOVED***
					return "FAILURE", err
				***REMOVED***
			***REMOVED***
		***REMOVED*** else if workerMessage.MessageType == "DEBUG" ***REMOVED***
			// TODO: make this configurable
			continue
		***REMOVED*** else ***REMOVED***
			libOrch.DispatchWorkerMessage(gs, workerMessage.WorkerId, workerMessage.ChildJobId, workerMessage.Message, workerMessage.MessageType)
		***REMOVED***
	***REMOVED***

	// Should never get here
	return "FAILURE", errors.New("an unexpected error occurred")
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

func findChildJob(childJobs *map[string]jobDistribution, location string, childJobId string) *libOrch.ChildJob ***REMOVED***
	for _, childJob := range (*childJobs)[location].Jobs ***REMOVED***
		if childJob.ChildJobId == childJobId ***REMOVED***
			return &childJob
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
