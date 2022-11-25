package agent

import (
	"context"
	"fmt"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

type jobArgs struct ***REMOVED***
	JobId      string `json:"jobId"`
	Source     string `json:"source"`
	SourceName string `json:"sourceName"`
***REMOVED***

type runningJob struct ***REMOVED***
	JobId      string
	Source     string
	SourceName string
	Messages   []string
***REMOVED***

func runAgentServer(
	quitChannel chan struct***REMOVED******REMOVED***,
	abortAllChannel chan struct***REMOVED******REMOVED***,
	setJobCount func(int),
) chan struct***REMOVED******REMOVED*** ***REMOVED***
	server := socketio.NewServer(nil)

	runningJobs := make(map[string]runningJob)

	server.OnConnect("/", func(s socketio.Conn) error ***REMOVED***
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	***REMOVED***)

	server.OnEvent("/", "newJob", func(s socketio.Conn, msg string) ***REMOVED***
		fmt.Println("newJob:", msg)

		handleJobCreation(msg, s, runningJobs, setJobCount)
	***REMOVED***)

	server.OnEvent("/", "abortJob", func(s socketio.Conn, msg string) ***REMOVED***
		fmt.Println("abortJob:", msg)
		abortJob(msg, runningJobs, setJobCount)
	***REMOVED***)

	server.OnEvent("/", "abortAllJobs", func(s socketio.Conn, msg string) ***REMOVED***
		fmt.Println("abortAllJobs:", msg)
		abortAllJobs(runningJobs, setJobCount)
	***REMOVED***)

	server.OnDisconnect("/", func(s socketio.Conn, reason string) ***REMOVED***
		fmt.Println("closed", reason)
	***REMOVED***)

	httpServer := http.Server***REMOVED***
		Addr:    "localhost:5000",
		Handler: server,
	***REMOVED***

	serverStoppedCh := make(chan struct***REMOVED******REMOVED***)

	go func() ***REMOVED***
		<-quitChannel
		fmt.Println("Shutting down agent server")
		httpServer.Shutdown(context.Background())
		server.Close()
		serverStoppedCh <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***()

	go func() ***REMOVED***
		<-abortAllChannel
		abortAllJobs(runningJobs, setJobCount)
	***REMOVED***()

	return serverStoppedCh
***REMOVED***
