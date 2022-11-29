package lib

type JobUserUpdate struct {
	UpdateType string `json:"updateType"`
}

type WrappedJobUserUpdate struct {
	Update JobUserUpdate `json:"update"`
	JobId  string        `json:"jobId"`
}
