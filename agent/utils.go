package agent

import (
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func broadcastMessage(message []byte, connections *map[string]*net.Conn) {
	for _, conn := range *connections {
		wsutil.WriteServerMessage(*conn, ws.OpText, message)
	}
}
