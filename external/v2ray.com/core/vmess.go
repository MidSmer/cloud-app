package core

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"time"
	"unsafe"
	"v2ray.com/core/app/proxyman"
	app_inbound "v2ray.com/core/app/proxyman/inbound"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/common/task"
	"v2ray.com/core/common/uuid"
	feature_inbound "v2ray.com/core/features/inbound"
	"v2ray.com/core/features/policy"
	"v2ray.com/core/features/routing"
	"v2ray.com/core/proxy/vmess"
	"v2ray.com/core/proxy/vmess/encoding"
	"v2ray.com/core/proxy/vmess/inbound"
	"v2ray.com/core/transport/internet"
)

//vmess

type userByEmail struct {
	sync.Mutex
	cache           map[string]*protocol.MemoryUser
	defaultLevel    uint32
	defaultAlterIDs uint16
}

func newUserByEmail(config *inbound.DefaultConfig) *userByEmail {
	return &userByEmail{
		cache:           make(map[string]*protocol.MemoryUser),
		defaultLevel:    config.Level,
		defaultAlterIDs: uint16(config.AlterId),
	}
}

func (v *userByEmail) addNoLock(u *protocol.MemoryUser) bool {
	email := strings.ToLower(u.Email)
	user, found := v.cache[email]
	if found {
		return false
	}
	v.cache[email] = user
	return true
}

func (v *userByEmail) Add(u *protocol.MemoryUser) bool {
	v.Lock()
	defer v.Unlock()

	return v.addNoLock(u)
}

func (v *userByEmail) Get(email string) (*protocol.MemoryUser, bool) {
	email = strings.ToLower(email)

	v.Lock()
	defer v.Unlock()

	user, found := v.cache[email]
	if !found {
		id := uuid.New()
		rawAccount := &vmess.Account{
			Id:      id.String(),
			AlterId: uint32(v.defaultAlterIDs),
		}
		account, err := rawAccount.AsAccount()
		common.Must(err)
		user = &protocol.MemoryUser{
			Level:   v.defaultLevel,
			Email:   email,
			Account: account,
		}
		v.cache[email] = user
	}
	return user, found
}

func (v *userByEmail) Remove(email string) bool {
	email = strings.ToLower(email)

	v.Lock()
	defer v.Unlock()

	if _, found := v.cache[email]; !found {
		return false
	}
	delete(v.cache, email)
	return true
}

// Handler is an inbound connection handler that handles messages in VMess protocol.
type VmessInboundHandler struct {
	policyManager         policy.Manager
	inboundHandlerManager feature_inbound.Manager
	clients               *vmess.TimedUserValidator
	usersByEmail          *userByEmail
	detours               *inbound.DetourConfig
	sessionHistory        *encoding.SessionHistory
	secure                bool
}

// New creates a new VMess inbound handler.
func NewVmessInboundHandler(ctx context.Context, config *inbound.Config) (*VmessInboundHandler, error) {
	inboundHandlerManager, err := app_inbound.New(ctx, &proxyman.InboundConfig{})
	if err != nil {
		return nil, err
	}

	handler := &VmessInboundHandler{
		policyManager:         policy.DefaultManager{},
		inboundHandlerManager: inboundHandlerManager,
		clients:               vmess.NewTimedUserValidator(protocol.DefaultIDHash),
		detours:               config.Detour,
		usersByEmail:          newUserByEmail(config.GetDefaultValue()),
		sessionHistory:        encoding.NewSessionHistory(),
		secure:                config.SecureEncryptionOnly,
	}

	return handler, nil
}

// Close implements common.Closable.
func (h *VmessInboundHandler) Close() error {
	return errors.Combine(
		h.clients.Close(),
		h.sessionHistory.Close(),
		common.Close(h.usersByEmail))
}

// Network implements proxy.Inbound.Network().
func (*VmessInboundHandler) Network() []net.Network {
	return []net.Network{net.Network_TCP}
}

func (h *VmessInboundHandler) GetUser(email string) *protocol.MemoryUser {
	user, existing := h.usersByEmail.Get(email)
	if !existing {
		h.clients.Add(user)
	}
	return user
}

func (h *VmessInboundHandler) AddUser(ctx context.Context, user *protocol.MemoryUser) error {
	if len(user.Email) > 0 && !h.usersByEmail.Add(user) {
		return newError("User ", user.Email, " already exists.")
	}
	return h.clients.Add(user)
}

func (h *VmessInboundHandler) RemoveUser(ctx context.Context, email string) error {
	if len(email) == 0 {
		return newError("Email must not be empty.")
	}
	if !h.usersByEmail.Remove(email) {
		return newError("User ", email, " not found.")
	}
	h.clients.Remove(email)
	return nil
}

func transferResponse(timer signal.ActivityUpdater, session *encoding.ServerSession, request *protocol.RequestHeader, response *protocol.ResponseHeader, input buf.Reader, output *buf.BufferedWriter) error {
	session.EncodeResponseHeader(response, output)

	bodyWriter := session.EncodeResponseBody(request, output)

	{
		// Optimize for small response packet
		data, err := input.ReadMultiBuffer()
		if err != nil {
			return err
		}

		if err := bodyWriter.WriteMultiBuffer(data); err != nil {
			return err
		}
	}

	if err := output.SetBuffered(false); err != nil {
		return err
	}

	if err := buf.Copy(input, bodyWriter, buf.UpdateActivity(timer)); err != nil {
		return err
	}

	if request.Option.Has(protocol.RequestOptionChunkStream) {
		if err := bodyWriter.WriteMultiBuffer(buf.MultiBuffer{}); err != nil {
			return err
		}
	}

	return nil
}

func isInsecureEncryption(s protocol.SecurityType) bool {
	return s == protocol.SecurityType_NONE || s == protocol.SecurityType_LEGACY || s == protocol.SecurityType_UNKNOWN
}

// Process implements proxy.Inbound.Process().
func (h *VmessInboundHandler) Process(ctx context.Context, network net.Network, connection internet.Connection, dispatcher routing.Dispatcher) error {
	sessionPolicy := h.policyManager.ForLevel(0)
	if err := connection.SetReadDeadline(time.Now().Add(sessionPolicy.Timeouts.Handshake)); err != nil {
		return newError("unable to set read deadline").Base(err).AtWarning()
	}

	reader := &buf.BufferedReader{Reader: buf.NewReader(connection)}
	svrSession := encoding.NewServerSession(h.clients, h.sessionHistory)
	request, err := svrSession.DecodeRequestHeader(reader)

	typeOfSess := reflect.TypeOf(svrSession)
	if typeOfSess.Kind() == reflect.Ptr {
		typeOfSess = typeOfSess.Elem()
	}
	if bodyKeyType, ok := typeOfSess.FieldByName("requestBodyKey"); ok {
		requestBodyKey := (*[16]byte)(unsafe.Pointer(uintptr(unsafe.Pointer(svrSession)) + bodyKeyType.Offset))
		fmt.Println(requestBodyKey[:])
	}

	if err != nil {
		if errors.Cause(err) != io.EOF {
			log.Record(&log.AccessMessage{
				From:   connection.RemoteAddr(),
				To:     "",
				Status: log.AccessRejected,
				Reason: err,
			})
			err = newError("invalid request from ", connection.RemoteAddr()).Base(err).AtInfo()
		}
		return err
	}

	if h.secure && isInsecureEncryption(request.Security) {
		log.Record(&log.AccessMessage{
			From:   connection.RemoteAddr(),
			To:     "",
			Status: log.AccessRejected,
			Reason: "Insecure encryption",
		})
		return newError("client is using insecure encryption: ", request.Security)
	}

	if request.Command != protocol.RequestCommandMux {
		log.Record(&log.AccessMessage{
			From:   connection.RemoteAddr(),
			To:     request.Destination(),
			Status: log.AccessAccepted,
			Reason: "",
		})
	}

	newError("received request for ", request.Destination()).WriteToLog(session.ExportIDToError(ctx))

	if err := connection.SetReadDeadline(time.Time{}); err != nil {
		newError("unable to set back read deadline").Base(err).WriteToLog(session.ExportIDToError(ctx))
	}

	inbound := session.InboundFromContext(ctx)
	if inbound == nil {
		panic("no inbound metadata")
	}
	inbound.User = request.User

	sessionPolicy = h.policyManager.ForLevel(request.User.Level)

	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, sessionPolicy.Timeouts.ConnectionIdle)

	ctx = policy.ContextWithBufferPolicy(ctx, sessionPolicy.Buffer)
	link, err := dispatcher.Dispatch(ctx, request.Destination())
	if err != nil {
		return newError("failed to dispatch request to ", request.Destination()).Base(err)
	}

	requestDone := func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.DownlinkOnly)

		bodyReader := svrSession.DecodeRequestBody(request, reader)
		if err := buf.Copy(bodyReader, link.Writer, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to transfer request").Base(err)
		}
		return nil
	}

	responseDone := func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.UplinkOnly)

		writer := buf.NewBufferedWriter(buf.NewWriter(connection))
		defer writer.Flush()

		response := &protocol.ResponseHeader{
			Command: h.generateCommand(ctx, request),
		}
		return transferResponse(timer, svrSession, request, response, link.Reader, writer)
	}

	var requestDonePost = task.OnSuccess(requestDone, task.Close(link.Writer))
	if err := task.Run(ctx, requestDonePost, responseDone); err != nil {
		common.Interrupt(link.Reader)
		common.Interrupt(link.Writer)
		return newError("connection ends").Base(err)
	}

	return nil
}

func (h *VmessInboundHandler) generateCommand(ctx context.Context, request *protocol.RequestHeader) protocol.ResponseCommand {
	if h.detours != nil {
		tag := h.detours.To
		if h.inboundHandlerManager != nil {
			handler, err := h.inboundHandlerManager.GetHandler(ctx, tag)
			if err != nil {
				newError("failed to get detour handler: ", tag).Base(err).AtWarning().WriteToLog(session.ExportIDToError(ctx))
				return nil
			}
			proxyHandler, port, availableMin := handler.GetRandomInboundProxy()
			inboundHandler, ok := proxyHandler.(*VmessInboundHandler)
			if ok && inboundHandler != nil {
				if availableMin > 255 {
					availableMin = 255
				}

				newError("pick detour handler for port ", port, " for ", availableMin, " minutes.").AtDebug().WriteToLog(session.ExportIDToError(ctx))
				user := inboundHandler.GetUser(request.User.Email)
				if user == nil {
					return nil
				}
				account := user.Account.(*vmess.MemoryAccount)
				return &protocol.CommandSwitchAccount{
					Port:     port,
					ID:       account.ID.UUID(),
					AlterIds: uint16(len(account.AlterIDs)),
					Level:    user.Level,
					ValidMin: byte(availableMin),
				}
			}
		}
	}

	return nil
}
