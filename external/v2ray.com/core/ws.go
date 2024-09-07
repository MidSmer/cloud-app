package core

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/v2fly/v2ray-core/v5/common/net"
	httpProto "github.com/v2fly/v2ray-core/v5/common/protocol/http"
)

type Connection interface {
	net.Conn
}

type ConnHandler func(Connection)

type websocketHandler struct {
	ctx     context.Context
	addConn ConnHandler
}

func (h *websocketHandler) Start(writer http.ResponseWriter, request *http.Request, conn *websocket.Conn) {

	forwardedAddrs := httpProto.ParseXForwardedFor(request.Header)
	remoteAddr := conn.RemoteAddr()
	if len(forwardedAddrs) > 0 && forwardedAddrs[0].Family().IsIP() {
		remoteAddr.(*net.TCPAddr).IP = forwardedAddrs[0].IP()
	}

	h.addConn(newConnection(conn, remoteAddr))
}

func NewWebsocketInboundHandler() (*websocketHandler, error) {
	ctx := context.Background()

	worker, err := NewInboundHandler(ctx)

	if err != nil {
		return nil, err
	}

	handler := &websocketHandler{
		ctx: ctx,
		addConn: func(conn Connection) {
			go worker.callback(conn)
		},
	}

	return handler, nil
}
