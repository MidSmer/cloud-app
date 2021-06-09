package route

import (
	"github.com/MidSmer/cloud-app/external/v2ray.com/core"
	"github.com/gorilla/websocket"
	"net/http"
	"time"

	_ "v2ray.com/core/transport/internet/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: time.Second * 4,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func WS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Logger.Error("ws connection fail!", err)
		return
	}

	handler, err := core.NewWebsocketInboundHandler()
	if err != nil {
		Logger.Error("web socket inbound handler create fail!", err)
		return
	}

	handler.Start(w, r, conn)
}
