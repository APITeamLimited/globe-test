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

	// server := socketio.NewServer(nil)

	// server.OnConnect("/", func(s socketio.Conn) error {
	// 	s.SetContext("")
	// 	fmt.Println("connected:", s.ID())
	// 	// Accept the connection

	// 	s.Emit("message", "Hello, world!")

	// 	return nil
	// })

	// server.OnEvent("/", "newJob", func(s socketio.Conn, msg string) {
	// 	fmt.Println("newJob:", msg)

	// 	handleNewJob(msg, conn, runningJobs, setJobCount)
	// })

	// server.OnEvent("/", "abortJob", func(s socketio.Conn, msg string) {
	// 	fmt.Println("abortJob:", msg)
	// 	abortJob(msg, runningJobs, setJobCount)
	// })

	// server.OnEvent("/", "abortAllJobs", func(s socketio.Conn, msg string) {
	// 	fmt.Println("abortAllJobs:", msg)
	// 	abortAllJobs(runningJobs, setJobCount)
	// })

	// server.OnEvent("/", "jobUpdate", func(s socketio.Conn, msg string) {
	// 	fmt.Println("jobUpdate:", msg)
	// 	// TODO
	// })

	// server.OnDisconnect("/", func(s socketio.Conn, reason string) {
	// 	fmt.Println("closed", reason)
	// })

	// go func() {
	// 	if err := server.Serve(); err != nil {
	// 		log.Fatalf("socketio listen error: %s\n", err)
	// 	}
	// }()
	// defer server.Close()

	serverAddress := fmt.Sprintf("localhost:%d", libAgent.AgentPort)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			fmt.Println("error upgrading to websocket:", err)
			return
		}

		randId := uuid.New().String()

		connections[randId] = &conn

		go func() {
			defer conn.Close()

			for {
				msg, op, err := wsutil.ReadClientData(conn)
				if err != nil {
					fmt.Println("read error:", err)
					return
				}

				// Handle closed connections
				if op == ws.OpClose {
					delete(connections, randId)
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
					handleNewJob(msg, &conn, &runningJobs, setJobCount, orchestratorClient, &connections)
				case "abortJob":
					handleAbortJob(msg, &conn, &runningJobs, setJobCount, orchestratorClient, &connections)
				case "abortAllJobs":
					handleAbortAllJobs(&runningJobs, setJobCount, orchestratorClient, &connections)
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
