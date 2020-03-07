package palapi

import "time"

type Feature string

const Definitions Feature = "definition"
const Synonyms Feature = "synonyms"
const Antonyms Feature = "antonyms"
const Examples Feature = "examples"
const Frequency Feature = "frequency"

type Report struct {
	Word               string
	At                 time.Time
	QueryDuration      time.Duration
	ExtractionDuration time.Duration

	Definitions *[]WordDefinition
	Frequency   *WordFrequency
	Synonyms    *[]string
	Antonyms    *[]string
	Examples    *[]string
}

type ProviderStatus string

const IDLE ProviderStatus = "idle"
const SCRAPING ProviderStatus = "scraping"
const PROCESSING ProviderStatus = "processing"

type Provider interface {
	Source() Source
	AvailableFeatures() []Feature
	FindWord(word string) (*Report, error)
	Status() ProviderStatus
}
