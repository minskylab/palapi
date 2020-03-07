package main

import (
	"time"

	"github.com/minskylab/palapi"
	"github.com/pkg/errors"
)

type Provider struct {
	currentStatus palapi.ProviderStatus
}

func (p *Provider) Source() palapi.Source {
	return  palapi.Source{
		ID:        "rae",
		Name:      "RAE Website",
		Relevancy: 0.7,
		URL:       "https://dle.rae.es",
		Metadata: map[string]string{
			"type": "scraper",
		},
	}
}

func (p *Provider) AvailableFeatures() []palapi.Feature {
	return []palapi.Feature{
		palapi.Definitions,
		palapi.Examples,
	}
}

func (p *Provider) FindWord(word string) (*palapi.Report, error) {
	p.currentStatus = palapi.SCRAPING
	definitions, examples, extractionDur, err := scrapRAE(word)
	if err != nil {
		return nil, errors.Wrap(err, "scrape not worked correctly")
	}
	p.currentStatus = palapi.PROCESSING
	p.currentStatus = palapi.IDLE
	return &palapi.Report{
		Word:               word,
		At:                 time.Now(),
		QueryDuration:      0,
		ExtractionDuration: extractionDur,
		Definitions:        &definitions,
		Frequency:          nil,
		Synonyms:           nil,
		Antonyms:           nil,
		Examples:           &examples,
	}, nil
}

func (p *Provider) Status() palapi.ProviderStatus {
	return p.currentStatus
}

