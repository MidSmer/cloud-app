package route

import (
	"github.com/MidSmer/cloud-app/external/v2ray.com/core"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
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

	//defer conn.Close()

	handler, err := core.NewWebsocketInboundHandler()
	handler.Start(w, r, conn)

}
