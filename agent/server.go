package agent

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/APITeamLimited/globe-test/agent/libAgent"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
	"github.com/rs/cors"
)

func runAgentServer(
	abortAllChannel chan struct***REMOVED******REMOVED***,
	setJobCount func(int),
) ***REMOVED***
	orchestratorClient := redis.NewClient(&redis.Options***REMOVED***
		Addr:     fmt.Sprintf("%s:%s", libAgent.OrchestratorRedisHost, libAgent.OrchestratorRedisPort),
		Username: "default",
		Password: "",
	***REMOVED***)

	runningJobs := make(map[string]libOrch.Job)
	connections := make(map[string]*net.Conn)

	serverAddress := fmt.Sprintf("localhost:%d", libAgent.AgentPort)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil ***REMOVED***
			fmt.Println("error upgrading to websocket:", err)
			return
		***REMOVED***

		randId := uuid.New().String()
		fmt.Printf("New connection assigned ID %s\n", randId)

		connections[randId] = &conn

		sendRunningJobsToClient(&conn, &runningJobs)

		go func() ***REMOVED***
			defer conn.Close()
			defer delete(connections, randId)

			for ***REMOVED***
				msg, op, err := wsutil.ReadClientData(conn)
				if err != nil ***REMOVED***
					fmt.Println("read error:", err)
					return
				***REMOVED***

				// Return ping messages
				if op == ws.OpPing ***REMOVED***
					wsutil.WriteServerMessage(conn, ws.OpPong, msg)
				***REMOVED***

				var parsedMessage libAgent.ClientLocalTestManagerMessage
				err = json.Unmarshal(msg, &parsedMessage)

				if err != nil ***REMOVED***
					fmt.Println("error parsing message:", err)
					return
				***REMOVED***

				switch parsedMessage.Type ***REMOVED***
				case "newJob":
					handleNewJob(msg, &conn, &runningJobs, setJobCount, orchestratorClient)
				case "abortJob":
					handleAbortJob(msg, &conn, &runningJobs, setJobCount, orchestratorClient)
				case "abortAllJobs":
					handleAbortAllJobs(&runningJobs, &conn, setJobCount, orchestratorClient)
				case "jobUpdate":
					handleJobUpdate(msg, &conn, &runningJobs, orchestratorClient)
				default:
					fmt.Println("unknown message type:", parsedMessage.Type)
				***REMOVED***
			***REMOVED***
		***REMOVED***()
	***REMOVED***)

	fmt.Printf("Starting agent server on %s\n", serverAddress)
	http.Handle("/agent", cors.AllowAll().Handler(mux))
	http.ListenAndServe(serverAddress, nil)
***REMOVED***
