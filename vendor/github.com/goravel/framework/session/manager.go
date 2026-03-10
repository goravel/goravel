package session

import (
	"fmt"
	"sync"
	"time"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/foundation"
	contractssession "github.com/goravel/framework/contracts/session"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/session/driver"
	"github.com/goravel/framework/support/color"
)

// Ensure interface implementation
var _ contractssession.Manager = (*Manager)(nil)

type Manager struct {
	sessionPool sync.Pool
	config      config.Config
	json        foundation.Json

	drivers map[string]contractssession.Driver

	cookie        string
	defaultDriver string
	files         string
	gcInterval    int
	lifetime      int

	mu sync.RWMutex
}

func NewManager(config config.Config, json foundation.Json) *Manager {
	cookie := config.GetString("session.cookie")
	defaultDriver := config.GetString("session.default", "file")
	files := config.GetString("session.files")
	gcInterval := config.GetInt("session.gc_interval", 30)
	lifetime := config.GetInt("session.lifetime", 120)

	manager := &Manager{
		config: config,
		json:   json,

		cookie:        cookie,
		defaultDriver: defaultDriver,
		files:         files,
		gcInterval:    gcInterval,
		lifetime:      lifetime,

		drivers: make(map[string]contractssession.Driver),
		sessionPool: sync.Pool{New: func() any {
			return NewSession("", nil, json)
		},
		},
	}

	err := manager.registerDriver(defaultDriver)
	if err != nil {
		color.Errorln(errors.SessionDriverRegisterFailed.Args(err).Error())
	}
	return manager
}

func (m *Manager) BuildSession(handler contractssession.Driver, sessionID ...string) (contractssession.Session, error) {
	if handler == nil {
		return nil, errors.SessionDriverIsNotSet
	}

	session := m.acquireSession()
	session.SetDriver(handler).SetName(m.cookie)

	if len(sessionID) > 0 {
		session.SetID(sessionID[0])
	} else {
		session.SetID("")
	}

	return session, nil
}

func (m *Manager) Driver(name ...string) (contractssession.Driver, error) {
	driverName := m.defaultDriver
	if len(name) > 0 {
		driverName = name[0]
	}

	m.mu.RLock()
	driverInstance, instanceExists := m.drivers[driverName]
	m.mu.RUnlock()
	if instanceExists {
		return driverInstance, nil
	}

	err := m.registerDriver(driverName)
	if err != nil {
		return nil, err
	}

	return m.drivers[driverName], nil
}

func (m *Manager) ReleaseSession(session contractssession.Session) {
	session.Flush().
		SetDriver(nil).
		SetName("").
		SetID("")
	m.sessionPool.Put(session)
}

func (m *Manager) acquireSession() contractssession.Session {
	session := m.sessionPool.Get().(contractssession.Session)
	return session
}

func (m *Manager) custom(driver string) (contractssession.Driver, error) {
	via := m.config.Get(fmt.Sprintf("session.drivers.%s.via", driver))
	if custom, ok := via.(contractssession.Driver); ok {
		return custom, nil
	}
	if custom, ok := via.(func() (contractssession.Driver, error)); ok {
		return custom()
	}

	return nil, errors.SessionDriverContractNotFulfilled.Args(driver)
}

func (m *Manager) file() contractssession.Driver {
	return driver.NewFile(m.files, m.lifetime)
}

func (m *Manager) registerDriver(name string) error {
	driver := m.config.GetString(fmt.Sprintf("session.drivers.%s.driver", name))

	switch driver {
	case "file":
		driverInstance := m.file()
		m.drivers[name] = driverInstance
		m.startGcTimer(driverInstance)
	case "custom":
		driverInstance, err := m.custom(name)
		if err != nil {
			return err
		}
		m.drivers[name] = driverInstance
		m.startGcTimer(driverInstance)
	default:
		return errors.SessionDriverNotSupported.Args(driver)
	}

	return nil
}

func (m *Manager) startGcTimer(driverInstance contractssession.Driver) {
	if m.gcInterval <= 0 {
		// No need to start the timer if the interval is zero or negative
		return
	}

	ticker := time.NewTicker(time.Duration(m.gcInterval) * time.Minute)

	go func() {
		for range ticker.C {
			if err := driverInstance.Gc(m.lifetime * 60); err != nil {
				color.Errorf("Error performing garbage collection: %s\n", err)
			}
		}
	}()
}
