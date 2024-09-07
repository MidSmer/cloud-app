package core

import (
	"sync"

	"github.com/v2fly/v2ray-core/v5/common/protocol"
	"github.com/v2fly/v2ray-core/v5/common/uuid"
	"github.com/v2fly/v2ray-core/v5/proxy/vmess"
)

type accountManager struct {
	sync.RWMutex
	users []*protocol.MemoryUser
}

var (
	amInstance *accountManager
	once       sync.Once
)

func GetAccountManagerInstance() *accountManager {
	once.Do(func() {
		amInstance = &accountManager{
			users: make([]*protocol.MemoryUser, 0, 16),
		}
	})

	return amInstance
}

func (m *accountManager) Generate() error {
	id := uuid.New()

	return m.Add(id.String())
}

func (m *accountManager) Add(idStr string) error {
	m.Lock()
	defer m.Unlock()

	id, err := uuid.ParseString(idStr)
	if err != nil {
		return err
	}

	for _, user := range m.users {
		userId := user.Account.(*vmess.MemoryAccount).ID.UUID()
		if id.Equals(&userId) {
			return newError("failed uuid exist")
		}
	}

	account, err := (&vmess.Account{Id: id.String(), AlterId: 16}).AsAccount()
	if err != nil {
		newError("failed create vmess account").Base(err).WriteToLog()
		return err
	}
	memoryUser := &protocol.MemoryUser{
		Account: account,
	}

	m.users = append(m.users, memoryUser)

	return nil
}

func (m *accountManager) Get() []*protocol.MemoryUser {
	defer m.RUnlock()
	m.RLock()

	return m.users
}

func (m *accountManager) GetID() []string {
	defer m.RUnlock()
	m.RLock()

	var id []string

	for _, a := range m.users {
		account := a.Account.(*vmess.MemoryAccount)
		id = append(id, account.ID.String())
	}

	return id
}
