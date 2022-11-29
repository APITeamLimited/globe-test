package agent

import (
	"net"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
)

func handleAbortAllJobs(runningJobs *map[string]libOrch.Job, setJobCount func(int),
	orchestratorClient *redis.Client, connections *map[string]*net.Conn) {
	// Loop through all running jobs and cancel them
	for _, job := range *runningJobs {
		processAbortion(job, runningJobs, setJobCount, orchestratorClient, connections)
	}
	setJobCount(len(*runningJobs))
}
