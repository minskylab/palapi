package palapi

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (m *Manager) reportWord(word string, deepest int64) (*Word, error) {
	if deepest > m.maxExploration {
		return nil, nil
	}

	log.WithField("deep", deepest).Infof("reporting: %s", word)

	w, err := m.persistence.GetWord(word)
	if err == nil && !time.Now().After(w.LastUpdate.Add(m.maxAntiquityOfWord)) {
		// saved word
		return w, nil
	}
	log.Debug("continue the analysis of the word in the storage")

	if err != storm.ErrNotFound {
		return nil, errors.Wrap(err, "persistence layer 'get' failed")
	}

	log.Debug("the word isn't in the database or need to be updated")
	// word not found or needs update
	// providers sorted
	sort.SliceStable(m.providers, func(i, j int) bool {
		return m.providers[j].Source().Relevancy < m.providers[i].Source().Relevancy
	})

	var syntheticWord Word

	definitions := make([]WordDefinition, 0)
	synonyms := make([]WordID, 0)
	antonyms := make([]WordID, 0)
	examples := make([]Sentence, 0)
	frequency := WordFrequency{}
	sources := make([]SourceID, 0)
	log.Debug("initialized vars ok")

	wg := new(sync.WaitGroup)

	for _, provider := range m.providers {
		source := provider.Source()

		log.WithFields(log.Fields{"source": source.Name, "relevancy": source.Relevancy}).Debug("walking providers")
		report, err := provider.FindWord(word)
		if err != nil {
			return nil, errors.Wrap(err, "provider "+source.Name+" failed to report '"+word+"' word")
		}

		log.WithFields(log.Fields{
			"word":             report.Word,
			"extract_duration": report.ExtractionDuration,
			"query_duration":   report.QueryDuration,
		}).Debug("report")

		for _, feature := range provider.AvailableFeatures() {
			switch feature {
			case Definitions:
				for _, def := range *report.Definitions {
					definitions = append(definitions, def)
				}
			case Synonyms:
				for _, syn := range *report.Synonyms {
					wg.Add(1)
					go func(syn string) {
						defer wg.Done()
						word, err := m.reportWord(syn, deepest+1)
						if err != nil {
							log.Error(err)
							return
						}

						if word != nil {
							synonyms = append(synonyms, WordID(syn))
						}
					}(syn)
				}
			case Antonyms:
				for _, ant := range *report.Antonyms {
					wg.Add(1)
					go func(ant string) {
						defer wg.Done()
						word, err := m.reportWord(ant, deepest+1)
						if err != nil {
							log.Error(err)
							return
						}
						if word != nil {
							antonyms = append(antonyms, WordID(ant))
						}
					}(ant)
				}
			case Frequency:
				frequency = *report.Frequency
			case Examples:
				for _, e := range *report.Examples {
					examples = append(examples, Sentence(e))
				}
			default:
				return nil, ErrFeatureNotExist
			}

		}

		sources = append(sources, source.ID)
	}

	wg.Wait()

	mapSynonyms := map[int64]WordID{}
	for i, s := range synonyms {
		mapSynonyms[int64(i)] = s
	}

	mapAntonyms := map[int64]WordID{}
	for i, s := range antonyms {
		mapAntonyms[int64(i)] = s
	}

	syntheticWord = Word{
		ID:          WordID(word),
		LastUpdate:  time.Now(),
		Sources:     sources,
		Definitions: definitions,
		Synonyms:    mapSynonyms,
		Antonyms:    mapAntonyms,
		Examples:    examples,
		Frequency:   frequency,
	}

	return m.persistence.SaveWord(syntheticWord)
}

func (m *Manager) ReportWord(word string) (*Word, error) {
	word = strings.ToLower(word)
	return m.reportWord(word, 0)
}
