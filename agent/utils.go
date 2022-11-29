package agent

import (
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func broadcastMessage(message []byte, connections *map[string]*net.Conn) ***REMOVED***
	for _, conn := range *connections ***REMOVED***
		wsutil.WriteServerMessage(*conn, ws.OpText, message)
	***REMOVED***
***REMOVED***
