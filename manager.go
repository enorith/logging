package logging

import (
	"fmt"
	"sync"

	"github.com/enorith/logging/writers"
	"go.uber.org/zap"
)

type ChannelRegister func(conf zap.Config) (*zap.Logger, error)

var DefaultChannel = "default"

var DefaultManager = NewManager()

type Manager struct {
	crs      map[string]ChannelRegister
	channels map[string]*zap.Logger
	rotates  []*writers.RotateFileWriter
	using    string
	mu       sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		crs:      make(map[string]ChannelRegister),
		channels: make(map[string]*zap.Logger),
		rotates:  make([]*writers.RotateFileWriter, 0),
		using:    DefaultChannel,
		mu:       sync.RWMutex{},
	}
}

func (m *Manager) Resolve(channel string, cr ChannelRegister) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.crs[channel] = cr
}

func (m *Manager) Channel(channel ...string) (*zap.Logger, error) {
	var ch string
	if len(channel) > 0 {
		ch = channel[0]
	}
	if ch == "" {
		ch = DefaultChannel
	}
	m.mu.RLock()
	logger, cok := m.channels[ch]
	m.mu.RUnlock()
	if cok {
		return logger, nil
	}

	m.mu.RLock()
	cr, ok := m.crs[ch]
	m.mu.RUnlock()
	if ok {
		logger, err := cr(zap.NewProductionConfig())
		if err != nil {
			return nil, err
		}
		m.mu.Lock()
		m.channels[ch] = logger
		m.mu.Unlock()
		return logger, nil
	}

	return nil, fmt.Errorf("unregisterd log channel [%s]", ch)
}

func (m *Manager) Using(defaultChannel string) *Manager {
	m.using = defaultChannel

	return m
}

func (m *Manager) AddRotate(r *writers.RotateFileWriter) *Manager {
	m.rotates = append(m.rotates, r)

	return m
}

func (m *Manager) Sync() {
	for _, l := range m.channels {
		l.Sync()
	}
}

func (m *Manager) Cleanup() {
	for _, r := range m.rotates {
		r.Cleanup()
	}
}
