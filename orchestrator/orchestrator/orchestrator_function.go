package orchestrator

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"sync"

// 	"github.com/APITeamLimited/globe-test/lib"
// 	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
// 	"github.com/google/uuid"
// 	"github.com/gorilla/websocket"
// )

// func RunOrchestratorServer() {
// 	port := lib.GetEnvVariableRaw("ORCHESTRATOR_SERVER_PORT", "8079", true)
// 	standalone := lib.GetEnvVariableBool("ORCHESTRATOR_STANDALONE", true)

// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		runOrchestrator(w, r, standalone)
// 	})

// 	fmt.Printf("Starting GlobeTest orchestrator server on port %s\n", port)

// 	if err := http.ListenAndServe(":"+port, nil); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func runOrchestrator(w http.ResponseWriter, r *http.Request, standalone bool) {
// 	ctx := context.Background()
// 	orchestratorId := uuid.NewString()

// 	upgrader := websocket.Upgrader{}

// 	// Upgrade the connection to a websocket
// 	conn, err := upgrader.Upgrade(w, r, headers)
// 	if err != nil {
// 		fmt.Printf("Error upgrading connection to websocket: %s", err.Error())
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte(err.Error()))
// 		return
// 	}

// 	connWriteMutex := &sync.Mutex{}
// 	connReadMutex := &sync.Mutex{}

// 	job, err := extractJobMessage(conn, connReadMutex)
// 	if err != nil {
// 		fmt.Printf("Error extracting event message: %s", err.Error())
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte(err.Error()))
// 		return
// 	}

// 	job.AssignedOrchestrator = orchestratorId

// }

// func extractJobMessage(conn *websocket.Conn, connReadMutex *sync.Mutex) (*libOrch.Job, error) {
// 	eventMessage := lib.EventMessage{}

// 	// Read the message from the connection
// 	connReadMutex.Lock()
// 	// Okay here as at the start
// 	err := conn.ReadJSON(&eventMessage)
// 	connReadMutex.Unlock()
// 	if err != nil {
// 		// If websocket is closed, return
// 		fmt.Printf("Error reading message from connection: %s", err.Error())
// 		return nil, err
// 	}

// 	jobMessage := &libOrch.Job{}

// 	// Unmarshal the message
// 	err = json.Unmarshal(byte(eventMessage.Data), jobMessage)
// 	if err != nil {
// 		fmt.Printf("Error unmarshalling message: %s", err.Error())
// 		return nil, err
// 	}

// 	return jobMessage, nil
// }
