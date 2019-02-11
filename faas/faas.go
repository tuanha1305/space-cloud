package faas

import (
	"encoding/json"
	"sync"
	"time"

	nats "github.com/nats-io/go-nats"

	"github.com/spaceuptech/space-cloud/config"
)

// Module is responsible for FaaS
type Module struct {
	sync.RWMutex
	nc      *nats.Conn
	enabled bool
}

// Init returns a new instance of the FaaS module
func Init() *Module {
	return new(Module)
}

// SetConfig set the config required by the FaaS module
func (m *Module) SetConfig(faas *config.FaaS) error {
	m.Lock()
	defer m.Unlock()

	if faas == nil || !faas.Enabled {
		m.enabled = false
		return nil
	}

	nc, err := nats.Connect(faas.Nats)
	if err != nil {
		return err
	}

	m.nc = nc
	m.enabled = true
	return nil
}

// Request calls a function on the provided engine
func (m *Module) Request(engine, function string, timeout int, val interface{}) ([]byte, error) {
	m.RLock()
	defer m.RUnlock()

	// Marshal the object into json
	data, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	// Send request over nats
	subject := "faas:" + engine + ":" + function
	msg, err := m.nc.Request(subject, data, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}

func (m *Module) isEnabled() bool {
	m.RLock()
	defer m.RUnlock()

	return m.enabled
}
