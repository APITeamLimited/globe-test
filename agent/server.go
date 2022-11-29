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

	// server := socketio.NewServer(nil)

	// server.OnConnect("/", func(s socketio.Conn) error ***REMOVED***
	// 	s.SetContext("")
	// 	fmt.Println("connected:", s.ID())
	// 	// Accept the connection

	// 	s.Emit("message", "Hello, world!")

	// 	return nil
	// ***REMOVED***)

	// server.OnEvent("/", "newJob", func(s socketio.Conn, msg string) ***REMOVED***
	// 	fmt.Println("newJob:", msg)

	// 	handleNewJob(msg, conn, runningJobs, setJobCount)
	// ***REMOVED***)

	// server.OnEvent("/", "abortJob", func(s socketio.Conn, msg string) ***REMOVED***
	// 	fmt.Println("abortJob:", msg)
	// 	abortJob(msg, runningJobs, setJobCount)
	// ***REMOVED***)

	// server.OnEvent("/", "abortAllJobs", func(s socketio.Conn, msg string) ***REMOVED***
	// 	fmt.Println("abortAllJobs:", msg)
	// 	abortAllJobs(runningJobs, setJobCount)
	// ***REMOVED***)

	// server.OnEvent("/", "jobUpdate", func(s socketio.Conn, msg string) ***REMOVED***
	// 	fmt.Println("jobUpdate:", msg)
	// 	// TODO
	// ***REMOVED***)

	// server.OnDisconnect("/", func(s socketio.Conn, reason string) ***REMOVED***
	// 	fmt.Println("closed", reason)
	// ***REMOVED***)

	// go func() ***REMOVED***
	// 	if err := server.Serve(); err != nil ***REMOVED***
	// 		log.Fatalf("socketio listen error: %s\n", err)
	// 	***REMOVED***
	// ***REMOVED***()
	// defer server.Close()

	serverAddress := fmt.Sprintf("localhost:%d", libAgent.AgentPort)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil ***REMOVED***
			fmt.Println("error upgrading to websocket:", err)
			return
		***REMOVED***

		randId := uuid.New().String()

		connections[randId] = &conn

		go func() ***REMOVED***
			defer conn.Close()

			for ***REMOVED***
				msg, op, err := wsutil.ReadClientData(conn)
				if err != nil ***REMOVED***
					fmt.Println("read error:", err)
					return
				***REMOVED***

				// Handle closed connections
				if op == ws.OpClose ***REMOVED***
					delete(connections, randId)
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
					handleNewJob(msg, &conn, &runningJobs, setJobCount, orchestratorClient, &connections)
				case "abortJob":
					handleAbortJob(msg, &conn, &runningJobs, setJobCount, orchestratorClient, &connections)
				case "abortAllJobs":
					handleAbortAllJobs(&runningJobs, setJobCount, orchestratorClient, &connections)
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
