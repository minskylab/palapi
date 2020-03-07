package palapi

import "time"

type Head struct {
	Word WordID
	Level int64
}

type Manager struct {
	Providers []Provider
	DeepExploration int64
	LastHeads []*Head

	Persistence Persistence

	MaxAntiquityOfWord time.Duration
}
