package palapi

import "time"

type Head struct {
	Word  WordID
	Level int64
}

type Manager struct {
	providers      []Provider
	maxExploration int64
	heads          []*Head

	persistence Persistence

	maxAntiquityOfWord time.Duration
}

func NewManager(persistence Persistence, exploration int64) (*Manager, error) {
	return &Manager{
		providers:          []Provider{},
		maxExploration:     exploration,
		heads:              []*Head{},
		persistence:        persistence,
		maxAntiquityOfWord: 5 * 24 * time.Hour,
	}, nil
}
