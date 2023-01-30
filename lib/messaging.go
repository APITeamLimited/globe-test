package lib

type JobUserUpdate struct {
	UpdateType string `json:"updateType"`
}

const (
	CHILD_JOB_INFO    = "childJobInfo"
	GO_MESSAGE_TYPE   = "go"
	CHILD_USER_UPDATE = "childUserUpdate"
)

type EventMessage struct {
	Variant string `json:"variant"`
	// JSON encoded string
	Data string `json:"data"`
}
