package libOrch

import "time"

type Distribution struct {
	LoadZone string `json:"loadZone"`
	Percent  int    `json:"percent"`
}

type APITeamOptions struct {
	Distribution `json:"distribution"`
}

type Job struct {
	Id         string `json:"id"`
	Source     string `json:"source"`
	SourceName string `json:"sourceName"`
	Options    string `json:"options"`
}

type OrchestratorMessage struct {
	JobId          string    `json:"jobId"`
	Time           time.Time `json:"time"`
	OrchestratorId string    `json:"orchestratorId"`
	Message        string    `json:"message"`
	MessageType    string    `json:"messageType"`
}

type WorkerMessage struct {
	JobId       string    `json:"jobId"`
	Time        time.Time `json:"time"`
	WorkerId    string    `json:"workerId"`
	Message     string    `json:"message"`
	MessageType string    `json:"messageType"`
}

type OrchestratorOrWorkerMessage struct {
	JobId          string    `json:"jobId"`
	Time           time.Time `json:"time"`
	OrchestratorId string    `json:"orchestratorId"`
	WorkerId       string    `json:"workerId"`
	Message        string    `json:"message"`
	MessageType    string    `json:"messageType"`
}

type MarkMessage struct {
	Mark    string      `json:"mark"`
	Message interface{} `json:"message"`
}
