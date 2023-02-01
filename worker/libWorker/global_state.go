package libWorker

import (
	"context"
	"sync"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/gorilla/websocket"
)

type WorkerState struct {
	Status   string `json:"status"`
	WorkerId string `json:"workerId"`
}

type BaseGlobalState interface {
	Ctx() context.Context
	Conn() *websocket.Conn
	ConnWriteMutex() *sync.Mutex
	ConnReadMutex() *sync.Mutex
	JobId() string
	ChildJobId() string
	WorkerId() string
	GetWorkerStatus() string
	SetWorkerStatus(status string)
	FuncModeEnabled() bool
	FuncModeInfo() *lib.FuncModeInfo
	MessageQueue() *MessageQueue

	SetRunAbortFunc(cancelFunc context.CancelFunc)
	GetRunAbortFunc() context.CancelFunc
}
