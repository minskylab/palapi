package wordreference

import "github.com/minskylab/palapi"

type Provider struct {
	currentStatus palapi.ProviderStatus
	baseURL       string
}

func NewProvider() (*Provider, error) {
	return &Provider{
		currentStatus: palapi.IDLE,
		baseURL:       "https://www.wordreference.com/sinonimos",
	}, nil
}
