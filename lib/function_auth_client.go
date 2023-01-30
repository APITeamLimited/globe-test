package lib

import (
	"net/http"

	"cloud.google.com/go/functions/apiv2/functionspb"
)

type LiveService struct {
	Location string
	Uri      string
	State    functionspb.Function_State
}

type FunctionResult struct {
	Response *http.Response
	Error    error
}

type RunAuthClient interface {
	Functions() []LiveService
	ExecuteService(location string) (*(chan FunctionResult), error)
	CheckServiceAvailability(location string) error
}
