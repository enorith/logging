package logging

import (
	"sync"

	"gopkg.in/yaml.v3"
)

type DriverRegister func(config *yaml.Node) (Logger, error)

type Manager struct {
	driverRegisters map[string]DriverRegister
	mu              sync.RWMutex
}

func (m *Manager) Resolve(driver string, fn DriverRegister) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.driverRegisters[driver] = fn
}
