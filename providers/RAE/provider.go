package rae

import "github.com/minskylab/palapi"

type Provider struct {
	currentStatus palapi.ProviderStatus
	baseURL       string
}

func NewProvider() (*Provider, error) {
	return &Provider{
		baseURL:       "https://dle.rae.es",
		currentStatus: palapi.IDLE,
	}, nil
}
