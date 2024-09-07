package core

import (
	"context"

	core "github.com/v2fly/v2ray-core/v5"
	"github.com/v2fly/v2ray-core/v5/common/serial"
	"github.com/v2fly/v2ray-core/v5/features/dns/localdns"
	"github.com/v2fly/v2ray-core/v5/features/policy"
	"github.com/v2fly/v2ray-core/v5/proxy/freedom"

	"github.com/v2fly/v2ray-core/v5/app/proxyman"
	proxymanOutbound "github.com/v2fly/v2ray-core/v5/app/proxyman/outbound"
	"github.com/v2fly/v2ray-core/v5/common"
	"github.com/v2fly/v2ray-core/v5/common/mux"
	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/common/session"
	"github.com/v2fly/v2ray-core/v5/features/outbound"
	"github.com/v2fly/v2ray-core/v5/proxy"
	"github.com/v2fly/v2ray-core/v5/transport"
	"github.com/v2fly/v2ray-core/v5/transport/internet"
	"github.com/v2fly/v2ray-core/v5/transport/pipe"
)

// OutboundHandler is an implements of outbound.OutboundHandler.
type OutboundHandler struct {
	tag             string
	senderSettings  *proxyman.SenderConfig
	streamSettings  *internet.MemoryStreamConfig
	proxy           proxy.Outbound
	outboundManager outbound.Manager
	mux             *mux.ClientManager
}

// NewHandler create a new OutboundHandler based on the given configuration.
func NewOutboundHandler(ctx context.Context, config *core.OutboundHandlerConfig) (outbound.Handler, error) {
	h := &OutboundHandler{
		tag: config.Tag,
	}

	if config.SenderSettings != nil {
		senderSettings, err := serial.GetInstanceOf(config.SenderSettings)
		if err != nil {
			return nil, err
		}
		switch s := senderSettings.(type) {
		case *proxyman.SenderConfig:
			h.senderSettings = s
			mss, err := internet.ToMemoryStreamConfig(s.StreamSettings)
			if err != nil {
				return nil, newError("failed to parse stream settings").Base(err).AtWarning()
			}
			h.streamSettings = mss
		default:
			return nil, newError("settings is not SenderConfig")
		}
	}

	proxyConfig, err := serial.GetInstanceOf(config.ProxySettings)
	if err != nil {
		return nil, err
	}

	proxyHandler := new(freedom.Handler)
	err = proxyHandler.Init(proxyConfig.(*freedom.Config), policy.DefaultManager{}, localdns.New())
	if err != nil {
		return nil, err
	}

	if h.senderSettings != nil && h.senderSettings.MultiplexSettings != nil && h.senderSettings.MultiplexSettings.Enabled {
		config := h.senderSettings.MultiplexSettings
		if config.Concurrency < 1 || config.Concurrency > 1024 {
			return nil, newError("invalid mux concurrency: ", config.Concurrency).AtWarning()
		}
		h.mux = &mux.ClientManager{
			Picker: &mux.IncrementalWorkerPicker{
				Factory: &mux.DialingWorkerFactory{
					Proxy:  proxyHandler,
					Dialer: h,
					Strategy: mux.ClientStrategy{
						MaxConcurrency: config.Concurrency,
						MaxConnection:  128,
					},
				},
			},
		}
	}

	h.proxy = proxyHandler
	return h, nil
}

func NewOutboundManager(ctx context.Context) (*proxymanOutbound.Manager, error) {
	manager, err := proxymanOutbound.New(ctx, &proxyman.OutboundConfig{})
	if err != nil {
		return nil, err
	}

	handler, err := NewOutboundHandler(ctx, &core.OutboundHandlerConfig{
		Tag:            "",
		SenderSettings: serial.ToTypedMessage(&proxyman.SenderConfig{}),
		ProxySettings:  serial.ToTypedMessage(&freedom.Config{}),
	})
	if err != nil {
		return nil, err
	}

	err = manager.AddHandler(ctx, handler)
	if err != nil {
		return nil, err
	}

	return manager, nil
}

// Tag implements outbound.OutboundHandler.
func (h *OutboundHandler) Tag() string {
	return h.tag
}

// Dispatch implements proxy.Outbound.Dispatch.
func (h *OutboundHandler) Dispatch(ctx context.Context, link *transport.Link) {
	if h.mux != nil {
		if err := h.mux.Dispatch(ctx, link); err != nil {
			newError("failed to process mux outbound traffic").Base(err).WriteToLog(session.ExportIDToError(ctx))
			common.Interrupt(link.Writer)
		}
	} else {
		if err := h.proxy.Process(ctx, link, h); err != nil {
			// Ensure outbound ray is properly closed.
			newError("failed to process outbound traffic").Base(err).WriteToLog(session.ExportIDToError(ctx))
			common.Interrupt(link.Writer)
		} else {
			common.Must(common.Close(link.Writer))
		}
		common.Interrupt(link.Reader)
	}
}

// Address implements internet.Dialer.
func (h *OutboundHandler) Address() net.Address {
	if h.senderSettings == nil || h.senderSettings.Via == nil {
		return nil
	}
	return h.senderSettings.Via.AsAddress()
}

// Dial implements internet.Dialer.
func (h *OutboundHandler) Dial(ctx context.Context, dest net.Destination) (internet.Connection, error) {
	if h.senderSettings != nil {
		if h.senderSettings.ProxySettings.HasTag() {
			tag := h.senderSettings.ProxySettings.Tag
			handler := h.outboundManager.GetHandler(tag)
			if handler != nil {
				newError("proxying to ", tag, " for dest ", dest).AtDebug().WriteToLog(session.ExportIDToError(ctx))
				ctx = session.ContextWithOutbound(ctx, &session.Outbound{
					Target: dest,
				})

				opts := pipe.OptionsFromContext(ctx)
				uplinkReader, uplinkWriter := pipe.New(opts...)
				downlinkReader, downlinkWriter := pipe.New(opts...)

				go handler.Dispatch(ctx, &transport.Link{Reader: uplinkReader, Writer: downlinkWriter})
				return net.NewConnection(net.ConnectionInputMulti(uplinkWriter), net.ConnectionOutputMulti(downlinkReader)), nil
			}

			newError("failed to get outbound handler with tag: ", tag).AtWarning().WriteToLog(session.ExportIDToError(ctx))
		}

		if h.senderSettings.Via != nil {
			outbound := session.OutboundFromContext(ctx)
			if outbound == nil {
				outbound = new(session.Outbound)
				ctx = session.ContextWithOutbound(ctx, outbound)
			}
			outbound.Gateway = h.senderSettings.Via.AsAddress()
		}
	}

	return internet.Dial(ctx, dest, h.streamSettings)
}

// GetOutbound implements proxy.GetOutbound.
func (h *OutboundHandler) GetOutbound() proxy.Outbound {
	return h.proxy
}

// Start implements common.Runnable.
func (h *OutboundHandler) Start() error {
	return nil
}

// Close implements common.Closable.
func (h *OutboundHandler) Close() error {
	common.Close(h.mux)
	return nil
}
