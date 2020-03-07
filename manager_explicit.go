package palapi

import "github.com/pkg/errors"

func (m *Manager) RegisterProvider(provider Provider) error {
	if m.providers == nil {
		m.providers = []Provider{}
	}

	m.providers = append(m.providers, provider)

	return nil
}

func (m *Manager) SetExplorationRange(deep int64) {
	m.maxExploration = deep
}

func (m *Manager) UnregisterProvider(id SourceID) error {
	for i, provider := range m.providers {
		if provider.Source().ID == id {
			m.providers = append(m.providers[:i], m.providers[i+1:]...)
			return nil
		}
	}

	return errors.New("provider id not found in our providers list")
}
