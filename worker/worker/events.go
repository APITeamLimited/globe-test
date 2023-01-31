package worker

import (
	"encoding/json"
	"fmt"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/gorilla/websocket"
)

type eventChannels struct {
	goMessageChannel    chan string
	childUpdatesChannel chan string
}

func getEventChannels(gs libWorker.BaseGlobalState) eventChannels {
	// Listen on websocket for messages from orchestrator
	channels := eventChannels{
		goMessageChannel:    make(chan string),
		childUpdatesChannel: make(chan string),
	}

	go func() {
		for {
			gs.ConnReadMutex().Lock()
			messageKind, message, err := gs.Conn().ReadMessage()
			gs.ConnReadMutex().Unlock()

			if messageKind == websocket.CloseMessage {
				return
			}

			if err != nil {
				gs.GetRunAbortFunc()()
				return
			}

			if messageKind != websocket.TextMessage {
				continue
			}

			eventMessage := lib.EventMessage{}

			err = json.Unmarshal(message, &eventMessage)
			if err != nil {
				fmt.Println("Error unmarshalling event message", err, string(message))
				return
			}

			switch eventMessage.Variant {
			case lib.GO_MESSAGE_TYPE:
				channels.goMessageChannel <- eventMessage.Data
			case lib.CHILD_USER_UPDATE:
				channels.childUpdatesChannel <- eventMessage.Data
			}
		}
	}()

	return channels
}
