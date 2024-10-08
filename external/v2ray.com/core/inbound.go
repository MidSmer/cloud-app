package core

import (
	"context"

	"github.com/v2fly/v2ray-core/v5/app/dispatcher"
	"github.com/v2fly/v2ray-core/v5/features/policy"
	"github.com/v2fly/v2ray-core/v5/proxy/vmess/inbound"

	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/common/session"
	"github.com/v2fly/v2ray-core/v5/features/routing"
	"github.com/v2fly/v2ray-core/v5/features/stats"
	"github.com/v2fly/v2ray-core/v5/transport/internet"
	"github.com/v2fly/v2ray-core/v5/transport/internet/tcp"

	_ "github.com/v2fly/v2ray-core/v5/transport/internet/websocket"
)

type worker interface {
	Start() error
	Close() error
}

type tcpWorker struct {
	address         net.Address
	port            net.Port
	stream          *internet.MemoryStreamConfig
	recvOrigDest    bool
	dispatcher      routing.Dispatcher
	uplinkCounter   stats.Counter
	downlinkCounter stats.Counter
}

func getTProxyType(s *internet.MemoryStreamConfig) internet.SocketConfig_TProxyMode {
	if s == nil || s.SocketSettings == nil {
		return internet.SocketConfig_Off
	}
	return s.SocketSettings.Tproxy
}

func (w *tcpWorker) callback(conn internet.Connection) {
	ctx, cancel := context.WithCancel(context.Background())
	sid := session.NewID()
	ctx = session.ContextWithID(ctx, sid)

	defer cancel()

	if w.recvOrigDest {
		var dest net.Destination
		switch getTProxyType(w.stream) {
		case internet.SocketConfig_Redirect:
			d, err := tcp.GetOriginalDestination(conn)
			if err != nil {
				newError("failed to get original destination").Base(err).WriteToLog(session.ExportIDToError(ctx))
				return
			} else {
				dest = d
			}
		case internet.SocketConfig_TProxy:
			dest = net.DestinationFromAddr(conn.LocalAddr())
		}
		if dest.IsValid() {
			ctx = session.ContextWithOutbound(ctx, &session.Outbound{
				Target: dest,
			})
		}
	}
	ctx = session.ContextWithInbound(ctx, &session.Inbound{
		Source:  net.DestinationFromAddr(conn.RemoteAddr()),
		Gateway: net.TCPDestination(w.address, w.port),
		Tag:     "inbound",
	})
	content := new(session.Content)
	ctx = session.ContextWithContent(ctx, content)
	if w.uplinkCounter != nil || w.downlinkCounter != nil {
		conn = &internet.StatCouterConnection{
			Connection:   conn,
			ReadCounter:  w.uplinkCounter,
			WriteCounter: w.downlinkCounter,
		}
	}

	vmessInbound, err := NewVmessInboundHandler(ctx, &inbound.Config{
		SecureEncryptionOnly: true,
	})
	if err != nil {
		newError("failed create vmess inbound").Base(err).WriteToLog(session.ExportIDToError(ctx))
		return
	}

	accountManager := GetAccountManagerInstance()
	for _, a := range accountManager.Get() {
		vmessInbound.AddUser(ctx, a)
	}

	err = vmessInbound.Process(ctx, net.Network_TCP, conn, w.dispatcher)
	if err != nil {
		newError("failed vmess inbound process").Base(err).WriteToLog(session.ExportIDToError(ctx))
		return
	}

	if err := conn.Close(); err != nil {
		newError("failed to close connection").Base(err).WriteToLog(session.ExportIDToError(ctx))
	}
}

func (w *tcpWorker) Start() error {
	return nil
}

func (w *tcpWorker) Close() error {
	return nil
}

func getStatCounter() (stats.Counter, stats.Counter) {
	var uplinkCounter stats.Counter
	var downlinkCounter stats.Counter

	return uplinkCounter, downlinkCounter
}

func NewInboundHandler(ctx context.Context) (*tcpWorker, error) {
	uplinkCounter, downlinkCounter := getStatCounter()

	mss, err := internet.ToMemoryStreamConfig(&internet.StreamConfig{
		ProtocolName: "websocket",
	})
	if err != nil {
		return nil, newError("failed to parse stream config").Base(err).AtWarning()
	}

	dp := &dispatcher.DefaultDispatcher{}
	om, err := NewOutboundManager(ctx)
	if err != nil {
		return nil, newError("failed create outbound manager").Base(err).AtWarning()
	}

	err = dp.Init(&dispatcher.Config{}, om, routing.DefaultRouter{}, policy.DefaultManager{}, stats.NoopManager{})
	if err != nil {
		return nil, newError("failed dispatcher init").Base(err).AtWarning()
	}

	worker := &tcpWorker{
		//address:         net.AnyIP,
		//port:            net.Port(80),
		stream:          mss,
		recvOrigDest:    false,
		dispatcher:      dp,
		uplinkCounter:   uplinkCounter,
		downlinkCounter: downlinkCounter,
	}

	return worker, nil
}
