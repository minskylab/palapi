package WordReference

import "github.com/minskylab/palapi"

type Provider struct {
	currentStatus palapi.ProviderStatus
}

func (p *Provider) Source() palapi.Source {
	return palapi.Source{ // TODO: Optimize this
		ID:        "wordreference",
		Name:      "Word Reference Spanish Synonyms",
		Relevancy: 0.68,
		URL:       "https://www.wordreference.com",
		Metadata: map[string]string{
			"type": "scraper",
		},
	}
}

func (p *Provider) AvailableFeatures() []palapi.Feature {
	panic("implement me")
}

func (p *Provider) FindWord(word string) (*palapi.Report, error) {
	panic("implement me")
}

func (p *Provider) Status() palapi.ProviderStatus {
	panic("implement me")
}
