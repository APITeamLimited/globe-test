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
	abortAllChannel chan struct{},
	setJobCount func(int),
) {
	orchestratorClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", libAgent.OrchestratorRedisHost, libAgent.OrchestratorRedisPort),
		Username: "default",
		Password: "",
	})

	runningJobs := make(map[string]libOrch.Job)
	connections := make(map[string]*net.Conn)

	serverAddress := fmt.Sprintf("localhost:%d", libAgent.AgentPort)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			fmt.Println("error upgrading to websocket:", err)
			return
		}

		randId := uuid.New().String()
		fmt.Printf("New connection %s\n", randId)

		connections[randId] = &conn

		sendRunningJobsToClient(&conn, &runningJobs)

		go func() {
			defer conn.Close()
			defer delete(connections, randId)

			for {
				msg, op, err := wsutil.ReadClientData(conn)
				if err != nil {
					fmt.Println("read error:", err)
					return
				}

				// Return ping messages
				if op == ws.OpPing {
					wsutil.WriteServerMessage(conn, ws.OpPong, msg)
				}

				var parsedMessage libAgent.ClientLocalTestManagerMessage
				err = json.Unmarshal(msg, &parsedMessage)

				if err != nil {
					fmt.Println("error parsing message:", err)
					return
				}

				switch parsedMessage.Type {
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
				}
			}
		}()
	})

	fmt.Printf("Starting agent server on %s\n", serverAddress)
	http.Handle("/agent", cors.AllowAll().Handler(mux))
	http.ListenAndServe(serverAddress, nil)
}
