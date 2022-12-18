package lib

import (
	"net/http"

	"cloud.google.com/go/functions/apiv2/functionspb"
)

type LiveFunction struct {
	Location string
	Uri      string
	State    functionspb.Function_State
}

type FunctionResult struct {
	Response *http.Response
	Error    error
}

type FunctionAuthClient interface {
	Functions() []LiveFunction
	ExecuteFunction(location string) (*(chan FunctionResult), error)
	CheckFunctionAvailability(location string) error
}
