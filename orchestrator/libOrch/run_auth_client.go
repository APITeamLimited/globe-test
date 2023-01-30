package libOrch

import (
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/gorilla/websocket"
)

type LiveService struct {
	Location string
	Uri      string
	State    runpb.Condition_State
}

type RunAuthClient interface {
	Services() []LiveService
	ExecuteService(location string) (*websocket.Conn, error)
	CheckServiceAvailability(location string) error
}
