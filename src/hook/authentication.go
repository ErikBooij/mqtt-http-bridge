package hook

import (
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"sync"
)

type AuthHookInterface interface {
	mqtt.Hook
	AddUser(username, password string)
}

func Authentication(open bool) AuthHookInterface {
	return &authHook{
		open:  open,
		users: make(map[string]user),
	}
}

type authHook struct {
	mqtt.HookBase

	// open indicates no authentication is required
	open bool
	// users contains a map of users keyed on username
	users map[string]user

	mu sync.RWMutex
}

type user struct {
	username string
	password string
}

func (a *authHook) AddUser(username, password string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.users[username] = user{username: username, password: password}
}

func (a *authHook) ID() string {
	return "auth-hook"
}

func (a *authHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	return true
}

func (a *authHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	if a.open {
		return true
	}

	a.mu.RLock()
	defer a.mu.RUnlock()

	u, ok := a.users[string(cl.Properties.Username)]

	if !ok {
		return false
	}

	return u.password == string(pk.Connect.Password)
}

func (a *authHook) Provides(b byte) bool {
	switch b {
	case
		mqtt.OnConnectAuthenticate,
		mqtt.OnACLCheck:
		return true
	default:
		return false
	}
}
