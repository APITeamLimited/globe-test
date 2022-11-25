package agent

import (
	"context"
	"fmt"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

type jobArgs struct {
	JobId      string `json:"jobId"`
	Source     string `json:"source"`
	SourceName string `json:"sourceName"`
}

type runningJob struct {
	JobId      string
	Source     string
	SourceName string
	Messages   []string
}

func runAgentServer(
	quitChannel chan struct{},
	abortAllChannel chan struct{},
	setJobCount func(int),
) chan struct{} {
	server := socketio.NewServer(nil)

	runningJobs := make(map[string]runningJob)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "newJob", func(s socketio.Conn, msg string) {
		fmt.Println("newJob:", msg)

		handleJobCreation(msg, s, runningJobs, setJobCount)
	})

	server.OnEvent("/", "abortJob", func(s socketio.Conn, msg string) {
		fmt.Println("abortJob:", msg)
		abortJob(msg, runningJobs, setJobCount)
	})

	server.OnEvent("/", "abortAllJobs", func(s socketio.Conn, msg string) {
		fmt.Println("abortAllJobs:", msg)
		abortAllJobs(runningJobs, setJobCount)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	httpServer := http.Server{
		Addr:    "localhost:5000",
		Handler: server,
	}

	serverStoppedCh := make(chan struct{})

	go func() {
		<-quitChannel
		fmt.Println("Shutting down agent server")
		httpServer.Shutdown(context.Background())
		server.Close()
		serverStoppedCh <- struct{}{}
	}()

	go func() {
		<-abortAllChannel
		abortAllJobs(runningJobs, setJobCount)
	}()

	return serverStoppedCh
}
