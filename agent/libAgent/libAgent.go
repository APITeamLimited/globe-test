package libAgent

import (
	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
)

const AgentPort = 59125

const OrchestratorRedisHost = "localhost"
const OrchestratorRedisPort = "59126"

const WorkerRedisHost = "localhost"
const WorkerRedisPort = "59127"

type (
	ClientLocalTestManagerMessage struct ***REMOVED***
		Type string `json:"type"`
	***REMOVED***

	ClientNewJobMessage struct ***REMOVED***
		Type    string      `json:"type"` // "newJob"
		Message libOrch.Job `json:"message"`
	***REMOVED***
	ClientAbortJobMessage struct ***REMOVED***
		Type    string `json:"type"` // "abortJob"
		Message string `json:"message"`
	***REMOVED***
	ClientJobUpdateMessage struct ***REMOVED***
		Type    string                   `json:"type"` // "jobUpdate"
		Message lib.WrappedJobUserUpdate `json:"message"`
	***REMOVED***
)

// Server relays some messages back when successful

type (
	ServerLocalTestManagerMessage struct ***REMOVED***
		Type string `json:"type"`
	***REMOVED***

	ServerAbortJobMessage struct ***REMOVED***
		Type    string `json:"type"` // "abortJob"
		Message string `json:"message"`
	***REMOVED***

	ServerNewJobMessage struct ***REMOVED***
		Type    string      `json:"type"` // "newJob"
		Message libOrch.Job `json:"message"`
	***REMOVED***

	ServerGlobeTestMessage struct ***REMOVED***
		Type    string `json:"type"` // "globeTestMessage"
		Message string `json:"message"`
	***REMOVED***

	ServerWrappedJobUserUpdate struct ***REMOVED***
		Type    string                   `json:"type"` // "jobUpdate"
		Message lib.WrappedJobUserUpdate `json:"message"`
	***REMOVED***
)
