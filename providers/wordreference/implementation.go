package wordreference

import (
	"strings"
	"time"

	"github.com/minskylab/palapi"
)

func (p *Provider) Source() palapi.Source {
	return palapi.Source{ // TODO: Optimize this
		ID:        "word-reference-synonyms",
		Name:      "Word Reference Spanish Synonyms",
		Relevancy: 0.68,
		URL:       "https://www.wordreference.com",
		Metadata: map[string]string{
			"type": "scraper",
		},
	}
}

func (p *Provider) AvailableFeatures() []palapi.Feature {
	return []palapi.Feature{
		palapi.Synonyms,
		palapi.Antonyms,
	}
}

func (p *Provider) FindWord(word string) (*palapi.Report, error) {
	word = strings.Trim(word, " ,.\\/'\"")
	t1 := time.Now()

	p.currentStatus = palapi.PROCESSING
	p.currentStatus = palapi.SCRAPING
	synonyms, antonyms, duration, err := p.extractByWord(word)
	if err != nil {
		return nil, err
	}

	return &palapi.Report{
		Word:               word,
		At:                 time.Now(),
		QueryDuration:      time.Since(t1),
		ExtractionDuration: duration,
		Definitions:        nil,
		Frequency:          nil,
		Synonyms:           &synonyms,
		Antonyms:           &antonyms,
		Examples:           nil,
	}, nil
}

func (p *Provider) Status() palapi.ProviderStatus {
	return p.currentStatus
}
