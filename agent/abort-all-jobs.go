package agent

import (
	"net"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
)

func handleAbortAllJobs(runningJobs *map[string]libOrch.Job, conn *net.Conn, setJobCount func(int),
	orchestratorClient *redis.Client) ***REMOVED***
	// Loop through all running jobs and cancel them
	for _, job := range *runningJobs ***REMOVED***
		processAbortion(job, runningJobs, setJobCount, orchestratorClient)
	***REMOVED***
	setJobCount(len(*runningJobs))

	displaySuccessMessage(conn, "Stopping all test runs")
***REMOVED***
